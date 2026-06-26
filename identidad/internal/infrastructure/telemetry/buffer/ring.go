package buffer

import (
	"errors"
	"sync"
)

// RingBuffer es ahora un buffer lineal simple (FIFO) para evitar
// errores de sincronización por sobre-ingeniería de prioridades.
type RingBuffer struct {
	capacity int
	buffer   chan []byte
	mu       sync.Mutex
	closed   bool
}

// NewRingBuffer crea un buffer simple basado en canales de Go.
func NewRingBuffer(cfg Config) *RingBuffer {
	cap := cfg.Capacity
	if cap <= 0 {
		cap = 10000
	}
	return &RingBuffer{
		capacity: cap,
		buffer:   make(chan []byte, cap),
	}
}

// Write inserta un evento en el buffer. Si está lleno, descarta el más antiguo.
func (r *RingBuffer) Write(event []byte, priority Prioridad) error {
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return errors.New("buffer cerrado")
	}
	r.mu.Unlock()

	select {
	case r.buffer <- event:
		return nil
	default:
		// Buffer lleno: Intentamos sacar uno viejo para meter el nuevo (best effort)
		select {
		case <-r.buffer:
			eventsDroppedCounter.Inc()
		default:
		}
		
		select {
		case r.buffer <- event:
			return nil
		default:
			return errors.New("buffer lleno y no se puede descartar")
		}
	}
}

// ReadBatch lee hasta 'max' eventos del buffer.
func (r *RingBuffer) ReadBatch(max int) [][]byte {
	var batch [][]byte
	for i := 0; i < max; i++ {
		select {
		case ev := <-r.buffer:
			batch = append(batch, ev)
		default:
			// No hay más mensajes por ahora
			if len(batch) == 0 {
				return nil
			}
			return batch
		}
	}
	return batch
}

// FillRatio retorna qué tan lleno está el buffer (0.0 a 1.0).
func (r *RingBuffer) FillRatio() float64 {
	return float64(len(r.buffer)) / float64(r.capacity)
}

// Close marca el buffer como cerrado.
func (r *RingBuffer) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closed = true
	// No cerramos el canal para evitar panics en escritores concurrentes,
	// simplemente dejamos que el GC lo limpie.
}
