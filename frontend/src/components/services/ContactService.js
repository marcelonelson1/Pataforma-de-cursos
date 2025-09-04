import axios from 'axios';

const CONTACT_API_URL = process.env.REACT_APP_CONTACT_API_URL || 'http://localhost:8004';

// Crear instancia para Contact Service
const contactApi = axios.create({
  baseURL: CONTACT_API_URL,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// Interceptor para añadir token a las peticiones
contactApi.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
}, error => {
  return Promise.reject(error);
});

// Interceptor para manejar errores de respuesta
contactApi.interceptors.response.use(
  response => response,
  error => {
    console.error('Contact API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

const ContactService = {
  // Enviar mensaje de contacto (publico)
  sendContactMessage: (messageData) => {
    return contactApi.post('/api/contact', messageData);
  },

  // Obtener todos los mensajes de contacto (admin)
  getAllMessages: () => {
    return contactApi.get('/api/admin/messages');
  },

  // Obtener un mensaje específico por ID (admin)
  getMessageById: (messageId) => {
    return contactApi.get(`/api/admin/messages/${messageId}`);
  },

  // Marcar mensaje como leído (admin)
  markAsRead: (messageId) => {
    return contactApi.patch(`/api/admin/messages/${messageId}/read`);
  },

  // Marcar mensaje como no leído (admin)
  markAsUnread: (messageId) => {
    return contactApi.patch(`/api/admin/messages/${messageId}/unread`);
  },

  // Eliminar mensaje (admin)
  deleteMessage: (messageId) => {
    return contactApi.delete(`/api/admin/messages/${messageId}`);
  },

  // Responder mensaje via email (admin)
  replyToMessage: (messageId, replyData) => {
    return contactApi.post(`/api/admin/messages/${messageId}/reply`, replyData);
  },

  // Obtener estadísticas de mensajes (admin)
  getContactStats: () => {
    return contactApi.get('/api/admin/messages/stats');
  },

  // Obtener mensajes por estado (admin)
  getMessagesByStatus: (status) => {
    return contactApi.get(`/api/admin/messages/status/${status}`);
  },

  // Buscar mensajes (admin)
  searchMessages: (query) => {
    return contactApi.get(`/api/admin/messages/search?q=${encodeURIComponent(query)}`);
  },

  // Obtener mensajes filtrados por fecha (admin)
  getMessagesByDateRange: (startDate, endDate) => {
    return contactApi.get(`/api/admin/messages/date-range?start=${startDate}&end=${endDate}`);
  },

  // Marcar múltiples mensajes como leídos (admin)
  markMultipleAsRead: (messageIds) => {
    return contactApi.patch('/api/admin/messages/bulk/read', { message_ids: messageIds });
  },

  // Eliminar múltiples mensajes (admin)
  deleteMultipleMessages: (messageIds) => {
    return contactApi.delete('/api/admin/messages/bulk/delete', { 
      data: { message_ids: messageIds }
    });
  }
};

export default ContactService;