---
title: "Spec 0 — Refactor: CorreoElectronico como Value Object"
version: 1.0
date_created: 2026-05-22
owner: Equipo Identidad
tags: refactor, dominio, value-object, correo, verificacion
---

# Spec 0 — Refactor: CorreoElectronico como Value Object

> **Propósito**: Refactorizar el modelo de dominio del correo electrónico en la entidad Usuario, extrayendo un Value Object `CorreoElectronico` que encapsule la dirección y su estado de verificación, y limpiando la máquina de estados eliminando `VERIFICACION_FALLIDA`.
>
> **No incluye**: El Value Object `PruebaVerificacion` (se maneja en la spec de registro), cambios en la lógica de negocio de registro o verificación, nuevos casos de uso.
>
> **Ejecutar antes que**: Cualquier spec de autorización. No es una feature nueva, es una limpieza/mejora del dominio existente.

---

## 1. Estado Actual del Código

Actualmente en el código:

- `Usuario` tiene `correo string` como campo plano y `estadoVerificacionCorreo EstadoVerificacionCorreo` como otro campo separado
- `EstadoVerificacionCorreo` tiene 5 estados: `PENDIENTE_VERIFICACION`, `VERIFICADO`, `ENLACE_EXPIRADO`, `REENVIO_SOLICITADO`, `VERIFICACION_FALLIDA`
- `VERIFICACION_FALLIDA` es un estado que se determinó que no debe existir (si un token de verificación es inválido, simplemente se rechaza sin cambiar el estado del dominio)
- No existe un Value Object `CorreoElectronico` que agrupe estos conceptos

---

## 2. Motivación

| Razón | Detalle |
|-------|---------|
| **Encapsulamiento** | El correo y su estado de verificación son conceptos cohesivos que deben estar juntos. Tenerlos como campos separados en Usuario permite estados inconsistentes (ej: correo vacío con estado VERIFICADO). |
| **Inmutabilidad** | Un Value Object `CorreoElectronico` garantiza que siempre se cree con un estado válido y que las transiciones de estado sean controladas. |
| **Limpieza** | `VERIFICACION_FALLIDA` es un estado espurio. La decisión de diseño fue: un intento de verificación con token inválido NO altera el estado del dominio. Este estado ya no se usa pero sigue en el código. |
| **Consistencia DDD** | Value Objects para conceptos de dominio enriquecidos es una práctica recomendada de Domain-Driven Design. |

---

## 3. Cambios en el Dominio

### 3.1 Value Object: `CorreoElectronico`

Se crea un nuevo Value Object que encapsula:

| Atributo | Tipo | Descripción |
|----------|------|-------------|
| `direccion` | string | La dirección de correo electrónico (ej: usuario@dominio.com) |
| `estado` | EstadoVerificacionCorreo | Estado actual de verificación |

**Comportamiento del VO:**

- **Constructor**: Recibe una dirección de correo válida. Crea el VO con estado `PENDIENTE_VERIFICACION`.
- **Validación**: La dirección debe ser un email válido (formato RFC 5322) y no vacío.
- **Transiciones de estado**: Solo permite las transiciones válidas definidas en la máquina de estados (sin `VERIFICACION_FALLIDA`).
- **Métodos de comportamiento**: `Verificar()`, `MarcarExpirado()`, `SolicitarReenvio()`, `EstaVerificado()`, `EstaPendiente()`.
- **Reconstrucción desde BD**: Constructor alternativo que recibe todos los valores sin validar (para hidratación desde persistencia).

### 3.2 Máquina de Estados (limpiada)

La máquina de estados limpiada tiene 4 estados (se elimina `VERIFICACION_FALLIDA`):

```
PENDIENTE_VERIFICACION ───→ VERIFICADO
PENDIENTE_VERIFICACION ───→ ENLACE_EXPIRADO
PENDIENTE_VERIFICACION ───→ REENVIO_SOLICITADO
ENLACE_EXPIRADO        ───→ REENVIO_SOLICITADO
REENVIO_SOLICITADO     ───→ VERIFICADO
REENVIO_SOLICITADO     ───→ ENLACE_EXPIRADO
VERIFICADO             ───→ (terminal, sin transiciones)
```

**Cambios respecto al estado actual:**

| Cambio | Antes | Después |
|--------|-------|---------|
| Estados | 5: PENDIENTE, VERIFICADO, ENLACE_EXPIRADO, REENVIO_SOLICITADO, VERIFICACION_FALLIDA | 4: Eliminar VERIFICACION_FALLIDA |
| Transiciones desde PENDIENTE | Permitía → VERIFICACION_FALLIDA | Ya no permite |
| Transiciones desde REENVIO | Permitía → VERIFICACION_FALLIDA | Ya no permite |

### 3.3 Entidad `Usuario` modificada

