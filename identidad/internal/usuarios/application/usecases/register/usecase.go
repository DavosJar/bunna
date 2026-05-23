package register

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type RegistrarUsuarioCasoDeUso struct {
	unitOfWork usuario.UnitOfWork
}

func NewRegistrarUsuarioCasoDeUso(unitOfWork usuario.UnitOfWork) *RegistrarUsuarioCasoDeUso {
	return &RegistrarUsuarioCasoDeUso{unitOfWork: unitOfWork}
}

func (uc *RegistrarUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoRegistrarUsuario) (*RespuestaRegistrarUsuario, error) {
	if err := validarComando(cmd); err != nil {
		return nil, err
	}

	var respuesta *RespuestaRegistrarUsuario

	err := uc.unitOfWork.Transaccional(ctx, func(tx usuario.UnitOfWork) error {
		nuevoID, err := tx.GeneradorID().NextID(ctx)
		if err != nil {
			return fmt.Errorf("error al generar ID: %w", err)
		}

		nuevoUsuario, err := usuario.NuevoUsuario(nuevoID, cmd.Correo, cmd.Nombre, cmd.Apellido, cmd.Telefono)
		if err != nil {
			return fmt.Errorf("error al crear usuario: %w", err)
		}

		usuarioCreado, err := tx.UsuarioRepository().Crear(ctx, nuevoUsuario)
		if err != nil {
			return fmt.Errorf("error al persistir usuario: %w", err)
		}

		passwordHash, err := tx.EncriptacionServicio().Hashear(cmd.Password)
		if err != nil {
			return fmt.Errorf("error al hashear password: %w", err)
		}

		nuevasCredenciales := domain.NuevaCredencialesUsuario(usuarioCreado.ID(), passwordHash)
		_, err = tx.CredencialesRepository().Crear(ctx, nuevasCredenciales)
		if err != nil {
			return fmt.Errorf("error al persistir credenciales: %w", err)
		}

		respuesta = &RespuestaRegistrarUsuario{
			UsuarioID: usuarioCreado.ID(),
			Correo:    usuarioCreado.Correo(),
			Estado:    string(usuarioCreado.Estado()),
			CreadoEn:  usuarioCreado.FechaCreacion(),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return respuesta, nil
}

func validarComando(cmd *ComandoRegistrarUsuario) error {
	if cmd.Correo == "" {
		return fmt.Errorf("correo no puede estar vacío")
	}

	if _, err := mail.ParseAddress(cmd.Correo); err != nil {
		return fmt.Errorf("formato de correo inválido: %w", err)
	}

	if cmd.Password == "" {
		return fmt.Errorf("password no puede estar vacío")
	}
	if cmd.Nombre == "" {
		return fmt.Errorf("nombre no puede estar vacío")
	}
	return nil
}
