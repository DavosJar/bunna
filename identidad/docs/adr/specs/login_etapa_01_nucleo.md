# Etapa 1: Núcleo de Autenticación (Login + JWT)

## Ubicación

Nueva bounded context `internal/autenticacion/` con capas `domain/`, `application/services/login/`.
Sigue el mismo patrón que `usuarios` y `seguridad`.

## Dependencias

- `autenticacion → seguridad` (repositorio de credenciales + servicio de encriptación)
- `autenticacion → usuarios` (repositorio de usuarios)
- `autenticacion → shared` (generador de IDs)

## Alcance

### Incluye
- Entidad `Sesion` con máquina de estados
- Value Object `TokenPair`
- Interfaces: `TokenServicio` (generar/validar access token), `SesionRepositorio` (CRUD + búsquedas)
- Servicio de aplicación `ServicioLogin` con su Comando y Respuesta
- Validaciones de comando
- Tests completos (ver TDD)
- Mocks para todas las dependencias

### NO incluye (vienen después)
- Implementación concreta de JWT (infrastructure)
- PostgreSQL (infrastructure)
- Refresh token rotation (etapa 2)
- Logout (etapa 2)
- Rate limiting (etapa 3)
- Auditoría (etapa 4)
- Handlers HTTP / presentación
- Integración con contexto usuarios para bloqueo permanente (etapa 3)

## Modelo de Dominio

### Decisiones de Diseño

**Sesion es Aggregate Root**
- Tiene su propio repositorio (SesionRepositorio)
- Máquina de estados propia (ACTIVA → REVOCADA | EXPIRADA)
- Lifecycle independiente: se crea en login, se revoca en logout, se expira por tiempo
- El repositorio opera directamente sobre Sesion, no a través de otro aggregate
- Conclusión: Sesion es un Aggregate Root en el bounded context autenticacion

### Sesión — Entidad

Estados: `ACTIVA → REVOCADA | EXPIRADA`. REVOCADA y EXPIRADA son terminales.

Condiciones de transición:
- `ACTIVA → REVOCADA`: solo por logout explícito del usuario
- `ACTIVA → EXPIRADA`: cuando `time.Now() > expiraRefresh`
- Transición desde terminal: prohibida, error de dominio

Comportamiento esperado:
- Creación: requiere ID, usuarioID, refresh token hash, dispositivo, IP, fechas de expiración
- Revocar: llama a `sesion.Revocar()` que marca estado + timestamp, luego `SesionRepositorio.Actualizar(sesion)`
- Consultar si expiró: compara `expiraRefresh` contra `time.Now()`
- Consultar si activa: estado == ACTIVA && no expiró

### TokenPair — Value Object

Compuesto por: access token (string), refresh token (string), expiración de cada uno (duración).
Inmutable, sin comportamiento.

### TokenServicio — Interfaz

Operaciones:
1. **GenerarAccessToken**: recibe usuarioID + sesionID → retorna `{ tokenString string, jti string, expira time.Time }`
2. **ValidarAccessToken**: recibe token string → retorna claims `{ usuarioID string, sesionID string, jti string, iat time.Time, exp time.Time }`
   - Si token expirado → error específico `ErrTokenExpirado`
   - Si firma inválida → error específico `ErrTokenFirmaInvalida`
   - Si token bien formado → claims

El campo `jti` es el JWT ID: un identificador único por token, generado por la implementación concreta. Se incluye en el retorno de `GenerarAccessToken` para que el servicio de aplicación pueda referenciarlo (rotación en etapa 2, auditoría en etapa 4).

### SesionRepositorio — Interfaz

Operaciones: `Crear`, `Actualizar`, `ObtenerPorID`, `ObtenerPorRefreshTokenHash`, `ObtenerActivasPorUsuarioID`.

> **Nota de diseño**: no existe método `Revocar` separado. La revocación es responsabilidad de la entidad (`sesion.Revocar()`); el repositorio persiste el estado resultante mediante `Actualizar`. Un método `Revocar` en el repositorio duplicaría lógica de dominio y acoplaría la infra al estado interno de la entidad.

### Errores de dominio

