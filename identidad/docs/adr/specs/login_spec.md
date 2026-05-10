# Especificación del Caso de Uso Login

> **Propósito**: Definir de forma incremental la especificación completa del login con JWT, refresh tokens, gestión de sesiones y seguridad perimetral, siguiendo la arquitectura limpia del proyecto (`dominio → aplicación → infraestructura`).  
> **No incluye**: presentación (handlers HTTP), auditoría, eventos de dominio, ni implementación concreta de infraestructura.  
> **Formato**: Una sola especificación fraccionada internamente en etapas incrementales. Cada etapa tiene sentido por sí misma pero encaja en el flujo global.

---

## Mapa de Dependencias entre Etapas

```
Etapa 1 (Sesiones - Dominio)
    └──→ Etapa 2 (Login - Aplicación)
            ├──→ Etapa 3 (Refresh Token - Aplicación)
            └──→ Etapa 4 (Logout - Aplicación)
                    └──→ Etapa 5 (Seguridad Perimetral)
                            └──→ Etapa 6 (Integración y Configuración)
```

Cada etapa lista sus prerrequisitos. No se debe iniciar una etapa sin completar la anterior que la antecede.

---

## Convenciones del Proyecto (Aplican a Todas las Etapas)

1. **Paquete de dominio**: `internal/sesiones/domain/` — Entidades, interfaces de repositorio, errores, value objects.
2. **Paquete de aplicación**: `internal/sesiones/application/services/<caso-uso>/` — Comando (input DTO), Respuesta (output DTO), servicio, test.
3. **Infraestructura**: `internal/sesiones/infrastructure/persistence/postgres/` — Implementación concreta del repositorio.
4. **Registry**: `internal/registry/registry.go` — Inyección de dependencias.
5. **Config**: `internal/config/env.go` — Variables de entorno.
6. **Transacciones**: Uso de `UnitOfWork` para operaciones que modifican multiples agregados (ej: login + sesión).
7. **Tests**: Mocks con `MockGeneradorID` y setups de BD real con `setupTestDB()`. Tests en el package de aplicación con BD PostgreSQL de prueba.
8. **Constructores**: `Nuevo<Entidad>(...)` para creación, `Nuevo<Entidad>DesdeBD(...)` para hidratación desde persistencia.
9. **Errores de dominio**: Variables exportadas estilo `ErrSesionExpirada`, usadas con `errors.New(...)` en el paquete de dominio.

---

# Etapa 1: Dominio de Sesiones

## Objetivo

Modelar la entidad `Sesion` en el dominio, con su ciclo de vida completo: creación, verificación de vigencia, extensión, revocación y expiración. Esta es la base de todo el login.

## Ubicación en el proyecto

```
internal/sesiones/domain/
├── sesion.go                    # Entidad raíz
├── sesion_repositorio.go       # Interfaz del repositorio
├── errores.go                   # Errores de dominio
└── tokens.go                    # Value Object: TokenPair (access + refresh)
```

## Consideraciones de Diseño

- `Sesion` es un agregado independiente. No cuelga de `Usuario` ni de `CredencialesUsuario`.
- Una sesión pertenece a un usuario (`usuarioID`), pero la entidad sólo necesita el ID, no el objeto completo.
- El ciclo de vida de una sesión es: `ACTIVA → EXPIRADA` (por tiempo) o `ACTIVA → REVOCADA` (por logout explícito o revocación administrativa).
- La sesión almacena tanto el **access token** como el **refresh token** (o sus hashes). Por seguridad se almacena el hash del refresh token, no el token en plano.
- Toda fecha se maneja como `time.Time`. La comparación contra el reloj del sistema se hace en el servicio de aplicación, no en el dominio. El dominio sólo pregunta "¿estás expirada?" contra una referencia de tiempo que recibe por parámetro.
- IP de origen se almacena como metadato. No se valida formato en dominio (es un string).
- El refresh token tiene su propia fecha de expiración, independiente del access token.
- `TokenPair` es un Value Object inmutable que agrupa access + refresh token + sus expiraciones.

## Contrato del Repositorio

```
SesionRepositorio (interfaz en dominio)

- Crear(ctx, sesion) → (Sesion, error)
- Actualizar(ctx, sesion) → (Sesion, error)
- ObtenerPorID(ctx, sesionID) → (Sesion, error)
- ObtenerPorRefreshTokenHash(ctx, refreshTokenHash) → (Sesion, error)
- ListarActivasPorUsuarioID(ctx, usuarioID, ahora time.Time) → ([]Sesion, error)
- Eliminar(ctx, sesionID) → error
```

## Escenarios de TDD para el Dominio de Sesiones

### Creación de sesión
1. **Happy path**: Crear sesión con datos válidos → sesión creada en estado ACTIVA, con fechas correctas, tokens asignados, no expirada.
2. **Happy path**: Crear sesión desde persistencia con `NuevaSesionDesdeBD` → reconstruye el estado exacto, no emite eventos, no valida.
3. **Edge: usuarioID vacío**: Crear sesión con `usuarioID = ""` → error `ErrUsuarioIDRequerido`.
4. **Edge: refreshTokenHash vacío**: Crear sesión con `refreshTokenHash = ""` → error `ErrRefreshTokenHashRequerido`.
5. **Edge: accessTokenHash vacío**: Crear sesión con `accessTokenHash = ""` → error `ErrAccessTokenHashRequerido`.
6. **Edge: expiración en el pasado**: Crear sesión con `fechaExpiracion` anterior a `fechaCreacion` → error del dominio (sesión inválida porque nunca podría estar activa).
7. **Edge: IP de origen vacía**: Debe permitirse (por ahora es metadato opcional).

