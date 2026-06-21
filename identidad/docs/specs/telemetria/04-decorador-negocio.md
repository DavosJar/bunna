# Mini-Spec 4: El Enchufe de Negocio — Decoradores sobre Facades

> **Propósito**: Capturar eventos de auditoría de NEGOCIO envolviendo los casos de uso concretos SIN modificarlos. Este es el punto de conexión más importante: el decorador se enchufa **entre la Facade y el caso de uso real**.

---

## 1. La metáfora del enchufe

El flujo actual es:

```
Handler → Facade → CasoDeUsoConcreto
```

El decorado interpone una capa invisible:

```
Handler → Facade → Decorador (telemetría) → CasoDeUsoConcreto
```

**La interfaz del caso de uso (ej: `RegistroUseCase`) es el "enchufe".** El decorador implementa EXACTAMENTE la misma interfaz, y la Facade no sabe si está llamando al concreto o al decorador.

---

## 2. ¿Dónde vive?

Tres opciones, por orden de preferencia:

| Opción | Ruta | Argumento |
|--------|------|-----------|
| ✅ Recomendada | `internal/infrastructure/telemetry/decorator/` | La telemetría es infraestructura. El decorador vive acá. |
| Alternativa | `internal/application/decorators/` | Si se considera capa de aplicación. Se descarta porque violaría que `application` no sepa de telemetría. |

**Decisión**: va en `internal/infrastructure/telemetry/decorator/` porque depende de `BufferWriter` (infraestructura).

---

## 3. Mecanismo

### 3.1 Un decorador por interfaz de caso de uso

Cada interfaz en `internal/presentation/facades/*_facade.go` define contratos como:

```go
type RegistroUseCase interface {
    Ejecutar(ctx, *ComandoRegistrarUsuario) (*RespuestaRegistrarUsuario, error)
}
type LoginUseCase interface {
    Ejecutar(ctx, ComandoIniciarSesion) (*RespuestaIniciarSesion, error)
}
```

Para cada una, se crea su decorador:

```
decorator/
├── decorator_registro.go   ← implementa RegistroUseCase
├── decorator_login.go       ← implementa LoginUseCase
├── decorator_refresh.go     ← implementa RefreshUseCase
├── decorator_logout.go      ← implementa LogoutUseCase
├── ...
└── decorator_base.go        ← lógica compartida (clasificación de nivel, safe-print)
```

### 3.2 Flujo del decorador

1. Recibe `ctx` y `comando`.
2. Toma timestamp.
3. Construye representación "safe" del comando:
   - Cada comando expone un método `ToLog() map[string]any` que retorna solo campos no sensibles.
   - Alternativa: convención de tags `log:"sensitive"` en los campos del comando.
   - Contraseñas, tokens, datos biométricos → `"***"`.
4. Extrae `trace_id` del `ctx` (inyectado por el middleware de Mini-Spec 3).
5. Extrae `user_id` del `ctx` (inyectado por el middleware JWT).
6. **Delega**: llama `uc.Ejecutar(ctx, cmd)` al caso de uso real.
7. Toma timestamp post-ejecución.
8. Clasifica el resultado:
   - `err == nil` → `"success"`
   - Error de validación → `"validation_error"`
   - Error de negocio → `"business_error"`
   - Error de infraestructura → `"infrastructure_error"`
   - Otro → `"unexpected_error"`
9. Define `level` según tabla de la sección 1.2.2 del spec original.
10. Construye payload `LogPayload` con `log_type = "NEGOCIO"`.
11. Envía a `BufferWriter.Write()`.
12. Retorna EXACTAMENTE los mismos valores `(respuesta, error)` que produjo el caso de uso real.

### 3.3 Transparencia absoluta

El caso de uso concreto (`RegistrarUsuarioCasoDeUso`):
- No sabe que existe telemetría.
- No importa `BufferWriter`.
- No tiene código de logging.
- Sus tests existentes funcionan sin cambios.

---

## 4. ¿Dónde se enchufa en el Registry?

**Hoy en Registry:**

```go
// Registry crea casos de uso concretos
registroUseCase := uc_register.NewRegistrarUsuarioCasoDeUso(...)
loginUseCase := uc_sesiones_login.NewIniciarSesionCasoDeUso(...)

// Luego construye la facade con los casos de uso
authFacade := facades.NewAuthFacade(registroUseCase, ..., loginUseCase, ...)
```

**Con telemetría habilitada:**

```go
// Se crea el BufferWriter (ver Mini-Spec 2)
bufferWriter := buffer.NewRingBuffer(cfg)

// Se envuelven los casos de uso con decoradores
registroUseCase = decorator.NewDecoradorRegistro(registroUseCase, bufferWriter)
loginUseCase = decorator.NewDecoradorLogin(loginUseCase, bufferWriter)
// ... etc

// La facade recibe los decorados — no nota la diferencia
authFacade := facades.NewAuthFacade(registroUseCase, ..., loginUseCase, ...)
```

**La facade NO cambia.** El handler NO cambia. El caso de uso NO cambia. Solo el Registry cambia para interponer el decorador.

---

## 5. Método auxiliar `ToLog()` en cada comando

Cada comando necesita exponer sus campos no sensibles. Se agrega un método:

```go
// En el archivo del comando (ej: ComandoRegistro)
func (c ComandoRegistro) ToLog() map[string]any {
    return map[string]any{
        "nombre":   c.Nombre,
        "apellido": c.Apellido,
        "correo":   c.Correo,
        // Password NO se incluye
    }
}
```

**Esto no es infraestructura.** Es un método del propio comando (dominio/aplicación). No viola la arquitectura porque es pura transformación de datos.

---

## 6. Archivos a crear

```
internal/infrastructure/telemetry/decorator/
├── decorator_base.go         ← lógica compartida (clasificador de nivel, safe-print helper)
├── decorator_registro.go     ← envuelve RegistroUseCase
├── decorator_login.go        ← envuelve LoginUseCase
├── decorator_refresh.go      ← envuelve RefreshUseCase
├── decorator_logout.go       ← envuelve LogoutUseCase
├── decorator_switchtenant.go ← envuelve CambiarTenantCasoDeUso
├── ... (1 por interfaz de caso de uso)
└── decorator_base_test.go    ← tests del clasificador
```
