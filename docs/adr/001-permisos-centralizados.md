# ADR-001: Centralización de permisos de fincas en el sistema de identidad

- **Fecha**: 2026-06-26
- **Estado**: Aprobado con modificaciones

## Contexto

Fincas tiene 10 casos de uso que verifican permisos con constantes locales.
Cada use case declara `const permisoRequerido = "CREAR_FINCA"` y llama a
`auth.TienePermiso()`. El `AuthMiddleware` actual extrae solo `UsuarioID`
del JWT y deja `TenantID` y `Permisos` vacíos.

Identidad tiene un sistema RBAC completo con roles de sistema
(sys_admin, administrador, agronomo, caficultor), repositorios de permisos,
y seed de permisos al startup. Actualmente solo maneja permisos del
módulo `identidad:*`. El spec existente en
`identidad/docs/specs/usuarios/spec_permisos_jwt.md` ya describía esta
integración pero nunca se implementó.

## Decisión

**APROBADO** con 5 modificaciones sobre la propuesta original.

Las secciones abajo detallan cada punto de la validación solicitada y los
cambios específicos que el developer debe implementar.

---

### 1. Catálogo compartido — NO en `bunna/internal/permisos/`

La propuesta original ubicaba las constantes en `bunna/internal/permisos/`
como paquete Go compartido. **Se rechaza esta ubicación.**

**Motivo**: `fincas` e `identidad` son módulos Go independientes
(`services/fincas` y `services/identidad`). Un paquete bajo
`bunna/internal/` NO es importable por ninguno de los dos sin añadir
directivas `replace` en `go.mod`, lo que crea acoplamiento en tooling
y dificulta el desarrollo independiente.

**Alternativa aprobada**:

| Módulo | Ubicación del catálogo |
|--------|------------------------|
| fincas | `fincas/internal/application/permisos.go` — constantes agrupadas en un solo archivo |
| identidad | `identidad/internal/rbac/domain/permisos.go` — añadir entradas a `TodosLosPermisos` |

Ambos servicios definen las mismas constantes, cada uno en su módulo.
No hay dependencia de compilación entre ellos. La coherencia se garantiza
por proceso de CI/review (no técnica).

**Cambios para el developer**:

- Crear `fincas/internal/application/permisos.go`:

```go
package application

// Permisos de fincas — catálogo centralizado.
// El formato es UPPER_SNAKE_CASE (short code).
// Identidad almacena el namespaced completo pero la API devuelve
// estos códigos que son los que TienePermiso() compara.
const (
    PermisoFincaCrear        = "CREAR_FINCA"
    PermisoFincaDesactivar   = "DESACTIVAR_FINCA"
    PermisoLoteCrear         = "CREAR_LOTE"
    PermisoLoteEliminar      = "ELIMINAR_LOTE"
    PermisoMuestraCrear      = "CREAR_MUESTRA"
    PermisoMuestraVer        = "VER_MUESTRAS"
    PermisoDiagnosticoSolicitar = "SOLICITAR_DIAGNOSTICO"
    PermisoDiagnosticoAceptar   = "ACEPTAR_DIAGNOSTICO"
    PermisoDiagnosticoRechazar  = "RECHAZAR_DIAGNOSTICO"
    PermisoReporteGenerar    = "GENERAR_REPORTE"
)
```

- Refactorizar los 10 use cases para usar estas constantes. Cada use case
  cambia de `const permisoRequerido = "CREAR_FINCA"` a
  `const permisoRequerido = application.PermisoFincaCrear`.

- En identidad, añadir a `TodosLosPermisos` en `permisos.go`:

```go
// ── Permisos de fincas ──
{PermisoFincaCrear, ...},
{PermisoFincaDesactivar, ...},
// ... los 10 permisos
```

- Asignar los nuevos permisos a los roles en `roles.go`. Los permisos de
  fincas deben asignarse a sys_admin y a agronomo (que opera fincas).

---

### 2. Comunicación fincas → identidad — API REST

**Aprobado**: API REST. Se rechaza DB directa.

La DB directa acopla los esquemas y viola bounded contexts. La API REST
mantiene la independencia de despliegue.

**Precisión**: No crear endpoint `POST /api/v1/permisos/sync`.

En su lugar:

- **Fincas middleware** llama a `GET /api/v1/mis-permisos` de identidad
  (endpoint ya existe). Envía el mismo JWT como Bearer token. Identidad
  responde con la lista de permisos del usuario autenticado.

- Alternativamente, si se necesita pasar `rol` y `tenantID` desde fincas
  (porque fincas extrajo el JWT localmente), crear
  `GET /api/v1/permisos/verificar?rol=X&tenantID=Y` en identidad.

- **Caché local obligatoria** en fincas (TTL 5 minutos, en memoria).
  Sin caché, cada request HTTP a fincas genera un request HTTP a identidad,
  duplicando latencia y creando presión innecesaria.

**Cambios para el developer**:

1. En `fincas/internal/infrastructure/` crear:
   - `permissionclient/` — HTTP client a identidad
   - `permissioncache/` — caché TTL en memoria (map[string][]string con
     expiración)

2. Modificar `AuthMiddleware` en `auth_middleware.go`:

