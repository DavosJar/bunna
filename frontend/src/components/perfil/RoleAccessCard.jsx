import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { usePermisos } from '../../hooks/usePermisos';
import { ROLES, ROL_PROFILES } from '../../utils/roleAccess';
import './RoleAccessCard.css';

const ROLES_PREVISUALIZABLES = [ROLES.CAFICULTOR, ROLES.AGRONOMO, ROLES.ADMINISTRADOR, ROLES.SYS_ADMIN];

export default function RoleAccessCard() {
  const navigate = useNavigate();
  const { setRolPreview } = useAuth();
  const {
    rolLabel,
    rolRealLabel,
    rolProfile,
    enVistaPrevia,
    puedePrevisualizar,
    esGlobalSysAdmin,
    rutaInicio,
    puedeAccederAdmin,
    adminTabs,
    effectiveUser,
  } = usePermisos();

  const handlePreview = (rol) => {
    if (rol === null) {
      setRolPreview(null);
      return;
    }
    setRolPreview(rol);
    navigate(rol === ROLES.CAFICULTOR || rol === ROLES.AGRONOMO ? '/fincas' : '/admin');
  };

  return (
    <div className={`role-card role-card--${rolProfile.accent}`}>
      <div className="role-card__header">
        <div>
          <span className="role-card__eyebrow">Tu rol actual</span>
          <h2 className="role-card__title">{rolLabel}</h2>
          <p className="role-card__tagline">{rolProfile.tagline}</p>
        </div>
        <span className={`role-card__badge role-card__badge--${rolProfile.accent}`}>{rolLabel}</span>
      </div>

      {/* Banner de vista previa activa */}
      {enVistaPrevia && (
        <div className="role-card__preview-banner">
          Vista simulada. Tu rol real es <strong>{rolRealLabel}</strong>.
          <button type="button" className="role-card__link-btn" onClick={() => handlePreview(null)}>
            Salir de simulación
          </button>
        </div>
      )}

      {/* Banner sys_admin global */}
      {esGlobalSysAdmin() && !enVistaPrevia && (
        <div className="role-card__sys-banner">
          Tienes permisos de <strong>Super Admin</strong> en el servidor.
          {rolRealLabel !== 'Super Admin' && ' El JWT puede mostrar otro rol; la UI usa tus permisos reales.'}
        </div>
      )}

      {/* Qué puedes hacer */}
      <div className="role-card__section">
        <h3>Qué puedes hacer</h3>
        <ul className="role-card__list">
          {rolProfile.capabilities.map((cap) => (
            <li key={cap}>{cap}</li>
          ))}
        </ul>
      </div>

      {/* Panel admin disponible */}
      {puedeAccederAdmin() && adminTabs.length > 0 && (
        <div className="role-card__section">
          <h3>Panel Admin</h3>
          <p className="role-card__meta">Pestañas disponibles: {adminTabs.join(', ')}</p>
        </div>
      )}

      {/* ── FLUJO REAL: Cómo invitar a tu equipo ── */}
      <div className="role-card__section role-card__section--guide">
        <h3>¿Cómo agregar a tu equipo?</h3>
        <div className="role-card__guide-grid">
          <div className="role-card__guide-item role-card__guide-item--highlight">
            <strong>📧 Invitar por correo (recomendado)</strong>
            <p>
              Ve a <strong>Panel Admin → Usuarios → + Invitar</strong>, ingresa el correo
              y selecciona el rol (caficultor o agrónomo). La persona recibirá un enlace
              para crear su cuenta y unirse a tu finca con ese rol.
            </p>
          </div>
          <div className="role-card__guide-item">
            <strong>🔄 Cambiar de finca</strong>
            <p>Si te invitan a otra finca, cambia de tenant en el selector superior. Cada finca puede tener un rol distinto para ti.</p>
          </div>
          <div className="role-card__guide-item">
            <strong>🔐 Super Admin</strong>
            <p>No se asigna desde la UI. Requiere asignación global en el servidor (API o base de datos).</p>
          </div>
        </div>
      </div>

      {/* ── SIMULACIÓN (secundario, solo para admins) ── */}
      {puedePrevisualizar() && (
        <div className="role-card__section">
          <h3>Vista previa de otro rol</h3>
          <p className="role-card__meta">
            Simula cómo se verían los menús y pantallas para cada rol. <em>No cambia tus permisos reales en el servidor.</em>
          </p>
          <div className="role-card__preview-btns">
            {ROLES_PREVISUALIZABLES.map((rol) => {
              const profile = ROL_PROFILES[rol];
              const activo = effectiveUser?.rol === rol;
              return (
                <button
                  key={rol}
                  type="button"
                  className={`role-card__preview-btn role-card__preview-btn--${profile.accent} ${activo ? 'role-card__preview-btn--active' : ''}`}
                  onClick={() => handlePreview(rol)}
                >
                  {profile.label}
                </button>
              );
            })}
            {enVistaPrevia && (
              <button type="button" className="role-card__preview-btn role-card__preview-btn--reset" onClick={() => handlePreview(null)}>
                Mi rol real
              </button>
            )}
          </div>
        </div>
      )}

      <div className="role-card__footer">
        <span>Inicio según rol: <strong>{rutaInicio}</strong></span>
      </div>
    </div>
  );
}
