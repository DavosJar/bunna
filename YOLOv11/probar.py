import json
from ultralytics import YOLO

# 1. Cargamos tu modelo ya entrenado
model = YOLO("runs/detect/train/weights/best.pt")

# 2. Ejecutamos la deteccion en una imagen nueva
# Cambia '3.jpg' por el nombre de una imagen que tengas
results = model.predict(source="1.jpg", save=True, conf=0.20)

# 3. Construimos un diagnostico estructurado en JSON
result = results[0]
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

print(json.dumps(payload, ensure_ascii=True, indent=2))
