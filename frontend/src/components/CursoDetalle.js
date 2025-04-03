import React, { useState, useEffect, useMemo } from 'react';
import { useParams } from 'react-router-dom';
import VideoPlayer from './VideoPlayer';
import './CursoDetalle.css';
import { useAuth } from '../context/AuthContext';
import PaymentModal from './PaymentModal';

function CursoDetalle() {
  const { id } = useParams();
  const { isLoggedIn } = useAuth();
  const [curso, setCurso] = useState(null);
  const [loading, setLoading] = useState(true);
  const [currentVideo, setCurrentVideo] = useState(0);
  const [completedVideos, setCompletedVideos] = useState({});
  const [overallProgress, setOverallProgress] = useState(0);
  const [paymentStatus, setPaymentStatus] = useState(null);
  const [showPaymentModal, setShowPaymentModal] = useState(false);
  const [error, setError] = useState(null);

  const cursos = useMemo(() => [
    {
      id: 1,
      titulo: 'Curso de React',
      descripcion: 'Aprende React desde cero hasta un nivel avanzado con proyectos prácticos.',
      contenido: 'En este curso aprenderás todos los conceptos fundamentales de React, componentes, estados, props, hooks, y mucho más.',
      precio: 29.99,
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

  useEffect(() => {
    const cursoEncontrado = cursos.find((c) => c.id === parseInt(id));
    if (cursoEncontrado) {
      setCurso(cursoEncontrado);
    } else {
      setError('Curso no encontrado');
    }
    setLoading(false);
  }, [id, cursos]);

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
        
        // Realizamos la petición con la URL completa del backend para evitar problemas de CORS
        const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
        const response = await fetch(`${apiUrl}/api/pagos/${curso.id}`, {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Accept': 'application/json'
          }
        });
        
        // Imprimir información de la respuesta para depuración
        console.log(`Respuesta recibida con estado: ${response.status} ${response.statusText}`);
        console.log(`Content-Type: ${response.headers.get('content-type')}`);
        
        // Obtener el texto de la respuesta primero para depurar
        const textResponse = await response.text();
        console.log('Respuesta como texto:', textResponse);
        
        // Intentar parsear manualmente el texto a JSON
        let data;
        try {
          data = JSON.parse(textResponse);
          console.log('Datos parseados:', data);
        } catch (jsonError) {
          console.error('Error al parsear JSON:', jsonError);
          console.error('Respuesta no procesable como JSON. Contenido:', textResponse);
          
          if (textResponse.includes('<!DOCTYPE html>') || textResponse.includes('<html>')) {
            setError('El servidor devolvió una página HTML en lugar de JSON. Esto puede indicar que el backend no está procesando correctamente la solicitud.');
          } else {
            setError('Formato de respuesta del servidor incorrecto');
          }
          
          setPaymentStatus('no_pagado');
          return;
        }
        
        if (!response.ok) {
          throw new Error(data.error || `Error ${response.status}`);
        }
        
        // Acceder directamente al estado del pago (sin anidamiento)
        const estadoPago = data.estado || 'no_pagado';
        setPaymentStatus(estadoPago);
        setError(null);
        
      } catch (err) {
        console.error('Error al verificar estado de pago:', err);
        
        let errorMessage = err.message;
        if (err.message.includes('Failed to fetch')) {
          errorMessage = 'Error de conexión con el servidor. Asegúrate de que el backend esté en ejecución.';
        } else if (err.message.includes('404')) {
          errorMessage = 'No se encontró registro de pago';
        } else if (err.message.includes('401')) {
          errorMessage = 'No autorizado - por favor inicia sesión nuevamente';
        } else if (err.message.includes('JSON')) {
          errorMessage = 'Formato de respuesta del servidor incorrecto';
        }
        
        setError(errorMessage);
        setPaymentStatus('no_pagado');
      }
    };

    checkPaymentStatus();
  }, [isLoggedIn, curso]);

  useEffect(() => {
    if (curso && Object.keys(completedVideos).length > 0) {
      const completedCount = Object.values(completedVideos).filter(completed => completed).length;
      const totalVideos = curso.videos.length;
      const newOverallProgress = Math.round((completedCount / totalVideos) * 100);
      setOverallProgress(newOverallProgress);
    }
  }, [completedVideos, curso]);

  const changeVideo = (index) => {
    if (index >= 0 && index < curso?.videos?.length) {
      setCurrentVideo(index);
    }
  };

  const markAsCompleted = (videoId) => {
    setCompletedVideos(prev => ({
      ...prev,
      [videoId]: !prev[videoId]
    }));
  };

  const handlePaymentSuccess = () => {
    setPaymentStatus('aprobado');
    setShowPaymentModal(false);
    setError(null);
  };

  if (loading) {
    return <div className="loading">Cargando curso...</div>;
  }

  if (error && !curso) {
    return <div className="error">{error}</div>;
  }

  if (!curso) {
    return <div className="error">Curso no encontrado</div>;
  }

  return (
    <div className="curso-detalle animate__animated animate__fadeIn">
      <h2>{curso.titulo}</h2>
      <p>{curso.descripcion}</p>

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

      {paymentStatus === 'aprobado' && (
        <>
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