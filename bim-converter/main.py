import os
import json
import uuid
import tempfile
import ifcopenshell
import ifcopenshell.util.element
import ifcopenshell.geom
import requests
from minio import Minio
from minio.error import S3Error

# Конфигурация из переменных окружения
MINIO_ENDPOINT = os.getenv("MINIO_ENDPOINT", "minio:9000")
MINIO_ACCESS_KEY = os.getenv("MINIO_ACCESS_KEY", "minioadmin")
MINIO_SECRET_KEY = os.getenv("MINIO_SECRET_KEY", "minioadmin")
MINIO_BUCKET = os.getenv("MINIO_BUCKET", "ar-models")
GO_BACKEND_URL = os.getenv("GO_BACKEND_URL", "http://app:8080/api/v1/internal/bim-processed")

# Инициализация MinIO клиента
minio_client = Minio(
    MINIO_ENDPOINT,
    access_key=MINIO_ACCESS_KEY,
    secret_key=MINIO_SECRET_KEY,
    secure=False
)

def process_ifc(project_id: str, ifc_s3_key: str):
    print(f"🚀 Начало обработки: {ifc_s3_key}")
    
    # 1. Скачиваем IFC из MinIO во временный файл
    with tempfile.NamedTemporaryFile(suffix=".ifc", delete=False) as tmp_ifc:
        minio_client.fget_object(MINIO_BUCKET, ifc_s3_key, tmp_ifc.name)
        ifc_file_path = tmp_ifc.name

    try:
        # 2. Открываем IFC файл
        ifc_file = ifcopenshell.open(ifc_file_path)
        
        # Настройки геометрии для извлечения (упрощенные bounding boxes для MVP)
        settings = ifcopenshell.geom.settings()
        settings.set(settings.USE_WORLD_COORDS, True)
        settings.set(settings.INCLUDE_CURVES, False)

        rooms_data = []
        elements_data = []

        # 3. Извлекаем помещения (IfcSpace)
        spaces = ifc_file.by_type("IfcSpace")
        for space in spaces:
            room_id = str(uuid.uuid4())
            bbox = extract_bbox(ifc_file, space, settings)
            
            rooms_data.append({
                "id": room_id,
                "project_id": project_id,
                "name": space.Name or "Без названия",
                "bbox": bbox
            })

            # 4. Извлекаем элементы внутри этого помещения (упрощенно: все IfcFlowSegment, IfcDistributionElement)
            # В реальном проекте здесь используется ifcopenshell.util.spatial.get_container
            elements = ifc_file.by_type("IfcFlowSegment") + ifc_file.by_type("IfcDistributionControlElement")
            for elem in elements:
                elem_id = str(uuid.uuid4())
                elem_bbox = extract_bbox(ifc_file, elem, settings)
                elem_type = elem.is_a() # например, "IfcPipeSegment"
                
                elements_data.append({
                    "id": elem_id,
                    "room_id": room_id,
                    "type": elem_type,
                    "world_coords": elem_bbox["center"], # {x, y, z}
                    "properties": {
                        "name": elem.Name,
                        "tag": elem.Tag,
                        "dimensions": elem_bbox["dimensions"]
                    }
                })

        # 5. Конвертация всей модели в GLB (упрощенно через ifcopenshell)
        # Для MVP сохраняем тот же IFC под именем .glb, чтобы не усложнять, 
        # но в продакшене здесь должен быть вызов ifcopenshell.geom.create_shape и экспорт в gltf
        glb_s3_key = f"projects/{project_id}/model.glb"
        
        # Загружаем "результат" обратно в MinIO (пока просто копируем для теста, 
        # реальный glb-экспорт требует ~20 строк кода с trimesh)
        minio_client.fput_object(MINIO_BUCKET, glb_s3_key, ifc_file_path)

        # 6. Отправляем метаданные обратно в Go-бэкенд
        payload = {
            "project_id": project_id,
            "glb_s3_key": glb_s3_key,
            "rooms": rooms_data,
            "elements": elements_data
        }
        
        print(f"📤 Отправка метаданных в Go: {len(rooms_data)} комнат, {len(elements_data)} элементов")
        response = requests.post(GO_BACKEND_URL, json=payload)
        response.raise_for_status()
        print("✅ Обработка успешно завершена!")

    except Exception as e:
        print(f"❌ Ошибка при обработке IFC: {e}")
    finally:
        os.remove(ifc_file_path)

def extract_bbox(ifc_file, element, settings):
    """Извлекает ограничивающий прямоугольник (Bounding Box) элемента"""
    try:
        shape = ifcopenshell.geom.create_shape(settings, element)
        bbox = shape.geometry.bounding_box()
        return {
            "center": {"x": bbox[0], "y": bbox[1], "z": bbox[2]},
            "dimensions": {"x": bbox[3], "y": bbox[4], "z": bbox[5]}
        }
    except:
        return {"center": {"x": 0, "y": 0, "z": 0}, "dimensions": {"x": 0, "y": 0, "z": 0}}

if __name__ == "__main__":
    # Для теста можно запускать так: python main.py <project_id> <ifc_s3_key>
    import sys
    if len(sys.argv) == 3:
        process_ifc(sys.argv[1], sys.argv[2])
    else:
        print("Использование: python main.py <project_id> <ifc_s3_key>")
