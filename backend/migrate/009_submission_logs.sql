SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `submission_logs`;

CREATE TABLE IF NOT EXISTS `submission_logs` (
    `log_id` int NOT NULL AUTO_INCREMENT,
    `data_element_id` varchar(100) NOT NULL,
    `internal_id` varchar(255) NOT NULL,
    `status` enum('submitted', 'failed') NOT NULL,
    `payload` json DEFAULT NULL,
    `error_message` text,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`log_id`),
    KEY `idx_data_element` (`data_element_id`),
    KEY `idx_internal_id` (`internal_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;