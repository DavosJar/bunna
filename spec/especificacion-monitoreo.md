---
title: Módulo de Monitoreo — servicio-monitoreo
version: 0.1
date_created: 2026-06-20
owner: Team Infraestructura
tags: monitoring, rust, clickhouse, grafana, kafka, time-series
---

# Introducción

Especificación para el módulo de monitoreo (`servicio-monitoreo`) que proporciona ingestión, almacenamiento y visualización de métricas de telemetría, métricas de hardware y alertas del sistema. El módulo utiliza un consumidor escrito en Rust que lee eventos desde Kafka (3 tópicos), los transforma y los persiste en ClickHouse como base de datos time-series. Grafana se encarga de la capa de visualización y dashboards.

## 1. Purpose & Scope

**Propósito:** Implementar un pipeline de monitoreo en tiempo real que capture, almacene y visualice métricas de telemetría del sistema, métricas de hardware de los nodos y alertas generadas por los servicios, permitiendo observar el estado operativo del ecosistema Bunna.

**Alcance:**
- Consumidor-ingestor en Rust que se suscribe a 3 tópicos de Kafka (`telemetry`, `hardware.metrics`, `hardware.alerts`)
- Escritura en ClickHouse con esquemas optimizados para consultas time-series
- Dashboards en Grafana para visualización de métricas y alertas
- Esquemas de retención de datos (TTL) por tipo de métrica
- Consultas analíticas sobre ventanas de tiempo

**Fuera de alcance:**
- Generación de alertas (las alertas son consumidas, no generadas por este módulo)
- UI de administración del módulo
- Autodiscovery de nodos (los nodos publican contra Kafka directamente)
- Orquestación del deployment de los agentes de monitoreo en los nodos

**Audiencia:** Desarrolladores, DevOps, SRE, administradores del sistema.

**Supuestos:**
- Los 3 tópicos de Kafka ya existen y reciben datos de los agentes/servicios del sistema
- ClickHouse está disponible como servicio (cluster o standalone)
- Grafana tiene acceso de red a ClickHouse como datasource
- Los mensajes en Kafka tienen formato JSON con una estructura conocida
- El módulo corre como un servicio independiente (no embebido en otros binarios)

## 2. Definitions

| Término | Definición |
|---------|------------|
| Telemetría | Datos de instrumentación de servicios: latencias, tasas de error, throughput, conteo de peticiones |
| Hardware Metrics | Métricas de recursos físicos/virtuales: CPU, memoria, disco, red por nodo |
| Hardware Alerts | Eventos de alerta generados por thresholds de hardware: CPU > 90%, disco lleno, nodo caído |
| Time-Series | Datos indexados por timestamp, optimizados para consultas sobre ventanas temporales |
| TTL | Time-To-Live — política de retención que elimina datos después de un período |
| Ingestor | Componente que consume datos de Kafka y los persiste en ClickHouse |
| Consumidor | Suscriptor a un tópico de Kafka que procesa mensajes |

## 3. Requirements, Constraints & Guidelines

### Requirements

- **REQ-001**: El consumidor DEBE suscribirse simultáneamente a los 3 tópicos de Kafka (`telemetry`, `hardware.metrics`, `hardware.alerts`).
- **REQ-002**: Por cada mensaje recibido, el ingestor DEBE parsear el JSON, validar la estructura y escribir el registro en la tabla ClickHouse correspondiente.
- **REQ-003**: Si un mensaje no cumple con el schema esperado, DEBE enviarse a un tópico `dead-letter` (`telemetry.dlq`) con el error de validación.
- **REQ-004**: El ingestor DEBE confirmar offsets (commit) solo después de que los datos hayan sido escritos exitosamente en ClickHouse (al menos una vez — at-least-once delivery).
- **REQ-005**: El módulo DEBE implementar un mecanismo de reintentos con backoff exponencial para fallos de escritura en ClickHouse.
- **REQ-006**: Cada tabla en ClickHouse DEBE tener una política de TTL definida para evitar crecimiento ilimitado de datos.
- **REQ-007**: El módulo DEBE exponer métricas de salud interna (mensajes consumidos, errores, latencia de escritura) en un endpoint `/health`.
- **REQ-008**: El módulo DEBE configurarse exclusivamente mediante variables de entorno (12-factor app).
- **REQ-009**: El módulo DEBE publicar métricas de su propio rendimiento (mensajes/segundo, errores/segundo) en un endpoint `/metrics` con formato Prometheus.
- **REQ-010**: El módulo DEBE soportar shutdown graceful: terminar consumo, esperar escrituras pendientes, y luego cerrar conexiones.

