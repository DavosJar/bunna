package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"gorm.io/gorm"
)

type sesionRepositorio struct {
	db *gorm.DB
}

// NewSesionRepositorio crea una nueva instancia del repositorio de sesiones.
func NewSesionRepositorio(db *gorm.DB) sesiones_domain.SesionRepositorio {
	return &sesionRepositorio{db: db}
}

// Crear persiste una nueva sesión en la base de datos.
func (r *sesionRepositorio) Crear(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	model := SesionFromDomain(s)
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

// Actualizar actualiza una sesión existente usando SQL crudo para evitar
// el problema de zero values de GORM con campos booleanos y strings vacíos.
func (r *sesionRepositorio) Actualizar(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	result := r.db.WithContext(ctx).Exec(
		`UPDATE sesiones SET
			access_token_hash = ?,
			refresh_token_hash = ?,
			estado = ?,
			fecha_actualizacion = ?,
			fecha_expiracion_access = ?,
			fecha_expiracion_refresh = ?,
			ultima_actividad = ?,
			contador_refrescos = ?
		WHERE id = ?`,
		s.AccessTokenHash(),
		s.RefreshTokenHash(),
		string(s.Estado()),
		s.FechaActualizacion(),
		s.FechaExpiracionAccess(),
		s.FechaExpiracionRefresh(),
		s.UltimaActividad(),
		s.ContadorRefrescos(),
		s.ID(),
	)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("sesión con ID %s no encontrada", s.ID())
	}
	return r.ObtenerPorID(ctx, s.ID())
}

// ObtenerPorID retorna una sesión por su ID.
func (r *sesionRepositorio) ObtenerPorID(ctx context.Context, sesionID string) (*sesiones_domain.Sesion, error) {
	var model SesionModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", sesionID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("sesión no encontrada")
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

// ObtenerPorRefreshTokenHash retorna una sesión por el hash de su refresh token.
func (r *sesionRepositorio) ObtenerPorRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*sesiones_domain.Sesion, error) {
	var model SesionModel
	result := r.db.WithContext(ctx).First(&model, "refresh_token_hash = ?", refreshTokenHash)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("sesión no encontrada")
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

// ListarActivasPorUsuarioID retorna todas las sesiones activas de un usuario
// cuyo access token no ha expirado en el momento dado.
func (r *sesionRepositorio) ListarActivasPorUsuarioID(ctx context.Context, usuarioID string, ahora time.Time) ([]*sesiones_domain.Sesion, error) {
	var models []SesionModel
	result := r.db.WithContext(ctx).
		Where("usuario_id = ? AND estado = ? AND fecha_expiracion_access > ?",
			usuarioID, string(sesiones_domain.EstadoActiva), ahora).
		Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}
	sesiones := make([]*sesiones_domain.Sesion, len(models))
	for i, m := range models {
		sesiones[i] = m.ToDomain()
	}
	return sesiones, nil
}

// Listar retorna sesiones filtradas y paginadas según la especificación.
func (r *sesionRepositorio) Listar(ctx context.Context, spec sesiones_domain.EspecificacionSesion, pag shareddomain.Paginacion) ([]*sesiones_domain.Sesion, error) {
	query := r.db.WithContext(ctx).Model(&SesionModel{})

	mapeoColumnas := map[string]string{
		"usuarioID":              "usuario_id",
		"estado":                 "estado",
		"ipOrigen":               "ip_origen",
		"fechaCreacion":          "fecha_creacion",
		"fechaActualizacion":     "fecha_actualizacion",
		"fechaExpiracionAccess":  "fecha_expiracion_access",
		"fechaExpiracionRefresh": "fecha_expiracion_refresh",
		"ultimaActividad":        "ultima_actividad",
		"contadorRefrescos":      "contador_refrescos",
	}

	for _, filtro := range spec.ListaFiltros {
		if !sesiones_domain.ColumnasPermitidas[filtro.Campo] {
			continue
		}

		columnaDB, ok := mapeoColumnas[filtro.Campo]
		if !ok {
			continue
		}

		switch filtro.Operador {
		case "=":
			query = query.Where(columnaDB+" = ?", filtro.Valor)
		case "!=":
			query = query.Where(columnaDB+" != ?", filtro.Valor)
		case "LIKE":
			query = query.Where(columnaDB+" LIKE ?", filtro.Valor)
		case ">":
			query = query.Where(columnaDB+" > ?", filtro.Valor)
		case "<":
			query = query.Where(columnaDB+" < ?", filtro.Valor)
		case ">=":
			query = query.Where(columnaDB+" >= ?", filtro.Valor)
		case "<=":
			query = query.Where(columnaDB+" <= ?", filtro.Valor)
		}
	}

	for _, ord := range pag.Ordenaciones {
		if !sesiones_domain.ColumnasPermitidas[ord.Campo] {
			continue
		}

		columnaDB, ok := mapeoColumnas[ord.Campo]
		if !ok {
			continue
		}

		orden := "ASC"
		if ord.Tipo == shareddomain.DESC {
			orden = "DESC"
		}
		query = query.Order(columnaDB + " " + orden)
	}

	offset := (pag.Pagina - 1) * pag.TamanoPagina
	if pag.Pagina < 1 {
		offset = 0
	}
	if pag.TamanoPagina < 1 {
		pag.TamanoPagina = 10
	}

	var models []SesionModel
	result := query.Offset(offset).Limit(pag.TamanoPagina).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	sesiones := make([]*sesiones_domain.Sesion, len(models))
	for i, m := range models {
		sesiones[i] = m.ToDomain()
	}
	return sesiones, nil
}

// InvalidarTodasPorUsuarioID revoca todas las sesiones activas de un usuario.
// Se usa en detección de robo de refresh token.
func (r *sesionRepositorio) InvalidarTodasPorUsuarioID(ctx context.Context, usuarioID string) error {
	result := r.db.WithContext(ctx).Exec(
		`UPDATE sesiones SET estado = ? WHERE usuario_id = ? AND estado = ?`,
		string(sesiones_domain.EstadoRevocada),
		usuarioID,
		string(sesiones_domain.EstadoActiva),
	)
	return result.Error
}

// Eliminar elimina una sesión por su ID.
func (r *sesionRepositorio) Eliminar(ctx context.Context, sesionID string) error {
	result := r.db.WithContext(ctx).Delete(&SesionModel{}, "id = ?", sesionID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("sesión no encontrada")
	}
	return nil
}
