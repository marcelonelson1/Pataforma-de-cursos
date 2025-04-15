import React, { useRef, useState, useEffect, useCallback } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faBackward, faForward, faPlay, faPause,
  faCheckCircle, faChevronLeft, faChevronRight
} from '@fortawesome/free-solid-svg-icons';
import { faCircle as faCircleRegular } from '@fortawesome/free-regular-svg-icons';
import './CursoDetalle.css';

function VideoPlayer({ videos, currentVideo, changeVideo, completedVideos, markAsCompleted }) {
  const videoRef = useRef(null);
  const progressRef = useRef(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [progress, setProgress] = useState(0);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [videoDurations, setVideoDurations] = useState({});
  const [progressHistory, setProgressHistory] = useState({}); // Almacenar el progreso actual de cada video

  // Umbral para marcar automáticamente como completado (90% del video)
  const AUTO_COMPLETE_THRESHOLD = 0.9;

  useEffect(() => {
    const videoElement = videoRef.current;

    if (videoElement) {
      const handleTimeUpdate = () => {
        const current = videoElement.currentTime;
        const videoDuration = videoElement.duration || 0;

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
            
            // Comprobar si el video ha sido visto en un 90% y no está marcado como completado
            const isNearEnd = progressPercent > (AUTO_COMPLETE_THRESHOLD * 100);
            if (isNearEnd && videoId && !completedVideos[videoId]) {
              markAsCompleted(videoId);
            }
          }
        }
      };

      const handleDurationChange = () => {
        const videoDuration = videoElement.duration || 0;
        setDuration(videoDuration);
        
        if (videos[currentVideo] && videos[currentVideo].id) {
          setVideoDurations(prev => ({
            ...prev,
            [videos[currentVideo].id]: videoDuration
          }));
        }
      };

      const handlePlay = () => {
        setIsPlaying(true);
      };

      const handlePause = () => {
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

      videoElement.addEventListener('timeupdate', handleTimeUpdate);
      videoElement.addEventListener('durationchange', handleDurationChange);
      videoElement.addEventListener('play', handlePlay);
      videoElement.addEventListener('pause', handlePause);
      videoElement.addEventListener('ended', handleEnded);

      return () => {
        videoElement.removeEventListener('timeupdate', handleTimeUpdate);
        videoElement.removeEventListener('durationchange', handleDurationChange);
        videoElement.removeEventListener('play', handlePlay);
        videoElement.removeEventListener('pause', handlePause);
        videoElement.removeEventListener('ended', handleEnded);
      };
    }
  }, [currentVideo, videos, completedVideos, markAsCompleted, changeVideo]);

  const togglePlay = () => {
    if (videoRef.current) {
      if (videoRef.current.paused) {
        videoRef.current.play();
      } else {
        videoRef.current.pause();
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
    
    tempVideo.src = videoUrl;
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
  }, [videos, videoDurations, loadVideoDuration]);

  useEffect(() => {
    if (videoRef.current) {
      setIsPlaying(false);
      
      // Restablece el progreso actual cuando cambia de video
      setProgress(0);
      setCurrentTime(0);
      
      // Si hay progreso guardado para este video, establece la posición
      if (videos[currentVideo] && videos[currentVideo].id) {
        const videoId = videos[currentVideo].id;
        if (progressHistory[videoId] > 0) {
          const savedProgressPercent = progressHistory[videoId];
          const savedTime = (savedProgressPercent / 100) * (videoRef.current.duration || 0);
          
          // Solo si el progreso es menor al 95% (para evitar siempre iniciar al final)
          if (savedProgressPercent < 95) {
            videoRef.current.currentTime = savedTime;
          }
        }
      }
    }
  }, [currentVideo, progressHistory, videos]);

  if (!videos || !Array.isArray(videos) || videos.length === 0) {
    return <div className="no-videos">No hay videos disponibles para este curso</div>;
  }

  // Obtener el video actual
  const currentVideoObj = videos[currentVideo] || {};
  
  // Adaptación para soportar ambas estructuras de datos (la original y la del backend)
  const videoUrl = currentVideoObj.url || currentVideoObj.video_url || '';
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
          poster="/path-to-poster-image.jpg"
          preload="metadata"
          controlsList="nodownload"
        ></video>
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
            onClick={() => currentVideoObj.id && markAsCompleted(currentVideoObj.id)}
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
          
          return (
            <div
              key={video.id || index}
              className={`lesson-item ${index === currentVideo ? 'active' : ''} ${video.id && completedVideos[video.id] ? 'completed' : ''}`}
              onClick={() => changeVideo(index)}
            >
              <span className="lesson-index">{lessonIndex}.</span>
              <span className="lesson-title">{lessonTitle}</span>
              <span className="lesson-duration">{lessonDuration}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default VideoPlayer;