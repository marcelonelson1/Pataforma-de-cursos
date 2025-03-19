import React from 'react';
import './SobreMiPage.css';

function SobreMiPage() {
  return (
    <div className="sobre-mi-page animate__animated animate__fadeIn">
      <h2>Sobre Mí</h2>
      <div className="sobre-mi-content">
        <img
          src="https://via.placeholder.com/150" // Reemplaza con la URL de tu imagen
          alt="Foto del creador"
          className="sobre-mi-image"
        />
        <p>
          Hola, soy [Tu Nombre], el creador de esta plataforma. Soy un apasionado
          de la tecnología y la educación, y me encanta compartir mis
          conocimientos con los demás.
        </p>
        <p>
          Con más de [X] años de experiencia en el desarrollo web, he creado
          esta plataforma para ayudarte a aprender de manera efectiva y
          divertida.
        </p>
        <p>
          ¡Espero que disfrutes de los cursos y que te sean de gran utilidad!
        </p>
      </div>
    </div>
  );
}

export default SobreMiPage;