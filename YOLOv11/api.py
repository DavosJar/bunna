from fastapi import FastAPI, File, UploadFile, HTTPException
from ultralytics import YOLO
from pathlib import Path
from uuid import uuid4
import json

APP_TITLE = "YOLOv11 API"
APP_VERSION = "1.0.0"

MODEL_PATH = Path("runs/detect/train/weights/best.pt")
UPLOADS_DIR = Path("uploads")
CONF_THRESHOLD = 0.20

app = FastAPI(title=APP_TITLE, version=APP_VERSION)
model = None


def get_model():
    global model
    if model is None:
        if not MODEL_PATH.exists():
            raise HTTPException(
                status_code=500,
                detail=f"Model not found at {MODEL_PATH}",
            )
        model = YOLO(str(MODEL_PATH))
    return model


def build_payload(result):
    names = result.names or {}
    detections = []

    if result.boxes is not None and len(result.boxes) > 0:
        for box in result.boxes:
            class_id = int(box.cls[0])
            detections.append(
                {
                    "class_id": class_id,
                    "class_name": names.get(class_id, str(class_id)),
                    "confidence": float(box.conf[0]),
                    "bbox_xyxy": [float(v) for v in box.xyxy[0].tolist()],
                }
            )

    payload = {
        "image_input": str(result.path),
        "image_output": str(result.save_dir),
        "num_detections": len(detections),
        "detections": detections,
    }

    return payload


@app.get("/health")
async def health():
    return {"status": "ok"}


@app.post("/api/v1/diagnostico")
async def diagnostico(archivo: UploadFile = File(...)):
    UPLOADS_DIR.mkdir(parents=True, exist_ok=True)
    suffix = Path(archivo.filename or "image").suffix
    image_path = UPLOADS_DIR / f"{uuid4().hex}{suffix}"

    content = await archivo.read()
    image_path.write_bytes(content)

    yolo = get_model()
    results = yolo.predict(source=str(image_path), save=True, conf=CONF_THRESHOLD)
    payload = build_payload(results[0])

    return json.loads(json.dumps(payload, ensure_ascii=True))


@app.get("/api/v1/diagnostico/{image_name}")
async def diagnostico_desde_archivo(image_name: str):
    image_path = UPLOADS_DIR / image_name
    if not image_path.exists():
        raise HTTPException(status_code=404, detail="Image not found")

    yolo = get_model()
    results = yolo.predict(source=str(image_path), save=True, conf=CONF_THRESHOLD)
    payload = build_payload(results[0])

    return json.loads(json.dumps(payload, ensure_ascii=True))
