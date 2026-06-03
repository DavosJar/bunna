---
title: CLI Tool — Creación del Primer SysAdmin
version: 0.1
date_created: 2026-06-03
owner: Team Identidad
tags: tool, process, identidad, sysadmin, cli
---

# Introducción

Especificación para un comando CLI que permita crear el primer usuario **sys_admin** del sistema cuando no exista ninguno. Actualmente existe un seed script (`cmd/seed/main.go`) que inserta un usuario "admin@bunna.com" directamente en la tabla `usuarios` sin pasar por el flujo de registro completo ni asignar el rol `sys_admin` a través de las tablas de RBAC (`usuario_roles`). Este spec define un comando robusto que reemplaza ese seed, creando un sys_admin completo, consistente y con rollback automático.

## 1. Purpose & Scope

**Propósito:** Proveer un entrypoint CLI único, idempotente, que cree el primer sys_admin del sistema siguiendo el flujo de registro estándar (hash de contraseña, creación en `usuarios`, `credenciales_usuarios`, y asignación del rol `sys_admin` vía `usuario_roles`).

**Alcance:**
- Comando ejecutable via `go run ./cmd/init_superuser` o binario compilado
- Creación del usuario con estado `ACTIVO` y correo `VERIFICADO`
- Asignación del rol `sys_admin` (global, sin tenant) en `usuario_roles`
- Solo se ejecuta si NO existe ningún usuario con rol `sys_admin` en la BD
- Soporta modo **full** (todos los parámetros en un solo comando) y modo **interactivo** (solicita datos paso a paso)
- Rollback automático en caso de fallo de cualquier paso
- Mensajes de consola claros (éxito con datos del usuario, error con causa específica)

**Audiencia:** Desarrolladores, DevOps, administradores del sistema.

**Supuestos:**
- La base de datos ya existe y las migraciones ya fueron ejecutadas
- El seed de roles y permisos ya fue ejecutado (el rol `sys_admin` existe en `roles`)
- El comando se ejecuta local o en un entorno con acceso a la BD

## 2. Definitions

| Término | Definición |
|---------|------------|
| SysAdmin | Usuario con el rol `sys_admin`, que tiene permisos globales sobre todo el sistema |
| RBAC | Role-Based Access Control — sistema de control de acceso basado en roles |
| Seed | Proceso de inserción inicial de datos necesarios para el funcionamiento del sistema |
| Idempotente | Operación que produce el mismo resultado sin importar cuántas veces se ejecute |
| Rollback | Reversión de todos los cambios realizados si alguna operación falla |
| modo full | Pasar todos los parámetros en una sola línea de comando |
| modo interactivo | El sistema solicita cada campo al usuario paso a paso |
| Tenant | Unidad de aislamiento multi-tenant (el sys_admin es global, no pertenece a ningún tenant) |
| TenantIDSistema | UUID sentinela `00000000-0000-0000-0000-000000000000` para asignaciones de sistema |

## 3. Requirements, Constraints & Guidelines

### Requirements

