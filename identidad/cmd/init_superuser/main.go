// Package main implementa el comando CLI `init_superuser` responsable de
// crear el primer sys_admin del sistema.
//
// Reemplaza al script `cmd/seed` (deprecado). A diferencia de aquel, este
// comando:
//   - Usa los repositorios y servicios del registry (sin SQL directo ni
//     modelos duplicados) — ver ADR-001 / CON-001 / VC-004.
//   - Crea el sys_admin en una transacción atómica de 3 tablas (usuarios,
//     credenciales_usuarios, usuario_roles) con rollback automático.
//   - Es idempotente: si ya existe un usuario con rol sys_admin, no hace
//     nada y muestra el ID del sys_admin preexistente.
//   - Soporta modo `--full` (no interactivo) y modo interactivo.
//
// La lógica de negocio vive en `internal/bootstrap/application/usecase`
// (testeable sin I/O de consola). Este binario es una capa fina de UI:
// parsea flags, presenta prompts interactivos, muestra banners y formatea
// errores.
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/mail"
	"os"
	"strings"

	"github.com/davosjar/bunna/services/identidad/internal/bootstrap/application/usecase"
	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
	"github.com/davosjar/bunna/services/identidad/internal/shared/application"
)

// Códigos de color ANSI (estilo cmd/test y cmd/seed).
const (
	reset   = "\033[0m"
	rojo    = "\033[31m"
	verde   = "\033[32m"
	amar    = "\033[33m"
	magenta = "\033[35m"
	cian    = "\033[36m"
)

// helpText se imprime con -h/--help (GUD-001).
const helpText = `init_superuser — Creación del primer sys_admin

Uso:
  go run ./cmd/init_superuser [flags]

Sin flags, entra en modo interactivo (solicita nombre, apellido, correo y
contraseña con confirmación).

Flags:
  --full       Modo completo no interactivo (requiere los 4 flags de datos).
  --nombre     Nombre del sys_admin (requerido con --full).
  --apellido   Apellido del sys_admin (requerido con --full).
  --correo     Correo electrónico del sys_admin (requerido con --full).
  --password   Contraseña del sys_admin (requerido con --full).
  -h, --help   Muestra esta ayuda.

Ejemplos:
  $ go run ./cmd/init_superuser
  $ go run ./cmd/init_superuser --full \
      --nombre "Admin" --apellido "Sistema" \
      --correo "admin@sistema.com" --password "S3gur0#2026"
`

func main() {
	// --- Flags (flag stdlib, GUD-001) ---
	full := flag.Bool("full", false, "Modo completo no interactivo (requiere todos los flags de datos).")
	nombre := flag.String("nombre", "", "Nombre del sys_admin.")
	apellido := flag.String("apellido", "", "Apellido del sys_admin.")
	correo := flag.String("correo", "", "Correo electrónico del sys_admin.")
	password := flag.String("password", "", "Contraseña del sys_admin.")
	help := flag.Bool("help", false, "Muestra la ayuda.")
	flag.BoolVar(help, "h", false, "Muestra la ayuda (alias de --help).")
	flag.Usage = func() { fmt.Print(helpText) }

	flag.Parse()

	if *help {
		fmt.Print(helpText)
		os.Exit(0)
	}

	// Banner inicial (GUD-002).
	imprimirBannerInicial()

	// --- Inicialización de dependencias (PAT-001) ---
	cfg, err := config.LoadConfig()
	if err != nil {
		falloFatal("cargando configuración: %v", err)
	}

	db, err := config.InitDB(cfg.GetDSN())
	if err != nil {
		falloFatal("conectando a la base de datos: %v", err)
	}

	reg := registry.NewRegistry(db, cfg)
	defer reg.Close()

	uc := reg.BootstrapSysAdminCasoDeUso()

	// --- Construcción del comando según el modo ---
	var cmd *usecase.ComandoCrearPrimerSysAdmin
	if *full {
		// Modo no interactivo: valida que todos los flags requeridos estén.
		if err := validarFlagsFull(*nombre, *apellido, *correo, *password); err != nil {
			falloFatal("flags incompletos: %v", err)
		}
		cmd = &usecase.ComandoCrearPrimerSysAdmin{
			Nombre:   *nombre,
			Apellido: *apellido,
			Correo:   *correo,
			Password: *password,
		}
		// En --full no aplicamos la política completa de password; el
		// caso de uso aplica el baseline (no-vacío + len>=8).
	} else {
		// Modo interactivo: prompts + política completa de password.
		c, err := promptInteractivo(os.Stdin)
		if err != nil {
			falloFatal("%v", err)
		}
		cmd = c
	}

	// --- Ejecución del caso de uso ---
	fmt.Printf("\n%s⏳ Creando sys_admin...%s\n", amar, reset)

	resp, err := uc.Ejecutar(context.Background(), cmd)
	if err != nil {
		falloFatal("error al crear sys_admin: %v\n   Todos los cambios fueron revertidos.", err)
	}

	// --- Render del resultado ---
	if resp.YaExistia {
		fmt.Printf("\n%s⚠️  Ya existe un usuario con rol sys_admin. No se realizarán cambios.%s\n", amar, reset)
		fmt.Printf("   ID del sys_admin existente: %s%s%s\n", cian, resp.ExistenteID, reset)
		os.Exit(0)
	}

	imprimirBannerExito(resp)
}

