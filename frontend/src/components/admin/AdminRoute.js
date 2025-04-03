import React, { useState, useEffect } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

function AdminRoute({ children }) {
  const { isLoggedIn } = useAuth();
  const [isAdmin, setIsAdmin] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Si no está autenticado, no hay necesidad de verificar rol de admin
    if (!isLoggedIn) {
      setIsAdmin(false);
      setLoading(false);
      return;
    }

    // Verificar si el usuario es administrador
    const verifyAdminRole = async () => {
      try {
        const token = localStorage.getItem('token');
        
        const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
        const response = await fetch(`${apiUrl}/api/auth/verify-admin`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });

        if (response.ok) {
          setIsAdmin(true);
        } else {
          setIsAdmin(false);
        }
      } catch (error) {
        console.error('Error al verificar rol de administrador:', error);
        setIsAdmin(false);
      } finally {
        setLoading(false);
      }
    };

    verifyAdminRole();
  }, [isLoggedIn]);

  // Mostrar indicador de carga mientras se verifica
  if (loading) {
    return (
      <div className="admin-loading">
        <div className="spinner"></div>
        <p>Verificando acceso...</p>
      </div>
    );
  }

  // Redirigir si no es admin
  if (!isAdmin) {
    return <Navigate to="/login" replace />;
  }

  // Permitir acceso si es admin
  return children;
}

export default AdminRoute;