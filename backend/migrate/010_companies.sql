SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `companies`;

CREATE TABLE IF NOT EXISTS `companies` (
    `company_id` varchar(50) NOT NULL,
    `tenant_id` varchar(50) NOT NULL,
    `company_name` varchar(255) NOT NULL,
    `uen` varchar(50) NOT NULL,
    `company_type` enum(
        'contractor',
        'offsite_fabricator'
    ) NOT NULL,
    `address` varchar(255) DEFAULT NULL,
    `latitude` decimal(9, 6) DEFAULT NULL,
    `longitude` decimal(9, 6) DEFAULT NULL,
    `status` enum('active', 'inactive') NOT NULL DEFAULT 'active',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`company_id`),
    UNIQUE KEY `idx_uen` (`uen`),
    KEY `tenant_id` (`tenant_id`),
    CONSTRAINT `companies_ibfk_1` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`tenant_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;