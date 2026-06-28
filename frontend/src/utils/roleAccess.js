/** Roles de sistema (alineados con backend identidad) */
export const ROLES = {
  SYS_ADMIN: 'sys_admin',
  ADMINISTRADOR: 'administrador',
  AGRONOMO: 'agronomo',
  CAFICULTOR: 'caficultor',
};

const ROL_LABELS = {
  sys_admin: 'Super Admin',
  administrador: 'Administrador',
  agronomo: 'Agrónomo',
  caficultor: 'Caficultor',
};

/** Perfil visual y funcional de cada rol */
export const ROL_PROFILES = {
  sys_admin: {
    label: 'Super Admin',
    accent: 'purple',
    homeRoute: '/admin',
    tagline: 'Control total del sistema',
    capabilities: [
      'Usuarios, roles y permisos en cualquier tenant',
      'Sesiones activas e IPs bloqueadas',
      'Configuración global del sistema',
      'Operación de fincas y análisis YOLO',
    ],
  },
  administrador: {
    label: 'Administrador',
    accent: 'green',
    homeRoute: '/admin',
    tagline: 'Dueño de tu finca / tenant',
    capabilities: [
      'Invitar usuarios y asignar roles (caficultor, agrónomo)',
      'Gestionar roles personalizados del tenant',
      'Operación: fincas, análisis YOLO, perfil',
    ],
  },
  agronomo: {
    label: 'Agrónomo',
    accent: 'blue',
    homeRoute: '/fincas',
    tagline: 'Soporte técnico en campo',
    capabilities: [
      'Consultar usuarios del tenant',
      'Operación: fincas, análisis YOLO, perfil',
    ],
  },
  caficultor: {
    label: 'Caficultor',
    accent: 'earth',
    homeRoute: '/fincas',
    tagline: 'Operación en tu finca',
    capabilities: [
      'Gestionar fincas y lotes',
      'Análisis YOLO de muestras',
      'Mi perfil (sin panel de administración)',
    ],
  },
};

/** Etiqueta legible del rol */
export function formatRol(rol) {
  return ROL_LABELS[rol] || rol || 'Usuario';
}

export function getRolProfile(rol) {
  return ROL_PROFILES[rol] || ROL_PROFILES.caficultor;
}

/** sys_admin global: JWT puede decir administrador pero permisos del servidor lo confirman */
export function isGlobalSysAdmin(user, permisos = []) {
  if (!user) return false;
  if (user.rol === ROLES.SYS_ADMIN) return true;
  return permisos.includes('identidad:sesion:consultar')
    && permisos.includes('identidad:ip:consultar');
}

/** Panel de administración de identidad */
export function puedeAccederAdmin(user, permisos = []) {
  if (!user) return false;
  if (user.rol === ROLES.SYS_ADMIN || user.rol === ROLES.ADMINISTRADOR) return true;
  if (user.rol === ROLES.CAFICULTOR) return false;

  const permisosAdmin = [
    'identidad:usuario:invitar',
    'identidad:usuario:modificar',
    'identidad:usuario:eliminar',
    'identidad:usuario:expulsar',
    'identidad:rol:crear',
    'identidad:rol:modificar',
    'identidad:rol:asignar',
    'identidad:rol:eliminar',
    'identidad:sesion:consultar',
    'identidad:ip:consultar',
  ];
  return permisosAdmin.some((p) => permisos.includes(p));
}

/** Configurar nombre/slug del tenant propio */
export function puedeConfigurarTenant(user, permisos = []) {
  if (!user) return false;
  return user.rol === ROLES.SYS_ADMIN || permisos.includes('identidad:tenant:configurar');
}

/** Ruta de inicio según rol */
export function getRutaInicio(user, permisos = []) {
  if (!user) return '/login';
  return getRolProfile(user.rol).homeRoute;
}

/** Pestañas visibles en Panel Admin */
export function getAdminTabs(user, permisos = []) {
  if (!puedeAccederAdmin(user, permisos)) return [];

  const tabs = [];

  const veUsuarios = user.rol === ROLES.SYS_ADMIN
    || user.rol === ROLES.ADMINISTRADOR
    || user.rol === ROLES.AGRONOMO
    || permisos.includes('identidad:usuario:consultar');
  if (veUsuarios) tabs.push('Usuarios');

  // Tab Invitaciones: solo para administradores y sys_admin (quienes pueden invitar)
  const veInvitaciones = user.rol === ROLES.SYS_ADMIN
    || user.rol === ROLES.ADMINISTRADOR
    || permisos.includes('identidad:usuario:invitar');
  if (veInvitaciones) tabs.push('Invitaciones');

  const veRoles = user.rol === ROLES.SYS_ADMIN
    || user.rol === ROLES.ADMINISTRADOR
    || ['identidad:rol:crear', 'identidad:rol:modificar', 'identidad:rol:asignar', 'identidad:permiso:consultar']
      .some((p) => permisos.includes(p));
  if (veRoles) tabs.push('Roles');

  if (user.rol === ROLES.SYS_ADMIN) {
    tabs.push('Sesiones', 'IPs Bloqueadas', 'Sistema');
  }

  return tabs;
}

/** Secciones del menú lateral según rol */
export function getNavItems(user, permisos = []) {
  const items = [
    { to: '/fincas', label: 'Mis Fincas', section: 'operacion' },
    { to: '/dashboard', label: 'Análisis YOLO', section: 'operacion' },
    { to: '/perfil', label: 'Mi Perfil', section: 'operacion' },
  ];

  if (puedeAccederAdmin(user, permisos)) {
    const adminLabel = user.rol === ROLES.AGRONOMO ? 'Panel Admin (limitado)' : 'Panel Admin';
    items.push({ to: '/admin', label: adminLabel, section: 'admin' });
  }

  if (puedeConfigurarTenant(user, permisos)) {
    items.push({ to: '/finca-config', label: 'Config. Tenant', section: 'admin' });
  }

  return items;
}
