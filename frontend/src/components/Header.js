import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import logo from '../img/header/logohead.png';
import './Header.css';

function Header() {
  const { isLoggedIn, logout } = useAuth();
  const [menuOpen, setMenuOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth > 480 && menuOpen) {
        setMenuOpen(false);
      }
    };

    const handleScroll = () => {
      const offset = window.scrollY;
      if (offset > 50) {
        setScrolled(true);
      } else {
        setScrolled(false);
      }
    };

    window.addEventListener('resize', handleResize);
    window.addEventListener('scroll', handleScroll);
    
    return () => {
      window.removeEventListener('resize', handleResize);
      window.removeEventListener('scroll', handleScroll);
    };
  }, [menuOpen]);

  const handleLogout = () => {
    logout();
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
    <header className={`header ${scrolled ? 'scrolled' : ''}`}>
      <div className="logo">
        <Link to="/" onClick={closeMenu}>
          <img 
            src={logo} 
            alt="Logo" 
            className="logo-image" 
          />
        </Link>
      </div>
      
      {/* Desktop Navigation */}
      <div className="main-nav">
        <div className="nav-links">
          <Link to="/cursos" className="nav-link" onClick={closeMenu}>Cursos</Link>
          <Link to="/servicios" className="nav-link" onClick={closeMenu}>Servicios</Link>
          <Link to="/sobre-mi" className="nav-link" onClick={closeMenu}>Sobre Mí</Link>
          <Link to="/contacto" className="nav-link" onClick={closeMenu}>Contacto</Link>
        </div>
      </div>
      
      <div className="auth-nav">
        {isLoggedIn ? (
          <button onClick={handleLogout} className="auth-link">Cerrar Sesión</button>
        ) : (
          <>
            <Link to="/login" className="auth-link" onClick={closeMenu}>Iniciar Sesión</Link>
            <Link to="/registro" className="auth-link" onClick={closeMenu}>Registrarse</Link>
          </>
        )}
      </div>
      
      {/* Mobile Menu Toggle */}
      <button className="menu-toggle" onClick={toggleMenu} aria-label="Abrir menú">
        {menuOpen ? '✕' : '☰'}
      </button>
      
      {/* Mobile Menu */}
      <div className={`mobile-menu ${menuOpen ? 'active' : ''}`}>
        <Link to="/cursos" className="nav-link" onClick={closeMenu}>Cursos</Link>
        <Link to="/servicios" className="nav-link" onClick={closeMenu}>Servicios</Link>
        <Link to="/sobre-mi" className="nav-link" onClick={closeMenu}>Sobre Mí</Link>
        <Link to="/contactopage" className="nav-link" onClick={closeMenu}>Contacto</Link>
        
        <div className="mobile-menu-auth">
          {isLoggedIn ? (
            <button onClick={handleLogout} className="auth-link">Cerrar Sesión</button>
          ) : (
            <>
              <Link to="/login" className="auth-link" onClick={closeMenu}>Iniciar Sesión</Link>
              <Link to="/registro" className="auth-link" onClick={closeMenu}>Registrarse</Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
}

export default Header;