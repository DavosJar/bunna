## 1. ENDPOINTS DEL MÓDULO IDENTIDAD

| Método | Ruta | Descripción | Request Body | Response Body | Requiere JWT |
|--------|------|-------------|--------------|---------------|--------------|
| GET | /health | Verifica que el servicio está activo y respondiendo. | — | `{ status: string }` | No |
| POST | /api/v1/auth/register | Crea un nuevo usuario con sus credenciales. Devuelve el ID y estado inicial. | `RegisterRequest` { nombre, apellido, correo, password, telefono? } | `ApiResponse[RegisterResponse]` { usuario_id, correo, estado } | No |
| POST | /api/v1/auth/login | Autentica al usuario y devuelve tokens JWT de acceso y refresco. | `LoginRequest` { correo, password } | `ApiResponse[LoginResponse]` { access_token, refresh_token, expires_in, token_type, usuario_id, tenant_id, rol } | No |
| POST | /api/v1/auth/refresh | Renueva el access token usando el refresh token. Aplica rotación de tokens. | `RefreshRequest` { refresh_token } | `ApiResponse[RefreshResponse]` { access_token, refresh_token, expires_in, token_type, usuario_id } | No |
| POST | /api/v1/auth/logout | Cierra la sesión actual del usuario autenticado. | — | `ApiResponse[LogoutResponse]` { sesiones_revocadas } | Sí |
| POST | /api/v1/auth/logout/all | Cierra todas las sesiones activas del usuario autenticado. | — | `ApiResponse[LogoutResponse]` { sesiones_revocadas } | Sí |
| POST | /api/v1/usuarios | Crea un nuevo usuario en el sistema (requiere permisos de administración). | `CrearUsuarioRequest` { correo, nombre, apellido, password } | `ApiResponse[CrearUsuarioResponse]` { id, correo, nombre, apellido, activo, creado_en } | Sí |
| GET | /api/v1/usuarios | Lista usuarios del sistema con filtros y paginación. | Query: pagina, tamano, correo, estado | `ApiResponse[ListarUsuariosResponse]` { usuarios[], total, pagina, tamano } | Sí |
| PUT | /api/v1/usuarios/{usuarioID} | Modifica los datos de un usuario existente. | `ModificarUsuarioRequest` { nombre, apellido } | `ApiResponse[ModificarUsuarioResponse]` { id, correo, nombre, apellido, modificado_en } | Sí |
| DELETE | /api/v1/usuarios/{usuarioID} | Desactiva un usuario del sistema (baja lógica). | `DarDeBajaUsuarioRequest` { motivo? } | `ApiResponse[DarDeBajaUsuarioResponse]` { usuario_id, estado, baja_en } | Sí |
| POST | /api/v1/usuarios/{usuarioID}/expulsar | Expulsa a un usuario del sistema, desactivándolo e invalidando todas sus sesiones. | — | `ApiResponse[ExpulsarUsuarioResponse]` { usuario_id, estado, sesiones_revocadas, expulsado_en } | Sí |
| GET | /api/v1/mi-perfil | Obtiene los datos del perfil del usuario autenticado. | — | `ApiResponse[VerMiPerfilResponse]` { id, correo, nombre, apellido, telefono, estado, creado_en } | Sí |
| PUT | /api/v1/mi-perfil | Actualiza los datos del perfil del usuario autenticado. | `ModificarMiPerfilRequest` { nombre, apellido } | `ApiResponse[ModificarMiPerfilResponse]` { id, correo, nombre, apellido, modificado_en } | Sí |
| PUT | /api/v1/mi-password | Cambia la contraseña del usuario autenticado. | `CambiarMiPasswordRequest` { password_actual, nueva_password } | `ApiResponse[CambiarMiPasswordResponse]` { modificado_en } | Sí |
| POST | /api/v1/usuarios/{usuarioID}/reset-password | Resetea la contraseña de un usuario (requiere permisos administrativos). | `ResetearPasswordRequest` { nueva_password } | `ApiResponse[ResetearPasswordResponse]` { usuario_id, modificado_en } | Sí |
| POST | /api/v1/usuarios/{usuarioID}/unlock | Desbloquea la cuenta de un usuario bloqueada por intentos fallidos. | — | `ApiResponse[DesbloquearCuentaResponse]` { usuario_id, desbloqueado_en } | Sí |
| GET | /api/v1/ips-bloqueadas | Lista las direcciones IP bloqueadas temporalmente por exceso de intentos. | Query: pagina, tamano | `ApiResponse[ListarIPsBloqueadasResponse]` { ips[], total, pagina } | Sí |
| DELETE | /api/v1/ips-bloqueadas/{ip} | Elimina el bloqueo de una dirección IP. | — | `ApiResponse[DesbloquearIPResponse]` { ip, desbloqueado_en } | Sí |
| GET | /api/v1/credenciales/{usuarioID} | Obtiene el estado de las credenciales de un usuario (bloqueo, intentos, verificación). | — | `ApiResponse[ConsultarCredencialesResponse]` { usuario_id, activo, correo_verificado, intentos_fallidos, bloqueado_hasta } | Sí |
| GET | /api/v1/sesiones | Lista las sesiones activas del sistema con paginación. | Query: pagina, tamano | `ApiResponse[ListarSesionesResponse]` { sesiones[], total, pagina } | Sí |
| DELETE | /api/v1/sesiones/{sesionID} | Fuerza el cierre de una sesión específica (requiere permisos administrativos). | — | `ApiResponse[ForzarCierreSesionResponse]` { sesion_id, estado, revocado_en } | Sí |
| GET | /api/v1/roles | Lista los roles del sistema con paginación. | Query: pagina, tamano | `ApiResponse[ListarRolesResponse]` { roles[], total, pagina } | Sí |
| POST | /api/v1/roles | Crea un nuevo rol en el sistema con permisos opcionales. | `CrearRolRequest` { nombre, descripcion, permisos[] } | `ApiResponse[CrearRolResponse]` { id, nombre, descripcion, es_sistema, creado_en } | Sí |
| PUT | /api/v1/roles/{rolID} | Actualiza el nombre y descripción de un rol. | `ModificarRolRequest` { nombre, descripcion } | `ApiResponse[ModificarRolResponse]` { id, nombre, descripcion, modificado_en } | Sí |
| DELETE | /api/v1/roles/{rolID} | Elimina un rol del sistema (no se pueden eliminar roles de sistema). | — | `ApiResponse[EliminarRolResponse]` { rol_id, eliminado_en } | Sí |
| POST | /api/v1/usuarios/{usuarioID}/roles | Asigna un rol a un usuario, opcionalmente en un tenant específico. | `AsignarRolRequest` { rol_id, tenant_id } | `ApiResponse[AsignarRolResponse]` { usuario_id, rol_id, tenant_id, asignado_en } | Sí |
| DELETE | /api/v1/usuarios/{usuarioID}/roles/{rolID} | Revoca un rol asignado a un usuario. | — | `ApiResponse[RevocarRolResponse]` { usuario_id, rol_id, tenant_id, revocado_en } | Sí |
| POST | /api/v1/roles/{rolID}/permisos | Asigna un permiso a un rol específico. | `AsignarPermisoRequest` { permiso_codigo } | `ApiResponse[AsignarPermisoResponse]` { rol_id, permiso_codigo, asignado_en } | Sí |
| DELETE | /api/v1/roles/{rolID}/permisos/{codigo} | Revoca un permiso previamente asignado a un rol. | — | `ApiResponse[RevocarPermisoResponse]` { rol_id, permiso_codigo, revocado_en } | Sí |
| GET | /api/v1/permisos | Lista todos los permisos disponibles en el sistema. | — | `ApiResponse[ListarPermisosResponse]` { permisos[], total } | Sí |
| PUT | /api/v1/tenants/{tenantID} | Actualiza la configuración de un tenant. | `ConfigurarTenantRequest` { nombre, slug } | `ApiResponse[ConfigurarTenantResponse]` { tenant_id, nombre, slug, modificado_en } | Sí |
| POST | /api/v1/verificacion/solicitar | Envía un enlace de verificación al correo del usuario autenticado. | — | `ApiResponse[SolicitarVerificacionResponse]` { mensaje } | Sí |
| POST | /api/v1/verificacion/confirmar | Confirma la verificación del correo electrónico usando el token recibido. | `ConfirmarVerificacionRequest` { token } | `ApiResponse[ConfirmarVerificacionResponse]` { mensaje } | No |
| POST | /api/v1/verificacion/reenviar | Reenvía el enlace de verificación al correo del usuario autenticado. | — | `ApiResponse[ReenviarVerificacionResponse]` { mensaje } | Sí |
| POST | /api/v1/recuperacion/solicitar | Envía un enlace de recuperación al correo electrónico proporcionado. | `SolicitarRecuperacionRequest` { correo } | `ApiResponse[SolicitarRecuperacionResponse]` { mensaje } | No |
| POST | /api/v1/recuperacion/validar | Valida si un token de recuperación es válido y devuelve el ID del usuario asociado. | `ValidarTokenRecuperacionRequest` { token } | `ApiResponse[ValidarTokenRecuperacionResponse]` { usuario_id, valido } | No |
| POST | /api/v1/recuperacion/confirmar | Restablece la contraseña usando el token de recuperación. | `ConfirmarRecuperacionRequest` { token, nueva_password } | `ApiResponse[ConfirmarRecuperacionResponse]` { mensaje } | No |

