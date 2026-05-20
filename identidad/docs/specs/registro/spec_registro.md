---
title: Especificación del Caso de Uso — Registro con Verificación de Correo Electrónico
version: 1.0
date_created: 2026-05-14
owner: Equipo Identidad
tags: registro, email-verification, usuario, credenciales
---

# Especificación del Caso de Uso — Registro con Verificación de Correo Electrónico

## 1. Propósito y Alcance

Definir la especificación completa del caso de uso "Registro de Usuario" con verificación de correo electrónico, siguiendo la arquitectura limpia del proyecto.

**Incluye:**
- Registro de nuevo usuario con creación atómica de usuario + credenciales
- Generación y envío del secreto de verificación por correo electrónico
- Confirmación de verificación mediante token
- Reenvío del secreto de verificación
- Expiración automática de secretos de verificación
- Estados de verificación y máquina de estados

**No incluye:**
- Login, refresh token, logout (ver `../sesiones/login_spec.md`)
- Seguridad perimetral (rate limiting, IP blocking)
- MFA/2FA, OAuth2/OpenID Connect
- Implementación concreta de handlers HTTP (presentación)
- Notificaciones push o SMS

## 2. Definiciones

| Término | Definición |
|---------|------------|
| **Secreto de verificación** | Código único (usualmente UUID) generado al registrar, enviado por email como token para confirmar la dirección de correo. En dominio se modela como `PruebaVerificacion`. |
| **Email personal como remitente** | Usamos una cuenta de Gmail personal del equipo (configurada vía SMTP) como remitente de los correos de verificación, hasta tener infraestructura de email corporativo |
| **Máquina de estados de verificación** | Modelo de estados por los que pasa el proceso de verificación (PENDIENTE_VERIFICACION → VERIFICADO / ENLACE_EXPIRADO / REENVIO_SOLICITADO). `VERIFICACION_FALLIDA` fue eliminado: un intento con secreto inválido no altera el estado. |
| **Eventos de dominio** | Eventos emitidos por la entidad Usuario cuando ocurren cambios de estado relevantes (UsuarioCreado, CorreoVerificado, etc.) |
| **SMTP directo** | Estrategia de envío de email sin usar servicios de terceros como SendGrid o Mailgun; el servidor se conecta directamente a un relay SMTP (Gmail) |

## 3. Requisitos, Restricciones y Guías

### Registro (ya implementado parcialmente)

- **REQ-REG-001**: El registro debe crear un usuario y sus credenciales en una sola transacción atómica (UnitOfWork).
- **REQ-REG-002**: El usuario se crea siempre en estado `NO_VERIFICADO` con verificación de correo en `PENDIENTE_VERIFICACION`.
- **REQ-REG-003**: El password se almacena como hash bcrypt, nunca en plano.
- **REQ-REG-004**: El correo electrónico se valida con formato estándar (RFC 5322) usando `net/mail`.
- **REQ-REG-005**: El email debe ser único en el sistema (restricción a nivel de BD + validación en aplicación).
- **REQ-REG-006**: Después de crear usuario + credenciales exitosamente, se debe disparar el flujo de verificación.
- **CON-REG-001**: El comando de registro valida campos ANTES de abrir transacción.
- **CON-REG-002**: Cualquier error dentro de la transacción causa rollback automático (no quedan usuarios huérfanos).

### Verificación de correo

- **REQ-VER-001**: Al completar el registro exitosamente, se genera un secreto de verificación único asociado al usuario (un UUID v4/v7 que funciona como token).
- **REQ-VER-002**: El secreto de verificación se almacena hasheado (SHA-256) en BD, nunca en plano.
- **REQ-VER-003**: Se envía un email al usuario con un enlace de verificación que incluye el token (el valor en plano).
- **REQ-VER-004**: El email de verificación se envía vía SMTP usando una cuenta de Gmail personal configurada vía variables de entorno (ver `Sistema de Email`).
- **REQ-VER-005**: El endpoint de confirmación recibe el token, lo valida (hasheándolo y comparando), y cambia el estado del usuario a `ACTIVO` y verificación a `VERIFICADO`.
- **REQ-VER-006**: Si el secreto ha expirado, se rechaza la verificación y el usuario puede solicitar reenvío.
- **REQ-VER-007**: Un usuario puede solicitar reenvío de verificación máximo N veces (configurable, default: 5) en una ventana de tiempo.
- **REQ-VER-008**: Un usuario con correo ya verificado no puede solicitar reenvío.
- **REQ-VER-009**: El secreto se genera y persiste DENTRO de la transacción de registro (no después).
- **REQ-SEC-001**: El secreto de verificación NUNCA debe exponerse en la respuesta HTTP del endpoint de registro ni en ninguna otra respuesta de la API. El token (valor en plano) solo se comunica al usuario a través del correo electrónico de verificación.

