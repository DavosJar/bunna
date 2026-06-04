import { useState } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
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

const IconChevron = ({ className }) => (
  <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <polyline points="15 18 9 12 15 6"/>
  </svg>
);

const IconMenu = ({ className }) => (
  <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="18" x2="21" y2="18"/>
  </svg>
);

export default function Sidebar({ collapsed, onToggle, mobileOpen, onMobileClose }) {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => { logout(); navigate('/login'); };

  const initials = user?.nombre
    ? user.nombre.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
    : '??';

  const sidebarClass = [
    'sidebar',
    collapsed ? 'sidebar--collapsed' : '',
    mobileOpen ? 'sidebar--mobile-open' : '',
  ].filter(Boolean).join(' ');

  return (
    <>
      {mobileOpen && <div className="sidebar-overlay" onClick={onMobileClose} />}
      <aside className={sidebarClass}>
        {/* Brand */}
        <div className="sidebar__brand">
          <div className="sidebar__brand-left">
            <LogoIcon className="sidebar__logo" imgClassName="sidebar__logo-img" />
            <span className="sidebar__brand-name">Bunna</span>
          </div>
          <button className="sidebar__toggle" onClick={onToggle} aria-label="Colapsar menú">
            <IconChevron className={undefined} />
          </button>
        </div>

        {/* User card */}
        <div className="sidebar__user-card">
          <div className="sidebar__avatar">{initials}</div>
          <div className="sidebar__user-info">
            <p className="sidebar__user-name">{user?.nombre || 'Usuario'}</p>
            <p className="sidebar__user-role">{user?.rol || 'Caficultor'}</p>
          </div>
        </div>

        {/* Nav */}
        <nav className="sidebar__nav">
          <span className="sidebar__section-label">General</span>

          <NavLink to="/dashboard" className={({ isActive }) => `sidebar__item ${isActive ? 'sidebar__item--active' : ''}`}>
            <IconDashboard className="sidebar__item-icon" />
            <span className="sidebar__item-text">Dashboard</span>
          </NavLink>

          <NavLink to="/perfil" className={({ isActive }) => `sidebar__item ${isActive ? 'sidebar__item--active' : ''}`}>
            <IconPerfil className="sidebar__item-icon" />
            <span className="sidebar__item-text">Mi Perfil</span>
          </NavLink>

          <span className="sidebar__section-label">Administración</span>

          <NavLink to="/admin" className={({ isActive }) => `sidebar__item ${isActive ? 'sidebar__item--active' : ''}`}>
            <IconAdmin className="sidebar__item-icon" />
            <span className="sidebar__item-text">Panel Admin</span>
          </NavLink>
        </nav>

        {/* Footer */}
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