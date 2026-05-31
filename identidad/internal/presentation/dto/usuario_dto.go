package dto

type CrearUsuarioRequest struct {
	Correo   string `json:"correo"   doc:"Correo electrónico del usuario" example:"juan@correo.com"`
	Nombre   string `json:"nombre"   doc:"Nombre del usuario"             example:"Juan"`
	Apellido string `json:"apellido" doc:"Apellido del usuario"           example:"Pérez"`
	Password string `json:"password" doc:"Contraseña (mínimo 8 caracteres)" example:"secreto123"`
}

type CrearUsuarioResponse struct {
	ID       string `json:"id"        doc:"ID único del usuario" example:"01926b1e-..."`
	Correo   string `json:"correo"    doc:"Correo electrónico"   example:"juan@correo.com"`
	Nombre   string `json:"nombre"    doc:"Nombre del usuario"   example:"Juan"`
	Apellido string `json:"apellido"  doc:"Apellido del usuario" example:"Pérez"`
	Activo   bool   `json:"activo"    doc:"Usuario activo"       example:"true"`
	CreadoEn string `json:"creado_en" doc:"Fecha de creación"    example:"2026-05-23T12:00:00Z"`
}

type ListarUsuariosRequest struct {
	Pagina       int    `query:"pagina"       doc:"Número de página (1-based)" example:"1"`
	TamanoPagina int    `query:"tamano"       doc:"Elementos por página"       example:"20"`
	Correo       string `query:"correo"       doc:"Filtrar por correo"         example:""`
	Estado       string `query:"estado"       doc:"Filtrar por estado"         example:"ACTIVO"`
}

type UsuarioItem struct {
	ID       string `json:"id"        doc:"ID del usuario"`
	Correo   string `json:"correo"    doc:"Correo electrónico"`
	Nombre   string `json:"nombre"    doc:"Nombre del usuario"`
	Apellido string `json:"apellido"  doc:"Apellido del usuario"`
	Estado   string `json:"estado"    doc:"Estado del usuario"`
	CreadoEn string `json:"creado_en" doc:"Fecha de creación"`
}

type ListarUsuariosResponse struct {
	Usuarios []UsuarioItem `json:"usuarios" doc:"Lista de usuarios"`
	Total    int           `json:"total"    doc:"Total de resultados"`
	Pagina   int           `json:"pagina"   doc:"Página actual"`
	Tamano   int           `json:"tamano"   doc:"Elementos por página"`
}

type ModificarUsuarioRequest struct {
	Nombre   string `json:"nombre"   doc:"Nuevo nombre"   example:"Juan Actualizado"`
	Apellido string `json:"apellido" doc:"Nuevo apellido" example:"Pérez Actualizado"`
}

type ModificarUsuarioResponse struct {
	ID           string `json:"id"            doc:"ID del usuario"`
	Correo       string `json:"correo"        doc:"Correo electrónico"`
	Nombre       string `json:"nombre"        doc:"Nombre actualizado"`
	Apellido     string `json:"apellido"      doc:"Apellido actualizado"`
	ModificadoEn string `json:"modificado_en" doc:"Fecha de modificación"`
}

type DarDeBajaUsuarioRequest struct {
	Motivo string `json:"motivo,omitempty" doc:"Motivo de la baja" example:"Cierre de cuenta"`
}

type DarDeBajaUsuarioResponse struct {
	UsuarioID string `json:"usuario_id" doc:"ID del usuario dado de baja"`
	Estado    string `json:"estado"     doc:"Nuevo estado del usuario"`
	BajaEn    string `json:"baja_en"    doc:"Fecha de baja"`
}

type ExpulsarUsuarioResponse struct {
	UsuarioID         string `json:"usuario_id"          doc:"ID del usuario expulsado"`
	Estado            string `json:"estado"              doc:"Nuevo estado del usuario"`
	SesionesRevocadas int    `json:"sesiones_revocadas"  doc:"Cantidad de sesiones cerradas"`
	ExpulsadoEn       string `json:"expulsado_en"        doc:"Fecha de expulsión"`
}

type VerMiPerfilResponse struct {
	ID       string `json:"id"        doc:"ID del usuario"`
	Correo   string `json:"correo"    doc:"Correo electrónico"`
	Nombre   string `json:"nombre"    doc:"Nombre del usuario"`
	Apellido string `json:"apellido"  doc:"Apellido del usuario"`
	Telefono string `json:"telefono"  doc:"Teléfono del usuario"`
	Estado   string `json:"estado"    doc:"Estado del usuario"`
	CreadoEn string `json:"creado_en" doc:"Fecha de creación"`
}

type ModificarMiPerfilRequest struct {
	Nombre   string `json:"nombre"   doc:"Nuevo nombre"   example:"Juan"`
	Apellido string `json:"apellido" doc:"Nuevo apellido" example:"Pérez"`
}

type ModificarMiPerfilResponse struct {
	ID           string `json:"id"            doc:"ID del usuario"`
	Correo       string `json:"correo"        doc:"Correo electrónico"`
	Nombre       string `json:"nombre"        doc:"Nombre actualizado"`
	Apellido     string `json:"apellido"      doc:"Apellido actualizado"`
	ModificadoEn string `json:"modificado_en" doc:"Fecha de modificación"`
}