### Máquina de estados de verificación

- **REQ-VER-010**: La máquina de estados de verificación se limpia eliminando `VERIFICACION_FALLIDA`. Un token inválido (no encontrado en BD al hashearlo) retorna error 404/400 sin alterar el estado del dominio. Esto evita que un atacante pueda sabotear a un usuario legítimo forzando un cambio de estado con tokens basura.

- **REQ-VER-011**: Las transiciones permitidas son:
  ```
  PENDIENTE_VERIFICACION ───→ VERIFICADO           (confirmación exitosa)
  PENDIENTE_VERIFICACION ───→ ENLACE_EXPIRADO       (pasó tiempo de expiración)
  PENDIENTE_VERIFICACION ───→ REENVIO_SOLICITADO    (usuario pide reenvío)
  ENLACE_EXPIRADO        ───→ REENVIO_SOLICITADO    (usuario pide reenvío)
  REENVIO_SOLICITADO     ───→ VERIFICADO            (confirmación exitosa)
  REENVIO_SOLICITADO     ───→ ENLACE_EXPIRADO       (expiró nuevo token)
  VERIFICADO             ───→ (terminal)
  ```

Nota: El viejo REQ-VER-010 pasa a REQ-VER-011 para mantener el orden. Elimina las líneas de VERIFICACION_FALLIDA.

### Sistema de Email

- **REQ-EMAIL-001**: El envío de email usa SMTP directo; no hay servicio externo de email transaccional.
- **REQ-EMAIL-002**: Las credenciales SMTP se configuran vía variables de entorno: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`.
- **REQ-EMAIL-003**: Por defecto se usa Gmail SMTP (smtp.gmail.com:587) con TLS.
- **REQ-EMAIL-004**: El remitente (`FROM`) se configura con el email personal del equipo (`SMTP_FROM`).
- **REQ-EMAIL-005**: Los correos se envían de forma asíncrona (goroutine o cola en memoria) para no bloquear la respuesta HTTP.
- **REQ-EMAIL-006**: Si el envío de email falla, el registro NO debe fallar (consistencia eventual). El usuario queda creado pero sin verificar, puede solicitar reenvío manualmente.
- **CON-EMAIL-001**: No hay logs de emails enviados en BD en esta etapa (se considera mejora futura).
- **CON-EMAIL-002**: No hay plantillas de email complejas; el cuerpo se construye con texto plano + enlace.
- **CON-EMAIL-003**: La contraseña de Gmail debe ser una "Contraseña de Aplicación" (no la contraseña principal de Google).

## 4. Interfaces y Contratos de Datos

### Ubicación en el proyecto

```
internal/usuarios/
├── domain/
│   └── usuario/
│       ├── usuario.go                                    # (existente) + campo correo CorreoElectronico
│       ├── correo_electronico.go                         # NUEVO: Value Object embebido
│       ├── prueba_verificacion.go                        # NUEVO: Value Object embebido en CorreoElectronico
│       ├── estado_verificacion_correo.go                 # (existente) limpiar VERIFICACION_FALLIDA
│       ├── errors.go                                     # (existente) + nuevos errores
│       ├── repositorio.go                                # (existente) + ObtenerPorCorreo
│       ├── unit_of_work.go                               # (existente)
│       └── eventos.go                                    # (existente)
├── application/
│   └── services/
│       └── registro/
│           ├── comando.go                                # (existente)
│           ├── respuesta.go                              # (existente)
│           ├── servicio_registro.go                      # (existente) + token DENTRO de transacción
│           └── servicio_registro_test.go                 # (existente)
├── infrastructure/
    └── persistence/
        └── postgres/
            ├── usuario_repositorio.go                    # (existente) + nuevos métodos
            ├── usuario_model.go                           # (existente) + campos embebidos planos
            └── unit_of_work.go                           # (existente)

internal/notificaciones/                                   # NUEVO módulo
├── domain/
│   ├── email_servicio.go                                 # Interfaz EmailServicio
│   └── errores.go                                        # Errores de envío
└── infrastructure/
    └── email/
        └── smtp_email_servicio.go                         # Implementación SMTP con Gmail

internal/config/
└── env.go                                                # (existente) + nuevas vars SMTP
```

### Estructuras de datos

#### EstadoVerificacionCorreo (actualizado — sin VERIFICACION_FALLIDA)
```go
const (
    PENDIENTE_VERIFICACION EstadoVerificacionCorreo = "PENDIENTE_VERIFICACION"
    VERIFICADO             EstadoVerificacionCorreo = "VERIFICADO"
    ENLACE_EXPIRADO        EstadoVerificacionCorreo = "ENLACE_EXPIRADO"
    REENVIO_SOLICITADO     EstadoVerificacionCorreo = "REENVIO_SOLICITADO"
)

