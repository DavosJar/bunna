# Especificación Técnica de Arquitectura — Sistema de Telemetría Asíncrona para Módulo IAM

> **Propósito**: Capturar eventos de telemetría de tres orígenes distintos (API, NEGOCIO, BD) con cero impacto en la latencia de las peticiones del usuario, y enviarlos a un único topic Kafka (`iam.logs`) mediante procesamiento asíncrono en segundo plano.

---

## 1. Estructura y Esquema del Payload JSON Unificado

### 1.1 Envoltorio Global (Campos Comunes)

Todo evento de log, independientemente de su origen, comparte los siguientes campos de identificación y trazabilidad:

| Campo | Tipo lógico | Propósito |
|---|---|---|
| `log_type` | Enumeración: `"API"`, `"NEGOCIO"`, `"BD"` | Filtro primario en Grafana. Define qué sección de datos esperar en el payload. |
| `level` | Enumeración: `"INFO"`, `"WARN"`, `"ERROR"` | Severidad del evento. Se determina por reglas específicas de cada tipo de log. |
| `timestamp` | ISO 8601 con zona horaria UTC | Momento exacto de la captura, generado en el punto de intercepción. |
| `trace_id` | UUID v7 | Correlación entre logs del mismo flujo. Se genera en el middleware HTTP y se propaga via `context.Context`. |
| `span_id` | UUID v7 | Identificador del paso específico dentro del trace. |
| `service_name` | Cadena fija: `"identidad"` | Permite filtrar por microservicio en Grafana multi-fuente. |
| `environment` | Cadena: `"dev"`, `"staging"`, `"production"` | Diferenciación de ambientes para filtros en dashboards. |

### 1.2 Sección Específica por Tipo de Log

#### 1.2.1 `log_type = "API"`

Se genera automáticamente en el middleware global del enrutador. Contiene métricas de red sin exponer datos sensibles del request.

| Campo | Tipo lógico | Regla de captura |
|---|---|---|
| `method` | Cadena | Verbo HTTP: `GET`, `POST`, `PUT`, `DELETE`, etc. Obtenido del request. |
| `path` | Cadena | Ruta enrutada (template con parámetros, no la URL raw). Ej: `/api/v1/auth/register`. |
| `status_code` | Entero | Código HTTP de respuesta (200, 201, 400, 500, etc.). |
| `duration_ms` | Flotante | Tiempo total desde que el request entra al middleware hasta que la respuesta se envía. Precisión de microsegundos. |
| `client_ip` | Cadena | IP del cliente. Anonimizada: último octeto enmascarado para IPv4 (ej: `192.168.1.xxx`). |
| `user_agent` | Cadena | Header `User-Agent`, truncado a 128 caracteres. |
| `content_length` | Entero | Tamaño del body de respuesta en bytes. |

**Reglas de determinación de `level` para `log_type = "API"`:**

| Condición | Level |
|---|---|
| Código de respuesta 2xx o 3xx | `INFO` |
| Código 4xx (error del cliente), excepto 401/403 en patrón normal | `INFO` |
| Código 401 o 403 con patrón de múltiples intentos (más de 5 en 1 minuto desde misma IP) | `WARN` |
| Código 4xx que no es 401/403 (ej: 404, 422, 429) | `WARN` |
| Código 5xx | `ERROR` |
| Duración mayor a `max_duration_warning_ms` (configurable, ej: 3s) | `WARN` |
| Duración mayor a `max_duration_error_ms` (configurable, ej: 10s) | `ERROR` |
| Pánico recuperado por Gin Recovery | `ERROR` |

#### 1.2.2 `log_type = "NEGOCIO"`

Se genera mediante el Decorador de Casos de Uso. Contiene datos de auditoría de dominio. Nunca incluye contraseñas, tokens, datos biométricos, ni información personal innecesaria.

