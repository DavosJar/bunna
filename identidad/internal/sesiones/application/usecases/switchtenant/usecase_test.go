package switchtenant_test

import (
	"context"
	"errors"
	"testing"
	"time"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/switchtenant"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type mockMembresiaRepo struct {
	existeMiembro bool
	errExiste     error
}

func (m *mockMembresiaRepo) Crear(ctx context.Context, memb *tenant.Membresia) error { return nil }
func (m *mockMembresiaRepo) Eliminar(ctx context.Context, uid, tid string) error { return nil }
func (m *mockMembresiaRepo) ExisteMiembro(ctx context.Context, uid, tid string) (bool, error) {
	return m.existeMiembro, m.errExiste
}
func (m *mockMembresiaRepo) ListarUsuariosPorTenant(ctx context.Context, tid string) ([]string, error) { return nil, nil }
func (m *mockMembresiaRepo) ListarTenantsPorUsuario(ctx context.Context, uid string) ([]string, error) { return nil, nil }

type mockUTRR struct {
	roles []*rbac.RolDB
	err   error
}

func (m *mockUTRR) Crear(ctx context.Context, uid, tid, rid string) error { return nil }
func (m *mockUTRR) Eliminar(ctx context.Context, uid, tid, rid string) error { return nil }
func (m *mockUTRR) ListarRolesPorUsuarioEnTenant(ctx context.Context, uid, tid string) ([]*rbac.RolDB, error) {
	return m.roles, m.err
}
func (m *mockUTRR) TieneRolEnTenant(ctx context.Context, uid, tid, rol string) (bool, error) { return false, nil }

type mockUnitOfWork struct {
	sesionRepo    sesiones_domain.SesionRepositorio
	credRepo      seguridad_domain.CredencialesRepositorio
	usuarioRepo   usuario_domain.UsuarioRepositorio
	encriptacion  seguridad_domain.EncriptacionServicio
	tokenServicio sesiones_domain.TokenServicio
	generadorID   shared_domain.GeneradorID
}

func (m *mockUnitOfWork) Transaccional(ctx context.Context, fn func(tx sesiones_domain.UnitOfWork) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fn(m)
}
func (m *mockUnitOfWork) SesionRepositorio() sesiones_domain.SesionRepositorio { return m.sesionRepo }
func (m *mockUnitOfWork) CredencialesRepositorio() seguridad_domain.CredencialesRepositorio { return m.credRepo }
func (m *mockUnitOfWork) UsuarioRepositorio() usuario_domain.UsuarioRepositorio { return m.usuarioRepo }
func (m *mockUnitOfWork) EncriptacionServicio() seguridad_domain.EncriptacionServicio { return m.encriptacion }
func (m *mockUnitOfWork) TokenServicio() sesiones_domain.TokenServicio { return m.tokenServicio }
func (m *mockUnitOfWork) GeneradorID() shared_domain.GeneradorID { return m.generadorID }

type mockSesionRepo struct {
	sesion       *sesiones_domain.Sesion
	errObtener   error
	errActualizar error
}

func (m *mockSesionRepo) Crear(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) { return s, nil }
func (m *mockSesionRepo) Actualizar(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	if m.errActualizar != nil { return nil, m.errActualizar }
	return s, nil
}
func (m *mockSesionRepo) ObtenerPorID(ctx context.Context, id string) (*sesiones_domain.Sesion, error) {
	return m.sesion, m.errObtener
}
func (m *mockSesionRepo) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones_domain.Sesion, error) { return nil, nil }
func (m *mockSesionRepo) ListarActivasPorUsuarioID(ctx context.Context, uid string, ahora time.Time) ([]*sesiones_domain.Sesion, error) { return nil, nil }
func (m *mockSesionRepo) Listar(ctx context.Context, spec sesiones_domain.EspecificacionSesion, pag shared_domain.Paginacion) ([]*sesiones_domain.Sesion, error) { return nil, nil }
func (m *mockSesionRepo) InvalidarTodasPorUsuarioID(ctx context.Context, uid string) error { return nil }
func (m *mockSesionRepo) Eliminar(ctx context.Context, id string) error { return nil }

type mockTokenServicio struct{}

func (m *mockTokenServicio) GenerarAccessToken(uid, sid, tid, rol string) (string, time.Time, error) {
	return "new-access", time.Now().Add(15 * time.Minute), nil
}
func (m *mockTokenServicio) GenerarRefreshToken(uid, sid string) (string, time.Time, error) {
	return "new-refresh", time.Now().Add(24 * time.Hour), nil
}
func (m *mockTokenServicio) ValidarAccessToken(tok string) (*sesiones_domain.TokenClaims, error) { return nil, nil }
func (m *mockTokenServicio) ValidarRefreshToken(tok string) (*sesiones_domain.TokenClaims, error) { return nil, nil }
func (m *mockTokenServicio) HashearToken(tok string) string { return "hash:" + tok }

