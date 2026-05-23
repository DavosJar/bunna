# Plan de Casos de Uso — Capacidades del Sistema

Basado en [`spec-capacidades-del-sistema.md`](./autorizacion/spec-capacidades-del-sistema.md) v4.1.

---

## Estructura por caso de uso

```
internal/{modulo}/application/usecases/{nombre_caso_uso}/
├── command.go      # Comando{campos} + validación
├── response.go     # Respuesta{campos}
└── usecase.go      # CasoDeUso{deps} + método Ejecutar
```

---

## I. Gestión de Usuarios

### 1. `createuser` — Crear Usuario

Protegido por permiso `Crear_Usuario`.

- **Comando**: correo, nombre, apellido, password (plano), CreatedBy (quién crea), EjecutorID (autenticado)
- **Respuesta**: ID, correo, nombre, apellido, activo, creado en
- **Dependencias**: repositorio de usuarios, servicio de autorización
- **Flujo**: validar email y campos obligatorios → autorizar permiso → hashear password → crear entidad Usuario + Credenciales (transaccional) → armar respuesta

### 2. `listusers` — Consultar Usuarios

Protegido por permiso `Consultar_Usuarios`.

- **Comando**: filtros opcionales por estado / correo / tenant, paginación, EjecutorID
- **Respuesta**: lista de usuarios con ID, correo, nombre, apellido, estado, fecha creación + total, página, tamaño
- **Dependencias**: repositorio de usuarios, servicio de autorización
- **Flujo**: autorizar permiso → aplicar filtros de seguridad → listar con paginación → armar respuesta

### 3. `updateuser` — Modificar Usuario

Protegido por permiso `Modificar_Usuario`.

- **Comando**: UsuarioID, nombre (opcional), apellido (opcional), teléfono (opcional), EjecutorID
- **Respuesta**: ID, correo, nombre, apellido, modificado en
- **Dependencias**: repositorio de usuarios, servicio de autorización
- **Flujo**: autorizar permiso → obtener usuario → modificar campos del agregado → persistir → armar respuesta

### 4. `deleteuser` — Dar de Baja Usuario

Protegido por permiso `Dar_De_Baja_Usuario`.

- **Comando**: UsuarioID, motivo (opcional), EjecutorID
- **Respuesta**: UsuarioID, estado (`PENDIENTE_DE_ELIMINACION`), fecha baja
- **Dependencias**: repositorio de usuarios, servicio de autorización
- **Flujo**: autorizar permiso → obtener usuario → marcar como pendiente de eliminación → persistir → armar respuesta

### 5. `expeluser` — Expulsar Usuario

Protegido por permiso `Expulsar_Usuario` (nueva constante: `identidad:usuario:expulsar`).

- **Comando**: UsuarioID, EjecutorID
- **Respuesta**: UsuarioID, estado final, cantidad de sesiones revocadas, fecha expulsión
- **Dependencias**: repositorio de usuarios, repositorio de sesiones (o UnitOfWork), servicio de autorización
- **Flujo**: autorizar permiso → transacción: bloquear usuario + invalidar todas sus sesiones activas → armar respuesta

---

## II. Gestión de Credenciales y Accesos

### 6. `viewcredentials` — Consultar Credenciales

Protegido por permiso `Consultar_Credenciales` (nueva constante: `identidad:credenciales:consultar`).

- **Comando**: UsuarioID, EjecutorID
- **Respuesta**: UsuarioID, activo, correo verificado, intentos fallidos, bloqueado hasta
- **Dependencias**: repositorio de credenciales, servicio de autorización
- **Flujo**: autorizar permiso → obtener credenciales → armar respuesta con estado de seguridad

### 7. `resetpassword` — Resetear Contraseña (admin)

Protegido por permiso `Resetear_Contrasena`.

- **Comando**: UsuarioID, nueva password, EjecutorID
- **Respuesta**: UsuarioID, modificado en
- **Dependencias**: UnitOfWork (credenciales + sesiones), servicio de autorización
- **Flujo**: autorizar permiso → transacción: obtener credenciales → hashear nueva password → actualizar credenciales → resetear intentos fallidos → invalidar sesiones activas del usuario → armar respuesta

