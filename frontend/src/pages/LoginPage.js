import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './LoginPage.css';

function LoginPage() {
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuth();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    
    try {
      const response = await fetch('http://localhost:5000/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData),
        credentials: 'include'
      });
      
      const contentType = response.headers.get('content-type');
      let data;
      
      if (contentType && contentType.includes('application/json')) {
        data = await response.json();
      } else {
        const text = await response.text();
        throw new Error(`Error en la respuesta del servidor: ${text}`);
      }
      
      if (!response.ok) {
        throw new Error(data.error || data.message || 'Error al iniciar sesión');
      }
      
      login(data.token);
      navigate('/', { replace: true });
    } catch (err) {
      setError(err.message || 'Error de red al intentar conectar con el servidor');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page animate__animated animate__fadeIn">
      <h2>Iniciar Sesión</h2>
      
      <form className="login-form" onSubmit={handleSubmit}>
        {error && <div className="error-message">{error}</div>}
        
        

        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            required
            className="form-input"
            placeholder=" "
          />
        </div>
        
        <div className="form-group">
          <label htmlFor="password">Contraseña</label>
          <div className="password-input-container">
            <input
              type={showPassword ? 'text' : 'password'}
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              className="form-input"
              placeholder=" "
            />
            <button
              type="button"
              className="password-toggle-btn"
              onClick={togglePasswordVisibility}
              aria-label={showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
            >
              {showPassword ? (
                <svg viewBox="0 0 24 24" width="20" height="20">
                  <path fill="currentColor" d="M12 9a3 3 0 0 1 3 3a3 3 0 0 1-3 3a3 3 0 0 1-3-3a3 3 0 0 1 3-3m0-4.5c5 0 9.27 3.11 11 7.5c-1.73 4.39-6 7.5-11 7.5S2.73 16.39 1 12c1.73-4.39 6-7.5 11-7.5M3.18 12a9.821 9.821 0 0 0 17.64 0a9.821 9.821 0 0 0-17.64 0Z"/>
                </svg>
              ) : (
                <svg viewBox="0 0 24 24" width="20" height="20">
                  <path fill="currentColor" d="M2 5.27L3.28 4L20 20.72L18.73 22l-3.08-3.08c-1.15.38-2.37.58-3.65.58c-5 0-9.27-3.11-11-7.5c.69-1.76 1.79-3.31 3.19-4.54L2 5.27M12 9a3 3 0 0 1 3 3a3 3 0 0 1-.17 1L11 9.17A3 3 0 0 1 12 9m0-4.5c5 0 9.27 3.11 11 7.5a11.79 11.79 0 0 1-4 5.19l-1.42-1.43A9.862 9.862 0 0 0 20.82 12A9.821 9.821 0 0 0 12 6.5c-1.09 0-2.16.18-3.16.5L7.3 5.47c1.44-.62 3.03-.97 4.7-.97M3.18 12A9.821 9.821 0 0 0 12 17.5c.69 0 1.37-.07 2-.2L11.72 15A3.064 3.064 0 0 1 9 12.28L5.6 8.87c-.99.85-1.82 1.91-2.42 3.13Z"/>
                </svg>
              )}
            </button>
          </div>
        </div>
        
        <div className="forgot-password">
          <a href="/recuperar-contrasena">¿Olvidaste tu contraseña?</a>
        </div>
        
        <button type="submit" className="submit-btn" disabled={loading}>
          {loading ? 'Procesando...' : 'Iniciar Sesión'}
        </button>
      </form>
      
      <div className="register-link">
        ¿No tienes cuenta? <a href="/registro">Regístrate aquí</a>
      </div>
    </div>
  );
}

export default LoginPage;