### Verificación de vigencia
8. **Happy path**: Sesión activa con `fechaExpiracionAccess` en el futuro → `EstaActiva(ahora)` retorna `true`.
9. **Sad path**: Sesión con `fechaExpiracionAccess` en el pasado → `EstaActiva(ahora)` retorna `false`. El estado NO cambia automáticamente (el servicio lo cambia).
10. **Sad path**: Sesión en estado REVOCADA → `EstaActiva(ahora)` retorna `false` aunque la fecha no haya expirado.

### Cambio de estado
11. **Happy path**: Marcar sesión como expirada → cambia a estado EXPIRADA, `EstaActiva()` retorna `false`.
12. **Happy path**: Marcar sesión como revocada → cambia a estado REVOCADA, `EstaActiva()` retorna `false`.
13. **Edge: marcar como expirada cuando ya está revocada** → no debe cambiar (error de transición de estado o no-op).
14. **Edge: marcar como revocada cuando ya está expirada** → debe permitirse (por seguridad, revocar algo ya expirado no debería romper nada, pero el estado final debe ser REVOCADA si se requiere registro de revocación). Decisión: se permite, cambia a REVOCADA.

### Refresh token
15. **Happy path**: Refresh token no expirado → `RefreshTokenValido(ahora)` retorna `true`.
16. **Sad path**: Refresh token expirado → `RefreshTokenValido(ahora)` retorna `false`.
17. **Edge**: Sesión revocada pero refresh token no expirado → `RefreshTokenValido()` debe retornar `false` porque la sesión no está activa.
18. **Edge**: Refresh token con fecha zero → interpretar como que expiró o nunca se generó, retornar `false`.

### Value Object TokenPair
19. **Happy path**: Crear TokenPair con todos los valores → getters retornan los valores correctos.
20. **Edge**: Access token string vacío → error.
21. **Edge**: Refresh token string vacío → error.

## Actividades de la Etapa 1

1. Crear el directorio `internal/sesiones/domain/`.
2. Definir las constantes de estado de sesión (`ACTIVA`, `EXPIRADA`, `REVOCADA`), tipo `EstadoSesion`.
3. Implementar entidad `Sesion` con sus constructores y métodos de comportamiento.
4. Implementar value object `TokenPair` inmutable.
5. Definir errores de dominio exportados.
6. Definir interfaz `SesionRepositorio`.
7. Escribir tests del dominio que cubran los 21 escenarios listados.
8. Verificar que todos los tests pasen (100% cobertura de los escenarios).

## Checklist de Validación de Etapa 1

- [ ] ¿La entidad Sesion tiene estado explícito (activa/expirada/revocada)?
- [ ] ¿El refresh token se almacena como hash, no en plano?
- [ ] ¿TokenPair es inmutable?
- [ ] ¿El repositorio es una interfaz en el dominio?
- [ ] ¿Los constructores `Nuevo...` y `Nuevo...DesdeBD` existen?
- [ ] ¿Hay getters públicos para todos los campos privados?
- [ ] ¿La sesión puede decir "estoy activa" contra una referencia de tiempo externa?
- [ ] ¿El refresh token puede decir "estoy válido" contra una referencia de tiempo externa?
- [ ] ¿Hay tests que cubren los 21 escenarios?
- [ ] ¿No hay dependencias de infraestructura en el dominio?

---

# Etapa 2: Login — Servicio de Aplicación

## Objetivo

Implementar el caso de uso "Iniciar Sesión" en la capa de aplicación. Recibe credenciales, las valida contra el repositorio de credenciales, y si son correctas crea una sesión con su par de tokens.

## Ubicación en el proyecto

```
internal/sesiones/application/services/login/
├── comando.go           # ComandoLogin
├── respuesta.go         # RespuestaLogin
├── servicio_login.go    # ServicioLogin
└── servicio_login_test.go
```

## Prerrequisitos

- Etapa 1 completada (dominio de sesiones).
- Ya existe: `CredencialesUsuario`, `CredencialesRepositorio`, `EncriptacionServicio` en `seguridad/domain/`.
- Ya existe: `UnitOfWork` en `usuarios/domain/usuario/unit_of_work.go`. Esta etapa requiere extenderlo o crear uno propio que incluya `SesionRepositorio`.

## Dependencias

- `CredencialesRepositorio` — para obtener credenciales y verificar password.
- `EncriptacionServicio` — para verificar password contra hash.
- `SesionRepositorio` — para persistir la sesión creada.
- `UnitOfWork` — para atomicidad: login exitoso = sesión creada + intentos reseteados; login fallido = intento incrementado.
- `GeneradorID` — para generar ID de sesión.
- **TokenServicio** (interfaz de infraestructura) — para generar el par de tokens JWT. Esta interfaz se define en dominio de sesiones pero su implementación concreta (JWT) está en infraestructura.

## Interfaces Requeridas (definir en dominio de sesiones o en un paquete compartido)

```
type TokenServicio interface {
    GenerarAccessToken(usuarioID, sesionID string) (tokenString string, expira time.Time, err error)
    GenerarRefreshToken(usuarioID, sesionID string) (tokenString string, expira time.Time, err error)
    ValidarAccessToken(tokenString string) (claims *TokenClaims, err error)
    ValidarRefreshToken(tokenString string) (claims *TokenClaims, err error)
    HashearToken(tokenString string) string
}

type TokenClaims struct {
    UsuarioID string
    SesionID  string
    Tipo      string // "access" o "refresh"
    Expira    time.Time
}
```

