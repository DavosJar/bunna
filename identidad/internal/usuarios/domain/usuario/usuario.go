package usuario

import (
	"time"
)

type EstadoUsuario string

const (
	NO_VERIFICADO            EstadoUsuario = "NO_VERIFICADO"
	ACTIVO                   EstadoUsuario = "ACTIVO"
	INACTIVO                 EstadoUsuario = "INACTIVO"
	PENDIENTE_DE_ELIMINACION EstadoUsuario = "PENDIENTE_DE_ELIMINACION"
	BLOQUEADO                EstadoUsuario = "BLOQUEADO"
)

var transiciones = map[EstadoUsuario]map[EstadoUsuario]bool{
	NO_VERIFICADO: {
		ACTIVO:                   true,
		PENDIENTE_DE_ELIMINACION: true,
	},
	ACTIVO: {
		INACTIVO:                 true,
		BLOQUEADO:                true,
		PENDIENTE_DE_ELIMINACION: true,
	},
	INACTIVO: {
		ACTIVO:                   true,
		PENDIENTE_DE_ELIMINACION: true,
	},
	BLOQUEADO: {
		ACTIVO:                   true,
		PENDIENTE_DE_ELIMINACION: true,
	},
	PENDIENTE_DE_ELIMINACION: {},
}

type Usuario struct {
	id                 string
	nombre             string
	apellido           string
	correoElectronico  *CorreoElectronico
	telefono           string
	estado             EstadoUsuario
	fechaCreacion      time.Time
	fechaActualizacion time.Time
	eventos            *EventosUsuario
}

func NuevoUsuario(id, correo, nombre, apellido, telefono string) (*Usuario, error) {
	correoVO, err := NuevoCorreoElectronico(correo)
	if err != nil {
		return nil, err
	}

	ahora := time.Now()
	u := &Usuario{
		id:                 id,
		nombre:             nombre,
		apellido:           apellido,
		correoElectronico:  correoVO,
		telefono:           telefono,
		estado:             NO_VERIFICADO,
		fechaCreacion:      ahora,
		fechaActualizacion: ahora,
		eventos:            NuevosEventosUsuario(),
	}

	u.eventos.RegistrarCreacion(u.id, u.correoElectronico.Direccion())
	return u, nil
}

func NewUsuarioFromPersistence(id, nombre, apellido string, correoElectronico *CorreoElectronico, telefono string, estado EstadoUsuario, fechaCreacion, fechaActualizacion time.Time) *Usuario {
	return &Usuario{
		id:                 id,
		nombre:             nombre,
		apellido:           apellido,
		correoElectronico:  correoElectronico,
		telefono:           telefono,
		estado:             estado,
		fechaCreacion:      fechaCreacion,
		fechaActualizacion: fechaActualizacion,
		eventos:            NuevosEventosUsuario(),
	}
}

// Getters
func (u *Usuario) ID() string            { return u.id }
func (u *Usuario) Nombre() string        { return u.nombre }
func (u *Usuario) Apellido() string      { return u.apellido }
func (u *Usuario) Correo() string        { return u.correoElectronico.Direccion() }
func (u *Usuario) Telefono() string      { return u.telefono }
func (u *Usuario) Estado() EstadoUsuario { return u.estado }
func (u *Usuario) EstadoVerificacionCorreo() EstadoVerificacionCorreo {
	return u.correoElectronico.Estado()
}
func (u *Usuario) FechaCreacion() time.Time      { return u.fechaCreacion }
func (u *Usuario) FechaActualizacion() time.Time { return u.fechaActualizacion }

func (u *Usuario) CambiarEstado(siguiente EstadoUsuario) error {
	if u.estado == siguiente {
		return nil
	}
	if !transiciones[u.estado][siguiente] {
		return ErrTransicionNoPermitida
	}
	u.estado = siguiente
	u.eventos.RegistrarCambioEstado(u.id, siguiente)
	return nil
}

func (u *Usuario) Bloquear() error {
	if err := u.CambiarEstado(BLOQUEADO); err != nil {
		return err
	}
	u.eventos.RegistrarBloqueo(u.id)
	return nil
}

func (u *Usuario) Activar() error {
	if err := u.CambiarEstado(ACTIVO); err != nil {
		return err
	}
	u.eventos.RegistrarActivacion(u.id)
	return nil
}

func (u *Usuario) Inactivar() error {
	if err := u.CambiarEstado(INACTIVO); err != nil {
		return err
	}
	u.eventos.RegistrarInactivacion(u.id)
	return nil
}

// ActualizarDatosPersonales actualiza nombre y apellido del usuario.
// No genera evento de dominio porque es una actualización de datos, no de estado.
func (u *Usuario) ActualizarDatosPersonales(nombre, apellido string) {
	u.nombre = nombre
	u.apellido = apellido
	u.fechaActualizacion = time.Now()
}

func (u *Usuario) VerificarCorreo() error {
	if err := u.correoElectronico.Verificar(); err != nil {
		return ErrTransicionVerificacionNoPermitida
	}
	u.eventos.RegistrarVerificacion(u.id)
	return nil
}

func (u *Usuario) SolicitarReenvioVerificacion() error {
	if err := u.correoElectronico.SolicitarReenvio(); err != nil {
		return ErrTransicionVerificacionNoPermitida
	}
	u.eventos.registrarEvento("ReenvioVerificacionSolicitado", map[string]string{"id": u.id})
	return nil
}

func (u *Usuario) MarcarEnlaceExpirado() error {
	if err := u.correoElectronico.MarcarExpirado(); err != nil {
		return ErrTransicionVerificacionNoPermitida
	}
	u.eventos.registrarEvento("EnlaceVerificacionExpirado", map[string]string{"id": u.id})
	return nil
}

func (u *Usuario) PullEventos() []EventoDominio {
	return u.eventos.Extraer()
}

func (u *Usuario) Eventos() *EventosUsuario {
	return u.eventos
}
