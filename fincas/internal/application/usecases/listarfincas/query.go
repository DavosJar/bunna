package listarfincas

// Query contiene los parámetros de búsqueda. Por simplicidad no paginamos aquí.
type Query struct {
}

func (q *Query) Validar() error {
	return nil
}
