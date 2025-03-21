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
          Hola, soy [Tu Nombre], el creador de esta plataforma. Soy un apasionado
          de la tecnología y la educación en entornos digitales, dedicado a compartir
          conocimientos de manera innovadora y accesible.
        </p>
        <p>
          Con más de [X] años de experiencia en desarrollo web y educación online,
          he creado esta plataforma con un enfoque único que combina técnicas de 
          aprendizaje efectivas con una experiencia visual inmersiva.
        </p>
        <p>
          Mi objetivo es transformar la manera en que aprendemos tecnología, 
          haciendo que el proceso sea tan emocionante como gratificante. Cada curso
          ha sido diseñado meticulosamente para ofrecerte no solo conocimientos,
          sino una experiencia de aprendizaje completa.
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