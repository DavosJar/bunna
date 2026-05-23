package usuario

import (
	"regexp"
)

var regexEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// CorreoElectronico es un Value Object que encapsula la dirección
// de correo y su estado de verificación.
type CorreoElectronico struct {
	direccion string
	estado    EstadoVerificacionCorreo
}

// NuevoCorreoElectronico crea un VO con estado inicial PENDIENTE_VERIFICACION.
// Retorna error si la dirección está vacía o tiene formato inválido.
func NuevoCorreoElectronico(direccion string) (*CorreoElectronico, error) {
	if direccion == "" {
		return nil, ErrCorreoRequerido
	}
	if !regexEmail.MatchString(direccion) {
		return nil, ErrCorreoFormatoInvalido
	}
	return &CorreoElectronico{
		direccion: direccion,
		estado:    PENDIENTE_VERIFICACION,
	}, nil
}

// NuevoCorreoElectronicoDesdeBD reconstruye el VO desde persistencia
// sin validar ni emitir eventos.
func NuevoCorreoElectronicoDesdeBD(direccion string, estado EstadoVerificacionCorreo) *CorreoElectronico {
	return &CorreoElectronico{
		direccion: direccion,
		estado:    estado,
	}
}

// Getters
func (c *CorreoElectronico) Direccion() string              { return c.direccion }
func (c *CorreoElectronico) Estado() EstadoVerificacionCorreo { return c.estado }
func (c *CorreoElectronico) EstaVerificado() bool           { return c.estado == VERIFICADO }
func (c *CorreoElectronico) EstaPendiente() bool            { return c.estado == PENDIENTE_VERIFICACION }

// Verificar transiciona el estado a VERIFICADO.
func (c *CorreoElectronico) Verificar() error {
	if !c.estado.PuedeTransicionarA(VERIFICADO) {
		return ErrTransicionVerificacionNoPermitida
	}
	c.estado = VERIFICADO
	return nil
}

// MarcarExpirado transiciona el estado a ENLACE_EXPIRADO.
func (c *CorreoElectronico) MarcarExpirado() error {
	if !c.estado.PuedeTransicionarA(ENLACE_EXPIRADO) {
		return ErrTransicionVerificacionNoPermitida
	}
	c.estado = ENLACE_EXPIRADO
	return nil
}

// SolicitarReenvio transiciona el estado a REENVIO_SOLICITADO.
func (c *CorreoElectronico) SolicitarReenvio() error {
	if !c.estado.PuedeTransicionarA(REENVIO_SOLICITADO) {
		return ErrTransicionVerificacionNoPermitida
	}
	c.estado = REENVIO_SOLICITADO
	return nil
}