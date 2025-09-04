import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import './Cursos.css';
import CursosService from './services/CursosService';

function Cursos() {
  const [cursos, setCursos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [cursosAcceso, setCursosAcceso] = useState({});
  const [verificandoAccesos, setVerificandoAccesos] = useState(false);

  // Cursos estáticos para fallback
  const cursosDefault = [
    { 
      id: 1, 
      titulo: 'Curso de React', 
      descripcion: 'Aprende React desde cero.', 
      imagen_url: 'https://via.placeholder.com/300', 
      duracion: '10 horas', 
      nivel: 'Principiante', 
      estado: 'Publicado', 
      precio: 29.99,
      gratuito: false
    },
    { 
      id: 2, 
      titulo: 'Curso de Node.js', 
      descripcion: 'Domina el backend con Node.js.', 
      imagen_url: 'https://via.placeholder.com/300', 
      duracion: '15 horas', 
      nivel: 'Intermedio', 
      estado: 'Publicado', 
      precio: 39.99,
      gratuito: false
    },
    { 
      id: 3, 
      titulo: 'Curso de Diseño Web', 
      descripcion: 'Crea diseños modernos y responsivos.', 
      imagen_url: 'https://via.placeholder.com/300', 
      duracion: '12 horas', 
      nivel: 'Principiante', 
      estado: 'Publicado', 
      precio: 24.99,
      gratuito: false
    },
  ];

  // Función para obtener la URL base de la API del Course Service
  const getCourseApiBaseUrl = () => {
    return process.env.REACT_APP_COURSE_API_URL || 'http://localhost:8003';
  };

  // Función mejorada para verificar y corregir URLs de imágenes
  const getImageUrl = (url) => {
    if (!url) return 'https://via.placeholder.com/300';
    
    if (url.startsWith('http') || url.startsWith('data:')) {
      return url;
    }
    
    if (url.startsWith('/static/')) {
      const baseUrl = getCourseApiBaseUrl();
      return `${baseUrl}${url}`;
    }
    
    return url;
  };
  
  // **FUNCIÓN ACTUALIZADA: Verificar acceso usando la nueva integración con Payment Service**
  const verificarAccesoCursos = async (cursosData) => {
    try {
      console.log('Verificando acceso a cursos con Payment Service...');
      setVerificandoAccesos(true);
      
      // Verificar si hay un token (usuario autenticado)
      const token = localStorage.getItem('token');
      
      if (!token) {
        console.log('Usuario no autenticado, solo cursos gratuitos tendrán acceso');
        // Para usuarios no autenticados, solo los cursos gratuitos tienen acceso
        const accesoObj = cursosData.reduce((acc, curso) => {
          acc[curso.id] = curso.gratuito || curso.precio === 0;
          return acc;
        }, {});
        return accesoObj;
      }
      
      // **USAR NUEVA FUNCIÓN INTEGRADA PARA VERIFICAR ACCESO**
      const accesoPromises = cursosData.map(async (curso) => {
        try {
          console.log(`Verificando acceso para curso ${curso.id}: ${curso.titulo}`);
          
          // Si el curso es gratuito, tiene acceso automáticamente
          if (curso.gratuito || curso.precio === 0) {
            return {
              cursoId: curso.id,
              tieneAcceso: true,
              esGratuito: true,
              estadoPago: 'gratuito'
            };
          }
          
          // **USAR SERVICIO ACTUALIZADO**
          const verificacion = await CursosService.verificarAccesoCurso(curso.id);
          console.log(`Resultado verificación curso ${curso.id}:`, verificacion);
          
          return {
            cursoId: curso.id,
            tieneAcceso: verificacion.has_access || false,
            esGratuito: false,
            estadoPago: verificacion.payment_status || 'no_verificado',
            error: verificacion.error || null
          };
        } catch (error) {
          console.log(`Error al verificar acceso del curso ${curso.id}:`, error);
          return { 
            cursoId: curso.id, 
            tieneAcceso: false, 
            esGratuito: curso.gratuito || curso.precio === 0,
            estadoPago: 'error',
            error: error.message
          };
        }
      });

      const resultadosAcceso = await Promise.all(accesoPromises);
      
      // Convertir el array de resultados a un objeto para fácil acceso
      const accesoObj = resultadosAcceso.reduce((acc, item) => {
        acc[item.cursoId] = {
          tiene_acceso: item.tieneAcceso,
          es_gratuito: item.esGratuito,
          estado_pago: item.estadoPago,
          error: item.error
        };
        return acc;
      }, {});
      
      console.log('Estado de acceso a cursos actualizado:', accesoObj);
      
      return accesoObj;
    } catch (error) {
      console.error('Error al verificar accesos:', error);
      // En caso de error, asumir que no hay acceso excepto para cursos gratuitos
      const accesoObj = cursosData.reduce((acc, curso) => {
        acc[curso.id] = {
          tiene_acceso: curso.gratuito || curso.precio === 0,
          es_gratuito: curso.gratuito || curso.precio === 0,
          estado_pago: 'error',
          error: error.message
        };
        return acc;
      }, {});
      return accesoObj;
    } finally {
      setVerificandoAccesos(false);
    }
  };

  useEffect(() => {
    const fetchCursos = async () => {
      try {
        setLoading(true);
        
        // Usar el servicio actualizado para obtener cursos
        const cursosData = await CursosService.getCursos(true);
        console.log('Cursos obtenidos:', cursosData);
        
        let cursosParaUsar = [];
        
        // Verificar si la respuesta contiene datos válidos
        if (Array.isArray(cursosData) && cursosData.length > 0) {
          cursosParaUsar = cursosData;
        } else {
          console.warn('No se obtuvieron cursos del servidor, usando datos por defecto');
          cursosParaUsar = cursosDefault;
        }
        
        // Filtrar solo los cursos publicados (por si acaso)
        const cursosPublicados = cursosParaUsar.filter(curso => 
          curso.estado === 'Publicado' || curso.estado === 'publicado'
        );
        
        console.log('Cursos publicados:', cursosPublicados);
        
        // Si no hay cursos publicados, mostrar mensaje informativo
        if (cursosPublicados.length === 0) {
          setError('No hay cursos publicados disponibles en este momento.');
        }
        
        setCursos(cursosPublicados);
        
        // **VERIFICAR ESTADO DE ACCESO USANDO NUEVA INTEGRACIÓN**
        if (cursosPublicados.length > 0) {
          const accesos = await verificarAccesoCursos(cursosPublicados);
          setCursosAcceso(accesos);
        }
        
        setError(null);
      } catch (err) {
        console.error('Error al cargar cursos:', err);
        setError('No se pudieron cargar los cursos. Por favor, intenta de nuevo más tarde.');
        // En caso de error, usar los cursos por defecto
        setCursos(cursosDefault);
        
        // Verificar acceso para cursos por defecto también
        const accesos = await verificarAccesoCursos(cursosDefault);
        setCursosAcceso(accesos);
      } finally {
        setLoading(false);
      }
    };

    fetchCursos();
  }, []);

  // **ACTUALIZAR VERIFICACIÓN PERIÓDICA PARA PAGOS PENDIENTES**
  useEffect(() => {
    const interval = setInterval(async () => {
      const token = localStorage.getItem('token');
      if (token && cursos.length > 0) {
        console.log('Verificación periódica de accesos...');
        
        // Solo verificar cursos que tienen pagos pendientes
        const cursosPendientes = cursos.filter(curso => {
          const acceso = cursosAcceso[curso.id];
          return acceso && acceso.estado_pago === 'pendiente';
        });
        
        if (cursosPendientes.length > 0) {
          console.log(`Verificando ${cursosPendientes.length} cursos con pagos pendientes`);
          const accesosActualizados = await verificarAccesoCursos(cursosPendientes);
          
          // Actualizar solo los cursos que cambiaron
          setCursosAcceso(prev => ({
            ...prev,
            ...accesosActualizados
          }));
        }
      }
    }, 30000); // Verificar cada 30 segundos
    
    return () => clearInterval(interval);
  }, [cursos, cursosAcceso]);

  // **FUNCIÓN HELPER: Obtener texto de estado de acceso**
  const getAccessStatusText = (curso) => {
    const acceso = cursosAcceso[curso.id];
    
    if (!acceso) {
      return 'Verificando...';
    }
    
    if (acceso.es_gratuito) {
      return 'Acceder al Curso';
    }
    
    if (acceso.tiene_acceso) {
      return 'Acceder al Curso';
    }
    
    switch (acceso.estado_pago) {
      case 'pendiente':
        return 'Pago Pendiente';
      case 'rechazado':
        return 'Pago Rechazado - Reintentar';
      case 'no_pagado':
      case 'no_verificado':
        return 'Más Información';
      case 'error':
        return 'Error - Más Información';
      default:
        return 'Más Información';
    }
  };

  // **FUNCIÓN HELPER: Obtener clase CSS para el botón**
  const getAccessButtonClass = (curso) => {
    const acceso = cursosAcceso[curso.id];
    
    if (!acceso) {
      return 'curso-cta';
    }
    
    if (acceso.es_gratuito || acceso.tiene_acceso) {
      return 'curso-cta access-granted';
    }
    
    switch (acceso.estado_pago) {
      case 'pendiente':
        return 'curso-cta payment-pending';
      case 'rechazado':
        return 'curso-cta payment-rejected';
      default:
        return 'curso-cta';
    }
  };

  if (loading) {
    return <div className="cursos-loading">Cargando cursos...</div>;
  }

  if (error && cursos.length === 0) {
    return <div className="cursos-error">{error}</div>;
  }

  if (cursos.length === 0) {
    return <div className="cursos-empty">No hay cursos publicados disponibles en este momento.</div>;
  }

  return (
    <div className="cursos-container">
      {verificandoAccesos && (
        <div className="verification-notice">
          🔄 Verificando acceso a cursos...
        </div>
      )}
      
      <ul className="cursos-list">
        {cursos.map((curso) => {
          const acceso = cursosAcceso[curso.id];
          
          return (
            <li key={curso.id} className="curso-item animate__animated animate__fadeInUp">
              <div className="curso-image-container">
                <img 
                  src={getImageUrl(curso.imagen_url)} 
                  alt={curso.titulo} 
                  className="curso-image" 
                  onError={(e) => {
                    console.log(`Error al cargar imagen para curso ${curso.id}:`, curso.imagen_url);
                    e.target.onerror = null;
                    e.target.src = 'https://via.placeholder.com/300';
                  }}
                />
                <span className="curso-badge">{curso.nivel || 'Todos los niveles'}</span>
                
                {/* **BADGES MEJORADOS PARA DIFERENTES ESTADOS** */}
                {(curso.gratuito || curso.precio === 0) && (
                  <span className="curso-badge curso-badge-free">Gratis</span>
                )}
                
                {acceso && acceso.estado_pago === 'pendiente' && (
                  <span className="curso-badge curso-badge-pending">Pago Pendiente</span>
                )}
                
                {acceso && acceso.tiene_acceso && !acceso.es_gratuito && (
                  <span className="curso-badge curso-badge-owned">Comprado</span>
                )}
                
                {acceso && acceso.estado_pago === 'rechazado' && (
                  <span className="curso-badge curso-badge-rejected">Pago Rechazado</span>
                )}
              </div>
              
              <Link to={`/curso/${curso.id}`} className="curso-link">
                {curso.titulo}
              </Link>
              <p className="curso-descripcion">{curso.descripcion}</p>
              
              <div className="curso-footer">
                <div className="curso-meta">
                  <span className="curso-meta-item">
                    <i className="curso-meta-icon">⏳</i> 
                    {curso.duracion || (curso.total_chapters ? `${curso.total_chapters} capítulos` : 'Curso completo')}
                  </span>
                  
                  {/* **MOSTRAR PRECIO O ESTADO SEGÚN ACCESO** */}
                  {!curso.gratuito && curso.precio > 0 && (
                    <>
                      {!acceso?.tiene_acceso && (
                        <span className="curso-meta-item curso-precio">
                          <i className="curso-meta-icon">💰</i> ${curso.precio?.toFixed(2) || '29.99'}
                        </span>
                      )}
                      
                      {acceso?.estado_pago === 'pendiente' && (
                        <span className="curso-meta-item curso-pending">
                          <i className="curso-meta-icon">⏳</i> Verificando pago...
                        </span>
                      )}
                      
                      {acceso?.tiene_acceso && (
                        <span className="curso-meta-item curso-owned">
                          <i className="curso-meta-icon">✅</i> Comprado
                        </span>
                      )}
                    </>
                  )}
                  
                  {(curso.gratuito || curso.precio === 0) && (
                    <span className="curso-meta-item curso-gratis">
                      <i className="curso-meta-icon">🎉</i> GRATIS
                    </span>
                  )}
                </div>
                
                <Link 
                  to={`/curso/${curso.id}`} 
                  className={getAccessButtonClass(curso)}
                >
                  {getAccessStatusText(curso)}
                </Link>
              </div>
              
              {/* **MOSTRAR INFORMACIÓN ADICIONAL PARA ESTADOS ESPECIALES** */}
              {acceso && acceso.estado_pago === 'pendiente' && (
                <div className="curso-status-info pending">
                  <small>⏳ Tu pago está siendo procesado. Te notificaremos cuando esté listo.</small>
                </div>
              )}
              
              {acceso && acceso.estado_pago === 'rechazado' && (
                <div className="curso-status-info rejected">
                  <small>❌ El pago fue rechazado. Puedes intentar con otro método.</small>
                </div>
              )}
              
              {acceso && acceso.error && acceso.estado_pago === 'error' && (
                <div className="curso-status-info error">
                  <small>⚠️ Error al verificar acceso. Intenta recargar la página.</small>
                </div>
              )}
            </li>
          );
        })}
      </ul>
      
      {error && <div className="cursos-error">{error}</div>}
    </div>
  );
}

export default Cursos;