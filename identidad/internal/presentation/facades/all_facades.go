package facades

import "github.com/davosjar/bunna/services/identidad/internal/registry"

// AllFacades agrupa todas las fachadas de la capa de presentación.
type AllFacades struct {
	Auth         AuthFacade
	Usuario      UsuarioFacade
	Sesion       SesionFacade
	Seguridad    SeguridadFacade
	Rbac         RbacFacade
	Tenant       TenantFacade
	Verificacion VerificacionFacade
	Recuperacion RecuperacionFacade
}

// NewAllFacades construye todas las fachadas a partir del Registry.
func NewAllFacades(reg *registry.Registry) *AllFacades {
	return &AllFacades{
		Auth: NewAuthFacade(
			reg.GetServicioRegistro(),
			reg.IniciarSesionCasoDeUso,
			reg.RenovarSesionCasoDeUso,
			reg.CerrarSesionCasoDeUso,
		),
		Usuario: NewUsuarioFacade(
			reg.CrearUsuarioCasoDeUso,
			reg.ListarUsuariosCasoDeUso,
			reg.ModificarUsuarioCasoDeUso,
			reg.DarDeBajaUsuarioCasoDeUso,
			reg.ExpulsarUsuarioCasoDeUso,
			reg.VerMiPerfilCasoDeUso,
			reg.ModificarMiPerfilCasoDeUso,
		),
		Sesion: NewSesionFacade(
			reg.ListarSesionesCasoDeUso,
			reg.ForzarCierreSesionCasoDeUso,
		),
		Seguridad: NewSeguridadFacade(
			reg.CambiarMiContrasenaCasoDeUso,
			reg.ResetearContrasenaCasoDeUso,
			reg.DesbloquearCuentaCasoDeUso,
			reg.ListarIPsBloqueadasCasoDeUso,
			reg.DesbloquearIPCasoDeUso,
			reg.ConsultarCredencialesCasoDeUso,
		),
		Rbac: NewRbacFacade(
			reg.ListarRolesCasoDeUso,
			reg.CrearRolCasoDeUso,
			reg.ModificarRolCasoDeUso,
			reg.EliminarRolCasoDeUso,
			reg.AsignarRolCasoDeUso,
			reg.RevocarRolCasoDeUso,
			reg.AsignarPermisoARolCasoDeUso,
			reg.RevocarPermisoDeRolCasoDeUso,
		),
		Tenant: NewTenantFacade(
			reg.ConfigurarTenantCasoDeUso,
		),
		Verificacion: NewVerificacionFacade(
			reg.VerificarCorreoCasoDeUso,
		),
		Recuperacion: NewRecuperacionFacade(
			reg.RecuperarContrasenaCasoDeUso,
		),
	}
}
