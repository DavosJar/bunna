# ADR-001: Caso de uso interno `CrearPrimerSysAdminCasoDeUso` (bootstrap)

**Date:** 2026-06-27
**Status:** Proposed
**Authors:** Arquitecto (GLM-5.2) — pendiente validación del equipo
**Spec:** `docs/specs/usuarios/spec-superuser.md`
**Reemplaza:** `cmd/seed/main.go` (a deprecar tras implementación)

---

## Context

El sistema no tiene forma de crear el **primer** `sys_admin` sin pasar por casos de uso
que exigen un `EjecutorID` con permisos (`CrearUsuarioCasoDeUso` requiere
`PermisoUsuarioCrear`; `AsignarRolCasoDeUso` requiere `PermisoRolAsignar`). Es el problema
del huevo y la gallina: nadie tiene permisos para crear al primer usuario con permisos.

El `cmd/seed/main.go` actual lo resuelve insertando SQL directo con modelos duplicados,
sin asignar el rol `sys_admin` en `usuario_roles` y sin rollback real integrado al
registry. El spec exige reemplazarlo por un comando CLI `cmd/init_superuser/` que cree
el sys_admin en una **única transacción** sobre 3 tablas (`usuarios`,
`credenciales_usuarios`, `usuario_roles`), idempotente, con estado `ACTIVO` + correo
`VERIFICADO`, y usando los repositorios del registry (CON-001: sin SQL directo).

El usuario ya decidió la "Opción B": crear un caso de uso de bootstrap interno que **no
requiere permisos**. Este ADR define su arquitectura.

### Hallazgos relevantes del código (verificados con gomcp)

- `internal/registry/registry.go`: NO expone getters públicos para `rolRepository` ni
  `usuarioRolRepository`. Solo `UsuarioTenantRolRepositorio()`. `cmd/test/main.go`
  construye esos repos directamente con `rbac_postgres.NewRolRepositorio(db)` /
  `NewUsuarioRolRepositorio(db)` como workaround.
- `internal/usuarios/infrastructure/persistence/postgres/unit_of_work.go`:
  `UnitOfWorkPostgres.Transaccional` reutiliza las mismas instancias de `usuarioRepo` y
  `credencialesRepo` (que retienen el `db` original) dentro de la tx, **no** reconstruidas
  con `txDB`. Esto es un bug latente: las escrituras no participan realmente de la tx.
  Solo cubre `usuarios` + `credenciales_usuarios`.
- `internal/sesiones/infrastructure/persistence/postgres/unit_of_work.go`:
  `SesionUnitOfWorkPostgres.Transaccional` **sí** reconstruye los repos con `txDB`
  (patrón correcto). Es el modelo a seguir.
- `internal/usuarios/application/usecases/register/usecase.go`: patrón más cercano, pero
  **no transaccional** (llama repos directamente) y deja el usuario en `NO_VERIFICADO`.
- `internal/rbac/domain/repositorios.go`: `UsuarioRolRepositorio` solo tiene
  `Crear/Eliminar/ListarRolesPorUsuario/TieneRol`. No existe "existe usuario con rol X".
- `usuario.NuevoUsuario(...)` → `NO_VERIFICADO`; `.Activar()` → `ACTIVO`;
  `.VerificarCorreo()` → correo verificado. `seguridad.NuevaCredencialesUsuario(...)`
  → `activo=true, correoVerificado=false`; `.VerificarCorreo()` → `true`.

## Alternatives considered

### Opción A: Ubicación bajo `internal/usuarios/application/usecases/bootstrap/`
**Pros:** reutiliza el BC de usuarios; el `register` ya cruza BCs (depende de rbac/seguridad/tenants).
**Cons:** bootstrap es una preocupación **operacional** (inicialización del sistema), no
una operación de dominio del BC usuarios. Meterlo en usuarios contamina el BC con
scripts de arranque y acopla usuarios a rbac+seguridad para algo que no es su
responsabilidad semántica.

