---
title: Filtros Seguros por Entidad de Dominio
version: 1.0
date_created: 2026-05-23
owner: Equipo Identidad
tags: seguridad, filtros, especificacion, dominio
---

# Filtros Seguros por Entidad de Dominio

## 1. Propósito y Alcance

Definir los campos de filtro y ordenación que pueden exponerse de forma segura para cada entidad de dominio en el sistema de identidad. Cada entidad tiene un `ColumnasPermitidas` que actúa como whitelist contra inyección SQL y exposición no autorizada de datos sensibles.

Queda explícitamente fuera del alcance la implementación técnica de los filtros (capa de infraestructura).

## 2. Definiciones

| Término | Definición |
|---------|------------|
| Entidad de dominio | Objeto con identidad y ciclo de vida dentro de un dominio acotado |
| Filtro seguro | Campo cuyo valor puede usarse como criterio de búsqueda sin riesgo de exponer datos sensibles o permitir enumeración maliciosa |
| PII | Personally Identifiable Information — datos que permiten identificar a una persona |
| Enumeration attack | Técnica donde un atacante deduce información mediante la variación sistemática de parámetros de búsqueda |
| Whitelist de columnas | Mapa de strings a boolean que define qué campos pueden usarse en cláusulas WHERE y ORDER BY |

## 3. Requisitos, Restricciones y Guías

- **SEC-001**: Ningún campo que contenga hashes, tokens, passwords o secretos criptográficos debe estar en `ColumnasPermitidas`.
- **SEC-002**: Campos que expongan información de seguridad interna (intentos fallidos, bloqueos) deben evaluarse caso por caso: pueden exponerse para administración pero no para autousuario.
- **SEC-003**: El `id` interno de una entidad no debe estar en `ColumnasPermitidas` si su exposición permite enumerar registros. Si se requiere búsqueda por ID, debe hacerse mediante un método dedicado (ej. `ObtenerPorID`).
- **SEC-004**: Datos PII (email, teléfono, IP) pueden filtrarse solo con operadores seguros (igualdad exacta, LIKE controlado). Nunca exponer directamente en listados públicos sin autorización explícita.
- **SEC-005**: Toda ordenación (ORDER BY) debe limitarse a las mismas columnas permitidas para filtros.
- **CON-001**: Los operadores permitidos son `=`, `!=`, `LIKE`. El operador `LIKE` no debe permitir comodín inicial (`%...`) para evitar búsquedas de fuerza bruta sobre PII.
- **CON-002**: Cada entidad define su propio `ColumnasPermitidas` en su archivo de especificación (ej. `especificacion_usuario.go`).
- **PAT-001**: Para cada entidad, el map de columnas permitidas se define como `var ColumnasPermitidas = map[string]bool{...}` en el mismo archivo que la `Especificacion*`.

## 4. Entidades y sus Filtros Seguros

### 4.1 Usuario (`internal/usuarios/domain/usuario/`)

Estado actual (implementado):

| Campo | Tipo | Operadores | Seguro | Notas |
|-------|------|-----------|--------|-------|
| `nombre` | string | `=`, `!=`, `LIKE` | ✅ | No PII directo |
| `apellido` | string | `=`, `!=`, `LIKE` | ✅ | No PII directo |
| `correo` | string | `=`, `!=`, `LIKE` | ✅ | Es PII pero es identificador lógico del usuario, necesario para búsqueda administrativa. Case-insensitive. |
| `fechaCreacion` | datetime | `=`, `!=`, `LIKE` | ✅ | Metadata temporal |
| `fechaActualizacion` | datetime | `=`, `!=`, `LIKE` | ✅ | Metadata temporal |
| `estado` | enum string | `=`, `!=` | ✅ | Valores controlados: `ACTIVO`, `INACTIVO`, `BLOQUEADO`, `PENDIENTE_ELIMINACION`, `NO_VERIFICADO` |
| `telefono` | string | `=`, `!=`, `LIKE` | ✅ | PII controlado. Útil para búsqueda administrativa. |
| `estadoVerificacionCorreo` | enum string | `=`, `!=` | ✅ | Valores controlados: `PENDIENTE_VERIFICACION`, `VERIFICADO`, `ENLACE_EXPIRADO`, `REENVIO_SOLICITADO` |

**No expuesto (y no debe estarlo):**
- `id` — No en `ColumnasPermitidas`. La búsqueda por ID se hace vía `ObtenerPorID()`. Previene enumeration attacks.
- `eventos`/`eventosUsuario` — No es campo de búsqueda, es cola de eventos de dominio.
- `correoElectronico` (objeto completo) — Se expone su `Direccion` vía mapeo a `correo`.

### 4.2 CredencialesUsuario (`internal/seguridad/domain/`)

Estado actual (implementado):

| Campo | Tipo | Operadores | Seguro | Notas |
|-------|------|-----------|--------|-------|
| `usuarioID` | string | `=`, `!=` | ✅ | Para busqueda por usuario padre |
| `activo` | bool | `=`, `!=` | ✅ | Estado booleano |
| `correoVerificado` | bool | `=`, `!=` | ✅ | Estado booleano |
| `intentosFallidos` | int | `=`, `!=`, `LIKE` | ⚠️ | Expone política de seguridad interna. Aceptable solo para administradores. |
| `bloqueadoHasta` | datetime | `=` | ⚠️ | Expone tiempo de bloqueo. Aceptable solo para administradores. |

