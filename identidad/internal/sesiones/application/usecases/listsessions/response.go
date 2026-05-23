package listsessions

import "time"

type SesionDTO struct {
	ID              string
	UsuarioID       string
	IPOrigen        string
	Estado          string
	UltimaActividad time.Time
}

type RespuestaListarSesiones struct {
	Sesiones []SesionDTO
	Total    int
	Pagina   int
}