```go
type AuthMiddleware struct {
    validator      *jwtvalidator.TokenValidator
    permissionSvc  PermissionService  // interfaz: ResolverPermisos(usuarioID, tenantID, rol) ([]string, error)
    cache          PermissionCache    // interfaz: Get/Set con TTL
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Validar JWT (existente)
        // 2. Extraer TenantID y Rol del JWT (YA existen en claims de identidad)
        // 3. Resolver permisos vía PermissionService con caché
        // 4. Poblar AuthContext:
        auth := &application.AuthContext{
            UsuarioID: usuarioID,
            TenantID:  tenantID,     // ← ANTES estaba vacío
            Permisos:  permisos,     // ← ANTES estaba nil
        }
        c.Set(ClaveAuthContext, auth)
    }
}
```

3. En `registry/container.go` inyectar `PermissionService` y `PermissionCache`.

---

### 3. Cola "permisos" RabbitMQ — SE ELIMINA

La propuesta original añadía una cola RabbitMQ donde cada caso de uso
publicaba un evento `permiso.ejecutado`. **Se rechaza este punto.**

**Motivos**:

- Permission checking es síncrono — necesito saber AHORA si el usuario
  puede ejecutar la acción. Una cola async no ayuda en la decisión.
- Es complejidad prematura: no hay consumidor definido ni caso de uso
  concreto (auditoría, dashboards, etc.).
- Fincas ya tiene `EventPublisher` para eventos de dominio (FincaCreada,
  MuestraTomada, etc.). Si en el futuro se necesita auditoría de permisos,
  ese mismo mecanismo sirve con un routing key `audit.v1.permiso.ejecutado`.

**Cambio**: No implementar cola "permisos". No tocar `EventPublisher`.

---

### 4. Riesgo de acoplamiento circular

**Riesgo bajo** con la aproximación REST.

- Fincas llama a identidad vía HTTP → sin imports de paquetes Go.
- Identidad NO llama a fincas para permisos.
- Si en el futuro alguien agrega una dependencia Go directa
  (ej. `fincas` importa paquetes de `identidad` o viceversa), el
  compilador de Go detecta ciclos inmediatamente porque ambos son
  módulos separados.

**Único riesgo real**: compartir tipos (DTOs) entre servicios. Mitigación:
cada servicio define sus propios DTOs de API, no se comparten structs.

---

### 5. Compatibilidad del formato `fincas:{recurso}:{accion}`

**Compatible**, pero con precisión sobre el formato real:

| Capa | Formato | Ejemplo |
|------|---------|---------|
| DB identidad | `fincas:{recurso}:{accion}` | `fincas:finca:crear` |
| API identidad → fincas | short code UPPER_SNAKE | `CREAR_FINCA` |
| Constante en use case de fincas | short code UPPER_SNAKE | `CREAR_FINCA` |
| Comparación `TienePermiso()` | short code UPPER_SNAKE | `CREAR_FINCA` |

Esto coincide con el esquema existente de identidad:
- `PermisoDB.Modulo` = `"fincas"` (nuevo valor, antes solo `"identidad"`)
- `PermisoDB.Codigo` = `"fincas:finca:crear"` (nombrado completo)
- La API responde con `Codigo` = `"CREAR_FINCA"` (stripped)

Identidad maneja la traducción internamente. Para los roles de sistema,
se añaden a `RolesDeSistema` usando el short code (igual que los demás
permisos). El seed y `TienePermisoEnRol` comparan contra el `Codigo`
completo almacenado en DB.

---

### 6. Cambios estructurales adicionales recomendados

1. **Extraer `rol` y `tenant_id` del JWT en fincas**: El JWT emitido por
   identidad YA incluye `rol` y `tenant_id` como claims (confirmado en
   `jwt_middleware.go` de identidad). Fincas debe parsearlos:

```go
// auth_middleware.go — extraer claims adicionales
type CustomClaims struct {
    jwt.RegisteredClaims
    UsuarioID string `json:"sub"`
    SesionID  string `json:"sid"`
    TenantID  string `json:"tenant_id"`
    Rol       string `json:"rol"`
}
```

2. **Endpoint `GET /api/v1/permisos/verificar` en identidad** (opcional,
   alternativa a reenviar el JWT): Recibe `rol` y `tenantID` como query
   params. No requiere autenticación adicional (el caller ya validó el JWT).
   Responde con `{ "permisos": ["CREAR_FINCA", ...] }`.

3. **Seed de permisos en identidad**: El `SeedServicio` ya itera
   `TodosLosPermisos`. Al añadir los 10 permisos de fincas, se seedearán
   automáticamente al próximo startup de identidad. No se requiere nueva
   infraestructura.

4. **Tests**: El `AuthContext` actual tiene `TienePermiso()`. Como ahora
   `Permisos` dejará de ser nil, los tests de handler que construyen
   `&application.AuthContext{}` manualmente deben actualizarse para incluir
   los permisos que necesitan.

---

## Resumen de cambios vs propuesta original

| Aspecto | Propuesto | Aprobado |
|---------|-----------|----------|
| Catálogo compartido | `bunna/internal/permisos/` | `fincas/internal/application/permisos.go` + identidad `permisos.go` |
| Sync permisos | `POST /sync` | No crear; usar seed existente de identidad |
| Middleware permisos | Consultar identidad | Consultar identidad + CACHÉ local TTL |
| Cola RabbitMQ | Sí, routing key `permisos.v1.ejecutado` | No implementar |
| Formato permiso | `fincas:{recurso}:{accion}` | Short code en fincas, namespaced en DB identidad |
| TenantID en JWT | No se mencionaba | Ya existe en claims de identidad, extraerlo |

## Veredicto

**APROBADO** con las 5 modificaciones detalladas arriba. El developer puede
implementar directamente siguiendo la tabla de cambios y los fragmentos de
código provistos en cada sección.
