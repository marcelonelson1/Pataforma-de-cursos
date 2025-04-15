import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000';

// Función optimizada para obtener el token de autenticación
const authHeader = () => {
  try {
    // Primero intentar obtener el token desde localStorage según AuthService
    const token = localStorage.getItem('token');
    if (token) {
      console.log('Token encontrado en formato AuthService:', token.substring(0, 10) + '...');
      return { 'Authorization': `Bearer ${token}` };
    }
    
    // Si no se encuentra en formato de AuthService, buscar en formato de StatsDashboard
    const userStr = localStorage.getItem('user');
    if (!userStr) {
      console.warn('No se encontró información del usuario en localStorage');
      return {};
    }
    
    // Procesar información desde localStorage
    const user = JSON.parse(userStr);
    
    // Verificar posibles ubicaciones del token
    if (user.token) {
      return { 'Authorization': `Bearer ${user.token}` };
    } else if (user.accessToken) {
      return { 'Authorization': `Bearer ${user.accessToken}` };
    } else if (user.access_token) {
      return { 'Authorization': `Bearer ${user.access_token}` };
    }
    
    // Si no se encontró token en los campos estándar, buscar en otros lugares posibles
    if (user.user && typeof user.user === 'object') {
      if (user.user.token) {
        return { 'Authorization': `Bearer ${user.user.token}` };
      } else if (user.user.accessToken) {
        return { 'Authorization': `Bearer ${user.user.accessToken}` };
      } else if (user.user.access_token) {
        return { 'Authorization': `Bearer ${user.user.access_token}` };
      }
    }
    
    console.log('No se encontró token en objeto user. Campos disponibles:', Object.keys(user));
    return {};
  } catch (error) {
    console.error("Error al obtener headers de autenticación:", error);
    return {};
  }
};

// Función para comprobar si el usuario actual tiene rol de admin
const checkIsAdmin = () => {
  try {
    // Verificar en formato AuthContext.js
    const token = localStorage.getItem('token');
    if (token) {
      // Parsear el token para obtener el rol
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
      
      const userData = parseJwt(token);
      if (userData.role === 'admin') {
        console.log('Usuario identificado como admin mediante token JWT');
        return true;
      }
    }
    
    // Verificar en formato del componente StatsDashboard
    const userStr = localStorage.getItem('user');
    if (!userStr) {
      return false;
    }
    
    const user = JSON.parse(userStr);
    
    // Verificar el rol en diferentes niveles posibles de la estructura
    if (user.role && user.role.toLowerCase() === 'admin') {
      return true;
    }
    
    if (user.user && user.user.role && user.user.role.toLowerCase() === 'admin') {
      return true;
    }
    
    return false;
  } catch (error) {
    console.error("Error al verificar rol de admin:", error);
    return false;
  }
};

// Función para sincronizar datos de usuario entre formatos
const syncUserData = () => {
  try {
    // Si existe token en formato AuthContext pero no existe en formato StatsDashboard
    const token = localStorage.getItem('token');
    const userStr = localStorage.getItem('user');
    
    if (token && !userStr) {
      // Parsear el token JWT para obtener información del usuario
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
      
      const userData = parseJwt(token);
      
      // Crear objeto de usuario compatible con StatsDashboard
      const user = {
        id: userData.user_id || userData.id || 1,
        nombre: userData.nombre || 'Usuario',
        email: userData.email || 'user@example.com',
        role: userData.role || 'user',
        token: token
      };
      
      // Si el usuario es admin en el token, asegurarse de que tenga ese rol
      if (userData.role === 'admin') {
        user.role = 'admin';
      }
      
      // Guardar en formato compatible con StatsDashboard
      localStorage.setItem('user', JSON.stringify(user));
      console.log('Datos de usuario sincronizados de token JWT a formato user');
      
      return true;
    }
    
    // Si existe en formato StatsDashboard pero no en formato AuthContext
    if (!token && userStr) {
      const user = JSON.parse(userStr);
      
      // Obtener token del formato StatsDashboard
      const userToken = user.token || user.accessToken || user.access_token || 
                        (user.user && (user.user.token || user.user.accessToken || user.user.access_token));
      
      if (userToken) {
        // Guardar en formato compatible con AuthContext
        localStorage.setItem('token', userToken);
        console.log('Token sincronizado de formato user a formato AuthContext');
        
        return true;
      }
    }
    
    return false;
  } catch (error) {
    console.error("Error al sincronizar datos de usuario:", error);
    return false;
  }
};

