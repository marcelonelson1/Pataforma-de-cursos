import React, { createContext, useState, useEffect, useContext } from 'react';

// Crear el contexto de autenticación
const AuthContext = createContext(null);

// Proveedor del contexto
export const AuthProvider = ({ children }) => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [loading, setLoading] = useState(true);

  // Función para verificar el estado de autenticación
  const checkAuthStatus = () => {
    const token = localStorage.getItem('token');
    setIsLoggedIn(!!token);
    setLoading(false);
  };

  // Función para iniciar sesión
  const login = (token) => {
    localStorage.setItem('token', token);
    setIsLoggedIn(true);
    window.dispatchEvent(new Event('loginSuccess'));
  };

  // Función para cerrar sesión
  const logout = () => {
    localStorage.removeItem('token');
    setIsLoggedIn(false);
    window.dispatchEvent(new Event('logoutSuccess'));
  };

  useEffect(() => {
    // Verificar estado de autenticación al montar el componente
    checkAuthStatus();

    // Escuchar eventos de almacenamiento (para sincronizar entre pestañas)
    const handleStorageChange = (e) => {
      if (e.key === 'token') {
        checkAuthStatus();
      }
    };

    window.addEventListener('storage', handleStorageChange);
    window.addEventListener('loginSuccess', checkAuthStatus);
    window.addEventListener('logoutSuccess', checkAuthStatus);
    
    return () => {
      window.removeEventListener('storage', handleStorageChange);
      window.removeEventListener('loginSuccess', checkAuthStatus);
      window.removeEventListener('logoutSuccess', checkAuthStatus);
    };
  }, []);

  return (
    <AuthContext.Provider value={{ isLoggedIn, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

// Hook personalizado para usar el contexto de autenticación
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth debe ser usado dentro de un AuthProvider');
  }
  return context;
};