## Flujo del Servicio de Login (Dentro de Transacción)

```
1. Validar comando (fuera de transacción):
   - email no vacío, formato válido
   - password no vacío

2. Transacción (UnitOfWork):
   a. Obtener credenciales por email (requiere búsqueda inversa: email → usuarioID → credenciales)
      NOTA: Se necesita un método nuevo en algún repositorio: ObtenerCredencialesPorEmail o un servicio que resuelva email → usuarioID.
   b. Si no existen credenciales → error "credenciales inválidas" (genérico, no revelar existencia)
   c. Verificar si credenciales están bloqueadas → error "cuenta bloqueada"
   d. Verificar si credenciales están activas → error "cuenta inactiva"
   e. Verificar password contra hash (usando EncriptacionServicio)
   f. Si password incorrecto:
      - Marcar intento fallido en credenciales
      - Actualizar credenciales
      - Error genérico "credenciales inválidas"
   g. Si password correcto:
      - Resetear intentos fallidos
      - Generar par de tokens (access + refresh) via TokenServicio
      - Hashear refresh token
      - Crear entidad Sesion con los hashes
      - Persistir sesión
      - Actualizar credenciales (intentos reseteados)
   h. Retornar DTO con tokens (en plano), usuarioID, sesionID, expiraciones
```

## Escenarios de TDD para Login

### Happy path
1. **Login exitoso**: Credenciales válidas, password correcto → retorna DTO con tokens, sesión creada en BD, intentos fallidos resetean a 0.
2. **Login exitoso después de reintentos**: Usuario con 3 intentos fallidos previos, ahora password correcto → login OK, intentos resetean a 0.

### Validaciones
3. **Email vacío**: Comando con email vacío → error de validación.
4. **Email mal formado**: Comando con formato inválido → error de validación.
5. **Password vacío**: Comando con password vacío → error de validación.

### Sad path: Estado de cuenta
6. **Credenciales no existen**: Email no registrado → error genérico "credenciales inválidas" (no revelar si el email existe o no).
7. **Cuenta bloqueada**: Credenciales con `bloqueadoHasta > ahora` → error específico "cuenta temporalmente bloqueada".
8. **Bloqueo expirado**: Credenciales con `bloqueadoHasta < ahora` → debe permitir login (el bloqueo ya venció), pero ¿los intentos se resetean? Decisión: NO se resetean automáticamente hasta que haya un intento exitoso. El servicio debe evaluar si el bloqueo expiró y tratar el intento actual normalmente.
9. **Cuenta inactiva**: Credenciales con `activo = false` → error "cuenta inactiva".
10. **Correo no verificado**: Credenciales con `correoVerificado = false` → ¿se permite login? Decisión: depende de la política de negocio. Inicialmente configurable, por defecto SÍ se permite login pero el servicio debe poder consultar esta política.

### Sad path: Password
11. **Password incorrecto**: Password no coincide → intento fallido incrementado, error genérico "credenciales inválidas".
12. **Password incorrecto y alcanza bloqueo**: 5to intento fallido consecutivo → se bloquea la cuenta por 15 minutos, error "cuenta temporalmente bloqueada".
13. **Password incorrecto con cuenta ya bloqueada**: Intentar login estando bloqueado → error "cuenta temporalmente bloqueada", NO se incrementa el contador (no penalizar intentos durante bloqueo).

### Edge: Transaccionalidad
14. **Fallo al crear sesión después de validar password**: Simular error en repositorio de sesión → rollback: el intento fallido NO debe persistirse, el password correcto NO debe quedar sin sesión.
15. **Fallo al actualizar credenciales después de crear sesión**: Simular error → rollback: sesión no creada, intentos no modificados.
16. **Context timeout**: Context cancelado durante la transacción → rollback completo, error de timeout.

### Edge: Tokens
17. **TokenServicio falla al generar access token**: Error desde infraestructura → rollback, error de generación de tokens.
18. **TokenServicio falla al generar refresh token**: Error desde infraestructura → rollback, error de generación de tokens.
19. **Hash de refresh token colisiona**: Caso extremo donde el hash de un refresh token nuevo coincide con uno existente. Decisión: por ahora no se maneja (los hashes criptográficos tienen colisión despreciable), pero se documenta como riesgo futuro.

## Actividades de la Etapa 2

1. Definir la interfaz `TokenServicio` en el dominio de sesiones (o en un paquete `internal/sesiones/domain/`).
2. Definir `TokenClaims` como struct de dominio.
3. Agregar al `SesionRepositorio` el método `ObtenerPorEmail` o crear un servicio de aplicación que resuelva email → usuarioID vía `UsuarioRepositorio`.
4. Crear directorio `internal/sesiones/application/services/login/`.
5. Implementar `ComandoLogin` con campos: `Email`, `Password`.
6. Implementar `RespuestaLogin` con campos: `AccessToken`, `RefreshToken`, `ExpiracionAccess`, `ExpiracionRefresh`, `UsuarioID`, `SesionID`.
7. Implementar `ServicioLogin` con el flujo descrito.
8. Escribir tests unitarios con mocks de `CredencialesRepositorio`, `SesionRepositorio`, `TokenServicio`, `EncriptacionServicio`.
9. Escribir tests de integración con BD real (setupTestDB).
10. Extender `UnitOfWork` para incluir `SesionRepositorio` y `TokenServicio`.

