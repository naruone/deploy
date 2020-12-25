# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 172.100.200.61 (MySQL 5.7.19-log)
# Database: titan_customer_2
# Generation Time: 2020-12-25 09:26:04 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table deploy_tasks
# ------------------------------------------------------------

CREATE TABLE `deploy_tasks` (
  `task_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `task_name` varchar(255) NOT NULL DEFAULT '',
  `deploy_type` tinyint(4) NOT NULL,
  `description` varchar(500) NOT NULL DEFAULT '',
  `output` text NOT NULL,
  `env_id` int(11) NOT NULL,
  `branch` varchar(255) NOT NULL DEFAULT '',
  `version` varchar(255) NOT NULL DEFAULT '',
  `after_script` text NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '1',
  `uuid` varchar(255) NOT NULL DEFAULT '',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table env_pro_servers
# ------------------------------------------------------------

CREATE TABLE `env_pro_servers` (
  `env_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `env_name` varchar(255) NOT NULL DEFAULT '',
  `project_id` int(11) NOT NULL,
  `server_ids` varchar(255) NOT NULL DEFAULT '',
  `jump_server` int(11) NOT NULL DEFAULT '0',
  `after_script` text NOT NULL,
  `last_ver` varchar(255) NOT NULL DEFAULT '',
  `uuid` varchar(255) NOT NULL DEFAULT '',
  `keep_version_cnt` int(11) NOT NULL DEFAULT '10',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`env_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table projects
# ------------------------------------------------------------

CREATE TABLE `projects` (
  `project_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `project_name` varchar(255) NOT NULL DEFAULT '',
  `repo_url` varchar(255) NOT NULL DEFAULT '',
  `dst` varchar(255) NOT NULL DEFAULT '',
  `web_root` varchar(255) NOT NULL DEFAULT '',
  `after_script` text NOT NULL,
  `err_msg` text NOT NULL,
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `status` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`project_id`),
  UNIQUE KEY `uniq-dst` (`dst`),
  UNIQUE KEY `uniq-repo_url` (`repo_url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table servers
# ------------------------------------------------------------

CREATE TABLE `servers` (
  `server_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `type` tinyint(4) NOT NULL DEFAULT '1',
  `ssh_addr` varchar(100) NOT NULL DEFAULT '',
  `ssh_port` int(11) NOT NULL,
  `ssh_user` varchar(100) NOT NULL DEFAULT '',
  `ssh_key_path` varchar(255) NOT NULL DEFAULT '',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`server_id`),
  UNIQUE KEY `uniq-type-address` (`type`,`ssh_addr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table users
# ------------------------------------------------------------

CREATE TABLE `users` (
  `user_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `nick_name` varchar(255) NOT NULL DEFAULT '',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;

INSERT INTO `users` (`user_id`, `user_name`, `password`, `nick_name`, `create_at`, `status`)
VALUES
	(1,'admin','e10adc3949ba59abbe56e057f20f883e','God','2020-12-08 11:55:22',1);

/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
