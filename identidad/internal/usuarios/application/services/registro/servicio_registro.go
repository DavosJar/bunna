package registro

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

// ServicioRegistro es el servicio de aplicación para el caso de uso de registro
type ServicioRegistro struct {
	unitOfWork usuario.UnitOfWork
}

// NuevoServicioRegistro crea una nueva instancia del servicio de registro
func NuevoServicioRegistro(unitOfWork usuario.UnitOfWork) *ServicioRegistro {
	return &ServicioRegistro{
		unitOfWork: unitOfWork,
	}
}

// Ejecutar ejecuta el caso de uso de registro dentro de una transacción
func (s *ServicioRegistro) Ejecutar(ctx context.Context, comando *ComandoRegistro) (*DtoRespuestaRegistro, error) {
	// 1. Validar comando (ANTES de transacción - costo bajo)
	if err := validarComando(comando); err != nil {
		return nil, err
	}

	var respuesta *DtoRespuestaRegistro

	err := s.unitOfWork.Transaccional(ctx, func(tx usuario.UnitOfWork) error {
		// 2a. Generar nuevo ID
		nuevoID, err := tx.GeneradorID().NextID(ctx)
		if err != nil {
			return fmt.Errorf("error al generar ID: %w", err)
		}

		// 2b. Crear usuario
		nuevoUsuario, err := usuario.NuevoUsuario(nuevoID, comando.Correo, comando.Nombre, comando.Apellido, comando.Telefono)
		if err != nil {
			return fmt.Errorf("error al crear usuario: %w", err)
		}

		usuarioCreado, err := tx.UsuarioRepository().Crear(ctx, nuevoUsuario)
		if err != nil {
			return fmt.Errorf("error al persister usuario: %w", err)
		}

		// 2c. Hashear password
		passwordHash, err := tx.EncriptacionServicio().Hashear(comando.Password)
		if err != nil {
			return fmt.Errorf("error al hashear password: %w", err)
		}

		// 2d. Crear credenciales
		nuevasCredenciales := domain.NuevaCredencialesUsuario(usuarioCreado.ID(), passwordHash)
		_, err = tx.CredencialesRepository().Crear(ctx, nuevasCredenciales)
		if err != nil {
			return fmt.Errorf("error al persister credenciales: %w", err)
		}

		// 2e. Preparar respuesta
		respuesta = &DtoRespuestaRegistro{
			UsuarioID: usuarioCreado.ID(),
			Correo:    usuarioCreado.Correo(),
			Estado:    string(usuarioCreado.Estado()),
			Timestamp: usuarioCreado.FechaCreacion(),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return respuesta, nil
}

// validarComando valida que el comando tenga los datos obligatorios
func validarComando(comando *ComandoRegistro) error {
	if comando.Correo == "" {
		return fmt.Errorf("correo no puede estar vacío")
	}

	// Validar formato de email
	if _, err := mail.ParseAddress(comando.Correo); err != nil {
		return fmt.Errorf("formato de correo inválido: %w", err)
	}

	// Validar que el dominio del correo exista via DNS
	if err := validarDominioCorreo(comando.Correo); err != nil {
		return fmt.Errorf("dominio de correo inválido: %w", err)
	}

	if comando.Password == "" {
		return fmt.Errorf("password no puede estar vacío")
	}

	if comando.Nombre == "" {
		return fmt.Errorf("nombre no puede estar vacío")
	}

	return nil
}