## Checklist de Validación de Etapa 2

- [ ] ¿El comando de login valida email y password antes de abrir transacción?
- [ ] ¿El error para "credenciales no existen" es genérico (no revela existencia)?
- [ ] ¿El error para "password incorrecto" es genérico (idem)?
- [ ] ¿Los intentos fallidos se incrementan SOLO si el password es incorrecto?
- [ ] ¿Los intentos fallidos NO se incrementan si la cuenta ya está bloqueada?
- [ ] ¿Los intentos fallidos se resetean AL CREAR la sesión exitosamente (no antes)?
- [ ] ¿La transacción es atómica: o se crea todo o no queda nada?
- [ ] ¿El refresh token se almacena hasheado, no en plano?
- [ ] ¿Hay tests para cada escenario listado?
- [ ] ¿Se usa `TokenServicio` como interfaz (no implementación concreta)?

---

# Etapa 3: Refresh Token — Servicio de Aplicación

## Objetivo

Implementar el caso de uso "Renovar Sesión" (refresh token). Recibe un refresh token válido, lo invalida (rotación) y genera un nuevo par de tokens para la misma sesión.

## Ubicación en el proyecto

```
internal/sesiones/application/services/refresh/
├── comando.go           # ComandoRefresh
├── respuesta.go         # RespuestaRefresh
├── servicio_refresh.go  # ServicioRefresh
└── servicio_refresh_test.go
```

## Prerrequisitos

- Etapa 1 y 2 completadas.
- `TokenServicio` implementado (o mockeable) con `ValidarRefreshToken` y `HashearToken`.

## Consideraciones de Diseño

- **Rotación obligatoria**: Cada vez que se usa un refresh token, este se invalida y se genera uno nuevo. Esto previene que un refresh token robado pueda ser reutilizado.
- **Detección de robo**: Si un refresh token ya rotado es presentado nuevamente, significa que alguien más lo tiene. En ese caso, se deben invalidar TODAS las sesiones del usuario (medida de seguridad severa). Esta es una decisión de negocio importante.
- **Límite de refrescos por sesión**: Opcional. Se puede configurar un máximo de refrescos por sesión (ej: 10) para forzar re-login periódico.
- **La sesión NO cambia de ID**: El refresh genera nuevos tokens pero la misma sesión (mismo `sesionID`). Esto permite rastrear la sesión a lo largo de sus refrescos.
- **Contador de refrescos**: La entidad `Sesion` debe llevar un contador de cuántas veces se ha refrescado.

## Flujo del Servicio de Refresh (Dentro de Transacción)

```
1. Validar comando (fuera de transacción):
   - refreshToken no vacío

2. Validar refresh token: (Siempre antes de transacción)
   a. Llamar a TokenServicio.ValidarRefreshToken(token)
   b. Si el token es inválido (expirado, mal formado, firma inválida) → error
   c. Si es válido → extraer claims: usuarioID, sesionID. Calcular refreshTokenHash del token recibido.

3. Transacción (UnitOfWork):
   a. Buscar sesión por refreshTokenHash (calculado en paso 2c). Opcionalmente filtrar por usuarioID de los claims para eficiencia.
   b. Si NO se encuentra la sesión → el token JWT es válido pero su hash no existe en BD (fue rotado). Esto es detección de robo:
      - Se tiene el usuarioID desde los claims del paso 2c
      - Invalidar TODAS las sesiones activas de ese usuarioID
      - Error genérico "token inválido" (no informar al atacante)
   c. Si la sesión existe pero está REVOCADA o EXPIRADA:
      - Error "sesión no válida"
   d. Si la sesión existe y está ACTIVA:
      - Verificar que el refresh token de la sesión no haya expirado
      - Si expiró → marcar sesión como EXPIRADA → error "refresh token expirado"
      - Generar nuevo par de tokens (access + refresh)
      - Hashear el nuevo refresh token
      - Actualizar sesión: nuevo accessTokenHash, nuevo refreshTokenHash, nuevas fechas de expiración, incrementar contador de refrescos
      - Persistir sesión actualizada
   e. Retornar nuevo DTO con tokens
```

## Escenarios de TDD para Refresh

### Happy path
1. **Refresh exitoso**: Refresh token válido, sesión activa → nuevo par de tokens, sesión actualizada, contador incrementado en 1.
2. **Refresh exitoso múltiples veces**: Refrescar 3 veces seguidas → cada vez nuevo par, contador = 3.

### Validaciones
3. **Refresh token vacío**: Comando con token vacío → error de validación.

### Sad path: Token inválido
4. **Refresh token expirado**: Token firmado pero con fecha expirada → error "token inválido o expirado".
5. **Refresh token mal formado**: Token que no puede ser parseado → error "token inválido o expirado".
6. **Refresh token con firma inválida**: Token alterado → error "token inválido o expirado".

### Sad path: Sesión inválida
7. **Refresh token válido pero sesión REVOCADA**: Token firmado correctamente, sesión existe pero fue revocada → error "sesión no válida".
8. **Refresh token válido pero sesión EXPIRADA**: Token firmado, sesión expiró → error "sesión no válida".
9. **Sesión no encontrada por hash**: Token válido pero su hash no corresponde a ninguna sesión (caso de rotación previa) → error "token inválido", y considerar detección de robo.

