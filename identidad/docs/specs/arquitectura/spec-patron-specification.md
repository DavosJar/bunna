---
title: Patrón Specification — Blueprint de Búsqueda y Filtrado
version: 2.0
date_created: 2026-05-23
owner: Equipo Identidad
tags: arquitectura, patron-specification, blueprint, integracion
---

# Patrón Specification — Blueprint de Búsqueda y Filtrado

## 1. Propósito y Alcance

Este documento define la **receta arquitectónica estándar** para implementar búsquedas complejas, filtrado dinámico, ordenación y paginación en **cualquier módulo nuevo** del sistema.

El objetivo es garantizar que cada vez que una nueva entidad requiera consultas complejas, se implemente de manera uniforme, transferible y estrictamente segura contra ataques de inyección SQL o enumeración, desacoplando los criterios de negocio del motor de base de datos.

## 2. Elementos Base (Shared Domain)

Todo nuevo módulo debe reutilizar los siguientes conceptos base ya existentes en la carpeta compartida (`shared/domain`):

*   **CriterioFiltro**: Define una condición atómica (`Campo`, `Operador`, `Valor`).
*   **Paginacion**: Define los límites y el orden de los resultados (`Pagina`, `TamanoPagina`, `Ordenaciones`).
*   **Ordenacion** y **TipoOrdenacion**: Define la dirección del orden (`ASC`, `DESC`) para un campo específico.

## 3. Pasos para Implementar en un Nuevo Módulo

Para dotar a una nueva entidad (ej. `Factura`, `Rol`, `Auditoria`) de capacidades de búsqueda, se deben seguir obligatoriamente estos 3 pasos en la capa de **Dominio**.

### Paso 1: Crear la Especificación de la Entidad
En el paquete de dominio de la entidad, crear un tipo `Especificacion[Entidad]` que envuelva una lista de criterios de filtro compartidos. Esto sirve como contenedor de reglas de negocio para la búsqueda.

### Paso 2: Definir el Whitelist de Seguridad (`ColumnasPermitidas`)
En el mismo archivo de dominio, se debe declarar de forma inmutable un mapa estático de strings a booleanos (`ColumnasPermitidas`). 
**Regla de Oro:** Solo los campos de la entidad que sean seguros para exponerse en búsquedas públicas/administrativas deben tener el valor `true`. Jamás incluir IDs internos si permiten enumeración, contraseñas, hashes o datos confidenciales no filtrables.

### Paso 3: Definir la Firma del Repositorio
La interfaz del repositorio en el dominio DEBE incluir una firma estandarizada para la consulta. El nombre recomendado es `Listar` o `Buscar`, pero la estructura de los parámetros debe ser exacta:

**Firma Estándar Obligatoria:**
`Listar(ctx context.Context, especificacion Especificacion[Entidad], paginacion domain.Paginacion) ([]*[Entidad], error)`

*(Reemplazar `[Entidad]` por el nombre real, ej. `Usuario`, `Rol`)*.

## 4. Reglas de Implementación en la Capa de Persistencia (Infraestructura)

El repositorio concreto (ej. Postgres, MongoDB) que implemente la firma anterior debe obedecer estrictamente este flujo de resolución:

1.  **Mapeo de Nombres (Domain to DB)**: El repositorio debe definir internamente un diccionario que traduzca el nombre del campo del dominio (ej. `fechaCreacion`) al nombre real de la columna en la base de datos (ej. `fecha_creacion`).
2.  **Validación de Whitelist (Seguridad)**: Antes de procesar CUALQUIER filtro u ordenación, se debe verificar que el campo solicitado exista y sea `true` en el mapa `ColumnasPermitidas` definido en el dominio.
3.  **Degradación Elegante (Graceful Degradation)**: Si un cliente envía un filtro con un campo no permitido, inexistente, o malicioso, el repositorio **DEBE ignorar el filtro silenciosamente** y continuar. No debe romper la consulta ni devolver error de base de datos.
4.  **Traducción de Operadores**: Soporte mínimo obligatorio para operadores de igualdad (`=`), diferencia (`!=`) y similitud (`LIKE`). Si el campo lo requiere (ej. fechas), implementar rangos.
5.  **Paginación Segura**: 
    *   Si `Pagina` < 1, forzar a 1.
    *   Si `TamanoPagina` < 1, forzar a un límite predeterminado seguro (ej. 10 o 50).
    *   Toda paginación debe traducirse a `LIMIT` y `OFFSET` (o su equivalente) en el motor subyacente de forma parametrizada.

## 5. Implementación de Referencia

Para ver la aplicación correcta de este blueprint, el módulo de **Usuarios** (`internal/usuarios/domain/usuario` e `internal/usuarios/infrastructure/persistence/postgres`) es la implementación canónica y debe usarse como guía de código.