| Campo | Tipo lógico | Regla de captura |
|---|---|---|
| `use_case` | Cadena | Nombre lógico del caso de uso ejecutado. Ej: `"RegistrarUsuario"`, `"IniciarSesion"`. |
| `command` | Cadena | Representación "safe-printable" del comando: se serializan los campos no sensibles. Los campos marcados como sensitivos (contraseñas, tokens) se reemplazan por `"***"`. |
| `result` | Cadena | Clasificación del resultado: `"success"`, `"validation_error"`, `"business_error"`, `"infrastructure_error"`. |
| `user_id` | Cadena | Identificador del usuario autenticado (si aplica). Vacío para flujos anónimos como registro o login inicial. |
| `details` | Mapa clave-valor | Campos adicionales contextuales que el caso de uso quiera registrar (ej: `tenant_id` creado, `rol_asignado`, `cantidad_sesiones_cerradas`). |
| `duration_usecase_ms` | Flotante | Tiempo de ejecución del método `Ejecutar` dentro del decorador. |

**Reglas de determinación de `level` para `log_type = "NEGOCIO"`:**

| Condición | Level |
|---|---|
| Ejecución exitosa sin anomalías (`result = "success"`) | `INFO` |
| Ejecución exitosa pero con condiciones atípicas documentadas (datos incompletos, ubicación inusual, recuperación de estado inconsistente) | `WARN` |
| Error de validación de reglas de negocio esperadas (contraseña no cumple política, correo duplicado) | `WARN` |
| Error de negocio grave (intento de acceso a recurso no autorizado que burló el middleware, manipulación de datos de otro usuario) | `ERROR` |
| Error de infraestructura (no pudo persistir, no pudo enviar email, timeout de repositorio) | `ERROR` |
| Error inesperado no clasificable (panic, error desconocido) | `ERROR` |

#### 1.2.3 `log_type = "BD"`

Se genera interceptando las operaciones del ORM (GORM) en la capa de infraestructura. Contiene métricas de queries sin exponer los datos de las filas.

| Campo | Tipo lógico | Regla de captura |
|---|---|---|
| `operation` | Cadena | Tipo de operación: `"SELECT"`, `"INSERT"`, `"UPDATE"`, `"DELETE"`, `"TRANSACTION_BEGIN"`, `"TRANSACTION_COMMIT"`, `"TRANSACTION_ROLLBACK"`. |
| `table` | Cadena | Nombre de la tabla involucrada. Obtenido del modelo GORM. |
| `duration_ms` | Flotante | Tiempo de ejecución de la query. |
| `rows_affected` | Entero | Cantidad de filas afectadas o retornadas. |
| `error_sql_state` | Cadena | Código de estado SQL si la query falló (ej: `"23505"` para unique violation). Solo presente si `level = ERROR`. |
| `query_hash` | Cadena | Hash SHA-256 de la consulta normalizada (sin valores literales) para agrupación y análisis de rendimiento. Nunca se registra la query raw. |

**Reglas de determinación de `level` para `log_type = "BD"`:**

| Condición | Level |
|---|---|
| Query exitosa con duración menor a `slow_query_threshold_ms` (ej: 200ms) | `INFO` |
| Query exitosa con duración entre `slow_query_threshold_ms` y `critical_query_threshold_ms` (ej: 200ms–1s) | `WARN` |
| Query exitosa que retorna un número inusualmente alto de filas (ej: > 1000 registros) | `WARN` |
| Transacción que dura más de `long_transaction_threshold_ms` (ej: > 5s) | `WARN` |
| Query que excede `critical_query_threshold_ms` (ej: > 1s) | `ERROR` |
| Query que falla con error SQL | `ERROR` |
| Transacción que hace rollback | `ERROR` |

---

## 2. Flujo de Datos y Desacoplamiento en Clean Architecture

### 2.1 Capa de Red (API) — Middleware Global del Enrutador

**Ubicación arquitectónica**: En el borde externo de la aplicación, dentro del paquete `internal/presentation/middleware`, al mismo nivel que el middleware JWT y Rate Limit existentes.

**Mecanismo de operación**:

1. **Registro temprano**: El middleware se registra usando `router.Use()` en la función `New()` del paquete `router`, **antes que cualquier otro middleware**. Esto garantiza que capture el tiempo de entrada lo más temprano posible en la cadena de middlewares de Gin.

