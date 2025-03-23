import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './LoginPage.css'; // Reutilizamos los estilos

function ForgotPasswordPage() {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setMessage('');
    
    try {
      const response = await fetch('http://localhost:5000/api/auth/forgot-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email }),
      });
      
      // Verificar si la respuesta es JSON antes de intentar parsearla
      const contentType = response.headers.get('content-type');
      let data;
      
      if (contentType && contentType.includes('application/json')) {
        data = await response.json();
      } else {
        // Si no es JSON, obtener el texto de la respuesta
        const text = await response.text();
        throw new Error(`Error en la respuesta del servidor: ${text}`);
      }
      
      // Verificar si hay error en la respuesta
      if (!response.ok) {
        throw new Error(data.error || data.message || 'Error al procesar la solicitud');
      }
      
      setMessage('Se ha enviado un enlace de recuperación a tu correo electrónico.');
      setSubmitted(true);
      
      // En un entorno de desarrollo, mostramos el token para facilitar las pruebas
      if (data.resetToken) {
        setMessage(prev => `${prev}\n\nToken (solo para desarrollo): ${data.resetToken}`);
      }
      
      if (data.resetLink) {
        setMessage(prev => `${prev}\n\nEnlace: ${data.resetLink}`);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page animate__animated animate__fadeIn">
      <h2>Recuperar Contraseña</h2>
      
      <form className="login-form" onSubmit={handleSubmit}>
        {error && <div className="error-message">{error}</div>}
        {message && <div className="success-message">{message}</div>}
        
        {!submitted ? (
          <>
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
                placeholder=" "
              />
            </div>
            
            <button type="submit" className="submit-btn" disabled={loading}>
              {loading ? 'Procesando...' : 'Enviar Enlace de Recuperación'}
            </button>
          </>
        ) : (
          <button 
            type="button" 
            className="submit-btn" 
            onClick={() => navigate('/login')}
          >
            Volver a Iniciar Sesión
          </button>
        )}
      </form>
      
      {!submitted && (
        <div className="register-link">
          <a href="/login">Volver a Iniciar Sesión</a>
        </div>
      )}
    </div>
  );
}

export default ForgotPasswordPage;