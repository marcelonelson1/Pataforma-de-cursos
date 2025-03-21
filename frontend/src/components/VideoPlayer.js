import React, { useRef, useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faBackward, faForward, faPlay, faPause,
  faCheckCircle, faChevronLeft, faChevronRight
} from '@fortawesome/free-solid-svg-icons';
import { faCircle as faCircleRegular } from '@fortawesome/free-regular-svg-icons'; // Corregido
import './CursoDetalle.css';

function VideoPlayer({ videos }) {
  const videoRef = useRef(null);
  const progressRef = useRef(null);
  const [currentVideo, setCurrentVideo] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const [progress, setProgress] = useState(0);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [completedVideos, setCompletedVideos] = useState({});

  // Manejar eventos del video
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
        }
      };

      const handleDurationChange = () => {
        setDuration(videoElement.duration || 0);
      };

      const handlePlay = () => {
        setIsPlaying(true);
      };

      const handlePause = () => {
        setIsPlaying(false);
      };

      const handleEnded = () => {
        setIsPlaying(false);
        setCompletedVideos(prev => ({
          ...prev,
          [videos[currentVideo].id]: true
        }));

        if (currentVideo < videos.length - 1) {
          setCurrentVideo(currentVideo + 1);
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
  }, [currentVideo, videos]);

  // Cambiar de video
  const changeVideo = (index) => {
    if (index >= 0 && index < videos.length) {
      setCurrentVideo(index);
    }
  };

  // Funciones de control de reproducción
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

  // Formatear tiempo
  const formatTime = (timeInSeconds) => {
    if (isNaN(timeInSeconds) || timeInSeconds < 0) return "0:00";

    const minutes = Math.floor(timeInSeconds / 60);
    const seconds = Math.floor(timeInSeconds % 60);
    return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
  };

  // Marcar video como completado manualmente
  const markAsCompleted = (videoId) => {
    setCompletedVideos(prev => ({
      ...prev,
      [videoId]: !prev[videoId]
    }));
  };

  return (
    <div className="video-section">
      <div className="video-container">
        <video
          ref={videoRef}
          className="video-player"
          src={videos[currentVideo]?.url}
          poster="/path-to-poster-image.jpg"
          preload="metadata"
          controlsList="nodownload"
        ></video>
      </div>

      {/* Controles personalizados */}
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

      {/* Navegación entre videos y botón de completado */}
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
            className={`complete-button ${completedVideos[videos[currentVideo].id] ? 'completed' : ''}`}
            onClick={() => markAsCompleted(videos[currentVideo].id)}
            title={completedVideos[videos[currentVideo].id] ? "Marcar como no completado" : "Marcar como completado"}
          >
            {completedVideos[videos[currentVideo].id] ?
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
    </div>
  );
}

export default VideoPlayer;