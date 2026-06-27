# Análisis del Microservicio `fincas`

> Fecha: 2026-06-27
> Propósito: Documentar la estructura actual, el esquema de base de datos y el manejo de IoT/nodos.

---

## 1. Estructura del Proyecto

```
fincas/
├── cmd/
│   ├── main.go                  ← Entry point (NO inicia el HTTP server — solo hace select {})
│   └── test/main.go             ← Test harness
├── internal/
│   ├── fincas/
│   │   ├── domain/              ← Entidades Finca, Lote, errores, repositorios, especificaciones
│   │   └── infrastructure/persistence/postgres/
│   │       ├── finca_model.go / lote_model.go     ← Modelos GORM
│   │       ├── finca_repositorio.go / lote_repositorio.go
│   │       └── unit_of_work.go
│   ├── diagnostico/
│   │   ├── domain/              ← Diagnostico, Muestra, Ubicacion, ResultadoInferencia, CandidatoReentrenamiento
│   │   └── infrastructure/persistence/postgres/
│   │       ├── diagnostico_model.go / muestra_model.go / candidato_model.go
│   │       ├── diagnostico_repositorio.go / muestra_repositorio.go / candidato_repositorio.go
│   │       └── unit_of_work.go
│   ├── application/
│   │   ├── auth_context.go      ← UsuarioID, TenantID, Permisos
│   │   ├── errores.go           ← ErrForbidden, ErrNotFound, ErrConflictoEstado, ErrValidacion
│   │   ├── event_publisher.go   ← Interfaz EventPublisher
│   │   ├── unit_of_work.go      ← Interfaz UnitOfWorkDiagnostico
│   │   └── usecases/            ← 17 casos de uso (command → usecase → salida)
│   ├── presentation/
│   │   ├── router/router.go     ← 17 rutas bajo /api/v1/
│   │   ├── handler/             ← 5 handlers (Finca, Lote, Muestra, Diagnostico, Reporte)
│   │   ├── facade/              ← 5 facades
│   │   ├── mapper/              ← 5 mappers (dominio ↔ DTO)
│   │   ├── dto/                 ← 6 DTOs
│   │   └── middleware/auth_middleware.go ← JWT Bearer
│   ├── registry/
│   │   ├── container.go         ← DI de toda la aplicación
│   │   └── migrate.go           ← AutoMigrate de las 5 tablas
│   ├── infrastructure/
│   │   ├── eventpublisher/console.go ← Solo escribe a stdout
│   │   ├── security/jwt/token_validator.go ← HS256, misma key que identidad
│   │   └── telemetry/           ← AOP decorators + Kafka (deshabilitado por defecto)
│   └── shared/
│       ├── domain/specifications.go  ← CriterioFiltro, Paginacion, Ordenacion, GeneradorID
│       ├── infrastructure/idgenerator/uuid_v7.go
│       └── presentation/             ← ApiResponse[T], ProblemaDetalle (RFC 9457)
├── NODO/Nodos/Nodos.ino         ← Firmware ESP32-CAM
├── docs/specs/                  ← 5 especificaciones técnicas
└── shared/presentation/api_response.go ← (duplicado, no usado)
```

Clean Architecture con 4 capas. Dos módulos de dominio (`fincas/` y `diagnostico/`) completamente aislados entre sí. La capa de aplicación es la única que los orquesta.

---

## 2. Esquema de Base de Datos

Auto-migración vía GORM en `internal/registry/migrate.go`. 5 tablas, sin FK físicas entre módulos.

### 2.1 `fincas`

| Columna | Tipo | Restricciones |
|---------|------|--------------|
| id | varchar(36) | PK |
| nombre | text | |
| ubicacion | text | |
| descripcion | text | |
| usuario_id | varchar(36) | INDEX |
| tenant_id | varchar(36) | INDEX, nullable |
| estado | varchar(20) | DEFAULT 'ACTIVA' |
| created_at | timestamp | |
| updated_at | timestamp | |

### 2.2 `lotes`

| Columna | Tipo | Restricciones |
|---------|------|--------------|
| id | varchar(36) | PK |
| finca_id | varchar(36) | INDEX |
| tenant_id | varchar(36) | INDEX |
| nombre | text | |
| area | double precision | |
| descripcion | text | |
| estado | varchar(20) | DEFAULT 'ACTIVO' |
| created_at | timestamp | |
| updated_at | timestamp | |

