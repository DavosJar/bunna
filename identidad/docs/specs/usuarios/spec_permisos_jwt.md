# Permisos — Estructura

## JWT (access token)

El servicio `identidad` emite el JWT con estos claims:

```json
{
  "iss": "identidad.bunna",
  "sub": "usuario-id",
  "sid": "sesion-id",
  "typ": "access",
  "tenant_id": "tenant-1",
  "rol": "agricultor",
  "iat": 1748383000,
  "exp": 1748386600
}
```

:warning: **No lleva permisos.** El token es liviano. Los permisos no van en claims.

---

## Cómo se resuelven los permisos

Fincas no parsea el JWT para sacar permisos. Apenas recibe el request:

```
1. Extrae tenant_id y rol del JWT
2. Pregunta a identidad: GET /roles/{rol}/permisos?servicio=fincas
3. Identidad responde: ["CREAR_FINCA", "CREAR_LOTE", "CREAR_MUESTRA", ...]
4. Los mete en AuthContext.Permisos y los usa directo
```

Cero parseo. El string que viaja desde la DB de identidad hasta el `TienePermiso()` es el mismo.

---

## Formato del permiso

| Dónde | Formato | Ejemplo |
|---|---|---|
| DB de identidad | `servicio.modulo.PERMISO` | `fincas.finca.CREAR_FINCA` |
| API response (identidad → fincas) | `PERMISO` | `CREAR_FINCA` |
| Constante en código de fincas | `PERMISO` | `CREAR_FINCA` |
| Comparación en `TienePermiso()` | `PERMISO` | `CREAR_FINCA` |

Identidad usa el namespacing `servicio.modulo` para organizar sus tablas y filtrar lo que le devuelve a cada servicio. Fincas nunca ve los puntos.

---

## Catálogo

```
CREAR_FINCA
DESACTIVAR_FINCA
CREAR_LOTE
ELIMINAR_LOTE
CREAR_MUESTRA
VER_MUESTRAS
SOLICITAR_DIAGNOSTICO
ACEPTAR_DIAGNOSTICO
RECHAZAR_DIAGNOSTICO
GENERAR_REPORTE
```

---

## Output cuando falta un permiso

```json
HTTP 403 Forbidden

{
    "error": "no tiene permisos para esta operación"
}
```
