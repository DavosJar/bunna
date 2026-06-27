# Plan de Pruebas — YOLOv11 + Módulo de Inferencia

> **Propósito:** Documentar el pipeline completo de inferencia (entrenamiento → detección → diagnóstico), los tests existentes y el plan de pruebas para los componentes YOLOv11, image-service y el módulo de inferencia de fincas.
>
> **Servicios involucrados:**
> - `YOLOv11/` — Python (Ultralytics, FastAPI, PyTorch)
> - `image-service/` — Go (procesamiento de imágenes + MQTT)
> - `fincas/` — Go (dominio de inferencia, registro de diagnóstico)
> - `frontend/` — JavaScript (dashboard de diagnóstico)
>
> **Última actualización:** 2026-06-25

---

## 1. Resumen Ejecutivo

| Componente | Lenguaje | Tests existentes | Cobertura | Estado |
|---|---|---|---|---|
| YOLOv11 - Entrenamiento | Python | 0 | 0% | ❌ |
| YOLOv11 - Inferencia (`probar.py`) | Python | 0 | 0% | ❌ |
| YOLOv11 - API (FastAPI) | Python | 0 | 0% | ❌ No implementado |
| **image-service** | Go 1.26 | 0 | 0% | ❌ |
| **fincas - domain inferencia** | Go 1.26 | 4 tests (handler) | Parcial | 🔶 |
| **fincas - registrarinferencia** | Go 1.26 | 0 (use case) | 0% | ❌ |
| **fincas - domain diagnóstico** | Go 1.26 | 0 (domain) | 0% | ❌ |
| **fincas - repositorios** | Go 1.26 | 0 | 0% | ❌ |
| **frontend - yoloApi** | JavaScript | 0 | 0% | ❌ |
| **frontend - DashboardPage** | JSX | 0 | 0% | ❌ |
| **Pipeline E2E** | — | 0 | 0% | ❌ |

**Conclusión:** El pipeline de inferencia es el módulo con **menor cobertura de pruebas** del monorepo. No existe ningún test automatizado.

---

## 2. Arquitectura del Pipeline de Inferencia

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        PIPELINE DE INFERENCIA                          │
│                                                                         │
│  FLUJO EDGE (automático)                                                │
│  ┌──────────┐    ┌──────────────┐    ┌──────────┐    ┌───────────────┐ │
│  │ Sensor   │───▶│ image-service │───▶│ YOLO     │───▶│ RabbitMQ      │ │
│  │ MQTT     │    │ (resize 640)  │    │ inference │    │ Consumer      │ │
│  └──────────┘    └──────────────┘    └──────────┘    └───────┬───────┘ │
│                                                               │         │
│  FLUJO MANUAL (on-demand)                                      │         │
│  ┌──────────┐    ┌──────────────┐    ┌──────────┐              │         │
│  │ Usuario  │───▶│ Frontend     │───▶│ YOLO API │              │         │
│  │ (imagen) │    │ Dashboard    │    │ (FastAPI)│              │         │
│  └──────────┘    └──────────────┘    └──────────┘              │         │
│                                                                 │         │
│  POST-PROCESAMIENTO (Fincas Service)                             │         │
│  ┌──────────────────────────────────────────────────────────────▼──────┐ │
│  │  RegistrarInferencia (use case)                                     │ │
│  │  ┌──────────┐    ┌──────────────┐    ┌───────────────────────────┐ │ │
│  │  │ Validar  │───▶│ Construir    │───▶│ Persistir Diagnóstico     │ │ │
│  │  │ Comando  │    │ ResultadoInf │    │ (PENDIENTE) + publicar    │ │ │
│  │  └──────────┘    └──────────────┘    │ evento DiagnosticoCreado │ │ │
│  │                                      └───────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                          │
│  FEEDBACK LOOP                                                           │
│  ┌────────────┐    ┌──────────────────┐    ┌──────────────────────────┐ │
│  │ Aceptar /  │───▶│ Candidato        │───▶│ Reentrenamiento YOLO     │ │
│  │ Rechazar   │    │ Reentrenamiento   │    │ (futuro ciclo)           │ │
│  └────────────┘    └──────────────────┘    └──────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 3. YOLOv11 — Modelo de Visión

