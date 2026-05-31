---
title: Especificación de Capa de Aplicación — Casos de Uso
version: 1.0
date_created: 2026-05-24
owner: Fincas Team
tags: aplicacion, casos-de-uso, fincas, lotes, muestras, diagnosticos, reportes
---

# Especificación de Capa de Aplicación — Casos de Uso

> **Propósito**: Definir los contratos de entrada, salida, validaciones, flujo de orquestación, política transaccional y eventos publicados para cada caso de uso de la capa de aplicación del microservicio `fincas`. No contiene código ni estructuras de implementación concretas.

---

## 1. Estructura de Carpetas

La capa de aplicación se organiza en un paquete por caso de uso dentro de `internal/application/usecases/`. Cada paquete contiene tres archivos:

- `command.go` → estructura de entrada con validaciones
- `usecase.go` → lógica de orquestación, validación de permisos, validación de tenant, lógica de dominio, persistencia y publicación de eventos
- `salida.go` → estructura de respuesta

```
internal/application/usecases/
├── registrarfinca/
│   ├── command.go
│   ├── usecase.go
│   └── salida.go
├── desactivarfinca/
│   ├── command.go
│   ├── usecase.go
│   └── salida.go
├── agregarlote/
│   ├── command.go
│   ├── usecase.go
│   └── salida.go
├── eliminarlote/
│   ├── command.go
│   ├── usecase.go
│   └── salida.go
├── tomarmuestra/
│   ├── command.go
│   ├── usecase.go
│   └── salida.go
├── listarmuestrasporlote/
│   ├── command.go
│   ├── usecase.go
│   └── salida.go
├── solicitardiagnosticomanual/
│   ├── command.go
│   ├── usecase.go
│   └── salida.go
├── registrarinferencia/
│   ├── command.go
│   ├── usecase.go
│   └── salida.go
├── aceptardiagnostico/
│   ├── command.go
│   ├── usecase.go
│   └── salida.go
├── rechazardiagnostico/
│   ├── command.go
│   ├── usecase.go
│   └── salida.go
└── generarreporteporlote/
    ├── command.go
    ├── usecase.go
    └── salida.go
```

---

## 2. Glosario

| Término | Definición |
|---------|-----------|
| **AuthContext** | Contexto de autenticación inyectado desde el handler. Contiene UsuarioID, TenantID (opcional en v1, obligatorio en v2+) y Permisos ([]string). |
| **EventPublisher** | Interfaz para publicar eventos en RabbitMQ. Expone un único método Publish(ctx, routingKey, event). |
| **Unit of Work** | Abstracción transaccional que agrupa múltiples operaciones de persistencia en una sola transacción atómica. Si la transacción falla, no se publica el evento. |
| **ErrForbidden** | Error de aplicación: el usuario no tiene el permiso requerido. Mapea a HTTP 403. |
| **ErrNotFound** | Error de aplicación: el recurso no existe en BD, o existe pero es de otro tenant. Mapea a HTTP 404. |
| **ErrConflictoEstado** | Error de aplicación: transición de estado no permitida (ej. aceptar un diagnóstico que ya no está PENDIENTE). Mapea a HTTP 409. |
| **ErrValidacion** | Error de aplicación: campo inválido (coordenadas fuera de rango, nombre vacío, área negativa, etc.). Mapea a HTTP 400. |
| **ErrFincaConLotes** | Error de aplicación específico: se intenta desactivar una finca con lotes activos sin confirmación explícita. Mapea a HTTP 409. |

---

## 3. Reglas Transversales

### 3.1 Validación de permisos

Todo caso de uso que recibe una solicitud de usuario (no operaciones internas) debe validar que el permiso requerido esté presente en AuthContext.Permisos. Si no está, retorna ErrForbidden.

### 3.2 Regla de tenant (404 vs 403)

Cuando un recurso existe en BD pero su TenantID no coincide con AuthContext.TenantID, el caso de uso retorna ErrNotFound (no ErrForbidden). Esto evita filtrar información sobre la existencia de recursos entre distintos tenants. Es responsabilidad del caso de uso verificar la tenencia después de cargar el recurso.

### 3.3 Política de Unit of Work

- Casos de uso con **una sola escritura** en BD: sin Unit of Work. Se persiste directamente; si falla, no se publica el evento.
- Casos de uso con **múltiples escrituras** en BD: se envuelven en Unit of Work. Solo si la transacción completa tiene éxito se publica el evento.
- Casos de uso de **solo lectura**: sin Unit of Work ni eventos.

La publicación del evento SIEMPRE ocurre después de confirmar la escritura en BD. Nunca antes ni en medio de una transacción.

