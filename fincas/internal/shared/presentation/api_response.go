package presentation

type Link struct {
	Href   string `json:"href" doc:"URL del recurso"`
	Method string `json:"method" doc:"Método HTTP"`
}

type ApiResponse[T any] struct {
	Data  T               `json:"data" doc:"Payload de la respuesta"`
	Links map[string]Link `json:"links" doc:"Enlaces de navegación"`
}

func NewApiResponse[T any](data T, links map[string]Link) ApiResponse[T] {
	return ApiResponse[T]{
		Data: data,
	}
}

func NewApiResponseWithLinks[T any](data T, links map[string]Link) ApiResponse[T] {
	return ApiResponse[T]{
		Data:  data,
		Links: links,
	}
}
