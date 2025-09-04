import React from 'react';
import './SobreMiPage.css';

function SobreMiPage() {
  const skills = [
    'JavaScript', 'React', 'Node.js', 'HTML5', 'CSS3', 
    'MongoDB', 'Express', 'Redux', 'Git', 'Responsive Design'
  ];

  return (
    <div className="sobre-mi-page animate__animated animate__fadeIn">
      <h2>Sobre Mí</h2>
      <div className="sobre-mi-content">
        <img
          src="https://via.placeholder.com/150"
          alt="Foto del creador"
          className="sobre-mi-image animate__animated animate__pulse animate__infinite animate__slow"
        />
        <p>
          Hola, soy Marcelo Nelson, el creador de esta plataforma. Soy un apasionado
          de la tecnología y la educación en entornos digitales, dedicado a compartir
          conocimientos de manera innovadora y accesible.
        </p>
        <p>
         
        </p>
        
        
        <div className="sobre-mi-skills">
          <h3>Tecnologías y Habilidades</h3>
          <div className="skill-badges">
            {skills.map((skill, index) => (
              <span key={index} className="skill-badge animate__animated animate__fadeInUp" 
                    style={{animationDelay: `${index * 0.1}s`}}>
                {skill}
              </span>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

export default SobreMiPage;