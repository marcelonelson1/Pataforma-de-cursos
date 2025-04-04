import React, { useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import { 
  FaPlus, 
  FaEdit, 
  FaTrash, 
  FaUpload, 
  FaClock, 
  FaTimes, 
  FaEye,
  FaSearch
} from 'react-icons/fa';
import './CoursesAdmin.css';

const CoursesAdmin = () => {
  const [courses, setCourses] = useState([
    { 
      id: uuidv4(),
      title: 'Curso de React Avanzado',
      price: 49.99,
      description: 'Domina React con proyectos reales y las últimas características de React 18. Aprenderás Hooks avanzados, patrones de optimización y arquitectura de aplicaciones profesionales.',
      image: null,
      publishedDate: '2023-06-15',
      status: 'Publicado',
      students: 120,
      rating: 4.8,
      chapters: [
        {
          id: uuidv4(),
          title: 'Introducción a React',
          description: 'Conceptos fundamentales de React y configuración del entorno',
          video: null,
          duration: '15:00',
          videoFile: null,
          isPublished: true
        },
        {
          id: uuidv4(),
          title: 'Hooks Avanzados',
          description: 'UseReducer, useMemo, useCallback y custom hooks',
          video: null,
          duration: '28:30',
          videoFile: null,
          isPublished: true
        }
      ]
    },
    { 
      id: uuidv4(),
      title: 'Desarrollo con Node.js',
      price: 59.99,
      description: 'Aprende a crear APIs RESTful y aplicaciones backend con Node.js, Express y MongoDB.',
      image: null,
      publishedDate: '2023-05-22',
      status: 'Borrador',
      students: 85,
      rating: 4.5,
      chapters: [
        {
          id: uuidv4(),
          title: 'Introducción a Node.js',
          description: 'Fundamentos y configuración del entorno',
          video: null,
          duration: '18:45',
          videoFile: null,
          isPublished: false
        }
      ]
    }
  ]);

  const [editingCourse, setEditingCourse] = useState(null);
  const [editingChapter, setEditingChapter] = useState(null);
  const [isFormVisible, setIsFormVisible] = useState(false);
  const [isChapterModalVisible, setIsChapterModalVisible] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [isLoading, setIsLoading] = useState(true);

  // Simulación de carga de datos
  useEffect(() => {
    setIsLoading(true);
    // Simulamos una carga asincrónica
    const timer = setTimeout(() => {
      setIsLoading(false);
    }, 800);
    
    return () => clearTimeout(timer);
  }, []);

  // Filtrado de cursos
  const filteredCourses = courses.filter(course => {
    const matchesSearch = course.title.toLowerCase().includes(searchTerm.toLowerCase()) || 
                         course.description.toLowerCase().includes(searchTerm.toLowerCase());
    
    const matchesStatus = statusFilter === 'all' || 
                          (statusFilter === 'published' && course.status === 'Publicado') ||
                          (statusFilter === 'draft' && course.status === 'Borrador');
    
    return matchesSearch && matchesStatus;
  });

  // Course handlers
  const handleCourseSubmit = (e) => {
    e.preventDefault();
    setIsLoading(true);

    const formData = new FormData(e.target);
    
    // Simulamos una operación de guardado
    setTimeout(() => {
      const newCourse = {
        id: editingCourse?.id || uuidv4(),
        title: formData.get('title'),
        price: parseFloat(formData.get('price')),
        description: formData.get('description'),
        status: formData.get('status'),
        publishedDate: formData.get('status') === 'Publicado' ? new Date().toISOString().split('T')[0] : null,
        chapters: editingCourse?.chapters || [],
        students: editingCourse?.students || 0,
        rating: editingCourse?.rating || 0,
        image: editingCourse?.image || null
      };
  
      setCourses(prev => editingCourse 
        ? prev.map(c => c.id === editingCourse.id ? newCourse : c)
        : [...prev, newCourse]
      );
      
      setIsFormVisible(false);
      setEditingCourse(null);
      setIsLoading(false);
    }, 800);
  };

  const handleDeleteCourse = (id) => {
    if (window.confirm('¿Estás seguro de que quieres eliminar este curso? Esta acción no se puede deshacer.')) {
      setIsLoading(true);
      
      // Simulamos operación asincrónica
      setTimeout(() => {
        setCourses(courses.filter(course => course.id !== id));
        setIsLoading(false);
      }, 500);
    }
  };

  const handleImageUpload = (e) => {
    const file = e.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (event) => {
        setEditingCourse(prev => ({
          ...prev,
          image: event.target.result
        }));
      };
      reader.readAsDataURL(file);
    }
  };

  // Chapter handlers
  const handleChapterSubmit = (e) => {
    e.preventDefault();
    setIsLoading(true);

    const formData = new FormData(e.target);
    
    // Simulamos una operación asincrónica
    setTimeout(() => {
      const newChapter = {
        id: editingChapter?.id || uuidv4(),
        title: formData.get('chapterTitle'),
        description: formData.get('chapterDescription'),
        duration: formData.get('duration'),
        isPublished: formData.get('isPublished') === 'on',
        videoFile: editingChapter?.videoFile || null,
        video: editingChapter?.video || null
      };
  
      const updatedChapters = editingChapter 
        ? editingCourse.chapters.map(ch => ch.id === editingChapter.id ? newChapter : ch)
        : [...editingCourse.chapters, newChapter];
  
      setEditingCourse(prev => ({
        ...prev,
        chapters: updatedChapters
      }));
      
      setEditingChapter(null);
      setIsChapterModalVisible(false);
      setIsLoading(false);
    }, 600);
  };

  const handleVideoUpload = (e) => {
    const file = e.target.files[0];
    if (file) {
      const videoUrl = URL.createObjectURL(file);
      
      // Obtener duración del video (en un entorno real)
      // Aquí simplemente simulamos la duración
      const minutes = Math.floor(Math.random() * 30) + 5;
      const seconds = Math.floor(Math.random() * 59);
      const duration = `${minutes}:${seconds.toString().padStart(2, '0')}`;
      
      setEditingChapter(prev => ({
        ...prev,
        video: videoUrl,
        videoFile: file,
        duration: duration
      }));
    }
  };

  const handleDeleteChapter = (chapterId) => {
    if (window.confirm('¿Estás seguro de que quieres eliminar este capítulo?')) {
      setEditingCourse(prev => ({
        ...prev,
        chapters: prev.chapters.filter(ch => ch.id !== chapterId)
      }));
    }
  };

  const handleNewCourse = () => {
    setEditingCourse({ 
      chapters: [],
      status: 'Borrador',
      publishedDate: null,
      students: 0,
      rating: 0
    });
    setIsFormVisible(true);
  };

  const handleEditCourse = (course) => {
    setEditingCourse(course);
    setIsFormVisible(true);
  };

  const handleNewChapter = () => {
    setEditingChapter({});
    setIsChapterModalVisible(true);
  };

  const handleEditChapter = (chapter) => {
    setEditingChapter(chapter);
    setIsChapterModalVisible(true);
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
    }
  };

  // Previsualización del curso
  const handlePreviewCourse = (course) => {
    console.log('Previsualizando curso:', course.title);
    // Aquí implementaríamos la navegación a la vista de previsualización
  };

  // Loading overlay
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

      {/* Listado de Cursos */}
      <div className="courses-list">
        {filteredCourses.length === 0 ? (
          <div className="no-results">
            <p>No se encontraron cursos que coincidan con tu búsqueda.</p>
          </div>
        ) : (
          filteredCourses.map(course => (
            <div key={course.id} className="course-card">
              <div className="course-header">
                <h3>{course.title}</h3>
                <span className="price">${course.price.toFixed(2)}</span>
              </div>
              
              <div className="course-meta">
                <span className={`status-badge ${course.status === 'Publicado' ? 'published' : 'draft'}`}>
                  {course.status}
                </span>
                {course.status === 'Publicado' && (
                  <span className="date-published">
                    Publicado: {course.publishedDate}
                  </span>
                )}
              </div>
              
              <p className="description">{course.description}</p>
              
              <div className="course-stats">
                <div className="stat">
                  <span className="stat-value">{course.students}</span>
                  <span className="stat-label">Estudiantes</span>
                </div>
                <div className="stat">
                  <span className="stat-value">{course.rating.toFixed(1)}</span>
                  <span className="stat-label">Valoración</span>
                </div>
                <div className="stat">
                  <span className="stat-value">{course.chapters.length}</span>
                  <span className="stat-label">Capítulos</span>
                </div>
              </div>
              
              <div className="chapters-summary">
                <h4>Capítulos <span>{course.chapters.length}</span></h4>
                {course.chapters.length > 0 ? (
                  <ul>
                    {course.chapters.slice(0, 3).map(chapter => (
                      <li key={chapter.id}>
                        <span>{chapter.title}</span>
                        <span className="chapter-duration">
                          <FaClock /> {chapter.duration}
                        </span>
                      </li>
                    ))}
                    {course.chapters.length > 3 && (
                      <li className="more-chapters">
                        +{course.chapters.length - 3} capítulos más
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

      {/* Modal para editar/crear curso */}
      {isFormVisible && (
        <div className="form-overlay">
          <div className="form-container">
            <div className="form-header">
              <h3>{editingCourse.id ? 'Editar Curso' : 'Nuevo Curso'}</h3>
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
                  defaultValue={editingCourse.title || ''}
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
                  defaultValue={editingCourse.description || ''}
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
                    defaultValue={editingCourse.price || 0}
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
                    defaultValue={editingCourse.status || 'Borrador'}
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
                {editingCourse.image && (
                  <div className="image-preview">
                    <img src={editingCourse.image} alt="Vista previa del curso" />
                  </div>
                )}
              </div>

              <div className="chapters-section">
                <div className="section-header">
                  <h3>Capítulos</h3>
                  <button
                    type="button"
                    onClick={handleNewChapter}
                    className="action-button small"
                  >
                    <FaPlus />
                    Añadir Capítulo
                  </button>
                </div>
                
                {editingCourse.chapters.length > 0 ? (
                  <div className="chapters-list">
                    {editingCourse.chapters.map((chapter, index) => (
                      <div key={chapter.id} className="chapter-preview">
                        <div className="chapter-info">
                          <div className="chapter-number">{index + 1}</div>
                          <div>
                            <h4>{chapter.title}</h4>
                            <p>{chapter.description.substring(0, 60)}{chapter.description.length > 60 ? '...' : ''}</p>
                            <div className="chapter-meta">
                              <span className="chapter-duration">
                                <FaClock /> {chapter.duration}
                              </span>
                              <span className={`chapter-status ${chapter.isPublished ? 'published' : 'draft'}`}>
                                {chapter.isPublished ? 'Publicado' : 'Borrador'}
                              </span>
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
                    onClick={handleNewChapter}
                    className="add-chapter-btn"
                  >
                    <FaPlus /> Añadir tu primer capítulo
                  </button>
                )}
              </div>

              <div className="form-actions">
                <button type="button" className="btn-cancel" onClick={closeForm}>
                  Cancelar
                </button>
                <button type="submit" className="btn-primary">
                  {editingCourse.id ? 'Actualizar' : 'Crear'} Curso
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Modal para editar/crear capítulo */}
      {isChapterModalVisible && (
        <div className="chapter-modal">
          <div className="chapter-modal-overlay" onClick={closeChapterModal} />
          <form onSubmit={handleChapterSubmit} className="form-container">
            <div className="form-header">
              <h3>{editingChapter.id ? 'Editar Capítulo' : 'Nuevo Capítulo'}</h3>
              <button type="button" className="form-close" onClick={closeChapterModal}>
                <FaTimes />
              </button>
            </div>
            
            <div className="form-group">
              <label htmlFor="chapterTitle">Título del Capítulo</label>
              <input
                id="chapterTitle"
                name="chapterTitle"
                type="text"
                className="form-control"
                defaultValue={editingChapter.title || ''}
                placeholder="Introduce el título del capítulo"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="chapterDescription">Descripción</label>
              <textarea
                id="chapterDescription"
                name="chapterDescription"
                className="form-control"
                defaultValue={editingChapter.description || ''}
                placeholder="Describe el contenido de este capítulo"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="chapterVideo">Video del Capítulo</label>
              <div className="file-upload">
                <input
                  id="chapterVideo"
                  type="file"
                  accept="video/mp4,video/webm"
                  onChange={handleVideoUpload}
                  className="form-control-file"
                />
                <label htmlFor="chapterVideo" className="upload-label">
                  <FaUpload /> Subir Video
                </label>
              </div>
              {editingChapter.video && (
                <div className="video-preview">
                  <video controls>
                    <source src={editingChapter.video} type="video/mp4" />
                    Tu navegador no soporta el elemento de video.
                  </video>
                </div>
              )}
            </div>

            <div className="form-row">
              <div className="form-group">
                <label htmlFor="duration">Duración (MM:SS)</label>
                <input
                  id="duration"
                  name="duration"
                  type="text"
                  pattern="[0-9]{1,2}:[0-9]{2}"
                  className="form-control"
                  defaultValue={editingChapter.duration || ''}
                  placeholder="Ej: 15:30"
                  required
                />
              </div>
              
              <div className="form-group">
                <label className="checkbox-label">
                  <input
                    type="checkbox"
                    name="isPublished"
                    defaultChecked={editingChapter.isPublished || false}
                  />
                  <span className="checkmark"></span>
                  Publicar este capítulo
                </label>
              </div>
            </div>

            <div className="form-actions">
              <button type="button" className="btn-cancel" onClick={closeChapterModal}>
                Cancelar
              </button>
              <button type="submit" className="btn-primary">
                Guardar Capítulo
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
};

export default CoursesAdmin;