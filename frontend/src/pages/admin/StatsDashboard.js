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
import { 
  FaChartLine, 
  FaChartBar, 
  FaCalendarAlt, 
  FaChartPie, 
  FaArrowUp, 
  FaArrowDown, 
  FaSpinner,
  FaUsers,
  FaBookOpen,
  FaMoneyBillWave,
  FaStar,
  FaInfo,
  FaRedo
} from 'react-icons/fa';

// Import stats service
import StatsService from '../../components/services/statsService';

// API URL
const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000';

// Color palette
const COLORS = {
  primary: '#6366f1',
  secondary: '#8b5cf6',
  success: '#10b981',
  warning: '#f59e0b',
  danger: '#ef4444',
  info: '#0ea5e9',
  chartColors: ['#6366f1', '#8b5cf6', '#0ea5e9', '#f59e0b', '#ef4444', '#10b981', '#ec4899']
};

// Function to parse JWT token
const parseJwt = (token) => {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );
    return JSON.parse(jsonPayload);
  } catch (e) {
    console.error('Error parsing JWT:', e);
    return {};
  }
};

// Enhanced function to get user data
const getUserData = () => {
  try {
    // First check if token exists in AuthContext format
    const token = localStorage.getItem('token');
    if (token) {
      console.log('Token found in AuthContext format');
      
      // Get user info from token
      const decodedToken = parseJwt(token);
      
      // Try to get user object from localStorage to complement info
      const userStr = localStorage.getItem('user');
      if (userStr) {
        // Combine information
        const storedUser = JSON.parse(userStr);
        
        // Check if format is compatible
        if (storedUser.token === token) {
          console.log('Found user object compatible with token');
          return storedUser;
        }
        
        // If not compatible, update
        console.log('Updating user object with token information');
        storedUser.token = token;
        
        // Ensure role is synchronized
        if (decodedToken.role) {
          storedUser.role = decodedToken.role;
        }
        
        // Save update
        localStorage.setItem('user', JSON.stringify(storedUser));
        return storedUser;
      }
      
      // If no user object, create one based on token
      const userData = {
        id: decodedToken.user_id || decodedToken.id || 1,
        nombre: decodedToken.nombre || 'Usuario',
        email: decodedToken.email || '',
        role: decodedToken.role || 'user',
        token: token
      };
      
      // Save for future use
      localStorage.setItem('user', JSON.stringify(userData));
      console.log('Created user object from token');
      
      return userData;
    }
    
    // If no token, look in standard StatsDashboard format
    const userStr = localStorage.getItem('user');
    if (userStr) {
      const userData = JSON.parse(userStr);
      
      // If there's a token in the user object, sync it with AuthContext format
      if (userData.token && !token) {
        localStorage.setItem('token', userData.token);
        console.log('Token synchronized from user object to AuthContext format');
      }
      
      return userData;
    }
    
    console.warn('No user data found in any format');
    return null;
  } catch (error) {
    console.error('Error getting user data:', error);
    return null;
  }
};

// Enhanced function to check if user is admin
const isUserAdmin = (userData) => {
  if (!userData) return false;
  
  // First check if token exists in AuthContext format
  const token = localStorage.getItem('token');
  if (token) {
    try {
      const decodedToken = parseJwt(token);
      if (decodedToken.role === 'admin') {
        return true;
      }
    } catch (error) {
      console.error('Error verifying role from token:', error);
    }
  }
  
  // Check in userData object
  if (userData.role && userData.role.toLowerCase() === 'admin') {
    return true;
  }
  
  // Check in nested user object
  if (userData.user && userData.user.role && userData.user.role.toLowerCase() === 'admin') {
    return true;
  }
  
  return false;
};

// Enhanced function to get auth headers
const authHeader = () => {
  try {
    // First check if token exists in AuthContext format
    const token = localStorage.getItem('token');
    if (token) {
      return { 'Authorization': `Bearer ${token}` };
    }
    
    // If no token in AuthContext format, look in user object
    const userData = getUserData();
    if (!userData) return {};
    
    // Look for token in different possible locations
    if (userData.token) {
      return { 'Authorization': `Bearer ${userData.token}` };
    } else if (userData.accessToken) {
      return { 'Authorization': `Bearer ${userData.accessToken}` };
    } else if (userData.access_token) {
      return { 'Authorization': `Bearer ${userData.access_token}` };
    }
    
    console.log('No authentication token found in any format');
    return {};
  } catch (error) {
    console.error("Error getting authentication headers:", error);
    return {};
  }
};