2. **Fase de pre-procesamiento** (antes de `c.Next()`):
   - a. Toma un timestamp de alta precisión (monotonic clock, `time.Now()`).
   - b. Extrae o genera el `trace_id`: busca en headers `X-Trace-ID` o `X-Cloud-Trace-Context`. Si no existe, genera un nuevo UUID v7 (reutilizando el generador existente en `shared/infrastructure/idgenerator`).
   - c. Extrae `client_ip` siguiendo la misma lógica que el router actual: primero `X-Forwarded-For`, luego `X-Real-IP`, finalmente `r.RemoteAddr`.
   - d. Almacena el `trace_id` en el `context.Context` del request usando una clave privada no exportada para que fluya hacia las capas inferiores.
   - e. Ejecuta `c.Next()` para continuar la cadena de middlewares y handlers.

3. **Fase de post-procesamiento** (después de `c.Next()`):
   - a. Calcula `duration_ms` como la diferencia contra el timestamp inicial.
   - b. Lee `c.Writer.Status()` para obtener el código HTTP de respuesta.
   - c. Obtiene la ruta enrutada via `c.FullPath()` (contiene el template, no la URL con parámetros).
   - d. Lee el header `User-Agent` y el `Content-Length` del response.
   - e. Aplica las reglas de nivel (INFO/WARN/ERROR según status code y duración).
   - f. Construye el payload con `log_type = "API"` y lo envía al **Buffer de Logs Asíncrono** (sección 3).
   - g. El envío al buffer es **no bloqueante**: si el buffer está lleno, aplica la política de descarte definida en 3.3, nunca bloquea la goroutine del request.

**Pureza arquitectónica**: El middleware depende únicamente de:
- La interfaz abstracta del Buffer de Logs (inyectada en su constructor).
- Tipos del sistema (`http`, `time`, `context`).
- No conoce Kafka, no conoce serialización concreta, no conoce modelos de dominio.

### 2.2 Capa de Datos (BD) — Interceptor del ORM (GORM Hooks)

**Ubicación arquitectónica**: En la capa de infraestructura. Se implementa como un `gorm.Plugin` que se registra en la instancia de `*gorm.DB` dentro del paquete `internal/registry` al momento de construir el Registry.

**Mecanismo de operación**:

1. Se implementa un plugin GORM que se suscribe a los siguientes hooks del ciclo de vida de una query:
   - `gorm:after_query` — captura `SELECT`.
   - `gorm:after_create` — captura `INSERT`.
   - `gorm:after_update` — captura `UPDATE`.
   - `gorm:after_delete` — captura `DELETE`.
   - `gorm:after_begin_transaction` — captura inicio de transacción.
   - `gorm:after_commit` — captura commit exitoso.
   - `gorm:after_rollback` — captura rollback.

2. **Flujo de cada callback**:
   - a. Toma timestamp al inicio del callback.
   - b. Deja que GORM complete la operación (el callback se ejecuta después de la operación real).
   - c. Toma timestamp después de la operación.
   - d. Calcula `duration_ms`.
   - e. Extrae el nombre de la tabla del `*gorm.DB.Statement.Table`.
   - f. Extrae `RowsAffected` del `*gorm.DB.Statement.RowsAffected`.
   - g. Extrae el error del `*gorm.DB.Statement.Error`.
   - h. Si hay error SQL, mapea el error GORM a un código estándar de estado SQL (PostgreSQL error code).
   - i. Calcula `query_hash` aplicando SHA-256 al `Statement.SQL` normalizado (sin valores literales).
   - j. Lee el `trace_id` desde el `context.Context` asociado a la sesión GORM (previamente propagado desde el middleware).
   - k. Aplica reglas de nivel según duración, éxito/fallo, y filas afectadas.
   - l. Construye payload con `log_type = "BD"` y envía al Buffer de Logs.

3. **Registro en el Registry**: El plugin se instancia una sola vez y se pasa a `db.Use(plugin)` dentro de la función `NewRegistry()`. Todos los repositorios que usen esa instancia de `*gorm.DB` quedan automáticamente interceptados sin necesidad de modificar ningún repositorio existente.

**Pureza arquitectónica**:
- El interceptor vive exclusivamente en infraestructura.
- Los repositorios de dominio (`UsuarioRepositorio`, `SesionRepositorio`, etc.) no tienen código de logging ni saben que están siendo cronometrados.
- El dominio permanece completamente puro: ninguna entidad, valor objeto o interfaz de repositorio se ve afectada.