func sesionActiva() *sesiones_domain.Sesion {
	now := time.Now()
	s, _ := sesiones_domain.NuevaSesion("ses-1", "user-1", "hash:old-access", "hash:old-refresh",
		"127.0.0.1", now, now.Add(15*time.Minute), now.Add(24*time.Hour))
	return s
}

func TestCambiarTenantExitoso(t *testing.T) {
	membRepo := &mockMembresiaRepo{existeMiembro: true}
	utrRepo := &mockUTRR{roles: []*rbac.RolDB{{Nombre: "editor"}}}
	sesionRepo := &mockSesionRepo{sesion: sesionActiva()}
	uow := &mockUnitOfWork{
		sesionRepo:    sesionRepo,
		tokenServicio: &mockTokenServicio{},
	}
	uc := switchtenant.NewCambiarTenantCasoDeUso(membRepo, utrRepo, uow)
	resp, err := uc.Ejecutar(context.Background(), switchtenant.ComandoCambiarTenant{
		UsuarioID: "user-1", SesionID: "ses-1", TenantID: "tenant-2",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.TenantID != "tenant-2" {
		t.Errorf("TenantID incorrecto: %v", resp.TenantID)
	}
	if resp.Rol != "editor" {
		t.Errorf("Rol incorrecto: %v", resp.Rol)
	}
	if resp.UsuarioID != "user-1" {
		t.Errorf("UsuarioID incorrecto: %v", resp.UsuarioID)
	}
	if resp.SesionID != "ses-1" {
		t.Errorf("SesionID incorrecto: %v", resp.SesionID)
	}
	if resp.AccessToken == "" {
		t.Error("AccessToken vacío")
	}
	if resp.RefreshToken == "" {
		t.Error("RefreshToken vacío")
	}
}

func TestCambiarTenantNoEresMiembro(t *testing.T) {
	membRepo := &mockMembresiaRepo{existeMiembro: false}
	uc := switchtenant.NewCambiarTenantCasoDeUso(membRepo, &mockUTRR{}, &mockUnitOfWork{})
	_, err := uc.Ejecutar(context.Background(), switchtenant.ComandoCambiarTenant{
		UsuarioID: "user-1", SesionID: "ses-1", TenantID: "tenant-2",
	})
	if !errors.Is(err, switchtenant.ErrNoEresMiembro) {
		t.Errorf("esperaba ErrNoEresMiembro, got %v", err)
	}
}

func TestCambiarTenantExisteMiembroError(t *testing.T) {
	membRepo := &mockMembresiaRepo{errExiste: errors.New("fallo bd")}
	uc := switchtenant.NewCambiarTenantCasoDeUso(membRepo, &mockUTRR{}, &mockUnitOfWork{})
	_, err := uc.Ejecutar(context.Background(), switchtenant.ComandoCambiarTenant{
		UsuarioID: "user-1", SesionID: "ses-1", TenantID: "tenant-2",
	})
	if err == nil {
		t.Fatal("esperaba error de repositorio")
	}
}

func TestCambiarTenantSinRol(t *testing.T) {
	membRepo := &mockMembresiaRepo{existeMiembro: true}
	utrRepo := &mockUTRR{roles: []*rbac.RolDB{}}
	uc := switchtenant.NewCambiarTenantCasoDeUso(membRepo, utrRepo, &mockUnitOfWork{})
	_, err := uc.Ejecutar(context.Background(), switchtenant.ComandoCambiarTenant{
		UsuarioID: "user-1", SesionID: "ses-1", TenantID: "tenant-2",
	})
	if !errors.Is(err, switchtenant.ErrSinRolEnTenant) {
		t.Errorf("esperaba ErrSinRolEnTenant, got %v", err)
	}
}

func TestCambiarTenantListarRolesError(t *testing.T) {
	membRepo := &mockMembresiaRepo{existeMiembro: true}
	utrRepo := &mockUTRR{err: errors.New("fallo bd")}
	uc := switchtenant.NewCambiarTenantCasoDeUso(membRepo, utrRepo, &mockUnitOfWork{})
	_, err := uc.Ejecutar(context.Background(), switchtenant.ComandoCambiarTenant{
		UsuarioID: "user-1", SesionID: "ses-1", TenantID: "tenant-2",
	})
	if err == nil {
		t.Fatal("esperaba error al listar roles")
	}
}
