import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './ResetPasswordPage.css';

function ResetPasswordPage() {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);
  const [validatingToken, setValidatingToken] = useState(true);
  const [tokenValid, setTokenValid] = useState(false);
  const [formSubmitted, setFormSubmitted] = useState(false);
  const { token } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    // Validar el token al cargar la página
    const validateToken = async () => {
      try {
        const response = await fetch(`http://localhost:5000/api/auth/reset-password/${token}/validate`);
        
        // Verificar si la respuesta es JSON
        const contentType = response.headers.get('content-type');
        let data;
        
        if (contentType && contentType.includes('application/json')) {
          data = await response.json();
        } else {
          // Si no es JSON, obtener el texto
          const text = await response.text();
          throw new Error(`Error en la respuesta del servidor: ${text}`);
        }
        
        if (response.ok && data.valid) {
          setTokenValid(true);
        } else {
          setError(data.error || 'El enlace de restablecimiento no es válido o ha expirado.');
        }
      } catch (err) {
        setError(err.message || 'Error al validar el enlace de restablecimiento.');
      } finally {
        setValidatingToken(false);
      }
    };

    validateToken();
  }, [token]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validar que las contraseñas coincidan
    if (password !== confirmPassword) {
      setError('Las contraseñas no coinciden.');
      return;
    }
    
    setLoading(true);
    setError('');
    setSuccess('');
    
    try {
      const response = await fetch('http://localhost:5000/api/auth/reset-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
          token,
          password 
        }),
      });
      
      // Verificar si la respuesta es JSON
      const contentType = response.headers.get('content-type');
      let data;
      
      if (contentType && contentType.includes('application/json')) {
        data = await response.json();
      } else {
        // Si no es JSON, obtener el texto
        const text = await response.text();
        throw new Error(`Error en la respuesta del servidor: ${text}`);
      }
      
      // Verificar el estado de la respuesta
      if (!response.ok) {
        throw new Error(data.error || 'Error al restablecer la contraseña');
      }
      
      setSuccess('¡Contraseña actualizada!');
      setFormSubmitted(true);
      
      // Ya no necesitamos el timeout para redirección, ahora mostramos un botón como en ForgotPasswordPage
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (validatingToken) {
    return (
      <div className="reset-password-page">
        <h2>Validando Enlace</h2>
        <p className="form-description">Por favor espera mientras validamos tu enlace de restablecimiento...</p>
      </div>
    );
  }

  if (!tokenValid) {
    return (
      <div className="reset-password-page">
        <h2>Enlace No Válido</h2>
        <div className="reset-password-error-message">{error}</div>
        <div className="login-link">
          <a href="/recuperar-contrasena">Solicitar un nuevo enlace</a>
        </div>
      </div>
    );
  }

  return (
    <div className="reset-password-page">
      <h2>Restablecer Contraseña</h2>
      
      {!formSubmitted ? (
        <form className="reset-password-form" onSubmit={handleSubmit}>
          {/* Mostrar el contenedor de mensajes sólo cuando hay algún mensaje */}
          {(error || success) && (
            <div className="reset-password-messages">
              {error && <div className="reset-password-error-message">{error}</div>}
              {success && <div className="reset-password-success-message">{success}</div>}
            </div>
          )}
          
          <p className="form-description">
            Introduce tu nueva contraseña para completar el proceso de restablecimiento.
          </p>
          
          <div className="form-group">
            <label htmlFor="password">Nueva Contraseña</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength="6"
              className="form-input"
              placeholder="Ingresa tu nueva contraseña"
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="confirmPassword">Confirmar Contraseña</label>
            <input
              type="password"
              id="confirmPassword"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              className="form-input"
              placeholder="Confirma tu nueva contraseña"
            />
          </div>
          
          <button type="submit" className="submit-btn" disabled={loading}>
            {loading ? 'Procesando...' : 'Restablecer Contraseña'}
          </button>
          
          <div className="login-link">
            <a href="/login">Volver a Iniciar Sesión</a>
          </div>
        </form>
      ) : (
        <>
          <div className="reset-password-success-container">
            <div className="reset-password-success-message">
              <div className="success-title">¡Contraseña actualizada!</div>
              <p>Tu contraseña ha sido restablecida correctamente.</p>
              <p>Ya puedes iniciar sesión con tu nueva contraseña.</p>
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

export default ResetPasswordPage;