### 2.3 `muestras`

| Columna | Tipo | Restricciones |
|---------|------|--------------|
| id | varchar(36) | PK |
| lote_id | varchar(36) | INDEX |
| tenant_id | varchar(36) | INDEX |
| latitud | double precision | |
| longitud | double precision | |
| created_at | timestamp | |
| updated_at | timestamp | |

### 2.4 `diagnosticos`

| Columna | Tipo | Restricciones |
|---------|------|--------------|
| id | varchar(36) | PK |
| nombre | varchar(200) | |
| muestras_id | varchar(36) | INDEX |
| tenant_id | varchar(36) | INDEX |
| estado | varchar(20) | DEFAULT 'PENDIENTE' |
| image_url | text | |
| tiene_clorosis | boolean | |
| confianza | decimal(5,4) | |
| procesado_at | timestamp | |
| created_at | timestamp | |
| updated_at | timestamp | |

### 2.5 `candidatos_reentrenamiento`

| Columna | Tipo | Restricciones |
|---------|------|--------------|
| id | varchar(36) | PK |
| diagnostico_id | varchar(36) | UNIQUE INDEX |
| image_url | text | |
| tiene_clorosis | boolean | |
| confianza | decimal(5,4) | |
| motivo | text | nullable |
| rechazado_por_usuario_id | varchar(36) | |
| created_at | timestamp | |

---

## 3. IoT / Nodos

### 3.1 Lo que existe

#### ESP32-CAM Firmware (`fincas/NODO/Nodos/Nodos.ino`)

```
Board:    AI-THINKER ESP32-CAM
Camera:   OV2640, VGA 640x480, JPEG quality 12
Red:      WiFi (SSID/PASS configurables)
Endpoint: https://bunna-yolo.duckdns.org/api/v1/diagnostico
Auth:     X-Node-Key: bunna-fincaPrueba (hardcoded, NO validado)
Envío:    HTTP POST multipart, field "archivo", filename "foto.jpg"
Frecuencia: 10s (dev) / 5 min (producción deseada)
Sleep:    delay(), sin deep sleep
SSL:      No manejo explícito de certificados
```

El ESP32 envía la foto **directamente al YOLO API** (Python/FastAPI en duckdns.org), no a fincas.

#### Image-service (`image-service/`)

```
Lenguaje: Go
Función:  Watcher de directorio, redimensiona imágenes a 640px
Output:   MQTT topic "images/processed/{filename}" via Mosquitto
Broker:   tcp://localhost:1883 (dockerizado)
Estado:   Standalone, NO conectado a fincas ni a YOLO
```

### 3.2 Pipeline actual vs deseado

```
── ACTUAL ──
ESP32 → POST → YOLO API → JSON response al ESP32 (nadie persiste nada en BD)

── DESEADO (specs) ──
FLUJO EDGE:
  ESP32 → HTTP → YOLO API → RabbitMQ → fincas.RegistrarInferencia → BD
  image-service → MQTT → [preprocesador] → YOLO → RabbitMQ → fincas.RegistrarInferencia

FLUJO MANUAL:
  App → TomarMuestra → SolicitarDiagnosticoManual → RabbitMQ → [preprocesador]
    → YOLO → RabbitMQ → fincas.RegistrarInferencia → BD
```

Ambos flujos convergen en `RegistrarInferencia`, que ya está implementado como use case pero **nadie lo invoca** (no hay RabbitMQ consumer).

### 3.3 Lo que falta (gaps)

| Componente | Estado | Impacto |
|---|---|---|
| Modelo `Nodo` en BD (id, nombre, ubicación, estado, batería, última conexión, lote asociado, firmware version) | ❌ No existe | No hay registro de qué dispositivos están desplegados |
| Autenticación de dispositivos (X-Node-Key validado vs BD) | ❌ No existe | Cualquiera puede enviar fotos |
| Consumer MQTT en fincas | ❌ No existe | image-service publica a MQTT pero nadie consume |
| Consumer RabbitMQ en fincas | ❌ No existe | `RegistrarInferencia` no tiene trigger |
| Endpoint HTTP para recibir inferencias de YOLO (callback/webhook) | ❌ No existe | YOLO responde a quien llamó, no hay push a fincas |
| Asociación nodo ↔ lote/finca | ❌ No existe | No se sabe qué nodo está en qué lote |
| Deep sleep en ESP32 | ❌ Usa `delay()` | Consume batería constante en campo |

