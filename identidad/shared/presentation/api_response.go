// Package presentation contiene estructuras compartidas de la capa de presentación.
package presentation

// Link representa un hipervínculo HATEOAS en una respuesta de la API.
type Link struct {
	Href   string `json:"href"             doc:"URL del recurso enlazado"`
	Method string `json:"method"           doc:"Método HTTP del enlace"`
}

// ApiResponse es la estructura genérica para todas las respuestas exitosas de la API.
// T representa el tipo del payload de datos de la respuesta.
type ApiResponse[T any] struct {
	Data  T               `json:"data"             doc:"Payload de la respuesta"`
	Links map[string]Link `json:"_links,omitempty" doc:"Links HATEOAS opcionales"`
}

// NewApiResponse construye un ApiResponse sin links.
func NewApiResponse[T any](data T) ApiResponse[T] {
	return ApiResponse[T]{Data: data}
}

// NewApiResponseWithLinks construye un ApiResponse con links HATEOAS.
func NewApiResponseWithLinks[T any](data T, links map[string]Link) ApiResponse[T] {
	return ApiResponse[T]{Data: data, Links: links}
}
