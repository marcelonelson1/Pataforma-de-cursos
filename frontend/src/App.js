import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Header from './components/Header';
import Footer from './components/Footer';
import Home from './pages/Home';
import Portfolio from './pages/portfolio';
import CursosPage from './pages/CursosPage';
import SobreMiPage from './pages/SobreMiPage';
import ContactoPage from './pages/ContactoPage';
import ServiciosPage from './pages/ServiciosPage'; // Importación del nuevo componente
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ForgotPasswordPage from './pages/ForgotPasswordPage';
import ResetPasswordPage from './pages/ResetPasswordPage';
import CursoDetalle from './components/CursoDetalle';
import { AuthProvider } from './context/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';
import AdminRoutes from './pages/admin/AdminRoutes';
import 'animate.css';
import './App.css';
import './pages/ContactoPage.css';
import './pages/ServiciosPage.css'; // Importación del CSS de la página de servicios

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="app">
          <Header />
          <main>
            <Routes>
              <Route path="/*" element={<MainRoutes />} />
              <Route path="/admin/*" element={<AdminRoutes />} />
            </Routes>
          </main>
          <Footer />
        </div>
      </Router>
    </AuthProvider>
  );
}

const MainRoutes = () => {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/portfolio" element={<Portfolio />} />
      <Route path="/cursos" element={<CursosPage />} />
      <Route path="/servicios" element={<ServiciosPage />} /> {/* Nueva ruta para la página de servicios */}
      <Route path="/sobre-mi" element={<SobreMiPage />} />
      <Route path="/contacto" element={<ContactoPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/registro" element={<RegisterPage />} />
      <Route path="/recuperar-contrasena" element={<ForgotPasswordPage />} />
      <Route path="/reset-password/:token" element={<ResetPasswordPage />} />
      <Route 
        path="/curso/:id" 
        element={
          <ProtectedRoute>
            <CursoDetalle />
          </ProtectedRoute>
        } 
      />
    </Routes>
  );
};

export default App;