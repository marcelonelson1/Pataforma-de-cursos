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
        const videoDuration = videoElement.duration || 0;
        setDuration(videoDuration);
        
        setVideoDurations(prev => ({
          ...prev,
          [videos[currentVideo].id]: videoDuration
        }));
      };

      const handlePlay = () => {
        setIsPlaying(true);
      };

      const handlePause = () => {
        setIsPlaying(false);
      };

      const handleEnded = () => {
        setIsPlaying(false);
        markAsCompleted(videos[currentVideo].id);

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
  }, [currentVideo, videos, markAsCompleted, changeVideo]);

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
    if (videoDurations[videoId]) return;

    const tempVideo = document.createElement('video');
    tempVideo.preload = 'metadata';
    
    tempVideo.onloadedmetadata = function() {
      setVideoDurations(prev => ({
        ...prev,
        [videoId]: tempVideo.duration
      }));
    };
    
    tempVideo.src = videoUrl;
  }, [videoDurations]);

  useEffect(() => {
    videos.forEach(video => {
      if (!videoDurations[video.id]) {
        loadVideoDuration(video.url, video.id);
      }
    });
  }, [videos, videoDurations, loadVideoDuration]);

  useEffect(() => {
    if (videoRef.current) {
      setIsPlaying(false);
    }
  }, [currentVideo]);

  if (!videos || videos.length === 0) {
    return <div className="no-videos">No hay videos disponibles para este curso</div>;
  }

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

      <div className="video-title">
        <h3>
          <span className="video-index">{videos[currentVideo].indice || (currentVideo + 1)}.</span> {videos[currentVideo].titulo}
        </h3>
      </div>

      <div className="lessons-list">
        <h3>Contenido del curso</h3>
        {videos.map((video, index) => (
          <div
            key={video.id}
            className={`lesson-item ${index === currentVideo ? 'active' : ''} ${completedVideos[video.id] ? 'completed' : ''}`}
            onClick={() => changeVideo(index)}
          >
            <span className="lesson-index">{video.indice || (index + 1)}.</span>
            <span className="lesson-title">{video.titulo}</span>
            <span className="lesson-duration">
              {videoDurations[video.id] ? formatTime(videoDurations[video.id]) : video.duracion}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}

export default VideoPlayer;