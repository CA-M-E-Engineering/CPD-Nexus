SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `tenants`;

CREATE TABLE IF NOT EXISTS `tenants` (
    `tenant_id` varchar(50) NOT NULL,
    `tenant_name` varchar(255) NOT NULL,
    `tenant_type` enum(
        'client',
        'vendor',
        'internal'
    ) NOT NULL,
    `username` varchar(255) DEFAULT NULL,
    `password_hash` varchar(255) DEFAULT NULL,
    `contact_email` varchar(255) NOT NULL,
    `contact_phone` varchar(50) DEFAULT NULL,
    `latitude` decimal(10, 8) DEFAULT NULL,
    `longitude` decimal(11, 8) DEFAULT NULL,
    `address` varchar(255) DEFAULT NULL,
    `status` enum('active', 'inactive') NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`tenant_id`),
    UNIQUE KEY `username` (`username`),
    KEY `idx_status` (`status`),
    KEY `idx_tenant_type` (`tenant_type`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;