### Opción B (elegida): Nuevo package `internal/bootstrap/`
**Pros:** aísla una preocupación cross-cutting (cruza usuarios+seguridad+rbac); hace
explícito el boundary de arranque; no ensucia ningún BC de dominio.
**Cons:** añade un 4º BC operacional. Aceptado: es pequeño y no crecerá.

### Transacción — Opción (a): Extender `usuario.UnitOfWork` existente
**Pros:** reutiliza infraestructura.
**Cons:** hereda el bug latente de no reconstruir repos con `txDB`; requiere añadir
`rolRepo`+`usuarioRolRepo` (el BC usuarios pasaría a depender de rbac); riesgo de
impactar consumidores existentes del UoW. Scope creep en el BC usuarios.

### Transacción — Opción (c): `db.Transaction` de GORM directamente en el caso de uso
**Pros:** mínimo código.
**Cons:** filtra GORM a la capa de aplicación, rompiendo la separación que mantiene el
resto del codebase (los casos de uso dependen de interfaces de dominio, no de GORM).

### Transacción — Opción (b) (elegida): Nuevo `bootstrap.UnitOfWork` específico
**Pros:** replica el patrón **correcto** de `SesionUnitOfWorkPostgres` (reconstruye
repos con `txDB`); aislado al nuevo package; transacción real de 3 tablas con
rollback/commit automáticos; sin tocar el UoW de usuarios (sin heredar su bug); sin
acoplar el BC usuarios a rbac.
**Cons:** un UoW más en el codebase. Aceptado por coherencia con `SesionUnitOfWork`.

### `ExisteUsuarioConRol` — Alternativa: query count directo en el caso de uso
**Pros:** sin tocar la interfaz rbac.
**Cons:** filtra SQL a la capa de aplicación; o requiere un método ad-hoc igualmente.

### `ExisteUsuarioConRol` (elegida): nuevo método en `rbac.UsuarioRolRepositorio`
**Pros:** minimal, intent-revealing, reutilizable, aditivo (no rompe nada).
**Cons:** toda implementación/mock de `UsuarioRolRepositorio` debe añadir el método.

## Decision

**Elegida: Opción B + Opción (b) + nuevo método en `rbac.UsuarioRolRepositorio`.**

Razón: bootstrap es operacional y cruza 3 BCs; aislarlo en `internal/bootstrap/` evita
contaminar el BC de usuarios y deja el boundary explícito. La transacción de 3 tablas
se resuelve con un `bootstrap.UnitOfWork` dedicado que replica el patrón correcto de
`SesionUnitOfWorkPostgres` (reconstruye los 3 repos con `txDB`), sin heredar el bug
latente del UoW de usuarios ni acoplar usuarios a rbac. La idempotencia se resuelve con
un método aditivo y reutilizable en `rbac.UsuarioRolRepositorio`. El caso de uso **no
requiere permisos** (es bootstrap) y es construido en `NewRegistry` como el resto de los
casos de uso, manteniendo al CLI como capa fina.

## Consequences

**Positive:**
- Transacción atómica real sobre 3 tablas con rollback automático (REQ-006/007, AC-005/009).
- Idempotencia por rol (no por correo) (REQ-001/002, AC-004/008).
- Sin SQL directo ni modelos duplicados (CON-001, VC-004).
- Caso de uso testeable sin I/O de consola (GUD-004, VC-003).
- CLI fino, consistente con PAT-001.
- No se modifica el UoW de usuarios (sin riesgo de regresión en `CrearUsuarioCasoDeUso`).

**Negative / accepted debt:**
- Nuevo BC operacional `internal/bootstrap/` (4º package top-level). Aceptado.
- Se duplica lógica mínima de creación (`NuevoUsuario`+`Activar`+`VerificarCorreo`+hash+
  credenciales) frente a `CrearUsuarioCasoDeUso`. Necesario: el caso de uso existente
  exige permisos y no soporta estado `ACTIVO`+`VERIFICADO` en una sola tx con rol.