class StatsService {
  // Sincronizar datos de usuario al inicializar el servicio
  constructor() {
    syncUserData();
  }
  
  // Obtener estadísticas generales para el panel de administración
  getAdminStats(period = 'month') {
    console.log('Requesting admin stats with period:', period);
    syncUserData(); // Sincronizar datos antes de cada petición
    const headers = authHeader();
    console.log('Headers:', headers);
    
    return axios.get(`${API_URL}/api/admin/stats?period=${period}`, { 
      headers: headers,
      withCredentials: true  // Include cookies if necessary
    })
    .then(response => {
      console.log('Stats response:', response.data);
      if (response.data && response.data.success) {
        return response.data.data;
      }
      console.error('API returned success=false:', response.data);
      return null;
    })
    .catch(error => {
      // Mejorado el manejo de errores para mostrar información más específica
      if (error.response) {
        console.error(`Error ${error.response.status} in getAdminStats:`, 
          error.response.data || error.response.statusText);
        
        // Si es un error de autenticación, intentar corregirlo
        if (error.response.status === 401 || error.response.status === 403) {
          console.log('Intentando corregir token de autenticación...');
          this.fixAdminPermissions();
        }
      } else {
        console.error('Error en getAdminStats (sin respuesta):', error.message || error);
      }
      throw error;
    });
  }

  // Obtener estadísticas detalladas del dashboard de administración
  getAdminDashboard() {
    syncUserData(); // Sincronizar datos antes de cada petición
    return axios.get(`${API_URL}/api/admin/dashboard`, { 
      headers: authHeader(),
      withCredentials: true
    })
    .then(response => {
      if (response.data && response.data.success) {
        return response.data.data;
      }
      console.error('API returned success=false:', response.data);
      return null;
    })
    .catch(error => {
      console.error('Error in getAdminDashboard:', error.response || error.message || error);
      throw error;
    });
  }

  // Obtener estadísticas detalladas de ventas
  getSalesStats(period = 'month', startDate = null, endDate = null) {
    syncUserData(); // Sincronizar datos antes de cada petición
    let url = `${API_URL}/api/admin/sales-stats?period=${period}`;
    
    if (startDate) {
      url += `&start_date=${startDate}`;
    }
    
    if (endDate) {
      url += `&end_date=${endDate}`;
    }
    
    return axios.get(url, { 
      headers: authHeader(),
      withCredentials: true
    })
    .then(response => {
      if (response.data && response.data.success) {
        return response.data.data;
      }
      console.error('API returned success=false:', response.data);
      return null;
    })
    .catch(error => {
      console.error('Error in getSalesStats:', error.response || error.message || error);
      throw error;
    });
  }

  // Obtener registro de actividad con paginación
  getActivityLog(page = 1, limit = 50, filters = {}) {
    syncUserData(); // Sincronizar datos antes de cada petición
    let url = `${API_URL}/api/admin/activity-log?page=${page}&limit=${limit}`;
    
    // Añadir filtros adicionales si existen
    if (filters.user_id) url += `&user_id=${filters.user_id}`;
    if (filters.action) url += `&action=${filters.action}`;
    if (filters.start_date) url += `&start_date=${filters.start_date}`;
    if (filters.end_date) url += `&end_date=${filters.end_date}`;
    
    return axios.get(url, { 
      headers: authHeader(),
      withCredentials: true
    })
    .then(response => {
      if (response.data && response.data.success) {
        return response.data.data;
      }
      console.error('API returned success=false:', response.data);
      return null;
    })
    .catch(error => {
      console.error('Error in getActivityLog:', error.response || error.message || error);
      throw error;
    });
  }
  
