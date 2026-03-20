SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `bridge_logs`;

CREATE TABLE IF NOT EXISTS `bridge_logs` (
    `id` int NOT NULL AUTO_INCREMENT,
    `user_id` varchar(255) NOT NULL,
    `action` varchar(100) NOT NULL,
    `request_id` varchar(100) NOT NULL,
    `request_payload` json DEFAULT NULL,
    `response_payload` json DEFAULT NULL,
    `status_code` int DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_request_id` (`request_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_action` (`action`),
    KEY `idx_created_at` (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;