**No expuesto (y no debe estarlo):**
- `passwordHash` — Secreto criptográfico. NUNCA exponer.

### 4.3 Sesion (`internal/sesiones/domain/`)

No implementado actualmente. Propuesta:

| Campo | Tipo | Operadores | Seguro | Notas |
|-------|------|-----------|--------|-------|
| `usuarioID` | string | `=`, `!=` | ✅ | Para listar sesiones de un usuario |
| `estado` | enum string | `=`, `!=` | ✅ | Valores controlados: `ACTIVA`, `EXPIRADA`, `REVOCADA` |
| `ipOrigen` | string | `=`, `!=` | ⚠️ | PII. Exponer solo a administradores. |
| `fechaCreacion` | datetime | `=`, `!=` | ✅ | Metadata temporal |
| `fechaActualizacion` | datetime | `=`, `!=` | ✅ | Metadata temporal |
| `fechaExpiracionAccess` | datetime | `=`, `!=` | ✅ | Metadata temporal |
| `fechaExpiracionRefresh` | datetime | `=`, `!=` | ✅ | Metadata temporal |
| `ultimaActividad` | datetime | `=`, `!=` | ✅ | Metadata temporal |
| `contadorRefrescos` | int | `=`, `!=` | ✅ | Contador de renovaciones |

**No expuesto (y no debe estarlo):**
- `id` — Búsqueda vía `ObtenerPorID()`. Previene enumeration.
- `accessTokenHash` — Secreto. NUNCA exponer.
- `refreshTokenHash` — Secreto. NUNCA exponer.

### 4.4 IntentoPorIP (`internal/seguridad/domain/`)

No implementado actualmente. Propuesta:

| Campo | Tipo | Operadores | Seguro | Notas |
|-------|------|-----------|--------|-------|
| `ip` | string | `=`, `!=` | ⚠️ | PII. Exponer solo a administradores de seguridad. |
| `contador` | int | `=`, `!=`, `>=`, `<=` | ⚠️ | Expone actividad de seguridad. Solo admins. |
| `ventanaInicio` | datetime | `=`, `!=` | ⚠️ | Expone patrones de ataque. Solo admins. |
| `bloqueadoHasta` | datetime | `=`, `!=` | ⚠️ | Estado de bloqueo. Solo admins. |

**No expuesto (y no debe estarlo):**
- `id` — Búsqueda vía `ObtenerPorID()`.
- Todos los campos de IntentoPorIP son sensiblemente de seguridad y solo deben exponerse a roles administrativos específicos.

## 5. Operadores y su Seguridad

| Operador | Uso seguro | Riesgos |
|----------|-----------|---------|
| `=` | ✅ Búsqueda exacta. Seguro con whitelist. | Bajo |
| `!=` | ✅ Negación exacta. Seguro con whitelist. | Bajo |
| `LIKE` | ⚠️ Requiere validación. No permitir `%...` inicial. | Medio: permite fuerza bruta si se usa en PII sin control |
| `>=`, `<=`, `BETWEEN` | ⚠️ Solo para campos numéricos o datetime. Requiere whitelist. | Bajo si se controla con whitelist |
| `<`, `>` | ⚠️ Solo para campos numéricos o datetime. | Bajo si se controla con whitelist |

## 6. Criterios de Aceptación

- **AC-001**: Dado un `ColumnasPermitidas` de Usuario, ningún campo que contenga `password`, `hash`, `token`, o `secret` en su nombre debe estar presente.
- **AC-002**: Dada una búsqueda por `id` de Usuario, el sistema debe usar `ObtenerPorID()` y no un filtro genérico por `ColumnasPermitidas`.
- **AC-003**: Dado un filtro `LIKE` sobre `correo`, la implementación debe ser case-insensitive (LOWER).
- **AC-004**: Dada una ordenación por un campo no listado en `ColumnasPermitidas`, la implementación debe ignorar esa ordenación (no fallar).
- **AC-005**: Dado un filtro por un campo no listado en `ColumnasPermitidas`, la implementación debe ignorar ese filtro (no fallar).
- **AC-006**: Cada entidad de dominio que requiera búsqueda debe tener su propio `ColumnasPermitidas` y su `Especificacion*` en el mismo archivo.

## 7. Estrategia de Automatización de Pruebas

- **Test por entidad**: Cada `ColumnasPermitidas` debe tener un test que verifique que los campos sensibles conocidos NO están en el map.
- **Test de integración**: Para cada repositorio, verificar que filtrar por campo no permitido retorna resultados sin aplicar ese filtro (graceful degradation, no error).
- **Test de seguridad**: Verificar que no hay campos que expongan `passwordHash`, `accessTokenHash`, o `refreshTokenHash` en ninguna especificación.

## 8. Dependencias e Integraciones Externas

- `internal/shared/domain/specification.go` — Tipos base `CriterioFiltro`, `Paginacion`, `Ordenacion`.
- `internal/usuarios/domain/usuario/especificacion_usuario.go` — Definición actual de filtros de Usuario.
- `internal/seguridad/domain/especificacion_credenciales.go` — Definición actual de filtros de Credenciales.
