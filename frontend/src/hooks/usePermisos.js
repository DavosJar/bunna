import { useAuth } from '../context/AuthContext';

export function usePermisos() {
  const { user, permisos } = useAuth();

  const puede = (permiso) => {
    if (!user || !permisos) return false;
    return permisos.includes(permiso);
  };

  const esAdmin = () => user?.rol === 'sys_admin' || user?.rol === 'administrador';
  const esSysAdmin = () => user?.rol === 'sys_admin';

  return { puede, esAdmin, esSysAdmin, permisos };
}
