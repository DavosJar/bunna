---
title: Especificación de Dominio — Fincas
version: 2.0
date_created: 2026-05-23
owner: Equipo Catastro
tags: fincas, lotes, dominio
---
 
# Especificación de Dominio — Fincas
 
> **Propósito**: Definir exclusivamente la capa de dominio del microservicio `fincas`: entidades, identidad, relaciones, reglas de tenencia y operaciones que alteran estado. No incluye validaciones de formato, contratos HTTP, persistencia, casos de uso de aplicación ni ninguna otra capa.
 
---
 
## 1. Glosario
 
| Término | Definición |
|---------|-----------|
| **Finca** | Unidad principal de gestión agrícola. Agregado raíz. Pertenece a un único usuario. |
| **Lote** | Subdivisión espacial de una Finca. No puede existir sin una Finca. |
| **usuarioID** | Identificador del usuario propietario de la Finca. Referencia lógica al sistema de identidad; no hay FK física. |
| **tenantID** | Identificador de la organización a la que pertenece la Finca. Opcional en versión inicial. |
| **Propietario** | Usuario cuyo `usuarioID` coincide con el `usuarioID` registrado en la Finca. |
| **Permiso atómico** | Capacidad indivisible con formato `modulo:recurso:verbo`. Definida en el catálogo de identidad. |
 
---
 
## 2. Entidades del Dominio
 
### 2.1 Finca
 
Agregado raíz. Representa la unidad principal de producción agrícola de un usuario.
 
| Atributo | Tipo | Notas |
|----------|------|-------|
| `id` | identidad | Generado al crear. Inmutable. |
| `nombre` | string | Requerido. |
| `ubicacion` | string | Requerido. |
| `descripcion` | string | Opcional. Vacío por defecto. |
| `usuarioID` | identidad | ID del propietario. Asignado desde el contexto de autenticación, nunca desde el cliente. Inmutable. |
| `tenantID` | identidad | Opcional en v1. Futuro: obligatorio para multi-tenant. |
| `createdAt` | timestamp UTC | Asignado al crear. Inmutable. |
| `updatedAt` | timestamp UTC | Actualizado en cada modificación. |
 
Una Finca tiene cero o más Lotes. Los Lotes **no** se almacenan como lista embebida en la Finca; se referencian desde el Lote mediante `fincaID`.
 
### 2.2 Lote
 
Entidad dependiente del agregado Finca. Representa una subdivisión espacial dentro de una Finca.
 
| Atributo | Tipo | Notas |
|----------|------|-------|
| `id` | identidad | Generado al crear. Inmutable. |
| `fincaID` | identidad | Referencia a la Finca propietaria. Obligatorio. Inmutable. |
| `nombre` | string | Requerido. |
| `area` | decimal | Requerido. Positivo. En hectáreas. |
| `descripcion` | string | Opcional. Vacío por defecto. |
| `createdAt` | timestamp UTC | Asignado al crear. Inmutable. |
| `updatedAt` | timestamp UTC | Actualizado en cada modificación. |
 
Un Lote siempre referencia una Finca existente. No existe Lote sin Finca.
 
---
 
## 3. Relación entre Entidades
 
```
Finca 1 ───< Lote
```
 
- Una Finca puede tener cero o muchos Lotes.
- Un Lote pertenece a exactamente una Finca (vía `fincaID`).
- La relación se navega **desde el Lote hacia la Finca**, no mediante una lista embebida en Finca.
- La eliminación de una Finca implica la eliminación en cascada de todos sus Lotes.
---
 
## 4. Reglas de Tenencia
 
Estas reglas gobiernan qué usuario puede operar sobre qué recurso. Son invariantes del dominio, independientes de cualquier capa externa.
 
| ID | Regla |
|----|-------|
| **TEN-01** | Solo el propietario de una Finca puede modificarla o eliminarla. |
| **TEN-02** | Solo el propietario de la Finca padre puede crear, modificar o eliminar sus Lotes. |
| **TEN-03** | El `usuarioID` de una Finca es inmutable tras la creación; nunca puede ser reasignado por operación directa. |
| **TEN-04** | El `fincaID` de un Lote es inmutable tras la creación; un Lote no puede moverse a otra Finca. |
| **TEN-05** | Un usuario con rol `sys_admin` actúa sobre cualquier Finca o Lote sin restricción de propiedad. Este rol es global y se verifica desde el contexto de autenticación. |
| **TEN-06** | La verificación de propiedad se realiza comparando el `usuarioID` del recurso con el `usuarioID` extraído del contexto de autenticación. |
 
---
 
## 5. Operaciones del Dominio
 
El dominio expone únicamente las operaciones que alteran el estado de sus entidades. Las operaciones que involucran una sola entidad se definen directamente sobre ella. Las operaciones que requieren coordinar dos entidades se resuelven en un **servicio de dominio**, que orquesta usando los repositorios definidos en el propio dominio.
 
### 5.1 Operaciones sobre Finca
 
Operaciones propias de la entidad Finca, sin necesidad de coordinar con Lote.
 