- Race TOCTOU en el chequeo de idempotencia (dos CLI concurrentes podrían ambos pasar
  el check y crear dos sys_admins). Ver "Riesgos".
- Se añade un método a una interfaz compartida (`rbac.UsuarioRolRepositorio`); los
  mocks/fakes de tests existentes deben implementarlo.

**Required actions:**
- [ ] Crear `internal/bootstrap/domain/` con la interfaz `UnitOfWork`.
- [ ] Crear `internal/bootstrap/application/usecase/` con `command.go`, `response.go`,
      `usecase.go` (contratos abajo).
- [ ] Crear `internal/bootstrap/infrastructure/persistence/postgres/` con
      `BootstrapUnitOfWorkPostgres` (patrón `SesionUnitOfWorkPostgres`).
- [ ] Añadir `ObtenerUsuarioConRol(ctx, rolNombre) (usuarioID string, encontrado bool, err error)`
      a `rbac.UsuarioRolRepositorio` (domain) e implementarlo en `rbac_postgres`
      (`SELECT u.id FROM usuario_roles ur JOIN roles r ON r.id=ur.rol_id JOIN usuarios u
      ON u.id=ur.usuario_id WHERE r.nombre=$1 LIMIT 1`).
- [ ] Actualizar mocks/fakes de `UsuarioRolRepositorio` en tests existentes.
- [ ] En `NewRegistry`: construir `bootstrapUoW` y `bootstrapUC`; añadir getter
      `BootstrapSysAdminCasoDeUso()`.
- [ ] Crear `cmd/init_superuser/main.go` (CLI fino).
- [ ] Marcar `cmd/seed/main.go` como deprecado (VC-008).

---

# Diseño detallado

## 1. Ubicación del package

```
internal/bootstrap/
├── domain/
│   └── unit_of_work.go          // interfaz bootstrap.UnitOfWork (cross-BC)
├── application/
│   └── usecase/
│       ├── command.go           // ComandoCrearPrimerSysAdmin + ToLog
│       ├── response.go          // RespuestaCrearPrimerSysAdmin
│       └── usecase.go           // CrearPrimerSysAdminCasoDeUso + Ejecutar
└── infrastructure/
    └── persistence/
        └── postgres/
            └── unit_of_work.go  // BootstrapUnitOfWorkPostgres
```

**Justificación:** bootstrap es una preocupación operacional cross-cutting (usuarios +
seguridad + rbac). Aislarla evita (a) contaminar el BC `usuarios` con scripts de
arranque y (b) acoplar `usuarios` a `rbac` dentro de su propio UoW. El precedente de
`register` (que cruza BCs) no aplica: `register` es una operación de dominio legítima
del BC usuarios (auto-registro); bootstrap no lo es.

## 2. Contratos (pseudocódigo Go)

```go
// internal/bootstrap/application/usecase/command.go
package usecase

type ComandoCrearPrimerSysAdmin struct {
    Nombre   string
    Apellido string
    Correo   string
    Password string
}

// ToLog excluye siempre el password.
func (c ComandoCrearPrimerSysAdmin) ToLog() map[string]any {
    return map[string]any{
        "correo":   c.Correo,
        "nombre":   c.Nombre,
        "apellido": c.Apellido,
    }
}
```

```go
// internal/bootstrap/application/usecase/response.go
package usecase

import "time"

type RespuestaCrearPrimerSysAdmin struct {
    UsuarioID   string
    Nombre      string
    Apellido    string
    Correo      string
    Estado      string    // "ACTIVO"
    Verificado  bool      // true
    CreadoEn    time.Time
    YaExistia   bool      // true → no se creó nada
    ExistenteID string    // ID del sys_admin preexistente si YaExistia
}
```

