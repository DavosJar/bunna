import { useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import Sidebar from './Sidebar';
import './Layout.css';

export default function Layout({ children, title, subtitle }) {
  const { user } = useAuth();
  const [collapsed, setCollapsed] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);

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
            <div className="layout__topbar-avatar">{initials}</div>
            <div className="layout__topbar-info">
              <span className="layout__topbar-name">{user?.nombre || 'Usuario'}</span>
              <span className="layout__topbar-role">{user?.rol || 'Caficultor'}</span>
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