package application

// ModuloFincas identifica el módulo para eventos de permisos
const ModuloFincas = "fincas"

// Permisos de fincas — códigos cortos (se mantienen igual que hoy para no romper contracts)
const (
	PermisoCrearFinca           = "fincas:finca:crear"
	PermisoDesactivarFinca      = "fincas:finca:desactivar"
	PermisoCrearLote            = "fincas:lote:crear"
	PermisoEliminarLote         = "fincas:lote:eliminar"
	PermisoCrearMuestra         = "fincas:muestra:crear"
	PermisoVerMuestras          = "fincas:muestra:consultar"
	PermisoSolicitarDiagnostico = "fincas:diagnostico:solicitar"
	PermisoAceptarDiagnostico   = "fincas:diagnostico:aceptar"
	PermisoRechazarDiagnostico  = "fincas:diagnostico:rechazar"
	PermisoGenerarReporte       = "fincas:reporte:generar"
)

// PermisoInfo contiene el metadata completo de un permiso
type PermisoInfo struct {
	Codigo      string `json:"codigo"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Modulo      string `json:"modulo"`
}

// CatalogoFincas lista todos los permisos del módulo fincas con metadata
var CatalogoFincas = []PermisoInfo{
	{Codigo: PermisoCrearFinca, Nombre: "Crear Finca", Descripcion: "Registrar una nueva finca en el sistema", Modulo: ModuloFincas},
	{Codigo: PermisoDesactivarFinca, Nombre: "Desactivar Finca", Descripcion: "Marcar una finca como pendiente de eliminación", Modulo: ModuloFincas},
	{Codigo: PermisoCrearLote, Nombre: "Crear Lote", Descripcion: "Agregar un lote a una finca existente", Modulo: ModuloFincas},
	{Codigo: PermisoEliminarLote, Nombre: "Eliminar Lote", Descripcion: "Eliminar un lote de una finca", Modulo: ModuloFincas},
	{Codigo: PermisoCrearMuestra, Nombre: "Tomar Muestra", Descripcion: "Registrar la toma de una muestra en un lote", Modulo: ModuloFincas},
	{Codigo: PermisoVerMuestras, Nombre: "Ver Muestras", Descripcion: "Listar las muestras asociadas a un lote", Modulo: ModuloFincas},
	{Codigo: PermisoSolicitarDiagnostico, Nombre: "Solicitar Diagnóstico", Descripcion: "Solicitar un diagnóstico manual para una muestra", Modulo: ModuloFincas},
	{Codigo: PermisoAceptarDiagnostico, Nombre: "Aceptar Diagnóstico", Descripcion: "Aceptar un diagnóstico pendiente", Modulo: ModuloFincas},
	{Codigo: PermisoRechazarDiagnostico, Nombre: "Rechazar Diagnóstico", Descripcion: "Rechazar un diagnóstico pendiente", Modulo: ModuloFincas},
	{Codigo: PermisoGenerarReporte, Nombre: "Generar Reporte", Descripcion: "Generar reporte de clorosis por lote", Modulo: ModuloFincas},
}
