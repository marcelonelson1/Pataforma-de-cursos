// ForgotPasswordPage.js
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './ForgotPasswordPage.css';

function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);
  const [formSubmitted, setFormSubmitted] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');
    
    try {
      const response = await fetch('http://localhost:5000/api/auth/forgot-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email }),
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
        throw new Error(data.error || data.message || 'Error al procesar la solicitud');
      }
      
      setSuccess(`¡Correo enviado!`);
      setFormSubmitted(true);
      
    } catch (err) {
      console.error("Error en la solicitud:", err);
      setError(err.message || 'Ocurrió un error al procesar la solicitud');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="forgot-password-page">
      <h2>Recuperar Contraseña</h2>
      
      {!formSubmitted ? (
        <form className="forgot-password-form" onSubmit={handleSubmit}>
          {/* Mostrar el contenedor de mensajes sólo cuando hay algún mensaje */}
          {(error || success) && (
            <div className="forgot-password-messages">
              {error && <div className="forgot-password-error-message">{error}</div>}
              {success && <div className="forgot-password-success-message">{success}</div>}
            </div>
          )}
          
          <p className="form-description">
            Introduce tu dirección de correo electrónico y te enviaremos un enlace para restablecer tu contraseña.
          </p>
          
          <div className="form-group">
            <label htmlFor="email">Email</label>
            <input
              type="email"
              id="email"
              name="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="form-input"
              placeholder="correo@ejemplo.com"
            />
          </div>
          
          <button type="submit" className="submit-btn" disabled={loading}>
            {loading ? 'Procesando...' : 'Enviar Enlace de Recuperación'}
          </button>
          
          <div className="login-link">
            <a href="/login">Volver a Iniciar Sesión</a>
          </div>
        </form>
      ) : (
        <>
          <div className="forgot-password-success-container">
            <div className="forgot-password-success-message">
              <div className="success-title">¡Correo enviado!</div>
              <p>Se ha enviado un enlace de recuperación a:</p>
              <p className="success-email"><strong>{email}</strong></p>
              <p>Revisa tu bandeja de entrada y sigue las instrucciones.</p>
            </div>
          </div>
          
          <button 
            type="button" 
            className="submit-btn login-button" 
            onClick={() => navigate('/login')}
          >
            Volver a Iniciar Sesión
          </button>
        </>
      )}
    </div>
  );
}

export default ForgotPasswordPage;