### Sad path: Límites
10. **Exceder límite de refrescos**: Si se configura un máximo y se alcanza → error "límite de refrescos alcanzado, inicie sesión nuevamente".
11. **Refresh token válido pero fecha de sesión expiró**: El token en sí no expiró pero la sesión como conjunto tiene una expiración absoluta (ej: 7 días desde creación) → error "sesión expirada".

### Edge: Detección de robo
12. **Reutilización de refresh token rotado (detección de robo)**: El token JWT es válido (firma correcta, no expirado), se extrae usuarioID y sesionID de los claims, pero al buscar el hash del token presentado en BD no se encuentra ninguna sesión (el hash fue reemplazado durante la rotación anterior). Con el usuarioID de los claims se invalidan todas las sesiones activas de ese usuario. Error genérico "token inválido".
13. **Usuario sin sesiones activas después de detección de robo**: La invalidación masiva deja 0 sesiones activas → correcto.

### Edge: Transaccionalidad
14. **Fallo al persistir sesión actualizada**: Error en BD después de generar nuevos tokens → rollback, tokens no válidos, sesión en estado anterior.
15. **Fallo en generación de tokens**: TokenServicio falla → rollback, sesión intacta.

## Actividades de la Etapa 3

1. Agregar campo `contadorRefrescos` a la entidad `Sesion` (entero, 0 por defecto).
2. Agregar método `IncrementarContadorRefrescos()` y getter `ContadorRefrescos()`.
3. Agregar método `RotarTokens(nuevoAccessHash, nuevoRefreshHash, nuevaExpiracionAccess, nuevaExpiracionRefresh)`.
4. Agregar al repositorio: `ObtenerPorRefreshTokenHash`.
5. Agregar al repositorio: `InvalidarTodasPorUsuarioID(usuarioID)` (para detección de robo).
6. Crear directorio `internal/sesiones/application/services/refresh/`.
7. Implementar `ComandoRefresh`, `RespuestaRefresh`, `ServicioRefresh`.
8. Implementar el algoritmo de detección de robo.
9. Escribir tests unitarios (mocks).
10. Escribir tests de integración (BD real).

## Checklist de Validación de Etapa 3

- [ ] ¿El refresh token se rota en cada uso (nuevo token, viejo invalidado)?
- [ ] ¿Hay detección de reutilización de tokens rotados?
- [ ] ¿La detección de robo invalida TODAS las sesiones del usuario?
- [ ] ¿El contador de refrescos se incrementa correctamente?
- [ ] ¿Hay límite configurable de refrescos por sesión?
- [ ] ¿La sesión tiene una expiración absoluta además de la del token?
- [ ] ¿Los tests cubren detección de robo?
- [ ] ¿El mensaje de error en detección de robo es genérico (no informar al atacante)?

---

# Etapa 4: Cierre de Sesión (Logout) — Servicio de Aplicación

## Objetivo

Implementar el cierre de sesión explícito (logout) y el cierre de sesión por inactividad (timeout).

## Ubicación en el proyecto

```
internal/sesiones/application/services/logout/
├── comando.go            # ComandoLogout
├── respuesta.go          # RespuestaLogout
├── servicio_logout.go    # ServicioLogout
└── servicio_logout_test.go
```

## Prerrequisitos

- Etapa 2 completada.
- El middleware de autenticación (en presentación, fuera del alcance actual) extrae `sesionID` del token.

## Consideraciones de Diseño

- **Logout explícito**: El usuario solicita cerrar sesión. Se revoca la sesión actual.
- **Logout de todas las sesiones**: El usuario puede cerrar todas sus sesiones activas (ej: "cerrar sesión en todos los dispositivos").
- **Logout por tiempo (timeout de inactividad)**: No es un servicio de aplicación en sí, es una validación que ocurre en el middleware o en el servicio de refresh. Se especifica aquí como concepto.
- **El logout no requiere contraseña**: Solo requiere el `sesionID` (que viene del token) o el `usuarioID`.

## Flujo del Servicio de Logout

### Logout de una sesión específica
```
1. Recibir sesionID (del token autenticado o del comando)
2. Validar que la sesión pertenezca al usuario autenticado
3. Marcar sesión como REVOCADA
4. Persistir cambio
```

### Logout de todas las sesiones
```
1. Recibir usuarioID
2. Listar todas las sesiones activas del usuario
3. Marcar todas como REVOCADAS
4. Persistir cambios (batch)
```

## Escenarios de TDD para Logout

### Happy path
1. **Logout sesión específica**: Sesión activa, se revoca → estado REVOCADA, followers ya no pueden usarla.
2. **Logout de todas las sesiones**: Usuario con 3 sesiones activas → las 3 pasan a REVOCADAS.
3. **Logout y luego intentar refresh**: Sesión revocada, se intenta refresh → error "sesión no válida".

### Sad path
4. **Logout de sesión ya expirada**: Sesión EXPIRADA, se solicita logout → se permite (cambio a REVOCADA o no-op). Decisión: no-op, no hay necesidad de cambiar estado.
5. **Logout de sesión ya revocada**: Sesión ya REVOCADA → no-op, no hay error.
6. **Logout de sesión que no pertenece al usuario**: SesionID de otro usuario → error "no autorizado".
7. **Sesión no encontrada**: SesionID inexistente → error "sesión no encontrada".