- **REQ-001**: El comando DEBE verificar si existe al menos un usuario con rol `sys_admin` en `usuario_roles` antes de crear uno nuevo.
- **REQ-002**: Si ya existe un `sys_admin`, el comando DEBE mostrar un mensaje claro y finalizar sin hacer ningún cambio.
- **REQ-003**: El comando DEBE crear un usuario en la tabla `usuarios` con estado `ACTIVO` y `estado_verificacion_correo = 'VERIFICADO'`.
- **REQ-004**: El comando DEBE crear las credenciales en `credenciales_usuarios` con contraseña hasheada (bcrypt), `activo = true`, `correo_verificado = true`.
- **REQ-005**: El comando DEBE asignar el rol `sys_admin` en la tabla `usuario_roles` (rol global, sin tenant).
- **REQ-006**: El comando DEBE usar una transacción de base de datos que abarque todos los pasos (usuario, credenciales, rol).
- **REQ-007**: Si cualquier paso falla dentro de la transacción, DEBE hacerse rollback completo y mostrar un mensaje de error con la causa específica.
- **REQ-008**: El comando DEBE existir como un package `cmd/init_superuser/` con su propio `main.go`.
- **REQ-009**: El comando DEBE soportar el flag `--full` (modo no interactivo) que acepta nombre, apellido, correo y contraseña como argumentos.
- **REQ-010**: Sin el flag `--full`, el comando DEBE entrar en modo interactivo solicitando cada campo al usuario.
- **REQ-011**: En modo interactivo, el comando DEBE validar que la contraseña cumpla con la política de seguridad (mínimo 8 caracteres, al menos una mayúscula, un número y un carácter especial).
- **REQ-012**: En modo interactivo, el comando DEBE solicitar confirmación de la contraseña (escribirla dos veces).
- **REQ-013**: El comando DEBE mostrar un resumen visual claro al finalizar con los datos del sys_admin creado (ID, nombre completo, correo).
- **REQ-014**: El comando DEBE usar los mismos repositorios y servicios del registry para mantener consistencia con el resto del sistema (NO insertar SQL directo ni usar modelos duplicados como hace `cmd/seed/main.go`).

### Constraints

- **CON-001**: No debe usar inserción directa de SQL ni modelos hardcodeados — debe usar la misma capa de aplicación que el registro HTTP (repositorios, casos de uso).
- **CON-002**: No debe modificar ni eliminar datos existentes — es exclusivamente aditivo.
- **CON-003**: El UUID del usuario debe generarse con UUIDv7 (como el resto del sistema) para compatibilidad con índices.
- **CON-004**: El comando no debe depender de variables de entorno o archivos de configuración externa más allá del `config.LoadConfig()` existente.

### Guidelines

- **GUD-001**: Usar `github.com/spf13/cobra` o flags nativos de `flag` para el parsing de argumentos CLI. Preferir `flag` estándar para mantener el ecosistema liviano.
- **GUD-002**: Los mensajes de éxito deben usar formato visual con borde/banner (como `cmd/seed/main.go`) para ser fácilmente identificables.
- **GUD-003**: Los mensajes de error deben incluir el prefijo `❌` y la causa raíz del fallo.
- **GUD-004**: Separar la lógica de negocio (crear sys_admin) del UI de consola (interacción con usuario) usando una función exportable `Ejecutar(ctx, params)` que pueda ser reutilizada en tests.

### Patterns

- **PAT-001**: Seguir el patrón de `cmd/test/main.go` para la inicialización: `config.LoadConfig()` → `config.InitDB()` → `registry.NewRegistry()`.
- **PAT-002**: Seguir el patrón de transacción del `cmd/seed/main.go` para rollback automático.
- **PAT-003**: Usar el `CrearUsuarioCasoDeUso` del registry para crear el usuario en lugar de insert SQL directo.

## 4. Interfaces & Data Contracts

### CLI Interface

```
Usage: go run ./cmd/init_superuser [flags]

Flags:
  --full     Modo completo no interactivo (requiere --nombre, --apellido, --correo, --password)
  --nombre   Nombre del sys_admin (requerido con --full)
  --apellido Apellido del sys_admin (requerido con --full)
  --correo   Correo electrónico del sys_admin (requerido con --full)
  --password Contraseña del sys_admin (requerido con --full)
  -h, --help Muestra esta ayuda
```

### Ejemplos

```bash
# Modo interactivo (solicita cada campo)
go run ./cmd/init_superuser

# Modo full (todos los parámetros)
go run ./cmd/init_superuser --full \
  --nombre "Admin" \
  --apellido "Sistema" \
  --correo "admin@sistema.com" \
  --password "S3gur0#2026"
```

### Internal Service Interface

