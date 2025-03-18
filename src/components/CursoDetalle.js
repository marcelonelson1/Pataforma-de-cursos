import React from 'react';
import { useParams } from 'react-router-dom';
import './CursoDetalle.css';

function CursoDetalle() {
  const { id } = useParams();
  const cursos = [
    { id: 1, titulo: 'Curso de React', descripcion: 'Aprende React desde cero.', contenido: 'Contenido del curso de React...' },
    { id: 2, titulo: 'Curso de Node.js', descripcion: 'Domina el backend con Node.js.', contenido: 'Contenido del curso de Node.js...' },
    { id: 3, titulo: 'Curso de Diseño Web', descripcion: 'Crea diseños modernos y responsivos.', contenido: 'Contenido del curso de Diseño Web...' },
  ];

  const curso = cursos.find((curso) => curso.id === parseInt(id));

  if (!curso) {
    return <div className="error">Curso no encontrado</div>;
  }

  return (
    <div className="curso-detalle animate__animated animate__fadeIn">
      <h2>{curso.titulo}</h2>
      <p>{curso.descripcion}</p>
      <div className="contenido">{curso.contenido}</div>
    </div>
  );
}

export default CursoDetalle;