### Edge: Timeout (validación en aplicación)
8. **Sesión con timeout de inactividad**: Sesión con `ultimaActividad` + `timeoutInactividad < ahora` → la sesión debe marcarse como EXPIRADA al detectarse.
9. **Timeout configurable**: El tiempo de inactividad es un parámetro de configuración (no hardcodeado).

## Actividades de la Etapa 4

1. Agregar campo `ultimaActividad` a la entidad `Sesion`.
2. Agregar método `RegistrarActividad(ahora)` a `Sesion`.
3. Agregar método `TimeoutExcedido(ahora, timeoutDuracion) bool`.
4. Implementar `ServicioLogout` con operaciones de revocación individual y masiva.
5. Implementar lógica de timeout de inactividad (validación en aplicación, no middleware).
6. Escribir tests para cada escenario.
7. Escribir test que verifique que después de logout no se puede hacer refresh.

## Checklist de Validación de Etapa 4

- [ ] ¿El logout explícito revoca la sesión (no la elimina)?
- [ ] ¿Hay operación de "cerrar todas las sesiones"?
- [ ] ¿El campo `ultimaActividad` se actualiza en cada operación de la sesión?
- [ ] ¿La validación de timeout se puede invocar desde aplicación?
- [ ] ¿Los tests cubren que un refresh post-logout falla?
- [ ] ¿No se puede hacer logout de una sesión que no pertenece al usuario?
- [ ] ¿Hacer logout de una sesión ya expirada es no-op (no cambia estado)?
- [ ] ¿Hacer logout de una sesión ya revocada es no-op (no cambia estado)?

---

# Etapa 5: Seguridad Perimetral

## Objetivo

Implementar las protecciones de seguridad a nivel de aplicación: bloqueo por IP, rate limiting, timeouts configurables. Esta etapa NO toca el dominio existente, sino que agrega servicios de aplicación e infraestructura.

## Ubicación en el proyecto

```
internal/seguridad/application/services/
├── bloqueo_ip/
│   ├── comando.go
│   ├── respuesta.go
│   └── servicio_bloqueo_ip.go
├── rate_limiter/
│   ├── comando.go
│   ├── respuesta.go
│   └── servicio_rate_limiter.go
```

O alternativamente como parte del servicio de login existente (ver consideraciones).

## Prerrequisitos

- Etapas 1-4 completadas.
- `CredencialesUsuario` ya tiene `intentosFallidos` y `bloqueadoHasta`.

## Consideraciones de Diseño

### IP Blocking
- **Propósito**: Prevenir ataques de fuerza bruta desde una misma dirección IP.
- **Alcance**: Es independiente del bloqueo por usuario. Un atacante puede probar 5 usuarios desde una IP y bloquear la IP, no los usuarios.
- **Umbral**: N intentos fallidos desde una misma IP en una ventana de tiempo (ej: 20 intentos en 15 minutos).
- **Duración del bloqueo**: Tiempo configurable (ej: 30 minutos).
- **Almacenamiento**: Preferiblemente en Redis (para expiración automática). Como etapa inicial, puede ser en PostgreSQL con limpieza periódica o en memoria (no recomendado para producción).
- **Decisión para esta etapa**: Se modela como un servicio de aplicación con interfaz de repositorio. La implementación concreta se decide en Etapa 6.

### Rate Limiting
- **Propósito**: Limitar la cantidad de requests de login desde una misma IP en una ventana de tiempo.
- **Diferenciación**: Rate limiting es preventivo (limita requests antes de procesarlos), IP blocking es reactivo (bloquea después de detectar patrón abusivo).
- **Umbral**: X requests por minuto por IP (ej: 10 requests/min).
- **Por usuario**: Opcionalmente, limitar intentos por usuario+IP combinado.

### Timeouts
- **De inactividad**: Tiempo máximo entre operaciones de una sesión (ya cubierto en Etapa 4).
- **Absoluto de sesión**: Tiempo máximo de vida de una sesión (ej: 7 días) independientemente de actividad.
- **De refresh token**: Tiempo de vida del refresh token (ej: 24 horas).

## Escenarios de TDD para Seguridad

### IP Blocking
1. **IP no bloqueada**: IP sin intentos fallidos → permitir login.
2. **IP bloqueada por exceder umbral**: 20 intentos fallidos desde misma IP en ventana de 15 min → bloqueada por 30 min.
3. **IP bloqueada pero con intentos exitosos parciales**: 15 intentos fallidos desde IP, login exitoso desde misma IP → ¿se resetea el contador de IP? Decisión: NO, el contador de IP solo se resetea cuando expira la ventana de tiempo.
4. **IP bloqueada no impide login de usuario legítimo**: Usuario legítimo desde IP bloqueada → debe fallar (esa es la gracia del bloqueo por IP). El usuario debe usar otra red o esperar.
5. **Bloqueo expirado**: IP bloqueada, pasa el tiempo de bloqueo → permitir login nuevamente.
6. **Limpieza de registros antiguos**: Los intentos fallidos por IP deben limpiarse (TTL o job).

### Rate Limiting
7. **Dentro del límite**: 5 requests en un minuto con límite de 10/min → permitir.
8. **Límite excedido**: 11 requests en un minuto → error "demasiados intentos, intente más tarde".
9. **Ventana de tiempo deslizante**: Request en minuto 0, luego 10 requests en minuto 1 → límite evaluado correctamente con ventana deslizante (no fija).
10. **Reset después de ventana**: 11 requests, esperar 1 minuto, nuevo request → permitir.

