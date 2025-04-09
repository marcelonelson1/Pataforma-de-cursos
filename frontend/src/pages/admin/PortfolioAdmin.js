import React, { useState, useEffect, useCallback } from 'react';
import { FaPlus, FaEdit, FaTrash, FaSpinner, FaExclamationTriangle, FaEye, FaCheck } from 'react-icons/fa';
import axios from 'axios';
import './PortfolioAdmin.css';

const PortfolioAdmin = () => {
  const [projects, setProjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showForm, setShowForm] = useState(false);
  const [editingProject, setEditingProject] = useState(null);
  const [formData, setFormData] = useState({
    title: '',
    category: 'architecture',
    image: null,
    description: ''
  });
  const [previewUrl, setPreviewUrl] = useState('');
  const [submitLoading, setSubmitLoading] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');
  const [submitSuccess, setSubmitSuccess] = useState(false);
  
  // API URL base
  const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';

  // Categorías disponibles para los proyectos
  const categories = [
    { name: 'Arquitectura', value: 'architecture' },
    { name: 'Interiores', value: 'interiors' },
    { name: 'Paisajismo', value: 'landscape' },
    { name: 'Comercial', value: 'commercial' }
  ];

  // Función para obtener todos los proyectos - wrapped with useCallback
  const fetchProjects = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await axios.get(`${apiUrl}/api/portfolio`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      if (response.data.success) {
        setProjects(response.data.data || []);
      } else {
        setError(response.data.error || 'Error al cargar proyectos');
      }
    } catch (err) {
      setError('Error de conexión al servidor: ' + (err.response?.data?.error || err.message));
      console.error('Error al cargar proyectos:', err);
    } finally {
      setLoading(false);
    }
  }, [apiUrl]);

  // Cargar proyectos del portfolio al montar el componente
  useEffect(() => {
    fetchProjects();
  }, [fetchProjects]);

  // Manejar cambios en el formulario
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  // Manejar cambios en el archivo de imagen
  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setFormData({
        ...formData,
        image: file
      });
      
      // Crear URL para previsualización
      const fileReader = new FileReader();
      fileReader.onload = () => {
        setPreviewUrl(fileReader.result);
      };
      fileReader.readAsDataURL(file);
    }
  };

  // Abrir formulario para crear nuevo proyecto
  const handleNewProject = () => {
    setFormData({
      title: '',
      category: 'architecture',
      image: null,
      description: ''
    });
    setPreviewUrl('');
    setEditingProject(null);
    setShowForm(true);
    setSubmitSuccess(false);
  };

  // Abrir formulario para editar proyecto existente
  const handleEditProject = (project) => {
    setFormData({
      title: project.title,
      category: project.category,
      description: project.description || '',
      image: null // No se puede pre-cargar la imagen existente como File
    });
    
    // Combinar URL base si la imagen no tiene URL completa
    const imageUrl = project.image_url.startsWith('http') 
      ? project.image_url 
      : `${apiUrl}${project.image_url}`;
    
    setPreviewUrl(imageUrl); // Mostrar la imagen actual
    setEditingProject(project.id);
    setShowForm(true);
    setSubmitSuccess(false);
  };

  // Enviar formulario (crear o actualizar)
  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitLoading(true);
    setError(null);
    setSuccessMessage('');
    setSubmitSuccess(false);
    
    try {
      const formDataObj = new FormData();
      formDataObj.append('title', formData.title);
      formDataObj.append('category', formData.category);
      formDataObj.append('description', formData.description || '');
      
      // Solo adjuntar la imagen si hay una nueva
      if (formData.image) {
        formDataObj.append('image', formData.image);
      }
      
      let response;
      
      if (editingProject) {
        // Actualizar proyecto existente
        response = await axios.put(
          `${apiUrl}/api/portfolio/${editingProject}`, 
          formDataObj,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
              Authorization: `Bearer ${localStorage.getItem('token')}`
            }
          }
        );
      } else {
        // Crear nuevo proyecto
        response = await axios.post(
          `${apiUrl}/api/portfolio`, 
          formDataObj,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
              Authorization: `Bearer ${localStorage.getItem('token')}`
            }
          }
        );
      }
      
      if (response.data.success) {
        setSuccessMessage(editingProject ? 'Proyecto actualizado con éxito' : 'Proyecto creado con éxito');
        setSubmitSuccess(true);
        
        // Recargar la lista de proyectos
        fetchProjects(); 
        
        // Cerrar el formulario después de 2 segundos
        setTimeout(() => {
          setShowForm(false);
          setSuccessMessage('');
          setSubmitSuccess(false);
        }, 2000);
      } else {
        setError(response.data.error || 'Error al procesar la solicitud');
      }
    } catch (err) {
      setError('Error al guardar: ' + (err.response?.data?.error || err.message));
      console.error('Error al guardar proyecto:', err);
    } finally {
      setSubmitLoading(false);
    }
  };

  // Eliminar proyecto
  const handleDeleteProject = async (projectId) => {
    if (window.confirm('¿Estás seguro de que deseas eliminar este proyecto?')) {
      setLoading(true);
      try {
        const response = await axios.delete(`${apiUrl}/api/portfolio/${projectId}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`
          }
        });
        
        if (response.data.success) {
          // Eliminar proyecto de la lista local
          setProjects(projects.filter(project => project.id !== projectId));
          setSuccessMessage('Proyecto eliminado con éxito');
          
          // Limpiar mensaje después de 2 segundos
          setTimeout(() => {
            setSuccessMessage('');
          }, 2000);
        } else {
          setError(response.data.error || 'Error al eliminar el proyecto');
        }
      } catch (err) {
        setError('Error de conexión: ' + (err.response?.data?.error || err.message));
        console.error('Error al eliminar proyecto:', err);
      } finally {
        setLoading(false);
      }
    }
  };

  // Ver proyecto en el sitio público
  const handleViewProject = (project) => {
    // Abrir en una nueva pestaña la página del portfolio filtrada por categoría
    window.open(`/portfolio?category=${project.category}`, '_blank');
  };

  return (
    <div className="portfolio-admin">
      <div className="portfolio-admin__header">
        <h2 className="portfolio-admin__title">Gestión de Portfolio</h2>
        <button 
          className="portfolio-admin__btn portfolio-admin__btn--primary"
          onClick={handleNewProject}
          disabled={loading}
        >
          <FaPlus /> Nuevo Proyecto
        </button>
      </div>

      {successMessage && (
        <div className="portfolio-admin__alert portfolio-admin__alert--success">
          {successMessage}
        </div>
      )}

      {error && (
        <div className="portfolio-admin__alert portfolio-admin__alert--danger">
          <FaExclamationTriangle /> {error}
        </div>
      )}

      {showForm && (
        <div className={`portfolio-admin__form-container ${submitSuccess ? 'portfolio-admin__form-container--success' : ''}`}>
          <h3 className="portfolio-admin__form-title">{editingProject ? 'Editar Proyecto' : 'Nuevo Proyecto'}</h3>
          <form onSubmit={handleSubmit} className="portfolio-admin__form">
            <div className="portfolio-admin__form-group">
              <label htmlFor="title" className="portfolio-admin__form-label">Título:</label>
              <input
                type="text"
                id="title"
                name="title"
                className="portfolio-admin__form-control"
                value={formData.title}
                onChange={handleChange}
                required
              />
            </div>

            <div className="portfolio-admin__form-group">
              <label htmlFor="category" className="portfolio-admin__form-label">Categoría:</label>
              <select
                id="category"
                name="category"
                className="portfolio-admin__form-control"
                value={formData.category}
                onChange={handleChange}
                required
              >
                {categories.map(cat => (
                  <option key={cat.value} value={cat.value}>
                    {cat.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="portfolio-admin__form-group">
              <label htmlFor="description" className="portfolio-admin__form-label">Descripción:</label>
              <textarea
                id="description"
                name="description"
                className="portfolio-admin__form-control portfolio-admin__form-textarea"
                value={formData.description}
                onChange={handleChange}
                rows="3"
              />
            </div>

            <div className="portfolio-admin__form-group">
              <label htmlFor="image" className="portfolio-admin__form-label">
                Imagen: {editingProject ? '(Subir solo si deseas cambiarla)' : '(Requerida)'}
              </label>
              <input
                type="file"
                id="image"
                name="image"
                accept="image/*"
                className="portfolio-admin__file-input"
                onChange={handleImageChange}
                required={!editingProject}
              />
              {previewUrl && (
                <div className="portfolio-admin__image-preview">
                  <img 
                    src={previewUrl} 
                    alt="Vista previa" 
                    className="portfolio-admin__thumbnail" 
                    style={{ maxHeight: '200px' }} 
                  />
                </div>
              )}
            </div>

            <div className="portfolio-admin__form-buttons">
              <button
                type="button"
                className="portfolio-admin__btn portfolio-admin__btn--secondary"
                onClick={() => setShowForm(false)}
                disabled={submitLoading}
              >
                Cancelar
              </button>
              <button
                type="submit"
                className={`portfolio-admin__btn ${submitSuccess ? 'portfolio-admin__btn--success' : 'portfolio-admin__btn--primary'}`}
                disabled={submitLoading}
              >
                {submitLoading ? (
                  <>
                    <FaSpinner className="portfolio-admin__spinner" /> Guardando...
                  </>
                ) : submitSuccess ? (
                  <>
                    <FaCheck /> ¡Listo!
                  </>
                ) : (
                  "Guardar"
                )}
              </button>
            </div>
          </form>
        </div>
      )}

      {loading && !showForm ? (
        <div className="portfolio-admin__loading">
          <FaSpinner className="portfolio-admin__spinner" /> Cargando proyectos...
        </div>
      ) : (
        <div className="portfolio-admin__grid">
          {projects.length === 0 ? (
            <div className="portfolio-admin__empty">
              No hay proyectos disponibles. Crea uno nuevo.
            </div>
          ) : (
            projects.map(project => {
              // Combinar URL base si la imagen no tiene URL completa
              const imageUrl = project.image_url.startsWith('http') 
                ? project.image_url 
                : `${apiUrl}${project.image_url}`;
              
              return (
                <div key={project.id} className="portfolio-admin__card">
                  <div className="portfolio-admin__card-image">
                    <img 
                      src={imageUrl} 
                      alt={project.title} 
                      className="portfolio-admin__card-img" 
                      onError={(e) => {
                        // Si la imagen falla al cargar, mostrar una imagen de respaldo
                        e.target.src = '/placeholder-image.jpg';
                        console.error(`Error al cargar imagen: ${imageUrl}`);
                      }}
                    />
                  </div>
                  <div className="portfolio-admin__card-body">
                    <h5 className="portfolio-admin__card-title">{project.title}</h5>
                    <span className="portfolio-admin__card-category">
                      {categories.find(cat => cat.value === project.category)?.name || project.category}
                    </span>
                    {project.description && (
                      <p className="portfolio-admin__card-text">{project.description}</p>
                    )}
                  </div>
                  <div className="portfolio-admin__card-footer">
                    <button
                      className="portfolio-admin__action-btn portfolio-admin__action-btn--view"
                      onClick={() => handleViewProject(project)}
                      title="Ver en sitio"
                    >
                      <FaEye />
                    </button>
                    <button
                      className="portfolio-admin__action-btn portfolio-admin__action-btn--edit"
                      onClick={() => handleEditProject(project)}
                      title="Editar"
                    >
                      <FaEdit />
                    </button>
                    <button
                      className="portfolio-admin__action-btn portfolio-admin__action-btn--delete"
                      onClick={() => handleDeleteProject(project.id)}
                      title="Eliminar"
                    >
                      <FaTrash />
                    </button>
                  </div>
                </div>
              );
            })
          )}
        </div>
      )}
    </div>
  );
};

export default PortfolioAdmin;