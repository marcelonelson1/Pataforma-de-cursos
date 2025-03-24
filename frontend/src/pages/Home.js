import React from 'react';
import './Home.css';
import 'animate.css';

function Home() {
  return (
    <div className="home animate__animated animate__fadeIn animate__faster">
      <div className="bg-gradient"></div>
      <div className="bg-accent"></div>
      
      <div className="hero-container">
        <h1 className="hero-title animate__animated animate__fadeInUp animate__faster">Bienvenido a Mi Plataforma de Cursos</h1>
        <p className="hero-subtitle animate__animated animate__fadeInUp animate__faster animate__delay-05s">
          Desarrolla tus habilidades con cursos diseñados por expertos. Aprende a tu propio ritmo y conviértete en un profesional en tu campo.
        </p>
        
        <div className="hero-cta-container animate__animated animate__fadeInUp animate__faster animate__delay-1s">
          <a href="/cursos" className="hero-cta">Explorar Cursos</a>
          <a href="/registro" className="hero-cta secondary">Registrarse</a>
        </div>
        
        <div className="features-grid animate__animated animate__fadeInUp animate__faster animate__delay-15s">
          <div className="feature-card">
            <div className="feature-icon">📚</div>
            <h3 className="feature-title">Amplia Biblioteca</h3>
            <p className="feature-description">Accede a cientos de cursos actualizados en diversas áreas de conocimiento</p>
          </div>
          
          <div className="feature-card">
            <div className="feature-icon">👨‍🏫</div>
            <h3 className="feature-title">Profesores Expertos</h3>
            <p className="feature-description">Aprende con instructores certificados con años de experiencia en la industria</p>
          </div>
          
          <div className="feature-card">
            <div className="feature-icon">🏆</div>
            <h3 className="feature-title">Certificaciones</h3>
            <p className="feature-description">Obtén certificados validados que potenciarán tu currículum profesional</p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Home;