### Constraints

- **CON-001**: El inestor debe escribirse en Rust por requerimientos de rendimiento y seguridad de memoria.
- **CON-002**: Las tablas ClickHouse deben usar el motor `MergeTree` con particionado por fecha para consultas eficientes.
- **CON-003**: No debe haber pérdida de datos aceptada: si ClickHouse no responde, los mensajes no deben committearse en Kafka.
- **CON-004**: El módulo no debe depender de servicios externos más allá de Kafka, ClickHouse y el sistema de archivos local para logs.
- **CON-005**: Los dashboards de Grafana deben versionarse dentro del repositorio como archivos JSON.

### Guidelines

- **GUD-001**: Usar `rdkafka` (librería librdkafka) para el consumo de Kafka en Rust.
- **GUD-002**: Usar `clickhouse-rs` o el driver nativo HTTP de ClickHouse para escritura.
- **GUD-003**: Usar `serde` y `serde_json` para deserialización de mensajes JSON.
- **GUD-004**: Usar `tracing` para logging estructurado.
- **GUD-005**: Usar `tokio` como runtime asíncrono para el consumo y escritura concurrentes.
- **GUD-006**: Separar la lógica en módulos: `consumer`, `parser`, `writer`, `config`, `health`.

### Patterns

- **PAT-001**: Seguir el patrón de consumidor con `Stream` de rdkafka para procesamiento asíncrono por tópico.
- **PAT-002**: Usar el patrón `fan-out` con tareas tokio independientes por tópico para procesamiento paralelo.
- **PAT-003**: Usar `Circuit Breaker` para la conexión a ClickHouse: si falla N veces consecutivas, esperar un período antes de reintentar.

## 4. Interfaces & Data Contracts

### Flujo de Datos

```
Agentes/Servicios → Kafka (3 tópicos) → servicio-monitoreo (Rust) → ClickHouse → Grafana
                                            │
                                            └→ /health (monitoreo interno)
                                            └→ /metrics (Prometheus)
```

### Tópicos de Kafka

| Tópico | Formato | Particiones | Retención | Descripción |
|--------|---------|-------------|-----------|-------------|
| `telemetry` | JSON | 3 | 7 días | Métricas de servicios: latencia, error_rate, throughput |
| `hardware.metrics` | JSON | 3 | 7 días | Métricas de nodos: cpu, memory, disk, network |
| `hardware.alerts` | JSON | 1 | 30 días | Eventos de alerta generados por thresholds |

### Esquemas de Mensajes Kafka

#### Tópico: `telemetry`

```json
{
  "service_name": "string",
  "service_instance": "string (UUID)",
  "timestamp": "string (ISO 8601)",
  "metrics": {
    "latency_p50_ms": "float",
    "latency_p95_ms": "float",
    "latency_p99_ms": "float",
    "error_rate": "float",
    "requests_per_second": "float",
    "active_connections": "uint32"
  },
  "tags": {
    "environment": "string",
    "version": "string",
    "host": "string"
  }
}
```

#### Tópico: `hardware.metrics`

```json
{
  "node_id": "string (UUID)",
  "node_name": "string",
  "timestamp": "string (ISO 8601)",
  "metrics": {
    "cpu_percent": "float",
    "memory_used_bytes": "uint64",
    "memory_total_bytes": "uint64",
    "disk_used_bytes": "uint64",
    "disk_total_bytes": "uint64",
    "disk_io_read_bytes_per_sec": "float",
    "disk_io_write_bytes_per_sec": "float",
    "network_rx_bytes_per_sec": "float",
    "network_tx_bytes_per_sec": "float",
    "load_average_1m": "float",
    "load_average_5m": "float",
    "load_average_15m": "float"
  },
  "tags": {
    "datacenter": "string",
    "rack": "string",
    "host": "string"
  }
}
```