```
ErrCredencialesInvalidas   — correo no existe o contraseña incorrecta (mensaje idéntico en ambos casos)
ErrCuentaBloqueada         — usuario en estado BLOQUEADO permanente
ErrCuentaInactiva          — usuario en estado INACTIVO
ErrCuentaNoDisponible      — usuario en estado PENDIENTE_DE_ELIMINACION
ErrCorreoNoVerificado      — usuario en estado NO_VERIFICADO
ErrBloqueadoTemporalmente  — credenciales con bloqueadoHasta > now
ErrTransicionInvalida      — intento de cambiar estado desde un estado terminal
ErrTokenExpirado           — access token expirado
ErrTokenFirmaInvalida      — firma JWT inválida
```

Todos son sentinel errors (`var Err... = errors.New(...)`). El servicio retorna siempre el error del dominio; nunca expone el error de infraestructura al llamador.

## Flujo del ServicioLogin

1. **Validar comando** (antes de tocar cualquier repositorio)
   - Correo vacío → error
   - Formato correo inválido → error
   - Password vacío → error

2. **Buscar usuario por correo**
   - No existe → `ErrCredencialesInvalidas`

3. **Validar estado del usuario**
   - BLOQUEADO → `ErrCuentaBloqueada`
   - INACTIVO → `ErrCuentaInactiva`
   - PENDIENTE_DE_ELIMINACION → `ErrCuentaNoDisponible`
   - NO_VERIFICADO → `ErrCorreoNoVerificado`

4. **Obtener credenciales del usuario**
   - No existen → `ErrCredencialesInvalidas`
   - *(En este paso se leen `intentosFallidos` previos — necesario para throttling en etapa 3)*

5. **Verificar si la cuenta está bloqueada temporalmente**
   - `bloqueadoHasta > now` → `ErrBloqueadoTemporalmente`

6. **Verificar contraseña** (contra hash con EncriptacionServicio)
   - Incorrecta → marcar intento fallido; si llega a 5 → bloquear 15 min; retornar `ErrCredencialesInvalidas`
   - Correcta → resetear `intentosFallidos = 0` y `bloqueadoHasta = zero`

7. **Generar tokens**: access token (JWT vía `TokenServicio`) + refresh token (UUID)

8. **Crear sesión** con estado ACTIVA y persistir vía `SesionRepositorio.Crear`

9. **Retornar** TokenPair + datos básicos del usuario (ID, correo, nombre)

### Reglas de seguridad críticas
- El mensaje de error debe ser EXACTAMENTE EL MISMO para "correo no existe" que para "contraseña incorrecta". Nunca revelar cuál falló.
- El bloqueo por intentos es TEMPORAL (15 minutos), no permanente (el permanente llega en etapa 3).
- Después de login exitoso, los intentos fallidos se resetean a 0.
- Las passwords mayores a 72 bytes se aceptan sin error; bcrypt trunca silenciosamente y eso es conocido y aceptado. El sistema no valida un máximo de longitud para no revelar que bcrypt tiene ese límite.

## Especificación TDD

### Happy Path (4 casos)

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 1 | Login exitoso básico | Usuario ACTIVO, credenciales válidas, verificadas, sin intentos | Login correcto | TokenPair no vacío, sesión ACTIVA creada (`SesionRepositorio.Crear` llamado una vez), datos de usuario correctos |
| 2 | Login con metadatos | Usuario ACTIVO | Login con dispositivo="Mozilla/5.0", IP="192.168.1.1" | Sesión persistida guarda dispositivo="Mozilla/5.0" e IP="192.168.1.1" |
| 3 | Login resetea intentos | Usuario con 3 intentos fallidos previos | Login correcto | Credenciales actualizadas con `intentosFallidos=0`, `bloqueadoHasta=zero` |
| 4 | Múltiples dispositivos | Usuario ACTIVO | Login desde iOS, luego login desde Android | `SesionRepositorio.Crear` llamado dos veces; `ObtenerActivasPorUsuarioID` retorna 2 sesiones para ese usuario |

### Sad Path — Validaciones (4 casos)

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 5 | Correo vacío | — | Comando con Correo="" | Error de validación; ninguna dependencia llamada |
| 6 | Formato inválido | — | Correo="notanemail" | Error de validación |
| 7 | Sin arroba | — | Correo="user.com" | Error de validación |
| 8 | Password vacío | — | Password="" | Error de validación; ninguna dependencia llamada |

