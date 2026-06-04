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

// CambiarHash actualiza el hash de la contraseña y resetea intentos fallidos
func (c *CredencialesUsuario) CambiarHash(nuevoHash string) {
	c.passwordHash = nuevoHash
	c.intentosFallidos = 0
	c.bloqueadoHasta = time.Time{}
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