#### Tópico: `hardware.alerts`

```json
{
  "alert_id": "string (UUID)",
  "node_id": "string (UUID)",
  "node_name": "string",
  "timestamp": "string (ISO 8601)",
  "alert_type": "string (cpu_high | memory_high | disk_full | node_down | network_high)",
  "severity": "string (info | warning | critical)",
  "message": "string",
  "current_value": "float",
  "threshold_value": "float",
  "tags": {
    "datacenter": "string",
    "rack": "string"
  }
}
```

### Tablas ClickHouse — DDL

#### Tabla: `telemetry_metrics`

```sql
CREATE TABLE telemetry_metrics (
    service_name    String,
    service_instance UUID,
    timestamp       DateTime64(3),
    latency_p50_ms  Float64,
    latency_p95_ms  Float64,
    latency_p99_ms  Float64,
    error_rate      Float64,
    requests_per_second Float64,
    active_connections  UInt32,
    env             String,
    service_version String,
    host            String,
    ingested_at     DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toDate(timestamp)
ORDER BY (service_name, toStartOfHour(timestamp), timestamp)
TTL toDate(timestamp) + INTERVAL 90 DAY
SETTINGS index_granularity = 8192;
```

#### Tabla: `hardware_metrics`

```sql
CREATE TABLE hardware_metrics (
    node_id        UUID,
    node_name      String,
    timestamp      DateTime64(3),
    cpu_percent    Float64,
    memory_used_bytes  UInt64,
    memory_total_bytes UInt64,
    disk_used_bytes    UInt64,
    disk_total_bytes   UInt64,
    disk_io_read_bytes_per_sec  Float64,
    disk_io_write_bytes_per_sec Float64,
    network_rx_bytes_per_sec    Float64,
    network_tx_bytes_per_sec    Float64,
    load_average_1m  Float64,
    load_average_5m  Float64,
    load_average_15m Float64,
    datacenter       String,
    rack             String,
    host             String,
    ingested_at     DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toDate(timestamp)
ORDER BY (node_id, toStartOfHour(timestamp), timestamp)
TTL toDate(timestamp) + INTERVAL 90 DAY
SETTINGS index_granularity = 8192;
```

#### Tabla: `hardware_alerts`

```sql
CREATE TABLE hardware_alerts (
    alert_id      UUID,
    node_id       UUID,
    node_name     String,
    timestamp     DateTime64(3),
    alert_type    String,
    severity      String,
    message       String,
    current_value Float64,
    threshold_value Float64,
    datacenter    String,
    rack          String,
    ingested_at   DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toDate(timestamp)
ORDER BY (toStartOfHour(timestamp), severity, alert_type)
TTL toDate(timestamp) + INTERVAL 365 DAY
SETTINGS index_granularity = 8192;
```

#### Tabla auxiliar: `dead_letter_queue` (mensajes no procesables)

```sql
CREATE TABLE dead_letter_queue (
    topic         String,
    partition     Int32,
    offset        Int64,
    raw_message   String,
    error         String,
    failed_at     DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toDate(failed_at)
ORDER BY (topic, failed_at)
TTL toDate(failed_at) + INTERVAL 30 DAY
SETTINGS index_granularity = 8192;
```

