---
title: "Reporte de Implementación — Recuperación de Contraseña vía Email"
version: 1.0
date: 2026-05-22
owner: Equipo Identidad
status: EN_PROGRESO
tags: reporte, password, recovery, implementacion
---

# Reporte de Implementación: Recuperación de Contraseña vía Email

> **Propósito**: Evaluar el estado actual del código contra lo especificado en `spec-4-password-recovery.md`.

## 1. Resumen Ejecutivo

| Dimensión | Resultado |
|-----------|-----------|
| **Estado general** | EN_PROGRESO |
| **Tabla tokens_recuperacion** | Modelo GORM creado, sin AutoMigrate |
| **Modelo GORM** | ✅ Existe (`TokenRecuperacionModel`) |
| **Repositorio** | ✅ Existe (`tokenRecuperacionRepositorio`) — Token |
| **Repositorio** | ❌ **NO existe implementación** — `UsuarioRecuperacionRepositorio` (interfaz definida, sin adapter concreto) |
| **Servicio de recuperación** | ✅ Implementado (Solicitar, Validar, Confirmar) |
| **Rate limiting** | ❌ No integrado — config definida pero nunca ejecutada |
| **Invalidación de sesiones** | ✅ Implementada en `ConfirmarRestablecimiento` |
| **Build** | ✅ Compila |
| **Tests de dominio** | ✅ Pasan (5 tests) |
| **Tests de servicio** | ❌ No existen (`servicio_recuperacion_test.go` no creado) |

## 2. Estado por Componente

### 2.1 Modelo de Datos

| Elemento | Estado |
|----------|--------|
| Tabla `tokens_recuperacion` | Modelo GORM existe con columnas: id, usuario_id, token_hash, expira_en, usado, creado_en, usado_en |
| Índices | No especificados en migraciones (no hay `AutoMigrate` para este modelo). La spec pide `idx_tokens_recuperacion_hash` e `idx_tokens_recuperacion_usuario`. |
| FK a `usuarios(id)` | No declarada explícitamente en el modelo GORM (no tiene `gorm:"constraint:..."`) |
| `ON DELETE CASCADE` | No implementado |

### 2.2 Capa de Dominio

| Elemento | Estado | Archivo |
|----------|--------|---------|
| Entidad `TokenRecuperacion` | ✅ Existe con todos los campos | `internal/recuperacion/domain/token_recuperacion.go` |
| Value Objects | ❌ No implementados (los campos son strings planos: id, usuarioID, tokenHash) | — |
| `NuevoTokenRecuperacion` (factory) | ✅ Crea token con hash SHA-256 | `token_recuperacion.go:21` |
| `NuevoTokenRecuperacionDesdeBD` (factory) | ✅ Reconstruye desde BD | `token_recuperacion.go:33` |
| `EsValido()` | ✅ Valida no usado y no expirado | `token_recuperacion.go:46` |
| `Usar()` | ✅ Marca como usado con timestamp | `token_recuperacion.go:57` |
| `HashearToken()` | ✅ SHA-256 | `token_recuperacion.go:63` |
| Getters | ✅ Todos los campos tienen getter | `token_recuperacion.go:69-75` |
| Interfaz `TokenRecuperacionRepositorio` | ✅ Crear, ObtenerPorHash, Actualizar | `repositorio.go:6-10` |
| Interfaz `UsuarioRecuperacionRepositorio` | ✅ ObtenerPorCorreo, ActualizarPassword | `repositorio.go:20-22` |
| Errores de dominio | ✅ Definidos (ErrEnlaceInvalido, ErrEnlaceExpirado, ErrEnlaceYaUtilizado, ErrDemasiadasSolicitudes, ErrPasswordDebil, ErrEmailRequerido, ErrEmailInvalido) | `errores.go` |
| Tests de dominio | ✅ 5 tests (creación, válido, expirado, usado, hash determinista, hash distinto, desde BD) | `token_recuperacion_test.go` |

### 2.3 Capa de Aplicación

| Elemento | Estado | Archivo |
|----------|--------|---------|
| `ServicioRecuperacion` | ✅ Implementado | `servicio_recuperacion.go` |
| `SolicitarRecuperacion` | ✅ Implementado (sin rate limiting) | `servicio_recuperacion.go:72` |
| `ValidarToken` (GET) | ✅ Implementado | `servicio_recuperacion.go:127` |
| `ConfirmarRestablecimiento` (POST) | ✅ Implementado (sin transacción, sin reset de intentos) | `servicio_recuperacion.go:149` |
| `ConfigRecuperacion` | ✅ Definido con defaults | `servicio_recuperacion.go:17-22` |
| Comandos DTO | ✅ ComandoSolicitarRecuperacion, ComandoValidarToken, ComandoConfirmarRestablecimiento | `comando.go` |
| Respuestas DTO | ✅ RespuestaSolicitar, RespuestaValidar, RespuestaConfirmar | `respuesta.go` |
| Tests del servicio | ❌ **No existen** — `servicio_recuperacion_test.go` no está creado | — |

