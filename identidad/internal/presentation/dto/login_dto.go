package dto

// LoginRequest es el body del request POST /api/v1/auth/login.
type LoginRequest struct {
	Correo   string `json:"correo"   doc:"Correo electrónico del usuario" example:"juan@correo.com"`
	Password string `json:"password" doc:"Contraseña del usuario"         example:"secreto123"`
}

// LoginResponse es el payload dentro de data en la respuesta 200.
type LoginResponse struct {
	AccessToken  string `json:"access_token"  doc:"JWT access token"                    example:"eyJhbGci..."`
	RefreshToken string `json:"refresh_token" doc:"JWT refresh token"                   example:"eyJhbGci..."`
	ExpiresIn    int64  `json:"expires_in"    doc:"Segundos hasta expiración del access" example:"900"`
	TokenType    string `json:"token_type"    doc:"Tipo de token, siempre Bearer"        example:"Bearer"`
	UsuarioID    string `json:"usuario_id"    doc:"ID del usuario autenticado"           example:"01926b1e-dead-beef-cafe-000000000001"`
	TenantID     string `json:"tenant_id"     doc:"ID del tenant del usuario"            example:"01926b1e-dead-beef-cafe-000000000002"`
	Rol          string `json:"rol"           doc:"Rol del usuario en el tenant"         example:"administrador"`
}
