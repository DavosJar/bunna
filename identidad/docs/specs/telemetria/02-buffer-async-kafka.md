# Mini-Spec 2: El Corazón — Buffer Asíncrono y Publicación Kafka

> **Propósito**: Construir el motor de telemetría que recibe eventos de los puntos de captura, los bufferiza en memoria y los publica en Kafka SIN bloquear jamás la petición del usuario.

---

## 1. ¿Qué problema resuelve?

Los 3 puntos de captura (API, NEGOCIO, BD) generan eventos de log. Publicar cada uno directo a Kafka añadiría latencia a la petición del usuario (~1-5ms por publicación). El buffer permite:
- Publicación asíncrona en batch.
- Protección contra picos de tráfico.
- Descarte controlado si Kafka está caído.
- Memoria fija, sin fugas.

---

## 2. Componentes

### 2.1 Ring Buffer acotado

Un arreglo circular de tamaño fijo en `internal/infrastructure/telemetry/buffer/ring.go`:

| Parámetro | Valor sugerido | Por qué |
|-----------|---------------|---------|
| `buffer_capacity` | 10,000 slots | ~5MB para payload promedio de 512 bytes |
| Segmento prioritario | 20% (2,000 slots) | Solo ERROR + WARN de negocio |
| Segmento general | 80% (8,000 slots) | INFO + WARN de API/BD |

**Operaciones**:
- `Write(evento, prioridad)`: intenta insertar. Si el segmento está lleno:
  - Prioridad Alta → espera hasta 100ms con timeout (es la ÚNICA situación que puede bloquear).
  - Prioridad Media → descarta el más antiguo del segmento general si fill_ratio > 0.85.
  - Prioridad Baja → descarta el más antiguo si fill_ratio > 0.70.
- `ReadBatch(n)`: el consumidor extrae hasta `n` eventos para publicar.

**Thread-safety**: operaciones atómicas con `sync/atomic`, sin locks en el path crítico.

### 2.2 Consumidor (goroutine de fondo)

En `internal/infrastructure/telemetry/buffer/consumer.go`:

| Parámetro | Valor sugerido |
|-----------|---------------|
| `batch_size` normal | 100 eventos |
| `batch_size` reducido (fill > 85%) | 25 eventos |
| `batch_size` crítico (fill > 95%) | 1 evento (envío individual) |
| `flush_interval_ms` | 500ms |
| `max_retries` | 3 |
| `backoff_base_ms` | 100ms |
| `backoff_max_ms` | 10s |

**Comportamiento**:
- Drena el ring buffer en lotes.
- Publica en Kafka.
- Si Kafka falla: reintenta con backoff exponencial. Tras `max_retries` intentos fallidos, **descarta el lote**.
- Si el buffer se llena (>85%): reduce tamaño de batch para drenar más rápido.
- Si el buffer está crítico (>95%): envía de 1 en 1.

### 2.3 Productor Kafka

En `internal/infrastructure/telemetry/buffer/kafka_producer.go`:

| Parámetro | Valor sugerido |
|-----------|---------------|
| Topic | `"iam.logs"` |
| Key | `trace_id` (UUID v7) |
| Compression | `snappy` |
| ACKS | `1` |
| Timeout por publish | 2s |
| Particiones | 6 |
| Replicación | 3 |

Usa `segmentio/kafka-go` o `IBM/sarama` (el que ya esté disponible en el proyecto).

### 2.4 Watchdog

Una goroutine separada monitorea al consumidor. Si el consumidor no ha publicado exitosamente en 30 segundos, lo reinicia y escribe una alerta en stdout/stderr (para que el orquestador Docker/K8s la capture).

---

## 3. Prioridades y política de descarte

| Prioridad | Qué eventos | Política de descarte |
|-----------|------------|---------------------|
| **Alta** (nunca descarta) | ERROR de cualquier tipo, WARN de NEGOCIO | Slot en segmento reservado. Si lleno, bloquea 100ms máximo. |
| **Media** (descarta si >85%) | WARN de API, WARN de BD | Sobrescribe el más antiguo del segmento general. |
| **Baja** (descarta si >70%) | INFO de cualquier tipo | Sobrescribe el más antiguo del segmento general. |

**Efecto práctico**: en una tormenta de logs, los ERROR siempre se preservan. Los INFO se descartan primero.

---

## 4. Métricas Prometheus expuestas

| Métrica | Tipo | Propósito |
|---------|------|-----------|
| `telemetry_buffer_fill_ratio` | Gauge (0-1) | % de ocupación del buffer |
| `telemetry_events_dropped_total` | Counter | Eventos descartados acumulados |
| `telemetry_kafka_publish_latency_ms` | Histograma | Latencia de publicación |
| `telemetry_kafka_publish_errors_total` | Counter | Errores de Kafka acumulados |
| `telemetry_buffer_enqueued_total` | Counter | Eventos encolados exitosamente |
| `telemetry_consumer_alive` | Gauge (0/1) | Salud del consumidor |

---

## 5. Graceful Shutdown

Al recibir SIGTERM/SIGINT:
1. Se marca el buffer como "cerrado" — los `Write()` nuevos devuelven error.
2. Se drena el contenido restante (timeout máximo 5 segundos).
3. Se cierra el productor Kafka.
4. Se retorna el control para completar el shutdown general.

---

## 6. Archivos a crear

```
internal/infrastructure/telemetry/buffer/
├── writer.go           ← interfaz BufferWriter
├── noop_writer.go      ← implementación noop (feature flag false)
├── ring.go             ← ring buffer con segmentación
├── consumer.go         ← goroutine consumidora
├── kafka_producer.go   ← productor Kafka
├── metrics.go          ← métricas Prometheus
└── config.go           ← configuración del buffer
```