### 3.1 Dataset

| Propiedad | Valor |
|---|---|
| Clases | `deficiencia_nitrogeno` (0), `hoja_sana` (1) |
| Imágenes entrenamiento | 165 |
| Imágenes validación | 37 |
| Fuente | Roboflow (`cesars-workspace-lzja3`) |
| Tamaño entrada | 640×640 px |
| Modelo base | YOLO11n (nano) → entrenado a `best.pt` |
| Épocas | 300 (early stopping con patience 50) |
| Batch | 16 |

### 3.2 Archivos de Entrenamiento e Inferencia

| Archivo | Propósito | Líneas | Tests |
|---|---|---|---|
| `entrenamieto.py` | Script de entrenamiento | 11 | ❌ 0 |
| `probar.py` | Script de inferencia (genera JSON) | 35 | ❌ 0 |
| `test.py` | Placeholder vacío | 0 | ❌ |
| `setup.sh` | Script de setup del entorno | 9 | ❌ |
| `requirements.txt` | Dependencias | 8 | ❌ |
| `data.yaml` | Configuración del dataset | 13 | ❌ |

### 3.3 Formato del Output de Inferencia (`probar.py`)

```json
{
  "image_input": "3.jpg",
  "image_output": "runs/detect/results",
  "num_detections": 2,
  "detections": [
    {
      "class_id": 0,
      "class_name": "deficiencia_nitrogeno",
      "confidence": 0.85,
      "bbox_xyxy": [100.5, 200.3, 300.7, 400.2]
    },
    {
      "class_id": 1,
      "class_name": "hoja_sana",
      "confidence": 0.92,
      "bbox_xyxy": [150.1, 250.4, 350.6, 450.3]
    }
  ]
}
```

**Campos mapeados al dominio:**
| JSON output | Domain VO | Tipo |
|---|---|---|
| `detections[].class_name == "deficiencia_nitrogeno"` | `ResultadoInferencia.TieneClorosis = true` | `bool` |
| `detections[].confidence` | `ResultadoInferencia.Confianza` | `float64` (0-1) |
| `image_input` | `ResultadoInferencia.ImageUrl` | `string` |
| timestamp de procesamiento | `ResultadoInferencia.ProcesadoAt` | `time.Time` |

### 3.4 Estado Actual del Modelo

| Métrica | Valor |
|---|---|
| Weights | `runs/detect/train/weights/best.pt` |
| mAP (no reportado) | ❌ Sin métricas de validación |
| Precisión (no reportado) | ❌ Sin métricas |
| Recall (no reportado) | ❌ Sin métricas |
| Confianza mínima usada | 0.20 (en `probar.py`) |

---

## 4. Image-Service — Procesamiento de Imágenes

### 4.1 Componentes

| Archivo | Propósito | Líneas | Tests |
|---|---|---|---|
| `main.go` | Entry point, loop de escaneo de directorio | 106 | ❌ 0 |
| `image_processor.go` | Redimensionamiento a 640px (Lanczos3) | 100 | ❌ 0 |
| `mqtt_publisher.go` | Publicación MQTT (topico `images/processed/{filename}`) | 57 | ❌ 0 |
| `docker-compose.yml` | Mosquitto MQTT broker | 14 | ❌ |
| `Makefile` | Build, run, test | 61 | ❌ |

### 4.2 Flujo del Image-Service

```
1. Escanea directorio cada 5s (ticker)
2. Filtra archivos: .jpg, .jpeg, .png
3. Procesa cada imagen:
   a. Decodifica JPEG o PNG
   b. Redimensiona a 640px manteniendo aspect ratio (nfnt/resize Lanczos3)
   c. Re-codifica a JPEG/PNG
4. Publica bytes procesados a MQTT topic: images/processed/{filename}
5. (Opcional) Elimina o mueve archivo original
```