```go
// SysAdminCreator es el contrato para la lógica de creación del primer sys_admin
type SysAdminCreator interface {
    // CrearPrimerSysAdmin crea el primer sys_admin si no existe.
    // Retorna el usuario creado, o nil si ya existía uno.
    // En caso de error, retorna error describiendo la causa.
    CrearPrimerSysAdmin(ctx context.Context, params ParametrosCreacion) (*SysAdminResultado, error)
}

type ParametrosCreacion struct {
    Nombre   string
    Apellido string
    Correo   string
    Password string
}

type SysAdminResultado struct {
    UsuarioID string
    Nombre    string
    Apellido  string
    Correo    string
    CreadoEn  string
}
```

### Database Objects Involved

| Tabla | Operación | Propósito |
|-------|-----------|-----------|
| `usuarios` | INSERT | Crear el usuario con estado ACTIVO y correo VERIFICADO |
| `credenciales_usuarios` | INSERT | Almacenar password hash, cuenta activa y verificada |
| `usuario_roles` | INSERT | Asignar rol global `sys_admin` al usuario |

### Validations

- **Correo**: Formato email válido (regex estándar). No debe existir previamente en `usuarios.correo`.
- **Password** (interactivo): mínimo 8 caracteres, al menos 1 mayúscula, 1 número, 1 carácter especial.
- **Password** (--full): mínimo 8 caracteres (validación menos estricta para scripts).
- **Nombre/Apellido**: No vacío, máximo 100 caracteres.

## 5. Acceptance Criteria

- **AC-001**: Given una BD sin ningún sys_admin, When se ejecuta el comando con datos válidos, Then se crea el usuario en `usuarios`, las credenciales en `credenciales_usuarios`, y se asigna `sys_admin` en `usuario_roles`. Se muestra mensaje de éxito con los datos del usuario.
- **AC-002**: Given una BD sin ningún sys_admin, When se ejecuta el comando en modo interactivo, Then el sistema solicita nombre, apellido, correo y contraseña (con confirmación), y al finalizar muestra el resumen con ID, nombre, apellido y correo.
- **AC-003**: Given una BD sin ningún sys_admin, When se ejecuta el comando con `--full` y todos los flags, Then crea el sys_admin sin solicitar entrada interactiva.
- **AC-004**: Given una BD que ya tiene un sys_admin, When se ejecuta el comando con datos válidos, Then NO se crea nada nuevo y se muestra un mensaje: "Ya existe un usuario con rol sys_admin. No se realizarán cambios."
- **AC-005**: Given una BD sin ningún sys_admin, When falla la inserción en cualquier tabla dentro de la transacción, Then todas las operaciones se revierten (rollback) y se muestra un mensaje de error con la causa específica.
- **AC-006**: Given una BD sin ningún sys_admin, When se ejecuta con `--password` que no cumple la política (menos de 8 caracteres), Then el comando rechaza la contraseña y muestra un error de validación.
- **AC-007**: Given una BD sin ningún sys_admin, When se ejecuta con `--correo` inválido, Then el comando rechaza el correo y muestra un error de validación.
- **AC-008**: Given el comando se ejecuta dos veces en una BD sin sys_admin, When la primera ejecución crea el sys_admin exitosamente, Then la segunda ejecución detecta que ya existe y no hace nada.
- **AC-009**: Given una BD sin ningún sys_admin, When falla la transacción (ej: violación de constraint), Then el usuario NO queda creado en ninguna tabla y se muestra el error específico de PostgreSQL.

## 6. Test Automation Strategy

