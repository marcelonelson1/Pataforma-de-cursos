/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19-11.4.5-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: cursos_db
-- ------------------------------------------------------
-- Server version	11.4.5-MariaDB-1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*M!100616 SET @OLD_NOTE_VERBOSITY=@@NOTE_VERBOSITY, NOTE_VERBOSITY=0 */;

--
-- Table structure for table `activity_logs`
--

DROP TABLE IF EXISTS `activity_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `activity_logs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `action` varchar(50) NOT NULL,
  `details` varchar(500) DEFAULT NULL,
  `ip` varchar(50) DEFAULT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=210 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `activity_logs`
--

LOCK TABLES `activity_logs` WRITE;
/*!40000 ALTER TABLE `activity_logs` DISABLE KEYS */;
INSERT INTO `activity_logs` VALUES
(1,3,'create_portfolio','Creado proyecto \'hola\' en categoría \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 15:16:51.616'),
(2,3,'create_portfolio','Creado proyecto \'tiki\' en categoría \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 15:18:32.551'),
(3,3,'create_portfolio','Creado proyecto \'hola\' en categoría \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 15:37:21.116'),
(4,3,'delete_portfolio','Eliminado proyecto ID: 1, Título: \'hola\', Categoría: \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 15:40:01.826'),
(5,3,'create_portfolio','Creado proyecto \'holaMMM\' en categoría \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 15:40:24.666'),
(6,3,'update_portfolio','Actualizado proyecto ID: 3, Título: \'hola\' -> \'hola\', Categoría: \'architecture\' -> \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 16:13:45.921'),
(7,3,'create_portfolio','Creado proyecto \'hola\' en categoría \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 16:20:25.127'),
(8,3,'update_portfolio','Actualizado proyecto ID: 5, Título: \'hola\' -> \'hola\', Categoría: \'architecture\' -> \'interiors\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 16:20:57.319'),
(9,3,'update_portfolio','Actualizado proyecto ID: 5, Título: \'hola\' -> \'hola\', Categoría: \'interiors\' -> \'interiors\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 16:22:59.372'),
(10,3,'update_portfolio','Actualizado proyecto ID: 5, Título: \'hola\' -> \'hola\', Categoría: \'interiors\' -> \'interiors\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 16:26:13.169'),
(11,3,'update_portfolio','Actualizado proyecto ID: 5, Título: \'hola\' -> \'hola\', Categoría: \'interiors\' -> \'interiors\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 16:27:55.372'),
(12,3,'update_portfolio','Actualizado proyecto ID: 5, Título: \'hola\' -> \'hola\', Categoría: \'interiors\' -> \'interiors\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 16:29:54.095'),
(13,3,'update_portfolio','Actualizado proyecto ID: 5, Título: \'hola\' -> \'hola\', Categoría: \'interiors\' -> \'interiors\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 16:30:13.245'),
(14,3,'delete_portfolio','Eliminado proyecto ID: 5, Título: \'hola\', Categoría: \'interiors\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 22:32:45.404'),
(15,3,'delete_portfolio','Eliminado proyecto ID: 4, Título: \'holaMMM\', Categoría: \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-08 22:32:53.371'),
(16,3,'create_portfolio','Creado proyecto \'render escaleras\' en categoría \'commercial\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-09 19:05:13.294'),
(17,3,'delete_portfolio','Eliminado proyecto ID: 6, Título: \'render escaleras\', Categoría: \'commercial\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-10 17:38:01.239'),
(18,3,'delete_portfolio','Eliminado proyecto ID: 2, Título: \'tiki\', Categoría: \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-10 17:38:05.548'),
(19,3,'delete_portfolio','Eliminado proyecto ID: 3, Título: \'hola\', Categoría: \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-11 17:50:31.264'),
(20,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:18:49.943'),
(21,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:18:49.997'),
(22,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:18:50.008'),
(23,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:18:51.636'),
(24,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:18:52.475'),
(25,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:18:52.624'),
(26,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:05.158'),
(27,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:05.401'),
(28,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:05.423'),
(29,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:05.478'),
(30,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:05.689'),
(31,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:05.866'),
(32,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:42.851'),
(33,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:43.023'),
(34,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:43.038'),
(35,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:43.057'),
(36,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:43.271'),
(37,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:19:43.428'),
(38,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:20:46.148'),
(39,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:20:46.241'),
(40,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:20:49.511'),
(41,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:20:49.642'),
(42,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:22:35.399'),
(43,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:22:35.598'),
(44,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:22:35.751'),
(45,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:22:35.757'),
(46,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:22:35.963'),
(47,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:22:36.187'),
(48,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:30:33.301'),
(49,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:30:33.541'),
(50,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:30:33.575'),
(51,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:30:33.611'),
(52,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:30:33.784'),
(53,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:30:33.912'),
(54,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:37:12.171'),
(55,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:37:12.288'),
(56,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:37:12.335'),
(57,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:37:12.347'),
(58,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:37:12.480'),
(59,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:37:12.629'),
(60,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:40:37.057'),
(61,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:40:37.142'),
(62,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:40:37.160'),
(63,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:40:37.176'),
(64,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:40:37.442'),
(65,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:40:37.576'),
(66,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:48:25.499'),
(67,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:48:25.948'),
(68,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:48:25.957'),
(69,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:48:26.117'),
(70,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:49:20.059'),
(71,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:49:20.328'),
(72,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:49:20.337'),
(73,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-13 23:49:20.488'),
(74,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:24:56.262'),
(75,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:24:56.395'),
(76,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:27:41.406'),
(77,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:27:41.531'),
(78,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:33.208'),
(79,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:33.463'),
(80,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:33.485'),
(81,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:33.512'),
(82,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:33.615'),
(83,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:33.756'),
(84,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:52.743'),
(85,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:53.085'),
(86,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:53.126'),
(87,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:53.141'),
(88,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:53.351'),
(89,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:37:53.601'),
(90,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:02.141'),
(91,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:02.220'),
(92,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:12.462'),
(93,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:12.636'),
(94,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:13.990'),
(95,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:14.130'),
(96,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:15.947'),
(97,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:16.105'),
(98,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:17.380'),
(99,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:17.551'),
(100,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:23.751'),
(101,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:23.926'),
(102,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:28.451'),
(103,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:38:28.635'),
(104,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:39:06.017'),
(105,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:39:06.090'),
(106,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:39:06.948'),
(107,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 00:39:07.058'),
(108,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:15.654'),
(109,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:15.889'),
(110,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:15.932'),
(111,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:15.937'),
(112,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:16.092'),
(113,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:16.220'),
(114,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:36.675'),
(115,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:36.828'),
(116,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:36.856'),
(117,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:36.876'),
(118,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:36.987'),
(119,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 18:42:37.109'),
(120,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:13:37.485'),
(121,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:13:37.737'),
(122,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:13:37.801'),
(123,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:13:37.817'),
(124,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:13:37.954'),
(125,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:13:38.101'),
(126,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:13:40.726'),
(127,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:13:40.944'),
(128,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:20:32.680'),
(129,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:20:32.791'),
(130,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:20:32.795'),
(131,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:20:32.900'),
(132,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:18.316'),
(133,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:18.524'),
(134,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:18.660'),
(135,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:18.749'),
(136,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:29.192'),
(137,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:29.316'),
(138,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:29.596'),
(139,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:29.609'),
(140,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:29.689'),
(141,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:29.694'),
(142,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:45.273'),
(143,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:28:45.419'),
(144,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:30:13.443'),
(145,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:30:13.456'),
(146,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:30:13.466'),
(147,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:30:13.567'),
(148,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:45:51.407'),
(149,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:45:51.623'),
(150,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:45:51.650'),
(151,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:45:51.831'),
(152,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:49:39.042'),
(153,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:49:39.275'),
(154,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:49:39.280'),
(155,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-14 23:49:39.357'),
(156,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:02:42.254'),
(157,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:02:42.337'),
(158,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:02:42.349'),
(159,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:02:42.368'),
(160,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:03:19.999'),
(161,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:03:20.012'),
(162,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:03:24.262'),
(163,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:03:24.273'),
(164,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:03:36.871'),
(165,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:03:36.896'),
(166,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:03:54.984'),
(167,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:03:54.990'),
(168,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:03:55.035'),
(169,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:03:55.048'),
(170,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:04:30.873'),
(171,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:04:30.872'),
(172,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:04:30.899'),
(173,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:04:30.901'),
(174,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:16.715'),
(175,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:16.733'),
(176,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:16.728'),
(177,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:16.780'),
(178,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:16.780'),
(179,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:16.820'),
(180,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:16.874'),
(181,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:16.880'),
(182,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:29.887'),
(183,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:29.889'),
(184,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:31.965'),
(185,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:31.967'),
(186,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:33.032'),
(187,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:33.033'),
(188,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:34.143'),
(189,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:34.145'),
(190,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:57.256'),
(191,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:57.289'),
(192,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:57.296'),
(193,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:57.312'),
(194,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:57.385'),
(195,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:15:57.491'),
(196,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:23:01.988'),
(197,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:23:02.095'),
(198,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:23:02.118'),
(199,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:23:02.123'),
(200,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:23:02.260'),
(201,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:23:02.389'),
(202,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:23:09.465'),
(203,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:23:09.521'),
(204,3,'view_stats','Visualización de estadísticas del panel admin','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:23:10.670'),
(205,3,'view_sales_stats','Consulta de estadísticas de ventas','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-15 00:23:10.726'),
(206,3,'create_portfolio','Creado proyecto \'florrr\' en categoría \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-17 21:09:01.663'),
(207,3,'delete_portfolio','Eliminado proyecto ID: 7, Título: \'florrr\', Categoría: \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-18 18:42:24.155'),
(208,3,'create_portfolio','Creado proyecto \'redner 1\' en categoría \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-18 18:42:43.519'),
(209,3,'create_portfolio','Creado proyecto \'2\' en categoría \'architecture\'','127.0.0.1','Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0','2025-04-18 18:43:07.810');
/*!40000 ALTER TABLE `activity_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `capitulos`
--

DROP TABLE IF EXISTS `capitulos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `capitulos` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `curso_id` bigint(20) unsigned NOT NULL,
  `titulo` varchar(200) NOT NULL,
  `descripcion` varchar(500) DEFAULT NULL,
  `duracion` varchar(10) DEFAULT NULL,
  `video_url` varchar(255) DEFAULT NULL,
  `orden` bigint(20) DEFAULT 0,
  `publicado` tinyint(1) DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `video_nombre` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_capitulos_curso_id` (`curso_id`),
  CONSTRAINT `fk_capitulos_curso` FOREIGN KEY (`curso_id`) REFERENCES `cursos` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_cursos_capitulos` FOREIGN KEY (`curso_id`) REFERENCES `cursos` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `capitulos`
--

LOCK TABLES `capitulos` WRITE;
/*!40000 ALTER TABLE `capitulos` DISABLE KEYS */;
INSERT INTO `capitulos` VALUES
(33,60,'test2','cap1','06:20','/static/videos/60/60-nuevo-2442d6e8-0c5f-482e-b7b4-7d916936c0e5.mp4',0,0,'2025-04-18 18:39:47.499','2025-04-18 18:39:47.499','60-nuevo-2442d6e8-0c5f-482e-b7b4-7d916936c0e5.mp4'),
(34,60,'aaaa','aaaaaaaa','06:14','/static/videos/60/60-nuevo-b9cd5155-220c-4f45-babf-8152907b6841.mp4',1,0,'2025-04-18 18:40:02.018','2025-04-18 18:40:02.018','60-nuevo-b9cd5155-220c-4f45-babf-8152907b6841.mp4'),
(35,60,'aerer','sefsefs','06:20','/static/videos/60/60-nuevo-01780c41-190f-4123-a0e9-58e55b517006.mp4',2,0,'2025-04-18 18:40:13.157','2025-04-18 18:40:13.157','60-nuevo-01780c41-190f-4123-a0e9-58e55b517006.mp4');
/*!40000 ALTER TABLE `capitulos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `contact_messages`
--

DROP TABLE IF EXISTS `contact_messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `contact_messages` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `message` text NOT NULL,
  `read` tinyint(1) DEFAULT 0,
  `starred` tinyint(1) DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_contact_messages_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `contact_messages`
--

LOCK TABLES `contact_messages` WRITE;
/*!40000 ALTER TABLE `contact_messages` DISABLE KEYS */;
INSERT INTO `contact_messages` VALUES
(1,'2025-04-07 18:24:41.871','2025-04-07 18:24:41.871','2025-04-07 18:32:23.933','marcelo nelson','marcelinho.nelson@gmail.com','','HOLA\n',0,0),
(2,'2025-04-07 18:32:40.962','2025-04-07 18:34:56.574','2025-04-07 19:50:37.363','hola','marcelinho.nelson@gmail.com','','hola\n',1,0),
(3,'2025-04-07 18:59:56.863','2025-04-07 19:04:47.293','2025-04-09 16:44:42.419','marcelo','marcelinho.nelson@gmail.com','','hola',1,0),
(4,'2025-04-07 19:27:38.979','2025-04-07 19:48:41.102',NULL,'marcelinho.nelson@gmail.com','marcelinho.nelson@gmail.com','hola','hola',1,0),
(5,'2025-04-08 22:26:43.471','2025-04-08 22:27:03.075',NULL,'flor mi vida','marcelinho.nelson@gmail.com','','flor te amo',1,0),
(6,'2025-04-09 16:44:29.451','2025-04-09 16:44:29.451',NULL,'hola','marcelinho.nelson@gmail.com','','aaaaaaaaaaaaaa',0,0),
(7,'2025-04-09 19:03:08.730','2025-04-09 19:03:08.730',NULL,'hola','marcelinho.nelson@gmail.com','','quiero un render de interior\n',0,0),
(8,'2025-04-09 19:03:35.669','2025-04-09 19:04:15.862',NULL,'marcelinho.nelson@gmail.com','rm.renders04@gmail.com','','quiero un render de interior',1,0),
(9,'2025-04-15 02:08:05.187','2025-04-15 02:08:05.187',NULL,'marcelo nelson','lavacasaturnosaturtinta@gmail.com','','hola\n',0,0),
(10,'2025-04-17 02:52:41.359','2025-04-17 02:53:26.984','2025-04-17 21:16:52.252','aaaa','marcelinho.nelson@gmail.com','','aaaaaa',1,0),
(11,'2025-04-17 21:06:53.115','2025-04-18 16:56:38.316',NULL,'marcelo nelson','marcelinho.nelson@gmail.com','','holaaaaaa',1,0),
(12,'2025-04-17 21:08:08.813','2025-04-17 21:16:24.141','2025-04-17 21:16:48.249','aaaaa','marcelinho.nelson@gmail.com','','holaaaaa',1,0),
(13,'2025-04-18 20:40:11.486','2025-04-18 20:40:11.486',NULL,'aaa','marcelinho.nelson@gmail.com','','a',0,0),
(14,'2025-04-19 06:02:01.516','2025-04-19 06:02:01.516',NULL,'aaaaaaa','marcelinho.nelson@gmail.com','','hola',0,0);
/*!40000 ALTER TABLE `contact_messages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cursos`
--

DROP TABLE IF EXISTS `cursos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `cursos` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `titulo` varchar(200) NOT NULL,
  `descripcion` varchar(500) DEFAULT NULL,
  `contenido` text DEFAULT NULL,
  `precio` decimal(10,2) DEFAULT 29.99,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `estado` varchar(20) DEFAULT 'Borrador',
  `imagen_url` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=62 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cursos`
--

LOCK TABLES `cursos` WRITE;
/*!40000 ALTER TABLE `cursos` DISABLE KEYS */;
INSERT INTO `cursos` VALUES
(60,'aaaaaaa','aaaaa','aaaaa',12.00,'2025-04-18 18:31:19.269','2025-04-18 18:40:55.943','Publicado','/static/images/96d09bda-0b58-4040-ab3a-98401149f067.jpeg'),
(61,'a','a','a',12.00,'2025-04-19 06:02:37.605','2025-04-19 06:02:37.605','Publicado','');
/*!40000 ALTER TABLE `cursos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `home_images`
--

DROP TABLE IF EXISTS `home_images`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `home_images` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `image_url` varchar(255) DEFAULT NULL,
  `title` varchar(100) DEFAULT NULL,
  `subtitle` varchar(200) DEFAULT NULL,
  `order` bigint(20) DEFAULT 0,
  `is_active` tinyint(1) DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `image_order` bigint(20) DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `home_images`
--

LOCK TABLES `home_images` WRITE;
/*!40000 ALTER TABLE `home_images` DISABLE KEYS */;
INSERT INTO `home_images` VALUES
(3,'/static/home/1744400331962961628.jpeg','Nueva Imagen','Descripción de la imagen',0,1,'2025-04-11 16:38:51.969','2025-04-11 16:39:33.277',2),
(5,'/static/home/1744400397592614611.jpeg','Nueva Imagen','Descripción de la imagen',0,1,'2025-04-11 16:39:57.612','2025-04-11 16:39:57.612',3);
/*!40000 ALTER TABLE `home_images` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pagos`
--

DROP TABLE IF EXISTS `pagos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `pagos` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `usuario_id` bigint(20) unsigned NOT NULL,
  `curso_id` bigint(20) unsigned NOT NULL,
  `monto` decimal(10,2) NOT NULL,
  `metodo` varchar(50) NOT NULL,
  `estado` varchar(20) NOT NULL DEFAULT 'pendiente',
  `transaccion_id` varchar(100) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `moneda` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_pagos_usuario_id` (`usuario_id`),
  KEY `idx_pagos_curso_id` (`curso_id`),
  CONSTRAINT `fk_pagos_curso` FOREIGN KEY (`curso_id`) REFERENCES `cursos` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_pagos_curso_1743551677` FOREIGN KEY (`curso_id`) REFERENCES `cursos` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_pagos_curso_id` FOREIGN KEY (`curso_id`) REFERENCES `cursos` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_pagos_usuario` FOREIGN KEY (`usuario_id`) REFERENCES `usuarios` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_pagos_usuario_1743551677` FOREIGN KEY (`usuario_id`) REFERENCES `usuarios` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_pagos_usuario_id` FOREIGN KEY (`usuario_id`) REFERENCES `usuarios` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=133 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pagos`
--

LOCK TABLES `pagos` WRITE;
/*!40000 ALTER TABLE `pagos` DISABLE KEYS */;
INSERT INTO `pagos` VALUES
(124,3,60,12.00,'paypal','pendiente','08752803R2827174U','2025-04-18 18:31:27.832','2025-04-18 18:31:33.981','USD'),
(125,3,60,12.00,'paypal','pendiente','40M535722F826993S','2025-04-18 18:36:40.147','2025-04-18 18:36:40.607','USD'),
(126,3,60,12.00,'paypal','pendiente','8EB96675R54714222','2025-04-18 18:37:14.801','2025-04-18 18:37:15.253','USD'),
(127,3,60,12.00,'paypal','pendiente','9LB11101VV666981P','2025-04-18 18:38:16.360','2025-04-18 18:38:17.795','USD'),
(128,3,60,12.00,'paypal','pendiente','59S26711Y4970725T','2025-04-18 18:38:50.900','2025-04-18 18:38:52.241','USD'),
(129,3,60,12.00,'paypal','pendiente','39M00571H5683584K','2025-04-18 18:44:33.988','2025-04-18 18:44:43.353','USD'),
(130,3,60,12.00,'paypal','pendiente','4SK34861JH636980N','2025-04-18 18:45:52.270','2025-04-18 18:46:01.688','USD'),
(131,3,60,12.00,'paypal','aprobado','2LT84893W85695433','2025-04-18 18:46:47.055','2025-04-18 18:46:57.200','USD'),
(132,3,61,12.00,'paypal','aprobado','15V817358W7869009','2025-04-19 06:02:45.751','2025-04-19 06:03:09.274','USD');
/*!40000 ALTER TABLE `pagos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `password_resets`
--

DROP TABLE IF EXISTS `password_resets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `password_resets` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(100) NOT NULL,
  `token` varchar(100) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `used` tinyint(1) DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_password_resets_token` (`token`),
  KEY `idx_password_resets_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=69 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `password_resets`
--

LOCK TABLES `password_resets` WRITE;
/*!40000 ALTER TABLE `password_resets` DISABLE KEYS */;
INSERT INTO `password_resets` VALUES
(68,'marcelinho.nelson@gmail.com','b812c9c2-ad97-436b-a754-11caf464cd9d','2025-04-17 21:06:17.745','2025-04-17 22:06:17.745',0);
/*!40000 ALTER TABLE `password_resets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `progreso_capitulos`
--

DROP TABLE IF EXISTS `progreso_capitulos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `progreso_capitulos` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `usuario_id` bigint(20) unsigned NOT NULL,
  `curso_id` bigint(20) unsigned NOT NULL,
  `capitulo_id` bigint(20) unsigned NOT NULL,
  `completado` tinyint(1) DEFAULT 0,
  `progreso` decimal(5,2) DEFAULT 0.00,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_usuario_capitulo` (`usuario_id`,`curso_id`,`capitulo_id`),
  KEY `fk_progreso_capitulo_curso` (`curso_id`),
  KEY `fk_progreso_capitulo_capitulo` (`capitulo_id`),
  CONSTRAINT `fk_progreso_capitulo_capitulo` FOREIGN KEY (`capitulo_id`) REFERENCES `capitulos` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_progreso_capitulo_curso` FOREIGN KEY (`curso_id`) REFERENCES `cursos` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_progreso_capitulo_usuario` FOREIGN KEY (`usuario_id`) REFERENCES `usuarios` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `progreso_capitulos`
--

LOCK TABLES `progreso_capitulos` WRITE;
/*!40000 ALTER TABLE `progreso_capitulos` DISABLE KEYS */;
INSERT INTO `progreso_capitulos` VALUES
(9,3,60,35,1,100.00,'2025-04-18 18:47:07.691','2025-04-18 18:47:07.691');
/*!40000 ALTER TABLE `progreso_capitulos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `progreso_usuarios`
--

DROP TABLE IF EXISTS `progreso_usuarios`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `progreso_usuarios` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `usuario_id` bigint(20) unsigned NOT NULL,
  `curso_id` bigint(20) unsigned NOT NULL,
  `porcentaje_total` decimal(5,2) DEFAULT 0.00,
  `ultimo_capitulo` bigint(20) unsigned DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_usuario_curso` (`usuario_id`,`curso_id`),
  KEY `fk_progreso_curso` (`curso_id`),
  CONSTRAINT `fk_progreso_curso` FOREIGN KEY (`curso_id`) REFERENCES `cursos` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_progreso_usuario` FOREIGN KEY (`usuario_id`) REFERENCES `usuarios` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `progreso_usuarios`
--

LOCK TABLES `progreso_usuarios` WRITE;
/*!40000 ALTER TABLE `progreso_usuarios` DISABLE KEYS */;
INSERT INTO `progreso_usuarios` VALUES
(7,3,60,33.33,34,'2025-04-18 18:47:04.545','2025-04-18 18:47:16.363');
/*!40000 ALTER TABLE `progreso_usuarios` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `project_portfolios`
--

DROP TABLE IF EXISTS `project_portfolios`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `project_portfolios` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(200) NOT NULL,
  `category` varchar(50) NOT NULL,
  `description` varchar(500) DEFAULT NULL,
  `image_url` varchar(255) DEFAULT NULL,
  `order_index` bigint(20) DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `project_portfolios`
--

LOCK TABLES `project_portfolios` WRITE;
/*!40000 ALTER TABLE `project_portfolios` DISABLE KEYS */;
INSERT INTO `project_portfolios` VALUES
(8,'redner 1','architecture','d5render','/static/portfolio/d2e03326-9ee8-48c6-a695-89a45d132d12.jpeg',1,'2025-04-18 18:42:43.517','2025-04-18 18:42:43.517'),
(9,'2','architecture','d5renders','/static/portfolio/43ea8a67-8461-462f-8b81-c0ad9304bd66.jpeg',2,'2025-04-18 18:43:07.808','2025-04-18 18:43:07.808');
/*!40000 ALTER TABLE `project_portfolios` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `usuarios`
--

DROP TABLE IF EXISTS `usuarios`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `usuarios` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `nombre` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(100) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `role` varchar(20) DEFAULT 'user',
  `phone` varchar(20) DEFAULT NULL,
  `image_url` varchar(255) DEFAULT NULL,
  `last_login` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_usuarios_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `usuarios`
--

LOCK TABLES `usuarios` WRITE;
/*!40000 ALTER TABLE `usuarios` DISABLE KEYS */;
INSERT INTO `usuarios` VALUES
(3,'marcelo nelson','marcelinho.nelson@gmail.com','$2a$10$GKdf3EABqKchbT9Mqwz8A.ViNLT0V9pUI5rk6l57kQeaUizcapJ5i','2025-03-31 15:28:52.256','2025-04-17 03:02:21.995','admin',NULL,'/static/profiles/user_.jpeg','2025-04-17 03:02:21.994'),
(38,'MARCELO','PILARlEO123@GMAIL.COM','$2a$10$w9PCxeODdSALmveqpu45m.GaUIwrXNJF7ZIgNb7eYUARassym7MNy','2025-04-11 16:49:12.795','2025-04-11 16:49:12.795','user','','','0000-00-00 00:00:00.000'),
(39,'aaaaaa','wdadawdaw@gmail.com','$2a$10$eRZu1mSUW3fF3qiK57VkROfGYqd0W3goCga0rFwCmlhtmrz9TUD.a','2025-04-13 23:30:13.553','2025-04-13 23:30:13.553','user','','','0000-00-00 00:00:00.000'),
(40,'marcelo','marce.nelson8@gmail.com','$2a$10$6DrxxklI5BxO7PXn8bGIYedeB3QpeO8rarNnUBgA1LG74G70sdl9a','2025-04-15 00:30:57.467','2025-04-15 00:31:07.315','user','','','2025-04-15 00:31:07.314'),
(41,'lala','marce.nelson4@gmail.com','$2a$10$bSQtgi2rb3815y./yXyT9.HfHFIibU2oKmMwPH2IIUE9GlsqHm8Pq','2025-04-15 00:32:52.796','2025-04-15 00:35:04.327','user','','','2025-04-15 00:35:04.327'),
(42,'aaaaa','bbbaaaaaa@gmail.com','$2a$10$zZNr85gCGQ4bRNS9d39lhelyJ7CwzlHwdmBor7PEzhVsDLiL1bB/2','2025-04-17 02:47:56.327','2025-04-17 02:47:56.327','user','','','0000-00-00 00:00:00.000'),
(43,'flor','flor@gmail.com','$2a$10$bIHLx6B1vnPUh5BtA8pk9ug.1iDRKRNu9jkY6moSjBRXRejZ3MT6W','2025-04-17 02:56:14.014','2025-04-17 02:56:24.060','user','','','2025-04-17 02:56:24.059');
/*!40000 ALTER TABLE `usuarios` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*M!100616 SET NOTE_VERBOSITY=@OLD_NOTE_VERBOSITY */;

-- Dump completed on 2025-04-19 17:20:53