## 2. ESTRUCTURA DE CARPETAS

```
identidad/
├── cmd
│   ├── main.go
│   ├── seed
│   │   └── main.go
│   └── test
│       └── main.go
├── internal
│   ├── config
│   │   ├── database.go
│   │   └── env.go
│   ├── handler
│   │   └── handler.go
│   ├── notificaciones
│   │   ├── domain
│   │   │   ├── email_servicio.go
│   │   │   ├── errores.go
│   │   │   ├── templates.go
│   │   │   ├── templates_test.go
│   │   │   └── tipos_email.go
│   │   └── infrastructure
│   │       └── email
│   │           └── smtp_servicio.go
│   ├── presentation
│   │   ├── dto
│   │   │   ├── login_dto.go
│   │   │   ├── rbac_dto.go
│   │   │   ├── recuperacion_dto.go
│   │   │   ├── register_dto.go
│   │   │   ├── seguridad_dto.go
│   │   │   ├── sesion_dto.go
│   │   │   ├── tenant_dto.go
│   │   │   ├── usuario_dto.go
│   │   │   └── verificacion_dto.go
│   │   ├── facades
│   │   │   ├── all_facades.go
│   │   │   ├── auth_facade.go
│   │   │   ├── auth_facade_impl.go
│   │   │   ├── auth_facade_test.go
│   │   │   ├── rbac_facade.go
│   │   │   ├── recuperacion_facade.go
│   │   │   ├── seguridad_facade.go
│   │   │   ├── sesion_facade.go
│   │   │   ├── tenant_facade.go
│   │   │   ├── usuario_facade.go
│   │   │   └── verificacion_facade.go
│   │   ├── handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── handlers_test.go
│   │   │   ├── health_handler.go
│   │   │   ├── login_handler.go
│   │   │   ├── rbac_handler.go
│   │   │   ├── recuperacion_handler.go
│   │   │   ├── register_handler.go
│   │   │   ├── seguridad_handler.go
│   │   │   ├── sesion_handler.go
│   │   │   ├── tenant_handler.go
│   │   │   ├── usuario_handler.go
│   │   │   └── verificacion_handler.go
│   │   ├── middleware
│   │   │   ├── jwt_middleware.go
│   │   │   └── jwt_middleware_test.go
│   │   └── router
│   │       └── router.go
│   ├── rbac
│   │   ├── application
│   │   │   ├── seed_servicio.go
│   │   │   ├── services
│   │   │   │   └── autorizacion
│   │   │   │       └── autorizacion_servicio.go
│   │   │   └── usecases
│   │   │       ├── assignpermissiontorole
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── assignrole
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── checkpermission
│   │   │       │   └── usecase.go
│   │   │       ├── createrole
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── deleterole
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── listpermisos
│   │   │       │   └── usecase.go
│   │   │       ├── listroles
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── revokepermissionfromrole
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── revokerole
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       └── updaterole
│   │   │           ├── command.go
│   │   │           ├── response.go
│   │   │           └── usecase.go
│   │   ├── domain
│   │   │   ├── errors.go
│   │   │   ├── especificacion_rol.go
│   │   │   ├── permisos.go
│   │   │   ├── permisos_test.go
│   │   │   ├── permiso_utils.go
│   │   │   ├── repositorios.go
│   │   │   ├── roles.go
│   │   │   └── roles_test.go
│   │   └── infrastructure
│   │       └── persistence
│   │           └── postgres
│   │               ├── rbac_models.go
│   │               └── rbac_repositorios.go
│   ├── recuperacion
│   │   ├── application
│   │   │   ├── services
│   │   │   │   └── recuperacion
│   │   │   │       ├── comando.go
│   │   │   │       ├── respuesta.go
│   │   │   │       └── servicio_recuperacion.go
│   │   │   └── usecases
│   │   │       └── forgotpassword
│   │   │           ├── command.go
│   │   │           ├── response.go
│   │   │           └── usecase.go
│   │   ├── domain
│   │   │   ├── errores.go
│   │   │   ├── repositorio.go
│   │   │   ├── token_recuperacion.go
│   │   │   └── token_recuperacion_test.go
│   │   └── infrastructure
│   │       └── persistence
│   │           └── postgres
│   │               ├── token_recuperacion_model.go
│   │               ├── token_recuperacion_repositorio.go
│   │               └── usuario_recuperacion_repositorio.go
│   ├── registry
│   │   └── registry.go
│   ├── seguridad
│   │   ├── application
│   │   │   ├── services
│   │   │   │   ├── bloqueo_ip
│   │   │   │   │   ├── servicio_bloqueo_ip.go
│   │   │   │   │   └── servicio_bloqueo_ip_test.go
│   │   │   │   └── rate_limiter
│   │   │   │       ├── servicio_rate_limiter.go
│   │   │   │       └── servicio_rate_limiter_test.go
│   │   │   └── usecases
│   │   │       ├── changemypassword
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── listblockedips
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── resetpassword
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── unblockip
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── unlockaccount
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       └── viewcredentials
│   │   │           ├── command.go
│   │   │           ├── response.go
│   │   │           └── usecase.go
│   │   ├── domain
│   │   │   ├── credenciales.go
│   │   │   ├── credenciales_repositorio.go
│   │   │   ├── credenciales_test.go
│   │   │   ├── encriptacion.go
│   │   │   ├── especificacion_credenciales.go
│   │   │   ├── especificacion_intento_ip.go
│   │   │   ├── intento_ip.go
│   │   │   └── intento_ip_repositorio.go
│   │   └── infrastructure
│   │       ├── persistence
│   │       │   └── postgres
│   │       │       ├── credenciales_model.go
│   │       │       ├── credenciales_model_test.go
│   │       │       ├── credenciales_repositorio.go
│   │       │       ├── credenciales_repositorio_test.go
│   │       │       ├── intento_ip_model.go
│   │       │       ├── intento_ip_repositorio.go
│   │       │       ├── rate_limit_model.go
│   │       │       └── rate_limit_repositorio.go
│   │       └── security
│   │           └── bcrypt
│   │               ├── encriptacion.go
│   │               └── encriptacion_test.go
│   ├── sesiones
│   │   ├── application
│   │   │   ├── services
│   │   │   │   ├── login
│   │   │   │   │   ├── comando.go
│   │   │   │   │   ├── ejecutor.go
│   │   │   │   │   ├── respuesta.go
│   │   │   │   │   ├── servicio_login.go
│   │   │   │   │   └── servicio_login_test.go
│   │   │   │   ├── logout
│   │   │   │   │   ├── comando.go
│   │   │   │   │   ├── respuesta.go
│   │   │   │   │   ├── servicio_logout.go
│   │   │   │   │   └── servicio_logout_test.go
│   │   │   │   └── refresh
│   │   │   │       ├── comando.go
│   │   │   │       ├── respuesta.go
│   │   │   │       ├── servicio_refresh.go
│   │   │   │       └── servicio_refresh_test.go
│   │   │   └── usecases
│   │   │       ├── listsessions
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── login
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   ├── usecase.go
│   │   │       │   └── usecase_test.go
│   │   │       ├── logout
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   ├── usecase.go
│   │   │       │   └── usecase_test.go
│   │   │       ├── refresh
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   ├── usecase.go
│   │   │       │   └── usecase_test.go
│   │   │       └── terminatesession
│   │   │           ├── command.go
│   │   │           ├── response.go
│   │   │           └── usecase.go
│   │   ├── domain
│   │   │   ├── errores.go
│   │   │   ├── especificacion_sesion.go
│   │   │   ├── sesion.go
│   │   │   ├── sesion_repositorio.go
│   │   │   ├── sesion_test.go
│   │   │   ├── token_servicio.go
│   │   │   ├── tokens.go
│   │   │   └── unit_of_work.go
│   │   └── infrastructure
│   │       ├── persistence
│   │       │   └── postgres
│   │       │       ├── sesion_model.go
│   │       │       ├── sesion_repositorio.go
│   │       │       └── unit_of_work.go
│   │       └── security
│   │           └── jwt
│   │               └── jwt_token_servicio.go
│   ├── shared
│   │   ├── domain
│   │   │   ├── idgenerator.go
│   │   │   └── specification.go
│   │   └── infrastructure
│   │       └── idgenerator
│   │           ├── uuid_v7.go
│   │           └── uuid_v7_test.go
│   ├── tenants
│   │   ├── application
│   │   │   ├── services
│   │   │   │   └── gestionar_tenant
│   │   │   │       ├── comando.go
│   │   │   │       ├── respuesta.go
│   │   │   │       └── servicio_tenant.go
│   │   │   └── usecases
│   │   │       └── updatetenant
│   │   │           ├── command.go
│   │   │           ├── response.go
│   │   │           └── usecase.go
│   │   ├── domain
│   │   │   └── tenant
│   │   │       ├── errors.go
│   │   │       ├── membresia.go
│   │   │       ├── repositorio.go
│   │   │       ├── tenant.go
│   │   │       └── tenant_test.go
│   │   └── infrastructure
│   │       └── persistence
│   │           └── postgres
│   │               ├── tenant_model.go
│   │               └── tenant_repositorio.go
│   ├── usuarios
│   │   ├── application
│   │   │   ├── services
│   │   │   │   └── registro
│   │   │   │       ├── comando.go
│   │   │   │       ├── ejecutor.go
│   │   │   │       ├── respuesta.go
│   │   │   │       ├── servicio_registro.go
│   │   │   │       ├── servicio_registro_test.go
│   │   │   │       └── validador_correo.go
│   │   │   └── usecases
│   │   │       ├── createuser
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── deleteuser
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── expeluser
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── listusers
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── register
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   ├── usecase.go
│   │   │       │   └── usecase_test.go
│   │   │       ├── updatemyprofile
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       ├── updateuser
│   │   │       │   ├── command.go
│   │   │       │   ├── response.go
│   │   │       │   └── usecase.go
│   │   │       └── viewmyprofile
│   │   │           ├── command.go
│   │   │           ├── response.go
│   │   │           └── usecase.go
│   │   ├── domain
│   │   │   └── usuario
│   │   │       ├── correo_electronico.go
│   │   │       ├── correo_electronico_test.go
│   │   │       ├── errors.go
│   │   │       ├── especificacion_usuario.go
│   │   │       ├── estado_verificacion_correo.go
│   │   │       ├── estado_verificacion_correo_test.go
│   │   │       ├── eventos.go
│   │   │       ├── eventos_test.go
│   │   │       ├── repositorio_test.go
│   │   │       ├── ropositorio.go
│   │   │       ├── unit_of_work.go
│   │   │       ├── usuario.go
│   │   │       └── usuario_test.go
│   │   └── infrastructure
│   │       └── persistence
│   │           └── postgres
│   │               ├── unit_of_work.go
│   │               ├── usuario_model.go
│   │               └── usuario_repositorio.go
│   └── verificacion
│       ├── application
│       │   ├── services
│       │   │   └── verificacion
│       │   │       ├── comando.go
│       │   │       ├── respuesta.go
│       │   │       └── servicio_verificacion.go
│       │   └── usecases
│       │       └── verifyemail
│       │           ├── command.go
│       │           ├── response.go
│       │           └── usecase.go
│       ├── domain
│       │   ├── errores.go
│       │   ├── prueba_verificacion.go
│       │   ├── prueba_verificacion_test.go
│       │   └── repositorio.go
│       └── infrastructure
│           └── persistence
│               └── postgres
│                   ├── verificacion_model.go
│                   └── verificacion_repositorio.go
├── main.go
└── shared
    └── presentation
        ├── api_response.go
        └── error_handler.go
```

