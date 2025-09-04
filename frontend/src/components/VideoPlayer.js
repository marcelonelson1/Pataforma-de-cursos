import React, { useRef, useState, useEffect, useCallback } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faBackward, faForward, faPlay, faPause,
  faCheckCircle, faChevronLeft, faChevronRight
} from '@fortawesome/free-solid-svg-icons';
import { faCircle as faCircleRegular } from '@fortawesome/free-regular-svg-icons';
import './CursoDetalle.css';

function VideoPlayer({ videos, currentVideo, changeVideo, completedVideos, markAsCompleted, courseId }) {
  const videoRef = useRef(null);
  const progressRef = useRef(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [progress, setProgress] = useState(0);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [videoDurations, setVideoDurations] = useState({});
  const [progressHistory, setProgressHistory] = useState({}); // Almacenar el progreso actual de cada video
  const [watchTimeHistory, setWatchTimeHistory] = useState({}); // Almacenar tiempo de visualización de cada video

  // Umbral para marcar automáticamente como completado (90% del video)
  const AUTO_COMPLETE_THRESHOLD = 0.9;

  // Función para guardar progreso localmente como respaldo
  const saveProgressLocally = (chapterId, timeWatched, progressPercent) => {
    try {
      const progressKey = `video_progress_${courseId}_${chapterId}`;
      const timeKey = `video_time_${courseId}_${chapterId}`;
      
      localStorage.setItem(progressKey, progressPercent.toString());
      localStorage.setItem(timeKey, timeWatched.toString());
      
      // También actualizar el estado local
      setProgressHistory(prev => ({ ...prev, [chapterId]: progressPercent }));
      setWatchTimeHistory(prev => ({ ...prev, [chapterId]: timeWatched }));
      
      console.log(`💾 Progreso guardado localmente: ${chapterId} - ${progressPercent}% - ${timeWatched}s`);
    } catch (error) {
      console.error('Error guardando progreso localmente:', error);
    }
  };

  // Función para cargar progreso local como respaldo
  const loadProgressLocally = () => {
    try {
      if (!courseId || !videos || videos.length === 0) return;
      
      const localProgress = {};
      const localWatchTime = {};
      
      videos.forEach(video => {
        if (video && video.id) {
          const progressKey = `video_progress_${courseId}_${video.id}`;
          const timeKey = `video_time_${courseId}_${video.id}`;
          
          const savedProgress = localStorage.getItem(progressKey);
          const savedTime = localStorage.getItem(timeKey);
          
          if (savedProgress !== null) {
            localProgress[video.id] = parseFloat(savedProgress);
          }
          if (savedTime !== null) {
            localWatchTime[video.id] = parseFloat(savedTime);
          }
        }
      });
      
      if (Object.keys(localProgress).length > 0) {
        setProgressHistory(localProgress);
        console.log('📥 Progreso local cargado:', localProgress);
      }
      
      if (Object.keys(localWatchTime).length > 0) {
        setWatchTimeHistory(localWatchTime);
        console.log('📥 Tiempo local cargado:', localWatchTime);
      }
    } catch (error) {
      console.error('Error cargando progreso local:', error);
    }
  };

  // Función para guardar progreso de visualización en el backend
  const saveWatchProgress = async (chapterId, timeWatched) => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        console.warn('No hay token para guardar progreso');
        return;
      }

      if (!courseId) {
        console.warn('No se pudo obtener curso ID para guardar progreso');
        return;
      }

      const baseUrl = process.env.REACT_APP_COURSE_API_URL || 'http://localhost:8003';
      const url = `${baseUrl}/api/progress/chapter/watch-time`;
      const payload = {
        curso_id: courseId,
        capitulo_id: chapterId,
        tiempo_visto: Math.max(1, timeWatched) // Asegurar que sea al menos 1 segundo
      };

      console.log(`💾 Guardando progreso:`, payload);

      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(payload)
      });

      console.log('📡 Respuesta guardado:', response.status, response.statusText);

      if (response.ok) {
        const responseData = await response.json();
        console.log(`✅ Progreso guardado exitosamente:`, responseData);
      } else {
        const errorText = await response.text();
        console.warn('❌ Error al guardar progreso:', response.status, errorText);
      }
    } catch (error) {
      console.error('❌ Error guardando progreso:', error);
    }
  };

  // Función para cargar progreso guardado del backend
  const loadSavedProgress = async (courseId) => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        console.warn('No hay token para cargar progreso');
        return;
      }

      const baseUrl = process.env.REACT_APP_COURSE_API_URL || 'http://localhost:8003';
      const url = `${baseUrl}/api/progress/course/${courseId}`;
      console.log('🌐 Llamando a:', url);
      
      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      console.log('📡 Respuesta del servidor:', response.status, response.statusText);

      if (response.ok) {
        const data = await response.json();
        console.log('📊 Datos recibidos:', data.success ? 'Progreso encontrado' : 'Sin progreso');
        
        // Verificar múltiples estructuras posibles
        let chaptersProgress = {};
        let foundProgress = false;
        
        // Estructura 1: data.progreso.capitulos_progreso
        if (data.success && data.progreso?.capitulos_progreso) {
          console.log('📋 Estructura 1: data.progreso.capitulos_progreso');
          Object.values(data.progreso.capitulos_progreso).forEach(chapter => {
            const progressPercent = chapter.progreso || 0;
            chaptersProgress[chapter.capitulo_id] = progressPercent;
          });
          foundProgress = true;
        }
        
        // Estructura 2: data.data.progreso.capitulos_progreso
        else if (data.success && data.data?.progreso?.capitulos_progreso) {
          console.log('📋 Estructura 2: data.data.progreso.capitulos_progreso');
          Object.values(data.data.progreso.capitulos_progreso).forEach(chapter => {
            const progressPercent = chapter.progreso || 0;
            chaptersProgress[chapter.capitulo_id] = progressPercent;
          });
          foundProgress = true;
        }
        
        // Estructura 3: data.data (directo)
        else if (data.success && data.data?.capitulos_progreso) {
          console.log('📋 Estructura 3: data.data.capitulos_progreso');
          Object.values(data.data.capitulos_progreso).forEach(chapter => {
            const progressPercent = chapter.progreso || 0;
            chaptersProgress[chapter.capitulo_id] = progressPercent;
          });
          foundProgress = true;
        }
        
        if (foundProgress) {
          setProgressHistory(chaptersProgress);
          console.log('📥 Progreso cargado desde backend:', chaptersProgress);
          
          // También actualizar el progreso local como respaldo
          Object.entries(chaptersProgress).forEach(([chapterId, progressPercent]) => {
            const progressKey = `video_progress_${courseId}_${chapterId}`;
            localStorage.setItem(progressKey, progressPercent.toString());
          });
        } else {
          console.log('📭 No hay progreso guardado o estructura no reconocida');
          console.log('🔍 Estructuras disponibles:', Object.keys(data));
          if (data.data) console.log('🔍 data.data keys:', Object.keys(data.data));
          
          // Si no hay progreso en backend, intentar cargar desde localStorage
          loadProgressLocally();
        }
      } else {
        const errorText = await response.text();
        console.error('❌ Error del servidor:', response.status, errorText);
      }
    } catch (error) {
      console.error('❌ Error cargando progreso:', error);
    }
  };

  // Función para marcar capítulo como completado en el backend
  const markChapterCompletedInBackend = async (chapterId, progressPercent) => {
    try {
      const token = localStorage.getItem('token');
      if (!token) return;

      if (!courseId) {
        console.warn('No se pudo obtener curso ID para marcar completado');
        return;
      }

      const baseUrl = process.env.REACT_APP_COURSE_API_URL || 'http://localhost:8003';
      const response = await fetch(`${baseUrl}/api/progress/chapter/complete`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          curso_id: courseId,
          capitulo_id: chapterId,
          completado: true,
          progreso: progressPercent,
          tiempo_visto: Math.floor(videoRef.current?.currentTime || 0)
        })
      });

      if (response.ok) {
        const data = await response.json();
        console.log(`🎉 Capítulo ${chapterId} marcado como completado en backend`);
      } else {
        console.warn('Error al marcar completado en backend:', response.status);
      }
    } catch (error) {
      console.error('Error marcando completado en backend:', error);
    }
  };

  // Función para procesar URLs de video del Course Service
  const getVideoUrl = (url) => {
    if (!url) return '';
    
    // Si ya es una URL completa, devolverla como está
    if (url.startsWith('http')) {
      return url;
    }
    
    // Si es una ruta relativa que empieza con /static, añadir la URL base del Course Service
    if (url.startsWith('/static/')) {
      const baseUrl = process.env.REACT_APP_COURSE_API_URL || 'http://localhost:8003';
      return `${baseUrl}${url}`;
    }
    
    // Si es una ruta del API de files
    if (url.startsWith('/api/files/')) {
      const baseUrl = process.env.REACT_APP_COURSE_API_URL || 'http://localhost:8003';
      return `${baseUrl}${url}`;
    }
    
    return url;
  };

  useEffect(() => {
    const videoElement = videoRef.current;
    console.log('🔧 useEffect ejecutándose, videoElement:', !!videoElement);

    if (videoElement) {
      console.log('✅ Agregando event listeners al video');
      const handleTimeUpdate = () => {
        const current = videoElement.currentTime;
        const videoDuration = videoElement.duration || 0;

        // Debug: verificar si la condición de guardado se cumple
        if (Math.floor(current) % 5 === 0 && Math.floor(current) !== Math.floor(current - 0.1)) {
          console.log(`⏰ Checkpoint cada 5s: ${current.toFixed(2)}s`);
          console.log(`🔍 Estado guardado:`, {
            currentVideo,
            videoId: videos[currentVideo]?.id,
            current: current,
            condicion: current > 0
          });
        }

        // Guardar progreso cada 10 segundos
        if (Math.floor(current) % 10 === 0 && Math.floor(current) !== Math.floor(current - 0.1)) {
          console.log(`💾 INTENTANDO GUARDAR: ${current.toFixed(2)}s / ${videoDuration.toFixed(2)}s`);
          // Guardar progreso en backend (solo si el tiempo es > 0 para evitar error de validación)
          if (videos[currentVideo] && videos[currentVideo].id && current > 0) {
            const videoId = videos[currentVideo].id;
            const progressPercent = videoDuration > 0 ? (current / videoDuration) * 100 : 0;
            
            console.log(`✅ GUARDANDO para currentVideo: ${currentVideo}, videoId: ${videoId}`);
            
            // Guardar en backend
            saveWatchProgress(videoId, Math.floor(current));
            
            // Guardar localmente como respaldo
            saveProgressLocally(videoId, Math.floor(current), progressPercent);
          } else {
            console.log(`❌ NO guardando:`, {
              hasVideos: !!videos[currentVideo],
              hasId: !!videos[currentVideo]?.id,
              currentGreaterThan0: current > 0
            });
          }
        }

        if (videoDuration > 0) {
          const progressPercent = (current / videoDuration) * 100;
          setProgress(progressPercent);
          setCurrentTime(current);
          
          // Actualizar progreso para el video actual
          if (videos[currentVideo] && videos[currentVideo].id) {
            const videoId = videos[currentVideo].id;
            setProgressHistory(prev => ({
              ...prev,
              [videoId]: progressPercent
            }));
            
            // Guardar progreso local cada 5 segundos (para reducir escrituras y evitar loops)
            if (Math.floor(current) % 5 === 0 && Math.floor(current) !== Math.floor(current - 0.1)) {
              saveProgressLocally(videoId, Math.floor(current), progressPercent);
            }
            
            // Comprobar si el video ha sido visto en un 90% y no está marcado como completado
            const isNearEnd = progressPercent > (AUTO_COMPLETE_THRESHOLD * 100);
            if (isNearEnd && videoId && !completedVideos[videoId]) {
              console.log(`✅ Auto-marcando capítulo ${videoId} como completado (${progressPercent.toFixed(1)}%)`);
              markAsCompleted(videoId);
              // También marcar en el backend
              markChapterCompletedInBackend(videoId, progressPercent);
            }
          }
        }
      };

      const handleDurationChange = () => {
        const videoDuration = videoElement.duration || 0;
        console.log('📏 Duración del video cargada:', videoDuration);
        setDuration(videoDuration);
        
        if (videos[currentVideo] && videos[currentVideo].id) {
          setVideoDurations(prev => ({
            ...prev,
            [videos[currentVideo].id]: videoDuration
          }));
        }
      };

      const handlePlay = () => {
        console.log('🎵 Video reproduciendo');
        setIsPlaying(true);
      };

      const handlePause = () => {
        console.log('⏸️ Video pausado');
        setIsPlaying(false);
      };

      const handleEnded = () => {
        setIsPlaying(false);
        if (videos[currentVideo] && videos[currentVideo].id) {
          markAsCompleted(videos[currentVideo].id);
        }

        if (currentVideo < videos.length - 1) {
          changeVideo(currentVideo + 1);
        }
      };

      const handleError = (e) => {
        console.error('Error al cargar video:', e);
        console.log('URL del video que falló:', videoElement.src);
      };

      videoElement.addEventListener('timeupdate', handleTimeUpdate);
      videoElement.addEventListener('durationchange', handleDurationChange);
      videoElement.addEventListener('play', handlePlay);
      videoElement.addEventListener('pause', handlePause);
      videoElement.addEventListener('ended', handleEnded);
      videoElement.addEventListener('error', handleError);

      return () => {
        console.log('🗑️ Removiendo event listeners del video');
        videoElement.removeEventListener('timeupdate', handleTimeUpdate);
        videoElement.removeEventListener('durationchange', handleDurationChange);
        videoElement.removeEventListener('play', handlePlay);
        videoElement.removeEventListener('pause', handlePause);
        videoElement.removeEventListener('ended', handleEnded);
        videoElement.removeEventListener('error', handleError);
      };
    } else {
      console.log('❌ No hay videoElement para agregar listeners');
    }
  }, [currentVideo]); // Solo depender de currentVideo para evitar re-renders constantes

  const togglePlay = async () => {
    if (videoRef.current) {
      console.log('🔍 Estado antes de toggle:', {
        paused: videoRef.current.paused,
        currentTime: videoRef.current.currentTime,
        duration: videoRef.current.duration,
        readyState: videoRef.current.readyState,
        networkState: videoRef.current.networkState,
        src: videoRef.current.src
      });
      
      try {
        if (videoRef.current.paused) {
          console.log('▶️ Iniciando reproducción...');
          await videoRef.current.play();
          console.log('✅ Reproducción iniciada');
        } else {
          console.log('⏸️ Pausando...');
          videoRef.current.pause();
          console.log('✅ Video pausado');
        }
        
        // Estado después del toggle
        setTimeout(() => {
          console.log('🔍 Estado después de toggle:', {
            paused: videoRef.current?.paused,
            currentTime: videoRef.current?.currentTime,
            duration: videoRef.current?.duration
          });
        }, 100);
        
      } catch (error) {
        console.error('❌ Error al reproducir:', error.message);
      }
    }
  };

  const handleBackward = () => {
    if (videoRef.current) {
      videoRef.current.currentTime = Math.max(0, videoRef.current.currentTime - 10);
    }
  };

  const handleForward = () => {
    if (videoRef.current) {
      videoRef.current.currentTime = Math.min(
        videoRef.current.duration || 0,
        videoRef.current.currentTime + 10
      );
    }
  };

  const handleProgressClick = (e) => {
    if (progressRef.current && videoRef.current) {
      const progressBar = progressRef.current;
      const rect = progressBar.getBoundingClientRect();
      const position = (e.clientX - rect.left) / progressBar.offsetWidth;
      const newTime = position * (videoRef.current.duration || 0);

      if (newTime >= 0 && newTime <= (videoRef.current.duration || 0)) {
        videoRef.current.currentTime = newTime;
      }
    }
  };

  const formatTime = (timeInSeconds) => {
    if (isNaN(timeInSeconds) || timeInSeconds < 0) return "0:00";

    const minutes = Math.floor(timeInSeconds / 60);
    const seconds = Math.floor(timeInSeconds % 60);
    return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
  };

  const loadVideoDuration = useCallback((videoUrl, videoId) => {
    if (videoDurations[videoId] || !videoUrl) return;

    const tempVideo = document.createElement('video');
    tempVideo.preload = 'metadata';
    
    tempVideo.onloadedmetadata = function() {
      setVideoDurations(prev => ({
        ...prev,
        [videoId]: tempVideo.duration
      }));
    };
    
    // Manejar errores para evitar console spam
    tempVideo.onerror = function() {
      console.warn(`No se pudo cargar metadatos para el video: ${videoUrl}`);
    };
    
    tempVideo.src = getVideoUrl(videoUrl);
  }, [videoDurations]);

  useEffect(() => {
    if (!Array.isArray(videos)) return;
    
    videos.forEach(video => {
      if (!video || !video.id) return;
      
      // Intentar obtener la URL del video desde diferentes estructuras de datos
      const videoUrl = video.url || video.video_url;
      if (!videoDurations[video.id] && videoUrl) {
        loadVideoDuration(videoUrl, video.id);
      }
    });

    // Cargar progreso guardado si tenemos courseId
    if (courseId && videos.length > 0) {
      console.log('🔍 Cargando progreso para courseId:', courseId);
      // Primero cargar progreso local (es más rápido)
      loadProgressLocally();
      // Luego cargar y sincronizar con backend
      loadSavedProgress(courseId);
    } else if (!courseId) {
      console.warn('⚠️ No se proporcionó courseId para cargar progreso');
    }
  }, [videos, videoDurations, loadVideoDuration, courseId]);

  useEffect(() => {
    if (videoRef.current) {
      setIsPlaying(false);
      
      // NO restablecer el progreso a 0 - solo actualizar el estado de reproducción
      // El progreso se restaurará en el siguiente useEffect
    }
  }, [currentVideo, videos]);

  // Efecto separado para restaurar posición (solo cuando duration está disponible)
  useEffect(() => {
    if (videoRef.current && videoRef.current.duration > 0 && videos[currentVideo]?.id) {
      const videoId = videos[currentVideo].id;
      let savedProgressPercent = progressHistory[videoId];
      let savedTime = watchTimeHistory[videoId];
      
      // Si no hay progreso en memoria, intentar cargar desde localStorage
      if (!savedProgressPercent && courseId) {
        const progressKey = `video_progress_${courseId}_${videoId}`;
        const timeKey = `video_time_${courseId}_${videoId}`;
        const localProgress = localStorage.getItem(progressKey);
        const localTime = localStorage.getItem(timeKey);
        
        if (localProgress) {
          savedProgressPercent = parseFloat(localProgress);
        }
        if (localTime) {
          savedTime = parseFloat(localTime);
        }
      }
      
      // Preferir tiempo exacto si está disponible
      if (savedTime > 0 && savedTime < (videoRef.current.duration * 0.95)) {
        console.log(`🔄 Restaurando tiempo exacto: ${savedTime.toFixed(2)}s`);
        videoRef.current.currentTime = savedTime;
        setCurrentTime(savedTime);
        setProgress((savedTime / videoRef.current.duration) * 100);
      }
      // Fallback a porcentaje de progreso
      else if (savedProgressPercent > 0 && savedProgressPercent < 95) {
        const calculatedTime = (savedProgressPercent / 100) * videoRef.current.duration;
        console.log(`🔄 Restaurando desde porcentaje: ${calculatedTime.toFixed(2)}s (${savedProgressPercent.toFixed(1)}%)`);
        videoRef.current.currentTime = calculatedTime;
        setCurrentTime(calculatedTime);
        setProgress(savedProgressPercent);
      }
    }
  }, [duration, currentVideo]); // SOLO depender de duration y currentVideo para evitar loops infinitos


  if (!videos || !Array.isArray(videos) || videos.length === 0) {
    return <div className="no-videos">No hay videos disponibles para este curso</div>;
  }

  // Obtener el video actual
  const currentVideoObj = videos[currentVideo] || {};
  
  // Adaptación para soportar ambas estructuras de datos (la original y la del backend)
  const videoUrl = getVideoUrl(currentVideoObj.url || currentVideoObj.video_url || '');
  const videoTitle = currentVideoObj.titulo || '';
  const videoIndex = currentVideoObj.indice || currentVideoObj.orden || (currentVideo + 1).toString();
  const videoDuration = currentVideoObj.duracion || '';

  return (
    <div className="video-section">
      <div className="video-container">
        <video
          ref={videoRef}
          className="video-player"
          src={videoUrl}
          preload="auto"
          playsInline
          crossOrigin="anonymous"
        ></video>
        
        {!videoUrl && (
          <div className="video-error">
            No hay video disponible para este capítulo
          </div>
        )}
      </div>

      <div className="video-controls">
        <button className="control-button" onClick={handleBackward} title="Retroceder 10 segundos">
          <FontAwesomeIcon icon={faBackward} />
        </button>

        <button className="control-button" onClick={togglePlay} title={isPlaying ? "Pausar" : "Reproducir"}>
          <FontAwesomeIcon icon={isPlaying ? faPause : faPlay} />
        </button>

        <button className="control-button" onClick={handleForward} title="Avanzar 10 segundos">
          <FontAwesomeIcon icon={faForward} />
        </button>

        <div
          className="video-progress"
          ref={progressRef}
          onClick={handleProgressClick}
        >
          <div
            className="progress-bar"
            style={{ width: `${progress}%` }}
          ></div>
        </div>

        <div className="time-display">
          {formatTime(currentTime)} / {formatTime(duration)}
        </div>
      </div>

      <div className="video-navigation">
        <div className="lesson-navigation">
          <button
            className="nav-button"
            onClick={() => changeVideo(currentVideo - 1)}
            disabled={currentVideo === 0}
          >
            <FontAwesomeIcon icon={faChevronLeft} /> Video anterior
          </button>

          <button
            className={`complete-button ${currentVideoObj.id && completedVideos[currentVideoObj.id] ? 'completed' : ''}`}
            onClick={() => {
              if (currentVideoObj.id) {
                console.log(`🖱️ Click manual en marcar completado para capítulo ${currentVideoObj.id}`);
                markAsCompleted(currentVideoObj.id);
                // Si se marca como completado manualmente, también actualizar en backend
                if (!completedVideos[currentVideoObj.id]) {
                  console.log(`📤 Marcando completado en backend (manual)`);
                  markChapterCompletedInBackend(currentVideoObj.id, 100);
                }
              }
            }}
            title={currentVideoObj.id && completedVideos[currentVideoObj.id] ? "Marcar como no completado" : "Marcar como completado"}
          >
            {currentVideoObj.id && completedVideos[currentVideoObj.id] ?
              <><FontAwesomeIcon icon={faCheckCircle} /> Completado</> :
              <><FontAwesomeIcon icon={faCircleRegular} /> Marcar como completado</>
            }
          </button>

          <button
            className="nav-button"
            onClick={() => changeVideo(currentVideo + 1)}
            disabled={currentVideo === videos.length - 1}
          >
            Video siguiente <FontAwesomeIcon icon={faChevronRight} />
          </button>
        </div>
      </div>

      <div className="video-title">
        <h3>
          <span className="video-index">{videoIndex}.</span> {videoTitle}
        </h3>
      </div>

      <div className="lessons-list">
        <h3>Contenido del curso</h3>
        {videos.map((video, index) => {
          if (!video) return null;
          
          // Adaptación para soportar ambas estructuras
          const lessonIndex = video.indice || video.orden || (index + 1).toString();
          const lessonTitle = video.titulo || '';
          const lessonDuration = videoDurations[video.id] ? 
            formatTime(videoDurations[video.id]) : 
            (video.duracion || '');
          
          // Variables no usadas pero necesarias para la lógica
          // eslint-disable-next-line no-unused-vars
          
          return (
            <div
              key={video.id || index}
              className={`lesson-item ${index === currentVideo ? 'active' : ''} ${video.id && completedVideos[video.id] ? 'completed' : ''}`}
              onClick={() => changeVideo(index)}
            >
              <span className="lesson-index">{lessonIndex}.</span>
              <span className="lesson-title">{lessonTitle}</span>
              <span className="lesson-duration">{lessonDuration}</span>
              {video.id && completedVideos[video.id] && (
                <FontAwesomeIcon icon={faCheckCircle} className="lesson-check" />
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default VideoPlayer;