### Endpoints del servicio

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/health` | GET | Estado del servicio: conexiones a Kafka y ClickHouse, últimos offsets consumidos |
| `/metrics` | GET | Métricas Prometheus: mensajes consumidos por tópico, errores, latencias |
| `/ready` | GET | Readiness probe: indica si el servicio está listo para recibir tráfico |
| `/live` | GET | Liveness probe: indica si el proceso está vivo |

### Respuesta esperada — `/health`

```json
{
  "status": "ok",
  "uptime_seconds": 123456,
  "kafka": {
    "connected": true,
    "topics": {
      "telemetry": { "last_offset": 1500, "lag": 0, "messages_total": 150000 },
      "hardware.metrics": { "last_offset": 3200, "lag": 5, "messages_total": 320000 },
      "hardware.alerts": { "last_offset": 450, "lag": 0, "messages_total": 45000 }
    }
  },
  "clickhouse": {
    "connected": true,
    "pending_writes": 0,
    "write_errors_total": 0
  }
}
```

### Estructura del Proyecto Rust

```
servicio-monitoreo/
├── Cargo.toml
├── Dockerfile
├── src/
│   ├── main.rs              # Entry point: init config, tokio runtime, signal handling
│   ├── config.rs            # Carga de variables de entorno (struct Config con serde)
│   ├── consumer.rs          # Consumidor Kafka: suscripción a 3 tópicos, dispatch
│   ├── parser.rs            # Parsers específicos por tópico (telemetry, hardware, alerts)
│   ├── writer.rs            # Writer a ClickHouse con batch y backpressure
│   ├── health.rs            # Endpoints HTTP /health, /ready, /live
│   ├── metrics.rs           # Métricas Prometheus internas
│   ├── error.rs             # Tipos de error del dominio
│   ├── models/
│   │   ├── telemetry.rs     # Struct TelemetryMetric (serde Deserialize)
│   │   ├── hardware.rs      # Struct HardwareMetric (serde Deserialize)
│   │   └── alert.rs         # Struct HardwareAlert (serde Deserialize)
│   └── db/
│       ├── mod.rs           # ClickHouse connection pool
│       └── queries.rs       # Queries parametrizadas de inserción
├── dashboards/
│   ├── telemetry-overview.json    # Dashboard Grafana: latencias, error rate, RPS
│   ├── hardware-overview.json     # Dashboard Grafana: CPU, memoria, disco por nodo
│   └── alerts-board.json          # Dashboard Grafana: alertas activas por severidad
└── tests/
    └── integration.rs      # Tests de integración con Kafka y ClickHouse mockeados
