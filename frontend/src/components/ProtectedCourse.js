import React, { useEffect, useState } from 'react';
import { Navigate, useParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import PaymentModal from './PaymentModal';

const ProtectedCourse = ({ children }) => {
  const { id } = useParams();
  const { isLoggedIn, loading, verificarAccesoCurso } = useAuth();
  const [tieneAcceso, setTieneAcceso] = useState(null);
  const [showPayment, setShowPayment] = useState(false);

  useEffect(() => {
    const checkAccess = async () => {
      if (!isLoggedIn) {
        setTieneAcceso(false);
        return;
      }
      
      const acceso = await verificarAccesoCurso(id);
      setTieneAcceso(acceso);
      if (!acceso) setShowPayment(true);
    };
    
    checkAccess();
  }, [id, isLoggedIn, verificarAccesoCurso]);

  if (loading || tieneAcceso === null) {
    return <div className="loading">Cargando...</div>;
  }

  if (!isLoggedIn) {
    return <Navigate to="/login" state={{ from: `/curso/${id}` }} replace />;
  }

  if (!tieneAcceso) {
    return (
      <>
        <PaymentModal 
          cursoId={id}
          show={showPayment}
          onClose={() => setShowPayment(false)}
          onSuccess={() => setTieneAcceso(true)}
        />
        <div className="access-denied">
          <h2>Acceso no autorizado</h2>
          <p>Debes comprar este curso para ver su contenido.</p>
          <button 
            className="btn btn-primary"
            onClick={() => setShowPayment(true)}
          >
            Comprar Curso
          </button>
        </div>
      </>
    );
  }

  return children;
};

export default ProtectedCourse;