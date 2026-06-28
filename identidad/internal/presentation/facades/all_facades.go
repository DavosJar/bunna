package facades

import "github.com/davosjar/bunna/services/identidad/internal/registry"

type AllFacades struct {
	Auth         AuthFacade
	Usuario      UsuarioFacade
	Sesion       SesionFacade
	Seguridad    SeguridadFacade
	Rbac         RbacFacade
	Tenant       TenantFacade
	Verificacion VerificacionFacade
	Recuperacion RecuperacionFacade
	Invitacion   InvitacionFacade
}

func NewAllFacades(reg *registry.Registry) *AllFacades {
	return &AllFacades{
		Auth: NewAuthFacade(
			reg.GetRegistrarUsuarioCasoDeUso(),
			reg.SolicitarVerificacionCasoDeUso,
			reg.IniciarSesionCasoDeUso,
			reg.RenovarSesionCasoDeUso,
			reg.CerrarSesionCasoDeUso,
			reg.CambiarTenantCasoDeUso,
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
			reg.ListarPermisosCasoDeUso,
			reg.ListarMisPermisosCasoDeUso,
			reg.CrearRolCasoDeUso,
			reg.ModificarRolCasoDeUso,
			reg.EliminarRolCasoDeUso,
			reg.AsignarRolCasoDeUso,
			reg.RevocarRolCasoDeUso,
			reg.AsignarPermisoARolCasoDeUso,
			reg.RevocarPermisoDeRolCasoDeUso,
		),
		Tenant: NewTenantFacade(
			reg.ListarMisTenantsCasoDeUso,
			reg.MembresiaRepository(),
			reg.UsuarioTenantRolRepositorio(),
		),
		Verificacion: NewVerificacionFacade(
			reg.SolicitarVerificacionCasoDeUso,
			reg.ConfirmarVerificacionCasoDeUso,
			reg.ReenviarVerificacionCasoDeUso,
		),
		Recuperacion: NewRecuperacionFacade(
			reg.SolicitarRecuperacionCasoDeUso,
			reg.ValidarTokenRecuperacionCasoDeUso,
			reg.ConfirmarRecuperacionCasoDeUso,
		),
		Invitacion: NewInvitacionFacade(
			reg.CrearInvitacionCasoDeUso,
			reg.AceptarInvitacionCasoDeUso,
			reg.ObtenerInvitacionCasoDeUso,
			reg.ListarInvitacionesCasoDeUso,
			reg.ReenviarInvitacionCasoDeUso,
			reg.EliminarInvitacionCasoDeUso,
		),
	}
}
