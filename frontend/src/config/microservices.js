// Configuración centralizada de microservicios
// Este archivo facilita el deployment independiente de cada microservicio

const config = {
  // URLs base de los microservicios
  services: {
    auth: process.env.REACT_APP_AUTH_API_URL || 'http://localhost:8001',
    payment: process.env.REACT_APP_PAYMENT_API_URL || 'http://localhost:8002', 
    course: process.env.REACT_APP_COURSE_API_URL || 'http://localhost:8003',
    contact: process.env.REACT_APP_CONTACT_API_URL || 'http://localhost:8004',
    portfolio: process.env.REACT_APP_PORTFOLIO_API_URL || 'http://localhost:8005',
    home: process.env.REACT_APP_HOME_API_URL || 'http://localhost:8006',
    analytics: process.env.REACT_APP_ANALYTICS_API_URL || 'http://localhost:8007'
  },

  // Configuración de timeout para las peticiones
  timeout: {
    default: 30000, // 30 segundos
    upload: 60000,  // 60 segundos para uploads
    auth: 15000     // 15 segundos para autenticación
  },

  // Headers comunes
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  },

  // Configuración de retry para peticiones fallidas
  retry: {
    attempts: 3,
    delay: 1000 // 1 segundo entre intentos
  }
};

// Función para obtener la URL completa de un endpoint
export const getServiceUrl = (service, endpoint = '') => {
  const baseUrl = config.services[service];
  if (!baseUrl) {
    console.warn(`Servicio "${service}" no encontrado en la configuración`);
    return '';
  }
  return `${baseUrl}${endpoint.startsWith('/') ? endpoint : '/' + endpoint}`;
};

// Función para verificar si un servicio está disponible
export const isServiceAvailable = async (service) => {
  try {
    const url = getServiceUrl(service, '/health');
    const response = await fetch(url, {
      method: 'GET',
      timeout: 5000
    });
    return response.ok;
  } catch (error) {
    console.warn(`Servicio ${service} no disponible:`, error.message);
    return false;
  }
};

// Configuración para diferentes entornos
export const getEnvironmentConfig = () => {
  const env = process.env.NODE_ENV || 'development';
  
  switch (env) {
    case 'production':
      return {
        ...config,
        services: {
          auth: process.env.REACT_APP_AUTH_API_URL || 'https://auth.midominio.com',
          payment: process.env.REACT_APP_PAYMENT_API_URL || 'https://payment.midominio.com',
          course: process.env.REACT_APP_COURSE_API_URL || 'https://course.midominio.com',
          contact: process.env.REACT_APP_CONTACT_API_URL || 'https://contact.midominio.com',
          portfolio: process.env.REACT_APP_PORTFOLIO_API_URL || 'https://portfolio.midominio.com',
          home: process.env.REACT_APP_HOME_API_URL || 'https://home.midominio.com',
          analytics: process.env.REACT_APP_ANALYTICS_API_URL || 'https://analytics.midominio.com'
        }
      };
    case 'staging':
      return {
        ...config,
        services: {
          auth: process.env.REACT_APP_AUTH_API_URL || 'https://auth-staging.midominio.com',
          payment: process.env.REACT_APP_PAYMENT_API_URL || 'https://payment-staging.midominio.com',
          course: process.env.REACT_APP_COURSE_API_URL || 'https://course-staging.midominio.com',
          contact: process.env.REACT_APP_CONTACT_API_URL || 'https://contact-staging.midominio.com',
          portfolio: process.env.REACT_APP_PORTFOLIO_API_URL || 'https://portfolio-staging.midominio.com',
          home: process.env.REACT_APP_HOME_API_URL || 'https://home-staging.midominio.com',
          analytics: process.env.REACT_APP_ANALYTICS_API_URL || 'https://analytics-staging.midominio.com'
        }
      };
    default:
      return config;
  }
};

export default config;