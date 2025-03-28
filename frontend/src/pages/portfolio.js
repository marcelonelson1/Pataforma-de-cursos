import React, { useState, useEffect, useRef } from 'react';
import './portfolio.css';

// Import portfolio images
import portfolio1 from '../img/home/render1.jpeg';
import portfolio2 from '../img/home/render2.jpeg';
import portfolio3 from '../img/home/render3.jpeg';
import portfolio4 from '../img/home/render4.jpeg';
import portfolio5 from '../img/home/render5.jpeg';
import portfolio6 from '../img/portfolio/render6.jpeg';  // Assume you'll add this image
import portfolio7 from '../img/portfolio/render7.jpeg'; // Assume you'll add this image
import portfolio8 from '../img/portfolio/render8.jpeg'; // Assume you'll add this image

function Portfolio() {
  const contentRefs = useRef([]);

  // Scroll animation effectp
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
  }, []);

  // Portfolio categories
  const categories = [
    { name: 'Arquitectura', value: 'architecture' },
    { name: 'Interiores', value: 'interiors' },
    { name: 'Paisajismo', value: 'landscape' },
    { name: 'Comercial', value: 'commercial' }
  ];

  const [selectedCategory, setSelectedCategory] = useState('all');
  const portfolioImages = [
    { src: portfolio1, category: 'architecture', title: 'Proyecto Residencial Moderno' },
    { src: portfolio2, category: 'interiors', title: 'Diseño de Interiores Contemporáneo' },
    { src: portfolio3, category: 'landscape', title: 'Jardín Urbano Minimalista' },
    { src: portfolio4, category: 'commercial', title: 'Oficina Corporativa' },
    { src: portfolio5, category: 'architecture', title: 'Casa Ecológica' },
    { src: portfolio6, category: 'interiors', title: 'Sala de Estar Elegante' },
    { src: portfolio7, category: 'landscape', title: 'Terraza Panorámica' },
    { src: portfolio8, category: 'commercial', title: 'Centro Comercial' }
  ];

  const filteredProjects = selectedCategory === 'all' 
    ? portfolioImages 
    : portfolioImages.filter(img => img.category === selectedCategory);

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
            onClick={() => setSelectedCategory('all')}
          >
            Todos
          </button>
          {categories.map(category => (
            <button
              key={category.value}
              className={`filter-btn ${selectedCategory === category.value ? 'active' : ''}`}
              onClick={() => setSelectedCategory(category.value)}
            >
              {category.name}
            </button>
          ))}
        </div>
      </div>

      <div 
        className="portfolio-gallery"
        data-scroll="fadeInUp"
      >
        {filteredProjects.map((project, index) => (
          <div 
            key={index} 
            className="portfolio-gallery-item"
          >
            <img 
              src={project.src} 
              alt={project.title} 
              className="portfolio-gallery-image" 
            />
            <div className="portfolio-gallery-overlay">
              <h3>{project.title}</h3>
              <p>{project.category.charAt(0).toUpperCase() + project.category.slice(1)}</p>
            </div>
          </div>
        ))}
      </div>

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