```go
// internal/bootstrap/application/usecase/usecase.go
package usecase

type CrearPrimerSysAdminCasoDeUso struct {
    uow bootstrap.UnitOfWork // única dependencia
}

func NewCrearPrimerSysAdminCasoDeUso(uow bootstrap.UnitOfWork) *CrearPrimerSysAdminCasoDeUso {
    return &CrearPrimerSysAdminCasoDeUso{uow: uow}
}

// Ejecutar crea el primer sys_admin si no existe ninguno.
// - Si ya existe: retorna Respuesta{YaExistia:true, ExistenteID:...}, nil.
// - Si no existe: lo crea en una tx de 3 tablas y retorna Respuesta{...}.
// - Error: rollback automático; el error describe la causa raíz.
func (uc *CrearPrimerSysAdminCasoDeUso) Ejecutar(
    ctx context.Context,
    cmd *ComandoCrearPrimerSysAdmin,
) (*RespuestaCrearPrimerSysAdmin, error)
```

```go
// internal/bootstrap/domain/unit_of_work.go
package bootstrap

type UnitOfWork interface {
    Transaccional(ctx context.Context, fn func(tx UnitOfWork) error) error
    UsuarioRepositorio()    usuario.UsuarioRepositorio
    CredencialesRepositorio() seguridad.CredencialesRepositorio
    UsuarioRolRepositorio() rbac.UsuarioRolRepositorio
    RolRepositorio()        rbac.RolRepositorio
    EncriptacionServicio()  seguridad.EncriptacionServicio
    GeneradorID()           shared_domain.GeneradorID
}
```

## 3. Dependencias a inyectar

El caso de uso recibe **una sola** dependencia: `bootstrap.UnitOfWork`, que encapsula:
- `usuario.UsuarioRepositorio` (Crear usuario)
- `seguridad.CredencialesRepositorio` (Crear credenciales)
- `rbac.UsuarioRolRepositorio` (Crear asignación + `ObtenerUsuarioConRol` idempotencia)
- `rbac.RolRepositorio` (`ObtenerPorNombre("sys_admin")` dentro de la tx)
- `seguridad.EncriptacionServicio` (bcrypt hash)
- `shared_domain.GeneradorID` (UUIDv7)

El caso de uso **no** depende de `checkpermission.AuthService` (no hay validación de
permisos: es bootstrap).

## 4. Decisión de transacción

**Opción (b): nuevo `bootstrap.UnitOfWork` dedicado.**

`BootstrapUnitOfWorkPostgres.Transaccional` replica el patrón **correcto** de
`SesionUnitOfWorkPostgres`: abre `db.WithContext(ctx).Transaction(func(txDB *gorm.DB)...)`
y reconstruye los 3 repos con `txDB` (`NewUsuarioRepositorio(txDB)`,
`NewCredencialesRepositorio(txDB)`, `NewUsuarioRolRepositorio(txDB)`,
`NewRolRepositorio(txDB)`). `EncriptacionServicio` y `GeneradorID` se reutilizan (no
tocan la BD). Si `fn` retorna error → ROLLBACK automático de GORM; nil → COMMIT.

Para el path de **lectura** (pre-check de idempotencia fuera de la tx), el UoW expone
los repos plain-db (los pasados al constructor), de modo que
`uow.UsuarioRolRepositorio().ObtenerUsuarioConRol(...)` se ejecuta sin abrir tx.

**Justificación:** (a) heredaría el bug latente del UoW de usuarios y acoplaría
`usuarios` a `rbac`; (c) filtraría GORM a la capa de aplicación. (b) es aislado, correcto
y consistente con `SesionUnitOfWork`.

## 5. Decisión sobre `ExisteUsuarioConRol`

**Agregar método aditivo a `rbac.UsuarioRolRepositorio`:**

```go
// Obtiene el ID del primer usuario con el rol dado (por nombre de rol).
// encontrado=false si no hay ninguno. No retorna error en ese caso.
ObtenerUsuarioConRol(ctx context.Context, rolNombre string) (usuarioID string, encontrado bool, err error)
```

Implementación en `rbac_postgres`: `SELECT u.id FROM usuario_roles ur JOIN roles r ON
r.id = ur.rol_id JOIN usuarios u ON u.id = ur.usuario_id WHERE r.nombre = $1 LIMIT 1`.
`encontrado=false` cuando `gorm.ErrRecordNotFound`.

