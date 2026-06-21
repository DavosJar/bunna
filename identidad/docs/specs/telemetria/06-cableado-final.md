# Mini-Spec 6: El Cableado Final — Registry, Feature Flag y Pruebas

> **Propósito**: Conectar todas las piezas (Mini-Specs 1–5) en el Registry, controlar todo con un feature flag, y definir la estrategia de pruebas.

---

## 1. Feature flag: `TELEMETRY_ENABLED`

**Dónde**: en `config.Config`:

```go
type Config struct {
    // ... config existente ...
    Telemetry struct {
        Enabled bool   `env:"TELEMETRY_ENABLED" default:"false"`
        // ... demás parámetros del buffer, Kafka, umbrales ...
    }
}
```

**Comportamiento en Registry:**

```
if cfg.Telemetry.Enabled:
    bufferWriter = buffer.NewRingBuffer(cfg.Telemetry)
    db.Use(gormplugin.NewTelemetryPlugin(bufferWriter, cfg))
    // decorar casos de uso
else:
    bufferWriter = buffer.NewNoopWriter()  // descarta todo silenciosamente
    // no registrar plugin
    // no decorar casos de uso
```

**El middleware Gin** se controla en el router:

```go
func NewRouter(handlers, cfg) *gin.Engine {
    // ...
    if cfg.Telemetry.Enabled {
        r.Use(middleware.NewTelemetryMiddleware(bufferWriter, cfg))
    }
    // ... otros middlewares ...
}
```

---

## 2. Cableado completo en Registry

### 2.1 Paso 1: Crear el BufferWriter

```go
func NewRegistry(db *gorm.DB, cfg *config.Config) *Registry {
    var telemetryWriter buffer.BufferWriter

    if cfg.Telemetry.Enabled {
        telemetryWriter = buffer.NewRingBuffer(cfg.Telemetry)
    } else {
        telemetryWriter = buffer.NewNoopWriter()
    }
    // ...
```

### 2.2 Paso 2: Registrar plugin GORM (si habilitado)

```go
    if cfg.Telemetry.Enabled {
        gormPlugin := gormplugin.NewTelemetryPlugin(telemetryWriter, cfg.Telemetry)
        db.Use(gormPlugin)
    }
```

### 2.3 Paso 3: Envolver casos de uso (si habilitado)

```go
    registroUseCase := uc_register.NewRegistrarUsuarioCasoDeUso(...)

    if cfg.Telemetry.Enabled {
        registroUseCase = decorator.NewDecoradorRegistro(registroUseCase, telemetryWriter)
    }

    // La facade se construye IGUAL
    authFacade := facades.NewAuthFacade(registroUseCase, ..., loginUseCase, ...)
```

**Todos los demás casos de uso siguen el mismo patrón.** El cableado de decoradores se hace justo después de crear cada caso de uso concreto, y antes de pasarlo a la facade.

### 2.4 Paso 4: Arrancar el consumidor

```go
    if cfg.Telemetry.Enabled {
        // El consumidor arranca en segundo plano al construir el buffer
        // El buffer.NewRingBuffer() inicia la goroutine internamente
    }
```

### 2.5 Paso 5: Graceful shutdown

```go
    // En main.go o donde se manejen señales:
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
    <-sigCh

    if cfg.Telemetry.Enabled {
        telemetryWriter.Shutdown(5 * time.Second) // drena y cierra Kafka
    }
```

---

## 3. Estrategia de pruebas

### 3.1 Pruebas unitarias del decorador

```go
// mock del caso de uso
mockUC := &mockRegistroUseCase{
    resultado: &RespuestaRegistrarUsuario{...},
    err: nil,
}
// mock del buffer
mockBuffer := &mockBufferWriter{}

decorator := NewDecoradorRegistro(mockUC, mockBuffer)

resp, err := decorator.Ejecutar(ctx, comando)

// Verificar:
// 1. mockUC fue llamado con los mismos argumentos
// 2. Se envió 1 evento al buffer
// 3. El evento tiene log_type = "NEGOCIO"
// 4. level = "INFO" (porque result = success)
// 5. duration_usecase_ms > 0
// 6. La respuesta es idéntica a la del mock
```

### 3.2 Pruebas unitarias del buffer

```go
// Verificar:
// 1. Inserción O(1) hasta capacidad máxima
// 2. Descarte por prioridad (llenario con INFOs, meter un ERROR → el ERROR se acepta)
// 3. Drenaje batch
// 4. NoopWriter descarta todo sin error
```

### 3.3 Pruebas de integración del middleware

```go
// Levantar servidor Gin con el middleware
// Hacer petición HTTP
// Verificar que se produjo 1 evento en el buffer con log_type = "API"
```

### 3.4 Pruebas de integración del plugin GORM

```go
// Abrir SQLite en memoria
// Registrar plugin
// Ejecutar query vía GORM
// Verificar que se produjo 1 evento con log_type = "BD"
```

### 3.5 Prueba de feature flag deshabilitado

```go
// Config con Telemetry.Enabled = false
// Registry produce facades sin decorar
// Middleware no está registrado
// Buffer es noop — cero riesgo de efectos secundarios
```

---

## 4. Checklist de implementación

| # | Mini-Spec | Archivos | Depende de |
|---|-----------|----------|------------|
| 1 | Fundación | Interfaces, payload, mapa de paquetes | Nada |
| 2 | Buffer + Kafka | ring.go, consumer.go, kafka_producer.go, metrics.go | Mini-Spec 1 |
| 3 | Middleware API | middleware.go + config.go | Mini-Spec 2 (BufferWriter) |
| 4 | Decorador negocio | decorator_*.go | Mini-Spec 2 (BufferWriter) + Mini-Spec 1 (ToLog en comandos) |
| 5 | Plugin GORM | plugin.go + config.go | Mini-Spec 2 (BufferWriter) |
| 6 | Cableado final | cambios en registry.go + config.go + router.go | Mini-Specs 1–5 |

---

## 5. Lo que NO cambia (principio de no regresión)

| Archivo | ¿Cambia? | Razón |
|---------|----------|-------|
| `internal/**/domain/*.go` | ❌ | El dominio no sabe de telemetría |
| `internal/**/application/usecases/*.go` | ❌ | Los casos de uso no se modifican |
| `internal/**/presentation/handler/*.go` | ❌ | Los handlers reciben la facade igual que antes |
| `internal/**/presentation/facades/*.go` | ❌ | Las facades reciben la interfaz, no les importa si es decorada o no |
| `internal/**/infrastructure/persistence/**/*.go` | ❌ | Los repositorios son interceptados via GORM automáticamente |
| `tests/` existentes | ❌ | Todos los tests existentes deben pasar sin cambios |

| Archivo | ¿Cambia? | Cambio |
|---------|----------|--------|
| `internal/registry/registry.go` | ✅ | Cableado condicional de decoradores, plugin y buffer |
| `internal/config/` (o como se llame) | ✅ | Nuevo campo `Telemetry.Enabled` + configuración del buffer |
| `internal/presentation/router/` | ✅ | Registro condicional del middleware |
| `main.go` | ✅ | Graceful shutdown del buffer |
```
