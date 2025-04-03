import React, { useState } from 'react';

function PaymentModal({ curso, onClose, onSuccess }) {
  const [paymentMethod, setPaymentMethod] = useState('tarjeta');
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState(null);
  const [cardDetails, setCardDetails] = useState({
    number: '',
    expiry: '',
    cvv: ''
  });
  // Usando una constante en lugar de un estado para el modo de desarrollo
  const devMode = false;

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsProcessing(true);
    setError(null);

    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No has iniciado sesión. Por favor inicia sesión para continuar.');
      }

      const paymentData = {
        curso_id: curso.id,
        monto: curso.precio || 29.99,
        metodo: devMode ? 'dev' : paymentMethod
      };

      if (paymentMethod === 'tarjeta' && !devMode) {
        paymentData.detalles_tarjeta = {
          numero: cardDetails.number,
          expiracion: cardDetails.expiry,
          cvv: cardDetails.cvv
        };
      }

      // Usamos la URL absoluta del backend para evitar problemas de CORS
      const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
      
      console.log("Enviando solicitud de pago a:", `${apiUrl}/api/pagos`);
      console.log("Datos del pago:", paymentData);
      
      const response = await fetch(`${apiUrl}/api/pagos`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(paymentData)
      });

      console.log(`Respuesta recibida con estado: ${response.status} ${response.statusText}`);
      console.log(`Content-Type: ${response.headers.get('content-type')}`);

      // Leer el cuerpo de la respuesta como texto primero
      const responseText = await response.text();
      console.log('Respuesta como texto:', responseText);
      
      // Verificar si podemos parsearlo como JSON
      let data;
      try {
        data = JSON.parse(responseText);
        console.log('Datos parseados:', data);
      } catch (jsonError) {
        // Si no es JSON válido, usamos el texto como mensaje de error
        console.error('Respuesta no JSON:', responseText);
        
        if (responseText.includes('<!DOCTYPE html>') || responseText.includes('<html>')) {
          throw new Error('El servidor devolvió una página HTML en lugar de JSON. Esto puede indicar que el backend no está procesando correctamente la solicitud.');
        } else {
          throw new Error(responseText || `Error ${response.status}: ${response.statusText}`);
        }
      }
      
      // Verificar el estado HTTP y la respuesta
      if (!response.ok) {
        throw new Error(data.error || `Error ${response.status}: ${response.statusText}`);
      }
      
      if (data.success === false) {
        throw new Error(data.error || 'Error al procesar el pago');
      }

      console.log('Respuesta de pago:', data);
      
      // Si estamos en modo desarrollo, simulamos una espera más corta
      if (devMode) {
        setTimeout(() => {
          onSuccess();
        }, 1000);
      } else {
        // En modo producción, esperamos a la respuesta real (simulada en el backend)
        setTimeout(() => {
          checkPaymentStatus(data.pago_id);
        }, 2000);
      }

    } catch (err) {
      console.error('Error completo:', err);
      
      let errorMessage = err.message || 'Error desconocido al procesar el pago';
      
      // Mensajes de error más específicos basados en el tipo de error
      if (err.message.includes('Failed to fetch')) {
        errorMessage = 'No se pudo conectar con el servidor. Verifica que el backend esté en ejecución.';
      } else if (err.message.includes('401')) {
        errorMessage = 'Tu sesión ha expirado. Por favor inicia sesión nuevamente.';
      } else if (err.message.includes('404')) {
        errorMessage = 'El servicio de pagos no está disponible en este momento. Inténtalo más tarde.';
      }
      
      setError(errorMessage);
      setIsProcessing(false);
    }
  };

  // Función para verificar el estado del pago
  const checkPaymentStatus = async (pagoId) => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No has iniciado sesión');
      }

      const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
      console.log(`Verificando estado del pago. ID Curso: ${curso.id}, ID Pago: ${pagoId}`);
      
      const response = await fetch(`${apiUrl}/api/pagos/${curso.id}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Accept': 'application/json'
        }
      });

      console.log(`Respuesta recibida con estado: ${response.status} ${response.statusText}`);
      console.log(`Content-Type: ${response.headers.get('content-type')}`);

      const responseText = await response.text();
      console.log('Respuesta como texto:', responseText);
      
      let data;
      try {
        data = JSON.parse(responseText);
        console.log('Datos parseados:', data);
      } catch (error) {
        console.error('Error al parsear respuesta:', responseText);
        if (responseText.includes('<!DOCTYPE html>') || responseText.includes('<html>')) {
          throw new Error('El servidor devolvió una página HTML en lugar de JSON.');
        } else {
          throw new Error('Respuesta no válida del servidor');
        }
      }

      // Acceder directamente al estado (sin anidamiento)
      const estadoPago = data.estado || 'pendiente';

      if (estadoPago === 'aprobado') {
        onSuccess();
      } else if (estadoPago === 'rechazado') {
        setError('El pago fue rechazado. Por favor intenta con otro método de pago.');
        setIsProcessing(false);
      } else if (estadoPago === 'pendiente') {
        // Si sigue pendiente, volvemos a verificar después de un tiempo
        setTimeout(() => checkPaymentStatus(pagoId), 2000);
      } else {
        console.error(`Estado de pago no reconocido: "${estadoPago}"`, data);
        setError(`Estado de pago desconocido: ${estadoPago}. Por favor contacta a soporte.`);
        setIsProcessing(false);
      }
    } catch (err) {
      console.error('Error al verificar estado:', err);
      
      let errorMessage = 'Error al verificar el estado del pago';
      if (err.message.includes('Failed to fetch')) {
        errorMessage = 'No se pudo conectar con el servidor. Verifica tu conexión a internet.';
      } else if (err.message.includes('HTML')) {
        errorMessage = 'El servidor no está respondiendo correctamente. Contacta al administrador.';
      }
      
      setError(errorMessage);
      setIsProcessing(false);
    }
  };

  const handleCardChange = (e) => {
    const { name, value } = e.target;
    
    let processedValue = value;
    if (name === 'number') {
      processedValue = value.replace(/\D/g, '').substring(0, 16);
    } else if (name === 'expiry') {
      const cleaned = value.replace(/\D/g, '').substring(0, 4);
      if (cleaned.length > 2) {
        processedValue = cleaned.substring(0, 2) + '/' + cleaned.substring(2);
      } else {
        processedValue = cleaned;
      }
    } else if (name === 'cvv') {
      processedValue = value.replace(/\D/g, '').substring(0, 4);
    }
    
    setCardDetails(prev => ({
      ...prev,
      [name]: processedValue
    }));
  };

  return (
    <div className="payment-modal-overlay">
      <div className="payment-modal">
        <div className="modal-header">
          <h3>Pagar curso: {curso.titulo}</h3>
          <button 
            onClick={onClose}
            className="close-button"
            disabled={isProcessing}
          >
            &times;
          </button>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        <div className="modal-content">
          <div className="payment-details">
            <p>Precio: ${curso.precio?.toFixed(2) || '29.99'}</p>
          </div>
          
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label>Método de pago:</label>
              <select 
                value={paymentMethod}
                onChange={(e) => setPaymentMethod(e.target.value)}
                disabled={isProcessing}
                required
              >
                <option value="tarjeta">Tarjeta de crédito/débito</option>
                <option value="paypal">PayPal</option>
                <option value="transferencia">Transferencia bancaria</option>
              </select>
            </div>

            {!devMode && paymentMethod === 'tarjeta' && (
              <div className="card-details">
                <div className="form-group">
                  <label>Número de tarjeta</label>
                  <input 
                    type="text" 
                    name="number"
                    placeholder="1234 5678 9012 3456" 
                    value={cardDetails.number}
                    onChange={handleCardChange}
                    disabled={isProcessing}
                    required
                  />
                </div>
                <div className="form-row">
                  <div className="form-group">
                    <label>Fecha expiración (MM/AA)</label>
                    <input 
                      type="text" 
                      name="expiry"
                      placeholder="MM/AA" 
                      value={cardDetails.expiry}
                      onChange={handleCardChange}
                      disabled={isProcessing}
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>CVV</label>
                    <input 
                      type="text" 
                      name="cvv"
                      placeholder="123" 
                      value={cardDetails.cvv}
                      onChange={handleCardChange}
                      disabled={isProcessing}
                      required
                    />
                  </div>
                </div>
              </div>
            )}

            <button 
              type="submit" 
              className="submit-button"
              disabled={isProcessing}
            >
              {isProcessing ? (
                <>
                  <span className="spinner"></span>
                  Procesando pago...
                </>
              ) : (
                'Pagar ahora'
              )}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}

export default PaymentModal;