### 4.3 Imágenes de Prueba

| Archivo | Ubicación |
|---|---|
| `test_1.jpg` | `image-service/images/` |
| `test_2.jpg` | `image-service/images/` |
| `test_3.jpg` | `image-service/images/` |
| `test_4.jpg` | `image-service/images/` |
| `test_5.jpg` | `image-service/images/` |

---

## 5. Fincas — Módulo de Inferencia y Diagnóstico

### 5.1 Capa de Dominio — Entidades

| Archivo | Entidad | Métodos clave | Tests |
|---|---|---|---|
| `inferencia_vo.go` | `ResultadoInferencia` (VO) | Constructor con validación (URL no vacía, confianza 0-1) | ❌ 0 |
| `diagnostico.go` | `Diagnostico` (Aggregate) | `NewDiagnostico`, `NewDiagnosticoFromStorage`, `MarcarComoAceptado`, `MarcarComoRechazado` | ❌ 0 |
| `muestras.go` | `Muestra`, `Ubicacion` | `NewMuestra`, validación coordenadas | ❌ 0 |
| `candidato_reentrenamiento.go` | `CandidatoReentrenamiento` | Constructor con todos los campos | ❌ 0 |
| `errores.go` | Errores de dominio | 7 tipos de error | ❌ 0 |
| `especificaciones.go` | Specification pattern | Filtros para búsqueda | ❌ 0 |
| `repositories.go` | Interfaces | 3 interfaces de repositorio | ❌ 0 |

#### Máquina de Estados del Diagnóstico

```
        ┌──────────┐
        │ PENDIENTE│  (estado inicial)
        └────┬─────┘
             │
      ┌──────┴──────┐
      ▼             ▼
┌──────────┐  ┌──────────┐
│ ACEPTADO │  │ RECHAZADO│
└──────────┘  └────┬─────┘
                   │
                   ▼
        ┌──────────────────┐
        │ Candidato        │
        │ Reentrenamiento   │
        └──────────────────┘
```

### 5.2 Capa de Aplicación — Use Case `RegistrarInferencia`

| Archivo | Propósito | Líneas | Tests |
|---|---|---|---|
| `command.go` | DTO de entrada + validación | 37 | ❌ 0 |
| `usecase.go` | Orquestación del registro de inferencia | 126 | ❌ 0 |
| `salida.go` | DTO de respuesta | 16 | ❌ 0 |

**Flujo del use case:**
```
1. Validar comando (MuestraID, ImageURL, TieneClorosis, Confianza, ProcesadoAt)
2. Cargar Muestra desde BD (MuestraRepositorio.ObtenerPorID)
3. Construir ResultadoInferencia VO
4. Generar ID único
5. Generar nombre auto: "INF-{YYYYMMDD}-{random}"
6. Extraer TenantID de la muestra
7. Construir Diagnostico entity (estado = PENDIENTE)
8. Persistir (DiagnosticoRepositorio.Crear)
9. Publicar evento DiagnosticoCreado (routing key: diagnosticos.v1.diagnostico.creado)
```

**Casos de prueba necesarios:**

| Escenario | Entrada esperada | Resultado esperado |
|---|---|---|
| Inferencia exitosa (clorosis detectada) | MuestraID válido, TieneClorosis=true, Confianza=0.85 | Diagnóstico PENDIENTE creado |
| Inferencia exitosa (hoja sana) | MuestraID válido, TieneClorosis=false, Confianza=0.92 | Diagnóstico PENDIENTE creado |
| MuestraID no existe | UUID inválido | ErrNotFound |
| Confianza fuera de rango (>1.0) | Confianza=1.5 | ErrValidacion |
| Confianza negativa | Confianza=-0.1 | ErrValidacion |
| ImageURL vacía | ImageURL="" | ErrValidacion |
| Fallo al persistir | Error de BD | Error propagado, sin evento |
| Fallo al publicar evento | Error de RabbitMQ | Error propagado |

