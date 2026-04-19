# Test de Repositorio de Usuario

Script standalone que testea directamente el repositorio de usuarios contra la base de datos PostgreSQL, sin HTTP.

## Requisitos

- PostgreSQL corriendo
- Variables de entorno en `.env`:
  - `DB_HOST=localhost`
  - `DB_PORT=5432`
  - `DB_USER=identidad_user`
  - `DB_PASSWORD=identidad_pass_dev`
  - `DB_NAME=identidad_db`
  - `DB_SSLMODE=disable`

## Uso

```bash
# Build
make build-test

# Run
make test

# O ejecutar directamente
./bin/test
```

## Arquitectura de Prueba: Flujo Temporal de Eventos

El test está estructurado en **5 fases** que demuestran la lógica correcta de eventos del dominio:

### FASE 1: CREACIÓN DE USUARIOS
- Crear N usuarios **SIN manipularlos** durante creación
- Cada usuario comienza en estado `NO_VERIFICADO`
- Se persisten en BD
- Los eventos quedan en la cola del usuario creado

```
[1] Crear Juan → estado NO_VERIFICADO ✓
[2] Crear María → estado NO_VERIFICADO ✓
[3] Crear Carlos → estado NO_VERIFICADO ✓
```

### FASE 2: EVENTOS INICIALES DE CREACIÓN
- Capturar eventos post-creación de cada usuario
- Los eventos están vacíos (0) porque `NuevoUsuario` no emite eventos
- Esto es **correcto en DDD**: eventos solo se emiten por cambios de estado

```
Juan eventos: Total: 0
María eventos: Total: 0
Carlos eventos: Total: 0
```

### FASE 3: MANIPULACIÓN DE ESTADOS
- Operar sobre usuarios ya creados
- Activar, bloquear, etc.
- Cada operación emite eventos

```
[4] Activar Juan → UsuarioActivado ✓
[5] Activar María → UsuarioActivado ✓
[6] Activar Carlos → UsuarioActivado ✓
[7] Bloquear Carlos → UsuarioBloqueado ✓
```

### FASE 4: EVENTOS DESPUÉS DE CAMBIOS
- Mostrar cola de eventos para cada operación de cambio de estado
- Cada método de estado emite **2 eventos**:
  1. `EstadoUsuarioCambiado` (emitido por `CambiarEstado()`)
  2. Evento específico (emitido por método como `Activar()` o `Bloquear()`)

```
Juan eventos (post-activación): 2
   1. EstadoUsuarioCambiado
   2. UsuarioActivado

Carlos eventos (post-bloqueo): 2
   1. EstadoUsuarioCambiado
   2. UsuarioBloqueado
```

### FASE 5: OPERACIONES BD
- Persistir cambios en base de datos
- Listar usuarios en diferentes estados
- Eliminar usuario y verificar final

```
[8] Persistir cambios en BD ✓
[9] Listar todos los usuarios: 3
[10] Listar solo ACTIVOS: 2
[11] Eliminar Carlos ✓
[12] Listar finales (todos): 2
```

## Qué Testea

1. **Creación**: 3 usuarios en estado NO_VERIFICADO (sin eventos)
2. **Eventos iniciales**: Verificar cola vacía post-creación
3. **Cambios de estado**: Activar y bloquear usuarios
4. **Eventos post-cambios**: Capturar y verificar eventos por operación
5. **Persistencia**: Actualizar estados en BD
6. **Listado**: Todos, filtrados por estado
7. **Eliminación**: Borrar usuario y verificar final

## Output de Eventos

El test muestra claramente la naturaleza temporal de los eventos:

```
📊 RESUMEN:
   • Usuarios creados: 3
   • Usuarios eliminados: 1
   • Usuarios finales: 2
   • Eventos de creación: 0 (Juan), 0 (María), 0 (Carlos)
   • Eventos de activación: 2 (Juan), 2 (María), 2 (Carlos)
   • Eventos de bloqueo: 2 (Carlos)
```

## Puntos Clave de la Arquitectura

### ✓ Correcto
- Crear primero, capturar eventos
- LUEGO modificar estados
- Cada `PullEventos()` limpia la cola
- Eventos reflejan cambios reales en el dominio

### ✗ Evitar
- NO hacer `.Activar()` durante creación
- NO manipular durante construcción
- NO asumir eventos en `NuevoUsuario()`

## Flow

1. Carga config desde `.env`
2. Conecta a PostgreSQL
3. Limpia tabla `usuarios` (DROP si existe)
4. Crea tabla con migraciones
5. Ejecuta 5 fases de prueba (creación → eventos → manipulación → eventos post-cambios → BD)
6. Muestra resumen de eventos
7. Cierra conexión automáticamente

