package registro

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/mail"

	"github.com/davosjar/bunna/services/identidad/internal/domain/usuario"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServicioRegistro es el servicio de aplicación para el caso de uso de registro
type ServicioRegistro struct {
	usuarioRepo      usuario.UsuarioRepositorio
	credencialesRepo domain.CredencialesRepositorio
	db               *gorm.DB
}

// NuevoServicioRegistro crea una nueva instancia del servicio de registro
func NuevoServicioRegistro(
	usuarioRepo usuario.UsuarioRepositorio,
	credencialesRepo domain.CredencialesRepositorio,
	db *gorm.DB,
) *ServicioRegistro {
	return &ServicioRegistro{
		usuarioRepo:      usuarioRepo,
		credencialesRepo: credencialesRepo,
		db:               db,
	}
}

// Ejecutar ejecuta el caso de uso de registro
func (s *ServicioRegistro) Ejecutar(ctx context.Context, comando *ComandoRegistro) (*DtoRespuestaRegistro, error) {
	// 1. Validar comando
	if err := validarComando(comando); err != nil {
		return nil, err
	}

	var respuesta *DtoRespuestaRegistro

	// 2. Ejecutar en transacción
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 3. Crear usuario en BD (dominio usuario)
		nuevoID := uuid.New().String()
		nuevoUsuario, err := usuario.NuevoUsuario(nuevoID, comando.Correo, comando.Nombre, comando.Apellido, comando.Telefono)
		if err != nil {
			return fmt.Errorf("error al crear usuario: %w", err)
		}

		usuarioCreado, err := s.usuarioRepo.Crear(ctx, nuevoUsuario)
		if err != nil {
			return fmt.Errorf("error al persister usuario: %w", err)
		}

		// 4. Hashear password
		passwordHash := hashPassword(comando.Password)

		// 5. Crear credenciales en BD (dominio credenciales)
		nuevasCredenciales := domain.NuevaCredencialesUsuario(usuarioCreado.ID(), passwordHash)
		_, err = s.credencialesRepo.Crear(ctx, nuevasCredenciales)
		if err != nil {
			return fmt.Errorf("error al persister credenciales: %w", err)
		}

		// 7. Preparar respuesta
		respuesta = &DtoRespuestaRegistro{
			UsuarioID: usuarioCreado.ID(),
			Correo:    usuarioCreado.Correo(),
			Estado:    string(usuarioCreado.Estado()),
			Timestamp: usuarioCreado.FechaCreacion(),
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error en la transacción de registro: %w", err)
	}

	// 6. Retornar DTO con usuarioID
	return respuesta, nil
}

// validarComando valida que el comando tenga los datos obligatorios
func validarComando(comando *ComandoRegistro) error {
	if comando.Correo == "" {
		return fmt.Errorf("correo no puede estar vacío")
	}

	// Validar formato de email con net.mail.ParseAddress
	if _, err := mail.ParseAddress(comando.Correo); err != nil {
		return fmt.Errorf("formato de correo inválido: %w", err)
	}

	if comando.Password == "" {
		return fmt.Errorf("password no puede estar vacío")
	}
	if comando.Nombre == "" {
		return fmt.Errorf("nombre no puede estar vacío")
	}
	return nil
}

// hashPassword realiza hash del password usando SHA256
// NOTA: Esto es temporal, se debe cambiar a bcrypt en producción
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
