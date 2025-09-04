import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import './ServiciosPage.css';

// Importar íconos (usando los mismos que se ven en el Header.js)
// Aquí puedes reemplazar estos por los íconos que prefieras
// Por ejemplo, de react-icons, lucide-react o font-awesome

const ServiciosPage = () => {
  // Estado para animar elementos al hacer scroll
  const [isVisible, setIsVisible] = useState({});

  useEffect(() => {
    // Función para manejar las animaciones al hacer scroll
    const handleScroll = () => {
      const scrollElements = document.querySelectorAll('[data-scroll]');
      
      scrollElements.forEach((element) => {
        const elementTop = element.getBoundingClientRect().top;
        const elementHeight = element.offsetHeight;
        const windowHeight = window.innerHeight;
        
        // Si el elemento está visible en el viewport
        if (elementTop + elementHeight * 0.3 <= windowHeight) {
          element.classList.add('visible');
          
          // Actualizar estado para elementos específicos si es necesario
          if (element.id) {
            setIsVisible(prev => ({
              ...prev,
              [element.id]: true
            }));
          }
        }
      });
    };
    
    // Ejecutar una vez al cargar para animar elementos ya visibles
    handleScroll();
    
    // Agregar event listener
    window.addEventListener('scroll', handleScroll);
    
    // Limpiar
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, []);

  // Datos de los servicios
  const servicios = [
    {
      id: 'visualizacion-3d',
      icon: '🏙️',
      title: 'Visualización 3D',
      description: 'Transformamos tus planos en renders fotorrealistas que muestran cada detalle de tu proyecto arquitectónico con precisión y belleza.',
      link: '/contacto'
    },
    {
      id: 'animacion-arquitectonica',
      icon: '🎬',
      title: 'Animación Arquitectónica',
      description: 'Creamos recorridos virtuales y animaciones dinámicas que permiten explorar tu proyecto desde cualquier ángulo y perspectiva.',
      link: '/contacto'
    },
    {
      id: 'realidad-virtual',
      icon: '🥽',
      title: 'Realidad Virtual',
      description: 'Experimenta tus espacios antes de construirlos con nuestra tecnología de realidad virtual inmersiva para arquitectura.',
      link: '/contacto'
    },
    {
      id: 'renders-interiores',
      icon: '🛋️',
      title: 'Renders de Interiores',
      description: 'Visualizaciones fotorrealistas de espacios interiores con atención al detalle en iluminación, materiales y ambientación.',
      link: '/contacto'
    },
    {
      id: 'renders-exteriores',
      icon: '🏢',
      title: 'Renders de Exteriores',
      description: 'Imágenes de alto impacto visual que muestran la volumetría, materiales y entorno de tu proyecto arquitectónico.',
      link: '/contacto'
    },
    {
      id: 'maquetas-digitales',
      icon: '📐',
      title: 'Maquetas Digitales',
      description: 'Modelos 3D detallados que permiten visualizar el proyecto completo y su relación con el entorno urbano o natural.',
      link: '/contacto'
    }
  ];

  // Datos de los planes/precios
  const planes = [
    {
      id: 'basico',
      name: 'Básico',
      price: '$299',
      period: 'por imagen',
      features: [
        '2 renders exteriores',
        'Resolución estándar (1920x1080)',
        'Hasta 2 revisiones',
        'Entrega en 7 días hábiles',
        'Formato JPG/PNG'
      ],
      featured: false
    },
    {
      id: 'profesional',
      name: 'Profesional',
      price: '$499',
      period: 'por imagen',
      features: [
        '4 renders (exteriores e interiores)',
        'Alta resolución (4K)',
        'Hasta 4 revisiones',
        'Entrega en 5 días hábiles',
        'Formato JPG/PNG/TIFF',
        'Vista aérea incluida'
      ],
      featured: true
    },
    {
      id: 'premium',
      name: 'Premium',
      price: '$999',
      period: 'por proyecto',
      features: [
        '6 renders (exteriores e interiores)',
        'Resolución ultra HD (8K)',
        'Revisiones ilimitadas',
        'Entrega en 3 días hábiles',
        'Todos los formatos disponibles',
        'Recorrido virtual 360°',
        'Archivos editables'
      ],
      featured: false
    }
  ];

  // Proyectos para la galería
  const proyectos = [
    {
      id: 'proyecto1',
      title: 'Complejo Residencial Sky Gardens',
      image: 'placeholder-residencial.jpg' // Reemplazar con imágenes reales
    },
    {
      id: 'proyecto2',
      title: 'Centro Comercial Luminous Plaza',
      image: 'placeholder-comercial.jpg'
    },
    {
      id: 'proyecto3',
      title: 'Museo de Arte Contemporáneo',
      image: 'placeholder-museo.jpg'
    },
    {
      id: 'proyecto4',
      title: 'Residencia Privada Lakeside',
      image: 'placeholder-residencia.jpg'
    },
    {
      id: 'proyecto5',
      title: 'Hotel & Spa Mountain View',
      image: 'placeholder-hotel.jpg'
    },
    {
      id: 'proyecto6',
      title: 'Oficinas Corporativas Tech Hub',
      image: 'placeholder-oficinas.jpg'
    }
  ];

  // Testimonios
  const testimonios = [
    {
      id: 'testimonio1',
      text: 'Los renders que desarrollaron para nuestro proyecto residencial fueron claves para conseguir inversionistas. El nivel de detalle y realismo superó nuestras expectativas.',
      autor: 'Carlos Méndez',
      empresa: 'Constructora Altavista'
    },
    {
      id: 'testimonio2',
      text: 'Trabajar con este equipo ha sido una experiencia excepcional. Entendieron exactamente lo que necesitábamos y entregaron en tiempo récord renders de calidad impresionante.',
      autor: 'Ana Gómez',
      empresa: 'Grupo Arquitectónico Moderna'
    },
    {
      id: 'testimonio3',
      text: 'Estamos sumamente satisfechos con los resultados. Las visualizaciones 3D nos permitieron hacer ajustes cruciales antes de la construcción, ahorrando tiempo y recursos.',
      autor: 'Miguel Rodríguez',
      empresa: 'Desarrollo Inmobiliario Elite'
    },
    {
      id: 'testimonio4',
      text: 'La calidad de los renders fue fundamental para nuestras campañas de marketing. Definitivamente seguiremos contando con sus servicios para futuros proyectos.',
      autor: 'Laura Torres',
      empresa: 'Inmobiliaria Horizonte'
    }
  ];

  // Render del componente
  return (
    <div className="servicios-page">
      {/* Hero Section */}
      <section className="servicios-hero">
        <div className="hero-background"></div>
        <div className="hero-overlay"></div>
        <div className="servicios-hero-content" data-scroll="fadeInUp">
          <h1 className="servicios-hero-title">Servicios de Visualización Arquitectónica</h1>
          <p className="servicios-hero-subtitle">Transformamos tus ideas en experiencias visuales impactantes con nuestros servicios de renderizado 3D y visualización arquitectónica de alta calidad</p>
        </div>
      </section>

      {/* Introducción */}
      <section className="servicios-intro">
        <div className="servicios-intro-container" data-scroll="fadeInUp">
          <h2 className="servicios-intro-title">Renderes fotorrealistas para proyectos arquitectónicos</h2>
          <p className="servicios-intro-text">
            En nuestra empresa de visualización arquitectónica, nos especializamos en crear representaciones visuales de alta calidad que comunican la esencia y visión de tus proyectos arquitectónicos. Utilizamos las últimas tecnologías y software especializado para ofrecer renders fotorrealistas, animaciones y experiencias de realidad virtual que destacan por su atención al detalle, iluminación precisa y materialización realista.
          </p>
        </div>
      </section>

      {/* Servicios Principales */}
      <section className="servicios-grid">
        {servicios.map((servicio, index) => (
          <div 
            key={servicio.id} 
            className="servicio-card" 
            data-scroll={index % 2 === 0 ? "fadeInLeft" : "fadeInRight"}
          >
            <div className="servicio-icon">{servicio.icon}</div>
            <h3 className="servicio-title">{servicio.title}</h3>
            <p className="servicio-description">{servicio.description}</p>
            <Link to={servicio.link} className="servicio-link">
              Más información 
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M5 12h14"></path>
                <path d="M12 5l7 7-7 7"></path>
              </svg>
            </Link>
          </div>
        ))}
      </section>

      {/* Proceso de Trabajo */}
      <section className="proceso-section">
        <div className="proceso-container">
          <h2 className="proceso-title" data-scroll="fadeInUp">Nuestro Proceso de Trabajo</h2>
          <div className="proceso-steps">
            <div className="proceso-step" data-scroll="fadeInLeft">
              <div className="proceso-step-number">1</div>
              <div className="proceso-step-content">
                <h3 className="proceso-step-title">Briefing y Consulta</h3>
                <p className="proceso-step-description">
                  Comenzamos con una reunión detallada para entender tu visión, requerimientos y expectativas. Analizamos planos, referencias y cualquier material que puedas proporcionar para comprender completamente el proyecto.
                </p>
              </div>
            </div>
            
            <div className="proceso-step" data-scroll="fadeInRight">
              <div className="proceso-step-number">2</div>
              <div className="proceso-step-content">
                <h3 className="proceso-step-title">Modelado 3D</h3>
                <p className="proceso-step-description">
                  Creamos modelos tridimensionales precisos de tu proyecto utilizando software especializado. Cada elemento arquitectónico se desarrolla con atención al detalle siguiendo exactamente los planos proporcionados.
                </p>
              </div>
            </div>
            
            <div className="proceso-step" data-scroll="fadeInLeft">
              <div className="proceso-step-number">3</div>
              <div className="proceso-step-content">
                <h3 className="proceso-step-title">Materiales e Iluminación</h3>
                <p className="proceso-step-description">
                  Aplicamos texturas, materiales y configuramos la iluminación para crear ambientes realistas. Este paso es crucial para lograr el fotorrealismo que caracteriza nuestro trabajo.
                </p>
              </div>
            </div>
            
            <div className="proceso-step" data-scroll="fadeInRight">
              <div className="proceso-step-number">4</div>
              <div className="proceso-step-content">
                <h3 className="proceso-step-title">Renderizado y Post-producción</h3>
                <p className="proceso-step-description">
                  Procesamos las imágenes con potentes motores de renderizado y aplicamos técnicas avanzadas de post-producción para mejorar la calidad visual, ajustando color, contraste y añadiendo elementos ambientales.
                </p>
              </div>
            </div>
            
            <div className="proceso-step" data-scroll="fadeInLeft">
              <div className="proceso-step-number">5</div>
              <div className="proceso-step-content">
                <h3 className="proceso-step-title">Revisión y Entrega</h3>
                <p className="proceso-step-description">
                  Presentamos los resultados para tu revisión y realizamos los ajustes necesarios hasta tu completa satisfacción. Finalmente, entregamos los archivos en los formatos requeridos para su uso inmediato.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Planes y Precios */}
      <section className="planes-section">
        <div className="planes-container">
          <h2 className="planes-title" data-scroll="fadeInUp">Planes y Tarifas</h2>
          <p className="planes-subtitle" data-scroll="fadeInUp">
            Ofrecemos diferentes paquetes adaptados a las necesidades específicas de cada proyecto, desde visualizaciones básicas hasta completos paquetes de renders fotorrealistas.
          </p>
          
          <div className="planes-grid">
            {planes.map((plan) => (
              <div 
                key={plan.id} 
                className={`plan-card ${plan.featured ? 'featured' : ''}`}
                data-scroll="fadeInUp"
              >
                <h3 className="plan-name">{plan.name}</h3>
                <div className="plan-price">
                  {plan.price}<span> {plan.period}</span>
                </div>
                <ul className="plan-features">
                  {plan.features.map((feature, index) => (
                    <li key={index} className="plan-feature">
                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M20 6L9 17l-5-5"></path>
                      </svg>
                      {feature}
                    </li>
                  ))}
                </ul>
                <Link to="/contacto" className="plan-cta">
                  Solicitar Ahora
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M5 12h14"></path>
                    <path d="M12 5l7 7-7 7"></path>
                  </svg>
                </Link>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Galería de Trabajos */}
      <section className="galeria-section">
        <div className="galeria-container">
          <h2 className="galeria-title" data-scroll="fadeInUp">Proyectos Destacados</h2>
          <p className="galeria-subtitle" data-scroll="fadeInUp">
            Explora algunas de nuestras visualizaciones arquitectónicas más destacadas, donde materializamos ideas y proyectos a través de imágenes de alto impacto visual.
          </p>
          
          <div className="galeria-grid">
            {proyectos.map((proyecto) => (
              <div key={proyecto.id} className="galeria-item" data-scroll="fadeInUp">
                <img 
                  src={`/api/placeholder/600/400`} // Usamos placeholder temporalmente
                  alt={proyecto.title} 
                  className="galeria-image"
                />
                <div className="galeria-overlay">
                  <h4 className="galeria-item-title">{proyecto.title}</h4>
                </div>
              </div>
            ))}
          </div>
          
          <div className="galeria-ver-mas" data-scroll="fadeInUp">
            <Link to="/portfolio" className="galeria-btn">
              Ver Portafolio Completo
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M5 12h14"></path>
                <path d="M12 5l7 7-7 7"></path>
              </svg>
            </Link>
          </div>
        </div>
      </section>

      {/* Testimonios */}
      <section className="testimonios-section">
        <div className="testimonios-container">
          <h2 className="testimonios-title" data-scroll="fadeInUp">Lo que dicen nuestros clientes</h2>
          
          <div className="testimonios-grid">
            {testimonios.map((testimonio) => (
              <div key={testimonio.id} className="testimonio-card" data-scroll="fadeInUp">
                <div className="testimonio-content">
                  <p className="testimonio-text">{testimonio.text}</p>
                  <div className="testimonio-autor">
                    <div className="testimonio-avatar"></div>
                    <div className="testimonio-info">
                      <div className="testimonio-nombre">{testimonio.autor}</div>
                      <div className="testimonio-empresa">{testimonio.empresa}</div>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Final */}
      <section className="cta-section">
        <div className="cta-container" data-scroll="fadeInUp">
          <h2 className="cta-title">¿Listo para dar vida a tu proyecto?</h2>
          <p className="cta-text">
            Contacta con nosotros hoy mismo y descubre cómo nuestros servicios de visualización arquitectónica pueden ayudarte a comunicar tu visión de forma impactante y efectiva.
          </p>
          <div className="cta-buttons">
            <Link to="/contacto" className="cta-btn primary">
              Solicitar Presupuesto
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M5 12h14"></path>
                <path d="M12 5l7 7-7 7"></path>
              </svg>
            </Link>
            <Link to="/portfolio" className="cta-btn secondary">
              Ver Proyectos
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                <circle cx="8.5" cy="8.5" r="1.5"></circle>
                <path d="M21 15l-5-5L5 21"></path>
              </svg>
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
};

export default ServiciosPage;