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

  // Detecta el scroll y el cambio de dimensiones de la ventana
  useEffect(() => {
    const handleScroll = () => {
      if (window.scrollY > 50) {
        setScrolled(true);
      } else {
        setScrolled(false);
      }
    };

    const handleResize = () => {
      if (window.innerWidth > 480 && menuOpen) {
        setMenuOpen(false);
      }
    };

    window.addEventListener('scroll', handleScroll);
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('scroll', handleScroll);
      window.removeEventListener('resize', handleResize);
    };
  }, [menuOpen]);

  const handleLogout = () => {
    logout();
    setMenuOpen(false);
    navigate('/');
  };

  const toggleMenu = () => {
    setMenuOpen(!menuOpen);
  };

  const closeMenu = () => {
    setMenuOpen(false);
  };

  return (
    <div className="header-container">
      {/* Logo centrado en la parte superior */}
      <div className="logo-container">
        <Link to="/" onClick={closeMenu}>
          <img src={logo} alt="Logo" className="logo" />
        </Link>
      </div>

      {/* Barra de navegación principal */}
      <header className={`main-header ${scrolled ? 'header-scrolled' : ''}`}>
        <div className="header-content">
          {/* Navegación central */}
          <nav className="main-nav">
            <ul className="nav-list">
              <li><Link to="/cursos" className="nav-item" onClick={closeMenu}>Cursos</Link></li>
              <li><Link to="/servicios" className="nav-item" onClick={closeMenu}>Servicios</Link></li>
              <li><Link to="/sobre-mi" className="nav-item" onClick={closeMenu}>Sobre Mí</Link></li>
              <li><Link to="/contacto" className="nav-item" onClick={closeMenu}>Contacto</Link></li>
            </ul>
          </nav>

          {/* Área de autenticación a la derecha */}
          <div className="auth-area">
            {isLoggedIn ? (
              <button onClick={handleLogout} className="auth-button">Cerrar Sesión</button>
            ) : (
              <>
                <Link to="/login" className="auth-button" onClick={closeMenu}>Iniciar Sesión</Link>
                <Link to="/registro" className="auth-button" onClick={closeMenu}>Registrarse</Link>
              </>
            )}
          </div>
        </div>

        {/* Botón para menú móvil */}
        <button 
          className="mobile-menu-button" 
          onClick={toggleMenu} 
          aria-label={menuOpen ? 'Cerrar menú' : 'Abrir menú'}
        >
          {menuOpen ? '✕' : '☰'}
        </button>

        {/* Menú móvil desplegable */}
        <div className={`mobile-nav ${menuOpen ? 'mobile-nav-active' : ''}`}>
          <Link to="/cursos" className="mobile-nav-item" onClick={closeMenu}>Cursos</Link>
          <Link to="/servicios" className="mobile-nav-item" onClick={closeMenu}>Servicios</Link>
          <Link to="/sobre-mi" className="mobile-nav-item" onClick={closeMenu}>Sobre Mí</Link>
          <Link to="/contacto" className="mobile-nav-item" onClick={closeMenu}>Contacto</Link>
          
          <div className="mobile-auth">
            {isLoggedIn ? (
              <button onClick={handleLogout} className="mobile-auth-button">Cerrar Sesión</button>
            ) : (
              <>
                <Link to="/login" className="mobile-auth-button" onClick={closeMenu}>Iniciar Sesión</Link>
                <Link to="/registro" className="mobile-auth-button" onClick={closeMenu}>Registrarse</Link>
              </>
            )}
          </div>
        </div>
      </header>
    </div>
  );
}

export default Header;