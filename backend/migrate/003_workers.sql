SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `workers`;

CREATE TABLE IF NOT EXISTS `workers` (
    `worker_id` varchar(50) NOT NULL,
    `user_id` varchar(50) NOT NULL,
    `name` varchar(255) NOT NULL,
    `email` varchar(255) DEFAULT NULL,
    `role` enum('pic', 'manager', 'worker') NOT NULL,
    `fin_nric` varchar(50) DEFAULT NULL,
    `trade_code` varchar(50) DEFAULT NULL,
    `current_project_id` varchar(50) DEFAULT NULL,
    `company_name` varchar(255) DEFAULT NULL,
    `status` enum(
        'active',
        'inactive',
        'suspended',
        'archived'
    ) NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`worker_id`),
    UNIQUE KEY `email` (`email`),
    KEY `user_id` (`user_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `workers_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;