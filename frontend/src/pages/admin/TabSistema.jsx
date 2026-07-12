import { useState, useEffect } from 'react';
import axios from 'axios';
import { verificarSaludYOLO } from '../../services/yoloApi';
import { verificarSaludFincas } from '../../services/fincasApi';
import '../admin/Admin.css';

async function pingIdentidad() {
  const base = import.meta.env.VITE_API_URL || 'http://localhost:8080';
  const res = await axios.get(`${base}/api/v1/identidad/health`, { timeout: 5000 });
  return res.data;
}

export default function TabSistema() {
  const [estado, setEstado] = useState({ identidad: null, yolo: null, fincas: null });
  const [loading, setLoading] = useState(false);

  const verificar = async () => {
    setLoading(true);
    const next = { identidad: null, yolo: null, fincas: null };
    try { next.identidad = await pingIdentidad(); } catch (e) { next.identidad = { error: e.message }; }
    try { next.yolo = await verificarSaludYOLO(); } catch (e) { next.yolo = { error: e.message }; }
    try { next.fincas = await verificarSaludFincas(); } catch (e) { next.fincas = { error: e.message }; }
    setEstado(next);
    setLoading(false);
  };

  useEffect(() => { verificar(); }, []);

  const Item = ({ nombre, data }) => (
    <div className="admin-card" style={{ marginBottom: '1rem' }}>
      <h3 className="admin-card__title">{nombre}</h3>
      {data?.error ? (
        <p style={{ color: '#991b1b' }}>No disponible: {data.error}</p>
      ) : data ? (
        <pre style={{ fontSize: '0.8rem', background: '#f8fafc', padding: '0.75rem', borderRadius: '0.5rem', overflow: 'auto' }}>
          {JSON.stringify(data, null, 2)}
        </pre>
      ) : (
        <p className="admin-empty">Sin datos</p>
      )}
    </div>
  );

  return (
    <>
      <div style={{ marginBottom: '1rem' }}>
        <button className="btn-add" onClick={verificar} disabled={loading}>
          {loading ? 'Verificando...' : '↻ Verificar servicios'}
        </button>
        <a href={`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/docs`} target="_blank" rel="noreferrer" className="btn-admin btn-admin--primary" style={{ marginLeft: '0.5rem' }}>
          API Docs (Identidad)
        </a>
      </div>
      <Item nombre="Identidad (:8080)" data={estado.identidad} />
      <Item nombre="YOLOv11" data={estado.yolo} />
      <Item nombre="Fincas (:8082)" data={estado.fincas} />
      <p style={{ fontSize: '0.8rem', color: 'var(--color-gray-500)' }}>
        Monitoreo (Kafka/ClickHouse/Grafana) no expone API REST al frontend; usa Grafana en infraestructura.
      </p>
    </>
  );
}
