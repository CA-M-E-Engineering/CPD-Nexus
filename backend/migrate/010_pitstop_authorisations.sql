SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `pitstop_authorisations`;

CREATE TABLE IF NOT EXISTS `pitstop_authorisations` (
    `pitstop_auth_id` varchar(50) NOT NULL,
    `dataset_id` varchar(50) NOT NULL,
    `dataset_name` varchar(255) NOT NULL,
    `user_id` varchar(50) DEFAULT NULL,
    `regulator_id` char(36) NOT NULL,
    `regulator_name` varchar(255) NOT NULL,
    `on_behalf_of_id` char(36) NOT NULL,
    `on_behalf_of_name` varchar(255) NOT NULL,
    `status` varchar(20) NOT NULL DEFAULT 'ACTIVE',
    `last_synced_at` datetime DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`pitstop_auth_id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;