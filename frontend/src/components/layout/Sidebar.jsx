import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { usePermisos } from '../../hooks/usePermisos';
import LogoIcon from '../LogoIcon';
import './Sidebar.css';

const IconDashboard = ({ className }) => (
  <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/>
    <rect x="14" y="14" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/>
  </svg>
);

const IconPerfil = ({ className }) => (
  <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
    <circle cx="12" cy="7" r="4"/>
  </svg>
);

const IconAdmin = ({ className }) => (
  <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
    <circle cx="9" cy="7" r="4"/>
    <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
    <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
  </svg>
);

const IconLogout = ({ className }) => (
  <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
    <polyline points="16 17 21 12 16 7"/>
    <line x1="21" y1="12" x2="9" y2="12"/>
  </svg>
);

const IconFincas = ({ className }) => (
  <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
    <polyline points="9 22 9 12 15 12 15 22"/>
  </svg>
);

const IconSettings = ({ className }) => (
  <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="12" cy="12" r="3"/>
    <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
  </svg>
);

const IconChevron = ({ className }) => (
  <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <polyline points="15 18 9 12 15 6"/>
  </svg>
);

const NAV_ICONS = {
  '/fincas': IconFincas,
  '/dashboard': IconDashboard,
  '/perfil': IconPerfil,
  '/admin': IconAdmin,
  '/finca-config': IconSettings,
};

export default function Sidebar({ collapsed, onToggle, mobileOpen, onMobileClose }) {
  const { user, logout } = useAuth();
  const { navItems, rolLabel, rolProfile } = usePermisos();
  const navigate = useNavigate();

  const handleLogout = async () => { await logout(); navigate('/login'); };

  const operacion = navItems.filter((i) => i.section === 'operacion');
  const admin = navItems.filter((i) => i.section === 'admin');

  const sidebarClass = [
    'sidebar',
    collapsed ? 'sidebar--collapsed' : '',
    mobileOpen ? 'sidebar--mobile-open' : '',
  ].filter(Boolean).join(' ');

  const renderLink = (item) => {
    const Icon = NAV_ICONS[item.to] || IconDashboard;
    return (
      <NavLink
        key={item.to}
        to={item.to}
        className={({ isActive }) => `sidebar__item ${isActive ? 'sidebar__item--active' : ''}`}
      >
        <Icon className="sidebar__item-icon" />
        <span className="sidebar__item-text">{item.label}</span>
      </NavLink>
    );
  };

  return (
    <>
      {mobileOpen && <div className="sidebar-overlay" onClick={onMobileClose} />}
      <aside className={sidebarClass}>
        <div className="sidebar__brand">
          <div className="sidebar__brand-left">
            <LogoIcon className="sidebar__logo" imgClassName="sidebar__logo-img" />
            <span className="sidebar__brand-name">Bunna</span>
          </div>
          <button className="sidebar__toggle" onClick={onToggle} aria-label="Colapsar menú">
            <IconChevron />
          </button>
        </div>

        {!collapsed && user?.rol && (
          <div className={`sidebar__role-badge sidebar__role-badge--${rolProfile.accent}`} title={`Rol: ${rolLabel}`}>
            {rolLabel}
          </div>
        )}

        <nav className="sidebar__nav">
          <span className="sidebar__section-label">Operación</span>
          {operacion.map(renderLink)}

          {admin.length > 0 && (
            <>
              <span className="sidebar__section-label">Administración</span>
              {admin.map(renderLink)}
            </>
          )}
        </nav>

        <div className="sidebar__footer">
          <button className="sidebar__item sidebar__item--danger" onClick={handleLogout}>
            <IconLogout className="sidebar__item-icon" />
            <span className="sidebar__item-text">Cerrar sesión</span>
          </button>
        </div>
      </aside>
    </>
  );
}