// validarFlagsFull valida que los 4 flags requeridos del modo --full
// estén presentes y no vacíos.
func validarFlagsFull(nombre, apellido, correo, password string) error {
	var faltan []string
	if nombre == "" {
		faltan = append(faltan, "--nombre")
	}
	if apellido == "" {
		faltan = append(faltan, "--apellido")
	}
	if correo == "" {
		faltan = append(faltan, "--correo")
	}
	if password == "" {
		faltan = append(faltan, "--password")
	}
	if len(faltan) > 0 {
		return fmt.Errorf("faltan flags requeridos en modo --full: %s", strings.Join(faltan, ", "))
	}
	return nil
}

// promptInteractivo solicita nombre, apellido y correo, y luego solicita
// la contraseña (con confirmación) aplicando la política completa de
// password (`shared/application.ValidarFormatoPassword`) como pre-check de
// UI. Repite el prompt de contraseña hasta que cumpla (REQ-011/012).
func promptInteractivo(in *os.File) (*usecase.ComandoCrearPrimerSysAdmin, error) {
	r := bufio.NewReader(in)

	fmt.Printf("\n%s📝 Ingrese los datos del sys_admin:%s\n", cian, reset)

	nombre, err := leerCampo(r, "  Nombre")
	if err != nil {
		return nil, err
	}
	apellido, err := leerCampo(r, "  Apellido")
	if err != nil {
		return nil, err
	}

	// Correo: valida formato estándar; repite si es inválido.
	var correo string
	for {
		correo, err = leerCampo(r, "  Correo")
		if err != nil {
			return nil, err
		}
		if _, err := mail.ParseAddress(correo); err == nil {
			break
		}
		fmt.Printf("  %s❌ formato de correo inválido. Intente de nuevo.%s\n", rojo, reset)
	}

	// Password: política completa + confirmación. Repite hasta cumplir ambos.
	var password string
	for {
		p1, err := leerPassword(r, "  Contraseña")
		if err != nil {
			return nil, err
		}
		if err := application.ValidarFormatoPassword(p1, "contraseña"); err != nil {
			fmt.Printf("  %s❌ %s%s\n", rojo, err, reset)
			continue
		}
		p2, err := leerPassword(r, "  Confirmar contraseña")
		if err != nil {
			return nil, err
		}
		if p1 != p2 {
			fmt.Printf("  %s❌ Las contraseñas no coinciden. Intente de nuevo.%s\n", rojo, reset)
			continue
		}
		password = p1
		break
	}

	return &usecase.ComandoCrearPrimerSysAdmin{
		Nombre:   nombre,
		Apellido: apellido,
		Correo:   correo,
		Password: password,
	}, nil
}

// leerCampo Lee una línea no vacía (trimada). Imprime el prompt y retorna
// el valor o error.
func leerCampo(r *bufio.Reader, prompt string) (string, error) {
	for {
		fmt.Printf("%s: ", prompt)
		linea, err := r.ReadString('\n')
		if err != nil && len(linea) == 0 {
			return "", fmt.Errorf("leyendo %s: %w", strings.ToLower(strings.TrimSpace(strings.TrimSuffix(prompt, ":"))), err)
		}
		val := strings.TrimSpace(linea)
		if val != "" {
			return val, nil
		}
		fmt.Printf("  %s❌ El valor no puede estar vacío.%s\n", rojo, reset)
	}
}

// leerPassword Lee una línea no vacía (trimada) para contraseña. No mira
// eco del terminal (no ocultamos; en producción se usaría golang.org/x/term,
// pero conservamos simpleza con el resto del repo).
func leerPassword(r *bufio.Reader, prompt string) (string, error) {
	return leerCampo(r, prompt)
}

// ---- Banners y errores ----

func imprimirBannerInicial() {
	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════╗%s\n", magenta, reset)
	fmt.Printf("%s║  CREACIÓN DEL PRIMER SYS_ADMIN                               ║%s\n", magenta, reset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════╝%s\n", magenta, reset)
}

func imprimirBannerExito(resp *usecase.RespuestaCrearPrimerSysAdmin) {
	verif := "NO"
	if resp.Verificado {
		verif = "SÍ"
	}
	nombreCompleto := fmt.Sprintf("%s %s", resp.Nombre, resp.Apellido)

	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════╗%s\n", verde, reset)
	fmt.Printf("%s║  ✅ SYS_ADMIN CREADO EXITOSAMENTE                            ║%s\n", verde, reset)
	fmt.Printf("%s╠══════════════════════════════════════════════════════════════╣%s\n", verde, reset)
	fmt.Printf("%s║  ID:         %s%-48s%s  ║%s\n", verde, cian, resp.UsuarioID, verde, reset)
	fmt.Printf("%s║  Nombre:     %-48s  ║%s\n", verde, nombreCompleto, reset)
	fmt.Printf("%s║  Correo:     %-48s  ║%s\n", verde, resp.Correo, reset)
	fmt.Printf("%s║  Estado:     %-48s  ║%s\n", verde, resp.Estado, reset)
	fmt.Printf("%s║  Verificado: %-48s  ║%s\n", verde, verif, reset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════╝%s\n", verde, reset)
	fmt.Println()
}

// falloFatal imprime un error con prefijo ❌ (GUD-003) y termina el proceso.
func falloFatal(format string, args ...interface{}) {
	fmt.Printf("\n%s❌ %s%s\n", rojo, fmt.Sprintf(format, args...), reset)
	os.Exit(1)
}

// (fin del archivo)
