import React, { useState, useEffect } from 'react';

function PaymentModal({ curso, onClose, onSuccess }) {
  const [paymentMethod, setPaymentMethod] = useState('tarjeta');
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState(null);
  const [cardDetails, setCardDetails] = useState({
    number: '',
    expiry: '',
    cvv: ''
  });
  // Estado para guardar la URL de checkout externa (para Coinbase, PayPal, etc.)
  const [checkoutUrl, setCheckoutUrl] = useState(null);
  
  // Determinar si estamos en la etapa de checkout externo
  const [externalCheckout, setExternalCheckout] = useState(false);
  
  // Usando una constante en lugar de un estado para el modo de desarrollo
  const devMode = process.env.NODE_ENV === 'development' && false; // Cambiar a true para forzar modo desarrollo

  // Opción de pre-cargar scripts de pasarelas de pago
  const [loadedScripts, setLoadedScripts] = useState({
    stripe: false,
    paypal: false,
    mercadopago: false
  });

  // Claves de API públicas para pasarelas (en producción usar variables de entorno)
  const apiKeys = {
    stripe: process.env.REACT_APP_STRIPE_PUBLIC_KEY || 'pk_test_yourkey',
    paypal: process.env.REACT_APP_PAYPAL_CLIENT_ID || 'sb-clientid',
    mercadopago: process.env.REACT_APP_MERCADOPAGO_PUBLIC_KEY || 'TEST-pubkey-123'
  };

  // Cargar scripts necesarios según método de pago
  useEffect(() => {
    const loadScript = (id, src) => {
      return new Promise((resolve, reject) => {
        if (document.getElementById(id)) {
          resolve();
          return;
        }
        
        const script = document.createElement('script');
        script.id = id;
        script.src = src;
        script.async = true;
        
        script.onload = () => resolve();
        script.onerror = () => reject(new Error(`Failed to load script: ${src}`));
        
        document.body.appendChild(script);
      });
    };
    
    const loadPaymentScripts = async () => {
      try {
        if (paymentMethod === 'stripe' && !loadedScripts.stripe) {
          await loadScript('stripe-js', 'https://js.stripe.com/v3/');
          setLoadedScripts(prev => ({ ...prev, stripe: true }));
        } 
        else if (paymentMethod === 'paypal' && !loadedScripts.paypal) {
          await loadScript('paypal-js', `https://www.paypal.com/sdk/js?client-id=${apiKeys.paypal}&currency=USD`);
          setLoadedScripts(prev => ({ ...prev, paypal: true }));
        }
        else if (paymentMethod === 'mercadopago' && !loadedScripts.mercadopago) {
          await loadScript('mercadopago-js', 'https://secure.mlstatic.com/sdk/javascript/v1/mercadopago.js');
          setLoadedScripts(prev => ({ ...prev, mercadopago: true }));
        }
      } catch (error) {
        console.error('Error loading payment script:', error);
      }
    };
    
    loadPaymentScripts();
  }, [paymentMethod, loadedScripts, apiKeys]);

  // Efecto para gestionar el checkout externo
  useEffect(() => {
    // Si tenemos una URL de checkout, abrir en una nueva ventana o redirigir
    if (checkoutUrl && !externalCheckout) {
      setExternalCheckout(true);
      
      const checkoutWindow = window.open(checkoutUrl, '_blank');
      
      // En caso de que el navegador bloquee el popup
      if (!checkoutWindow || checkoutWindow.closed || typeof checkoutWindow.closed === 'undefined') {
        // Mostrar mensaje alternativo con enlace
        console.log("El navegador bloqueó la ventana de pago. Usar enlace alternativo.");
      }
    }
  }, [checkoutUrl, externalCheckout]);

  // Manejar el regreso después de un pago externo
  useEffect(() => {
    if (externalCheckout) {
      let intervalId;
      
      const checkPaymentStatus = async () => {
        if (curso && curso.id) {
          try {
            const token = localStorage.getItem('token');
            if (!token) return;
            
            const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
            const response = await fetch(`${apiUrl}/api/pagos/${curso.id}`, {
              headers: {
                'Authorization': `Bearer ${token}`,
                'Accept': 'application/json'
              }
            });
            
            if (!response.ok) return;
            
            const data = await response.json();
            
            if (data.estado === 'aprobado') {
              clearInterval(intervalId);
              onSuccess();
            } else if (data.estado === 'rechazado') {
              clearInterval(intervalId);
              setError('El pago fue rechazado. Por favor intenta con otro método de pago.');
              setIsProcessing(false);
              setExternalCheckout(false);
            }
          } catch (error) {
            console.error('Error checking payment status:', error);
          }
        }
      };
      
      // Verificar cada 3 segundos
      intervalId = setInterval(checkPaymentStatus, 3000);
      
      return () => clearInterval(intervalId);
    }
  }, [externalCheckout, curso, onSuccess]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsProcessing(true);
    setError(null);

    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No has iniciado sesión. Por favor inicia sesión para continuar.');
      }

      // Construir datos básicos de pago
      const paymentData = {
        curso_id: curso.id,
        monto: curso.precio || 29.99,
        metodo: devMode ? 'dev' : paymentMethod,
        redirect_url: window.location.origin + '/payment-callback'
      };

      // Agregar datos específicos según el método de pago
      if (paymentMethod === 'tarjeta' && !devMode) {
        paymentData.detalles_tarjeta = {
          numero: cardDetails.number.replace(/\s/g, ''),
          expiracion: cardDetails.expiry,
          cvv: cardDetails.cvv
        };
      }
      
      // Para métodos que requieren token (en producción se obtendría del SDK)
      if (paymentMethod === 'stripe' && loadedScripts.stripe) {
        // En producción, integrar con Stripe Elements o Checkout
        const stripe = window.Stripe(apiKeys.stripe);
        // Ejemplo: crear token con Stripe (en producción usar elementos seguros)
        // const { token } = await stripe.createToken(...);
        // paymentData.payment_token = token.id;
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

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || `Error ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();
      console.log('Respuesta de pago:', data);
      
      // Guardar ID de pago para referencia
      if (data.pago_id) {
        localStorage.setItem('current_payment_id', data.pago_id);
      }
      
      // Manejar respuesta según el método de pago
      if (data.checkout_url) {
        // Para pagos externos (PayPal, Coinbase, etc.)
        setCheckoutUrl(data.checkout_url);
        // No cerramos el modal, solo mostramos mensaje de redirección
      } else {
        // Para métodos de pago procesados internamente
        if (devMode) {
          setTimeout(() => {
            onSuccess();
          }, 1000);
        } else {
          // En modo producción, esperamos a la respuesta real
          setTimeout(() => {
            checkPaymentStatus(curso.id);
          }, 2000);
        }
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
  const checkPaymentStatus = async (cursoId) => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No has iniciado sesión');
      }

      const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
      console.log(`Verificando estado del pago. ID Curso: ${cursoId}`);
      
      const response = await fetch(`${apiUrl}/api/pagos/${cursoId}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Accept': 'application/json'
        }
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || `Error ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();
      console.log('Estado del pago:', data);

      // Acceder directamente al estado
      const estadoPago = data.estado || 'pendiente';

      if (estadoPago === 'aprobado') {
        onSuccess();
      } else if (estadoPago === 'rechazado') {
        setError('El pago fue rechazado. Por favor intenta con otro método de pago.');
        setIsProcessing(false);
      } else if (estadoPago === 'pendiente') {
        // Si sigue pendiente, volvemos a verificar después de un tiempo
        setTimeout(() => checkPaymentStatus(cursoId), 2000);
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
      // Formatear número de tarjeta en grupos de 4 dígitos
      const cleaned = value.replace(/\D/g, '').substring(0, 16);
      processedValue = cleaned.replace(/(\d{4})(?=\d)/g, '$1 ');
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

  // Determinar si mostrar mensajes de redirección
  const showRedirectMessage = externalCheckout && checkoutUrl;
  
  // Renderizar el componente del método de pago seleccionado
  const renderPaymentMethodComponent = () => {
    if (devMode) return null;
    
    switch (paymentMethod) {
      case 'tarjeta':
        return (
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
            <div className="card-brands">
              <img src="/images/visa.svg" alt="Visa" />
              <img src="/images/mastercard.svg" alt="Mastercard" />
              <img src="/images/amex.svg" alt="American Express" />
            </div>
          </div>
        );
        
      case 'paypal':
        return (
          <div className="payment-info">
            <div className="payment-brand">
              <img src="/images/paypal.svg" alt="PayPal" className="payment-logo" />
            </div>
            <p>Serás redirigido a PayPal para completar el pago de forma segura.</p>
          </div>
        );
        
      case 'stripe':
        return (
          <div className="payment-info">
            <div className="payment-brand">
              <img src="/images/stripe.svg" alt="Stripe" className="payment-logo" />
            </div>
            <p>Serás redirigido a una página segura de Stripe para completar el pago.</p>
          </div>
        );
        
      case 'coinbase':
        return (
          <div className="payment-info">
            <div className="payment-brand">
              <img src="/images/coinbase.svg" alt="Coinbase" className="payment-logo" />
            </div>
            <p>Serás redirigido a Coinbase Commerce para pagar con criptomonedas.</p>
            <p className="crypto-info">Aceptamos Bitcoin, Ethereum, Litecoin y más.</p>
          </div>
        );
        
      case 'mercadopago':
        return (
          <div className="payment-info">
            <div className="payment-brand">
              <img src="/images/mercadopago.svg" alt="Mercado Pago" className="payment-logo" />
            </div>
            <p>Serás redirigido a Mercado Pago para completar la transacción.</p>
          </div>
        );
        
      case 'transferencia':
        return (
          <div className="payment-info">
            <div className="bank-info">
              <h4>Datos bancarios:</h4>
              <p><strong>Banco:</strong> Nombre del Banco</p>
              <p><strong>Titular:</strong> Nombre de la Empresa</p>
              <p><strong>Cuenta:</strong> XXXX-XXXX-XXXX-XXXX</p>
              <p><strong>CBU/CLABE:</strong> XXXXXXXXXXXXXXXXX</p>
              <p className="transfer-note">Importante: Incluye tu número de pedido en la referencia de la transferencia.</p>
            </div>
          </div>
        );
        
      default:
        return null;
    }
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
        
        {/* Mostrar mensaje cuando se ha redirigido a un checkout externo */}
        {showRedirectMessage && (
          <div className="redirect-message">
            <div className="redirect-icon">
              <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                <polyline points="15 3 21 3 21 9"></polyline>
                <line x1="10" y1="14" x2="21" y2="3"></line>
              </svg>
            </div>
            <h4>Procesando tu pago</h4>
            <p>Te hemos redirigido a la página de pago. Por favor completa el proceso en la ventana abierta.</p>
            <p>Esta ventana verificará automáticamente cuando el pago esté completado.</p>
            <button 
              onClick={() => window.open(checkoutUrl, '_blank')}
              className="checkout-button"
            >
              Abrir página de pago nuevamente
            </button>
          </div>
        )}

        {!showRedirectMessage && (
          <div className="modal-content">
            <div className="payment-details">
              <p className="course-price">Precio: ${curso.precio?.toFixed(2) || '29.99'}</p>
            </div>
            
            <form onSubmit={handleSubmit}>
              <div className="form-group payment-method-selector">
                <label>Método de pago:</label>
                <select 
                  value={paymentMethod}
                  onChange={(e) => setPaymentMethod(e.target.value)}
                  disabled={isProcessing}
                  required
                >
                  <option value="tarjeta">Tarjeta de crédito/débito</option>
                  <option value="paypal">PayPal</option>
                  <option value="stripe">Stripe</option>
                  <option value="coinbase">Coinbase (Crypto)</option>
                  <option value="mercadopago">Mercado Pago</option>
                  <option value="transferencia">Transferencia bancaria</option>
                </select>
              </div>

              {renderPaymentMethodComponent()}

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
            
            <div className="secure-payment-info">
              <div className="secure-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                </svg>
              </div>
              <p>Pago 100% seguro. Tus datos están protegidos.</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default PaymentModal;