### 3.4 Errores de dominio vs errores de aplicación

Los casos de uso invocan errores definidos en la capa de dominio (ej. ErrTransicionEstadoNoPermitida, ErrFincaConLotes, ErrDiagnosticoNoEncontrado). El caso de uso traduce esos errores de dominio a errores de aplicación (ErrConflictoEstado, ErrNotFound, ErrValidacion) según la tabla del mapa de errores. No se propagan errores de dominio directamente a la capa de presentación.

---

## 4. Casos de Uso

---

### 4.1 fincas:RegistrarFinca

**Permiso requerido:** CREAR_FINCA
**Ubicación:** `internal/application/usecases/registrarfinca/`

#### 4.1.1 command.go

Campos de entrada:
- Nombre (string): requerido. Entre 3 y 200 caracteres. Sin saltos de línea.
- Ubicacion (string): requerido. Entre 1 y 500 caracteres.
- Descripcion (string): opcional. Máximo 1000 caracteres si se provee.

Validaciones:
- Nombre vacío o menor a 3 caracteres → ErrValidacion.
- Nombre mayor a 200 caracteres → ErrValidacion.
- Ubicacion vacío → ErrValidacion.
- Ubicacion mayor a 500 caracteres → ErrValidacion.
- Descripcion mayor a 1000 caracteres → ErrValidacion.

#### 4.1.2 usecase.go

Dependencias inyectadas: FincaRepositorio, GeneradorID, EventPublisher.

Flujo:
1. Validar que CREAR_FINCA está en AuthContext.Permisos. Si no → ErrForbidden.
2. Generar ID con GeneradorID.
3. Construir entidad Finca vía el constructor del dominio, pasando nombre, ubicacion, descripcion, UsuarioID del AuthContext.
4. Asignar a la finca el ID generado y el TenantID desde AuthContext (puede ser nil).
5. Persistir la finca con FincaRepositorio.Crear.
6. Si la persistencia falla, retornar error. No publicar evento.
7. Si la persistencia es exitosa, publicar evento FincaCreada.

Transaccionalidad: una sola escritura → sin Unit of Work.

#### 4.1.3 salida.go

Campos de respuesta: ID, Nombre, Ubicacion, Descripcion, Estado (siempre "ACTIVA"), CreatedAt.

#### 4.1.4 Evento RabbitMQ

Routing key: fincas.v1.finca.creada

Campos del evento: EventID, FincaID, Nombre, UsuarioID, TenantID (opcional), OcurredAt.

---

### 4.2 fincas:DesactivarFinca

**Permiso requerido:** DESACTIVAR_FINCA
**Transición de estado:** ACTIVA → PENDIENTE_ELIMINACION
**Ubicación:** `internal/application/usecases/desactivarfinca/`

#### 4.2.1 command.go

Campos de entrada:
- FincaID (string): UUID válido, requerido.
- Confirmar (booleano): por defecto false. Indica si el usuario confirma la desactivación cuando la finca tiene lotes activos.

Validaciones:
- FincaID UUID no vacío → si no, ErrValidacion.

#### 4.2.2 usecase.go

Dependencias inyectadas: FincaRepositorio, FincaService (servicio de dominio), GeneradorID, EventPublisher.

Flujo:
1. Validar que DESACTIVAR_FINCA está en AuthContext.Permisos. Si no → ErrForbidden.
2. Cargar finca con FincaRepositorio.ObtenerPorID usando FincaID del comando.
3. Si no existe → ErrNotFound.
4. Validar tenencia: comparar TenantID de la finca vs AuthContext.TenantID. Si no coincide → ErrNotFound (el recurso existe pero no es de este tenant).
5. Contar lotes activos de la finca con FincaRepositorio.ContarLotes.
6. Invocar FincaService.EliminarFincaConLotes pasando la finca, la cantidad de lotes y el booleano Confirmar.
   - Si la finca tiene lotes y Confirmar es false → el servicio retorna error de dominio ErrFincaConLotes. El caso de uso lo traduce a ErrValidacion con el mensaje "La finca tiene N lote(s) asociado(s). Confirma la eliminación".
   - Si la transición ACTIVA → PENDIENTE_ELIMINACION no es válida (la finca ya está en otro estado) → el servicio retorna error de dominio ErrTransicionEstadoNoPermitida. El caso de uso lo traduce a ErrConflictoEstado.
7. Persistir el cambio de estado con FincaRepositorio.Actualizar.
8. Si la persistencia falla, retornar error. No publicar evento.
9. Si la persistencia es exitosa, publicar evento FincaDesactivada.

