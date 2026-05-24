package domain

import "time"

type EstadoDiagnostico string

const (
	EstadoDiagnosticoPendiente EstadoDiagnostico = "PENDIENTE"
	EstadoDiagnosticoAceptado  EstadoDiagnostico = "ACEPTADO"
	EstadoDiagnosticoRechazado EstadoDiagnostico = "RECHAZADO"
)

// Mapa transiciones
var transiciones = map[EstadoDiagnostico]map[EstadoDiagnostico]bool{
	EstadoDiagnosticoPendiente: {
		EstadoDiagnosticoAceptado:  true,
		EstadoDiagnosticoRechazado: true,
	},
}

func (e EstadoDiagnostico) EsValido() bool {
	_, ok := transiciones[e]
	return ok
}

func (e EstadoDiagnostico) PuedeTransicionarA(nuevo EstadoDiagnostico) bool {
	return transiciones[e][nuevo]
}

type Diagnostico struct {
	id                  string
	nombre              string
	muestrasId          string
	tenantID            string
	estado              EstadoDiagnostico
	resultadoInferencia *ResultadoInferencia
	updatedAt           time.Time
	createdAt           time.Time
}

func NewDiagnostico(id, nombre, muestrasId, tenantID string, resultadoInferencia *ResultadoInferencia) (*Diagnostico, error) {
	if nombre == "" {
		return nil, ErrNombreRequerido
	}
	if muestrasId == "" {
		return nil, ErrMuestrasIdRequerido
	}
	if tenantID == "" {
		return nil, ErrTenantIdRequerido
	}
	if resultadoInferencia == nil {
		return nil, ErrResultadoInferenciaRequerido
	}
	return &Diagnostico{
		id:                  id,
		nombre:              nombre,
		muestrasId:          muestrasId,
		tenantID:            tenantID,
		resultadoInferencia: resultadoInferencia,
		estado:              EstadoDiagnosticoPendiente,
		createdAt:           time.Now(),
		updatedAt:           time.Now(),
	}, nil
}

func NewDiagnosticoFromStorage(id, nombre, muestrasId, tenantID string, resultadoInferencia *ResultadoInferencia, createdAt, updatedAt time.Time, estado EstadoDiagnostico) (*Diagnostico, error) {
	return &Diagnostico{
		id:                  id,
		nombre:              nombre,
		muestrasId:          muestrasId,
		tenantID:            tenantID,
		resultadoInferencia: resultadoInferencia,
		createdAt:           createdAt,
		updatedAt:           updatedAt,
		estado:              estado,
	}, nil
}

func (d *Diagnostico) MarcarComoAceptado() error {
	if !d.estado.PuedeTransicionarA(EstadoDiagnosticoAceptado) {
		return ErrDiagnosticoNoEncontrado
	}
	d.estado = EstadoDiagnosticoAceptado
	d.updatedAt = time.Now()
	return nil
}

func (d *Diagnostico) MarcarComoRechazado() error {
	if !d.estado.PuedeTransicionarA(EstadoDiagnosticoRechazado) {
		return ErrDiagnosticoNoEncontrado
	}
	d.estado = EstadoDiagnosticoRechazado
	d.updatedAt = time.Now()
	return nil
}

// getters sin setters
func (d *Diagnostico) ID() string                                { return d.id }
func (d *Diagnostico) Nombre() string                            { return d.nombre }
func (d *Diagnostico) MuestrasId() string                        { return d.muestrasId }
func (d *Diagnostico) TenantID() string                          { return d.tenantID }
func (d *Diagnostico) Estado() EstadoDiagnostico                 { return d.estado }
func (d *Diagnostico) ResultadoInferencia() *ResultadoInferencia { return d.resultadoInferencia }
func (d *Diagnostico) CreatedAt() time.Time                      { return d.createdAt }
func (d *Diagnostico) UpdatedAt() time.Time                      { return d.updatedAt }
