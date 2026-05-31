package listblockedips

type IPBloqueadaDTO struct {
	IP             string
	Intentos       int
	BloqueadoHasta string
}

type RespuestaListarIPsBloqueadas struct {
	IPs    []IPBloqueadaDTO
	Total  int
	Pagina int
}
