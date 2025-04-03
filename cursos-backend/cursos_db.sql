-- Crear base de datos si no existe
CREATE DATABASE IF NOT EXISTS cursos_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE cursos_db;

-- Crear tabla de usuarios
CREATE TABLE IF NOT EXISTS usuarios (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(100) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Crear tabla de cursos
CREATE TABLE IF NOT EXISTS cursos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  titulo VARCHAR(200) NOT NULL,
  descripcion VARCHAR(500),
  contenido TEXT,
  precio DECIMAL(10,2) DEFAULT 29.99,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Crear tabla de pagos
CREATE TABLE IF NOT EXISTS pagos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  usuario_id INT NOT NULL,
  curso_id INT NOT NULL,
  monto DECIMAL(10,2) NOT NULL,
  metodo VARCHAR(50) NOT NULL,
  estado VARCHAR(20) NOT NULL DEFAULT 'pendiente',
  transaccion_id VARCHAR(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
  FOREIGN KEY (curso_id) REFERENCES cursos(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Insertar algunos cursos de muestra
INSERT INTO cursos (titulo, descripcion, contenido, precio) VALUES 
('Curso de React', 'Aprende React desde cero.', 'En este curso aprenderás los fundamentos de React, incluyendo componentes, estado, props, hooks, y más. Crearás aplicaciones web modernas y reactivas utilizando las mejores prácticas de la industria.', 29.99),
('Curso de Node.js', 'Domina el backend con Node.js.', 'Aprende a construir APIs RESTful con Node.js y Express. Este curso cubre desde lo básico hasta conceptos avanzados como autenticación, manejo de errores, y despliegue.', 39.99),
('Curso de Diseño Web', 'Crea diseños modernos y responsivos.', 'Domina HTML, CSS y JavaScript para crear sitios web atractivos y funcionales. Aprenderás sobre diseño responsivo, animaciones, y las últimas tendencias en diseño web.', 24.99);