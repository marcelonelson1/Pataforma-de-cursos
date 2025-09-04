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
  const [categories, setCategories] = useState([]);
  const [editingCourse, setEditingCourse] = useState(null);
  const [editingChapter, setEditingChapter] = useState(null);
  const [isFormVisible, setIsFormVisible] = useState(false);
  const [isChapterModalVisible, setIsChapterModalVisible] = useState(false);
  const [isCategoryModalVisible, setIsCategoryModalVisible] = useState(false);
  const [newCategoryName, setNewCategoryName] = useState('');
  const [newCategoryDescription, setNewCategoryDescription] = useState('');
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [isLoading, setIsLoading] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [videoError, setVideoError] = useState('');
  const [imageFile, setImageFile] = useState(null);
  const [imageError, setImageError] = useState('');
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

  // Obtener URL base de la API
  const getApiBaseUrl = () => {
    if (process.env.REACT_APP_COURSE_API_URL) {
      return process.env.REACT_APP_COURSE_API_URL;
    }
    
    const currentUrl = window.location.origin;
    
    if (currentUrl.includes('localhost')) {
      return 'http://localhost:8003';
    }
    
    return currentUrl;
  };

  // Función para corregir URLs de imágenes
  const getImageUrl = (url) => {
    if (!url) return null;
    
    // Si ya es una URL completa o una data URL, devolverla como está
    if (url.startsWith('http') || url.startsWith('data:')) {
      return url;
    }
    
    // Si es una ruta relativa que empieza con /static, añadir la URL base
    if (url.startsWith('/static/')) {
      const baseUrl = getApiBaseUrl();
      return `${baseUrl}${url}`;
    }
    
    // Por defecto, devolver la URL sin cambios
    return url;
  };

  const getAuthHeader = useCallback(() => {
    const token = localStorage.getItem('token');
    return {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'multipart/form-data'
      }
    };
  }, []);

  const fetchCourses = useCallback(async () => {
    try {
      setIsLoading(true);
      const response = await axios.get(`${getApiBaseUrl()}/api/admin/courses`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      // Procesar los datos para asegurar que las URLs de imágenes son correctas
      const coursesWithChapters = response.data.data.courses.map(course => ({
        ...course,
        imagen_url: getImageUrl(course.imagen_url),
        capitulos: Array.isArray(course.capitulos) ? course.capitulos : []
      }));
      
      console.log('Cursos cargados con éxito:', coursesWithChapters);
      setCourses(coursesWithChapters);
    } catch (error) {
      console.error('Error al obtener cursos:', error);
      alert('Error al cargar los cursos');
    } finally {
      setIsLoading(false);
    }
  }, []);

  const fetchCategories = useCallback(async () => {
    try {
      const response = await axios.get(`${getApiBaseUrl()}/api/categories`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      if (response.data.success) {
        setCategories(response.data.data.categories || []);
      }
    } catch (error) {
      console.error('Error al obtener categorías:', error);
    }
  }, []);

  useEffect(() => {
    fetchCourses();
    fetchCategories();
  }, [fetchCourses, fetchCategories]);

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
    setImageError('');

    try {
      const formData = new FormData();
      
      // Agregar datos del curso
      formData.append('titulo', e.target.title.value);
      formData.append('descripcion', e.target.description.value);
      formData.append('contenido', e.target.description.value);
      formData.append('precio', e.target.price.value);
      formData.append('estado', e.target.status.value);
      
      // Agregar categoría si se seleccionó
      if (e.target.category && e.target.category.value) {
        formData.append('categoria_id', e.target.category.value);
      }
      
      // Agregar archivo de imagen si existe
      if (imageFile) {
        formData.append('imagen', imageFile);
        console.log('Archivo de imagen añadido al FormData:', imageFile.name);
      } else if (editingCourse?.imagen_url && !editingCourse.imagen_url.startsWith('data:')) {
        // Si hay una URL de imagen existente (no es una vista previa local)
        const imagenUrl = editingCourse.imagen_url;
        // Si la URL comienza con la URL base, extraer solo la parte relativa
        const baseUrl = getApiBaseUrl();
        const relativePath = imagenUrl.startsWith(baseUrl) 
          ? imagenUrl.substring(baseUrl.length) 
          : imagenUrl;
        
        formData.append('imagen_url', relativePath);
        console.log('URL de imagen existente añadida al FormData:', relativePath);
      }

      // Imprimir el contenido del FormData para depuración
      for (let pair of formData.entries()) {
        console.log(pair[0] + ': ' + pair[1]);
      }

      let response;
      const token = localStorage.getItem('token');
      const authHeader = {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'multipart/form-data'
        }
      };

      if (editingCourse?.id) {
        console.log(`Actualizando curso ID ${editingCourse.id}`);
        response = await axios.put(
          `${getApiBaseUrl()}/api/courses/${editingCourse.id}`,
          formData,
          authHeader
        );
        
        const updatedCourse = {
          ...response.data.data,
          imagen_url: getImageUrl(response.data.data.imagen_url),
          capitulos: editingCourse.capitulos || []
        };
        
        console.log('Curso actualizado:', updatedCourse);
        setCourses(courses.map(c => c.id === editingCourse.id ? updatedCourse : c));
        setEditingCourse(updatedCourse);
      } else {
        console.log('Creando nuevo curso');
        response = await axios.post(
          `${getApiBaseUrl()}/api/courses`,
          formData,
          authHeader
        );
        
        const newCourse = {
          ...response.data.data,
          imagen_url: getImageUrl(response.data.data.imagen_url),
          capitulos: response.data.data.capitulos || []
        };
        
        console.log('Nuevo curso creado:', newCourse);
        setCourses([...courses, newCourse]);
        setEditingCourse(newCourse);
      }

      setIsFormVisible(false);
      setImageFile(null);
    } catch (error) {
      console.error('Error al guardar curso:', error);
      setImageError('Error al subir la imagen. Por favor, intenta con una imagen más pequeña o en formato JPG/PNG.');
      alert('Error al guardar el curso: ' + (error.response?.data?.error || error.message));
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeleteCourse = async (id) => {
    if (window.confirm('¿Estás seguro de que quieres eliminar este curso?')) {
      try {
        setIsLoading(true);
        await axios.delete(`${getApiBaseUrl()}/api/courses/${id}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
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
    if (!file) return;

    setImageError('');
    
    // Validar tipo de archivo
    const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
    if (!validTypes.includes(file.type)) {
      setImageError('Por favor selecciona una imagen válida (JPEG, PNG, GIF o WEBP)');
      return;
    }

    // Validar tamaño del archivo (máximo 5MB)
    const maxSize = 5 * 1024 * 1024; // 5MB
    if (file.size > maxSize) {
      setImageError('La imagen es demasiado grande. Por favor sube una imagen de menos de 5MB.');
      return;
    }

    setImageFile(file);
    console.log('Archivo de imagen seleccionado:', file.name);
    
    // Crear vista previa local
    const reader = new FileReader();
    reader.onload = (event) => {
      setEditingCourse(prev => ({
        ...prev,
        imagen_url: event.target.result
      }));
    };
    reader.readAsDataURL(file);
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
    setImageFile(null);
    setImageError('');
    setIsFormVisible(true);
  };

  const handleEditCourse = (course) => {
    setEditingCourse({
      ...course,
      capitulos: Array.isArray(course.capitulos) ? [...course.capitulos] : []
    });
    setImageFile(null);
    setImageError('');
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
      let videoUrl = chapterForm.video_url;
      
      if (chapterForm.videoFile && chapterForm.videoPreview) {
        setIsUploading(true);
        
        const formData = new FormData();
        formData.append('video', chapterForm.videoFile);
        formData.append('curso_id', editingCourse.id);
        formData.append('capitulo_id', editingChapter?.id || 'nuevo');
        
        const uploadResponse = await axios.post(
          `${getApiBaseUrl()}/api/files/upload-video`,
          formData,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            onUploadProgress: (progressEvent) => {
              const percentCompleted = Math.round(
                (progressEvent.loaded * 100) / progressEvent.total
              );
              setUploadProgress(percentCompleted);
            }
          }
        );
        
        console.log('Respuesta de subida de video:', uploadResponse.data);
        videoUrl = uploadResponse.data.data.video_url;
        console.log('Video URL extraída:', videoUrl);
        setIsUploading(false);
      }
      
      const chapterData = {
        ...chapterForm,
        video_url: videoUrl,
        curso_id: parseInt(editingCourse.id, 10)
      };
      
      delete chapterData.videoFile;
      delete chapterData.videoPreview;

      let response;
      if (editingChapter?.id) {
        response = await axios.put(
          `${getApiBaseUrl()}/api/chapters/${editingChapter.id}`,
          chapterData,
          {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`,
              'Content-Type': 'application/json'
            }
          }
        );
      } else {
        response = await axios.post(
          `${getApiBaseUrl()}/api/chapters`,
          chapterData,
          {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`,
              'Content-Type': 'application/json'
            }
          }
        );
      }

      const updatedChapters = editingChapter?.id
        ? editingCourse.capitulos.map(ch => 
            ch.id === editingChapter.id ? response.data.data : ch
          )
        : [...(editingCourse.capitulos || []), response.data.data];

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

      setIsChapterModalVisible(false);
      setEditingChapter(null);
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

    const validTypes = ['video/mp4', 'video/webm', 'video/ogg'];
    if (!validTypes.includes(file.type)) {
      setVideoError('Por favor selecciona un archivo de video válido (MP4, WebM, OGG)');
      return;
    }

    const maxSize = 100 * 1024 * 1024;
    if (file.size > maxSize) {
      setVideoError('El archivo es demasiado grande. Por favor sube un video de menos de 100MB.');
      return;
    }

    try {
      const videoUrl = URL.createObjectURL(file);
      setChapterForm(prev => ({
        ...prev,
        video_url: videoUrl,
        videoFile: file,
        videoPreview: true
      }));

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
        await axios.delete(`${getApiBaseUrl()}/api/chapters/${chapterId}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        
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
      setImageFile(null);
      setImageError('');
    }
  };

  const closeChapterModal = () => {
    if (window.confirm('¿Estás seguro? Los cambios no guardados se perderán.')) {
      setIsChapterModalVisible(false);
      setEditingChapter(null);
      setVideoError('');
      
      if (chapterForm.videoPreview && chapterForm.video_url) {
        URL.revokeObjectURL(chapterForm.video_url);
      }
    }
  };

  const handlePreviewCourse = (course) => {
    console.log('Previsualizando curso:', course.titulo);
    window.open(`/curso/${course.id}`, '_blank');
  };

  // Funciones para manejar categorías
  const handleCreateCategory = async (e) => {
    e.preventDefault();
    if (!newCategoryName.trim()) return;

    try {
      setIsLoading(true);
      const response = await axios.post(`${getApiBaseUrl()}/api/categories`, {
        nombre: newCategoryName,
        descripcion: newCategoryDescription
      }, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        }
      });

      if (response.data.success) {
        await fetchCategories(); // Recargar categorías
        setNewCategoryName('');
        setNewCategoryDescription('');
        setIsCategoryModalVisible(false);
        alert('Categoría creada exitosamente');
      }
    } catch (error) {
      console.error('Error al crear categoría:', error);
      alert('Error al crear la categoría');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeleteCategory = async (categoryId) => {
    if (!window.confirm('¿Estás seguro de que quieres eliminar esta categoría?')) return;

    try {
      setIsLoading(true);
      await axios.delete(`${getApiBaseUrl()}/api/categories/${categoryId}`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });

      await fetchCategories(); // Recargar categorías
      alert('Categoría eliminada exitosamente');
    } catch (error) {
      console.error('Error al eliminar categoría:', error);
      alert('Error al eliminar la categoría: ' + (error.response?.data?.error || error.message));
    } finally {
      setIsLoading(false);
    }
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
          
          <button onClick={() => setIsCategoryModalVisible(true)} className="action-button secondary">
            <FaEdit />
            Gestionar Categorías
          </button>
          
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
                <label htmlFor="category">Categoría</label>
                <select
                  id="category"
                  name="category"
                  className="form-control"
                  defaultValue={editingCourse?.categoria_id || ''}
                >
                  <option value="">Seleccionar categoría</option>
                  {categories.map(category => (
                    <option key={category.id} value={category.id}>
                      {category.nombre}
                    </option>
                  ))}
                </select>
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
                {imageError && <div className="error-message">{imageError}</div>}
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

      {/* Modal de Gestión de Categorías */}
      {isCategoryModalVisible && (
        <div className="modal-overlay" onClick={() => setIsCategoryModalVisible(false)}>
          <div className="modal-container category-modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Gestionar Categorías</h3>
              <button className="modal-close" onClick={() => setIsCategoryModalVisible(false)}>
                <FaTimes />
              </button>
            </div>
            
            <div className="modal-content">
              {/* Formulario para nueva categoría */}
              <form onSubmit={handleCreateCategory} className="category-form">
                <h4>Crear Nueva Categoría</h4>
                <div className="form-group">
                  <label htmlFor="categoryName">Nombre de la Categoría</label>
                  <input
                    id="categoryName"
                    type="text"
                    className="form-control"
                    value={newCategoryName}
                    onChange={(e) => setNewCategoryName(e.target.value)}
                    placeholder="Ej: Programación"
                    required
                  />
                </div>
                
                <div className="form-group">
                  <label htmlFor="categoryDescription">Descripción</label>
                  <textarea
                    id="categoryDescription"
                    className="form-control"
                    value={newCategoryDescription}
                    onChange={(e) => setNewCategoryDescription(e.target.value)}
                    placeholder="Describe esta categoría..."
                    rows="2"
                  />
                </div>
                
                <button type="submit" className="btn-primary" disabled={!newCategoryName.trim()}>
                  <FaPlus /> Crear Categoría
                </button>
              </form>
              
              {/* Lista de categorías existentes */}
              <div className="categories-list">
                <h4>Categorías Existentes</h4>
                {categories.length === 0 ? (
                  <p>No hay categorías creadas</p>
                ) : (
                  <div className="category-items">
                    {categories.map(category => (
                      <div key={category.id} className="category-item">
                        <div className="category-info">
                          <strong>{category.nombre}</strong>
                          {category.descripcion && <p>{category.descripcion}</p>}
                        </div>
                        <button
                          onClick={() => handleDeleteCategory(category.id)}
                          className="btn-danger small"
                        >
                          <FaTrash />
                        </button>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default CoursesAdmin;