var TransicionesVerificacion = map[EstadoVerificacionCorreo][]EstadoVerificacionCorreo{
    PENDIENTE_VERIFICACION: {VERIFICADO, ENLACE_EXPIRADO, REENVIO_SOLICITADO},
    ENLACE_EXPIRADO:        {REENVIO_SOLICITADO},
    REENVIO_SOLICITADO:     {VERIFICADO, ENLACE_EXPIRADO},
    VERIFICADO:             {},
}
```

#### PruebaVerificacion (NUEVO - Value Object)
```go
// PruebaVerificacion encapsula el secreto que demuestra posesión del correo.
// El nombre "secretoHash" es del dominio. El algoritmo concreto (SHA-256) es
// un detalle de infraestructura resuelto en el mapeador de persistencia.
type PruebaVerificacion struct {
    secretoHash string    // hash del secreto (nunca el secreto en plano)
    expiraEn    time.Time // momento a partir del cual el secreto expira
}

func NuevaPruebaVerificacion(secretoHash string, expiraEn time.Time) (*PruebaVerificacion, error)
func (p *PruebaVerificacion) SecretoHash() string
func (p *PruebaVerificacion) ExpiraEn() time.Time
func (p *PruebaVerificacion) Expiro(ahora time.Time) bool
func (p *PruebaVerificacion) EstaPendiente() bool  // secretoHash != ""
```

#### CorreoElectronico (NUEVO - Value Object embebido en Usuario)
```go
type CorreoElectronico struct {
    direccion            string
    estado               EstadoVerificacionCorreo
    pruebaPosesion       PruebaVerificacion
    reenviosContador     int
    reenviosVentanaInicio time.Time
}

func NuevoCorreoElectronico(direccion string) (*CorreoElectronico, error)
func CorreoElectronicoDesdeBD(direccion string, estado EstadoVerificacionCorreo, prueba PruebaVerificacion, reenviosContador int, reenviosVentanaInicio time.Time) *CorreoElectronico

// Getters
func (c *CorreoElectronico) Direccion() string
func (c *CorreoElectronico) Estado() EstadoVerificacionCorreo
func (c *CorreoElectronico) PruebaPosesion() PruebaVerificacion
func (c *CorreoElectronico) ReenviosContador() int
func (c *CorreoElectronico) ReenviosVentanaInicio() time.Time

// Comportamiento
func (c *CorreoElectronico) EstaVerificado() bool
func (c *CorreoElectronico) EstaPendiente() bool
func (c *CorreoElectronico) AsignarPrueba(secretoHash string, expiraEn time.Time) error
func (c *CorreoElectronico) LimpiarPrueba()
func (c *CorreoElectronico) Verificar(ahora time.Time) error
func (c *CorreoElectronico) MarcarExpirado(ahora time.Time) error
func (c *CorreoElectronico) SolicitarReenvio(ahora time.Time, maxReenvios int, ventana time.Duration) error
func (c *CorreoElectronico) PuedeSolicitarReenvio(maxReenvios int, ventana time.Duration, ahora time.Time) bool
```

#### Usuario actualizado
```go
type Usuario struct {
    id                       string
    nombre                   string
    apellido                 string
    correo                   CorreoElectronico  // Value Object embebido
    telefono                 string
    estado                   EstadoUsuario
    fechaCreacion            time.Time
    fechaActualizacion       time.Time
    eventos                  *EventosUsuario
}

// Constructor actualizado
func NuevoUsuario(id, correo, nombre, apellido, telefono string) (*Usuario, error)

