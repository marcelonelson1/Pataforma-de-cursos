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
  FaSave,
  FaExclamationTriangle
} from 'react-icons/fa';
import { useAuth } from '../../context/AuthContext';
import ProfileService from '../../components/services/ProfileService';
import './ProfileAdmin.css';

const ProfileAdmin = () => {
  const { user, isLoggedIn, updateUser } = useAuth();
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('profile');
  const [successMessage, setSuccessMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [profileImage, setProfileImage] = useState(null);
  
  // Datos de perfil
  const [profileData, setProfileData] = useState({
    name: '',
    email: '',
    role: '',
    joinDate: '',
    phone: '',
    lastLogin: ''
  });
  
  // Configuración de notificaciones
  const [notificationSettings, setNotificationSettings] = useState({
    emailNotifications: true,
    newMessages: true,
    newStudents: true,
    salesReports: true,
    systemUpdates: false
  });
  
  // Cargar datos del perfil
  useEffect(() => {
    if (isLoggedIn && user) {
      setLoading(true);
      
      // Cargar perfil desde el contexto primero
      const userData = {
        name: user.nombre || 'Usuario',
        email: user.email || '',
        role: user.role === 'admin' ? 'Administrador' : 'Usuario',
        joinDate: formatDate(user.created_at) || new Date().toISOString().split('T')[0],
        phone: user.phone || '',
        lastLogin: formatDateTime(user.last_login || new Date())
      };
      
      setProfileData(userData);
      
      // Luego cargar datos actualizados del backend
      ProfileService.getProfile()
        .then(response => {
          if (response.data.success) {
            const updatedData = {
              name: response.data.user.nombre || userData.name,
              email: response.data.user.email || userData.email,
              role: response.data.user.role === 'admin' ? 'Administrador' : 'Usuario',
              joinDate: formatDate(response.data.user.created_at) || userData.joinDate,
              phone: response.data.user.phone || '',
              lastLogin: formatDateTime(response.data.user.last_login || userData.lastLogin)
            };
            
            setProfileData(updatedData);
            
            // Si hay una imagen de perfil, establecerla
            if (response.data.user.image_url) {
              const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
              setProfileImage(`${apiUrl}${response.data.user.image_url}`);
            }
          }
        })
        .catch(error => {
          console.error('Error al cargar perfil:', error);
          // No mostramos error porque ya estamos usando datos del contexto
        })
        .finally(() => {
          // Cargar preferencias de notificación
          fetchNotificationSettings();
        });
    }
  }, [isLoggedIn, user]);
  
  // Formatear fecha
  const formatDate = (dateString) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toISOString().split('T')[0];
  };
  
  // Formatear fecha y hora
  const formatDateTime = (date) => {
    if (!date) return '';
    if (typeof date === 'string') {
      date = new Date(date);
    }
    return date.toISOString().replace('T', ' ').substring(0, 19);
  };
  
  // Obtener preferencias de notificación
  const fetchNotificationSettings = async () => {
    try {
      const response = await ProfileService.getNotificationSettings();
      
      if (response.data.success) {
        setNotificationSettings(response.data.settings);
      }
      
      setLoading(false);
    } catch (error) {
      console.error('Error al cargar preferencias:', error);
      setErrorMessage('No se pudieron cargar las preferencias de notificación');
      setLoading(false);
    }
  };
  
  // Mostrar mensaje de éxito
  const showSuccessMessage = (message) => {
    setSuccessMessage(message);
    setErrorMessage('');
    setTimeout(() => {
      setSuccessMessage('');
    }, 3000);
  };
  
  // Mostrar mensaje de error
  const showErrorMessage = (message) => {
    setErrorMessage(message);
    setSuccessMessage('');
    setTimeout(() => {
      setErrorMessage('');
    }, 5000);
  };
  
  // Manejar cambio de imagen
  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      if (file.size > 2 * 1024 * 1024) {
        showErrorMessage('La imagen no debe superar los 2MB');
        return;
      }
      
      const reader = new FileReader();
      reader.onload = (event) => {
        setProfileImage(event.target.result);
      };
      reader.readAsDataURL(file);
      
      // Subir imagen al servidor
      uploadProfileImage(file);
    }
  };
  
  // Subir imagen de perfil
  const uploadProfileImage = async (file) => {
    try {
      setLoading(true);
      
      const formData = new FormData();
      formData.append('image', file);
      
      const response = await ProfileService.uploadProfileImage(formData);
      
      if (response.data.success) {
        showSuccessMessage('Imagen de perfil actualizada');
        
        // Actualizar URL de imagen en el contexto
        if (updateUser && typeof updateUser === 'function') {
          updateUser({ image_url: response.data.imageUrl });
        }
      } else {
        throw new Error(response.data.error || 'Error al subir la imagen');
      }
      
    } catch (error) {
      console.error('Error:', error);
      showErrorMessage('No se pudo subir la imagen. Intente de nuevo.');
    } finally {
      setLoading(false);
    }
  };
  
  // Manejar actualización de perfil
  const handleProfileUpdate = async (e) => {
    e.preventDefault();
    setLoading(true);
    
    try {
      const updateData = {
        nombre: e.target.name.value,
        email: e.target.email.value,
        phone: e.target.phone.value
      };
      
      const response = await ProfileService.updateProfile(updateData);
      
      if (response.data.success) {
        // Actualizar datos locales
        setProfileData({
          ...profileData,
          name: updateData.nombre,
          email: updateData.email,
          phone: updateData.phone
        });
        
        // Actualizar usuario en el contexto
        if (updateUser && typeof updateUser === 'function') {
          updateUser({
            nombre: updateData.nombre,
            email: updateData.email,
            phone: updateData.phone
          });
        }
        
        showSuccessMessage('Perfil actualizado correctamente');
      } else {
        throw new Error(response.data.error || 'Error al actualizar el perfil');
      }
      
    } catch (error) {
      console.error('Error:', error);
      showErrorMessage(error.response?.data?.error || 'No se pudo actualizar el perfil. Verifique sus datos.');
    } finally {
      setLoading(false);
    }
  };
  
  // Manejar cambio de contraseña
  const handlePasswordChange = async (e) => {
    e.preventDefault();
    setLoading(true);
    
    const currentPassword = e.target.currentPassword.value;
    const newPassword = e.target.newPassword.value;
    const confirmPassword = e.target.confirmPassword.value;
    
    // Validar que las contraseñas coincidan
    if (newPassword !== confirmPassword) {
      showErrorMessage('Las contraseñas no coinciden');
      setLoading(false);
      return;
    }
    
    // Validar requisitos de la contraseña
    const passwordRegex = /^(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;
    if (!passwordRegex.test(newPassword)) {
      showErrorMessage('La contraseña no cumple con los requisitos de seguridad');
      setLoading(false);
      return;
    }
    
    try {
      const passwordData = {
        currentPassword,
        newPassword
      };
      
      const response = await ProfileService.changePassword(passwordData);
      
      if (response.data.success) {
        showSuccessMessage('Contraseña actualizada correctamente');
        e.target.reset();
      } else {
        throw new Error(response.data.error || 'Error al cambiar la contraseña');
      }
      
    } catch (error) {
      console.error('Error:', error);
      showErrorMessage(error.response?.data?.error || 'No se pudo actualizar la contraseña. Verifique que la contraseña actual sea correcta.');
    } finally {
      setLoading(false);
    }
  };
  
  // Manejar actualización de notificaciones
  const handleNotificationChange = (setting) => {
    setNotificationSettings(prev => ({
      ...prev,
      [setting]: !prev[setting]
    }));
  };
  
  // Manejar guardado de notificaciones
  const handleSaveNotifications = async () => {
    setLoading(true);
    
    try {
      const response = await ProfileService.updateNotificationSettings(notificationSettings);
      
      if (response.data.success) {
        showSuccessMessage('Preferencias de notificación actualizadas');
      } else {
        throw new Error(response.data.error || 'Error al guardar preferencias');
      }
      
    } catch (error) {
      console.error('Error:', error);
      showErrorMessage(error.response?.data?.error || 'No se pudieron guardar las preferencias de notificación');
    } finally {
      setLoading(false);
    }
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
      
      {errorMessage && (
        <div className="error-message">
          <FaExclamationTriangle />
          {errorMessage}
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
            <p className="join-date">Miembro desde: {profileData.joinDate}</p>
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
                    <span className="info-value">{profileData.joinDate}</span>
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