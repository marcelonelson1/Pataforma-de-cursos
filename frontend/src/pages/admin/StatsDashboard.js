import React, { useState, useEffect } from 'react';
import { 
  BarChart, 
  Bar, 
  XAxis, 
  YAxis, 
  Tooltip, 
  CartesianGrid,
  Legend,
  PieChart, 
  Pie, 
  Cell, 
  LineChart,
  Line,
  ResponsiveContainer,
  AreaChart,
  Area
} from 'recharts';
import { FaChartLine, FaChartBar, FaCalendarAlt, FaChartPie, FaArrowUp, FaArrowDown } from 'react-icons/fa';
import './StatsDashboard.css';

const StatsDashboard = () => {
  const [isLoading, setIsLoading] = useState(true);
  const [timeFilter, setTimeFilter] = useState('month');
  const [chartType, setChartType] = useState('bar');

  // Datos de estadísticas generales
  const statsData = [
    { 
      id: 'estudiantes', 
      title: 'Estudiantes Activos', 
      value: 1250, 
      prevValue: 1120, 
      change: '+11.6%', 
      trend: 'up',
      icon: "👥"
    },
    { 
      id: 'cursos', 
      title: 'Cursos Publicados', 
      value: 24, 
      prevValue: 20, 
      change: '+20%', 
      trend: 'up',
      icon: "📚"
    },
    { 
      id: 'ingresos', 
      title: 'Ingresos Mensuales', 
      value: '$6,540', 
      prevValue: '$5,890', 
      change: '+11.0%', 
      trend: 'up',
      icon: "💰"
    },
    { 
      id: 'valoracion', 
      title: 'Valoración Media', 
      value: '4.8', 
      prevValue: '4.7', 
      change: '+2.1%', 
      trend: 'up',
      icon: "⭐"
    }
  ];

  // Datos de ventas por curso
  const salesData = [
    { name: 'React Avanzado', ventas: 150, ingresos: 7485, porcentaje: 35 },
    { name: 'Node.js', ventas: 90, ingresos: 5394, porcentaje: 21 },
    { name: 'TypeScript', ventas: 75, ingresos: 4500, porcentaje: 17 },
    { name: 'JavaScript', ventas: 65, ingresos: 3250, porcentaje: 15 },
    { name: 'Vue.js', ventas: 55, ingresos: 2750, porcentaje: 12 }
  ];

  // Datos de ventas por mes
  const monthlyData = [
    { month: 'Ene', ventas: 45, usuarios: 580 },
    { month: 'Feb', ventas: 52, usuarios: 620 },
    { month: 'Mar', ventas: 48, usuarios: 710 },
    { month: 'Abr', ventas: 61, usuarios: 800 },
    { month: 'May', ventas: 55, usuarios: 880 },
    { month: 'Jun', ventas: 67, usuarios: 940 },
    { month: 'Jul', ventas: 71, usuarios: 1020 },
    { month: 'Ago', ventas: 85, usuarios: 1150 },
    { month: 'Sep', ventas: 91, usuarios: 1250 }
  ];

  // Datos para el gráfico de usuarios
  const userData = [
    { name: 'Nuevos', value: 650, color: '#3b82f6' },
    { name: 'Recurrentes', value: 410, color: '#10b981' },
    { name: 'Premium', value: 190, color: '#f97316' }
  ];

  // Datos para el gráfico de métodos de pago
  const paymentData = [
    { name: 'PayPal', value: 55, color: '#3b82f6' },
    { name: 'Tarjeta', value: 35, color: '#10b981' },
    { name: 'Transferencia', value: 10, color: '#f97316' }
  ];
  
  // Efecto para simular carga de datos
  useEffect(() => {
    const timer = setTimeout(() => {
      setIsLoading(false);
    }, 1000);
    return () => clearTimeout(timer);
  }, []);

  // CustomTooltip para gráficos
  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      return (
        <div className="custom-tooltip">
          <p className="label">{`${label}`}</p>
          {payload.map((entry, index) => (
            <p key={index} style={{ color: entry.color }}>
              {`${entry.name}: ${entry.value}`}
            </p>
          ))}
        </div>
      );
    }
    return null;
  };

  // Función para cambiar el período de tiempo
  const handleTimeFilterChange = (filter) => {
    setIsLoading(true);
    setTimeFilter(filter);
    setTimeout(() => {
      setIsLoading(false);
    }, 500);
  };

  // Función para cambiar el tipo de gráfico
  const handleChartTypeChange = (type) => {
    setChartType(type);
  };

  // Loading overlay
  const LoadingOverlay = () => (
    <div className="loading-overlay">
      <div className="spinner"></div>
    </div>
  );

  return (
    <div className="stats-dashboard">
      {isLoading && <LoadingOverlay />}

      <div className="section-header">
        <h2 className="section-title">Estadísticas</h2>
        <div className="time-filters">
          <button 
            className={`filter-btn ${timeFilter === 'week' ? 'active' : ''}`}
            onClick={() => handleTimeFilterChange('week')}
          >
            Semanal
          </button>
          <button 
            className={`filter-btn ${timeFilter === 'month' ? 'active' : ''}`}
            onClick={() => handleTimeFilterChange('month')}
          >
            Mensual
          </button>
          <button 
            className={`filter-btn ${timeFilter === 'year' ? 'active' : ''}`}
            onClick={() => handleTimeFilterChange('year')}
          >
            Anual
          </button>
        </div>
      </div>

      {/* Tarjetas de estadísticas */}
      <div className="stats-grid">
        {statsData.map((stat) => (
          <div key={stat.id} className="stat-card">
            <div className="stat-icon">{stat.icon}</div>
            <div className="stat-info">
              <span className="stat-title">{stat.title}</span>
              <span className="stat-value">{stat.value}</span>
              <div className={`stat-change ${stat.trend}`}>
                {stat.trend === 'up' ? <FaArrowUp /> : <FaArrowDown />}
                <span>{stat.change}</span>
                <span className="period">vs mes anterior</span>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Gráficos */}
      <div className="charts-container">
        {/* Gráfico de ventas por curso */}
        <div className="chart">
          <div className="chart-header">
            <h3 className="chart-title">Ventas por Curso</h3>
            <div className="chart-controls">
              <button 
                className={`chart-type-btn ${chartType === 'bar' ? 'active' : ''}`}
                onClick={() => handleChartTypeChange('bar')}
              >
                <FaChartBar />
              </button>
              <button 
                className={`chart-type-btn ${chartType === 'line' ? 'active' : ''}`}
                onClick={() => handleChartTypeChange('line')}
              >
                <FaChartLine />
              </button>
            </div>
          </div>
          
          <ResponsiveContainer width="100%" height={300}>
            {chartType === 'bar' ? (
              <BarChart data={salesData} margin={{ top: 20, right: 30, left: 20, bottom: 10 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.1)" />
                <XAxis dataKey="name" stroke="rgba(255,255,255,0.7)" />
                <YAxis stroke="rgba(255,255,255,0.7)" />
                <Tooltip content={<CustomTooltip />} />
                <Legend />
                <Bar dataKey="ventas" name="Ventas" fill="#3b82f6" radius={[4, 4, 0, 0]} />
              </BarChart>
            ) : (
              <LineChart data={salesData} margin={{ top: 20, right: 30, left: 20, bottom: 10 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.1)" />
                <XAxis dataKey="name" stroke="rgba(255,255,255,0.7)" />
                <YAxis stroke="rgba(255,255,255,0.7)" />
                <Tooltip content={<CustomTooltip />} />
                <Legend />
                <Line type="monotone" dataKey="ventas" name="Ventas" stroke="#3b82f6" strokeWidth={2} dot={{ r: 4 }} activeDot={{ r: 6 }} />
              </LineChart>
            )}
          </ResponsiveContainer>
        </div>
        
        {/* Gráfico de tendencia mensual */}
        <div className="chart">
          <div className="chart-header">
            <h3 className="chart-title">Tendencia Mensual</h3>
            <div className="chart-period">
              <FaCalendarAlt /> {timeFilter === 'week' ? 'Últimos 7 días' : timeFilter === 'month' ? 'Últimos 30 días' : 'Este año'}
            </div>
          </div>
          
          <ResponsiveContainer width="100%" height={300}>
            <AreaChart data={monthlyData} margin={{ top: 20, right: 30, left: 20, bottom: 10 }}>
              <defs>
                <linearGradient id="colorVentas" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.8}/>
                  <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorUsuarios" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#10b981" stopOpacity={0.8}/>
                  <stop offset="95%" stopColor="#10b981" stopOpacity={0}/>
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.1)" />
              <XAxis dataKey="month" stroke="rgba(255,255,255,0.7)" />
              <YAxis stroke="rgba(255,255,255,0.7)" />
              <Tooltip content={<CustomTooltip />} />
              <Legend />
              <Area 
                type="monotone" 
                dataKey="ventas" 
                name="Ventas" 
                stroke="#3b82f6" 
                fillOpacity={1} 
                fill="url(#colorVentas)" 
              />
              <Area 
                type="monotone" 
                dataKey="usuarios" 
                name="Usuarios" 
                stroke="#10b981" 
                fillOpacity={1} 
                fill="url(#colorUsuarios)" 
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
        
        {/* Gráfico de distribución de usuarios */}
        <div className="chart">
          <div className="chart-header">
            <h3 className="chart-title">Distribución de Usuarios</h3>
            <div className="chart-type">
              <FaChartPie />
            </div>
          </div>
          
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={userData}
                cx="50%"
                cy="50%"
                innerRadius={60}
                outerRadius={90}
                paddingAngle={5}
                dataKey="value"
                label={({name, percent}) => `${name}: ${(percent * 100).toFixed(0)}%`}
              >
                {userData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip formatter={(value) => [`${value} usuarios`, 'Cantidad']} />
              <Legend />
            </PieChart>
          </ResponsiveContainer>
        </div>
        
        {/* Gráfico de métodos de pago */}
        <div className="chart">
          <div className="chart-header">
            <h3 className="chart-title">Métodos de Pago</h3>
            <div className="chart-type">
              <FaChartPie />
            </div>
          </div>
          
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={paymentData}
                cx="50%"
                cy="50%"
                innerRadius={60}
                outerRadius={90}
                paddingAngle={5}
                dataKey="value"
                label={({name, percent}) => `${name}: ${(percent * 100).toFixed(0)}%`}
              >
                {paymentData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip formatter={(value) => [`${value}%`, 'Porcentaje']} />
              <Legend />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
};

export default StatsDashboard;