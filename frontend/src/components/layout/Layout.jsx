import { useState, useRef, useEffect, useCallback } from 'react';
import { useAuth } from '../../context/AuthContext';
import Sidebar from './Sidebar';
import TenantSwitcher from './TenantSwitcher';
import './Layout.css';

export default function Layout({ children, title, subtitle }) {
  const { user, logout } = useAuth();
  const [collapsed, setCollapsed] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const dropdownRef = useRef(null);

  const handleClickOutside = useCallback((e) => {
    if (dropdownRef.current && !dropdownRef.current.contains(e.target)) {
      setDropdownOpen(false);
    }
  }, []);

  useEffect(() => {
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [handleClickOutside]);

  const handleLogout = () => {
    setDropdownOpen(false);
    logout();
  };

  const initials = user?.nombre
    ? user.nombre.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
    : '??';

  return (
    <div className="layout">
      <Sidebar
        collapsed={collapsed}
        onToggle={() => setCollapsed(c => !c)}
        mobileOpen={mobileOpen}
        onMobileClose={() => setMobileOpen(false)}
      />

      <div className={`layout__sidebar-space ${collapsed ? 'layout__sidebar-space--collapsed' : ''}`} />

      <div className="layout__content">
        <header className="layout__topbar">
          <div className="layout__topbar-left">
            <button className="layout__mobile-toggle" onClick={() => setMobileOpen(o => !o)} aria-label="Abrir menú">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="18" x2="21" y2="18"/>
              </svg>
            </button>
            <div className="layout__breadcrumb">
              {title && <span className="layout__page-title">{title}</span>}
              {subtitle && <span className="layout__page-subtitle">{subtitle}</span>}
            </div>
          </div>
          <div className="layout__topbar-right">
            <TenantSwitcher />
            <div className="layout__topbar-profile" ref={dropdownRef}>
              <button
                className="layout__topbar-profile-btn"
                onClick={() => setDropdownOpen(o => !o)}
                aria-label="Menú de perfil"
                aria-expanded={dropdownOpen}
              >
                <div className="layout__topbar-avatar">{initials}</div>
                <div className="layout__topbar-info">
                  <span className="layout__topbar-name">{user?.nombre || 'Usuario'}</span>
                  <span className="layout__topbar-role">{user?.rol || 'Caficultor'}</span>
                </div>
              </button>

              {dropdownOpen && (
                <div className="layout__topbar-dropdown">
                  <button className="layout__topbar-dropdown-item" onClick={handleLogout}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" width="16" height="16">
                      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                      <polyline points="16 17 21 12 16 7" />
                      <line x1="21" y1="12" x2="9" y2="12" />
                    </svg>
                    Cerrar sesión
                  </button>
                </div>
              )}
            </div>
          </div>
        </header>

        <main className="layout__main">
          {children}
        </main>
      </div>
    </div>
  );
}