from fastapi import FastAPI, File, UploadFile, HTTPException
from ultralytics import YOLO
from pathlib import Path
from uuid import uuid4
import base64
import json
import cv2

APP_TITLE = "YOLOv11 API"
APP_VERSION = "1.0.0"

MODEL_PATH = Path("runs/detect/train/weights/best.pt")
UPLOADS_DIR = Path("uploads")
PREDICT_DIR = Path("runs/predict")
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
                    "confidence": round(float(box.conf[0]), 4),
                    "bbox_xyxy": [round(float(v), 2) for v in box.xyxy[0].tolist()],
                }
            )

    return detections


def build_feedback(detections):
    """Genera feedback de nitrógeno basado en las detecciones reales del modelo."""
    total = len(detections)

    if total == 0:
        return {
            "level": "unknown",
            "label": "Sin detectar",
            "percentage": 0,
            "recommendation": "No se detectaron hojas en la imagen. Intenta con otra foto más clara.",
        }

    deficiencias = sum(1 for d in detections if d["class_name"] == "deficiencia_nitrogeno")
    sanas = total - deficiencias
    porcentaje_sanas = round((sanas / total) * 100)

    if porcentaje_sanas >= 70:
        level = "high"
        label = "Alto"
        recommendation = (
            "El nivel de nitrógeno es óptimo. Mantén el plan de fertilización actual "
            "y realiza un nuevo análisis en 30 días para monitorear la evolución."
        )
    elif porcentaje_sanas >= 40:
        level = "medium"
        label = "Medio"
        recommendation = (
            "Se recomienda aplicar urea (46-0-0) a razón de 150g por planta. "
            "Repite el análisis en 15 días para evaluar la respuesta del cultivo."
        )
    else:
        level = "low"
        label = "Bajo"
        recommendation = (
            "Nivel crítico de nitrógeno detectado. Aplica fertilizante nitrogenado de forma "
            "urgente (200g/planta de urea). Consulta con un agrónomo para un plan de recuperación."
        )

    return {
        "level": level,
        "label": label,
        "percentage": porcentaje_sanas,
        "recommendation": recommendation,
    }


def encode_image_base64(result):
    """Convierte la imagen anotada por YOLO a base64 para enviar al frontend."""
    annotated = result.plot()  # numpy array BGR
    _, buffer = cv2.imencode(".jpg", annotated, [cv2.IMWRITE_JPEG_QUALITY, 85])
    b64 = base64.b64encode(buffer).decode("utf-8")
    return f"data:image/jpeg;base64,{b64}"


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
    results = yolo.predict(
        source=str(image_path),
        save=True,
        conf=CONF_THRESHOLD,
        project=str(PREDICT_DIR),
        name="results",
        exist_ok=True,
    )

    result = results[0]
    detections = build_payload(result)
    feedback = build_feedback(detections)
    image_b64 = encode_image_base64(result)

    avg_confidence = 0.0
    if detections:
        avg_confidence = round(sum(d["confidence"] for d in detections) / len(detections), 4)

    response = {
        "num_detections": len(detections),
        "detections": detections,
        "feedback": feedback,
        "avg_confidence": avg_confidence,
        "image": image_b64,
    }

    # Imprimir JSON en consola para monitoreo
    log = {k: v for k, v in response.items() if k != "image"}
    print("\n" + "=" * 60)
    print("📋 RESULTADO DIAGNÓSTICO")
    print("=" * 60)
    print(json.dumps(log, indent=2, ensure_ascii=False))
    print("=" * 60 + "\n")

    return response


@app.get("/api/v1/diagnostico/{image_name}")
async def diagnostico_desde_archivo(image_name: str):
    image_path = UPLOADS_DIR / image_name
    if not image_path.exists():
        raise HTTPException(status_code=404, detail="Image not found")

    yolo = get_model()
    results = yolo.predict(
        source=str(image_path),
        save=True,
        conf=CONF_THRESHOLD,
        project=str(PREDICT_DIR),
        name="results",
        exist_ok=True,
    )

    result = results[0]
    detections = build_payload(result)
    feedback = build_feedback(detections)
    image_b64 = encode_image_base64(result)

    avg_confidence = 0.0
    if detections:
        avg_confidence = round(sum(d["confidence"] for d in detections) / len(detections), 4)

    response = {
        "num_detections": len(detections),
        "detections": detections,
        "feedback": feedback,
        "avg_confidence": avg_confidence,
        "image": image_b64,
    }

    log = {k: v for k, v in response.items() if k != "image"}
    print("\n" + "=" * 60)
    print("📋 RESULTADO DIAGNÓSTICO")
    print("=" * 60)
    print(json.dumps(log, indent=2, ensure_ascii=False))
    print("=" * 60 + "\n")

    return response
