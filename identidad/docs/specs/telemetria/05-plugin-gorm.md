# Mini-Spec 5: El Enchufe de Base de Datos — GORM Plugin

> **Propósito**: Capturar métricas de TODAS las operaciones de BD (SELECT, INSERT, UPDATE, DELETE, transacciones) sin modificar ningún repositorio existente.

---

## 1. ¿Dónde se enchufa?

Se implementa como un `gorm.Plugin` en `internal/infrastructure/telemetry/gormplugin/`. Se registra UNA sola vez en la instancia global de `*gorm.DB` dentro de `registry.NewRegistry()` con `db.Use(plugin)`.

**El "enchufe" es el sistema de hooks de GORM.** Cualquier repositorio que use esa instancia de `*gorm.DB` queda automáticamente interceptado. No se necesita cambiar ningún repositorio.

---

## 2. Hooks interceptados

El plugin se suscribe a estos callbacks de GORM:

| Hook GORM | Operación capturada |
|-----------|-------------------|
| `gorm:after_query` | SELECT |
| `gorm:after_create` | INSERT |
| `gorm:after_update` | UPDATE |
| `gorm:after_delete` | DELETE |
| `gorm:after_begin_transaction` | Inicio de transacción |
| `gorm:after_commit` | Commit exitoso |
| `gorm:after_rollback` | Rollback |

---

## 3. Flujo de cada callback

1. GORM ejecuta la operación real.
2. El callback se dispara (después de la operación).
3. Se calcula `duration_ms`.
4. Se extrae `table` del `Statement.Table`.
5. Se extrae `rows_affected` del `Statement.RowsAffected`.
6. Si hay error: se mapea el error GORM a código SQL state (PostgreSQL).
7. Se calcula `query_hash` = SHA-256 de la consulta normalizada sin valores literales (NUNCA se registra la query raw).
8. Se lee `trace_id` del `context.Context` asociado a la sesión GORM.
9. Se aplican reglas de nivel:
   - Duración < 200ms → INFO
   - Duración 200ms–1s → WARN
   - Duración > 1s → ERROR
   - Filas > 1000 → WARN
   - Transacción > 5s → WARN
   - Rollback → ERROR
10. Se construye payload `LogPayload` con `log_type = "BD"`.
11. Se envía a `BufferWriter.Write()`.

---

## 4. Pureza arquitectónica

- Los repositorios de dominio (`UsuarioRepositorio`, `SesionRepositorio`) **no tienen código de logging**.
- No necesitan inyectar `BufferWriter`.
- No saben que están siendo cronometrados.
- **Sus tests unitarios no se ven afectados.**

---

## 5. El `query_hash`

Se calcula con SHA-256 de `Statement.SQL` después de reemplazar los valores literales por placeholders (`?`). Esto permite:
- Agrupar queries semánticamente iguales en Grafana.
- Detectar queries lentas sin exponer datos sensibles.
- Comparar rendimiento entre ejecuciones de la misma query.

---

## 6. Archivos a crear

```
internal/infrastructure/telemetry/gormplugin/
├── plugin.go          ← implementación de gorm.Plugin
├── plugin_test.go     ← tests con SQLite en memoria
└── config.go          ← umbrales (slow_query_threshold_ms, etc.)
```
