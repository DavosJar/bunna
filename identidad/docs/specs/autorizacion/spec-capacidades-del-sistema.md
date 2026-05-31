---
title: Capacidades del Sistema (Permisos en Lenguaje Ubicuo)
version: 4.1
date_created: 2026-05-23
owner: Equipo Identidad
tags: autorizacion, capacidades, lenguaje-ubicuo, permisos
---

# Capacidades del Sistema

## 1. Concepto Core

En este sistema, **Capacidad = Permiso = Caso de Uso**. 
No usamos formatos técnicos complejos. Hablamos estrictamente el **Lenguaje Ubicuo** del dominio.

Si un rol puede hacer algo, ese algo es el nombre exacto del caso de uso. Estas capacidades son **constantes estáticas de dominio** en el código. Nadie puede inventar permisos en runtime; solo pueden asignarse los roles a los usuarios, o construir nuevos roles personalizados agrupando estas constantes.

Existen roles de sistema (`SYS_ADMIN`, `ADMINISTRADOR`) que son inborrables y no se les pueden alterar sus permisos, pero los usuarios con la capacidad adecuada pueden crear roles nuevos.

## 2. Catálogo de Capacidades (Permisos Constantes)

Esta es la lista definitiva de permisos/casos de uso del sistema. Estos son los valores exactos que componen los roles.

### Gestión de Usuarios
- `Consultar_Usuarios` : Permite listar y ver los detalles de los usuarios de la organización.
- `Crear_Usuario` : Permite dar de alta a un usuario nuevo manualmente.
- `Modificar_Usuario` : Permite editar la información básica de un usuario.
- `Dar_De_Baja_Usuario` : Permite marcar a un usuario para eliminación (baja lógica).
- `Expulsar_Usuario` : Permite dar de baja inmediata y revocar todas las sesiones activas de un usuario al instante.
- `Agregar_A_Equipo` : Permite asignar un usuario a un equipo de trabajo u organización.

### Gestión de Credenciales y Accesos
- `Consultar_Credenciales` : Permite ver el estado de seguridad de un usuario (bloqueos, intentos fallidos).
- `Resetear_Contrasena` : Permite forzar el reseteo de la contraseña de otro usuario.
- `Desbloquear_Cuenta` : Permite quitar el bloqueo de una cuenta inhabilitada temporalmente por intentos fallidos.

### Gestión de Roles y Permisos
- `Consultar_Roles` : Permite ver la lista de roles y sus permisos, así como ver qué roles tiene un usuario.
- `Crear_Rol` : Permite crear un rol personalizado nuevo en la organización.
- `Modificar_Rol` : Permite cambiar el nombre o descripción de un rol personalizado.
- `Eliminar_Rol` : Permite eliminar un rol personalizado (los de sistema no se pueden borrar).
- `Asignar_Rol` : Permite otorgar un rol (de sistema o personalizado) a un usuario.
- `Revocar_Rol` : Permite quitarle un rol previamente asignado a un usuario.
- `Asignar_Permiso_A_Rol` : Permite agregar una capacidad/permiso del catálogo a un rol personalizado.
- `Revocar_Permiso_De_Rol` : Permite quitar una capacidad/permiso de un rol personalizado.

### Control de Sesiones y Seguridad
- `Consultar_Sesiones` : Permite ver los dispositivos y sesiones desde donde están conectados los usuarios.
- `Forzar_Cierre_Sesion` : Permite desconectar/matar la sesión de un usuario de forma remota.
- `Configurar_Tenant` : Permite cambiar la configuración global y parámetros de la organización (tenant).
- `Consultar_IPs_Bloqueadas` : Permite ver la lista de direcciones IP que el sistema ha bloqueado por ataques.
- `Desbloquear_IP` : Permite quitar una dirección IP de la lista negra.

## 3. Casos de Uso de Autogestión (Implícitos)

Las siguientes acciones son inherentes al propio usuario autenticado. **No son permisos que se asignen a roles**. El sistema las autoriza implícitamente siempre que el usuario opere sobre sus propios datos o sean rutas de acceso público:

- `Registrarse` (Público)
- `Iniciar_Sesion` (Público)
- `Verificar_Correo` (Público)
- `Renovar_Sesion` (Público)
- `Cerrar_Mi_Sesion` (Self-service)
- `Ver_Mi_Perfil` (Self-service)
- `Modificar_Mi_Perfil` (Self-service)
- `Cambiar_Mi_Contrasena` (Self-service)