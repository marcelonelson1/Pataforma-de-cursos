import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import './Cursos.css';
import CursosService from './services/CursosService';

function Cursos() {
  const [cursos, setCursos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Cursos estáticos para fallback (todos marcados como publicados)
  const cursosDefault = [
    { id: 1, titulo: 'Curso de React', descripcion: 'Aprende React desde cero.', imagen_url: 'https://via.placeholder.com/300', duracion: '10 horas', nivel: 'Principiante', estado: 'Publicado' },
    { id: 2, titulo: 'Curso de Node.js', descripcion: 'Domina el backend con Node.js.', imagen_url: 'https://via.placeholder.com/300', duracion: '15 horas', nivel: 'Intermedio', estado: 'Publicado' },
    { id: 3, titulo: 'Curso de Diseño Web', descripcion: 'Crea diseños modernos y responsivos.', imagen_url: 'https://via.placeholder.com/300', duracion: '12 horas', nivel: 'Principiante', estado: 'Publicado' },
  ];

  useEffect(() => {
    const fetchCursos = async () => {
      try {
        setLoading(true);
        const response = await CursosService.getCursos();
        console.log('Respuesta del servidor:', response);
        
        let cursosData = [];
        
        // Verificar si la respuesta es un array
        if (Array.isArray(response)) {
          cursosData = response;
        } else if (response && Array.isArray(response.data)) {
          cursosData = response.data;
        } else {
          // Si la respuesta no es un array, usar los cursos por defecto
          console.warn('La respuesta del servidor no es un array:', response);
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
                src={curso.imagen_url || 'https://via.placeholder.com/300'} 
                alt={curso.titulo} 
                className="curso-image" 
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
              </div>
              <button className="curso-cta">Más Información</button>
            </div>
          </li>
        ))}
      </ul>
      {error && <div className="cursos-error">{error}</div>}
    </div>
  );
}

export default Cursos;