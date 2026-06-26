import { useState, useEffect } from 'react';
import { verificarSaludYOLO } from '../../services/yoloApi';
import { fincasApiDisponible } from '../../services/fincasApi';
import { IconRefresh } from '../icons/Icons';
import './ServiceStatusBar.css';

async function pingIdentidad() {
  const base = import.meta.env.VITE_API_URL || '/api';
  const res = await fetch(`${base}/health`, { signal: AbortSignal.timeout(5000) });
  if (!res.ok) throw new Error('No responde');
  return res.json();
}

export default function ServiceStatusBar({ compact = false }) {
  const [services, setServices] = useState({ identidad: null, yolo: null, fincas: null });
  const [loading, setLoading] = useState(true);

  const check = async () => {
    setLoading(true);
    const next = { identidad: 'down', yolo: 'down', fincas: 'down' };
    try { await pingIdentidad(); next.identidad = 'up'; } catch { /* */ }
    try { await verificarSaludYOLO(); next.yolo = 'up'; } catch { /* */ }
    try { const ok = await fincasApiDisponible(); next.fincas = ok ? 'up' : 'down'; } catch { /* */ }
    setServices(next);
    setLoading(false);
  };

  useEffect(() => { check(); }, []);

  const items = [
    { key: 'identidad', label: 'Identidad', port: ':8080' },
    { key: 'yolo', label: 'YOLO IA', port: '' },
    { key: 'fincas', label: 'Fincas', port: ':8082' },
  ];

  return (
    <div className={`service-bar ${compact ? 'service-bar--compact' : ''}`}>
      {items.map(({ key, label, port }) => (
        <div key={key} className={`service-bar__item service-bar__item--${services[key] || 'down'}`}>
          <span className="service-bar__dot" />
          <span className="service-bar__label">{label}</span>
          {!compact && port && <span className="service-bar__port">{port}</span>}
        </div>
      ))}
      <button type="button" className="service-bar__refresh" onClick={check} disabled={loading} title="Verificar servicios">
        <IconRefresh size={14} className={loading ? 'service-bar__refresh--spin' : ''} />
      </button>
    </div>
  );
}
