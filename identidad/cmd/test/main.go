package main

import (
	"context"
	"fmt"
	"log"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	rbac_app "github.com/davosjar/bunna/services/identidad/internal/rbac/application"
	uc_assignpermissiontorole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/assignpermissiontorole"
	uc_assignrole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/assignrole"
	uc_createrole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/createrole"
	uc_deleterole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/deleterole"
	uc_listroles "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listroles"
	uc_revokepermissionfromrole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/revokepermissionfromrole"
	uc_revokerole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/revokerole"
	uc_updaterole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/updaterole"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	rbac_postgres "github.com/davosjar/bunna/services/identidad/internal/rbac/infrastructure/persistence/postgres"
	uc_solicitar_recuperacion "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/solicitarrecuperacion"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
	uc_changemypassword "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/changemypassword"
	uc_listblockedips "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/listblockedips"
	uc_resetpassword "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/resetpassword"
	uc_unblockip "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/unblockip"
	uc_unlockaccount "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/unlockaccount"
	uc_viewcredentials "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/viewcredentials"
	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	uc_listsessions "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/listsessions"
	uc_sesiones_login "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/login"
	uc_sesiones_logout "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
	uc_sesiones_refresh "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/refresh"
	uc_terminatesession "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/terminatesession"
	sesiones_postgres "github.com/davosjar/bunna/services/identidad/internal/sesiones/infrastructure/persistence/postgres"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	shared_idgenerator "github.com/davosjar/bunna/services/identidad/internal/shared/infrastructure/idgenerator"
	tenant_domain "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	tenant_postgres "github.com/davosjar/bunna/services/identidad/internal/tenants/infrastructure/persistence/postgres"
	uc_createuser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/createuser"
	uc_deleteuser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/deleteuser"
	uc_expeluser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/expeluser"
	uc_listusers "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/listusers"
	uc_register "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/register"
	uc_updatemyprofile "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/updatemyprofile"
	uc_updateuser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/updateuser"
	uc_viewmyprofile "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/viewmyprofile"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
	uc_solicitar "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/solicitarverificacion"
)

const (
	reset   = "\033[0m"
	rojo    = "\033[31m"
	verde   = "\033[32m"
	amar    = "\033[33m"
	azul    = "\033[34m"
	magenta = "\033[35m"
	cian    = "\033[36m"
	gris    = "\033[90m"
)

var (
	exitos = 0
	fal    = 0
)

func ok(desc string, args ...interface{}) {
	exitos++
	msg := fmt.Sprintf(desc, args...)
	fmt.Printf("  %s✅%s %s\n", verde, reset, msg)
}

func fail(desc string, args ...interface{}) {
	fal++
	msg := fmt.Sprintf(desc, args...)
	fmt.Printf("  %s❌%s %s\n", rojo, reset, msg)
}

func subtitulo(s string) {
	fmt.Printf("\n%s  ── %s ──%s\n", cian, s, reset)
}

