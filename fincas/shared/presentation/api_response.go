package presentation

// Link representa un link HATEOAS
type Link struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}

// ApiResponse es la respuesta genérica para todos los endpoints.
type ApiResponse[T any] struct {
	Data  T                `json:"data"`
	Links map[string]Link  `json:"_links,omitempty"`
}

// NewResponse crea una ApiResponse con data y links HATEOAS.
func NewResponse[T any](data T, links map[string]Link) ApiResponse[T] {
	return ApiResponse[T]{Data: data, Links: links}
}
