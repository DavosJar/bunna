import './StatCard.css';

export default function StatCard({ icon, label, value, sub, trend, accent = 'green', onClick }) {
  const Tag = onClick ? 'button' : 'div';
  return (
    <Tag
      type={onClick ? 'button' : undefined}
      className={`stat-card stat-card--${accent} ${onClick ? 'stat-card--clickable' : ''}`}
      onClick={onClick}
    >
      {icon && <div className="stat-card__icon">{icon}</div>}
      <div className="stat-card__body">
        <span className="stat-card__label">{label}</span>
        <span className="stat-card__value">{value}</span>
        {sub && <span className="stat-card__sub">{sub}</span>}
        {trend != null && (
          <span className={`stat-card__trend ${trend >= 0 ? 'stat-card__trend--up' : 'stat-card__trend--down'}`}>
            {trend >= 0 ? '↑' : '↓'} {Math.abs(trend)}%
          </span>
        )}
      </div>
    </Tag>
  );
}