// Métodos delegados al VO
func (u *Usuario) Correo() CorreoElectronico
func (u *Usuario) EstaVerificado() bool              // u.correo.EstaVerificado()
func (u *Usuario) AsignarPruebaVerificacion(secretoHash string, expiraEn time.Time) error  // u.correo.AsignarPrueba(...)
func (u *Usuario) VerificarCorreo(ahora time.Time) error  // delega en correo.Verificar(ahora)
func (u *Usuario) SolicitarReenvioVerificacion(ahora time.Time, maxReenvios int, ventana time.Duration) error  // delega en correo.SolicitarReenvio(...)
```

#### EmailServicio (NUEVO - Interfaz en dominio)
```go
type EmailServicio interface {
    Enviar(ctx context.Context, para, asunto, cuerpo string) error
}
```

#### Extensiones a UsuarioRepositorio
```go
type UsuarioRepositorio interface {
    Crear(ctx context.Context, usuario *Usuario) (*Usuario, error)
    Actualizar(ctx context.Context, usuario *Usuario) (*Usuario, error)
    Eliminar(ctx context.Context, id string) error
    ObtenerPorID(ctx context.Context, id string) (*Usuario, error)
    ObtenerPorCorreo(ctx context.Context, correo string) (*Usuario, error)        // NUEVO
    ObtenerPorSecretoHash(ctx context.Context, hash string) (*Usuario, error)      // NUEVO (reemplaza ObtenerPorTokenVerificacionHash)
    Listar(ctx context.Context, especificacion EspecificacionUsuario, paginacion Paginacion) ([]*Usuario, error)
}
```

#### UsuarioModel (GORM) — con campos embebidos planos
```go
type UsuarioModel struct {
    ID                       string    `gorm:"type:uuid;primaryKey"`
    Nombre                   string    `gorm:"column:nombre"`
    Apellido                 string    `gorm:"column:apellido"`
    Correo                   string    `gorm:"column:correo;uniqueIndex"`
    CorreoEstado             string    `gorm:"column:correo_estado;default:PENDIENTE_VERIFICACION"`
    SecretoHash              string    `gorm:"column:secreto_hash"`
    SecretoExpiraEn          time.Time `gorm:"column:secreto_expira_en"`
    ReenviosContador         int       `gorm:"column:reenvios_contador;default:0"`
    ReenviosVentanaInicio    time.Time `gorm:"column:reenvios_ventana_inicio"`
    Telefono                 string    `gorm:"column:telefono"`
    Estado                   string    `gorm:"column:estado;default:NO_VERIFICADO"`
    FechaCreacion            time.Time `gorm:"column:fecha_creacion"`
    FechaActualizacion       time.Time `gorm:"column:fecha_actualizacion"`
}

func (UsuarioModel) TableName() string { return "usuarios" }
```

### Extensiones a UnitOfWork (existente)
```go
type UnitOfWork interface {
    // ... métodos existentes ...
    EmailServicio() notificaciones.EmailServicio  // NUEVO
}
```

### Variables de Entorno (NUEVAS)

| Variable | Default | Descripción |
|----------|---------|-------------|
| `SMTP_HOST` | `smtp.gmail.com` | Host del servidor SMTP |
| `SMTP_PORT` | `587` | Puerto SMTP (587 para TLS) |
| `SMTP_USER` | (requerido) | Email de la cuenta remitente (Gmail personal del equipo) |
| `SMTP_PASSWORD` | (requerido) | Contraseña de aplicación de Gmail |
| `SMTP_FROM` | mismo que SMTP_USER | Dirección FROM en los correos |
| `VERIFICACION_TOKEN_EXPIRACION` | `24h` | Duración del secreto de verificación |
| `VERIFICACION_MAX_REENVIOS` | `5` | Máximo de reenvíos permitidos |
| `VERIFICACION_VENTANA_REENVIOS` | `24h` | Ventana de tiempo para contar reenvíos |

### Contratos HTTP

#### POST /api/v1/auth/register (existente - sin cambios)
Request:
```json
{
    "nombre": "Juan",
    "apellido": "Pérez",
    "correo": "juan@correo.com",
    "password": "secreto123",
    "telefono": "0999999999"
}
```

Response 201:
```json
{
    "data": {
        "usuario_id": "01926b1e-...",
        "correo": "juan@correo.com",
        "estado": "NO_VERIFICADO"
    },
    "_links": {
        "self": { "href": "/api/v1/usuarios/{id}", "method": "GET" },
        "resend": { "href": "/api/v1/auth/verify/resend", "method": "POST" }
    }
}
```

#### NUEVO: POST /api/v1/auth/verify/resend
Request:
```json
{
    "correo": "juan@correo.com"
}
```

Response 200:
```json
{
    "data": {
        "message": "Nuevo enlace de verificación enviado al correo registrado",
        "email": "j***@correo.com"
    }
}
```

#### NUEVO: GET /api/v1/auth/verify?token={token}
Response 200:
```json
{
    "data": {
        "message": "Correo verificado exitosamente",
        "usuario_id": "01926b1e-...",
        "correo": "juan@correo.com"
    }
}
```

Response 410 (token expirado):
```json
{
    "title": "Token Expirado",
    "status": 410,
    "detail": "El enlace de verificación ha expirado. Solicite uno nuevo."
}
```

## 5. Flujo del Servicio de Registro (Actualizado)

### Flujo principal de registro
```
1. Validar comando (ANTES de transacción):
   - correo no vacío, formato válido, no duplicado
   - password no vacío
   - nombre no vacío

2. Transacción (UnitOfWork):
   a. Generar nuevo ID de usuario
   b. Crear entidad Usuario (NuevoUsuario → estado NO_VERIFICADO, verificación PENDIENTE_VERIFICACION)
   c. Asignar PruebaVerificacion al CorreoElectronico del usuario (secretoHash + expiraEn)
   d. Persistir usuario (en un solo INSERT: usuario + credenciales + hash del secreto + expiración)
   e. Hashear password con bcrypt
   f. Crear CredencialesUsuario
   g. Persistir credenciales
   h. Preparar respuesta DTO (sin token)
   i. → COMMIT

