package invitaciones

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Invitacion struct {
	id              string
	tenantID        string
	rolID           string
	email           string
	nombre          string
	tokenHash       string
	expiracion      time.Time
	aceptada        bool
	fechaCreacion   time.Time
	fechaAceptacion *time.Time
}

func NuevaInvitacion(id, tenantID, rolID, email, nombre, tokenEnPlano string, expiracion time.Time) *Invitacion {
	hash := hashearToken(tokenEnPlano)
	return &Invitacion{
		id:            id,
		tenantID:      tenantID,
		rolID:         rolID,
		email:         email,
		nombre:        nombre,
		tokenHash:     hash,
		expiracion:    expiracion,
		aceptada:      false,
		fechaCreacion: time.Now(),
	}
}

func NuevaInvitacionDesdeBD(id, tenantID, rolID, email, nombre, tokenHash string, expiracion time.Time, aceptada bool, fechaCreacion time.Time, fechaAceptacion *time.Time) *Invitacion {
	return &Invitacion{
		id:              id,
		tenantID:        tenantID,
		rolID:           rolID,
		email:           email,
		nombre:          nombre,
		tokenHash:       tokenHash,
		expiracion:      expiracion,
		aceptada:        aceptada,
		fechaCreacion:   fechaCreacion,
		fechaAceptacion: fechaAceptacion,
	}
}

func (i *Invitacion) ActualizarToken(tokenEnPlano string) {
	i.tokenHash = hashearToken(tokenEnPlano)
}

func (i *Invitacion) Expiro(ahora time.Time) bool {
	return ahora.After(i.expiracion)
}

func (i *Invitacion) CoincideToken(tokenEnPlano string) bool {
	return i.tokenHash == hashearToken(tokenEnPlano)
}

func (i *Invitacion) Aceptar() error {
	if i.aceptada {
		return ErrYaAceptada
	}
	if i.Expiro(time.Now()) {
		return ErrEnlaceExpirado
	}
	i.aceptada = true
	ahora := time.Now()
	i.fechaAceptacion = &ahora
	return nil
}

func (i *Invitacion) ID() string              { return i.id }
func (i *Invitacion) TenantID() string         { return i.tenantID }
func (i *Invitacion) RolID() string            { return i.rolID }
func (i *Invitacion) Email() string            { return i.email }
func (i *Invitacion) Nombre() string           { return i.nombre }
func (i *Invitacion) TokenHash() string        { return i.tokenHash }
func (i *Invitacion) Expiracion() time.Time    { return i.expiracion }
func (i *Invitacion) EstaAceptada() bool       { return i.aceptada }
func (i *Invitacion) FechaCreacion() time.Time { return i.fechaCreacion }
func (i *Invitacion) FechaAceptacion() *time.Time { return i.fechaAceptacion }

func hashearToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

func HashearTokenPublico(token string) string {
	return hashearToken(token)
}
