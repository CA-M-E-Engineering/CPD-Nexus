SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `workers`;

CREATE TABLE IF NOT EXISTS `workers` (
    `worker_id` varchar(50) NOT NULL,
    `user_id` varchar(50) NOT NULL,
    `name` varchar(255) NOT NULL,
    `email` varchar(255) DEFAULT NULL,
    `role` enum('pic', 'manager', 'worker') NOT NULL,
    `user_type` enum(
        'user',
        'visitor',
        'blocklist'
    ) NOT NULL DEFAULT 'user',
    `person_id_no` varchar(50) DEFAULT NULL,
    `person_id_and_work_pass_type` enum(
        'SP',
        'SB',
        'EP',
        'SPASS',
        'WP',
        'ENTREPASS',
        'LTVP'
    ) DEFAULT NULL,
    `person_nationality` char(2) DEFAULT NULL,
    `person_trade` varchar(10) DEFAULT NULL,
    `current_project_id` varchar(50) DEFAULT NULL,
    `status` enum(
        'active',
        'inactive',
        'suspended',
        'archived'
    ) NOT NULL,
    `auth_start_time` datetime DEFAULT NULL,
    `auth_end_time` datetime DEFAULT NULL,
    `fdid` int NOT NULL DEFAULT 1,
    `face_img_loc` varchar(255) DEFAULT NULL,
    `card_number` varchar(100) DEFAULT NULL,
    `card_type` varchar(50) DEFAULT NULL,
    `is_synced` tinyint(1) NOT NULL DEFAULT 0,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`worker_id`),
    UNIQUE KEY `email` (`email`),
    KEY `user_id` (`user_id`),
    KEY `idx_status` (`status`),
    CONSTRAINT `workers_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;