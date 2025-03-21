import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './Header.css';

function Header() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem('token');
    setIsLoggedIn(!!token);

    const handleResize = () => {
      if (window.innerWidth > 480 && menuOpen) {
        setMenuOpen(false);
      }
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [menuOpen]);

  const handleLogout = () => {
    localStorage.removeItem('token');
    setIsLoggedIn(false);
    setMenuOpen(false);
    navigate('/', { replace: true });
  };

  const toggleMenu = () => {
    setMenuOpen(!menuOpen);
  };

  const closeMenu = () => {
    if (menuOpen) setMenuOpen(false);
  };

  return (
    <header className="header animate__animated animate__fadeInDown">
      <div className="logo">
        <Link to="/" onClick={closeMenu}>Mi Plataforma de Cursos</Link>
      </div>

      <button className="menu-toggle" onClick={toggleMenu} aria-label="Abrir menú">
        {menuOpen ? '✕' : '☰'}
      </button>

      <nav className={`nav ${menuOpen ? 'active' : ''}`}>
        <Link to="/cursos" className="nav-link" onClick={closeMenu}>Cursos</Link>
        <Link to="/sobre-mi" className="nav-link" onClick={closeMenu}>Sobre Mí</Link>
        {isLoggedIn ? (
          <button onClick={handleLogout} className="auth-link">Cerrar Sesión</button>
        ) : (
          <>
            <Link to="/login" className="auth-link" onClick={closeMenu}>Iniciar Sesión</Link>
            <Link to="/registro" className="auth-link" onClick={closeMenu}>Registrarse</Link>
          </>
        )}
      </nav>
    </header>
  );
}

export default Header;