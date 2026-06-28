import { useState, useEffect, useCallback, useMemo } from 'react';
import { useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { useDiagnosis } from '../../context/DiagnosisContext';
import Layout from '../../components/layout/Layout';
import VincularDiagnosticoPanel from '../../components/fincas/VincularDiagnosticoPanel';
import AnalizarMuestraEnFinca from '../../components/fincas/AnalizarMuestraEnFinca';
import {
  registrarFinca, desactivarFinca, agregarLote, eliminarLote,
  listarMuestras, generarReporteLote,
  aceptarDiagnostico, rechazarDiagnostico,
  fincasApiDisponible, isServicioFincasNoDisponible, parseErrorFincas,
  registrarNodo, listarNodos, desactivarNodo,
} from '../../services/fincasApi';
import {
  loadFincasLocal, saveFincaLocal, updateFincaLocal,
  loadLotesLocal, saveLoteLocal, updateLoteLocal,
  loadMuestrasLocal, saveMuestraLocal,
  loadDiagnosticosLocal, saveDiagnosticoLocal, updateDiagnosticoLocal,
  getDiagnosticoPorMuestra, saveWorkflowState, loadWorkflowState,
  loadNodosLocal, saveNodoLocal, updateNodoLocal,
} from '../../services/localStore';
import { buildReporteLocal, tieneClorosisFromYolo } from '../../utils/yoloDiagnostico';
import { generateReportePDF } from '../../utils/generateReportePDF';
import { generateUUID } from '../../utils/uuid';
import { getGlobalFarmStats, getReportChartData, getMuestrasMapPoints } from '../../utils/farmAnalytics';
import StatCard from '../../components/ui/StatCard';
import ServiceStatusBar from '../../components/ui/ServiceStatusBar';
import ReportCharts from '../../components/charts/ReportCharts';
import MuestrasMap from '../../components/fincas/MuestrasMap';
import { IconFarm, IconGrid, IconSample, IconClipboard, IconChart, IconCheck } from '../../components/icons/Icons';
import '../../components/ui/StatCard.css';
import '../../components/charts/ReportCharts.css';
import '../../components/fincas/MuestrasMap.css';
import '../../components/fincas/VincularDiagnosticoPanel.css';
import '../../components/fincas/AnalizarMuestraEnFinca.css';
import '../admin/Admin.css';

const TABS = [
  { id: 'fincas', label: 'Fincas' },
  { id: 'lotes', label: 'Lotes' },
  { id: 'nodos', label: 'Cámaras (IoT)' },
  { id: 'muestras', label: 'Muestras' },
  { id: 'diagnosticos', label: 'Diagnósticos' },
  { id: 'reporte', label: 'Reporte' },
];

function DiagBadge({ estado }) {
  const map = {
    PENDIENTE: 'diag-badge--pendiente',
    ACEPTADO: 'diag-badge--aceptado',
    RECHAZADO: 'diag-badge--rechazado',
    SIN_DIAGNOSTICO: 'diag-badge--sin',
  };
  const labels = {
    PENDIENTE: 'Pendiente',
    ACEPTADO: 'Aceptado',
    RECHAZADO: 'Rechazado',
    SIN_DIAGNOSTICO: 'Sin diagnóstico',
  };
  return <span className={`diag-badge ${map[estado] || 'diag-badge--sin'}`}>{labels[estado] || estado}</span>;
}

export default function FincasPage() {
  const { user } = useAuth();
  const location = useLocation();
  const { results: yoloPendiente, activeHistorialId, marcarHistorialVinculado } = useDiagnosis();

  const [fincas, setFincas] = useState([]);
  const [lotes, setLotes] = useState([]);
  const [nodos, setNodos] = useState([]);
  const [muestras, setMuestras] = useState([]);
  const [diagnosticos, setDiagnosticos] = useState([]);
  const [reporte, setReporte] = useState(null);
  const [fincaSel, setFincaSel] = useState('');
  const [loteSel, setLoteSel] = useState('');
  const [nodoSel, setNodoSel] = useState('');
  const [muestraSel, setMuestraSel] = useState('');
  const [msg, setMsg] = useState(null);
  const [loading, setLoading] = useState(false);
  const [apiOnline, setApiOnline] = useState(null);
  const [motivoRechazo, setMotivoRechazo] = useState('');
  const [activeTab, setActiveTab] = useState('fincas');
  const [workflowRestored, setWorkflowRestored] = useState(false);

  const [formFinca, setFormFinca] = useState({ nombre: '', ubicacion: '', descripcion: '' });
  const [formLote, setFormLote] = useState({ nombre: '', area: '', descripcion: '' });
  const [formNodo, setFormNodo] = useState({ nombre: '', node_key: '', lote_id: '' });

  const fincaActiva = fincas.find((f) => f.id === fincaSel);
  const loteActivo = lotes.find((l) => l.id === loteSel);
  const globalStats = useMemo(() => getGlobalFarmStats(user?.id), [user?.id, fincas, lotes, muestras, diagnosticos]);
  const reportCharts = useMemo(() => getReportChartData(reporte), [reporte]);
  const mapPoints = useMemo(() => getMuestrasMapPoints(reporte), [reporte]);
  const diagActivo = muestraSel ? getDiagnosticoPorMuestra(user?.id, fincaSel, loteSel, muestraSel)
    || diagnosticos.find((d) => d.muestraID === muestraSel) : null;

  const tabHabilitado = (tabId) => {
    if (tabId === 'fincas') return true;
    if (tabId === 'lotes' || tabId === 'nodos' || tabId === 'muestras') return !!fincaSel;
    return !!fincaSel && !!loteSel;
  };

  const irATab = (tabId) => {
    if (!tabHabilitado(tabId)) return;
    setActiveTab(tabId);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const cargarLocal = useCallback(() => {
    if (!user?.id) return;
    const f = loadFincasLocal(user.id).filter((x) => x.estado !== 'PENDIENTE_ELIMINACION');
    setFincas(f);
    if (fincaSel) {
      setLotes(loadLotesLocal(user.id, fincaSel).filter((l) => l.estado !== 'ELIMINADO'));
      setNodos(loadNodosLocal(user.id, fincaSel).filter((n) => n.estado !== 'INACTIVO'));
      setMuestras(loadMuestrasLocal(user.id, fincaSel, loteSel || ''));
      setDiagnosticos(loadDiagnosticosLocal(user.id, fincaSel, loteSel || ''));
    }
  }, [user?.id, fincaSel, loteSel]);

  useEffect(() => { cargarLocal(); }, [cargarLocal]);

  // Restaurar finca/lote/tab guardados al volver a la página
  useEffect(() => {
    if (!user?.id || workflowRestored) return;
    const saved = loadWorkflowState(user.id);
    if (!saved) {
      setWorkflowRestored(true);
      return;
    }
    const fincasLocal = loadFincasLocal(user.id).filter((x) => x.estado !== 'PENDIENTE_ELIMINACION');
    if (saved.fincaSel && fincasLocal.some((f) => f.id === saved.fincaSel)) {
      setFincaSel(saved.fincaSel);
      setLotes(loadLotesLocal(user.id, saved.fincaSel).filter((l) => l.estado !== 'ELIMINADO'));
      setNodos(loadNodosLocal(user.id, saved.fincaSel).filter((n) => n.estado !== 'INACTIVO'));
      if (saved.loteSel) {
        const lotesLocal = loadLotesLocal(user.id, saved.fincaSel).filter((l) => l.estado !== 'ELIMINADO');
        if (lotesLocal.some((l) => l.id === saved.loteSel)) {
          setLoteSel(saved.loteSel);
        }
      }
      if (saved.muestraSel) setMuestraSel(saved.muestraSel);
      if (saved.activeTab && saved.activeTab !== 'fincas') setActiveTab(saved.activeTab);
    }
    setWorkflowRestored(true);
  }, [user?.id, workflowRestored]);

  useEffect(() => {
    const st = location.state;
    if (!st || !user?.id) return;
    if (st.activeTab) setActiveTab(st.activeTab);
  }, [location.state?.activeTab, user?.id]);

  // Guardar selección actual para recuperarla después
  useEffect(() => {
    if (!user?.id || !workflowRestored) return;
    saveWorkflowState(user.id, { fincaSel, loteSel, muestraSel, activeTab });
  }, [user?.id, fincaSel, loteSel, muestraSel, activeTab, workflowRestored]);

  useEffect(() => {
    fincasApiDisponible().then(setApiOnline);
  }, []);

  useEffect(() => {
    const st = location.state;
    if (!st || !user?.id) return;
    if (st.fincaId) {
      setFincaSel(st.fincaId);
      setLotes(loadLotesLocal(user.id, st.fincaId).filter((l) => l.estado !== 'ELIMINADO'));
      setNodos(loadNodosLocal(user.id, st.fincaId).filter((n) => n.estado !== 'INACTIVO'));
      setActiveTab(st.loteId ? 'muestras' : 'lotes');
    }
    if (st.loteId) setLoteSel(st.loteId);
    window.history.replaceState({}, document.title);
  }, [location.state, user?.id]);

  useEffect(() => {
    if (!user?.id || !fincaSel || !loteSel) {
      setMuestras([]);
      setDiagnosticos([]);
      return;
    }
    const localM = loadMuestrasLocal(user.id, fincaSel, loteSel);
    const localD = loadDiagnosticosLocal(user.id, fincaSel, loteSel);
    setDiagnosticos(localD);

    if (apiOnline === false) {
      setMuestras(localM);
      return;
    }
    
    const fetchBackend = async () => {
      try {
        const [apiM, reporte] = await Promise.all([
          listarMuestras(fincaSel, loteSel),
          generarReporteLote(fincaSel, loteSel).catch(() => null)
        ]);
        
        const merged = [...apiM];
        localM.forEach((lm) => {
          if (!merged.find((m) => m.id === lm.id)) merged.push(lm);
        });
        setMuestras(merged);

        if (reporte && reporte.muestras) {
          const syncD = [...localD];
          reporte.muestras.forEach((rm) => {
            if (rm.diagnosticoID) {
              const idx = syncD.findIndex(d => d.muestraID === rm.id);
              const API_BASE = import.meta.env.VITE_YOLO_API_URL || (import.meta.env.PROD ? 'https://i2u6hsbhf1.execute-api.us-east-2.amazonaws.com' : '/yolo');
              const newD = {
                id: rm.diagnosticoID,
                muestraID: rm.id,
                estado: rm.estadoDiagnostico,
                origen: 'iot', // Debe ser 'iot' para que la UI muestre 📷 Cámara IoT
                tieneClorosis: rm.tieneClorosis,
                yolo: {
                  image: rm.imageURL ? (rm.imageURL.startsWith('http') ? rm.imageURL : `${API_BASE}/${rm.imageURL}`) : null
                }
              };
              if (idx >= 0) {
                syncD[idx] = { ...syncD[idx], ...newD };
              } else {
                syncD.push(newD);
              }
              // Hacer persistente el diagnóstico sincronizado para que 'Ver Diagnóstico' funcione
              saveDiagnosticoLocal(user.id, fincaSel, loteSel, newD);
            }
          });
          setDiagnosticos(syncD);
        }
      } catch (err) {
        setMuestras(localM);
      }
    };
    
    fetchBackend();
    
    // Auto-polling para ver nuevas muestras (webhook de YOLO)
    const interval = setInterval(fetchBackend, 5000);
    return () => clearInterval(interval);
  }, [user?.id, fincaSel, loteSel, apiOnline]);

  const showMsg = (texto, tipo = 'exito') => {
    setMsg({ texto, tipo });
    setTimeout(() => setMsg(null), 6000);
  };

  const seleccionarFinca = (id) => {
    setFincaSel(id);
    setLoteSel('');
    setNodoSel('');
    setMuestraSel('');
    setReporte(null);
    setLotes(loadLotesLocal(user.id, id).filter((l) => l.estado !== 'ELIMINADO'));
    setNodos(loadNodosLocal(user.id, id).filter((n) => n.estado !== 'INACTIVO'));
  };

  const seleccionarLote = (id) => {
    setLoteSel(id);
    setMuestraSel('');
    setReporte(null);
  };

  const trasCrearFinca = (fincaId) => {
    seleccionarFinca(fincaId);
    setActiveTab('lotes');
  };

  const trasCrearLote = (loteId) => {
    seleccionarLote(loteId);
    setActiveTab('muestras');
  };

  // Auto-seleccionar finca/lote si solo hay uno (facilita continuar el flujo)
  useEffect(() => {
    if (!user?.id || fincas.length !== 1) return;
    const id = fincas[0].id;
    if (!fincaSel) setFincaSel(id);
    const lotesLocal = loadLotesLocal(user.id, id).filter((l) => l.estado !== 'ELIMINADO');
    setLotes(lotesLocal);
    if (lotesLocal.length === 0) {
      setActiveTab('lotes');
    } else if (!loteSel) {
      setLoteSel(lotesLocal[0].id);
      setActiveTab('muestras');
    }
  }, [user?.id, fincas, fincaSel, loteSel]);

  const handleCrearFinca = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const data = await registrarFinca(formFinca);
      saveFincaLocal(user.id, data);
      setFincas(loadFincasLocal(user.id));
      setFormFinca({ nombre: '', ubicacion: '', descripcion: '' });
      setApiOnline(true);
      trasCrearFinca(data.id);
      showMsg('Finca registrada. Ahora agrega un lote.');
    } catch (err) {
      if (isServicioFincasNoDisponible(err)) {
        setApiOnline(false);
        const local = {
          id: generateUUID(),
          ...formFinca,
          estado: 'ACTIVA',
          createdAt: new Date().toISOString(),
          _offline: true,
        };
        saveFincaLocal(user.id, local);
        setFincas(loadFincasLocal(user.id));
        setFormFinca({ nombre: '', ubicacion: '', descripcion: '' });
        trasCrearFinca(local.id);
        showMsg('Finca guardada. Siguiente paso: agregar un lote.', 'warn');
      } else {
        showMsg(parseErrorFincas(err), 'error');
      }
    } finally { setLoading(false); }
  };

  const handleDesactivarFinca = async (id) => {
    if (!confirm('¿Desactivar esta finca?')) return;
    setLoading(true);
    try {
      const finca = fincas.find((f) => f.id === id);
      if (finca?._offline || apiOnline === false) {
        updateFincaLocal(user.id, id, { estado: 'PENDIENTE_ELIMINACION' });
        setFincas(loadFincasLocal(user.id));
        if (fincaSel === id) { setFincaSel(''); setLoteSel(''); }
        showMsg('Finca desactivada (local)');
        return;
      }
      const data = await desactivarFinca(id, { confirmar: true });
      updateFincaLocal(user.id, id, { estado: data.estado });
      setFincas(loadFincasLocal(user.id));
      showMsg('Finca desactivada');
    } catch (err) {
      if (isServicioFincasNoDisponible(err)) {
        setApiOnline(false);
        updateFincaLocal(user.id, id, { estado: 'PENDIENTE_ELIMINACION' });
        setFincas(loadFincasLocal(user.id));
        showMsg('Finca desactivada localmente', 'warn');
      } else {
        showMsg(parseErrorFincas(err), 'error');
      }
    } finally { setLoading(false); }
  };

  const handleCrearLote = async (e) => {
    e.preventDefault();
    if (!fincaSel) return;
    setLoading(true);
    try {
      const finca = fincas.find((f) => f.id === fincaSel);
      const crearLocal = () => {
        const local = {
          id: generateUUID(),
          fincaID: fincaSel,
          nombre: formLote.nombre,
          area: parseFloat(formLote.area),
          descripcion: formLote.descripcion,
          estado: 'ACTIVO',
          createdAt: new Date().toISOString(),
          _offline: true,
        };
        saveLoteLocal(user.id, fincaSel, local);
        setLotes(loadLotesLocal(user.id, fincaSel).filter((l) => l.estado !== 'ELIMINADO'));
        setFormLote({ nombre: '', area: '', descripcion: '' });
        return local;
      };

      if (finca?._offline || apiOnline === false) {
        const local = crearLocal();
        trasCrearLote(local.id);
        showMsg('Lote guardado. Ve a Muestras para analizar una imagen.', 'warn');
        return;
      }
      const data = await agregarLote(fincaSel, { ...formLote, area: parseFloat(formLote.area) });
      saveLoteLocal(user.id, fincaSel, data);
      setLotes(loadLotesLocal(user.id, fincaSel).filter((l) => l.estado !== 'ELIMINADO'));
      setFormLote({ nombre: '', area: '', descripcion: '' });
      trasCrearLote(data.id);
      showMsg('Lote agregado. Sube una imagen en la pestaña Muestras.');
    } catch (err) {
      if (isServicioFincasNoDisponible(err)) {
        setApiOnline(false);
        const local = {
          id: generateUUID(),
          fincaID: fincaSel,
          nombre: formLote.nombre,
          area: parseFloat(formLote.area),
          descripcion: formLote.descripcion,
          estado: 'ACTIVO',
          createdAt: new Date().toISOString(),
          _offline: true,
        };
        saveLoteLocal(user.id, fincaSel, local);
        setLotes(loadLotesLocal(user.id, fincaSel).filter((l) => l.estado !== 'ELIMINADO'));
        setFormLote({ nombre: '', area: '', descripcion: '' });
        trasCrearLote(local.id);
        showMsg('Lote guardado. Ve a Muestras para analizar una imagen.', 'warn');
      } else {
        showMsg(parseErrorFincas(err), 'error');
      }
    } finally { setLoading(false); }
  };

  const handleEliminarLote = async (id) => {
    if (!confirm('¿Eliminar este lote?')) return;
    setLoading(true);
    try {
      const lote = lotes.find((l) => l.id === id);
      if (lote?._offline || apiOnline === false) {
        updateLoteLocal(user.id, fincaSel, id, { estado: 'ELIMINADO' });
        setLotes(loadLotesLocal(user.id, fincaSel).filter((l) => l.estado !== 'ELIMINADO'));
        if (loteSel === id) setLoteSel('');
        showMsg('Lote eliminado (local)');
        return;
      }
      await eliminarLote(id);
      updateLoteLocal(user.id, fincaSel, id, { estado: 'ELIMINADO' });
      setLotes(loadLotesLocal(user.id, fincaSel).filter((l) => l.estado !== 'ELIMINADO'));
      showMsg('Lote eliminado');
    } catch (err) {
      showMsg(parseErrorFincas(err), 'error');
    } finally { setLoading(false); }
  };

  const handleCrearNodo = async (e) => {
    e.preventDefault();
    if (!fincaSel) return;
    setLoading(true);
    try {
      const finca = fincas.find((f) => f.id === fincaSel);
      const crearLocal = () => {
        const local = {
          id: generateUUID(),
          fincaID: fincaSel,
          loteID: formNodo.lote_id || null,
          nombre: formNodo.nombre,
          nodeKey: formNodo.node_key,
          estado: 'ACTIVO',
          createdAt: new Date().toISOString(),
          _offline: true,
        };
        saveNodoLocal(user.id, fincaSel, local);
        setNodos(loadNodosLocal(user.id, fincaSel).filter((n) => n.estado !== 'INACTIVO'));
        setFormNodo({ nombre: '', node_key: '', lote_id: '' });
        return local;
      };

      if (finca?._offline || apiOnline === false) {
        crearLocal();
        showMsg('Cámara IoT guardada (modo offline).', 'warn');
        return;
      }
      const data = await registrarNodo(fincaSel, {
        nombre: formNodo.nombre,
        node_key: formNodo.node_key,
        lote_id: formNodo.lote_id || undefined,
      });
      saveNodoLocal(user.id, fincaSel, data);
      setNodos(loadNodosLocal(user.id, fincaSel).filter((n) => n.estado !== 'INACTIVO'));
      setFormNodo({ nombre: '', node_key: '', lote_id: '' });
      showMsg('Cámara IoT agregada exitosamente.');
    } catch (err) {
      if (isServicioFincasNoDisponible(err)) {
        setApiOnline(false);
        const local = {
          id: generateUUID(),
          fincaID: fincaSel,
          loteID: formNodo.lote_id || null,
          nombre: formNodo.nombre,
          nodeKey: formNodo.node_key,
          estado: 'ACTIVO',
          createdAt: new Date().toISOString(),
          _offline: true,
        };
        saveNodoLocal(user.id, fincaSel, local);
        setNodos(loadNodosLocal(user.id, fincaSel).filter((n) => n.estado !== 'INACTIVO'));
        setFormNodo({ nombre: '', node_key: '', lote_id: '' });
        showMsg('Cámara IoT guardada localmente.', 'warn');
      } else {
        showMsg(parseErrorFincas(err), 'error');
      }
    } finally { setLoading(false); }
  };

  const handleDesactivarNodo = async (id) => {
    if (!confirm('¿Desactivar esta cámara?')) return;
    setLoading(true);
    try {
      const nodo = nodos.find((n) => n.id === id);
      if (nodo?._offline || apiOnline === false) {
        updateNodoLocal(user.id, fincaSel, id, { estado: 'INACTIVO' });
        setNodos(loadNodosLocal(user.id, fincaSel).filter((n) => n.estado !== 'INACTIVO'));
        if (nodoSel === id) setNodoSel('');
        showMsg('Cámara desactivada (local)');
        return;
      }
      await desactivarNodo(id);
      updateNodoLocal(user.id, fincaSel, id, { estado: 'INACTIVO' });
      setNodos(loadNodosLocal(user.id, fincaSel).filter((n) => n.estado !== 'INACTIVO'));
      if (nodoSel === id) setNodoSel('');
      showMsg('Cámara desactivada exitosamente');
    } catch (err) {
      showMsg(parseErrorFincas(err), 'error');
    } finally { setLoading(false); }
  };


  const handleReporte = async () => {
    if (!fincaSel || !loteSel || !loteActivo) return;
    setLoading(true);
    setReporte(null);
    try {
      if (apiOnline !== false && !loteActivo._offline) {
        const data = await generarReporteLote(fincaSel, loteSel);
        setReporte({
          ...data,
          finca: fincaActiva ? { nombre: fincaActiva.nombre, ubicacion: fincaActiva.ubicacion } : data.finca,
          generadoEn: new Date().toISOString(),
        });
        showMsg('Reporte generado desde el servidor');
      } else {
        throw new Error('offline');
      }
    } catch (err) {
      const localM = loadMuestrasLocal(user.id, fincaSel, loteSel);
      const localD = loadDiagnosticosLocal(user.id, fincaSel, loteSel);
      const rep = buildReporteLocal({
        finca: fincaActiva,
        lote: loteActivo,
        muestras: localM,
        diagnosticos: localD,
      });
      setReporte(rep);
      showMsg(apiOnline === false || err.message === 'offline'
        ? 'Reporte generado localmente (datos YOLO vinculados)'
        : parseErrorFincas(err), apiOnline === false ? 'warn' : 'error');
    } finally { setLoading(false); }
  };

  const handleAceptarDiag = async () => {
    if (!diagActivo || !fincaSel || !loteSel) return;
    setLoading(true);
    try {
      if (diagActivo.diagnosticoID && apiOnline !== false) {
        await aceptarDiagnostico(diagActivo.diagnosticoID);
      }
      updateDiagnosticoLocal(user.id, fincaSel, loteSel, diagActivo.id, { estado: 'ACEPTADO' });
      setDiagnosticos(loadDiagnosticosLocal(user.id, fincaSel, loteSel));
      showMsg('Diagnóstico aceptado');
      handleReporte();
    } catch (err) {
      if (isServicioFincasNoDisponible(err)) {
        updateDiagnosticoLocal(user.id, fincaSel, loteSel, diagActivo.id, { estado: 'ACEPTADO' });
        setDiagnosticos(loadDiagnosticosLocal(user.id, fincaSel, loteSel));
        showMsg('Diagnóstico aceptado localmente', 'warn');
      } else {
        showMsg(parseErrorFincas(err), 'error');
      }
    } finally { setLoading(false); }
  };

  const handleRechazarDiag = async () => {
    if (!diagActivo || !fincaSel || !loteSel) return;
    setLoading(true);
    try {
      if (diagActivo.diagnosticoID && apiOnline !== false) {
        await rechazarDiagnostico(diagActivo.diagnosticoID, { motivo: motivoRechazo });
      }
      updateDiagnosticoLocal(user.id, fincaSel, loteSel, diagActivo.id, {
        estado: 'RECHAZADO',
        motivoRechazo,
      });
      setDiagnosticos(loadDiagnosticosLocal(user.id, fincaSel, loteSel));
      setMotivoRechazo('');
      showMsg('Diagnóstico rechazado');
    } catch (err) {
      if (isServicioFincasNoDisponible(err)) {
        updateDiagnosticoLocal(user.id, fincaSel, loteSel, diagActivo.id, {
          estado: 'RECHAZADO',
          motivoRechazo,
        });
        setDiagnosticos(loadDiagnosticosLocal(user.id, fincaSel, loteSel));
        showMsg('Diagnóstico rechazado localmente', 'warn');
      } else {
        showMsg(parseErrorFincas(err), 'error');
      }
    } finally { setLoading(false); }
  };

  const handleVincularYoloAMuestra = async () => {
    if (!muestraSel || !yoloPendiente?.feedback || yoloPendiente.feedback.label === 'Error') return;
    const diag = {
      id: generateUUID(),
      muestraID: muestraSel,
      estado: 'PENDIENTE',
      origen: 'yolo',
      yolo: {
        feedback: yoloPendiente.feedback,
        num_detections: yoloPendiente.num_detections,
        avg_confidence: yoloPendiente.avg_confidence,
        detections: yoloPendiente.detections || [],
        image: yoloPendiente.image,
      },
      tieneClorosis: tieneClorosisFromYolo(yoloPendiente),
      createdAt: new Date().toISOString(),
      _offline: true,
    };
    saveDiagnosticoLocal(user.id, fincaSel, loteSel, diag);
    setDiagnosticos(loadDiagnosticosLocal(user.id, fincaSel, loteSel));
    if (activeHistorialId) {
      marcarHistorialVinculado(activeHistorialId, {
        fincaId: fincaSel,
        loteId: loteSel,
        muestraId: muestraSel,
        diagnosticoId: diag.id,
        fincaNombre: fincaActiva?.nombre,
        loteNombre: loteActivo?.nombre,
      });
    }
    showMsg(`Diagnóstico YOLO vinculado a la muestra. Nitrógeno: ${yoloPendiente.feedback.label}`);
  };

  const handleDescargarPDF = () => {
    if (!reporte) {
      showMsg('Primero genera el reporte', 'warn');
      return;
    }
    try {
      const nombreUsuario = user?.nombre
        ? `${user.nombre} ${user.apellido || ''}`.trim()
        : user?.email;
      const filename = generateReportePDF(reporte, { usuario: nombreUsuario });
      showMsg(`PDF descargado: ${filename}`);
    } catch (err) {
      showMsg(`Error al generar PDF: ${err.message}`, 'error');
    }
  };

  const refrescarMuestrasYDiag = () => {
    setMuestras(loadMuestrasLocal(user.id, fincaSel, loteSel));
    setDiagnosticos(loadDiagnosticosLocal(user.id, fincaSel, loteSel));
  };

  const getDiagEstado = (muestraId) => {
    const d = diagnosticos.find((x) => x.muestraID === muestraId);
    return d?.estado || 'SIN_DIAGNOSTICO';
  };

  const getDiagOrigen = (muestraId) => {
    const d = diagnosticos.find((x) => x.muestraID === muestraId);
    return d?.origen || 'manual';
  };

  return (
    <Layout title="Mis Fincas" subtitle="Gestiona fincas, lotes, muestras y reportes.">
      <ServiceStatusBar compact />

      {fincas.length > 0 && (
        <div className="stat-grid">
          <StatCard icon={<IconFarm />} label="Fincas" value={globalStats.fincas} accent="green" />
          <StatCard icon={<IconGrid />} label="Lotes" value={globalStats.lotes} accent="earth" />
          <StatCard icon={<IconSample />} label="Muestras" value={globalStats.muestras} accent="blue" />
          <StatCard
            icon={<IconClipboard />}
            label="Diagnósticos"
            value={globalStats.diagnosticos}
            sub={globalStats.pendientes > 0 ? `${globalStats.pendientes} pendientes` : undefined}
            accent="amber"
            onClick={globalStats.pendientes > 0 && fincaSel && loteSel ? () => irATab('diagnosticos') : undefined}
          />
        </div>
      )}

      {apiOnline === false && (
        <div style={{
          marginBottom: '1rem', padding: '0.85rem 1rem', borderRadius: '0.5rem', fontSize: '0.9rem',
          background: '#fffbeb', color: '#92400e', border: '1px solid #fcd34d',
        }}>
          <strong>Modo local:</strong> fincas no responde en <strong>:8082</strong>.
          Los diagnósticos YOLO se guardan localmente y el reporte se calcula en el navegador.
        </div>
      )}

      <div className="fincas-steps">
        {TABS.map((t) => {
          const habilitado = tabHabilitado(t.id);
          const activo = activeTab === t.id;
          const indice = TABS.findIndex((x) => x.id === t.id);
          const completado = indice < TABS.findIndex((x) => x.id === activeTab);
          return (
            <button
              key={t.id}
              type="button"
              className={`fincas-step fincas-step--clickable ${activo ? 'fincas-step--active' : ''} ${completado ? 'fincas-step--done' : ''} ${!habilitado ? 'fincas-step--disabled' : ''}`}
              onClick={() => irATab(t.id)}
              disabled={!habilitado}
              title={!habilitado ? 'Completa el paso anterior primero' : ''}
            >
              {completado && !activo ? <IconCheck size={14} className="fincas-step__check" /> : null}{t.label}
            </button>
          );
        })}
      </div>

      {fincas.length > 0 && !fincaSel && activeTab === 'fincas' && (
        <div className="fincas-guide">
          <strong>Paso siguiente:</strong> haz clic en el nombre de tu finca en la tabla de abajo (fila verde) o usa el botón <strong>Seleccionar</strong> para continuar a Lotes.
        </div>
      )}

      {fincaSel && !loteSel && activeTab === 'lotes' && lotes.length === 0 && (
        <div className="fincas-guide">
          <strong>Finca «{fincaActiva?.nombre}» seleccionada.</strong> Agrega al menos un lote con nombre y área (hectáreas) para poder subir imágenes.
        </div>
      )}

      {fincaSel && activeTab === 'muestras' && (
        <div className="fincas-guide">
          <strong>{loteActivo ? `Lote «${loteActivo.nombre}» listo.` : `Finca «${fincas.find(f => f.id === fincaSel)?.nombre}» lista.`}</strong> Sube una o varias fotos de hoja, usa GPS y presiona <strong>Analizar</strong>. Puedes agregar más imágenes después de cada análisis.
        </div>
      )}

      {fincaSel && (
        <div className="fincas-breadcrumb">
          Finca: <strong>{fincaActiva?.nombre}</strong>
          {loteSel && <> → Lote: <strong>{loteActivo?.nombre}</strong></>}
          {muestraSel && <> → Muestra seleccionada</>}
        </div>
      )}

      {msg && (
        <div style={{
          padding: '0.75rem 1rem', marginBottom: '1rem', borderRadius: '0.5rem', fontSize: '0.9rem',
          background: msg.tipo === 'error' ? '#fef2f2' : msg.tipo === 'warn' ? '#fffbeb' : '#f0fdf4',
          color: msg.tipo === 'error' ? '#991b1b' : msg.tipo === 'warn' ? '#92400e' : '#166534',
          border: `1px solid ${msg.tipo === 'error' ? '#fecaca' : msg.tipo === 'warn' ? '#fcd34d' : '#bbf7d0'}`,
        }}>{msg.texto}</div>
      )}

      {/* ── Tab: Fincas ── */}
      {activeTab === 'fincas' && (
      <>
      <div className="admin-card" style={{ marginBottom: '1.5rem' }}>
        <h2 className="admin-card__title">1. Registrar finca</h2>
        <form onSubmit={handleCrearFinca} style={{ display: 'grid', gap: '0.75rem', maxWidth: 480 }}>
          <input className="form-input" placeholder="Nombre" value={formFinca.nombre} onChange={(e) => setFormFinca((f) => ({ ...f, nombre: e.target.value }))} required />
          <input className="form-input" placeholder="Ubicación" value={formFinca.ubicacion} onChange={(e) => setFormFinca((f) => ({ ...f, ubicacion: e.target.value }))} required />
          <input className="form-input" placeholder="Descripción" value={formFinca.descripcion} onChange={(e) => setFormFinca((f) => ({ ...f, descripcion: e.target.value }))} />
          <button type="submit" className="btn-add" disabled={loading}>+ Crear finca</button>
        </form>
      </div>

      <div className="admin-card" style={{ marginBottom: '1.5rem' }}>
        <h2 className="admin-card__title">Fincas ({fincas.length})</h2>
        {fincas.length === 0 ? <p className="admin-empty">No hay fincas. Crea una arriba.</p> : (
          <table className="admin-table">
            <thead><tr><th>Nombre</th><th>Ubicación</th><th>Estado</th><th>Acciones</th></tr></thead>
            <tbody>
              {fincas.map((f) => (
                <tr key={f.id} style={fincaSel === f.id ? { background: '#f0fdf4' } : undefined}>
                  <td>
                    <strong>{f.nombre}{f._offline ? ' (local)' : ''}</strong>
                  </td>
                  <td>{f.ubicacion}</td>
                  <td>{f.estado}</td>
                  <td style={{ display: 'flex', gap: '0.35rem', flexWrap: 'wrap' }}>
                    <button
                      type="button"
                      className={`btn-admin btn-with-icon ${fincaSel === f.id ? 'btn-admin--primary' : ''}`}
                      onClick={() => { seleccionarFinca(f.id); irATab('lotes'); }}
                    >
                      {fincaSel === f.id ? <><IconCheck size={14} /> Seleccionada</> : 'Seleccionar'}
                    </button>
                    {f.estado === 'ACTIVA' && (
                      <button type="button" className="btn-admin btn-admin--danger" onClick={() => handleDesactivarFinca(f.id)}>Desactivar</button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        {fincaSel && (
          <div className="fincas-siguiente">
            <button type="button" className="btn-add" onClick={() => irATab('lotes')}>
              Continuar a Lotes →
            </button>
          </div>
        )}
      </div>
      </>
      )}

      {/* ── Tab: Lotes ── */}
      {activeTab === 'lotes' && fincaSel && (
        <div className="admin-card" style={{ marginBottom: '1.5rem' }}>
          <h2 className="admin-card__title">2. Lotes — {fincaActiva?.nombre}</h2>
          <form onSubmit={handleCrearLote} style={{ display: 'grid', gap: '0.75rem', maxWidth: 480, marginBottom: '1rem' }}>
            <input className="form-input" placeholder="Nombre del lote" value={formLote.nombre} onChange={(e) => setFormLote((f) => ({ ...f, nombre: e.target.value }))} required />
            <input className="form-input" type="number" step="0.01" placeholder="Área (ha)" value={formLote.area} onChange={(e) => setFormLote((f) => ({ ...f, area: e.target.value }))} required />
            <input className="form-input" placeholder="Descripción" value={formLote.descripcion} onChange={(e) => setFormLote((f) => ({ ...f, descripcion: e.target.value }))} />
            <button type="submit" className="btn-add" disabled={loading}>+ Agregar lote</button>
          </form>
          {lotes.length === 0 ? (
            <p className="admin-empty">Agrega un lote para tomar muestras.</p>
          ) : (
            <table className="admin-table">
              <thead><tr><th>Nombre</th><th>Área (ha)</th><th>Estado</th><th>Acciones</th></tr></thead>
              <tbody>
                {lotes.map((l) => (
                  <tr key={l.id} style={loteSel === l.id ? { background: '#f0fdf4' } : undefined}>
                    <td><strong>{l.nombre}</strong></td>
                    <td>{l.area}</td>
                    <td>{l.estado}</td>
                    <td style={{ display: 'flex', gap: '0.35rem', flexWrap: 'wrap' }}>
                      <button
                        type="button"
                        className={`btn-admin btn-with-icon ${loteSel === l.id ? 'btn-admin--primary' : ''}`}
                        onClick={() => { seleccionarLote(l.id); irATab('muestras'); }}
                      >
                        {loteSel === l.id ? <><IconCheck size={14} /> Seleccionado</> : 'Seleccionar'}
                      </button>
                      {l.estado === 'ACTIVO' && (
                        <button type="button" className="btn-admin btn-admin--danger" onClick={() => handleEliminarLote(l.id)}>Eliminar</button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          {loteSel && (
            <div className="fincas-siguiente">
              <button type="button" className="btn-add" onClick={() => irATab('muestras')}>
                Continuar a Muestras (subir imagen) →
              </button>
            </div>
          )}
        </div>
      )}

      {/* ── Tab: Nodos (Cámaras) ── */}
      {activeTab === 'nodos' && fincaSel && (
        <div className="admin-card" style={{ marginBottom: '1.5rem' }}>
          <h2 className="admin-card__title">Cámaras IoT — {fincaActiva?.nombre}</h2>
          <form onSubmit={handleCrearNodo} style={{ display: 'grid', gap: '0.75rem', maxWidth: 480, marginBottom: '1rem' }}>
            <input className="form-input" placeholder="Nombre de la cámara" value={formNodo.nombre} onChange={(e) => setFormNodo((f) => ({ ...f, nombre: e.target.value }))} required />
            <input className="form-input" placeholder="Node Key (ej: cam-001)" value={formNodo.node_key} onChange={(e) => setFormNodo((f) => ({ ...f, node_key: e.target.value }))} required />
            <select className="form-input" value={formNodo.lote_id} onChange={(e) => setFormNodo((f) => ({ ...f, lote_id: e.target.value }))} required>
              <option value="">-- Selecciona un lote (Requerido) --</option>
              {lotes.filter((l) => l.estado === 'ACTIVO').map((l) => (
                <option key={l.id} value={l.id}>{l.nombre}</option>
              ))}
            </select>
            <button type="submit" className="btn-add" disabled={loading}>+ Agregar cámara</button>
          </form>
          {nodos.length === 0 ? (
            <p className="admin-empty">Aún no has agregado ninguna cámara IoT.</p>
          ) : (
            <table className="admin-table">
              <thead><tr><th>Nombre</th><th>Node Key</th><th>Lote</th><th>Estado</th><th>Acciones</th></tr></thead>
              <tbody>
                {nodos.map((n) => (
                  <tr key={n.id} style={nodoSel === n.id ? { background: '#f0fdf4' } : undefined}>
                    <td><strong>{n.nombre}</strong></td>
                    <td>{n.nodeKey || n.node_key}</td>
                    <td>{lotes.find((l) => l.id === (n.loteID || n.lote_id))?.nombre || 'Ninguno'}</td>
                    <td>{n.estado}</td>
                    <td style={{ display: 'flex', gap: '0.35rem', flexWrap: 'wrap' }}>
                      <button
                        type="button"
                        className={`btn-admin btn-with-icon ${nodoSel === n.id ? 'btn-admin--primary' : ''}`}
                        onClick={() => { setNodoSel(n.id); }}
                      >
                        {nodoSel === n.id ? <><IconCheck size={14} /> Seleccionada</> : 'Seleccionar'}
                      </button>
                      {n.estado === 'ACTIVO' && (
                        <button type="button" className="btn-admin btn-admin--danger" onClick={() => handleDesactivarNodo(n.id)}>Desactivar</button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}

      {/* ── Tab: Muestras ── */}
      {activeTab === 'muestras' && fincaSel && (
        <div className="admin-card" style={{ marginBottom: '1.5rem' }}>
          <h2 className="admin-card__title">Muestras — {loteActivo ? loteActivo.nombre : 'General de la Finca'}</h2>

          <AnalizarMuestraEnFinca
            userId={user.id}
            fincaId={fincaSel}
            loteId={loteSel}
            apiOnline={apiOnline}
            onCompletado={({ muestraId, label, esUltima, totalEnLote }) => {
              refrescarMuestrasYDiag();
              setMuestraSel(muestraId);
              if (esUltima) {
                showMsg(
                  totalEnLote > 1
                    ? `${totalEnLote} muestras analizadas. Última: Nitrógeno ${label}. Puedes subir más imágenes o ir a Diagnósticos.`
                    : `Muestra analizada. Nitrógeno: ${label}. Puedes subir otra imagen o ir a Diagnósticos.`,
                );
              }
            }}
            onError={(texto) => showMsg(texto, 'error')}
          />

          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', margin: '1.5rem 0 0.75rem' }}>
            <h3 style={{ fontSize: '0.95rem', margin: 0 }}>Muestras registradas ({muestras.length})</h3>
            <span style={{ fontSize: '0.75rem', color: '#6b7280' }}>Actualizando en vivo...</span>
          </div>
          {muestras.length === 0 ? (
            <p className="admin-empty">Aún no hay muestras. Usa el formulario de arriba para subir una imagen o espera que la cámara IoT envíe una inferencia.</p>
          ) : (
            <table className="admin-table">
              <thead>
                <tr><th>Coordenadas</th><th>Fecha</th><th>Origen</th><th>Diagnóstico</th><th>Acciones</th></tr>
              </thead>
              <tbody>
                {muestras.map((m) => {
                  const estado = getDiagEstado(m.id);
                  const origen = getDiagOrigen(m.id);
                  return (
                    <tr key={m.id} className={muestraSel === m.id ? 'muestra-row--con-diag' : ''}>
                      <td>{m.latitud?.toFixed?.(4) ?? m.latitud}, {m.longitud?.toFixed?.(4) ?? m.longitud}</td>
                      <td>{new Date(m.createdAt).toLocaleString()}</td>
                      <td>
                        {origen === 'iot' ? (
                          <span style={{ background: '#dbeafe', color: '#1e40af', padding: '0.2rem 0.5rem', borderRadius: '4px', fontSize: '0.8rem', fontWeight: 600 }}>📷 Cámara IoT</span>
                        ) : (
                          <span style={{ background: '#f3f4f6', color: '#374151', padding: '0.2rem 0.5rem', borderRadius: '4px', fontSize: '0.8rem', fontWeight: 600 }}>👤 Manual</span>
                        )}
                      </td>
                      <td><DiagBadge estado={estado} /></td>
                      <td>
                        <button
                          type="button"
                          className="btn-admin btn-admin--primary"
                          onClick={() => { setMuestraSel(m.id); irATab('diagnosticos'); }}
                        >
                          Ver diagnóstico
                        </button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          )}
          {muestras.length > 0 && (
            <div className="fincas-siguiente">
              <button type="button" className="btn-add" onClick={() => irATab('diagnosticos')}>
                Ir a Diagnósticos →
              </button>
            </div>
          )}
        </div>
      )}

      {/* ── Tab: Diagnósticos ── */}
      {activeTab === 'diagnosticos' && fincaSel && loteSel && (
        <div className="admin-card" style={{ marginBottom: '1.5rem' }}>
          <h2 className="admin-card__title">Diagnósticos — {loteActivo?.nombre}</h2>

          {!muestraSel && muestras.length > 0 && (
            <p className="admin-empty" style={{ marginBottom: '1rem' }}>
              Selecciona una muestra en la pestaña <button type="button" className="vincular-panel__link" onClick={() => irATab('muestras')}>Muestras</button> o elige una de la lista:
            </p>
          )}

          {!muestraSel && muestras.length > 0 && (
            <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap', marginBottom: '1rem' }}>
              {muestras.map((m) => (
                <button
                  key={m.id}
                  type="button"
                  className="btn-admin"
                  onClick={() => setMuestraSel(m.id)}
                >
                  {new Date(m.createdAt).toLocaleDateString()} — <DiagBadge estado={getDiagEstado(m.id)} />
                </button>
              ))}
            </div>
          )}

          {yoloPendiente?.feedback && yoloPendiente.feedback.label !== 'Error' && !diagActivo && (
            <div style={{ marginBottom: '1rem', padding: '1rem', background: '#eff6ff', borderRadius: '0.5rem', border: '1px solid #93c5fd' }}>
              <p style={{ margin: '0 0 0.5rem', fontSize: '0.9rem' }}>
                Tienes un diagnóstico YOLO pendiente: <strong>Nitrógeno {yoloPendiente.feedback.label}</strong>
              </p>
              {muestraSel ? (
                <button type="button" className="btn-add" onClick={handleVincularYoloAMuestra}>Vincular a muestra seleccionada</button>
              ) : (
                <>
                  <p style={{ margin: '0 0 0.75rem', fontSize: '0.85rem', color: '#6b7280' }}>
                    Selecciona una muestra arriba o crea una nueva con GPS:
                  </p>
                  <VincularDiagnosticoPanel
                    yoloResults={yoloPendiente}
                    historialId={activeHistorialId}
                    compact
                    onVinculado={(v) => {
                      setMuestras(loadMuestrasLocal(user.id, fincaSel, loteSel));
                      setDiagnosticos(loadDiagnosticosLocal(user.id, fincaSel, loteSel));
                      if (v?.muestraId) setMuestraSel(v.muestraId);
                    }}
                  />
                </>
              )}
            </div>
          )}

          {diagActivo && (
            <div style={{ display: 'grid', gap: '1rem' }}>
              <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center', flexWrap: 'wrap' }}>
                <DiagBadge estado={diagActivo.estado} />
                <span style={{ fontSize: '0.9rem' }}>
                  Nitrógeno <strong>{diagActivo.yolo?.feedback?.label || '—'}</strong>
                  {diagActivo.tieneClorosis != null && (
                    <> · Clorosis: <strong>{diagActivo.tieneClorosis ? 'Sí' : 'No'}</strong></>
                  )}
                </span>
              </div>

              {diagActivo.yolo?.image && (
                <img 
                  src={diagActivo.yolo.image} 
                  alt="Diagnóstico YOLO" 
                  style={{ maxWidth: 400, borderRadius: '0.5rem' }} 
                  onError={(e) => {
                    e.target.style.display = 'none';
                    if (!e.target.nextElementSibling) {
                      const msg = document.createElement('p');
                      msg.style.fontSize = '0.85rem';
                      msg.style.color = '#6b7280';
                      msg.style.margin = '0.5rem 0';
                      msg.innerHTML = '📸 <i>La imagen fue procesada en el servidor perimetral YOLO y no está disponible para previsualización web.</i>';
                      e.target.parentNode.insertBefore(msg, e.target.nextSibling);
                    }
                  }}
                />
              )}

              {diagActivo.yolo?.feedback?.recommendation && (
                <p style={{ fontSize: '0.9rem', color: '#4b5563', margin: 0 }}>
                  {diagActivo.yolo.feedback.recommendation}
                </p>
              )}

              {diagActivo.estado === 'PENDIENTE' && (
                <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap', alignItems: 'center' }}>
                  <input
                    className="form-input"
                    placeholder="Motivo de rechazo (opcional)"
                    value={motivoRechazo}
                    onChange={(e) => setMotivoRechazo(e.target.value)}
                    style={{ flex: 1, minWidth: 200 }}
                  />
                  <button type="button" className="btn-admin btn-admin--primary" onClick={handleAceptarDiag} disabled={loading}>Aceptar diagnóstico</button>
                  <button type="button" className="btn-admin btn-admin--danger" onClick={handleRechazarDiag} disabled={loading}>Rechazar</button>
                </div>
              )}
            </div>
          )}

          {!muestraSel && diagnosticos.length > 0 && (
            <p className="admin-empty" style={{ marginTop: '0.5rem' }}>
              {diagnosticos.length} diagnóstico(s) en este lote. Selecciona una muestra para ver detalle.
            </p>
          )}

          {!muestraSel && diagnosticos.length === 0 && !yoloPendiente?.feedback && (
            <p className="admin-empty">
              Primero sube una imagen en la pestaña <button type="button" className="vincular-panel__link" onClick={() => irATab('muestras')}>Muestras</button>.
            </p>
          )}

          {diagActivo?.estado === 'PENDIENTE' && (
            <div className="fincas-siguiente">
              <button type="button" className="btn-add" onClick={() => { handleAceptarDiag(); setActiveTab('reporte'); }}>
                Aceptar y ver reporte →
              </button>
            </div>
          )}
        </div>
      )}

      {/* ── Tab: Reporte ── */}
      {activeTab === 'reporte' && fincaSel && loteSel && (
        <div className="admin-card" style={{ marginBottom: '1.5rem' }}>
          <h2 className="admin-card__title">Reporte — {loteActivo?.nombre}</h2>
          <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap', marginBottom: '1rem' }}>
            <button type="button" className="btn-add" onClick={handleReporte} disabled={loading}>
              {loading ? 'Generando...' : 'Generar reporte'}
            </button>
            {reporte && (
              <button type="button" className="btn-admin btn-admin--primary" onClick={handleDescargarPDF}>
                Descargar PDF técnico
              </button>
            )}
          </div>

          {reporte && (
            <>
              {reporte._local && (
                <p style={{ fontSize: '0.85rem', color: '#92400e', marginBottom: '0.75rem' }}>
                  Reporte calculado localmente a partir de diagnósticos YOLO vinculados.
                </p>
              )}
              <div className="reporte-metricas">
                <div className="reporte-metrica">
                  <div className="reporte-metrica__valor">{reporte.metricas?.totalMuestras ?? 0}</div>
                  <div className="reporte-metrica__label">Muestras</div>
                </div>
                <div className="reporte-metrica">
                  <div className="reporte-metrica__valor">{reporte.metricas?.diagnosticosAceptados ?? 0}</div>
                  <div className="reporte-metrica__label">Aceptados</div>
                </div>
                <div className="reporte-metrica">
                  <div className="reporte-metrica__valor">{reporte.metricas?.diagnosticosPendientes ?? 0}</div>
                  <div className="reporte-metrica__label">Pendientes</div>
                </div>
                <div className="reporte-metrica">
                  <div className="reporte-metrica__valor">{(reporte.metricas?.porcentajeAfectado ?? 0).toFixed(1)}%</div>
                  <div className="reporte-metrica__label">Área afectada</div>
                </div>
              </div>

              <p style={{ fontSize: '0.9rem' }}>
                Lote <strong>{reporte.nombre}</strong> · Área total: <strong>{reporte.areaTotal} ha</strong>
                {reporte.finca?.nombre && <> · Finca: <strong>{reporte.finca.nombre}</strong></>}
              </p>

              <ReportCharts
                nitrogenData={reportCharts.nitrogenData}
                estadoData={reportCharts.estadoData}
                serieTemporal={reportCharts.serieTemporal}
              />

              <MuestrasMap points={mapPoints} />

              <table className="admin-table" style={{ marginTop: '1rem' }}>
                <thead>
                  <tr><th>Muestra</th><th>Diagnóstico</th><th>Estado</th><th>Clorosis</th><th>Nitrógeno</th></tr>
                </thead>
                <tbody>
                  {(reporte.muestras || []).map((m) => (
                    <tr key={m.id}>
                      <td style={{ fontSize: '0.75rem', fontFamily: 'monospace' }}>{m.id?.slice?.(0, 8)}…</td>
                      <td style={{ fontSize: '0.75rem' }}>{m.diagnosticoID?.slice?.(0, 8) || '—'}…</td>
                      <td><DiagBadge estado={m.estadoDiagnostico} /></td>
                      <td>{m.tieneClorosis == null ? '—' : m.tieneClorosis ? 'Sí' : 'No'}</td>
                      <td>{m.yoloLabel || '—'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </>
          )}

          {!reporte && !loading && (
            <div className="fincas-report-empty">
              <span className="fincas-report-empty__icon"><IconChart size={40} /></span>
              <h3>Genera tu reporte técnico</h3>
              <p>Consolida las muestras y diagnósticos del lote en métricas, gráficos, mapa GPS y PDF descargable.</p>
              <button type="button" className="btn-add" onClick={handleReporte}>
                Generar reporte ahora
              </button>
            </div>
          )}
        </div>
      )}
    </Layout>
  );
}
