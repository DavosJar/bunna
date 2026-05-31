package rbac

// Nombres de roles de sistema — inmutables
const (
	RolSysAdmin      = "sys_admin"
	RolAdministrador = "administrador"
	RolAgronomo      = "agronomo"
	RolCaficultor    = "caficultor"
)

// RolInfo contiene los metadatos de un rol de sistema
type RolInfo struct {
	Nombre      string
	Descripcion string
	EsSistema   bool
	Permisos    []string
}

// RolesDeSistema define los 4 roles inmutables con sus permisos
var RolesDeSistema = []RolInfo{
	{
		Nombre:      RolSysAdmin,
		Descripcion: "Super administrador global con todos los permisos",
		EsSistema:   true,
		Permisos: []string{
			PermisoUsuarioCrear,
			PermisoUsuarioModificar,
			PermisoUsuarioEliminar,
			PermisoUsuarioConsultar,
			PermisoUsuarioResetearPassword,
			PermisoUsuarioExpulsar,
			PermisoCredencialesConsultar,
			PermisoCredencialesDesbloquear,
			PermisoIPBloqueadaConsultar,
			PermisoIPDesbloquear,
			PermisoSesionConsultar,
			PermisoSesionForzarCierre,
			PermisoRolCrear,
			PermisoRolModificar,
			PermisoRolEliminar,
			PermisoRolAsignar,
			PermisoRolRevocar,
			PermisoRolPermisoAsignar,
			PermisoRolPermisoRevocar,
			PermisoPermisoConsultar,
			PermisoTenantConfigurar,
		},
	},
	{
		Nombre:      RolAdministrador,
		Descripcion: "Administrador del tenant con permisos completos dentro de su tenant",
		EsSistema:   true,
		Permisos: []string{
			PermisoUsuarioCrear,
			PermisoUsuarioModificar,
			PermisoUsuarioEliminar,
			PermisoUsuarioConsultar,
			PermisoUsuarioResetearPassword,
			PermisoUsuarioExpulsar,
			PermisoCredencialesConsultar,
			PermisoCredencialesDesbloquear,
			PermisoIPBloqueadaConsultar,
			PermisoIPDesbloquear,
			PermisoSesionConsultar,
			PermisoSesionForzarCierre,
			PermisoRolCrear,
			PermisoRolModificar,
			PermisoRolEliminar,
			PermisoRolAsignar,
			PermisoRolRevocar,
			PermisoRolPermisoAsignar,
			PermisoRolPermisoRevocar,
			PermisoPermisoConsultar,
			PermisoTenantConfigurar,
		},
	},
	{
		Nombre:      RolAgronomo,
		Descripcion: "Permisos intermedios: crear y modificar usuarios, consultar",
		EsSistema:   true,
		Permisos: []string{
			PermisoUsuarioCrear,
			PermisoUsuarioModificar,
			PermisoUsuarioConsultar,
			PermisoPermisoConsultar,
		},
	},
	{
		Nombre:      RolCaficultor,
		Descripcion: "Solo consulta de usuarios",
		EsSistema:   true,
		Permisos: []string{
			PermisoUsuarioConsultar,
		},
	},
}
