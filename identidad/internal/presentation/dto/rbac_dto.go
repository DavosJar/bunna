package dto

type RolItem struct {
	ID          string   `json:"id"           doc:"ID del rol"`
	Nombre      string   `json:"nombre"       doc:"Nombre del rol"`
	Descripcion string   `json:"descripcion"  doc:"Descripción del rol"`
	EsSistema   bool     `json:"es_sistema"   doc:"Rol de sistema (no modificable)"`
	Permisos    []string `json:"permisos"     doc:"Códigos de permisos asociados"`
}

type ListarRolesResponse struct {
	Roles  []RolItem `json:"roles"  doc:"Lista de roles"`
	Total  int       `json:"total"  doc:"Total de resultados"`
	Pagina int       `json:"pagina" doc:"Página actual"`
}

type CrearRolRequest struct {
	Nombre      string   `json:"nombre"      doc:"Nombre del rol"          example:"Editor"`
	Descripcion string   `json:"descripcion" doc:"Descripción del rol"     example:"Puede editar contenidos"`
	Permisos    []string `json:"permisos"    doc:"Códigos de permisos iniciales" example:"[\"usuario:leer\"]"`
}

type CrearRolResponse struct {
	ID          string `json:"id"           doc:"ID del rol creado"`
	Nombre      string `json:"nombre"       doc:"Nombre del rol"`
	Descripcion string `json:"descripcion"  doc:"Descripción del rol"`
	EsSistema   bool   `json:"es_sistema"   doc:"Rol de sistema"`
	CreadoEn    string `json:"creado_en"    doc:"Fecha de creación"`
}

type ModificarRolRequest struct {
	Nombre      string `json:"nombre"      doc:"Nuevo nombre"       example:"Editor Senior"`
	Descripcion string `json:"descripcion" doc:"Nueva descripción"  example:"Puede editar y publicar contenidos"`
}

type ModificarRolResponse struct {
	ID           string `json:"id"            doc:"ID del rol"`
	Nombre       string `json:"nombre"        doc:"Nombre actualizado"`
	Descripcion  string `json:"descripcion"   doc:"Descripción actualizada"`
	ModificadoEn string `json:"modificado_en" doc:"Fecha de modificación"`
}

type EliminarRolResponse struct {
	RolID       string `json:"rol_id"       doc:"ID del rol eliminado"`
	EliminadoEn string `json:"eliminado_en" doc:"Fecha de eliminación"`
}

type AsignarRolRequest struct {
	RolID    string `json:"rol_id"     doc:"ID del rol a asignar"  example:"rol-123"`
	TenantID string `json:"tenant_id"  doc:"ID del tenant (opcional, vacío = global)" example:""`
}

type AsignarRolResponse struct {
	UsuarioID  string `json:"usuario_id"  doc:"ID del usuario"`
	RolID      string `json:"rol_id"      doc:"ID del rol asignado"`
	TenantID   string `json:"tenant_id"   doc:"ID del tenant"`
	AsignadoEn string `json:"asignado_en" doc:"Fecha de asignación"`
}

type RevocarRolResponse struct {
	UsuarioID  string `json:"usuario_id"  doc:"ID del usuario"`
	RolID      string `json:"rol_id"      doc:"ID del rol revocado"`
	TenantID   string `json:"tenant_id"   doc:"ID del tenant"`
	RevocadoEn string `json:"revocado_en" doc:"Fecha de revocación"`
}

type AsignarPermisoRequest struct {
	PermisoCodigo string `json:"permiso_codigo" doc:"Código del permiso" example:"usuario:crear"`
}

type AsignarPermisoResponse struct {
	RolID         string `json:"rol_id"          doc:"ID del rol"`
	PermisoCodigo string `json:"permiso_codigo"  doc:"Código del permiso asignado"`
	AsignadoEn    string `json:"asignado_en"     doc:"Fecha de asignación"`
}

type RevocarPermisoResponse struct {
	RolID         string `json:"rol_id"          doc:"ID del rol"`
	PermisoCodigo string `json:"permiso_codigo"  doc:"Código del permiso revocado"`
	RevocadoEn    string `json:"revocado_en"     doc:"Fecha de revocación"`
}

type PermisoItem struct {
	ID          string `json:"id"          doc:"ID del permiso"`
	Codigo      string `json:"codigo"      doc:"Código del permiso"`
	Nombre      string `json:"nombre"      doc:"Nombre del permiso"`
	Descripcion string `json:"descripcion" doc:"Descripción del permiso"`
	Modulo      string `json:"modulo"      doc:"Módulo al que pertenece"`
}

type ListarPermisosResponse struct {
	Permisos []PermisoItem `json:"permisos" doc:"Lista de permisos"`
	Total    int           `json:"total"    doc:"Total de resultados"`
}

type ListarRolesDeUsuarioResponse struct {
	Roles []RolDeUsuarioItem `json:"roles"`
}

type RolDeUsuarioItem struct {
	RolID  string `json:"rol_id"`
	Nombre string `json:"nombre"`
}