- **Test Levels**: Integración (conexión real a BD de prueba)
- **Framework**: No se requiere framework específico — las pruebas de integración en `cmd/test/main.go` son el modelo a seguir.
- **Test Data Management**: El test debe limpiar después de ejecutarse (eliminar el usuario y asignaciones creadas).
- **CI/CD Integration**: Se puede agregar a los pipelines como paso opcional de bootstrap.
- **Coverage Requirements**: La función exportable `CrearPrimerSysAdmin` debe tener cobertura del 100% de casos de error.
- **Scenarios to test**:
  1. Creación exitosa (verificar las 3 tablas: usuarios, credenciales, usuario_roles)
  2. Idempotencia (segunda llamada no hace nada)
  3. Rollback forzado (simular error, verificar que no queden registros huérfanos)
  4. Validación de parámetros inválidos (correo mal formado, password débil)
  5. Ejecución en paralelo (solo uno debe crear el sys_admin)

## 7. Rationale & Context

### Por qué no usar el seed actual

El script `cmd/seed/main.go` actual:
- Inserta SQL directo con modelos hardcodeados duplicados de las tablas
- No usa el caso de uso `CrearUsuarioCasoDeUso` del registry
- No asigna el rol `sys_admin` en `usuario_roles` — el usuario queda sin relación RBAC
- La verificación de existencia solo chequea por correo, no por rol
- No tiene rollback transaccional real (lo tiene, pero sin integración con el registry)

### Por qué este enfoque

- **Consistencia**: Usar los mismos repositorios y casos de uso que el flujo HTTP garantiza que el sys_admin se crea exactamente igual que cualquier otro usuario, pero con estado pre-verificado y rol asignado.
- **Idempotencia**: El chequeo por rol `sys_admin` en `usuario_roles` (no por correo) garantiza que aunque se intente migrar de un sistema a otro, no se duplique.
- **Transaccionalidad**: Una transacción única asegura que no queden registros huérfanos si algo falla.
- **Seguridad**: Usar el password hasher del sistema (bcrypt con costo 12) en lugar de hardcodear un hash.
- **Mantenibilidad**: El comando es autónomo, no depende de scripts externos ni archivos de configuración adicionales.

### Por qué modo interactivo + full

El modo interactivo es más amigable para desarrolladores que ejecutan el comando localmente por primera vez. El modo full es necesario para automatización (scripts de deploy, CI/CD, entornos cloud).

## 8. Dependencies & External Integrations

### Application Dependencies
- **APP-001**: `registry.Registry` — necesario para obtener los casos de uso y repositorios con todas sus dependencias ya inyectadas.
- **APP-002**: `config.LoadConfig()` — para cargar la configuración de BD desde variables de entorno (mismo mecanismo que el servidor HTTP).
- **APP-003**: `config.InitDB(dsn)` — para inicializar la conexión GORM.

### Use Cases Required from Registry
| Caso de Uso | Propósito |
|-------------|-----------|
| `CrearUsuarioCasoDeUso` | Crear el usuario en `usuarios` con hash de password |
| `AsignarRolCasoDeUso` | Asignar el rol `sys_admin` en `usuario_roles` |

### Repositories Required
| Repositorio | Propósito |
|-------------|-----------|
| `RolRepositorio.ObtenerPorNombre(ctx, "sys_admin")` | Obtener el ID del rol sys_admin |
| `UsuarioRolRepositorio.Crear(ctx, usuarioID, rolID)` | Asignar rol global |

### Domain Constants
- `rbac.RolSysAdmin = "sys_admin"` — nombre del rol de sistema.

## 9. Examples & Edge Cases

### Ejemplo 1: Ejecución exitosa (modo interactivo)
```
$ go run ./cmd/init_superuser

╔══════════════════════════════════════════════════════════════╗
║  CREACIÓN DEL PRIMER SYS_ADMIN                              ║
╚══════════════════════════════════════════════════════════════╝

🔍 Verificando si ya existe un sys_admin...
✅ No se encontró ningún sys_admin. Procediendo con la creación.

📝 Ingrese los datos del sys_admin:
  Nombre: Admin
  Apellido: Sistema
  Correo: admin@sistema.com
  Contraseña: [oculto]
  Confirmar contraseña: [oculto]

⏳ Creando sys_admin...

╔══════════════════════════════════════════════════════════════╗
║  ✅ SYS_ADMIN CREADO EXITOSAMENTE                           ║
╠══════════════════════════════════════════════════════════════╣
║  ID:       019e7cad-7087-77fe-8f6a-015e554730bb            ║
║  Nombre:   Admin Sistema                                    ║
║  Correo:   admin@sistema.com                                ║
║  Estado:   ACTIVO                                           ║
║  Verificado: SÍ                                             ║
╚══════════════════════════════════════════════════════════════╝
```

