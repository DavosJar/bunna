from ultralytics import YOLO

# 1. Cargamos TU modelo ya entrenado
model = YOLO('runs/detect/train/weights/best.pt')

# 2. Ejecutamos la detección en una imagen nueva
# Cambia 'foto_prueba.jpg' por el nombre de una imagen que tengas
results = model.predict(source='1.jpg', save=True, conf=0.20)

# 3. Esto te dirá dónde se guardó la imagen con los cuadros dibujados
print("La imagen resultante está en: runs/detect/predict/")
