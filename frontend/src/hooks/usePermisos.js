import { useMemo } from 'react';
import { useAuth } from '../context/AuthContext';
import {
  puedeAccederAdmin,
  puedeConfigurarTenant,
  getRutaInicio,
  getAdminTabs,
  getNavItems,
  formatRol,
  getRolProfile,
  buildEffectiveUser,
  getPermisosEfectivos,
  puedePrevisualizarRol,
  isGlobalSysAdmin,
} from '../utils/roleAccess';

export function usePermisos() {
  const { user, permisos, rolPreview } = useAuth();

  const effectiveUser = useMemo(
    () => buildEffectiveUser(user, rolPreview, permisos),
    [user, rolPreview, permisos],
  );

  const permisosEfectivos = useMemo(
    () => getPermisosEfectivos(user, permisos, rolPreview),
    [user, permisos, rolPreview],
  );

  const puede = (permiso) => {
    if (!effectiveUser) return false;
    if (effectiveUser.rol === 'sys_admin') return true;
    return permisosEfectivos.includes(permiso);
  };

  const esAdmin = () => effectiveUser?.rol === 'sys_admin' || effectiveUser?.rol === 'administrador';
  const esSysAdmin = () => effectiveUser?.rol === 'sys_admin';
  const esCaficultor = () => effectiveUser?.rol === 'caficultor';
  const esAgronomo = () => effectiveUser?.rol === 'agronomo';

  return {
    puede,
    esAdmin,
    esSysAdmin,
    esCaficultor,
    esAgronomo,
    permisos: permisosEfectivos,
    permisosReales: permisos,
    effectiveUser,
    rolPreview,
    enVistaPrevia: Boolean(rolPreview),
    puedePrevisualizar: () => puedePrevisualizarRol(user, permisos),
    esGlobalSysAdmin: () => isGlobalSysAdmin(user, permisos),
    puedeAccederAdmin: () => puedeAccederAdmin(effectiveUser, permisosEfectivos),
    puedeConfigurarTenant: () => puedeConfigurarTenant(effectiveUser, permisosEfectivos),
    rutaInicio: getRutaInicio(effectiveUser, permisosEfectivos),
    adminTabs: getAdminTabs(effectiveUser, permisosEfectivos),
    navItems: getNavItems(effectiveUser, permisosEfectivos),
    rolLabel: formatRol(effectiveUser?.rol),
    rolRealLabel: formatRol(effectiveUser?.rolReal || user?.rol),
    rolProfile: getRolProfile(effectiveUser?.rol),
  };
}
