// Package domain contiene las entidades y contratos del dominio de seguridad.
package domain

import "time"

// IntentoPorIP registra los intentos de login fallidos desde una dirección IP.
// Se usa para implementar bloqueo por IP independiente del bloqueo por usuario.
type IntentoPorIP struct {
	id           string
	ip           string
	contador     int
	ventanaInicio time.Time
	bloqueadoHasta time.Time
}

// NuevoIntentoPorIP crea un nuevo registro de intento para una IP.
func NuevoIntentoPorIP(id, ip string, ahora time.Time) *IntentoPorIP {
	return &IntentoPorIP{
		id:            id,
		ip:            ip,
		contador:      1,
		ventanaInicio: ahora,
		bloqueadoHasta: time.Time{},
	}
}

// NuevoIntentoPorIPDesdeBD reconstruye la entidad desde persistencia sin validaciones.
func NuevoIntentoPorIPDesdeBD(id, ip string, contador int, ventanaInicio, bloqueadoHasta time.Time) *IntentoPorIP {
	return &IntentoPorIP{
		id:             id,
		ip:             ip,
		contador:       contador,
		ventanaInicio:  ventanaInicio,
		bloqueadoHasta: bloqueadoHasta,
	}
}

// IncrementarContador suma un intento fallido al contador.
func (i *IntentoPorIP) IncrementarContador() {
	i.contador++
}

// Bloquear establece el tiempo de bloqueo para esta IP.
func (i *IntentoPorIP) Bloquear(hasta time.Time) {
	i.bloqueadoHasta = hasta
}

// EstaBloqueada retorna true si la IP está bloqueada en el momento dado.
func (i *IntentoPorIP) EstaBloqueada(ahora time.Time) bool {
	return ahora.Before(i.bloqueadoHasta)
}

// VentanaExpirada retorna true si la ventana de tiempo para contar intentos ya venció.
func (i *IntentoPorIP) VentanaExpirada(ahora time.Time, ventana time.Duration) bool {
	return ahora.After(i.ventanaInicio.Add(ventana))
}

// ID retorna el identificador del registro.
func (i *IntentoPorIP) ID() string { return i.id }

// IP retorna la dirección IP.
func (i *IntentoPorIP) IP() string { return i.ip }

// Contador retorna el número de intentos fallidos en la ventana actual.
func (i *IntentoPorIP) Contador() int { return i.contador }

// VentanaInicio retorna el inicio de la ventana de conteo.
func (i *IntentoPorIP) VentanaInicio() time.Time { return i.ventanaInicio }

// BloqueadoHasta retorna la fecha hasta la que está bloqueada la IP.
func (i *IntentoPorIP) BloqueadoHasta() time.Time { return i.bloqueadoHasta }