3. DESPUÉS de transacción exitosa (FUERA de la transacción — post-commit):
   a. Enviar email de verificación vía SMTP con el secreto en plano (asíncrono)
   b. Si el email falla: loguear error, no fallar el registro. El usuario tiene el secreto persistido y puede solicitar reenvío.

4. Retornar respuesta con usuarioID, correo, estado
```

### Flujo de verificación
```
1. Recibir token del enlace (GET /api/v1/auth/verify?token=xxx)
2. Hashear el secreto recibido (SHA-256)
3. Buscar usuario por secretoHash en BD
4. Si no se encuentra → error 404/400 "enlace inválido" (NO altera estado del dominio)
5. Si el secreto expiró (expiraEn < ahora):
   a. Marcar estadoVerificacionCorreo → ENLACE_EXPIRADO (transición desde estado actual)
   b. Limpiar prueba de verificación
   c. Error 410 "enlace expirado, solicite reenvío"
6. Si el secreto es válido:
   a. Cambiar estadoVerificacionCorreo → VERIFICADO
   b. Cambiar estado → ACTIVO (transición válida: NO_VERIFICADO → ACTIVO)
   c. Emitir evento CorreoVerificado
   d. Limpiar prueba de verificación
   e. Retornar éxito 200
```

### Flujo de reenvío
```
1. Validar comando: correo no vacío
2. Buscar usuario por correo
3. Si no existe → error genérico "si el correo está registrado, recibirá un enlace"
4. Si el correo ya está verificado (VERIFICADO) → error "el correo ya fue verificado"
5. Si excedió límite de reenvíos en ventana de tiempo
   → error "has solicitado demasiados reenvíos, intenta más tarde"
6. Si el estado actual permite reenvío (PENDIENTE_VERIFICACION, ENLACE_EXPIRADO):
   a. Cambiar estadoVerificacionCorreo → REENVIO_SOLICITADO
   b. Generar nuevo secreto de verificación
   c. Asignar PruebaVerificacion al CorreoElectronico (nuevo secretoHash + expiraEn)
   d. Persistir cambios en usuario
   e. Enviar nuevo email de verificación con el secreto en plano (asíncrono)
   f. Retornar éxito