### 5.3 Capa de Presentación — Handler Tests (EXISTENTES)

**Archivo:** `diagnostico_handler_test.go` (245 líneas) — **ÚNICOS TESTS DEL MÓDULO**

#### `TestDiagnosticoHandler_SolicitarManual` — 3 sub-tests

| # | Test | Escenario | Status Esperado | Resultado |
|---|---|---|---|---|
| 1 | SolicitarManual/Exito | Facade retorna SolicitudID | 201 | ✅ |
| 2 | SolicitarManual/CuerpoInvalido | JSON mal formado | 400 | ✅ |
| 3 | SolicitarManual/ErrorFacade | Facade retorna error | 422 | ✅ |

#### `TestDiagnosticoHandler_Aceptar` — 2 sub-tests

| # | Test | Escenario | Status Esperado | Resultado |
|---|---|---|---|---|
| 4 | Aceptar/Exito | Facade retorna EstadoCambioResponse | 200 | ✅ |
| 5 | Aceptar/ErrorFacade | Facade retorna error | 422 | ✅ |

#### `TestDiagnosticoHandler_Rechazar` — 3 sub-tests

| # | Test | Escenario | Status Esperado | Resultado |
|---|---|---|---|---|
| 6 | Rechazar/Exito | Facade retorna EstadoCambioResponse | 200 | ✅ |
| 7 | Rechazar/CuerpoInvalido | JSON mal formado | 400 | ✅ |
| 8 | Rechazar/ErrorFacade | Facade retorna error | 422 | ✅ |

### 5.4 Capa de Infraestructura — Repositorios

| Archivo | Operaciones | Tests |
|---|---|---|
| `diagnostico_model.go` | GORM model + ToDomain/FromDomain | ❌ 0 |
| `diagnostico_repositorio.go` | CRUD + ListarPorFinca + Buscar con specs | ❌ 0 |
| `muestra_model.go` | GORM model | ❌ 0 |
| `muestra_repositorio.go` | CRUD + ListarPorDiagnostico | ❌ 0 |
| `candidato_model.go` | GORM model | ❌ 0 |
| `candidato_repositorio.go` | Crear + ListarPendientes | ❌ 0 |
| `unit_of_work.go` | Transaccional para rechazo | ❌ 0 |

### 5.5 Event Publisher

| Archivo | Implementación | Tests |
|---|---|---|
| `internal/application/event_publisher.go` | Interface `EventPublisher` | ❌ 0 |
| `internal/infrastructure/eventpublisher/console.go` | ConsolePublisher (stdout) | ❌ 0 |

**Nota:** No existe implementación de RabbitMQ. El `ConsolePublisher` solo loggea a stdout. Tampoco existe un consumer de RabbitMQ que dispare `RegistrarInferencia`.

---

## 6. Frontend — Dashboard de Diagnóstico

### 6.1 `yoloApi.js`

| Función | Método HTTP | Endpoint | Tests |
|---|---|---|---|
| `verificarSaludYOLO()` | GET | `/health` | ❌ 0 |
| `diagnosticar(imagen)` | POST | `/api/v1/diagnostico` (multipart, campo `archivo`) | ❌ 0 |

### 6.2 `DashboardPage.jsx`

| Funcionalidad | Estado | Tests |
|---|---|---|
| Selección de imagen (drag & drop / file input) | ✅ Implementado | ❌ 0 |
| Llamada a `diagnosticar()` | ✅ Implementado | ❌ 0 |
| Visualización de resultados (badge, gauge, chips) | ✅ Implementado | ❌ 0 |
| Manejo de errores (fallback UI) | ✅ Implementado | ❌ 0 |
| Loading state durante análisis | ✅ Implementado | ❌ 0 |

---

## 7. Análisis de Brechas de Pruebas

### 7.1 Resumen de Cobertura

