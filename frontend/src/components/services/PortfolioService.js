import axios from 'axios';

const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';

const PortfolioService = {
  // Obtener todos los proyectos del portfolio
  getAllProjects: () => {
    return axios.get(`${apiUrl}/api/portfolio`);
  },

  // Obtener un proyecto específico por ID
  getProjectById: (projectId) => {
    return axios.get(`${apiUrl}/api/portfolio/${projectId}`);
  },

  // Crear un nuevo proyecto
  createProject: (projectData) => {
    const formData = new FormData();
    Object.keys(projectData).forEach(key => {
      formData.append(key, projectData[key]);
    });

    return axios.post(`${apiUrl}/api/portfolio`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
  },

  // Actualizar un proyecto existente
  updateProject: (projectId, projectData) => {
    const formData = new FormData();
    Object.keys(projectData).forEach(key => {
      if (projectData[key] !== null) {
        formData.append(key, projectData[key]);
      }
    });

    return axios.put(`${apiUrl}/api/portfolio/${projectId}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
  },

  // Eliminar un proyecto
  deleteProject: (projectId) => {
    return axios.delete(`${apiUrl}/api/portfolio/${projectId}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
  },

  // Obtener proyectos por categoría
  getProjectsByCategory: (category) => {
    return axios.get(`${apiUrl}/api/portfolio/category/${category}`);
  },
  
  // Cambiar orden de proyectos
  updateProjectOrder: (projectIds) => {
    return axios.post(`${apiUrl}/api/portfolio/reorder`, { projectIds }, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
  }
};

export default PortfolioService;