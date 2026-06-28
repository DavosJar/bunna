import { useAuth } from '../context/AuthContext';
import {
  puedeAccederAdmin,
  puedeConfigurarTenant,
  getRutaInicio,
  getAdminTabs,
  getNavItems,
  formatRol,
  getRolProfile,
} from '../utils/roleAccess';

export function usePermisos() {
  const { user, permisos } = useAuth();

  const puede = (permiso) => {
    if (!user) return false;
    if (user.rol === 'sys_admin') return true;
    return permisos.includes(permiso);
  };

  const esAdmin = () => user?.rol === 'sys_admin' || user?.rol === 'administrador';
  const esSysAdmin = () => user?.rol === 'sys_admin';
  const esCaficultor = () => user?.rol === 'caficultor';
  const esAgronomo = () => user?.rol === 'agronomo';

  return {
    puede,
    esAdmin,
    esSysAdmin,
    esCaficultor,
    esAgronomo,
    permisos,
    permisosReales: permisos,
    puedeAccederAdmin: () => puedeAccederAdmin(user, permisos),
    puedeConfigurarTenant: () => puedeConfigurarTenant(user, permisos),
    rutaInicio: getRutaInicio(user, permisos),
    adminTabs: getAdminTabs(user, permisos),
    navItems: getNavItems(user, permisos),
    rolLabel: formatRol(user?.rol),
    rolProfile: getRolProfile(user?.rol),
  };
}
