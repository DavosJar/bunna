package rbac

import "errors"

var (
	ErrPermisoDenegado           = errors.New("permiso denegado")
	ErrRolNoEncontrado           = errors.New("rol no encontrado")
	ErrRolYaAsignado             = errors.New("el usuario ya tiene ese rol en este contexto")
	ErrRolNoAsignado             = errors.New("el usuario no tiene ese rol en este contexto")
	ErrRolInmutable              = errors.New("el rol es de sistema y no puede modificarse ni eliminarse")
	ErrPasswordActualIncorrecto  = errors.New("la contraseña actual no coincide")
	ErrUsuarioNoPerteneceTenant  = errors.New("el usuario no es miembro del tenant")
	ErrSysAdminRequiereTenantVacio = errors.New("el rol sys_admin es global y no acepta tenant_id")
)