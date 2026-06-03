import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../../context/AuthContext';
import {
  getUsuarios, createUsuario, updateUsuario, bajaUsuario, expulsarUsuario, unlockUsuario,
  getRoles, createRol, deleteRol,
  asignarRol, revocarRol,
  getPermisos, asignarPermisoARol, revocarPermisoDeRol,
  getSesiones, cerrarSesion,
  getIPsBloqueadas, desbloquearIP,
} from '../../services/identidadApi';
import Layout from '../../components/layout/Layout';
import { usePermisos } from '../../hooks/usePermisos';
import './Admin.css';

function Modal({ title, onClose, onConfirm, confirmLabel = 'Confirmar', danger = false, children }) {
  return (
    <div className="admin-modal-backdrop" onClick={onClose}>
      <div className="admin-modal" onClick={(e) => e.stopPropagation()}>
        <h3 className="admin-modal__title">{title}</h3>
        {children}
        <div className="admin-modal__actions">
          <button className="btn-modal-cancel" onClick={onClose}>Cancelar</button>
          <button className={`btn-modal-confirm ${danger ? 'btn-modal-confirm--danger' : ''}`} onClick={onConfirm}>
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
}

function ModalRoles({ usuario, onClose }) {
  const [roles, setRoles] = useState([]);
  const [rolesUsuario, setRolesUsuario] = useState([]);
  const [loading, setLoading] = useState(false);
  const [mensaje, setMensaje] = useState('');

  useEffect(() => {
    setLoading(true);
    Promise.all([
      getRoles(),
    ]).then(([rolesData]) => {
      setRoles(rolesData.roles || []);
    }).finally(() => setLoading(false));
  }, []);

  const handleAsignar = async (rolID) => {
    try {
      await asignarRol(usuario.id, { rol_id: rolID, tenant_id: '' });
      setMensaje('Rol asignado exitosamente');
      setRolesUsuario(prev => [...prev, rolID]);
      setTimeout(() => setMensaje(''), 3000);
    } catch (e) {
      setMensaje('Error: ' + (e.response?.data?.detail || 'No se pudo asignar'));
      setTimeout(() => setMensaje(''), 3000);
    }
  };

  const handleRevocar = async (rolID) => {
    try {
      await revocarRol(usuario.id, rolID);
      setMensaje('Rol revocado exitosamente');
      setRolesUsuario(prev => prev.filter(id => id !== rolID));
      setTimeout(() => setMensaje(''), 3000);
    } catch (e) {
      setMensaje('Error: ' + (e.response?.data?.detail || 'No se pudo revocar'));
      setTimeout(() => setMensaje(''), 3000);
    }
  };

  return (
    <div className="admin-modal-backdrop" onClick={onClose}>
      <div className="admin-modal" style={{ maxWidth: 520 }} onClick={e => e.stopPropagation()}>
        <h3 className="admin-modal__title">Gestionar roles — {usuario.nombre} {usuario.apellido}</h3>
        <p style={{ fontSize: '0.85rem', color: 'var(--color-gray-500)', marginBottom: '1rem' }}>{usuario.correo}</p>
        {mensaje && (
          <div style={{ padding: '0.5rem 0.75rem', marginBottom: '1rem', borderRadius: '0.5rem',
            background: mensaje.startsWith('Error') ? '#fef2f2' : '#f0fdf4',
            color: mensaje.startsWith('Error') ? '#991b1b' : '#166534',
            fontSize: '0.85rem', fontWeight: 500 }}>
            {mensaje}
          </div>
        )}
        {loading ? <p>Cargando roles...</p> : (
          <table className="admin-table">
            <thead>
              <tr><th>Rol</th><th>Descripción</th><th>Acción</th></tr>
            </thead>
            <tbody>
              {roles.map(r => (
                <tr key={r.id}>
                  <td><strong>{r.nombre}</strong></td>
                  <td style={{ fontSize: '0.8rem', color: 'var(--color-gray-500)' }}>{r.descripcion}</td>
                  <td>
                    <button className="btn-admin btn-admin--primary" style={{ marginRight: 4 }}
                      onClick={() => handleAsignar(r.id)}>
                      Asignar
                    </button>
                    <button className="btn-admin btn-admin--danger"
                      onClick={() => handleRevocar(r.id)}>
                      Revocar
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        <div className="admin-modal__actions" style={{ marginTop: '1.5rem' }}>
          <button className="btn-modal-cancel" onClick={onClose}>Cerrar</button>
        </div>
      </div>
    </div>
  );
}

function TabUsuarios() {
  const { puede } = usePermisos();
  const [usuarios, setUsuarios] = useState([]);
  const [total, setTotal] = useState(0);
  const [pagina, setPagina] = useState(1);
  const [filtroCorreo, setFiltroCorreo] = useState('');
  const [filtroEstado, setFiltroEstado] = useState('');
  const [loading, setLoading] = useState(false);
  const [modal, setModal] = useState(null);
  const [form, setForm] = useState({ correo: '', nombre: '', apellido: '', password: '' });

  const cargar = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getUsuarios({ pagina, correo: filtroCorreo, estado: filtroEstado });
      setUsuarios(data.usuarios || []);
      setTotal(data.total || 0);
    } catch { } finally { setLoading(false); }
  }, [pagina, filtroCorreo, filtroEstado]);

  useEffect(() => { cargar(); }, [cargar]);

  const handleAccion = async () => {
    try {
      if (modal.tipo === 'baja') await bajaUsuario(modal.usuario.id);
      if (modal.tipo === 'expulsar') await expulsarUsuario(modal.usuario.id);
      if (modal.tipo === 'unlock') await unlockUsuario(modal.usuario.id);
      if (modal.tipo === 'crear') await createUsuario(form);
      if (modal.tipo === 'editar') await updateUsuario(modal.usuario.id, { nombre: form.nombre, apellido: form.apellido });
    } catch { } finally {
      setModal(null);
      setForm({ correo: '', nombre: '', apellido: '', password: '' });
      cargar();
    }
  };

  const abrirEditar = (u) => {
    setForm({ correo: u.correo, nombre: u.nombre, apellido: u.apellido, password: '' });
    setModal({ tipo: 'editar', usuario: u });
  };

  const totalPaginas = Math.ceil(total / 20);

  return (
    <>
      <div className="admin-card">
        <div className="admin-card__top">
          <h2 className="admin-card__title">Usuarios ({total})</h2>
          {puede('identidad:usuario:crear') && <button className="btn-add" onClick={() => setModal({ tipo: 'crear' })}>+ Nuevo usuario</button>}
        </div>
        <div className="admin-search">
          <input type="text" placeholder="Filtrar por correo..." value={filtroCorreo} onChange={(e) => { setFiltroCorreo(e.target.value); setPagina(1); }} />
          <select value={filtroEstado} onChange={(e) => { setFiltroEstado(e.target.value); setPagina(1); }}>
            <option value="">Todos los estados</option>
            <option value="ACTIVO">Activo</option>
            <option value="INACTIVO">Inactivo</option>
            <option value="EXPULSADO">Expulsado</option>
            <option value="NO_VERIFICADO">No verificado</option>
          </select>
        </div>
        {loading ? <div className="admin-empty">Cargando...</div>
          : usuarios.length === 0 ? <div className="admin-empty">No hay usuarios.</div>
          : (
            <div className="admin-table-wrapper">
              <table className="admin-table">
                <thead>
                  <tr><th>Nombre</th><th>Correo</th><th>Estado</th><th>Creado</th><th>Acciones</th></tr>
                </thead>
                <tbody>
                  {usuarios.map((u) => (
                    <tr key={u.id}>
                      <td>{u.nombre} {u.apellido}</td>
                      <td>{u.correo}</td>
                      <td><span className={`admin-badge admin-badge--${u.estado?.toLowerCase()}`}>{u.estado}</span></td>
                      <td>{new Date(u.creado_en).toLocaleDateString()}</td>
                      <td>
                        <div className="admin-actions">
                          {puede('identidad:usuario:modificar') && <button className="btn-admin btn-admin--primary" onClick={() => abrirEditar(u)}>Editar</button>}
                          {puede('identidad:rol:asignar') && <button className="btn-admin btn-admin--primary" onClick={() => setModal({ tipo: 'roles', usuario: u })}>Roles</button>}
                          {puede('identidad:credenciales:desbloquear') && <button className="btn-admin btn-admin--warning" onClick={() => setModal({ tipo: 'unlock', usuario: u })}>Desbloquear</button>}
                          {puede('identidad:usuario:eliminar') && <button className="btn-admin btn-admin--danger" onClick={() => setModal({ tipo: 'baja', usuario: u })}>Dar baja</button>}
                          {puede('identidad:usuario:expulsar') && <button className="btn-admin btn-admin--danger" onClick={() => setModal({ tipo: 'expulsar', usuario: u })}>Expulsar</button>}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        {totalPaginas > 1 && (
          <div className="admin-pagination">
            <button className="btn-page" disabled={pagina === 1} onClick={() => setPagina(p => p - 1)}>← Anterior</button>
            <span>Página {pagina} de {totalPaginas}</span>
            <button className="btn-page" disabled={pagina === totalPaginas} onClick={() => setPagina(p => p + 1)}>Siguiente →</button>
          </div>
        )}
      </div>

      {modal?.tipo === 'baja' && (
        <Modal title="Dar de baja usuario" onClose={() => setModal(null)} onConfirm={handleAccion} confirmLabel="Dar baja" danger>
          <p style={{ color: 'var(--color-gray-600)', fontSize: '0.9rem' }}>¿Desactivar a <strong>{modal.usuario.nombre} {modal.usuario.apellido}</strong>?</p>
        </Modal>
      )}
      {modal?.tipo === 'expulsar' && (
        <Modal title="Expulsar usuario" onClose={() => setModal(null)} onConfirm={handleAccion} confirmLabel="Expulsar" danger>
          <p style={{ color: 'var(--color-gray-600)', fontSize: '0.9rem' }}>¿Expulsar y cerrar todas las sesiones de <strong>{modal.usuario.nombre}</strong>?</p>
        </Modal>
      )}
      {modal?.tipo === 'unlock' && (
        <Modal title="Desbloquear cuenta" onClose={() => setModal(null)} onConfirm={handleAccion} confirmLabel="Desbloquear">
          <p style={{ color: 'var(--color-gray-600)', fontSize: '0.9rem' }}>¿Desbloquear la cuenta de <strong>{modal.usuario.nombre}</strong>?</p>
        </Modal>
      )}
      {modal?.tipo === 'roles' && (
        <ModalRoles
          usuario={modal.usuario}
          onClose={() => { setModal(null); cargar(); }}
        />
      )}
      {(modal?.tipo === 'crear' || modal?.tipo === 'editar') && (
        <Modal title={modal.tipo === 'crear' ? 'Nuevo usuario' : 'Editar usuario'} onClose={() => setModal(null)} onConfirm={handleAccion} confirmLabel={modal.tipo === 'crear' ? 'Crear' : 'Guardar'}>
          <div className="form-row">
            <div className="form-group">
              <label className="form-label">Nombre</label>
              <input className="form-input" type="text" value={form.nombre} onChange={(e) => setForm(f => ({ ...f, nombre: e.target.value }))} />
            </div>
            <div className="form-group">
              <label className="form-label">Apellido</label>
              <input className="form-input" type="text" value={form.apellido} onChange={(e) => setForm(f => ({ ...f, apellido: e.target.value }))} />
            </div>
          </div>
          {modal.tipo === 'crear' && (
            <>
              <div className="form-group">
                <label className="form-label">Correo</label>
                <input className="form-input" type="email" value={form.correo} onChange={(e) => setForm(f => ({ ...f, correo: e.target.value }))} />
              </div>
              <div className="form-group">
                <label className="form-label">Contraseña</label>
                <input className="form-input" type="password" placeholder="Mínimo 8 caracteres" value={form.password} onChange={(e) => setForm(f => ({ ...f, password: e.target.value }))} />
              </div>
            </>
          )}
        </Modal>
      )}
    </>
  );
}

function ModalPermisosRol({ rol, onClose }) {
  const [permisos, setPermisos] = useState([]);
  const [loading, setLoading] = useState(false);
  const [guardando, setGuardando] = useState(null);
  const [mensaje, setMensaje] = useState(null);
  const [permisosActivos, setPermisosActivos] = useState(new Set(rol.permisos || []));

  useEffect(() => {
    setLoading(true);
    getPermisos()
      .then(data => setPermisos(data.permisos || []))
      .catch(() => setPermisos([]))
      .finally(() => setLoading(false));
  }, []);

  const mostrarMensaje = (texto, tipo = 'exito') => {
    setMensaje({ texto, tipo });
    setTimeout(() => setMensaje(null), 3000);
  };

  const handleToggle = async (codigo) => {
    if (rol.es_sistema) return;
    setGuardando(codigo);
    try {
      if (permisosActivos.has(codigo)) {
        await revocarPermisoDeRol(rol.id, codigo);
        setPermisosActivos(prev => { const s = new Set(prev); s.delete(codigo); return s; });
        mostrarMensaje('Permiso revocado correctamente');
      } else {
        await asignarPermisoARol(rol.id, codigo);
        setPermisosActivos(prev => new Set([...prev, codigo]));
        mostrarMensaje('Permiso asignado correctamente');
      }
    } catch (e) {
      mostrarMensaje(e.response?.data?.detail || 'Error al actualizar permiso', 'error');
    } finally {
      setGuardando(null);
    }
  };

  const modulos = [...new Set(permisos.map(p => p.modulo))];

  return (
    <div className="admin-modal-backdrop" onClick={onClose} role="dialog" aria-modal="true" aria-label={`Gestionar permisos del rol ${rol.nombre}`}>
      <div className="admin-modal" style={{ maxWidth: 680, width: '95vw', maxHeight: '90vh', display: 'flex', flexDirection: 'column' }} onClick={e => e.stopPropagation()}>
        
        {/* Header */}
        <div style={{ borderBottom: '1px solid var(--color-gray-200)', paddingBottom: '1rem', marginBottom: '1rem', flexShrink: 0 }}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
            <h3 className="admin-modal__title" style={{ margin: 0 }}>
              Permisos del rol
            </h3>
            <button onClick={onClose} aria-label="Cerrar" style={{ background: 'none', border: 'none', cursor: 'pointer', fontSize: '1.25rem', color: 'var(--color-gray-400)', padding: '0.25rem' }}>✕</button>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', flexWrap: 'wrap' }}>
            <span style={{ fontWeight: 700, fontSize: '1.1rem', color: 'var(--color-gray-900)' }}>{rol.nombre}</span>
            {rol.es_sistema && (
              <span style={{ background: '#fef3c7', color: '#92400e', fontSize: '0.75rem', fontWeight: 600, padding: '0.2rem 0.6rem', borderRadius: '9999px' }}>
                Rol de sistema — solo lectura
              </span>
            )}
          </div>
          {!rol.es_sistema && (
            <p style={{ margin: '0.5rem 0 0', fontSize: '0.85rem', color: 'var(--color-gray-500)' }}>
              Activa o desactiva los permisos haciendo clic en el interruptor de cada uno.
            </p>
          )}
        </div>

        {/* Mensaje */}
        {mensaje && (
          <div style={{
            padding: '0.6rem 1rem', marginBottom: '0.75rem', borderRadius: '0.5rem', fontSize: '0.85rem', fontWeight: 500, flexShrink: 0,
            background: mensaje.tipo === 'error' ? '#fef2f2' : '#f0fdf4',
            color: mensaje.tipo === 'error' ? '#991b1b' : '#166534',
            border: `1px solid ${mensaje.tipo === 'error' ? '#fecaca' : '#bbf7d0'}`,
          }}>
            {mensaje.tipo === 'error' ? '⚠️ ' : '✓ '}{mensaje.texto}
          </div>
        )}

        {/* Contenido scrolleable */}
        <div style={{ overflowY: 'auto', flex: 1 }}>
          {loading ? (
            <div style={{ textAlign: 'center', padding: '2rem', color: 'var(--color-gray-500)' }}>Cargando permisos...</div>
          ) : (
            modulos.map(modulo => (
              <div key={modulo} style={{ marginBottom: '1.25rem' }}>
                <div style={{ fontSize: '0.7rem', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em', color: 'var(--color-gray-400)', marginBottom: '0.5rem', paddingBottom: '0.25rem', borderBottom: '1px solid var(--color-gray-100)' }}>
                  {modulo}
                </div>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                  {permisos.filter(p => p.modulo === modulo).map(p => {
                    const activo = permisosActivos.has(p.codigo);
                    const cargando = guardando === p.codigo;
                    return (
                      <div key={p.codigo} style={{
                        display: 'flex', alignItems: 'center', justifyContent: 'space-between',
                        padding: '0.65rem 0.85rem', borderRadius: '0.5rem', gap: '1rem',
                        background: activo ? '#f0fdf4' : '#f8fafc',
                        border: `1px solid ${activo ? '#bbf7d0' : 'var(--color-gray-200)'}`,
                        transition: 'all 0.15s',
                        opacity: cargando ? 0.7 : 1,
                      }}>
                        <div style={{ flex: 1, minWidth: 0 }}>
                          <div style={{ fontSize: '0.8rem', fontWeight: 600, color: 'var(--color-gray-800)', fontFamily: 'monospace' }}>{p.codigo}</div>
                          <div style={{ fontSize: '0.78rem', color: 'var(--color-gray-500)', marginTop: '0.1rem' }}>{p.nombre}</div>
                        </div>
                        {!rol.es_sistema ? (
                          <button
                            onClick={() => handleToggle(p.codigo)}
                            disabled={cargando}
                            aria-label={activo ? `Revocar permiso ${p.nombre}` : `Asignar permiso ${p.nombre}`}
                            aria-pressed={activo}
                            style={{
                              flexShrink: 0, width: 44, height: 24, borderRadius: '9999px', border: 'none', cursor: cargando ? 'wait' : 'pointer',
                              background: activo ? 'var(--color-green-600, #16a34a)' : '#cbd5e1',
                              position: 'relative', transition: 'background 0.2s',
                            }}
                          >
                            <span style={{
                              position: 'absolute', top: 2, left: activo ? 22 : 2, width: 20, height: 20,
                              borderRadius: '50%', background: '#fff', transition: 'left 0.2s',
                              boxShadow: '0 1px 3px rgba(0,0,0,0.2)',
                            }} />
                          </button>
                        ) : (
                          <span style={{ fontSize: '0.75rem', fontWeight: 600, color: activo ? '#16a34a' : '#94a3b8' }}>
                            {activo ? 'Asignado' : '—'}
                          </span>
                        )}
                      </div>
                    );
                  })}
                </div>
              </div>
            ))
          )}
        </div>

        {/* Footer */}
        <div style={{ borderTop: '1px solid var(--color-gray-200)', paddingTop: '1rem', marginTop: '1rem', flexShrink: 0, textAlign: 'right' }}>
          <button className="btn-modal-cancel" onClick={onClose}>Cerrar</button>
        </div>
      </div>
    </div>
  );
}


function TabRoles() {
  const { puede } = usePermisos();
  const [roles, setRoles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modal, setModal] = useState(null);
  const [form, setForm] = useState({ nombre: '', descripcion: '' });

  const cargar = async () => {
    setLoading(true);
    try {
      const data = await getRoles();
      setRoles(data.roles || []);
    } catch { } finally { setLoading(false); }
  };

  useEffect(() => { cargar(); }, []);

  const handleAccion = async () => {
    try {
      if (modal.tipo === 'crear') await createRol(form);
      if (modal.tipo === 'eliminar') await deleteRol(modal.rol.id);
    } catch { } finally {
      setModal(null);
      setForm({ nombre: '', descripcion: '' });
      cargar();
    }
  };

  return (
    <>
      <div className="admin-card">
        <div className="admin-card__top">
          <h2 className="admin-card__title">Roles</h2>
          {puede('identidad:rol:crear') && <button className="btn-add" onClick={() => setModal({ tipo: 'crear' })}>+ Nuevo rol</button>}
        </div>
        {loading ? <div className="admin-empty">Cargando...</div>
          : roles.length === 0 ? <div className="admin-empty">No hay roles.</div>
          : (
            <div className="admin-table-wrapper">
              <table className="admin-table">
                <thead>
                  <tr><th>Nombre</th><th>Descripción</th><th>Sistema</th><th>Acciones</th></tr>
                </thead>
                <tbody>
                  {roles.map((r) => (
                    <tr key={r.id}>
                      <td><strong>{r.nombre}</strong></td>
                      <td>{r.descripcion}</td>
                      <td>{r.es_sistema ? '✅ Sí' : '—'}</td>
                      <td>
                        <div className="admin-actions">
                          <button className="btn-admin btn-admin--primary" onClick={() => setModal({ tipo: 'permisos', rol: r })}>Permisos</button>
                          {!r.es_sistema && <button className="btn-admin btn-admin--danger" onClick={() => setModal({ tipo: 'eliminar', rol: r })}>Eliminar</button>}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
      </div>
      {modal?.tipo === 'permisos' && (
        <ModalPermisosRol
          rol={modal.rol}
          onClose={() => setModal(null)}
        />
      )}
      {modal?.tipo === 'crear' && (
        <Modal title="Nuevo rol" onClose={() => setModal(null)} onConfirm={handleAccion} confirmLabel="Crear">
          <div className="form-group">
            <label className="form-label">Nombre</label>
            <input className="form-input" type="text" value={form.nombre} onChange={(e) => setForm(f => ({ ...f, nombre: e.target.value }))} />
          </div>
          <div className="form-group">
            <label className="form-label">Descripción</label>
            <input className="form-input" type="text" value={form.descripcion} onChange={(e) => setForm(f => ({ ...f, descripcion: e.target.value }))} />
          </div>
        </Modal>
      )}
      {modal?.tipo === 'eliminar' && (
        <Modal title="Eliminar rol" onClose={() => setModal(null)} onConfirm={handleAccion} confirmLabel="Eliminar" danger>
          <p style={{ color: 'var(--color-gray-600)', fontSize: '0.9rem' }}>¿Eliminar el rol <strong>{modal.rol.nombre}</strong>?</p>
        </Modal>
      )}
    </>
  );
}

function TabSesiones() {
  const [sesiones, setSesiones] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modal, setModal] = useState(null);

  const cargar = async () => {
    setLoading(true);
    try {
      const data = await getSesiones();
      setSesiones(data.sesiones || []);
    } catch { } finally { setLoading(false); }
  };

  useEffect(() => { cargar(); }, []);

  const handleCerrar = async () => {
    try { await cerrarSesion(modal.sesion.id); } catch { } finally { setModal(null); cargar(); }
  };

  return (
    <>
      <div className="admin-card">
        <div className="admin-card__top">
          <h2 className="admin-card__title">Sesiones activas</h2>
          <button className="btn-add" onClick={cargar}>↻ Refrescar</button>
        </div>
        {loading ? <div className="admin-empty">Cargando...</div>
          : sesiones.length === 0 ? <div className="admin-empty">No hay sesiones activas.</div>
          : (
            <div className="admin-table-wrapper">
              <table className="admin-table">
                <thead>
                  <tr><th>ID Sesión</th><th>Usuario ID</th><th>IP Origen</th><th>Estado</th><th>Última actividad</th><th>Acciones</th></tr>
                </thead>
                <tbody>
                  {sesiones.map((s) => (
                    <tr key={s.id}>
                      <td style={{ fontFamily: 'monospace', fontSize: '0.75rem' }}>{s.id}</td>
                      <td style={{ fontFamily: 'monospace', fontSize: '0.75rem' }}>{s.usuario_id}</td>
                      <td>{s.ip_origen}</td>
                      <td><span className="admin-badge admin-badge--activo">{s.estado}</span></td>
                      <td>{new Date(s.ultima_actividad).toLocaleString()}</td>
                      <td><button className="btn-admin btn-admin--danger" onClick={() => setModal({ sesion: s })}>Cerrar</button></td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
      </div>
      {modal && (
        <Modal title="Cerrar sesión" onClose={() => setModal(null)} onConfirm={handleCerrar} confirmLabel="Cerrar sesión" danger>
          <p style={{ color: 'var(--color-gray-600)', fontSize: '0.9rem' }}>¿Forzar el cierre de la sesión <strong>{modal.sesion.id}</strong>?</p>
        </Modal>
      )}
    </>
  );
}

function TabIPsBloqueadas() {
  const [ips, setIps] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modal, setModal] = useState(null);

  const cargar = async () => {
    setLoading(true);
    try {
      const data = await getIPsBloqueadas();
      setIps(data.ips || []);
    } catch { } finally { setLoading(false); }
  };

  useEffect(() => { cargar(); }, []);

  const handleDesbloquear = async () => {
    try { await desbloquearIP(modal.ip.ip); } catch { } finally { setModal(null); cargar(); }
  };

  return (
    <>
      <div className="admin-card">
        <div className="admin-card__top">
          <h2 className="admin-card__title">IPs bloqueadas</h2>
          <button className="btn-add" onClick={cargar}>↻ Refrescar</button>
        </div>
        {loading ? <div className="admin-empty">Cargando...</div>
          : ips.length === 0 ? <div className="admin-empty">No hay IPs bloqueadas.</div>
          : (
            <div className="admin-table-wrapper">
              <table className="admin-table">
                <thead>
                  <tr><th>IP</th><th>Intentos</th><th>Bloqueada hasta</th><th>Acciones</th></tr>
                </thead>
                <tbody>
                  {ips.map((ip) => (
                    <tr key={ip.ip}>
                      <td style={{ fontFamily: 'monospace' }}>{ip.ip}</td>
                      <td>{ip.intentos}</td>
                      <td>{new Date(ip.bloqueado_hasta).toLocaleString()}</td>
                      <td><button className="btn-admin btn-admin--primary" onClick={() => setModal({ ip })}>Desbloquear</button></td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
      </div>
      {modal && (
        <Modal title="Desbloquear IP" onClose={() => setModal(null)} onConfirm={handleDesbloquear} confirmLabel="Desbloquear">
          <p style={{ color: 'var(--color-gray-600)', fontSize: '0.9rem' }}>¿Desbloquear la IP <strong>{modal.ip.ip}</strong>?</p>
        </Modal>
      )}
    </>
  );
}

export default function AdminPage() {
  const { user } = useAuth();
  const [tabActiva, setTabActiva] = useState('Usuarios');

  const isSysAdmin = user?.rol === 'sys_admin';
  const TABS = isSysAdmin
    ? ['Usuarios', 'Roles', 'Sesiones', 'IPs Bloqueadas']
    : ['Usuarios', 'Roles'];
  const subtitle = isSysAdmin
    ? 'Gestión de usuarios, roles, sesiones e IPs bloqueadas.'
    : 'Gestión de usuarios y roles de tu finca.';

  return (
    <Layout title="Panel de administración" subtitle={subtitle}>
      <div className="admin-tabs">
        {TABS.map((tab) => (
          <button key={tab} className={`admin-tab ${tabActiva === tab ? 'admin-tab--active' : ''}`} onClick={() => setTabActiva(tab)}>
            {tab}
          </button>
        ))}
      </div>
      {tabActiva === 'Usuarios' && <TabUsuarios />}
      {tabActiva === 'Roles' && <TabRoles />}
      {tabActiva === 'Sesiones' && <TabSesiones />}
      {tabActiva === 'IPs Bloqueadas' && <TabIPsBloqueadas />}
    </Layout>
  );
}