### 8. `unlockaccount` — Desbloquear Cuenta

Protegido por permiso `Desbloquear_Cuenta` (nueva constante: `identidad:credenciales:desbloquear`).

- **Comando**: UsuarioID, EjecutorID
- **Respuesta**: UsuarioID, fecha desbloqueo
- **Dependencias**: repositorio de credenciales, servicio de autorización
- **Flujo**: autorizar permiso → obtener credenciales → limpiar bloqueado hasta e intentos fallidos → persistir → armar respuesta

---

## III. Gestión de Roles y Permisos

### 9. `listroles` — Consultar Roles

Protegido por permiso `Consultar_Roles`.

- **Comando**: filtros opcionales, paginación, EjecutorID
- **Respuesta**: lista de roles con sus permisos, total, página
- **Dependencias**: repositorio de roles, servicio de autorización
- **Flujo**: autorizar permiso → listar roles del sistema + personalizados → armar respuesta

### 10. `createrole` — Crear Rol

Protegido por permiso `Crear_Rol` (nueva constante: `identidad:rol:crear`).

- **Comando**: nombre, descripción, permisos iniciales (opcional), EjecutorID
- **Respuesta**: ID, nombre, descripción, es sistema (false), creado en
- **Dependencias**: repositorio de roles, servicio de autorización
- **Flujo**: autorizar permiso → validar nombre único → crear rol personalizado → asignar permisos iniciales → armar respuesta

### 11. `updaterole` — Modificar Rol

Protegido por permiso `Modificar_Rol` (nueva constante: `identidad:rol:modificar`).

- **Comando**: RolID, nombre (opcional), descripción (opcional), EjecutorID
- **Respuesta**: ID, nombre, descripción, modificado en
- **Dependencias**: repositorio de roles, servicio de autorización
- **Flujo**: autorizar permiso → validar que no sea rol de sistema → actualizar campos → armar respuesta

### 12. `deleterole` — Eliminar Rol

Protegido por permiso `Eliminar_Rol` (nueva constante: `identidad:rol:eliminar`).

- **Comando**: RolID, EjecutorID
- **Respuesta**: RolID, eliminado en
- **Dependencias**: repositorio de roles, repositorio de usuario-rol, servicio de autorización
- **Flujo**: autorizar permiso → validar que no sea rol de sistema → reasignar o migrar usuarios con este rol → eliminar rol → armar respuesta

### 13. `assignrole` — Asignar Rol

Protegido por permiso `Asignar_Rol`.

- **Comando**: UsuarioID, RolID, TenantID (opcional, vacío = rol global), EjecutorID
- **Respuesta**: UsuarioID, RolID, TenantID, asignado en
- **Dependencias**: repositorio de usuario-rol (o usuario-tenant-rol), servicio de autorización
- **Flujo**: autorizar permiso → validar que rol exista → validar consistencia tenant → asignar → armar respuesta

### 14. `revokerole` — Revocar Rol

Protegido por permiso `Revocar_Rol`.

- **Comando**: UsuarioID, RolID, TenantID (opcional), EjecutorID
- **Respuesta**: UsuarioID, RolID, TenantID, revocado en
- **Dependencias**: repositorio de usuario-rol, servicio de autorización
- **Flujo**: autorizar permiso → validar que no sea el último `SYS_ADMIN` → revocar → armar respuesta

### 15. `assignpermissiontorole` — Asignar Permiso a Rol

Protegido por permiso `Asignar_Permiso_A_Rol` (nueva constante: `identidad:rol:permiso:asignar`).

- **Comando**: RolID, código de permiso, EjecutorID
- **Respuesta**: RolID, código de permiso, asignado en
- **Dependencias**: repositorio de rol-permiso, servicio de autorización
- **Flujo**: autorizar permiso → validar que rol sea personalizado (no de sistema) → validar que permiso exista en catálogo → asignar → armar respuesta

### 16. `revokepermissionfromrole` — Revocar Permiso de Rol

