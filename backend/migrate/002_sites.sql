SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `sites`;

CREATE TABLE IF NOT EXISTS `sites` (
    `site_id` varchar(50) NOT NULL,
    `user_id` varchar(50) NOT NULL,
    `site_name` varchar(255) NOT NULL,
    `location` varchar(255) DEFAULT NULL,
    `latitude` float DEFAULT NULL,
    `longitude` float DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`site_id`),
    KEY `user_id` (`user_id`),
    CONSTRAINT `sites_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;