Transaccionalidad: una sola escritura → sin Unit of Work. La consulta de cantidad de lotes es de solo lectura y ocurre antes de la escritura.

#### 4.2.3 salida.go

Campos de respuesta: ID, Estado ("PENDIENTE_ELIMINACION"), UpdatedAt.

#### 4.2.4 Evento RabbitMQ

Routing key: fincas.v1.finca.desactivada

Campos del evento: EventID, FincaID, EstadoAnterior ("ACTIVA"), EstadoNuevo ("PENDIENTE_ELIMINACION"), UsuarioID, TenantID (opcional), OcurredAt.

---

### 4.3 lotes:AgregarLote

**Permiso requerido:** CREAR_LOTE
**Ubicación:** `internal/application/usecases/agregarlote/`

#### 4.3.1 command.go

Campos de entrada:
- FincaID (string): UUID válido, requerido. Identifica la finca a la que pertenecerá el lote.
- Nombre (string): requerido. Entre 3 y 150 caracteres.
- Area (float64): requerido. Mayor a 0. Máximo 99999.99. Representa hectáreas con 2 decimales de precisión.
- Descripcion (string): opcional. Máximo 1000 caracteres si se provee.

Validaciones:
- FincaID UUID no vacío → ErrValidacion.
- Nombre vacío o menor a 3 caracteres → ErrValidacion.
- Nombre mayor a 150 caracteres → ErrValidacion.
- Area menor o igual a 0 → ErrValidacion.
- Area mayor a 99999.99 → ErrValidacion.
- Descripcion mayor a 1000 caracteres → ErrValidacion.

#### 4.3.2 usecase.go

Dependencias inyectadas: FincaRepositorio, LoteRepositorio, GeneradorID, EventPublisher.

Flujo:
1. Validar que CREAR_LOTE está en AuthContext.Permisos. Si no → ErrForbidden.
2. Cargar finca con FincaRepositorio.ObtenerPorID usando FincaID del comando.
3. Si la finca no existe → ErrNotFound.
4. Validar tenencia de la finca: comparar TenantID de la finca vs AuthContext.TenantID. Si no coincide → ErrNotFound.
5. Verificar que la finca esté en estado ACTIVA. Si está PENDIENTE_ELIMINACION → ErrConflictoEstado ("No se pueden agregar lotes a una finca en estado PENDIENTE_ELIMINACION").
6. Generar ID con GeneradorID.
7. Construir entidad Lote vía el constructor del dominio (NuevoLote). El TenantID del lote se asigna desde AuthContext.TenantID. Si la finca tiene TenantID, el lote debe heredarlo.
8. Persistir el lote con LoteRepositorio.Crear.
9. Si la persistencia falla, retornar error. No publicar evento.
10. Si la persistencia es exitosa, publicar evento LoteCreado.

Transaccionalidad: una sola escritura → sin Unit of Work.

#### 4.3.3 salida.go

Campos de respuesta: ID, FincaID, Nombre, Area, Descripcion, Estado ("ACTIVO"), CreatedAt.

#### 4.3.4 Evento RabbitMQ

Routing key: fincas.v1.lote.creado

Campos del evento: EventID, LoteID, FincaID, Nombre, Area, TenantID, OcurredAt.

---

### 4.4 lotes:EliminarLote

**Permiso requerido:** ELIMINAR_LOTE
**Transición de estado:** ACTIVO → ELIMINADO
**Ubicación:** `internal/application/usecases/eliminarlote/`

#### 4.4.1 command.go

Campos de entrada:
- LoteID (string): UUID válido, requerido.

Validaciones:
- LoteID UUID no vacío → ErrValidacion.

#### 4.4.2 usecase.go

Dependencias inyectadas: LoteRepositorio, EventPublisher.

Flujo:
1. Validar que ELIMINAR_LOTE está en AuthContext.Permisos. Si no → ErrForbidden.
2. Cargar lote con LoteRepositorio.ObtenerPorID usando LoteID del comando.
3. Si no existe → ErrNotFound.
4. Validar tenencia: comparar TenantID del lote vs AuthContext.TenantID. Si no coincide → ErrNotFound.
5. Verificar estado actual del lote. Si ya está ELIMINADO → ErrConflictoEstado.
6. Ejecutar cambio de estado en la entidad Lote vía su método de dominio (ACTIVO → ELIMINADO). Si la transición no es válida → ErrConflictoEstado.
7. Persistir el cambio de estado con LoteRepositorio.Actualizar.
8. Si la persistencia falla, retornar error. No publicar evento.
9. Si la persistencia es exitosa, publicar evento LoteEliminado.

