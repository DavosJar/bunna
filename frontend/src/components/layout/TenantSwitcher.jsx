import { useState, useRef, useEffect } from 'react';
import { useAuth } from '../../context/AuthContext';
import './TenantSwitcher.css';

export default function TenantSwitcher() {
  const { user, availableTenants, currentTenant, switchTenant, loading } = useAuth();
  const [open, setOpen] = useState(false);
  const ref = useRef(null);

  // Close on click outside
  useEffect(() => {
    const handleClick = (e) => {
      if (ref.current && !ref.current.contains(e.target)) setOpen(false);
    };
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, []);

  // Don't show if no tenants or only one tenant
  if (!availableTenants?.length || availableTenants.length <= 1) return null;

  const displayName = currentTenant?.nombre || availableTenants[0]?.nombre || 'Tenant';
  const displayRole = currentTenant?.rol || availableTenants[0]?.rol || '';

  const handleSwitch = async (tenantId) => {
    await switchTenant(tenantId);
    setOpen(false);
  };

  return (
    <div className="tenant-switcher" ref={ref}>
      <button
        className="tenant-switcher__trigger"
        onClick={() => setOpen(o => !o)}
        disabled={loading}
        aria-label="Cambiar de finca"
        aria-expanded={open}
      >
        <svg className="tenant-switcher__icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
          <line x1="8" y1="21" x2="16" y2="21"/>
          <line x1="12" y1="17" x2="12" y2="21"/>
        </svg>
        <span className="tenant-switcher__label">Cambiar finca</span>
        <span className="tenant-switcher__current">
          <span className="tenant-switcher__name">{displayName}</span>
          {displayRole && <span className="tenant-switcher__role">{displayRole}</span>}
        </span>
        <span className={`tenant-switcher__arrow ${open ? 'tenant-switcher__arrow--open' : ''}`}>
          <svg width="10" height="6" viewBox="0 0 10 6" fill="none">
            <path d="M1 1L5 5L9 1" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
          </svg>
        </span>
      </button>

      {open && (
        <div className="tenant-switcher__dropdown">
          <div className="tenant-switcher__dropdown-header">
            <span className="tenant-switcher__dropdown-title">Cambiar de finca</span>
          </div>
          {availableTenants.map(t => {
            const isActive = t.id === user?.tenantID;
            const isOwn = t.id === user?.ownTenantID;
            return (
              <button
                key={t.id}
                className={`tenant-switcher__option ${isActive ? 'tenant-switcher__option--active' : ''}`}
                onClick={() => handleSwitch(t.id)}
                disabled={loading || isActive}
              >
                <div className="tenant-switcher__option-left">
                  <span className="tenant-switcher__option-name">
                    {t.nombre}
                    {isOwn && <span className="tenant-switcher__own-badge" title="Mi finca">Propia</span>}
                  </span>
                  <span className="tenant-switcher__option-role">{t.rol}</span>
                </div>
                {isActive && (
                  <span className="tenant-switcher__check">
                    <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                      <path d="M2 7L5.5 10.5L12 4" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                    </svg>
                  </span>
                )}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
