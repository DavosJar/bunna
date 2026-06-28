package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	bootstrap "github.com/davosjar/bunna/services/identidad/internal/bootstrap/domain"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

// errSeedRolesNoEjecutado se retorna cuando el rol sys_admin no existe en la BD.
// Indica que el seed de roles/permisos no fue ejecutado antes del bootstrap.
var errSeedRolesNoEjecutado = errors.New(
	"el rol sys_admin no existe en la BD: el seed de roles/permisos no fue ejecutado",
)

// CrearPrimerSysAdminCasoDeUso crea el primer sys_admin del sistema si no
// existe ninguno. Es un caso de uso de bootstrap: NO requiere permisos
// (no verifica EjecutorID ni perms RBAC), porque resuelve el problema del
// huevo y la gallina del primer usuario con privilegios.
//
// Depende de una sola abstracción: `bootstrap.UnitOfWork`, que encapsula
// los 4 repos + 2 servicios necesarios para crear al sys_admin en una
// transacción atómica de 3 tablas (usuarios, credenciales_usuarios,
// usuario_roles).
type CrearPrimerSysAdminCasoDeUso struct {
	uow bootstrap.UnitOfWork
}

// NewCrearPrimerSysAdminCasoDeUso construye el caso de uso con su única
// dependencia: el UnitOfWork de bootstrap.
func NewCrearPrimerSysAdminCasoDeUso(uow bootstrap.UnitOfWork) *CrearPrimerSysAdminCasoDeUso {
	return &CrearPrimerSysAdminCasoDeUso{uow: uow}
}

