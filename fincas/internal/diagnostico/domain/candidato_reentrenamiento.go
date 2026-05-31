package domain

import "time"

// CandidatoReentrenamiento representa un diagnóstico rechazado cuya imagen
// debe ser revisada para mejorar el modelo de inferencia.
// No es un aggregate raíz; es una entidad dependiente de Diagnostico.
type CandidatoReentrenamiento struct {
	id                    string
	diagnosticoID         string
	imageURL              string
	tieneClorosis         bool
	confianza             float64
	motivo                string
	rechazadoPorUsuarioID string
	createdAt             time.Time
}

func NewCandidatoReentrenamiento(
	id, diagnosticoID, imageURL string,
	tieneClorosis bool,
	confianza float64,
	motivo, rechazadoPorUsuarioID string,
) (*CandidatoReentrenamiento, error) {
	if id == "" {
		return nil, ErrIdRequerido
	}
	if diagnosticoID == "" {
		return nil, ErrDiagnosticoNoEncontrado
	}
	if imageURL == "" {
		return nil, ErrImageUrlRequerida
	}
	return &CandidatoReentrenamiento{
		id:                    id,
		diagnosticoID:         diagnosticoID,
		imageURL:              imageURL,
		tieneClorosis:         tieneClorosis,
		confianza:             confianza,
		motivo:                motivo,
		rechazadoPorUsuarioID: rechazadoPorUsuarioID,
		createdAt:             time.Now(),
	}, nil
}

func NewCandidatoReentrenamientoFromStorage(
	id, diagnosticoID, imageURL string,
	tieneClorosis bool,
	confianza float64,
	motivo, rechazadoPorUsuarioID string,
	createdAt time.Time,
) *CandidatoReentrenamiento {
	return &CandidatoReentrenamiento{
		id:                    id,
		diagnosticoID:         diagnosticoID,
		imageURL:              imageURL,
		tieneClorosis:         tieneClorosis,
		confianza:             confianza,
		motivo:                motivo,
		rechazadoPorUsuarioID: rechazadoPorUsuarioID,
		createdAt:             createdAt,
	}
}

func (c *CandidatoReentrenamiento) ID() string                    { return c.id }
func (c *CandidatoReentrenamiento) DiagnosticoID() string         { return c.diagnosticoID }
func (c *CandidatoReentrenamiento) ImageURL() string              { return c.imageURL }
func (c *CandidatoReentrenamiento) TieneClorosis() bool           { return c.tieneClorosis }
func (c *CandidatoReentrenamiento) Confianza() float64            { return c.confianza }
func (c *CandidatoReentrenamiento) Motivo() string                { return c.motivo }
func (c *CandidatoReentrenamiento) RechazadoPorUsuarioID() string { return c.rechazadoPorUsuarioID }
func (c *CandidatoReentrenamiento) CreatedAt() time.Time          { return c.createdAt }
