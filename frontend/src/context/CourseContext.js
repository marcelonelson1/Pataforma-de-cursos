import React, { createContext, useState, useContext, useEffect, useCallback } from 'react';
import { useAuth } from './AuthContext';

// Crear contexto para compartir datos de cursos en toda la aplicación
const CourseContext = createContext();

export const CourseProvider = ({ children }) => {
  const { isLoggedIn, user, isAdmin } = useAuth();
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [paymentStatus, setPaymentStatus] = useState({});
  
  // URL base de la API
  const apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:5000';

  // Función para obtener todos los cursos
  const fetchCourses = useCallback(async () => {
    setLoading(true);
    try {
      const response = await fetch(`${apiUrl}/api/cursos`, {
        headers: {
          'Accept': 'application/json'
        }
      });

      if (!response.ok) {
        throw new Error(`Error ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();
      
      // Transformar los datos de la API al formato requerido por los componentes
      const formattedCourses = data.map(curso => ({
        id: curso.id,
        titulo: curso.titulo,
        descripcion: curso.descripcion,
        contenido: curso.contenido,
        precio: curso.precio,
        // Inicializamos videos como un array vacío que se completará más tarde
        videos: [], 
        publishedDate: curso.created_at ? new Date(curso.created_at).toISOString().split('T')[0] : null,
        status: 'Publicado', // Asumimos que todos los cursos de la API están publicados
        students: Math.floor(Math.random() * 200), // Datos de ejemplo
        rating: (Math.random() * 2 + 3).toFixed(1), // Calificación aleatoria entre 3.0 y 5.0
      }));

      // En una implementación real, haríamos otra petición para obtener los videos de cada curso
      // Por ahora, utilizamos datos de ejemplo
      const coursesWithVideos = formattedCourses.map(course => {
        // Generamos videos de ejemplo basados en el ID del curso
        const mockVideos = generateMockVideos(course.id);
        return {
          ...course,
          videos: mockVideos,
          chapters: mockVideos.map(video => ({
            id: video.id,
            title: video.titulo,
            description: `Descripción para ${video.titulo}`,
            duration: video.duracion,
            video: video.url,
            isPublished: true
          }))
        };
      });

      setCourses(coursesWithVideos);
      setError(null);
    } catch (err) {
      console.error('Error al cargar cursos:', err);
      setError('No se pudieron cargar los cursos. Por favor, inténtalo de nuevo más tarde.');
      // Si falla la carga desde la API, usamos datos de demostración
      setCourses(getDemoCourses());
    } finally {
      setLoading(false);
    }
  }, [apiUrl]);

  // Función para generar videos de ejemplo basados en el ID del curso
  const generateMockVideos = (courseId) => {
    const videoCount = Math.floor(Math.random() * 3) + 2; // Entre 2 y 4 videos por curso
    const videoSources = [
      'https://www.learningcontainer.com/wp-content/uploads/2020/05/sample-mp4-file.mp4',
      'https://filesamples.com/samples/video/mp4/sample_640x360.mp4',
      'https://filesamples.com/samples/video/mp4/sample_960x540.mp4',
      'https://filesamples.com/samples/video/mp4/sample_1280x720.mp4'
    ];
    
    return Array.from({ length: videoCount }, (_, i) => ({
      id: parseInt(`${courseId}${i + 1}`),
      indice: `${i + 1}`,
      titulo: `Capítulo ${i + 1} del Curso ${courseId}`,
      url: videoSources[i % videoSources.length],
      duracion: `${Math.floor(Math.random() * 20) + 5}:${Math.floor(Math.random() * 59).toString().padStart(2, '0')}`,
      isPublished: true
    }));
  };

  // Datos de demostración para usar cuando la API no está disponible
  const getDemoCourses = () => [
    {
      id: 1,
      titulo: 'Curso de React',
      descripcion: 'Aprende React desde cero hasta un nivel avanzado con proyectos prácticos.',
      contenido: 'En este curso aprenderás todos los conceptos fundamentales de React, componentes, estados, props, hooks, y mucho más.',
      precio: 29.99,
      publishedDate: '2023-06-15',
      status: 'Publicado',
      students: 120,
      rating: 4.8,
      videos: [
        {
          id: 11,
          indice: "1",
          titulo: 'Introducción a React',
          url: 'https://www.learningcontainer.com/wp-content/uploads/2020/05/sample-mp4-file.mp4',
          duracion: '08:45'
        },
        {
          id: 12,
          indice: "2",
          titulo: 'Componentes y Props',
          url: 'https://filesamples.com/samples/video/mp4/sample_640x360.mp4',
          duracion: '10:15'
        },
        {
          id: 13,
          indice: "3",
          titulo: 'Estado y Ciclo de Vida',
          url: 'https://filesamples.com/samples/video/mp4/sample_960x540.mp4',
          duracion: '12:50'
        },
        {
          id: 14,
          indice: "4",
          titulo: 'Hooks en React',
          url: 'https://filesamples.com/samples/video/mp4/sample_1280x720.mp4',
          duracion: '15:30'
        }
      ],
      chapters: [
        {
          id: 11,
          title: 'Introducción a React',
          description: 'Introducción a los conceptos básicos de React',
          duration: '08:45',
          video: 'https://www.learningcontainer.com/wp-content/uploads/2020/05/sample-mp4-file.mp4',
          isPublished: true
        },
        {
          id: 12,
          title: 'Componentes y Props',
          description: 'Aprende a crear y usar componentes con props',
          duration: '10:15',
          video: 'https://filesamples.com/samples/video/mp4/sample_640x360.mp4',
          isPublished: true
        },
        {
          id: 13,
          title: 'Estado y Ciclo de Vida',
          description: 'Manejo del estado y ciclo de vida de componentes',
          duration: '12:50',
          video: 'https://filesamples.com/samples/video/mp4/sample_960x540.mp4',
          isPublished: true
        },
        {
          id: 14,
          title: 'Hooks en React',
          description: 'Uso de hooks para manejar estado y efectos',
          duration: '15:30',
          video: 'https://filesamples.com/samples/video/mp4/sample_1280x720.mp4',
          isPublished: true
        }
      ]
    },
    {
      id: 2,
      titulo: 'Desarrollo con Node.js',
      descripcion: 'Aprende a crear APIs RESTful y aplicaciones backend con Node.js, Express y MongoDB.',
      contenido: 'Este curso te enseñará a construir servicios web robustos y escalables utilizando Node.js y las tecnologías más demandadas del mercado.',
      precio: 39.99,
      publishedDate: '2023-05-10',
      status: 'Publicado',
      students: 85,
      rating: 4.5,
      videos: [
        {
          id: 21,
          indice: "1",
          titulo: 'Introducción a Node.js',
          url: 'https://filesamples.com/samples/video/mp4/sample_960x540.mp4',
          duracion: '12:30'
        },
        {
          id: 22,
          indice: "2",
          titulo: 'Express y Middleware',
          url: 'https://filesamples.com/samples/video/mp4/sample_640x360.mp4',
          duracion: '14:45'
        }
      ],
      chapters: [
        {
          id: 21,
          title: 'Introducción a Node.js',
          description: 'Fundamentos de Node.js y su ecosistema',
          duration: '12:30',
          video: 'https://filesamples.com/samples/video/mp4/sample_960x540.mp4',
          isPublished: true
        },
        {
          id: 22,
          title: 'Express y Middleware',
          description: 'Creación de APIs con Express y uso de middleware',
          duration: '14:45',
          video: 'https://filesamples.com/samples/video/mp4/sample_640x360.mp4',
          isPublished: true
        }
      ]
    }
  ];

  // Función para verificar el estado de pago de un curso
  const checkCoursePaymentStatus = useCallback(async (courseId) => {
    if (!isLoggedIn) {
      setPaymentStatus(prev => ({ ...prev, [courseId]: 'no_pagado' }));
      return 'no_pagado';
    }

    try {
      const token = localStorage.getItem('token');
      if (!token) {
        setPaymentStatus(prev => ({ ...prev, [courseId]: 'no_pagado' }));
        return 'no_pagado';
      }

      const response = await fetch(`${apiUrl}/api/pagos/${courseId}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Accept': 'application/json'
        }
      });

      // Si el servidor responde con un error, asumimos que no está pagado
      if (!response.ok) {
        setPaymentStatus(prev => ({ ...prev, [courseId]: 'no_pagado' }));
        return 'no_pagado';
      }

      const data = await response.json();
      const estado = data.estado || 'no_pagado';
      
      setPaymentStatus(prev => ({ ...prev, [courseId]: estado }));
      return estado;
    } catch (err) {
      console.error('Error al verificar estado de pago:', err);
      setPaymentStatus(prev => ({ ...prev, [courseId]: 'no_pagado' }));
      return 'no_pagado';
    }
  }, [apiUrl, isLoggedIn]);

  // Función para procesar el pago de un curso
  const processCoursePayment = useCallback(async (courseId, paymentMethod, cardDetails = null) => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No has iniciado sesión');
      }

      const course = courses.find(c => c.id === parseInt(courseId));
      if (!course) {
        throw new Error('Curso no encontrado');
      }

      const paymentData = {
        curso_id: courseId,
        monto: course.precio,
        metodo: paymentMethod
      };

      if (paymentMethod === 'tarjeta' && cardDetails) {
        paymentData.detalles_tarjeta = {
          numero: cardDetails.number,
          expiracion: cardDetails.expiry,
          cvv: cardDetails.cvv
        };
      }

      const response = await fetch(`${apiUrl}/api/pagos`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(paymentData)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || `Error ${response.status}`);
      }

      const data = await response.json();
      
      // Simulamos un proceso de verificación del pago
      // En un entorno real, esto se haría con polling al servidor o websockets
      return new Promise((resolve, reject) => {
        const checkStatus = async () => {
          const status = await checkCoursePaymentStatus(courseId);
          if (status === 'aprobado') {
            resolve({ success: true, message: 'Pago aprobado' });
          } else if (status === 'rechazado') {
            reject(new Error('El pago fue rechazado'));
          } else {
            // Si sigue pendiente, volvemos a verificar después de un tiempo
            setTimeout(checkStatus, 2000);
          }
        };
        
        // Iniciamos la verificación después de 2 segundos
        setTimeout(checkStatus, 2000);
      });
    } catch (err) {
      console.error('Error al procesar pago:', err);
      throw err;
    }
  }, [apiUrl, courses, checkCoursePaymentStatus]);

  // Función para crear un nuevo curso (solo para administradores)
  const createCourse = useCallback(async (courseData) => {
    try {
      setLoading(true);
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No has iniciado sesión');
      }
      
      // Verificar si el usuario es administrador
      if (!isAdmin()) {
        throw new Error('No tienes permisos para crear cursos');
      }

      // Preparar los datos para la API
      const apiCourseData = {
        titulo: courseData.title,
        descripcion: courseData.description,
        contenido: courseData.content || courseData.description,
        precio: parseFloat(courseData.price)
      };

      const response = await fetch(`${apiUrl}/api/cursos`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(apiCourseData)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || `Error ${response.status}`);
      }

      const newCourse = await response.json();
      
      // Convertir el formato de la API al formato de nuestra aplicación
      const formattedCourse = {
        id: newCourse.id,
        titulo: newCourse.titulo,
        descripcion: newCourse.descripcion,
        contenido: newCourse.contenido,
        precio: newCourse.precio,
        publishedDate: new Date().toISOString().split('T')[0],
        status: 'Publicado',
        students: 0,
        rating: 0,
        videos: [],
        chapters: courseData.chapters || []
      };

      // Actualizar el estado local
      setCourses(prevCourses => [...prevCourses, formattedCourse]);
      
      return formattedCourse;
    } catch (err) {
      console.error('Error al crear curso:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [apiUrl, isAdmin]);

  // Función para actualizar un curso existente (solo para administradores)
  const updateCourse = useCallback(async (courseId, courseData) => {
    try {
      setLoading(true);
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No has iniciado sesión');
      }
      
      // Verificar si el usuario es administrador
      if (!isAdmin()) {
        throw new Error('No tienes permisos para actualizar cursos');
      }

      // Preparar los datos para la API
      const apiCourseData = {
        titulo: courseData.title || courseData.titulo,
        descripcion: courseData.description || courseData.descripcion,
        contenido: courseData.content || courseData.contenido || courseData.description || courseData.descripcion,
        precio: parseFloat(courseData.price || courseData.precio)
      };

      const response = await fetch(`${apiUrl}/api/cursos/${courseId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(apiCourseData)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || `Error ${response.status}`);
      }

      const updatedCourse = await response.json();
      
      // Actualizar el estado local manteniendo los videos/chapters existentes
      setCourses(prevCourses => prevCourses.map(course => {
        if (course.id === parseInt(courseId)) {
          return {
            ...course,
            titulo: updatedCourse.titulo,
            descripcion: updatedCourse.descripcion,
            contenido: updatedCourse.contenido,
            precio: updatedCourse.precio,
            status: courseData.status || course.status,
            chapters: courseData.chapters || course.chapters,
            videos: courseData.videos || course.videos,
            // Si hay videos, actualizar los chapters, y viceversa
            ...(courseData.videos && !courseData.chapters ? { 
              chapters: courseData.videos.map(v => ({
                id: v.id,
                title: v.titulo,
                description: `Descripción para ${v.titulo}`,
                duration: v.duracion,
                video: v.url,
                isPublished: true
              }))
            } : {}),
            ...(courseData.chapters && !courseData.videos ? { 
              videos: courseData.chapters.map(c => ({
                id: c.id,
                indice: c.index || courseData.chapters.indexOf(c) + 1,
                titulo: c.title,
                url: c.video,
                duracion: c.duration,
                isPublished: c.isPublished
              }))
            } : {})
          };
        }
        return course;
      }));
      
      return updatedCourse;
    } catch (err) {
      console.error('Error al actualizar curso:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [apiUrl, isAdmin]);

  // Función para eliminar un curso (solo para administradores)
  const deleteCourse = useCallback(async (courseId) => {
    try {
      setLoading(true);
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No has iniciado sesión');
      }
      
      // Verificar si el usuario es administrador
      if (!isAdmin()) {
        throw new Error('No tienes permisos para eliminar cursos');
      }

      const response = await fetch(`${apiUrl}/api/cursos/${courseId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || `Error ${response.status}`);
      }

      // Actualizar el estado local
      setCourses(prevCourses => prevCourses.filter(course => course.id !== parseInt(courseId)));
      
      return { success: true };
    } catch (err) {
      console.error('Error al eliminar curso:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [apiUrl, isAdmin]);

  // Función para agregar un capítulo/video a un curso
  const addChapter = useCallback((courseId, chapterData) => {
    // Verificar permisos de administrador
    if (!isAdmin()) {
      throw new Error('No tienes permisos para esta acción');
    }
    
    setCourses(prevCourses => prevCourses.map(course => {
      if (course.id === parseInt(courseId)) {
        // Generar un nuevo ID para el capítulo
        const newChapterId = Math.max(0, ...course.chapters.map(c => c.id)) + 1;
        
        // Crear el nuevo capítulo
        const newChapter = {
          id: chapterData.id || newChapterId,
          title: chapterData.title,
          description: chapterData.description,
          duration: chapterData.duration,
          video: chapterData.video,
          isPublished: chapterData.isPublished
        };
        
        // Crear el nuevo video correspondiente
        const newVideo = {
          id: chapterData.id || newChapterId,
          indice: String(course.videos.length + 1),
          titulo: chapterData.title,
          url: chapterData.video,
          duracion: chapterData.duration
        };
        
        return {
          ...course,
          chapters: [...course.chapters, newChapter],
          videos: [...course.videos, newVideo]
        };
      }
      return course;
    }));
  }, [isAdmin]);

  // Función para actualizar un capítulo/video existente
  const updateChapter = useCallback((courseId, chapterId, chapterData) => {
    // Verificar permisos de administrador
    if (!isAdmin()) {
      throw new Error('No tienes permisos para esta acción');
    }
    
    setCourses(prevCourses => prevCourses.map(course => {
      if (course.id === parseInt(courseId)) {
        // Actualizar el capítulo
        const updatedChapters = course.chapters.map(chapter => {
          if (chapter.id === chapterId) {
            return {
              ...chapter,
              title: chapterData.title,
              description: chapterData.description,
              duration: chapterData.duration,
              video: chapterData.video || chapter.video,
              isPublished: chapterData.isPublished
            };
          }
          return chapter;
        });
        
        // Actualizar el video correspondiente
        const updatedVideos = course.videos.map(video => {
          if (video.id === chapterId) {
            return {
              ...video,
              titulo: chapterData.title,
              url: chapterData.video || video.url,
              duracion: chapterData.duration
            };
          }
          return video;
        });
        
        return {
          ...course,
          chapters: updatedChapters,
          videos: updatedVideos
        };
      }
      return course;
    }));
  }, [isAdmin]);

  // Función para eliminar un capítulo/video
  const deleteChapter = useCallback((courseId, chapterId) => {
    // Verificar permisos de administrador
    if (!isAdmin()) {
      throw new Error('No tienes permisos para esta acción');
    }
    
    setCourses(prevCourses => prevCourses.map(course => {
      if (course.id === parseInt(courseId)) {
        return {
          ...course,
          chapters: course.chapters.filter(chapter => chapter.id !== chapterId),
          videos: course.videos.filter(video => video.id !== chapterId)
        };
      }
      return course;
    }));
  }, [isAdmin]);

  // Cargar los cursos al iniciar
  useEffect(() => {
    fetchCourses();
  }, [fetchCourses]);
  
  // Verificar permisos de administrador para funciones administrativas
  const isAdminUser = isAdmin();

  const contextValue = {
    courses,
    loading,
    error,
    paymentStatus,
    fetchCourses,
    checkCoursePaymentStatus,
    processCoursePayment,
    // Exponer todas las funciones, pero las de administrador solo funcionarán si el usuario es admin
    createCourse,
    updateCourse,
    deleteCourse,
    addChapter,
    updateChapter,
    deleteChapter,
    isAdminUser
  };

  return (
    <CourseContext.Provider value={contextValue}>
      {children}
    </CourseContext.Provider>
  );
};

// Hook para acceder al contexto
export const useCourses = () => {
  const context = useContext(CourseContext);
  if (!context) {
    throw new Error('useCourses debe ser usado dentro de un CourseProvider');
  }
  return context;
};