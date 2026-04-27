"""
Configuración de preprocesamiento para YOLO.
Parámetros recomendados y formatos de entrada/salida.
"""

# Parámetros de entrada de imagen
TAMAÑO_ENTRADA = 640  # Recomendado para YOLOv8
FORMATO_IMAGEN = "RGB"  # YOLOv8 usa RGB
CANALES = 3

# Normalización
MEDIA = [0.485, 0.456, 0.406]  # ImageNet normalization
DESVIACION = [0.229, 0.224, 0.225]  # ImageNet normalization

# Parámetros de preprocesamiento
INTERPOLACION = "bilinear"  # Interpolación al redimensionar
ESCALA_VALORES = 255.0  # Normalizar valores a [0, 1]

# Validación de imagen
FORMATOS_PERMITIDOS = ["jpeg", "jpg", "png", "bmp", "webp"]
TAMAÑO_MAX_MB = 50
TAMAÑO_MIN_PIXELES = 32

# Parámetros de detección
CONFIANZA_MINIMA = 0.25
IOU_THRESHOLD = 0.45

# Configuración de predicción
BATCH_SIZE = 1
DEVICE = "cpu"  # "cuda" si hay GPU disponible