### Timeouts (configurables)
11. **Timeout de sesía excedido**: Sesión con más de 7 días de creada → debe ser inválida (el servicio de refresh debe rechazar).
12. **Timeout de inactividad excedido**: Sesión sin actividad por más de 30 min → debe marcarse como expirada al siguiente uso.
13. **Timeout de refresh token excedido**: Refresh token con más de 24 horas → inválido, forzar re-login.

## Actividades de la Etapa 5

1. Definir entidad o tabla de `IntentoPorIP` (IP, contador, ventana_inicio, bloqueado_hasta).
2. Definir repositorio `IntentoIPRepositorio` en dominio de seguridad.
3. Implementar `ServicioBloqueoIP` que verifique y registre intentos por IP.
4. Integrar el servicio de bloqueo IP en el flujo de login (antes de validar credenciales).
5. Definir servicio de rate limiting (puede ser un middleware o un servicio inyectado).
6. Hacer configurables todos los valores: umbrales, ventanas, duraciones.
7. Implementar limpieza de registros antiguos (job programado o TTL).
8. Escribir tests para cada escenario de seguridad.

## Integración con el Login (Flujo completo con seguridad)

```
1. Rate limiting: verificar si la IP ha excedido el límite de requests/min
2. IP blocking: verificar si la IP está bloqueada
3. Bloqueo por usuario: verificar si las credenciales están bloqueadas (Etapa 2)
4. Validar password y crear sesión (Etapa 2)
5. Si falla: registrar intento fallido por IP y por usuario
6. Si el intento fallido por IP excede umbral: bloquear IP
```

## Checklist de Validación de Etapa 5

- [ ] ¿El bloqueo por IP es independiente del bloqueo por usuario?
- [ ] ¿Los umbrales son configurables (no hardcode)?
- [ ] ¿Las ventanas de tiempo son configurables?
- [ ] ¿El rate limiting es preventivo (antes de procesar)?
- [ ] ¿Los registros de IP tienen limpieza automática?
- [ ] ¿Los errores de rate limiting y bloqueo son informativos pero sin revelar detalles?
- [ ] ¿Los tests cubren ventanas deslizantes (no solo fijas)?
- [ ] ¿El timeout absoluto de sesión se valida en cada operación?

---

# Etapa 6: Integración y Configuración

## Objetivo

Conectar todos los componentes: implementar `TokenServicio` con JWT, crear las tablas en PostgreSQL, configurar variables de entorno, registrar todo en el Registry, y exponer los servicios para que la capa de presentación los consuma.

## Prerrequisitos

- Etapas 1-5 completadas.
- Decisiones de infraestructura tomadas.

---

## Decisiones de Infraestructura (Fijo)

Estas decisiones son vinculantes para la implementación. No deben cambiarse sin una revisión explícita del equipo.

### Algoritmo y Formato JWT

- **Algoritmo**: HMAC-SHA256 (HS256).
- **Clave**: Secreta, configurable vía `JWT_SECRET`.
- **Claims del access token**: `sub` (usuarioID), `sid` (sesionID), `iat`, `exp`, `typ` ("access").
- **Claims del refresh token**: `sub` (usuarioID), `sid` (sesionID), `iat`, `exp`, `typ` ("refresh"), `jti` (ID único para detección de reuso).
- **Expiración access token**: 15 minutos (configurable).
- **Expiración refresh token**: 24 horas (configurable).
- **Hash del refresh token**: SHA-256 del token antes de almacenar en BD. El token en plano nunca se persiste.

### Esquema de Tablas PostgreSQL

```sql
CREATE TABLE sesiones (
    id VARCHAR(36) PRIMARY KEY,
    usuario_id VARCHAR(36) NOT NULL,
    access_token_hash VARCHAR(64) NOT NULL,
    refresh_token_hash VARCHAR(64) NOT NULL,
    estado VARCHAR(20) NOT NULL DEFAULT 'ACTIVA',
    ip_origen VARCHAR(45),
    fecha_creacion TIMESTAMP NOT NULL,
    fecha_actualizacion TIMESTAMP NOT NULL,
    fecha_expiracion_access TIMESTAMP NOT NULL,
    fecha_expiracion_refresh TIMESTAMP NOT NULL,
    ultima_actividad TIMESTAMP NOT NULL,
    contador_refrescos INT NOT NULL DEFAULT 0
);

CREATE TABLE intentos_por_ip (
    id VARCHAR(36) PRIMARY KEY,
    ip VARCHAR(45) NOT NULL,
    contador INT NOT NULL DEFAULT 1,
    ventana_inicio TIMESTAMP NOT NULL,
    bloqueado_hasta TIMESTAMP,
    fecha_creacion TIMESTAMP NOT NULL,
    fecha_actualizacion TIMESTAMP NOT NULL
);

CREATE INDEX idx_intentos_ip ON intentos_por_ip(ip);
```

### Variables de Entorno

| Variable | Default | Descripción |
|---|---|---|
| `JWT_SECRET` | (requerido) | Clave HMAC para firmar JWT |
| `JWT_ACCESS_EXPIRACION` | `15m` | Duración del access token |
| `JWT_REFRESH_EXPIRACION` | `24h` | Duración del refresh token |
| `SESION_TIMEOUT_INACTIVIDAD` | `30m` | Tiempo máximo entre operaciones |
| `SESION_TIMEOUT_ABSOLUTO` | `168h` (7d) | Tiempo máximo de vida de la sesión |
| `SESION_MAX_REFRESCOS` | `0` (sin límite) | 0 = sin límite |
| `BLOQUEO_IP_MAX_INTENTOS` | `20` | Intentos fallidos antes de bloquear IP |
| `BLOQUEO_IP_VENTANA` | `15m` | Ventana para contar intentos por IP |
| `BLOQUEO_IP_DURACION` | `30m` | Duración del bloqueo de IP |
| `RATE_LIMIT_MAX_REQUESTS` | `10` | Máximo requests por ventana por IP |
| `RATE_LIMIT_VENTANA` | `1m` | Ventana de rate limiting |

