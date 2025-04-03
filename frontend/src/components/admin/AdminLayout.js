import React, { useState } from 'react';
import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { 
  faTachometerAlt, faGraduationCap, faComments, 
  faChartBar, faSignOutAlt, faBars, faTimes, faUser
} from '@fortawesome/free-solid-svg-icons';
import { useAuth } from '../../context/AuthContext';
import './AdminLayout.css';

function AdminLayout() {
  const { logout, user } = useAuth();
  const navigate = useNavigate();
  const [sidebarOpen, setSidebarOpen] = useState(true);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  return (
    <div className={`admin-layout ${sidebarOpen ? '' : 'sidebar-collapsed'}`}>
      {/* Sidebar */}
      <div className="admin-sidebar">
        <div className="sidebar-header">
          <h2>Panel Admin</h2>
          <button className="sidebar-toggle" onClick={toggleSidebar}>
            <FontAwesomeIcon icon={sidebarOpen ? faTimes : faBars} />
          </button>
        </div>

        <div className="sidebar-user">
          <div className="user-avatar">
            <FontAwesomeIcon icon={faUser} />
          </div>
          {sidebarOpen && (
            <div className="user-info">
              <div className="user-name">{user?.nombre || 'Administrador'}</div>
              <div className="user-role">Administrador</div>
            </div>
          )}
        </div>

        <nav className="sidebar-nav">
          <ul>
            <li>
              <NavLink to="/admin" end>
                <FontAwesomeIcon icon={faTachometerAlt} />
                {sidebarOpen && <span>Dashboard</span>}
              </NavLink>
            </li>
            <li>
              <NavLink to="/admin/cursos">
                <FontAwesomeIcon icon={faGraduationCap} />
                {sidebarOpen && <span>Cursos</span>}
              </NavLink>
            </li>
            <li>
              <NavLink to="/admin/consultas">
                <FontAwesomeIcon icon={faComments} />
                {sidebarOpen && <span>Consultas</span>}
              </NavLink>
            </li>
            <li>
              <NavLink to="/admin/reportes">
                <FontAwesomeIcon icon={faChartBar} />
                {sidebarOpen && <span>Reportes</span>}
              </NavLink>
            </li>
          </ul>
        </nav>

        <div className="sidebar-footer">
          <button onClick={handleLogout} className="logout-btn">
            <FontAwesomeIcon icon={faSignOutAlt} />
            {sidebarOpen && <span>Cerrar Sesión</span>}
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="admin-main">
        <header className="admin-header">
          <div className="header-toggle">
            <button onClick={toggleSidebar}>
              <FontAwesomeIcon icon={faBars} />
            </button>
          </div>
          <div className="header-title">
            Panel de Administración
          </div>
          <div className="header-actions">
            <NavLink to="/" className="view-site-btn">
              Ver sitio
            </NavLink>
            <button onClick={handleLogout} className="logout-btn">
              <FontAwesomeIcon icon={faSignOutAlt} />
              <span>Cerrar Sesión</span>
            </button>
          </div>
        </header>

        <main className="admin-content">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

export default AdminLayout;