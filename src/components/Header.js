import React from 'react';
import { Link } from 'react-router-dom';
import './Header.css';

function Header() {
  return (
    <header className="header animate__animated animate__fadeInDown">
      <div className="logo">
        <Link to="/">Mi Plataforma de Cursos</Link>
      </div>
      <nav className="nav">
        <Link to="/cursos">Cursos</Link>
        <Link to="/sobre-mi">Sobre Mí</Link>
      </nav>
    </header>
  );
}

export default Header;