### Sad Path — Usuario (6 casos)

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 9 | Correo no existe | No hay usuario con ese correo | Login cualquiera | `ErrCredencialesInvalidas` (mensaje idéntico al caso 10) |
| 10 | Password incorrecta | Usuario existe con pass "X" | Login con pass "Y" | `ErrCredencialesInvalidas` (mensaje idéntico al caso 9), `intentosFallidos` incrementado en 1 |
| 11 | Usuario BLOQUEADO | Estado=BLOQUEADO | Credenciales válidas | `ErrCuentaBloqueada` |
| 12 | Usuario INACTIVO | Estado=INACTIVO | Credenciales válidas | `ErrCuentaInactiva` |
| 13 | PENDIENTE_ELIMINACION | Estado=PENDIENTE_DE_ELIMINACION | Credenciales válidas | `ErrCuentaNoDisponible` |
| 14 | NO_VERIFICADO | Estado=NO_VERIFICADO, verificación=PENDIENTE | Credenciales válidas | `ErrCorreoNoVerificado` |

### Sad Path — Bloqueo temporal (3 casos)

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 15 | Bloqueado aún activo | 5 intentos fallidos, bloqueadoHasta > now | Password correcta | `ErrBloqueadoTemporalmente`; credenciales no se modifican |
| 16 | Bloqueo expiró | 5 intentos fallidos, bloqueadoHasta < now | Password correcta | Login exitoso; `intentosFallidos=0`, `bloqueadoHasta=zero` |
| 17 | 5to intento activa bloqueo | 4 intentos fallidos, sin bloqueo | Password incorrecta | `ErrCredencialesInvalidas`; `intentosFallidos=5`, `bloqueadoHasta=now+15min` |

### Edge Cases (4 casos)

| # | Caso | Descripción |
|---|------|-------------|
| 18 | Concurrencia en intentos | 2 login fallidos simultáneos con 4 intentos previos. El dominio es puro; la atomicidad la garantiza infra (optimistic lock). El test verifica que el dominio aplica la lógica correctamente sobre el estado que recibe, sin asumir aislamiento de BD. |
| 19 | Unicode en correo | Correo con acentos o emojis válido según `mail.ParseAddress` → login procede normalmente. |
| 20 | Password >72 bytes | bcrypt trunca silenciosamente; el servicio no retorna error. El test envía una password de 100 bytes y verifica que el login procede si el hash fue generado con los mismos 100 bytes. |
| 21 | Claims del access token | El token retornado por el mock de `TokenServicio` debe contener `sub=usuarioID`, `jti` no vacío, `iat` y `exp`. El test verifica que `ServicioLogin` pasa `usuarioID` y `sesionID` correctos al llamar `GenerarAccessToken`. |

## Criterios de Aceptación

1. Todos los test cases pasan
2. Error idéntico (mismo tipo y mensaje) para "correo no existe" y "contraseña incorrecta"
3. Bloqueo temporal en el 5to intento fallido consecutivo
4. Bloqueo temporal = 15 minutos exactos
5. Intentos fallidos se resetean tras login exitoso
6. TokenPair con access y refresh no vacíos
7. Sesión creada con estado ACTIVA + metadatos dispositivo/IP
8. Validaciones ocurren antes de cualquier llamada a repositorio
9. `SesionRepositorio` no expone método `Revocar`; la revocación pasa siempre por `sesion.Revocar()` + `Actualizar`
10. Cobertura > 90% en servicio de aplicación

## Tareas

### Dominio
- [ ] Crear package `internal/autenticacion/domain/`
- [ ] Definir sentinel errors de dominio (`ErrCredencialesInvalidas`, `ErrBloqueadoTemporalmente`, etc.)
- [ ] Implementar entidad `Sesion` con máquina de estados
- [ ] Implementar value object `TokenPair`
- [ ] Definir interfaz `TokenServicio` (con retorno de `jti`)
- [ ] Definir interfaz `SesionRepositorio` (sin método `Revocar`)
- [ ] Tests de dominio: transiciones, creación, revocación, expiración

### Aplicación
- [ ] Crear `internal/autenticacion/application/services/login/`
- [ ] Implementar `ComandoLogin` y `DtoRespuestaLogin`
- [ ] Implementar `ServicioLogin` con el flujo completo de 9 pasos
- [ ] Implementar validaciones de comando
- [ ] Implementar mocks para tests
- [ ] Tests: 21 casos (happy + sad + edge)

### Integración
- [ ] Actualizar `registry/registry.go` con dependencias de autenticación
- [ ] `go build ./...` compila sin errores
- [ ] `go test ./internal/autenticacion/...` todos pasan
