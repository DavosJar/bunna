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

## Qué Testea

1. **Crear**: 5 usuarios de prueba (3 ACTIVO, 2 NO_VERIFICADO)
2. **Obtener por ID**: Recupera un usuario existente
3. **Listar sin filtros**: Todos los usuarios
4. **Listar con filtro**: Por estado (ACTIVO), por nombre (LIKE)
5. **Listar con ordenación**: ASC por nombre, DESC por apellido
6. **Listar con paginación**: Página 1 y 2 (tamaño 2)
7. **Listar combinado**: Filtro + ordenación + paginación
8. **Actualizar**: Cambiar estado de usuario (ACTIVO → INACTIVO)
9. **Eliminar**: Borrar un usuario y verificar eliminación

## Output

Cada operación imprime mensajes claros:
- `[CREADO]` — Usuario creado exitosamente
- `[RESULTADO]` — Resultado de operación exitosa
- `[ERROR]` — Error durante operación
- `[INFO]` — Información adicional (antes/después)
- `[TEST n]` — Inicio de test específico
- `[SETUP]` — Preparación (limpieza, migración)

## Flow

1. Carga config desde `.env`
2. Conecta a PostgreSQL
3. Limpia tabla `usuarios` (DROP si existe)
4. Crea tabla con migraciones
5. Ejecuta 12 tests CRUD
6. Cierra conexión automáticamente