| Componente | Nivel | Tests existentes | Cobertura requerida | Déficit |
|---|---|---|---|---|
| YOLOv11 - Entrenamiento | Unit/Integration | 0 | Validar entrenamiento, métricas, overfitting | 100% |
| YOLOv11 - Inferencia | Unit | 0 | Precisión de detección en imágenes de prueba | 100% |
| YOLOv11 - API (FastAPI) | Unit/Integration | 0 (no existe) | Endpoints, validación, formato output | — |
| image-service image_processor | Unit | 0 | Redimensionamiento, codificación, errores | 100% |
| image-service mqtt_publisher | Integration | 0 | Conexión, publicación QoS 1, reconexión | 100% |
| image-service main loop | Integration | 0 | Escaneo directorio, pipeline completo | 100% |
| fincas - inferencia_vo | Unit | 0 | Validación campos, rangos | 100% |
| fincas - diagnostico domain | Unit | 0 | Máquina de estados, transiciones | 100% |
| fincas - diagnostico repositorio | Integration | 0 | CRUD, filtros, paginación | 100% |
| fincas - registrarinferencia use case | Unit | 0 | Flujo completo, errores, eventos | 100% |
| fincas - diagnostico handler | Unit | 8 tests | Edge cases, auth, validación (parcial) | ~50% |
| frontend - yoloApi | Unit | 0 | Llamadas HTTP, manejo de respuestas | 100% |
| frontend - DashboardPage | Component | 0 | Renderizado, estados, interacciones | 100% |
| Pipeline E2E | E2E | 0 | Flujo completo imagen → diagnóstico | 100% |

### 7.2 Brechas Críticas

| # | Brecha | Riesgo |
|---|---|---|
| 1 | **Sin tests de inferencia** — No se valida que `probar.py` produzca JSON válido | Rotura silenciosa del pipeline |
| 2 | **Sin tests de image-service** — No se verifica redimensionamiento correcto | Imágenes corruptas o mal formateadas |
| 3 | **Sin tests de `ResultadoInferencia` VO** — Validación de confianza 0-1 no testeada | Datos inválidos en BD |
| 4 | **Sin tests de máquina de estados** — Transiciones PENDIENTE→ACEPTADO/RECHAZADO no validadas | Estados inconsistentes |
| 5 | **Sin tests de repositorio** — CRUD de diagnósticos no probado | Fallos de persistencia en producción |
| 6 | **Sin Event Publisher implementado** — RabbitMQ no existe, solo console.log | Eventos perdidos en producción |
| 7 | **Sin RabbitMQ Consumer** — `RegistrarInferencia` no tiene quien lo ejecute | Pipeline roto |
| 8 | **Sin tests frontend** — UI del dashboard no probada | Regresión visual o funcional |

---

## 8. Plan de Pruebas Propuesto

### 8.1 Prioridad 🔴 CRÍTICA (Semana 1-2)

| # | Tarea | Componente | Tipo | Esfuerzo | Dependencias |
|---|---|---|---|---|---|
| 1 | Test `ResultadoInferencia` VO | `fincas/internal/diagnostico/domain/` | Unit (Go) | 1h | — |
| 2 | Test `Diagnostico` aggregate (máquina estados) | `fincas/internal/diagnostico/domain/` | Unit (Go) | 2h | — |
| 3 | Test `Muestra` + `Ubicacion` | `fincas/internal/diagnostico/domain/` | Unit (Go) | 1h | — |
| 4 | Test `CandidatoReentrenamiento` | `fincas/internal/diagnostico/domain/` | Unit (Go) | 1h | — |
| 5 | Test `RegistrarInferencia` use case | `fincas/internal/application/usecases/registrarinferencia/` | Unit (Go) | 3h | #1, #2 |

### 8.2 Prioridad 🟡 ALTA (Semana 3-4)

