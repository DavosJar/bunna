package dto

// ErrorResponse es la respuesta estándar para errores.
type ErrorResponse struct {
	Error   string `json:"error"`
	Detalle string `json:"detalle,omitempty"`
}

// ListaResponse es un wrapper genérico para respuestas paginadas.
type ListaResponse[T any] struct {
	Data         []T `json:"data"`
	Total        int `json:"total,omitempty"`
	Pagina       int `json:"pagina,omitempty"`
	TotalPaginas int `json:"totalPaginas,omitempty"`
}

// PaginacionRequest son los params de paginación en query string.
type PaginacionRequest struct {
	Pagina       int    `form:"pagina"`
	TamanoPagina int    `form:"tamanoPagina"`
	OrdenarPor   string `form:"ordenarPor"`
	Orden        string `form:"orden"`
}