**Justificación:** es la query mínima e intent-revealing; satisface REQ-001/002 y el
mensaje de "ya existe" con ID (ejemplo 2 del spec). Es aditiva (no rompe
implementaciones existentes salvo mocks, que deben actualizarse). Evita `ListarRolesPor
Usuario` (requiere un usuarioID que no tenemos) y evita SQL ad-hoc en el caso de uso.

## 6. Cómo el CLI obtiene las dependencias

**Construir el caso de uso (y su UoW) dentro de `NewRegistry` y exponerlo por getter**,
igual que todos los demás casos de uso. El CLI queda fino (PAT-001):

```go
// cmd/init_superuser/main.go (esqueleto)
cfg := config.LoadConfig()
db := config.InitDB(config.GetDSN(cfg))
reg := registry.NewRegistry(db, cfg)
defer reg.Close()
uc := reg.BootstrapSysAdminCasoDeUso()
resp, err := uc.Ejecutar(ctx, &usecase.ComandoCrearPrimerSysAdmin{...})
// render banner / error
```

**Cambios en `internal/registry/registry.go`:**
- Nuevo campo privado `bootstrapUoW bootstrap.UnitOfWork` (opcional, si se quiere
  exponer el UoW; si no, solo el UC).
- Construcción en `NewRegistry`:
  `bootstrapUoW := bootstrap_postgres.NewBootstrapUnitOfWork(db, usuarioRepo,
  credencialesRepo, rolRepo, usuarioRolRepo, encriptacion, generadorID)`
  `bootstrapUC := bootstrap_uc.NewCrearPrimerSysAdminCasoDeUso(bootstrapUoW)`
- Nuevo getter:
  `func (r *Registry) BootstrapSysAdminCasoDeUso() *bootstrap_uc.CrearPrimerSysAdminCasoDeUso`

**No** se añaden getters públicos para `rolRepository`/`usuarioRolRepository` al
Registry: no hacen falta porque el UC se construye dentro de `NewRegistry`. Esto
minimiza la superficie pública y evita replicar el workaround de `cmd/test/main.go`.

## 7. Dónde vive la validación de password

- **Caso de uso:** valida el **baseline** (no vacío + `len >= 8`). Es el contrato para
  cualquier caller y satisface AC-006 (rechaza `< 8` en `--full`).
- **CLI interactivo:** aplica la **política completa**
  (`shared/application.ValidarFormatoPassword`: mayúscula, minúscula, número,
  no-alfanumérico) como pre-check de UI **antes** de invocar al caso de uso, y repite
  el prompt hasta que cumpla (REQ-011, REQ-012 confirmación).
- **CLI `--full`:** no aplica la política completa; confía en el baseline del caso de
  uso (spec §4: `--full` solo mínimo 8).

**Justificación:** mantiene el caso de uso agnóstico a la política (reutilizable) y
ubica la política estricta en el único sitio donde es un concern de UX (el prompt
interactivo).

## 8. Separación UI / lógica (GUD-004)

- `internal/bootstrap/application/usecase/` → **lógica pura**, sin I/O de consola.
  Retorna `*RespuestaCrearPrimerSysAdmin`. Testeable con un fake `bootstrap.UnitOfWork`.
- `cmd/init_superuser/main.go` → **UI**: parsing de flags (`flag` stdlib, GUD-001),
  prompts interactivos (`bufio`), banners (GUD-002), prefijo `❌` (GUD-003), invoca
  `Ejecutar` y formatea la `Respuesta` como banner. **Nunca** el caso de uso escribe a
  stdout/stderr.

## 9. Diagrama de flujo de `Ejecutar`