| # | Tarea | Componente | Tipo | Esfuerzo |
|---|---|---|---|---|
| 6 | Test `image_processor` (resize + encode) | `image-service/` | Unit (Go) | 2h |
| 7 | Test `mqtt_publisher` (connect + publish) | `image-service/` | Integration (Go) | 2h |
| 8 | Test `diagnostico_repositorio` (CRUD + filters) | `fincas/` | Integration (Go) | 3h |
| 9 | Test `muestra_repositorio` | `fincas/` | Integration (Go) | 2h |
| 10 | Test `candidato_repositorio` | `fincas/` | Integration (Go) | 1h |
| 11 | Test `yoloApi.js` (axios mock) | `frontend/` | Unit (Vitest) | 1h |
| 12 | Test `DashboardPage` (render + states) | `frontend/` | Component (RTL) | 3h |

### 8.3 Prioridad 🟢 MEDIA (Semana 5-6)

| # | Tarea | Componente | Tipo | Esfuerzo |
|---|---|---|---|---|
| 13 | Implementar RabbitMQ Event Publisher | `fincas/` | Integration | 4h |
| 14 | Implementar RabbitMQ Consumer (para `RegistrarInferencia`) | `fincas/` | Integration | 4h |
| 15 | Test pipeline image-service → MQTT → consumer | `image-service/` + `fincas/` | E2E | 4h |
| 16 | Test YOLO API (FastAPI endpoint wrapper) | `YOLOv11/` | Integration | 3h |
| 17 | Test YOLO inference con imágenes de prueba conocidas | `YOLOv11/` | Integration | 2h |
| 18 | Test frontend → YOLO API → diagnóstico | `frontend/` + `YOLOv11/` | E2E | 4h |

### 8.4 Prioridad 🔵 BAJA (Semana 7-8)

| # | Tarea | Componente | Tipo | Esfuerzo |
|---|---|---|---|---|
| 19 | Validación de métricas del modelo (mAP, precisión, recall) | `YOLOv11/` | Data validation | 3h |
| 20 | Test de estrés: 100 imágenes simultáneas | `image-service/` | Performance | 3h |
| 21 | Test de integración continua para YOLO (GitHub Actions) | `YOLOv11/` | CI/CD | 2h |

### 8.5 Timeline

```
Sprint 1 (🔴 Crítica)          Sprint 2 (🟡 Alta)           Sprint 3 (🟢 Media/Baja)
┌──────────────────────┐      ┌──────────────────────┐      ┌───────────────────────┐
│ #1  ResultadoInf VO  │      │ #6  image-processor  │      │ #13 RabbitMQ publisher│
│ #2  Diagnóstico st   │      │ #7  mqtt-publisher   │      │ #14 RabbitMQ consumer │
│ #3  Muestra domain   │ ──▶  │ #8  diag-repositorio │ ──▶  │ #15 Pipeline MQTT E2E │
│ #4  Candidato domain │      │ #9  muestra-repo     │      │ #16 YOLO API tests    │
│ #5  RegistrarInf UC  │      │ #10 candidato-repo   │      │ #17 YOLO inference    │
│                      │      │ #11 yoloApi.js (FE)  │      │ #18 Front→YOLO E2E    │
│                      │      │ #12 DashboardPage FE │      │ #19-21 Modelo + CI    │
└──────────────────────┘      └──────────────────────┘      └───────────────────────┘
```

---

## 9. Casos de Prueba Detallados

### 9.1 Domain: `ResultadoInferencia` VO

| # | Test | Input | Expected |
|---|---|---|---|
| 1 | Crear resultado válido (clorosis) | imageURL="https://img.jpg", tieneClorosis=true, confianza=0.85, procesadoAt=now | VO creado sin error |
| 2 | Crear resultado válido (sana) | imageURL="https://img.jpg", tieneClorosis=false, confianza=0.95, procesadoAt=now | VO creado sin error |
| 3 | Confianza exactamente 0.0 | confianza=0.0 | Válido (límite inferior) |
| 4 | Confianza exactamente 1.0 | confianza=1.0 | Válido (límite superior) |
| 5 | Confianza negativa | confianza=-0.01 | Error |
| 6 | Confianza > 1.0 | confianza=1.01 | Error |
| 7 | ImageURL vacía | imageURL="" | Error |
| 8 | ImageURL solo espacios | imageURL="   " | Error |
| 9 | ProcesadoAt zero time | procesadoAt=time.Time{} | Error |
| 10 | ProcesadoAt en futuro | procesadoAt=futureDate | Error |

