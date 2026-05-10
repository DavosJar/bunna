# Etapa 4: Auditoría y Monitoreo

## Prerrequisito

Etapas 1, 2 y 3 completas (login + sesiones + seguridad funcionando).

## Alcance

### Incluye
- Eventos de dominio de autenticación (login exitoso, login fallido, refresh, logout, bloqueo)
- Publicación de eventos para consumo externo (notificaciones, logging, métricas)
- Trazabilidad completa: quién, cuándo, desde dónde, qué dispositivo
- Evento de seguridad ante actividad sospechosa (anti-replay, IP desconocida, horario atípico)
- Interfaz de auditoría (repositorio de eventos de autenticación)
- Integración con el sistema de eventos existente (eventos.go pattern)

### NO incluye
- Dashboard / UI de monitoreo
- Sistema de notificaciones al usuario (emails, push)
- Implementación concreta del bus de eventos (infrastructure)
- Geolocalización de IPs (infrastructure externa)
- Machine learning para detección de anomalías

## Decisiones de Diseño

### Emisión de SesionExpirada es Lazy

Las sesiones no se auto-expiran. El evento SesionExpirada se emite cuando:
1. ServicioRefreshToken intenta usar una sesión y detecta que `time.Now() > expiraRefresh`
2. En ese momento, emite el evento y retorna error "sesión expirada"

Un job periódico (infraestructura, etapa posterior) puede marcar sesiones expiradas en batch, pero no es requerido para esta etapa.

### EventoAutenticacion es alias de EventoDominio

No hay tipo separado. Se usa el mismo `EventoDominio` existente. La diferencia es semántica: los eventos emitidos por el contexto autenticacion son eventos de autenticacion.

## Comportamiento esperado

### Eventos de Autenticación

Cada acción de autenticación debe emitir un evento de dominio con:

| Evento | Disparador | Payload mínimo |
|--------|-----------|----------------|
| `LoginExitoso` | Login correcto | usuarioID, sesionID, timestamp, IP, dispositivo |
| `LoginFallido` | Login con credenciales inválidas | correo (no usuarioID si no existe), IP, dispositivo, intentosFallidos actuales |
| `LoginBloqueado` | Login a cuenta bloqueada temporal o permanente | usuarioID, IP, tipoBloqueo (temporal/permanente), tiempoRestante |
| `RefreshExitoso` | Refresh token exitoso | usuarioID, sesionID, timestamp, viejoJTI, nuevoJTI |
| `RefreshReplay` | Intento de refresh con token ya rotado | usuarioID, sesionID, IP, timestamp |
| `LogoutExitoso` | Logout individual | usuarioID, sesionID, timestamp |
| `LogoutAllEjecutado` | Cierre de todas las sesiones | usuarioID, cantidadSesionesRevocadas, timestamp |
| `SesionExpirada` | Sesión alcanzó expiraRefresh (lazy: cuando se intenta usar) | usuarioID, sesionID, timestamp |
| `LimiteSesionesAlcanzado` | Se superó límite de sesiones concurrentes | usuarioID, sesionIDRevocada (si aplica), políticaAplicada |
| `BloqueoPermanenteActivado` | Usuario pasa a BLOQUEADO por umbral | usuarioID, intentosFallidosTotal, timestamp |
| `RateLimitExcedido` | IP o correo supera rate limit | clave (IP o correo), tipo, timestamp |

### Formato de Eventos

Los eventos usan el tipo `EventoDominio` existente (no hay tipo separado EventoAutenticacion).

Estructura:

```
Nombre: string (ej: "LoginExitoso")
Payload: map[string]interface{} con los datos del evento
Ocurrido: time.Time
```

Todos los eventos de autenticacion son EventoDominio con Nombre prefijado por el contexto (ej: "autenticacion.LoginExitoso") o sin prefijo si se usa el patrón existente.

### Consumidores de Eventos (futuro)

Los eventos se publican para que otros subsistemas los consuman:

1. **Logger**: persiste todos los eventos en tabla `eventos_autenticacion` (PostgreSQL en infra)
2. **Notificador**: cuando hay LoginFallido desde IP desconocida o LoginExitoso en horario atípico → enviar alerta
3. **Métricas**: contadores para dashboard (login success rate, intentos fallidos, etc.)
4. **Security Information**: detectar patrones anómalos (muchos logins fallidos desde varias IPs a un mismo usuario)

### Repositorio de Auditoría

Interfaz para persistir eventos de autenticación:

```
RegistrarEvento(ctx, evento EventoAutenticacion) error
ObtenerEventosPorUsuarioID(ctx, usuarioID string, desde time.Time, hasta time.Time, paginacion) ([]EventoAutenticacion, error)
ObtenerEventosPorTipo(ctx, tipo string, desde time.Time, hasta time.Time, paginacion) ([]EventoAutenticacion, error)
ObtenerIntentosFallidosPorIP(ctx, ip string, desde time.Time) (int, error)
```

## Especificación TDD

