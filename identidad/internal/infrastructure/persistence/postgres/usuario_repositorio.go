package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/davosjar/bunna/services/identidad/internal/domain/usuario"
	"gorm.io/gorm"
)

type usuarioRepositorio struct {
	db *gorm.DB
}

func NewUsuarioRepositorio(db *gorm.DB) usuario.UsuarioRepositorio {
	return &usuarioRepositorio{db: db}
}

func (r *usuarioRepositorio) Crear(ctx context.Context, u *usuario.Usuario) (*usuario.Usuario, error) {
	model, err := FromDomain(u)
	if err != nil {
		return nil, fmt.Errorf("error al convertir usuario a modelo: %w", err)
	}
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}
	// Retorna usuario con ID asignado por BD
	return model.ToDomain(), nil
}

func (r *usuarioRepositorio) Actualizar(ctx context.Context, u *usuario.Usuario) (*usuario.Usuario, error) {
	model, err := FromDomain(u)
	if err != nil {
		return nil, fmt.Errorf("error al convertir usuario a modelo: %w", err)
	}

	// Actualizar explícitamente solo los campos relevantes
	// GORM automáticamente actualiza fecha_actualizacion gracias al tag autoUpdateTime
	result := r.db.WithContext(ctx).Model(&UsuarioModel{}).
		Where("id = ?", u.ID()).
		Updates(map[string]interface{}{
			"nombre":                     model.Nombre,
			"apellido":                   model.Apellido,
			"correo":                     model.Correo,
			"telefono":                   model.Telefono,
			"estado":                     model.Estado,
			"estado_verificacion_correo": model.EstadoVerificacionCorreo,
		})

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("usuario con ID %s no encontrado", u.ID())
	}
	// Retornar el usuario actualizado desde BD
	return r.ObtenerPorID(ctx, u.ID())
}

func (r *usuarioRepositorio) Eliminar(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&UsuarioModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("usuario no encontrado")
	}
	return nil
}

func (r *usuarioRepositorio) ObtenerPorID(ctx context.Context, id string) (*usuario.Usuario, error) {
	var model UsuarioModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

func (r *usuarioRepositorio) Listar(ctx context.Context, spec usuario.EspecificacionUsuario, pag usuario.Paginacion) ([]*usuario.Usuario, error) {
	query := r.db.WithContext(ctx).Model(&UsuarioModel{})

	mapeoColumnas := map[string]string{
		"nombre":                   "nombre",
		"apellido":                 "apellido",
		"correo":                   "correo",
		"fechaCreacion":            "fecha_creacion",
		"fechaActualizacion":       "fecha_actualizacion",
		"estado":                   "estado",
		"telefono":                 "telefono",
		"estadoVerificacionCorreo": "estado_verificacion_correo",
	}

	for _, filtro := range spec.ListaLiltros {
		if !usuario.ColumnasPermitidas[filtro.Campo] {
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
		}
	}

	for _, ord := range pag.Ordenaciones {
		if !usuario.ColumnasPermitidas[ord.Campo] {
			continue
		}

		columnaDB, ok := mapeoColumnas[ord.Campo]
		if !ok {
			continue
		}

		orden := "ASC"
		if ord.Tipo == usuario.DESC {
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

	var models []UsuarioModel
	result := query.Offset(offset).Limit(pag.TamanoPagina).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	usuarios := make([]*usuario.Usuario, len(models))
	for i, m := range models {
		usuarios[i] = m.ToDomain()
	}

	return usuarios, nil
}
