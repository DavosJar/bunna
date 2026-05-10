package domain

import "time"

type EstadoSesion string

const (
	EstadoActiva   EstadoSesion = "ACTIVA"
	EstadoExpirada EstadoSesion = "EXPIRADA"
	EstadoRevocada EstadoSesion = "REVOCADA"
)

type Sesion struct {
	id                    string
	usuarioID             string
	accessTokenHash       string
	refreshTokenHash      string
	estado                EstadoSesion
	ipOrigen              string
	fechaCreacion         time.Time
	fechaActualizacion    time.Time
	fechaExpiracionAccess  time.Time
	fechaExpiracionRefresh time.Time
	ultimaActividad       time.Time
	contadorRefrescos     int
}

// NuevaSesion crea una sesión nueva con validaciones.
// Los tokens recibidos son HASHES, no tokens en plano.
func NuevaSesion(
	id string,
	usuarioID string,
	accessTokenHash string,
	refreshTokenHash string,
	ipOrigen string,
	fechaCreacion time.Time,
	fechaExpiracionAccess time.Time,
	fechaExpiracionRefresh time.Time,
) (*Sesion, error) {
	if usuarioID == "" {
		return nil, ErrUsuarioIDRequerido
	}
	if accessTokenHash == "" {
		return nil, ErrAccessTokenHashRequerido
	}
	if refreshTokenHash == "" {
		return nil, ErrRefreshTokenHashRequerido
	}
	if !fechaExpiracionAccess.After(fechaCreacion) {
		return nil, ErrFechaExpiracionInvalida
	}

	return &Sesion{
		id:                    id,
		usuarioID:             usuarioID,
		accessTokenHash:       accessTokenHash,
		refreshTokenHash:      refreshTokenHash,
		estado:                EstadoActiva,
		ipOrigen:              ipOrigen,
		fechaCreacion:         fechaCreacion,
		fechaActualizacion:    fechaCreacion,
		fechaExpiracionAccess:  fechaExpiracionAccess,
		fechaExpiracionRefresh: fechaExpiracionRefresh,
		ultimaActividad:       fechaCreacion,
		contadorRefrescos:     0,
	}, nil
}

// NuevaSesionDesdeBD reconstruye la entidad desde persistencia sin validaciones.
func NuevaSesionDesdeBD(
	id string,
	usuarioID string,
	accessTokenHash string,
	refreshTokenHash string,
	estado EstadoSesion,
	ipOrigen string,
	fechaCreacion time.Time,
	fechaActualizacion time.Time,
	fechaExpiracionAccess time.Time,
	fechaExpiracionRefresh time.Time,
	ultimaActividad time.Time,
	contadorRefrescos int,
) *Sesion {
	return &Sesion{
		id:                    id,
		usuarioID:             usuarioID,
		accessTokenHash:       accessTokenHash,
		refreshTokenHash:      refreshTokenHash,
		estado:                estado,
		ipOrigen:              ipOrigen,
		fechaCreacion:         fechaCreacion,
		fechaActualizacion:    fechaActualizacion,
		fechaExpiracionAccess:  fechaExpiracionAccess,
		fechaExpiracionRefresh: fechaExpiracionRefresh,
		ultimaActividad:       ultimaActividad,
		contadorRefrescos:     contadorRefrescos,
	}
}

// --- Comportamiento ---

// EstaActiva retorna true solo si el estado es ACTIVA y el access token no ha expirado.
// La referencia de tiempo viene de afuera: el dominio no llama a time.Now().
func (s *Sesion) EstaActiva(ahora time.Time) bool {
	if s.estado != EstadoActiva {
		return false
	}
	return ahora.Before(s.fechaExpiracionAccess)
}

// RefreshTokenValido retorna true si la sesión está ACTIVA y el refresh no ha expirado.
func (s *Sesion) RefreshTokenValido(ahora time.Time) bool {
	if s.estado != EstadoActiva {
		return false
	}
	if s.fechaExpiracionRefresh.IsZero() {
		return false
	}
	return ahora.Before(s.fechaExpiracionRefresh)
}

// MarcarExpirada transiciona a EXPIRADA.
// No está permitido si ya está REVOCADA.
func (s *Sesion) MarcarExpirada() error {
	if s.estado == EstadoRevocada {
		return ErrTransicionEstadoInvalida
	}
	s.estado = EstadoExpirada
	return nil
}

// Revocar transiciona a REVOCADA desde cualquier estado.
func (s *Sesion) Revocar() {
	s.estado = EstadoRevocada
}

// RegistrarActividad actualiza la última actividad de la sesión.
func (s *Sesion) RegistrarActividad(ahora time.Time) {
	s.ultimaActividad = ahora
	s.fechaActualizacion = ahora
}

// TimeoutExcedido retorna true si el tiempo sin actividad supera el timeout configurado.
func (s *Sesion) TimeoutExcedido(ahora time.Time, timeout time.Duration) bool {
	return ahora.After(s.ultimaActividad.Add(timeout))
}

// RotarTokens reemplaza los hashes y actualiza expiraciones. Incrementa el contador.
func (s *Sesion) RotarTokens(
	nuevoAccessHash string,
	nuevoRefreshHash string,
	nuevaExpiracionAccess time.Time,
	nuevaExpiracionRefresh time.Time,
	ahora time.Time,
) {
	s.accessTokenHash = nuevoAccessHash
	s.refreshTokenHash = nuevoRefreshHash
	s.fechaExpiracionAccess = nuevaExpiracionAccess
	s.fechaExpiracionRefresh = nuevaExpiracionRefresh
	s.contadorRefrescos++
	s.fechaActualizacion = ahora
	s.ultimaActividad = ahora
}

// --- Getters ---

func (s *Sesion) ID() string                         { return s.id }
func (s *Sesion) UsuarioID() string                  { return s.usuarioID }
func (s *Sesion) AccessTokenHash() string            { return s.accessTokenHash }
func (s *Sesion) RefreshTokenHash() string           { return s.refreshTokenHash }
func (s *Sesion) Estado() EstadoSesion               { return s.estado }
func (s *Sesion) IPOrigen() string                   { return s.ipOrigen }
func (s *Sesion) FechaCreacion() time.Time           { return s.fechaCreacion }
func (s *Sesion) FechaActualizacion() time.Time      { return s.fechaActualizacion }
func (s *Sesion) FechaExpiracionAccess() time.Time   { return s.fechaExpiracionAccess }
func (s *Sesion) FechaExpiracionRefresh() time.Time  { return s.fechaExpiracionRefresh }
func (s *Sesion) UltimaActividad() time.Time         { return s.ultimaActividad }
func (s *Sesion) ContadorRefrescos() int             { return s.contadorRefrescos }