```

## 6. Escenarios de TDD

### Registro

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 1 | Registro exitoso | Datos válidos | `ServicioRegistro.Ejecutar(cmd)` | Usuario creado en NO_VERIFICADO, credenciales creadas, CorreoElectronico con PENDIENTE_VERIFICACION + PruebaVerificacion asignada, email encolado |
| 2 | Registro con correo duplicado | Correo ya existe | validación previa | Error "correo ya registrado", sin transacción |
| 3 | Email vacío | correo="" | validación | Error |
| 4 | Email mal formado | "invalido" | validación | Error |
| 5 | Password vacío | password="" | validación | Error |
| 6 | Nombre vacío | nombre="" | validación | Error |
| 7 | Rollback si falla persistencia | Usuario OK, credenciales falla | `ServicioRegistro.Ejecutar(cmd)` | Usuario no existe en BD (rollback) |
| 8 | Context timeout | Context cancelado | `ServicioRegistro.Ejecutar(cmd)` | Rollback completo |
| 9 | PruebaVerificacion persistida en transacción | Registro exitoso | Dentro de transacción | Usuario tiene secretoHash + expiraEn en BD |
| 10 | Email enviado post-commit | Commit exitoso | Post-commit | EmailServicio.Enviar() llamado con secreto en plano |
| 11 | Error de email no bloquea | Commit exitoso, email falla | Post-commit | Usuario creado con PruebaVerificacion persistida, sin error al usuario |

### Verificación de correo

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 12 | Verificación exitosa | PruebaVerificacion válida, no expirada | `Verificar(secreto)` | Estado → ACTIVO, verificación → VERIFICADO, PruebaVerificacion limpiada, evento emitido |
| 13 | Secreto inválido (no existe en BD) | Hash del secreto no coincide con ningún registro | `Verificar(secreto)` | Error 404 "enlace inválido". Estado del dominio NO se modifica |
| 14 | Secreto expirado | PruebaVerificacion con expiraEn < ahora | `Verificar(secreto)` | Verificación → ENLACE_EXPIRADO, error 410 |
| 15 | Doble verificación | Usuario ya VERIFICADO | `Verificar(secreto)` | Error "correo ya verificado" |
| 16 | Secreto mal formado | String que no es UUID | `Verificar(secreto)` | Error 400 "enlace inválido" (validación temprana) |

### Reenvío de verificación

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 17 | Reenvío exitoso | Usuario PENDIENTE_VERIFICACION | `SolicitarReenvio(correo)` | Nuevo secreto, reenviosContador+1, reenviosVentanaInicio actualizada, email enviado |
| 18 | Reenvío después de expiración | Estado ENLACE_EXPIRADO | `SolicitarReenvio(correo)` | Nuevo secreto, reenviosContador+1, email enviado |
| 19 | Reenvío cuando ya verificado | Estado VERIFICADO | `SolicitarReenvio(correo)` | Error "correo ya verificado" |
| 20 | Reenvío con correo inexistente | Correo no registrado | `SolicitarReenvio(correo)` | Error genérico "si está registrado..." |
| 21 | Excede límite de reenvíos | reenviosContador=5, dentro de ventana 24h | `SolicitarReenvio(correo)` | Error "demasiados intentos, intente más tarde" |
| 22 | Reenvío después de ventana | reenviosContador=5, ventana expiró | `SolicitarReenvio(correo)` | Permitido, reenviosContador reseteado a 1, nueva ventana |
| 23 | Reenvío con secretos anteriores rotados | Reenvío exitoso con viejo secreto expirado | Nuevo secreto generado | Viejo secreto ya no es válido (reemplazado en BD) |

### PruebaVerificacion (Value Object)

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 24 | Creación de prueba | secretoHash + expiraEn válidos | `NuevaPruebaVerificacion(hash, expira)` | VO creado, EstaPendiente() = true |
| 25 | Prueba expirada | expiraEn en el pasado | `prueba.Expiro(ahora)` | true |
| 26 | Prueba vigente | expiraEn en el futuro | `prueba.Expiro(ahora)` | false |
| 27 | Prueba vacía (sin secretos) | secretoHash = "" | `prueba.EstaPendiente()` | false |

### CorreoElectronico (Value Object)

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 28 | Creación con dirección válida | email válido | `NuevoCorreoElectronico(email)` | VO creado, estado PENDIENTE_VERIFICACION, reenviosContador=0 |
| 29 | Dirección vacía | "" | Constructor | Error del dominio |
| 30 | Asignar prueba | PruebaVerificacion válida | `correo.AsignarPrueba(hash, expira)` | PruebaPosesion actualizada |
| 31 | Verificación exitosa | Prueba válida y no expirada | `correo.Verificar(ahora)` | Estado → VERIFICADO, prueba limpiada |
| 32 | Verificación con prueba expirada | Prueba expirada | `correo.Verificar(ahora)` | Estado → ENLACE_EXPIRADO, error |
| 33 | Verificación cuando ya VERIFICADO | Estado ya terminal | `correo.Verificar(ahora)` | Error |
| 34 | Solicitar reenvío dentro de límite | reenviosContador=3, max=5 | `correo.SolicitarReenvio(ahora, 5, 24h)` | reenviosContador=4, nuevo secreto asignado |
| 35 | Solicitar reenvío excede límite | reenviosContador=5, dentro de ventana | `correo.SolicitarReenvio(ahora, 5, 24h)` | Error, estado no cambia |
| 36 | Solicitar reenvío después de ventana | reenviosContador=5, ventana expiró | `correo.SolicitarReenvio(ahora, 5, 24h)` | Permitido, contador reseteado a 1 |

## 7. Actividades de Implementación

### Fase 1: Value Object PruebaVerificacion y CorreoElectronico
1. Crear `internal/usuarios/domain/usuario/prueba_verificacion.go`
2. Implementar generación de secreto UUID v4
3. Implementar hashing SHA-256
4. Implementar validación de expiración
5. Crear `internal/usuarios/domain/usuario/correo_electronico.go` con Value Object
6. Implementar métodos de comportamiento: Verificar, MarcarExpirado, SolicitarReenvio, AsignarPrueba
7. Escribir tests de ambos value objects

### Fase 2: Actualizar Usuario con CorreoElectronico embebido
1. Reemplazar campos planos `tokenVerificacionHash`/`tokenVerificacionExpira` por `correo CorreoElectronico`
2. Agregar métodos delegados: `AsignarPruebaVerificacion()`, `VerificarCorreo()`, `SolicitarReenvioVerificacion()`
3. Actualizar `NuevoUsuario` y constructores de persistencia
4. Escribir tests de unidad

### Fase 3: Crear módulo de notificaciones
1. Crear `internal/notificaciones/domain/email_servicio.go` con interfaz
2. Crear `internal/notificaciones/infrastructure/email/smtp_email_servicio.go` con implementación SMTP
3. Configurar conexión SMTP con Gmail (TLS, puerto 587)
4. Implementar construcción de cuerpo de email (texto plano con enlace)
5. Escribir tests de integración SMTP (usando servidor SMTP de prueba o mock)

### Fase 4: Extender el registro
1. Agregar validación de correo duplicado en `validarComando()` (antes de transacción: consultar si existe)
2. Generar secreto de verificación DENTRO de la transacción y asignarlo al CorreoElectronico
3. Enviar email post-commit con el secreto en plano
4. Agregar `EmailServicio` a `UnitOfWork` (o inyectarlo directamente en `ServicioRegistro`)
5. Escribir tests de integración

### Fase 5: Endpoint de verificación
1. Crear `internal/usuarios/application/services/verificacion/` con comando, respuesta y servicio
2. Implementar flujo de confirmación de verificación
3. Agregar handler HTTP: GET /api/v1/auth/verify?token=xxx
4. Agregar facade o usar AuthFacade existente
5. Escribir tests de integración

### Fase 6: Endpoint de reenvío
1. Implementar flujo de reenvío en el servicio de verificación
2. Agregar handler HTTP: POST /api/v1/auth/verify/resend
3. Implementar control de límite de reenvíos en CorreoElectronico.SolicitarReenvio()
4. Escribir tests de integración

### Fase 7: Integración y configuración
1. Agregar variables SMTP y de verificación a `config/env.go`
2. Agregar `smtp_email_servicio` al `registry`
3. Agregar migración para nuevos campos en tabla `usuarios` (secreto_hash, secreto_expira_en, reenvios_contador, reenvios_ventana_inicio, correo_estado)
4. Verificar que todos los tests pasen

## 8. Criterios de Aceptación

- **AC-REG-001**: Dado un registro exitoso, Cuando se completa la transacción, Entonces el usuario existe en BD en estado NO_VERIFICADO con PruebaVerificacion asignada (secretoHash + expiraEn).
- **AC-REG-002**: Dado un registro exitoso, Cuando se completa, Entonces se intenta enviar un email de verificación (aunque falle, el registro no falla y el secreto ya está persistido).
- **AC-REG-003**: Dado un email duplicado, Cuando se intenta registrar, Entonces se retorna error sin crear nada.
- **AC-VER-001**: Dado un secreto de verificación válido, Cuando se confirma, Entonces el usuario pasa a estado ACTIVO y verificación VERIFICADO.
- **AC-VER-002**: Dado un secreto expirado, Cuando se confirma, Entonces se retorna error 410 y el estado de verificación pasa a ENLACE_EXPIRADO.
- **AC-VER-003**: Dado un correo no verificado, Cuando se solicita reenvío, Entonces se genera un nuevo secreto y se envía un nuevo email.
- **AC-VER-004**: Dado un correo ya verificado, Cuando se solicita reenvío, Entonces se retorna error informando que ya está verificado.
- **AC-EMAIL-001**: Dado que SMTP está configurado, Cuando se envía un email, Entonces el mensaje llega al destinatario con el enlace de verificación.
- **AC-EMAIL-002**: Dado que SMTP falla (credenciales inválidas, red caída), Cuando se envía un email, Entonces se loguea el error y el registro continúa sin afectar al usuario.

## 9. Estrategia de Automatización de Pruebas

- **Niveles**: Unitarias (dominio), Integración (BD real + mock SMTP), End-to-End (servidor HTTP)
- **Framework**: `testing` estándar + `testify`
- **Mocks**: `MockEmailServicio` implementando `EmailServicio` interfaz para pruebas
- **BD de prueba**: PostgreSQL con `setupTestDB()` (estilo existente)
- **SMTP de prueba**: Usar servidor SMTP mock (ej: MailHog, Papercut, o servidor SMTP de prueba embebido) para pruebas de integración
- **Cobertura mínima**: 90% en dominio, 80% en aplicación
- **Pruebas de envío real**: Solo en entorno de staging con Gmail real configurado

## 10. Justificación y Contexto

### ¿Por qué email de verificación?
- Cumplimiento de buenas prácticas de seguridad: validar que el usuario tiene acceso al correo registrado
- Prevención de registros con correos falsos o temporales
- Base para flujos futuros: recuperación de contraseña, notificaciones

### ¿Por qué SMTP directo con Gmail personal?
- **Cero costo**: Mientras el proyecto no tenga ingresos ni infraestructura de email de negocio
- **Simplicidad**: Gmail SMTP es gratuito, confiable y bien soportado
- **Limitación conocida**: 500 emails/día en cuentas gratuitas de Gmail
- **Migración futura**: La interfaz `EmailServicio` permite cambiar a SendGrid, Mailgun, AWS SES, etc. sin tocar el dominio

### ¿Por qué el email se envía fuera de la transacción?
- El envío de email es una operación lenta y propensa a fallos
- Si el email falla, no debe revertir el registro del usuario (el usuario puede solicitar reenvío)
- Diseño consistente con el principio de "consistencia eventual" para operaciones de notificación

### ¿Por qué el secreto se persiste en la tabla de usuarios?
- Simplicidad: no requiere tabla adicional ni nuevas migraciones complejas
- Un usuario tiene exactamente un secreto de verificación activo a la vez
- Si en el futuro se requiere soporte multi-secreto, se migra a tabla separada

## 11. Dependencias e Integraciones Externas

### Dependencias del proyecto (ya existen)
- `net/smtp` (stdlib) — para el envío SMTP
- `crypto/sha256` (stdlib) — para hashear secretos
- `github.com/google/uuid` (o `github.com/davosjar/...`) — para generar UUIDs
- GORM — para persistencia

### Dependencias de infraestructura
- Cuenta de Gmail personal del equipo para SMTP
- Contraseña de aplicación de Gmail (no la contraseña principal)
- Servidor SMTP mock para pruebas (MailHog o similar)

### Variables de entorno requeridas
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=equipo@email.com
SMTP_PASSWORD=xxxx xxxx xxxx xxxx  (contraseña de aplicación)
SMTP_FROM=Identidad <equipo@email.com>
VERIFICACION_TOKEN_EXPIRACION=24h
VERIFICACION_MAX_REENVIOS=5
VERIFICACION_VENTANA_REENVIOS=24h
```