```
Ejecutar(ctx, cmd)
  │
  ├─ 1. Validar cmd (baseline)
  │     - nombre/apellido no vacíos, <= 100 chars
  │     - correo: mail.ParseAddress (formato)
  │     - password: no vacío, len >= 8
  │     → error de validación si falla
  │
  ├─ 2. Pre-check idempotencia (FUERA de tx, read-only)
  │     uow.UsuarioRolRepositorio().ObtenerUsuarioConRol(ctx, rbac.RolSysAdmin)
  │     → si encontrado: return Respuesta{YaExistia:true, ExistenteID:...}, nil
  │
  ├─ 3. uow.Transaccional(ctx, func(tx) error {
  │     a. tx.RolRepositorio().ObtenerPorNombre(ctx, rbac.RolSysAdmin) → rolID
  │          (si no existe → error "seed de roles no ejecutado")
  │     b. tx.GeneradorID().NextID(ctx) → usuarioID (UUIDv7)
  │     c. usuario.NuevoUsuario(usuarioID, correo, nombre, apellido, "")
  │        → u.Activar(); u.VerificarCorreo()   // ACTIVO + VERIFICADO
  │     d. tx.UsuarioRepositorio().Crear(ctx, u) → persisted
  │     e. tx.EncriptacionServicio().Hashear(password) → hash
  │     f. seguridad.NuevaCredencialesUsuario(usuarioID, hash)
  │        → c.VerificarCorreo()                 // correo_verificado=true
  │     g. tx.CredencialesRepositorio().Crear(ctx, c)
  │     h. tx.UsuarioRolRepositorio().Crear(ctx, usuarioID, rolID)  // global, sin tenant
  │     i. return nil  → COMMIT
  │  })
  │
  └─ 4. On tx error → return wrapped error (rollback automático de GORM)
     On success → return Respuesta{UsuarioID, Nombre, Apellido, Correo,
                                   Estado:"ACTIVO", Verificado:true, CreadoEn, YaExistia:false}
```

## 10. Riesgos y trade-offs

1. **Race TOCTOU en idempotencia.** Dos ejecuciones CLI concurrentes podrían ambas pasar
   el pre-check (paso 2) y crear dos sys_admins. El spec (AC paralelo "solo uno debe
   crear") sugiere robustez. **Recomendación:** aceptar para CLI operator-run
   (secuencial); como endurecimiento opcional, tomar `pg_try_advisory_xact_lock` al
   inicio de la tx (paso 3) — requiere una llamada SQL puntual en el UoW, no en el caso
   de uso. **Decisión diferida al equipo.**

2. **Bug latente en `usuario.UnitOfWork`.** Fuera de scope. El nuevo `bootstrap.UoW`
   sigue el patrón correcto (`SesionUnitOfWorkPostgres`) y no lo hereda. Dejar como
   deuda técnica documentada para un futuro ADR de cleanup.

3. **Duplicación mínima de lógica de creación** frente a `CrearUsuarioCasoDeUso`.
   Aceptado: necesaria para (a) no requerir permisos y (b) fijar `ACTIVO`+`VERIFICADO` y
   asignar rol en una sola tx.

4. **Método nuevo en interfaz compartida.** Cualquier mock/fake de
   `rbac.UsuarioRolRepositorio` debe implementar `ObtenerUsuarioConRol`. Acción
   requerida listada arriba.

5. **Unicidad de correo.** No hay pre-check de correo en el caso de uso; se delega al
   constraint `uq_usuarios_correo` de la BD, que dispara rollback (AC-009). Aceptado:
   mantiene el caso de uso simple y consistente con `register`.

6. **Mostrar ID del sys_admin existente.** El spec (ejemplo 2) lo muestra. Se satisface
   vía `ObtenerUsuarioConRol` (retorna el ID). Si el equipo prefiere no exponerlo, caer
   a un método `ExisteUsuarioConRol(ctx, rolNombre) (bool, error)` bool-only. **Menor.**

## Unresolved

- **Race handling definitivo** (aceptar vs advisory lock vs unique partial index).
  Recomendación del arquitecto: **aceptar + documentar**; endurecer con advisory lock
  solo si el comando se usa en automatización paralela.
- **Si exponer el ID del sys_admin existente** en el mensaje de "ya existe".
  Recomendación: **sí** (ya lo devuelve `ObtenerUsuarioConRol`).