## 3. CÓDIGO CLAVE


### 3.1 Login Handler

```go
package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/dto"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	presentation "github.com/davosjar/bunna/services/identidad/shared/presentation"
)

// LoginInput es el input del endpoint POST /api/v1/auth/login.
type LoginInput struct {
	RealIP string `header:"X-Real-IP"`
	Body   dto.LoginRequest
}

// LoginOutput es el output del endpoint POST /api/v1/auth/login.
type LoginOutput struct {
	Body presentation.ApiResponse[dto.LoginResponse]
}

// LoginHandler maneja el inicio de sesión.
// Sigue CON-PRES-002: no importa domain ni mapper.
type LoginHandler struct {
	facade facades.AuthFacade
}

// NewLoginHandler construye el handler con su facade.
func NewLoginHandler(facade facades.AuthFacade) *LoginHandler {
	return &LoginHandler{facade: facade}
}

// Register registra el endpoint en la API Huma.
func (h *LoginHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "login-usuario",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/login",
		Summary:     "Inicio de sesión",
		Description: "Autentica al usuario y devuelve tokens JWT de acceso y refresco.",
		Tags:        []string{"Autenticación"},
	}, h.handle)
}

func (h *LoginHandler) handle(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	resp, err := h.facade.Login(ctx, facades.ComandoLogin{
		Email:    input.Body.Correo,
		Password: input.Body.Password,
		IPOrigen: input.RealIP,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	// CON-PRES-005: links HATEOAS los construye el handler
	links := map[string]presentation.Link{
		"self": {
			Href:   "/api/v1/usuarios/" + resp.UsuarioID,
			Method: http.MethodGet,
		},
		"refresh": {
			Href:   "/api/v1/auth/refresh",
			Method: http.MethodPost,
		},
	}

	out := &LoginOutput{}
	out.Body = presentation.NewApiResponseWithLinks(dto.LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		TokenType:    resp.TokenType,
		UsuarioID:    resp.UsuarioID,
		TenantID:     resp.TenantID,
		Rol:          resp.Rol,
	}, links)

	return out, nil
}
```

