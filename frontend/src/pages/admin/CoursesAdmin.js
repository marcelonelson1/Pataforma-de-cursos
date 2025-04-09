import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import { 
  FaPlus, 
  FaEdit, 
  FaTrash, 
  FaUpload, 
  FaClock, 
  FaTimes, 
  FaEye,
  FaSearch,
  FaSave,
  FaVideo,
  FaExclamationTriangle
} from 'react-icons/fa';
import './CoursesAdmin.css';

const CoursesAdmin = () => {
  const [courses, setCourses] = useState([]);
  const [editingCourse, setEditingCourse] = useState(null);
  const [editingChapter, setEditingChapter] = useState(null);
  const [isFormVisible, setIsFormVisible] = useState(false);
  const [isChapterModalVisible, setIsChapterModalVisible] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [isLoading, setIsLoading] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [videoError, setVideoError] = useState('');
  const [chapterForm, setChapterForm] = useState({
    titulo: '',
    descripcion: '',
    duracion: '00:00',
    video_url: '',
    publicado: false,
    orden: 0,
    curso_id: null,
    videoFile: null,
    videoPreview: false
  });

  const getAuthHeader = useCallback(() => {
    const token = localStorage.getItem('token');
    return {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    };
  }, []);

  const fetchCourses = useCallback(async () => {
    try {
      setIsLoading(true);
      const response = await axios.get('/api/cursos', getAuthHeader());
      
      // Para cada curso, aseguramos que tenga un array de capítulos
      const coursesWithChapters = response.data.map(course => ({
        ...course,
        capitulos: Array.isArray(course.capitulos) ? course.capitulos : []
      }));
      
      setCourses(coursesWithChapters);
    } catch (error) {
      console.error('Error al obtener cursos:', error);
      alert('Error al cargar los cursos');
    } finally {
      setIsLoading(false);
    }
  }, [getAuthHeader]);

  useEffect(() => {
    fetchCourses();
  }, [fetchCourses]);

  const filteredCourses = courses.filter(course => {
    const matchesSearch = course.titulo.toLowerCase().includes(searchTerm.toLowerCase()) || 
                         course.descripcion.toLowerCase().includes(searchTerm.toLowerCase());
    
    const matchesStatus = statusFilter === 'all' || 
                          (statusFilter === 'published' && course.estado === 'Publicado') ||
                          (statusFilter === 'draft' && course.estado === 'Borrador');
    
    return matchesSearch && matchesStatus;
  });

  const handleCourseSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const formData = new FormData(e.target);
      const courseData = {
        titulo: formData.get('title'),
        descripcion: formData.get('description'),
        contenido: formData.get('description'),
        precio: parseFloat(formData.get('price')),
        estado: formData.get('status'),
        imagen_url: editingCourse?.imagen_url || ''
      };

      let response;
      if (editingCourse?.id) {
        response = await axios.put(
          `/api/cursos/${editingCourse.id}`,
          courseData,
          getAuthHeader()
        );
        
        // Aseguramos que los capítulos se mantengan después de la actualización
        const updatedCourse = {
          ...response.data,
          capitulos: editingCourse.capitulos || []
        };
        
        setCourses(courses.map(c => c.id === editingCourse.id ? updatedCourse : c));
        setEditingCourse(updatedCourse); // Actualizamos el curso en edición
      } else {
        response = await axios.post(
          '/api/cursos',
          courseData,
          getAuthHeader()
        );
        
        const newCourse = {
          ...response.data,
          capitulos: []
        };
        
        setCourses([...courses, newCourse]);
        setEditingCourse(newCourse); // Actualizamos para poder añadir capítulos inmediatamente
      }

      setIsFormVisible(false);
      
      // No resetear el curso en edición para poder seguir trabajando con él
      // setEditingCourse(null);
    } catch (error) {
      console.error('Error al guardar curso:', error);
      alert('Error al guardar el curso');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeleteCourse = async (id) => {
    if (window.confirm('¿Estás seguro de que quieres eliminar este curso?')) {
      try {
        setIsLoading(true);
        await axios.delete(`/api/cursos/${id}`, getAuthHeader());
        setCourses(courses.filter(course => course.id !== id));
      } catch (error) {
        console.error('Error al eliminar curso:', error);
        alert('Error al eliminar el curso');
      } finally {
        setIsLoading(false);
      }
    }
  };

  const handleImageUpload = (e) => {
    const file = e.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (event) => {
        setEditingCourse(prev => ({
          ...prev,
          imagen_url: event.target.result
        }));
      };
      reader.readAsDataURL(file);
    }
  };

  const handleNewCourse = () => {
    setEditingCourse({ 
      titulo: '',
      descripcion: '',
      contenido: '',
      precio: 0,
      estado: 'Borrador',
      imagen_url: '',
      capitulos: []
    });
    setIsFormVisible(true);
  };

  const handleEditCourse = (course) => {
    // Hacemos una copia profunda para evitar problemas de referencia
    setEditingCourse({
      ...course,
      capitulos: Array.isArray(course.capitulos) ? [...course.capitulos] : []
    });
    setIsFormVisible(true);
  };

  const handleAddChapter = () => {
    if (!editingCourse || !editingCourse.id) {
      alert('Por favor, guarde el curso primero antes de añadir capítulos.');
      return;
    }
    
    setChapterForm({
      titulo: '',
      descripcion: '',
      duracion: '00:00',
      video_url: '',
      publicado: false,
      orden: editingCourse.capitulos?.length || 0,
      curso_id: editingCourse.id,
      videoFile: null,
      videoPreview: false
    });
    setVideoError('');
    setIsChapterModalVisible(true);
    setEditingChapter(null);
  };

  const handleEditChapter = (chapter) => {
    setChapterForm({
      titulo: chapter.titulo,
      descripcion: chapter.descripcion,
      duracion: chapter.duracion,
      video_url: chapter.video_url,
      publicado: chapter.publicado,
      orden: chapter.orden,
      curso_id: editingCourse.id,
      videoFile: null,
      videoPreview: false
    });
    setVideoError('');
    setIsChapterModalVisible(true);
    setEditingChapter(chapter);
  };

  const handleChapterSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setUploadProgress(0);

    try {
      // Verificar si hay un nuevo archivo de video para subir
      let videoUrl = chapterForm.video_url;
      
      if (chapterForm.videoFile && chapterForm.videoPreview) {
        setIsUploading(true);
        
        // Crear FormData para subir el archivo
        const formData = new FormData();
        formData.append('video', chapterForm.videoFile);
        formData.append('curso_id', editingCourse.id);
        formData.append('capitulo_id', editingChapter?.id || 'nuevo');
        
        // Realizar la subida del archivo
        const uploadResponse = await axios.post(
          '/api/videos/upload',
          formData,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
              ...getAuthHeader().headers
            },
            onUploadProgress: (progressEvent) => {
              const percentCompleted = Math.round(
                (progressEvent.loaded * 100) / progressEvent.total
              );
              setUploadProgress(percentCompleted);
            }
          }
        );
        
        // Obtener la URL del video subido
        videoUrl = uploadResponse.data.video_url;
        setIsUploading(false);
      }
      
      // Datos del capítulo para guardar en la base de datos
      const chapterData = {
        ...chapterForm,
        video_url: videoUrl, // Usar la URL del servidor si se subió un nuevo video
        curso_id: parseInt(editingCourse.id, 10)
      };
      
      // Eliminar propiedades auxiliares
      delete chapterData.videoFile;
      delete chapterData.videoPreview;

      console.log('Enviando datos de capítulo:', chapterData);

      let response;
      if (editingChapter?.id) {
        response = await axios.put(
          `/api/capitulos/${editingChapter.id}`,
          chapterData,
          getAuthHeader()
        );
        
        console.log('Respuesta de actualización de capítulo:', response.data);
      } else {
        response = await axios.post(
          '/api/capitulos',
          chapterData,
          getAuthHeader()
        );
        
        console.log('Respuesta de creación de capítulo:', response.data);
      }

      // Actualizamos los capítulos del curso en edición
      const updatedChapters = editingChapter?.id
        ? editingCourse.capitulos.map(ch => 
            ch.id === editingChapter.id ? response.data : ch
          )
        : [...(editingCourse.capitulos || []), response.data];

      // Actualizamos el curso en edición con los nuevos capítulos
      const updatedEditingCourse = {
        ...editingCourse,
        capitulos: updatedChapters
      };
      
      setEditingCourse(updatedEditingCourse);

      // Actualizamos la lista de cursos con el curso actualizado
      setCourses(courses.map(c => 
        c.id === editingCourse.id 
          ? updatedEditingCourse
          : c
      ));

      // Aseguramos que los cambios se vean reflejados en la UI
      setIsChapterModalVisible(false);
      setEditingChapter(null);
      
      // Recargamos los cursos para asegurar que tenemos los datos actualizados
      fetchCourses();
    } catch (error) {
      console.error('Error al guardar capítulo:', error);
      const errorMessage = error.response?.data?.error || error.message;
      alert('Error al guardar el capítulo: ' + errorMessage);
    } finally {
      setIsLoading(false);
      setIsUploading(false);
    }
  };

  const handleChapterFormChange = (e) => {
    const { name, value, type, checked } = e.target;
    setChapterForm(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleVideoUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;
    
    setVideoError('');

    // Validar el tipo de archivo
    const validTypes = ['video/mp4', 'video/webm', 'video/ogg'];
    if (!validTypes.includes(file.type)) {
      setVideoError('Por favor selecciona un archivo de video válido (MP4, WebM, OGG)');
      return;
    }

    // Validar el tamaño del archivo (máximo 100MB)
    const maxSize = 100 * 1024 * 1024; // 100MB en bytes
    if (file.size > maxSize) {
      setVideoError('El archivo es demasiado grande. Por favor sube un video de menos de 100MB.');
      return;
    }

    try {
      // Crear una vista previa local temporal
      const videoUrl = URL.createObjectURL(file);
      setChapterForm(prev => ({
        ...prev,
        video_url: videoUrl,
        videoFile: file, // Guardar el archivo para subirlo cuando se envíe el formulario
        videoPreview: true // Indicar que es una vista previa
      }));

      // Extraer duración del video para autocompletar el campo
      const video = document.createElement('video');
      video.preload = 'metadata';
      video.onloadedmetadata = function() {
        window.URL.revokeObjectURL(video.src);
        const duration = video.duration;
        const minutes = Math.floor(duration / 60);
        const seconds = Math.floor(duration % 60);
        const formattedDuration = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
        
        setChapterForm(prev => ({
          ...prev,
          duracion: formattedDuration
        }));
      };
      video.src = videoUrl;
      
    } catch (error) {
      console.error('Error al procesar el video:', error);
      setVideoError('Error al preparar el video para la subida');
    }
  };

  const handleDeleteChapter = async (chapterId) => {
    if (window.confirm('¿Estás seguro de que quieres eliminar este capítulo?')) {
      try {
        setIsLoading(true);
        await axios.delete(`/api/capitulos/${chapterId}`, getAuthHeader());
        
        const updatedChapters = editingCourse.capitulos.filter(ch => ch.id !== chapterId);
        
        const updatedEditingCourse = {
          ...editingCourse,
          capitulos: updatedChapters
        };
        
        setEditingCourse(updatedEditingCourse);
        
        setCourses(courses.map(c => 
          c.id === editingCourse.id 
            ? updatedEditingCourse
            : c
        ));
        
        // Recargamos los cursos para asegurar que tenemos los datos actualizados
        fetchCourses();
      } catch (error) {
        console.error('Error al eliminar capítulo:', error);
        alert('Error al eliminar el capítulo');
      } finally {
        setIsLoading(false);
      }
    }
  };

  const closeForm = () => {
    if (window.confirm('¿Estás seguro? Los cambios no guardados se perderán.')) {
      setIsFormVisible(false);
      setEditingCourse(null);
    }
  };

  const closeChapterModal = () => {
    if (window.confirm('¿Estás seguro? Los cambios no guardados se perderán.')) {
      setIsChapterModalVisible(false);
      setEditingChapter(null);
      setVideoError('');
      
      // Revocar URL de objeto si es una vista previa
      if (chapterForm.videoPreview && chapterForm.video_url) {
        URL.revokeObjectURL(chapterForm.video_url);
      }
    }
  };

  const handlePreviewCourse = (course) => {
    console.log('Previsualizando curso:', course.titulo);
  };

  const LoadingOverlay = () => (
    <div className="loading-overlay">
      <div className="spinner"></div>
    </div>
  );

  return (
    <div className="courses-admin">
      {isLoading && <LoadingOverlay />}

      <div className="section-header">
        <h2 className="section-title">Gestión de Cursos</h2>
        <div className="filters-section">
          <div className="search-bar course-search">
            <FaSearch />
            <input 
              type="text" 
              placeholder="Buscar cursos..." 
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>
          
          <div className="status-filters">
            <button 
              className={`filter-btn ${statusFilter === 'all' ? 'active' : ''}`}
              onClick={() => setStatusFilter('all')}
            >
              Todos
            </button>
            <button 
              className={`filter-btn ${statusFilter === 'published' ? 'active' : ''}`}
              onClick={() => setStatusFilter('published')}
            >
              Publicados
            </button>
            <button 
              className={`filter-btn ${statusFilter === 'draft' ? 'active' : ''}`}
              onClick={() => setStatusFilter('draft')}
            >
              Borradores
            </button>
          </div>
          
          <button onClick={handleNewCourse} className="action-button">
            <FaPlus />
            Nuevo Curso
          </button>
        </div>
      </div>

      <div className="courses-list">
        {filteredCourses.length === 0 ? (
          <div className="no-results">
            <p>No se encontraron cursos que coincidan con tu búsqueda.</p>
          </div>
        ) : (
          filteredCourses.map(course => (
            <div key={course.id} className="course-card">
              <div className="course-header">
                <h3>{course.titulo}</h3>
                <span className="price">${course.precio.toFixed(2)}</span>
              </div>
              
              <div className="course-meta">
                <span className={`status-badge ${course.estado === 'Publicado' ? 'published' : 'draft'}`}>
                  {course.estado}
                </span>
                {course.created_at && (
                  <span className="date-published">
                    Creado: {new Date(course.created_at).toLocaleDateString()}
                  </span>
                )}
              </div>
              
              <p className="description">{course.descripcion}</p>
              
              <div className="course-stats">
                <div className="stat">
                  <span className="stat-value">{course.capitulos?.length || 0}</span>
                  <span className="stat-label">Capítulos</span>
                </div>
              </div>
              
              <div className="chapters-summary">
                <h4>Capítulos <span>{course.capitulos?.length || 0}</span></h4>
                {course.capitulos && course.capitulos.length > 0 ? (
                  <ul>
                    {course.capitulos.slice(0, 3).map((chapter, index) => (
                      <li key={chapter.id || index}>
                        <span>{chapter.titulo}</span>
                        <div className="chapter-meta">
                          <span className="chapter-duration">
                            <FaClock /> {chapter.duracion || '00:00'}
                          </span>
                          {chapter.video_url && (
                            <span className="chapter-video">
                              <FaVideo /> Video
                            </span>
                          )}
                        </div>
                      </li>
                    ))}
                    {course.capitulos.length > 3 && (
                      <li className="more-chapters">
                        +{course.capitulos.length - 3} capítulos más
                      </li>
                    )}
                  </ul>
                ) : (
                  <p className="no-chapters">No hay capítulos creados aún.</p>
                )}
              </div>
              
              <div className="course-actions">
                <button onClick={() => handlePreviewCourse(course)} className="btn-preview">
                  <FaEye />
                  Previsualizar
                </button>
                <button onClick={() => handleEditCourse(course)} className="btn-edit">
                  <FaEdit />
                  Editar
                </button>
                <button onClick={() => handleDeleteCourse(course.id)} className="btn-delete">
                  <FaTrash />
                  Eliminar
                </button>
              </div>
            </div>
          ))
        )}
      </div>

      {isFormVisible && (
        <div className="form-overlay">
          <div className="form-container">
            <div className="form-header">
              <h3>{editingCourse?.id ? 'Editar Curso' : 'Nuevo Curso'}</h3>
              <button className="form-close" onClick={closeForm}>
                <FaTimes />
              </button>
            </div>
            
            <form onSubmit={handleCourseSubmit}>
              <div className="form-group">
                <label htmlFor="title">Título del Curso</label>
                <input
                  id="title"
                  name="title"
                  type="text"
                  className="form-control"
                  defaultValue={editingCourse?.titulo || ''}
                  placeholder="Introduce el título del curso"
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="description">Descripción del Curso</label>
                <textarea
                  id="description"
                  name="description"
                  className="form-control"
                  defaultValue={editingCourse?.descripcion || ''}
                  placeholder="Describe el contenido y objetivos del curso"
                  required
                />
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label htmlFor="price">Precio ($)</label>
                  <input
                    id="price"
                    name="price"
                    type="number"
                    step="0.01"
                    className="form-control"
                    defaultValue={editingCourse?.precio || 0}
                    min="0"
                    required
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="status">Estado</label>
                  <select
                    id="status"
                    name="status"
                    className="form-control"
                    defaultValue={editingCourse?.estado || 'Borrador'}
                  >
                    <option value="Borrador">Borrador</option>
                    <option value="Publicado">Publicado</option>
                  </select>
                </div>
              </div>

              <div className="form-group">
                <label htmlFor="courseImage">Imagen del Curso</label>
                <div className="file-upload">
                  <input
                    id="courseImage"
                    type="file"
                    accept="image/*"
                    onChange={handleImageUpload}
                    className="form-control-file"
                  />
                  <label htmlFor="courseImage" className="upload-label">
                    <FaUpload /> Subir Imagen
                  </label>
                </div>
                {editingCourse?.imagen_url && (
                  <div className="image-preview">
                    <img src={editingCourse.imagen_url} alt="Vista previa del curso" />
                  </div>
                )}
              </div>

              <div className="form-actions">
                <button type="button" className="btn-cancel" onClick={closeForm}>
                  Cancelar
                </button>
                <button type="submit" className="btn-primary">
                  {editingCourse?.id ? 'Actualizar' : 'Crear'} Curso
                </button>
              </div>
            </form>

            {editingCourse?.id && (
              <div className="chapters-section">
                <div className="section-header">
                  <h3>Capítulos</h3>
                  <button
                    type="button"
                    onClick={handleAddChapter}
                    className="action-button small"
                  >
                    <FaPlus />
                    Añadir Capítulo
                  </button>
                </div>
                
                {editingCourse?.capitulos && editingCourse.capitulos.length > 0 ? (
                  <div className="chapters-list">
                    {editingCourse.capitulos.map((chapter, index) => (
                      <div key={chapter.id || index} className="chapter-preview">
                        <div className="chapter-info">
                          <div className="chapter-number">{index + 1}</div>
                          <div>
                            <h4>{chapter.titulo}</h4>
                            <p>{chapter.descripcion?.substring(0, 60)}{chapter.descripcion?.length > 60 ? '...' : ''}</p>
                            <div className="chapter-meta">
                              <span className="chapter-duration">
                                <FaClock /> {chapter.duracion || '00:00'}
                              </span>
                              <span className={`chapter-status ${chapter.publicado ? 'published' : 'draft'}`}>
                                {chapter.publicado ? 'Publicado' : 'Borrador'}
                              </span>
                              {chapter.video_url && (
                                <span className="chapter-video">
                                  <FaVideo /> Video disponible
                                </span>
                              )}
                            </div>
                          </div>
                        </div>
                        <div className="chapter-actions">
                          <button
                            type="button"
                            onClick={() => handleEditChapter(chapter)}
                            className="btn-action"
                          >
                            <FaEdit />
                          </button>
                          <button
                            type="button"
                            onClick={() => handleDeleteChapter(chapter.id)}
                            className="btn-action danger"
                          >
                            <FaTrash />
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <button
                    type="button"
                    onClick={handleAddChapter}
                    className="add-chapter-btn"
                  >
                    <FaPlus /> Añadir tu primer capítulo
                  </button>
                )}
              </div>
            )}
          </div>
        </div>
      )}

      {isChapterModalVisible && (
        <div className="chapter-modal">
          <div className="chapter-modal-overlay" onClick={closeChapterModal} />
          <div className="chapter-modal-content">
            <div className="modal-header">
              <h3>{editingChapter ? 'Editar Capítulo' : 'Nuevo Capítulo'}</h3>
              <button className="close-button" onClick={closeChapterModal}>
                <FaTimes />
              </button>
            </div>
            
            <form onSubmit={handleChapterSubmit}>
              <input
                type="hidden"
                name="curso_id"
                value={editingCourse?.id || ''}
              />
              
              <div className="form-group">
                <label>Título del Capítulo</label>
                <input
                  type="text"
                  name="titulo"
                  value={chapterForm.titulo}
                  onChange={handleChapterFormChange}
                  placeholder="Introduce el título del capítulo"
                  required
                />
              </div>

              <div className="form-group">
                <label>Descripción</label>
                <textarea
                  name="descripcion"
                  value={chapterForm.descripcion}
                  onChange={handleChapterFormChange}
                  placeholder="Describe el contenido de este capítulo"
                  required
                />
              </div>

              <div className="form-group">
                <label>Video del Capítulo</label>
                <div className="file-upload-wrapper">
                  <input
                    type="file"
                    id="chapterVideo"
                    accept="video/*"
                    onChange={handleVideoUpload}
                    className="form-control-file"
                  />
                  <label htmlFor="chapterVideo" className="upload-label">
                    <FaUpload /> {chapterForm.video_url && !chapterForm.videoPreview ? 'Cambiar Video' : 'Subir Video'}
                  </label>
                  
                  {videoError && (
                    <div className="error-message">
                      <FaExclamationTriangle /> {videoError}
                    </div>
                  )}
                  
                  {chapterForm.video_url && (
                    <div className="video-preview">
                      <video controls>
                        <source src={chapterForm.video_url} type="video/mp4" />
                        Tu navegador no soporta el elemento de video.
                      </video>
                      {chapterForm.videoPreview && (
                        <div className="preview-badge">Vista previa - No guardado</div>
                      )}
                    </div>
                  )}
                  
                  {isUploading && (
                    <div className="upload-progress">
                      <div 
                        className="progress-bar" 
                        style={{ width: `${uploadProgress}%` }}
                      ></div>
                      <span className="progress-text">{uploadProgress}% Completado</span>
                    </div>
                  )}
                </div>
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label>Duración (MM:SS)</label>
                  <input
                    type="text"
                    name="duracion"
                    value={chapterForm.duracion}
                    onChange={handleChapterFormChange}
                    placeholder="Ej: 15:30"
                    pattern="[0-9]{2}:[0-9]{2}"
                    required
                  />
                </div>

                <div className="form-group checkbox-group">
                  <label>
                    <input
                      type="checkbox"
                      name="publicado"
                      checked={chapterForm.publicado}
                      onChange={handleChapterFormChange}
                    />
                    <span className="checkmark"></span>
                    Publicado
                  </label>
                </div>
              </div>

              <div className="form-actions">
                <button
                  type="button"
                  className="btn-cancel"
                  onClick={closeChapterModal}
                >
                  Cancelar
                </button>
                <button 
                  type="submit" 
                  className="btn-primary"
                  disabled={isLoading || isUploading}
                >
                  <FaSave /> {isLoading || isUploading ? 'Guardando...' : 'Guardar Capítulo'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default CoursesAdmin;