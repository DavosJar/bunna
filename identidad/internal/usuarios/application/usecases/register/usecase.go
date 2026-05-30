package register

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type RegistrarUsuarioCasoDeUso struct {
	userRepo           usuario.UsuarioRepositorio
	credRepo           seguridad.CredencialesRepositorio
	encSvc             seguridad.EncriptacionServicio
	idGen              shareddomain.GeneradorID
	tenantRepo         tenant.TenantRepositorio
	membresiaRepo      tenant.MembresiaRepositorio
	rolRepo            rbac.RolRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
}

func NewRegistrarUsuarioCasoDeUso(
	userRepo usuario.UsuarioRepositorio,
	credRepo seguridad.CredencialesRepositorio,
	encSvc seguridad.EncriptacionServicio,
	idGen shareddomain.GeneradorID,
	tenantRepo tenant.TenantRepositorio,
	membresiaRepo tenant.MembresiaRepositorio,
	rolRepo rbac.RolRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
) *RegistrarUsuarioCasoDeUso {
	return &RegistrarUsuarioCasoDeUso{
		userRepo:             userRepo,
		credRepo:             credRepo,
		encSvc:               encSvc,
		idGen:                idGen,
		tenantRepo:           tenantRepo,
		membresiaRepo:        membresiaRepo,
		rolRepo:              rolRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
	}
}

func (uc *RegistrarUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoRegistrarUsuario) (*RespuestaRegistrarUsuario, error) {
	if err := validarComando(cmd); err != nil {
		return nil, err
	}

	usuarioID, err := uc.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID de usuario: %w", err)
	}

	nuevoUsuario, err := usuario.NuevoUsuario(usuarioID, cmd.Correo, cmd.Nombre, cmd.Apellido, cmd.Telefono)
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

	tenantID, err := uc.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID de tenant: %w", err)
	}

	nombreTenant := fmt.Sprintf("%s %s", cmd.Nombre, cmd.Apellido)
	slugTenant := generarSlug(nombreTenant)

	slugTenant = uc.asegurarSlugUnico(ctx, slugTenant, tenantID)

	nuevoTenant, err := tenant.NuevoTenant(tenantID, nombreTenant, slugTenant)
	if err != nil {
		return nil, fmt.Errorf("error al crear tenant: %w", err)
	}

	tenantCreado, err := uc.tenantRepo.Crear(ctx, nuevoTenant)
	if err != nil {
		return nil, fmt.Errorf("error al persistir tenant: %w", err)
	}

	membresia, err := tenant.NuevaMembresia(usuarioCreado.ID(), tenantCreado.ID())
	if err != nil {
		return nil, fmt.Errorf("error al crear membresía: %w", err)
	}

	if err := uc.membresiaRepo.Crear(ctx, membresia); err != nil {
		return nil, fmt.Errorf("error al persistir membresía: %w", err)
	}

	rolAdmin, err := uc.rolRepo.ObtenerPorNombre(ctx, rbac.RolAdministrador)
	if err != nil {
		return nil, fmt.Errorf("error al obtener rol administrador: %w", err)
	}

	if err := uc.usuarioTenantRolRepo.Crear(ctx, usuarioCreado.ID(), tenantCreado.ID(), rolAdmin.ID); err != nil {
		return nil, fmt.Errorf("error al asignar rol administrador: %w", err)
	}

	return &RespuestaRegistrarUsuario{
		UsuarioID: usuarioCreado.ID(),
		TenantID:  tenantCreado.ID(),
		Correo:    usuarioCreado.Correo(),
		Estado:    string(usuarioCreado.Estado()),
		CreadoEn:  usuarioCreado.FechaCreacion(),
	}, nil
}

func (uc *RegistrarUsuarioCasoDeUso) asegurarSlugUnico(ctx context.Context, slug, tenantID string) string {
	if _, err := uc.tenantRepo.ObtenerPorSlug(ctx, slug); err != nil {
		return slug
	}
	maxLen := len(tenantID)
	if maxLen > 8 {
		maxLen = 8
	}
	return fmt.Sprintf("%s-%s", slug, tenantID[:maxLen])
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

func generarSlug(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)

	replacer := strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
		"ä", "a", "ë", "e", "ï", "i", "ö", "o", "ü", "u",
		"ñ", "n", "ç", "c",
	)
	s = replacer.Replace(s)

	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "tenant"
	}
	if len(s) > 100 {
		s = s[:100]
	}
	s = strings.TrimRight(s, "-")
	return s
}
