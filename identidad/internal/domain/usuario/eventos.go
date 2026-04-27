package usuario

import (
	"time"
)

// EventoDominio representa algo que sucedió en el núcleo
type EventoDominio struct {
	Nombre   string
	Payload  interface{}
	Ocurrido time.Time
}

// EventosUsuario maneja todos los eventos del usuario
// Expone métodos específicos según el caso de uso
type EventosUsuario struct {
	eventos []EventoDominio
}

// NuevosEventosUsuario crea una nueva instancia de gestor de eventos
func NuevosEventosUsuario() *EventosUsuario {
	return &EventosUsuario{
		eventos: make([]EventoDominio, 0),
	}
}

// RegistrarCreacion registra cuando un usuario es creado
func (e *EventosUsuario) RegistrarCreacion(id, correo string) {
	e.registrarEvento("UsuarioCreado", map[string]string{
		"id":     id,
		"correo": correo,
	})
}

// RegistrarCambioEstado registra cuando cambia el estado del usuario
func (e *EventosUsuario) RegistrarCambioEstado(id string, nuevoEstado EstadoUsuario) {
	e.registrarEvento("EstadoUsuarioCambiado", map[string]interface{}{
		"id":    id,
		"nuevo": nuevoEstado,
	})
}

func (e *EventosUsuario) RegistrarBloqueo(id string) {
	e.registrarEvento("UsuarioBloqueado", map[string]string{"id": id})
}

func (e *EventosUsuario) RegistrarActivacion(id string) {
	e.registrarEvento("UsuarioActivado", map[string]string{"id": id})
}

func (e *EventosUsuario) RegistrarInactivacion(id string) {
	e.registrarEvento("UsuarioInactivado", map[string]string{"id": id})
}

func (e *EventosUsuario) RegistrarVerificacion(id string) {
	e.registrarEvento("CorreoVerificado", map[string]string{"id": id})
}

func (e *EventosUsuario) Extraer() []EventoDominio {
	eventos := e.eventos
	e.eventos = make([]EventoDominio, 0)
	return eventos
}

func (e *EventosUsuario) Cantidad() int {
	return len(e.eventos)
}

// RegistrarEvento permite que capas externas registren eventos personalizados
func (e *EventosUsuario) RegistrarEvento(nombre string, payload interface{}) {
	e.registrarEvento(nombre, payload)
}

func (e *EventosUsuario) registrarEvento(nombre string, payload interface{}) {
	e.eventos = append(e.eventos, EventoDominio{
		Nombre:   nombre,
		Payload:  payload,
		Ocurrido: time.Now(),
	})
}
