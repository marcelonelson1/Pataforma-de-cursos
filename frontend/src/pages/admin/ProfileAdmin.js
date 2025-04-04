import React, { useState, useEffect } from 'react';
import { 
  FaUser, 
  FaEnvelope, 
  FaKey, 
  FaBell, 
  FaShieldAlt, 
  FaCheckCircle,
  FaUserCog,
  FaUpload,
  FaSave
} from 'react-icons/fa';
import './ProfileAdmin.css';

const ProfileAdmin = () => {
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('profile');
  const [successMessage, setSuccessMessage] = useState('');
  const [profileImage, setProfileImage] = useState(null);
  
  // Datos de perfil simulados
  const [profileData, setProfileData] = useState({
    name: 'Admin Usuario',
    email: 'admin@ejemplo.com',
    role: 'Administrador',
    joinDate: '2022-01-15',
    phone: '+34 612 345 678',
    lastLogin: '2023-09-18 10:35:42'
  });
  
  // Configuración de notificaciones simulada
  const [notificationSettings, setNotificationSettings] = useState({
    emailNotifications: true,
    newMessages: true,
    newStudents: true,
    salesReports: true,
    systemUpdates: false
  });
  
  // Efecto para simular carga
  useEffect(() => {
    const timer = setTimeout(() => {
      setLoading(false);
    }, 800);
    return () => clearTimeout(timer);
  }, []);
  
  // Mostrar mensaje de éxito
  const showSuccessMessage = (message) => {
    setSuccessMessage(message);
    setTimeout(() => {
      setSuccessMessage('');
    }, 3000);
  };
  
  // Manejar cambio de imagen
  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (event) => {
        setProfileImage(event.target.result);
      };
      reader.readAsDataURL(file);
    }
  };
  
  // Manejar actualización de perfil
  const handleProfileUpdate = (e) => {
    e.preventDefault();
    setLoading(true);
    const formData = new FormData(e.target);
    
    // Simulamos una operación asincrónica
    setTimeout(() => {
      const updatedProfile = {
        name: formData.get('name'),
        email: formData.get('email'),
        phone: formData.get('phone'),
        role: profileData.role,
        joinDate: profileData.joinDate,
        lastLogin: profileData.lastLogin
      };
      
      setProfileData(updatedProfile);
      setLoading(false);
      showSuccessMessage('Perfil actualizado correctamente.');
    }, 800);
  };
  
  // Manejar cambio de contraseña
  const handlePasswordChange = (e) => {
    e.preventDefault();
    setLoading(true);
    
    // Simulamos una operación asincrónica
    setTimeout(() => {
      setLoading(false);
      showSuccessMessage('Contraseña actualizada correctamente.');
      e.target.reset();
    }, 800);
  };
  
  // Manejar actualización de notificaciones
  const handleNotificationChange = (setting) => {
    setNotificationSettings(prev => ({
      ...prev,
      [setting]: !prev[setting]
    }));
  };
  
  // Manejar guardado de notificaciones
  const handleSaveNotifications = () => {
    setLoading(true);
    
    // Simulamos una operación asincrónica
    setTimeout(() => {
      setLoading(false);
      showSuccessMessage('Preferencias de notificación actualizadas.');
    }, 600);
  };
  
  // Loading overlay
  const LoadingOverlay = () => (
    <div className="loading-overlay">
      <div className="spinner"></div>
    </div>
  );
  
  return (
    <div className="profile-admin">
      {loading && <LoadingOverlay />}
      
      {successMessage && (
        <div className="success-message">
          <FaCheckCircle />
          {successMessage}
        </div>
      )}
      
      <div className="section-header">
        <h2 className="section-title">Mi Perfil</h2>
      </div>
      
      <div className="profile-container">
        <div className="profile-sidebar">
          <div className="profile-image-container">
            {profileImage ? (
              <img src={profileImage} alt="Foto de perfil" className="profile-image" />
            ) : (
              <div className="profile-image-placeholder">
                {profileData.name.charAt(0).toUpperCase()}
              </div>
            )}
            <div className="profile-image-upload">
              <label htmlFor="profileImage" className="upload-label">
                <FaUpload /> Cambiar foto
              </label>
              <input
                id="profileImage"
                type="file"
                accept="image/*"
                onChange={handleImageChange}
                className="image-input"
              />
            </div>
          </div>
          
          <div className="profile-info">
            <h3>{profileData.name}</h3>
            <div className="role-badge">{profileData.role}</div>
            <p className="join-date">Miembro desde: {new Date(profileData.joinDate).toLocaleDateString()}</p>
          </div>
          
          <ul className="profile-tabs">
            <li 
              className={activeTab === 'profile' ? 'active' : ''}
              onClick={() => setActiveTab('profile')}
            >
              <FaUser /> Información Personal
            </li>
            <li 
              className={activeTab === 'security' ? 'active' : ''}
              onClick={() => setActiveTab('security')}
            >
              <FaShieldAlt /> Seguridad
            </li>
            <li 
              className={activeTab === 'notifications' ? 'active' : ''}
              onClick={() => setActiveTab('notifications')}
            >
              <FaBell /> Notificaciones
            </li>
          </ul>
        </div>
        
        <div className="profile-content">
          {activeTab === 'profile' && (
            <div className="tab-content">
              <h3 className="tab-title"><FaUserCog /> Editar Información Personal</h3>
              
              <form onSubmit={handleProfileUpdate} className="profile-form">
                <div className="form-group">
                  <label htmlFor="name">
                    <FaUser /> Nombre Completo
                  </label>
                  <input
                    type="text"
                    id="name"
                    name="name"
                    className="form-control"
                    defaultValue={profileData.name}
                    required
                  />
                </div>
                
                <div className="form-group">
                  <label htmlFor="email">
                    <FaEnvelope /> Correo Electrónico
                  </label>
                  <input
                    type="email"
                    id="email"
                    name="email"
                    className="form-control"
                    defaultValue={profileData.email}
                    required
                  />
                </div>
                
                <div className="form-group">
                  <label htmlFor="phone">
                    Teléfono
                  </label>
                  <input
                    type="tel"
                    id="phone"
                    name="phone"
                    className="form-control"
                    defaultValue={profileData.phone}
                  />
                </div>
                
                <div className="form-actions">
                  <button type="submit" className="btn-primary">
                    <FaSave /> Guardar Cambios
                  </button>
                </div>
              </form>
              
              <div className="account-info">
                <h4>Información de la Cuenta</h4>
                <ul>
                  <li>
                    <span className="info-label">Rol</span>
                    <span className="info-value">{profileData.role}</span>
                  </li>
                  <li>
                    <span className="info-label">Miembro desde</span>
                    <span className="info-value">{new Date(profileData.joinDate).toLocaleDateString()}</span>
                  </li>
                  <li>
                    <span className="info-label">Último acceso</span>
                    <span className="info-value">{profileData.lastLogin}</span>
                  </li>
                </ul>
              </div>
            </div>
          )}
          
          {activeTab === 'security' && (
            <div className="tab-content">
              <h3 className="tab-title"><FaKey /> Cambiar Contraseña</h3>
              
              <form onSubmit={handlePasswordChange} className="password-form">
                <div className="form-group">
                  <label htmlFor="currentPassword">Contraseña Actual</label>
                  <input
                    type="password"
                    id="currentPassword"
                    name="currentPassword"
                    className="form-control"
                    required
                  />
                </div>
                
                <div className="form-group">
                  <label htmlFor="newPassword">Nueva Contraseña</label>
                  <input
                    type="password"
                    id="newPassword"
                    name="newPassword"
                    className="form-control"
                    required
                  />
                  <div className="password-requirements">
                    <p>La contraseña debe cumplir los siguientes requisitos:</p>
                    <ul>
                      <li>Al menos 8 caracteres</li>
                      <li>Al menos una letra mayúscula</li>
                      <li>Al menos un número</li>
                      <li>Al menos un carácter especial</li>
                    </ul>
                  </div>
                </div>
                
                <div className="form-group">
                  <label htmlFor="confirmPassword">Confirmar Nueva Contraseña</label>
                  <input
                    type="password"
                    id="confirmPassword"
                    name="confirmPassword"
                    className="form-control"
                    required
                  />
                </div>
                
                <div className="form-actions">
                  <button type="submit" className="btn-primary">
                    <FaKey /> Actualizar Contraseña
                  </button>
                </div>
              </form>
              
              <div className="security-info">
                <h4>Seguridad de la Cuenta</h4>
                <div className="security-settings">
                  <div className="setting-item">
                    <div className="setting-content">
                      <h5>Verificación en dos pasos</h5>
                      <p>Añade una capa adicional de seguridad a tu cuenta.</p>
                    </div>
                    <label className="toggle-switch">
                      <input type="checkbox" />
                      <span className="toggle-slider"></span>
                    </label>
                  </div>
                  
                  <div className="setting-item">
                    <div className="setting-content">
                      <h5>Sesiones activas</h5>
                      <p>Revisa y cierra sesiones abiertas en otros dispositivos.</p>
                    </div>
                    <button className="btn-text">Ver sesiones</button>
                  </div>
                </div>
              </div>
            </div>
          )}
          
          {activeTab === 'notifications' && (
            <div className="tab-content">
              <h3 className="tab-title"><FaBell /> Preferencias de Notificación</h3>
              
              <div className="notification-settings">
                <div className="setting-group">
                  <h4>General</h4>
                  
                  <div className="setting-item">
                    <div className="setting-content">
                      <h5>Notificaciones por correo electrónico</h5>
                      <p>Recibir notificaciones generales por correo.</p>
                    </div>
                    <label className="toggle-switch">
                      <input 
                        type="checkbox" 
                        checked={notificationSettings.emailNotifications}
                        onChange={() => handleNotificationChange('emailNotifications')}
                      />
                      <span className="toggle-slider"></span>
                    </label>
                  </div>
                </div>
                
                <div className="setting-group">
                  <h4>Plataforma</h4>
                  
                  <div className="setting-item">
                    <div className="setting-content">
                      <h5>Nuevos mensajes</h5>
                      <p>Recibir notificaciones cuando lleguen nuevos mensajes.</p>
                    </div>
                    <label className="toggle-switch">
                      <input 
                        type="checkbox" 
                        checked={notificationSettings.newMessages}
                        onChange={() => handleNotificationChange('newMessages')}
                      />
                      <span className="toggle-slider"></span>
                    </label>
                  </div>
                  
                  <div className="setting-item">
                    <div className="setting-content">
                      <h5>Nuevos estudiantes</h5>
                      <p>Notificar cuando nuevos estudiantes se registren en la plataforma.</p>
                    </div>
                    <label className="toggle-switch">
                      <input 
                        type="checkbox" 
                        checked={notificationSettings.newStudents}
                        onChange={() => handleNotificationChange('newStudents')}
                      />
                      <span className="toggle-slider"></span>
                    </label>
                  </div>
                  
                  <div className="setting-item">
                    <div className="setting-content">
                      <h5>Informes de ventas</h5>
                      <p>Recibir informes periódicos de ventas y métricas.</p>
                    </div>
                    <label className="toggle-switch">
                      <input 
                        type="checkbox" 
                        checked={notificationSettings.salesReports}
                        onChange={() => handleNotificationChange('salesReports')}
                      />
                      <span className="toggle-slider"></span>
                    </label>
                  </div>
                  
                  <div className="setting-item">
                    <div className="setting-content">
                      <h5>Actualizaciones del sistema</h5>
                      <p>Notificar sobre nuevas características y actualizaciones.</p>
                    </div>
                    <label className="toggle-switch">
                      <input 
                        type="checkbox" 
                        checked={notificationSettings.systemUpdates}
                        onChange={() => handleNotificationChange('systemUpdates')}
                      />
                      <span className="toggle-slider"></span>
                    </label>
                  </div>
                </div>
                
                <div className="form-actions">
                  <button 
                    className="btn-primary"
                    onClick={handleSaveNotifications}
                  >
                    <FaSave /> Guardar Preferencias
                  </button>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProfileAdmin;