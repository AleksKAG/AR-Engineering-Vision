import os
import requests
import tempfile
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from ultralytics import YOLO
import urllib.request

app = FastAPI(title="AR Engineering Vision AI Service")

# Загружаем модель YOLOv8 (для MVP используем nano-модель. 
# В продакшене здесь будет путь к вашей кастомной модели, обученной на трубах/лотках: YOLO("best.pt"))
model = YOLO("yolov8n.pt")

class AnalysisRequest(BaseModel):
    image_url: str
    expected_type: str  # Например: "duct", "pipe", "socket"

@app.post("/api/v1/ai/analyze")
async def analyze_image(req: AnalysisRequest):
    try:
        # 1. Скачиваем изображение по Presigned URL из MinIO
        with tempfile.NamedTemporaryFile(suffix=".jpg", delete=False) as tmp_file:
            urllib.request.urlretrieve(req.image_url, tmp_file.name)
            image_path = tmp_file.name

        # 2. Запускаем инференс YOLO
        results = model(image_path)
        
        detections = []
        for r in results:
            boxes = r.boxes
            for box in boxes:
                cls_id = int(box.cls[0])
                class_name = model.names[cls_id].lower()
                confidence = float(box.conf[0])
                
                # Фильтруем только то, что похоже на инженерные сети (упрощенно для MVP)
                # В реальной модели здесь будут классы: 'duct', 'cable_tray', 'pipe', 'socket'
                if confidence > 0.5:
                    detections.append({
                        "class": class_name,
                        "confidence": round(confidence, 2),
                        "bbox": box.xyxy[0].tolist() # [x1, y1, x2, y2]
                    })

        # 3. Простая логика сравнения для MVP
        ai_detected_type = detections[0]["class"] if detections else "unknown"
        is_match = req.expected_type.lower() in ai_detected_type or ai_detected_type == "unknown"
        
        return {
            "status": "success",
            "expected_type": req.expected_type,
            "ai_detected_type": ai_detected_type,
            "is_match": is_match,
            "detections": detections,
            "message": "OK" if is_match else f"Warning: Expected {req.expected_type}, but AI detected {ai_detected_type}"
        }

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        if 'image_path' in locals():
            os.remove(image_path)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
