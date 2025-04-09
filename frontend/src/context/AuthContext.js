import React, { createContext, useState, useEffect, useContext, useCallback } from 'react';

// Crear el contexto de autenticación
const AuthContext = createContext(null);

// Proveedor del contexto
export const AuthProvider = ({ children }) => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState(null);
  const [authError, setAuthError] = useState(null);

  // URL base de la API
  const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';

  // Función para decodificar JWT
  const parseJwt = (token) => {
    try {
      const base64Url = token.split('.')[1];
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join('')
      );
      return JSON.parse(jsonPayload);
    } catch (e) {
      console.error('Error al parsear JWT:', e);
      return {};
    }
  };

  // Función para verificar el estado de autenticación
  const checkAuthStatus = useCallback(async () => {
    const token = localStorage.getItem('token');
    
    if (!token) {
      setIsLoggedIn(false);
      setUser(null);
      setLoading(false);
      return;
    }
    
    try {
      // Intentar obtener información del usuario desde el token
      const userData = parseJwt(token);
      
      // Verificar si el token ha expirado
      const expirationTime = userData.exp * 1000; // Convertir a milisegundos
      if (expirationTime < Date.now()) {
        // Token expirado
        localStorage.removeItem('token');
        setIsLoggedIn(false);
        setUser(null);
        setAuthError('La sesión ha expirado. Por favor, inicie sesión nuevamente.');
        setLoading(false);
        return;
      }
      
      setIsLoggedIn(true);
      
      // Establecer datos básicos del usuario desde el token
      setUser({
        id: userData.user_id,
        role: userData.role || 'user',
        nombre: userData.nombre || 'Usuario',
        email: userData.email || ''
      });
      
      // Verificar si el token está por expirar (si queda menos de 1 hora)
      const timeRemaining = expirationTime - Date.now();
      
      // Si queda poco tiempo, intentamos obtener datos actualizados del usuario
      if (timeRemaining < 3600000) { // Menos de 1 hora
        try {
          const response = await fetch(`${apiUrl}/api/auth/check-admin`, {
            headers: {
              'Authorization': `Bearer ${token}`,
              'Accept': 'application/json'
            }
          });
          
          if (response.ok) {
            const data = await response.json();
            if (data.success && data.isAdmin) {
              setUser(prev => ({ ...prev, role: 'admin' }));
            }
          }
        } catch (err) {
          console.error('Error al verificar estado de administrador:', err);
        }
      }
      
    } catch (error) {
      console.error('Error al procesar token:', error);
      setAuthError('Error al procesar la información de autenticación.');
    } finally {
      setLoading(false);
    }
  }, [apiUrl]);

  // Función para iniciar sesión
  const login = async (token, userData) => {
    localStorage.setItem('token', token);
    setIsLoggedIn(true);
    setAuthError(null);
    
    // Si recibimos datos del usuario, los guardamos
    if (userData) {
      setUser(userData);
    } else {
      // Si no, intentamos extraerlos del token
      try {
        const tokenData = parseJwt(token);
        setUser({
          id: tokenData.user_id,
          role: tokenData.role || 'user',
          nombre: tokenData.nombre || 'Usuario',
          email: tokenData.email || ''
        });
      } catch (error) {
        console.error('Error al procesar token en login:', error);
      }
    }
    
    window.dispatchEvent(new Event('loginSuccess'));
  };

  // Función para cerrar sesión
  const logout = () => {
    localStorage.removeItem('token');
    setIsLoggedIn(false);
    setUser(null);
    setAuthError(null);
    window.dispatchEvent(new Event('logoutSuccess'));
  };

  // Función para comprobar si un usuario es administrador (localmente)
  const isAdmin = () => {
    return user && user.role === 'admin';
  };

  // Verificar si el usuario es administrador mediante API
  const checkAdminStatus = async () => {
    if (!isLoggedIn) return false;
    
    const token = localStorage.getItem('token');
    if (!token) return false;
    
    try {
      // Usar AbortController para manejar timeouts
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 5000); // 5 segundos de timeout
      
      const response = await fetch(`${apiUrl}/api/auth/check-admin`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Accept': 'application/json'
        },
        signal: controller.signal
      });
      
      clearTimeout(timeoutId);
      
      if (!response.ok) {
        console.error('Error verificando estado de administrador:', response.status);
        if (response.status === 401 || response.status === 403) {
          // Si el token no es válido o el usuario no tiene permisos, cerramos la sesión
          logout();
        }
        return false;
      }
      
      const data = await response.json();
      
      // Actualizar el rol del usuario si es necesario
      if (data.success) {
        if (data.isAdmin && user && user.role !== 'admin') {
          setUser({...user, role: 'admin'});
        } else if (!data.isAdmin && user && user.role === 'admin') {
          setUser({...user, role: 'user'});
        }
        
        return data.isAdmin;
      } else {
        // Si la respuesta no tiene éxito, asumimos que no es admin
        if (user && user.role === 'admin') {
          setUser({...user, role: 'user'});
        }
        return false;
      }
    } catch (err) {
      console.error('Error al verificar estado de administrador:', err);
      // Si es un error de timeout o de red, no cambiamos el rol local
      if (err.name !== 'AbortError') {
        // Para otros errores, asumimos el peor caso (no es admin)
        if (user && user.role === 'admin') {
          setUser({...user, role: 'user'});
        }
      }
      return false;
    }
  };

  // Nueva función para actualizar los datos del usuario
  const updateUser = async (newUserData) => {
    if (!user) return;
    
    // Actualizar en el estado local
    setUser(prevUser => ({
      ...prevUser,
      ...newUserData
    }));
  };

  // Nueva función para refrescar el token
  const refreshToken = async () => {
    const token = localStorage.getItem('token');
    if (!token) return false;
    
    try {
      const response = await fetch(`${apiUrl}/api/auth/refresh-token`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });
      
      if (!response.ok) {
        if (response.status === 401) {
          logout(); // Si el token ya no es válido, cerramos sesión
        }
        return false;
      }
      
      const data = await response.json();
      if (data.success && data.token) {
        localStorage.setItem('token', data.token);
        return true;
      }
      
      return false;
    } catch (error) {
      console.error('Error al refrescar token:', error);
      return false;
    }
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
  }, [checkAuthStatus]);

  return (
    <AuthContext.Provider value={{ 
      isLoggedIn, 
      loading, 
      user, 
      authError,
      login, 
      logout, 
      isAdmin, 
      checkAdminStatus,
      updateUser,  // Añadida nueva función
      refreshToken // Añadida nueva función 
    }}>
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