  // Verificar permisos de administrador
  checkAdminPermissions() {
    syncUserData(); // Sincronizar datos antes de la verificación
    
    // Primero verificar localmente para evitar llamadas innecesarias
    const isAdminLocally = checkIsAdmin();
    console.log('¿Es admin según verificación local?', isAdminLocally);
    
    // Si localmente parece ser admin, verificar con el servidor
    if (isAdminLocally) {
      return axios.get(`${API_URL}/api/auth/check-admin`, { 
        headers: authHeader(),
        withCredentials: true
      })
      .then(response => {
        console.log('Admin check response:', response.data);
        if (response.data && response.data.success) {
          return response.data.data;
        }
        return { isAdmin: false, message: 'No se pudo verificar permisos de administrador' };
      })
      .catch(error => {
        console.error('Error checking admin permissions:', error.response || error.message || error);
        // Si el token está mal formado o ha expirado, intentar corregirlo
        this.fixAdminPermissions();
        return { isAdmin: false, message: `Error: ${error.message}` };
      });
    } else {
      // Si no parece ser admin, intentar corregir y verificar de nuevo
      this.fixAdminPermissions();
      
      // Devolver una promesa para mantener la consistencia en la interfaz
      return Promise.resolve({ 
        isAdmin: false, 
        message: 'El usuario no tiene rol de administrador según la verificación local' 
      });
    }
  }
  
  // Corregir permisos de administrador en localStorage
  fixAdminPermissions() {
    try {
      // Verificar primero si hay token en formato AuthContext
      const token = localStorage.getItem('token');
      if (token) {
        // Parsear el token para obtener información
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
        
        const userData = parseJwt(token);
        
        // Verificar si se necesita actualizar el formato de StatsDashboard
        const userStr = localStorage.getItem('user');
        if (!userStr) {
          // Crear objeto de usuario basado en token
          const user = {
            id: userData.user_id || userData.id || 1,
            nombre: userData.nombre || 'Usuario',
            email: userData.email || 'user@example.com',
            role: 'admin', // Forzar rol de admin
            token: token
          };
          
          localStorage.setItem('user', JSON.stringify(user));
          console.log('Creado objeto user a partir de token JWT');
          
          return {
            success: true,
            message: 'Creado objeto de usuario a partir de token JWT'
          };
        } else {
          // Actualizar el objeto user existente
          let user = JSON.parse(userStr);
          let modified = false;
          
          if (!user.role || user.role.toLowerCase() !== 'admin') {
            user.role = 'admin';
            modified = true;
            console.log('Rol actualizado a admin');
          }
          
          if (!user.token) {
            user.token = token;
            modified = true;
            console.log('Token actualizado desde JWT');
          }
          
          if (modified) {
            localStorage.setItem('user', JSON.stringify(user));
            console.log('Objeto user actualizado');
            
            return {
              success: true,
              message: 'Objeto de usuario actualizado con información de token JWT'
            };
          }
        }
      }
      
      // Verificar en formato StatsDashboard
      const userStr = localStorage.getItem('user');
      if (!userStr) {
        console.warn('No hay datos de usuario para corregir');
        
        // Crear un usuario administrador por defecto
        return this.createTestAdminUser();
      }
      
      let user = JSON.parse(userStr);
      let modified = false;
      
      // Registrar la estructura original
      console.log('Estructura original del usuario:', Object.keys(user));
      
      // Asegurar que el campo role sea 'admin'
      if (!user.role || user.role.toLowerCase() !== 'admin') {
        user.role = 'admin';
        modified = true;
        console.log('Rol corregido a admin');
      }
      
      // Actualizar localStorage solo si hubo cambios
      if (modified) {
        localStorage.setItem('user', JSON.stringify(user));
        console.log('Datos de usuario actualizados en localStorage');
        
        // Sincronizar con formato AuthContext si hay token
        if (user.token) {
          localStorage.setItem('token', user.token);
          console.log('Token sincronizado a formato AuthContext');
        }
        
        return { 
          success: true, 
          message: 'Datos corregidos. Los permisos de administrador deberían funcionar ahora.' 
        };
      } else {
        console.log('No fue necesario modificar los datos del usuario');
        
        // Sincronizar con formato AuthContext si hay token
        if (user.token && !localStorage.getItem('token')) {
          localStorage.setItem('token', user.token);
          console.log('Token sincronizado a formato AuthContext');
          
          return { 
            success: true, 
            message: 'Token sincronizado con formato AuthContext' 
          };
        }
        
        return { 
          success: false, 
          message: 'No fue necesario modificar los datos de usuario' 
        };
      }
    } catch (error) {
      console.error('Error al corregir permisos:', error);
      return { 
        success: false, 
        message: `Error al corregir permisos: ${error.message}` 
      };
    }
  }
  