### Ejemplo 2: Ya existe un sys_admin
```
$ go run ./cmd/init_superuser

🔍 Verificando si ya existe un sys_admin...
⚠️  Ya existe un usuario con rol sys_admin. No se realizarán cambios.
   ID del sys_admin existente: 019e7cad-7087-77fe-8f6a-015e554730ba
```

### Ejemplo 3: Error de validación de contraseña
```
$ go run ./cmd/init_superuser

📝 Ingrese los datos del sys_admin:
  Nombre: Admin
  Apellido: Sistema
  Correo: admin@sistema.com
  Contraseña: 123
❌ La contraseña debe tener al menos 8 caracteres, una mayúscula, un número y un carácter especial.
```

### Ejemplo 4: Error de transacción
```
$ go run ./cmd/init_superuser

⏳ Creando sys_admin...
❌ Error al crear sys_admin: ERROR: duplicate key value violates unique constraint "uq_usuarios_correo" (SQLSTATE 23505)
   Todos los cambios fueron revertidos.
```

### Ejemplo 5: Fallo en el paso de asignación de rol (rollback automático)
```
Al insertar en `usuarios` y `credenciales_usuarios` OK, pero falla al asignar el rol
en `usuario_roles`. La transacción revierte TODO (usuario + credenciales + rol).
El usuario no queda registrado en ninguna tabla.
```

## 10. Validation Criteria

Para dar por cumplida esta especificación:

1. **VC-001**: Existe el directorio `cmd/init_superuser/` con un `main.go` compilable.
2. **VC-002**: `go run ./cmd/init_superuser` se ejecuta sin errores y produce la salida esperada.
3. **VC-003**: La función `CrearPrimerSysAdmin` está en un package separado del UI y es testeable unitariamente.
4. **VC-004**: No existe código que inserte SQL directo ni modelos hardcodeados — todo pasa por repositorios/use cases del registry.
5. **VC-005**: Ejecutar el comando dos veces seguidas: la primera crea el sys_admin, la segunda muestra el mensaje de "ya existe" y no hace nada.
6. **VC-006**: El comando funciona tanto en modo interactivo como en modo `--full`.
7. **VC-007**: No se rompen pruebas existentes ni el funcionamiento del servidor HTTP principal (`cmd/main.go`).
8. **VC-008**: El script `cmd/seed/main.go` existente puede ser marcado como deprecado o eliminado después de la implementación.

## 11. Related Specifications / Further Reading

- [Cmd Seed actual](/home/alexis/procesos_software/bunna/identidad/cmd/seed/main.go) — Script que este spec reemplaza
- [Cmd Test integración](/home/alexis/procesos_software/bunna/identidad/cmd/test/main.go) — Patrón de inicialización a seguir
- [Registry del sistema](/home/alexis/procesos_software/bunna/identidad/internal/registry/)
- [RBAC Domain — roles.go](/home/alexis/procesos_software/bunna/identidad/internal/rbac/domain/roles.go) — Constantes de roles del sistema
- [CreateUser use case](/home/alexis/procesos_software/bunna/identidad/internal/usuarios/application/usecases/createuser/) — Caso de uso de creación de usuarios
- [AssignRol use case](/home/alexis/procesos_software/bunna/identidad/internal/rbac/application/usecases/assignrole/) — Caso de uso de asignación de roles