Nota: las muestras asociadas al lote no se eliminan en cascada. Permanecen con el loteID como referencia histórica. Si en el futuro se requiere validar dependencias, se puede consultar el repositorio de muestras antes de eliminar.

Transaccionalidad: una sola escritura → sin Unit of Work.

#### 4.4.3 salida.go

Campos de respuesta: ID, Estado ("ELIMINADO"), UpdatedAt.

#### 4.4.4 Evento RabbitMQ

Routing key: fincas.v1.lote.eliminado

Campos del evento: EventID, LoteID, FincaID, EstadoAnterior ("ACTIVO"), EstadoNuevo ("ELIMINADO"), TenantID, OcurredAt.

---

### 4.5 muestras:TomarMuestra

**Permiso requerido:** CREAR_MUESTRA
**Ubicación:** `internal/application/usecases/tomarmuestra/`

#### 4.5.1 command.go

Campos de entrada:
- LoteID (string): UUID válido, requerido. Lote al que pertenece la muestra.
- Latitud (float64): requerido. Rango [-90, 90]. Coordenada GPS.
- Longitud (float64): requerido. Rango [-180, 180]. Coordenada GPS.

Validaciones:
- LoteID UUID no vacío → ErrValidacion.
- Latitud fuera del rango [-90, 90] → ErrValidacion.
- Longitud fuera del rango [-180, 180] → ErrValidacion.

#### 4.5.2 usecase.go

Dependencias inyectadas: LoteRepositorio (de fincas/domain), MuestraRepositorio (de diagnostico/domain), GeneradorID, EventPublisher.

Flujo:
1. Validar que CREAR_MUESTRA está en AuthContext.Permisos. Si no → ErrForbidden.
2. Cargar lote con LoteRepositorio.ObtenerPorID usando LoteID del comando.
3. Si el lote no existe → ErrNotFound.
4. Validar tenencia del lote: comparar TenantID del lote vs AuthContext.TenantID. Si no coincide → ErrNotFound.
5. Verificar que el lote esté en estado ACTIVO. Si está ELIMINADO → ErrConflictoEstado ("No se pueden tomar muestras en un lote eliminado").
6. Construir Ubicacion VO con el constructor del dominio (NewUbicacion). Si las coordenadas son inválidas → ErrValidacion.
7. Generar ID con GeneradorID.
8. Construir entidad Muestra vía el constructor del dominio (NewMuestra). El TenantID se asigna desde el AuthContext.
9. Persistir la muestra con MuestraRepositorio.Crear.
10. Si la persistencia falla, retornar error. No publicar evento.
11. Si la persistencia es exitosa, publicar evento MuestraTomada.

Transaccionalidad: una sola escritura → sin Unit of Work.

#### 4.5.3 salida.go

Campos de respuesta: ID, LoteID, Latitud, Longitud, CreatedAt.

#### 4.5.4 Evento RabbitMQ

Routing key: diagnosticos.v1.muestra.tomada

Campos del evento: EventID, MuestraID, LoteID, Latitud, Longitud, TenantID, OcurredAt.

---

### 4.6 muestras:ListarMuestrasPorLote

**Permiso requerido:** VER_MUESTRAS
**Ubicación:** `internal/application/usecases/listarmuestrasporlote/`

#### 4.6.1 command.go

Campos de entrada:
- LoteID (string): UUID válido, requerido.

Validaciones:
- LoteID UUID no vacío → ErrValidacion.

#### 4.6.2 usecase.go

Dependencias inyectadas: LoteRepositorio, MuestraRepositorio.

Flujo:
1. Validar que VER_MUESTRAS está en AuthContext.Permisos. Si no → ErrForbidden.
2. Cargar lote con LoteRepositorio.ObtenerPorID usando LoteID del comando.
3. Si el lote no existe → ErrNotFound.
4. Validar tenencia del lote: comparar TenantID del lote vs AuthContext.TenantID. Si no coincide → ErrNotFound.
5. Consultar muestras del lote usando MuestraRepositorio.Buscar. Filtrar por loteID y tenantID.
6. Mapear resultados a estructuras de salida.
7. Retornar lista.

Solo lectura → sin evento, sin Unit of Work.

#### 4.6.3 salida.go

Arreglo de estructuras, cada una con: ID, LoteID, Latitud, Longitud, CreatedAt.

Las coordenadas se incluyen para renderizar mapa. No se exponen IDs de recursos internos que permitan enumeración.

---

### 4.7 muestras:SolicitarDiagnosticoManual

**Permiso requerido:** SOLICITAR_DIAGNOSTICO
**Ubicación:** `internal/application/usecases/solicitardiagnosticomanual/`

#### 4.7.1 command.go