### 2.3 Capa de Aplicación (NEGOCIO) — Patrón Decorador (Wrapper)

**Ubicación arquitectónica**: Nuevo paquete `internal/application/decorators`. No se modifica ningún caso de uso existente.

**Fundamento del patrón**: El Decorador es una envoltura que implementa la **misma interfaz** que el caso de uso real, pero agrega comportamiento transversal (logging, métricas) antes y después de delegar en la instancia concreta. Como todos los casos de uso en la arquitectura actual dependen de interfaces definidas en el paquete `facades` (ej: `RegistroUseCase`, `LoginUseCase`), el decorador implementa exactamente esas interfaces y se interpone entre la facade y el concreto.

**Flujo arquitectónico completo**:

1. **Definición de la interfaz objetivo**: Cada caso de uso expone su contrato mediante una interfaz en el paquete `facades`:
   - `RegistroUseCase` — `Ejecutar(ctx, *ComandoRegistrarUsuario) (*RespuestaRegistrarUsuario, error)`
   - `LoginUseCase` — `Ejecutar(ctx, ComandoIniciarSesion) (*RespuestaIniciarSesion, error)`
   - Ídem para `RefreshUseCase`, `LogoutUseCase`, etc.

2. **Construcción del Decorador**: Se crea una estructura `DecoradorLogs[T any]` parametrizable, o una estructura concreta por cada interfaz. Esta estructura:
   - Implementa la misma interfaz que el caso de uso real.
   - Recibe en su constructor la instancia concreta del caso de uso y una referencia al Buffer de Logs.
   - Se construye en el Registry, envolviendo el caso de uso real:
     ```
     Antes: facadeImpl ← RegistroUseCase (concreto)
     Después: facadeImpl ← DecoradorLogsRegistro ← RegistroUseCase (concreto)
     ```

3. **Ejecución interceptada**: Secuencia dentro del método `Ejecutar` del decorador:
   - a. Recibe `ctx` y `cmd`.
   - b. Toma timestamp de alta precisión.
   - c. Construye una representación "safe" del comando: serializa a mapa/claves planas, excluyendo campos marcados como sensitivos. Cada comando expone un método `ToLog()` o se usa una convención de tags (ej: `json:"-"` + anotación personalizada `log:"sensitive"`).
   - d. Extrae el `trace_id` del `ctx`.
   - e. Extrae el `user_id` del `ctx` (previamente inyectado por el middleware JWT).
   - f. **Delega**: invoca `uc.Ejecutar(ctx, cmd)` — llama al caso de uso real.
   - g. Toma timestamp después de la ejecución.
   - h. Calcula `duration_usecase_ms`.
   - i. **Clasifica el `result`**:
      - `err == nil` → `"success"`
      - Error de validación (`ErrValidation`) → `"validation_error"`
      - Error de negocio (`ErrBusiness`) → `"business_error"`
      - Error de infraestructura (error de BD, error de red, timeout) → `"infrastructure_error"`
      - Cualquier otro error → `"unexpected_error"`
   - j. **Define el `level`** según las reglas de la sección 1.2.2.
   - k. Si `result = "success"`, extrae del objeto de respuesta campos no sensibles para poblar `details`.
   - l. Construye el payload con `log_type = "NEGOCIO"` y lo envía al Buffer de Logs.
   - m. Retorna exactamente los mismos valores `(respuesta, error)` que produjo el caso de uso real, **sin modificación alguna**.

4. **Transparencia absoluta para el caso de uso**: El `RegistrarUsuarioCasoDeUso` dentro de `internal/usuarios/application/usecases/register/`:
   - No tiene dependencia de Kafka.
   - No tiene dependencia del Buffer de Logs.
   - No tiene conocimiento de que está siendo cronometrado.
   - No contiene ninguna línea de código de logging o telemetría.
   - Su prueba unitaria existente sigue funcionando sin cambios porque el decorador es un componente separado que se prueba aparte.

5. **Extensibilidad**: El patrón se replica para todas las interfaces de casos de uso. Como todas siguen el mismo contrato general `Ejecutar(ctx, Comando) (Respuesta, error)`, se puede implementar un decorador genérico usando `reflect` o, más idiomáticamente en Go, un decorador concreto por cada interfaz con una estructura interna compartida para la lógica de logging.

