package dto

type CambiarMiPasswordRequest struct {
	PasswordActual string `json:"password_actual" doc:"Contraseña actual"       example:"vieja123"`
	NuevaPassword  string `json:"nueva_password" minLength:"8" doc:"Nueva contraseña (mínimo 8 caracteres, 1 mayúscula, 1 minúscula, 1 número, 1 no alfanumérico)" example:"Secreto1!"`
}

type CambiarMiPasswordResponse struct {
	ModificadoEn string `json:"modificado_en" doc:"Fecha de modificación"`
}

type ResetearPasswordRequest struct {
	NuevaPassword string `json:"nueva_password" minLength:"8" doc:"Nueva contraseña (mínimo 8 caracteres, 1 mayúscula, 1 minúscula, 1 número, 1 no alfanumérico)" example:"Secreto1!"`
}

type ResetearPasswordResponse struct {
	UsuarioID    string `json:"usuario_id"   doc:"ID del usuario"`
	ModificadoEn string `json:"modificado_en" doc:"Fecha de modificación"`
}

type DesbloquearCuentaResponse struct {
	UsuarioID      string `json:"usuario_id"      doc:"ID del usuario desbloqueado"`
	DesbloqueadoEn string `json:"desbloqueado_en" doc:"Fecha de desbloqueo"`
}

type IPBloqueadaItem struct {
	IP             string `json:"ip"               doc:"Dirección IP bloqueada"`
	Intentos       int    `json:"intentos"          doc:"Cantidad de intentos fallidos"`
	BloqueadoHasta string `json:"bloqueado_hasta"  doc:"Fecha hasta la que está bloqueada"`
}

type ListarIPsBloqueadasResponse struct {
	IPs    []IPBloqueadaItem `json:"ips"    doc:"Lista de IPs bloqueadas"`
	Total  int               `json:"total"  doc:"Total de resultados"`
	Pagina int               `json:"pagina" doc:"Página actual"`
}

type DesbloquearIPResponse struct {
	IP             string `json:"ip"               doc:"IP desbloqueada"`
	DesbloqueadoEn string `json:"desbloqueado_en"  doc:"Fecha de desbloqueo"`
}

type ConsultarCredencialesResponse struct {
	UsuarioID        string `json:"usuario_id"         doc:"ID del usuario"`
	Activo           bool   `json:"activo"             doc:"Cuenta activa"`
	CorreoVerificado bool   `json:"correo_verificado"  doc:"Correo verificado"`
	IntentosFallidos int    `json:"intentos_fallidos"  doc:"Intentos fallidos consecutivos"`
	BloqueadoHasta   string `json:"bloqueado_hasta"    doc:"Fecha de desbloqueo (vacío si no bloqueada)"`
}
