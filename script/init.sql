CREATE DATABASE IF NOT EXISTS `yaranai`;
USE `yaranai`;
CREATE USER yaranai IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON yaranai.* TO 'yaranai'@'%';
CREATE TABLE IF NOT EXISTS `task` (
  `user` text NOT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` text NOT NULL,
  `description` text NOT NULL,
  `condition_id` int(11) NOT NULL,
  `difficulty` int(11) NOT NULL,
  `created_at` date NOT NULL,
  `updated_at` date NOT NULL,
  `dueDate` date NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `condition` (
  `condition_id` int(11) NOT NULL AUTO_INCREMENT,
  `user` text NOT NULL,
  `condition` text NOT NULL,
  PRIMARY KEY (`condition_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `deleted_task` (
  `user` text NOT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `condition_id` int(11) NOT NULL,
  `created_at` date NOT NULL,
  `due_date` date NOT NULL,
  `deleted_at_unix` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
ALTER TABLE `task` ADD FOREIGN KEY (`condition_id`) REFERENCES `condition`(`condition_id`);

ALTER TABLE `deleted_task` ADD FOREIGN KEY (`condition_id`) REFERENCES `condition`(`condition_id`);
INSERT INTO `condition` (`condition_id`,`user`,`condition`) VALUES (1,'ramdos','anywhere');
-- INSERT INTO `task` (`user`,`id`,`title`,`description`,`condition_id`,`difficulty`,`created_at`,`updated_at`,`dueDate`) VALUES ('ramdos',1,'電磁気学の課題','第二回の講義までにやる',2,3,'2023-12-01','2023-12-05','2023-12-10');

