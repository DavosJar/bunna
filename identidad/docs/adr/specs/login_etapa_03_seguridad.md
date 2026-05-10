# Etapa 3: Seguridad y Rate Limiting

## Prerrequisito

Etapas 1 y 2 completas (login + refresh tokens + logout funcionando).

## Alcance

### Incluye
- Rate limiting por dirección IP en login
- Rate limiting por correo/usuario
- Throttling progresivo (retraso artificial en login)
- Límite de sesiones concurrentes por usuario
- Bloqueo permanente por umbral de bloqueos temporales
- Política de expiración configurable para access/refresh tokens
- Interfaz `Clock` para tests deterministas
- Tests de seguridad

### NO incluye (etapa 4)
- Detección de anomalías (geolocalización, patrones de horario)
- Notificaciones al usuario (email de "nuevo login detectado")
- Dashboard de monitoreo
- Integración con sistemas externos de seguridad

## Decisiones de Diseño

### Bloqueo Permanente cruza Bounded Contexts

El bloqueo permanente requiere cambiar el estado del Usuario (en contexto `usuarios`) desde el contexto `autenticacion`. Para mantener la inversión de dependencias:

**Solución: BloqueoUsuarioPort**

El contexto `autenticacion` define una interfaz (puerto) en su capa de dominio:

```go
type BloqueoUsuarioPort interface {
    BloquearPermanentemente(ctx context.Context, usuarioID string) error
}
```

El contexto `usuarios` implementa este puerto. El contexto `autenticacion` depende de la interfaz, no de la implementación. La inyección se realiza en `registry/registry.go`.

Flujo:
1. `ServicioLogin` detecta que `bloqueosConsecutivos >= umbralBloqueosPermanente`
2. Llama a `bloqueoPort.BloquearPermanentemente(ctx, usuarioID)`
3. El puerto (implementado en `usuarios`) cambia el estado del Usuario a BLOQUEADO
4. Si falla, `ServicioLogin` retorna error

### Interfaz Clock para tests deterministas

El throttling, el rate limiting y el cálculo de bloqueos dependen de `time.Now()`. Para que los tests sean deterministas sin `time.Sleep` reales, se introduce una interfaz `Clock` inyectable:

```go
type Clock interface {
    Now() time.Time
    Sleep(d time.Duration)
}
```

- La implementación real usa `time.Now()` y `time.Sleep()`
- Los mocks en tests controlan el tiempo y registran si `Sleep` fue llamado (y con qué duración), sin bloquear
- `ServicioLogin`, `ServicioRefreshToken` y el `RateLimiter` reciben `Clock` como dependencia

### Throttling implementado dentro de ServicioLogin

El throttling vive DENTRO de `ServicioLogin`, no como decorador externo, porque necesita acceso a `intentosFallidos` previos de `CredencialesUsuario` que ya se leyeron en el paso 4 del flujo. Un decorador externo requeriría acceso adicional al repositorio de credenciales, creando acoplamiento innecesario.

### Convención de RateLimiter

La interfaz `RateLimiter` retorna `(bool, time.Duration)` (permitido, tiempoRestante). Los tests inyectan un mock que devuelve valores predefinidos sin tocar estado externo.

## Comportamiento esperado

### Rate Limiting por IP

- Máximo N intentos de login por minuto desde una misma IP
- Configurable: `MAX_INTENTOS_POR_IP`, `VENTANA_MINUTOS`
- Al superar: `ErrRateLimitExcedido` (HTTP 429 en la capa de presentación)
- La implementación concreta del store (Redis, memoria, etc.) es infraestructura
- A nivel de aplicación, se define la interfaz `RateLimiter` inyectable

Política por defecto:
- 10 intentos/minuto desde misma IP en login
- 20 intentos/minuto desde misma IP en refresh (menos restrictivo porque requiere token válido)
- Bloqueo de 5 minutos al superar el límite

### Rate Limiting por Correo/Usuario

- Máximo N intentos fallidos por minuto sobre un mismo correo (independientemente de la IP)
- Complementa el bloqueo por intentos de `CredencialesUsuario` (que es persistente en BD)
- Esta es una capa adicional en aplicación; no modifica entidades de dominio
- Previene ataques distribuidos donde múltiples IPs atacan un solo usuario

Política por defecto:
- 5 intentos/minuto por correo en login
- Bloqueo de 1 minuto al superar

### Throttling Progresivo (Retraso Artificial)

