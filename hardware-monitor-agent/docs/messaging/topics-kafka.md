# Topics Kafka — Hardware Monitor Agent

Topics del agente de monitoreo de hardware: quién publica y JSON exacto que viaja a Kafka.

---

## Topics

| Topic | Estado | Propósito |
|-------|--------|-----------|
| `hardware.metrics` | ✅ Implementado | Snapshots periódicos de métricas del nodo (CPU, RAM, disco, red, contenedores) |
| `hardware.alerts` | ✅ Implementado | Alertas CRITICAL y eventos de recovery |

---

## Publicadores

| Emisor | Servicio | Componente | Se dispara en |
|--------|----------|------------|---------------|
| Agregador | `hardware-monitor-agent` | `Aggregator` → `KafkaPublisher` | Cada intervalo de recolección (default 10s) |
| Guardian | `hardware-monitor-agent` | `Guardian` → `KafkaPublisher` | Inmediatamente al detectar umbral CRITICAL o recovery |

**No hay consumidor dentro del agente.** Los topics se publican para sistemas externos de visualización y alertamiento.

> El agente corre como sidecar en Docker Swarm y publica métricas del host anfitrión.

---

## JSON que viaja a Kafka

### Topic: `hardware.metrics`

Emitido por el **Agregador** cada intervalo de recolección. Contiene todas las métricas del nodo en un snapshot.

```json
{
  "node_id": "swarm-node-01",
  "timestamp": "2026-06-20T14:30:00Z",
  "interval_ms": 10000,
  "cpu": {
    "usage_percent": 45.2,
    "cores": 8
  },
  "ram": {
    "total_mb": 16384,
    "used_mb": 8192,
    "available_mb": 8192,
    "usage_percent": 50.0
  },
  "disks": [
    {
      "mount": "/",
      "total_gb": 256.0,
      "used_gb": 128.0,
      "available_gb": 128.0,
      "usage_percent": 50.0
    }
  ],
  "net": {
    "interfaces": [
      {
        "name": "eth0",
        "received_bytes": 1024000,
        "transmitted_bytes": 512000,
        "received_bytes_per_sec": 102400.0,
        "transmitted_bytes_per_sec": 51200.0
      }
    ]
  },
  "containers": [
    {
      "container_id": "abc123",
      "cpu_shares": 512,
      "memory_limit_mb": 1024
    }
  ]
}
```

### Topic: `hardware.alerts`

Emitido por el **Guardian** inmediatamente cuando una métrica supera el umbral CRITICAL o cuando se recupera.

#### Alerta CRITICAL

```json
{
  "node_id": "swarm-node-01",
  "timestamp": "2026-06-20T14:30:05Z",
  "metric": "cpu",
  "severity": "Critical",
  "value": 92.5,
  "threshold": 90.0,
  "message": "CPU usage at 92.5% exceeds CRITICAL threshold of 90.0%",
  "previous_state": "normal",
  "event_type": "alert"
}
```

#### Evento de Recovery

```json
{
  "node_id": "swarm-node-01",
  "timestamp": "2026-06-20T14:35:00Z",
  "metric": "cpu",
  "severity": "Info",
  "value": 45.2,
  "threshold": 90.0,
  "message": "cpu usage returned to normal: 45.2%",
  "previous_state": "critical",
  "event_type": "recovery"
}
```

---

## Reglas de publicación

| Condición | Acción | Topic |
|-----------|--------|-------|
| Métrica supera umbral CRITICAL | Publica alerta (una sola vez, control de inundación) | `hardware.alerts` |
| Métrica en CRITICAL vuelve a nivel normal | Publica recovery | `hardware.alerts` |
| Métrica supera umbral WARN pero no CRITICAL | Solo log local, NO publica a Kafka | — |
| Cada intervalo de recolección | Publica snapshot completo | `hardware.metrics` |

---

## Lo que NO está en el JSON

El mensaje **no incluye**:
- Versión del agente
- Nombre del servicio (`hardware-monitor-agent`)
- Entorno (dev/staging/prod)
- Metadatos del sidecar

El único identificador del emisor es `node_id`, que identifica el nodo Swarm, no el servicio que publica.

---

## Referencias

| Archivo | Qué contiene |
|---------|-------------|
| `src/types.rs` | Estructuras `Snapshot` y `AlertEvent` (serializadas a JSON) |
| `src/kafka_publisher.rs` | Publicador Kafka usando `rdkafka` |
| `src/guardian.rs` | Evaluación de umbrales y control de inundación |
| `src/aggregator.rs` | Consolidación de métricas en snapshot |
| `src/main.rs` | Cableado del sistema (topics: "hardware.metrics", "hardware.alerts") |
| `spec/spec-hardware-monitor-agent.md` | Especificación completa del agente |
| `docker-compose.dev.yml` | Servicio Kafka para desarrollo |