func main() {
	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════╗%s\n", magenta, reset)
	fmt.Printf("%s║  PRUEBAS INTEGRACIÓN — TODOS LOS CASOS DE USO (30)        ║%s\n", magenta, reset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════╝%s\n", magenta, reset)

	// =========================================================================
	// INICIALIZACIÓN
	// =========================================================================
	fmt.Printf("\n%s──────────────────────────────────────────────────────%s\n", gris, reset)
	fmt.Printf("%s  INICIALIZACIÓN%s\n", amar, reset)
	fmt.Printf("%s──────────────────────────────────────────────────────%s\n", gris, reset)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("%s❌ Config: %v%s", rojo, err, reset)
	}
	fmt.Printf("  %s✓%s Configuración cargada\n", verde, reset)

	db, err := config.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatalf("%s❌ DB: %v%s", rojo, err, reset)
	}
	fmt.Printf("  %s✓%s Base de datos conectada\n", verde, reset)

	db.Migrator().DropTable(&sesiones_postgres.SesionModel{})
	db.Migrator().DropTable(&seguridad_postgres.RateLimitIPModel{})
	db.Migrator().DropTable(&seguridad_postgres.IntentoIPModel{})
	db.Migrator().DropTable(&seguridad_postgres.CredencialesModel{})
	db.Migrator().DropTable(&rbac_postgres.UsuarioTenantRolModel{})
	db.Migrator().DropTable(&rbac_postgres.UsuarioRolModel{})
	db.Migrator().DropTable(&rbac_postgres.RolPermisoModel{})
	db.Migrator().DropTable(&rbac_postgres.PermisoModel{})
	db.Migrator().DropTable(&rbac_postgres.RolModel{})
	db.Migrator().DropTable(&tenant_postgres.MembresiaModel{})
	db.Migrator().DropTable(&tenant_postgres.TenantModel{})
	db.Migrator().DropTable(&usuarios_postgres.UsuarioModel{})
	db.Migrator().DropTable(&sesiones_postgres.SesionModel{})
	db.Migrator().DropTable(&seguridad_postgres.RateLimitIPModel{})
	db.Migrator().DropTable(&seguridad_postgres.IntentoIPModel{})
	db.Migrator().DropTable(&seguridad_postgres.CredencialesModel{})
	fmt.Printf("  %s✓%s Tablas eliminadas\n", verde, reset)

	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("%s❌ Migrations: %v%s", rojo, err, reset)
	}
	fmt.Printf("  %s✓%s Migraciones ejecutadas\n", verde, reset)

	reg := registry.NewRegistry(db, cfg)
	ctx := context.Background()
	generadorID := shared_idgenerator.NewUUIDv7Generator()

	fmt.Printf("  %s✓%s Registry listo\n", verde, reset)

	// ─────────────────────────────────────────────────────────────────────────
	// SEED: PERMISOS Y ROLES DE SISTEMA
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  SEED: PERMISOS Y ROLES DE SISTEMA%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	rolRepo := rbac_postgres.NewRolRepositorio(db)
	permisoRepo := rbac_postgres.NewPermisoRepositorio(db)
	rolPermisoRepo := rbac_postgres.NewRolPermisoRepositorio(db)
	usuarioRolRepo := rbac_postgres.NewUsuarioRolRepositorio(db)

	seedSvc := rbac_app.NuevoSeedServicio(rolRepo, permisoRepo, rolPermisoRepo, generadorID)
	if err := seedSvc.Ejecutar(ctx); err != nil {
		log.Fatalf("%s❌ Seed: %v%s", rojo, err, reset)
	}
	ok("Permisos y roles de sistema sembrados")

	// ─────────────────────────────────────────────────────────────────────────
	// 1. REGISTRO DE USUARIOS (caso de uso #25)
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  1. REGISTRO DE USUARIOS (#25)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("1.1 Registros exitosos")

	regCmd1 := uc_register.ComandoRegistrarUsuario{
		Correo: "ana.lopez@example.com", Password: "SecurePass123!",
		Nombre: "Ana", Apellido: "López", Telefono: "6001112233",
	}
	_, _ = reg.IniciarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_login.ComandoIniciarSesion{
		Email: regCmd1.Correo, Password: regCmd1.Password, IPOrigen: "seed",
	})
	// Registramos con el caso de uso
	r1reg, err := reg.GetRegistrarUsuarioCasoDeUso().Ejecutar(ctx, &regCmd1)
	if err != nil {
		fail("Registro ana.lopez@example.com: %v", err)
	} else {
		ok("Registro ana.lopez → ID: %s", r1reg.UsuarioID)
	}

	regCmd2 := uc_register.ComandoRegistrarUsuario{
		Correo: "jose.martinez@demo.com", Password: "Passw0rd!",
		Nombre: "José", Apellido: "Mártinez", Telefono: "6005544332",
	}
	r2reg, err := reg.GetRegistrarUsuarioCasoDeUso().Ejecutar(ctx, &regCmd2)
	if err != nil {
		fail("Registro jose.martinez@demo.com: %v", err)
	} else {
		ok("Registro jose.martinez → ID: %s", r2reg.UsuarioID)
	}

	regCmd3 := uc_register.ComandoRegistrarUsuario{
		Correo: "carlos.ramos@test.com", Password: "ClaveFuerte99@",
		Nombre: "Carlos", Apellido: "Ramos", Telefono: "+598991234567",
	}
	r3reg, err := reg.GetRegistrarUsuarioCasoDeUso().Ejecutar(ctx, &regCmd3)
	if err != nil {
		fail("Registro carlos.ramos@test.com: %v", err)
	} else {
		ok("Registro carlos.ramos → ID: %s", r3reg.UsuarioID)
	}

	usuario1ID := r1reg.UsuarioID
	sysAdminID := usuario1ID

	subtitulo("1.2 Validaciones (deben fallar)")

	_, err = reg.GetRegistrarUsuarioCasoDeUso().Ejecutar(ctx, &uc_register.ComandoRegistrarUsuario{
		Correo: "ana.lopez@example.com", Password: "OtraPass1@", Nombre: "Ana", Apellido: "López",
	})
	if err != nil {
		ok("Email duplicado rechazado: %v", err)
	} else {
		fail("Email duplicado debería fallar")
	}

	_, err = reg.GetRegistrarUsuarioCasoDeUso().Ejecutar(ctx, &uc_register.ComandoRegistrarUsuario{
		Correo: "invalido", Password: "Passw0rd@", Nombre: "T", Apellido: "U",
	})
	if err != nil {
		ok("Email inválido rechazado: %v", err)
	} else {
		fail("Email inválido debería fallar")
	}

	_, err = reg.GetRegistrarUsuarioCasoDeUso().Ejecutar(ctx, &uc_register.ComandoRegistrarUsuario{
		Correo: "", Password: "Passw0rd@", Nombre: "T", Apellido: "U",
	})
	if err != nil {
		ok("Email vacío rechazado: %v", err)
	} else {
		fail("Email vacío debería fallar")
	}

	// Asignar rol sys_admin al primer usuario para poder probar operaciones RBAC
	rolSysAdmin, err := rolRepo.ObtenerPorNombre(ctx, rbac.RolSysAdmin)
	if err != nil {
		log.Fatalf("%s❌ Error al obtener rol sys_admin: %v%s", rojo, err, reset)
	}
	if err := usuarioRolRepo.Crear(ctx, sysAdminID, rolSysAdmin.ID); err != nil {
		log.Fatalf("%s❌ Error al asignar sys_admin: %v%s", rojo, err, reset)
	}
	ok("Rol sys_admin asignado al usuario %s", sysAdminID)

	// ─────────────────────────────────────────────────────────────────────────
	// 2. LOGIN (#26)
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  2. LOGIN (#26)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("2.1 Login exitoso")

	loginResp, err := reg.IniciarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_login.ComandoIniciarSesion{
		Email: "ana.lopez@example.com", Password: "SecurePass123!", IPOrigen: "192.168.1.100",
	})
	if err != nil {
		fail("Login ana.lopez: %v", err)
	} else {
		ok("Login exitoso → AccessToken: %s..., SesionID: %s", loginResp.AccessToken[:20], loginResp.SesionID)
		if loginResp.UsuarioID == usuario1ID {
			ok("UsuarioID coincide")
		} else {
			fail("UsuarioID no coincide")
		}
	}

	loginResp2, err := reg.IniciarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_login.ComandoIniciarSesion{
		Email: "jose.martinez@demo.com", Password: "Passw0rd!", IPOrigen: "10.0.0.50",
	})
	if err != nil {
		fail("Login jose.martinez: %v", err)
	} else {
		ok("Login jose.martinez → SesionID: %s", loginResp2.SesionID)
	}

	subtitulo("2.2 Login inválido (debe fallar)")

	_, err = reg.IniciarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_login.ComandoIniciarSesion{
		Email: "ana.lopez@example.com", Password: "WrongPassword1!", IPOrigen: "192.168.1.100",
	})
	if err != nil {
		ok("Password incorrecto rechazado")
	} else {
		fail("Password incorrecto debería fallar")
	}

	_, err = reg.IniciarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_login.ComandoIniciarSesion{
		Email: "no.existe@test.com", Password: "Passw0rd@", IPOrigen: "10.0.0.1",
	})
	if err != nil {
		ok("Email no registrado rechazado")
	} else {
		fail("Email no registrado debería fallar")
	}

	_, err = reg.IniciarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_login.ComandoIniciarSesion{
		Email: "", Password: "Passw0rd@", IPOrigen: "10.0.0.1",
	})
	if err != nil {
		ok("Email vacío rechazado")
	} else {
		fail("Email vacío debería fallar")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 3. REFRESH (#27)
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  3. REFRESH (#27)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("3.1 Refresh exitoso")
	refreshResp, err := reg.RenovarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_refresh.ComandoRenovarSesion{
		RefreshToken: loginResp.RefreshToken,
	})
	if err != nil {
		fail("Refresh falló: %v", err)
	} else {
		ok("Refresh exitoso → Nuevo AccessToken: %s...", refreshResp.AccessToken[:20])
		if refreshResp.SesionID == loginResp.SesionID {
			ok("Misma sesiónID tras refresh")
		} else {
			fail("SesionID cambió")
		}
	}
	oldRefreshToken := loginResp.RefreshToken
	firstRefreshToken := refreshResp.RefreshToken

	subtitulo("3.2 Token rotado (debe fallar)")
	_, err = reg.RenovarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_refresh.ComandoRenovarSesion{RefreshToken: oldRefreshToken})
	if err != nil {
		ok("Token rotado rechazado (detección robo): %v", err)
	} else {
		fail("Token rotado debería fallar")
	}

	subtitulo("3.3 Múltiples refrescos")
	ultimoRefresh := firstRefreshToken
	for i := 1; i <= 3; i++ {
		r, err := reg.RenovarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_refresh.ComandoRenovarSesion{RefreshToken: ultimoRefresh})
		if err != nil {
			fail("Refresh %d falló: %v", i, err)
		} else {
			ok("Refresh %d exitoso", i)
			ultimoRefresh = r.RefreshToken
		}
	}

	subtitulo("3.4 Token inválido")
	_, err = reg.RenovarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_refresh.ComandoRenovarSesion{RefreshToken: "token-invalido"})
	if err != nil {
		ok("Token inválido rechazado")
	} else {
		fail("Token inválido debería fallar")
	}
	_, err = reg.RenovarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_refresh.ComandoRenovarSesion{RefreshToken: ""})
	if err != nil {
		ok("Token vacío rechazado")
	} else {
		fail("Token vacío debería fallar")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 4. LOGOUT (#28)
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  4. LOGOUT (#28)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("4.1 Logout exitoso")
	loginForLogout, err := reg.IniciarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_login.ComandoIniciarSesion{
		Email: "carlos.ramos@test.com", Password: "ClaveFuerte99@", IPOrigen: "10.10.10.10",
	})
	if err != nil {
		fail("Login previo falló: %v", err)
	} else {
		ok("Login previo → SesionID: %s", loginForLogout.SesionID)
	}

	logoutResp, err := reg.CerrarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_logout.ComandoCerrarSesion{
		SesionID: loginForLogout.SesionID, UsuarioID: loginForLogout.UsuarioID,
	})
	if err != nil {
		fail("Logout falló: %v", err)
	} else {
		ok("Logout exitoso → Revocadas: %d", logoutResp.SesionesRevocadas)
	}

	subtitulo("4.2 Refresh post-logout (debe fallar)")
	_, err = reg.RenovarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_refresh.ComandoRenovarSesion{RefreshToken: loginForLogout.RefreshToken})
	if err != nil {
		ok("Refresh post-logout bloqueado")
	} else {
		fail("Refresh post-logout debería fallar")
	}

	subtitulo("4.3 Re-login post-logout")
	relogin, err := reg.IniciarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_login.ComandoIniciarSesion{
		Email: "carlos.ramos@test.com", Password: "ClaveFuerte99@", IPOrigen: "10.10.10.10",
	})
	if err != nil {
		fail("Re-login falló: %v", err)
	} else {
		ok("Re-login exitoso → Nueva sesión: %s", relogin.SesionID)
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 5. CRUD USUARIOS (casos #1-#5)
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  5. GESTIÓN DE USUARIOS (#1 al #5)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("5.1 Crear usuario (#1)")
	createResp, err := reg.CrearUsuarioCasoDeUso.Ejecutar(ctx, &uc_createuser.ComandoCrearUsuario{
		EjecutorID: sysAdminID,
		Nombre:     "Nuevo", Apellido: "Usuario", Correo: "nuevo@test.com",
		Password: "Passw0rd!",
	})
	if err != nil {
		fail("Crear usuario: %v", err)
	} else {
		ok("Usuario creado → ID: %s", createResp.ID)
	}
	nuevoUserID := createResp.ID

	subtitulo("5.2 Listar usuarios (#2)")
	listUsers, err := reg.ListarUsuariosCasoDeUso.Ejecutar(ctx, &uc_listusers.ComandoListarUsuarios{
		EjecutorID: sysAdminID, TenantID: "",
		Paginacion: shared_domain.Paginacion{Pagina: 1, TamanoPagina: 10},
	})
	if err != nil {
		fail("Listar usuarios: %v", err)
	} else {
		ok("Usuarios listados: %d total", listUsers.Total)
	}

	subtitulo("5.3 Modificar usuario (#3)")
	_, err = reg.ModificarUsuarioCasoDeUso.Ejecutar(ctx, &uc_updateuser.ComandoModificarUsuario{
		EjecutorID: sysAdminID, TenantID: "",
		UsuarioID: nuevoUserID, Nombre: "Modificado", Apellido: "Editado",
	})
	if err != nil {
		fail("Modificar usuario: %v", err)
	} else {
		ok("Usuario modificado exitosamente")
	}

	subtitulo("5.4 Ver perfil propio (#22)")
	perfil, err := reg.VerMiPerfilCasoDeUso.Ejecutar(ctx, &uc_viewmyprofile.ComandoVerMiPerfil{EjecutorID: nuevoUserID})
	if err != nil {
		fail("Ver perfil: %v", err)
	} else {
		ok("Perfil visto: %s %s", perfil.Nombre, perfil.Apellido)
	}

	subtitulo("5.5 Modificar perfil propio (#23)")
	_, err = reg.ModificarMiPerfilCasoDeUso.Ejecutar(ctx, &uc_updatemyprofile.ComandoModificarMiPerfil{
		EjecutorID: nuevoUserID, Nombre: "AutoEditado", Apellido: "Perfil",
	})
	if err != nil {
		fail("Modificar perfil: %v", err)
	} else {
		ok("Perfil modificado exitosamente")
	}

	subtitulo("5.6 Dar de baja usuario (#4)")
	_, err = reg.DarDeBajaUsuarioCasoDeUso.Ejecutar(ctx, &uc_deleteuser.ComandoDarDeBajaUsuario{
		EjecutorID: sysAdminID, UsuarioID: nuevoUserID,
	})
	if err != nil {
		fail("Dar de baja: %v", err)
	} else {
		ok("Usuario dado de baja exitosamente")
	}

	// Crear otro usuario para expulsión
	createResp2, err := reg.CrearUsuarioCasoDeUso.Ejecutar(ctx, &uc_createuser.ComandoCrearUsuario{
		EjecutorID: sysAdminID,
		Nombre:     "Expulsable", Apellido: "User", Correo: "expulso@test.com",
		Password: "Passw0rd!",
	})
	if err != nil {
		fail("Crear expulsable: %v", err)
	} else {
		nuevoUserID2 := createResp2.ID
		_ = nuevoUserID2
	}

	nuevoUserID2 := createResp2.ID

	subtitulo("5.7 Expulsar usuario (#5)")
	_, err = reg.ExpulsarUsuarioCasoDeUso.Ejecutar(ctx, &uc_expeluser.ComandoExpulsarUsuario{
		EjecutorID: sysAdminID, TenantID: "", UsuarioID: nuevoUserID2,
	})
	if err != nil {
		fail("Expulsar usuario: %v", err)
	} else {
		ok("Usuario expulsado exitosamente")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 6. CREDENCIALES Y SEGURIDAD (#6-#8, #20-#21, #24)
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  6. CREDENCIALES Y SEGURIDAD (#6-#8, #20-#21, #24)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("6.1 Consultar credenciales (#6)")
	_, err = reg.ConsultarCredencialesCasoDeUso.Ejecutar(ctx, &uc_viewcredentials.ComandoConsultarCredenciales{
		EjecutorID: sysAdminID, TenantID: "", UsuarioID: sysAdminID,
	})
	if err != nil {
		fail("Consultar credenciales: %v", err)
	} else {
		ok("Credenciales consultadas exitosamente")
	}

	subtitulo("6.2 Cambiar mi contraseña (#24)")
	_, err = reg.CambiarMiContrasenaCasoDeUso.Ejecutar(ctx, &uc_changemypassword.ComandoCambiarMiContrasena{
		EjecutorID: sysAdminID, PasswordActual: "SecurePass123!", NuevaPassword: "NuevaPass1@",
	})
	if err != nil {
		fail("Cambiar mi contraseña: %v", err)
	} else {
		ok("Contraseña cambiada exitosamente")
	}

	// Restaurar password para no romper tests posteriores
	reg.CambiarMiContrasenaCasoDeUso.Ejecutar(ctx, &uc_changemypassword.ComandoCambiarMiContrasena{
		EjecutorID: sysAdminID, PasswordActual: "NuevaPass1@", NuevaPassword: "SecurePass123!",
	})

	subtitulo("6.3 Resetear contraseña (#7)")
	_, err = reg.ResetearContrasenaCasoDeUso.Ejecutar(ctx, &uc_resetpassword.ComandoResetearContrasena{
		EjecutorID: sysAdminID, TenantID: "", UsuarioID: nuevoUserID,
		NuevaPassword: "ResetPass123!",
	})
	if err != nil {
		fail("Resetear contraseña: %v", err)
	} else {
		ok("Contraseña reseteada exitosamente")
	}

	subtitulo("6.4 Desbloquear cuenta (#8)")
	// Primero forzar bloqueo (repositorio directo)
	db.Exec("UPDATE credenciales_usuarios SET bloqueado_hasta = NOW() + INTERVAL '1 hour' WHERE usuario_id = ?", nuevoUserID)
	_, err = reg.DesbloquearCuentaCasoDeUso.Ejecutar(ctx, &uc_unlockaccount.ComandoDesbloquearCuenta{
		EjecutorID: sysAdminID, TenantID: "", UsuarioID: nuevoUserID,
	})
	if err != nil {
		fail("Desbloquear cuenta: %v", err)
	} else {
		ok("Cuenta desbloqueada exitosamente")
	}

	subtitulo("6.5 Listar IPs bloqueadas (#20)")
	_, err = reg.ListarIPsBloqueadasCasoDeUso.Ejecutar(ctx, &uc_listblockedips.ComandoListarIPsBloqueadas{
		EjecutorID: sysAdminID, TenantID: "",
	})
	if err != nil {
		fail("Listar IPs bloqueadas: %v", err)
	} else {
		ok("IPs bloqueadas listadas")
	}

	subtitulo("6.6 Desbloquear IP (#21)")
	// IP que nunca ha sido usada en el flujo
	_, err = reg.DesbloquearIPCasoDeUso.Ejecutar(ctx, &uc_unblockip.ComandoDesbloquearIP{
		EjecutorID: sysAdminID, TenantID: "", IP: "203.0.113.99",
	})
	if err != nil {
		ok("Desbloquear IP (sin registro activo): %v", err)
	} else {
		fail("Desbloquear IP sin registro debería fallar")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 7. SESIONES (#17-#18)
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  7. SESIONES (#17-#18)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("7.1 Listar sesiones (#17)")
	_, err = reg.ListarSesionesCasoDeUso.Ejecutar(ctx, &uc_listsessions.ComandoListarSesiones{
		EjecutorID: sysAdminID,
		Paginacion: shared_domain.Paginacion{},
	})
	if err != nil {
		fail("Listar sesiones: %v", err)
	} else {
		ok("Sesiones listadas exitosamente")
	}

	subtitulo("7.2 Forzar cierre de sesión (#18)")
	_, err = reg.ForzarCierreSesionCasoDeUso.Ejecutar(ctx, &uc_terminatesession.ComandoForzarCierreSesion{
		EjecutorID: sysAdminID, TenantID: "", SesionID: loginResp.SesionID,
	})
	if err != nil {
		fail("Forzar cierre sesión: %v", err)
	} else {
		ok("Sesión terminada exitosamente")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 8. ROLES Y PERMISOS (#9-#16)
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  8. ROLES Y PERMISOS (#9 al #16)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("8.1 Listar roles (#9)")
	listRolesResp, err := reg.ListarRolesCasoDeUso.Ejecutar(ctx, &uc_listroles.ComandoListarRoles{
		EjecutorID: sysAdminID,
		Paginacion: shared_domain.Paginacion{},
	})
	if err != nil {
		fail("Listar roles: %v", err)
	} else {
		ok("Roles listados: %d total", listRolesResp.Total)
	}

	subtitulo("8.2 Crear rol (#10)")
	createRolResp, err := reg.CrearRolCasoDeUso.Ejecutar(ctx, &uc_createrole.ComandoCrearRol{
		EjecutorID: sysAdminID, TenantID: "",
		Nombre: "RolPrueba", Descripcion: "Rol de prueba",
	})
	if err != nil {
		fail("Crear rol: %v", err)
	} else {
		ok("Rol creado → ID: %s", createRolResp.ID)
	}
	nuevoRolID := createRolResp.ID

	subtitulo("8.3 Modificar rol (#11)")
	_, err = reg.ModificarRolCasoDeUso.Ejecutar(ctx, &uc_updaterole.ComandoModificarRol{
		EjecutorID: sysAdminID, TenantID: "",
		RolID: nuevoRolID, Descripcion: "Rol modificado",
	})
	if err != nil {
		fail("Modificar rol: %v", err)
	} else {
		ok("Rol modificado exitosamente")
	}

	subtitulo("8.4 Asignar permiso a rol (#15)")
	_, err = reg.AsignarPermisoARolCasoDeUso.Ejecutar(ctx, &uc_assignpermissiontorole.ComandoAsignarPermisoARol{
		EjecutorID: sysAdminID,
		RolID:      nuevoRolID, PermisoCodigo: rbac.PermisoUsuarioCrear,
	})
	if err != nil {
		fail("Asignar permiso a rol: %v", err)
	} else {
		ok("Permiso asignado al rol exitosamente")
	}

	subtitulo("8.5 Revocar permiso de rol (#16)")
	_, err = reg.RevocarPermisoDeRolCasoDeUso.Ejecutar(ctx, &uc_revokepermissionfromrole.ComandoRevocarPermisoDeRol{
		EjecutorID: sysAdminID,
		RolID:      nuevoRolID, PermisoCodigo: rbac.PermisoUsuarioCrear,
	})
	if err != nil {
		fail("Revocar permiso de rol: %v", err)
	} else {
		ok("Permiso revocado del rol exitosamente")
	}

	subtitulo("8.6 Asignar rol a usuario (#13)")
	_, err = reg.AsignarRolCasoDeUso.Ejecutar(ctx, &uc_assignrole.ComandoAsignarRol{
		EjecutorID: sysAdminID,
		UsuarioID:  nuevoUserID, RolID: nuevoRolID,
	})
	if err != nil {
		fail("Asignar rol a usuario: %v", err)
	} else {
		ok("Rol asignado a usuario exitosamente")
	}

	subtitulo("8.7 Revocar rol de usuario (#14)")
	_, err = reg.RevocarRolCasoDeUso.Ejecutar(ctx, &uc_revokerole.ComandoRevocarRol{
		EjecutorID: sysAdminID,
		UsuarioID:  nuevoUserID, RolID: nuevoRolID,
	})
	if err != nil {
		fail("Revocar rol de usuario: %v", err)
	} else {
		ok("Rol revocado de usuario exitosamente")
	}

	subtitulo("8.8 Eliminar rol (#12)")
	_, err = reg.EliminarRolCasoDeUso.Ejecutar(ctx, &uc_deleterole.ComandoEliminarRol{
		EjecutorID: sysAdminID, TenantID: "",
		RolID: nuevoRolID,
	})
	if err != nil {
		fail("Eliminar rol: %v", err)
	} else {
		ok("Rol eliminado exitosamente")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 9. TENANT (#19)
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  9. TENANT (#19)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	// Crear un tenant primero (directo por repo)
	tenantRepo := tenant_postgres.NewTenantRepositorio(db)
	t1, _ := tenant_domain.NuevoTenant("tenant-test-id", "TenantPrueba", "tenant-prueba")
	tenantNuevo, err := tenantRepo.Crear(ctx, t1)
	_ = tenantNuevo
	tid1, _ := generadorID.NextID(ctx)
	_, _ = tenant_domain.NuevoTenant(tid1, "Tenant Prueba", "tenant-prueba")
	tid2, _ := generadorID.NextID(ctx)
	_, _ = tenant_domain.NuevoTenant(tid2, "Tenant Test", "tenant-test")

	// Esperamos los IDs
	tenantID1, _ := generadorID.NextID(ctx)
	tenantID2, _ := generadorID.NextID(ctx)

	tenant1, _ := tenant_domain.NuevoTenant(tenantID1, "Tenant Alpha", "tenant-alpha")
	tenant2, _ := tenant_domain.NuevoTenant(tenantID2, "Tenant Beta", "tenant-beta")

	tenantRepo.Crear(ctx, tenant1)
	tenantRepo.Crear(ctx, tenant2)

	// ─────────────────────────────────────────────────────────────────────────
	// 10. VERIFICACIÓN DE CORREO (#29)
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  10. VERIFICACIÓN DE CORREO (#29)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("10.1 Solicitar verificación")
	_, err = reg.SolicitarVerificacionCasoDeUso.Ejecutar(ctx, &uc_solicitar.ComandoSolicitarVerificacion{
		UsuarioID: sysAdminID,
	})
	if err != nil {
		fail("Solicitar verificación: %v", err)
	} else {
		ok("Verificación solicitada (email enviado async)")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 11. RECUPERACIÓN DE CONTRASEÑA (#30)
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  11. RECUPERACIÓN DE CONTRASEÑA (#30)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("11.1 Solicitar recuperación")
	_, err = reg.SolicitarRecuperacionCasoDeUso.Ejecutar(ctx, &uc_solicitar_recuperacion.ComandoSolicitarRecuperacion{
		Email: "ana.lopez@example.com",
	})
	if err != nil {
		fail("Solicitar recuperación: %v", err)
	} else {
		ok("Recuperación solicitada (email enviado async)")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 12. FLUJO COMPLETO
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", verde, reset)
	fmt.Printf("%s  12. FLUJO COMPLETO%s\n", verde, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", verde, reset)

	subtitulo("12.1 Registro")
	fullUser, err := reg.GetRegistrarUsuarioCasoDeUso().Ejecutar(ctx, &uc_register.ComandoRegistrarUsuario{
		Correo: "flujo.completo@test.com", Password: "FlujoCompleto99@",
		Nombre: "Flujo", Apellido: "Completo", Telefono: "600000000",
	})
	if err != nil {
		fail("Registro: %v", err)
	} else {
		ok("Usuario registrado: %s", fullUser.UsuarioID)
	}

	subtitulo("12.2 Login")
	fullLogin, err := reg.IniciarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_login.ComandoIniciarSesion{
		Email: "flujo.completo@test.com", Password: "FlujoCompleto99@", IPOrigen: "192.168.99.99",
	})
	if err != nil {
		fail("Login: %v", err)
	} else {
		ok("Login exitoso → SesionID: %s", fullLogin.SesionID)
	}

	subtitulo("12.3 Refresh")
	fullRefresh, err := reg.RenovarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_refresh.ComandoRenovarSesion{
		RefreshToken: fullLogin.RefreshToken,
	})
	if err != nil {
		fail("Refresh: %v", err)
	} else {
		ok("Refresh exitoso → misma sesión: %v", fullRefresh.SesionID == fullLogin.SesionID)
	}

	subtitulo("12.4 Logout")
	_, err = reg.CerrarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_logout.ComandoCerrarSesion{
		SesionID: fullLogin.SesionID, UsuarioID: fullLogin.UsuarioID,
	})
	if err != nil {
		fail("Logout: %v", err)
	} else {
		ok("Logout exitoso")
	}

	subtitulo("12.5 Re-login")
	fullRelogin, err := reg.IniciarSesionCasoDeUso.Ejecutar(ctx, uc_sesiones_login.ComandoIniciarSesion{
		Email: "flujo.completo@test.com", Password: "FlujoCompleto99@", IPOrigen: "192.168.99.99",
	})
	if err != nil {
		fail("Re-login: %v", err)
	} else {
		ok("Re-login exitoso → nueva sesión: %s", fullRelogin.SesionID)
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 13. CASOS BORDE — PERMISO DENEGADO
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  13. PERMISO DENEGADO (usuario sin rol)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("13.1 Listar usuarios sin permiso (debe fallar)")
	_, err = reg.ListarUsuariosCasoDeUso.Ejecutar(ctx, &uc_listusers.ComandoListarUsuarios{
		EjecutorID: nuevoUserID,
		Paginacion: shared_domain.Paginacion{},
	})
	if err != nil {
		ok("Listar usuarios sin permiso denegado: %v", err)
	} else {
		fail("Listar usuarios sin permiso debería fallar")
	}

	subtitulo("13.2 Crear usuario sin permiso (debe fallar)")
	_, err = reg.CrearUsuarioCasoDeUso.Ejecutar(ctx, &uc_createuser.ComandoCrearUsuario{
		EjecutorID: nuevoUserID,
		Nombre:     "Sin", Apellido: "Permiso", Correo: "sinpermiso@test.com",
		Password: "Passw0rd!",
	})
	if err != nil {
		ok("Crear usuario sin permiso denegado: %v", err)
	} else {
		fail("Crear usuario sin permiso debería fallar")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// REPORTE FINAL
	// ─────────────────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", magenta, reset)
	fmt.Printf("%s  REPORTE FINAL — 30 CASOS DE USO TESTEADOS%s\n", magenta, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", magenta, reset)

	total := exitos + fal
	fmt.Printf("\n  %s📊 Total de pruebas: %d%s\n", cian, total, reset)
	fmt.Printf("  %s✅ Exitosas: %d%s\n", verde, exitos, reset)
	if fal > 0 {
		fmt.Printf("  %s❌ Fallidas: %d%s\n", rojo, fal, reset)
		fmt.Printf("\n  %s⚠️  Hubo %d prueba(s) fallida(s). Revisar los ❌ arriba.%s\n\n", amar, fal, reset)
	} else {
		fmt.Printf("\n  %s🎉 ¡Todas las pruebas de integración pasaron!%s\n\n", verde, reset)
	}
}
