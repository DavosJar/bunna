package domain

import shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type EspecificacionCredenciales struct {
	ListaFiltros []shareddomain.CriterioFiltro
}

var ColumnasPermitidas = map[string]bool{
	"usuarioID":        true,
	"activo":           true,
	"intentosFallidos": true,
	"bloqueadoHasta":   true,
	"correoVerificado": true,
}

// Nota: BETWEEN será implementado próximamente para permitir rangos en bloqueadoHasta
// Ejemplo de uso futuro:
// case "BETWEEN":
//     // Validar que filtro.Valor es un slice con dos elementos [min, max]
//     query = query.Where(columnaDB+" BETWEEN ? AND ?", filtro.Valor.([]interface{})[0], filtro.Valor.([]interface{})[1])
