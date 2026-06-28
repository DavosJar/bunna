import { usePermisos } from '../../hooks/usePermisos';
import './RoleAccessCard.css';

export default function RoleAccessCard() {
  const {
    rolLabel,
    rolProfile,
    puedeAccederAdmin,
    adminTabs,
  } = usePermisos();

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
    </div>
  );
}