- El retraso se calcula sobre `intentosFallidos` AL INICIAR la petición (valor leído de BD en paso 4, antes de validar la contraseña)
- Fórmula: `retraso = min(base * 2^(intentosFallidosPrevios - 1), maximo)`
  - Con `base=200ms` y `maximo=3000ms`: 0 fallos → 0ms, 1 fallo → 200ms, 2 → 400ms, 3 → 800ms, 4 → 1600ms, 5+ → 3000ms
- El retraso se aplica SIEMPRE al final, independientemente de si el login es exitoso o fallido
- Esto evita ataques de timing: el atacante no puede distinguir éxito de fallo por la duración de la respuesta
- El orden de ejecución dentro del servicio:
  1. Leer `intentosFallidos` previos (paso 4 del flujo)
  2. Calcular `retraso` sobre ese valor
  3. Ejecutar el resto del flujo de login (validar contraseña, crear sesión, etc.)
  4. Aplicar `clock.Sleep(retraso)` **antes** de retornar la respuesta
  5. Retornar resultado (éxito o error)
- El sleep usa la interfaz `Clock` para que los tests registren el valor sin bloquearse

> **Seguridad**: el retraso se aplica sobre los intentos *previos*, no el estado final. Si el login es exitoso, `intentosFallidos` se resetea a 0 después del sleep — el SIGUIENTE login tendrá 0ms de retraso, no el actual.

### Límite de Sesiones Concurrentes

- Máximo M sesiones activas simultáneas por usuario
- Configurable: `MAX_SESIONES_POR_USUARIO`
- Se evalúa en `ServicioLogin`, antes de crear la nueva sesión (paso 8 del flujo)
- Al superar el límite:
  - **Opción A (defecto)**: revocar la sesión más antigua (por `creadaEn`), crear la nueva
  - **Opción B**: retornar `ErrLimiteSesionesAlcanzado`
- La política es configurable; el defecto es Opción A

### Bloqueo Permanente

- Cuando `CredencialesUsuario` acumula N bloqueos temporales consecutivos, el usuario pasa a estado BLOQUEADO permanente
- Contador: campo `bloqueosConsecutivos` en `CredencialesUsuario`
- Umbral configurable: `UMBRAL_BLOQUEOS_PERMANENTE` (defecto: 3)
- Comportamiento de `bloqueosConsecutivos`:
  - Se **incrementa** cada vez que `intentosFallidos` llega a 5 (bloqueo temporal activado)
  - Se **resetea a 0** cuando el usuario hace login exitoso
  - Si llega al umbral en el mismo evento que lo incrementa → se dispara el bloqueo permanente
- Al alcanzar el umbral: `ServicioLogin` llama a `BloqueoUsuarioPort.BloquearPermanentemente(usuarioID)`

### Política de Expiración Configurable

- Access token: TTL configurable (defecto: 15 minutos)
- Refresh token: TTL configurable (defecto: 7 días)
- La configuración vive en el servicio de aplicación, no en dominio
- `TokenServicio.GenerarAccessToken` recibe el TTL como parámetro; la sesión recibe `expiraRefresh = now + refreshTTL`

### Errores de dominio (nuevos en esta etapa)

```
ErrRateLimitExcedido         — IP o correo superó el límite; incluye tiempo restante
ErrLimiteSesionesAlcanzado   — solo cuando política=rechazar (Opción B)
```

## Especificación TDD

### Rate Limiting por IP

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 1 | IP dentro del límite | Mock `RateLimiter` retorna `(true, 0)` para clave IP | Login desde 1.1.1.1 | Login procede normalmente |
| 2 | IP excede el límite | Mock `RateLimiter` retorna `(false, 5min)` para clave IP | Login desde 1.1.1.1 | `ErrRateLimitExcedido` con `tiempoRestante=5min`; ningún repositorio de usuario o credenciales llamado |
| 3 | IPs diferentes no se afectan | Mock retorna `(false, ...)` para IP-A y `(true, 0)` para IP-B | Login desde IP-B | Login procede |
| 4 | Rate limit de IP aplica antes que todo lo demás | IP bloqueada | Login con correo vacío | `ErrRateLimitExcedido` (rate limit se verifica antes de validar el comando) |
| 5 | Rate limit de IP no afecta refresh | Mock bloqueado para clave "login:1.1.1.1" | `ServicioRefreshToken` desde misma IP | No afectado (el `RateLimiter` de login no se inyecta en refresh) |

