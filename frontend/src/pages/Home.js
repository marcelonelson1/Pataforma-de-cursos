import React, { useState, useEffect, useRef } from 'react';
import './Home.css';
import './portfolio';

// Importa tus imágenes locales
import img1 from '../img/home/render7.jpeg';
import img2 from '../img/home/render8.jpeg';

// Nuevas imágenes para la sección de portafolio
import portfolio1 from '../img/home/render1.jpeg';
import portfolio2 from '../img/home/render2.jpeg';
import portfolio3 from '../img/home/render3.jpeg';

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

      {/* Sección de Portafolio mejor centrada */}
      <div className="portfolio-section">
        <div className="portfolio-container">
          <h2 
            className="portfolio-title"
            data-scroll="fadeInUp"
          >
            Nuestro Portafolio
          </h2>
          <p 
            className="portfolio-subtitle"
            data-scroll="fadeInUp"
          >
            Explora algunos de nuestros mejores renders y proyectos destacados
          </p>

          <div 
            className="portfolio-grid"
            data-scroll="fadeInUp"
          >
            {[portfolio1, portfolio2, portfolio3].map((img, index) => (
              <div key={index} className="portfolio-item">
                <img 
                  src={img} 
                  alt={`Proyecto ${index + 1}`} 
                  className="portfolio-image" 
                />
                <div className="portfolio-overlay">
                  <h4>Proyecto {index + 1}</h4>
                  <p>Render de alta calidad</p>
                </div>
              </div>
            ))}
          </div>

          <div 
            className="portfolio-cta"
            data-scroll="fadeInUp"
          >
            <a href="/portfolio" className="hero-cta">Ver Portafolio Completo</a>
          </div>
        </div>
      </div>

      {/* Sección de Testimonios */}
      <div className="testimonials-section">
        <div className="testimonials-container">
          <h2 
            className="testimonials-title"
            data-scroll="fadeInUp"
          >
            Lo Que Dicen Nuestros Clientes
          </h2>

          <div 
            className="testimonials-grid"
            data-scroll="fadeInUp"
          >
            <div className="testimonial-card">
              <p className="testimonial-text">
                "RM Renders transformó completamente mi presentación de proyecto. La calidad de los renders es impresionante."
              </p>
              <div className="testimonial-author">
                <span className="author-name">Carlos Mendoza</span>
                <span className="author-role">Arquitecto</span>
              </div>
            </div>

            <div className="testimonial-card">
              <p className="testimonial-text">
                "Increíble servicio y atención al detalle. Superaron mis expectativas en cada render."
              </p>
              <div className="testimonial-author">
                <span className="author-name">María Fernández</span>
                <span className="author-role">Diseñadora de Interiores</span>
              </div>
            </div>

            <div className="testimonial-card">
              <p className="testimonial-text">
                "Los renders de RM me ayudaron a cerrar importantes contratos. Altamente recomendados."
              </p>
              <div className="testimonial-author">
                <span className="author-name">Juan López</span>
                <span className="author-role">Constructor</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Home;