Protegido por permiso `Revocar_Permiso_De_Rol` (nueva constante: `identidad:rol:permiso:revocar`).

- **Comando**: RolID, código de permiso, EjecutorID
- **Respuesta**: RolID, código de permiso, revocado en
- **Dependencias**: repositorio de rol-permiso, servicio de autorización
- **Flujo**: autorizar permiso → validar que rol sea personalizado → revocar → armar respuesta

---

## IV. Control de Sesiones y Seguridad

### 17. `listsessions` — Consultar Sesiones

Protegido por permiso `Consultar_Sesiones` (nueva constante: `identidad:sesion:consultar`).

- **Comando**: UsuarioID (opcional, filtro), paginación, EjecutorID
- **Respuesta**: lista de sesiones activas con IP, última actividad, dispositivo
- **Dependencias**: repositorio de sesiones, servicio de autorización
- **Flujo**: autorizar permiso → listar sesiones activas → armar respuesta

### 18. `terminatesession` — Forzar Cierre de Sesión

Protegido por permiso `Forzar_Cierre_Sesion` (nueva constante: `identidad:sesion:forzar_cierre`).

- **Comando**: SesionID, UsuarioID (propietario), EjecutorID
- **Respuesta**: SesionID, estado (`REVOCADA`), fecha revocación
- **Dependencias**: repositorio de sesiones, servicio de autorización
- **Flujo**: autorizar permiso → obtener sesión → revocar → persistir → armar respuesta

### 19. `updatetenant` — Configurar Tenant

Protegido por permiso `Configurar_Tenant` (nueva constante: `identidad:tenant:configurar`).

- **Comando**: TenantID, configuración a modificar (campos por definir según entidad Tenant), EjecutorID
- **Respuesta**: TenantID, campos modificados, fecha modificación
- **Dependencias**: repositorio de tenants, servicio de autorización
- **Flujo**: autorizar permiso → validar que tenant exista → actualizar configuración → armar respuesta
- **Nota**: el dominio actual de `Tenant` solo tiene nombre/slug/activo; se requiere extender la entidad con campos de configuración.

### 20. `listblockedips` — Consultar IPs Bloqueadas

Protegido por permiso `Consultar_IPs_Bloqueadas` (nueva constante: `identidad:ip:consultar`).

- **Comando**: paginación, EjecutorID
- **Respuesta**: lista de IPs bloqueadas con fecha bloqueo, contador de intentos
- **Dependencias**: repositorio de intento por IP, servicio de autorización
- **Flujo**: autorizar permiso → listar IPs bloqueadas vigentes → armar respuesta

### 21. `unblockip` — Desbloquear IP

Protegido por permiso `Desbloquear_IP` (nueva constante: `identidad:ip:desbloquear`).

- **Comando**: dirección IP, EjecutorID
- **Respuesta**: IP, fecha desbloqueo
- **Dependencias**: repositorio de intento por IP, servicio de autorización
- **Flujo**: autorizar permiso → obtener intento por IP → limpiar bloqueo → persistir → armar respuesta

---

## V. Autogestión (Implícitos — sin permiso de rol)

Estos casos de uso no requieren permiso del RBAC. Se autorizan por coincidencia de identidad (`EjecutorID == propietario`) o son rutas públicas.

### 22. `viewmyprofile` — Ver Mi Perfil

- **Comando**: EjecutorID
- **Respuesta**: ID, correo, nombre, apellido, teléfono, estado, fecha creación
- **Dependencias**: repositorio de usuarios
- **Flujo**: obtener usuario por ID del ejecutor → armar respuesta (solo datos propios)

### 23. `updatemyprofile` — Modificar Mi Perfil

- **Comando**: EjecutorID, nombre (opcional), apellido (opcional), teléfono (opcional)
- **Respuesta**: ID, correo, nombre, apellido, modificado en
- **Dependencias**: repositorio de usuarios
- **Flujo**: obtener usuario → modificar campos permitidos (no email) → persistir → armar respuesta
- **Nota**: el email no es modificable por autogestión; cambiaría la identidad del usuario y requeriría nueva verificación.

### 24. `changemypassword` — Cambiar Mi Contraseña

