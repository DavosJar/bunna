import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { useDiagnosis } from '../../context/DiagnosisContext';
import {
  loadFincasLocal, loadLotesLocal, saveMuestraLocal, saveDiagnosticoLocal,
} from '../../services/localStore';
import {
  tomarMuestra, fincasApiDisponible, guardarResultadoManual,
  isServicioFincasNoDisponible,
} from '../../services/fincasApi';
import { tieneClorosisFromYolo } from '../../utils/yoloDiagnostico';
import { generateUUID } from '../../utils/uuid';
import { IconMapPin } from '../icons/Icons';
import './VincularDiagnosticoPanel.css';

/**
 * Panel para vincular un resultado YOLO a una finca → lote → muestra.
 * Funciona en modo local (siempre) y sincroniza con API fincas cuando está disponible.
 */
export default function VincularDiagnosticoPanel({ yoloResults, historialId, onVinculado, compact = false }) {
  const { user } = useAuth();
  const { marcarHistorialVinculado } = useDiagnosis();
  const navigate = useNavigate();

  const [fincas, setFincas] = useState([]);
  const [lotes, setLotes] = useState([]);
  const [fincaId, setFincaId] = useState('');
  const [loteId, setLoteId] = useState('');
  const [latitud, setLatitud] = useState('-0.123456');
  const [longitud, setLongitud] = useState('-78.123456');
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState(null);
  const [apiOnline, setApiOnline] = useState(null);

  useEffect(() => {
    if (!user?.id) return;
    setFincas(loadFincasLocal(user.id).filter((f) => f.estado === 'ACTIVA'));
    fincasApiDisponible().then(setApiOnline);
  }, [user?.id]);

  useEffect(() => {
    if (!user?.id || !fincaId) { setLotes([]); return; }
    setLotes(loadLotesLocal(user.id, fincaId).filter((l) => l.estado === 'ACTIVO'));
  }, [user?.id, fincaId]);

  const usarUbicacion = () => {
    if (!navigator.geolocation) {
      setMsg({ tipo: 'error', texto: 'Geolocalización no disponible en este navegador.' });
      return;
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setLatitud(pos.coords.latitude.toFixed(6));
        setLongitud(pos.coords.longitude.toFixed(6));
      },
      () => setMsg({ tipo: 'error', texto: 'No se pudo obtener la ubicación.' }),
    );
  };

  const handleVincular = async (e) => {
    e.preventDefault();
    if (!user?.id || !fincaId || !loteId || !yoloResults?.feedback) return;

    setLoading(true);
    setMsg(null);

    const finca = fincas.find((f) => f.id === fincaId);
    const lote = lotes.find((l) => l.id === loteId);
    const coords = { latitud: parseFloat(latitud), longitud: parseFloat(longitud) };

    let muestra = {
      id: generateUUID(),
      loteID: loteId,
      ...coords,
      createdAt: new Date().toISOString(),
      _offline: true,
    };

    let diagnostico = {
      id: generateUUID(),
      muestraID: muestra.id,
      estado: 'PENDIENTE',
      origen: 'yolo',
      yolo: {
        feedback: yoloResults.feedback,
        num_detections: yoloResults.num_detections,
        avg_confidence: yoloResults.avg_confidence,
        detections: yoloResults.detections || [],
        image: yoloResults.image,
      },
      tieneClorosis: tieneClorosisFromYolo(yoloResults),
      createdAt: new Date().toISOString(),
      _offline: true,
    };

    const offline = apiOnline === false || finca?._offline || lote?._offline;

    try {
      if (!offline) {
        const muestraApi = await tomarMuestra(fincaId, loteId, coords);
        muestra = { ...muestra, ...muestraApi, _offline: false };
        
        const payload = {
          imageURL: yoloResults.image || '',
          tieneClorosis: diagnostico.tieneClorosis || false,
          confianza: yoloResults.avg_confidence || 0,
          procesadoAt: diagnostico.createdAt
        };
        const apiD = await guardarResultadoManual(muestra.id, payload);
        
        diagnostico.id = apiD.id;
        diagnostico.diagnosticoID = apiD.id;
        diagnostico.muestraID = muestra.id;
        diagnostico._offline = false;
      }

      saveMuestraLocal(user.id, fincaId, loteId, muestra);
      saveDiagnosticoLocal(user.id, fincaId, loteId, diagnostico);

      const vinculado = {
        fincaId,
        loteId,
        muestraId: muestra.id,
        diagnosticoId: diagnostico.id,
        fincaNombre: finca?.nombre,
        loteNombre: lote?.nombre,
      };

      if (historialId) {
        marcarHistorialVinculado(historialId, vinculado);
      }

      const modo = offline ? 'localmente' : 'en el lote (muestra en servidor, diagnóstico YOLO local)';
      setMsg({
        tipo: 'exito',
        texto: `Diagnóstico vinculado ${modo}: ${finca?.nombre} → ${lote?.nombre}. Nitrógeno ${yoloResults.feedback.label}.`,
      });

      onVinculado?.(vinculado);
    } catch (err) {
      if (isServicioFincasNoDisponible(err)) {
        saveMuestraLocal(user.id, fincaId, loteId, muestra);
        saveDiagnosticoLocal(user.id, fincaId, loteId, diagnostico);
        if (historialId) {
          marcarHistorialVinculado(historialId, {
            fincaId, loteId, muestraId: muestra.id, diagnosticoId: diagnostico.id,
            fincaNombre: finca?.nombre, loteNombre: lote?.nombre,
          });
        }
        setMsg({
          tipo: 'warn',
          texto: `Guardado localmente en ${finca?.nombre} → ${lote?.nombre}. Sincroniza cuando fincas (:8082) esté activo.`,
        });
        onVinculado?.({ fincaId, loteId, muestraId: muestra.id });
      } else {
        setMsg({ tipo: 'error', texto: err.response?.data?.error || 'No se pudo vincular el diagnóstico.' });
      }
    } finally {
      setLoading(false);
    }
  };

  if (!yoloResults?.feedback || (yoloResults.feedback.level === 'unknown' && yoloResults.feedback.label === 'Error')) {
    return null;
  }

  if (fincas.length === 0) {
    return (
      <div className={`vincular-panel ${compact ? 'vincular-panel--compact' : ''}`}>
        <p className="vincular-panel__hint">
          Crea una finca en <button type="button" className="vincular-panel__link" onClick={() => navigate('/fincas')}>Mis Fincas</button> para vincular este diagnóstico.
        </p>
      </div>
    );
  }

  return (
    <div className={`vincular-panel ${compact ? 'vincular-panel--compact' : ''}`}>
      <h3 className="vincular-panel__title">Vincular a finca y lote</h3>
      <p className="vincular-panel__desc">
        Registra una muestra con GPS y asocia este diagnóstico YOLO al lote seleccionado.
      </p>

      {msg && (
        <div className={`vincular-panel__msg vincular-panel__msg--${msg.tipo}`}>{msg.texto}</div>
      )}

      <form onSubmit={handleVincular} className="vincular-panel__form">
        <div className="vincular-panel__row">
          <select className="form-input" value={fincaId} onChange={(e) => { setFincaId(e.target.value); setLoteId(''); }} required>
            <option value="">Seleccionar finca</option>
            {fincas.map((f) => (
              <option key={f.id} value={f.id}>{f.nombre}{f._offline ? ' (local)' : ''}</option>
            ))}
          </select>
          <select className="form-input" value={loteId} onChange={(e) => setLoteId(e.target.value)} required disabled={!fincaId}>
            <option value="">Seleccionar lote</option>
            {lotes.map((l) => (
              <option key={l.id} value={l.id}>{l.nombre} ({l.area} ha)</option>
            ))}
          </select>
        </div>

        {fincaId && lotes.length === 0 && (
          <p className="vincular-panel__hint">
            Esta finca no tiene lotes.{' '}
            <button type="button" className="vincular-panel__link" onClick={() => navigate('/fincas', { state: { fincaId } })}>
              Agregar lote
            </button>
          </p>
        )}

        <div className="vincular-panel__row">
          <input className="form-input" type="number" step="any" placeholder="Latitud" value={latitud} onChange={(e) => setLatitud(e.target.value)} required />
          <input className="form-input" type="number" step="any" placeholder="Longitud" value={longitud} onChange={(e) => setLongitud(e.target.value)} required />
          <button type="button" className="btn-admin btn-admin--icon" onClick={usarUbicacion}>
            <IconMapPin size={16} /> GPS
          </button>
        </div>

        <button type="submit" className="btn-analyze" disabled={loading || !fincaId || !loteId} style={{ marginTop: '0.5rem' }}>
          {loading ? 'Vinculando...' : 'Guardar en muestra'}
        </button>
      </form>
    </div>
  );
}
