import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faSearch, faPlus, faEdit, faTrashAlt, faArchive,
  faEye, faSort, faShoppingCart, faVideo, faCalendarAlt
} from '@fortawesome/free-solid-svg-icons';
import './AdminPages.css';

function AdminCursosPage() {
  const [cursos, setCursos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [ordenarPor, setOrdenarPor] = useState('fecha');
  const [mostrarArchivados, setMostrarArchivados] = useState(false);

  useEffect(() => {
    const fetchCursos = async () => {
      try {
        const token = localStorage.getItem('token');
        const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
        
        const url = new URL(`${apiUrl}/api/admin/cursos`);
        if (searchTerm) url.searchParams.append('titulo', searchTerm);
        url.searchParams.append('ordenar', ordenarPor);
        if (mostrarArchivados) url.searchParams.append('archivados', 'true');
        
        const response = await fetch(url, {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });

        if (!response.ok) {
          throw new Error('Error al obtener los cursos');
        }

        const data = await response.json();
        setCursos(data);
        setError(null);
      } catch (err) {
        console.error('Error al cargar cursos:', err);
        setError('No se pudieron cargar los cursos. Por favor, intenta de nuevo más tarde.');
      } finally {
        setLoading(false);
      }
    };

    fetchCursos();
  }, [searchTerm, ordenarPor, mostrarArchivados]);

  // Datos de ejemplo mientras se implementa la API
  const demoData = [
    {
      id: 1,
      titulo: 'Curso de React',
      descripcion: 'Aprende React desde cero hasta un nivel avanzado con proyectos prácticos.',
      precio: 29.99,
      total_ventas: 45,
      total_ingresos: 1349.55,
      cant_capitulos: 12,
      fecha_creacion: '2023-08-25T14:30:00Z',
      ultima_edicion: '2023-09-10T09:15:00Z',
      editado_por: 'Admin',
      archivado: false
    },
    {
      id: 2,
      titulo: 'Curso de Node.js',
      descripcion: 'Domina el backend con Node.js. Aprende a crear APIs RESTful, autenticación y más.',
      precio: 39.99,
      total_ventas: 38,
      total_ingresos: 1519.62,
      cant_capitulos: 10,
      fecha_creacion: '2023-07-15T10:00:00Z',
      ultima_edicion: '2023-09-05T15:30:00Z',
      editado_por: 'Admin',
      archivado: false
    },
    {
      id: 3,
      titulo: 'Curso de Diseño Web',
      descripcion: 'Aprende HTML, CSS y JavaScript para crear sitios web atractivos y responsivos.',
      precio: 24.99,
      total_ventas: 27,
      total_ingresos: 674.73,
      cant_capitulos: 8,
      fecha_creacion: '2023-09-01T08:45:00Z',
      ultima_edicion: '2023-09-12T11:20:00Z',
      editado_por: 'Admin',
      archivado: true
    }
  ];

  // Usar datos de la API o datos de ejemplo si aún está cargando
  const data = cursos.length > 0 ? cursos : demoData;
  
  // Filtrar cursos si estamos usando datos de demostración
  const filteredData = data
    .filter(curso => 
      (mostrarArchivados || !curso.archivado) && 
      curso.titulo.toLowerCase().includes(searchTerm.toLowerCase())
    );

  // Formatear fecha
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-ES', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric'
    });
  };

  const handleSearch = (e) => {
    e.preventDefault();
    // Si estamos usando datos de demostración, la búsqueda ya se aplica con el filtro
    // Si estamos usando la API, el cambio en searchTerm desencadenará una nueva solicitud
  };

  const archivarCurso = async (id) => {
    if (window.confirm('¿Estás seguro de que deseas archivar este curso?')) {
      try {
        const token = localStorage.getItem('token');
        const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
        
        const response = await fetch(`${apiUrl}/api/admin/cursos/${id}/archivar`, {
          method: 'PATCH',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });

        if (!response.ok) {
          throw new Error('Error al archivar el curso');
        }

        // Actualizar la lista de cursos
        setCursos(cursos.map(curso => 
          curso.id === id ? { ...curso, archivado: true } : curso
        ));
        
        alert('Curso archivado correctamente');
      } catch (err) {
        console.error('Error al archivar curso:', err);
        alert('Error al archivar el curso. Por favor, intenta de nuevo.');
      }
    }
  };

  const eliminarCurso = async (id) => {
    if (window.confirm('¿Estás COMPLETAMENTE seguro de eliminar permanentemente este curso? Esta acción no se puede deshacer.')) {
      if (window.confirm('Confirmar nuevamente: ¿Eliminar permanentemente este curso?')) {
        try {
          const token = localStorage.getItem('token');
          const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
          
          const response = await fetch(`${apiUrl}/api/admin/cursos/${id}?confirmacion=confirmar_eliminacion`, {
            method: 'DELETE',
            headers: {
              'Authorization': `Bearer ${token}`,
              'Content-Type': 'application/json'
            }
          });

          if (!response.ok) {
            throw new Error('Error al eliminar el curso');
          }

          // Actualizar la lista de cursos
          setCursos(cursos.filter(curso => curso.id !== id));
          
          alert('Curso eliminado permanentemente');
        } catch (err) {
          console.error('Error al eliminar curso:', err);
          alert('Error al eliminar el curso. Por favor, intenta de nuevo.');
        }
      }
    }
  };

  if (loading && cursos.length === 0) {
    return (
      <div className="admin-loading-container">
        <div className="spinner"></div>
        <p>Cargando cursos...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="admin-error-container">
        <h2>Error</h2>
        <p>{error}</p>
        <button className="admin-btn admin-btn-primary" onClick={() => window.location.reload()}>
          Intentar de nuevo
        </button>
      </div>
    );
  }

  return (
    <div className="admin-cursos-page">
      <div className="admin-page-title">
        <h1>Gestión de Cursos</h1>
        <p>Administra todos los cursos de la plataforma</p>
      </div>

      <div className="admin-actions-bar">
        <form className="admin-search-bar" onSubmit={handleSearch}>
          <input 
            type="text" 
            className="admin-search-input" 
            placeholder="Buscar cursos..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
          <button type="submit" className="admin-search-button">
            <FontAwesomeIcon icon={faSearch} />
          </button>
        </form>

        <div className="admin-filters">
          <select 
            className="admin-filter-select" 
            value={ordenarPor}
            onChange={(e) => setOrdenarPor(e.target.value)}
          >
            <option value="fecha">Ordenar por fecha</option>
            <option value="ventas">Ordenar por ventas</option>
            <option value="ingresos">Ordenar por ingresos</option>
            <option value="titulo">Ordenar por título</option>
          </select>
          
          <div className="admin-filter-checkbox">
            <input 
              type="checkbox" 
              id="mostrarArchivados" 
              checked={mostrarArchivados}
              onChange={(e) => setMostrarArchivados(e.target.checked)}
            />
            <label htmlFor="mostrarArchivados">Mostrar archivados</label>
          </div>
        </div>

        <div className="admin-actions">
          <Link to="/admin/cursos/nuevo" className="admin-btn admin-btn-primary">
            <FontAwesomeIcon icon={faPlus} /> Nuevo Curso
          </Link>
        </div>
      </div>

      <div className="admin-cursos-list">
        {filteredData.length === 0 ? (
          <div className="admin-empty-state">
            <p>No se encontraron cursos</p>
            <Link to="/admin/cursos/nuevo" className="admin-btn admin-btn-primary">
              <FontAwesomeIcon icon={faPlus} /> Crear Nuevo Curso
            </Link>
          </div>
        ) : (
          filteredData.map(curso => (
            <div 
              key={curso.id} 
              className={`admin-curso-card ${curso.archivado ? 'curso-archivado' : ''}`}
            >
              <div className="admin-curso-image">
                <img src={`/curso-thumbnail-${curso.id}.jpg`} alt={curso.titulo} onError={(e) => { e.target.src = '/curso-default.jpg'; }} />
              </div>
              <div className="admin-curso-content">
                <div className="admin-curso-header">
                  <h3 className="admin-curso-title">
                    {curso.titulo}
                    {curso.archivado && <span className="admin-badge admin-badge-warning ml-2">Archivado</span>}
                  </h3>
                  <div className="admin-curso-price">${curso.precio.toFixed(2)}</div>
                </div>
                <p className="admin-curso-description">{curso.descripcion}</p>
                <div className="admin-curso-footer">
                  <div className="admin-curso-stats">
                    <div className="admin-curso-stat">
                      <FontAwesomeIcon icon={faShoppingCart} /> {curso.total_ventas} ventas
                    </div>
                    <div className="admin-curso-stat">
                      <FontAwesomeIcon icon={faVideo} /> {curso.cant_capitulos} capítulos
                    </div>
                    <div className="admin-curso-stat">
                      <FontAwesomeIcon icon={faCalendarAlt} /> Actualizado: {formatDate(curso.ultima_edicion)}
                    </div>
                  </div>
                  <div className="admin-curso-actions">
                    <Link to={`/curso/${curso.id}`} className="admin-btn admin-btn-sm admin-btn-secondary" title="Ver curso">
                      <FontAwesomeIcon icon={faEye} />
                    </Link>
                    <Link to={`/admin/cursos/${curso.id}/capitulos`} className="admin-btn admin-btn-sm admin-btn-secondary" title="Gestionar capítulos">
                      <FontAwesomeIcon icon={faSort} />
                    </Link>
                    <Link to={`/admin/cursos/editar/${curso.id}`} className="admin-btn admin-btn-sm admin-btn-secondary" title="Editar curso">
                      <FontAwesomeIcon icon={faEdit} />
                    </Link>
                    {!curso.archivado ? (
                      <button 
                        className="admin-btn admin-btn-sm admin-btn-secondary" 
                        onClick={() => archivarCurso(curso.id)}
                        title="Archivar curso"
                      >
                        <FontAwesomeIcon icon={faArchive} />
                      </button>
                    ) : (
                      <button 
                        className="admin-btn admin-btn-sm admin-btn-danger" 
                        onClick={() => eliminarCurso(curso.id)}
                        title="Eliminar permanentemente"
                      >
                        <FontAwesomeIcon icon={faTrashAlt} />
                      </button>
                    )}
                  </div>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

export default AdminCursosPage;