### Rate Limiting por Correo

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 6 | Correo dentro del límite | Mock retorna `(true, 0)` para clave correo | Login con "a@b.com" | Login procede |
| 7 | Correo excede límite | Mock retorna `(false, 1min)` para clave correo | Login con "a@b.com" | `ErrRateLimitExcedido` con `tiempoRestante=1min` |
| 8 | Correo diferente no afectado | Mock retorna `(false,...)` para "a@b.com" y `(true,0)` para "c@d.com" | Login con "c@d.com" | Login procede |
| 9 | Rate limit de correo se verifica después de validar formato | Correo inválido | Login con "notanemail" | Error de validación de formato, no `ErrRateLimitExcedido` |

### Throttling Progresivo

Los tests de throttling inyectan un mock de `Clock` que registra las llamadas a `Sleep` sin bloquear.

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 10 | Sin intentos fallidos previos | `intentosFallidos=0` en BD | Login (exitoso o fallido) | `clock.Sleep` no llamado |
| 11 | 1 intento fallido previo | `intentosFallidos=1` | Login | `clock.Sleep(200ms)` llamado |
| 12 | 3 intentos fallidos previos | `intentosFallidos=3` | Login | `clock.Sleep(800ms)` llamado |
| 13 | 5+ intentos fallidos previos | `intentosFallidos=5` | Login | `clock.Sleep(3000ms)` llamado (tope máximo) |
| 14 | Login exitoso CON fallos previos | `intentosFallidos=2`, login correcto | Login exitoso | `clock.Sleep(400ms)` llamado; `intentosFallidos` reseteado a 0 en credenciales; próximo login tendrá 0ms de retraso |
| 15 | Login exitoso SIN fallos previos | `intentosFallidos=0`, login correcto | Login exitoso | `clock.Sleep` no llamado |
| 16 | Throttling NO aplica en refresh | — | `ServicioRefreshToken` | `clock.Sleep` no llamado |

### Límite de Sesiones Concurrentes

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 17 | Dentro del límite | 2 sesiones activas, `MAX_SESIONES=5` | Nuevo login | Sesión creada; 3 sesiones activas totales |
| 18 | Excede límite — Opción A | 5 sesiones activas, `MAX_SESIONES=5`, política=revocar_antigua | Nuevo login | Sesión más antigua revocada (`sesion.Revocar()` + `Actualizar`); nueva sesión creada; total=5 |
| 19 | Excede límite — Opción B | 5 sesiones activas, `MAX_SESIONES=5`, política=rechazar | Nuevo login | `ErrLimiteSesionesAlcanzado`; ninguna sesión creada ni revocada |

### Bloqueo Permanente

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 20 | Primer bloqueo temporal | `bloqueosConsecutivos=0`, `intentosFallidos=4` | 5to intento fallido | `intentosFallidos=5`, `bloqueadoHasta=now+15min`, `bloqueosConsecutivos=1`; `BloqueoUsuarioPort` NO llamado |
| 21 | Segundo bloqueo temporal | `bloqueosConsecutivos=1`, `intentosFallidos=4` | 5to intento fallido | `bloqueosConsecutivos=2`; `BloqueoUsuarioPort` NO llamado |
| 22 | Tercer bloqueo alcanza umbral | `bloqueosConsecutivos=2`, `intentosFallidos=4`, umbral=3 | 5to intento fallido | `bloqueosConsecutivos=3`; `BloqueoUsuarioPort.BloquearPermanentemente` llamado; error "cuenta bloqueada" |
| 23 | Login exitoso resetea bloqueosConsecutivos | `bloqueosConsecutivos=2`, login correcto | Login exitoso | `bloqueosConsecutivos=0`, `intentosFallidos=0` en credenciales actualizadas |
| 24 | Usuario ya en estado BLOQUEADO | Usuario con estado=BLOQUEADO | Login | `ErrCuentaBloqueada` (retornado en paso 3 del flujo, antes de verificar credenciales) |

### Política de Expiración

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 25 | Access token expira según TTL | Config: `accessTTL=15min` | `GenerarAccessToken` | `exp = now + 15min` en el retorno del mock de `TokenServicio` |
| 26 | Refresh token expira según TTL | Config: `refreshTTL=7d` | Login exitoso | `sesion.expiraRefresh = now + 7d` |
| 27 | Sesión expirada requiere re-login | `expiraRefresh < now` | Refresh | `ErrSesionExpirada` |

## Criterios de Aceptación