Campos de entrada:
- MuestraID (string): UUID válido, requerido. Muestra ya registrada a la que se asociará el diagnóstico.
- ImageURL (string): requerido. URL de la imagen ya subida al almacenamiento. Debe ser HTTPS.

Validaciones:
- MuestraID UUID no vacío → ErrValidacion.
- ImageURL vacío → ErrValidacion.
- ImageURL no comienza con https:// → ErrValidacion.
- Opcional: validar extensión del archivo (.jpg, .jpeg, .png) como recomendación.

#### 4.7.2 usecase.go

Dependencias inyectadas: MuestraRepositorio, GeneradorID, EventPublisher.

Flujo:
1. Validar que SOLICITAR_DIAGNOSTICO está en AuthContext.Permisos. Si no → ErrForbidden.
2. Cargar la muestra con MuestraRepositorio.ObtenerPorID usando MuestraID del comando.
3. Si la muestra no existe → ErrNotFound.
4. Validar tenencia: comparar TenantID de la muestra vs AuthContext.TenantID. Si no coincide → ErrNotFound.
5. (Opcional) Validación de negocio: si la muestra ya tiene un diagnóstico en estado PENDIENTE, se permite igualmente. El pipeline decide si reemplaza o acumula.
6. Generar SolicitudID con GeneradorID para correlación futura.
7. Publicar evento SolicitudDiagnosticoManual en RabbitMQ.
8. Si la publicación falla, retornar error al usuario (la solicitud no entró al pipeline).
9. Si la publicación es exitosa, retornar confirmación.

Este caso de uso NO persiste nada en base de datos. No crea un Diagnóstico ni una Solicitud en BD. El Diagnóstico se creará cuando RegistrarInferencia consuma el resultado de YOLO, independientemente del origen (edge o manual). Esto mantiene el pipeline unificado y a YOLO desacoplado.

Transaccionalidad: no aplica (no hay escritura en BD). Sin Unit of Work.

#### 4.7.3 salida.go

Campos de respuesta: SolicitudID (UUID generado para tracking), MuestraID, SolicitadoEn (timestamp de publicación del evento).

No retorna datos del diagnóstico porque aún no existe. El frontend debe suscribirse al evento de resultado o implementar polling sobre el listado de diagnósticos.

#### 4.7.4 Evento RabbitMQ

Routing key: diagnosticos.v1.solicitud.diagnostico.manual

Campos del evento: SolicitudID, MuestraID, ImageURL, TenantID, UsuarioID, OcurredAt.

Nota: este evento lo consume el preprocesador del pipeline, no RegistrarInferencia. El preprocesador descarga la imagen de ImageURL, la prepara y la envía a YOLO. YOLO procesa y publica el resultado en la misma cola que consume RegistrarInferencia, independientemente del origen. SolicitudID permite correlación si el pipeline lo propaga.

---

### 4.8 diagnosticos:RegistrarInferencia

**Sin permiso de usuario:** lo dispara el consumer de RabbitMQ al recibir el resultado de YOLO, venga del flujo edge o del flujo on-demand.
**Ubicación:** `internal/application/usecases/registrarinferencia/`

#### 4.8.1 command.go

Campos de entrada:
- MuestraID (string): UUID válido, requerido. Muestra a la que se asocia el diagnóstico.
- ImageURL (string): requerido. URL de la imagen procesada por YOLO.
- TieneClorosis (booleano): requerido. Resultado de la inferencia.
- Confianza (float64): requerido. Valor entre 0.0 y 1.0 inclusive.
- ProcesadoAt (timestamp): requerido. Momento en que YOLO procesó la imagen. No puede ser zero time ni fecha futura.

Validaciones:
- MuestraID UUID no vacío → ErrValidacion.
- ImageURL vacío → ErrValidacion.
- Confianza fuera del rango [0.0, 1.0] → ErrValidacion.
- ProcesadoAt zero o futuro → ErrValidacion.

#### 4.8.2 usecase.go

Dependencias inyectadas: MuestraRepositorio, DiagnosticoRepositorio, GeneradorID, EventPublisher.

Flujo:
1. Sin verificación de permiso. Es una operación interna del sistema disparada por el consumer de RabbitMQ.
2. Cargar la muestra con MuestraRepositorio.ObtenerPorID usando MuestraID del comando.
3. Si la muestra no existe → ErrNotFound. (Loggear el error pero no colapsar el consumer).
4. Construir ResultadoInferencia VO con el constructor del dominio. Si falla validación → ErrValidacion.
5. Generar ID con GeneradorID.
6. Generar nombre autogenerado para el diagnóstico con formato "INF-{YYYYMMDD}-{random}".
7. Construir entidad Diagnostico vía el constructor del dominio. El TenantID se extrae de la muestra cargada. El estado se asigna como PENDIENTE por el constructor.
8. Persistir el diagnóstico con DiagnosticoRepositorio.Crear.
9. Si la persistencia falla, retornar error. No publicar evento.
10. Si la persistencia es exitosa, publicar evento DiagnosticoCreado.

