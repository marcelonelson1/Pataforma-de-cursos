import axios from 'axios';

const PORTFOLIO_API_URL = process.env.REACT_APP_PORTFOLIO_API_URL || 'http://localhost:8005';

// Crear instancia para Portfolio Service
const portfolioApi = axios.create({
  baseURL: PORTFOLIO_API_URL,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// Interceptor para añadir token a las peticiones
portfolioApi.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  
  // Si estamos enviando FormData, remover Content-Type para que axios lo establezca automáticamente
  if (config.data instanceof FormData) {
    delete config.headers['Content-Type'];
  }
  
  return config;
}, error => {
  return Promise.reject(error);
});

// Interceptor para manejar errores de respuesta
portfolioApi.interceptors.response.use(
  response => response,
  error => {
    console.error('Portfolio API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

const PortfolioService = {
  // Obtener todos los proyectos del portfolio (publico)
  getAllProjects: () => {
    return portfolioApi.get('/api/portfolio');
  },

  // Obtener todos los proyectos para admin
  getAllProjectsAdmin: () => {
    return portfolioApi.get('/api/admin/portfolio');
  },

  // Obtener un proyecto específico por ID
  getProjectById: (projectId) => {
    return portfolioApi.get(`/api/portfolio/${projectId}`);
  },

  // Obtener proyectos por categoría
  getProjectsByCategory: (category) => {
    return portfolioApi.get(`/api/portfolio/category/${category}`);
  },

  // Obtener categorías disponibles
  getCategories: () => {
    return portfolioApi.get('/api/portfolio/categories');
  },

  // Crear un nuevo proyecto (admin)
  createProject: (projectData) => {
    return portfolioApi.post('/api/admin/portfolio', projectData);
  },

  // Actualizar un proyecto existente (admin)
  updateProject: (projectId, projectData) => {
    return portfolioApi.put(`/api/admin/portfolio/${projectId}`, projectData);
  },

  // Eliminar un proyecto (admin)
  deleteProject: (projectId) => {
    return portfolioApi.delete(`/api/admin/portfolio/${projectId}`);
  },

  // Subir imagen para un proyecto (admin)
  uploadProjectImage: (projectId, imageFile) => {
    const formData = new FormData();
    formData.append('project_id', projectId);
    formData.append('image', imageFile);

    return portfolioApi.post('/api/admin/portfolio/upload-image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
  },

  // Eliminar imagen de un proyecto (admin)
  deleteProjectImage: (projectId) => {
    return portfolioApi.delete(`/api/admin/portfolio/${projectId}/image`);
  },

  // Reordenar proyectos (admin)
  reorderProjects: (projectIds) => {
    return portfolioApi.post('/api/admin/portfolio/reorder', { project_ids: projectIds });
  },

  // Activar/desactivar proyecto (admin)
  toggleProjectStatus: (projectId) => {
    return portfolioApi.patch(`/api/admin/portfolio/${projectId}/toggle`);
  },

  // Obtener estadísticas del portfolio (admin)
  getPortfolioStats: () => {
    return portfolioApi.get('/api/admin/portfolio/stats');
  }
};

export default PortfolioService;