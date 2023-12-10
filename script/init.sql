CREATE DATABASE IF NOT EXISTS `yaranai`;
USE `yaranai`;
CREATE USER yaranai IDENTIFIED BY 'password';
CREATE TABLE IF NOT EXISTS `task` (
  `user` text NOT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` text NOT NULL,
  `description` text NOT NULL,
  `possibility` int(11) NOT NULL,
  `difficulty` int(11) NOT NULL,
  `created_at` date NOT NULL,
  `updated_at` date NOT NULL,
  `due_date` date NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
INSERT INTO `task` (`user`,`id`,`title`,`description`,`possibility`,`difficulty`,`created_at`,`updated_at`,`due_date`) VALUES ('ramdos',1,'電磁気学の課題','第二回の講義までにやる',2,3,'2023-12-01','2023-12-05','2023-12-10');