---

## 4. Estado por Capa

| Capa | Estado |
|---|---|
| **Dominio (fincas + diagnostico)** | ✅ 100% completo |
| **Casos de uso (17)** | ✅ 100% implementados |
| **Infraestructura (modelos, repos, UoW)** | ✅ 100% implementada |
| **Presentación (router, handlers, facades, mappers, DTOs)** | ✅ 100% implementada |
| **Error handler RFC 9457** | ✅ |
| **ApiResponse[T] con HATEOAS** | ✅ |
| **Auth JWT middleware** | ✅ (pero TenantID y Permisos siempre vacíos) |
| **EventPublisher** (ConsolePublisher) | ⚠️ Solo stdout, no RabbitMQ |
| **Servidor HTTP** (`router.Run()`) | ❌ Roto — `cmd/main.go` hace `select {}` |
| **Frontend React+Vite** | ❌ No existe |
| **RabbitMQ consumer** | ❌ No existe |
| **Integración con YOLO** (callback) | ❌ No existe |
| **Modelo de nodos IoT** | ❌ No existe |
| **Tests de integración** (PostgreSQL real) | ❌ No existen |
| **OpenAPI docs** | ❌ No existen |

---

## 5. Eventos del Sistema (Routing Keys)

Todos definidos en los use cases, publicados por `ConsolePublisher` (stdout):

| Evento | Routing Key | Use Case |
|--------|-------------|----------|
| FincaCreada | `fincas.v1.finca.creada` | registrarfinca |
| FincaDesactivada | `fincas.v1.finca.desactivada` | desactivarfinca |
| LoteCreado | `fincas.v1.lote.creado` | agregarlote |
| LoteEliminado | `fincas.v1.lote.eliminado` | eliminarlote |
| MuestraTomada | `diagnosticos.v1.muestra.tomada` | tomarmuestra |
| SolicitudDiagnosticoManual | `diagnosticos.v1.solicitud.diagnostico.manual` | solicitardiagnosticomanual |
| DiagnosticoCreado | `diagnosticos.v1.diagnostico.creado` | registrarinferencia |
| DiagnosticoAceptado | `diagnosticos.v1.diagnostico.aceptado` | aceptardiagnostico |
| DiagnosticoRechazado | `diagnosticos.v1.diagnostico.rechazado` | rechazardiagnostico |

---

## 6. Endpoints REST

Agrupados bajo `/api/v1/`, todos requieren JWT (Bearer), excepto `/health`.

| Método | Ruta | Handler |
|--------|------|---------|
| GET | `/health` | Salud |
| POST | `/api/v1/fincas` | RegistrarFinca |
| GET | `/api/v1/fincas` | ListarFincas |
| GET | `/api/v1/fincas/:id` | ObtenerFinca |
| PUT | `/api/v1/fincas/:id` | EditarFinca |
| POST | `/api/v1/fincas/:id/desactivar` | DesactivarFinca |
| POST | `/api/v1/fincas/:fincaID/lotes` | AgregarLote |
| GET | `/api/v1/fincas/:fincaID/lotes` | ListarLotesPorFinca |
| GET | `/api/v1/fincas/:fincaID/lotes/:id` | ObtenerLote |
| PUT | `/api/v1/fincas/:fincaID/lotes/:id` | EditarLote |
| POST | `/api/v1/fincas/:fincaID/lotes/:loteID/muestras` | TomarMuestra |
| GET | `/api/v1/fincas/:fincaID/lotes/:loteID/muestras` | ListarMuestrasPorLote |
| GET | `/api/v1/fincas/:fincaID/lotes/:loteID/reporte` | GenerarReportePorLote |
| POST | `/api/v1/lotes/:id/eliminar` | EliminarLote |
| GET | `/api/v1/lotes/:loteID/muestras` | ListarMuestrasPorLote |
| POST | `/api/v1/lotes/:loteID/muestras` | TomarMuestra |
| POST | `/api/v1/muestras/:muestraID/diagnosticos/manual` | SolicitarDiagnosticoManual |
| POST | `/api/v1/diagnosticos/:id/aceptar` | AceptarDiagnostico |
| POST | `/api/v1/diagnosticos/:id/rechazar` | RechazarDiagnostico |
