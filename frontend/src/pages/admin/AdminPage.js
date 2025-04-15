import React, { useState, useEffect, useCallback } from 'react';
import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import ProfileService from '../../components/services/ProfileService';
import './Admin.css';
import { 
  FaChalkboardTeacher, 
  FaChartPie, 
  FaEnvelope, 
  FaBell, 
  FaUserCircle, 
  FaSignOutAlt, 
  FaSearch,
  FaSpinner,
  FaImages,
  FaHome
} from 'react-icons/fa';

const AdminPage = () => {
  const navigate = useNavigate();
  const { user, logout, checkAdminStatus, refreshToken } = useAuth();
  const [notifications, setNotifications] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(false);
  const [profileImage, setProfileImage] = useState(null);
  const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';

  // Cargar datos del perfil
  const loadProfileData = useCallback(async () => {
    try {
      const response = await ProfileService.getProfile();
      
      if (response.data.success) {
        if (response.data.user.image_url) {
          setProfileImage(`${apiUrl}${response.data.user.image_url}`);
        }
        fetchNotifications();
      }
    } catch (error) {
      console.error('Error al cargar datos del perfil:', error);
    } finally {
      setLoading(false);
    }
  }, [apiUrl]);

  // Obtener notificaciones
  const fetchNotifications = async () => {
    try {
      setTimeout(() => {
        setNotifications(2);
      }, 500);
    } catch (error) {
      console.error('Error al cargar notificaciones:', error);
    }
  };

  // Verificar estado de administrador
  useEffect(() => {
    const verifyAdmin = async () => {
      setLoading(true);
      const isAdminUser = await checkAdminStatus();
      
      if (!isAdminUser) {
        logout();
        navigate('/login', { state: { message: 'Acceso denegado. Se requieren privilegios de administrador.' } });
        return;
      }
      
      await refreshToken();
      await loadProfileData();
    };
    
    verifyAdmin();
  }, [checkAdminStatus, logout, navigate, refreshToken, loadProfileData]);

  // Manejar búsqueda
  const handleSearch = (e) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      setLoading(true);
      setTimeout(() => {
        console.log(`Buscando: ${searchQuery}`);
        setLoading(false);
        setSearchQuery('');
      }, 500);
    }
  };

  // Cerrar sesión
  const handleLogout = () => {
    setLoading(true);
    setTimeout(() => {
      logout();
      navigate('/login', { state: { message: 'Sesión cerrada correctamente' } });
    }, 300);
  };

  // Ver todas las notificaciones
  const handleViewAllNotifications = () => {
    navigate('/admin/notifications');
  };

  return (
    <div className="admin-container">
      {loading && (
        <div className="loading-overlay">
          <FaSpinner className="spinner" />
        </div>
      )}
      
      <nav className="admin-sidebar">
        <h2>Panel Admin</h2>
        
        <div className="nav-section">
          <div className="nav-section-title">Gestión Principal</div>
          <NavLink to="dashboard" className={({ isActive }) => isActive ? 'active' : ''}>
            <FaHome />
            Dashboard
          </NavLink>
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
          <div className="nav-section-title">Contenido</div>
          <NavLink to="home-images" className={({ isActive }) => isActive ? 'active' : ''}>
            <FaImages />
            Slider Home
          </NavLink>
          <NavLink to="portfolio" className={({ isActive }) => isActive ? 'active' : ''}>
            <FaImages />
            Portfolio
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
          <h1>Panel de Administración</h1>
          
          <div className="user-actions">
            <form className="search-bar" onSubmit={handleSearch}>
              {loading ? <FaSpinner className="spinner" /> : <FaSearch />}
              <input 
                type="text" 
                placeholder="Buscar cursos, usuarios..." 
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                disabled={loading}
              />
            </form>
            
            <div 
              className="notification-bell" 
              onClick={handleViewAllNotifications}
              data-tooltip="Notificaciones no leídas"
            >
              <FaBell />
              {notifications > 0 && (
                <span className="notification-indicator">{notifications}</span>
              )}
            </div>
            
            <div 
              className="user-profile" 
              onClick={() => navigate('/admin/profile')}
              data-tooltip="Editar perfil"
            >
              {profileImage ? (
                <img 
                  src={profileImage} 
                  alt="Avatar del usuario" 
                  className="user-avatar-img"
                  onError={(e) => {
                    e.target.onerror = null;
                    setProfileImage(null);
                  }}
                />
              ) : (
                <div className="user-avatar-fallback">
                  {user?.nombre?.charAt(0).toUpperCase() || 'A'}
                </div>
              )}
              <div className="user-info">
                <span className="user-name">{user?.nombre || 'Administrador'}</span>
                <span className="user-email">{user?.email}</span>
              </div>
            </div>
          </div>
        </div>
        
        <div className="admin-main-content">
          <Outlet />
        </div>
      </div>
    </div>
  );
};

export default AdminPage;