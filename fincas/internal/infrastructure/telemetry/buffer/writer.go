package buffer

import "errors"

// Prioridad defines the priority levels for telemetry events.
type Prioridad int

const (
    Alta Prioridad = iota
    Media
    Baja
)

var (
	ErrBufferLleno          = errors.New("buffer lleno: el segmento correspondiente está saturado")
	ErrPrioritarioTimeout   = errors.New("timeout en segmento prioritario: no se pudo insertar el evento de alta prioridad")
)

// BufferWriter is the contract for writing telemetry events to a buffer.
// The event is provided as a byte slice and a priority level.
type BufferWriter interface {
    Write(event []byte, priority Prioridad) error
}