// Ejecutar crea el primer sys_admin si no existe ninguno con rol `sys_admin`.
//
// Flujo (ver ADR-001 §9):
//  1. Valida cmd (baseline: nombre/apellido no vacíos <=100; correo con
//     mail.ParseAddress; password no vacío len>=8).
//  2. Pre-check idempotencia FUERA de la tx (read-only): si ya existe un
//     usuario con rol sys_admin → retorna Respuesta{YaExistia:true,
//     ExistenteID:...}, nil.
//  3. Abre tx y:
//     a. Obtiene el rol sys_admin (si no existe → errSeedRolesNoEjecutado).
//     b. Genera usuarioID (UUIDv7).
//     c. Crea usuario con estado NO_VERIFICADO → Activar() + VerificarCorreo()
//     (queda ACTIVO + VERIFICADO).
//     d. Persiste el usuario.
//     e. Hashea el password con EncriptacionServicio (bcrypt).
//     f. Crea credenciales → VerificarCorreo() (correo_verificado=true).
//     g. Persiste credenciales.
//     h. Asigna rol sys_admin (global, sin tenant) en usuario_roles.
//     i. COMMIT (o ROLLBACK automático si fn retorna error).
//  4. Si la tx falla → retorna error envuelto describiendo la causa raíz.
//     Si tuvo éxito → retorna Respuesta{...} con los datos del nuevo sys_admin.
//
// El caso de uso NO hace I/O de consola; es puro y testeable con un fake
// `bootstrap.UnitOfWork` (GUD-004).
func (uc *CrearPrimerSysAdminCasoDeUso) Ejecutar(
	ctx context.Context,
	cmd *ComandoCrearPrimerSysAdmin,
) (*RespuestaCrearPrimerSysAdmin, error) {
	if err := validarComando(cmd); err != nil {
		return nil, fmt.Errorf("validación de comando: %w", err)
	}

	// 2. Pre-check idempotencia (read-only, FUERA de la tx).
	existenteID, existe, err := uc.uow.UsuarioRolRepositorio().ObtenerUsuarioConRol(ctx, rbac.RolSysAdmin)
	if err != nil {
		return nil, fmt.Errorf("-verificando sys_admin existente: %w", err)
	}
	if existe {
		return &RespuestaCrearPrimerSysAdmin{
			YaExistia:   true,
			ExistenteID: existenteID,
		}, nil
	}

	// 3. Transacción atómica de 3 tablas.
	var creado *usuario.Usuario

	err = uc.uow.Transaccional(ctx, func(tx bootstrap.UnitOfWork) error {
		// a. Obtener rol sys_admin (debe existir por el seed).
		rol, err := tx.RolRepositorio().ObtenerPorNombre(ctx, rbac.RolSysAdmin)
		if err != nil {
			if errors.Is(err, rbac.ErrRolNoEncontrado) {
				return errSeedRolesNoEjecutado
			}
			return fmt.Errorf("obteniendo rol sys_admin: %w", err)
		}

		// b. Generar usuarioID (UUIDv7).
		usuarioID, err := tx.GeneradorID().NextID(ctx)
		if err != nil {
			return fmt.Errorf("generando ID de usuario: %w", err)
		}

		// c. Nuevo usuario → ACTIVO + VERIFICADO.
		u, err := usuario.NuevoUsuario(usuarioID, cmd.Correo, cmd.Nombre, cmd.Apellido, "")
		if err != nil {
			return fmt.Errorf("construyendo usuario: %w", err)
		}
		if err := u.Activar(); err != nil {
			return fmt.Errorf("activando usuario: %w", err)
		}
		if err := u.VerificarCorreo(); err != nil {
			return fmt.Errorf("verificando correo del usuario: %w", err)
		}

		// d. Persistir usuario.
		uCreado, err := tx.UsuarioRepositorio().Crear(ctx, u)
		if err != nil {
			return fmt.Errorf("persistiendo usuario: %w", err)
		}
		creado = uCreado

		// e. Hashear password con bcrypt (EncriptacionServicio del registry).
		hash, err := tx.EncriptacionServicio().Hashear(cmd.Password)
		if err != nil {
			return fmt.Errorf("hasheando password: %w", err)
		}

		// f. Crear credenciales → correo verificado.
		c := seguridad.NuevaCredencialesUsuario(usuarioID, hash)
		c.VerificarCorreo()

		// g. Persistir credenciales.
		if _, err := tx.CredencialesRepositorio().Crear(ctx, c); err != nil {
			return fmt.Errorf("persistiendo credenciales: %w", err)
		}

		// h. Asignar rol sys_admin (global, sin tenant).
		if err := tx.UsuarioRolRepositorio().Crear(ctx, usuarioID, rol.ID); err != nil {
			return fmt.Errorf("asignando rol sys_admin: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// 4. Respuesta de éxito.
	return &RespuestaCrearPrimerSysAdmin{
		UsuarioID:  creado.ID(),
		Nombre:     creado.Nombre(),
		Apellido:   creado.Apellido(),
		Correo:     creado.Correo(),
		Estado:     string(creado.Estado()),
		Verificado: string(creado.EstadoVerificacionCorreo()) == string(usuario.VERIFICADO),
		CreadoEn:   creado.FechaCreacion(),
		YaExistia:  false,
	}, nil
}

// validarComando aplica el baseline de validación del caso de uso.
// Es deliberadamente menos estricto que la política completa de password
// (`shared/application.ValidarFormatoPassword`): el caso de uso solo exige
// no-vacío + len>=8 para mantenerse reutilizable (la política completa es
// un concern de UX y la aplica el CLI interactivo).
func validarComando(cmd *ComandoCrearPrimerSysAdmin) error {
	if cmd == nil {
		return errors.New("comando no puede ser nil")
	}
	if cmd.Nombre == "" {
		return errors.New("nombre no puede estar vacío")
	}
	if len(cmd.Nombre) > 100 {
		return errors.New("nombre no puede superar los 100 caracteres")
	}
	if cmd.Apellido == "" {
		return errors.New("apellido no puede estar vacío")
	}
	if len(cmd.Apellido) > 100 {
		return errors.New("apellido no puede superar los 100 caracteres")
	}
	if cmd.Correo == "" {
		return errors.New("correo no puede estar vacío")
	}
	if _, err := mail.ParseAddress(cmd.Correo); err != nil {
		return fmt.Errorf("formato de correo inválido: %w", err)
	}
	if cmd.Password == "" {
		return errors.New("password no puede estar vacío")
	}
	if len(cmd.Password) < 8 {
		return errors.New("password debe tener al menos 8 caracteres")
	}
	return nil
}
