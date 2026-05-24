package domain

import "time"

type EstadoSesion string

const (
	EstadoActiva   EstadoSesion = "ACTIVA"
	EstadoExpirada EstadoSesion = "EXPIRADA"
	EstadoRevocada EstadoSesion = "REVOCADA"
)

type Sesion struct {
	id                     string
	usuarioID              string
	accessTokenHash        string
	refreshTokenHash       string
	estado                 EstadoSesion
	ipOrigen               string
	fechaCreacion          time.Time
	fechaActualizacion     time.Time
	fechaExpiracionAccess  time.Time
	fechaExpiracionRefresh time.Time
	ultimaActividad        time.Time
	contadorRefrescos      int
}

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
		id:                     id,
		usuarioID:              usuarioID,
		accessTokenHash:        accessTokenHash,
		refreshTokenHash:       refreshTokenHash,
		estado:                 EstadoActiva,
		ipOrigen:               ipOrigen,
		fechaCreacion:          fechaCreacion,
		fechaActualizacion:     fechaCreacion,
		fechaExpiracionAccess:  fechaExpiracionAccess,
		fechaExpiracionRefresh: fechaExpiracionRefresh,
		ultimaActividad:        fechaCreacion,
		contadorRefrescos:      0,
	}, nil
}

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
		id:                     id,
		usuarioID:              usuarioID,
		accessTokenHash:        accessTokenHash,
		refreshTokenHash:       refreshTokenHash,
		estado:                 estado,
		ipOrigen:               ipOrigen,
		fechaCreacion:          fechaCreacion,
		fechaActualizacion:     fechaActualizacion,
		fechaExpiracionAccess:  fechaExpiracionAccess,
		fechaExpiracionRefresh: fechaExpiracionRefresh,
		ultimaActividad:        ultimaActividad,
		contadorRefrescos:      contadorRefrescos,
	}
}

// --- Comportamiento ---

func (s *Sesion) EstaActiva(ahora time.Time) bool {
	if s.estado != EstadoActiva {
		return false
	}
	return ahora.Before(s.fechaExpiracionAccess)
}

func (s *Sesion) RefreshTokenValido(ahora time.Time) bool {
	if s.estado != EstadoActiva {
		return false
	}
	if s.fechaExpiracionRefresh.IsZero() {
		return false
	}
	return ahora.Before(s.fechaExpiracionRefresh)
}

// MarcarExpirada: no permitido desde REVOCADA ni desde EXPIRADA.
func (s *Sesion) MarcarExpirada() error {
	if s.estado == EstadoRevocada || s.estado == EstadoExpirada {
		return ErrTransicionEstadoInvalida
	}
	s.estado = EstadoExpirada
	return nil
}

// Revocar: permitido desde cualquier estado (idempotente si ya está revocada).
func (s *Sesion) Revocar() {
	s.estado = EstadoRevocada
}

func (s *Sesion) RegistrarActividad(ahora time.Time) {
	s.ultimaActividad = ahora
	s.fechaActualizacion = ahora
}

func (s *Sesion) TimeoutExcedido(ahora time.Time, timeout time.Duration) bool {
	return ahora.After(s.ultimaActividad.Add(timeout))
}

// RotarTokens: ahora retorna error si la sesión no está ACTIVA.
func (s *Sesion) RotarTokens(
	nuevoAccessHash string,
	nuevoRefreshHash string,
	nuevaExpiracionAccess time.Time,
	nuevaExpiracionRefresh time.Time,
	ahora time.Time,
) error {
	if s.estado != EstadoActiva {
		return ErrTransicionEstadoInvalida
	}
	s.accessTokenHash = nuevoAccessHash
	s.refreshTokenHash = nuevoRefreshHash
	s.fechaExpiracionAccess = nuevaExpiracionAccess
	s.fechaExpiracionRefresh = nuevaExpiracionRefresh
	s.contadorRefrescos++
	s.fechaActualizacion = ahora
	s.ultimaActividad = ahora
	return nil
}

// --- Getters ---

func (s *Sesion) ID() string                        { return s.id }
func (s *Sesion) UsuarioID() string                 { return s.usuarioID }
func (s *Sesion) AccessTokenHash() string           { return s.accessTokenHash }
func (s *Sesion) RefreshTokenHash() string          { return s.refreshTokenHash }
func (s *Sesion) Estado() EstadoSesion              { return s.estado }
func (s *Sesion) IPOrigen() string                  { return s.ipOrigen }
func (s *Sesion) FechaCreacion() time.Time          { return s.fechaCreacion }
func (s *Sesion) FechaActualizacion() time.Time     { return s.fechaActualizacion }
func (s *Sesion) FechaExpiracionAccess() time.Time  { return s.fechaExpiracionAccess }
func (s *Sesion) FechaExpiracionRefresh() time.Time { return s.fechaExpiracionRefresh }
func (s *Sesion) UltimaActividad() time.Time        { return s.ultimaActividad }
func (s *Sesion) ContadorRefrescos() int            { return s.contadorRefrescos }
