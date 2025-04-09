// src/services/AuthService.js
import axios from 'axios';

class AuthService {
  constructor() {
    this.api = axios.create({
      baseURL: process.env.REACT_APP_API_URL || 'http://localhost:5000/api',
      headers: {
        'Content-Type': 'application/json'
      }
    });

    // Interceptor para añadir token a todas las peticiones
    this.api.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token');
        if (token) {
          config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Interceptor para manejar errores de autenticación
    this.api.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response && error.response.status === 401) {
          // Token expirado o inválido
          localStorage.removeItem('token');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // Métodos para el perfil de usuario
  getProfile() {
    return this.api.get('/auth/profile');
  }

  updateProfile(profileData) {
    return this.api.put('/auth/profile', profileData);
  }

  changePassword(passwordData) {
    return this.api.post('/auth/change-password', passwordData);
  }

  uploadProfileImage(formData) {
    return this.api.post('/auth/profile/image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
  }

  // Métodos para las preferencias de notificación
  getNotificationSettings() {
    return this.api.get('/auth/notification-settings');
  }

  updateNotificationSettings(settings) {
    return this.api.put('/auth/notification-settings', settings);
  }

  // Método para refrescar el token
  refreshToken() {
    return this.api.post('/auth/refresh-token');
  }
}

const authService = new AuthService();
export default authService;