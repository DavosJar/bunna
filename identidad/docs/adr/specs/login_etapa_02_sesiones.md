# Etapa 2: Gestión de Sesiones (Refresh Tokens + Logout)

## Prerrequisito

Esta etapa asume completa la Etapa 1 (núcleo de login con JWT + entidad Sesión + interfaces TokenServicio y SesionRepositorio).

## Alcance

### Incluye
- Servicio `RefreshToken`: rotación de refresh tokens
- Servicio `Logout`: cierre de sesión individual y cierre de todas las sesiones
- Lógica de expiración y renovación en la entidad Sesión
- Actualización de `ServicioLogin` para que el refresh token se almacene como hash (no en plano)
- Tests completos para refresh y logout

### NO incluye (etapas posteriores)
- Implementación concreta de JWT (infrastructure)
- PostgreSQL (infrastructure)
- Rate limiting / IP blocking (etapa 3)
- Notificaciones de seguridad (etapa 4)
- Handlers HTTP

## Comportamiento esperado

### Refresh Token — Flujo

1. Cliente envía refresh token (string en plano)
2. Calcular SHA256 hash del refresh token recibido
3. Buscar sesión por ese hash en el repositorio (`ObtenerPorRefreshTokenHash`)
4. Si no se encuentra → `ErrRefreshTokenInvalido` (cubre: no existe, ya rotado, hash corrupto)
5. Validar estado de la sesión encontrada:
   - REVOCADA → `ErrSesionRevocada`
   - EXPIRADA (`expiraRefresh < now`) → `ErrSesionExpirada`
   - ACTIVA → continuar
6. Generar NUEVO access token (con nuevo `jti`) vía `TokenServicio.GenerarAccessToken`
7. Generar NUEVO refresh token (UUID)
8. Calcular hash del nuevo refresh token
9. **Rotación**: llamar `sesion.ActualizarRefreshToken(nuevoHash, nuevaExpiraRefresh)` y persistir con `SesionRepositorio.Actualizar`
   - El refresh token anterior queda inválido inmediatamente
   - `expiraAccess` de la sesión se actualiza con la nueva expiración del access token
10. Retornar nuevo TokenPair

> **Refresh voluntario**: si la sesión está ACTIVA y el access token aún no expiró, el refresh se permite igualmente. El cliente decide cuándo refrescar; el servidor no penaliza el refresh anticipado.

### Rotación de Refresh Tokens — Reglas

- Cada refresh genera un refresh token NUEVO
- El refresh token anterior queda inválido (el hash ya no existe en la sesión)
- **Protección anti-replay**: un token ya rotado no se encuentra por hash → `ErrRefreshTokenInvalido`. Sin side effects adicionales. El servicio no distingue entre "token ya rotado", "token nunca existió" o "token de otro usuario" — todos resultan en no encontrado.
- Si un atacante usa un token viejo antes que el cliente legítimo lo use, el cliente legítimo recibirá `ErrRefreshTokenInvalido` en su próximo intento. Este escenario es detectado en etapa 4 (auditoría) mediante el evento `RefreshReplay`; en esta etapa, el comportamiento es simplemente retornar error sin side effects.

### Logout — Flujo individual

**Entrada**: `ComandoLogout` acepta exactamente UNO de los dos campos: `refreshToken` (string) o `sesionID` (string). Si ambos están presentes, se ignora `refreshToken` y se usa `sesionID`. Si ninguno está presente → error de validación.

1. Si se recibe `refreshToken`: calcular hash → `ObtenerPorRefreshTokenHash`
   Si se recibe `sesionID`: `SesionRepositorio.ObtenerPorID`
2. Si no se encuentra → `ErrRefreshTokenInvalido` o `ErrSesionNoEncontrada` según el campo usado
3. Validar estado: si REVOCADA → `ErrSesionRevocada`; si EXPIRADA → `ErrSesionExpirada`
4. Llamar `sesion.Revocar()` + `SesionRepositorio.Actualizar(sesion)`
5. El access token emitido sigue siendo válido hasta su expiración natural (los JWTs no pueden revocarse)
6. Retornar confirmación

### Logout — Cierre de todas las sesiones

**Entrada**: `ComandoLogoutAll` con `usuarioID` (requerido).