### 9.2 Domain: `Diagnostico` Aggregate

| # | Test | Transición | Expected |
|---|---|---|---|
| 1 | Crear diagnóstico PENDIENTE | NewDiagnostico() | estado == PENDIENTE |
| 2 | Aceptar desde PENDIENTE | MarcarComoAceptado() | estado == ACEPTADO |
| 3 | Rechazar desde PENDIENTE | MarcarComoRechazado() | estado == RECHAZADO |
| 4 | Rechazar con motivo | MarcarComoRechazado("imagen borrosa") | estado == RECHAZADO + motivo |
| 5 | Aceptar desde ACEPTADO | MarcarComoAceptado() | Error (transición inválida) |
| 6 | Rechazar desde ACEPTADO | intentar cambiar a RECHAZADO | Error |
| 7 | Aceptar desde RECHAZADO | intentar cambiar a ACEPTADO | Error |
| 8 | Rechazar desde RECHAZADO | MarcarComoRechazado() | Error |
| 9 | Nombre auto-generado formato | "INF-20260625-XXXX" | Coincide patrón `^INF-\d{8}-[A-Z0-9]{4}$` |
| 10 | Reconstruir desde BD | NewDiagnosticoFromStorage() | Mismos valores |

### 9.3 Use Case: `RegistrarInferencia`

| # | Test | Setup | Expected |
|---|---|---|---|
| 1 | Registro exitoso (clorosis) | Muestra existe, datos válidos | Diagnóstico creado PENDIENTE, evento publicado |
| 2 | Registro exitoso (sana) | Muestra existe, tieneClorosis=false | Diagnóstico creado PENDIENTE |
| 3 | Muestra no encontrada | MuestraID inexistente | ErrNotFound |
| 4 | Confianza inválida | confianza=2.0 | ErrValidacion, sin persistencia |
| 5 | Error al persistir | repo.Crear retorna error | Error propagado, sin evento |
| 6 | Error al publicar evento | publisher.Publish retorna error | Error propagado |
| 7 | MuestraID vacío | MuestraID = "" | ErrValidacion |
| 8 | Contexto cancelado | ctx cancelado antes de persistir | Error de contexto |

### 9.4 Image-Service: `image_processor.go`

| # | Test | Input | Expected |
|---|---|---|---|
| 1 | Redimensionar JPEG 1920×1080 | test_1.jpg | Output 640px ancho, aspect ratio preservado |
| 2 | Redimensionar PNG 800×600 | test_2.jpg | Output 640px ancho |
| 3 | Redimensionar imagen cuadrada 500×500 | test_3.jpg | Output 640×640 |
| 4 | Archivo no existe | "/no/existe.jpg" | Error |
| 5 | Formato no soportado | "imagen.gif" | Error o skip |
| 6 | Archivo corrupto | datos binarios inválidos | Error de decodificación |
| 7 | Directorio vacío | dir sin imágenes | Lista vacía, sin error |

### 9.5 Frontend: `yoloApi.diagnosticar()`

| # | Test | Mock | Expected |
|---|---|---|---|
| 1 | Diagnóstico exitoso | 200 + JSON detecciones | Retorna data parseada |
| 2 | Error 400 (bad request) | 400 | Error lanzado |
| 3 | Error 500 (server error) | 500 | Error lanzado |
| 4 | Network error | axios lanza error de red | Error lanzado |
| 5 | Timeout excedido | respuesta > 30s | Error de timeout |
| 6 | Archivo vacío | File size 0 | Error de validación previo |

### 9.6 Frontend: `DashboardPage.jsx`

