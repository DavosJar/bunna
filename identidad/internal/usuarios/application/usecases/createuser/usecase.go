package createuser

import (
	"context"
	"fmt"
	"net"
	"net/mail"
	"strings"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/davosjar/bunna/services/identidad/internal/shared/application"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type CrearUsuarioCasoDeUso struct {
	userRepo usuario.UsuarioRepositorio
	credRepo seguridad.CredencialesRepositorio
	encSvc   seguridad.EncriptacionServicio
	authSvc  rbac.AuthorizationService
	idGen    shareddomain.GeneradorID
}

func NewCrearUsuarioCasoDeUso(
	userRepo usuario.UsuarioRepositorio,
	credRepo seguridad.CredencialesRepositorio,
	encSvc seguridad.EncriptacionServicio,
	authSvc rbac.AuthorizationService,
	idGen shareddomain.GeneradorID,
) *CrearUsuarioCasoDeUso {
	return &CrearUsuarioCasoDeUso{
		userRepo: userRepo,
		credRepo: credRepo,
		encSvc:   encSvc,
		authSvc:  authSvc,
		idGen:    idGen,
	}
}

func (uc *CrearUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoCrearUsuario) (*RespuestaCrearUsuario, error) {
	if cmd.Correo == "" {
		return nil, fmt.Errorf("correo no puede estar vacío")
	}
	if _, err := mail.ParseAddress(cmd.Correo); err != nil {
		return nil, fmt.Errorf("formato de correo inválido: %w", err)
	}
	if err := emailTieneDominioValido(cmd.Correo); err != nil {
		return nil, err
	}
	if cmd.Nombre == "" {
		return nil, fmt.Errorf("nombre no puede estar vacío")
	}
	if cmd.Password == "" {
		return nil, fmt.Errorf("password no puede estar vacío")
	}
	if err := application.ValidarFormatoPassword(cmd.Password, "password"); err != nil {
		return nil, err
	}

	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, "", rbac.PermisoUsuarioCrear)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	nuevoID, err := uc.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID: %w", err)
	}

	nuevoUsuario, err := usuario.NuevoUsuario(nuevoID, cmd.Correo, cmd.Nombre, cmd.Apellido, "")
	if err != nil {
		return nil, fmt.Errorf("error al crear usuario: %w", err)
	}

	usuarioCreado, err := uc.userRepo.Crear(ctx, nuevoUsuario)
	if err != nil {
		return nil, fmt.Errorf("error al persistir usuario: %w", err)
	}

	passwordHash, err := uc.encSvc.Hashear(cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("error al hashear password: %w", err)
	}

	nuevasCredenciales := seguridad.NuevaCredencialesUsuario(usuarioCreado.ID(), passwordHash)
	_, err = uc.credRepo.Crear(ctx, nuevasCredenciales)
	if err != nil {
		return nil, fmt.Errorf("error al persistir credenciales: %w", err)
	}

	return &RespuestaCrearUsuario{
		ID:       usuarioCreado.ID(),
		Correo:   usuarioCreado.Correo(),
		Nombre:   usuarioCreado.Nombre(),
		Apellido: usuarioCreado.Apellido(),
		Activo:   true,
		CreadoEn: usuarioCreado.FechaCreacion().Format("2006-01-02T15:04:05Z"),
	}, nil
}

func emailTieneDominioValido(correo string) error {
	parts := strings.SplitN(correo, "@", 2)
	if len(parts) != 2 {
		return fmt.Errorf("formato de correo inválido")
	}
	dominio := parts[1]
	mxRecords, err := net.LookupMX(dominio)
	if err != nil || len(mxRecords) == 0 {
		return fmt.Errorf("el dominio del correo no existe o no acepta correos")
	}
	return nil
}