---

## Tareas de Implementación (Orientativo)

El orden y la forma exacta de estas tareas puede ajustarse durante la implementación. Lo fijo son las interfaces y contratos definidos en etapas anteriores.

### Ubicación en el proyecto

```
internal/sesiones/infrastructure/
├── persistence/postgres/
│   ├── sesion_repositorio.go
│   └── sesion_model.go
├── security/jwt/
│   ├── jwt_token_servicio.go
│   └── jwt_token_servicio_test.go

internal/config/
├── env.go           # Nuevas variables de entorno

internal/registry/
├── registry.go      # Nuevas dependencias
```

### Registry

Actualizar `internal/registry/registry.go` para incluir:

- `SesionRepositorio` → instancia de `postgres.NewSesionRepositorio(db)`
- `TokenServicio` → instancia de `jwt.NewJWTTokenServicio(cfg)`
- `IntentoIPRepositorio` → instancia de `postgres.NewIntentoIPRepositorio(db)` (o Redis, según decisión del equipo)
- `ServicioLogin` → instancia con repositorios + unit of work + token servicio
- `ServicioRefresh` → instancia con repositorios + unit of work + token servicio
- `ServicioLogout` → instancia con repositorios + unit of work
- `ServicioBloqueoIP` → instancia con repositorio de IP

### Migraciones

Agregar migraciones automáticas para las nuevas tablas en el proceso de inicialización de BD (donde ya se ejecutan las migraciones de usuarios y credenciales).

### Actividades de la Etapa 6

1. Implementar `postgres.SesionRepositorio` con GORM, respetando la interfaz `SesionRepositorio` del dominio.
2. Implementar `postgres.IntentoIPRepositorio` con GORM (o Redis si se opta por esa vía).
3. Implementar `jwt.JWTTokenServicio` usando `github.com/golang-jwt/jwt/v5`.
4. Agregar las nuevas variables de entorno a `config/env.go` y `config.LoadConfig()`.
5. Agregar migraciones para `sesiones` e `intentos_por_ip` en el inicializador de BD.
6. Extender o crear `UnitOfWork` que incluya `SesionRepositorio`.
7. Actualizar `Registry` con todas las nuevas dependencias y servicios de aplicación.
8. Escribir tests de integración para cada repositorio (estilo `setupTestDB`).
9. Escribir tests para `JWTTokenServicio`: generación, validación, expiración, firma inválida, claims correctos.
10. Verificar que todos los tests de todas las etapas pasan integradamente.

## Checklist de Validación de Etapa 6

- [ ] ¿La implementación JWT implementa la interfaz `TokenServicio` definida en dominio?
- [ ] ¿La clave secreta JWT se lee de variable de entorno, no está hardcodeada?
- [ ] ¿Las tablas de sesiones e intentos por IP tienen los índices definidos en el esquema?
- [ ] ¿El refresh token se almacena como hash SHA-256 (no en plano)?
- [ ] ¿Hay migraciones automáticas al iniciar la aplicación?
- [ ] ¿Todas las nuevas dependencias están registradas en el Registry?
- [ ] ¿Los tests de integración usan BD real (setupTestDB) como los existentes?
- [ ] ¿Se respeta la estructura `dominio → aplicación → infraestructura` sin fugas de dependencias?
- [ ] ¿Los mensajes de error en la capa de infraestructura no exponen detalles internos?

---

# Resumen de Etapas y Prioridades

| Etapa | Depende de | Esfuerzo estimado | Prioridad |
|---|---|---|---|
| 1. Dominio de Sesiones | — | Medio | 🔴 Crítica (base de todo) |
| 2. Login (Aplicación) | Etapa 1 | Alto | 🔴 Crítica |
| 3. Refresh Token | Etapa 1, 2 | Alto | 🔴 Crítica |
| 4. Logout | Etapa 2 | Medio | 🟡 Alta |
| 5. Seguridad Perimetral | Etapa 2 | Alto | 🟡 Alta |
| 6. Integración | Etapas 1-5 | Alto | 🟢 Media (después de las demás) |

## Notas Finales

- **No se incluye en esta spec**: Implementación concreta de handlers HTTP (presentación), eventos de dominio, auditoría de sesiones, historial de inicios de sesión, notificaciones de seguridad (email de "nuevo inicio de sesión"), OAuth2/OpenID Connect, MFA/2FA.
- **Decisiones abiertas para el equipo** (marcar antes de implementar):
  - ¿El rate limiting se implementa en la aplicación o en un reverse proxy (nginx, Cloudflare)?
  - ¿El almacenamiento de intentos por IP se hace en PostgreSQL o en Redis?
  - ¿Se implementa detección de robo de refresh token desde la etapa 3 o se difiere?
  - ¿La política de correo no verificado bloquea el login o solo advierte?
- **Evolución futura post-esta-spec**: Estas omisiones pueden ser abordadas en especificaciones posteriores sin romper el diseño actual, gracias a la separación por capas.
