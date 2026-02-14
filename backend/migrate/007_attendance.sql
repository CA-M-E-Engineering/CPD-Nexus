SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `attendance`;

CREATE TABLE IF NOT EXISTS `attendance` (
    `attendance_id` char(36) NOT NULL,
    `device_id` varchar(50) NOT NULL,
    `worker_id` varchar(50) NOT NULL,
    `site_id` varchar(50) DEFAULT NULL,
    `tenant_id` varchar(50) DEFAULT NULL,
    `time_in` timestamp NULL DEFAULT NULL,
    `time_out` timestamp NULL DEFAULT NULL,
    `direction` enum('entry', 'exit', 'unknown') NOT NULL,
    `trade_code` varchar(10) NOT NULL,
    `status` enum(
        'pending',
        'submitted',
        'failed'
    ) NOT NULL DEFAULT 'pending',
    `submission_date` date NOT NULL,
    `batch_id` char(36) DEFAULT NULL,
    `response_payload` json DEFAULT NULL,
    `retry_count` int DEFAULT '0',
    `error_message` text,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`attendance_id`),
    KEY `worker_id` (`worker_id`),
    CONSTRAINT `attendance_ibfk_1` FOREIGN KEY (`worker_id`) REFERENCES `users` (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;