**¿Por qué en la interfaz de la facade y no en el struct concreto?**
Porque `authFacadeImpl` depende de **interfaces** (`RegistroUseCase`). El decorador implementa la misma interfaz y se interpone entre la facade y el struct concreto. Esto respeta el Principio de Inversión de Dependencias (DIP): las abstracciones no dependen de los detalles, los detalles dependen de las abstracciones. La facade nunca sabe si está llamando al caso de uso real o a un decorador.

---

## 3. Mecanismo de Seguridad contra Saturación de Memoria (Backpressure)

### 3.1 Arquitectura del Buffer Asíncrono

Se introduce un componente central llamado **Buffer de Logs Asíncrono** que actúa como cola de mensajes en memoria entre los puntos de captura (middleware, decorador, callbacks BD) y el productor de Kafka. Este buffer vive en un nuevo paquete de infraestructura: `internal/infrastructure/telemetry/buffer`.

**Componentes del buffer**:

1. **Ring Buffer acotado**: Un arreglo de tamaño fijo, predefinido en configuración. Cada slot contiene un evento de log (el payload JSON completo, ya serializado como `[]byte`). Tamaño típico inicial: 10,000 slots. Para un payload promedio de ~512 bytes, esto representa aproximadamente 5 MB de RAM — **memoria fija, sin crecimiento dinámico**.

2. **Productores (puntos de captura)**: El middleware API, el decorador NEGOCIO y los callbacks BD intentan insertar su evento en el ring buffer mediante una operación **no bloqueante**. Si el buffer está lleno, aplican la política de descarte por prioridad (sección 3.3). **Nunca esperan ni bloquean la goroutine del request**.

3. **Consumidor (Goroutine de fondo)**: Una o más Goroutines que:
   - Drenan el ring buffer en lotes (batch de hasta `batch_size` eventos, ej: 100).
   - Publican los lotes en Kafka mediante el productor Sarama/Confluent.
   - Si Kafka está caído o lento, reintentan con backoff exponencial (base: 100ms, techo: 10s).
   - Si el buffer supera el `high_watermark_ratio` (85%), el consumidor **reduce el tamaño del batch** (de 100 a 25 eventos) para drenar más rápido y en lotes más pequeños.
   - Si el buffer supera el `critical_watermark_ratio` (95%), el consumidor **pasa a modo de envío individual** (batch size = 1) para maximizar la tasa de drenaje.

4. **Métricas de salud expuestas**: El buffer expone las siguientes métricas (via Prometheus):
   - `telemetry_buffer_fill_ratio` — porcentaje de ocupación (0.0 a 1.0).
   - `telemetry_events_dropped_total` — contador de eventos descartados.
   - `telemetry_kafka_publish_latency_ms` — latencia de publicación en Kafka.
   - `telemetry_kafka_publish_errors_total` — errores de publicación.
   - `telemetry_buffer_enqueued_total` — eventos encolados exitosamente.
   - `telemetry_consumer_alive` — gauge binario (1 = vivo, 0 = muerto).

### 3.2 Configuración del Buffer

| Parámetro | Valor sugerido | Efecto |
|---|---|---|
| `buffer_capacity` | 10,000 | Número máximo de eventos en memoria. |
| `batch_size` | 100 | Eventos por lote de publicación Kafka. |
| `flush_interval_ms` | 500 | Tiempo máximo que un lote puede esperar antes de publicarse (reduce latencia en baja carga). |
| `max_retries` | 3 | Reintentos de publicación Kafka antes de descartar el lote. |
| `kafka_publish_timeout_ms` | 2,000 | Timeout por operación de publicación individual. |
| `backoff_base_ms` | 100 | Tiempo base para backoff exponencial en fallos de Kafka. |
| `backoff_max_ms` | 10,000 | Tiempo máximo de backoff. |
| `high_watermark_ratio` | 0.85 | Porcentaje de llenado que activa "drenaje rápido" (batch reducido). |
| `critical_watermark_ratio` | 0.95 | Porcentaje que activa modo de envío individual. |
| `prioritary_slots_ratio` | 0.20 | Porcentaje del buffer reservado exclusivamente para eventos de alta prioridad (ERROR + WARN de NEGOCIO). |