### Dependencias de configuración
- `.env` local de desarrollo con credenciales SMTP
- Docker Compose con MailHog para pruebas locales
- CI sin SMTP real (mock), staging con SMTP real

## 12. Ejemplos y Casos Borde

### Registro exitoso con verificación
```go
comando := &ComandoRegistro{
    Correo:   "nuevo@usuario.com",
    Password: "MiPassword123!",
    Nombre:   "María",
    Apellido: "González",
    Telefono: "0999999999",
}

respuesta, err := servicio.Ejecutar(ctx, comando)
// respuesta.UsuarioID != ""
// respuesta.Estado == "NO_VERIFICADO"
// Email enviado a nuevo@usuario.com con enlace de verificación
```

### Verificación con secreto expirado
```go
// Usuario recibe email 48h después del registro (secreto expira en 24h)
token := "01926b1e-aaaa-bbbb-cccc-000000000001"
err := servicioVerificacion.Confirmar(ctx, token)
// err.Error() contiene "enlace expirado"
// usuario.Correo().Estado() == ENLACE_EXPIRADO
```

### Reenvío excede límite
```go
// Usuario solicita reenvío 6 veces en 12h (límite: 5 en 24h)
err := servicioVerificacion.Reenviar(ctx, "nuevo@usuario.com")
// err.Error() contiene "demasiados intentos"
// No se envía email
```

