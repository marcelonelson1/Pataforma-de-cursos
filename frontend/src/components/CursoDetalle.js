import React, { useState, useEffect, useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import VideoPlayer from './VideoPlayer';
import './CursoDetalle.css';
import { useAuth } from '../context/AuthContext';
import PaymentModal from './PaymentModal';
import CursosService from './services/CursosService';

function CursoDetalle() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isLoggedIn, user } = useAuth();
  const [curso, setCurso] = useState(null);
  const [loading, setLoading] = useState(true);
  const [currentVideo, setCurrentVideo] = useState(0);
  const [completedVideos, setCompletedVideos] = useState({});
  const [overallProgress, setOverallProgress] = useState(0);
  const [paymentStatus, setPaymentStatus] = useState(null);
  const [showPaymentModal, setShowPaymentModal] = useState(false);
  const [error, setError] = useState(null);

  // Cursos estáticos para fallback
  const cursosFallback = useMemo(() => [
    {
      id: 1,
      titulo: 'Curso de React',
      descripcion: 'Aprende React desde cero hasta un nivel avanzado con proyectos prácticos.',
      contenido: 'En este curso aprenderás todos los conceptos fundamentales de React, componentes, estados, props, hooks, y mucho más.',
      precio: 29.99,
      estado: 'Publicado',
      videos: [
        {
          id: 1,
          indice: "1",
          titulo: 'Introducción a React',
          url: 'https://www.learningcontainer.com/wp-content/uploads/2020/05/sample-mp4-file.mp4',
          duracion: '08:45'
        },
        {
          id: 2,
          indice: "2",
          titulo: 'Componentes y Props',
          url: 'https://filesamples.com/samples/video/mp4/sample_640x360.mp4',
          duracion: '10:15'
        },
        {
          id: 3,
          indice: "3",
          titulo: 'Estado y Ciclo de Vida',
          url: 'https://filesamples.com/samples/video/mp4/sample_960x540.mp4',
          duracion: '12:50'
        },
        {
          id: 4,
          indice: "4",
          titulo: 'Hooks en React',
          url: 'https://filesamples.com/samples/video/mp4/sample_1280x720.mp4',
          duracion: '15:30'
        }
      ]
    },
  ], []);

  // Carga inicial del curso desde el backend
  useEffect(() => {
    const fetchCurso = async () => {
      try {
        setLoading(true);
        console.log(`Obteniendo curso con ID: ${id}`);
        const response = await CursosService.getCursoById(id);
        console.log('Respuesta del servidor:', response);
        
        // Procesar la respuesta
        let cursoData = response;
        
        // Si la respuesta tiene una propiedad data, usar esa
        if (response && response.data) {
          cursoData = response.data;
        }
        
        // Verificar si el curso está publicado
        const esPublicado = cursoData.estado === 'Publicado' || cursoData.estado === 'publicado';
        const esAdmin = user && user.role === 'admin'; // Suponiendo que tienes un rol de usuario en tu autenticación
        
        // Si el curso no está publicado y el usuario no es admin, redirigir o mostrar error
        if (!esPublicado && !esAdmin) {
          setError('Este curso no está disponible actualmente.');
          setCurso(null);
          setLoading(false);
          return;
        }
        
        // Si el curso tiene capítulos, adaptarlos al formato esperado por VideoPlayer
        if (cursoData && cursoData.capitulos) {
          // Transformar capitulos a formato compatible con videos
          const videosFormateados = cursoData.capitulos.map(capitulo => ({
            id: capitulo.id,
            indice: capitulo.orden?.toString() || "",
            titulo: capitulo.titulo,
            url: capitulo.video_url,
            duracion: capitulo.duracion || "00:00"
          }));
          
          // Mantener ambos formatos para compatibilidad
          cursoData.videos = videosFormateados;
        }
        
        setCurso(cursoData);
        setError(null);
      } catch (err) {
        console.error('Error al cargar curso:', err);
        setError('No se pudo cargar el curso. Por favor, intenta de nuevo más tarde.');
        
        // Si hay error, usar el curso estático como fallback
        const cursoEncontrado = cursosFallback.find((c) => c.id === parseInt(id));
        if (cursoEncontrado) {
          setCurso(cursoEncontrado);
        }
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchCurso();
    }
  }, [id, cursosFallback, user]);

  // Verificar estado de pago
  useEffect(() => {
    if (!isLoggedIn || !curso) {
      setPaymentStatus('no_pagado');
      return;
    }

    const checkPaymentStatus = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          setPaymentStatus('no_pagado');
          return;
        }

        console.log(`Verificando estado de pago para curso ID: ${curso.id}`);
        
        // Usamos el servicio para verificar el pago
        const response = await CursosService.verificarPago(curso.id);
        console.log('Respuesta de verificación de pago:', response);
        
        // Extraer el estado del pago
        let estadoPago = 'no_pagado';
        
        if (response && response.estado) {
          estadoPago = response.estado;
        } else if (response && response.data && response.data.estado) {
          estadoPago = response.data.estado;
        }
        
        setPaymentStatus(estadoPago);
        setError(null);
        
      } catch (err) {
        console.error('Error al verificar estado de pago:', err);
        
        let errorMessage = 'Error al verificar el estado del pago.';
        
        if (err.message) {
          if (err.message.includes('Failed to fetch') || err.message.includes('Network Error')) {
            errorMessage = 'Error de conexión con el servidor. Asegúrate de que el backend esté en ejecución.';
          } else if (err.message.includes('404')) {
            errorMessage = 'No se encontró registro de pago';
          } else if (err.message.includes('401')) {
            errorMessage = 'No autorizado - por favor inicia sesión nuevamente';
          } else if (err.message.includes('JSON')) {
            errorMessage = 'Formato de respuesta del servidor incorrecto';
          }
        }
        
        setError(errorMessage);
        setPaymentStatus('no_pagado');
      }
    };

    checkPaymentStatus();
  }, [isLoggedIn, curso]);

  // Calcular progreso del curso
  useEffect(() => {
    if (curso) {
      const videos = curso.videos || curso.capitulos || [];
      if (videos.length > 0 && Object.keys(completedVideos).length > 0) {
        const completedCount = Object.values(completedVideos).filter(completed => completed).length;
        const totalVideos = videos.length;
        const newOverallProgress = Math.round((completedCount / totalVideos) * 100);
        setOverallProgress(newOverallProgress);
      }
    }
  }, [completedVideos, curso]);

  const changeVideo = (index) => {
    const videos = curso?.videos || curso?.capitulos || [];
    if (index >= 0 && index < videos.length) {
      setCurrentVideo(index);
    }
  };

  const markAsCompleted = (videoId) => {
    // Actualizar estado local
    const newCompletedState = !completedVideos[videoId];
    setCompletedVideos(prev => ({
      ...prev,
      [videoId]: newCompletedState
    }));
    
    // Opcionalmente, podríamos sincronizar con el backend
    if (curso) {
      CursosService.marcarCapituloCompletado(curso.id, videoId, newCompletedState)
        .catch(err => console.error('Error al sincronizar estado completado:', err));
    }
  };

  const handlePaymentSuccess = async () => {
    try {
      setPaymentStatus('aprobado');
      setShowPaymentModal(false);
      setError(null);
      
      // Opcionalmente, recargar datos del curso
      const updatedCurso = await CursosService.getCursoById(id);
      setCurso(updatedCurso);
    } catch (err) {
      console.error('Error al actualizar después del pago:', err);
    }
  };

  if (loading) {
    return <div className="loading">Cargando curso...</div>;
  }

  if (error && !curso) {
    return (
      <div className="error-container">
        <div className="error">{error}</div>
        <button className="btn-back" onClick={() => navigate('/cursos')}>
          Volver a cursos
        </button>
      </div>
    );
  }

  if (!curso) {
    return (
      <div className="error-container">
        <div className="error">Curso no encontrado</div>
        <button className="btn-back" onClick={() => navigate('/cursos')}>
          Volver a cursos
        </button>
      </div>
    );
  }

  // Si el curso está en estado borrador y no es un administrador, redirigir
  if ((curso.estado === 'Borrador' || curso.estado === 'borrador') && 
      !(user && user.role === 'admin')) {
    return (
      <div className="error-container">
        <div className="error">Este curso no está disponible actualmente.</div>
        <button className="btn-back" onClick={() => navigate('/cursos')}>
          Volver a cursos
        </button>
      </div>
    );
  }

  // Determinar qué videos mostrar (compatibilidad con ambos formatos)
  const videosToShow = curso.videos || curso.capitulos || [];

  return (
    <div className="curso-detalle animate__animated animate__fadeIn">
      <h2>{curso.titulo}</h2>
      <p>{curso.descripcion}</p>

      {curso.estado === 'Borrador' && (
        <div className="draft-badge">
          Borrador - Este curso aún no está publicado
        </div>
      )}

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

      {error && paymentStatus !== 'aprobado' && (
        <div className="payment-error global-error">
          {error}
          {paymentStatus === 'no_pagado' && (
            <p>Puedes intentar comprar el curso de todos modos.</p>
          )}
        </div>
      )}

      {paymentStatus !== 'aprobado' && (
        <div className="payment-required">
          <h3>Para acceder a este curso necesitas comprarlo</h3>
          <p>Precio: ${curso.precio?.toFixed(2) || '29.99'}</p>
          
          {isLoggedIn ? (
            <>
              <button 
                onClick={() => setShowPaymentModal(true)}
                className="btn-pay"
              >
                Comprar Curso
              </button>
              
              {paymentStatus === 'rechazado' && (
                <div className="payment-error">
                  Tu pago fue rechazado. Por favor intenta con otro método.
                </div>
              )}
            </>
          ) : (
            <div className="login-required">
              Por favor inicia sesión para comprar este curso
            </div>
          )}
        </div>
      )}

      {paymentStatus === 'aprobado' && videosToShow.length > 0 && (
        <>
          <VideoPlayer
            videos={videosToShow}
            currentVideo={currentVideo}
            changeVideo={changeVideo}
            completedVideos={completedVideos}
            markAsCompleted={markAsCompleted}
          />

          <div className="contenido">
            <h3>Descripción del curso</h3>
            <p>{curso.contenido}</p>
          </div>
        </>
      )}

      {showPaymentModal && (
        <PaymentModal
          curso={curso}
          onClose={() => {
            setShowPaymentModal(false);
            setError(null);
          }}
          onSuccess={handlePaymentSuccess}
        />
      )}
    </div>
  );
}

export default CursoDetalle;