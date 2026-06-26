import { ScatterChart, Scatter, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, ZAxis } from 'recharts';
import { IconMapPin } from '../icons/Icons';
import './MuestrasMap.css';

const ESTADO_COLORS = {
  ACEPTADO: '#22c55e',
  PENDIENTE: '#eab308',
  RECHAZADO: '#ef4444',
  SIN_DIAGNOSTICO: '#9ca3af',
};

function MapDot({ cx, cy, payload }) {
  if (cx == null || cy == null) return null;
  return (
    <circle
      cx={cx}
      cy={cy}
      r={9}
      fill={payload.fill}
      stroke="#fff"
      strokeWidth={2}
    />
  );
}

export default function MuestrasMap({ points }) {
  if (!points?.length) {
    return (
      <div className="muestras-map muestras-map--empty">
        <span className="muestras-map__empty-icon"><IconMapPin size={36} /></span>
        <p>No hay coordenadas GPS registradas en las muestras de este lote.</p>
      </div>
    );
  }

  const data = points.map((p) => ({
    ...p,
    x: p.lng,
    y: p.lat,
    z: 80,
    fill: ESTADO_COLORS[p.estado] || ESTADO_COLORS.SIN_DIAGNOSTICO,
  }));

  const lngs = data.map((d) => d.x);
  const lats = data.map((d) => d.y);
  const lngPad = (Math.max(...lngs) - Math.min(...lngs)) * 0.15 || 0.001;
  const latPad = (Math.max(...lats) - Math.min(...lats)) * 0.15 || 0.001;

  return (
    <div className="muestras-map">
      <div className="muestras-map__header">
        <h3 className="muestras-map__title">Mapa de muestras</h3>
        <div className="muestras-map__legend">
          {Object.entries(ESTADO_COLORS).map(([estado, color]) => (
            <span key={estado} className="muestras-map__legend-item">
              <span className="muestras-map__legend-dot" style={{ background: color }} />
              {estado.replace('_', ' ')}
            </span>
          ))}
        </div>
      </div>
      <ResponsiveContainer width="100%" height={280}>
        <ScatterChart margin={{ top: 10, right: 20, bottom: 10, left: 10 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
          <XAxis
            type="number"
            dataKey="x"
            name="Longitud"
            domain={[Math.min(...lngs) - lngPad, Math.max(...lngs) + lngPad]}
            tick={{ fontSize: 10 }}
            tickFormatter={(v) => v.toFixed(4)}
            label={{ value: 'Longitud', position: 'bottom', fontSize: 11, offset: -5 }}
          />
          <YAxis
            type="number"
            dataKey="y"
            name="Latitud"
            domain={[Math.min(...lats) - latPad, Math.max(...lats) + latPad]}
            tick={{ fontSize: 10 }}
            tickFormatter={(v) => v.toFixed(4)}
            label={{ value: 'Latitud', angle: -90, position: 'insideLeft', fontSize: 11 }}
          />
          <ZAxis type="number" dataKey="z" range={[60, 200]} />
          <Tooltip
            cursor={{ strokeDasharray: '3 3' }}
            content={({ active, payload }) => {
              if (!active || !payload?.[0]) return null;
              const p = payload[0].payload;
              return (
                <div className="muestras-map__tooltip">
                  <strong>Nitrógeno: {p.label}</strong>
                  <span>Estado: {p.estado}</span>
                  <span>Clorosis: {p.clorosis == null ? '—' : p.clorosis ? 'Sí' : 'No'}</span>
                  <span>{p.lat.toFixed(5)}, {p.lng.toFixed(5)}</span>
                </div>
              );
            }}
          />
          <Scatter data={data} shape={MapDot} />
        </ScatterChart>
      </ResponsiveContainer>
      <p className="muestras-map__note">{points.length} punto{points.length > 1 ? 's' : ''} GPS · Coordenadas WGS84</p>
    </div>
  );
}
