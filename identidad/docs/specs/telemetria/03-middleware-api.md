# Mini-Spec 3: El Enchufe HTTP — Middleware API (Gin)

> **Propósito**: Interceptar TODAS las peticiones HTTP al borde del router Gin para capturar métricas de red sin tocar ningún handler.

---

## 1. ¿Dónde se enchufa?

En `internal/infrastructure/telemetry/middleware/` se crea un middleware Gin que se registra con `router.Use()` **antes que cualquier otro middleware** (JWT, rate limit, etc.).

**El "enchufe" es la cadena de middlewares de Gin.** Gin ejecuta los middlewares en orden de registro. El de telemetría debe ser el primero para:
1. Capturar el timestamp de entrada lo antes posible.
2. Generar el `trace_id` y ponerlo en el `context.Context`.
3. Que todos los middlewares y handlers posteriores tengan acceso al `trace_id`.

```
Request → TelemetryMiddleware → JWTMiddleware → Handler → Response
         (toma timestamp)      (lee trace_id)  (usa ctx)  (vuelve al middleware)
                                                          (calcula duración, envía evento)
```

---

## 2. Flujo del middleware

### Fase antes de `c.Next()` (pre-procesamiento):
1. Tomar timestamp de alta precisión (`time.Now()`).
2. Extraer o generar `trace_id`:
   - Si el header `X-Trace-ID` o `X-Cloud-Trace-Context` existe → reusarlo.
   - Si no → generar UUID v7 (reusar `shared/infrastructure/idgenerator`).
3. Generar `span_id` nuevo (UUID v7).
4. Extraer `client_ip`: `X-Forwarded-For` → `X-Real-IP` → `r.RemoteAddr`.
5. Almacenar `trace_id` en `context.Context` con clave privada no exportada.
6. Ejecutar `c.Next()`.

### Fase después de `c.Next()` (post-procesamiento):
1. Calcular `duration_ms` (diferencia contra timestamp inicial).
2. Leer `c.Writer.Status()` (código HTTP).
3. Obtener ruta enrutada con `c.FullPath()` (template, no URL raw).
4. Leer `User-Agent` (truncado a 128 chars) y `Content-Length`.
5. Aplicar reglas de nivel:
   - 2xx/3xx → INFO
   - 4xx (excepto 401/403) → WARN
   - 5xx → ERROR
   - Duración > 3s → WARN; > 10s → ERROR
   - Pánico recuperado → ERROR
6. Construir payload `LogPayload` con `log_type = "API"`.
7. Enviar a `BufferWriter.Write()`.

---

## 3. Conexión al buffer

El middleware recibe `BufferWriter` por constructor:

```go
func NewTelemetryMiddleware(writer buffer.BufferWriter, cfg Config) gin.HandlerFunc
```

**No conoce Kafka.** No conoce el ring buffer. Solo conoce la interfaz.

---

## 4. Archivos a crear

```
internal/infrastructure/telemetry/middleware/
├── middleware.go       ← construcción del middleware Gin
├── middleware_test.go  ← tests con servidor Gin de prueba
└── config.go           ← umbrales (max_duration_warning_ms, etc.)
```