Transaccionalidad: una sola escritura → sin Unit of Work.

#### 4.8.3 salida.go

Campos de respuesta: ID, MuestraID, Nombre, Estado ("PENDIENTE"), TieneClorosis, Confianza, ImageURL, ProcesadoAt, CreatedAt.

#### 4.8.4 Evento RabbitMQ

Routing key: diagnosticos.v1.diagnostico.creado

Campos del evento: EventID, DiagnosticoID, MuestraID, Estado ("PENDIENTE"), TieneClorosis, Confianza, TenantID, OcurredAt.

---

### 4.9 diagnosticos:AceptarDiagnostico

**Permiso requerido:** ACEPTAR_DIAGNOSTICO
**Transición:** PENDIENTE → ACEPTADO
**Ubicación:** `internal/application/usecases/aceptardiagnostico/`

#### 4.9.1 command.go

Campos de entrada:
- DiagnosticoID (string): UUID válido, requerido.

Validaciones:
- DiagnosticoID UUID no vacío → ErrValidacion.

#### 4.9.2 usecase.go

Dependencias inyectadas: DiagnosticoRepositorio, EventPublisher.

Flujo:
1. Validar que ACEPTAR_DIAGNOSTICO está en AuthContext.Permisos. Si no → ErrForbidden.
2. Cargar diagnóstico con DiagnosticoRepositorio.ObtenerPorID usando DiagnosticoID del comando.
3. Si no existe → ErrNotFound.
4. Validar tenencia: comparar TenantID del diagnóstico vs AuthContext.TenantID. Si no coincide → ErrNotFound.
5. Validar estado actual del diagnóstico. Si no es PENDIENTE → ErrConflictoEstado ("El diagnóstico no está pendiente. Estado actual: {estado}").
6. Ejecutar cambio de estado en la entidad Diagnostico vía su método de dominio (PENDIENTE → ACEPTADO). Si la transición no es válida → ErrConflictoEstado.
7. Persistir el cambio de estado con DiagnosticoRepositorio.Actualizar.
8. Si la persistencia falla, retornar error. No publicar evento.
9. Si la persistencia es exitosa, publicar evento DiagnosticoAceptado.

Transaccionalidad: una sola escritura → sin Unit of Work.

#### 4.9.3 salida.go

Campos de respuesta: ID, Estado ("ACEPTADO"), UpdatedAt.

#### 4.9.4 Evento RabbitMQ

Routing key: diagnosticos.v1.diagnostico.aceptado

Campos del evento: EventID, DiagnosticoID, MuestraID, EstadoAnterior ("PENDIENTE"), EstadoNuevo ("ACEPTADO"), UsuarioID, TenantID, OcurredAt.

---

### 4.10 diagnosticos:RechazarDiagnostico

**Permiso requerido:** RECHAZAR_DIAGNOSTICO
**Transición:** PENDIENTE → RECHAZADO. Además marca como candidato a reentrenamiento.
**Ubicación:** `internal/application/usecases/rechazardiagnostico/`

#### 4.10.1 command.go

Campos de entrada:
- DiagnosticoID (string): UUID válido, requerido.
- Motivo (string): opcional. Máximo 500 caracteres. Razón del rechazo.

Validaciones:
- DiagnosticoID UUID no vacío → ErrValidacion.
- Motivo mayor a 500 caracteres → ErrValidacion.

#### 4.10.2 usecase.go

Dependencias inyectadas: DiagnosticoRepositorio, CandidatoReentrenamientoRepositorio, UnitOfWork, EventPublisher.

Flujo:
1. Validar que RECHAZAR_DIAGNOSTICO está en AuthContext.Permisos. Si no → ErrForbidden.
2. Cargar diagnóstico con DiagnosticoRepositorio.ObtenerPorID usando DiagnosticoID del comando.
3. Si no existe → ErrNotFound.
4. Validar tenencia: comparar TenantID del diagnóstico vs AuthContext.TenantID. Si no coincide → ErrNotFound.
5. Validar estado actual. Si no es PENDIENTE → ErrConflictoEstado ("El diagnóstico no está pendiente. Estado actual: {estado}").
6. Ejecutar cambio de estado en la entidad Diagnostico vía su método de dominio (PENDIENTE → RECHAZADO). Si la transición no es válida → ErrConflictoEstado.
7. Ejecutar dentro de Unit of Work (múltiples escrituras atómicas):
   a. Persistir cambio de estado del diagnóstico con DiagnosticoRepositorio.Actualizar.
   b. Crear y persistir un nuevo CandidatoReentrenamiento con: ID del diagnóstico, ResultadoInferencia asociado, Motivo del rechazo y UsuarioID que rechazó.
