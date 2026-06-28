package expeluser_test

import (
	"context"
	"time"

	sesiones "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuariodomain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type mockUsuarioRepoExpel struct {
	obtenerPorID func(ctx context.Context, id string) (*usuariodomain.Usuario, error)
	actualizar   func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error)
}

func (m *mockUsuarioRepoExpel) Crear(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return u, nil
}
func (m *mockUsuarioRepoExpel) Actualizar(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	if m.actualizar != nil {
		return m.actualizar(ctx, u)
	}
	return u, nil
}
func (m *mockUsuarioRepoExpel) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockUsuarioRepoExpel) ObtenerPorID(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
	if m.obtenerPorID != nil {
		return m.obtenerPorID(ctx, id)
	}
	return nil, nil
}
func (m *mockUsuarioRepoExpel) ObtenerPorCorreo(ctx context.Context, correo string) (*usuariodomain.Usuario, error) {
	return nil, nil
}
func (m *mockUsuarioRepoExpel) Listar(ctx context.Context, _ usuariodomain.EspecificacionUsuario, _ shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
	return nil, nil
}

type mockSesionRepoExpel struct {
	errInvalidar error
}

func (m *mockSesionRepoExpel) Crear(ctx context.Context, s *sesiones.Sesion) (*sesiones.Sesion, error) {
	return s, nil
}
func (m *mockSesionRepoExpel) Actualizar(ctx context.Context, s *sesiones.Sesion) (*sesiones.Sesion, error) {
	return s, nil
}
func (m *mockSesionRepoExpel) ObtenerPorID(ctx context.Context, id string) (*sesiones.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepoExpel) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepoExpel) ListarActivasPorUsuarioID(ctx context.Context, uid string, ahora time.Time) ([]*sesiones.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepoExpel) Listar(ctx context.Context, spec sesiones.EspecificacionSesion, pag shareddomain.Paginacion) ([]*sesiones.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepoExpel) InvalidarTodasPorUsuarioID(ctx context.Context, uid string) error {
	return m.errInvalidar
}
func (m *mockSesionRepoExpel) Eliminar(ctx context.Context, id string) error { return nil }

type mockAuthSvcExpel struct {
	ok  bool
	err error
}

func (m *mockAuthSvcExpel) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}