1. `SesionRepositorio.ObtenerActivasPorUsuarioID(usuarioID)`
2. Por cada sesión: `sesion.Revocar()` + `SesionRepositorio.Actualizar(sesion)`
3. Si no hay sesiones activas → retornar éxito (operación idempotente, no es error)
4. Incluye la sesión del dispositivo actual si está entre las activas (no se excluye la sesión propia)

### Expiración de Sesiones

La sesión tiene dos fechas de expiración:
- `expiraAccess`: cuándo expira el access token actual (se actualiza en cada refresh)
- `expiraRefresh`: cuándo expira la sesión completa (se fija en el login; no se extiende con cada refresh)

Política:
- Access token expirado → el cliente DEBE usar refresh token para obtener uno nuevo
- Refresh token expirado (`expiraRefresh < now`) → la sesión se marca EXPIRADA, el cliente DEBE hacer login de nuevo
- Una sesión EXPIRADA no se puede refrescar ni reactivar

### Refresh Token como hash

- El refresh token se almacena SIEMPRE como hash SHA256 (no bcrypt)
  - Razón: bcrypt es lento (500ms+) y penaliza cada refresh
  - Los refresh tokens son UUIDs de alta entropía, no necesitan stretching
  - SHA256 es suficiente para tokens de alta entropía
- Nunca en plano en persistencia ni en logs
- La comparación se hace contra el hash almacenado
- El hash se calcula en el servicio de aplicación; la entidad `Sesion` recibe el hash ya calculado

### Errores de dominio (nuevos en esta etapa)

```
ErrRefreshTokenInvalido  — sesión no encontrada por hash (token no existe, ya rotado, o corrupto)
ErrSesionRevocada        — sesión en estado REVOCADA
ErrSesionExpirada        — sesión en estado EXPIRADA o expiraRefresh < now
ErrSesionNoEncontrada    — búsqueda por ID sin resultado (usado en logout por sesionID)
```

## Especificación TDD

### Refresh Token — Happy Path

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 1 | Refresh exitoso básico | Sesión ACTIVA, refresh token válido, access expirado | Refresh con refresh token | Nuevo TokenPair no vacío; `SesionRepositorio.Actualizar` llamado con el nuevo hash; viejo hash ya no produce resultado en `ObtenerPorRefreshTokenHash` |
| 2 | Refresh antes de que expire access | Sesión ACTIVA, `expiraAccess > now` | Refresh | También se permite; nuevo TokenPair retornado; `expiraAccess` de la sesión actualizado |
| 3 | Múltiples refrescos consecutivos | Sesión ACTIVA | Refresh → Refresh → Refresh | Cada rotación exitosa; solo el último hash persiste en la sesión |

### Refresh Token — Sad Path

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 4 | Refresh token no existe | No hay sesión con ese hash | Refresh | `ErrRefreshTokenInvalido` |
| 5 | Sesión REVOCADA | Sesión en estado REVOCADA | Refresh | `ErrSesionRevocada` |
| 6 | Sesión EXPIRADA | `expiraRefresh < now` | Refresh | `ErrSesionExpirada` |
| 7 | Refresh token ya rotado (anti-replay) | Sesión actualizada con nuevo hash | Refresh con el hash viejo | `ErrRefreshTokenInvalido` (búsqueda por hash no encuentra nada; comportamiento idéntico al caso 4) |
| 8 | Refresh token con hash corrupto | Hash recibido no coincide con ninguna sesión | Refresh | `ErrRefreshTokenInvalido` |
| 9 | Context timeout durante refresh | — | Refresh con contexto cancelado | Error de contexto propagado; `SesionRepositorio.Actualizar` no llamado |

### Logout — Happy Path

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 10 | Logout individual por refresh token | Sesión ACTIVA | Logout con `refreshToken` | `sesion.Revocar()` llamado; `SesionRepositorio.Actualizar` llamado; confirmación exitosa |
| 11 | Logout individual por sesionID | Sesión ACTIVA | Logout con `sesionID` | Mismo resultado que caso 10 |
| 12 | Logout de todas las sesiones | Usuario con 3 sesiones activas | LogoutAll | Las 3 sesiones en estado REVOCADA; `SesionRepositorio.Actualizar` llamado 3 veces |

