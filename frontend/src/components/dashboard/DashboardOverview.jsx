import { useNavigate } from 'react-router-dom';
import StatCard from '../ui/StatCard';
import ReportCharts from '../charts/ReportCharts';
import {
  IconFarm, IconGrid, IconSample, IconAlert, IconList, IconFarm as IconHome, IconInsight,
} from '../icons/Icons';
import './DashboardOverview.css';

export default function DashboardOverview({ stats, onNavigate }) {
  const navigate = useNavigate();

  const goFincas = (tab) => navigate('/fincas', { state: tab ? { activeTab: tab } : undefined });

  return (
    <div className="dash-overview">
      <div className="stat-grid">
        <StatCard icon={<IconFarm />} label="Fincas activas" value={stats.fincas} sub="Registradas en el sistema" accent="green" onClick={() => goFincas()} />
        <StatCard icon={<IconGrid />} label="Lotes" value={stats.lotes} sub="Áreas de cultivo" accent="earth" onClick={() => goFincas()} />
        <StatCard icon={<IconSample />} label="Muestras" value={stats.muestras} sub="Imágenes analizadas" accent="blue" />
        <StatCard
          icon={<IconAlert />}
          label="Pendientes"
          value={stats.pendientes}
          sub={stats.pendientes > 0 ? 'Requieren revisión' : 'Todo al día'}
          accent="amber"
          onClick={stats.pendientes > 0 ? () => goFincas('diagnosticos') : undefined}
        />
      </div>

      {stats.diagnosticos > 0 && (
        <ReportCharts
          nitrogenData={stats.nitrogenData}
          estadoData={stats.estadoData}
          serieTemporal={stats.serieTemporal}
        />
      )}

      <div className="dash-overview__grid">
        <div className="dash-panel">
          <div className="dash-panel__header">
            <h3 className="dash-panel__title">Actividad reciente</h3>
            <button type="button" className="dash-panel__link" onClick={() => onNavigate?.('historial')}>
              Ver historial →
            </button>
          </div>
          {stats.recientes.length === 0 ? (
            <div className="dash-panel__empty">
              <span className="dash-panel__empty-icon"><IconList size={32} /></span>
              <p>Aún no hay análisis. Sube una imagen en la pestaña Análisis o registra una finca.</p>
              <button type="button" className="btn-add" onClick={() => onNavigate?.('analisis')}>
                Analizar imagen
              </button>
            </div>
          ) : (
            <ul className="activity-list">
              {stats.recientes.map((item) => (
                <li key={`${item.tipo}-${item.id}`} className="activity-list__item">
                  <div className="activity-list__dot activity-list__dot--green" />
                  <div className="activity-list__body">
                    <span className="activity-list__title">{item.titulo}</span>
                    <span className="activity-list__sub">{item.subtitulo}</span>
                  </div>
                  <div className="activity-list__meta">
                    <span className={`activity-badge activity-badge--${item.estado?.toLowerCase()}`}>
                      {item.estado}
                    </span>
                    <time>{new Date(item.fecha).toLocaleDateString('es-EC')}</time>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="dash-panel">
          <div className="dash-panel__header">
            <h3 className="dash-panel__title">Mis fincas</h3>
            <button type="button" className="dash-panel__link" onClick={() => goFincas()}>
              Gestionar →
            </button>
          </div>
          {stats.fincasDetalle.length === 0 ? (
            <div className="dash-panel__empty">
              <span className="dash-panel__empty-icon"><IconHome size={32} /></span>
              <p>Registra tu primera finca para organizar lotes, muestras y reportes.</p>
              <button type="button" className="btn-add" onClick={() => goFincas()}>
                Crear finca
              </button>
            </div>
          ) : (
            <ul className="finca-mini-list">
              {stats.fincasDetalle.map((f) => (
                <li key={f.id}>
                  <button
                    type="button"
                    className="finca-mini-card"
                    onClick={() => navigate('/fincas', { state: { fincaId: f.id } })}
                  >
                    <div className="finca-mini-card__name">
                      {f.nombre}
                      {f._offline && <span className="finca-mini-card__badge">local</span>}
                    </div>
                    <div className="finca-mini-card__stats">
                      <span>{f.lotes} lote{f.lotes !== 1 ? 's' : ''}</span>
                      <span>·</span>
                      <span>{f.muestras} muestra{f.muestras !== 1 ? 's' : ''}</span>
                      {f.pendientes > 0 && (
                        <>
                          <span>·</span>
                          <span className="finca-mini-card__warn">{f.pendientes} pendiente{f.pendientes !== 1 ? 's' : ''}</span>
                        </>
                      )}
                    </div>
                    <span className="finca-mini-card__loc">{f.ubicacion}</span>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>

      {stats.aceptados > 0 && (
        <div className="dash-insight">
          <div className="dash-insight__icon"><IconInsight size={24} /></div>
          <div>
            <strong>Resumen agronómico</strong>
            <p>
              {stats.porcentajeAfectado >= 50
                ? `El ${stats.porcentajeAfectado.toFixed(0)}% de diagnósticos aceptados muestran clorosis. Revisa fertilización nitrogenada.`
                : stats.porcentajeAfectado > 0
                  ? `Clorosis detectada en ${stats.porcentajeAfectado.toFixed(0)}% de muestras aceptadas. Monitoreo recomendado.`
                  : 'Los niveles de nitrógeno son óptimos en las muestras aceptadas. Mantén el plan actual.'}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
