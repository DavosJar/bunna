---
title: Estandarización de Use Cases para Telemetría Uniforme
version: 0.1
date_created: 2026-06-26
owner: Team Identidad
tags: architecture, identidad, use-cases, telemetry, refactor
---

# Introducción

Estandarizar el 100% de los casos de uso (~34) para que sean decorables con Wrap[Cmd,Resp], eliminando excepciones y código legacy que no sigue el patrón.

## 1. Propósito & Alcance

**Propósito:** Que todos los casos de uso del sistema tengan una firma Ejecutar(ctx, cmd) uniforme, permitiendo telemetría automática vía Wrap sin excepciones.

**Alcance:**
- Use cases estándar ya decorados en registry.go — verificar que todos están envueltos
- Use cases multi-método (VerificarCorreo, RecuperarContrasena) — dividir en casos individuales
- ServicioTenant — eliminar (código muerto)
- internal/handler/ — eliminar (legado, reemplazado por presentation/handlers/)
- shared/presentation/ — mover a internal/shared/presentation/

**Audiencia:** Desarrolladores mid-level del equipo de identidad.

## 2. Definiciones

| Término | Definición |
|---------|------------|
| Use case estándar | Struct con un único método Ejecutar(ctx, cmd) (Resp, error) |
| Use case multi-método | Struct con 2+ métodos de ejecución (Solicitar, Confirmar, etc.) |
| Wrap | Decorador genérico que agrega telemetría automática a cualquier UseCase[Cmd,Resp] |
| Servicio legacy | Capa de aplicación anterior que coordina lógica sin usar casos de uso individuales |
| Registry | Composition root que instancia y cablea todas las dependencias |

## 3. Requisitos y Restricciones

- **REQ-001**: Todo caso de uso debe tener un único método Ejecutar(ctx, cmd) (Resp, error)
- **REQ-002**: Todo comando debe implementar ToLog() para exponer campos no sensibles
- **REQ-003**: El registry debe envolver con Wrap todos los use cases estándar si cfg.TelemetryEnabled = true
- **REQ-004**: No debe existir código muerto que compile sin usarse (ServicioTenant, handler legacy)
- **REQ-005**: Todo lo compartido debe estar bajo internal/ para mantener encapsulación del módulo
- **CON-001**: No modificar LogPayload, buffer, consumer.go, kafka_producer.go
- **CON-002**: No modificar las interfaces públicas de las facades (AuthFacade, UsuarioFacade, etc.)

## 4. Plan de Acción

### Fase 1 — Limpieza Inmediata (prioridad alta)

**Tarea 1.1:** Eliminar ServicioTenant
- Borrar el directorio internal/tenants/application/services/gestionar_tenant/ completo
- Borrar el handler ConfigurarTenantHandler en presentation/handlers/
- Borrar la ruta /api/v1/tenants/configurar del router
- Verificar que ningún otro archivo lo importa (el grep no debe encontrar referencias)

**Tarea 1.2:** Eliminar internal/handler/ legacy
- Borrar el directorio internal/handler/ completo
- Verificar que ningún otro archivo importa sus símbolos (Handler, NewHandler, RegisterRoutes, Health)

**Tarea 1.3:** Mover shared/presentation/ a internal/
- Mover el directorio shared/presentation/ a internal/shared/presentation/
- Actualizar todos los imports que referencian shared/presentation/ a internal/shared/presentation/
- Verificar que nada fuera de internal/ importa el paquete movido

### Fase 2 — Estandarización de Use Cases Multi-Método (prioridad media)

Hay 2 casos de uso que no siguen el patrón Ejecutar único:

- VerificarCorreoCasoDeUso (Solicitar, Confirmar, Reenviar)
- RecuperarContrasenaCasoDeUso (Solicitar, ValidarToken, Confirmar)

Se procesan uno a uno, en cualquier orden.

