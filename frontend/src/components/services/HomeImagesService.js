// src/components/services/HomeImagesService.js
import axios from 'axios';

class HomeImagesService {
  constructor() {
    this.baseURL = process.env.REACT_APP_HOME_API_URL || 'http://localhost:8006';
    this.apiUrl = `${this.baseURL}/api`;
    
    // Crear instancia para Home Service
    this.homeApi = axios.create({
      baseURL: this.baseURL,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    });

    // Interceptor para añadir token a las peticiones
    this.homeApi.interceptors.request.use(config => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers['Authorization'] = `Bearer ${token}`;
      }
      return config;
    }, error => {
      return Promise.reject(error);
    });

    // Interceptor para manejar errores de respuesta
    this.homeApi.interceptors.response.use(
      response => response,
      error => {
        console.error('Home API Error:', error.response?.data || error.message);
        return Promise.reject(error);
      }
    );
  }

  // Obtener todas las imágenes (requiere autenticación)
  async getHomeImages() {
    return this.homeApi.get('/api/admin/home/images');
  }

  // Obtener imágenes públicas (no requiere autenticación)
  async getPublicHomeImages() {
    return this.homeApi.get('/api/home/images');
  }

  // Subir una nueva imagen (proceso de dos pasos)
  async uploadHomeImage(file, title, subtitle) {
    try {
      // Paso 1: Crear el registro de imagen con metadata
      const imageResponse = await this.homeApi.post('/api/admin/home/images', {
        title: title || '',
        subtitle: subtitle || '',
        is_active: true
      });

      if (!imageResponse.data.success) {
        throw new Error(imageResponse.data.error || 'Error al crear el registro de imagen');
      }

      // Obtener el ID de la imagen creada
      const imageId = imageResponse.data.data.id;

      // Paso 2: Subir el archivo de imagen
      const formData = new FormData();
      formData.append('image', file);
      formData.append('image_id', imageId.toString());

      const uploadResponse = await this.homeApi.post('/api/admin/home/upload-image', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      });

      return uploadResponse;
    } catch (error) {
      console.error('Error uploading home image:', error);
      throw error;
    }
  }

  // Actualizar una imagen existente
  async updateHomeImage(id, updateData) {
    return this.homeApi.put(`/api/admin/home/images/${id}`, updateData);
  }

  // Eliminar una imagen
  async deleteHomeImage(id) {
    return this.homeApi.delete(`/api/admin/home/images/${id}`);
  }

  // Reordenar imágenes
  async reorderHomeImages(orderData) {
    return this.homeApi.patch('/api/admin/home/images/reorder', orderData);
  }
}

const homeImagesService = new HomeImagesService();
export default homeImagesService;