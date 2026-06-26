import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { DiagnosisProvider } from './context/DiagnosisContext';
import { PrivateRoute, PublicRoute, AdminRoute, TenantConfigRoute } from './components/auth/RoleGuards';
import LandingPage from './pages/landing/LandingPage';
import LoginPage from './pages/auth/LoginPage';
import RegisterPage from './pages/auth/RegisterPage';
import ForgotPasswordPage from './pages/auth/ForgotPasswordPage';
import ResetPasswordPage from './pages/auth/ResetPasswordPage';
import VerificacionCorreoPage from './pages/auth/VerificacionCorreoPage';
import AceptarInvitacionPage from './pages/auth/AceptarInvitacionPage';
import DashboardPage from './pages/dashboard/DashboardPage';
import ToastContainer from './components/toast/ToastContainer';
import PerfilPage from './pages/perfil/PerfilPage';
import FincasPage from './pages/fincas/FincasPage';
import FincaConfigPage from './pages/finca-config/FincaConfigPage';
import AdminPage from './pages/admin/AdminPage';

function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<LandingPage />} />
      <Route path="/login" element={<PublicRoute><LoginPage /></PublicRoute>} />
      <Route path="/register" element={<PublicRoute><RegisterPage /></PublicRoute>} />
      <Route path="/forgot-password" element={<PublicRoute><ForgotPasswordPage /></PublicRoute>} />
      <Route path="/reset-password" element={<ResetPasswordPage />} />
      <Route path="/verificar-correo" element={<VerificacionCorreoPage />} />
      <Route path="/aceptar-invitacion" element={<AceptarInvitacionPage />} />
      <Route path="/dashboard" element={<PrivateRoute><DashboardPage /></PrivateRoute>} />
      <Route path="/perfil" element={<PrivateRoute><PerfilPage /></PrivateRoute>} />
      <Route path="/fincas" element={<PrivateRoute><FincasPage /></PrivateRoute>} />
      <Route path="/finca-config" element={<TenantConfigRoute><FincaConfigPage /></TenantConfigRoute>} />
      <Route path="/admin" element={<AdminRoute><AdminPage /></AdminRoute>} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <DiagnosisProvider>
          <AppRoutes />
          <ToastContainer />
        </DiagnosisProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}
