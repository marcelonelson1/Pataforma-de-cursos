// src/services/ProfileService.js
import axios from 'axios';

class ProfileService {
  constructor() {
    this.apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:8001';
    this.analyticsUrl = process.env.REACT_APP_ANALYTICS_API_URL || 'http://localhost:8007';
    this.api = axios.create({
      baseURL: `${this.apiUrl}/api`,
    });
    this.analyticsApi = axios.create({
      baseURL: `${this.analyticsUrl}/api`,
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

    // Interceptor para analytics API
    this.analyticsApi.interceptors.request.use(
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
  }

  // Obtener perfil del usuario
  getProfile() {
    return this.api.get('/users/profile');
  }

  // Actualizar perfil del usuario
  updateProfile(profileData) {
    return this.api.put('/users/profile', profileData);
  }

  // Cambiar contraseña
  changePassword(passwordData) {
    return this.api.post('/auth/change-password', passwordData);
  }

  // Subir imagen de perfil
  uploadProfileImage(formData) {
    return this.api.post('/users/profile/image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
  }

  // Obtener preferencias de notificación
  getNotificationSettings() {
    return this.api.get('/users/notification-settings');
  }

  // Actualizar preferencias de notificación
  updateNotificationSettings(settings) {
    return this.api.put('/users/notification-settings', settings);
  }
}

const profileService = new ProfileService();
export default profileService;