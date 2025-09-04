-- Script de inicialización de base de datos para microservicios
-- Este script crea las bases de datos separadas para cada microservicio

-- Crear bases de datos para cada microservicio
CREATE DATABASE IF NOT EXISTS auth_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS payment_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS course_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS contact_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS portfolio_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS home_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS analytics_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Crear usuario para microservicios si no existe
CREATE USER IF NOT EXISTS 'microuser'@'%' IDENTIFIED BY 'micropass';

-- Otorgar permisos a cada base de datos
GRANT ALL PRIVILEGES ON auth_db.* TO 'microuser'@'%';
GRANT ALL PRIVILEGES ON payment_db.* TO 'microuser'@'%';
GRANT ALL PRIVILEGES ON course_db.* TO 'microuser'@'%';
GRANT ALL PRIVILEGES ON contact_db.* TO 'microuser'@'%';
GRANT ALL PRIVILEGES ON portfolio_db.* TO 'microuser'@'%';
GRANT ALL PRIVILEGES ON home_db.* TO 'microuser'@'%';
GRANT ALL PRIVILEGES ON analytics_db.* TO 'microuser'@'%';

-- Aplicar cambios
FLUSH PRIVILEGES;