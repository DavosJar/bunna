# Modelo YOLO API

API simple para preprocesar imágenes con YOLO.

## Instalar

```bash
pip install -r requirements.txt
```

## Ejecutar

```bash
uvicorn main:app --reload
```

## Endpoints

- `GET /health` - Verificar que está activa
- `POST /api/v1/diagnostico` - Enviar imagen para preprocesar
