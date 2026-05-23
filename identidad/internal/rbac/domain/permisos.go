package rbac

// Permisos atómicos del sistema — constantes de dominio, no modificables en runtime
const (
	PermisoUsuarioCrear           = "identidad:usuario:crear"
	PermisoUsuarioModificar       = "identidad:usuario:modificar"
	PermisoUsuarioEliminar        = "identidad:usuario:eliminar"
	PermisoUsuarioConsultar       = "identidad:usuario:consultar"
	PermisoUsuarioResetearPassword = "identidad:usuario:resetear_password"
	PermisoRolAsignar             = "identidad:rol:asignar"
	PermisoRolRevocar             = "identidad:rol:revocar"
	PermisoPermisoConsultar       = "identidad:permiso:consultar"
)

// TodosLosPermisos lista todos los permisos del sistema
var TodosLosPermisos = []PermisoInfo{
	{Codigo: PermisoUsuarioCrear, Nombre: "Crear Usuario", Descripcion: "Crear nuevos usuarios con asignación opcional de rol", Modulo: "identidad"},
	{Codigo: PermisoUsuarioModificar, Nombre: "Modificar Usuario", Descripcion: "Modificar datos personales de cualquier usuario", Modulo: "identidad"},
	{Codigo: PermisoUsuarioEliminar, Nombre: "Eliminar Usuario", Descripcion: "Marcar un usuario como pendiente de eliminación", Modulo: "identidad"},
	{Codigo: PermisoUsuarioConsultar, Nombre: "Consultar Usuario", Descripcion: "Listar y ver detalles de cualquier usuario", Modulo: "identidad"},
	{Codigo: PermisoUsuarioResetearPassword, Nombre: "Resetear Contraseña", Descripcion: "Resetear la contraseña de otro usuario", Modulo: "identidad"},
	{Codigo: PermisoRolAsignar, Nombre: "Asignar Rol", Descripcion: "Asignar un rol a un usuario", Modulo: "identidad"},
	{Codigo: PermisoRolRevocar, Nombre: "Revocar Rol", Descripcion: "Revocar un rol de un usuario", Modulo: "identidad"},
	{Codigo: PermisoPermisoConsultar, Nombre: "Consultar Permisos", Descripcion: "Listar permisos de un rol y roles de un usuario", Modulo: "identidad"},
}

// PermisoInfo contiene los metadatos de un permiso
type PermisoInfo struct {
	Codigo      string
	Nombre      string
	Descripcion string
	Modulo      string
}