import React, { useState, useEffect, useRef } from 'react';
import './portfolio.css';
import axios from 'axios';

function Portfolio() {
  const contentRefs = useRef([]);
  const [portfolioImages, setPortfolioImages] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';

  // Cargar proyectos del portfolio desde la API
  useEffect(() => {
    const fetchPortfolioProjects = async () => {
      try {
        setLoading(true);
        const response = await axios.get(`${apiUrl}/api/portfolio`);
        
        if (response.data.success) {
          const projects = (response.data.data || []).map(project => ({
            id: project.id,
            src: project.image_url.startsWith('http') 
              ? project.image_url 
              : `${apiUrl}${project.image_url}`,
            category: project.category,
            title: project.title,
            description: project.description || ''
          }));
          
          setPortfolioImages(projects);
        } else {
          setError('Error al cargar los proyectos del portfolio');
        }
      } catch (err) {
        console.error('Error al cargar los proyectos:', err);
        setError('Error de conexión al servidor');
      } finally {
        setLoading(false);
      }
    };
    
    fetchPortfolioProjects();
  }, [apiUrl]);

  // Obtener categoría de la URL si existe
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const categoryParam = params.get('category');
    if (categoryParam) {
      setSelectedCategory(categoryParam);
    }
  }, []);

  // Scroll animation effect
  useEffect(() => {
    const handleScroll = () => {
      contentRefs.current.forEach(el => {
        if (el) {
          const rect = el.getBoundingClientRect();
          const isVisible = rect.top < window.innerHeight * 0.85;
          if (isVisible) {
            el.classList.add('visible');
          }
        }
      });
    };

    contentRefs.current = Array.from(document.querySelectorAll('[data-scroll]'));
    
    window.addEventListener('scroll', handleScroll);
    handleScroll(); // Check initial state

    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, [portfolioImages]); // Ejecutar cuando se carguen las imágenes

  // Portfolio categories
  const categories = [
    { name: 'Arquitectura', value: 'architecture' },
    { name: 'Interiores', value: 'interiors' },
    { name: 'Paisajismo', value: 'landscape' },
    { name: 'Comercial', value: 'commercial' }
  ];

  const [selectedCategory, setSelectedCategory] = useState('all');

  const filteredProjects = selectedCategory === 'all' 
    ? portfolioImages 
    : portfolioImages.filter(img => img.category === selectedCategory);

  // Manejar cambio de categoría y actualizar URL
  const handleCategoryChange = (category) => {
    setSelectedCategory(category);
    
    // Actualizar URL sin recargar la página
    const url = new URL(window.location);
    if (category === 'all') {
      url.searchParams.delete('category');
    } else {
      url.searchParams.set('category', category);
    }
    window.history.pushState({}, '', url);
  };

  return (
    <div className="portfolio-page">
      <div className="portfolio-hero">
        <div className="portfolio-hero-content">
          <h1 
            className="portfolio-hero-title"
            data-scroll="fadeInUp"
          >
            Portafolio de Renders
          </h1>
          <p 
            className="portfolio-hero-subtitle"
            data-scroll="fadeInUp"
          >
            Explora nuestra colección de renders de alta calidad que dan vida a ideas y proyectos únicos.
          </p>
        </div>
      </div>

      <div className="portfolio-filters">
        <div 
          className="filter-container"
          data-scroll="fadeInUp"
        >
          <button 
            className={`filter-btn ${selectedCategory === 'all' ? 'active' : ''}`}
            onClick={() => handleCategoryChange('all')}
          >
            Todos
          </button>
          {categories.map(category => (
            <button
              key={category.value}
              className={`filter-btn ${selectedCategory === category.value ? 'active' : ''}`}
              onClick={() => handleCategoryChange(category.value)}
            >
              {category.name}
            </button>
          ))}
        </div>
      </div>

      {loading ? (
        <div className="portfolio-loading">
          <div className="loading-spinner"></div>
          <p>Cargando proyectos...</p>
        </div>
      ) : error ? (
        <div className="portfolio-error">
          <p>{error}</p>
          <button 
            className="filter-btn"
            onClick={() => window.location.reload()}
          >
            Reintentar
          </button>
        </div>
      ) : (
        <div 
          className="portfolio-gallery"
          data-scroll="fadeInUp"
        >
          {filteredProjects.length === 0 ? (
            <div className="portfolio-empty">
              <p>No hay proyectos disponibles en esta categoría.</p>
            </div>
          ) : (
            filteredProjects.map((project, index) => (
              <div 
                key={project.id || index} 
                className="portfolio-gallery-item"
              >
                <img 
                  src={project.src} 
                  alt={project.title} 
                  className="portfolio-gallery-image" 
                  onError={(e) => {
                    // Si la imagen falla al cargar, mostrar una imagen de respaldo
                    e.target.src = '/placeholder-image.jpg';
                    console.error(`Error al cargar imagen: ${project.src}`);
                  }}
                />
                <div className="portfolio-gallery-overlay">
                  <h3>{project.title}</h3>
                  <p>{categories.find(cat => cat.value === project.category)?.name || project.category}</p>
                  {project.description && <p className="portfolio-description">{project.description}</p>}
                </div>
              </div>
            ))
          )}
        </div>
      )}

      <div 
        className="portfolio-cta-section"
        data-scroll="fadeInUp"
      >
        <h2 className="portfolio-cta-title">Transforma Tus Ideas</h2>
        <p className="portfolio-cta-subtitle">
          ¿Tienes un proyecto en mente? Contáctanos y convirtamos tu visión en un render impresionante.
        </p>
        <div className="portfolio-cta-buttons">
          <a href="/contacto" className="hero-cta">Contáctanos</a>
          <a href="/servicios" className="hero-cta secondary">Conoce Nuestros Servicios</a>
        </div>
      </div>
    </div>
  );
}

export default Portfolio;