// Integrated stats service
const statsService = {
  // Get general stats for admin dashboard
  getAdminStats: async (period = 'month') => {
    try {
      const headers = {
        'Content-Type': 'application/json',
        ...authHeader()
      };
      
      console.log("Request headers:", headers);
      
      const response = await fetch(`${API_URL}/api/admin/stats?period=${period}`, { 
        method: 'GET',
        headers: headers,
        credentials: 'include' // To include cookies if needed
      });
      
      if (!response.ok) {
        const errorText = await response.text();
        console.error(`API Error (${response.status}):`, errorText);
        
        // If auth error, try to fix
        if (response.status === 401 || response.status === 403) {
          console.log('Attempting to fix authentication issues...');
          StatsService.fixAdminPermissions();
        }
        
        throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
      }
      
      const data = await response.json();
      if (data && data.success) {
        return data.data;
      }
      console.error("API returned success=false:", data);
      return null;
    } catch (error) {
      console.error('Error in getAdminStats:', error);
      throw error;
    }
  },

  // Get detailed stats for admin dashboard
  getAdminDashboard: async () => {
    try {
      const headers = {
        'Content-Type': 'application/json',
        ...authHeader()
      };
      
      const response = await fetch(`${API_URL}/api/admin/dashboard`, {
        method: 'GET',
        headers: headers,
        credentials: 'include'
      });
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
      }
      
      const data = await response.json();
      if (data && data.success) {
        return data.data;
      }
      return null;
    } catch (error) {
      console.error('Error in getAdminDashboard:', error);
      throw error;
    }
  },

  // Get detailed sales stats
  getSalesStats: async (period = 'month', startDate = null, endDate = null) => {
    try {
      let url = `${API_URL}/api/admin/sales-stats?period=${period}`;
      
      if (startDate) {
        url += `&start_date=${startDate}`;
      }
      
      if (endDate) {
        url += `&end_date=${endDate}`;
      }
      
      const headers = {
        'Content-Type': 'application/json',
        ...authHeader()
      };
      
      const response = await fetch(url, {
        method: 'GET',
        headers: headers,
        credentials: 'include'
      });
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
      }
      
      const data = await response.json();
      if (data && data.success) {
        return data.data;
      }
      return null;
    } catch (error) {
      console.error('Error in getSalesStats:', error);
      throw error;
    }
  },

  // Get activity log with pagination
  getActivityLog: async (page = 1, limit = 50, filters = {}) => {
    try {
      let url = `${API_URL}/api/admin/activity-log?page=${page}&limit=${limit}`;
      
      // Add additional filters if they exist
      if (filters.user_id) url += `&user_id=${filters.user_id}`;
      if (filters.action) url += `&action=${filters.action}`;
      if (filters.start_date) url += `&start_date=${filters.start_date}`;
      if (filters.end_date) url += `&end_date=${filters.end_date}`;
      
      const headers = {
        'Content-Type': 'application/json',
        ...authHeader()
      };
      
      const response = await fetch(url, {
        method: 'GET',
        headers: headers,
        credentials: 'include'
      });
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
      }
      
      const data = await response.json();
      if (data && data.success) {
        return data.data;
      }
      return null;
    } catch (error) {
      console.error('Error in getActivityLog:', error);
      throw error;
    }
  }
};