- **Comando**: EjecutorID, password actual, nueva password
- **Respuesta**: EjecutorID, modificado en
- **Dependencias**: repositorio de credenciales, servicio de encriptación
- **Flujo**: obtener credenciales → verificar password actual contra hash → hashear nueva password → actualizar → armar respuesta
- **Opcional**: invalidar otras sesiones (excepto la actual) como medida de seguridad.

---

## VI. Casos de Autogestión y Públicos Existentes (reestructurados)

### 25. `register` — Registrarse (público)

- Sin permiso (ruta pública)
- Comando: correo, password, nombre, apellido, teléfono
- Respuesta: UsuarioID, correo, estado (`NO_VERIFICADO`), creado en
- Dependencias: UnitOfWork (repositorio de usuarios + credenciales, encriptación, generador ID)
- Flujo: validar email y datos obligatorios → UnitOfWork: generar ID → crear usuario → hashear password → crear credenciales → armar respuesta

### 26. `login` — Iniciar Sesión (público)

- Sin permiso (ruta pública)
- Comando: email, password, IP origen
- Respuesta: access token, refresh token, expiración access, expiración refresh, UsuarioID, SesionID
- Dependencias: UnitOfWork de sesiones, servicio de bloqueo IP, rate limiter
- Flujo: validar email y password → rate limiting preventivo → verificar bloqueo por IP → UnitOfWork: resolver email → validar credenciales (contra hash) → verificar bloqueo cuenta → generar tokens → crear sesión → resetear intentos fallidos → armar respuesta. Si password incorrecto: registrar intento fallido por IP.

### 27. `refresh` — Renovar Sesión

- Sin permiso (usa refresh token válido)
- Comando: refresh token
- Respuesta: nuevo access token, nuevo refresh token, expiraciones, SesionID, UsuarioID
- Dependencias: UnitOfWork de sesiones, configuración (máximo refrescos, timeout absoluto)
- Flujo: validar token no vacío → validar JWT → hashear token → UnitOfWork: buscar sesión por hash → si no existe: invalidar todas las sesiones del usuario (posible robo) → verificar estado sesión → verificar timeout absoluto → verificar expiración refresh token → verificar límite de refrescos → rotar tokens → actualizar sesión → armar respuesta

### 28. `logout` — Cerrar Mi Sesión

- Self-service (autorización por token actual)
- **Dos comandos**: cerrar sesión específica (SesionID + UsuarioID) o cerrar todas las sesiones (UsuarioID)
- Respuesta: cantidad de sesiones revocadas
- Dependencias: UnitOfWork de sesiones
- Flujo (específica): validar IDs → obtener sesión → verificar pertenencia al usuario → si activa: revocar y persistir → armar respuesta. Si ya revocada/expirada: no-op.
- Flujo (todas): listar sesiones activas del usuario → revocar cada una → persistir → armar respuesta

### 29. `verifyemail` — Verificar Correo (público + token)

- Sin permiso (usa token de verificación)
- **Tres operaciones**: solicitar verificación, confirmar verificación con token, reenviar verificación
- Comandos: solicitar (UsuarioID), confirmar (token), reenviar (UsuarioID)
- Respuestas: mensaje de confirmación
- Dependencias: repositorio de verificación, servicio de email, generador ID, configuración (expiración token, máximo reenvíos)
- Flujo solicitar: generar token único → crear prueba de verificación con hash → persistir → enviar email con token (asíncrono, best-effort)
- Flujo confirmar: hashear token recibido → buscar por hash → verificar expiración → marcar correo como verificado → limpiar prueba
- Flujo reenviar: verificar límite de reenvíos → reutilizar solicitar

### 30. `forgotpassword` — Recuperar Contraseña (público + email)