**Brechas detectadas en el servicio:**

| Aspecto | Estado | Detalle |
|---------|--------|---------|
| Rate limiting por IP | ❌ No implementado | `ComandoSolicitarRecuperacion.IPOrigen` existe pero nunca se usa |
| Rate limiting por usuario | ❌ No implementado | Config definida pero nunca ejecutada |
| Transacción atómica en confirmación | ❌ No implementado | Las 3 operaciones (update password, mark token, invalidar sesiones) son secuenciales sin rollback |
| Reset de intentos fallidos | ❌ No implementado | La spec indica RN-REC-05: resetear `intentos_fallidos = 0` al confirmar, pero no se llama a `credRepo` ni `ResetearIntentos()` |
| Validación de password (fortaleza real) | ⚠️ Parcial | Solo valida `len < 8`, no usa las mismas políticas de fortaleza que registro |

### 2.4 Capa de Infraestructura

| Elemento | Estado | Archivo |
|----------|--------|---------|
| Modelo GORM `TokenRecuperacionModel` | ✅ Existe con columnas correctas | `token_recuperacion_model.go` |
| `TableName()` | ✅ `tokens_recuperacion` | `token_recuperacion_model.go:20` |
| `ToDomain()` / `FromDomain()` | ✅ Conversiones implementadas | `token_recuperacion_model.go:24-46` |
| Repositorio PostgreSQL | ✅ `Crear`, `ObtenerPorHash`, `Actualizar` | `token_recuperacion_repositorio.go` |
| `AutoMigrate` para `tokens_recuperacion` | ❌ **No incluido** en `config/database.go:RunMigrations` | `database.go:32-62` |
| Implementación de `UsuarioRecuperacionRepositorio` | ❌ **No existe** — no hay adapter que implemente `ObtenerPorCorreo` ni `ActualizarPassword` | — |

### 2.5 Seguridad

| Aspecto | Estado | Evidencia |
|---------|--------|-----------|
| Token hasheado (SHA-256) | ✅ Sí | `HashearToken()` in `token_recuperacion.go:63` |
| Token de un solo uso | ✅ Sí | `Usar()` setea `usado = true`, `EsValido()` rechaza si usado |
| Expiración configurable | ✅ Sí | `ConfigRecuperacion.TokenExpiracion` (default 1h) |
| No revelar existencia de email | ✅ Sí | `SolicitarRecuperacion` retorna mensaje genérico siempre |
| Rate limiting por IP | ❌ No implementado | Config existe con default 3/hora, pero el servicio nunca lo ejecuta |
| Rate limiting por usuario | ❌ No implementado | Config existe con default 1/15min, pero el servicio nunca lo ejecuta |
| Invalidación de sesiones | ✅ Sí | `sesionRepo.InvalidarTodasPorUsuarioID` llamado en confirm |
| Reset de intentos fallidos | ❌ No implementado | No se llama a `credRepo` ni `ResetearIntentos()` |
| Validación de fortaleza de password | ⚠️ Parcial | Solo verifica `len >= 8`, no usa `EspecificacionCredenciales` |

### 2.6 Integración

| Elemento | Estado | Detalle |
|----------|--------|---------|
| Template email RECUPERACION_CONTRASENA | ✅ Existe | `templates.go:25-40` con variables `{{nombre}}`, `{{token}}`, `{{expiracion_horas}}` |
| Constante TipoRecuperacionContrasena | ✅ Existe | `tipos_email.go:8` |
| Servicio de recuperación usa EmailServicio | ✅ Sí | `servicio_recuperacion.go:111-121` llama `EnviarTemplate` |
| Registro en Registry | ❌ No registrado | `registry.go` no incluye `ServicioRecuperacion` |
| Variables de entorno definidas | ❌ No definidas | `env.go` no tiene `RECUPERACION_TOKEN_EXPIRACION`, `RECUPERACION_RATE_LIMIT_IP`, etc. |
| HTTP Handlers | ❌ No implementados | No hay handlers para `/auth/recover`, `/auth/reset`, `/auth/reset/confirm` |
| Facade | ❌ No implementada | `AuthFacade` solo tiene `Registrar` y `Login` |
| Router | ❌ No registrado | `router.go:30-33` solo registra health, register, login |
| AutoMigrate | ❌ No incluido | `database.go:RunMigrations` no migra `TokenRecuperacionModel` |

## 3. Checklist de Validación

