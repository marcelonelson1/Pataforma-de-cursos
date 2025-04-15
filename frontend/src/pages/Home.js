import React, { useState, useEffect, useRef } from 'react';
import homeImagesService from '../components/services/HomeImagesService';
import './Home.css';
import './portfolio';

// Mantenemos las imágenes del portfolio como importaciones estáticas por ahora
import portfolio1 from '../img/home/render1.jpeg';
import portfolio2 from '../img/home/render2.jpeg';
import portfolio3 from '../img/home/render3.jpeg';

// Imágenes de fallback (se usarán si no hay imágenes disponibles en la API)
import imgFallback1 from '../img/home/render7.jpeg';
import imgFallback2 from '../img/home/render8.jpeg';

function Home() {
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [scrolled, setScrolled] = useState(false);
  const [sliderImages, setSliderImages] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const sliderRef = useRef(null);
  const contentRefs = useRef([]);
  const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
  
  // Obtener las imágenes del Home desde la API
  useEffect(() => {
    const fetchHomeImages = async () => {
      try {
        setIsLoading(true);
        console.log("Fetching public images...");
        
        // Implementar timeout para evitar carga infinita
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 10000);
        
        const response = await homeImagesService.getPublicHomeImages();
        
        clearTimeout(timeoutId);
        
        if (!response || !response.data || !response.data.data) {
          console.log("No valid data in response, using fallback images");
          throw new Error('No se pudieron cargar las imágenes');
        }
        
        const data = response.data;
        console.log("Received data:", data);
        
        // Usar imágenes fallback si no hay imágenes o datos inválidos
        if (!data.data || !Array.isArray(data.data) || data.data.length === 0) {
          console.log("No valid images found, using fallback images");
          setSliderImages([
            { url: imgFallback1, title: "RM RENDERS", subtitle: "Calidad y detalle" },
            { url: imgFallback2, title: "RM RENDERS", subtitle: "Diseños impactantes" }
          ]);
        } else {
          // Formatear las imágenes para su uso en el slider
          const formattedImages = data.data.map(img => ({
            url: `${apiUrl}${img.image_url}`,
            title: img.title || "RM RENDERS",
            subtitle: img.subtitle || "Diseños de alta calidad"
          }));
          console.log("Formatted images:", formattedImages);
          setSliderImages(formattedImages);
        }
      } catch (err) {
        console.error("Error al cargar imágenes:", err);
        setError(err.message);
        // Cargar imágenes de fallback en caso de error
        setSliderImages([
          { url: imgFallback1, title: "RM RENDERS", subtitle: "Calidad y detalle" },
          { url: imgFallback2, title: "RM RENDERS", subtitle: "Diseños impactantes" }
        ]);
      } finally {
        setIsLoading(false);
      }
    };

    fetchHomeImages();
  }, [apiUrl]);

  // Configuración del slider de imágenes
  useEffect(() => {
    if (sliderImages.length === 0) return;
    
    const interval = setInterval(() => {
      setCurrentImageIndex(prev => (prev + 1) % sliderImages.length);
    }, 3000);

    return () => clearInterval(interval);
  }, [sliderImages.length]);

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

  // Si hay un error de carga, mostrar imágenes de fallback
  useEffect(() => {
    if (error && sliderImages.length === 0) {
      setSliderImages([
        { url: imgFallback1, title: "RM RENDERS", subtitle: "Calidad y detalle" },
        { url: imgFallback2, title: "RM RENDERS", subtitle: "Diseños impactantes" }
      ]);
      setIsLoading(false);
    }
  }, [error, sliderImages.length]);

  return (
    <div className="home">
      {/* Slider de imágenes con enfoque en parte inferior */}
      <div 
        ref={sliderRef}
        className={`image-slider ${scrolled ? 'scroll-down' : ''}`}
      >
        <div className="image-overlay"></div>
        {isLoading ? (
          <div className="slider-loading">Cargando imágenes...</div>
        ) : sliderImages.length > 0 ? (
          sliderImages.map((img, index) => (
            <img
              key={index}
              src={img.url}
              alt={img.title || `Imagen ${index + 1}`}
              className={`slider-image ${
                index === currentImageIndex ? 'active' : 'next'
              }`}
              loading={index === 0 ? "eager" : "lazy"}
              onError={(e) => {
                console.error("Error loading image:", img.url);
                e.target.src = index % 2 === 0 ? imgFallback1 : imgFallback2;
              }}
            />
          ))
        ) : (
          // Mostrar imágenes de fallback si no hay imágenes disponibles
          <>
            <img
              src={imgFallback1}
              alt="RM RENDERS"
              className={`slider-image ${currentImageIndex === 0 ? 'active' : 'next'}`}
              loading="eager"
            />
            <img
              src={imgFallback2}
              alt="RM RENDERS"
              className={`slider-image ${currentImageIndex === 1 ? 'active' : 'next'}`}
              loading="lazy"
            />
          </>
        )}
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