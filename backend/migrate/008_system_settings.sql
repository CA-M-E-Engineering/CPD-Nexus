SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `system_settings`;

CREATE TABLE IF NOT EXISTS `system_settings` (
    `id` int NOT NULL DEFAULT '1',
    `device_sync_interval` TIME NOT NULL DEFAULT '00:01:00',
    `cpd_submission_time` TIME DEFAULT '09:00:00' COMMENT 'Daily time for CPD submission',
    `response_size_limit` bigint DEFAULT '1048576' COMMENT 'Maximum response size in bytes',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

INSERT INTO
    `system_settings` (
        `id`,
        `device_sync_interval`,
        `cpd_submission_time`,
        `response_size_limit`
    )
VALUES (
        1,
        '00:01:00',
        '09:00:00',
        1048576
    )
ON DUPLICATE KEY UPDATE
    id = id;

SET FOREIGN_KEY_CHECKS = 1;