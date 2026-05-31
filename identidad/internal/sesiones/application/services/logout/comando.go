// Package logout implementa el caso de uso de cierre de sesión explícito.
package logout

// ComandoLogout contiene los datos necesarios para cerrar una sesión específica.
type ComandoLogout struct {
	// SesionID es el identificador de la sesión a cerrar.
	SesionID string

	// UsuarioID es el identificador del usuario autenticado que solicita el logout.
	// Se usa para verificar que la sesión pertenece al usuario.
	UsuarioID string
}

// ComandoCerrarTodas contiene los datos para cerrar todas las sesiones de un usuario.
type ComandoCerrarTodas struct {
	// UsuarioID es el identificador del usuario cuyas sesiones se cerrarán.
	UsuarioID string
}