### 3.3 Política de Descarte por Prioridad

Cuando el ring buffer se llena porque Kafka no puede consumir al ritmo que se generan los logs, se debe proteger la disponibilidad de la aplicación principal. La política de descarte es **por prioridad jerárquica**:

#### Clasificación de prioridad de eventos

| Prioridad | Nivel | Tipo de log | Política |
|---|---|---|---|
| **Alta** (nunca se descartan) | `ERROR` | API, NEGOCIO, BD | Slot reservado en segmento prioritario. Si el segmento prioritario está lleno, la goroutine productora **espera hasta 100ms** (con timeout) a que el consumidor libere espacio. Esto es la ÚNICA situación que puede generar una mínima latencia adicional, y está diseñada para ser extremadamente rara (solo ocurre si incluso los errores no pueden drenarse). |
| **Alta** (nunca se descartan) | `WARN` | NEGOCIO | Misma política que ERROR. Estos eventos tienen el mayor valor de auditoría de negocio. |
| **Media** (se descartan bajo presión) | `WARN` | API, BD | Se descartan cuando `fill_ratio > 0.85`. El evento nuevo reemplaza al más antiguo dentro del mismo segmento. |
| **Baja** (se descartan primero) | `INFO` | API, BD, NEGOCIO | Se descartan cuando `fill_ratio > 0.70`. El evento nuevo reemplaza al más antiguo dentro del mismo segmento. |

#### Segmentación del Ring Buffer

El buffer se divide en dos segmentos lógicos:
- **Segmento prioritario** (20% de la capacidad total): Exclusivo para eventos de prioridad Alta. Soporta inserción con bloqueo controlado (timeout 100ms).
- **Segmento general** (80% de la capacidad total): Para eventos de prioridad Media y Baja. Soporta inserción circular con sobrescritura del más antiguo (overwrite).

Esta segmentación garantiza que:
- Los errores nunca se pierdan por saturación de INFOs.
- El consumo de memoria sea fijo y predecible (~5-10 MB).
- La aplicación principal nunca se vea forzada a hacer paging ni GC excesivo.

### 3.4 Protección de la Aplicación Principal

El diseño completo garantiza **impacto de 0ms en las peticiones del usuario** bajo cualquier condición:

1. **Inserción O(1) no bloqueante**: La inserción en el ring buffer usa operaciones atómicas (atomicos en Go: `sync/atomic`) y un diseño libre de locks o con sharded mutexes de muy baja contención. El tiempo de inserción es del orden de ~100-500 nanosegundos, despreciable frente a cualquier operación de dominio.

2. **Memoria fija y predecible**: El buffer tiene un tamaño fijo de 10,000 slots (~5-10 MB). No hay crecimiento dinámico del heap por retención de eventos. No hay fugas de memoria posibles por acumulación de logs no publicados.

3. **El consumidor no acumula**: Si Kafka falla, la Goroutine consumidora reintenta con backoff pero **descarta el lote después de `max_retries` intentos**. Esto evita que los mensajes se acumulen en un canal interno o buffer secundario.

4. **Timeouts estrictos**: Todas las operaciones contra Kafka tienen timeout configurable (2 segundos por defecto). Si Kafka no responde en ese tiempo, la operación se aborta y se reintenta.

5. **Watchdog de salud**: Una Goroutine separada monitorea el estado del consumidor principal. Si detecta que el consumidor lleva más de `consumer_heartbeat_timeout` (30 segundos) sin publicar exitosamente un lote, reinicia la Goroutine consumidora y registra una alerta en los logs estándar (stdout/stderr) para que el orquestador (Docker/K8s) la capture.

6. **Desconexión segura (graceful shutdown)**: Al recibir señal de terminación (SIGTERM/SIGINT), el sistema:
   - Detiene la aceptación de nuevos eventos en los puntos de captura.
   - Drena los eventos restantes del buffer hacia Kafka (con timeout máximo de 5 segundos).
   - Cierra el productor de Kafka.
   - Retorna el control al main para completar el shutdown graceful.

### 3.5 Integración con Prometheus y Grafana

Los siguientes indicadores se exponen como métricas Prometheus para configurar dashboards y alertas en Grafana:

| Métrica | Tipo | Umbral de alerta | Acción recomendada |
|---|---|---|---|
| `telemetry_buffer_fill_ratio` | Gauge (0-1) | > 0.85 → Warning<br>> 0.95 → Critical | Revisar salud del cluster Kafka. Escalar número de particiones del topic `iam.logs`. |
| `telemetry_events_dropped_total` | Counter | > 0 en última ventana de 5 min | Aumentar `buffer_capacity` o escalar Kafka. Revisar si el rate de eventos es sostenido. |
| `telemetry_kafka_publish_latency_ms` | Histograma | p99 > 1000 ms | Kafka está respondiendo lento. Revisar brokers, red, y carga. |
| `telemetry_kafka_publish_errors_total` | Counter | > 3 en 1 minuto | Kafka podría estar caído o inaccesible. Verificar conectividad de red. |
| `telemetry_consumer_alive` | Gauge (0/1) | == 0 | El consumidor de logs se detuvo. Activar reinicio automático y alertar al equipo. |

---

## 4. Diagrama de Flujo Completo de Datos

```
Petición HTTP
     │
     ▼
┌─────────────────────────────────────┐
│ Middleware API (Gin)                │
│   • Toma timestamp, genera trace_id │
│   • Extrae client_ip, method, path  │
│   • Inyecta trace_id en ctx         │
│   • c.Next() ──────────────────────────┼──┐
│   • Después de c.Next():            │  │
│     • Calcula duration_ms            │  │
│     • Lee status_code, content_length│  │
│     • Determina level (INFO/WARN/ERR)│  │
│     • Envía payload "API" al Buffer  │  │
└─────────────────────────────────────┘  │
                                         ▼
                               ┌───────────────────────────────────┐
                               │ Handler → Facade (authFacadeImpl) │
                               │ • Sin cambios en esta capa        │
                               └──────────────┬────────────────────┘
                                              ▼
                               ┌───────────────────────────────────────┐
                               │ Decorador NEGOCIO (wrapper)           │
                               │ • Toma timestamp                      │
                               │ • Serializa cmd (safe-print)          │
                               │ • Extrae trace_id, user_id del ctx    │
                               │ • Delegar: uc.Ejecutar(ctx, cmd)      │
                               │ • Toma timestamp post-ejecución       │
                               │ • Clasifica result y level            │
                               │ • Envía payload "NEGOCIO" al Buffer   │
                               │ • Retorna (respuesta, error) intactos │
                               └──────────────┬────────────────────────┘
                                              ▼
                               ┌──────────────────────────────────┐
                               │ Caso de Uso real                 │
                               │ (RegistrarUsuarioCasoDeUso,      │
                               │  IniciarSesionCasoDeUso, etc.)   │
                               │ • Sin cambios. No sabe del log.  │
                               └──────────────┬───────────────────┘
                                              │
                                              ▼ (invoca repositorios)
                               ┌──────────────────────────────────┐
                               │ GORM Hooks (BD interceptor)      │
                               │ • after_query, after_create, etc. │
                               │ • Toma timestamp                 │
                               │ • Mide duration_ms               │
                               │ • Extrae table, rows_affected    │
                               │ • Si error → error_sql_state     │
                               │ • query_hash (SHA-256 sin valores)│
                               │ • Determina level                │
                               │ • Envía payload "BD" al Buffer    │
                               └──────────────┬───────────────────┘
                                              │
                                              ▼
                               ┌──────────────────────────────────────┐
                               │                                      │
                               │   Ring Buffer Acotado                │
                               │   ┌──────────────────────────────┐   │
                               │   │ Segmento Prioritario (20%)    │   │
                               │   │ • ERROR de todos los tipos    │   │
                               │   │ • WARN de NEGOCIO             │   │
                               │   │ • Inserción con timeout 100ms │   │
                               │   └──────────────────────────────┘   │
                               │   ┌──────────────────────────────┐   │
                               │   │ Segmento General (80%)        │   │
                               │   │ • INFO de todos los tipos     │   │
                               │   │ • WARN de API y BD            │   │
                               │   │ • Sobrescritura circular      │   │
                               │   └──────────────────────────────┘   │
                               │                                      │
                               │   Capacidad: 10,000 slots (~5 MB)    │
                               │   Inserción: O(1) no bloqueante      │
                               └──────────────┬───────────────────────┘
                                              │
                                              │ (Goroutine consumidora)
                                              ▼
                               ┌──────────────────────────────────────┐
                               │ Productor Kafka (batch async)        │
                               │ • Batch: 100 eventos o 500ms         │
                               │ • Topic: "iam.logs"                  │
                               │ • Key: trace_id (orden por request)  │
                               │ • Backoff exp: 100ms → 10s           │
                               │ • Timeout por publish: 2s            │
                               │ • Si falla tras max_retries, descarta│
                               └──────────────────────────────────────┘
                                              │
                                              ▼
                                     ┌─────────────────┐
                                     │    Kafka         │
                                     │  topic: iam.logs │
                                     │  particiones: 6  │
                                     └─────────────────┘
                                              │
                                              ▼
                                     ┌─────────────────┐
                                     │    Grafana       │
                                     │  (fuente: Kafka) │
                                     │  filtros:        │
                                     │  • log_type      │
                                     │  • level         │
                                     │  • service_name  │
                                     │  • trace_id      │
                                     └─────────────────┘
```