### 3.2 Register Usecase

```go
package register

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type RegistrarUsuarioCasoDeUso struct {
	userRepo           usuario.UsuarioRepositorio
	credRepo           seguridad.CredencialesRepositorio
	encSvc             seguridad.EncriptacionServicio
	idGen              shareddomain.GeneradorID
	tenantRepo         tenant.TenantRepositorio
	membresiaRepo      tenant.MembresiaRepositorio
	rolRepo            rbac.RolRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
}

func NewRegistrarUsuarioCasoDeUso(
	userRepo usuario.UsuarioRepositorio,
	credRepo seguridad.CredencialesRepositorio,
	encSvc seguridad.EncriptacionServicio,
	idGen shareddomain.GeneradorID,
	tenantRepo tenant.TenantRepositorio,
	membresiaRepo tenant.MembresiaRepositorio,
	rolRepo rbac.RolRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
) *RegistrarUsuarioCasoDeUso {
	return &RegistrarUsuarioCasoDeUso{
		userRepo:             userRepo,
		credRepo:             credRepo,
		encSvc:               encSvc,
		idGen:                idGen,
		tenantRepo:           tenantRepo,
		membresiaRepo:        membresiaRepo,
		rolRepo:              rolRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
	}
}

func (uc *RegistrarUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoRegistrarUsuario) (*RespuestaRegistrarUsuario, error) {
	if err := validarComando(cmd); err != nil {
		return nil, err
	}

	usuarioID, err := uc.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID de usuario: %w", err)
	}

	nuevoUsuario, err := usuario.NuevoUsuario(usuarioID, cmd.Correo, cmd.Nombre, cmd.Apellido, cmd.Telefono)
```

