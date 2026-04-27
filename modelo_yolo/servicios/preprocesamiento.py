from PIL import Image
import io
import numpy as np
from config_yolo import TAMAÑO_ENTRADA, MEDIA, DESVIACION


def preprocesar_imagen(contenido_imagen: bytes):
    """Preprocesa imagen para YOLO."""
    img = Image.open(io.BytesIO(contenido_imagen))
    
    # Convertir a RGB
    if img.mode != "RGB":
        img = img.convert("RGB")
    
    # Redimensionar a 640x640
    img = img.resize((TAMAÑO_ENTRADA, TAMAÑO_ENTRADA))
    
    # Convertir a array
    imagen_array = np.array(img, dtype=np.float32) / 255.0
    
    # Normalizar con ImageNet
    imagen_array[..., 0] = (imagen_array[..., 0] - MEDIA[0]) / DESVIACION[0]
    imagen_array[..., 1] = (imagen_array[..., 1] - MEDIA[1]) / DESVIACION[1]
    imagen_array[..., 2] = (imagen_array[..., 2] - MEDIA[2]) / DESVIACION[2]
    
    return imagen_array
