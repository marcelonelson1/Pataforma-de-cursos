import React, { useState, useEffect, useRef } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faPlay, faEdit, faTrashAlt, faPlus, faGripLines,
  faSave, faTimes, faVideo, faImage, faArrowLeft, faCheck
} from '@fortawesome/free-solid-svg-icons';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';
import './AdminPages.css';

function AdminCapitulos() {
  const { id } = useParams(); // ID del curso
  const navigate = useNavigate();
  
  const [curso, setCurso] = useState(null);
  const [capitulos, setCapitulos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  
  // Estado para el formulario modal de capítulo
  const [showModal, setShowModal] = useState(false);
  const [modalMode, setModalMode] = useState('create'); // 'create' o 'edit'
  const [currentCapitulo, setCurrentCapitulo] = useState({
    id: null,
    titulo: '',
    descripcion: '',
    video_url: '',
    thumbnail_url: ''
  });
  const [videoFile, setVideoFile] = useState(null);
  const [thumbnailFile, setThumbnailFile] = useState(null);
  const [videoPreview, setVideoPreview] = useState('');
  const [thumbnailPreview, setThumbnailPreview] = useState('');
  const [uploadProgress, setUploadProgress] = useState(0);
  const [uploading, setUploading] = useState(false);
  const [formErrors, setFormErrors] = useState({});

  // Referencia para el elemento de video
  const videoRef = useRef(null);
  
  useEffect(() => {
    const fetchCursoYCapitulos = async () => {
      try {
        const token = localStorage.getItem('token');
        const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
        
        // Obtener curso
        const cursoResponse = await fetch(`${apiUrl}/api/admin/cursos/${id}`, {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });

        if (!cursoResponse.ok) {
          throw new Error('Error al obtener el curso');
        }

        const cursoData = await cursoResponse.json();
        setCurso(cursoData.curso);
        setCapitulos(cursoData.capitulos || []);
        
        setError(null);
      } catch (err) {
        console.error('Error al cargar curso y capítulos:', err);
        setError('No se pudieron cargar los datos. Por favor, intenta de nuevo más tarde.');
      } finally {
        setLoading(false);
      }
    };

    fetchCursoYCapitulos();
  }, [id]);

  // Datos de ejemplo mientras se desarrolla la API
  useEffect(() => {
    if (!loading && capitulos.length === 0 && id === "1") {
      // Curso de ejemplo
      setCurso({
        id: 1,
        titulo: 'Curso de React',
        descripcion: 'Aprende React desde cero hasta un nivel avanzado.'
      });
      
      // Capítulos de ejemplo
      setCapitulos([
        {
          id: 1,
          curso_id: 1,
          indice: "1",
          titulo: 'Introducción a React',
          descripcion: 'Fundamentos básicos de React y su ecosistema.',
          video_url: 'https://www.learningcontainer.com/wp-content/uploads/2020/05/sample-mp4-file.mp4',
          thumbnail_url: '/preview-1.jpg',
          duracion: 540, // en segundos (9 minutos)
          orden: 1
        },
        {
          id: 2,
          curso_id: 1,
          indice: "2",
          titulo: 'Componentes y Props',
          descripcion: 'Creación de componentes y comunicación a través de props.',
          video_url: 'https://filesamples.com/samples/video/mp4/sample_640x360.mp4',
          thumbnail_url: '/preview-2.jpg',
          duracion: 615, // en segundos (10:15 minutos)
          orden: 2
        },
        {
          id: 3,
          curso_id: 1,
          indice: "3",
          titulo: 'Estado y Ciclo de Vida',
          descripcion: 'Gestión del estado interno y ciclo de vida de componentes.',
          video_url: 'https://filesamples.com/samples/video/mp4/sample_960x540.mp4',
          thumbnail_url: '/preview-3.jpg',
          duracion: 770, // en segundos (12:50 minutos)
          orden: 3
        },
        {
          id: 4,
          curso_id: 1,
          indice: "4",
          titulo: 'Hooks en React',
          descripcion: 'Uso de Hooks para gestionar estado y efectos secundarios.',
          video_url: 'https://filesamples.com/samples/video/mp4/sample_1280x720.mp4',
          thumbnail_url: '/preview-4.jpg',
          duracion: 930, // en segundos (15:30 minutos)
          orden: 4
        }
      ]);
    }
  }, [loading, capitulos.length, id]);

  // Formatear duración en formato mm:ss
  const formatDuracion = (segundos) => {
    const minutos = Math.floor(segundos / 60);
    const segs = segundos % 60;
    return `${minutos}:${segs < 10 ? '0' : ''}${segs}`;
  };

  // Manejar reordenamiento de capítulos
  const handleDragEnd = async (result) => {
    if (!result.destination) return;
    
    const items = Array.from(capitulos);
    const [reorderedItem] = items.splice(result.source.index, 1);
    items.splice(result.destination.index, 0, reorderedItem);
    
    // Actualizar órdenes
    const updatedItems = items.map((item, index) => ({
      ...item,
      orden: index + 1
    }));
    
    setCapitulos(updatedItems);
    
    // Guardar el nuevo orden en el servidor
    try {
      const token = localStorage.getItem('token');
      const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
      
      const ordenesData = updatedItems.map(cap => ({
        capitulo_id: cap.id,
        nuevo_orden: cap.orden
      }));
      
      const response = await fetch(`${apiUrl}/api/admin/cursos/${id}/capitulos/reordenar`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ ordenes: ordenesData })
      });

      if (!response.ok) {
        throw new Error('Error al reordenar capítulos');
      }
      
    } catch (err) {
      console.error('Error al guardar nuevo orden:', err);
      // Mostrar mensaje de error, pero mantener la UI actualizada
    }
  };

  // Abrir modal para crear un nuevo capítulo
  const handleCrearCapitulo = () => {
    setCurrentCapitulo({
      id: null,
      titulo: '',
      descripcion: '',
      video_url: '',
      thumbnail_url: ''
    });
    setVideoFile(null);
    setThumbnailFile(null);
    setVideoPreview('');
    setThumbnailPreview('');
    setUploadProgress(0);
    setFormErrors({});
    setModalMode('create');
    setShowModal(true);
  };

  // Abrir modal para editar un capítulo existente
  const handleEditarCapitulo = (capitulo) => {
    setCurrentCapitulo({
      id: capitulo.id,
      titulo: capitulo.titulo,
      descripcion: capitulo.descripcion,
      video_url: capitulo.video_url,
      thumbnail_url: capitulo.thumbnail_url
    });
    setVideoFile(null);
    setThumbnailFile(null);
    setVideoPreview(capitulo.video_url);
    setThumbnailPreview(capitulo.thumbnail_url);
    setUploadProgress(0);
    setFormErrors({});
    setModalMode('edit');
    setShowModal(true);
  };

  // Validar formulario del capítulo
  const validateCapituloForm = () => {
    const errors = {};
    
    if (!currentCapitulo.titulo.trim()) {
      errors.titulo = 'El título es obligatorio';
    }
    
    if (modalMode === 'create' && !videoFile) {
      errors.video = 'Debes subir un video para el capítulo';
    }
    
    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  // Manejar cambios en los campos del formulario
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setCurrentCapitulo({ ...currentCapitulo, [name]: value });
    
    // Limpiar error
    if (formErrors[name]) {
      setFormErrors({ ...formErrors, [name]: null });
    }
  };

  // Manejar cambio de archivo de video
  const handleVideoChange = (e) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      
      // Validar tipo de archivo
      if (!file.type.startsWith('video/')) {
        setFormErrors({ ...formErrors, video: 'Por favor, selecciona un archivo de video válido (MP4, WebM)' });
        return;
      }
      
      // Validar tamaño (máx. 500MB)
      if (file.size > 500 * 1024 * 1024) {
        setFormErrors({ ...formErrors, video: 'El video debe ser menor a 500MB' });
        return;
      }
      
      setVideoFile(file);
      setVideoPreview(URL.createObjectURL(file));
      
      // Limpiar error
      if (formErrors.video) {
        setFormErrors({ ...formErrors, video: null });
      }
      
      // Obtener duración del video
      const video = document.createElement('video');
      video.preload = 'metadata';
      video.onloadedmetadata = function() {
        // Actualizar duración en el estado
        setCurrentCapitulo(prev => ({ 
          ...prev, 
          duracion: Math.round(video.duration)
        }));
        
        window.URL.revokeObjectURL(video.src);
      };
      video.src = URL.createObjectURL(file);
    }
  };

  // Manejar cambio de archivo de thumbnail
  const handleThumbnailChange = (e) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      
      // Validar tipo de archivo
      if (!file.type.startsWith('image/')) {
        setFormErrors({ ...formErrors, thumbnail: 'Por favor, selecciona una imagen válida (JPG, PNG, WebP)' });
        return;
      }
      
      // Validar tamaño (máx. 5MB)
      if (file.size > 5 * 1024 * 1024) {
        setFormErrors({ ...formErrors, thumbnail: 'La imagen debe ser menor a 5MB' });
        return;
      }
      
      setThumbnailFile(file);
      setThumbnailPreview(URL.createObjectURL(file));
      
      // Limpiar error
      if (formErrors.thumbnail) {
        setFormErrors({ ...formErrors, thumbnail: null });
      }
    }
  };

  // Generar thumbnail automáticamente desde el video
  const generateThumbnail = () => {
    if (videoRef.current && videoPreview) {
      // Establecer un tiempo a la mitad del video
      videoRef.current.currentTime = videoRef.current.duration / 2;
      
      // Cuando el video se ubique en ese tiempo, capturar frame
      videoRef.current.onseeked = () => {
        const canvas = document.createElement('canvas');
        canvas.width = videoRef.current.videoWidth;
        canvas.height = videoRef.current.videoHeight;
        
        const ctx = canvas.getContext('2d');
        ctx.drawImage(videoRef.current, 0, 0, canvas.width, canvas.height);
        
        // Convertir a Blob y luego a File
        canvas.toBlob((blob) => {
          const file = new File([blob], 'thumbnail.jpg', { type: 'image/jpeg' });
          setThumbnailFile(file);
          setThumbnailPreview(URL.createObjectURL(blob));
        }, 'image/jpeg', 0.8);
      };
    }
  };

  // Guardar un nuevo capítulo o actualizar uno existente
  const handleSaveCapitulo = async () => {
    if (!validateCapituloForm()) {
      return;
    }
    
    setUploading(true);
    setUploadProgress(0);
    
    try {
      const token = localStorage.getItem('token');
      const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
      
      // Crear FormData para enviar archivos
      const formData = new FormData();
      formData.append('titulo', currentCapitulo.titulo);
      formData.append('descripcion', currentCapitulo.descripcion);
      
      if (videoFile) {
        formData.append('video', videoFile);
      } else if (modalMode === 'edit') {
        formData.append('video_url', currentCapitulo.video_url);
      }
      
      if (thumbnailFile) {
        formData.append('thumbnail', thumbnailFile);
      } else if (modalMode === 'edit' && currentCapitulo.thumbnail_url) {
        formData.append('thumbnail_url', currentCapitulo.thumbnail_url);
      }
      
      if (currentCapitulo.duracion) {
        formData.append('duracion', currentCapitulo.duracion);
      }
      
      // URL y método según sea creación o edición
      const url = modalMode === 'edit'
        ? `${apiUrl}/api/admin/cursos/${id}/capitulos/${currentCapitulo.id}`
        : `${apiUrl}/api/admin/cursos/${id}/capitulos`;
        
      const method = modalMode === 'edit' ? 'PUT' : 'POST';
      
      // Crear XMLHttpRequest para poder monitorizar el progreso
      const xhr = new XMLHttpRequest();
      
      xhr.upload.addEventListener('progress', (event) => {
        if (event.lengthComputable) {
          const progress = Math.round((event.loaded / event.total) * 100);
          setUploadProgress(progress);
        }
      });
      
      xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
          if (xhr.status >= 200 && xhr.status < 300) {
            // Éxito
            const response = JSON.parse(xhr.responseText);
            
            // Actualizar lista de capítulos
            if (modalMode === 'edit') {
              const updatedCapitulo = {
                ...currentCapitulo,
                video_url: response.video_url || currentCapitulo.video_url,
                thumbnail_url: response.thumbnail_url || currentCapitulo.thumbnail_url,
                duracion: currentCapitulo.duracion
              };
              
              setCapitulos(capitulos.map(cap => 
                cap.id === updatedCapitulo.id ? updatedCapitulo : cap
              ));
            } else {
              // Agregar el nuevo capítulo a la lista
              const nuevoCapitulo = {
                id: response.id,
                curso_id: parseInt(id),
                titulo: currentCapitulo.titulo,
                descripcion: currentCapitulo.descripcion,
                video_url: response.video_url,
                thumbnail_url: response.thumbnail_url,
                duracion: currentCapitulo.duracion,
                orden: capitulos.length + 1
              };
              
              setCapitulos([...capitulos, nuevoCapitulo]);
            }
            
            // Cerrar modal
            setShowModal(false);
            
            // Mostrar mensaje de éxito
            alert(modalMode === 'edit' 
              ? 'Capítulo actualizado correctamente' 
              : 'Capítulo creado correctamente');
          } else {
            // Error
            console.error('Error al guardar capítulo:', xhr.responseText);
            
            try {
              const errorData = JSON.parse(xhr.responseText);
              setError(errorData.error || 'Error al guardar el capítulo');
            } catch (e) {
              setError('Error al guardar el capítulo');
            }
          }
          
          setUploading(false);
        }
      };
      
      xhr.open(method, url, true);
      xhr.setRequestHeader('Authorization', `Bearer ${token}`);
      xhr.send(formData);
      
    } catch (err) {
      console.error('Error al guardar capítulo:', err);
      setError(`Error al ${modalMode === 'edit' ? 'actualizar' : 'crear'} el capítulo: ${err.message}`);
      setUploading(false);
    }
  };

  // Eliminar un capítulo
  const handleDeleteCapitulo = async (capituloId) => {
    if (window.confirm('¿Estás seguro de que deseas eliminar este capítulo? Esta acción no se puede deshacer.')) {
      try {
        const token = localStorage.getItem('token');
        const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
        
        const response = await fetch(`${apiUrl}/api/admin/cursos/${id}/capitulos/${capituloId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });

        if (!response.ok) {
          throw new Error('Error al eliminar el capítulo');
        }

        // Eliminar de la lista y reordenar los restantes
        setCapitulos(capitulos.filter(cap => cap.id !== capituloId));
        
        alert('Capítulo eliminado correctamente');
      } catch (err) {
        console.error('Error al eliminar capítulo:', err);
        alert('Error al eliminar el capítulo. Por favor, intenta de nuevo.');
      }
    }
  };

  if (loading) {
    return (
      <div className="admin-loading-container">
        <div className="spinner"></div>
        <p>Cargando información del curso...</p>
      </div>
    );
  }

  if (error && !curso) {
    return (
      <div className="admin-error-container">
        <h2>Error</h2>
        <p>{error}</p>
        <button className="admin-btn admin-btn-primary" onClick={() => window.location.reload()}>
          Intentar de nuevo
        </button>
      </div>
    );
  }

  if (!curso) {
    return (
      <div className="admin-error-container">
        <h2>Error</h2>
        <p>No se encontró el curso</p>
        <Link to="/admin/cursos" className="admin-btn admin-btn-primary">
          Volver a la lista de cursos
        </Link>
      </div>
    );
  }

  return (
    <div className="admin-capitulos-page">
      <div className="admin-page-title">
        <div className="admin-title-with-back">
          <Link to="/admin/cursos" className="admin-back-button">
            <FontAwesomeIcon icon={faArrowLeft} />
          </Link>
          <div>
            <h1>Capítulos de Curso</h1>
            <p>{curso.titulo}</p>
          </div>
        </div>
        
        <button 
          className="admin-btn admin-btn-primary"
          onClick={handleCrearCapitulo}
          disabled={capitulos.length >= 50}
        >
          <FontAwesomeIcon icon={faPlus} /> Nuevo Capítulo
        </button>
      </div>

      {capitulos.length === 0 ? (
        <div className="admin-empty-state">
          <div className="admin-empty-icon">
            <FontAwesomeIcon icon={faVideo} />
          </div>
          <h2>No hay capítulos</h2>
          <p>Comienza a crear contenido para tu curso</p>
          <button className="admin-btn admin-btn-primary" onClick={handleCrearCapitulo}>
            <FontAwesomeIcon icon={faPlus} /> Añadir Primer Capítulo
          </button>
        </div>
      ) : (
        <>
          <div className="admin-card">
            <div className="admin-card-header">
              <h2 className="admin-card-title">Capítulos ({capitulos.length}/50)</h2>
              <p className="admin-card-subtitle">
                Arrastra los capítulos para cambiar su orden
              </p>
            </div>
            
            <DragDropContext onDragEnd={handleDragEnd}>
              <Droppable droppableId="capitulos">
                {(provided) => (
                  <div
                    {...provided.droppableProps}
                    ref={provided.innerRef}
                    className="admin-capitulos-list"
                  >
                    {capitulos.map((capitulo, index) => (
                      <Draggable key={capitulo.id} draggableId={capitulo.id.toString()} index={index}>
                        {(provided) => (
                          <div
                            className="admin-capitulo-item"
                            ref={provided.innerRef}
                            {...provided.draggableProps}
                          >
                            <div className="admin-capitulo-drag-handle" {...provided.dragHandleProps}>
                              <FontAwesomeIcon icon={faGripLines} />
                            </div>
                            
                            <div className="admin-capitulo-index">
                              {index + 1}
                            </div>
                            
                            <div className="admin-capitulo-thumbnail">
                              {capitulo.thumbnail_url ? (
                                <img 
                                  src={capitulo.thumbnail_url} 
                                  alt={capitulo.titulo} 
                                  onError={(e) => {e.target.src = '/thumbnail-default.jpg'}} 
                                />
                              ) : (
                                <div className="admin-thumbnail-placeholder">
                                  <FontAwesomeIcon icon={faImage} />
                                </div>
                              )}
                            </div>
                            
                            <div className="admin-capitulo-content">
                              <h3 className="admin-capitulo-title">{capitulo.titulo}</h3>
                              <p className="admin-capitulo-description">{capitulo.descripcion}</p>
                            </div>
                            
                            <div className="admin-capitulo-duration">
                              {capitulo.duracion ? formatDuracion(capitulo.duracion) : '--:--'}
                            </div>
                            
                            <div className="admin-capitulo-actions">
                              <button 
                                className="admin-btn admin-btn-sm admin-btn-secondary"
                                title="Reproducir video"
                                onClick={() => window.open(capitulo.video_url, '_blank')}
                              >
                                <FontAwesomeIcon icon={faPlay} />
                              </button>
                              
                              <button 
                                className="admin-btn admin-btn-sm admin-btn-secondary"
                                title="Editar capítulo"
                                onClick={() => handleEditarCapitulo(capitulo)}
                              >
                                <FontAwesomeIcon icon={faEdit} />
                              </button>
                              
                              <button 
                                className="admin-btn admin-btn-sm admin-btn-danger"
                                title="Eliminar capítulo"
                                onClick={() => handleDeleteCapitulo(capitulo.id)}
                              >
                                <FontAwesomeIcon icon={faTrashAlt} />
                              </button>
                            </div>
                          </div>
                        )}
                      </Draggable>
                    ))}
                    {provided.placeholder}
                  </div>
                )}
              </Droppable>
            </DragDropContext>
          </div>
          
          <div className="admin-actions-bottom">
            <Link to={`/admin/cursos/editar/${id}`} className="admin-btn admin-btn-secondary">
              <FontAwesomeIcon icon={faEdit} /> Editar información del curso
            </Link>
            
            <Link to="/admin/cursos" className="admin-btn admin-btn-primary">
              <FontAwesomeIcon icon={faCheck} /> Finalizar
            </Link>
          </div>
        </>
      )}

      {/* Modal para crear/editar capítulo */}
      {showModal && (
        <div className="admin-modal-overlay">
          <div className="admin-modal">
            <div className="admin-modal-header">
              <h2>{modalMode === 'create' ? 'Nuevo Capítulo' : 'Editar Capítulo'}</h2>
              <button className="admin-modal-close" onClick={() => setShowModal(false)}>
                <FontAwesomeIcon icon={faTimes} />
              </button>
            </div>
            
            <div className="admin-modal-body">
              {error && (
                <div className="admin-error-container">
                  <p>{error}</p>
                </div>
              )}
              
              <div className="admin-form-group">
                <label className="admin-form-label" htmlFor="titulo">
                  Título del Capítulo *
                </label>
                <input
                  type="text"
                  id="titulo"
                  name="titulo"
                  className={`admin-form-control ${formErrors.titulo ? 'is-invalid' : ''}`}
                  value={currentCapitulo.titulo}
                  onChange={handleInputChange}
                  placeholder="Ej: Introducción a React"
                  required
                />
                {formErrors.titulo && <span className="admin-form-error">{formErrors.titulo}</span>}
              </div>
              
              <div className="admin-form-group">
                <label className="admin-form-label" htmlFor="descripcion">
                  Descripción breve
                </label>
                <textarea
                  id="descripcion"
                  name="descripcion"
                  className="admin-form-control admin-form-textarea"
                  value={currentCapitulo.descripcion}
                  onChange={handleInputChange}
                  placeholder="Descripción opcional del capítulo"
                  rows="3"
                />
              </div>
              
              <div className="admin-form-group">
                <label className="admin-form-label">
                  Video {modalMode === 'create' ? '*' : ''}
                </label>
                <div className={`admin-dropzone ${formErrors.video ? 'is-invalid' : ''}`}>
                  <input
                    type="file"
                    id="video"
                    accept="video/*"
                    onChange={handleVideoChange}
                    style={{ display: 'none' }}
                  />
                  <label htmlFor="video" className="admin-dropzone-inner">
                    {videoPreview ? (
                      <div className="admin-video-preview">
                        <video 
                          ref={videoRef}
                          src={videoPreview} 
                          controls 
                          className="admin-video-player"
                        />
                        <div className="admin-video-actions">
                          <label htmlFor="video" className="admin-btn admin-btn-sm admin-btn-secondary">
                            Cambiar video
                          </label>
                          {videoPreview && !thumbnailPreview && (
                            <button 
                              type="button" 
                              className="admin-btn admin-btn-sm admin-btn-secondary"
                              onClick={generateThumbnail}
                            >
                              Generar miniatura
                            </button>
                          )}
                        </div>
                      </div>
                    ) : (
                      <>
                        <div className="admin-dropzone-icon">
                          <FontAwesomeIcon icon={faVideo} />
                        </div>
                        <div className="admin-dropzone-text">
                          Haz clic para subir un video
                        </div>
                        <div className="admin-dropzone-help">
                          Formatos: MP4, WebM | Máx: 500MB
                        </div>
                      </>
                    )}
                  </label>
                </div>
                {formErrors.video && <span className="admin-form-error">{formErrors.video}</span>}
              </div>
              
              <div className="admin-form-group">
                <label className="admin-form-label">
                  Miniatura (Opcional)
                </label>
                <div className={`admin-dropzone ${formErrors.thumbnail ? 'is-invalid' : ''}`}>
                  <input
                    type="file"
                    id="thumbnail"
                    accept="image/*"
                    onChange={handleThumbnailChange}
                    style={{ display: 'none' }}
                  />
                  <label htmlFor="thumbnail" className="admin-dropzone-inner">
                    {thumbnailPreview ? (
                      <div className="admin-thumbnail-preview">
                        <img 
                          src={thumbnailPreview} 
                          alt="Vista previa" 
                          className="admin-thumbnail-image" 
                        />
                        <div className="admin-thumbnail-text">
                          Haz clic para cambiar la imagen
                        </div>
                      </div>
                    ) : (
                      <>
                        <div className="admin-dropzone-icon">
                          <FontAwesomeIcon icon={faImage} />
                        </div>
                        <div className="admin-dropzone-text">
                          Haz clic para subir una miniatura
                        </div>
                        <div className="admin-dropzone-help">
                          Formatos: JPG, PNG, WebP | Máx: 5MB
                        </div>
                      </>
                    )}
                  </label>
                </div>
                {formErrors.thumbnail && <span className="admin-form-error">{formErrors.thumbnail}</span>}
              </div>
              
              {uploading && (
                <div className="admin-progress-container">
                  <div className="admin-progress-bar">
                    <div 
                      className="admin-progress-fill" 
                      style={{ width: `${uploadProgress}%` }}
                    ></div>
                  </div>
                  <div className="admin-progress-text">
                    Subiendo... {uploadProgress}%
                  </div>
                </div>
              )}
            </div>
            
            <div className="admin-modal-footer">
              <button 
                type="button" 
                className="admin-btn admin-btn-secondary"
                onClick={() => setShowModal(false)}
                disabled={uploading}
              >
                <FontAwesomeIcon icon={faTimes} /> Cancelar
              </button>
              
              <button 
                type="button" 
                className="admin-btn admin-btn-primary"
                onClick={handleSaveCapitulo}
                disabled={uploading}
              >
                <FontAwesomeIcon icon={faSave} /> 
                {uploading 
                  ? 'Guardando...' 
                  : modalMode === 'create' 
                    ? 'Crear Capítulo' 
                    : 'Guardar Cambios'
                }
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default AdminCapitulos;