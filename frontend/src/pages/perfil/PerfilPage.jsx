import { useState, useEffect } from 'react';
import { useAuth } from '../../context/AuthContext';
import { getMiPerfil, updateMiPerfil, updateMiPassword } from '../../services/identidadApi';
import { validarPassword } from '../../services/validacionPassword';
import Layout from '../../components/layout/Layout';
import './Perfil.css';

export default function PerfilPage() {
  const { user } = useAuth();

  const [perfil, setPerfil] = useState(null);
  const [nombre, setNombre] = useState('');
  const [apellido, setApellido] = useState('');
  const [passwordActual, setPasswordActual] = useState('');
  const [nuevaPassword, setNuevaPassword] = useState('');
  const [loadingPerfil, setLoadingPerfil] = useState(false);
  const [loadingPassword, setLoadingPassword] = useState(false);
  const [msgPerfil, setMsgPerfil] = useState(null);
  const [msgPassword, setMsgPassword] = useState(null);

  useEffect(() => {
    getMiPerfil().then((d) => {
      setPerfil(d);
      setNombre(d.nombre);
      setApellido(d.apellido);
    }).catch(() => {});
  }, []);

  const handleUpdatePerfil = async (e) => {
    e.preventDefault();
    setLoadingPerfil(true);
    setMsgPerfil(null);
    try {
      await updateMiPerfil({ nombre, apellido });
      setMsgPerfil({ tipo: 'exito', texto: 'Perfil actualizado correctamente.' });
    } catch {
      setMsgPerfil({ tipo: 'error', texto: 'Error al actualizar el perfil.' });
    } finally {
      setLoadingPerfil(false);
    }
  };

  const handleUpdatePassword = async (e) => {
    e.preventDefault();
    const validacion = validarPassword(nuevaPassword);
    if (!validacion.valida) {
      setMsgPassword({ tipo: 'error', texto: validacion.errores.join('. ') });
      return;
    }
    setLoadingPassword(true);
    setMsgPassword(null);
    try {
      await updateMiPassword({ password_actual: passwordActual, nueva_password: nuevaPassword });
      setMsgPassword({ tipo: 'exito', texto: 'Contraseña actualizada correctamente.' });
      setPasswordActual('');
      setNuevaPassword('');
    } catch {
      setMsgPassword({ tipo: 'error', texto: 'Contraseña actual incorrecta.' });
    } finally {
      setLoadingPassword(false);
    }
  };

  const initials = perfil
    ? `${perfil.nombre?.[0] || ''}${perfil.apellido?.[0] || ''}`.toUpperCase()
    : '??';

  return (
    <Layout title="Mi Perfil" subtitle="Gestiona tu información personal y contraseña.">
      {/* Info personal */}
      <div className="perfil-card">
        <h2 className="perfil-card__title">Información personal</h2>
        <div className="perfil-avatar">
          <div className="perfil-avatar__circle">{initials}</div>
          <div className="perfil-avatar__info">
            <p className="perfil-avatar__name">{perfil ? `${perfil.nombre} ${perfil.apellido}` : '...'}</p>
            <p>{perfil?.correo}</p>
            <p>Estado: {perfil?.estado}</p>
          </div>
        </div>

        {msgPerfil && (
          <div className={msgPerfil.tipo === 'exito' ? 'perfil-success' : 'perfil-error'}>
            {msgPerfil.texto}
          </div>
        )}

        <form onSubmit={handleUpdatePerfil}>
          <div className="form-row">
            <div className="form-group">
              <label className="form-label" htmlFor="p-nombre">Nombre</label>
              <input id="p-nombre" className="form-input" type="text" value={nombre} onChange={(e) => setNombre(e.target.value)} required disabled={loadingPerfil} />
            </div>
            <div className="form-group">
              <label className="form-label" htmlFor="p-apellido">Apellido</label>
              <input id="p-apellido" className="form-input" type="text" value={apellido} onChange={(e) => setApellido(e.target.value)} required disabled={loadingPerfil} />
            </div>
          </div>
          <button type="submit" className="btn-perfil btn-perfil--primary" disabled={loadingPerfil}>
            {loadingPerfil ? 'Guardando...' : 'Guardar cambios'}
          </button>
        </form>
      </div>

      {/* Cambiar contraseña */}
      <div className="perfil-card">
        <h2 className="perfil-card__title">Cambiar contraseña</h2>

        {msgPassword && (
          <div className={msgPassword.tipo === 'exito' ? 'perfil-success' : 'perfil-error'}>
            {msgPassword.texto}
          </div>
        )}

        <form onSubmit={handleUpdatePassword}>
          <div className="form-group">
            <label className="form-label" htmlFor="p-actual">Contraseña actual</label>
            <input id="p-actual" className="form-input" type="password" placeholder="••••••••" value={passwordActual} onChange={(e) => setPasswordActual(e.target.value)} required disabled={loadingPassword} />
          </div>
          <div className="form-group">
            <label className="form-label" htmlFor="p-nueva">Nueva contraseña</label>
            <input id="p-nueva" className="form-input" type="password" placeholder="Mínimo 8 caracteres" value={nuevaPassword} onChange={(e) => setNuevaPassword(e.target.value)} required disabled={loadingPassword} />
          </div>
          <button type="submit" className="btn-perfil btn-perfil--primary" disabled={loadingPassword}>
            {loadingPassword ? 'Actualizando...' : 'Actualizar contraseña'}
          </button>
        </form>
      </div>
    </Layout>
  );
}