SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `activity_logs`;

CREATE TABLE IF NOT EXISTS `activity_logs` (
    `id` int NOT NULL AUTO_INCREMENT,
    `user_id` varchar(255) NOT NULL,
    `user_name` varchar(255) DEFAULT NULL,
    `action` varchar(255) NOT NULL,
    `target_type` varchar(50) DEFAULT NULL,
    `target_id` varchar(255) DEFAULT NULL,
    `details` text,
    `ip_address` varchar(45) DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_action` (`action`),
    KEY `idx_created_at` (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;