### 3.3 JWT Middleware

```go
// Package middleware contiene los middlewares HTTP de la capa de presentación.
package middleware

import (
	"context"
	"net/http"
	"strings"

	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	"github.com/gin-gonic/gin"
)

// Claves usadas en gin.Context.Set/Get (tipo string para compatibilidad gin).
const (
	ClaveUsuarioID = "usuarioID"
	ClaveSesionID  = "sesionID"
	ClaveTenantID  = "tenantID"
	ClaveRol       = "rol"
)

// ctxKeyUsuarioID y ctxKeySesionID son claves con tipo para context.Context (evita colisiones).
type ctxKey string

const (
	ctxKeyUsuarioID ctxKey = "usuarioID"
	ctxKeySesionID  ctxKey = "sesionID"
	ctxKeyTenantID  ctxKey = "tenantID"
	ctxKeyRol       ctxKey = "rol"
)

// GetUsuarioIDFromCtx extrae el usuarioID del context.Context (útil para handlers Huma).
func GetUsuarioIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyUsuarioID).(string); ok {
		return v
	}
	return ""
}

// GetSesionIDFromCtx extrae el sesionID del context.Context (útil para handlers Huma).
func GetSesionIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeySesionID).(string); ok {
		return v
	}
	return ""
}

// GetTenantIDFromCtx extrae el tenantID del context.Context (útil para handlers Huma).
func GetTenantIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyTenantID).(string); ok {
		return v
	}
	return ""
}

// GetRolFromCtx extrae el rol del context.Context (útil para handlers Huma).
func GetRolFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyRol).(string); ok {
		return v
	}
	return ""
}

// JWTMiddleware valida el token Bearer del header Authorization.
// Si el token es válido, inyecta usuarioID y sesionID en el contexto Gin y en el context.Context.
// Si el token es inválido o ausente, responde 401 y aborta.
func JWTMiddleware(tokenSvc sesiones_domain.TokenServicio) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"title":  "Unauthorized",
				"detail": "header Authorization ausente",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"title":  "Unauthorized",
				"detail": "formato de Authorization inválido, se esperaba: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]
		claims, err := tokenSvc.ValidarAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"title":  "Unauthorized",
				"detail": "token inválido o expirado",
			})
			return
		}

		// Inyectar claims en contexto Gin para handlers gin nativos
		c.Set(ClaveUsuarioID, claims.UsuarioID)
		c.Set(ClaveSesionID, claims.SesionID)
		c.Set(ClaveTenantID, claims.TenantID)
		c.Set(ClaveRol, claims.Rol)

		// Inyectar claims en context.Context para handlers Huma
		reqCtx := context.WithValue(c.Request.Context(), ctxKeyUsuarioID, claims.UsuarioID)
		reqCtx = context.WithValue(reqCtx, ctxKeySesionID, claims.SesionID)
		reqCtx = context.WithValue(reqCtx, ctxKeyTenantID, claims.TenantID)
		reqCtx = context.WithValue(reqCtx, ctxKeyRol, claims.Rol)
		c.Request = c.Request.WithContext(reqCtx)

		c.Next()
	}
}
```

