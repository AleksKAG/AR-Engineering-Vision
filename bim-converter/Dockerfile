FROM python:3.10-slim

# Установка системных зависимостей для IfcOpenShell
RUN apt-get update && apt-get install -y \
    libx11-6 \
    libgl1-mesa-glx \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY main.py .

CMD ["python", "main.py"]
