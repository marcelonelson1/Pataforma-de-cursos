import React, { useState } from 'react';
import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import './Admin.css';
// Importamos los iconos necesarios
import { FaChalkboardTeacher, FaChartPie, FaEnvelope, FaBell, FaUserCircle, FaSignOutAlt, FaSearch } from 'react-icons/fa';

const AdminPage = () => {
  const navigate = useNavigate();
  const [notifications] = useState(2); // Ejemplo de notificaciones
  const [searchQuery, setSearchQuery] = useState('');

  // Simulación de búsqueda
  const handleSearch = (e) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      // Aquí implementaríamos la búsqueda real
      console.log(`Buscando: ${searchQuery}`);
      // Reinicia el campo después de la búsqueda
      setSearchQuery('');
    }
  };

  // Simulación de logout
  const handleLogout = () => {
    // Aquí implementaríamos la lógica real de cierre de sesión
    console.log('Cerrando sesión...');
    // Redirigir al usuario a la página de inicio de sesión
    navigate('/login');
  };

  return (
    <div className="admin-container">
      <nav className="admin-sidebar">
        <h2>Panel Admin</h2>
        
        <div className="nav-section">
          <div className="nav-section-title">Gestión</div>
          <NavLink to="courses" className={({ isActive }) => isActive ? 'active' : ''}>
            <FaChalkboardTeacher />
            Cursos
          </NavLink>
          <NavLink to="stats" className={({ isActive }) => isActive ? 'active' : ''}>
            <FaChartPie />
            Estadísticas
          </NavLink>
          <NavLink to="messages" className={({ isActive }) => isActive ? 'active' : ''}>
            <FaEnvelope />
            Mensajes
            {notifications > 0 && (
              <span className="notification-badge">{notifications}</span>
            )}
          </NavLink>
        </div>

        <div className="nav-section">
          <div className="nav-section-title">Cuenta</div>
          <NavLink to="profile" className={({ isActive }) => isActive ? 'active' : ''}>
            <FaUserCircle />
            Mi Perfil
          </NavLink>
          <button onClick={handleLogout} className="logout-link">
            <FaSignOutAlt />
            Cerrar Sesión
          </button>
        </div>
      </nav>
      
      <div className="admin-content">
        <div className="admin-header">
          <h1>Dashboard</h1>
          
          <div className="user-actions">
            <form className="search-bar" onSubmit={handleSearch}>
              <FaSearch />
              <input 
                type="text" 
                placeholder="Buscar..." 
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </form>
            
            <div className="notification-bell" data-tooltip="Notificaciones">
              <FaBell />
              {notifications > 0 && (
                <span className="notification-indicator"></span>
              )}
            </div>
            
            <div className="user-profile">
              <div className="user-avatar">
                A
              </div>
              <div className="user-info">
                <span className="user-name">Admin</span>
                <span className="user-role">Administrador</span>
              </div>
            </div>
          </div>
        </div>
        
        <Outlet />
      </div>
    </div>
  );
};

export default AdminPage;