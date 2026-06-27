# Arquitectura de Autorización Distribuida

Este documento describe el estado actual y los requerimientos para completar el sistema de permisos entre `identidad` y los servicios de negocio (ej. `fincas`).

## 1. Estado Actual (Catálogo de Capacidades)
- [x] **Sincronización de Catálogo**: `fincas` publica sus permisos estandarizados (`fincas:recurso:accion`) a Kafka al arrancar.
- [x] **Registro Central**: `identidad` consume el catálogo y actualiza su registro global (UPSERT).
- [x] **Infraestructura**: Kafka compartido y red Docker unificada.

## 2. Requerimientos: Identidad (IAM Service)

Para que el sistema sea funcional, `identidad` debe actuar como la fuente de verdad de las asignaciones:

1. **Auto-Asignación Proactiva**:
   - Al registrar permisos de un nuevo módulo (ej. `fincas`), `identidad` debe asociar automáticamente esos permisos al rol **Administrador** de todos los tenants existentes.
   - *Justificación*: El Admin del tenant es el dueño del espacio de nombres y debe poder delegar estas nuevas capacidades sin intervención del sistema central.

2. **Publicación de Cambios en Roles**:
   - Cada vez que se modifique la tabla `rol_permisos` (vía API), se debe publicar un evento a Kafka.
   - **Tópico**: `dev.iam.roles` (o similar).
   - **Evento**: `permisos.rol.actualizado`
   - **Payload**:
     ```json
     {
       "rol_id": "nombre-del-rol",
       "tenant_id": "uuid-del-tenant",
       "permisos": ["fincas:finca:crear", "fincas:muestra:crear"],
       "modulo": "fincas"
     }
     ```

## 3. Requerimientos: Fincas (Service Provider)

Para autorizar peticiones de forma local y rápida:

1. **Tabla de Autorización Local**:
   - Necesita una tabla (ej. `rol_permisos_tenant`) que replique la asociación enviada por `identidad`.
   
2. **Consumidor de Roles**:
   - Debe escuchar el tópico de roles de `identidad` y mantener actualizada su tabla local.

3. **Middleware de Autorización (PDP)**:
   - Al recibir un request, el middleware extrae del **JWT**: `roles` (lista) y `tenant_id`.
   - Cruza estos datos con la tabla local para verificar si alguno de los roles del usuario posee el permiso requerido para el endpoint.
   - **Ventaja**: Autorización en < 1ms sin llamadas inter-servicio.

---
**Nota**: No es necesario replicar la tabla de usuarios ni las asociaciones usuario-rol, ya que esa información viaja de forma segura y verificada dentro del JWT en cada petición.
