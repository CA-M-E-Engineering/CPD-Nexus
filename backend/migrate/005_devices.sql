SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `devices`;

CREATE TABLE IF NOT EXISTS `devices` (
    `device_id` varchar(50) NOT NULL,
    `sn` varchar(100) NOT NULL,
    `tenant_id` varchar(50) NOT NULL,
    `site_id` varchar(50) DEFAULT NULL,
    `model` varchar(100) NOT NULL,
    `status` enum(
        'online',
        'offline',
        'unknown',
        'inactive'
    ) NOT NULL,
    `last_heartbeat` timestamp NULL DEFAULT NULL,
    `last_online_check` datetime DEFAULT NULL,
    `battery` int DEFAULT '100',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`device_id`),
    KEY `tenant_id` (`tenant_id`),
    KEY `site_id` (`site_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `devices_ibfk_1` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`tenant_id`),
    CONSTRAINT `devices_ibfk_2` FOREIGN KEY (`site_id`) REFERENCES `sites` (`site_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;