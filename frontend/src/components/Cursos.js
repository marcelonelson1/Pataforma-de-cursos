import React from 'react';
import { Link } from 'react-router-dom';
import './Cursos.css';

function Cursos() {
  const cursos = [
    { id: 1, titulo: 'Curso de React', descripcion: 'Aprende React desde cero.', imagen: 'https://via.placeholder.com/300', duracion: '10 horas', nivel: 'Principiante' },
    { id: 2, titulo: 'Curso de Node.js', descripcion: 'Domina el backend con Node.js.', imagen: 'https://via.placeholder.com/300', duracion: '15 horas', nivel: 'Intermedio' },
    { id: 3, titulo: 'Curso de Diseño Web', descripcion: 'Crea diseños modernos y responsivos.', imagen: 'https://via.placeholder.com/300', duracion: '12 horas', nivel: 'Principiante' },
  ];

  return (
    <div className="cursos-container">
      <ul className="cursos-list">
        {cursos.map((curso) => (
          <li key={curso.id} className="curso-item animate__animated animate__fadeInUp">
            <div className="curso-image-container">
              <img src={curso.imagen} alt={curso.titulo} className="curso-image" />
              <span className="curso-badge">{curso.nivel}</span>
            </div>
            <Link to={`/curso/${curso.id}`} className="curso-link">
              {curso.titulo}
            </Link>
            <p className="curso-descripcion">{curso.descripcion}</p>
            <div className="curso-footer">
              <div className="curso-meta">
                <span className="curso-meta-item">
                  <i className="curso-meta-icon">⏳</i> {curso.duracion}
                </span>
              </div>
              <button className="curso-cta">Más Información</button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default Cursos;