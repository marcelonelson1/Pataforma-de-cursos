import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import './Cursos.css';
import CursosService from './services/CursosService';
import axios from 'axios';

function Cursos() {
  const [cursos, setCursos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [cursosPagados, setCursosPagados] = useState({});

  // Cursos estáticos para fallback (todos marcados como publicados)
  const cursosDefault = [
    { id: 1, titulo: 'Curso de React', descripcion: 'Aprende React desde cero.', imagen_url: 'https://via.placeholder.com/300', duracion: '10 horas', nivel: 'Principiante', estado: 'Publicado', precio: 29.99 },
    { id: 2, titulo: 'Curso de Node.js', descripcion: 'Domina el backend con Node.js.', imagen_url: 'https://via.placeholder.com/300', duracion: '15 horas', nivel: 'Intermedio', estado: 'Publicado', precio: 39.99 },
    { id: 3, titulo: 'Curso de Diseño Web', descripcion: 'Crea diseños modernos y responsivos.', imagen_url: 'https://via.placeholder.com/300', duracion: '12 horas', nivel: 'Principiante', estado: 'Publicado', precio: 24.99 },
  ];

  // Función para obtener la URL base de la API
  const getApiBaseUrl = () => {
    // Usar la variable de entorno si está definida
    if (process.env.REACT_APP_API_URL) {
      return process.env.REACT_APP_API_URL;
    }
    
    // Si no está definida, intentar determinarla basándose en la URL actual
    const currentUrl = window.location.origin; // e.g. http://localhost:3000
    
    // Si estamos en desarrollo local, asumimos que la API está en el puerto 5000
    if (currentUrl.includes('localhost')) {
      return 'http://localhost:5000';
    }
    
    // En producción, asumimos que la API está en el mismo origen
    return currentUrl;
  };

  // Función mejorada para verificar y corregir URLs de imágenes
  const getImageUrl = (url) => {
    if (!url) return 'https://via.placeholder.com/300';
    
    // Si ya es una URL completa o una data URL, devolverla como está
    if (url.startsWith('http') || url.startsWith('data:')) {
      return url;
    }
    
    // Si es una ruta relativa que empieza con /static, añadir la URL base
    if (url.startsWith('/static/')) {
      const baseUrl = getApiBaseUrl();
      // Asegurar que no hay barras duplicadas entre baseUrl y url
      return `${baseUrl}${url}`;
    }
    
    // Por defecto, devolver la URL sin cambios
    return url;
  };
  
  // Función para verificar si el usuario ha pagado los cursos
  const verificarPagosCursos = async (cursosData) => {
    try {
      // Verificamos si hay un token (usuario autenticado)
      const token = localStorage.getItem('token');
      
      if (!token) {
        console.log('Usuario no autenticado, no se verifican pagos');
        return {};
      }
      
      const pagosPromises = cursosData.map(async (curso) => {
        try {
          const response = await axios.get(
            `${getApiBaseUrl()}/api/pagos/${curso.id}`,
            {
              headers: {
                Authorization: `Bearer ${token}`
              }
            }
          );
          
          // Determinamos si el curso está pagado basado en la respuesta
          return {
            cursoId: curso.id,
            pagado: response.data.estado === 'aprobado'
          };
        } catch (error) {
          console.log(`Error al verificar pago del curso ${curso.id}:`, error);
          return { cursoId: curso.id, pagado: false };
        }
      });

      const resultadosPagos = await Promise.all(pagosPromises);
      
      // Convertimos el array de resultados a un objeto para fácil acceso
      const pagosObj = resultadosPagos.reduce((acc, item) => {
        acc[item.cursoId] = item.pagado;
        return acc;
      }, {});
      
      console.log('Estado de pagos de cursos:', pagosObj);
      
      return pagosObj;
    } catch (error) {
      console.error('Error al verificar pagos:', error);
      return {};
    }
  };

  useEffect(() => {
    const fetchCursos = async () => {
      try {
        setLoading(true);
        
        // Usar axios directamente para tener más control sobre la respuesta
        const response = await axios.get(`${getApiBaseUrl()}/api/cursos`);
        console.log('Respuesta del servidor:', response);
        
        let cursosData = [];
        
        // Verificar si la respuesta contiene datos
        if (response && response.data) {
          cursosData = Array.isArray(response.data) ? response.data : [];
        } else {
          console.warn('La respuesta del servidor no contiene datos válidos');
          cursosData = cursosDefault;
        }
        
        // Filtrar solo los cursos publicados
        const cursosPublicados = cursosData.filter(curso => 
          curso.estado === 'Publicado' || curso.estado === 'publicado'
        );
        
        console.log('Cursos publicados:', cursosPublicados);
        
        // Si no hay cursos publicados, mostrar mensaje informativo
        if (cursosPublicados.length === 0) {
          setError('No hay cursos publicados disponibles en este momento.');
        }
        
        // Añadir un log para depurar las URLs de las imágenes
        cursosPublicados.forEach(curso => {
          console.log(`Curso ${curso.id} - URL de imagen original:`, curso.imagen_url);
          console.log(`Curso ${curso.id} - URL de imagen procesada:`, getImageUrl(curso.imagen_url));
        });
        
        // Verificar estado de pagos para los cursos
        const pagos = await verificarPagosCursos(cursosPublicados);
        setCursosPagados(pagos);
        
        setCursos(cursosPublicados);
        setError(null);
      } catch (err) {
        console.error('Error al cargar cursos:', err);
        setError('No se pudieron cargar los cursos. Por favor, intenta de nuevo más tarde.');
        // En caso de error, usar los cursos publicados por defecto
        setCursos(cursosDefault);
      } finally {
        setLoading(false);
      }
    };

    fetchCursos();
    
    // Volver a verificar cuando el componente se monte o cuando cambie el token
    const interval = setInterval(() => {
      const token = localStorage.getItem('token');
      if (token && cursos.length > 0) {
        verificarPagosCursos(cursos).then(pagos => {
          setCursosPagados(pagos);
        });
      }
    }, 60000); // Verificar cada minuto
    
    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return <div className="cursos-loading">Cargando cursos...</div>;
  }

  if (error && cursos.length === 0) {
    return <div className="cursos-error">{error}</div>;
  }

  if (cursos.length === 0) {
    return <div className="cursos-empty">No hay cursos publicados disponibles en este momento.</div>;
  }

  return (
    <div className="cursos-container">
      <ul className="cursos-list">
        {cursos.map((curso) => (
          <li key={curso.id} className="curso-item animate__animated animate__fadeInUp">
            <div className="curso-image-container">
              <img 
                src={getImageUrl(curso.imagen_url)} 
                alt={curso.titulo} 
                className="curso-image" 
                onError={(e) => {
                  console.log(`Error al cargar imagen para curso ${curso.id}:`, curso.imagen_url);
                  e.target.onerror = null; // Prevenir bucle infinito
                  e.target.src = 'https://via.placeholder.com/300';
                }}
              />
              <span className="curso-badge">{curso.nivel || 'Todos los niveles'}</span>
            </div>
            <Link to={`/curso/${curso.id}`} className="curso-link">
              {curso.titulo}
            </Link>
            <p className="curso-descripcion">{curso.descripcion}</p>
            <div className="curso-footer">
              <div className="curso-meta">
                <span className="curso-meta-item">
                  <i className="curso-meta-icon">⏳</i> {curso.duracion || (curso.capitulos ? `${curso.capitulos.length} capítulos` : 'Curso completo')}
                </span>
                
                {/* Mostrar precio solo si el curso NO está pagado */}
                {!cursosPagados[curso.id] && (
                  <span className="curso-meta-item curso-precio">
                    <i className="curso-meta-icon">💰</i> ${curso.precio?.toFixed(2) || '29.99'}
                  </span>
                )}
              </div>
              <Link to={`/curso/${curso.id}`} className="curso-cta">
                {cursosPagados[curso.id] ? 'Acceder al Curso' : 'Más Información'}
              </Link>
            </div>
          </li>
        ))}
      </ul>
      {error && <div className="cursos-error">{error}</div>}
    </div>
  );
}

export default Cursos;