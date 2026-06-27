# Módulo de Inferencia (YOLOv11) - CafeScan

## CU-RF4: Diagnóstico de Deficiencia de Nitrógeno
| Caso de Uso | Diagnóstico de Deficiencia de Nitrógeno |
| --- | --- |
| **Código** | CU-RF4 |
| **Actores** | Usuario registrado (Caficultor, Agrónomo) |
| **Tipo** | Esencial |
| **Autor** | Equipo Backend |

**Precondición:** El usuario tiene una imagen válida de una hoja de café y se encuentra autenticado en el sistema.
**Postcondición:** El sistema procesa la imagen, genera cajas delimitadoras (bounding boxes), devuelve la evaluación del nivel de nitrógeno junto con la imagen renderizada, y guarda el registro temporalmente.
**Propósito:** Permitir al usuario conocer de manera automática y en tiempo real el nivel de nitrógeno de su cultivo a través de visión por computadora.
**Resumen:** El usuario sube una imagen al sistema. La API de inferencia recibe la imagen, la procesa mediante el modelo pre-entrenado YOLO, calcula el porcentaje de áreas sanas vs deficientes, y devuelve un diagnóstico (Alto, Medio, Bajo) junto con recomendaciones agronómicas y la imagen procesada en base64.

### Curso Normal:
| Paso | Actor | Sistema (CafeScan / Inferencia) |
| --- | --- | --- |
| 1 | El usuario selecciona una imagen de una hoja y presiona "Analizar". | |
| 2 | | El sistema recibe el archivo `multipart/form-data` y lo guarda localmente en `uploads/`. |
| 3 | | El sistema carga el modelo YOLO (si no está en memoria) y ejecuta la predicción sobre la imagen. |
| 4 | | El sistema extrae las detecciones (`class_id`, `confidence`, `bbox`). |
| 5 | | El sistema ejecuta la lógica de evaluación (sanas vs deficiencia) para determinar el nivel (high, medium, low). |
| 6 | | El sistema dibuja las cajas sobre la imagen y la codifica en Base64. |
| 7 | | El sistema responde con el JSON final conteniendo detecciones, feedback, nivel de confianza y la imagen Base64. |
| 8 | El usuario visualiza la imagen procesada y las recomendaciones. | |

### Flujos Alternativos:
| ID | Condición | Acción del sistema |
| --- | --- | --- |
| FA-2a | Paso 2 — El archivo no es una imagen o está vacío. | El sistema (FastAPI) retorna error 422 de validación de entidad. |
| FA-4a | Paso 4 — YOLO no detecta ninguna hoja (0 detecciones). | El sistema genera un feedback con nivel "unknown" indicando que no se detectaron hojas y recomienda usar otra foto. |
| FA-3a | Paso 3 — El modelo `best.pt` no se encuentra en el servidor. | El sistema retorna HTTP 500 indicando "Model not found". |

---

## 1.5 Diagrama de Casos de Uso — Módulo Inferencia
*(Espacio reservado para diagrama UML de casos de uso)*

---

## SECCIÓN 2: DISEÑO

### 2.1 Tabla de Endpoints del Módulo Inferencia
| Método | Ruta | Descripción | Request Body | Response Body | JWT |
| --- | --- | --- | --- | --- | --- |
| GET | `/health` | Verifica que la API de inferencia está activa. | — | `{ status: "ok" }` | No |
| POST | `/api/v1/diagnostico` | Recibe una imagen, ejecuta el modelo YOLO y retorna resultados. | `multipart/form-data` (archivo) | `{ num_detections, detections[], feedback{}, avg_confidence, image }` | No (CORS Proxy) |
| GET | `/api/v1/diagnostico/{image_name}` | Ejecuta el diagnóstico usando una imagen previamente guardada en el disco. | — | `{ num_detections, detections[], feedback{}, avg_confidence, image }` | No |

### 2.2 Diagrama de Clases: Módulo Inferencia
*(Espacio reservado para el diagrama de clases)*

---

## SECCIÓN 3: IMPLEMENTACIÓN

### 3.1 Estructura de Carpetas — Módulo YOLOv11/
El módulo de inferencia está construido como un microservicio independiente en Python usando FastAPI y Ultralytics.

```text
YOLOv11/
├── api.py               → Servidor principal FastAPI y lógica de inferencia
├── requirements.txt     → Dependencias (fastapi, ultralytics, opencv, etc.)
├── runs/                → Resultados y pesos del modelo
│   ├── detect/
│   │   └── train/
│   │       └── weights/
│   │           └── best.pt  → Modelo YOLO pre-entrenado
│   └── predict/         → Imágenes cacheadas temporalmente por YOLO
└── uploads/             → Imágenes originales subidas por los usuarios
```

### 3.2 Código Clave del Módulo

#### 3.2.1 Inicialización del Modelo Singleton
Carga perezosa del modelo (Lazy Loading) para optimizar el inicio del servidor y mantener el modelo en memoria RAM:
```python
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
```

#### 3.2.2 Lógica de Negocio: Evaluación Agronómica (`build_feedback`)
Evalúa el porcentaje de hojas sanas vs deficiencia y retorna la recomendación adecuada:
```python
def build_feedback(detections):
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
        recommendation = "El nivel de nitrógeno es óptimo. Mantén el plan de fertilización actual..."
    elif porcentaje_sanas >= 40:
        level = "medium"
        label = "Medio"
        # ...
    else:
        level = "low"
        label = "Bajo"
        # ...

    return { "level": level, "label": label, "percentage": porcentaje_sanas, "recommendation": recommendation }
```

#### 3.2.3 Handler Principal de Inferencia (`/api/v1/diagnostico`)
Endpoint que gestiona la carga del archivo, ejecuta la predicción, formatea el JSON e imprime los logs:
```python
@app.post("/api/v1/diagnostico")
async def diagnostico(archivo: UploadFile = File(...)):
    UPLOADS_DIR.mkdir(parents=True, exist_ok=True)
    suffix = Path(archivo.filename or "image").suffix
    image_path = UPLOADS_DIR / f"{uuid4().hex}{suffix}"

    # Guardar archivo
    content = await archivo.read()
    image_path.write_bytes(content)

    # Inferencia con YOLO
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
    
    # ... Cálculo de promedio de confianza y retorno de respuesta HTTP ...
```