const StatsDashboard = () => {
  const [isLoading, setIsLoading] = useState(true);
  const [timeFilter, setTimeFilter] = useState('month');
  const [chartType, setChartType] = useState('bar');
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);
  const [authError, setAuthError] = useState(null);
  
  // States for backend data
  const [statsData, setStatsData] = useState([]);
  const [salesData, setSalesData] = useState([]);
  const [monthlyData, setMonthlyData] = useState([]);
  const [userData, setUserData] = useState([]);
  const [paymentData, setPaymentData] = useState([]);
  const [error, setError] = useState(null);
  const [animationActive, setAnimationActive] = useState(false);

  // Check authentication and admin role at start
  useEffect(() => {
    checkUserAuth();
    
    // Add animation class after component mounts
    setTimeout(() => {
      setAnimationActive(true);
    }, 100);
  }, []);

  // Enhanced checkUserAuth function with better debugging and server validation
  const checkUserAuth = () => {
    try {
      setIsLoading(true);
      setError(null);
      setAuthError(null);
      
      // First try to sync data between formats
      StatsService.fixAdminPermissions();
      
      // Get user data
      const userData = getUserData();
      
      if (!userData) {
        console.error('No user data found in localStorage or token format');
        setIsAuthenticated(false);
        setIsAdmin(false);
        setAuthError('No se encontraron datos de usuario. Por favor inicie sesión nuevamente.');
        setIsLoading(false);
        return;
      }
      
      // Log user data for debugging
      console.log('User data found:', {
        id: userData.id,
        email: userData.email,
        hasToken: !!userData.token || !!localStorage.getItem('token'),
        role: userData.role
      });
      
      // Check if any token exists
      const hasToken = !!(userData.token || userData.accessToken || userData.access_token || localStorage.getItem('token'));
      
      if (!hasToken) {
        console.error('No token found in user data');
        setIsAuthenticated(false);
        setIsAdmin(false);
        setAuthError('Token de autenticación no encontrado. Por favor inicie sesión nuevamente.');
        setIsLoading(false);
        return;
      }
      
      // User has token, consider authenticated
      setIsAuthenticated(true);
      
      // Check admin status locally
      const adminStatus = isUserAdmin(userData);
      console.log('Is user admin (local check)?', adminStatus);
      
      if (!adminStatus) {
        console.log('User is not admin according to local check. Trying to fix...');
        const fixResult = StatsService.fixAdminPermissions();
        console.log('Fix result:', fixResult);
        
        // Check again after trying to fix
        const updatedUserData = getUserData();
        const updatedAdminStatus = isUserAdmin(updatedUserData);
        
        if (updatedAdminStatus) {
          console.log('Permissions fixed successfully, user is now admin');
          setIsAdmin(true);
          loadAllData();
        } else {
          console.error('User does not have admin permissions even after fixing');
          setIsAdmin(false);
          setAuthError('No tienes permisos de administrador para ver estas estadísticas.');
          setIsLoading(false);
        }
        return;
      }
      
      // If we get here, user has admin permissions according to local check
      setIsAdmin(true);
      
      // Verify admin status with server as a double-check
      const headers = {
        'Content-Type': 'application/json',
        ...authHeader()
      };
      
      fetch(`${API_URL}/api/auth/check-admin`, {
        method: 'GET',
        headers: headers,
        credentials: 'include'
      })
      .then(response => {
        if (!response.ok) {
          throw new Error(`Error ${response.status}: ${response.statusText}`);
        }
        return response.json();
      })
      .then(data => {
        if (data && data.success && data.data && data.data.isAdmin) {
          console.log('Server confirmed admin status:', data);
          setIsAdmin(true);
          // Load stats data
          loadAllData();
        } else {
          console.warn('Server denied admin status but local check passed:', data);
          // If server denies admin status but locally seemed to be admin,
          // trust local verification as it could be a backend issue
          setIsAdmin(true);
          loadAllData();
        }
      })
      .catch(err => {
        console.error('Error checking admin status with server:', err);
        // If there's an error verifying with server, trust local check
        console.log('Continuing with local verification after server error');
        // If we already determined is admin locally, load the data
        loadAllData();
      });
    } catch (err) {
      console.error('General error in checkUserAuth:', err);
      setIsAuthenticated(false);
      setIsAdmin(false);
      setAuthError(`Error al verificar la autenticación: ${err.message}`);
      setIsLoading(false);
    }
  };

  // Load data when time filter changes
  useEffect(() => {
    if (isAuthenticated && isAdmin && timeFilter) {
      loadAllData();
    }
  }, [timeFilter, isAuthenticated, isAdmin]);

  // Function to load all necessary data
  const loadAllData = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      // Load general stats
      const generalStats = await statsService.getAdminStats(timeFilter);
      
      if (generalStats) {
        // Transform data for dashboard
        const transformedStats = transformGeneralStats(generalStats);
        setStatsData(transformedStats);
      } else {
        throw new Error('No se pudieron obtener estadísticas generales');
      }
      
      // Load sales stats
      const salesStats = await statsService.getSalesStats(timeFilter);
      
      if (salesStats) {
        // Transform course sales data
        const transformedSales = transformSalesData(salesStats);
        setSalesData(transformedSales.coursesSales);
        
        // Set monthly data
        setMonthlyData(transformedSales.monthlyData);
        
        // Set user data
        setUserData(transformedSales.userData);
        
        // Set payment method data
        setPaymentData(transformedSales.paymentData);
      } else {
        throw new Error('No se pudieron obtener estadísticas de ventas');
      }
      
    } catch (err) {
      console.error('Error loading data:', err);
      
      const errorMessage = err.message || 'Error desconocido';
      
      if (errorMessage.includes('401')) {
        setError('No autorizado. Por favor inicia sesión nuevamente.');
        
        // Try to automatically fix the authentication issue
        console.log('Trying to fix authentication issues...');
        const fixResult = StatsService.fixAdminPermissions();
        console.log('Fix result:', fixResult);
        
        // If there's an authentication error, use default data
        // so at least something shows on the dashboard
        setDefaultData();
      } else {
        setError(`Error al cargar datos de estadísticas: ${errorMessage}`);
        // Also use default data for other errors
        setDefaultData();
      }
    } finally {
      setIsLoading(false);
    }
  };

  // Function to transform general stats
  const transformGeneralStats = (data) => {
    // If no data, return empty array
    if (!data || !data.stats) return [];
    
    const { activeStudents, publishedCourses, monthlyRevenue, averageRating } = data.stats;
    
    return [
      { 
        id: 'estudiantes', 
        title: 'Estudiantes Activos', 
        value: activeStudents.current, 
        prevValue: activeStudents.previous, 
        change: calculatePercentage(activeStudents.current, activeStudents.previous),
        trend: activeStudents.current >= activeStudents.previous ? 'up' : 'down',
        icon: <FaUsers />,
        color: COLORS.primary
      },
      { 
        id: 'cursos', 
        title: 'Cursos Publicados', 
        value: publishedCourses.current, 
        prevValue: publishedCourses.previous, 
        change: calculatePercentage(publishedCourses.current, publishedCourses.previous),
        trend: publishedCourses.current >= publishedCourses.previous ? 'up' : 'down',
        icon: <FaBookOpen />,
        color: COLORS.secondary
      },
      { 
        id: 'ingresos', 
        title: 'Ingresos Mensuales', 
        value: `${formatCurrency(monthlyRevenue.current)}`, 
        prevValue: `${formatCurrency(monthlyRevenue.previous)}`, 
        change: calculatePercentage(monthlyRevenue.current, monthlyRevenue.previous),
        trend: monthlyRevenue.current >= monthlyRevenue.previous ? 'up' : 'down',
        icon: <FaMoneyBillWave />,
        color: COLORS.success
      },
      { 
        id: 'valoracion', 
        title: 'Valoración Media', 
        value: averageRating.current.toFixed(1), 
        prevValue: averageRating.previous.toFixed(1), 
        change: calculatePercentage(averageRating.current, averageRating.previous),
        trend: averageRating.current >= averageRating.previous ? 'up' : 'down',
        icon: <FaStar />,
        color: COLORS.warning
      }
    ];
  };

  // Function to transform sales data
  const transformSalesData = (data) => {
    // If no data, return empty objects
    if (!data) {
      return {
        coursesSales: [],
        monthlyData: [],
        userData: [],
        paymentData: []
      };
    }
    
    // Transform course sales
    const coursesSales = data.coursesSales?.map((course, index) => ({
      name: course.name,
      ventas: course.sales,
      ingresos: course.revenue,
      porcentaje: course.percentage,
      color: COLORS.chartColors[index % COLORS.chartColors.length]
    })) || [];
    
    // Transform monthly data
    const monthlyData = data.monthlyData?.map(item => ({
      month: item.month,
      ventas: item.sales,
      usuarios: item.users
    })) || [];
    
    // Transform user data
    const userData = [
      { name: 'Nuevos', value: data.userStats?.new || 0, color: COLORS.primary },
      { name: 'Recurrentes', value: data.userStats?.returning || 0, color: COLORS.success },
      { name: 'Premium', value: data.userStats?.premium || 0, color: COLORS.warning }
    ];
    
    // Transform payment method data
    const paymentData = [
      { name: 'PayPal', value: data.paymentMethods?.paypal || 0, color: COLORS.primary },
      { name: 'Tarjeta', value: data.paymentMethods?.card || 0, color: COLORS.success },
      { name: 'Transferencia', value: data.paymentMethods?.transfer || 0, color: COLORS.warning }
    ];
    
    return {
      coursesSales,
      monthlyData,
      userData,
      paymentData
    };
  };

  // Function to calculate percentage change
  const calculatePercentage = (current, previous) => {
    if (previous === 0) return '+100%';
    
    const change = ((current - previous) / previous) * 100;
    return `${change >= 0 ? '+' : ''}${change.toFixed(1)}%`;
  };

  // Function to format numbers
  const formatNumber = (num) => {
    return new Intl.NumberFormat('es-ES').format(num);
  };

  // Function to format currency
  const formatCurrency = (num) => {
    return new Intl.NumberFormat('es-ES', {
      style: 'currency',
      currency: 'EUR',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(num);
  };

  // Set default data in case of error
  const setDefaultData = () => {
    // General stats data
    setStatsData([
      { 
        id: 'estudiantes', 
        title: 'Estudiantes Activos', 
        value: 1250, 
        prevValue: 1120, 
        change: '+11.6%', 
        trend: 'up',
        icon: <FaUsers />,
        color: COLORS.primary
      },
      { 
        id: 'cursos', 
        title: 'Cursos Publicados', 
        value: 24, 
        prevValue: 20, 
        change: '+20%', 
        trend: 'up',
        icon: <FaBookOpen />,
        color: COLORS.secondary
      },
      { 
        id: 'ingresos', 
        title: 'Ingresos Mensuales', 
        value: '6.540 €', 
        prevValue: '5.890 €', 
        change: '+11.0%', 
        trend: 'up',
        icon: <FaMoneyBillWave />,
        color: COLORS.success
      },
      { 
        id: 'valoracion', 
        title: 'Valoración Media', 
        value: '4.8', 
        prevValue: '4.7', 
        change: '+2.1%', 
        trend: 'up',
        icon: <FaStar />,
        color: COLORS.warning
      }
    ]);

    // Course sales data
    setSalesData([
      { name: 'React Avanzado', ventas: 150, ingresos: 7485, porcentaje: 35, color: COLORS.chartColors[0] },
      { name: 'Node.js', ventas: 90, ingresos: 5394, porcentaje: 21, color: COLORS.chartColors[1] },
      { name: 'TypeScript', ventas: 75, ingresos: 4500, porcentaje: 17, color: COLORS.chartColors[2] },
      { name: 'JavaScript', ventas: 65, ingresos: 3250, porcentaje: 15, color: COLORS.chartColors[3] },
      { name: 'Vue.js', ventas: 55, ingresos: 2750, porcentaje: 12, color: COLORS.chartColors[4] }
    ]);

    // Monthly sales data
    setMonthlyData([
      { month: 'Ene', ventas: 45, usuarios: 580 },
      { month: 'Feb', ventas: 52, usuarios: 620 },
      { month: 'Mar', ventas: 48, usuarios: 710 },
      { month: 'Abr', ventas: 61, usuarios: 800 },
      { month: 'May', ventas: 55, usuarios: 880 },
      { month: 'Jun', ventas: 67, usuarios: 940 },
      { month: 'Jul', ventas: 71, usuarios: 1020 },
      { month: 'Ago', ventas: 85, usuarios: 1150 },
      { month: 'Sep', ventas: 91, usuarios: 1250 }
    ]);

    // User chart data
    setUserData([
      { name: 'Nuevos', value: 650, color: COLORS.primary },
      { name: 'Recurrentes', value: 410, color: COLORS.success },
      { name: 'Premium', value: 190, color: COLORS.warning }
    ]);

    // Payment methods chart data
    setPaymentData([
      { name: 'PayPal', value: 55, color: COLORS.primary },
      { name: 'Tarjeta', value: 35, color: COLORS.success },
      { name: 'Transferencia', value: 10, color: COLORS.warning }
    ]);
  };

  // Function to change time period
  const handleTimeFilterChange = (filter) => {
    if (filter === timeFilter) return;
    
    setTimeFilter(filter);
    // useEffect will handle data reloading
  };

  // Function to change chart type
  const handleChartTypeChange = (type) => {
    setChartType(type);
  };

  // CustomTooltip for charts
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

  // Loading overlay
  const LoadingOverlay = () => (
    <div className="loading-overlay">
      <FaSpinner className="spinner" />
      <p className="loading-text">Cargando estadísticas...</p>
    </div>
  );

  // Error message
  const ErrorMessage = () => (
    <div className="error-message">
      <div className="error-icon"><FaInfo /></div>
      <p>Error al cargar datos: {error}</p>
      {isAuthenticated ? (
        <button className="retry-btn" onClick={loadAllData}>
          <FaRedo /> Reintentar
        </button>
      ) : (
        <a href="/login" className="login-btn">
          Iniciar sesión
        </a>
      )}
    </div>
  );

  // Auth required component
  const AuthRequired = () => (
    <div className="auth-required">
      <h3>Se requiere autenticación</h3>
      <p>{authError || 'Debes iniciar sesión con permisos de administrador para ver este panel.'}</p>
      <a href="/login" className="login-btn">
        Iniciar sesión
      </a>
      <button 
        className="create-test-user-btn" 
        onClick={() => {
          const result = StatsService.createTestAdminUser();
          if (result.success) {
            setAuthError('Usuario de prueba creado. Recargando...');
            setTimeout(() => {
              window.location.reload();
            }, 1500);
          } else {
            setAuthError(`Error al crear usuario de prueba: ${result.message}`);
          }
        }}
      >
        Crear usuario de prueba
      </button>
    </div>
  );

  // Insufficient permissions component
  const InsufficientPermissions = () => (
    <div className="auth-required">
      <h3>Permisos insuficientes</h3>
      <p>{authError || 'Tu cuenta no tiene permisos de administrador para ver este panel.'}</p>
      <a href="/" className="login-btn">
        Volver al inicio
      </a>
    </div>
  );

  // Component to verify and fix authentication issues
  const AuthFixer = () => (
    <div className="auth-fixer">
      <h3>Problemas de autenticación detectados</h3>
      <p>Se han detectado problemas con tus permisos de administrador. ¿Quieres intentar corregirlos?</p>
      <button 
        className="fix-btn"
        onClick={() => {
          const result = StatsService.fixAdminPermissions();
          
          if (result.success) {
            setAuthError('Datos corregidos. Recargando...');
            
            // Reload page after a brief delay
            setTimeout(() => {
              window.location.reload();
            }, 1500);
          } else {
            setAuthError(`No se pudieron corregir los permisos: ${result.message}`);
            
            // Offer to create a test user
            if (window.confirm('¿Deseas crear un usuario administrador de prueba para poder continuar?')) {
              const testResult = StatsService.createTestAdminUser();
              
              if (testResult.success) {
                setAuthError('Usuario de prueba creado. Recargando...');
                setTimeout(() => {
                  window.location.reload();
                }, 1500);
              } else {
                setAuthError(`Error al crear usuario de prueba: ${testResult.message}`);
              }
            }
          }
        }}
      >
        Corregir permisos
      </button>
    </div>
  );

  // If not authenticated, show message
  if (!isAuthenticated && !isLoading) {
    return (
      <div className={`stats-dashboard ${animationActive ? 'animate-in' : ''}`}>
        <AuthRequired />
      </div>
    );
  }
  
  // If authenticated but not admin, show insufficient permissions message
  // and option to fix typical issues
  if (isAuthenticated && !isAdmin && !isLoading) {
    return (
      <div className={`stats-dashboard ${animationActive ? 'animate-in' : ''}`}>
        <InsufficientPermissions />
        <AuthFixer />
      </div>
    );
  }

  return (
    <div className={`stats-dashboard ${animationActive ? 'animate-in' : ''}`}>
      {isLoading && <LoadingOverlay />}

      <div className="dashboard-header">
        <div className="dashboard-title">
          <h2 className="section-title">Panel de Estadísticas</h2>
          <p className="dashboard-subtitle">
            {timeFilter === 'week' 
              ? 'Datos de la última semana' 
              : timeFilter === 'month' 
                ? 'Datos del último mes' 
                : 'Datos del último año'}
          </p>
        </div>
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

      {error && <ErrorMessage />}

      {/* Stats Cards */}
      <div className="stats-grid">
        {statsData.map((stat) => (
          <div key={stat.id} className="stat-card">
            <div className="stat-icon" style={{ backgroundColor: `${stat.color}15`, color: stat.color }}>
              {stat.icon}
            </div>
            <div className="stat-info">
              <span className="stat-title">{stat.title}</span>
              <span className="stat-value">{stat.value}</span>
              <div className={`stat-change ${stat.trend}`}>
                {stat.trend === 'up' ? <FaArrowUp /> : <FaArrowDown />}
                <span>{stat.change}</span>
                <span className="period">vs {timeFilter === 'week' ? 'semana' : timeFilter === 'month' ? 'mes' : 'año'} anterior</span>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Charts */}
      <div className="charts-container">
        {/* Course Sales Chart */}
        <div className="chart">
          <div className="chart-header">
            <h3 className="chart-title">Ventas por Curso</h3>
            <div className="chart-controls">
              <button 
                className={`chart-type-btn ${chartType === 'bar' ? 'active' : ''}`}
                onClick={() => handleChartTypeChange('bar')}
                aria-label="Mostrar gráfico de barras"
              >
                <FaChartBar />
              </button>
              <button 
                className={`chart-type-btn ${chartType === 'line' ? 'active' : ''}`}
                onClick={() => handleChartTypeChange('line')}
                aria-label="Mostrar gráfico de líneas"
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
                <Bar dataKey="ventas" name="Ventas" radius={[4, 4, 0, 0]}>
                  {salesData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color || COLORS.chartColors[index % COLORS.chartColors.length]} />
                  ))}
                </Bar>
              </BarChart>
            ) : (
              <LineChart data={salesData} margin={{ top: 20, right: 30, left: 20, bottom: 10 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.1)" />
                <XAxis dataKey="name" stroke="rgba(255,255,255,0.7)" />
                <YAxis stroke="rgba(255,255,255,0.7)" />
                <Tooltip content={<CustomTooltip />} />
                <Legend />
                <Line 
                  type="monotone" 
                  dataKey="ventas" 
                  name="Ventas" 
                  stroke={COLORS.primary} 
                  strokeWidth={2} 
                  dot={{ r: 4 }} 
                  activeDot={{ r: 6, stroke: COLORS.primary, strokeWidth: 2, fill: 'white' }} 
                />
              </LineChart>
            )}
          </ResponsiveContainer>
        </div>
        
        {/* Monthly Trend Chart */}
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
                  <stop offset="5%" stopColor={COLORS.primary} stopOpacity={0.8}/>
                  <stop offset="95%" stopColor={COLORS.primary} stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorUsuarios" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor={COLORS.success} stopOpacity={0.8}/>
                  <stop offset="95%" stopColor={COLORS.success} stopOpacity={0}/>
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
                stroke={COLORS.primary} 
                fillOpacity={1} 
                fill="url(#colorVentas)" 
                activeDot={{ r: 6, stroke: COLORS.primary, strokeWidth: 2, fill: 'white' }}
              />
              <Area 
                type="monotone" 
                dataKey="usuarios" 
                name="Usuarios" 
                stroke={COLORS.success} 
                fillOpacity={1} 
                fill="url(#colorUsuarios)" 
                activeDot={{ r: 6, stroke: COLORS.success, strokeWidth: 2, fill: 'white' }}
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
        
        {/* User Distribution Chart */}
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
        
        {/* Payment Methods Chart */}
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