import React, { useState, useEffect, useCallback } from 'react';
import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import ProfileService from '../../components/services/ProfileService';
import './Admin.css';
// Importamos los iconos necesarios
import { 
  FaChalkboardTeacher, 
  FaChartPie, 
  FaEnvelope, 
  FaBell, 
  FaUserCircle, 
  FaSignOutAlt, 
  FaSearch,
  FaSpinner,
  FaImages
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
        // Cargar imagen de perfil si existe
        if (response.data.user.image_url) {
          setProfileImage(`${apiUrl}${response.data.user.image_url}`);
        }
        
        // Cargar notificaciones
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
      // Aquí conectaríamos con el endpoint de notificaciones
      // Por ahora simulamos una carga
      setTimeout(() => {
        setNotifications(2); // Simular 2 notificaciones
      }, 500);
    } catch (error) {
      console.error('Error al cargar notificaciones:', error);
    }
  };

  // Verificar estado de administrador y cargar datos
  useEffect(() => {
    const verifyAdmin = async () => {
      setLoading(true);
      
      // Verificar que el usuario tenga permisos de administrador
      const isAdminUser = await checkAdminStatus();
      
      if (!isAdminUser) {
        // Si no es administrador, redirigir al inicio
        logout();
        navigate('/login', { state: { message: 'Acceso denegado. Se requieren privilegios de administrador.' } });
        return;
      }
      
      // Intentar refrescar el token si es necesario
      await refreshToken();
      
      // Cargar datos adicionales
      await loadProfileData();
    };
    
    verifyAdmin();
  }, [checkAdminStatus, logout, navigate, refreshToken, loadProfileData]);

  // Simulación de búsqueda
  const handleSearch = (e) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      // Implementación de búsqueda real
      setLoading(true);
      
      // Simular llamada a API
      setTimeout(() => {
        console.log(`Buscando: ${searchQuery}`);
        setLoading(false);
        
        // Aquí redirigirías a los resultados
        // navigate(`/admin/search?q=${encodeURIComponent(searchQuery)}`);
        
        // Reinicia el campo después de la búsqueda
        setSearchQuery('');
      }, 500);
    }
  };

  // Manejo de cierre de sesión
  const handleLogout = () => {
    setLoading(true);
    
    // Registrar última actividad antes de cerrar sesión
    try {
      // Aquí se podría hacer una llamada al backend para registrar el cierre de sesión
      // Por ahora simplemente esperamos un poco
      setTimeout(() => {
        logout();
        navigate('/login', { state: { message: 'Sesión cerrada correctamente' } });
      }, 300);
    } catch (error) {
      console.error('Error al cerrar sesión:', error);
      logout();
      navigate('/login');
    }
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
          <h1>Dashboard</h1>
          
          <div className="user-actions">
            <form className="search-bar" onSubmit={handleSearch}>
              {loading ? <FaSpinner className="spinner" /> : <FaSearch />}
              <input 
                type="text" 
                placeholder="Buscar..." 
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                disabled={loading}
              />
            </form>
            
            <div 
              className="notification-bell" 
              data-tooltip="Notificaciones"
              onClick={handleViewAllNotifications}
            >
              <FaBell />
              {notifications > 0 && (
                <span className="notification-indicator"></span>
              )}
            </div>
            
            <div className="user-profile" onClick={() => navigate('/admin/profile')}>
              {profileImage ? (
                <img 
                  src={profileImage} 
                  alt="Perfil" 
                  className="user-avatar-img"
                />
              ) : (
                <div className="user-avatar">
                  {user?.nombre?.charAt(0) || 'A'}
                </div>
              )}
              <div className="user-info">
                <span className="user-name">{user?.nombre || 'Admin'}</span>
                <span className="user-role">{user?.role === 'admin' ? 'Administrador' : 'Usuario'}</span>
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