#### `CrearFinca(nombre, ubicacion, descripcion, usuarioID) → Finca`
 
Produce una nueva Finca en estado válido. Asigna un nuevo `id`, registra `usuarioID` como propietario y establece `createdAt` y `updatedAt`.
 
**Postcondiciones:** la Finca resultante tiene todos sus campos obligatorios asignados y su `id` es único.
 
---
 
#### `ActualizarFinca(nombre, ubicacion, descripcion) → Finca`
 
Actualiza los campos editables de la Finca. Actualiza `updatedAt`.
 
**Campos no modificables:** `id`, `usuarioID`, `tenantID`, `createdAt`.
 
---
 
#### `EsPropietario(usuarioID) → bool`
 
Verifica si el `usuarioID` provisto coincide con el propietario de la Finca. Las operaciones que exigen propiedad (reglas **TEN-01** y **TEN-02**) invocan este método antes de proceder.
 
---
 
### 5.2 Operaciones sobre Lote
 
Operaciones propias de la entidad Lote, sin necesidad de coordinar con Finca.
 
#### `ActualizarLote(nombre, area, descripcion) → Lote`
 
Actualiza los campos editables del Lote. Actualiza `updatedAt`.
 
**Campos no modificables:** `id`, `fincaID`, `createdAt`.
 
---
 
### 5.3 Servicio de Dominio: `FincaService`
 
Resuelve las operaciones que requieren coordinar `Finca` y `Lote`. Recibe `FincaRepositorio` y `LoteRepositorio` como dependencias, ambas interfaces del dominio.
 
---
 
#### `RegistrarLoteEnFinca(finca, nombre, area, descripcion) → Lote`
 
Crea un nuevo Lote asociado a la Finca dada. La Finca debe existir y ya fue resuelta por el llamador.
 
**Precondición:** `finca` es una instancia válida y persistida. La verificación de propiedad (regla **TEN-02**) se realiza antes de invocar esta operación usando `finca.EsPropietario(usuarioID)`.
 
**Comportamiento:** produce un Lote con `fincaID` igual al `id` de la Finca. Asigna `id`, `createdAt` y `updatedAt`.
 
**Postcondición:** el Lote resultante referencia una Finca existente y tiene todos sus campos obligatorios asignados.
 
---
 
#### `EliminarFincaConLotes(finca, confirmado) → error`
 
Elimina una Finca y todos sus Lotes. Antes de eliminar, consulta cuántos Lotes tiene la Finca usando `LoteRepositorio`.
 
**Regla:** si la Finca tiene Lotes y `confirmado == false` → error `ErrFincaConLotes`.
 
**Regla:** si `confirmado == true` o la Finca no tiene Lotes → elimina primero todos los Lotes asociados y luego la Finca.
 
**Postcondición:** ni la Finca ni ninguno de sus Lotes permanecen en el sistema.
 
---
 
## 6. Errores del Dominio
 
Los errores del dominio representan violaciones de invariantes o estados inválidos. No son errores de formato ni de infraestructura.
 
| Error | Significado |
|-------|-------------|
| `ErrFincaNoEncontrada` | La Finca referenciada no existe. |
| `ErrLoteNoEncontrado` | El Lote referenciado no existe. |
| `ErrNoPropietario` | El usuario no es propietario del recurso que intenta operar. |
| `ErrFincaConLotes` | Se intenta eliminar una Finca con Lotes sin confirmación explícita. |
 
---
 
## 7. Interfaces de Repositorio
 
Las interfaces de repositorio pertenecen al dominio. Definen el contrato de persistencia sin presuponer ninguna implementación. Se implementan en la capa de infraestructura.
 
### FincaRepositorio
 
| Método | Descripción |
|--------|-------------|
| `Crear(finca)` | Persiste una Finca nueva. |
| `ObtenerPorID(id)` | Retorna la Finca con ese ID, o indica ausencia. |
| `ListarPorUsuario(usuarioID)` | Retorna todas las Fincas del usuario. |
| `ListarTodas()` | Retorna todas las Fincas del sistema (uso exclusivo de `sys_admin`). |
| `Actualizar(finca)` | Persiste los cambios de una Finca existente. |
| `Eliminar(id)` | Elimina la Finca. |
 
### LoteRepositorio
 
| Método | Descripción |
|--------|-------------|
| `Crear(lote)` | Persiste un Lote nuevo. |
| `ObtenerPorID(id)` | Retorna el Lote con ese ID, o indica ausencia. |
| `ListarPorFinca(fincaID)` | Retorna todos los Lotes de la Finca. |
| `ContarPorFinca(fincaID)` | Retorna la cantidad de Lotes de la Finca. Usado por `FincaService` para evaluar la eliminación. |
| `EliminarPorFinca(fincaID)` | Elimina todos los Lotes de la Finca. Usado por `FincaService` en eliminación con cascada. |
| `Actualizar(lote)` | Persiste los cambios de un Lote existente. |
| `Eliminar(id)` | Elimina el Lote. |
 