### 3.4 Modelo Credenciales

```go
package domain

import "time"

type CredencialesUsuario struct {
	usuarioID        string
	passwordHash     string
	activo           bool
	correoVerificado bool
	intentosFallidos int
	bloqueadoHasta   time.Time
}

func NuevaCredencialesUsuario(usuarioID, passwordHash string) *CredencialesUsuario {
	return &CredencialesUsuario{
		usuarioID:        usuarioID,
		passwordHash:     passwordHash,
		activo:           true,
		correoVerificado: false,
		intentosFallidos: 0,
		bloqueadoHasta:   time.Time{},
	}
}

func NuevaCredencialesUsuarioDesdeBD(usuarioID, passwordHash string, activo, correoVerificado bool, intentosFallidos int, bloqueadoHasta time.Time) *CredencialesUsuario {
	return &CredencialesUsuario{
		usuarioID:        usuarioID,
		passwordHash:     passwordHash,
		activo:           activo,
		correoVerificado: correoVerificado,
		intentosFallidos: intentosFallidos,
		bloqueadoHasta:   bloqueadoHasta,
	}
}

func (c *CredencialesUsuario) VerificarPassword(passwordHash string) bool {
	return c.passwordHash == passwordHash
}

func (c *CredencialesUsuario) IncrementarIntentoFallido() {
	c.intentosFallidos++
}

func (c *CredencialesUsuario) BloquearHasta(hasta time.Time) {
	c.bloqueadoHasta = hasta
}

// MarcarIntentoFallido exists for backward compatibility.
// New code should use IncrementarIntentoFallido + BloquearHasta directly.
func (c *CredencialesUsuario) MarcarIntentoFallido(ahora time.Time) {
	c.IncrementarIntentoFallido()
	if c.intentosFallidos >= 5 {
		c.BloquearHasta(ahora.Add(15 * time.Minute))
	}
}

func (c *CredencialesUsuario) ResetearIntentos() {
	c.intentosFallidos = 0
	c.bloqueadoHasta = time.Time{}
}

func (c *CredencialesUsuario) EstaBloqueado(ahora time.Time) bool {
	return ahora.Before(c.bloqueadoHasta)
}

func (c *CredencialesUsuario) VerificarCorreo() {
	c.correoVerificado = true
}

func (c *CredencialesUsuario) Desactivar() {
	c.activo = false
}

func (c *CredencialesUsuario) Activar() {
	c.activo = true
}

// Getters públicos para acceso de lectura
func (c *CredencialesUsuario) UsuarioID() string {
	return c.usuarioID
}

func (c *CredencialesUsuario) PasswordHash() string {
	return c.passwordHash
}

func (c *CredencialesUsuario) Activo() bool {
	return c.activo
}

func (c *CredencialesUsuario) CorreoVerificado() bool {
	return c.correoVerificado
}

func (c *CredencialesUsuario) IntentosFallidos() int {
	return c.intentosFallidos
}

func (c *CredencialesUsuario) BloqueadoHasta() time.Time {
	return c.bloqueadoHasta
}
```