- Sin permiso (ruta pública + enlace por email)
- **Tres operaciones**: solicitar recuperación, validar token, confirmar restablecimiento
- Comandos: solicitar (email + IP origen), validar token (token), confirmar (token + nueva password)
- Respuestas: solicitar (mensaje genérico), validar token (UsuarioID + válido), confirmar (mensaje)
- Dependencias: repositorio de tokens de recuperación, repositorio de usuarios, repositorio de sesiones, repositorio de credenciales, servicio de encriptación, servicio de email, generador ID, configuración
- Flujo solicitar: validar email → buscar usuario (silencioso si no existe) → generar token → persistir → enviar email (asíncrono) → respuesta genérica siempre
- Flujo validar token: hashear token → buscar por hash → verificar vigencia → responder
- Flujo confirmar: validar nueva password (longitud mínima) → validar token → hashear nueva password → actualizar credenciales → marcar token como usado → invalidar todas las sesiones del usuario → respuesta

---

## Nuevas constantes de permiso

Además de las 8 existentes en `permisos.go`, se requieren las siguientes para cubrir todas las capacidades del spec:

| Capacidad | Constante propuesta |
|---|---|
| `Expulsar_Usuario` | `identidad:usuario:expulsar` |
| `Consultar_Credenciales` | `identidad:credenciales:consultar` |
| `Desbloquear_Cuenta` | `identidad:credenciales:desbloquear` |
| `Crear_Rol` | `identidad:rol:crear` |
| `Modificar_Rol` | `identidad:rol:modificar` |
| `Eliminar_Rol` | `identidad:rol:eliminar` |
| `Asignar_Permiso_A_Rol` | `identidad:rol:permiso:asignar` |
| `Revocar_Permiso_De_Rol` | `identidad:rol:permiso:revocar` |
| `Consultar_Sesiones` | `identidad:sesion:consultar` |
| `Forzar_Cierre_Sesion` | `identidad:sesion:forzar_cierre` |
| `Configurar_Tenant` | `identidad:tenant:configurar` |
| `Consultar_IPs_Bloqueadas` | `identidad:ip:consultar` |
| `Desbloquear_IP` | `identidad:ip:desbloquear` |

---

## Resumen de carpetas

```
internal/usuarios/application/usecases/
├── createuser/          Crear Usuario (nuevo)
├── listusers/           Consultar Usuarios (nuevo)
├── updateuser/          Modificar Usuario (nuevo)
├── deleteuser/          Dar de Baja Usuario (nuevo)
├── expeluser/           Expulsar Usuario (nuevo)
├── viewmyprofile/       Ver Mi Perfil (nuevo)
├── updatemyprofile/     Modificar Mi Perfil (nuevo)
└── register/            Registrarse (reestructurado)

internal/seguridad/application/usecases/
├── viewcredentials/     Consultar Credenciales (nuevo)
├── resetpassword/       Resetear Contraseña (nuevo)
├── unlockaccount/       Desbloquear Cuenta (nuevo)
├── listblockedips/      Consultar IPs Bloqueadas (nuevo)
├── unblockip/           Desbloquear IP (nuevo)
└── changemypassword/    Cambiar Mi Contraseña (nuevo)

internal/sesiones/application/usecases/
├── listsessions/        Consultar Sesiones (nuevo)
├── terminatesession/    Forzar Cierre de Sesión (nuevo)
├── login/               Iniciar Sesión (reestructurado)
├── refresh/             Renovar Sesión (reestructurado)
└── logout/              Cerrar Mi Sesión (reestructurado)

internal/rbac/application/usecases/
├── listroles/               Consultar Roles (nuevo)
├── createrole/              Crear Rol (nuevo)
├── updaterole/              Modificar Rol (nuevo)
├── deleterole/              Eliminar Rol (nuevo)
├── assignrole/              Asignar Rol (nuevo)
├── revokerole/              Revocar Rol (nuevo)
├── assignpermissiontorole/  Asignar Permiso a Rol (nuevo)
└── revokepermissionfromrole/ Revocar Permiso de Rol (nuevo)

internal/tenants/application/usecases/
└── updatetenant/        Configurar Tenant (nuevo)

internal/verificacion/application/usecases/
└── verifyemail/         Verificar Correo (reestructurado)

internal/recuperacion/application/usecases/
└── forgotpassword/      Recuperar Contraseña (reestructurado)
```

**Total: 24 casos de uso — 18 nuevos + 6 reestructurados.**
