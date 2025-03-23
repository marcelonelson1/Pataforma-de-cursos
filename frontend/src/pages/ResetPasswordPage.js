import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './LoginPage.css'; // Reutilizamos los estilos

function ResetPasswordPage() {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [validatingToken, setValidatingToken] = useState(true);
  const [tokenValid, setTokenValid] = useState(false);
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
    setMessage('');
    
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
      
      setMessage('Tu contraseña ha sido actualizada correctamente.');
      
      // Redirigir al login después de 3 segundos
      setTimeout(() => {
        navigate('/login');
      }, 3000);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (validatingToken) {
    return (
      <div className="login-page animate__animated animate__fadeIn">
        <h2>Validando Enlace</h2>
        <p>Por favor espera mientras validamos tu enlace de restablecimiento...</p>
      </div>
    );
  }

  if (!tokenValid) {
    return (
      <div className="login-page animate__animated animate__fadeIn">
        <h2>Enlace No Válido</h2>
        <div className="error-message">{error}</div>
        <div className="register-link">
          <a href="/recuperar-contrasena">Solicitar un nuevo enlace</a>
        </div>
      </div>
    );
  }

  return (
    <div className="login-page animate__animated animate__fadeIn">
      <h2>Restablecer Contraseña</h2>
      
      <form className="login-form" onSubmit={handleSubmit}>
        {error && <div className="error-message">{error}</div>}
        {message && <div className="success-message">{message}</div>}
        
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
            placeholder=" "
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
            placeholder=" "
          />
        </div>
        
        <button type="submit" className="submit-btn" disabled={loading}>
          {loading ? 'Procesando...' : 'Restablecer Contraseña'}
        </button>
      </form>
    </div>
  );
}

export default ResetPasswordPage;