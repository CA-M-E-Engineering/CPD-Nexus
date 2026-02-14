SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `projects`;

CREATE TABLE IF NOT EXISTS `projects` (
    `project_id` varchar(50) NOT NULL,
    `user_id` varchar(50) DEFAULT NULL,
    `site_id` varchar(50) DEFAULT NULL,
    `project_reference_number` varchar(50) NOT NULL,
    `project_title` varchar(255) NOT NULL,
    `project_contract_number` varchar(100) DEFAULT NULL,
    `project_contract_name` varchar(100) DEFAULT NULL,
    `project_location_description` varchar(255) DEFAULT NULL,
    `hdb_precinct_name` varchar(100) DEFAULT NULL,
    `main_contractor_name` varchar(255) DEFAULT NULL,
    `main_contractor_uen` varchar(50) DEFAULT NULL,
    `offsite_fabricator_name` varchar(255) DEFAULT NULL,
    `offsite_fabricator_uen` varchar(50) DEFAULT NULL,
    `offsite_fabricator_location` varchar(255) DEFAULT NULL,
    `worker_company_name` varchar(255) DEFAULT NULL,
    `worker_company_uen` varchar(50) DEFAULT NULL,
    `worker_company_client_name` varchar(255) DEFAULT NULL,
    `worker_company_client_uen` varchar(50) DEFAULT NULL,
    `status` enum(
        'active',
        'completed',
        'archived',
        'inactive'
    ) NOT NULL DEFAULT 'active',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`project_id`),
    KEY `site_id` (`site_id`),
    KEY `idx_status` (`status`),
    KEY `fk_projects_user` (`user_id`),
    CONSTRAINT `fk_projects_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
    CONSTRAINT `projects_ibfk_1` FOREIGN KEY (`site_id`) REFERENCES `sites` (`site_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;