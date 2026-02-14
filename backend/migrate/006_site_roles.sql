SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `site_roles`;

CREATE TABLE IF NOT EXISTS `site_roles` (
    `site_role_id` char(36) NOT NULL,
    `site_id` varchar(50) NOT NULL,
    `worker_id` varchar(50) NOT NULL,
    `role` enum(
        'main_contractor',
        'sub_contractor',
        'client',
        'consultant'
    ) NOT NULL,
    `is_primary` tinyint(1) DEFAULT '0',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`site_role_id`),
    KEY `site_id` (`site_id`),
    KEY `worker_id` (`worker_id`),
    CONSTRAINT `site_roles_ibfk_1` FOREIGN KEY (`site_id`) REFERENCES `sites` (`site_id`),
    CONSTRAINT `site_roles_ibfk_2` FOREIGN KEY (`worker_id`) REFERENCES `workers` (`worker_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;