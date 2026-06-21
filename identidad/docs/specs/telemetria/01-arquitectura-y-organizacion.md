# Mini-Spec 1: Fundación — Arquitectura y Organización del Código

> **Propósito**: Definir DÓNDE vive cada pieza del sistema de telemetría y CÓMO se relacionan, antes de implementar nada.

---

## 1. Árbol de paquetes (el mapa)

```
internal/
└── infrastructure/
    └── telemetry/
        ├── buffer/           ← Mini-Spec 2 (Ring Buffer + Kafka)
        ├── middleware/       ← Mini-Spec 3 (Middleware Gin)
        ├── decorator/        ← Mini-Spec 4 (Decoradores de casos de uso)
        └── gormplugin/       ← Mini-Spec 5 (Plugin GORM)
```

**Regla fundamental**: El `domain` NO importa `telemetry`. La `application` NO importa `telemetry` directamente. Solo el `registry/` y `infrastructure/` saben de telemetría.

---

## 2. El contrato fundamental: `BufferWriter`

Es la **única interfaz** que cruza capas. Todo punto de captura (middleware, decorador, plugin) depende de esta interfaz, NO del struct concreto:

```
BufferWriter <<interface>>
─────────────────────────────
+ Write(evento []byte, prioridad Prioridad) error
  → nil si se encoló
  → ErrBufferLleno si se descartó
  → ErrPrioritarioTimeout si el segmento prioritario está saturado
```

**Prioridad** es un enumerado simple: `Alta (0)`, `Media (1)`, `Baja (2)`.

**Esta interfaz vive en**: `internal/infrastructure/telemetry/buffer/writer.go`

**Nadie fuera de `buffer/` necesita saber que existe un ring buffer, Kafka, ni Prometheus.**

---

## 3. El payload unificado (el formato del dato)

Vive en `internal/infrastructure/telemetry/payload.go` como un struct plano:

```
LogPayload
├── log_type: "API" | "NEGOCIO" | "BD"
├── level: "INFO" | "WARN" | "ERROR"
├── timestamp: ISO8601 UTC
├── trace_id: UUID v7
├── span_id: UUID v7
├── service_name: "identidad"
├── environment: "dev" | "staging" | "production"
│
├── [si API] → ApiFields
│   ├── method, path, status_code, duration_ms
│   ├── client_ip (anonimizada), user_agent, content_length
│
├── [si NEGOCIO] → NegocioFields
│   ├── use_case, command (safe-print), result
│   ├── user_id, details (map), duration_usecase_ms
│
└── [si BD] → BdFields
      ├── operation, table, duration_ms
      ├── rows_affected, error_sql_state, query_hash
```

Las reglas de `level` (INFO/WARN/ERROR) las define cada punto de captura según sus propias condiciones (ver minispecs 3, 4, 5). No hay un clasificador central.

---

## 4. Mapa de conexiones (quién depende de quién)

```
middleware (Gin) ──────┐
decorator (app)  ──────┤──→ BufferWriter ──→ RingBuffer ──→ KafkaProducer
gormplugin (BD)  ──────┘                      │
                                              └──→ métricas Prometheus
```

**Ningún punto de captura conoce a Kafka.** Ningún punto de captura sabe si el buffer es un ring buffer, un channel, o una cola Redis. Solo conocen `BufferWriter`.

---

## 5. Hook de inicialización (dónde se prende todo)

La telemetría se inicializa en **dos puntos**, ambos en `internal/registry/`:

| Qué | Dónde se declara |
|-----|-----------------|
| Crear buffer + productor Kafka | `NewRegistry()` o nueva función `initTelemetry()` |
| Registrar middleware en Gin | En `router.New()` (router existente) |
| Registrar plugin GORM | En `NewRegistry()`, `db.Use(plugin)` |
| Envolver casos de uso | En `NewRegistry()`, al construir facades |

---

## 6. Feature flag: el interruptor general

Todo el sistema se controla con `config.Telemetry.Enabled` (bool). Cuando es `false`:
- No se crea el buffer ni el productor Kafka.
- `BufferWriter` se reemplaza por un **noop** que implementa la misma interfaz y descarta silenciosamente.
- No se registra middleware, no se registra plugin, no se envuelven casos de uso.
- **Cero impacto en el código existente.**

El noop writer vive en: `internal/infrastructure/telemetry/buffer/noop_writer.go`

---

## 7. Secuencia de implementación recomendada

1. Mini-Spec 1 (este documento) — definir paquetes, interfaces, payload
2. Mini-Spec 2 — buffer + Kafka (el motor, sin puntos de captura todavía)
3. Mini-Spec 3 — middleware API (el más fácil, primer log visible)
4. Mini-Spec 4 — decoradores de negocio (el que toca las facades)
5. Mini-Spec 5 — plugin GORM
6. Mini-Spec 6 — registry, feature flag, pruebas integradas
