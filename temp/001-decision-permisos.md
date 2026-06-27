# Decision: Centralizar permisos de fincas en identidad

## Veredicto: APROBADO con modificaciones

### Punto 1: Catálogo compartido
❌ **Rechazado** `bunna/internal/permisos/` — módulos Go separados no pueden importarlo.
✅ Crear `fincas/internal/application/permisos.go` con constantes.
✅ Añadir entradas a `identidad/internal/rbac/domain/permisos.go` en `TodosLosPermisos`.

### Punto 2: Comunicación fincas→identidad
✅ API REST (correcto). NO crear `POST /sync` — usar seed existente de identidad.
✅ Cache TTL local en fincas para evitar N+1 HTTP.
✅ Endpoint existente `GET /api/v1/mis-permisos` o crear `GET /api/v1/permisos/verificar`.

### Punto 3: Cola RabbitMQ "permisos"
❌ **Rechazado**. Permission check es síncrono. Auditoría futura puede usar el `EventPublisher` existente.

### Punto 4: Acoplamiento circular
✅ Riesgo bajo con REST. Sin imports Go directos. DTOs duplicados por servicio.

### Punto 5: Formato `fincas:{recurso}:{accion}`
✅ Compatible. En fincas se usa short code (`CREAR_FINCA`), en DB de identidad se usa namespaced (`fincas:finca:crear`). API traduce.

### Punto 6: Cambios estructurales
- Extraer `rol` y `tenant_id` del JWT (YA existen como claims).
- Refactorizar 10 use cases para usar constantes centralizadas.
- Añadir client HTTP + caché TTL en middleware.
- Actualizar tests que construyen `AuthContext` manual.

Ver `docs/adr/001-permisos-centralizados.md` para el detalle completo con fragmentos de código.
