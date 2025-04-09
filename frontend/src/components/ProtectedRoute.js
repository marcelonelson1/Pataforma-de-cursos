import React, { useState, useEffect } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { FaExclamationTriangle } from 'react-icons/fa';
import './ProtectedRoutes.css';


// Componente de carga
const LoadingSpinner = () => (
  <div className="loading-overlay">
    <div className="spinner"></div>
    <p>Verificando permisos...</p>
  </div>
);

// Componente de error de acceso
const AccessDenied = ({ message, onGoBack }) => (
  <div className="error-container">
    <FaExclamationTriangle size={40} color="#ff4d4f" />
    <h2>Acceso Denegado</h2>
    <p>{message || "No tienes permisos para acceder a esta página"}</p>
    <button onClick={onGoBack} className="btn-primary">
      Ir a Inicio
    </button>
  </div>
);

const ProtectedRoute = ({ children, adminOnly = false }) => {
  const { isLoggedIn, loading, user, checkAdminStatus } = useAuth();
  const [isVerifying, setIsVerifying] = useState(true);
  const [hasAdminAccess, setHasAdminAccess] = useState(false);
  const [error, setError] = useState("");
  const location = useLocation();

  useEffect(() => {
    const verifyAccess = async () => {
      setIsVerifying(true);
      
      try {
        // Si no hay autenticación, no necesitamos verificar nada más
        if (!isLoggedIn) {
          setIsVerifying(false);
          return;
        }
        
        // Si la ruta requiere permisos de administrador
        if (adminOnly) {
          // Verificación local rápida
          const localCheck = user && user.role === 'admin';
          
          if (!localCheck) {
            setHasAdminAccess(false);
            setError("No tienes permisos de administrador");
            setIsVerifying(false);
            return;
          }
          
          // Verificación en el servidor (más segura)
          const serverCheck = await checkAdminStatus();
          
          if (!serverCheck) {
            setHasAdminAccess(false);
            setError("No tienes permisos de administrador según el servidor");
            setIsVerifying(false);
            return;
          }
          
          // Si pasa ambas verificaciones
          setHasAdminAccess(true);
        }
      } catch (err) {
        console.error("Error al verificar permisos:", err);
        setError("Error al verificar permisos. Por favor, intenta nuevamente.");
        setHasAdminAccess(false);
      } finally {
        setIsVerifying(false);
      }
    };

    if (!loading) {
      verifyAccess();
    }
  }, [isLoggedIn, adminOnly, user, loading, checkAdminStatus]);

  // Si está cargando datos de autenticación o verificando permisos
  if (loading || isVerifying) {
    return <LoadingSpinner />;
  }

  // Si el usuario no está autenticado, redirigir al login
  if (!isLoggedIn) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // Si la ruta requiere admin pero el usuario no tiene permisos
  if (adminOnly && !hasAdminAccess) {
    return <AccessDenied message={error} onGoBack={() => window.location.href = '/'} />;
  }

  // Si todo está en orden, mostrar el contenido protegido
  return children;
};

export default ProtectedRoute;