// src/components/services/HomeImagesService.js
import axios from 'axios';

class HomeImagesService {
  constructor() {
    this.baseURL = process.env.REACT_APP_API_URL || 'http://localhost:5000';
    this.apiUrl = `${this.baseURL}/api`;
  }

  // Obtener todas las imágenes (requiere autenticación)
  async getHomeImages() {
    const token = localStorage.getItem('token');
    return axios.get(`${this.apiUrl}/home-images`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
  }

  // Obtener imágenes públicas (no requiere autenticación)
  async getPublicHomeImages() {
    return axios.get(`${this.apiUrl}/home-images/public`);
  }

  // Subir una nueva imagen
  async uploadHomeImage(file, title, subtitle) {
    const token = localStorage.getItem('token');
    const formData = new FormData();
    formData.append('image', file);
    formData.append('title', title || '');
    formData.append('subtitle', subtitle || '');

    return axios.post(`${this.apiUrl}/home-images`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': `Bearer ${token}`
      }
    });
  }

  // Actualizar una imagen existente
  async updateHomeImage(id, updateData) {
    const token = localStorage.getItem('token');
    return axios.put(`${this.apiUrl}/home-images/${id}`, updateData, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
  }

  // Eliminar una imagen
  async deleteHomeImage(id) {
    const token = localStorage.getItem('token');
    return axios.delete(`${this.apiUrl}/home-images/${id}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
  }

  // Reordenar imágenes
  async reorderHomeImages(orderData) {
    const token = localStorage.getItem('token');
    return axios.patch(`${this.apiUrl}/home-images/reorder`, orderData, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
  }
}

const homeImagesService = new HomeImagesService();
export default homeImagesService;