package register

import (
	"context"
	"fmt"
	"net"
	"net/mail"
	"regexp"
	"strings"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/davosjar/bunna/services/identidad/internal/shared/application"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type RegistrarUsuarioCasoDeUso struct {
	uow         UnitOfWork
	encSvc      seguridad.EncriptacionServicio
	idGen       shareddomain.GeneradorID
	rolPublisher rbac.RolPublisher
}

func NewRegistrarUsuarioCasoDeUso(
	uow UnitOfWork,
	encSvc seguridad.EncriptacionServicio,
	idGen shareddomain.GeneradorID,
	rolPublisher rbac.RolPublisher,
) *RegistrarUsuarioCasoDeUso {
	return &RegistrarUsuarioCasoDeUso{
		uow:         uow,
		encSvc:      encSvc,
		idGen:       idGen,
		rolPublisher: rolPublisher,
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

	tenantID, err := uc.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID de tenant: %w", err)
	}

	passwordHash, err := uc.encSvc.Hashear(cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("error al hashear password: %w", err)
	}

	nombreTenant := fmt.Sprintf("%s %s", cmd.Nombre, cmd.Apellido)
	slugTenant := generarSlug(nombreTenant)

	var respuesta *RespuestaRegistrarUsuario

	err = uc.uow.Transaccional(ctx, func(tx UnitOfWork) error {
		nuevoUsuario, err := usuario.NuevoUsuario(usuarioID, cmd.Correo, cmd.Nombre, cmd.Apellido, cmd.Telefono)
		if err != nil {
			return fmt.Errorf("error al crear usuario: %w", err)
		}

		usuarioCreado, err := tx.UsuarioRepository().Crear(ctx, nuevoUsuario)
		if err != nil {
			return fmt.Errorf("error al persistir usuario: %w", err)
		}

		nuevasCredenciales := seguridad.NuevaCredencialesUsuario(usuarioCreado.ID(), passwordHash)
		if _, err = tx.CredencialesRepository().Crear(ctx, nuevasCredenciales); err != nil {
			return fmt.Errorf("error al persistir credenciales: %w", err)
		}

		slugUnico := uc.asegurarSlugUnico(ctx, slugTenant, tenantID, tx)

		nuevoTenant, err := tenant.NuevoTenant(tenantID, nombreTenant, slugUnico)
		if err != nil {
			return fmt.Errorf("error al crear tenant: %w", err)
		}

		tenantCreado, err := tx.TenantRepository().Crear(ctx, nuevoTenant)
		if err != nil {
			return fmt.Errorf("error al persistir tenant: %w", err)
		}

		membresia, err := tenant.NuevaMembresia(usuarioCreado.ID(), tenantCreado.ID())
		if err != nil {
			return fmt.Errorf("error al crear membresía: %w", err)
		}

		if err := tx.MembresiaRepository().Crear(ctx, membresia); err != nil {
			return fmt.Errorf("error al persistir membresía: %w", err)
		}

		rolAdmin, err := tx.RolRepository().ObtenerPorNombre(ctx, rbac.RolAdministrador)
		if err != nil {
			return fmt.Errorf("error al obtener rol administrador: %w", err)
		}

		if err := tx.UsuarioTenantRolRepository().Crear(ctx, usuarioCreado.ID(), tenantCreado.ID(), rolAdmin.ID); err != nil {
			return fmt.Errorf("error al asignar rol administrador: %w", err)
		}

		permisosAdmin, err := tx.RolPermisoRepository().ListarPorRolYTenant(ctx, rolAdmin.ID, rbac.TenantIDSistema)
		if err != nil {
			return fmt.Errorf("error al obtener permisos base del administrador: %w", err)
		}

		var codigosPermisos []string
		for _, p := range permisosAdmin {
			codigosPermisos = append(codigosPermisos, p.Codigo)
		}

		respuesta = &RespuestaRegistrarUsuario{
			UsuarioID: usuarioCreado.ID(),
			TenantID:  tenantCreado.ID(),
			Correo:    usuarioCreado.Correo(),
			Estado:    string(usuarioCreado.Estado()),
			CreadoEn:  usuarioCreado.FechaCreacion(),
		}

		if err := uc.rolPublisher.PublicarRolActualizado(ctx, rbac.RolAdministrador, tenantCreado.ID(), codigosPermisos); err != nil {
			return fmt.Errorf("error al publicar rol administrador para el nuevo tenant: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return respuesta, nil
}

func (uc *RegistrarUsuarioCasoDeUso) asegurarSlugUnico(ctx context.Context, slug, tenantID string, tx UnitOfWork) string {
	if _, err := tx.TenantRepository().ObtenerPorSlug(ctx, slug); err != nil {
		return slug
	}
	maxLen := len(tenantID)
	if maxLen > 8 {
		maxLen = 8
	}
	return fmt.Sprintf("%s-%s", slug, tenantID[:maxLen])
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

func validarComando(cmd *ComandoRegistrarUsuario) error {
	if cmd.Correo == "" {
		return fmt.Errorf("correo no puede estar vacío")
	}

	if _, err := mail.ParseAddress(cmd.Correo); err != nil {
		return fmt.Errorf("formato de correo inválido: %w", err)
	}

	if err := emailTieneDominioValido(cmd.Correo); err != nil {
		return err
	}

	if cmd.Password == "" {
		return fmt.Errorf("password no puede estar vacío")
	}
	if err := application.ValidarFormatoPassword(cmd.Password, "password"); err != nil {
		return err
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
