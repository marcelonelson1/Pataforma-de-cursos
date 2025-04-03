import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { 
  faGraduationCap, faShoppingCart, faUsers, faComments,
  faChartLine, faArrowUp, faArrowDown
} from '@fortawesome/free-solid-svg-icons';
import {
  LineChart, Line, BarChart, Bar, XAxis, YAxis, CartesianGrid, 
  Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell
} from 'recharts';
import './AdminPages.css';

function AdminDashboard() {
  const [dashboardData, setDashboardData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const COLORS = ['#0acf97', '#4a6cf7', '#ffbc00', '#fa5c7c'];

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        const token = localStorage.getItem('token');
        const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';
        
        const response = await fetch(`${apiUrl}/api/admin/dashboard`, {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });

        if (!response.ok) {
          throw new Error('Error al obtener los datos del dashboard');
        }

        const data = await response.json();
        setDashboardData(data);
        setError(null);
      } catch (err) {
        console.error('Error al cargar datos del dashboard:', err);
        setError('No se pudieron cargar los datos del dashboard. Por favor, intenta de nuevo más tarde.');
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  // Datos de ejemplo mientras se carga la API real
  const demoData = {
    stats: {
      total_ventas: 156,
      total_ingresos: 4589.50,
      ventas_ultima_semana: 32,
      ingresos_ultima_semana: 948.75
    },
    top_cursos: [
      { id: 1, titulo: 'Curso de React', total_ventas: 45, ingresos: 1349.55 },
      { id: 2, titulo: 'Curso de Node.js', total_ventas: 38, ingresos: 1519.62 },
      { id: 3, titulo: 'Curso de Diseño Web', total_ventas: 27, ingresos: 674.73 },
      { id: 4, titulo: 'Curso de Python', total_ventas: 24, ingresos: 599.76 },
      { id: 5, titulo: 'Curso de JavaScript', total_ventas: 22, ingresos: 439.78 }
    ],
    metodos_pago: [
      { metodo: 'Stripe', total: 78 },
      { metodo: 'PayPal', total: 52 },
      { metodo: 'Transferencia', total: 15 },
      { metodo: 'Otro', total: 11 }
    ],
    ultimas_consultas: [
      { id: 1, nombre: 'Ana García', email: 'ana@ejemplo.com', asunto: 'Consulta sobre curso React', estado: 'no_leido', created_at: '2023-09-15T14:30:00Z' },
      { id: 2, nombre: 'Juan Pérez', email: 'juan@ejemplo.com', asunto: 'Problema con video del curso', estado: 'en_proceso', created_at: '2023-09-14T10:15:00Z' },
      { id: 3, nombre: 'María López', email: 'maria@ejemplo.com', asunto: 'Duda sobre certificado', estado: 'respondido', created_at: '2023-09-13T18:45:00Z' },
      { id: 4, nombre: 'Carlos Ruiz', email: 'carlos@ejemplo.com', asunto: 'Solicitud de información', estado: 'no_leido', created_at: '2023-09-13T09:20:00Z' },
      { id: 5, nombre: 'Laura Díaz', email: 'laura@ejemplo.com', asunto: 'Descuento para grupo', estado: 'archivado', created_at: '2023-09-12T15:10:00Z' }
    ],
    ventas_diarias: [
      { fecha: '2023-09-01', ventas: 5, ingresos: 149.95 },
      { fecha: '2023-09-02', ventas: 3, ingresos: 89.97 },
      { fecha: '2023-09-03', ventas: 8, ingresos: 239.92 },
      { fecha: '2023-09-04', ventas: 4, ingresos: 119.96 },
      { fecha: '2023-09-05', ventas: 7, ingresos: 209.93 },
      { fecha: '2023-09-06', ventas: 6, ingresos: 179.94 },
      { fecha: '2023-09-07', ventas: 10, ingresos: 299.90 },
      { fecha: '2023-09-08', ventas: 8, ingresos: 239.92 },
      { fecha: '2023-09-09', ventas: 5, ingresos: 149.95 },
      { fecha: '2023-09-10', ventas: 9, ingresos: 269.91 },
      { fecha: '2023-09-11', ventas: 12, ingresos: 359.88 },
      { fecha: '2023-09-12', ventas: 7, ingresos: 209.93 },
      { fecha: '2023-09-13', ventas: 11, ingresos: 329.89 },
      { fecha: '2023-09-14', ventas: 14, ingresos: 419.86 },
      { fecha: '2023-09-15', ventas: 9, ingresos: 269.91 }
    ]
  };

  // Usar datos de la API o datos de ejemplo si aún está cargando
  const data = dashboardData || demoData;

  const estadoConsultaBadge = (estado) => {
    switch (estado) {
      case 'no_leido':
        return <span className="admin-badge admin-badge-danger">No leído</span>;
      case 'en_proceso':
        return <span className="admin-badge admin-badge-warning">En proceso</span>;
      case 'respondido':
        return <span className="admin-badge admin-badge-success">Respondido</span>;
      case 'archivado':
        return <span className="admin-badge admin-badge-info">Archivado</span>;
      default:
        return <span className="admin-badge">{estado}</span>;
    }
  };

  // Formatear fecha
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-ES', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric'
    });
  };

  if (loading && !dashboardData) {
    return (
      <div className="admin-loading-container">
        <div className="spinner"></div>
        <p>Cargando datos del dashboard...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="admin-error-container">
        <h2>Error</h2>
        <p>{error}</p>
        <button className="admin-btn admin-btn-primary" onClick={() => window.location.reload()}>
          Intentar de nuevo
        </button>
      </div>
    );
  }

  return (
    <div className="admin-dashboard">
      <div className="admin-page-title">
        <h1>Dashboard</h1>
        <p>Bienvenido al panel de administración</p>
      </div>

      {/* Tarjetas de estadísticas */}
      <div className="admin-stats-grid">
        <div className="admin-stat-card">
          <div className="stat-icon">
            <FontAwesomeIcon icon={faShoppingCart} />
          </div>
          <div className="stat-content">
            <h3>Ventas Totales</h3>
            <div className="stat-value">{data.stats.total_ventas}</div>
            <div className="stat-change positive">
              <FontAwesomeIcon icon={faArrowUp} /> {data.stats.ventas_ultima_semana} esta semana
            </div>
          </div>
        </div>

        <div className="admin-stat-card">
          <div className="stat-icon">
            <FontAwesomeIcon icon={faChartLine} />
          </div>
          <div className="stat-content">
            <h3>Ingresos Totales</h3>
            <div className="stat-value">${data.stats.total_ingresos.toFixed(2)}</div>
            <div className="stat-change positive">
              <FontAwesomeIcon icon={faArrowUp} /> ${data.stats.ingresos_ultima_semana.toFixed(2)} esta semana
            </div>
          </div>
        </div>

        <div className="admin-stat-card">
          <div className="stat-icon">
            <FontAwesomeIcon icon={faGraduationCap} />
          </div>
          <div className="stat-content">
            <h3>Cursos Vendidos</h3>
            <div className="stat-value">{data.top_cursos.length}</div>
            <div className="stat-link">
              <Link to="/admin/cursos">Ver todos</Link>
            </div>
          </div>
        </div>

        <div className="admin-stat-card">
          <div className="stat-icon">
            <FontAwesomeIcon icon={faComments} />
          </div>
          <div className="stat-content">
            <h3>Consultas Nuevas</h3>
            <div className="stat-value">
              {data.ultimas_consultas.filter(c => c.estado === 'no_leido').length}
            </div>
            <div className="stat-link">
              <Link to="/admin/consultas">Ver todas</Link>
            </div>
          </div>
        </div>
      </div>

      {/* Gráficos y tablas */}
      <div className="admin-dashboard-content">
        <div className="admin-row">
          <div className="admin-col-8">
            <div className="admin-card">
              <div className="admin-card-header">
                <h2 className="admin-card-title">Ventas por día (últimos 15 días)</h2>
              </div>
              <div className="admin-card-body">
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart
                    data={data.ventas_diarias}
                    margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="fecha" />
                    <YAxis />
                    <Tooltip formatter={(value, name) => [value, name === 'ingresos' ? 'Ingresos ($)' : 'Ventas']} />
                    <Legend />
                    <Line type="monotone" dataKey="ventas" stroke="#4a6cf7" name="Ventas" />
                    <Line type="monotone" dataKey="ingresos" stroke="#0acf97" name="Ingresos ($)" />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </div>
          </div>

          <div className="admin-col-4">
            <div className="admin-card">
              <div className="admin-card-header">
                <h2 className="admin-card-title">Métodos de Pago</h2>
              </div>
              <div className="admin-card-body">
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={data.metodos_pago}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="total"
                      nameKey="metodo"
                      label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                    >
                      {data.metodos_pago.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip formatter={(value, name, props) => [value, props.payload.metodo]} />
                  </PieChart>
                </ResponsiveContainer>
              </div>
            </div>
          </div>
        </div>

        <div className="admin-row">
          <div className="admin-col-6">
            <div className="admin-card">
              <div className="admin-card-header">
                <h2 className="admin-card-title">Top 5 Cursos</h2>
                <Link to="/admin/cursos" className="admin-btn admin-btn-sm admin-btn-secondary">
                  Ver todos
                </Link>
              </div>
              <div className="admin-card-body">
                <div className="admin-table-responsive">
                  <table className="admin-table">
                    <thead>
                      <tr>
                        <th>Curso</th>
                        <th>Ventas</th>
                        <th>Ingresos</th>
                      </tr>
                    </thead>
                    <tbody>
                      {data.top_cursos.map(curso => (
                        <tr key={curso.id}>
                          <td>
                            <Link to={`/admin/cursos/editar/${curso.id}`}>
                              {curso.titulo}
                            </Link>
                          </td>
                          <td>{curso.total_ventas}</td>
                          <td>${curso.ingresos.toFixed(2)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>

          <div className="admin-col-6">
            <div className="admin-card">
              <div className="admin-card-header">
                <h2 className="admin-card-title">Últimas Consultas</h2>
                <Link to="/admin/consultas" className="admin-btn admin-btn-sm admin-btn-secondary">
                  Ver todas
                </Link>
              </div>
              <div className="admin-card-body">
                <div className="admin-table-responsive">
                  <table className="admin-table">
                    <thead>
                      <tr>
                        <th>Usuario</th>
                        <th>Asunto</th>
                        <th>Estado</th>
                        <th>Fecha</th>
                      </tr>
                    </thead>
                    <tbody>
                      {data.ultimas_consultas.map(consulta => (
                        <tr key={consulta.id}>
                          <td>{consulta.nombre}</td>
                          <td>
                            <Link to={`/admin/consultas?id=${consulta.id}`}>
                              {consulta.asunto}
                            </Link>
                          </td>
                          <td>{estadoConsultaBadge(consulta.estado)}</td>
                          <td>{formatDate(consulta.created_at)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AdminDashboard;