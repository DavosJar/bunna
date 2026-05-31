package tenant

import "errors"

var (
	ErrNombreRequerido     = errors.New("nombre del tenant requerido")
	ErrSlugRequerido       = errors.New("slug del tenant requerido")
	ErrSlugInvalido        = errors.New("slug inválido: solo se permiten letras minúsculas, números y guiones")
	ErrSlugDuplicado       = errors.New("ya existe un tenant con ese slug")
	ErrTenantNoEncontrado  = errors.New("tenant no encontrado")
	ErrTenantInactivo      = errors.New("el tenant está inactivo")
	ErrTenantYaActivo      = errors.New("el tenant ya está activo")
	ErrTenantYaInactivo    = errors.New("el tenant ya está inactivo")
	ErrUsuarioYaMiembro    = errors.New("el usuario ya es miembro del tenant")
	ErrUsuarioNoEsMiembro  = errors.New("el usuario no es miembro del tenant")
	ErrUltimoAdministrador = errors.New("no se puede remover al último administrador del tenant")
	ErrSinPermiso          = errors.New("no tiene permisos para realizar esta operación")
)