| Campo | Estado actual | Después del refactor |
|-------|---------------|---------------------|
| `correo` | `string` (campo plano) | ❌ Eliminado |
| `estadoVerificacionCorreo` | `EstadoVerificacionCorreo` (campo plano) | ❌ Eliminado |
| `correoElectronico` | — | ✅ NUEVO: `CorreoElectronico` (VO) |

**Métodos que cambian en Usuario:**

| Método | Estado actual | Después |
|--------|---------------|---------|
| `Correo() string` | Retorna `correo` plano | Retorna `correoElectronico.Direccion()` |
| `EstadoVerificacionCorreo()` | Retorna `estadoVerificacionCorreo` | Retorna `correoElectronico.Estado()` |
| `VerificarCorreo()` | Lógica inline | Delega en `correoElectronico.Verificar()` |
| `SolicitarReenvioVerificacion()` | Lógica inline | Delega en `correoElectronico.SolicitarReenvio()` |
| `MarcarEnlaceExpirado()` | Lógica inline | Delega en `correoElectronico.MarcarExpirado()` |
| `MarcarVerificacionFallida()` | Existe | ❌ ELIMINADO |
| `NuevoUsuario(...)` | Recibe `correo string` | Recibe `correo string`, crea `CorreoElectronico` internamente |
| `NewUsuarioFromPersistence(...)` | Recibe `correo string` + `estadoVerificacion` | Recibe un `CorreoElectronico` ya construido |

**Constructor `NuevoUsuario`:**

- Sigue recibiendo `correo string` como parámetro
- Internamente crea `NuevoCorreoElectronico(correo)`
- Si el correo es inválido, propaga el error

**Constructor `NewUsuarioFromPersistence`:**

- Ahora recibe `correoElectronico *CorreoElectronico` en lugar de `correo string` + `estadoVerificacion`
- Usa `NuevoCorreoElectronicoDesdeBD(direccion, estado)` para reconstruir

---

## 4. Archivos afectados

```
CAPA DE DOMINIO:
  internal/usuarios/domain/usuario/
  ├── correo_electronico.go              ✅ NUEVO — Value Object CorreoElectronico
  ├── correo_electronico_test.go          ✅ NUEVO — Tests del VO
  ├── estado_verificacion_correo.go       🔄 MODIFICAR — Eliminar VERIFICACION_FALLIDA
  ├── estado_verificacion_correo_test.go  🔄 MODIFICAR — Tests actualizados
  ├── usuario.go                          🔄 MODIFICAR — Reemplazar campos por VO
  ├── usuario_test.go                     🔄 MODIFICAR — Tests actualizados
  ├── eventos.go                          🔄 MODIFICAR — Evento VerificacionFallida eliminado
  ├── errors.go                           🔄 SIN CAMBIOS (los errores siguen siendo válidos)
  ├── repositorio.go                      🔄 SIN CAMBIOS (interfaz no cambia)
  └── especificacion_usuario.go           🔄 SIN CAMBIOS (filtros siguen siendo sobre "correo")

CAPA DE APLICACIÓN:
  internal/usuarios/application/services/registro/
  ├── servicio_registro.go    🔄 MODIFICAR — usuarioCreado.Correo() debe seguir funcionando
  ├── comando.go              🔄 SIN CAMBIOS
  ├── respuesta.go            🔄 SIN CAMBIOS
  ├── ejecutor.go             🔄 SIN CAMBIOS
  └── servicio_registro_test.go  🔄 POSIBLEMENTE — Si usa Correo() o EstadoVerificacionCorreo()

CAPA DE INFRAESTRUCTURA:
  internal/usuarios/infrastructure/persistence/postgres/
  ├── usuario_model.go         🔄 MODIFICAR — ToDomain() y FromDomain() deben usar el nuevo VO
  └── usuario_repositorio.go   🔄 MODIFICAR — Si referencia campos directos de usuario

CAPA DE PRESENTACIÓN:
  internal/presentation/
  ├── dto/register_dto.go       🔄 SIN CAMBIOS (solo serialización)
  ├── handlers/register_handler.go  🔄 SIN CAMBIOS
  └── facades/auth_facade_impl.go   🔄 SIN CAMBIOS

OTROS:
  internal/registry/registry.go     🔄 SIN CAMBIOS
```

---

## 5. Cambios en BD

No hay cambios en el esquema de BD. Las columnas existentes se mantienen:

| Columna actual | Estado |
|----------------|--------|
| `correo` (VARCHAR) | Se mantiene |
| `estado_verificacion_correo` (VARCHAR) | Se mantiene |

**Migración de datos**:

Si existen registros con `estado_verificacion_correo = 'VERIFICACION_FALLIDA'`, deben migrarse a `PENDIENTE_VERIFICACION` (es el estado más seguro: el usuario puede solicitar reenvío). Se ejecuta un UPDATE en la migración.

