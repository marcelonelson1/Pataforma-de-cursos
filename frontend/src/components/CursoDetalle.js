import React, { useState, useEffect, useMemo } from 'react';
import { useParams } from 'react-router-dom';
import VideoPlayer from './VideoPlayer';
import './CursoDetalle.css';

function CursoDetalle() {
  const { id } = useParams();
  const [curso, setCurso] = useState(null);
  const [loading, setLoading] = useState(true);
  const [currentVideo, setCurrentVideo] = useState(0);
  const [completedVideos, setCompletedVideos] = useState({});
  const [overallProgress, setOverallProgress] = useState(0);

  // Memorizar el array 'cursos' para evitar recrearlo en cada renderizado
  const cursos = useMemo(() => [
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
    // Más cursos...
  ], []); // <-- Array de dependencias vacío para que solo se calcule una vez

  // Buscar el curso correspondiente al ID
  useEffect(() => {
    const cursoEncontrado = cursos.find((c) => c.id === parseInt(id));
    setCurso(cursoEncontrado);
    setLoading(false);
  }, [id, cursos]); // <-- 'cursos' ya no cambia en cada renderizado

  // Calcular el progreso general del curso
  useEffect(() => {
    if (curso && Object.keys(completedVideos).length > 0) {
      const completedCount = Object.values(completedVideos).filter(completed => completed).length;
      const totalVideos = curso.videos.length;
      const newOverallProgress = Math.round((completedCount / totalVideos) * 100);
      setOverallProgress(newOverallProgress);
    }
  }, [completedVideos, curso]);

  // Cambiar de video
  const changeVideo = (index) => {
    if (index >= 0 && index < curso.videos.length) {
      setCurrentVideo(index);
    }
  };

  // Marcar video como completado
  const markAsCompleted = (videoId) => {
    setCompletedVideos(prev => ({
      ...prev,
      [videoId]: !prev[videoId]
    }));
  };

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

      {/* Barra de progreso del curso */}
      <div className="course-progress">
        <div className="progress-header">
          <span className="progress-title">Progreso del curso</span>
          <span className="progress-percentage">{overallProgress}%</span>
        </div>
        <div className="progress-bar-container">
          <div
            className="progress-bar-fill"
            style={{ width: `${overallProgress}%` }}
          ></div>
        </div>
      </div>

      {/* Reproductor de video con contenido del curso integrado */}
      <VideoPlayer
        videos={curso.videos}
        currentVideo={currentVideo}
        changeVideo={changeVideo}
        completedVideos={completedVideos}
        markAsCompleted={markAsCompleted}
      />

      <div className="contenido">
        <h3>Descripción del curso</h3>
        <p>{curso.contenido}</p>
      </div>
    </div>
  );
}

export default CursoDetalle;