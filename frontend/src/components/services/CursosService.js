// Servicio para manejar todas las operaciones relacionadas con cursos
import axios from 'axios';

// Configuración de la URL base sin incluir '/api'
const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000';

// Crear una instancia de axios con la configuración base
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// Interceptor para añadir el token a las peticiones
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
}, error => {
  return Promise.reject(error);
});

// Funciones del servicio
const CursosService = {
  // Obtener todos los cursos
  getCursos: async (soloPublicados = true) => {
    try {
      // Añadimos el prefijo '/api' a todas las rutas
      const response = await api.get('/api/cursos');
      
      // Procesar la respuesta
      let data = response.data;
      
      // Si se solicita solo cursos publicados, filtrarlos
      if (soloPublicados && Array.isArray(data)) {
        data = data.filter(curso => 
          curso.estado === 'Publicado' || curso.estado === 'publicado'
        );
      }
      
      return data;
    } catch (error) {
      console.error('Error al obtener cursos:', error);
      throw error;
    }
  },

  // Obtener un curso por su ID
  getCursoById: async (id) => {
    try {
      const response = await api.get(`/api/cursos/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error al obtener curso con ID ${id}:`, error);
      throw error;
    }
  },

  // Obtener capítulos de un curso
  getCapitulosByCurso: async (cursoId) => {
    try {
      const response = await api.get(`/api/capitulos/curso/${cursoId}`);
      return response.data;
    } catch (error) {
      console.error(`Error al obtener capítulos del curso ${cursoId}:`, error);
      throw error;
    }
  },

  // Verificar estado de pago de un curso
  verificarPago: async (cursoId) => {
    try {
      const response = await api.get(`/api/pagos/${cursoId}`);
      return response.data;
    } catch (error) {
      console.error(`Error al verificar pago del curso ${cursoId}:`, error);
      if (error.response && error.response.status === 404) {
        return { estado: 'no_pagado' };
      }
      throw error;
    }
  },

  // Procesar un pago
  procesarPago: async (datosPago) => {
    try {
      const response = await api.post('/api/pagos', datosPago);
      return response.data;
    } catch (error) {
      console.error('Error al procesar pago:', error);
      throw error;
    }
  },

  // Marcar un capítulo como completado
  marcarCapituloCompletado: async (cursoId, capituloId, completado, progreso = 100) => {
    try {
      const data = {
        curso_id: cursoId,
        capitulo_id: capituloId,
        completado: completado,
        progreso: progreso
      };
      const response = await api.post('/api/progreso/capitulo/completado', data);
      return response.data;
    } catch (error) {
      console.error(`Error al marcar capítulo ${capituloId} como ${completado ? 'completado' : 'no completado'}:`, error);
      throw error;
    }
  },

  // Obtener el progreso del usuario en un curso
  getProgresoUsuario: async (cursoId) => {
    try {
      const response = await api.get(`/api/progreso/curso/${cursoId}`);
      return response.data;
    } catch (error) {
      console.error(`Error al obtener progreso del curso ID ${cursoId}:`, error);
      throw error;
    }
  },

  // Guardar el último capítulo visto
  guardarUltimoCapitulo: async (cursoId, capituloId) => {
    try {
      const data = {
        curso_id: cursoId,
        capitulo_id: capituloId
      };
      const response = await api.post('/api/progreso/ultimo-capitulo', data);
      return response.data;
    } catch (error) {
      console.error(`Error al guardar último capítulo visto ${capituloId}:`, error);
      throw error;
    }
  }
};

export default CursosService;