---

## 5. Puntos de Extensión y Decisiones Arquitectónicas

### 5.1 Trazabilidad Distribuida (trace_id)
- El `trace_id` se genera **una sola vez por petición HTTP** en el middleware API.
- Se propaga a través del `context.Context` usando una clave de contexto privada (tipo no exportado) para evitar colisiones.
- Los decoradores y los hooks BD leen el `trace_id` del `ctx`. Si no existe (caso de pruebas o invocación directa del caso de uso sin HTTP), generan uno nuevo o usan un valor por defecto.
- Se usa como **key de Kafka** para garantizar que todos los eventos de una misma petición caigan en la misma partición y se consuman en orden.

### 5.2 Despliegue y Feature Flag
- El sistema completo se controla mediante una bandera `TELEMETRY_ENABLED` en la configuración (`config.Config`).
- Cuando está deshabilitada:
  - El middleware API se omite del router.
  - El plugin GORM no se registra en la base de datos.
  - Los decoradores no se construyen en el Registry; las facades reciben los casos de uso directamente.
  - El buffer y el consumidor no se inicializan.
- Esto permite habilitar/deshabilitar la telemetría sin cambiar una sola línea de código de dominio o aplicación.

### 5.3 Pruebas
- **Pruebas unitarias del decorador**: Se construye un mock del caso de uso y un mock del buffer. Se verifica que el decorador llame al caso de uso, capture la duración, clasifique correctamente el nivel y envíe el payload al buffer.
- **Pruebas unitarias del buffer**: Se verifica la inserción, la política de descarte, y el drenaje.
- **Pruebas de integración del middleware**: Se levanta un servidor de pruebas Gin, se ejecutan peticiones, y se verifica que el middleware produzca eventos en el buffer.
- **Pruebas de los hooks BD**: Se usa una base de datos de prueba (SQLite en memoria o PostgreSQL en contenedor Docker) y se verifica que las queries produzcan eventos.

### 5.4 Kafka — Topología
- **Topic único**: `"iam.logs"`.
- **Particiones**: 6 (mínimo recomendado para tolerancia a fallo y paralelismo).
- **Factor de replicación**: 3 (estándar en AWS MSK).
- **Key del mensaje**: `trace_id` (UUID v7) — garantiza orden por request.
- **Retención**: 7 días (configurable en el cluster).
- **Compresión**: `snappy` (balance entre ratio de compresión y velocidad).
- **ACKS**: `1` (el líder confirma — balance entre durabilidad y velocidad para logs).

### 5.5 Stack Tecnológico Propuesto
| Componente | Tecnología |
|---|---|
| Ring Buffer | Implementación propia en Go (paquete `buffer`) |
| Productor Kafka | `github.com/segmentio/kafka-go` o `github.com/IBM/sarama` |
| Métricas | `github.com/prometheus/client_golang` |
| Serialización | `encoding/json` estándar (sin reflection costs excesivos) |
| Generación de IDs | UUID v7 (reutilizar `shared/infrastructure/idgenerator`) |
