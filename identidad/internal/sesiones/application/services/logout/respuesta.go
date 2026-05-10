package logout

// RespuestaLogout confirma el resultado del cierre de sesión.
type RespuestaLogout struct {
	// SesionesRevocadas es la cantidad de sesiones que fueron revocadas.
	SesionesRevocadas int
}