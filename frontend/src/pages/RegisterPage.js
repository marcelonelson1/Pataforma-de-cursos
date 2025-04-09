import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './RegisterPage.css';

function RegisterPage() {
  const [formData, setFormData] = useState({
    nombre: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    let redirectTimer;
    if (success) {
      redirectTimer = setTimeout(() => {
        navigate('/login');
      }, 3000);
    }
    return () => clearTimeout(redirectTimer);
  }, [success, navigate]);

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');
    
    if (formData.password !== formData.confirmPassword) {
      setError('Las contraseñas no coinciden');
      setLoading(false);
      return;
    }
    
    try {
      const response = await fetch('http://localhost:5000/api/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          nombre: formData.nombre,
          email: formData.email,
          password: formData.password
        })
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.error || 'Error al registrarse');
      }
      
      // Establecer el éxito pero sin redirigir automáticamente
      setSuccess(`¡Registro exitoso, ${formData.nombre}! Redirigiendo...`);
      
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="register-container">
      <div className="register-page">
        <h2>Crear Cuenta</h2>
        
        <form className="register-form" onSubmit={handleSubmit}>
          {/* Mostrar el contenedor de mensajes sólo cuando hay algún mensaje */}
          {(error || success) && (
            <div className="register-messages">
              {error && <div className="register-error-message">{error}</div>}
              {success && <div className="register-success-message">{success}</div>}
            </div>
          )}
          
          <div className="form-group">
            <label htmlFor="nombre">Nombre</label>
            <input
              type="text"
              id="nombre"
              name="nombre"
              value={formData.nombre}
              onChange={handleChange}
              required
              className="form-input"
              placeholder="Tu nombre"
            />
          </div>
          
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
              placeholder="ejemplo@correo.com"
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="password">Contraseña</label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              className="form-input"
              minLength="6"
              placeholder="Al menos 6 caracteres"
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="confirmPassword">Confirmar Contraseña</label>
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={handleChange}
              required
              className="form-input"
              minLength="6"
              placeholder="Confirma tu contraseña"
            />
          </div>
          
          <button type="submit" className="submit-btn" disabled={loading}>
            {loading ? 'Procesando...' : 'Registrarse'}
          </button>
        </form>
        
        <div className="login-link">
          ¿Ya tienes una cuenta? <Link to="/login">Inicia sesión aquí</Link>
        </div>
      </div>
    </div>
  );
}

export default RegisterPage;