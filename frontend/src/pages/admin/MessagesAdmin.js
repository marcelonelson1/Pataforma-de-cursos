import React, { useState, useEffect, useCallback } from 'react';
import { 
  FaEnvelope, 
  FaEnvelopeOpen, 
  FaTrash, 
  FaReply, 
  FaStar, 
  FaSearch, 
  FaArrowUp, 
  FaArrowDown 
} from 'react-icons/fa';
import './MessagesAdmin.css';

const MessagesAdmin = () => {
  const [messages, setMessages] = useState([
    { 
      id: 1, 
      name: 'Juan Pérez', 
      email: 'juan@example.com', 
      message: '¿Tienen cursos sobre TypeScript? Me interesa conocer si ofrecen algún curso especializado en TypeScript para desarrollo web frontend. Ya tengo algo de experiencia con JavaScript pero quiero aprender a usar tipos.',
      read: false,
      starred: false,
      date: '2023-09-15',
      time: '14:30'
    },
    { 
      id: 2, 
      name: 'María García', 
      email: 'maria@example.com', 
      message: 'Necesito ayuda con mi suscripción. Hace dos días realicé el pago para el plan anual pero aún no tengo acceso a todos los cursos que deberían estar incluidos en mi suscripción.',
      read: true,
      starred: true,
      date: '2023-09-14',
      time: '10:15'
    },
    { 
      id: 3, 
      name: 'Carlos Rodríguez', 
      email: 'carlos@example.com', 
      message: 'Excelente plataforma, me han encantado los cursos. Quería felicitarles por la calidad del contenido y la metodología de enseñanza. Es muy clara y práctica.',
      read: true,
      starred: false,
      date: '2023-09-12',
      time: '16:45'
    },
    { 
      id: 4, 
      name: 'Ana López', 
      email: 'ana@example.com', 
      message: 'Me gustaría conocer si tienen algún descuento para estudiantes universitarios. Estoy en el último año de Ingeniería Informática y me interesa complementar mi formación con sus cursos.',
      read: false,
      starred: false,
      date: '2023-09-10',
      time: '09:22'
    },
    { 
      id: 5, 
      name: 'Roberto Sánchez', 
      email: 'roberto@example.com', 
      message: 'Tengo problemas para ver los videos del curso de Node.js. He intentado desde diferentes navegadores pero siempre se queda cargando indefinidamente.',
      read: false,
      starred: true,
      date: '2023-09-08',
      time: '13:10'
    }
  ]);

  const [filteredMessages, setFilteredMessages] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [filter, setFilter] = useState('all');
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedMessage, setSelectedMessage] = useState(null);
  const [sortConfig, setSortConfig] = useState({
    key: 'date',
    direction: 'desc'
  });
  const [replyText, setReplyText] = useState('');

  // Memoiza la función applyFilters con useCallback
  const applyFilters = useCallback(() => {
    let result = [...messages];
    
    // Aplicar filtro por estado
    if (filter === 'unread') {
      result = result.filter(message => !message.read);
    } else if (filter === 'read') {
      result = result.filter(message => message.read);
    } else if (filter === 'starred') {
      result = result.filter(message => message.starred);
    }
    
    // Aplicar búsqueda
    if (searchTerm) {
      const searchLower = searchTerm.toLowerCase();
      result = result.filter(
        message =>
          message.name.toLowerCase().includes(searchLower) ||
          message.email.toLowerCase().includes(searchLower) ||
          message.message.toLowerCase().includes(searchLower)
      );
    }
    
    // Aplicar ordenación
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
  }, [filter, messages, searchTerm, sortConfig]);

  // Efecto para inicializar los mensajes filtrados
  useEffect(() => {
    const timer = setTimeout(() => {
      applyFilters();
      setIsLoading(false);
    }, 800);
    
    return () => clearTimeout(timer);
  }, [applyFilters]);

  // Efecto para aplicar filtros cuando cambian
  useEffect(() => {
    applyFilters();
  }, [applyFilters]);

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

  // Marcar como leído/no leído
  const toggleReadStatus = (id) => {
    setMessages(
      messages.map(msg =>
        msg.id === id ? { ...msg, read: !msg.read } : msg
      )
    );
  };

  // Marcar/desmarcar favorito
  const toggleStarred = (id) => {
    setMessages(
      messages.map(msg =>
        msg.id === id ? { ...msg, starred: !msg.starred } : msg
      )
    );
  };

  // Eliminar mensaje
  const deleteMessage = (id) => {
    if (window.confirm('¿Estás seguro de que quieres eliminar este mensaje?')) {
      setMessages(messages.filter(msg => msg.id !== id));
      if (selectedMessage && selectedMessage.id === id) {
        setSelectedMessage(null);
      }
    }
  };

  // Ver detalle de mensaje
  const viewMessage = (message) => {
    // Si el mensaje no está leído, marcarlo como leído
    if (!message.read) {
      toggleReadStatus(message.id);
    }
    setSelectedMessage(message);
  };

  // Enviar respuesta
  const sendReply = (e) => {
    e.preventDefault();
    
    if (!replyText.trim()) {
      alert('Por favor, escribe un mensaje antes de enviar.');
      return;
    }
    
    // Aquí se implementaría el envío real de la respuesta
    alert(`Respuesta enviada a ${selectedMessage.email}: ${replyText}`);
    setReplyText('');
    setSelectedMessage(null);
  };

  // Loading overlay
  const LoadingOverlay = () => (
    <div className="loading-overlay">
      <div className="spinner"></div>
    </div>
  );

  // Formatear fecha
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
              />
              <div className="form-actions">
                <button type="button" className="btn-cancel" onClick={() => setSelectedMessage(null)}>
                  Cancelar
                </button>
                <button type="submit" className="btn-primary">
                  Enviar Respuesta
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