```sql
UPDATE usuarios
SET estado_verificacion_correo = 'PENDIENTE_VERIFICACION'
WHERE estado_verificacion_correo = 'VERIFICACION_FALLIDA';
```

---

## 6. Escenarios TDD

### 6.1 CorreoElectronico VO

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 1 | Crear VO con dirección válida | email válido | `NuevoCorreoElectronico(email)` | VO creado, estado PENDIENTE_VERIFICACION |
| 2 | Crear VO con dirección vacía | email = "" | `NuevoCorreoElectronico(email)` | Error de dominio |
| 3 | Crear VO con formato inválido | email = "invalido" | `NuevoCorreoElectronico(email)` | Error de dominio |
| 4 | Reconstruir desde BD | dirección + estado | `NuevoCorreoElectronicoDesdeBD(dir, estado)` | VO sin validar ni emitir eventos |
| 5 | Verificar correo | estado PENDIENTE | `VO.Verificar()` | estado → VERIFICADO |
| 6 | Verificar correo ya VERIFICADO | estado VERIFICADO | `VO.Verificar()` | Error de transición |
| 7 | Verificar correo en ENLACE_EXPIRADO | estado ENLACE_EXPIRADO | `VO.Verificar()` | Error (transición no permitida) |
| 8 | Marcar enlace expirado | estado PENDIENTE | `VO.MarcarExpirado()` | estado → ENLACE_EXPIRADO |
| 9 | Solicitar reenvío desde PENDIENTE | estado PENDIENTE | `VO.SolicitarReenvio()` | estado → REENVIO_SOLICITADO |
| 10 | Solicitar reenvío desde ENLACE_EXPIRADO | estado ENLACE_EXPIRADO | `VO.SolicitarReenvio()` | estado → REENVIO_SOLICITADO |
| 11 | Estado terminal VERIFICADO | estado VERIFICADO | Cualquier transición | Error |
| 12 | Getter direccion() | VO creado | `VO.Direccion()` | La misma dirección del constructor |
| 13 | Getter estado() | VO creado | `VO.Estado()` | El estado actual |
| 14 | EstaVerificado() true | estado VERIFICADO | `VO.EstaVerificado()` | true |
| 15 | EstaVerificado() false | estado PENDIENTE | `VO.EstaVerificado()` | false |
| 16 | EstaPendiente() true | estado PENDIENTE | `VO.EstaPendiente()` | true |
| 17 | EstaPendiente() false | estado VERIFICADO | `VO.EstaPendiente()` | false |

### 6.2 Máquina de estados limpiada

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 18 | VERIFICACION_FALLIDA no existe | Constantes de estado | Buscar VERIFICACION_FALLIDA | No debe existir |
| 19 | PENDIENTE no puede ir a VERIFICACION_FALLIDA | Estado PENDIENTE | `PuedeTransicionarA(VERIFICACION_FALLIDA)` | false (la constante ya no existe) |
| 20 | Mapa de transiciones tiene 4 estados | — | Verificar claves del mapa | 4 entradas (sin VERIFICACION_FALLIDA) |
| 21 | VERIFICADO sigue siendo terminal | — | Verificar transiciones de VERIFICADO | Lista vacía |

### 6.3 Usuario con VO embebido

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 22 | Crear usuario con correo válido | email válido | `NuevoUsuario(id, email, ...)` | Usuario creado, correoElectronico con estado PENDIENTE |
| 23 | Crear usuario con correo inválido | email mal formado | `NuevoUsuario(id, email, ...)` | Error propagado desde el VO |
| 24 | Correo() getter funciona | Usuario creado | `usuario.Correo()` | Retorna la dirección del VO |
| 25 | EstadoVerificacionCorreo() getter funciona | Usuario creado | `usuario.EstadoVerificacionCorreo()` | Retorna el estado del VO |
| 26 | VerificarCorreo delega en VO | Usuario PENDIENTE | `usuario.VerificarCorreo()` | estado → VERIFICADO vía VO |
| 27 | MarcarEnlaceExpirado delega en VO | Usuario PENDIENTE | `usuario.MarcarEnlaceExpirado()` | estado → ENLACE_EXPIRADO vía VO |
| 28 | MarcarVerificacionFallida NO existe | — | Llamar a `usuario.MarcarVerificacionFallida()` | COMPILE ERROR (método eliminado) |
| 29 | Reconstruir desde BD con VO | Todos los campos | `NewUsuarioFromPersistence(id, ..., correoElectronico, ...)` | Usuario con el mismo estado |
| 30 | Evento UsuarioCreado incluye correo | Creación exitosa | `usuario.PullEventos()` | Evento con correo en payload |

---

