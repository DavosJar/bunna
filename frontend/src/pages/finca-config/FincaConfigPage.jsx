import { useState, useEffect } from 'react';
import { useAuth } from '../../context/AuthContext';
import { usePermisos } from '../../hooks/usePermisos';
import { configurarTenant } from '../../services/identidadApi';
import Layout from '../../components/layout/Layout';
import '../perfil/Perfil.css';

export default function FincaConfigPage() {
  const { ownTenantId, fetchMisTenants, availableTenants } = useAuth();
  const { rolLabel } = usePermisos();
  const tenantPropio = availableTenants?.find((t) => t.id === ownTenantId) || availableTenants?.find((t) => t.es_propio);
  const tenantId = ownTenantId || tenantPropio?.id;
  const [nombre, setNombre] = useState('');
  const [slug, setSlug] = useState('');
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState(null);

  useEffect(() => {
    if (tenantPropio) {
      setNombre(tenantPropio.nombre || '');
      setSlug(tenantPropio.slug || '');
    }
  }, [tenantPropio]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!tenantId) return;
    setLoading(true);
    setMsg(null);
    try {
      await configurarTenant(tenantId, { nombre, slug });
      await fetchMisTenants();
      setMsg({ tipo: 'exito', texto: 'Configuración de la finca actualizada.' });
    } catch {
      setMsg({ tipo: 'error', texto: 'No se pudo actualizar. Verifica permisos identidad:tenant:configurar.' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout title="Configurar tenant" subtitle={`Datos de tu organización (${rolLabel}).`}>
      <div className="perfil-card">
        <h2 className="perfil-card__title">Datos del tenant</h2>
        {msg && (
          <div style={{
            padding: '0.75rem', marginBottom: '1rem', borderRadius: '0.5rem',
            background: msg.tipo === 'error' ? '#fef2f2' : '#f0fdf4',
            color: msg.tipo === 'error' ? '#991b1b' : '#166534',
          }}>{msg.texto}</div>
        )}
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label">Nombre</label>
            <input className="form-input" value={nombre} onChange={e => setNombre(e.target.value)} required />
          </div>
          <div className="form-group">
            <label className="form-label">Slug</label>
            <input className="form-input" value={slug} onChange={e => setSlug(e.target.value)} placeholder="mi-finca" />
          </div>
          <button type="submit" className="btn-primary" disabled={loading || !tenantId}>
            {loading ? 'Guardando...' : 'Guardar cambios'}
          </button>
        </form>
      </div>
    </Layout>
  );
}
