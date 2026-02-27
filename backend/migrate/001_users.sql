SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `users`;

CREATE TABLE IF NOT EXISTS `users` (
    `user_id` varchar(50) NOT NULL,
    `user_name` varchar(255) NOT NULL,
    `user_type` enum(
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
    `bridge_ws_url` varchar(255) DEFAULT NULL COMMENT 'WebSocket URL for the user''s IoT Bridge',
    `bridge_auth_token` varchar(255) DEFAULT NULL COMMENT 'Optional auth token for the bridge',
    `bridge_status` enum('active', 'inactive') NOT NULL DEFAULT 'inactive' COMMENT 'Whether the bridge connection should be active',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `username` (`username`),
    KEY `idx_status` (`status`),
    KEY `idx_user_type` (`user_type`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;