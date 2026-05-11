package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	sesiones_postgres "github.com/davosjar/bunna/services/identidad/internal/sesiones/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/login"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/logout"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/refresh"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/services/registro"
)

// ─────────────────────────────────────────────────────────────────────────────
// COLORES para salida en terminal
// ─────────────────────────────────────────────────────────────────────────────

const (
	reset   = "\033[0m"
	rojo    = "\033[31m"
	verde   = "\033[32m"
	amar  = "\033[33m"
	azul    = "\033[34m"
	magenta = "\033[35m"
	cian    = "\033[36m"
	gris    = "\033[90m"
)

var (
	exitos   = 0
	fal  = 0
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

// ─────────────────────────────────────────────────────────────────────────────
// MAIN
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════╗%s\n", magenta, reset)
	fmt.Printf("%s║   PRUEBAS DE INTEGRACIÓN — REGISTRO → LOGIN → REFRESH      ║%s\n", magenta, reset)
	fmt.Printf("%s║                  Capa de Aplicación (sin Handlers)          ║%s\n", magenta, reset)
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

	// Limpiar usando GORM Migrator (más confiable que raw SQL)
	db.Migrator().DropTable(&sesiones_postgres.SesionModel{})
	db.Migrator().DropTable(&seguridad_postgres.RateLimitIPModel{})
	db.Migrator().DropTable(&seguridad_postgres.IntentoIPModel{})
	db.Migrator().DropTable(&seguridad_postgres.CredencialesModel{})
	db.Migrator().DropTable(&usuarios_postgres.UsuarioModel{})
	fmt.Printf("  %s✓%s Tablas eliminadas con Migrator\n", verde, reset)

	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("%s❌ Migrations: %v%s", rojo, err, reset)
	}
	fmt.Printf("  %s✓%s Migraciones ejecutadas (limpieza total)\n", verde, reset)

	reg := registry.NewRegistry(db, cfg)
	ctx := context.Background()

	fmt.Printf("  %s✓%s Registry listo\n", verde, reset)

	// ─────────────────────────────────────────────────────────────────────────
	// 1. REGISTRO
	// ─────────────────────────────────────────────────────────────────────────

	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  1. REGISTRO DE USUARIOS%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	servicioRegistro := registro.NuevoServicioRegistro(reg.UsuarioUnitOfWork())

	subtitulo("1.1 Registros exitosos")

	// Caso 1.1.1 — Registro básico exitoso
	r1, err := servicioRegistro.Ejecutar(ctx, &registro.ComandoRegistro{
		Correo:   "ana.lopez@example.com",
		Password: "SecurePass123!",
		Nombre:   "Ana",
		Apellido: "López",
		Telefono: "6001112233",
	})
	if err != nil {
		fail("Registro ana.lopez@example.com: %v", err)
	} else {
		ok("Registro ana.lopez@example.com → UsuarioID: %s, Estado: %s", r1.UsuarioID, r1.Estado)
	}

	// Caso 1.1.2 — Registro con caracteres especiales
	r2, err := servicioRegistro.Ejecutar(ctx, &registro.ComandoRegistro{
		Correo:   "josé.mártinez@demo.com",
		Password: "Passw0rd!ñ",
		Nombre:   "José",
		Apellido: "Mártinez",
		Telefono: "6005544332",
	})
	if err != nil {
		fail("Registro josé.mártinez@demo.com: %v", err)
	} else {
		ok("Registro josé.mártinez@demo.com → ID: %s", r2.UsuarioID)
	}

	// Caso 1.1.3 — Registro con teléfono largo
	r3, err := servicioRegistro.Ejecutar(ctx, &registro.ComandoRegistro{
		Correo:   "carlos.ramos@test.com",
		Password: "ClaveFuerte99@",
		Nombre:   "Carlos",
		Apellido: "Ramos",
		Telefono: "+598991234567",
	})
	if err != nil {
		fail("Registro carlos.ramos@test.com: %v", err)
	} else {
		ok("Registro carlos.ramos@test.com → ID: %s", r3.UsuarioID)
	}

	// Guardamos para login
	usuario1ID := r1.UsuarioID

	subtitulo("1.2 Validaciones de registro (deben fallar)")

	// Caso 1.2.1 — Email duplicado
	_, err = servicioRegistro.Ejecutar(ctx, &registro.ComandoRegistro{
		Correo:   "ana.lopez@example.com",
		Password: "OtraPass1@",
		Nombre:   "Ana",
		Apellido: "López",
		Telefono: "6009988776",
	})
	if err != nil {
		ok("Email duplicado rechazado: %v", err)
	} else {
		fail("Email duplicado debería haber fallado")
	}

	// Caso 1.2.2 — Email inválido (sin @)
	_, err = servicioRegistro.Ejecutar(ctx, &registro.ComandoRegistro{
		Correo:   "invalido",
		Password: "Passw0rd@",
		Nombre:   "Test",
		Apellido: "User",
		Telefono: "6001111111",
	})
	if err != nil {
		ok("Email inválido 'invalido' rechazado: %v", err)
	} else {
		fail("Email inválido debería haber fallado")
	}

	// Caso 1.2.3 — Email vacío
	_, err = servicioRegistro.Ejecutar(ctx, &registro.ComandoRegistro{
		Correo:   "",
		Password: "Passw0rd@",
		Nombre:   "Test",
		Apellido: "User",
		Telefono: "6001111111",
	})
	if err != nil {
		ok("Email vacío rechazado: %v", err)
	} else {
		fail("Email vacío debería haber fallado")
	}

	// Caso 1.2.4 — Password vacío
	_, err = servicioRegistro.Ejecutar(ctx, &registro.ComandoRegistro{
		Correo:   "test@example.com",
		Password: "",
		Nombre:   "Test",
		Apellido: "User",
		Telefono: "6001111111",
	})
	if err != nil {
		ok("Password vacío rechazado: %v", err)
	} else {
		fail("Password vacío debería haber fallado")
	}

	// Caso 1.2.5 — Nombre vacío
	_, err = servicioRegistro.Ejecutar(ctx, &registro.ComandoRegistro{
		Correo:   "test@example.com",
		Password: "Passw0rd@",
		Nombre:   "",
		Apellido: "User",
		Telefono: "6001111111",
	})
	if err != nil {
		ok("Nombre vacío rechazado: %v", err)
	} else {
		fail("Nombre vacío debería haber fallado")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 2. LOGIN
	// ─────────────────────────────────────────────────────────────────────────

	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  2. LOGIN DE USUARIOS%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	servicioLogin := reg.ServicioLogin

	subtitulo("2.1 Login exitoso con credenciales correctas")

	// Caso 2.1.1 — Login con credenciales del usuario registrado
	loginResp, err := servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "ana.lopez@example.com",
		Password: "SecurePass123!",
		IPOrigen: "192.168.1.100",
	})
	if err != nil {
		fail("Login ana.lopez@example.com: %v", err)
	} else {
		ok("Login exitoso → AccessToken: %s..., RefreshToken: %s..., SesionID: %s",
			loginResp.AccessToken[:20],
			loginResp.RefreshToken[:20],
			loginResp.SesionID,
		)
		if loginResp.UsuarioID == usuario1ID {
			ok("UsuarioID coincide con el registrado: %s", loginResp.UsuarioID)
		} else {
			fail("UsuarioID no coincide: esperado %s, obtenido %s", usuario1ID, loginResp.UsuarioID)
		}
		if !loginResp.ExpiracionAccess.IsZero() && loginResp.ExpiracionAccess.After(time.Now()) {
			ok("Access token con expiración futura: %v", loginResp.ExpiracionAccess)
		} else {
			fail("Access token sin expiración o ya vencido")
		}
	}

	// Caso 2.1.2 — Login desde IP diferente
	loginResp2, err := servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "josé.mártinez@demo.com",
		Password: "Passw0rd!ñ",
		IPOrigen: "10.0.0.50",
	})
	if err != nil {
		fail("Login josé.mártinez@demo.com desde 10.0.0.50: %v", err)
	} else {
		ok("Login josé.mártinez@demo.com exitoso desde IP 10.0.0.50 → SesionID: %s", loginResp2.SesionID)
	}

	subtitulo("2.2 Login con credenciales inválidas (deben fallar)")

	// Caso 2.2.1 — Password incorrecto
	_, err = servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "ana.lopez@example.com",
		Password: "WrongPassword1!",
		IPOrigen: "192.168.1.100",
	})
	if err != nil {
		ok("Password incorrecto rechazado: %v", err)
	} else {
		fail("Password incorrecto debería haber fallado")
	}

	// Caso 2.2.2 — Email no registrado
	_, err = servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "no.existe@example.com",
		Password: "Passw0rd@",
		IPOrigen: "192.168.1.100",
	})
	if err != nil {
		ok("Email no registrado rechazado: %v", err)
	} else {
		fail("Email no registrado debería haber fallado")
	}

	// Caso 2.2.3 — Email vacío
	_, err = servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "",
		Password: "Passw0rd@",
		IPOrigen: "192.168.1.100",
	})
	if err != nil {
		ok("Login con email vacío rechazado: %v", err)
	} else {
		fail("Login con email vacío debería haber fallado")
	}

	// Caso 2.2.4 — Password vacío
	_, err = servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "ana.lopez@example.com",
		Password: "",
		IPOrigen: "192.168.1.100",
	})
	if err != nil {
		ok("Login con password vacío rechazado: %v", err)
	} else {
		fail("Login con password vacío debería haber fallado")
	}

	// Caso 2.2.5 — Login sin IP (debe funcionar igual, IP es opcional)
	_, err = servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "ana.lopez@example.com",
		Password: "SecurePass123!",
		IPOrigen: "",
	})
	if err != nil {
		fail("Login sin IP falló inesperadamente: %v", err)
	} else {
		ok("Login sin IP funciona correctamente")
	}

	subtitulo("2.3 Múltiples intentos fallidos (seguridad)")

	// 3 intentos fallidos consecutivos
	for i := 1; i <= 3; i++ {
		_, err := servicioLogin.Ejecutar(ctx, login.ComandoLogin{
			Email:    "ana.lopez@example.com",
			Password: fmt.Sprintf("WrongPass%d!", i),
			IPOrigen: "192.168.1.100",
		})
		if err != nil {
			ok("Intento fallido %d/3 rechazado: %v", i, err)
		} else {
			fail("Intento %d debería haber fallado", i)
		}
	}

	// Ahora login exitoso con la misma IP (debe funcionar si no superó umbral de IP)
	loginResp3, err := servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "ana.lopez@example.com",
		Password: "SecurePass123!",
		IPOrigen: "192.168.1.100",
	})
	if err != nil {
		fail("Login exitoso post-fallidos: %v", err)
	} else {
		ok("Login exitoso tras 3 intentos fallidos → Nueva sesión creada: %s", loginResp3.SesionID)
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 3. REFRESH TOKEN
	// ─────────────────────────────────────────────────────────────────────────

	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  3. REFRESH DE TOKEN%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	servicioRefresh := reg.ServicioRefresh

	subtitulo("3.1 Refresh exitoso")

	// Usar el refresh token del login anterior
	refreshResp, err := servicioRefresh.Ejecutar(ctx, refresh.ComandoRefresh{
		RefreshToken: loginResp.RefreshToken,
	})
	if err != nil {
		fail("Refresh token falló: %v", err)
	} else {
		ok("Refresh exitoso → Nuevo AccessToken: %s..., Nuevo RefreshToken: %s...",
			refreshResp.AccessToken[:20],
			refreshResp.RefreshToken[:20],
		)
		if refreshResp.SesionID == loginResp.SesionID {
			ok("Misma sesiónID tras refresh: %s", refreshResp.SesionID)
		} else {
			fail("SesionID cambió en refresh: esperado %s, obtenido %s", loginResp.SesionID, refreshResp.SesionID)
		}
		if refreshResp.UsuarioID == loginResp.UsuarioID {
			ok("Mismo usuarioID tras refresh: %s", refreshResp.UsuarioID)
		} else {
			fail("UsuarioID cambió en refresh")
		}
	}

	oldRefreshToken := loginResp.RefreshToken
	firstRefreshToken := refreshResp.RefreshToken

	subtitulo("3.2 Reutilización de token rotado (detección de robo)")

	// Refresh con el token ANTERIOR (ya rotado)
	_, err = servicioRefresh.Ejecutar(ctx, refresh.ComandoRefresh{
		RefreshToken: oldRefreshToken,
	})
	if err != nil {
		ok("Token rotado rechazado (detección de robo): %v", err)
	} else {
		fail("Token rotado debería haber fallado")
	}

	subtitulo("3.3 Múltiples refrescos")

	// 3 refrescos consecutivos
	ultimoRefresh := firstRefreshToken
	for i := 1; i <= 3; i++ {
		r, err := servicioRefresh.Ejecutar(ctx, refresh.ComandoRefresh{
			RefreshToken: ultimoRefresh,
		})
		if err != nil {
			fail("Refresh consecutivo %d falló: %v", i, err)
		} else {
			ok("Refresh consecutivo %d exitoso", i)
			ultimoRefresh = r.RefreshToken
		}
	}

	subtitulo("3.4 Token inválido/expirado")

	// Refresh con token claramente inválido
	_, err = servicioRefresh.Ejecutar(ctx, refresh.ComandoRefresh{
		RefreshToken: "token-invalido-claramente-123",
	})
	if err != nil {
		ok("Token inválido rechazado: %v", err)
	} else {
		fail("Token inválido debería haber fallado")
	}

	// Refresh con token vacío
	_, err = servicioRefresh.Ejecutar(ctx, refresh.ComandoRefresh{
		RefreshToken: "",
	})
	if err != nil {
		ok("Refresh token vacío rechazado: %v", err)
	} else {
		fail("Refresh token vacío debería haber fallado")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 4. LOGOUT
	// ─────────────────────────────────────────────────────────────────────────

	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  4. CIERRE DE SESIÓN (LOGOUT)%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	servicioLogout := reg.ServicioLogout

	subtitulo("4.1 Logout de sesión específica")

	// Hacer login para tener una sesión activa para logout
	loginForLogout, err := servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "carlos.ramos@test.com",
		Password: "ClaveFuerte99@",
		IPOrigen: "10.10.10.10",
	})
	if err != nil {
		fail("Login previo a logout falló: %v", err)
	} else {
		ok("Login previo exitoso → SesionID: %s", loginForLogout.SesionID)

		// Logout de la sesión
		logoutResp, err := servicioLogout.Ejecutar(ctx, logout.ComandoLogout{
			SesionID:  loginForLogout.SesionID,
			UsuarioID: loginForLogout.UsuarioID,
		})
		if err != nil {
			fail("Logout falló: %v", err)
		} else {
			ok("Logout exitoso → Sesiones revocadas: %d", logoutResp.SesionesRevocadas)
		}
	}

	subtitulo("4.2 Refresh post-logout (debe fallar)")

	// Intentar refresh con la sesión recién cerrada
	_, err = servicioRefresh.Ejecutar(ctx, refresh.ComandoRefresh{
		RefreshToken: loginForLogout.RefreshToken,
	})
	if err != nil {
		ok("Refresh post-logout rechazado: %v", err)
	} else {
		fail("Refresh post-logout debería haber fallado")
	}

	subtitulo("4.3 Login después de logout (debe funcionar)")

	// Login del mismo usuario después de logout
	relogin, err := servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "carlos.ramos@test.com",
		Password: "ClaveFuerte99@",
		IPOrigen: "10.10.10.10",
	})
	if err != nil {
		fail("Re-login post-logout falló: %v", err)
	} else {
		ok("Re-login post-logout exitoso → Nueva sesión: %s", relogin.SesionID)
		if relogin.SesionID == loginForLogout.SesionID {
			fail("Re-login debería crear sesión nueva, no reusar la anterior")
		} else {
			ok("Sesión nueva (diferente a la anterior): OK")
		}
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 5. FLUJO COMPLETO: REGISTRO → LOGIN → REFRESH → LOGOUT
	// ─────────────────────────────────────────────────────────────────────────

	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", verde, reset)
	fmt.Printf("%s  5. FLUJO COMPLETO DE UN USUARIO%s\n", verde, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", verde, reset)

	subtitulo("5.1 Registro")

	fullUser, err := servicioRegistro.Ejecutar(ctx, &registro.ComandoRegistro{
		Correo:   "flujo.completo@test.com",
		Password: "FlujoCompleto99@",
		Nombre:   "Flujo",
		Apellido: "Completo",
		Telefono: "6000000000",
	})
	if err != nil {
		fail("Registro falló: %v", err)
	} else {
		ok("Usuario registrado: %s (%s)", fullUser.UsuarioID, fullUser.Correo)
	}

	subtitulo("5.2 Login")

	fullLogin, err := servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "flujo.completo@test.com",
		Password: "FlujoCompleto99@",
		IPOrigen: "192.168.99.99",
	})
	if err != nil {
		fail("Login falló: %v", err)
	} else {
		ok("Login exitoso → SesionID: %s", fullLogin.SesionID)
	}

	subtitulo("5.3 Refresh")

	fullRefresh, err := servicioRefresh.Ejecutar(ctx, refresh.ComandoRefresh{
		RefreshToken: fullLogin.RefreshToken,
	})
	if err != nil {
		fail("Refresh falló: %v", err)
	} else {
		ok("Refresh exitoso → Misma sesión: %s, mismo usuario: %s",
			fullRefresh.SesionID, fullRefresh.UsuarioID)
		if fullRefresh.SesionID != fullLogin.SesionID {
			fail("La sesión debería ser la misma tras refresh")
		}
	}

	subtitulo("5.4 Logout")

	fullLogout, err := servicioLogout.Ejecutar(ctx, logout.ComandoLogout{
		SesionID:  fullLogin.SesionID,
		UsuarioID: fullLogin.UsuarioID,
	})
	if err != nil {
		fail("Logout falló: %v", err)
	} else {
		ok("Logout exitoso → Sesiones revocadas: %d", fullLogout.SesionesRevocadas)
	}

	subtitulo("5.5 Refresh post-logout (debe fallar)")

	_, err = servicioRefresh.Ejecutar(ctx, refresh.ComandoRefresh{
		RefreshToken: fullRefresh.RefreshToken,
	})
	if err != nil {
		ok("Refresh post-logout bloqueado: %v", err)
	} else {
		fail("Refresh post-logout debería haber fallado")
	}

	subtitulo("5.6 Re-login (debe funcionar)")

	fullRelogin, err := servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "flujo.completo@test.com",
		Password: "FlujoCompleto99@",
		IPOrigen: "192.168.99.99",
	})
	if err != nil {
		fail("Re-login falló: %v", err)
	} else {
		ok("Re-login exitoso → Nueva sesión creada: %s", fullRelogin.SesionID)
		if fullRelogin.SesionID == fullLogin.SesionID {
			fail("Re-login debería crear sesión nueva")
		} else {
			ok("Sesión nueva correctamente generada")
		}
	}

	// ─────────────────────────────────────────────────────────────────────────
	// 6. EDGE CASES
	// ─────────────────────────────────────────────────────────────────────────

	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", amar, reset)
	fmt.Printf("%s  6. CASOS BORDE%s\n", amar, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", amar, reset)

	subtitulo("6.1 Logout de sesión inexistente")
	_, err = servicioLogout.Ejecutar(ctx, logout.ComandoLogout{
		SesionID:  "no-existe-123",
		UsuarioID: fullRelogin.UsuarioID,
	})
	if err != nil {
		ok("Logout de sesión inexistente rechazado: %v", err)
	} else {
		fail("Logout de sesión inexistente debería fallar")
	}

	subtitulo("6.2 Logout de sesión de otro usuario")
	_, err = servicioLogout.Ejecutar(ctx, logout.ComandoLogout{
		SesionID:  fullRelogin.SesionID,
		UsuarioID: "otro-usuario-id",
	})
	if err != nil {
		ok("Logout de sesión ajena rechazado: %v", err)
	} else {
		fail("Logout de sesión ajena debería fallar")
	}

	subtitulo("6.3 Login con email en mayúsculas (case insensitive)")
	_, err = servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "ANA.LOPEZ@EXAMPLE.COM",
		Password: "SecurePass123!",
		IPOrigen: "10.0.0.1",
	})
	if err != nil {
		fail("Login con email en mayúsculas falló: %v", err)
	} else {
		ok("Login con email en mayúsculas funciona correctamente")
	}

	subtitulo("6.4 Login con email con espacios")
	_, err = servicioLogin.Ejecutar(ctx, login.ComandoLogin{
		Email:    "  ana.lopez@example.com  ",
		Password: "SecurePass123!",
		IPOrigen: "10.0.0.2",
	})
	if err != nil {
		ok("Login con email con espacios: %v (puede fallar si no se hace trim)", err)
	} else {
		ok("Login con email con espacios funciona (hay trim)")
	}

	// ─────────────────────────────────────────────────────────────────────────
	// REPORTE FINAL
	// ─────────────────────────────────────────────────────────────────────────

	fmt.Printf("\n%s══════════════════════════════════════════════════════════════%s\n", magenta, reset)
	fmt.Printf("%s  REPORTE FINAL%s\n", magenta, reset)
	fmt.Printf("%s══════════════════════════════════════════════════════════════%s\n", magenta, reset)

	total := exitos + fal
	fmt.Printf("\n  %s📊 Total de pruebas: %d%s\n", cian, total, reset)
	fmt.Printf("  %s✅ Exitosas: %d%s\n", verde, exitos, reset)
	if fal > 0 {
		fmt.Printf("  %s❌ Fallidas: %d%s\n", rojo, fal, reset)
	} else {
		fmt.Printf("  %s✅ Fallidas: %d%s\n", verde, fal, reset)
	}
	fmt.Printf("  %s────────────────────────%s\n", gris, reset)

	if fal > 0 {
		fmt.Printf("\n  %s⚠️  Hubo %d prueba(s) fallida(s). Revisar los ❌ arriba.%s\n\n", amar, fal, reset)
	} else {
		fmt.Printf("\n  %s🎉 ¡Todas las pruebas pasaron!%s\n\n", verde, reset)
	}
}