## 7. Tareas de implementación (orden sugerido)

```
Fase 1: Dominio — Value Object
  1.1 Limpiar EstadoVerificacionCorreo (eliminar VERIFICACION_FALLIDA + transiciones)
  1.2 Actualizar tests de estado_verificacion_correo_test.go
  1.3 Crear correo_electronico.go con el VO
  1.4 Escribir tests del VO

Fase 2: Dominio — Usuario
  2.1 Reemplazar campos `correo` y `estadoVerificacionCorreo` por `correoElectronico`
  2.2 Actualizar NuevoUsuario() para crear el VO internamente
  2.3 Actualizar NewUsuarioFromPersistence() para recibir VO
  2.4 Actualizar métodos delegados (VerificarCorreo, SolicitarReenvio, etc.)
  2.5 Eliminar MarcarVerificacionFallida()
  2.6 Actualizar eventos.go (quitar evento VerificacionFallida si existe)
  2.7 Actualizar tests de usuario_test.go

Fase 3: Infraestructura
  3.1 Actualizar UsuarioModel.ToDomain() para construir el VO
  3.2 Actualizar UsuarioModel.FromDomain() para extraer datos del VO
  3.3 Actualizar usuario_repositorio.go si es necesario
  3.4 Agregar migración de datos: UPDATE usuarios SET estado_verificacion_correo = 'PENDIENTE_VERIFICACION' WHERE estado_verificacion_correo = 'VERIFICACION_FALLIDA'

Fase 4: Verificación
  4.1 Compilar todo el proyecto (go build ./...)
  4.2 Ejecutar todos los tests (go test ./...)
  4.3 Verificar que ningún test referencia VERIFICACION_FALLIDA
```

---

## 8. Criterios de aceptación

| # | Criterio |
|---|----------|
| AC-1 | `CorreoElectronico` es un Value Object inmutable en el paquete de dominio |
| AC-2 | La máquina de estados de verificación tiene exactamente 4 estados (sin VERIFICACION_FALLIDA) |
| AC-3 | `Usuario` ya no tiene los campos planos `correo` y `estadoVerificacionCorreo` |
| AC-4 | `Usuario.Correo()` retorna la dirección desde el VO |
| AC-5 | `NewUsuarioFromPersistence` recibe un `*CorreoElectronico` completo |
| AC-6 | `MarcarVerificacionFallida()` ya no existe en `Usuario` |
| AC-7 | Todos los tests existentes siguen pasando (excepto los que probaban VERIFICACION_FALLIDA) |
| AC-8 | El build del proyecto compila sin errores |
| AC-9 | Los registros existentes en BD con VERIFICACION_FALLIDA se migran a PENDIENTE_VERIFICACION |

---

## 9. Checklist de Validación

- [ ] ¿VERIFICACION_FALLIDA fue eliminado de las constantes?
- [ ] ¿VERIFICACION_FALLIDA fue eliminado del mapa de transiciones?
- [ ] ¿Los tests de estado_verificacion_correo ya no referencian VERIFICACION_FALLIDA?
- [ ] ¿CorreoElectronico encapsula `direccion` y `estado` como campos privados?
- [ ] ¿CorreoElectronico valida el formato del email en su constructor?
- [ ] ¿Usuario ya no tiene campos planos `correo` ni `estadoVerificacionCorreo`?
- [ ] ¿Usuario.Correo() delega en correoElectronico.Direccion()?
- [ ] ¿Usuario.EstadoVerificacionCorreo() delega en correoElectronico.Estado()?
- [ ] ¿Usuario.VerificarCorreo() delega en correoElectronico.Verificar()?
- [ ] ¿Usuario.MarcarVerificacionFallida() fue eliminado?
- [ ] ¿NewUsuarioFromPersistence recibe un *CorreoElectronico?
- [ ] ¿UsuarioModel.ToDomain() construye el VO correctamente?
- [ ] ¿UsuarioModel.FromDomain() extrae datos del VO correctamente?
- [ ] ¿Los getters de CorreoElectronico son públicos (exportados)?
- [ ] ¿Hay tests para cada escenario de la tabla TDD?
- [ ] ¿La migración de datos existe para VERIFICACION_FALLIDA → PENDIENTE_VERIFICACION?
- [ ] ¿El proyecto compila y todos los tests pasan?

---

## 10. Especificaciones Relacionadas

| Especificación | Relación |
|----------------|----------|
| `registro/spec_registro.md` | Esta spec afecta la entidad Usuario que el registro consume |
| `sesiones/login_spec.md` | El login obtiene el correo del usuario (getter se mantiene) |
| `autorizacion/spec-tenant-management.md` | Los tenants referencian usuarios (no se afecta) |
| `autorizacion/spec-rbac-authorization.md` | RBAC construye sobre Usuario existente |