### Emisión de Eventos

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 1 | Login exitoso emite evento | Usuario ACTIVO, credenciales válidas | Login exitoso | Evento `LoginExitoso` emitido con usuarioID, sesionID, IP, dispositivo |
| 2 | Login fallido emite evento | Credenciales inválidas | Login fallido | Evento `LoginFallido` emitido con correo, IP, dispositivo, intentos actuales |
| 3 | Bloqueo temporal emite evento | 4 intentos previos | 5to fallo → bloqueo | Evento `LoginBloqueado` con tipo=temp, tiempoRestante≈15min |
| 4 | Bloqueo permanente emite evento | 3 bloqueos temporales previos | Siguiente bloqueo | Evento `BloqueoPermanenteActivado`, usuario cambia a estado BLOQUEADO |
| 5 | Refresh exitoso emite evento | Sesión activa | Refresh | Evento `RefreshExitoso` con viejoJTI y nuevoJTI |
| 6 | Replay detectado emite evento | Token ya rotado | Intento con token viejo | Evento `RefreshReplay` con usuarioID, sesionID, IP |
| 7 | Logout emite evento | Sesión activa | Logout | Evento `LogoutExitoso` con usuarioID, sesionID |
| 8 | LogoutAll emite evento | 3 sesiones activas | Logout all | Evento `LogoutAllEjecutado`, cantidadSesionesRevocadas=3 |
| 9 | Sesión expira emite evento | Sesión activa que alcanzó expiraRefresh | Check de expiración | Evento `SesionExpirada` con usuarioID, sesionID |
| 10 | Rate limit excedido emite evento | IP excede límite | Login bloqueado por rate limit | Evento `RateLimitExcedido` con clave=IP, tipo="ip" |

### Payload de Eventos

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 11 | LoginExitoso contiene todos los campos requeridos | — | Login exitoso | Payload tiene: usuarioID(string), sesionID(string), timestamp(time.Time), IP(string), dispositivo(string) |
| 12 | LoginFallido NO incluye usuarioID si correo no existe | Correo no registrado | Login | Payload tiene correo pero NO usuarioID |
| 13 | LoginFallido SÍ incluye usuarioID si correo existe | Correo existe, password incorrecta | Login | Payload tiene usuarioID Y correo |
| 14 | Eventos tienen timestamp correcto | — | Cualquier evento | Evento.Ocurrido ≈ time.Now() (con tolerancia de 1s) |

### Repositorio de Auditoría

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 15 | Registrar evento de login | Evento LoginExitoso | Repo.Crear | Evento persistido sin error |
| 16 | Consultar eventos por usuario y rango | 5 eventos para usuario X en los últimos 7 días | Repo.ObtenerPorUsuarioID | 5 eventos, ordenados por timestamp descendente |
| 17 | Consultar eventos por tipo | 10 LoginExitoso, 3 LoginFallido | Repo.ObtenerPorTipo("LoginFallido") | 3 eventos |
| 18 | Consultar intentos fallidos por IP | IP 1.1.1.1 tiene 8 intentos hoy | Repo.IntentosPorIP("1.1.1.1", inicioDelDia) | 8 intentos |

### Integración con Eventos Existentes

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 19 | Los eventos se integran con el patrón PullEventos | Servicio emite eventos | Se llama sesion.PullEventos() | Los eventos están disponibles y se limpian |
| 20 | Múltiples eventos en un solo flujo | Login exitoso | Flujo completo | LoginExitoso + SesionCreada (o evento interno) emitidos secuencialmente |

## Criterios de Aceptación

1. [ ] Cada acción de autenticación emite al menos un evento de dominio
2. [ ] Los eventos tienen payload mínimo completo (usuarioID, IP, timestamp, dispositivo)
3. [ ] LoginFallido NO expone usuarioID si el correo no existe (seguridad)
4. [ ] Los eventos siguen el patrón EventoDominio existente en el proyecto
5. [ ] El repositorio de auditoría permite consultar por usuario, tipo y rango de tiempo
6. [ ] El evento RefreshReplay permite detectar potenciales ataques
7. [ ] Todos los test cases pasan

## Tareas

### Eventos de Dominio
- [ ] Definir todos los tipos de evento como constantes (LoginExitoso, LoginFallido, etc.)
- [ ] Implementar generación de eventos en ServicioLogin (exitoso, fallido, bloqueado)
- [ ] Implementar generación de eventos en ServicioRefreshToken (exitoso, replay, sesionExpirada si aplica)
- [ ] Implementar generación de eventos en ServicioLogout (individual, all)
- [ ] SesionExpirada se emite lazy: cuando ServicioRefreshToken detecta que expiraRefresh < now
- [ ] Tests: casos 1-14

### Repositorio de Auditoría
- [ ] Definir interfaz `AuditoriaRepositorio` con métodos de registro y consulta
- [ ] La interfaz vive en `domain/` de autenticación
- [ ] Tests: casos 15-18

### Integración
- [ ] Conectar eventos de autenticación con el sistema PullEventos existente
- [ ] Los eventos deben ser accesibles desde el aggregate root de sesión
- [ ] Tests: casos 19-20
