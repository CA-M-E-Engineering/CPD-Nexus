SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `system_settings`;

CREATE TABLE IF NOT EXISTS `system_settings` (
    `id` int NOT NULL DEFAULT '1',
    `attendance_sync_time` TIME NOT NULL DEFAULT '23:00:00',
    `cpd_submission_time` TIME DEFAULT '09:00:00' COMMENT 'Daily time for CPD submission',
    `max_payload_size_kb` int DEFAULT '256' COMMENT 'Maximum SGBuildex payload size in KB',
    `max_workers_per_request` int DEFAULT '100' COMMENT 'Max workers per API request',
    `max_requests_per_minute` int DEFAULT '150' COMMENT 'API rate limit safety threshold',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

INSERT INTO
    `system_settings` (
        `id`,
        `attendance_sync_time`,
        `cpd_submission_time`,
        `max_payload_size_kb`,
        `max_workers_per_request`,
        `max_requests_per_minute`
    )
VALUES (
        1,
        '23:00:00',
        '09:00:00',
        256,
        100,
        150
    )
ON DUPLICATE KEY UPDATE
    id = id;

SET FOREIGN_KEY_CHECKS = 1;