package rbac

import "context"

// Permisos atómicos del sistema — constantes de dominio, no modificables en runtime
const (
	PermisoUsuarioCrear            = "identidad:usuario:crear"
	PermisoUsuarioModificar        = "identidad:usuario:modificar"
	PermisoUsuarioEliminar         = "identidad:usuario:eliminar"
	PermisoUsuarioConsultar        = "identidad:usuario:consultar"
	PermisoUsuarioResetearPassword = "identidad:usuario:resetear_password"
	PermisoUsuarioExpulsar         = "identidad:usuario:expulsar"
	PermisoUsuarioInvitar          = "identidad:usuario:invitar"
	PermisoCredencialesConsultar   = "identidad:credenciales:consultar"
	PermisoCredencialesDesbloquear = "identidad:credenciales:desbloquear"
	PermisoRolAsignar              = "identidad:rol:asignar"
	PermisoRolRevocar              = "identidad:rol:revocar"
	PermisoRolCrear                = "identidad:rol:crear"
	PermisoRolModificar            = "identidad:rol:modificar"
	PermisoRolEliminar             = "identidad:rol:eliminar"
	PermisoRolPermisoAsignar       = "identidad:rol:permiso:asignar"
	PermisoRolPermisoRevocar       = "identidad:rol:permiso:revocar"
	PermisoPermisoConsultar        = "identidad:permiso:consultar"
	PermisoSesionConsultar         = "identidad:sesion:consultar"
	PermisoSesionForzarCierre      = "identidad:sesion:forzar_cierre"
	PermisoTenantConfigurar        = "identidad:tenant:configurar"
	PermisoIPBloqueadaConsultar    = "identidad:ip:consultar"
	PermisoIPDesbloquear           = "identidad:ip:desbloquear"
)

// TodosLosPermisos lista todos los permisos del sistema
var TodosLosPermisos = []PermisoInfo{
	{Codigo: PermisoUsuarioCrear, Nombre: "Crear Usuario", Descripcion: "Crear nuevos usuarios con asignación opcional de rol", Modulo: "identidad"},
	{Codigo: PermisoUsuarioModificar, Nombre: "Modificar Usuario", Descripcion: "Modificar datos personales de cualquier usuario", Modulo: "identidad"},
	{Codigo: PermisoUsuarioEliminar, Nombre: "Eliminar Usuario", Descripcion: "Marcar un usuario como pendiente de eliminación", Modulo: "identidad"},
	{Codigo: PermisoUsuarioConsultar, Nombre: "Consultar Usuario", Descripcion: "Listar y ver detalles de cualquier usuario", Modulo: "identidad"},
	{Codigo: PermisoUsuarioResetearPassword, Nombre: "Resetear Contraseña", Descripcion: "Resetear la contraseña de otro usuario", Modulo: "identidad"},
	{Codigo: PermisoUsuarioExpulsar, Nombre: "Expulsar Usuario", Descripcion: "Dar de baja inmediata y revocar todas las sesiones activas", Modulo: "identidad"},
	{Codigo: PermisoUsuarioInvitar, Nombre: "Invitar Usuario", Descripcion: "Crear invitaciones para que nuevos usuarios se unan al tenant", Modulo: "identidad"},
	{Codigo: PermisoCredencialesConsultar, Nombre: "Consultar Credenciales", Descripcion: "Ver el estado de seguridad de un usuario", Modulo: "identidad"},
	{Codigo: PermisoCredencialesDesbloquear, Nombre: "Desbloquear Cuenta", Descripcion: "Quitar el bloqueo de una cuenta inhabilitada por intentos fallidos", Modulo: "identidad"},
	{Codigo: PermisoRolAsignar, Nombre: "Asignar Rol", Descripcion: "Asignar un rol a un usuario", Modulo: "identidad"},
	{Codigo: PermisoRolRevocar, Nombre: "Revocar Rol", Descripcion: "Revocar un rol de un usuario", Modulo: "identidad"},
	{Codigo: PermisoRolCrear, Nombre: "Crear Rol", Descripcion: "Crear un rol personalizado en la organización", Modulo: "identidad"},
	{Codigo: PermisoRolModificar, Nombre: "Modificar Rol", Descripcion: "Cambiar el nombre o descripción de un rol personalizado", Modulo: "identidad"},
	{Codigo: PermisoRolEliminar, Nombre: "Eliminar Rol", Descripcion: "Eliminar un rol personalizado", Modulo: "identidad"},
	{Codigo: PermisoRolPermisoAsignar, Nombre: "Asignar Permiso a Rol", Descripcion: "Agregar un permiso a un rol personalizado", Modulo: "identidad"},
	{Codigo: PermisoRolPermisoRevocar, Nombre: "Revocar Permiso de Rol", Descripcion: "Quitar un permiso de un rol personalizado", Modulo: "identidad"},
	{Codigo: PermisoPermisoConsultar, Nombre: "Consultar Permisos", Descripcion: "Listar permisos de un rol y roles de un usuario", Modulo: "identidad"},
	{Codigo: PermisoSesionConsultar, Nombre: "Consultar Sesiones", Descripcion: "Ver los dispositivos y sesiones activas de los usuarios", Modulo: "identidad"},
	{Codigo: PermisoSesionForzarCierre, Nombre: "Forzar Cierre de Sesión", Descripcion: "Desconectar la sesión de un usuario de forma remota", Modulo: "identidad"},
	{Codigo: PermisoTenantConfigurar, Nombre: "Configurar Tenant", Descripcion: "Cambiar la configuración global del tenant", Modulo: "identidad"},
	{Codigo: PermisoIPBloqueadaConsultar, Nombre: "Consultar IPs Bloqueadas", Descripcion: "Ver la lista de IPs bloqueadas por ataques", Modulo: "identidad"},
	{Codigo: PermisoIPDesbloquear, Nombre: "Desbloquear IP", Descripcion: "Quitar una IP de la lista negra", Modulo: "identidad"},
}

// AuthorizationService define el contrato para la verificación de permisos
type AuthorizationService interface {
	TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error)
}

// PermisoInfo contiene los metadatos de un permiso
type PermisoInfo struct {
	Codigo      string
	Nombre      string
	Descripcion string
	Modulo      string
}
