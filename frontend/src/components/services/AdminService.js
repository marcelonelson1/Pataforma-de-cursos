// src/services/AdminService.js
import axios from 'axios';

class AdminService {
  constructor() {
    this.apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
    this.api = axios.create({
      baseURL: `${this.apiUrl}/api/admin`,
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
        if (error.response && (error.response.status === 401 || error.response.status === 403)) {
          // Si hay error de autenticación, redirigir al login
          localStorage.removeItem('token');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // Obtener estadísticas para el dashboard
  getStats() {
    return this.api.get('/stats');
  }

  // Obtener datos para el dashboard
  getDashboard() {
    return this.api.get('/dashboard');
  }

  // Obtener logs de actividad
  getActivityLog() {
    return this.api.get('/activity-log');
  }

  // Obtener estadísticas de ventas
  getSalesStats() {
    return this.api.get('/sales-stats');
  }

  // Listar usuarios
  getUsers() {
    return this.api.get('/users');
  }

  // Obtener usuario por ID
  getUserById(id) {
    return this.api.get(`/users/${id}`);
  }

  // Actualizar usuario
  updateUser(id, userData) {
    return this.api.put(`/users/${id}`, userData);
  }

  // Eliminar usuario
  deleteUser(id) {
    return this.api.delete(`/users/${id}`);
  }

  // Cambiar rol de usuario
  changeUserRole(id, role) {
    return this.api.put(`/users/${id}/role`, { role });
  }

  // Búsqueda general (cursos, usuarios, etc.)
  search(query) {
    return this.api.get(`/search?q=${encodeURIComponent(query)}`);
  }

  // Obtener notificaciones
  getNotifications() {
    return this.api.get('/notifications');
  }

  // Marcar notificación como leída
  markNotificationAsRead(id) {
    return this.api.put(`/notifications/${id}/read`);
  }

  // Obtener mensajes
  getMessages() {
    return this.api.get('/messages');
  }
}

const adminService = new AdminService();
export default adminService;