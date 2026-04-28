package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"gorm.io/gorm"
)

type credencialesRepositorio struct {
	db *gorm.DB
}

func NewCredencialesRepositorio(db *gorm.DB) domain.CredencialesRepositorio {
	return &credencialesRepositorio{db: db}
}

func (r *credencialesRepositorio) Crear(ctx context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
	model, err := CredencialesFromDomain(c)
	if err != nil {
		return nil, fmt.Errorf("error al convertir credenciales a modelo: %w", err)
	}
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}
	// Retorna credenciales con datos desde BD
	return model.ToDomain(), nil
}

func (r *credencialesRepositorio) Actualizar(ctx context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
	model, err := CredencialesFromDomain(c)
	if err != nil {
		return nil, fmt.Errorf("error al convertir credenciales a modelo: %w", err)
	}

	// Actualizar explícitamente los campos relevantes
	result := r.db.WithContext(ctx).Model(&CredencialesModel{}).
		Where("usuario_id = ?", c.UsuarioID()).
		Updates(map[string]interface{}{
			"password_hash":     model.PasswordHash,
			"activo":            model.Activo,
			"correo_verificado": model.CorreoVerificado,
			"intentos_fallidos": model.IntentosFallidos,
			"bloqueado_hasta":   model.BloqueadoHasta,
		})

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("credenciales con usuarioID %s no encontradas", c.UsuarioID())
	}
	// Retornar las credenciales actualizadas desde BD
	return r.ObtenerPorUsuarioID(ctx, c.UsuarioID())
}

func (r *credencialesRepositorio) ObtenerPorUsuarioID(ctx context.Context, usuarioID string) (*domain.CredencialesUsuario, error) {
	var model CredencialesModel
	result := r.db.WithContext(ctx).First(&model, "usuario_id = ?", usuarioID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("credenciales no encontradas")
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

func (r *credencialesRepositorio) Eliminar(ctx context.Context, usuarioID string) error {
	result := r.db.WithContext(ctx).Delete(&CredencialesModel{}, "usuario_id = ?", usuarioID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("credenciales no encontradas")
	}
	return nil
}

// Método Find implementa búsqueda avanzada con filtros y paginación
func (r *credencialesRepositorio) Find(ctx context.Context, spec domain.EspecificacionCredenciales, pag domain.Paginacion) ([]*domain.CredencialesUsuario, error) {
	query := r.db.WithContext(ctx).Model(&CredencialesModel{})

	mapeoColumnas := map[string]string{
		"usuarioID":        "usuario_id",
		"activo":           "activo",
		"intentosFallidos": "intentos_fallidos",
		"bloqueadoHasta":   "bloqueado_hasta",
		"correoVerificado": "correo_verificado",
	}

	// Aplicar filtros
	for _, filtro := range spec.ListaFiltros {
		if !domain.ColumnasPermitidas[filtro.Campo] {
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
			// BETWEEN será implementado próximamente
			// case "BETWEEN":
			//     valores := filtro.Valor.([]interface{})
			//     if len(valores) == 2 {
			//         query = query.Where(columnaDB+" BETWEEN ? AND ?", valores[0], valores[1])
			//     }
		}
	}

	// Aplicar ordenaciones
	for _, ord := range pag.Ordenaciones {
		if !domain.ColumnasPermitidas[ord.Campo] {
			continue
		}

		columnaDB, ok := mapeoColumnas[ord.Campo]
		if !ok {
			continue
		}

		orden := "ASC"
		if ord.Tipo == domain.DESC {
			orden = "DESC"
		}
		query = query.Order(columnaDB + " " + orden)
	}

	// Aplicar paginación
	offset := (pag.Pagina - 1) * pag.TamanoPagina
	if pag.Pagina < 1 {
		offset = 0
	}
	if pag.TamanoPagina < 1 {
		pag.TamanoPagina = 10
	}

	var models []CredencialesModel
	result := query.Offset(offset).Limit(pag.TamanoPagina).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	credenciales := make([]*domain.CredencialesUsuario, len(models))
	for i, m := range models {
		credenciales[i] = m.ToDomain()
	}

	return credenciales, nil
}

// Método auxiliar para listar con filtros (puede ser extendido en el futuro)
// Nota: Este método no está incluido en la interfaz CredencialesRepositorio,
// pero puede agregarse si se requiere búsqueda avanzada con filtros
func (r *credencialesRepositorio) Listar(ctx context.Context, spec domain.EspecificacionCredenciales) ([]*domain.CredencialesUsuario, error) {
	query := r.db.WithContext(ctx).Model(&CredencialesModel{})

	mapeoColumnas := map[string]string{
		"usuarioID":        "usuario_id",
		"activo":           "activo",
		"intentosFallidos": "intentos_fallidos",
		"bloqueadoHasta":   "bloqueado_hasta",
		"correoVerificado": "correo_verificado",
	}

	for _, filtro := range spec.ListaFiltros {
		if !domain.ColumnasPermitidas[filtro.Campo] {
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
			// BETWEEN será implementado próximamente
			// case "BETWEEN":
			//     valores := filtro.Valor.([]interface{})
			//     if len(valores) == 2 {
			//         query = query.Where(columnaDB+" BETWEEN ? AND ?", valores[0], valores[1])
			//     }
		}
	}

	var models []CredencialesModel
	result := query.Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	credenciales := make([]*domain.CredencialesUsuario, len(models))
	for i, m := range models {
		credenciales[i] = m.ToDomain()
	}

	return credenciales, nil
}
