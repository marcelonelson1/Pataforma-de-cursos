import React, { useState, useEffect, useCallback } from 'react';
import { FaPlus, FaEdit, FaTrash, FaSpinner, FaExclamationTriangle, FaEye, FaCheck } from 'react-icons/fa';
import PortfolioService from '../../components/services/PortfolioService';
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
    category: 'arquitectura-residencial',
    customCategory: '',
    image: null,
    description: ''
  });
  const [previewUrl, setPreviewUrl] = useState('');
  const [submitLoading, setSubmitLoading] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');
  const [submitSuccess, setSubmitSuccess] = useState(false);
  

  // Categorías disponibles para renders arquitectónicos
  const categories = [
    { name: 'Arquitectura Residencial', value: 'arquitectura-residencial' },
    { name: 'Interiores', value: 'interiores' },
    { name: 'Paisajismo', value: 'paisajismo' },
    { name: 'Comercial', value: 'comercial' },
    { name: 'Hospitalario', value: 'hospitalario' },
    { name: 'Educativo', value: 'educativo' },
    { name: 'Industrial', value: 'industrial' },
    { name: 'Urbanismo', value: 'urbanismo' },
    { name: 'Personalizada', value: 'custom' }
  ];

  // Función para obtener todos los proyectos - wrapped with useCallback
  const fetchProjects = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await PortfolioService.getAllProjectsAdmin();
      
      if (response.data.success) {
        const projectsData = response.data.data?.projects || response.data.data || [];
        
        // Construct full URLs for portfolio images
        const projectsWithFullUrls = (Array.isArray(projectsData) ? projectsData : []).map(project => {
          let imageUrl = 'https://via.placeholder.com/300x200?text=Sin+imagen';
          
          if (project.image_url) {
            // Si la URL ya incluye /static/portfolio/, no añadir el prefijo
            if (project.image_url.startsWith('/static/portfolio/')) {
              imageUrl = `http://localhost:8005${project.image_url}`;
            } else {
              // Si es solo el filename, añadir el path completo
              imageUrl = `http://localhost:8005/static/portfolio/${project.image_url}`;
            }
          }
          
          return {
            ...project,
            image_url: imageUrl
          };
        });
        
        setProjects(projectsWithFullUrls);
      } else {
        setError(response.data.error || 'Error al cargar proyectos');
      }
    } catch (err) {
      setError('Error de conexión al servidor: ' + (err.response?.data?.error || err.message));
      console.error('Error al cargar proyectos:', err);
    } finally {
      setLoading(false);
    }
  }, []);

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
      category: 'arquitectura-residencial',
      customCategory: '',
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
    const isCustomCategory = !categories.some(cat => cat.value === project.category);
    setFormData({
      title: project.title,
      category: isCustomCategory ? 'custom' : project.category,
      customCategory: isCustomCategory ? project.category : '',
      description: project.description || '',
      image: null // No se puede pre-cargar la imagen existente como File
    });
    
    // Usar URL directa de la imagen
    const imageUrl = project.image_url;
    
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
      // Validar categoría personalizada
      const categoryToUse = formData.category === 'custom' ? formData.customCategory : formData.category;
      if (!categoryToUse || categoryToUse.trim() === '') {
        setError('La categoría es requerida');
        return;
      }

      // Preparar datos del proyecto (solo datos, sin imagen)
      const projectData = {
        title: formData.title,
        category: categoryToUse.trim(),
        description: formData.description || ''
      };
      
      let response;
      
      if (editingProject) {
        // Actualizar proyecto existente (solo datos)
        response = await PortfolioService.updateProject(editingProject, projectData);
      } else {
        // Crear nuevo proyecto (solo datos)
        response = await PortfolioService.createProject(projectData);
      }
      
      if (response.data.success) {
        const projectId = editingProject || response.data.data.id;
        
        // Si hay imagen nueva, subirla usando el endpoint dedicado
        if (formData.image) {
          try {
            console.log('Uploading image for project:', projectId);
            const imageResponse = await PortfolioService.uploadProjectImage(projectId, formData.image);
            
            if (imageResponse.data.success) {
              console.log('Image uploaded successfully:', imageResponse.data.data.image_url);
            } else {
              console.warn('Image upload failed:', imageResponse.data.error);
              setError('Proyecto guardado pero falló la subida de imagen: ' + imageResponse.data.error);
            }
          } catch (imageError) {
            console.error('Error uploading image:', imageError);
            setError('Proyecto guardado pero falló la subida de imagen: ' + (imageError.response?.data?.error || imageError.message));
          }
        }
        
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
        const response = await PortfolioService.deleteProject(projectId);
        
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

            {formData.category === 'custom' && (
              <div className="portfolio-admin__form-group">
                <label htmlFor="customCategory" className="portfolio-admin__form-label">Categoría Personalizada:</label>
                <input
                  type="text"
                  id="customCategory"
                  name="customCategory"
                  className="portfolio-admin__form-control"
                  value={formData.customCategory}
                  onChange={handleChange}
                  placeholder="Escriba su categoría personalizada"
                  required
                />
              </div>
            )}

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
              // Usar URL directa de la imagen
              const imageUrl = project.image_url;
              
              return (
                <div key={project.id} className="portfolio-admin__card">
                  <div className="portfolio-admin__card-image">
                    <img 
                      src={imageUrl} 
                      alt={project.title} 
                      className="portfolio-admin__card-img" 
                      onError={(e) => {
                        // Si la imagen falla al cargar, mostrar una imagen de respaldo
                        if (!e.target.src.includes('via.placeholder.com')) {
                          e.target.src = 'https://via.placeholder.com/300x200?text=Error+de+imagen';
                          console.error(`Error al cargar imagen del proyecto ${project.id}: ${imageUrl}`);
                        }
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