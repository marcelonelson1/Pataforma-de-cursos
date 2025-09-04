// services/CursosService.js - Actualizado para integración con Payment Service
import axios from 'axios';

// Configuración de URLs - Separar Course Service y Payment Service
const COURSE_API_URL = process.env.REACT_APP_COURSE_API_URL || 'http://localhost:8003';
const PAYMENT_API_URL = process.env.REACT_APP_PAYMENT_API_URL || 'http://localhost:8002';

// Crear instancia para Course Service
const courseApi = axios.create({
  baseURL: COURSE_API_URL,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// Crear instancia para Payment Service
const paymentApi = axios.create({
  baseURL: PAYMENT_API_URL,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// Interceptor para añadir el token a las peticiones del Course Service
courseApi.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
}, error => {
  return Promise.reject(error);
});

// Interceptor para añadir el token a las peticiones del Payment Service
paymentApi.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
}, error => {
  return Promise.reject(error);
});

// Interceptor para manejar errores de respuesta
[courseApi, paymentApi].forEach(api => {
  api.interceptors.response.use(
    response => response,
    error => {
      console.error('API Error:', error.response?.data || error.message);
      return Promise.reject(error);
    }
  );
});

// Funciones del servicio
const CursosService = {
  // Obtener todos los cursos
  getCursos: async (soloPublicados = true) => {
    try {
      console.log('Obteniendo cursos del Course Service...');
      
      const response = await courseApi.get('/api/courses');
      console.log('Respuesta completa del servidor:', response);
      
      let data = response.data;
      console.log('Datos extraídos:', data);
      
      let courses = [];
      
      // Manejar la estructura de respuesta correcta
      if (data && data.success === true && data.data) {
        if (Array.isArray(data.data.courses)) {
          courses = data.data.courses;
          console.log('Cursos extraídos del formato data.courses:', courses);
        } else if (Array.isArray(data.data)) {
          courses = data.data;
          console.log('Cursos extraídos del formato data:', courses);
        }
      } else if (data && Array.isArray(data.courses)) {
        courses = data.courses;
        console.log('Cursos extraídos del formato courses:', courses);
      } else if (Array.isArray(data)) {
        courses = data;
        console.log('Respuesta directa como array:', courses);
      }
      
      console.log('Cursos procesados final:', courses);
      
      // Si se solicita solo cursos publicados, filtrarlos
      if (soloPublicados && Array.isArray(courses)) {
        const cursosPublicados = courses.filter(curso => 
          curso.estado === 'Publicado' || curso.estado === 'publicado'
        );
        console.log('Cursos publicados filtrados:', cursosPublicados);
        return cursosPublicados;
      }
      
      return Array.isArray(courses) ? courses : [];
    } catch (error) {
      console.error('Error detallado al obtener cursos:', error);
      console.error('Status:', error.response?.status);
      console.error('Data:', error.response?.data);
      console.error('Message:', error.message);
      throw error;
    }
  },

  // Obtener un curso por su ID
  getCursoById: async (id) => {
    try {
      console.log(`Obteniendo curso con ID: ${id} desde /api/courses/${id}`);
      const response = await courseApi.get(`/api/courses/${id}`);
      console.log('Respuesta completa del backend:', response.data);
      
      let courseData = response.data;
      if (courseData && courseData.success && courseData.data) {
        courseData = courseData.data;
      }
      
      console.log('Datos del curso procesados:', courseData);
      console.log('¿Tiene capítulos?:', Array.isArray(courseData.capitulos), courseData.capitulos?.length || 0);
      
      return courseData;
    } catch (error) {
      console.error(`Error al obtener curso con ID ${id}:`, error);
      throw error;
    }
  },

  // Obtener capítulos de un curso
  getCapitulosByCurso: async (cursoId) => {
    try {
      const response = await courseApi.get(`/api/chapters/course/${cursoId}`);
      
      let chaptersData = response.data;
      if (chaptersData && chaptersData.success && chaptersData.chapters) {
        chaptersData = chaptersData.chapters;
      } else if (chaptersData && chaptersData.data) {
        chaptersData = chaptersData.data;
      }
      
      return chaptersData;
    } catch (error) {
      console.error(`Error al obtener capítulos del curso ${cursoId}:`, error);
      throw error;
    }
  },

  // **FUNCIÓN PRINCIPAL ACTUALIZADA: Verificar acceso combinando Course Service y Payment Service**
  verificarAccesoCurso: async (cursoId) => {
    try {
      console.log(`Verificando acceso al curso ${cursoId}...`);
      
      // Primero obtener información del curso del Course Service
      let cursoInfo;
      try {
        const cursoResponse = await courseApi.get(`/api/courses/${cursoId}`);
        cursoInfo = cursoResponse.data;
        if (cursoInfo && cursoInfo.success && cursoInfo.data) {
          cursoInfo = cursoInfo.data;
        }
      } catch (error) {
        console.error('Error al obtener información del curso:', error);
        throw error;
      }

      // Si el curso es gratuito, el usuario tiene acceso automáticamente
      if (cursoInfo && (cursoInfo.gratuito || cursoInfo.precio === 0)) {
        console.log('Curso gratuito - acceso permitido');
        return { has_access: true, is_free: true };
      }

      // Verificar si el usuario está autenticado
      const token = localStorage.getItem('token');
      if (!token) {
        console.log('Usuario no autenticado');
        return { has_access: false, requires_auth: true };
      }

      // **VERIFICAR PAGO EN PAYMENT SERVICE**
      try {
        const paymentResponse = await paymentApi.get(`/api/pagos/${cursoId}`);
        console.log('Respuesta del Payment Service:', paymentResponse.data);
        
        const paymentData = paymentResponse.data;
        
        // Analizar respuesta del Payment Service
        let tieneAcceso = false;
        let estadoPago = 'no_pagado';

        if (paymentData && paymentData.success) {
          estadoPago = paymentData.estado || 'no_pagado';
          
          // Si hay un pago aprobado, tiene acceso
          if (estadoPago === 'aprobado') {
            tieneAcceso = true;
          }
        } else if (paymentData && paymentData.estado) {
          estadoPago = paymentData.estado;
          if (estadoPago === 'aprobado') {
            tieneAcceso = true;
          }
        }

        return { 
          has_access: tieneAcceso, 
          payment_status: estadoPago,
          payment_info: paymentData.pago || paymentData.data || null
        };

      } catch (paymentError) {
        console.error('Error al verificar pago:', paymentError);
        
        // Si es error 404, significa que no hay pago registrado
        if (paymentError.response && paymentError.response.status === 404) {
          return { has_access: false, payment_status: 'no_pagado' };
        }
        
        // Para otros errores, no asumir nada
        return { has_access: false, payment_status: 'error', error: paymentError.message };
      }

    } catch (error) {
      console.error(`Error al verificar acceso al curso ${cursoId}:`, error);
      
      // En caso de error general, no permitir acceso
      return { has_access: false, error: error.message };
    }
  },

  // **NUEVA FUNCIÓN: Verificar estado de pago específicamente**
  verificarEstadoPago: async (cursoId) => {
    try {
      const response = await paymentApi.get(`/api/pagos/${cursoId}`);
      return response.data;
    } catch (error) {
      console.error(`Error al verificar estado de pago del curso ${cursoId}:`, error);
      if (error.response && error.response.status === 404) {
        return { estado: 'no_pagado', success: true };
      }
      throw error;
    }
  },

  // **FUNCIÓN ACTUALIZADA: Procesar un pago usando Payment Service**
  procesarPago: async (datosPago) => {
    try {
      console.log('Procesando pago en Payment Service:', datosPago);
      const response = await paymentApi.post('/api/pagos', datosPago);
      return response.data;
    } catch (error) {
      console.error('Error al procesar pago:', error);
      throw error;
    }
  },

  // **NUEVA FUNCIÓN: Obtener historial de pagos del usuario**
  obtenerHistorialPagos: async () => {
    try {
      const response = await paymentApi.get('/api/pagos/user/history');
      return response.data;
    } catch (error) {
      console.error('Error al obtener historial de pagos:', error);
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
      const response = await courseApi.post('/api/progress/chapter/complete', data);
      return response.data;
    } catch (error) {
      console.error(`Error al marcar capítulo ${capituloId} como ${completado ? 'completado' : 'no completado'}:`, error);
      throw error;
    }
  },

  // Obtener el progreso del usuario en un curso
  getProgresoUsuario: async (cursoId) => {
    try {
      const response = await courseApi.get(`/api/progress/course/${cursoId}`);
      
      let progressData = response.data;
      if (progressData && progressData.success && progressData.data) {
        progressData = progressData.data;
      }
      
      return progressData;
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
      const response = await courseApi.post('/api/progress/last-chapter', data);
      return response.data;
    } catch (error) {
      console.error(`Error al guardar último capítulo visto ${capituloId}:`, error);
      throw error;
    }
  },

  // Obtener resumen de progreso del usuario
  getResumenProgreso: async () => {
    try {
      const response = await courseApi.get('/api/progress/user/summary');
      
      let summaryData = response.data;
      if (summaryData && summaryData.success && summaryData.data) {
        summaryData = summaryData.data;
      }
      
      return summaryData;
    } catch (error) {
      console.error('Error al obtener resumen de progreso:', error);
      throw error;
    }
  },

  // Actualizar tiempo visto de un capítulo
  actualizarTiempoVisto: async (cursoId, capituloId, tiempoVisto) => {
    try {
      const data = {
        curso_id: cursoId,
        capitulo_id: capituloId,
        tiempo_visto: tiempoVisto
      };
      const response = await courseApi.post('/api/progress/chapter/watch-time', data);
      return response.data;
    } catch (error) {
      console.error('Error al actualizar tiempo visto:', error);
      throw error;
    }
  },

  // Obtener estadísticas del curso
  getEstadisticasCurso: async (cursoId) => {
    try {
      const response = await courseApi.get(`/api/progress/course/${cursoId}/stats`);
      
      let statsData = response.data;
      if (statsData && statsData.success && statsData.data) {
        statsData = statsData.data;
      }
      
      return statsData;
    } catch (error) {
      console.error(`Error al obtener estadísticas del curso ${cursoId}:`, error);
      throw error;
    }
  },

  // **NUEVAS FUNCIONES PARA INTEGRACIÓN CON PAYMENT SERVICE**

  // Verificar si un curso específico está pagado por el usuario actual
  esCursoPagado: async (cursoId) => {
    try {
      const verificacion = await CursosService.verificarAccesoCurso(cursoId);
      return verificacion.has_access === true;
    } catch (error) {
      console.error(`Error al verificar si curso ${cursoId} está pagado:`, error);
      return false;
    }
  },

  // Obtener información completa de acceso (curso + pago)
  getInformacionAccesoCompleta: async (cursoId) => {
    try {
      // Obtener información del curso
      const cursoInfo = await CursosService.getCursoById(cursoId);
      
      // Verificar acceso (incluye verificación de pago)
      const accesoInfo = await CursosService.verificarAccesoCurso(cursoId);
      
      return {
        curso: cursoInfo,
        acceso: accesoInfo,
        es_gratuito: cursoInfo.gratuito || cursoInfo.precio === 0,
        tiene_acceso: accesoInfo.has_access,
        estado_pago: accesoInfo.payment_status || 'no_verificado'
      };
    } catch (error) {
      console.error(`Error al obtener información completa de acceso para curso ${cursoId}:`, error);
      throw error;
    }
  },

  // Validar acceso antes de mostrar contenido
  validarAccesoContenido: async (cursoId) => {
    try {
      const informacion = await CursosService.getInformacionAccesoCompleta(cursoId);
      
      // Si es gratuito, permitir acceso
      if (informacion.es_gratuito) {
        return { permitido: true, razon: 'curso_gratuito' };
      }
      
      // Si tiene acceso pagado, permitir
      if (informacion.tiene_acceso) {
        return { permitido: true, razon: 'pago_aprobado' };
      }
      
      // Verificar si requiere autenticación
      if (informacion.acceso.requires_auth) {
        return { permitido: false, razon: 'requiere_login' };
      }
      
      // Por defecto, requiere pago
      return { 
        permitido: false, 
        razon: 'requiere_pago',
        precio: informacion.curso.precio,
        estado_pago: informacion.estado_pago
      };
      
    } catch (error) {
      console.error(`Error al validar acceso al contenido del curso ${cursoId}:`, error);
      return { permitido: false, razon: 'error', error: error.message };
    }
  }
};

export default CursosService;