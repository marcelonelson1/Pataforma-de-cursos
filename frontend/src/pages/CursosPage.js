import React from 'react';
import Cursos from '../components/Cursos';
import './CursosPage.css';

function CursosPage() {
  return (
    <div className="cursos-page">
      <h2>Nuestros Cursos</h2>
      <Cursos />
    </div>
  );
}

export default CursosPage;