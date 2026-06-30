package dto

// RegisterRequest es el body del request POST /api/v1/identidad/auth/register.
type RegisterRequest struct {
	Nombre   string  `json:"nombre"             doc:"Nombre del usuario"                example:"Juan"`
	Apellido string  `json:"apellido"           doc:"Apellido del usuario"              example:"Pérez"`
	Correo   string  `json:"correo"             doc:"Correo electrónico válido"         example:"juan@correo.com"`
	Password string  `json:"password" minLength:"8" doc:"Contraseña (mínimo 8 caracteres, 1 mayúscula, 1 minúscula, 1 número, 1 no alfanumérico)" example:"Secreto1!"`
	Telefono *string `json:"telefono,omitempty" doc:"Teléfono de contacto (opcional)"  example:"0999999999"`
}

// RegisterResponse es el payload dentro de data en la respuesta 201.
type RegisterResponse struct {
	UsuarioID string `json:"usuario_id" doc:"ID único del usuario creado"   example:"01926b1e-dead-beef-cafe-000000000001"`
	Correo    string `json:"correo"     doc:"Correo electrónico registrado"  example:"juan@correo.com"`
	Estado    string `json:"estado"     doc:"Estado inicial del usuario"     example:"NO_VERIFICADO"`
}
