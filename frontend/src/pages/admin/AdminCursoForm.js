import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faSave, faTimes, faImage, faEye, faMobile, faDesktop
} from '@fortawesome/free-solid-svg-icons';
import { Editor } from '@tinymce/tinymce-react';
import './AdminPages.css';

function AdminCursoForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isEditing = !!id;
  
  const [formData, setFormData] = useState({
    titulo: '',
    descripcion: '',
    contenido: '',
    precio: 29.99
  });
  
  const [thumbnail, setThumbnail] = useState(null);
  const [thumbnailPreview, setThumbnailPreview] = useState('');
  const [loading, setLoading] = useState(isEditing);
  const [error, setError] = useState(null);
  const [previewMode, setPreviewMode] = useState(null); // null, 'mobile', 'desktop'
  const [errors, setErrors] = useState({});
  const [saving, setSaving] = useState(false);

  // Cargar datos del curso si estamos en modo edición
  useEffect(() => {
    if (isEditing) {
      const fetchCurso = async () => {
        try {
          const token = localStorage.getItem('token');
          const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
          
          const response = await fetch(`${apiUrl}/api/admin/cursos/${id}`, {
            headers: {
              'Authorization': `Bearer ${token}`,
              'Content-Type': 'application/json'
            }
          });

          if (!response.ok) {
            throw new Error('Error al obtener el curso');
          }

          const data = await response.json();
          setFormData({
            titulo: data.curso.titulo,
            descripcion: data.curso.descripcion,
            contenido: data.curso.contenido,
            precio: data.curso.precio
          });
          
          // Si hay una URL de thumbnail, establecer la vista previa
          if (data.curso.thumbnail_url) {
            setThumbnailPreview(data.curso.thumbnail_url);
          }
          
          setError(null);
        } catch (err) {
          console.error('Error al cargar curso:', err);
          setError('No se pudo cargar la información del curso. Por favor, intenta de nuevo.');
        } finally {
          setLoading(false);
        }
      };

      fetchCurso();
    }
  }, [id, isEditing]);

  // Datos de ejemplo para modo edición
  useEffect(() => {
    if (isEditing && !loading) {
      // Si no hay datos de la API (en desarrollo), usar datos de ejemplo
      if (formData.titulo === '' && id === "1") {
        setFormData({
          titulo: 'Curso de React',
          descripcion: 'Aprende React desde cero hasta un nivel avanzado con proyectos prácticos.',
          contenido: '<p>En este curso aprenderás los fundamentos de React, incluyendo:</p><ul><li>Componentes</li><li>Estado y Props</li><li>Hooks</li><li>Rutas con React Router</li><li>Redux para gestión de estado</li></ul><p>Y mucho más...</p>',
          precio: 29.99
        });
        setThumbnailPreview('/curso-thumbnail-1.jpg');
      }
    }
  }, [isEditing, loading, formData.titulo, id]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    
    if (name === 'precio') {
      // Validar que el precio sea un número positivo
      const precio = parseFloat(value);
      if (!isNaN(precio) && precio >= 0) {
        setFormData(prev => ({ ...prev, [name]: precio }));
      }
    } else {
      setFormData(prev => ({ ...prev, [name]: value }));
    }
    
    // Limpiar error del campo
    if (errors[name]) {
      setErrors(prev => ({ ...prev, [name]: null }));
    }
  };

  const handleContentChange = (content) => {
    setFormData(prev => ({ ...prev, contenido: content }));
    
    // Limpiar error del campo
    if (errors.contenido) {
      setErrors(prev => ({ ...prev, contenido: null }));
    }
  };

  const handleThumbnailChange = (e) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      
      // Validar tipo de archivo
      if (!file.type.startsWith('image/')) {
        setErrors(prev => ({ 
          ...prev, 
          thumbnail: 'Por favor, selecciona una imagen válida (JPG, PNG, WebP)'
        }));
        return;
      }
      
      // Validar tamaño (máx. 5MB)
      if (file.size > 5 * 1024 * 1024) {
        setErrors(prev => ({ 
          ...prev, 
          thumbnail: 'La imagen debe ser menor a 5MB'
        }));
        return;
      }
      
      setThumbnail(file);
      setThumbnailPreview(URL.createObjectURL(file));
      
      // Limpiar error
      if (errors.thumbnail) {
        setErrors(prev => ({ ...prev, thumbnail: null }));
      }
    }
  };

  const validateForm = () => {
    const newErrors = {};
    
    if (!formData.titulo.trim()) {
      newErrors.titulo = 'El título es obligatorio';
    }
    
    if (!formData.descripcion.trim()) {
      newErrors.descripcion = 'La descripción es obligatoria';
    }
    
    if (!formData.contenido.trim()) {
      newErrors.contenido = 'El contenido es obligatorio';
    }
    
    if (isNaN(formData.precio) || formData.precio <= 0) {
      newErrors.precio = 'Ingresa un precio válido mayor a 0';
    }
    
    if (!isEditing && !thumbnail && !thumbnailPreview) {
      newErrors.thumbnail = 'La imagen del curso es obligatoria';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }
    
    setSaving(true);
    
    try {
      const token = localStorage.getItem('token');
      const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
      
      // Crear FormData para enviar archivos
      const formDataObj = new FormData();
      formDataObj.append('titulo', formData.titulo);
      formDataObj.append('descripcion', formData.descripcion);
      formDataObj.append('contenido', formData.contenido);
      formDataObj.append('precio', formData.precio);
      
      if (thumbnail) {
        formDataObj.append('thumbnail', thumbnail);
      }
      
      // URL y método según sea creación o edición
      const url = isEditing 
        ? `${apiUrl}/api/admin/cursos/${id}`
        : `${apiUrl}/api/admin/cursos`;
      
      const method = isEditing ? 'PUT' : 'POST';
      
      const response = await fetch(url, {
        method,
        headers: {
          'Authorization': `Bearer ${token}`,
          // No establecer Content-Type para permitir que el navegador configure el límite multipart
        },
        body: formDataObj
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Error al guardar el curso');
      }

      const data = await response.json();
      
      alert(isEditing ? 'Curso actualizado correctamente' : 'Curso creado correctamente');
      
      // Redirigir a la gestión de capítulos si es un curso nuevo
      if (!isEditing && data.id) {
        navigate(`/admin/cursos/${data.id}/capitulos`);
      } else {
        navigate('/admin/cursos');
      }
      
    } catch (err) {
      console.error('Error al guardar curso:', err);
      setError(`Error al ${isEditing ? 'actualizar' : 'crear'} el curso: ${err.message}`);
    } finally {
      setSaving(false);
    }
  };

  const handlePreview = (mode) => {
    if (previewMode === mode) {
      setPreviewMode(null);
    } else {
      setPreviewMode(mode);
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

  return (
    <div className="admin-curso-form-page">
      <div className="admin-page-title">
        <h1>{isEditing ? 'Editar Curso' : 'Crear Nuevo Curso'}</h1>
        <p>{isEditing ? 'Actualiza la información del curso' : 'Completa el formulario para crear un nuevo curso'}</p>
      </div>

      {error && (
        <div className="admin-error-container">
          <p>{error}</p>
        </div>
      )}

      {previewMode && (
        <div className="admin-preview-container">
          <div className={`admin-preview-frame ${previewMode}`}>
            <div className="admin-preview-header">
              <button 
                className="admin-btn admin-btn-secondary admin-btn-sm"
                onClick={() => setPreviewMode(null)}
              >
                <FontAwesomeIcon icon={faTimes} /> Cerrar Vista Previa
              </button>
            </div>
            <div className="admin-preview-content">
              <h1>{formData.titulo}</h1>
              <div className="admin-preview-thumbnail">
                {thumbnailPreview ? (
                  <img src={thumbnailPreview} alt={formData.titulo} />
                ) : (
                  <div className="admin-preview-placeholder">
                    <FontAwesomeIcon icon={faImage} />
                    <p>Sin imagen</p>
                  </div>
                )}
              </div>
              <p className="admin-preview-description">{formData.descripcion}</p>
              <div 
                className="admin-preview-contenido"
                dangerouslySetInnerHTML={{ __html: formData.contenido }}
              />
              <div className="admin-preview-price">
                <p>Precio: ${parseFloat(formData.precio).toFixed(2)}</p>
              </div>
            </div>
          </div>
        </div>
      )}

      <form className="admin-form" onSubmit={handleSubmit}>
        <div className="admin-form-card">
          <div className="admin-form-section">
            <h2 className="admin-form-section-title">Información Básica</h2>
            
            <div className="admin-form-row">
              <div className="admin-form-group col-12">
                <label className="admin-form-label" htmlFor="titulo">
                  Título del Curso *
                </label>
                <input
                  type="text"
                  id="titulo"
                  name="titulo"
                  className={`admin-form-control ${errors.titulo ? 'is-invalid' : ''}`}
                  value={formData.titulo}
                  onChange={handleChange}
                  placeholder="Ej: Curso Completo de React"
                  required
                />
                {errors.titulo && <span className="admin-form-error">{errors.titulo}</span>}
              </div>
            </div>
            
            <div className="admin-form-row">
              <div className="admin-form-group col-12">
                <label className="admin-form-label" htmlFor="descripcion">
                  Descripción Breve *
                </label>
                <textarea
                  id="descripcion"
                  name="descripcion"
                  className={`admin-form-control admin-form-textarea ${errors.descripcion ? 'is-invalid' : ''}`}
                  value={formData.descripcion}
                  onChange={handleChange}
                  placeholder="Describe brevemente el curso en 2-3 frases"
                  rows="3"
                  required
                />
                {errors.descripcion && <span className="admin-form-error">{errors.descripcion}</span>}
              </div>
            </div>
            
            <div className="admin-form-row">
              <div className="admin-form-group col-6">
                <label className="admin-form-label" htmlFor="precio">
                  Precio ($) *
                </label>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  id="precio"
                  name="precio"
                  className={`admin-form-control ${errors.precio ? 'is-invalid' : ''}`}
                  value={formData.precio}
                  onChange={handleChange}
                  required
                />
                {errors.precio && <span className="admin-form-error">{errors.precio}</span>}
              </div>
              
              <div className="admin-form-group col-6">
                <label className="admin-form-label" htmlFor="thumbnail">
                  Imagen del Curso {isEditing ? '' : '*'}
                </label>
                <div className={`admin-dropzone ${errors.thumbnail ? 'is-invalid' : ''}`}>
                  <input
                    type="file"
                    id="thumbnail"
                    name="thumbnail"
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
                          Haz clic para subir una imagen
                        </div>
                        <div className="admin-dropzone-help">
                          Formatos: JPG, PNG, WebP | Máx: 5MB
                        </div>
                      </>
                    )}
                  </label>
                </div>
                {errors.thumbnail && <span className="admin-form-error">{errors.thumbnail}</span>}
              </div>
            </div>
          </div>
          
          <div className="admin-form-section">
            <h2 className="admin-form-section-title">Contenido Detallado</h2>
            
            <div className="admin-form-row">
              <div className="admin-form-group col-12">
                <label className="admin-form-label">
                  Contenido del Curso *
                </label>
                <Editor
                  apiKey="tu-api-key-de-tinymce" // Reemplazar con API key real o usar editor alternativo
                  value={formData.contenido}
                  init={{
                    height: 400,
                    menubar: false,
                    plugins: [
                      'advlist autolink lists link image charmap print preview anchor',
                      'searchreplace visualblocks code fullscreen',
                      'insertdatetime media table paste code help wordcount'
                    ],
                    toolbar:
                      'undo redo | formatselect | bold italic backcolor | \
                      alignleft aligncenter alignright alignjustify | \
                      bullist numlist outdent indent | removeformat | help'
                  }}
                  onEditorChange={handleContentChange}
                />
                {errors.contenido && <span className="admin-form-error">{errors.contenido}</span>}
              </div>
            </div>
          </div>
          
          <div className="admin-form-actions">
            <div className="admin-preview-buttons">
              <button 
                type="button" 
                className={`admin-btn admin-btn-secondary ${previewMode === 'mobile' ? 'active' : ''}`}
                onClick={() => handlePreview('mobile')}
              >
                <FontAwesomeIcon icon={faMobile} /> Vista Móvil
              </button>
              <button 
                type="button" 
                className={`admin-btn admin-btn-secondary ${previewMode === 'desktop' ? 'active' : ''}`}
                onClick={() => handlePreview('desktop')}
              >
                <FontAwesomeIcon icon={faDesktop} /> Vista Desktop
              </button>
            </div>
            
            <div>
              <Link to="/admin/cursos" className="admin-btn admin-btn-secondary">
                <FontAwesomeIcon icon={faTimes} /> Cancelar
              </Link>
              <button 
                type="submit" 
                className="admin-btn admin-btn-primary"
                disabled={saving}
              >
                <FontAwesomeIcon icon={faSave} /> {saving ? 'Guardando...' : 'Guardar Curso'}
              </button>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
}

export default AdminCursoForm;