| # | Ítem | Estado | Evidencia |
|---|------|--------|-----------|
| 1 | ¿El token de recuperación se almacena como hash SHA-256 en tabla dedicada? | ✅ Sí | `HashearToken()` con SHA-256; tabla `tokens_recuperacion` con columna `token_hash` |
| 2 | ¿El token en plano nunca se expone en respuestas HTTP? | ✅ Sí | Solo va en el email (body del template) y en request body/query |
| 3 | ¿Los tokens son de un solo uso (usado = true)? | ✅ Sí | `Usar()` setea `usado = true`; `EsValido()` rechaza si `usado == true` |
| 4 | ¿La expiración del token es configurable? | ✅ Sí | `ConfigRecuperacion.TokenExpiracion` con default 1h |
| 5 | ¿La solicitud no revela si el email existe o no? | ✅ Sí | Mensaje genérico "Si el email existe..." siempre |
| 6 | ¿Hay rate limiting por IP y por usuario? | ❌ No | `IPOrigen` y `ConfigRecuperacion` definidos pero nunca usados en el servicio |
| 7 | ¿Al confirmar el restablecimiento se invalidan todas las sesiones activas? | ✅ Sí | `sesionRepo.InvalidarTodasPorUsuarioID` llamado |
| 8 | ¿La nueva contraseña pasa las mismas validaciones que en registro? | ⚠️ Parcial | Solo verifica `len >= 8`; no usa `EspecificacionCredenciales` |
| 9 | ¿Los tokens expirados o usados retornan error informativo pero seguro? | ✅ Sí | `ErrEnlaceExpirado`, `ErrEnlaceYaUtilizado`, `ErrEnlaceInvalido` |
| 10 | ¿El endpoint de confirmar cambia la contraseña en transacción atómica? | ❌ No | Operaciones secuenciales sin rollback ante fallos |
| 11 | ¿Los intentos fallidos se resetean al restablecer la contraseña? | ❌ No | No se invoca `ResetearIntentos()` ni `credRepo.Actualizar()` |

## 4. Brechas Detectadas

Las brechas están ordenadas por prioridad de implementación:

### Prioridad ALTA (bloqueante para funcionamiento)

| # | Brecha | Descripción | Componente |
|---|--------|-------------|------------|
| B1 | **`UsuarioRecuperacionRepositorio` sin implementación** | La interfaz está definida en `domain/repositorio.go` pero no hay ningún adapter concreto que la implemente. El servicio nunca puede ser instanciado. Se necesita implementar `ObtenerPorCorreo` (buscar usuario por email) y `ActualizarPassword` (actualizar hash en credenciales). | Infraestructura |
| B2 | **Sin HTTP handlers ni rutas** | No existen endpoints `/auth/recover`, `/auth/reset`, ni `/auth/reset/confirm`. No hay handlers, facade, ni registro en router. | Presentación |
| B3 | **Sin registro en Registry** | `ServicioRecuperacion` no está en `registry.go`, no se construye en `NewRegistry()`, no se pasa a facade ni a router. | DI / Registry |
| B4 | **Sin AutoMigrate para `TokenRecuperacionModel`** | La tabla `tokens_recuperacion` no se migra automáticamente en `config/database.go:RunMigrations`. | Config / BD |

### Prioridad MEDIA (incumplimiento de spec)

| # | Brecha | Descripción | Regla de Negocio |
|---|--------|-------------|------------------|
| B5 | **Rate limiting no implementado** | El servicio tiene `IPOrigen` en el comando y `ConfigRecuperacion` con defaults, pero `SolicitarRecuperacion` nunca verifica límites. Debe validar por IP (3/hora) y por usuario (1/15min). | RN-REC-08, RN-REC-09 |
| B6 | **Sin transacción atómica en confirmación** | `ConfirmarRestablecimiento` ejecuta 3 operaciones (actualizar password, marcar token usado, invalidar sesiones) sin transacción. Si falla la invalidación de sesiones, el cambio de password ya persiste. | Sección 6.3 paso 4 |
| B7 | **Sin reset de intentos fallidos** | Al confirmar restablecimiento, la spec exige resetear `intentos_fallidos = 0`, pero el servicio no lo hace. | RN-REC-05 |
| B8 | **Validación de password insuficiente** | Solo verifica `len >= 8`. Debe usar las mismas políticas de fortaleza que el registro (vía `CredencialesRepositorio` o `EspecificacionCredenciales`). | RN-REC-07 |

### Prioridad BAJA (mejora continua)

