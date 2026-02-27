SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE attendance;

TRUNCATE TABLE devices;

TRUNCATE TABLE workers;

TRUNCATE TABLE sites;

TRUNCATE TABLE users;

TRUNCATE TABLE projects;

-- ======================
-- Users (Keep only Vendor Admin)
-- ======================
INSERT INTO
    users (
        user_id,
        user_name,
        user_type,
        username,
        password_hash,
        contact_email,
        contact_phone,
        address,
        latitude,
        longitude,
        status,
        bridge_ws_url,
        bridge_status,
        created_at,
        updated_at
    )
VALUES (
        'Owner_001',
        'CA M&E Account',
        'vendor',
        'came_admin',
        '$2a$10$YGl2KZOrJ8oAtuyu5l59JuLCAeHZMfm15blSCSLwGAkfIU04c.F6G',
        'admin@came.com',
        '60000001',
        '120 Lower Delta Road, Cendex Centre',
        1.3521,
        103.8198,
        'active',
        'ws://localhost:8081/ws',
        'active',
        NOW(),
        NOW()
    );

-- ======================
-- System Settings
-- ======================
INSERT INTO
    system_settings (
        id,
        attendance_sync_time,
        cpd_submission_time,
        max_payload_size_kb,
        max_workers_per_request,
        max_requests_per_minute
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
    attendance_sync_time = '23:00:00',
    cpd_submission_time = '09:00:00',
    max_payload_size_kb = 256,
    max_workers_per_request = 100,
    max_requests_per_minute = 150;

SET FOREIGN_KEY_CHECKS = 1;