## 13. Validación de Cumplimiento

- [ ] ¿CorreoElectronico es un Value Object inmutable (o con comportamiento controlado)?
- [ ] ¿PruebaVerificacion encapsula el secreto de posesión del correo?
- [ ] ¿El secreto se almacena como hash SHA-256 (nunca en plano)?
- [ ] ¿El secreto se genera y persiste DENTRO de la transacción de registro?
- [ ] ¿El envío del email queda FUERA de la transacción (post-commit)?
- [ ] ¿El fallo de email no causa rollback del registro?
- [ ] ¿VERIFICACION_FALLIDA fue eliminado de la máquina de estados?
- [ ] ¿Un intento con secreto inválido NO altera el estado del dominio?
- [ ] ¿reenviosContador y reenviosVentanaInicio están en CorreoElectronico?
- [ ] ¿Los métodos de Usuario delegan en CorreoElectronico (no hay lógica duplicada)?
- [ ] ¿La interfaz EmailServicio está en el dominio de notificaciones?
- [ ] ¿La implementación SMTP está en infraestructura?
- [ ] ¿Los límites de reenvío son configurables vía entorno?
- [ ] ¿La expiración del secreto es configurable vía entorno?
- [ ] ¿Los errores de verificación no revelan información sensible?
- [ ] ¿El correo verificado es estado terminal en la máquina de estados?
- [ ] ¿Hay tests para cada escenario de TDD listado?
- [ ] ¿Las variables SMTP se leen de entorno, no están hardcodeadas?
- [ ] ¿El secreto de verificación NUNCA se expone en respuestas HTTP de la API?

## 14. Especificaciones Relacionadas

- `../sesiones/login_spec.md` — Login, refresh, logout, seguridad perimetral e integración JWT
- `../presentacion/spec-presentation-layer.md` — Capa de presentación con Gin + Huma + OpenAPI
- `../../adr/architecture-context.md` — Contexto arquitectónico y flujo de capas
- `../../adr/feature-template.md` — Template para nuevas features