| # | Brecha | Descripción |
|---|--------|-------------|
| B9 | **Variables de entorno no definidas** | `RECUPERACION_TOKEN_EXPIRACION`, `RECUPERACION_RATE_LIMIT_IP`, `RECUPERACION_RATE_LIMIT_USUARIO`, `RECUPERACION_RATE_LIMIT_VENTANA` no están en `env.go`. El servicio usa defaults, pero no son configurables vía env. |
| B10 | **Sin tests de servicio** | `servicio_recuperacion_test.go` no existe. Los 20 casos TDD de la spec (sección 9) no están cubiertos. |
| B11 | **Sin tests de repositorio** | No hay tests para `tokenRecuperacionRepositorio` (crear, obtener por hash, actualizar). |
| B12 | **Email asíncrono sin manejo de errores real** | El envío de email es `go func()` con `fmt.Printf` de error. No hay mecanismo de reintento ni cola. |
| B13 | **FK y constraint no declarados** | `usuario_id` no tiene constraint explícito (`REFERENCES usuarios(id) ON DELETE CASCADE`). Los índices `idx_tokens_recuperacion_hash` e `idx_tokens_recuperacion_usuario` no se crean. |
| B14 | **Sin Value Objects** | `id`, `usuarioID`, `tokenHash` son strings planos. No hay VO que encapsule validación de formato. |

## 5. Recomendaciones

### Inmediatas (orden de implementación)

1. **Implementar `UsuarioRecuperacionRepositorio`** — Crear un adapter que implemente `ObtenerPorCorreo` (consulta a `usuarios` por email, retorna proyección ligera) y `ActualizarPassword` (actualiza `password_hash` en `credenciales`). Puede reutilizar `CredencialesRepositorio.Actualizar()` y el filtro de `Listar()` del repositorio de usuarios, o implementar consultas directas.

2. **Registrar en Registry** — Añadir `ServicioRecuperacion` al `Registry` struct, construirlo en `NewRegistry()` inyectando `TokenRecuperacionRepositorio`, `UsuarioRecuperacionRepositorio`, `SesionRepositorio`, `CredencialesRepositorio`, `EncriptacionServicio`, `EmailServicio`, `GeneradorID`, y `ConfigRecuperacion`.

3. **Añadir AutoMigrate** — Incluir `TokenRecuperacionModel` en `RunMigrations()` en `database.go`.

4. **Crear HTTP handlers, facade y rutas** — Implementar `RecoverHandler`, `ResetHandler`, `ResetConfirmHandler` con Huma v2, facade, y registrar en router.

### De spec (para cumplir requisitos)

5. **Integrar rate limiting** — Antes de procesar `SolicitarRecuperacion`, validar rate limit por IP (usando `IntentoIPRepositorio` existente o el `ServicioRateLimit` genérico) y por usuario (nueva lógica con contador por `usuarioID`).

6. **Envolver confirmación en transacción** — Usar `gorm.DB` transaction o `UnitOfWork` para que las 3 operaciones sean atómicas.

7. **Resetear intentos fallidos** — En `ConfirmarRestablecimiento`, obtener credenciales del usuario vía `credRepo.ObtenerPorUsuarioID()`, llamar `ResetearIntentos()`, y persistir con `credRepo.Actualizar()`.

8. **Fortalecer validación de password** — Reutilizar la `EspecificacionCredenciales` del módulo de seguridad o centralizar la validación.

### De calidad

9. **Implementar tests TDD** — Cubrir los 20 escenarios de la sección 9 de la spec.

10. **Definir variables de entorno** — Añadir a `env.go` las 4 variables de recuperación, y pasarlas a `ConfigRecuperacion` desde `config.Config`.

### Archivos afectados por las recomendaciones

```
internal/recuperacion/
├── application/services/recuperacion/
│   ├── servicio_recuperacion.go              ← MEJORAR: rate limiting, transacción, reset intentos
│   └── servicio_recuperacion_test.go         ← CREAR: tests TDD
├── infrastructure/persistence/postgres/
│   ├── token_recuperacion_repositorio.go     ← Existe (sin cambios)
│   └── usuario_recuperacion_repositorio.go   ← CREAR: implementa UsuarioRecuperacionRepositorio

internal/presentation/
├── handlers/
│   ├── recover_handler.go                    ← CREAR
│   ├── reset_handler.go                      ← CREAR
│   └── reset_confirm_handler.go              ← CREAR
├── facades/
│   ├── auth_facade.go                        ← MEJORAR: añadir métodos de recuperación
│   └── auth_facade_impl.go                   ← MEJORAR: implementar nuevos métodos
├── router/
│   └── router.go                             ← MEJORAR: registrar nuevos handlers
└── dto/
    ├── recover_dto.go                        ← CREAR
    └── reset_dto.go                          ← CREAR

internal/registry/
└── registry.go                               ← MEJORAR: añadir ServicioRecuperacion

internal/config/
├── env.go                                    ← MEJORAR: añadir vars de recuperación
└── database.go                               ← MEJORAR: añadir AutoMigrate

cmd/
└── main.go                                   ← MEJORAR: pasar EmailServicio al registry
```