8. Si la transacción falla → retornar error. No publicar evento.
9. Si la transacción es exitosa → publicar evento DiagnosticoRechazado.

El candidato a reentrenamiento permite que futuros ciclos de entrenamiento del modelo YOLO utilicen estos rechazos como datos de entrenamiento adicionales. Su entidad tiene: id, diagnosticoID, imageURL, tieneClorosis, confianza, motivo, rechazadoPorUsuarioID, createdAt.

Transaccionalidad: múltiples escrituras → requiere Unit of Work.

#### 4.10.3 salida.go

Campos de respuesta: ID, Estado ("RECHAZADO"), Motivo (opcional), UpdatedAt.

#### 4.10.4 Evento RabbitMQ

Routing key: diagnosticos.v1.diagnostico.rechazado

Campos del evento: EventID, DiagnosticoID, MuestraID, EstadoAnterior ("PENDIENTE"), EstadoNuevo ("RECHAZADO"), Motivo (opcional), EsCandidatoRetrain (siempre true), UsuarioID, TenantID, OcurredAt.

---

### 4.11 reportes:GenerarReportePorLote

**Permiso requerido:** GENERAR_REPORTE
**Ubicación:** `internal/application/usecases/generarreporteporlote/`

#### 4.11.1 command.go

Campos de entrada:
- LoteID (string): UUID válido, requerido.

Validaciones:
- LoteID UUID no vacío → ErrValidacion.

#### 4.11.2 usecase.go

Dependencias inyectadas: LoteRepositorio, MuestraRepositorio, DiagnosticoRepositorio.

Flujo:
1. Validar que GENERAR_REPORTE está en AuthContext.Permisos. Si no → ErrForbidden.
2. Cargar lote con LoteRepositorio.ObtenerPorID usando LoteID del comando.
3. Si el lote no existe → ErrNotFound.
4. Validar tenencia: comparar TenantID del lote vs AuthContext.TenantID. Si no coincide → ErrNotFound.
5. Obtener todas las muestras del lote mediante MuestraRepositorio.Buscar filtrando por loteID y tenantID.
6. Para cada muestra, obtener su diagnóstico asociado mediante DiagnosticoRepositorio.Buscar filtrando por muestraID.
7. Calcular zonas afectadas con radio de influencia de 2 metros:
   - Identificar muestras con diagnóstico ACEPTADO y tieneClorosis true.
   - Para cada muestra positiva, definir un círculo de 2m de radio centrado en sus coordenadas GPS.
   - Calcular el área total afectada como la unión de estos círculos (aproximación).
   - Calcular el porcentaje del área del lote afectado en relación al área total del lote.
8. Ensamblar estructura de reporte.
9. Retornar reporte.

Solo lectura → sin evento, sin Unit of Work.

#### 4.11.3 salida.go

Estructura compuesta con las siguientes secciones:

- **Lote**: ID, Nombre, AreaTotal (hectáreas), Estado.
- **Muestras**: arreglo de objetos, cada uno con ID, Latitud, Longitud, DiagnosticoID (opcional, solo si existe diagnóstico), EstadoDiagnostico (opcional), TieneClorosis (opcional), Confianza (opcional).
- **ZonasAfectadas**: arreglo de objetos, cada uno con Latitud, Longitud, RadioMts (siempre 2.0).
- **Metricas**: TotalMuestras (entero), ConClorosis (entero), SinClorosis (entero), Pendientes (entero), AreaAfectadaEstimada (float64, en metros cuadrados), PorcentajeAfectado (float64, entre 0 y 100).
- **GeneradoEn**: timestamp de generación del reporte.

---

## 5. Matriz Resumen

