import React, { useState, useEffect, useCallback } from 'react';
import { 
  FaEnvelope, 
  FaEnvelopeOpen, 
  FaTrash, 
  FaReply, 
  FaStar, 
  FaSearch, 
  FaArrowUp, 
  FaArrowDown,
  FaCheckCircle,
  FaExclamationTriangle,
  FaTimes
} from 'react-icons/fa';
import axios from 'axios';
import './MessagesAdmin.css';

const MessagesAdmin = () => {
  const [messages, setMessages] = useState([]);
  const [filteredMessages, setFilteredMessages] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filter, setFilter] = useState('all');
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedMessage, setSelectedMessage] = useState(null);
  const [sortConfig, setSortConfig] = useState({
    key: 'date',
    direction: 'desc'
  });
  const [replyText, setReplyText] = useState('');
  const [sendingReply, setSendingReply] = useState(false);
  const [notification, setNotification] = useState(null);
  
  const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000';

  const fetchMessages = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await axios.get(`${API_URL}/api/admin/messages`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      // Manejo robusto de la respuesta
      let messagesData = [];
      
      if (response.data && response.data.success) {
        // Caso 1: Respuesta con estructura {success, data}
        messagesData = response.data.data || [];
      } else if (Array.isArray(response.data)) {
        // Caso 2: Respuesta es directamente un array
        messagesData = response.data;
      }
      
      const formattedMessages = messagesData.map(msg => ({
        id: msg.id || msg.ID,
        name: msg.name || msg.Name,
        email: msg.email || msg.Email,
        phone: msg.phone || msg.Phone || '',
        message: msg.message || msg.Message,
        read: msg.read || msg.Read || false,
        starred: msg.starred || msg.Starred || false,
        date: (msg.createdAt || msg.CreatedAt || new Date().toISOString()).split('T')[0],
        time: (msg.createdAt || msg.CreatedAt || new Date().toISOString()).split('T')[1]?.substring(0, 5) || '00:00'
      }));
      
      setMessages(formattedMessages);
    } catch (err) {
      console.error('Error al obtener mensajes:', err);
      setError('Error al obtener mensajes. Por favor, intente nuevamente.');
      
      // Datos de ejemplo para desarrollo
      const exampleMessages = [
        { 
          id: 1, 
          name: 'Juan Pérez', 
          email: 'juan@example.com', 
          phone: '+5493511234567',
          message: '¿Tienen cursos sobre TypeScript? Me interesa conocer si ofrecen algún curso especializado en TypeScript para desarrollo web frontend.',
          read: false,
          starred: false,
          date: '2023-09-15',
          time: '14:30'
        },
        { 
          id: 2, 
          name: 'María García', 
          email: 'maria@example.com', 
          phone: '+5493517654321',
          message: 'Necesito ayuda con mi suscripción. Hace dos días realicé el pago pero aún no tengo acceso.',
          read: true,
          starred: true,
          date: '2023-09-14',
          time: '10:15'
        }
      ];
      
      setMessages(exampleMessages);
    } finally {
      setIsLoading(false);
    }
  }, [API_URL]);

  useEffect(() => {
    fetchMessages();
  }, [fetchMessages]);

  useEffect(() => {
    applyFilters();
  }, [filter, messages, searchTerm, sortConfig]);

  // Función para mostrar una notificación
  const showNotification = (type, message, duration = 5000) => {
    setNotification({ type, message });
    
    // Auto cerrar después de duration ms
    if (duration) {
      setTimeout(() => {
        setNotification(null);
      }, duration);
    }
  };

  const applyFilters = () => {
    let result = [...messages];
    
    if (filter === 'unread') {
      result = result.filter(message => !message.read);
    } else if (filter === 'read') {
      result = result.filter(message => message.read);
    } else if (filter === 'starred') {
      result = result.filter(message => message.starred);
    }
    
    if (searchTerm) {
      const searchLower = searchTerm.toLowerCase();
      result = result.filter(
        message =>
          message.name.toLowerCase().includes(searchLower) ||
          message.email.toLowerCase().includes(searchLower) ||
          (message.phone && message.phone.toLowerCase().includes(searchLower)) ||
          message.message.toLowerCase().includes(searchLower)
      );
    }
    
    if (sortConfig.key) {
      result.sort((a, b) => {
        if (a[sortConfig.key] < b[sortConfig.key]) {
          return sortConfig.direction === 'asc' ? -1 : 1;
        }
        if (a[sortConfig.key] > b[sortConfig.key]) {
          return sortConfig.direction === 'asc' ? 1 : -1;
        }
        return 0;
      });
    }
    
    setFilteredMessages(result);
  };

  const handleFilterChange = (newFilter) => {
    setFilter(newFilter);
  };

  const handleSearchChange = (e) => {
    setSearchTerm(e.target.value);
  };

  const handleSort = (key) => {
    setSortConfig({
      key,
      direction:
        sortConfig.key === key && sortConfig.direction === 'asc'
          ? 'desc'
          : 'asc'
    });
  };

  const toggleReadStatus = async (id) => {
    try {
      await axios.patch(`${API_URL}/api/admin/messages/${id}/read`, null, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      setMessages(
        messages.map(msg =>
          msg.id === id ? { ...msg, read: !msg.read } : msg
        )
      );
    } catch (err) {
      console.error('Error al actualizar estado del mensaje:', err);
      showNotification('error', 'Error al actualizar el estado del mensaje');
    }
  };

  const toggleStarred = async (id) => {
    try {
      await axios.patch(`${API_URL}/api/admin/messages/${id}/star`, null, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      setMessages(
        messages.map(msg =>
          msg.id === id ? { ...msg, starred: !msg.starred } : msg
        )
      );
    } catch (err) {
      console.error('Error al actualizar estado de estrella:', err);
      showNotification('error', 'Error al actualizar el mensaje');
    }
  };

  const deleteMessage = async (id) => {
    if (window.confirm('¿Estás seguro de que quieres eliminar este mensaje?')) {
      try {
        await axios.delete(`${API_URL}/api/admin/messages/${id}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        
        setMessages(messages.filter(msg => msg.id !== id));
        if (selectedMessage && selectedMessage.id === id) {
          setSelectedMessage(null);
        }
        
        showNotification('success', 'Mensaje eliminado correctamente');
      } catch (err) {
        console.error('Error al eliminar mensaje:', err);
        showNotification('error', 'Error al eliminar el mensaje');
      }
    }
  };

  const viewMessage = (message) => {
    if (!message.read) {
      toggleReadStatus(message.id);
    }
    setSelectedMessage(message);
  };

  const sendReply = async (e) => {
    e.preventDefault();
    
    if (!replyText.trim()) {
      showNotification('error', 'Por favor, escribe un mensaje antes de enviar');
      return;
    }
    
    setSendingReply(true);
    
    try {
      // Incluir todos los headers necesarios
      const config = {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        }
      };
      
      // Asegurarse de que el cuerpo de la solicitud sea un objeto válido
      const requestBody = { message: replyText };
      
      // Usar try/catch específico para la petición
      try {
        const response = await axios.post(
          `${API_URL}/api/admin/messages/${selectedMessage.id}/reply`, 
          requestBody,
          config
        );
        
        console.log('Respuesta del servidor:', response.data);
        
        // Verificar la respuesta
        if (response.data && response.data.success) {
          // Mostrar notificación de éxito
          showNotification('success', `Respuesta enviada a ${selectedMessage.email}`);
          
          setReplyText('');
          
          // Actualizar mensaje como leído si no lo estaba
          setMessages(
            messages.map(msg =>
              msg.id === selectedMessage.id ? { ...msg, read: true } : msg
            )
          );
          
          // Cerrar el detalle del mensaje después de un breve retraso
          setTimeout(() => {
            setSelectedMessage(null);
          }, 1500);
        } else {
          throw new Error('La respuesta no indica éxito');
        }
      } catch (apiError) {
        console.error('Error en la petición API:', apiError);
        
        // Información detallada para depuración
        if (apiError.response) {
          console.error('Datos de respuesta de error:', apiError.response.data);
          console.error('Estado de error:', apiError.response.status);
          console.error('Headers de error:', apiError.response.headers);
        }
        
        throw apiError; // Re-lanzar para el manejador principal
      }
    } catch (err) {
      console.error('Error completo al enviar respuesta:', err);
      
      // Mensaje más detallado para el usuario en una notificación
      const errorMsg = err.response && err.response.data && err.response.data.error
        ? `Error: ${err.response.data.error}`
        : 'Error al enviar el correo electrónico de respuesta. Por favor, verifique la configuración del servidor de correo.';
      
      showNotification('error', errorMsg);
    } finally {
      setSendingReply(false);
    }
  };

  const LoadingOverlay = () => (
    <div className="loading-overlay">
      <div className="spinner"></div>
    </div>
  );

  const formatDate = (dateStr) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('es-ES', {
      day: '2-digit',
      month: 'short',
      year: 'numeric'
    });
  };

  return (
    <div className="messages-admin">
      {isLoading && <LoadingOverlay />}
      
      {notification && (
        <div className={`notification notification-${notification.type}`}>
          <div className="notification-icon">
            {notification.type === 'success' ? <FaCheckCircle /> : <FaExclamationTriangle />}
          </div>
          <div className="notification-message">{notification.message}</div>
          <button className="notification-close" onClick={() => setNotification(null)}>
            <FaTimes />
          </button>
        </div>
      )}
      
      <div className="section-header">
        <h2 className="section-title">Mensajes de Contacto</h2>
        <div className="message-stats">
          <div className="stat">
            <span className="stat-value">{messages.length}</span>
            <span className="stat-label">Total</span>
          </div>
          <div className="stat">
            <span className="stat-value">{messages.filter(m => !m.read).length}</span>
            <span className="stat-label">No leídos</span>
          </div>
        </div>
      </div>
      
      <div className="messages-header">
        <div className="message-filters">
          <button 
            className={`filter-btn ${filter === 'all' ? 'active' : ''}`}
            onClick={() => handleFilterChange('all')}
          >
            Todos
          </button>
          <button 
            className={`filter-btn ${filter === 'unread' ? 'active' : ''}`}
            onClick={() => handleFilterChange('unread')}
          >
            <FaEnvelope /> No leídos
          </button>
          <button 
            className={`filter-btn ${filter === 'read' ? 'active' : ''}`}
            onClick={() => handleFilterChange('read')}
          >
            <FaEnvelopeOpen /> Leídos
          </button>
          <button 
            className={`filter-btn ${filter === 'starred' ? 'active' : ''}`}
            onClick={() => handleFilterChange('starred')}
          >
            <FaStar /> Destacados
          </button>
        </div>
        
        <div className="search-bar">
          <FaSearch />
          <input 
            type="text" 
            placeholder="Buscar en mensajes..." 
            value={searchTerm}
            onChange={handleSearchChange}
          />
        </div>
      </div>
      
      {error && (
        <div className="error-message">
          {error}
        </div>
      )}
      
      {selectedMessage ? (
        <div className="message-detail">
          <div className="message-detail-header">
            <button className="back-btn" onClick={() => setSelectedMessage(null)}>
              <span>←</span> Volver a la lista
            </button>
            <div className="message-actions">
              <button 
                className={`action-btn ${selectedMessage.starred ? 'active' : ''}`}
                onClick={() => toggleStarred(selectedMessage.id)}
              >
                <FaStar />
              </button>
              <button 
                className="action-btn danger"
                onClick={() => deleteMessage(selectedMessage.id)}
              >
                <FaTrash />
              </button>
            </div>
          </div>
          
          <div className="message-detail-content">
            <div className="sender-info">
              <div className="sender-avatar">
                {selectedMessage.name.charAt(0).toUpperCase()}
              </div>
              <div>
                <h3>{selectedMessage.name}</h3>
                <a href={`mailto:${selectedMessage.email}`}>
                  {selectedMessage.email}
                </a>
                {selectedMessage.phone && (
                  <div className="sender-phone">{selectedMessage.phone}</div>
                )}
                <div className="message-date">
                  {formatDate(selectedMessage.date)} - {selectedMessage.time}
                </div>
              </div>
            </div>
            
            <div className="message-body">
              {selectedMessage.message}
            </div>
            
            <form className="reply-form" onSubmit={sendReply}>
              <h4><FaReply /> Responder</h4>
              <textarea
                placeholder="Escribe tu respuesta aquí..."
                value={replyText}
                onChange={(e) => setReplyText(e.target.value)}
                required
                disabled={sendingReply}
              />
              <div className="form-actions">
                <button 
                  type="button" 
                  className="btn-cancel" 
                  onClick={() => setSelectedMessage(null)}
                  disabled={sendingReply}
                >
                  Cancelar
                </button>
                <button 
                  type="submit" 
                  className="btn-primary"
                  disabled={sendingReply}
                >
                  {sendingReply ? 'Enviando...' : 'Enviar Respuesta'}
                </button>
              </div>
            </form>
          </div>
        </div>
      ) : (
        <div className="messages-container">
          <div className="messages-table-header">
            <div className="message-column sender" onClick={() => handleSort('name')}>
              Remitente
              {sortConfig.key === 'name' && (
                sortConfig.direction === 'asc' ? <FaArrowUp /> : <FaArrowDown />
              )}
            </div>
            <div className="message-column preview">
              Mensaje
            </div>
            <div className="message-column date" onClick={() => handleSort('date')}>
              Fecha
              {sortConfig.key === 'date' && (
                sortConfig.direction === 'asc' ? <FaArrowUp /> : <FaArrowDown />
              )}
            </div>
            <div className="message-column actions">
              Acciones
            </div>
          </div>
          
          <div className="messages-list">
            {filteredMessages.length === 0 ? (
              <div className="no-messages">
                <p>No hay mensajes que coincidan con los criterios de búsqueda.</p>
              </div>
            ) : (
              filteredMessages.map(msg => (
                <div 
                  key={msg.id} 
                  className={`message-row ${!msg.read ? 'unread' : ''} ${msg.starred ? 'starred' : ''}`}
                  onClick={() => viewMessage(msg)}
                >
                  <div className="message-column sender">
                    <div className="sender-avatar">
                      {msg.name.charAt(0).toUpperCase()}
                    </div>
                    <div>
                      <div className="sender-name">{msg.name}</div>
                      <div className="sender-email">{msg.email}</div>
                    </div>
                  </div>
                  
                  <div className="message-column preview">
                    <div className="message-preview">
                      {msg.message.substring(0, 80)}
                      {msg.message.length > 80 ? '...' : ''}
                    </div>
                  </div>
                  
                  <div className="message-column date">
                    {formatDate(msg.date)}
                    <div className="message-time">{msg.time}</div>
                  </div>
                  
                  <div 
                    className="message-column actions" 
                    onClick={(e) => e.stopPropagation()}
                  >
                    <button 
                      title={msg.read ? "Marcar como no leído" : "Marcar como leído"}
                      className="action-btn"
                      onClick={() => toggleReadStatus(msg.id)}
                    >
                      {msg.read ? <FaEnvelopeOpen /> : <FaEnvelope />}
                    </button>
                    
                    <button 
                      title={msg.starred ? "Quitar estrella" : "Destacar"}
                      className={`action-btn ${msg.starred ? 'active' : ''}`}
                      onClick={() => toggleStarred(msg.id)}
                    >
                      <FaStar />
                    </button>
                    
                    <button 
                      title="Eliminar mensaje"
                      className="action-btn danger"
                      onClick={() => deleteMessage(msg.id)}
                    >
                      <FaTrash />
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default MessagesAdmin;