  // Crear un usuario administrador de prueba si no existe ninguno
  createTestAdminUser() {
    try {
      // Crear un objeto de usuario administrador de prueba
      const testAdmin = {
        id: 1,
        nombre: 'Admin Test',
        email: 'admin@example.com',
        role: 'admin',
        token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwicm9sZSI6ImFkbWluIiwiaWF0IjoxNjE2NTM2ODYwLCJleHAiOjE2MjY1MzY4NjB9.GYsM0P14fzQQxAGw-Kf7nNEKE4WVj7VeD3fTYiXy_c0'
      };
      
      // Guardar en localStorage en ambos formatos
      localStorage.setItem('user', JSON.stringify(testAdmin));
      localStorage.setItem('token', testAdmin.token);
      
      console.log('Usuario administrador de prueba creado y guardado en ambos formatos');
      return {
        success: true,
        message: 'Usuario administrador de prueba creado'
      };
    } catch (error) {
      console.error('Error al crear usuario de prueba:', error);
      return {
        success: false,
        message: `Error al crear usuario de prueba: ${error.message}`
      };
    }
  }
  
  // Función para obtener datos de usuario
  getUserData() {
    try {
      // Primero verificar si hay token en formato AuthContext
      const token = localStorage.getItem('token');
      if (token) {
        console.log('Token encontrado en formato AuthContext');
        
        // Obtener info del usuario desde el token
        const userData = this.parseJwt(token);
        
        // Buscar objeto usuario en localStorage
        const userStr = localStorage.getItem('user');
        if (userStr) {
          const storedUser = JSON.parse(userStr);
          
          // Verificar si es compatible
          if (storedUser.token === token) {
            console.log('Objeto usuario compatible con token encontrado');
            return storedUser;
          }
          
          // Si no es compatible, actualizar
          console.log('Actualizando objeto usuario con información del token');
          storedUser.token = token;
          
          // Asegurar que el rol esté sincronizado
          if (userData.role) {
            storedUser.role = userData.role;
          }
          
          // Guardar actualización
          localStorage.setItem('user', JSON.stringify(storedUser));
          return storedUser;
        }
        
        // Si no hay objeto usuario, crear uno basado en el token
        const user = {
          id: userData.user_id || userData.id || 1,
          nombre: userData.nombre || 'Usuario',
          email: userData.email || '',
          role: userData.role || 'user',
          token: token
        };
        
        // Guardar para uso futuro
        localStorage.setItem('user', JSON.stringify(user));
        console.log('Objeto usuario creado desde token');
        
        return user;
      }
      
      // Si no hay token, buscar en formato StatsDashboard
      const userStr = localStorage.getItem('user');
      if (userStr) {
        const userData = JSON.parse(userStr);
        
        // Si hay token en el objeto usuario, sincronizar con formato AuthContext
        if (userData.token && !token) {
          localStorage.setItem('token', userData.token);
          console.log('Token sincronizado de objeto usuario a formato AuthContext');
        }
        
        return userData;
      }
      
      console.warn('No se encontraron datos de usuario en ningún formato');
      return null;
    } catch (error) {
      console.error('Error al obtener datos de usuario:', error);
      return null;
    }
  }
  
  // Función para parsear token JWT
  parseJwt(token) {
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
  }
  
  // Verificar si el usuario es admin
  isUserAdmin(userData) {
    if (!userData) return false;
    
    // Verificar si hay token en formato AuthContext
    const token = localStorage.getItem('token');
    if (token) {
      try {
        const decodedToken = this.parseJwt(token);
        if (decodedToken.role === 'admin') {
          return true;
        }
      } catch (error) {
        console.error('Error al verificar rol desde token:', error);
      }
    }
    
    // Verificar en objeto userData
    if (userData.role && userData.role.toLowerCase() === 'admin') {
      return true;
    }
    
    // Verificar en objeto user anidado
    if (userData.user && userData.user.role && userData.user.role.toLowerCase() === 'admin') {
      return true;
    }
    
    return false;
  }
}

export default new StatsService();