**Tarea 2.1:** Dividir VerificarCorreoCasoDeUso en 3 use cases estándar
- Crear SolicitarVerificacionCasoDeUso con su comando, respuesta y único Ejecutar
- Crear ConfirmarVerificacionCasoDeUso con su comando, respuesta y único Ejecutar
- Crear ReenviarVerificacionCasoDeUso con su comando, respuesta y único Ejecutar
- Cada comando debe tener ToLog()
- Actualizar VerificacionFacade para recibir 3 UseCase[Cmd,Resp] en lugar de 1 struct con 3 métodos
- Agregar 3 Wrap en registry.go
- Eliminar el struct VerificarCorreoCasoDeUso original
- Tests: redistribuir los tests del caso original

**Tarea 2.2:** Dividir RecuperarContrasenaCasoDeUso en 3 use cases estándar
- Crear SolicitarRecuperacionCasoDeUso con su comando, respuesta y único Ejecutar
- Crear ValidarTokenRecuperacionCasoDeUso con su comando, respuesta y único Ejecutar
- Crear ConfirmarRecuperacionCasoDeUso con su comando, respuesta y único Ejecutar
- Cada comando debe tener ToLog()
- Actualizar RecuperacionFacade para recibir 3 UseCase[Cmd,Resp] en lugar de 1 struct con 3 métodos
- Agregar 3 Wrap en registry.go
- Eliminar el struct RecuperarContrasenaCasoDeUso original
- Tests: redistribuir los tests del caso original

### Fase 3 — Verificación Final (prioridad alta, después de cada fase)

**Tarea 3.1:** Para cada fase, ejecutar:
- go build ./... (sin errores)
- go test ./... (solo deben fallar los tests de integración PostgreSQL pre-existentes)
- Revisar que el coverage de los nuevos use cases incluya casos de éxito y error

## 5. Criterios de Aceptación

- **AC-001**: El paquete internal/tenants/application/services/gestionar_tenant/ no existe
- **AC-002**: El paquete internal/handler/ no existe
- **AC-003**: Todos los imports de shared/presentation/ apuntan a internal/shared/presentation/
- **AC-004**: VerificarCorreoCasoDeUso fue reemplazado por 3 casos de uso individuales, cada uno con Ejecutar único
- **AC-005**: RecuperarContrasenaCasoDeUso fue reemplazado por 3 casos de uso individuales, cada uno con Ejecutar único
- **AC-006**: Todos los casos de uso en registry.go están envueltos con Wrap si cfg.TelemetryEnabled = true
- **AC-007**: go build ./... compila sin errores
- **AC-008**: go test ./... pasa (excepto integración PostgreSQL)

## 6. Casos de Ejemplo

### Caso estándar: ExpulsarUsuarioCasoDeUso
- Ubicación: internal/usuarios/application/usecases/expeluser/
- Firma: Ejecutar(ctx, *ComandoExpulsarUsuario) (*RespuestaExpulsarUsuario, error)
- Estado actual: YA estandarizado, YA decorado con Wrap("ExpulsarUsuario")
- Acción requerida: Ninguna. Es el modelo a seguir.

### Caso multi-método: ServicioTenant (legacy)
- Ubicación: internal/tenants/application/services/gestionar_tenant/
- Firma: 10 métodos públicos, sin Ejecutar
- Estado actual: Código muerto desde que se eliminó ConfigurarTenant
- Acción requerida: Eliminación completa (Fase 1)
- Ruta de convergencia: No aplica — no se estandariza, se elimina. Sus funcionalidades ya fueron reemplazadas por casos de uso individuales (ListarMisTenantsCasoDeUso, ObtenerTenantPorIDCasoDeUso, etc.)

## 7. Dependencias

- Registry: todos los nuevos use cases deben registrarse en NewRegistry
- Facades: actualizar constructores para recibir UseCase[Cmd,Resp] en lugar de tipos concretos
- Decorator: Wrap[Cmd,Resp] está implementado y funcionando. No requiere cambios.
- ToLog(): ya implementado en todos los comandos existentes. Los comandos nuevos deben implementarlo.

## 8. Notas

- Los use cases con 2 parámetros en Ejecutar (ListarPermisos, ListarMisPermisos, VerificarPermiso) NO se estandarizan — su firma es intencional y no requieren telemetría
- El orden sugerido es Fase 1 → build+test → Fase 2 → build+test → Fase 3
- Cada tarea de la Fase 2 puede hacerse en sesiones separadas, no bloquean producción