| # | Caso de uso | Permiso requerido | ¿Usa UoW? | Evento que publica | Routing key del evento |
|---|---|---|---|---|---|
| 1 | RegistrarFinca | CREAR_FINCA | No | FincaCreada | fincas.v1.finca.creada |
| 2 | DesactivarFinca | DESACTIVAR_FINCA | No | FincaDesactivada | fincas.v1.finca.desactivada |
| 3 | AgregarLote | CREAR_LOTE | No | LoteCreado | fincas.v1.lote.creado |
| 4 | EliminarLote | ELIMINAR_LOTE | No | LoteEliminado | fincas.v1.lote.eliminado |
| 5 | TomarMuestra | CREAR_MUESTRA | No | MuestraTomada | diagnosticos.v1.muestra.tomada |
| 6 | ListarMuestrasPorLote | VER_MUESTRAS | — | — | — |
| 7 | SolicitarDiagnosticoManual | SOLICITAR_DIAGNOSTICO | No | SolicitudDiagnosticoManual | diagnosticos.v1.solicitud.diagnostico.manual |
| 8 | RegistrarInferencia | *(ninguno)* | No | DiagnosticoCreado | diagnosticos.v1.diagnostico.creado |
| 9 | AceptarDiagnostico | ACEPTAR_DIAGNOSTICO | No | DiagnosticoAceptado | diagnosticos.v1.diagnostico.aceptado |
| 10 | RechazarDiagnostico | RECHAZAR_DIAGNOSTICO | Sí | DiagnosticoRechazado | diagnosticos.v1.diagnostico.rechazado |
| 11 | GenerarReportePorLote | GENERAR_REPORTE | — | — | — |

---

## 6. Mapa de Errores

| Error de aplicación | Causa | Código HTTP |
|---|---|---|
| ErrForbidden | El permiso requerido no está presente en AuthContext.Permisos | 403 |
| ErrNotFound | Recurso no existe en BD, o existe pero su TenantID no coincide con AuthContext.TenantID | 404 |
| ErrConflictoEstado | Se intenta una transición de estado no permitida (ej. aceptar diagnóstico que ya está ACEPTADO, agregar lote a finca PENDIENTE_ELIMINACION, tomar muestra en lote ELIMINADO) | 409 |
| ErrValidacion | Campo inválido en el comando (coordenadas fuera de rango, nombre vacío, área negativa o cero, URL mal formada, confianza fuera de rango) | 400 |
| ErrFincaConLotes | Se intenta desactivar una finca que tiene lotes activos sin enviar Confirmar=true | 409 |

---

## 7. Interface Unit of Work

Solo la necesita el caso de uso RechazarDiagnostico (caso 10). Su interfaz expone:

- Un método Transaccional que recibe una función callback. Dentro del callback se ejecutan todas las operaciones de escritura.
- Acceso a los repositorios necesarios dentro de la transacción: DiagnosticoRepositorio y CandidatoReentrenamientoRepositorio.
- La transacción inicia al entrar al callback, se confirma si el callback retorna nil, se hace rollback si retorna error.
- La publicación del evento ocurre solo si la transacción fue exitosa.

La implementación concreta vive en la capa de infraestructura (GORM + PostgreSQL). El resto de casos de uso con escritura escriben una sola entidad y no requieren Unit of Work.

---

## 8. Flujo de Eventos — Pipeline de Diagnóstico

### Flujo edge (automático)

```
Sensor MQTT → Preprocesador → YOLO → [RabbitMQ] → RegistrarInferencia → Diagnostico creado
```

### Flujo on-demand (manual)

```
Usuario (móvil) → TomarMuestra → [sube foto a storage] → SolicitarDiagnosticoManual → [publica evento]
       ↓
Preprocesador consume evento → descarga imagen → envía a YOLO → [YOLO publica resultado]
       ↓
RegistrarInferencia consume resultado (misma cola, mismo formato) → Diagnostico creado
```

Ambos flujos convergen en RegistrarInferencia. YOLO no conoce el origen de la imagen. El SolicitudID en el evento de SolicitarDiagnosticoManual permite correlacionar la solicitud original con el diagnóstico final si el pipeline propaga el ID como metadata.

---

## 9. Nuevas Entidades de Dominio Requeridas

### CandidatoReentrenamiento

Nueva entidad en `diagnostico/domain/` para persistir diagnósticos rechazados que alimentarán futuros ciclos de reentrenamiento del modelo YOLO.

Atributos: id, diagnosticoID, imageURL, tieneClorosis, confianza, motivo, rechazadoPorUsuarioID, createdAt.

Requiere un nuevo repositorio: CandidatoReentrenamientoRepositorio con métodos Crear y ListarPendientes.

---

## 10. Relación con Otras Especificaciones

- `spec-arquitectura-fincas.md` — Arquitectura general del microservicio.
- `spec-fincas-domain.md` — Entidades de dominio Finca y Lote, reglas de tenencia, errores de dominio.
- `spec-infrastructure.md` — Implementación de repositorios, Unit of Work, Specification pattern.
- `spec-presentation-backend.md` (futuro) — Handlers HTTP, DTOs, mappers, facades.