1. [ ] Rate limiting por IP se verifica antes de cualquier otra lógica, incluyendo validación de comando
2. [ ] Rate limiting por correo se verifica después de validar el formato del correo, antes de tocar repositorios
3. [ ] Throttling progresivo usa `clock.Sleep` (no `time.Sleep`) para ser testeable sin bloqueo
4. [ ] El retraso se calcula sobre los intentos PREVIOS (leídos antes de validar contraseña)
5. [ ] El retraso se aplica al final, tanto en éxito como en fallo
6. [ ] `bloqueosConsecutivos` se incrementa con cada bloqueo temporal y se resetea a 0 en login exitoso
7. [ ] Bloqueo permanente se dispara cuando `bloqueosConsecutivos >= umbral`, no antes
8. [ ] `BloqueoUsuarioPort` se llama en el mismo request que activa el umbral
9. [ ] Límite de sesiones concurrentes es configurable (política revocar_antigua o rechazar)
10. [ ] TTL de access y refresh token configurables externamente
11. [ ] Rate limiting y throttling no afectan `ServicioRefreshToken`
12. [ ] Todos los test cases pasan

## Tareas

### Interfaz Clock
- [ ] Definir interfaz `Clock` con `Now() time.Time` y `Sleep(d time.Duration)` en `internal/autenticacion/domain/`
- [ ] Implementar `RealClock` en infraestructura
- [ ] Implementar `MockClock` en tests (registra llamadas a `Sleep`, no bloquea)

### Rate Limiting — Interfaz
- [ ] Definir interfaz `RateLimiter` con `Permitir(clave string) (bool, time.Duration)` en `domain/`
- [ ] El servicio de aplicación recibe dos instancias de `RateLimiter`: una para IP, otra para correo
- [ ] Definir error `ErrRateLimitExcedido` con campo `TiempoRestante time.Duration`

### Rate Limiting — Integración en ServicioLogin
- [ ] Verificar rate limit por IP al inicio del método (antes de validar comando)
- [ ] Verificar rate limit por correo después de validar formato del correo (antes de tocar repositorios)
- [ ] Tests: casos 1-9

### Throttling Progresivo
- [ ] Leer `intentosFallidos` en paso 4 del flujo; guardar como `fallosPrevios`
- [ ] Calcular `retraso` con fórmula exponencial y tope
- [ ] Aplicar `clock.Sleep(retraso)` al final del método, antes de `return`
- [ ] Tests: casos 10-16

### Límite de Sesiones Concurrentes
- [ ] En `ServicioLogin`, paso 8: `ObtenerActivasPorUsuarioID`; si `len >= MAX_SESIONES`, aplicar política
- [ ] Política configurable mediante tipo enumerado o constante (`PoliticaRevocarAntigua`, `PoliticaRechazar`)
- [ ] Agregar `ErrLimiteSesionesAlcanzado`
- [ ] Tests: casos 17-19

### Bloqueo Permanente
- [ ] Agregar campo `bloqueosConsecutivos int` a entidad `CredencialesUsuario` en contexto `seguridad`
- [ ] Modificar `CredencialesUsuario.MarcarIntentoFallido`: si activa bloqueo temporal, incrementar `bloqueosConsecutivos`
- [ ] Agregar método `CredencialesUsuario.ResetearContadores()` que resetea `intentosFallidos=0`, `bloqueadoHasta=zero`, `bloqueosConsecutivos=0`
- [ ] `ServicioLogin`: en login exitoso, llamar `ResetearContadores()` en lugar de solo resetear `intentosFallidos`
- [ ] `ServicioLogin`: después de incrementar `bloqueosConsecutivos`, si `>= umbral`, llamar `BloqueoUsuarioPort.BloquearPermanentemente`
- [ ] Definir interfaz `BloqueoUsuarioPort` en `internal/autenticacion/domain/`
- [ ] El contexto `usuarios` implementa el puerto (inyectado en `registry`)
- [ ] Tests: casos 20-24

### Política de Expiración
- [ ] `ServicioLogin` y `ServicioRefreshToken` reciben `accessTTL time.Duration` y `refreshTTL time.Duration` como configuración
- [ ] Pasar `accessTTL` a `TokenServicio.GenerarAccessToken`
- [ ] Calcular `expiraRefresh = clock.Now().Add(refreshTTL)` al crear sesión
- [ ] Tests: casos 25-27

### Infraestructura (interfaces únicamente)
- [ ] La implementación concreta de `RateLimiter` (en memoria, Redis, etc.) se realizará en la etapa de infraestructura
- [ ] `RealClock` se registra en `registry` como implementación de `Clock`
