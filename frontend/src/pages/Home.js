import React, { useState, useEffect, useRef } from 'react';
import './Home.css';

// Importa tus imágenes locales
import img1 from '../img/home/render7.jpeg';
import img2 from '../img/home/render8.jpeg';

function Home() {
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [scrolled, setScrolled] = useState(false);
  const sliderRef = useRef(null);
  const contentRefs = useRef([]);
  const images = [img1, img2];

  // Configuración del slider de imágenes
  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentImageIndex(prev => (prev + 1) % images.length);
    }, 3000);

    return () => clearInterval(interval);
  }, [images.length]);

  // Configuración del scroll observer
  useEffect(() => {
    const handleScroll = () => {
      // Efecto para el slider
      if (sliderRef.current) {
        const sliderPosition = sliderRef.current.getBoundingClientRect().top;
        setScrolled(sliderPosition < -50);
      }

      // Animaciones para elementos con data-scroll
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

    // Solo observamos elementos con data-scroll
    contentRefs.current = Array.from(document.querySelectorAll('[data-scroll]'));
    
    window.addEventListener('scroll', handleScroll);
    handleScroll(); // Comprobar estado inicial

    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, []);

  return (
    <div className="home">
      {/* Slider de imágenes con enfoque en parte inferior */}
      <div 
        ref={sliderRef}
        className={`image-slider ${scrolled ? 'scroll-down' : ''}`}
      >
        <div className="image-overlay"></div>
        {images.map((img, index) => (
          <img
            key={index}
            src={img}
            alt={`Imagen ${index + 1}`}
            className={`slider-image ${
              index === currentImageIndex ? 'active' : 'next'
            }`}
            loading="eager"
          />
        ))}
      </div>
      
      {/* Contenido principal */}
      <div className="hero-content">
        <div className="hero-container">
          {/* Título y subtítulo estáticos */}
          <h1 className="hero-title">
            RM RENDERS
          </h1>
          
          <p className="hero-subtitle">
          Transformamos tus ideas en renders de alta calidad, con realismo y detalle que dan vida a cada proyecto.
          </p>
          
          {/* Elementos con animación al scroll */}
          <div 
            className="hero-cta-container"
            data-scroll="fadeInUp"
          >
            <a href="/cursos" className="hero-cta">Explorar Cursos</a>
            <a href="/registro" className="hero-cta secondary">Prueba Gratis</a>
          </div>
          
          <div className="features-grid">
            <div 
              className="feature-card"
              data-scroll="fadeInLeft"
            >
              <div className="feature-icon">🎓</div>
              <h3 className="feature-title">Certificación</h3>
              <p className="feature-description">Obtén certificados reconocidos</p>
            </div>
            
            <div 
              className="feature-card"
              data-scroll="fadeInUp"
            >
              <div className="feature-icon">👨‍💻</div>
              <h3 className="feature-title">Proyectos Reales</h3>
              <p className="feature-description">Aprende con casos prácticos</p>
            </div>
            
            <div 
              className="feature-card"
              data-scroll="fadeInRight"
            >
              <div className="feature-icon">📱</div>
              <h3 className="feature-title">Acceso Móvil</h3>
              <p className="feature-description">Aprende desde cualquier dispositivo</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Home;