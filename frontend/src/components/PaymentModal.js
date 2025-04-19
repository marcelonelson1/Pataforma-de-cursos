import React, { useState, useEffect, useCallback } from 'react';
import './PaymentModal.css';

function PaymentModal({ curso, onClose, onSuccess }) {
  // Estados principales
  const [paymentMethod, setPaymentMethod] = useState('tarjeta');
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState(null);
  const [paymentId, setPaymentId] = useState(null);
  const [cardDetails, setCardDetails] = useState({
    number: '',
    expiry: '',
    cvv: ''
  });
  
  // Estado para scripts de pasarelas de pago
  const [loadedScripts, setLoadedScripts] = useState({
    stripe: false,
    paypal: false,
    mercadopago: false,
    coinbase: false
  });

  // Controlando el modo de desarrollo
  const devMode = process.env.NODE_ENV === 'development' && false;

  // Claves de API públicas para pasarelas (usando variables de entorno)
  const apiKeys = {
    stripe: process.env.REACT_APP_STRIPE_PUBLIC_KEY || 'pk_test_actualTestKey',
    paypal: process.env.REACT_APP_PAYPAL_CLIENT_ID || 'ASYN839bjb4gjMr6nRCc-7YYR8HutdM48kFMWhq-Sxp-PgB5c5R38yGiLBEPwDBIptFj8IJ71OPVXVUt',
    mercadopago: process.env.REACT_APP_MERCADOPAGO_PUBLIC_KEY || 'TEST-pubkey-actual',
    coinbase: process.env.REACT_APP_COINBASE_API_KEY || 'test-coinbase-key'
  };

  // Función para obtener el token de autenticación
  const getAuthToken = useCallback(() => {
    // Verificamos el token en diferentes ubicaciones
    const possibleKeys = [
      'token', 
      'auth_token', 
      'authToken', 
      'jwt'
    ];
    
    // Buscar en localStorage y sessionStorage
    for (const key of possibleKeys) {
      const token = localStorage.getItem(key) || sessionStorage.getItem(key);
      if (token) {
        console.log("Token de autenticación encontrado");
        return token.startsWith('Bearer ') ? token : `Bearer ${token}`;
      }
    }
    
    console.warn("No se encontró token de autenticación");
    return null;
  }, []);

  // Verificar si estamos regresando de un proceso de pago y manejar parámetros de URL
  useEffect(() => {
    const checkReturnFromPayment = async () => {
      const returnFromPayment = localStorage.getItem('returning_from_payment');
      const savedCursoId = localStorage.getItem('payment_curso_id');
      const storedPaymentMethod = localStorage.getItem('payment_method');
      
      // Obtener parámetros de URL para detectar cancelación o estado
      const urlParams = new URLSearchParams(window.location.search);
      const paymentStatus = urlParams.get('payment_status');
      const paymentCancel = urlParams.get('cancel') || urlParams.get('canceled');
      
      // Verificar si hay parámetros específicos de PayPal que indiquen cancelación
      const paypalCancel = urlParams.get('token') && !urlParams.get('PayerID');
      
      // Si tenemos señales de cancelación o fallo desde la URL, limpiar el estado
      if (paymentStatus === 'cancelled' || paymentStatus === 'failed' || paymentCancel || paypalCancel) {
        console.log("Detección de cancelación de pago o fallo mediante URL");
        localStorage.removeItem('returning_from_payment');
        localStorage.removeItem('payment_curso_id');
        localStorage.removeItem('payment_check_count');
        localStorage.removeItem('payment_start_time');
        localStorage.removeItem('payment_method');
        setIsProcessing(false);
        setError('El proceso de pago fue cancelado o no se completó.');
        return;
      }
      
      if (returnFromPayment === 'true' && savedCursoId && savedCursoId === String(curso.id)) {
        console.log("Detectado retorno de proceso de pago. Verificando estado...");
        setIsProcessing(true);
        
        // Si venimos de PayPal y no hay un PayerID, considerarlo cancelado
        if (storedPaymentMethod === 'paypal' && urlParams.get('token') && !urlParams.get('PayerID')) {
          console.log("Regresando de PayPal sin PayerID - Pago cancelado");
          clearPaymentData();
          setIsProcessing(false);
          setError('El pago con PayPal fue cancelado. Por favor, intenta nuevamente.');
          return;
        }
        
        // Verificar si ya se ha realizado una verificación después de regresar
        const paymentVerified = localStorage.getItem('payment_verified_after_return');
        
        if (paymentVerified === 'true') {
          // Ya se ha verificado una vez, asumir que fue cancelado
          console.log("Ya se verificó el pago después de regresar, sin cambios en estado.");
          clearPaymentData();
          setIsProcessing(false);
          setError('El pago no pudo ser procesado. Por favor, intenta con otro método de pago.');
          return;
        }
        
        // Marcar que estamos verificando después de regresar
        localStorage.setItem('payment_verified_after_return', 'true');
        
        // Verificar el estado del pago
        await checkPaymentStatus(curso.id);
      }
    };
    
    checkReturnFromPayment();
  }, [curso.id]);

  // Cargar scripts según el método de pago seleccionado
  useEffect(() => {
    const loadScript = (id, src) => {
      return new Promise((resolve) => {
        // Evitar cargar el mismo script múltiples veces
        if (document.getElementById(id)) {
          resolve();
          return;
        }
        
        const script = document.createElement('script');
        script.id = id;
        script.src = src;
        script.async = true;
        
        script.onload = () => {
          console.log(`Script cargado: ${id}`);
          resolve();
        };
        
        script.onerror = (error) => {
          console.error(`Error al cargar script ${id}:`, error);
          resolve(); // Resolvemos para no bloquear la UI
        };
        
        document.body.appendChild(script);
      });
    };
    
    const loadPaymentScripts = async () => {
      try {
        // Mapeo de métodos de pago con sus respectivos scripts
        const scriptConfigs = {
          stripe: {
            loaded: loadedScripts.stripe,
            src: 'https://js.stripe.com/v3/',
            id: 'stripe-js'
          },
          paypal: {
            loaded: loadedScripts.paypal,
            src: `https://www.paypal.com/sdk/js?client-id=${apiKeys.paypal}&currency=USD&intent=capture`,
            id: 'paypal-js'
          },
          mercadopago: {
            loaded: loadedScripts.mercadopago,
            src: 'https://secure.mlstatic.com/sdk/javascript/v1/mercadopago.js',
            id: 'mercadopago-js'
          },
          coinbase: {
            loaded: loadedScripts.coinbase,
            src: 'https://commerce.coinbase.com/v1/checkout.js',
            id: 'coinbase-js'
          }
        };
        
        // Cargar script según el método seleccionado
        const config = scriptConfigs[paymentMethod];
        if (config && !config.loaded) {
          console.log(`Cargando script para ${paymentMethod}...`);
          await loadScript(config.id, config.src);
          setLoadedScripts(prev => ({ ...prev, [paymentMethod]: true }));
        }
      } catch (error) {
        console.error('Error al cargar script de pago:', error);
        setError('Error al cargar métodos de pago. Por favor, intenta nuevamente.');
      }
    };
    
    if (paymentMethod && !devMode) {
      loadPaymentScripts();
    }
  }, [paymentMethod, loadedScripts, apiKeys, devMode]);

  // Manejador para el envío del formulario de pago
  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsProcessing(true);
    setError(null);
    setPaymentId(null);

    try {
      // Obtener token actualizado
      const token = getAuthToken();
      
      if (!token) {
        throw new Error('No has iniciado sesión. Por favor inicia sesión para continuar.');
      }

      // Método de pago efectivo (usar 'dev' en modo desarrollo para todos excepto PayPal)
      const effectivePaymentMethod = (devMode && paymentMethod !== 'paypal') ? 'dev' : paymentMethod;

      // Construir datos de pago
      const paymentData = {
        curso_id: curso.id,
        monto: curso.precio || 29.99,
        metodo: effectivePaymentMethod,
        moneda: 'USD',
      };

      // Agregar datos específicos según el método de pago
      if (paymentMethod === 'tarjeta' && !devMode) {
        paymentData.detalles_tarjeta = {
          numero: cardDetails.number.replace(/\s/g, ''),
          expiracion: cardDetails.expiry,
          cvv: cardDetails.cvv
        };
      }
      
      // Configurar URLs de retorno con parámetros mejorados
      const cursoUrl = `/cursos/${curso.id}`;
      
      // Configuración específica para PayPal
      if (paymentMethod === 'paypal') {
        const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
        const callbackUrl = `${apiUrl}/api/pagos/paypal/callback`;
        
        // Añadir parámetros explícitos para el callback
        paymentData.return_url = `${callbackUrl}?pago_id=${paymentId || 'pending'}&curso_id=${curso.id}`;
        paymentData.cancel_url = window.location.origin + cursoUrl + '?payment_status=cancelled';
        
        console.log("URLs de PayPal configuradas:", {
          return_url: paymentData.return_url,
          cancel_url: paymentData.cancel_url
        });
      } else {
        // URLs estándar para otros métodos
        paymentData.return_url = window.location.origin + cursoUrl + '?payment_status=success';
        paymentData.cancel_url = window.location.origin + cursoUrl + '?payment_status=cancelled';
      }
      
      // Inicializar SDK para métodos que lo requieren
      initializePaymentSDK(paymentMethod);

      // Usar la URL absoluta del backend
      const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
      
      console.log("Enviando solicitud a:", `${apiUrl}/api/pagos`);
      console.log("Datos:", paymentData);
      
      // Preparar headers para la solicitud
      const headers = {
        'Content-Type': 'application/json',
        'Authorization': token
      };
      
      // Enviar solicitud al backend
      const response = await fetch(`${apiUrl}/api/pagos`, {
        method: 'POST',
        headers: headers,
        body: JSON.stringify(paymentData),
        credentials: 'include' // Incluir cookies
      });

      // Obtener y parsear respuesta
      const responseText = await response.text();
      console.log("Respuesta completa:", responseText);
      
      let data;
      try {
        data = JSON.parse(responseText);
      } catch (e) {
        console.error("Error al parsear respuesta JSON:", e);
        throw new Error("La respuesta del servidor no es un JSON válido");
      }

      // Verificar errores
      if (!response.ok || (data && data.error)) {
        const errorMessage = data && (data.message || data.error) 
          ? data.message || data.error 
          : `Error ${response.status}: ${response.statusText}`;
        throw new Error(errorMessage);
      }
      
      // Guardar ID de pago
      const paymentIdFromResponse = 
        (data.data && data.data.pago_id) ? data.data.pago_id : 
        data.pago_id ? data.pago_id : null;
        
      if (paymentIdFromResponse) {
        localStorage.setItem('current_payment_id', paymentIdFromResponse);
        setPaymentId(paymentIdFromResponse);
      }
      
      // Guardar ID del curso
      localStorage.setItem('course_redirect_id', curso.id);
      
      // Buscar URL de checkout
      const checkoutUrl = extractCheckoutUrl(data);
      
      if (checkoutUrl) {
        console.log("✅ URL de checkout encontrada:", checkoutUrl);
        
        // Guardar info antes de redirigir y resetear el contador de verificaciones
        localStorage.setItem('payment_curso_id', curso.id);
        localStorage.setItem('returning_from_payment', 'true');
        localStorage.setItem('payment_check_count', '0');
        localStorage.setItem('payment_start_time', Date.now().toString());
        localStorage.setItem('payment_method', paymentMethod); // Guardar el método de pago
        
        // Uso de setTimeout para prevenir posibles bloqueos de redirección
        console.log("Redirigiendo a:", checkoutUrl);
        setTimeout(() => {
          // Intentar redirección de forma segura
          try {
            window.location.assign(checkoutUrl);
          } catch (e) {
            console.error("Error al redirigir usando assign:", e);
            window.location.href = checkoutUrl;
          }
        }, 100);
        return;
      } else {
        console.warn("⚠️ No se encontró URL de checkout:", data);
        
        // Si el pago está en proceso, verificar estado
        const isPending = 
          (data.data && data.data.estado === "pendiente") || 
          (data.estado && data.estado === "pendiente");
          
        if (isPending) {
          console.log("Pago pendiente. Iniciando verificación...");
          return;
        }
        
        // Para métodos procesados internamente
        const isInternalMethod = 
          devMode || 
          effectivePaymentMethod === 'dev' || 
          effectivePaymentMethod === 'tarjeta' || 
          effectivePaymentMethod === 'transferencia';
          
        if (isInternalMethod) {
          // Consultar estado después de un breve retardo
          setTimeout(() => checkPaymentStatus(curso.id), 3000);
        } else {
          // Error: no hay URL y no es método interno
          setError("No se pudo iniciar el proceso. Intenta con otro método de pago.");
          setIsProcessing(false);
        }
      }

    } catch (err) {
      console.error('Error completo:', err);
      
      // Mensajes de error específicos
      let errorMessage = err.message || 'Error desconocido al procesar el pago';
      
      if (err.message.includes('Failed to fetch') || err.message.includes('NetworkError')) {
        errorMessage = 'No se pudo conectar con el servidor. Verifica tu conexión a internet.';
      } else if (err.message.includes('401') || err.message.includes('token') || err.message.includes('autorización')) {
        errorMessage = 'Tu sesión ha expirado. Por favor inicia sesión nuevamente.';
      } else if (err.message.includes('404')) {
        errorMessage = 'El servicio de pagos no está disponible. Inténtalo más tarde.';
      }
      
      setError(errorMessage);
      setIsProcessing(false);
    }
  };

  // Inicializar SDK específico según método de pago
  const initializePaymentSDK = (method) => {
    if (devMode) return;
    
    if (method === 'stripe' && loadedScripts.stripe && window.Stripe) {
      try {
        window.Stripe(apiKeys.stripe);
        console.log("Stripe inicializado correctamente");
      } catch (err) {
        console.error("Error al inicializar Stripe:", err);
      }
    }
    
    if (method === 'coinbase' && loadedScripts.coinbase && window.CoinbaseCommerce) {
      try {
        console.log("Coinbase Commerce disponible");
      } catch (err) {
        console.error("Error con Coinbase Commerce:", err);
      }
    }
  };

  // Función mejorada para extraer URL de checkout
  const extractCheckoutUrl = (data) => {
    // Imprimir toda la estructura de datos para depuración
    console.log("Estructura completa de la respuesta:", JSON.stringify(data, null, 2));
    
    // Buscar específicamente para PayPal primero si ese es el método seleccionado
    if (paymentMethod === 'paypal') {
      // Buscar enlaces específicos de PayPal
      if (data.links && Array.isArray(data.links)) {
        const approvalLink = data.links.find(link => 
          link.rel === 'approval_url' || 
          link.rel === 'approve' || 
          (link.href && link.href.includes('paypal.com/checkoutnow'))
        );
        
        if (approvalLink && approvalLink.href) {
          console.log("Enlace de aprobación de PayPal encontrado:", approvalLink.href);
          return approvalLink.href;
        }
      }
    }
    
    // Buscar directamente la estructura conocida para cada pasarela
    if (paymentMethod === 'paypal' && data.data && data.data.checkout_url) {
      return data.data.checkout_url;
    }
    
    // Buscar en diferentes posibles ubicaciones
    const possiblePaths = [
      'data.checkout_url',
      'data.redirectUrl',
      'data.redirect_url',
      'data.url',
      'checkout_url',
      'redirectUrl',
      'redirect_url',
      'url',
      'hosted_url',
      'data.hosted_url',
      'approval_url',
      'data.approval_url'
    ];
    
    for (const path of possiblePaths) {
      const parts = path.split('.');
      let value = data;
      
      // Navegar por la estructura del objeto
      for (const part of parts) {
        if (value && typeof value === 'object' && part in value) {
          value = value[part];
        } else {
          value = undefined;
          break;
        }
      }
      
      if (value && typeof value === 'string') {
        console.log(`URL de checkout encontrada en '${path}':`, value);
        return value;
      }
    }
    
    // Última verificación: buscar cualquier propiedad que parezca una URL en el primer nivel
    if (data && typeof data === 'object') {
      for (const key in data) {
        if (typeof data[key] === 'string' && 
           (data[key].startsWith('http') && 
            (data[key].includes('checkout') || 
             data[key].includes('pay') || 
             data[key].includes('paypal.com')))) {
          console.log(`URL potencial encontrada en propiedad '${key}':`, data[key]);
          return data[key];
        }
      }
    }
    
    console.warn("No se pudo encontrar URL de checkout en la respuesta");
    return null;
  };

  // Función mejorada para verificar el estado del pago con límite de tiempo absoluto
  const checkPaymentStatus = async (cursoId) => {
    try {
      // Verificar si ha pasado demasiado tiempo desde el inicio del proceso (5 minutos)
      const startTime = parseInt(localStorage.getItem('payment_start_time') || '0', 10);
      const currentTime = Date.now();
      const MAX_PAYMENT_TIME = 5 * 60 * 1000; // 5 minutos en milisegundos
      
      // Si el tiempo límite se ha excedido, considerarlo como pago cancelado
      if (startTime > 0 && (currentTime - startTime) > MAX_PAYMENT_TIME) {
        console.log("Tiempo máximo de pago excedido (5 minutos)");
        setError('El tiempo para completar el pago ha expirado. Por favor, intenta nuevamente.');
        setIsProcessing(false);
        clearPaymentData();
        return;
      }
      
      // Verificar si el método de pago era PayPal y aplicar lógica especial
      const storedPaymentMethod = localStorage.getItem('payment_method');
      const paymentVerified = localStorage.getItem('payment_verified_after_return');
      
      if (storedPaymentMethod === 'paypal') {
        // Para PayPal, verificar si hay parámetros en la URL que indiquen éxito
        const urlParams = new URLSearchParams(window.location.search);
        const paypalToken = urlParams.get('token');
        const payerID = urlParams.get('PayerID');
        
        // Si hay token pero no PayerID, es señal de cancelación
        if (paypalToken && !payerID) {
          console.log("Detectada cancelación de PayPal: token presente pero no PayerID");
          setError('El pago con PayPal fue cancelado. Intenta nuevamente o elige otro método de pago.');
          setIsProcessing(false);
          clearPaymentData();
          return;
        }
      }
      
      // Obtener token actualizado
      const token = getAuthToken();
      
      if (!token) {
        throw new Error('No has iniciado sesión');
      }

      const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
      console.log(`Verificando estado. ID Curso: ${cursoId}`);
      
      // Headers con token
      const headers = {
        'Authorization': token,
        'Accept': 'application/json'
      };
      
      // Realizar la solicitud
      const response = await fetch(`${apiUrl}/api/pagos/${cursoId}`, {
        headers: headers,
        credentials: 'include'
      });

      // Verificar errores
      if (!response.ok) {
        const errorText = await response.text();
        console.error("Error en verificación:", errorText);
        throw new Error(errorText || `Error ${response.status}: ${response.statusText}`);
      }

      // Procesar respuesta
      const data = await response.json();
      console.log('Estado del pago:', data);

      // Obtener estado
      const estadoPago = data.estado || 'pendiente';

      // Manejar según estado
      if (estadoPago === 'aprobado') {
        // Pago exitoso: limpiar datos y notificar
        clearPaymentData();
        setIsProcessing(false);
        onSuccess();
        
        // Redirigir si es necesario
        if (paymentMethod === 'paypal') {
          const redirectCursoId = localStorage.getItem('course_redirect_id') || cursoId;
          localStorage.removeItem('course_redirect_id');
          window.location.href = `/cursos/${redirectCursoId}`;
        }
      } else if (estadoPago === 'rechazado') {
        // Pago rechazado
        setError('El pago fue rechazado. Por favor intenta con otro método.');
        setIsProcessing(false);
        clearPaymentData();
      } else if (estadoPago === 'pendiente') {
        // Si estamos verificando después de regresar de la pasarela de pago y sigue pendiente,
        // consideramos que fue cancelado o no procesado correctamente
        if (paymentVerified === 'true') {
          console.log("Pago sigue pendiente después de verificación post-retorno. Considerando cancelado.");
          setError('El pago no pudo ser procesado. Por favor, intenta con otro método de pago.');
          setIsProcessing(false);
          clearPaymentData();
          return;
        }
        
        // Verificar número de intentos y tiempo transcurrido
        const maxChecks = parseInt(localStorage.getItem('payment_check_count') || '0', 10);
        const MAX_CHECKS = 5; // Reducido de 10 a 5 intentos máximos
        
        // Para PayPal, aplicamos un límite más estricto cuando regresamos de redirección
        const isReturningFromPayment = localStorage.getItem('returning_from_payment') === 'true';
        const isPayPal = storedPaymentMethod === 'paypal';
        
        if ((isPayPal && isReturningFromPayment && maxChecks >= 2) || maxChecks >= MAX_CHECKS) {
          // Si es PayPal con redirección y ya verificamos 2 veces, o si cualquier método supera MAX_CHECKS
          setError('El pago fue cancelado o no se completó. Por favor intenta nuevamente.');
          setIsProcessing(false);
          clearPaymentData();
        } else {
          // Incrementar contador y programar nueva verificación con tiempo incremental
          const nextCount = maxChecks + 1;
          localStorage.setItem('payment_check_count', String(nextCount));
          
          // Incrementar gradualmente el tiempo entre verificaciones (3s, 4s, 5s, etc.)
          const delay = 3000 + (nextCount * 1000);
          setTimeout(() => checkPaymentStatus(cursoId), delay);
        }
      } else {
        // Estado desconocido
        console.error(`Estado no reconocido: "${estadoPago}"`, data);
        setError(`Estado de pago desconocido: ${estadoPago}. Contacta a soporte.`);
        setIsProcessing(false);
        clearPaymentData();
      }
    } catch (err) {
      console.error('Error en verificación:', err);
      
      // Mensajes específicos
      let errorMessage = 'Error al verificar el pago';
      
      if (err.message.includes('Failed to fetch') || err.message.includes('NetworkError')) {
        errorMessage = 'No se pudo conectar con el servidor. Verifica tu conexión.';
      } else if (err.message.includes('HTML')) {
        errorMessage = 'El servidor no está respondiendo correctamente. Contacta al administrador.';
      } else if (err.message.includes('token') || err.message.includes('autorización') || err.message.includes('401')) {
        errorMessage = 'Tu sesión ha expirado. Por favor inicia sesión nuevamente.';
      }
      
      setError(errorMessage);
      setIsProcessing(false);
      clearPaymentData();
    }
  };

  // Función mejorada para limpiar datos de pago en localStorage
  const clearPaymentData = () => {
    const keysToRemove = [
      'payment_check_count',
      'current_payment_id',
      'returning_from_payment',
      'payment_curso_id',
      'course_redirect_id',
      'payment_start_time',
      'payment_method',
      'payment_verified_after_return'
    ];
    
    keysToRemove.forEach(key => localStorage.removeItem(key));
  };

  // Formateo automático de inputs de tarjeta
  const handleCardChange = (e) => {
    const { name, value } = e.target;
    
    let processedValue = value;
    
    if (name === 'number') {
      // Formatear número de tarjeta en grupos de 4 dígitos
      const cleaned = value.replace(/\D/g, '').substring(0, 16);
      processedValue = cleaned.replace(/(\d{4})(?=\d)/g, '$1 ');
    } else if (name === 'expiry') {
      // Formatear fecha de expiración como MM/AA
      const cleaned = value.replace(/\D/g, '').substring(0, 4);
      processedValue = cleaned.length > 2 
        ? cleaned.substring(0, 2) + '/' + cleaned.substring(2) 
        : cleaned;
    } else if (name === 'cvv') {
      // Limitar CVV a 3-4 dígitos
      processedValue = value.replace(/\D/g, '').substring(0, 4);
    }
    
    setCardDetails(prev => ({
      ...prev,
      [name]: processedValue
    }));
  };
  
  // Renderizar mensaje durante procesamiento con opción para cancelar
  const renderProcessingMessage = () => {
    return (
      <div className="redirect-message">
        <div className="redirect-icon">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
            <polyline points="15 3 21 3 21 9"></polyline>
            <line x1="10" y1="14" x2="21" y2="3"></line>
          </svg>
        </div>
        <h4>Verificando pago</h4>
        <p>Estamos verificando el estado de tu pago, espera un momento...</p>
        <div className="loading-spinner"></div>
        <div className="processing-actions">
          <button 
            type="button"
            className="retry-button"
            onClick={() => {
              localStorage.removeItem('payment_check_count');
              checkPaymentStatus(curso.id);
            }}
          >
            Verificar nuevamente
          </button>
          <button 
            type="button"
            className="cancel-button"
            onClick={() => {
              clearPaymentData();
              setIsProcessing(false);
              setError('El proceso de pago ha sido cancelado.');
            }}
          >
            Cancelar proceso
          </button>
        </div>
      </div>
    );
  };

  // Renderizar ícono para el método de pago
  const renderPaymentMethodIcon = (method) => {
    const icons = {
      tarjeta: (
        <svg className="payment-method-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <rect x="1" y="4" width="22" height="16" rx="2" ry="2"></rect>
          <line x1="1" y1="10" x2="23" y2="10"></line>
        </svg>
      ),
      paypal: (
        <svg className="payment-method-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M6.5 8c1.5 0 2.5 1.1 2.5 2.5 0 1.4-1 2.5-2.5 2.5h-3"></path>
          <path d="M6.5 13v6"></path>
          <path d="M13 8c2.2 0 3.5 1.1 3.5 2.5 0 1.4-1.3 2.5-3.5 2.5h-3"></path>
          <path d="M13 13v6"></path>
          <path d="M18.5 8c1.5 0 2.5 1.1 2.5 2.5 0 1.4-1 2.5-2.5 2.5h-3"></path>
          <path d="M18.5 13v6"></path>
        </svg>
      ),
      stripe: (
        <svg className="payment-method-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M2 12h20"></path>
          <path d="M2 7h20"></path>
          <path d="M2 17h20"></path>
        </svg>
      ),
      coinbase: (
        <svg className="payment-method-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="8" x2="12" y2="16"></line>
          <line x1="8" y1="12" x2="16" y2="12"></line>
        </svg>
      ),
      mercadopago: (
        <svg className="payment-method-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <rect x="2" y="5" width="20" height="14" rx="2"></rect>
          <line x1="2" y1="10" x2="22" y2="10"></line>
        </svg>
      ),
      transferencia: (
        <svg className="payment-method-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <line x1="12" y1="2" x2="12" y2="22"></line>
          <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
        </svg>
      )
    };
    
    return icons[method] || null;
  };
  
  // Renderizar componente específico según método de pago
  const renderPaymentMethodComponent = () => {
    if (devMode) return null;
    
    const methodComponents = {
      tarjeta: (
        <div className="card-details">
          <div className="form-group">
            <label>Número de tarjeta</label>
            <div className="input-with-icon">
              <input 
                type="text" 
                name="number"
                placeholder="1234 5678 9012 3456" 
                value={cardDetails.number}
                onChange={handleCardChange}
                disabled={isProcessing}
                required
              />
              <div className="input-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <rect x="1" y="4" width="22" height="16" rx="2" ry="2"></rect>
                  <line x1="1" y1="10" x2="23" y2="10"></line>
                </svg>
              </div>
            </div>
          </div>
          <div className="form-row">
            <div className="form-group">
              <label>Fecha expiración (MM/AA)</label>
              <div className="input-with-icon">
                <input 
                  type="text" 
                  name="expiry"
                  placeholder="MM/AA" 
                  value={cardDetails.expiry}
                  onChange={handleCardChange}
                  disabled={isProcessing}
                  required
                />
                <div className="input-icon">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                    <line x1="16" y1="2" x2="16" y2="6"></line>
                    <line x1="8" y1="2" x2="8" y2="6"></line>
                    <line x1="3" y1="10" x2="21" y2="10"></line>
                  </svg>
                </div>
              </div>
            </div>
            <div className="form-group">
              <label>CVV</label>
              <div className="input-with-icon">
                <input 
                  type="text" 
                  name="cvv"
                  placeholder="123" 
                  value={cardDetails.cvv}
                  onChange={handleCardChange}
                  disabled={isProcessing}
                  required
                />
                <div className="input-icon">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M12 2a10 10 0 1 0 10 10H12V2z"></path>
                    <path d="M21.2 8.4c.5 1.2.8 2.5.8 3.8"></path>
                  </svg>
                </div>
              </div>
            </div>
          </div>
          <div className="card-brands">
            <img src="/images/visa.svg" alt="Visa" />
            <img src="/images/mastercard.svg" alt="Mastercard" />
            <img src="/images/amex.svg" alt="American Express" />
          </div>
        </div>
      ),
      
      paypal: (
        <div className="payment-info">
          <div className="payment-brand">
            <img src="/images/paypal.svg" alt="PayPal" className="payment-logo" />
          </div>
          <p>Serás redirigido a PayPal para completar el pago de forma segura.</p>
          <div className="payment-advantage">
            <div className="advantage-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                <polyline points="22 4 12 14.01 9 11.01"></polyline>
              </svg>
            </div>
            <span>Rápido, seguro y sin compartir datos de tu tarjeta</span>
          </div>
          <p className="redirect-note"><strong>Nota:</strong> Si cancelas el proceso o cierras la ventana de PayPal, deberás iniciar el proceso nuevamente.</p>
        </div>
      ),
      
      stripe: (
        <div className="payment-info">
          <div className="payment-brand">
            <img src="/images/stripe.svg" alt="Stripe" className="payment-logo" />
          </div>
          <p>Serás redirigido a una página segura de Stripe para completar el pago.</p>
          <div className="payment-advantage">
            <div className="advantage-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
            </div>
            <span>Pago encriptado y certificado PCI DSS nivel 1</span>
          </div>
          <p className="redirect-note"><strong>Nota:</strong> Después de completar el pago, volverás automáticamente a esta página.</p>
        </div>
      ),
      
      coinbase: (
        <div className="payment-info">
          <div className="payment-brand">
            <img src="/images/coinbase.svg" alt="Coinbase" className="payment-logo" />
          </div>
          <p>Serás redirigido a Coinbase Commerce para pagar con criptomonedas.</p>
          <div className="payment-advantage">
            <div className="advantage-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10"></circle>
                <path d="M8 14s1.5 2 4 2 4-2 4-2"></path>
                <line x1="9" y1="9" x2="9.01" y2="9"></line>
                <line x1="15" y1="9" x2="15.01" y2="9"></line>
              </svg>
            </div>
            <span>Acepta Bitcoin, Ethereum, Litecoin y otras cryptos</span>
          </div>
          <p className="crypto-info">Transacciones rápidas y seguras con verificación en blockchain</p>
          <p className="redirect-note"><strong>Nota:</strong> Después de completar el pago, volverás automáticamente a esta página.</p>
        </div>
      ),
      
      mercadopago: (
        <div className="payment-info">
          <div className="payment-brand">
            <img src="/images/mercadopago.svg" alt="Mercado Pago" className="payment-logo" />
          </div>
          <p>Serás redirigido a Mercado Pago para completar la transacción.</p>
          <div className="payment-advantage">
            <div className="advantage-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>
              </svg>
            </div>
            <span>Pago rápido, cuotas y múltiples métodos locales</span>
          </div>
          <p className="redirect-note"><strong>Nota:</strong> Después de completar el pago, volverás automáticamente a esta página.</p>
        </div>
      ),
      
      transferencia: (
        <div className="payment-info">
          <div className="bank-info">
            <h4>Datos bancarios:</h4>
            <div className="bank-field">
              <div className="bank-field-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect>
                  <line x1="2" y1="7" x2="22" y2="7"></line>
                </svg>
              </div>
              <p><strong>Banco:</strong> Nombre del Banco</p>
            </div>
            <div className="bank-field">
              <div className="bank-field-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
              </div>
              <p><strong>Titular:</strong> Nombre de la Empresa</p>
            </div>
            <div className="bank-field">
              <div className="bank-field-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <rect x="3" y="5" width="18" height="14" rx="2"></rect>
                  <line x1="3" y1="10" x2="21" y2="10"></line>
                </svg>
              </div>
              <p><strong>Cuenta:</strong> XXXX-XXXX-XXXX-XXXX</p>
            </div>
            <div className="bank-field">
              <div className="bank-field-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <line x1="16.5" y1="9.4" x2="7.5" y2="4.21"></line>
                  <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
                  <polyline points="3.29 7 12 12 20.71 7"></polyline>
                  <line x1="12" y1="22" x2="12" y2="12"></line>
                </svg>
              </div>
              <p><strong>CBU/CLABE:</strong> XXXXXXXXXXXXXXXXX</p>
            </div>
            <div className="transfer-note">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
              <p>Importante: Incluye tu número de pedido en la referencia de la transferencia.</p>
            </div>
          </div>
        </div>
      )
    };
    
    return methodComponents[paymentMethod] || null;
  };

  return (
    <div className="payment-modal-overlay">
      <div className={`payment-modal ${isProcessing ? 'processing' : ''}`}>
        <div className="modal-header">
          <h3>Pagar curso: {curso.titulo}</h3>
          <button 
            onClick={onClose}
            className="close-button"
            disabled={isProcessing}
            aria-label="Cerrar"
          >
            &times;
          </button>
        </div>

        {error && (
          <div className="error-message">
            <div className="error-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
            </div>
            <p>{error}</p>
            {error.includes('sesión') && (
              <button 
                className="relogin-button"
                onClick={() => {
                  window.location.href = '/login?redirect=' + encodeURIComponent(window.location.pathname);
                }}
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="button-icon">
                  <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path>
                  <polyline points="10 17 15 12 10 7"></polyline>
                  <line x1="15" y1="12" x2="3" y2="12"></line>
                </svg>
                Iniciar sesión
              </button>
            )}
            {error.includes('cancelado') || error.includes('no se completó') ? (
              <button 
                className="retry-payment-button"
                onClick={() => {
                  setError(null);
                  setIsProcessing(false);
                }}
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="button-icon">
                  <polyline points="23 4 23 10 17 10"></polyline>
                  <polyline points="1 20 1 14 7 14"></polyline>
                  <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
                </svg>
                Intentar nuevamente
              </button>
            ) : null}
          </div>
        )}
        
        {isProcessing && !error && renderProcessingMessage()}

        {!isProcessing && !error && (
          <div className="modal-content">
            <div className="payment-details">
              <div className="course-title">{curso.titulo}</div>
              <p className="course-price">${curso.precio?.toFixed(2) || '29.99'}</p>
              {devMode && (
                <div className="dev-mode-notice">
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="dev-icon">
                    <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
                  </svg>
                  <p><strong>Modo Desarrollo Activo</strong>: Los pagos serán simulados.</p>
                </div>
              )}
            </div>
            
            <form onSubmit={handleSubmit}>
              <div className="form-group payment-method-selector">
                <label>Método de pago:</label>
                <div className="select-wrapper">
                  <select 
                    className="payment-method-dropdown"
                    value={paymentMethod}
                    onChange={(e) => {
                      setPaymentMethod(e.target.value);
                      localStorage.removeItem('payment_check_count');
                    }}
                    disabled={isProcessing}
                    required
                  >
                    <option value="tarjeta">💳 Tarjeta de crédito/débito</option>
                    <option value="paypal">🔄 PayPal</option>
                    <option value="stripe">💵 Stripe</option>
                    <option value="coinbase">₿ Coinbase (Crypto)</option>
                    <option value="mercadopago">📱 Mercado Pago</option>
                    <option value="transferencia">🏦 Transferencia bancaria</option>
                  </select>
                  <div className="select-arrow"></div>
                  {renderPaymentMethodIcon(paymentMethod)}
                </div>
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
                    <span>Procesando pago...</span>
                  </>
                ) : (
                  <>
                    <span>Pagar ahora</span>
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="arrow-icon">
                      <line x1="5" y1="12" x2="19" y2="12"></line>
                      <polyline points="12 5 19 12 12 19"></polyline>
                    </svg>
                  </>
                )}
              </button>
            </form>
            
            <div className="secure-payment-info">
              <div className="secure-icon">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                </svg>
              </div>
              <p>Pago 100% seguro. Tus datos están protegidos con encriptación SSL.</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default PaymentModal;