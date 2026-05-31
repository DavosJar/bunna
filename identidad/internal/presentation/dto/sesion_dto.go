package dto

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" doc:"Token de refresco" example:"eyJhbGci..."`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"  doc:"Nuevo JWT access token"             example:"eyJhbGci..."`
	RefreshToken string `json:"refresh_token" doc:"Nuevo JWT refresh token (rotado)"  example:"eyJhbGci..."`
	ExpiresIn    int64  `json:"expires_in"    doc:"Segundos hasta expiración del access" example:"900"`
	TokenType    string `json:"token_type"    doc:"Tipo de token"                      example:"Bearer"`
	UsuarioID    string `json:"usuario_id"    doc:"ID del usuario"                     example:"01926b1e-..."`
}

type LogoutResponse struct {
	SesionesRevocadas int `json:"sesiones_revocadas" doc:"Cantidad de sesiones cerradas" example:"1"`
}

type SesionItem struct {
	ID              string `json:"id"               doc:"ID de la sesión"`
	UsuarioID       string `json:"usuario_id"       doc:"ID del usuario"`
	IPOrigen        string `json:"ip_origen"        doc:"IP de origen"`
	Estado          string `json:"estado"            doc:"Estado de la sesión"`
	UltimaActividad string `json:"ultima_actividad"  doc:"Última actividad"`
}

type ListarSesionesResponse struct {
	Sesiones []SesionItem `json:"sesiones" doc:"Lista de sesiones"`
	Total    int          `json:"total"    doc:"Total de resultados"`
	Pagina   int          `json:"pagina"   doc:"Página actual"`
}

type ForzarCierreSesionResponse struct {
	SesionID   string `json:"sesion_id"   doc:"ID de la sesión cerrada"`
	Estado     string `json:"estado"      doc:"Nuevo estado"`
	RevocadoEn string `json:"revocado_en" doc:"Fecha de revocación"`
}