| # | Test | Escenario | Expected |
|---|---|---|---|
| 1 | Render inicial | Sin imagen seleccionada | Drop zone visible |
| 2 | Seleccionar imagen | File input change | Previsualización mostrada |
| 3 | Click "Analizar" | Imagen seleccionada + click | Loading state + llamada API |
| 4 | Diagnóstico exitoso (clorosis) | API retorna detección | Badge, gauge, chips visibles |
| 5 | Diagnóstico exitoso (sana) | API retorna sin clorosis | Badge verde, mensaje "saludable" |
| 6 | Error en diagnóstico | API retorna error | Mensaje de error, fallback UI |
| 7 | Multiple detecciones | API retorna 2+ detecciones | Chips múltiples, recomendación |
| 8 | Imagen no soportada | Archivo no JPEG/PNG | Validación previa, error UI |

---

## 10. Esquema de Base de Datos (Diagnóstico)

```sql
-- Tabla: diagnosticos (fincas)
id                  UUID PRIMARY KEY,
nombre              VARCHAR(50) NOT NULL,       -- "INF-20260625-XXXX"
muestras_id         UUID NOT NULL REFERENCES muestras(id),
tenant_id           UUID,
estado              VARCHAR(20) DEFAULT 'PENDIENTE',  -- PENDIENTE|ACEPTADO|RECHAZADO
image_url           TEXT NOT NULL,
tiene_clorosis      BOOLEAN NOT NULL,
confianza           DECIMAL(5,4) NOT NULL,      -- 0.0000 a 1.0000
procesado_at        TIMESTAMP NOT NULL,
created_at          TIMESTAMP DEFAULT NOW(),
updated_at          TIMESTAMP DEFAULT NOW()

-- Tabla: muestras (fincas)
id                  UUID PRIMARY KEY,
lote_id             UUID NOT NULL REFERENCES lotes(id),
tenant_id           UUID,
latitud             DECIMAL(10,7) NOT NULL,
longitud            DECIMAL(10,7) NOT NULL,
created_at          TIMESTAMP DEFAULT NOW(),
updated_at          TIMESTAMP DEFAULT NOW()

-- Tabla: candidatos_reentrenamiento (fincas)
id                  UUID PRIMARY KEY,
diagnostico_id      UUID UNIQUE NOT NULL REFERENCES diagnosticos(id),
image_url           TEXT NOT NULL,
tiene_clorosis      BOOLEAN NOT NULL,
confianza           DECIMAL(5,4) NOT NULL,
motivo              TEXT,
rechazado_por_usuario_id UUID,
created_at          TIMESTAMP DEFAULT NOW()
```

---

## 11. Comandos

```bash
# Test del handler de diagnóstico (únicos tests existentes)
cd fincas && go test ./internal/presentation/handler/... -v -run "TestDiagnosticoHandler"

# Todos los tests de fincas
cd fincas && go test ./... -v

# Suite de integración de fincas
cd fincas && make run-test

# Probar inferencia YOLO manualmente
cd YOLOv11 && python probar.py

# Entrenar modelo YOLO
cd YOLOv11 && python entrenamieto.py

# Ejecutar image-service
cd image-service && make run

# Build + test image-service
cd image-service && make build && make test
```

---

## 12. Referencias

| Documento | Ubicación |
|---|---|
| Especificación Casos de Uso Fincas | `fincas/docs/specs/spec-application-casos-de-uso.md` |
| Arquitectura Fincas | `fincas/docs/specs/spec-arquitectura-fincas.md` |
| Domain Fincas | `fincas/docs/specs/spec-fincas-domain.md` |
| Infraestructura Diagnóstico | `fincas/docs/specs/spec-infrastructure-diagnostico.md` |
| Especificación Monitoreo | `spec/especificacion-monitoreo.md` |
| Dashboard Frontend | `frontend/src/pages/dashboard/DashboardPage.jsx` |
| YOLO API Service | `frontend/src/services/yoloApi.js` |