### Logout — Sad Path

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 13 | Logout de sesión ya revocada | Sesión en estado REVOCADA | Logout | `ErrSesionRevocada` |
| 14 | Logout de sesión expirada | Sesión en estado EXPIRADA | Logout | `ErrSesionExpirada` |
| 15 | Logout con token inexistente | No hay sesión para ese hash | Logout con `refreshToken` | `ErrRefreshTokenInvalido` |
| 16 | LogoutAll con usuario sin sesiones activas | `ObtenerActivasPorUsuarioID` retorna vacío | LogoutAll | Éxito; `SesionRepositorio.Actualizar` no llamado |
| 17 | ComandoLogout sin ningún campo | refreshToken="" y sesionID="" | Logout | Error de validación; ningún repositorio llamado |

### Rotación — Anti-Replay (escenario de ataque)

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 18 | Replay después de rotación legítima | Cliente legítimo hizo refresh (hash rotado). Atacante tiene el token viejo | Atacante intenta refresh con token viejo | `ErrRefreshTokenInvalido`. Sesión NO se revoca (el legítimo sigue activo). Sin side effects. |
| 19 | Replay antes de uso legítimo | Atacante usa token válido antes que el cliente legítimo | El atacante obtiene nuevo TokenPair; el cliente legítimo intenta refresh después | Para el cliente legítimo: `ErrRefreshTokenInvalido`. Comportamiento idéntico al caso 18 desde la perspectiva del servicio. El servicio no puede ni necesita distinguir entre ambos escenarios. |

## Criterios de Aceptación

1. [ ] Refresh token genera rotación: el viejo hash queda inválido inmediatamente
2. [ ] Refresh token almacenado como hash SHA256, nunca en plano
3. [ ] `expiraAccess` de la sesión se actualiza en cada refresh
4. [ ] `expiraRefresh` de la sesión NO se extiende con cada refresh (se fija en el login)
5. [ ] Logout individual revoca UNA sesión específica
6. [ ] LogoutAll revoca TODAS las sesiones activas del usuario, incluyendo la propia
7. [ ] LogoutAll es idempotente: sin sesiones activas retorna éxito
8. [ ] Sesión EXPIRADA no se puede refrescar
9. [ ] Sesión REVOCADA no se puede refrescar
10. [ ] Anti-replay: token ya rotado devuelve `ErrRefreshTokenInvalido` sin side effects
11. [ ] `ComandoLogout` requiere al menos un campo; si ambos presentes, `sesionID` tiene precedencia
12. [ ] Todos los test cases pasan

## Tareas

### Aplicación — RefreshToken
- [ ] Crear `internal/autenticacion/application/services/refresh/`
- [ ] Implementar `ComandoRefreshToken` (`refreshToken string`)
- [ ] Implementar `DtoRespuestaRefreshToken`
- [ ] Implementar `ServicioRefreshToken` con flujo: hash → buscar → validar estado → rotar → actualizar → responder
- [ ] Implementar cálculo de hash SHA256 del refresh token recibido
- [ ] Tests: casos 1-9 + 18-19

### Aplicación — Logout
- [ ] Crear `internal/autenticacion/application/services/logout/`
- [ ] Implementar `ComandoLogout` (`refreshToken string`, `sesionID string`; al menos uno requerido; `sesionID` con precedencia)
- [ ] Implementar `ComandoLogoutAll` (`usuarioID string`)
- [ ] Implementar `DtoRespuestaLogout` (confirmación booleana)
- [ ] Implementar `ServicioLogout` con flujo individual
- [ ] Implementar `ServicioLogoutAll`
- [ ] Tests: casos 10-17

### Dominio — Sesion (extensión)
- [ ] Agregar método `ActualizarRefreshToken(nuevoHash string, nuevaExpiraAccess time.Time)` a la entidad `Sesion`
  - Solo válido desde estado ACTIVA; error de dominio si se llama desde terminal
  - Actualiza `refreshTokenHash` y `expiraAccess`; `expiraRefresh` no cambia

### Integración
- [ ] Actualizar `ServicioLogin` de etapa 1: guardar SHA256 hash del refresh token, no el token en plano
- [ ] `go build ./...` compila
- [ ] `go test ./internal/autenticacion/...` todos pasan
