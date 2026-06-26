package buffer

import "io"

// noopWriter is a BufferWriter implementation that discards all events.
type noopWriter struct{}

// NewNoopWriter returns a BufferWriter that does nothing with the events.
func NewNoopWriter() BufferWriter {
    return &noopWriter{}
}

// Write implements BufferWriter but discards the event.
func (n *noopWriter) Write(event []byte, priority Prioridad) error {
    // Discard the event; optionally could write to io.Discard to avoid unused variable.
    _, _ = io.Discard.Write(event)
    return nil
}
