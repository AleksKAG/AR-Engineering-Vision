FROM python:3.10-slim

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY main.py .

# Предварительно загружаем веса YOLO, чтобы контейнер стартовал быстрее
RUN python -c "from ultralytics import YOLO; YOLO('yolov8n.pt')"

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