```

### Variables de Entorno

| Variable | Requerida | Default | Descripción |
|----------|-----------|---------|-------------|
| `KAFKA_BROKERS` | Sí | — | Lista de brokers Kafka (`host1:9092,host2:9092`) |
| `KAFKA_GROUP_ID` | No | `servicio-monitoreo` | Grupo de consumidores Kafka |
| `KAFKA_TOPICS` | No | `telemetry,hardware.metrics,hardware.alerts` | Tópicos a consumir |
| `KAFKA_AUTO_OFFSET_RESET` | No | `earliest` | Política de reset de offset (`earliest` / `latest`) |
| `CLICKHOUSE_HOST` | Sí | — | Host de ClickHouse |
| `CLICKHOUSE_PORT` | No | `8123` | Puerto HTTP de ClickHouse |
| `CLICKHOUSE_USER` | No | `default` | Usuario de ClickHouse |
| `CLICKHOUSE_PASSWORD` | No | — | Contraseña de ClickHouse |
| `CLICKHOUSE_DATABASE` | No | `monitoreo` | Base de datos en ClickHouse |
| `CLICKHOUSE_POOL_SIZE` | No | `10` | Tamaño del pool de conexiones |
| `CLICKHOUSE_BATCH_SIZE` | No | `1000` | Número de mensajes por batch de escritura |
| `CLICKHOUSE_FLUSH_INTERVAL_MS` | No | `5000` | Intervalo máximo de flush (ms) |
| `HTTP_HOST` | No | `0.0.0.0` | Host del servidor HTTP interno |
| `HTTP_PORT` | No | `8080` | Puerto del servidor HTTP interno |
| `RUST_LOG` | No | `info` | Nivel de logging (tracing) |
| `DEAD_LETTER_TOPIC` | No | `telemetry.dlq` | Tópico para mensajes no procesables |

### Validaciones

- **timestamp**: Debe ser ISO 8601 con milisegundos. Si no es válido, se rechaza el mensaje.
- **service_name / node_id**: No vacíos, máximo 255 caracteres.
- **metricas numéricas**: Deben ser valores finitos (no NaN, no Inf). Si no, el mensaje va a DLQ.
- **severity**: Solo acepta `info`, `warning`, `critical`.
- **alert_type**: Solo acepta valores del enum definido.

## 5. Acceptance Criteria

- **AC-001**: Given un mensaje válido en el tópico `telemetry`, When el consumidor lo procesa, Then el registro aparece en la tabla `telemetry_metrics` de ClickHouse con todos los campos correctos.
- **AC-002**: Given mensajes en los 3 tópicos simultáneamente, When el consumidor corre, Then los mensajes se procesan concurrentemente y se persisten en sus respectivas tablas sin interferencia.
- **AC-003**: Given un mensaje con formato JSON inválido en cualquier tópico, When el consumidor lo recibe, Then el mensaje se envía al tópico `dead-letter` y NO se persiste en ClickHouse.
- **AC-004**: Given ClickHouse está caído, When llegan mensajes a Kafka, Then el consumidor NO confirma los offsets (no hace commit), y reintenta la escritura con backoff exponencial.
- **AC-005**: Given el servicio recibe una señal SIGTERM, When hay escrituras pendientes, Then el servicio espera a que todas las escrituras se completen antes de cerrar conexiones (shutdown graceful).
- **AC-006**: Given el servicio corriendo, When se consulta `GET /health`, Then retorna HTTP 200 con el JSON de estado detallado mostrando conexiones activas.
- **AC-007**: Given el servicio corriendo, When se consulta `GET /metrics`, Then retorna métricas en formato Prometheus con mensajes consumidos, errores y latencias.
- **AC-008**: Given datos en `telemetry_metrics` de los últimos 7 días, When se consulta en Grafana el dashboard `telemetry-overview`, Then se visualizan correctamente las series de tiempo de latencia, error_rate y RPS.
- **AC-009**: Given la tabla `hardware_metrics` tiene datos con más de 90 días, When ClickHouse ejecuta su limpieza periódica, Then los registros viejos son eliminados según la política TTL.
- **AC-010**: Given Grafana configurado con datasource ClickHouse, When se importan los dashboards del directorio `dashboards/`, Then las 3 visualizaciones cargan sin errores y muestran datos.

## 6. Test Automation Strategy

- **Test Levels**: Unitarios (parser, validaciones), Integración (Kafka + ClickHouse reales o mockeados)
- **Framework**: `cargo test` estándar, `rust-integration-test` para tests de integración, `mockall` para mocks
- **Test Data Management**: Docker Compose con Kafka y ClickHouse para tests de integración locales
- **CI/CD Integration**: El pipeline debe ejecutar `cargo test` y `cargo clippy` antes de buildear la imagen Docker
- **Coverage Requirements**: Cobertura mínima del 80% en módulos `parser`, `writer`, `config`
- **Scenarios to test**:
  1. Consumo y escritura exitosa de cada tipo de mensaje
  2. Mensajes inválidos van a DLQ
  3. Reconexión a Kafka tras caída del broker
  4. Reconexión a ClickHouse tras caída del servicio
  5. Shutdown graceful con escrituras pendientes
  6. Backpressure cuando ClickHouse está lento
  7. Concurrencia: múltiples mensajes del mismo tópico en paralelo

## 7. Rationale & Context

### Por qué Rust

- **Rendimiento**: El consumo de Kafka y escritura a ClickHouse debe ser predecible y de baja latencia. Rust ofrece rendimiento de C/C++ con garantías de seguridad de memoria.
- **Concurrencia**: El modelo de async/await de Tokio permite manejar decenas de miles de mensajes/segundo con un uso eficiente de CPU.
- **Tipado**: Los schemas de Kafka se modelan con structs tipados de Serde, eliminando errores de parseo en runtime.
- **Despliegue**: Binario estáticamente linkeado (~10 MB) sin dependencias de runtime, ideal para contenedores Docker ligeros.

### Por qué ClickHouse

- **Inserción por lotes**: ClickHouse está optimizado para inserts en batch (1000+ filas por segundo por lote), ideal para el throughput de Kafka.
- **Compresión**: Las tablas MergeTree comprimen datos 5-10x, reduciendo costos de almacenamiento.
- **Consultas analíticas**: Las funciones de ventana y agregación de ClickHouse permiten consultar percentiles, promedios y tendencias en milisegundos.
- **TTL nativo**: La retención de datos se configura a nivel de tabla sin procesos externos de limpieza.

### Por qué Grafana

- **Estándar de la industria**: Grafana es el visualizador más adoptado para time-series y se integra nativamente con ClickHouse mediante plugin.
- **Alertas visuales**: Permite configurar alertas sobre las métricas sin necesidad de un sistema adicional.
- **Dashboards versionables**: Los JSON de dashboards se almacenan en el repositorio, permitiendo revisión y deployment declarativo.

### Arquitectura de consumidor por tópico

Cada tópico se procesa en una tarea Tokio independiente con su propio stream de Kafka. Esto permite:
- Aislar fallos: si un tópico tiene mensajes corruptos, no afecta a los otros
- Escalabilidad: cada tópico puede tener su propia configuración de batch size y flush interval
- Monitoreo granular: se trackean métricas por separado por tópico

## 8. Dependencies & External Integrations

### Runtime Dependencies

| Dependencia | Versión Mínima | Propósito |
|-------------|----------------|-----------|
| Apache Kafka | 3.0+ | Sistema de mensajería para eventos de monitoreo |
| ClickHouse | 22.0+ | Base de datos time-series |
| Grafana | 9.0+ | Visualización y dashboards |

### Rust Crates (Cargo.toml)

| Crate | Propósito |
|-------|-----------|
| `rdkafka` | Cliente Kafka con soporte de consumidor Stream |
| `clickhouse` | Driver HTTP nativo para ClickHouse (async) |
| `serde` / `serde_json` | Deserialización JSON de mensajes Kafka |
| `tokio` | Runtime asíncrono con soporte de señales |
| `tracing` / `tracing-subscriber` | Logging estructurado |
| `anyhow` | Manejo de errores ergonómico |
| `thiserror` | Definición de errores personalizados |
| `prometheus` / `axum-prometheus` | Métricas para endpoint /metrics |
| `axum` | Servidor HTTP para endpoints health/metrics |
| `uuid` | Parseo de UUIDs en mensajes |
| `chrono` | Manejo de timestamps ISO 8601 |
| `config` | Carga de configuración desde entorno |
| `deadpool` | Pool de conexiones a ClickHouse |

### Grafana Configuration

- Datasource: ClickHouse (plugin `grafana-clickhouse-datasource`)
- Default datasource name: `clickhouse-monitoreo`
- Dashboards importados desde directorio `dashboards/`

## 9. Examples & Edge Cases

### Ejemplo 1: Consumo y escritura exitosa

```
[consumer] Mensaje recibido del tópico 'telemetry' (offset 1500)
[parser]   Parseado OK: service_name="auth-service", latency_p99=245.3
[writer]   Batch de 1000 registros escrito en telemetry_metrics (12ms)
[consumer] Commit offset 1500 -> 2500
```

### Ejemplo 2: Mensaje inválido → Dead Letter Queue

```
[consumer] Mensaje recibido del tópico 'hardware.metrics' (offset 3200)
[parser]   ERROR: campo 'cpu_percent' es NaN
[writer]   Publicando en tópico 'telemetry.dlq' (offset original: 3200)
[consumer] Commit offset 3200 -> 3201 (el mensaje no se pierde, se redirige)
```

### Ejemplo 3: ClickHouse caído temporalmente

```
[consumer] Mensaje recibido del tópico 'telemetry' (offset 1800)
[writer]   ERROR: conexión a ClickHouse rechazada
[writer]   Reintento 1/5 en 2s... falla
[writer]   Reintento 2/5 en 4s... falla
[writer]   Reintento 3/5 en 8s... ClickHouse responde
[writer]   Batch escrito exitosamente (2300ms total)
[consumer] Commit offset 1800 -> 2800
```
**Importante**: Durante la caída de ClickHouse, NO se hace commit de offsets. Al reestablecerse, se reprocesan los mensajes no committeados. Esto garantiza at-least-once delivery.

### Ejemplo 4: Shutdown graceful

```
$ kill -SIGTERM <pid>
[main]     Señal SIGTERM recibida. Iniciando shutdown graceful...
[consumer] Deteniendo consumo de nuevos mensajes...
[consumer] 150 mensajes en cola de escritura pendientes
[writer]   Escribiendo 150 registros restantes...
[writer]   OK: 150 registros escritos
[consumer] Cerrando conexión Kafka...
[writer]   Cerrando pool de ClickHouse...
[main]     Shutdown completado. Bye.
```

### Edge Cases

| Caso | Comportamiento Esperado |
|------|------------------------|
| Llegada de mensajes duplicados | ClickHouse no tiene unique constraints. El consumidor es at-least-once; el dashboard debe usar deduplicación en consulta (`SELECT DISTINCT` o agregación por ventana) |
| Kafka sin particiones asignadas | El consumidor espera bloqueado hasta que se asignen particiones. /health reporta `lag: null` |
| ClickHouse lento (latencia > 1s por batch) | Se reduce el batch_size dinámicamente para evitar timeouts. Se reporta la métrica `clickhouse_write_latency` |
| Mensaje con timestamp futuro (> 1 hora) | Se rechaza el mensaje y se envía a DLQ con error "timestamp en el futuro" |
| Threshold de alerta con valores extremos | Se persiste igual. El dashboard puede tener thresholds visuales pero no hay filtrado en el ingestor |
| Ingesta de recuperación tras caída de Kafka | Al reconectar, el consumidor procesa los mensajes acumulados. Se debe monitorear el lag inicial |

## 10. Validation Criteria

Para dar por cumplida esta especificación:

1. **VC-001**: Existe el directorio `servicio-monitoreo/` con `Cargo.toml` compilable y estructura de módulos definida.
2. **VC-002**: `cargo build` compila sin errores ni warnings.
3. **VC-003**: `cargo test` pasa todas las pruebas unitarias y de integración.
4. **VC-004**: `cargo clippy` pasa sin errores (warnings permitidos solo si documentados).
5. **VC-005**: El binario Docker se construye con `Dockerfile` y el contenedor arranca con las variables de entorno configuradas.
6. **VC-006**: El endpoint `/health` responde HTTP 200 con JSON válido cuando Kafka y ClickHouse están disponibles.
7. **VC-007**: El endpoint `/metrics` responde HTTP 200 con formato Prometheus.
8. **VC-008**: Enviando un mensaje válido a cualquier tópico Kafka, el mensaje aparece en la tabla ClickHouse correspondiente en menos de 10 segundos.
9. **VC-009**: Enviando un mensaje inválido (JSON mal formado), el mensaje aparece en el tópico `telemetry.dlq`.
10. **VC-010**: Los 3 dashboards de Grafana se importan sin errores y muestran datos al consultar.
11. **VC-011**: El TTL está configurado en todas las tablas y los datos viejos se eliminan automáticamente.

## 11. Related Specifications / Further Reading

- [Arquitectura general del sistema](/home/alexis/procesos_software/bunna/docs/architecture.md) — Documento de arquitectura global de Bunna
- [Kafka topics definition](/home/alexis/procesos_software/bunna/infra/kafka/topics.md) — Definición de tópicos y schemas en el clúster Kafka
- [Grafana datasource setup](/home/alexis/procesos_software/bunna/infra/grafana/datasources.md) — Configuración del datasource ClickHouse en Grafana
- [Plugin Grafana ClickHouse](https://grafana.com/grafana/plugins/grafana-clickhouse-datasource/) — Plugin oficial de Grafana para ClickHouse
- [ClickHouse MergeTree documentation](https://clickhouse.com/docs/en/engines/table-engines/mergetree-family/mergetree) — Documentación del motor MergeTree y TTL
- [rdkafka crate](https://docs.rs/rdkafka/) — Documentación del cliente Kafka para Rust
