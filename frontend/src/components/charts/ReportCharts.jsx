import {
  PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, CartesianGrid,
  Tooltip, Legend, ResponsiveContainer, LineChart, Line,
} from 'recharts';
import './ReportCharts.css';

const NITROGEN_COLORS = {
  Alto: '#22c55e',
  Medio: '#eab308',
  Bajo: '#ef4444',
  Error: '#9ca3af',
};

const DEFAULT_COLORS = ['#3a7d4f', '#6bb07e', '#eab308', '#ef4444', '#3b82f6'];

function ChartCard({ title, subtitle, children, empty }) {
  return (
    <div className="chart-card">
      <div className="chart-card__header">
        <h3 className="chart-card__title">{title}</h3>
        {subtitle && <p className="chart-card__subtitle">{subtitle}</p>}
      </div>
      {empty ? (
        <div className="chart-card__empty">Sin datos suficientes para graficar</div>
      ) : (
        <div className="chart-card__body">{children}</div>
      )}
    </div>
  );
}

export default function ReportCharts({ nitrogenData, estadoData, serieTemporal }) {
  const hasNitrogen = nitrogenData?.length > 0;
  const hasEstado = estadoData?.length > 0;
  const hasSerie = serieTemporal?.length > 0;

  return (
    <div className="charts-grid">
      <ChartCard title="Distribución de nitrógeno" subtitle="Por nivel detectado" empty={!hasNitrogen}>
        <ResponsiveContainer width="100%" height={220}>
          <PieChart>
            <Pie
              data={nitrogenData}
              dataKey="value"
              nameKey="name"
              cx="50%"
              cy="50%"
              innerRadius={50}
              outerRadius={80}
              paddingAngle={3}
              label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
            >
              {nitrogenData.map((entry) => (
                <Cell key={entry.name} fill={NITROGEN_COLORS[entry.name] || DEFAULT_COLORS[entry.name?.length % DEFAULT_COLORS.length]} />
              ))}
            </Pie>
            <Tooltip />
            <Legend />
          </PieChart>
        </ResponsiveContainer>
      </ChartCard>

      <ChartCard title="Estado de diagnósticos" subtitle="Aceptados, pendientes y rechazados" empty={!hasEstado}>
        <ResponsiveContainer width="100%" height={220}>
          <BarChart data={estadoData} layout="vertical" margin={{ left: 20 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
            <XAxis type="number" allowDecimals={false} />
            <YAxis type="category" dataKey="name" width={80} tick={{ fontSize: 12 }} />
            <Tooltip />
            <Bar dataKey="value" radius={[0, 4, 4, 0]}>
              {estadoData.map((entry) => (
                <Cell key={entry.name} fill={entry.color} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </ChartCard>

      <ChartCard
        title="Muestras en el tiempo"
        subtitle="Evolución de tomas de muestra"
        empty={!hasSerie}
      >
        <ResponsiveContainer width="100%" height={220}>
          <LineChart data={serieTemporal}>
            <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
            <XAxis dataKey="fecha" tick={{ fontSize: 11 }} />
            <YAxis allowDecimals={false} tick={{ fontSize: 11 }} />
            <Tooltip />
            <Line
              type="monotone"
              dataKey="total"
              stroke="#3a7d4f"
              strokeWidth={2.5}
              dot={{ fill: '#3a7d4f', r: 4 }}
              activeDot={{ r: 6 }}
            />
          </LineChart>
        </ResponsiveContainer>
      </ChartCard>
    </div>
  );
}
