import React, { useState, useEffect } from 'react';
import homeImagesService from '../../components/services/HomeImagesService';
import { 
  FaTrash, 
  FaEdit, 
  FaArrowUp, 
  FaArrowDown, 
  FaUpload, 
  FaSave, 
  FaHome,
  FaSpinner,
  FaEye,
  FaEyeSlash
} from 'react-icons/fa';
import './HomeImagesAdmin.css';

const HomeImagesAdmin = () => {
  const [images, setImages] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [editingId, setEditingId] = useState(null);
  const [editForm, setEditForm] = useState({ 
    title: '', 
    subtitle: '' 
  });
  const [uploading, setUploading] = useState(false);
  const [submitStatus, setSubmitStatus] = useState(null);
  const [uploadForm, setUploadForm] = useState({
    title: 'Nueva Imagen',
    subtitle: 'Descripción de la imagen'
  });
  const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';

  useEffect(() => {
    const initializeComponent = async () => {
      try {
        setLoading(true);
        
        // Verificar si el usuario está autenticado
        const token = localStorage.getItem('token');
        if (!token) {
          throw new Error("No hay una sesión activa. Por favor inicia sesión.");
        }
        
        // Cargar imágenes
        await fetchImages();
      } catch (err) {
        console.error("Error de verificación o carga:", err);
        setError(err.message || "Error al verificar permisos o cargar datos.");
        setLoading(false);
      }
    };
    
    initializeComponent();
  }, []);

  const fetchImages = async () => {
    try {
      setLoading(true);
      console.log("Fetching images with token...");
      
      const response = await homeImagesService.getHomeImages();
      
      if (!response || !response.data || !response.data.data) {
        throw new Error("Formato de respuesta inválido");
      }
      
      console.log("Imágenes cargadas:", response.data.data);
      setImages(Array.isArray(response.data.data) ? response.data.data : []);
      setError(null);
    } catch (err) {
      console.error("Error al cargar imágenes:", err);
      setError(err.message || "Error al cargar las imágenes. Inténtalo de nuevo.");
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('¿Estás seguro de eliminar esta imagen?')) return;
    
    try {
      await homeImagesService.deleteHomeImage(id);
      setImages(prev => prev.filter(img => img.id !== id));
      setSubmitStatus({ type: 'success', message: 'Imagen eliminada correctamente' });
      setTimeout(() => setSubmitStatus(null), 3000);
    } catch (err) {
      console.error("Error al eliminar imagen:", err);
      setError(err.response?.data?.message || err.message || "Error al eliminar la imagen");
      setSubmitStatus({ type: 'error', message: `Error: ${err.response?.data?.message || err.message || "Error al eliminar la imagen"}` });
      setTimeout(() => setSubmitStatus(null), 3000);
    }
  };

  const handleEdit = (image) => {
    setEditingId(image.id);
    setEditForm({
      title: image.title || '',
      subtitle: image.subtitle || ''
    });
  };

  const handleEditChange = (e) => {
    const { name, value } = e.target;
    setEditForm(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const saveEdit = async (id) => {
    try {
      const response = await homeImagesService.updateHomeImage(id, editForm);
      
      // Actualizar imagen en el estado local
      setImages(prev => 
        prev.map(img => 
          img.id === id ? response.data.data : img
        )
      );
      setEditingId(null);
      setSubmitStatus({ type: 'success', message: 'Imagen actualizada correctamente' });
      setTimeout(() => setSubmitStatus(null), 3000);
    } catch (err) {
      console.error("Error al actualizar imagen:", err);
      setError(err.response?.data?.message || err.message || "Error al actualizar la imagen");
      setSubmitStatus({ type: 'error', message: `Error: ${err.response?.data?.message || err.message || "Error al actualizar la imagen"}` });
      setTimeout(() => setSubmitStatus(null), 3000);
    }
  };

  const toggleActive = async (id, currentStatus) => {
    try {
      const response = await homeImagesService.updateHomeImage(id, { 
        is_active: !currentStatus 
      });
      
      // Actualizar imagen en el estado local
      setImages(prev => 
        prev.map(img => 
          img.id === id ? response.data.data : img
        )
      );
      setSubmitStatus({ 
        type: 'success', 
        message: `Imagen ${!currentStatus ? 'activada' : 'desactivada'} correctamente` 
      });
      setTimeout(() => setSubmitStatus(null), 3000);
    } catch (err) {
      console.error("Error al cambiar estado de imagen:", err);
      setError(err.response?.data?.message || err.message || "Error al cambiar el estado de la imagen");
      setSubmitStatus({ type: 'error', message: `Error: ${err.response?.data?.message || err.message || "Error al cambiar el estado de la imagen"}` });
      setTimeout(() => setSubmitStatus(null), 3000);
    }
  };

  const moveImage = async (id, direction) => {
    const currentIndex = images.findIndex(img => img.id === id);
    const newIndex = direction === 'up' ? currentIndex - 1 : currentIndex + 1;
    
    if (newIndex < 0 || newIndex >= images.length) return;
    
    const newOrder = [...images];
    [newOrder[currentIndex], newOrder[newIndex]] = 
      [newOrder[newIndex], newOrder[currentIndex]];
    
    try {
      const orderData = newOrder.map((img, idx) => ({ id: img.id, order: idx }));
      console.log("Sending reorder data:", orderData);
      
      await homeImagesService.reorderHomeImages(orderData);
      
      // Actualizar orden en estado local
      setImages(newOrder);
      setSubmitStatus({ type: 'success', message: 'Orden actualizado correctamente' });
      setTimeout(() => setSubmitStatus(null), 3000);
    } catch (err) {
      console.error("Error al reordenar imágenes:", err);
      setError(err.response?.data?.message || err.message || "Error al reordenar las imágenes");
      setSubmitStatus({ type: 'error', message: `Error: ${err.response?.data?.message || err.message || "Error al reordenar las imágenes"}` });
      setTimeout(() => setSubmitStatus(null), 3000);
    }
  };

  const handleUploadFormChange = (e) => {
    const { name, value } = e.target;
    setUploadForm(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;
    
    if (!file.type.startsWith('image/')) {
      setError('El archivo debe ser una imagen');
      setSubmitStatus({ type: 'error', message: 'El archivo debe ser una imagen' });
      setTimeout(() => setSubmitStatus(null), 3000);
      return;
    }
    
    setUploading(true);
    try {
      // Usar el servicio para subir la imagen
      await homeImagesService.uploadHomeImage(file, uploadForm.title, uploadForm.subtitle);
      
      // Recargar todas las imágenes para obtener la lista actualizada
      await fetchImages();
      
      e.target.value = '';
      setSubmitStatus({ type: 'success', message: 'Imagen subida correctamente' });
      setTimeout(() => setSubmitStatus(null), 3000);
    } catch (err) {
      console.error("Error al subir imagen:", err);
      setError(err.response?.data?.message || err.message || "Error al subir la imagen");
      setSubmitStatus({ type: 'error', message: `Error: ${err.response?.data?.message || err.message || "Error al subir la imagen"}` });
      setTimeout(() => setSubmitStatus(null), 3000);
    } finally {
      setUploading(false);
      setUploadForm({ title: 'Nueva Imagen', subtitle: 'Descripción de la imagen' });
    }
  };

  const viewImageInHome = () => {
    window.open('/', '_blank');
  };

  // Si está cargando por más de lo normal, mostrar mensaje y opciones
  const [showRetryOption, setShowRetryOption] = useState(false);
  
  useEffect(() => {
    let timer;
    if (loading) {
      timer = setTimeout(() => {
        setShowRetryOption(true);
      }, 5000); // Mostrar opción de reintentar después de 5 segundos
    } else {
      setShowRetryOption(false);
    }
    
    return () => {
      clearTimeout(timer);
    };
  }, [loading]);

  if (loading) {
    return (
      <div className="loading-container">
        <FaSpinner className="spinner" />
        <p>Cargando imágenes...</p>
        {showRetryOption && (
          <div className="loading-options">
            <p>Está tomando más tiempo de lo normal</p>
            <button className="retry-button" onClick={fetchImages}>
              Reintentar carga
            </button>
            <button 
              className="retry-button"
              onClick={() => setLoading(false)}
            >
              Continuar sin cargar
            </button>
          </div>
        )}
      </div>
    );
  }

  return (
    <div className="home-images-admin-container">
      <header className="admin-header">
        <h1>
          <FaHome className="header-icon" />
          Administrar Imágenes del Inicio
        </h1>
        <p>Gestiona las imágenes que aparecen en el slider principal de la página de inicio</p>
      </header>

      {error && (
        <div className="error-container">
          <p className="error-message">Error: {error}</p>
          <button 
            className="retry-button"
            onClick={() => {
              setError(null);
              fetchImages();
            }}
          >
            Reintentar
          </button>
        </div>
      )}

      {submitStatus && (
        <div className={`status-message ${submitStatus.type}`}>
          {submitStatus.message}
        </div>
      )}

      <div className="admin-actions">
        <button 
          className="preview-button"
          onClick={viewImageInHome}
        >
          <FaEye /> Ver en Inicio
        </button>
        
        <div className="upload-form">
          <label className="upload-button">
            {uploading ? (
              <>
                <FaSpinner className="spinner" />
                Subiendo...
              </>
            ) : (
              <>
                <FaUpload />
                Subir Nueva Imagen
              </>
            )}
            <input
              type="file"
              accept="image/*"
              onChange={handleUpload}
              disabled={uploading}
            />
          </label>
          <input
            type="text"
            placeholder="Título de la imagen"
            name="title"
            value={uploadForm.title}
            onChange={handleUploadFormChange}
            disabled={uploading}
          />
          <input
            type="text"
            placeholder="Subtítulo de la imagen"
            name="subtitle"
            value={uploadForm.subtitle}
            onChange={handleUploadFormChange}
            disabled={uploading}
          />
        </div>
      </div>

      <section className="images-list">
        {images.length === 0 ? (
          <div className="empty-state">
            <p>No hay imágenes configuradas</p>
            <p>Sube imágenes para mostrarlas en el slider de la página de inicio</p>
          </div>
        ) : (
          images.map((image, index) => (
            <article 
              key={image.id} 
              className={`image-card ${!image.is_active ? 'inactive' : ''}`}
            >
              <div className="image-preview">
                <img
                  src={`${apiUrl}${image.image_url}`}
                  alt={image.title || "Imagen sin título"}
                  loading="lazy"
                  onError={(e) => {
                    console.error("Error loading image in admin:", image.image_url);
                    e.target.src = "https://via.placeholder.com/300x200?text=Error+de+imagen";
                  }}
                />
                {!image.is_active && (
                  <div className="inactive-overlay">
                    <FaEyeSlash />
                    <span>Inactiva</span>
                  </div>
                )}
              </div>

              <div className="image-details">
                {editingId === image.id ? (
                  <div className="edit-form">
                    <input
                      type="text"
                      name="title"
                      value={editForm.title}
                      onChange={handleEditChange}
                      placeholder="Título"
                    />
                    <input
                      type="text"
                      name="subtitle"
                      value={editForm.subtitle}
                      onChange={handleEditChange}
                      placeholder="Subtítulo"
                    />
                    <div className="edit-actions">
                      <button
                        className="save-button"
                        onClick={() => saveEdit(image.id)}
                      >
                        <FaSave /> Guardar
                      </button>
                      <button
                        className="cancel-button"
                        onClick={() => setEditingId(null)}
                      >
                        Cancelar
                      </button>
                    </div>
                  </div>
                ) : (
                  <>
                    <h3>{image.title || 'Sin título'}</h3>
                    <p>{image.subtitle || 'Sin descripción'}</p>
                  </>
                )}

                <div className="image-actions">
                  <button
                    className={`status-button ${image.is_active ? 'active' : 'inactive'}`}
                    onClick={() => toggleActive(image.id, image.is_active)}
                  >
                    {image.is_active ? (
                      <>
                        <FaEyeSlash /> Desactivar
                      </>
                    ) : (
                      <>
                        <FaEye /> Activar
                      </>
                    )}
                  </button>

                  <div className="order-controls">
                    <button
                      className="move-button"
                      onClick={() => moveImage(image.id, 'up')}
                      disabled={index === 0}
                      title="Mover arriba"
                    >
                      <FaArrowUp />
                    </button>
                    <button
                      className="move-button"
                      onClick={() => moveImage(image.id, 'down')}
                      disabled={index === images.length - 1}
                      title="Mover abajo"
                    >
                      <FaArrowDown />
                    </button>
                  </div>

                  <button
                    className="edit-button"
                    onClick={() => handleEdit(image)}
                    title="Editar"
                  >
                    <FaEdit />
                  </button>

                  <button
                    className="delete-button"
                    onClick={() => handleDelete(image.id)}
                    title="Eliminar"
                  >
                    <FaTrash />
                  </button>
                </div>
              </div>
            </article>
          ))
        )}
      </section>
    </div>
  );
};

export default HomeImagesAdmin;