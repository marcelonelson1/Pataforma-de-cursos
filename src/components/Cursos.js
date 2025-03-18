import React from 'react';
import { Link } from 'react-router-dom';
import './Cursos.css';

function Cursos() {
  const cursos = [
    { id: 1, titulo: 'Curso de React', descripcion: 'Aprende React desde cero.' },
    { id: 2, titulo: 'Curso de Node.js', descripcion: 'Domina el backend con Node.js.' },
    { id: 3, titulo: 'Curso de Diseño Web', descripcion: 'Crea diseños modernos y responsivos.' },
  ];

  return (
    <ul className="cursos-list">
      {cursos.map((curso) => (
        <li key={curso.id} className="curso-item animate__animated animate__fadeInUp">
          <Link to={`/curso/${curso.id}`} className="curso-link">
            {curso.titulo}
          </Link>
          <p className="curso-descripcion">{curso.descripcion}</p>
        </li>
      ))}
    </ul>
  );
}

export default Cursos;