import { useState, useRef } from 'react';
import { diagnosticar } from '../../services/yoloApi';
import { tomarMuestra, guardarResultadoManual, isServicioFincasNoDisponible } from '../../services/fincasApi';
import { saveMuestraLocal, saveDiagnosticoLocal } from '../../services/localStore';
import { tieneClorosisFromYolo } from '../../utils/yoloDiagnostico';
import { generateUUID } from '../../utils/uuid';
import { IconCamera, IconMapPin, IconX } from '../../components/icons/Icons';
import './AnalizarMuestraEnFinca.css';

/**
 * Flujo todo-en-uno: subir imagen(es) → GPS → YOLO → muestra + diagnóstico local.
 * Soporta múltiples imágenes en una sesión.
 */
export default function AnalizarMuestraEnFinca({
  userId,
  fincaId,
  loteId,
  apiOnline,
  onCompletado,
  onError,
  onProgreso,
}) {
  const fileRef = useRef(null);
  const [cola, setCola] = useState([]);
  const [latitud, setLatitud] = useState('-0.123456');
  const [longitud, setLongitud] = useState('-78.123456');
  const [analyzing, setAnalyzing] = useState(false);
  const [progreso, setProgreso] = useState(null);
  const [dragActive, setDragActive] = useState(false);
  const [ultimosResultados, setUltimosResultados] = useState([]);

  const usarUbicacion = () => {
    if (!navigator.geolocation) {
      onError?.('Geolocalización no disponible.');
      return;
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setLatitud(pos.coords.latitude.toFixed(6));
        setLongitud(pos.coords.longitude.toFixed(6));
      },
      () => onError?.('No se pudo obtener la ubicación.'),
    );
  };

  const agregarArchivos = (fileList) => {
    const nuevos = Array.from(fileList || [])
      .filter((f) => f.type?.startsWith('image/'))
      .map((file) => {
        const id = generateUUID();
        const preview = URL.createObjectURL(file);
        return { id, file, preview, nombre: file.name };
      });
    if (!nuevos.length) return;
    setCola((prev) => [...prev, ...nuevos]);
  };

  const quitarDeCola = (id) => {
    setCola((prev) => {
      const item = prev.find((x) => x.id === id);
      if (item?.preview) URL.revokeObjectURL(item.preview);
      return prev.filter((x) => x.id !== id);
    });
  };

  const limpiarCola = () => {
    cola.forEach((item) => { if (item.preview) URL.revokeObjectURL(item.preview); });
    setCola([]);
    if (fileRef.current) fileRef.current.value = '';
  };

  const procesarUnaImagen = async (item, coords) => {
    let muestra = {
      id: generateUUID(),
      loteID: loteId,
      ...coords,
      nombreArchivo: item.file.name,
      createdAt: new Date().toISOString(),
      _offline: true,
    };

    if (apiOnline === true) {
      try {
        const apiM = await tomarMuestra(fincaId, loteId, coords);
        muestra = { ...muestra, ...apiM, _offline: false };
      } catch (err) {
        if (!isServicioFincasNoDisponible(err)) throw err;
      }
    }

    const yolo = await diagnosticar(item.file);

    if (yolo?.feedback?.label === 'Error') {
      throw new Error(yolo.feedback.recommendation || 'Error al analizar la imagen.');
    }

    let diagnostico = {
      id: generateUUID(),
      muestraID: muestra.id,
      estado: 'PENDIENTE',
      origen: 'yolo',
      yolo: {
        feedback: yolo.feedback,
        num_detections: yolo.num_detections,
        avg_confidence: yolo.avg_confidence,
        detections: yolo.detections || [],
        image: yolo.image,
      },
      tieneClorosis: tieneClorosisFromYolo(yolo),
      createdAt: new Date().toISOString(),
      _offline: true,
    };

    if (apiOnline === true && !muestra._offline) {
      try {
        const payload = {
          imageURL: yolo.image || '',
          tieneClorosis: diagnostico.tieneClorosis || false,
          confianza: yolo.avg_confidence || 0,
          procesadoAt: diagnostico.createdAt
        };
        const apiD = await guardarResultadoManual(muestra.id, payload);
        diagnostico = { 
          ...diagnostico,
          id: apiD.id,
          diagnosticoID: apiD.id,
          nombre: apiD.nombre,
          _offline: false 
        };
      } catch (err) {
        if (!isServicioFincasNoDisponible(err)) console.error('Error al guardar diagnóstico manual:', err);
      }
    }

    saveMuestraLocal(userId, fincaId, loteId, muestra);
    saveDiagnosticoLocal(userId, fincaId, loteId, diagnostico);

    return { muestra, diagnostico, yolo };
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!cola.length || !latitud || !longitud) return;

    setAnalyzing(true);
    const coords = { latitud: parseFloat(latitud), longitud: parseFloat(longitud) };
    const resultados = [];
    const total = cola.length;

    try {
      for (let i = 0; i < cola.length; i++) {
        const item = cola[i];
        setProgreso({ actual: i + 1, total, nombre: item.nombre });
        onProgreso?.({ actual: i + 1, total });

        const resultado = await procesarUnaImagen(item, coords);
        resultados.push(resultado);
        onCompletado?.({
          muestraId: resultado.muestra.id,
          diagnosticoId: resultado.diagnostico.id,
          label: resultado.yolo.feedback?.label,
          esUltima: i === cola.length - 1,
          totalEnLote: i + 1,
        });
      }

      setUltimosResultados((prev) => [...resultados, ...prev].slice(0, 10));
      limpiarCola();
      setProgreso(null);
    } catch (err) {
      const msg = err.response?.data?.detail || err.response?.data?.error || err.message;
      onError?.(`No se pudo completar el análisis: ${msg}`);
      if (resultados.length) {
        setUltimosResultados((prev) => [...resultados, ...prev].slice(0, 10));
        limpiarCola();
      }
    } finally {
      setAnalyzing(false);
      setProgreso(null);
    }
  };

  return (
    <div className="analizar-finca">
      <p className="analizar-finca__intro">
        Sube una o varias fotos de hoja de café, indica la ubicación GPS y el sistema creará las muestras y ejecutará el diagnóstico YOLO automáticamente.
      </p>

      <form onSubmit={handleSubmit} className="analizar-finca__form">
        <div
          className={`analizar-finca__dropzone ${dragActive ? 'analizar-finca__dropzone--active' : ''}`}
          onDragEnter={(e) => { e.preventDefault(); setDragActive(true); }}
          onDragLeave={(e) => { e.preventDefault(); setDragActive(false); }}
          onDragOver={(e) => e.preventDefault()}
          onDrop={(e) => {
            e.preventDefault();
            setDragActive(false);
            agregarArchivos(e.dataTransfer.files);
          }}
          onClick={() => fileRef.current?.click()}
        >
          <input
            ref={fileRef}
            type="file"
            accept="image/jpeg,image/png,.jpg,.jpeg,.png"
            multiple
            className="analizar-finca__input"
            onChange={(e) => {
              agregarArchivos(e.target.files);
              if (fileRef.current) fileRef.current.value = '';
            }}
          />
          <span className="analizar-finca__icon"><IconCamera size={32} /></span>
          <p>Haz clic o arrastra imágenes de hoja</p>
          <p className="analizar-finca__hint">PNG o JPG · Puedes seleccionar varias a la vez</p>
        </div>

        {cola.length > 0 && (
          <div className="analizar-finca__cola">
            <div className="analizar-finca__cola-header">
              <strong>{cola.length} imagen{cola.length > 1 ? 'es' : ''} en cola</strong>
              <button type="button" className="btn-admin btn-admin--danger" onClick={limpiarCola}>
                Limpiar todo
              </button>
            </div>
            <div className="analizar-finca__cola-grid">
              {cola.map((item) => (
                <div key={item.id} className="analizar-finca__cola-item">
                  <img src={item.preview} alt={item.nombre} className="analizar-finca__cola-thumb" />
                  <span className="analizar-finca__cola-name" title={item.nombre}>{item.nombre}</span>
                  <button
                    type="button"
                    className="analizar-finca__cola-remove"
                    onClick={(ev) => { ev.stopPropagation(); quitarDeCola(item.id); }}
                    aria-label="Quitar imagen"
                  >
                    <IconX size={12} />
                  </button>
                </div>
              ))}
            </div>
          </div>
        )}

        <div className="analizar-finca__coords">
          <input
            className="form-input"
            type="number"
            step="any"
            placeholder="Latitud"
            value={latitud}
            onChange={(e) => setLatitud(e.target.value)}
            required
          />
          <input
            className="form-input"
            type="number"
            step="any"
            placeholder="Longitud"
            value={longitud}
            onChange={(e) => setLongitud(e.target.value)}
            required
          />
          <button type="button" className="btn-admin btn-admin--icon" onClick={usarUbicacion}>
            <IconMapPin size={16} /> Usar GPS
          </button>
        </div>

        <button
          type="submit"
          className="btn-add analizar-finca__submit"
          disabled={analyzing || !cola.length || !latitud || !longitud}
        >
          {analyzing
            ? (progreso ? `Analizando ${progreso.actual}/${progreso.total}: ${progreso.nombre}` : 'Analizando con YOLO...')
            : cola.length > 1
              ? `Analizar ${cola.length} imágenes`
              : 'Tomar muestra y analizar imagen'}
        </button>
      </form>

      {ultimosResultados.length > 0 && (
        <div className="analizar-finca__resultado">
          <p><strong>Últimos análisis ({ultimosResultados.length}):</strong></p>
          <div className="analizar-finca__resultados-grid">
            {ultimosResultados.map((r) => (
              <div key={r.muestra.id} className="analizar-finca__resultado-item">
                <p className="analizar-finca__resultado-label">
                  Nitrógeno {r.yolo.feedback?.label}
                  {r.diagnostico.tieneClorosis != null && (
                    <> · Clorosis: {r.diagnostico.tieneClorosis ? 'Sí' : 'No'}</>
                  )}
                </p>
                {r.yolo.image && (
                  <img src={r.yolo.image} alt="Resultado YOLO" className="analizar-finca__yolo-img" />
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
