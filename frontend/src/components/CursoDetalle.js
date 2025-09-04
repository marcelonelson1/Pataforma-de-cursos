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
  const [hasAccess, setHasAccess] = useState(false);
  const [showPaymentModal, setShowPaymentModal] = useState(false);
  const [error, setError] = useState(null);
  const [loadingProgress, setLoadingProgress] = useState(false);
  
  // **NUEVOS ESTADOS PARA MANEJO DE PAGOS**
  const [paymentStatus, setPaymentStatus] = useState(null);
  const [isPaymentVerification, setIsPaymentVerification] = useState(false);

  // Cursos estáticos para fallback
  const cursosFallback = useMemo(() => [
    {
      id: 1,
      titulo: 'Curso de React',
      descripcion: 'Aprende React desde cero hasta un nivel avanzado con proyectos prácticos.',
      contenido: 'En este curso aprenderás todos los conceptos fundamentales de React, componentes, estados, props, hooks, y mucho más.',
      precio: 29.99,
      estado: 'Publicado',
      gratuito: false,
      capitulos: [
        {
          id: 1,
          orden: 1,
          titulo: 'Introducción a React',
          video_url: 'https://www.learningcontainer.com/wp-content/uploads/2020/05/sample-mp4-file.mp4',
          duracion: '08:45',
          publicado: true
        },
        {
          id: 2,
          orden: 2,
          titulo: 'Componentes y Props',
          video_url: 'https://filesamples.com/samples/video/mp4/sample_640x360.mp4',
          duracion: '10:15',
          publicado: true
        }
      ]
    },
  ], []);

  // Función para obtener la URL base del Course Service para videos
  const getCourseApiBaseUrl = () => {
    return process.env.REACT_APP_COURSE_API_URL || 'http://localhost:8003';
  };

  // Función para procesar URLs de video
  const getVideoUrl = (url) => {
    if (!url) return '';
    
    if (url.startsWith('http')) {
      return url;
    }
    
    if (url.startsWith('/static/')) {
      const baseUrl = getCourseApiBaseUrl();
      return `${baseUrl}${url}`;
    }
    
    if (url.startsWith('/api/files/')) {
      const baseUrl = getCourseApiBaseUrl();
      return `${baseUrl}${url}`;
    }
    
    return url;
  };

  // **FUNCIÓN ACTUALIZADA: Verificar parámetros de URL para pagos completados**
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const paymentStatus = urlParams.get('payment_status');
    const paymentSuccess = urlParams.get('success');
    const paymentCancelled = urlParams.get('cancelled');
    
    if (paymentStatus === 'success' || paymentSuccess === 'true') {
      console.log('Pago exitoso detectado en URL');
      setIsPaymentVerification(true);
      // Limpiar URL
      window.history.replaceState({}, document.title, window.location.pathname);
    } else if (paymentStatus === 'cancelled' || paymentCancelled === 'true') {
      console.log('Pago cancelado detectado en URL');
      setError('El pago fue cancelado. Puedes intentar nuevamente.');
      // Limpiar URL
      window.history.replaceState({}, document.title, window.location.pathname);
    }
  }, []);

  // **FUNCIÓN ACTUALIZADA: Carga inicial del curso con verificación de acceso mejorada**
  useEffect(() => {
    const fetchCurso = async () => {
      try {
        setLoading(true);
        console.log(`Obteniendo curso con ID: ${id}`);
        
        // **USAR NUEVA FUNCIÓN INTEGRADA**
        const informacionCompleta = await CursosService.getInformacionAccesoCompleta(id);
        console.log('Información completa del curso:', informacionCompleta);
        
        const cursoData = informacionCompleta.curso;
        const accesoInfo = informacionCompleta.acceso;
        
        // Verificar si el curso está publicado
        const esPublicado = cursoData.estado === 'Publicado' || cursoData.estado === 'publicado';
        const esAdmin = user && (user.role === 'admin' || user.role === 'instructor');
        
        if (!esPublicado && !esAdmin) {
          setError('Este curso no está disponible actualmente.');
          setCurso(null);
          setLoading(false);
          return;
        }
        
        // Procesar capítulos que ya vienen en el curso
        if (cursoData && Array.isArray(cursoData.capitulos)) {
          console.log('Capítulos encontrados en el curso:', cursoData.capitulos);
          
          const videosFormateados = cursoData.capitulos
            .filter(capitulo => capitulo.publicado || esAdmin)
            .sort((a, b) => (a.orden || 0) - (b.orden || 0))
            .map(capitulo => {
              console.log(`Procesando capítulo ${capitulo.id}: "${capitulo.titulo}"`);
              console.log(`  - video_url original:`, capitulo.video_url);
              console.log(`  - video_url procesada:`, getVideoUrl(capitulo.video_url));
              
              return {
                id: capitulo.id,
                indice: capitulo.orden?.toString() || "",
                titulo: capitulo.titulo,
                url: getVideoUrl(capitulo.video_url),
                video_url: getVideoUrl(capitulo.video_url),
                duracion: capitulo.duracion || "00:00",
                orden: capitulo.orden || 0
              };
            });
          
          cursoData.videos = videosFormateados;
          console.log('Videos formateados:', videosFormateados);
        } else {
          console.log('No se encontraron capítulos en el curso o no es un array:', cursoData.capitulos);
          cursoData.videos = [];
        }
        
        // **ESTABLECER ESTADO DE ACCESO BASADO EN LA NUEVA VERIFICACIÓN**
        setHasAccess(informacionCompleta.tiene_acceso);
        setPaymentStatus(informacionCompleta.estado_pago);
        
        setCurso(cursoData);
        setError(null);
        
        // Cargar progreso (funciona tanto con acceso como sin él)
        await fetchUserProgress();
        
      } catch (err) {
        console.error('Error al cargar curso:', err);
        setError('No se pudo cargar el curso. Por favor, intenta de nuevo más tarde.');
        
        // Fallback a curso estático
        const cursoEncontrado = cursosFallback.find((c) => c.id === parseInt(id));
        if (cursoEncontrado) {
          setCurso(cursoEncontrado);
        }
      } finally {
        setLoading(false);
        setIsPaymentVerification(false);
      }
    };

    if (id) {
      fetchCurso();
    }
  }, [id, cursosFallback, user, isPaymentVerification]);

  // Función para cargar progreso local como respaldo
  const loadLocalProgress = () => {
    try {
      if (!id) return;
      
      // Cargar progreso general del curso
      const overallKey = `course_overall_progress_${id}`;
      const savedOverallProgress = localStorage.getItem(overallKey);
      if (savedOverallProgress !== null) {
        setOverallProgress(parseInt(savedOverallProgress, 10));
        console.log('📥 Progreso general local cargado:', savedOverallProgress);
      }
      
      // Cargar último video visto
      const lastVideoKey = `course_last_video_${id}`;
      const savedLastVideo = localStorage.getItem(lastVideoKey);
      if (savedLastVideo !== null) {
        setCurrentVideo(parseInt(savedLastVideo, 10));
        console.log('📥 Último video local cargado:', savedLastVideo);
      }
      
      // Cargar videos completados
      const completedKey = `course_completed_videos_${id}`;
      const savedCompleted = localStorage.getItem(completedKey);
      if (savedCompleted) {
        try {
          const parsed = JSON.parse(savedCompleted);
          setCompletedVideos(parsed);
          console.log('📥 Videos completados local cargados:', parsed);
        } catch (e) {
          console.warn('Error parsing completed videos:', e);
        }
      }
    } catch (error) {
      console.error('Error cargando progreso local:', error);
    }
  };
  
  // Función para guardar progreso local
  const saveLocalProgress = (overallProg, completedVids, currentVid) => {
    try {
      if (!id) return;
      
      // Guardar progreso general
      if (overallProg !== undefined) {
        localStorage.setItem(`course_overall_progress_${id}`, overallProg.toString());
      }
      
      // Guardar videos completados
      if (completedVids) {
        localStorage.setItem(`course_completed_videos_${id}`, JSON.stringify(completedVids));
      }
      
      // Guardar último video visto
      if (currentVid !== undefined) {
        localStorage.setItem(`course_last_video_${id}`, currentVid.toString());
      }
      
      console.log('💾 Progreso local guardado');
    } catch (error) {
      console.error('Error guardando progreso local:', error);
    }
  };

  // **FUNCIÓN ACTUALIZADA: Cargar progreso del usuario**
  const fetchUserProgress = async () => {
    if (!isLoggedIn || !curso || !curso.id) {
      // Si no está logueado, cargar progreso local
      loadLocalProgress();
      return;
    }
    
    try {
      setLoadingProgress(true);
      console.log(`Obteniendo progreso para curso ID: ${curso.id}`);
      
      // Primero cargar progreso local (más rápido)
      loadLocalProgress();
      
      const response = await CursosService.getProgresoUsuario(curso.id);
      console.log('Respuesta de progreso:', response);
      
      if (response) {
        let updatedOverallProgress = overallProgress;
        let updatedCompletedVideos = { ...completedVideos };
        let updatedCurrentVideo = currentVideo;
        
        if (response.progreso_total !== undefined) {
          updatedOverallProgress = Math.round(response.progreso_total);
          setOverallProgress(updatedOverallProgress);
        }
        
        if (response.capitulos_progreso) {
          const completados = {};
          
          Object.values(response.capitulos_progreso).forEach(progreso => {
            if (progreso && progreso.capitulo_id) {
              completados[progreso.capitulo_id] = progreso.completado;
            }
          });
          
          updatedCompletedVideos = completados;
          setCompletedVideos(completados);
        }
        
        if (response.ultimo_capitulo && curso.videos && curso.videos.length > 0) {
          const videoIndex = curso.videos.findIndex(v => v.id === response.ultimo_capitulo);
          if (videoIndex !== -1) {
            updatedCurrentVideo = videoIndex;
            setCurrentVideo(videoIndex);
          }
        }
        
        // Guardar todo localmente como respaldo
        saveLocalProgress(updatedOverallProgress, updatedCompletedVideos, updatedCurrentVideo);
      }
    } catch (err) {
      console.error('Error al cargar progreso:', err);
      // En caso de error, mantener el progreso local
    } finally {
      setLoadingProgress(false);
    }
  };

  // Calcular progreso del curso
  useEffect(() => {
    if (curso) {
      const videos = curso.videos || curso.capitulos || [];
      if (videos.length > 0) {
        const completedCount = Object.values(completedVideos).filter(completed => completed).length;
        const totalVideos = videos.length;
        const newOverallProgress = Math.round((completedCount / totalVideos) * 100);
        
        if (newOverallProgress !== overallProgress) {
          setOverallProgress(newOverallProgress);
          // Guardar progreso general localmente
          saveLocalProgress(newOverallProgress, completedVideos, currentVideo);
        }
      }
    }
  }, [completedVideos, curso]); // Removido overallProgress de dependencias para evitar loops

  const changeVideo = async (index) => {
    const videos = curso?.videos || curso?.capitulos || [];
    if (index >= 0 && index < videos.length) {
      setCurrentVideo(index);
      
      // Guardar localmente siempre (funciona sin login)
      saveLocalProgress(overallProgress, completedVideos, index);
      
      if (isLoggedIn && hasAccess && videos[index] && videos[index].id) {
        try {
          await CursosService.guardarUltimoCapitulo(curso.id, videos[index].id);
        } catch (err) {
          console.error('Error al guardar último capítulo:', err);
        }
      }
    }
  };

  const markAsCompleted = async (videoId) => {
    const newCompletedState = !completedVideos[videoId];
    const updatedCompletedVideos = {
      ...completedVideos,
      [videoId]: newCompletedState
    };
    
    setCompletedVideos(updatedCompletedVideos);
    
    // Calcular nuevo progreso
    if (curso && (curso.videos || curso.capitulos)) {
      const videos = curso.videos || curso.capitulos;
      const completedCount = Object.values(updatedCompletedVideos).filter(completed => completed).length;
      const totalVideos = videos.length;
      const newOverallProgress = Math.round((completedCount / totalVideos) * 100);
      setOverallProgress(newOverallProgress);
      
      // Guardar localmente siempre
      saveLocalProgress(newOverallProgress, updatedCompletedVideos, currentVideo);
    }
    
    if (isLoggedIn && curso && hasAccess) {
      try {
        await CursosService.marcarCapituloCompletado(curso.id, videoId, newCompletedState);
      } catch (err) {
        console.error('Error al sincronizar estado completado:', err);
        // Revertir cambio local si falla el backend
        setCompletedVideos(prev => ({
          ...prev,
          [videoId]: !newCompletedState
        }));
        
        // Recalcular progreso con el estado revertido
        if (curso && (curso.videos || curso.capitulos)) {
          const videos = curso.videos || curso.capitulos;
          const revertedCompleted = { ...updatedCompletedVideos, [videoId]: !newCompletedState };
          const completedCount = Object.values(revertedCompleted).filter(completed => completed).length;
          const totalVideos = videos.length;
          const revertedProgress = Math.round((completedCount / totalVideos) * 100);
          setOverallProgress(revertedProgress);
          saveLocalProgress(revertedProgress, revertedCompleted, currentVideo);
        }
      }
    }
  };

  // **FUNCIÓN ACTUALIZADA: Manejo exitoso de pago**
  const handlePaymentSuccess = async () => {
    try {
      console.log('Pago exitoso - actualizando estado...');
      setShowPaymentModal(false);
      setError(null);
      setIsPaymentVerification(true);
      
      // **VERIFICAR NUEVAMENTE EL ACCESO DESPUÉS DEL PAGO**
      const informacionActualizada = await CursosService.getInformacionAccesoCompleta(id);
      console.log('Información actualizada después del pago:', informacionActualizada);
      
      // Actualizar estados
      setHasAccess(informacionActualizada.tiene_acceso);
      setPaymentStatus(informacionActualizada.estado_pago);
      setCurso(informacionActualizada.curso);
      
      // Cargar progreso actualizado
      await fetchUserProgress();
      
    } catch (err) {
      console.error('Error al actualizar después del pago:', err);
      setError('Pago procesado, pero hubo un error al actualizar. Recarga la página.');
    } finally {
      setIsPaymentVerification(false);
    }
  };

  // **FUNCIÓN NUEVA: Verificar estado de pago periódicamente si está pendiente**
  useEffect(() => {
    let intervalId;
    
    if (paymentStatus === 'pendiente' && isLoggedIn) {
      console.log('Pago pendiente detectado - iniciando verificación periódica');
      
      intervalId = setInterval(async () => {
        try {
          const estadoPago = await CursosService.verificarEstadoPago(id);
          console.log('Verificación periódica de pago:', estadoPago);
          
          if (estadoPago.estado === 'aprobado') {
            console.log('Pago aprobado detectado en verificación periódica');
            clearInterval(intervalId);
            await handlePaymentSuccess();
          } else if (estadoPago.estado === 'rechazado') {
            console.log('Pago rechazado detectado en verificación periódica');
            clearInterval(intervalId);
            setPaymentStatus('rechazado');
            setError('El pago fue rechazado. Por favor, intenta con otro método de pago.');
          }
        } catch (err) {
          console.error('Error en verificación periódica:', err);
        }
      }, 10000); // Verificar cada 10 segundos
    }
    
    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [paymentStatus, isLoggedIn, id]);

  if (loading || loadingProgress || isPaymentVerification) {
    return (
      <div className="loading">
        {isPaymentVerification ? 'Verificando pago...' : 'Cargando curso...'}
      </div>
    );
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

  // Verificar si es borrador
  if ((curso.estado === 'Borrador' || curso.estado === 'borrador') && 
      !(user && (user.role === 'admin' || user.role === 'instructor'))) {
    return (
      <div className="error-container">
        <div className="error">Este curso no está disponible actualmente.</div>
        <button className="btn-back" onClick={() => navigate('/cursos')}>
          Volver a cursos
        </button>
      </div>
    );
  }

  const videosToShow = curso.videos || curso.capitulos || [];
  const isFree = curso.gratuito || curso.precio === 0;

  return (
    <div className="curso-detalle animate__animated animate__fadeIn">
      <h2>{curso.titulo}</h2>
      <p>{curso.descripcion}</p>

      {curso.estado === 'Borrador' && (
        <div className="draft-badge">
          Borrador - Este curso aún no está publicado
        </div>
      )}

      {isFree && (
        <div className="free-badge">
          🎉 Este curso es completamente GRATIS
        </div>
      )}

      {/* **NUEVO: Mostrar estado de pago si está pendiente** */}
      {paymentStatus === 'pendiente' && (
        <div className="payment-pending-badge">
          ⏳ Pago pendiente - Verificando estado automáticamente...
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

      {error && !hasAccess && !isFree && (
        <div className="payment-error global-error">
          {error}
        </div>
      )}

      {!hasAccess && !isFree && (
        <div className="payment-required">
          <h3>Para acceder a este curso necesitas comprarlo</h3>
          <p>Precio: ${curso.precio?.toFixed(2) || '29.99'}</p>
          
          {/* **MEJORADO: Mostrar información específica según el estado del pago** */}
          {paymentStatus === 'pendiente' && (
            <div className="payment-status-info">
              <p>⏳ Tienes un pago pendiente para este curso. Estamos verificando el estado automáticamente.</p>
            </div>
          )}
          
          {paymentStatus === 'rechazado' && (
            <div className="payment-status-info payment-rejected">
              <p>❌ Tu último pago fue rechazado. Puedes intentar con otro método de pago.</p>
            </div>
          )}
          
          {isLoggedIn ? (
            <>
              <button 
                onClick={() => setShowPaymentModal(true)}
                className="btn-pay"
                disabled={paymentStatus === 'pendiente'}
              >
                {paymentStatus === 'pendiente' ? 'Verificando Pago...' : 'Comprar Curso'}
              </button>
            </>
          ) : (
            <div className="login-required">
              Por favor inicia sesión para comprar este curso
            </div>
          )}
        </div>
      )}

      {(hasAccess || isFree) && videosToShow.length > 0 && (
        <>
          <VideoPlayer
            videos={videosToShow}
            currentVideo={currentVideo}
            changeVideo={changeVideo}
            completedVideos={completedVideos}
            markAsCompleted={markAsCompleted}
            courseId={curso?.id}
          />

          <div className="contenido">
            <h3>Descripción del curso</h3>
            <p>{curso.contenido}</p>
          </div>
        </>
      )}

      {(hasAccess || isFree) && videosToShow.length === 0 && (
        <div className="no-content">
          Este curso aún no tiene contenido disponible.
        </div>
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