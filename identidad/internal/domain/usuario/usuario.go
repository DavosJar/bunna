package usuario

import (
	"time"
)

// Errores de Dominio (Exportados para que la capa de Infraestructura los mapee a HTTP/gRPC)
type EstadoUsuario string

const (
	NO_VERIFICADO            EstadoUsuario = "NO_VERIFICADO"
	ACTIVO                   EstadoUsuario = "ACTIVO"
	INACTIVO                 EstadoUsuario = "INACTIVO"
	PENDIENTE_DE_ELIMINACION EstadoUsuario = "PENDIENTE_DE_ELIMINACION"
	BLOQUEADO                EstadoUsuario = "BLOQUEADO"
)

// transiciones encapsula la máquina de estados
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

// Usuario es la Entidad Raíz (Aggregate Root)
type Usuario struct {
	id                       string
	nombre                   string
	apellido                 string
	correo                   string
	telefono                 string
	estado                   EstadoUsuario
	estadoVerificacionCorreo EstadoVerificacionCorreo
	fechaCreacion            time.Time
	fechaActualizacion       time.Time
	eventos                  *EventosUsuario // Para consistencia eventual
}

// NuevoUsuario crea un usuario a partir de los datos proporcionados.
// id puede estar vacío.
func NuevoUsuario(id, correo, nombre, apellido, telefono string) (*Usuario, error) {
	ahora := time.Now()
	u := &Usuario{
		id:                       id,
		nombre:                   nombre,
		apellido:                 apellido,
		correo:                   correo,
		telefono:                 telefono,
		estado:                   NO_VERIFICADO,
		estadoVerificacionCorreo: PENDIENTE_VERIFICACION,
		fechaCreacion:            ahora,
		fechaActualizacion:       ahora,
		eventos:                  NuevosEventosUsuario(),
	}

	if err := u.validar(); err != nil {
		return nil, err
	}

	// Registrar evento de creación
	u.eventos.RegistrarCreacion(u.id, u.correo)
	return u, nil

}

// Reglas de validación internas
func (u *Usuario) validar() error {
	if u.correo == "" {
		return ErrCorreoRequerido
	}
	return nil
}

// Getters (Solo lectura desde el exterior)
func (u *Usuario) ID() string            { return u.id }
func (u *Usuario) Nombre() string        { return u.nombre }
func (u *Usuario) Apellido() string      { return u.apellido }
func (u *Usuario) Correo() string        { return u.correo }
func (u *Usuario) Telefono() string      { return u.telefono }
func (u *Usuario) Estado() EstadoUsuario { return u.estado }
func (u *Usuario) EstadoVerificacionCorreo() EstadoVerificacionCorreo {
	return u.estadoVerificacionCorreo
}
func (u *Usuario) FechaCreacion() time.Time      { return u.fechaCreacion }
func (u *Usuario) FechaActualizacion() time.Time { return u.fechaActualizacion }

// CambiarEstado centraliza la lógica de transiciones
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

// Métodos de conveniencia (Ubiquitous Language)
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

// Métodos de verificación de correo
func (u *Usuario) VerificarCorreo() error {
	if !u.estadoVerificacionCorreo.PuedeTransicionarA(VERIFICADO) {
		return ErrTransicionVerificacionNoPermitida
	}
	u.estadoVerificacionCorreo = VERIFICADO
	u.eventos.RegistrarVerificacion(u.id)
	return nil
}

func (u *Usuario) SolicitarReenvioVerificacion() error {
	if !u.estadoVerificacionCorreo.PuedeTransicionarA(REENVIO_SOLICITADO) {
		return ErrTransicionVerificacionNoPermitida
	}
	u.estadoVerificacionCorreo = REENVIO_SOLICITADO
	u.eventos.registrarEvento("ReenvioVerificacionSolicitado", map[string]string{"id": u.id})
	return nil
}

func (u *Usuario) MarcarEnlaceExpirado() error {
	if !u.estadoVerificacionCorreo.PuedeTransicionarA(ENLACE_EXPIRADO) {
		return ErrTransicionVerificacionNoPermitida
	}
	u.estadoVerificacionCorreo = ENLACE_EXPIRADO
	u.eventos.registrarEvento("EnlaceVerificacionExpirado", map[string]string{"id": u.id})
	return nil
}

func (u *Usuario) MarcarVerificacionFallida() error {
	if !u.estadoVerificacionCorreo.PuedeTransicionarA(VERIFICACION_FALLIDA) {
		return ErrTransicionVerificacionNoPermitida
	}
	u.estadoVerificacionCorreo = VERIFICACION_FALLIDA
	u.eventos.registrarEvento("VerificacionFallida", map[string]string{"id": u.id})
	return nil
}

// NewUsuarioFromPersistence reconstruye un Usuario desde la base de datos.
// No valida ni emite eventos - se asume que los datos provienen de una fuente confiable.
func NewUsuarioFromPersistence(id, nombre, apellido, correo, telefono string, estado EstadoUsuario, estadoVerificacion EstadoVerificacionCorreo, fechaCreacion, fechaActualizacion time.Time) *Usuario {
	return &Usuario{
		id:                       id,
		nombre:                   nombre,
		apellido:                 apellido,
		correo:                   correo,
		telefono:                 telefono,
		estado:                   estado,
		estadoVerificacionCorreo: estadoVerificacion,
		fechaCreacion:            fechaCreacion,
		fechaActualizacion:       fechaActualizacion,
		eventos:                  NuevosEventosUsuario(),
	}
}

// PullEventos retorna todos los eventos pendientes y los limpia
func (u *Usuario) PullEventos() []EventoDominio {
	return u.eventos.Extraer()
}

// Eventos retorna la instancia de eventos para acceso directo si es necesario
func (u *Usuario) Eventos() *EventosUsuario {
	return u.eventos
}
