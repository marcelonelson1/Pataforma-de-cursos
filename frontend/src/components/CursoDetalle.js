import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import VideoPlayer from './VideoPlayer';
import './CursoDetalle.css';

function CursoDetalle() {
  const { id } = useParams();
  const [curso, setCurso] = useState(null);
  const [loading, setLoading] = useState(true);

  // Datos de prueba para los cursos
  const cursos = [
    {
      id: 1,
      titulo: 'Curso de React',
      descripcion: 'Aprende React desde cero hasta un nivel avanzado con proyectos prácticos.',
      contenido: 'En este curso aprenderás todos los conceptos fundamentales de React, componentes, estados, props, hooks, y mucho más.',
      videos: [
        {
          id: 1,
          titulo: 'Introducción a React',
          url: 'https://www.learningcontainer.com/wp-content/uploads/2020/05/sample-mp4-file.mp4',
          duracion: '8:24'
        },
        {
          id: 2,
          titulo: 'Componentes y Props',
          url: 'https://filesamples.com/samples/video/mp4/sample_640x360.mp4',
          duracion: '10:15'
        },
        {
          id: 3,
          titulo: 'Estado y Ciclo de Vida',
          url: 'https://filesamples.com/samples/video/mp4/sample_960x540.mp4',
          duracion: '12:50'
        },
        {
          id: 4,
          titulo: 'Hooks en React',
          url: 'https://filesamples.com/samples/video/mp4/sample_1280x720.mp4',
          duracion: '15:30'
        }
      ]
    },
    {
      id: 2,
      titulo: 'Curso de Node.js',
      descripcion: 'Domina el backend con Node.js y construye APIs robustas.',
      contenido: 'Aprenderás a utilizar Node.js para crear aplicaciones de servidor, APIs RESTful y mucho más.',
      videos: [
        {
          id: 1,
          titulo: 'Introducción a Node.js',
          url: 'https://filesamples.com/samples/video/mp4/sample_1920x1080.mp4',
          duracion: '7:15'
        },
        {
          id: 2,
          titulo: 'Módulos y NPM',
          url: 'https://www.learningcontainer.com/wp-content/uploads/2020/05/sample-mp4-file.mp4',
          duracion: '9:30'
        },
        {
          id: 3,
          titulo: 'Express.js',
          url: 'https://filesamples.com/samples/video/mp4/sample_640x360.mp4',
          duracion: '11:45'
        }
      ]
    },
    {
      id: 3,
      titulo: 'Curso de Diseño Web',
      descripcion: 'Crea diseños modernos y responsivos utilizando HTML5 y CSS3.',
      contenido: 'Aprenderás los fundamentos del diseño web moderno, incluyendo flex, grid y animaciones CSS.',
      videos: [
        {
          id: 1,
          titulo: 'Fundamentos de HTML5',
          url: 'https://filesamples.com/samples/video/mp4/sample_960x540.mp4',
          duracion: '6:45'
        },
        {
          id: 2,
          titulo: 'CSS Moderno',
          url: 'https://filesamples.com/samples/video/mp4/sample_1280x720.mp4',
          duracion: '10:20'
        },
        {
          id: 3,
          titulo: 'Diseño Responsivo',
          url: 'https://www.learningcontainer.com/wp-content/uploads/2020/05/sample-mp4-file.mp4',
          duracion: '8:55'
        },
        {
          id: 4,
          titulo: 'Animaciones CSS',
          url: 'https://filesamples.com/samples/video/mp4/sample_1920x1080.mp4',
          duracion: '7:40'
        },
        {
          id: 5,
          titulo: 'Flexbox y Grid',
          url: 'https://filesamples.com/samples/video/mp4/sample_640x360.mp4',
          duracion: '12:10'
        }
      ]
    }
  ];

  useEffect(() => {
    const cursoEncontrado = cursos.find((c) => c.id === parseInt(id));
    setCurso(cursoEncontrado);
    setLoading(false);
  }, [id]);

  if (loading) {
    return <div className="loading">Cargando curso...</div>;
  }

  if (!curso) {
    return <div className="error">Curso no encontrado</div>;
  }

  return (
    <div className="curso-detalle animate__animated animate__fadeIn">
      <h2>{curso.titulo}</h2>
      <p>{curso.descripcion}</p>

      {/* Reproductor de video */}
      <VideoPlayer videos={curso.videos} />

      {/* Lista de videos/lecciones */}
      <div className="lessons-list">
        <h3>Contenido del curso</h3>
        {curso.videos.map((video, index) => (
          <div
            key={video.id}
            className={`lesson-item ${index === 0 ? 'active' : ''}`}
          >
            <span className="lesson-title">{video.titulo}</span>
            <span className="lesson-duration">{video.duracion}</span>
          </div>
        ))}
      </div>

      <div className="contenido">
        <h3>Descripción del curso</h3>
        <p>{curso.contenido}</p>
      </div>
    </div>
  );
}

export default CursoDetalle;