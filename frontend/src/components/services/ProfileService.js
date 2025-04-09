// src/services/ProfileService.js
import axios from 'axios';

class ProfileService {
  constructor() {
    this.apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
    this.api = axios.create({
      baseURL: `${this.apiUrl}/api`,
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
  }

  // Obtener perfil del usuario
  getProfile() {
    return this.api.get('/auth/profile');
  }

  // Actualizar perfil del usuario
  updateProfile(profileData) {
    return this.api.put('/auth/profile', profileData);
  }

  // Cambiar contraseña
  changePassword(passwordData) {
    return this.api.post('/auth/change-password', passwordData);
  }

  // Subir imagen de perfil
  uploadProfileImage(formData) {
    return this.api.post('/auth/profile/image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
  }

  // Obtener preferencias de notificación
  getNotificationSettings() {
    return this.api.get('/auth/notification-settings');
  }

  // Actualizar preferencias de notificación
  updateNotificationSettings(settings) {
    return this.api.put('/auth/notification-settings', settings);
  }
}

const profileService = new ProfileService();
export default profileService;