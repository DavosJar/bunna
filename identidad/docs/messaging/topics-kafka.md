# Topics Kafka — Módulo Identidad

Topics del módulo identidad: quién publica, JSON que viaja a Kafka y eventos de dominio.

---

## Topics

| Topic | Estado | Propósito |
|-------|--------|-----------|
| `telemetry` | ✅ Implementado | Logs de telemetría unificados |
| `iam.logs` | 📄 Especificado | Alias del mismo topic |

> El nombre real del topic se configura vía `KAFKA_TOPIC`. Por defecto: `"telemetry"`.

---

## Publicadores

| Emisor | Servicio | Componente | Se dispara en |
|--------|----------|------------|---------------|
| Middleware HTTP | `identidad` | Middleware Gin | Cada request HTTP |
| Decorador negocio | `identidad` | Wrapper de casos de uso | login, registro, refresh, logout |
| Plugin GORM | `identidad` | Hook de BD | Cada operación sobre base de datos |

**No hay consumidor dentro de `identidad`.** El topic se publica para sistemas de observabilidad externos.

> ⚠️ **El JSON que viaja a Kafka no tiene ningún campo que identifique al emisor.** El campo `service_name` siempre es `"identidad"` sin importar qué componente emitió el mensaje.

---

## JSON que viaja a Kafka

### log_type = "API" (emitido por Middleware HTTP)

```json
{
  "log_type": "API",
  "level": "INFO",
  "timestamp": "2026-06-20T14:30:00.123Z",
  "trace_id": "0191f9c2-3a7e-7b00-8c4d-1e2f3a4b5c6d",
  "span_id": "0191f9c2-3a7e-7b00-8c4d-1e2f3a4b5c6e",
  "service_name": "identidad",
  "environment": "dev",
  "api": {
    "method": "POST",
    "path": "/api/v1/auth/login",
    "status_code": 200,
    "duration_ms": 45.23,
    "client_ip": "192.168.1.xxx",
    "user_agent": "Mozilla/5.0...",
    "content_length": 342
  }
}
```

### log_type = "NEGOCIO" (emitido por Decorador)

```json
{
  "log_type": "NEGOCIO",
  "level": "INFO",
  "timestamp": "2026-06-20T14:30:00.456Z",
  "trace_id": "0191f9c2-3a7e-7b00-8c4d-1e2f3a4b5c6d",
  "span_id": "0191f9c2-3a7e-7b00-8c4d-1e2f3a4b5c6f",
  "service_name": "identidad",
  "environment": "dev",
  "negocio": {
    "use_case": "IniciarSesion",
    "command": {
      "email": "user@example.com",
      "password": "***"
    },
    "result": "success",
    "user_id": "u_2a3b4c5d-6e7f-8a9b-0c1d-2e3f4a5b6c7d",
    "details": {
      "tenant_id": "t_abc123"
    },
    "duration_usecase_ms": 120.5
  }
}
```

### log_type = "BD" (emitido por Plugin GORM)

```json
{
  "log_type": "BD",
  "level": "INFO",
  "timestamp": "2026-06-20T14:30:00.789Z",
  "trace_id": "0191f9c2-3a7e-7b00-8c4d-1e2f3a4b5c6d",
  "span_id": "0191f9c2-3a7e-7b00-8c4d-1e2f3a4b5c70",
  "service_name": "identidad",
  "environment": "dev",
  "bd": {
    "operation": "SELECT",
    "table": "usuarios",
    "duration_ms": 3.2,
    "rows_affected": 1,
    "query_hash": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1"
  }
}
```

---

## Lo que NO está en el JSON

**No hay campo que indique quién emitió el mensaje.** Los únicos identificadores en el payload son:

| Campo | Valor | ¿Identifica al emisor? |
|-------|-------|------------------------|
| `service_name` | `"identidad"` | ❌ Siempre el mismo, no distingue componente |
| `log_type` | `"API"` / `"NEGOCIO"` / `"BD"` | ⚠️ Solo por convención se asocia a un componente |

Si otro servicio (ej: `fincas`, `image-service`) publicara al mismo topic, no habría forma de saber de qué servicio vino el mensaje.

---

## Eventos de dominio (futuro — aún NO se publican a Kafka)

| Evento | Disparador |
|--------|-----------|
| `UsuarioCreado` | Alta de usuario |
| `EstadoUsuarioCambiado` | Cambio de estado |
| `UsuarioBloqueado` | Bloqueo por intentos fallidos |
| `UsuarioActivado` | Activación de cuenta |
| `UsuarioInactivado` | Desactivación de cuenta |
| `CorreoVerificado` | Verificación de email |
| `ReenvioVerificacionSolicitado` | Reenvío de token |
| `EnlaceVerificacionExpirado` | Link de verificación expirado |

Actualmente viajan solo en memoria dentro de `identidad`. No se publican a Kafka.

---

## Referencias

| Archivo | Qué contiene |
|---------|-------------|
| `internal/infrastructure/telemetry/payload.go` | Estructura `LogPayload` |
| `internal/infrastructure/telemetry/buffer/kafka_producer.go` | Publicador Kafka |
| `internal/registry/registry.go` | Cableado del sistema |
| `docs/specs/SPEC-telemetria.md` | Especificación completa |
| `internal/usuarios/domain/usuario/eventos.go` | Eventos de dominio |