### 3.5 Modelo Usuario (dominio)

```go
package usuario

import (
	"time"
)

type EstadoUsuario string

const (
	NO_VERIFICADO            EstadoUsuario = "NO_VERIFICADO"
	ACTIVO                   EstadoUsuario = "ACTIVO"
	INACTIVO                 EstadoUsuario = "INACTIVO"
	PENDIENTE_DE_ELIMINACION EstadoUsuario = "PENDIENTE_DE_ELIMINACION"
	BLOQUEADO                EstadoUsuario = "BLOQUEADO"
)

var transiciones = map[EstadoUsuario]map[EstadoUsuario]bool{
	NO_VERIFICADO: {
		ACTIVO:                   true,
		PENDIENTE_DE_ELIMINACION: true,
	},
	ACTIVO: {
		INACTIVO:                 true,
		BLOQUEADO:                true,
		PENDIENTE_DE_ELIMINACION: true,
	},
	INACTIVO: {
		ACTIVO:                   true,
		PENDIENTE_DE_ELIMINACION: true,
	},
	BLOQUEADO: {
		ACTIVO:                   true,
		PENDIENTE_DE_ELIMINACION: true,
	},
	PENDIENTE_DE_ELIMINACION: {},
}

type Usuario struct {
	id                 string
	nombre             string
	apellido           string
	correoElectronico  *CorreoElectronico
	telefono           string
	estado             EstadoUsuario
	fechaCreacion      time.Time
	fechaActualizacion time.Time
	eventos            *EventosUsuario
}

func NuevoUsuario(id, correo, nombre, apellido, telefono string) (*Usuario, error) {
	correoVO, err := NuevoCorreoElectronico(correo)
	if err != nil {
		return nil, err
	}

	ahora := time.Now()
	u := &Usuario{
		id:                 id,
		nombre:             nombre,
		apellido:           apellido,
		correoElectronico:  correoVO,
		telefono:           telefono,
		estado:             NO_VERIFICADO,
		fechaCreacion:      ahora,
		fechaActualizacion: ahora,
		eventos:            NuevosEventosUsuario(),
	}

	u.eventos.RegistrarCreacion(u.id, u.correoElectronico.Direccion())
	return u, nil
}

func NewUsuarioFromPersistence(id, nombre, apellido string, correoElectronico *CorreoElectronico, telefono string, estado EstadoUsuario, fechaCreacion, fechaActualizacion time.Time) *Usuario {
	return &Usuario{
		id:                 id,
		nombre:             nombre,
		apellido:           apellido,
		correoElectronico:  correoElectronico,
		telefono:           telefono,
		estado:             estado,
		fechaCreacion:      fechaCreacion,
		fechaActualizacion: fechaActualizacion,
		eventos:            NuevosEventosUsuario(),
	}
}

// Getters
func (u *Usuario) ID() string            { return u.id }
func (u *Usuario) Nombre() string        { return u.nombre }
func (u *Usuario) Apellido() string      { return u.apellido }
func (u *Usuario) Correo() string        { return u.correoElectronico.Direccion() }
func (u *Usuario) Telefono() string      { return u.telefono }
func (u *Usuario) Estado() EstadoUsuario { return u.estado }
func (u *Usuario) EstadoVerificacionCorreo() EstadoVerificacionCorreo {
	return u.correoElectronico.Estado()
}
func (u *Usuario) FechaCreacion() time.Time      { return u.fechaCreacion }
func (u *Usuario) FechaActualizacion() time.Time { return u.fechaActualizacion }

func (u *Usuario) CambiarEstado(siguiente EstadoUsuario) error {
	if u.estado == siguiente {
		return nil
	}
	if !transiciones[u.estado][siguiente] {
		return ErrTransicionNoPermitida
	}
	u.estado = siguiente
	u.eventos.RegistrarCambioEstado(u.id, siguiente)
	return nil
}

func (u *Usuario) Bloquear() error {
	if err := u.CambiarEstado(BLOQUEADO); err != nil {
		return err
	}
	u.eventos.RegistrarBloqueo(u.id)
	return nil
}

func (u *Usuario) Activar() error {
	if err := u.CambiarEstado(ACTIVO); err != nil {
		return err
	}
	u.eventos.RegistrarActivacion(u.id)
	return nil
}

func (u *Usuario) Inactivar() error {
	if err := u.CambiarEstado(INACTIVO); err != nil {
		return err
	}
	u.eventos.RegistrarInactivacion(u.id)
	return nil
}

func (u *Usuario) VerificarCorreo() error {
	if err := u.correoElectronico.Verificar(); err != nil {
		return ErrTransicionVerificacionNoPermitida
	}
	u.eventos.RegistrarVerificacion(u.id)
	return nil
}

func (u *Usuario) SolicitarReenvioVerificacion() error {
	if err := u.correoElectronico.SolicitarReenvio(); err != nil {
		return ErrTransicionVerificacionNoPermitida
	}
	u.eventos.registrarEvento("ReenvioVerificacionSolicitado", map[string]string{"id": u.id})
	return nil
}

func (u *Usuario) MarcarEnlaceExpirado() error {
	if err := u.correoElectronico.MarcarExpirado(); err != nil {
		return ErrTransicionVerificacionNoPermitida
	}
	u.eventos.registrarEvento("EnlaceVerificacionExpirado", map[string]string{"id": u.id})
	return nil
}

func (u *Usuario) PullEventos() []EventoDominio {
	return u.eventos.Extraer()
}

func (u *Usuario) Eventos() *EventosUsuario {
	return u.eventos
}
```
