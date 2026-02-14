SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE site_roles;

TRUNCATE TABLE attendance;

TRUNCATE TABLE devices;

TRUNCATE TABLE users;

TRUNCATE TABLE sites;

TRUNCATE TABLE tenants;

TRUNCATE TABLE projects;

-- ======================
-- Tenants (Accounts)
-- ======================
INSERT INTO
    tenants (
        tenant_id,
        tenant_name,
        tenant_type,
        username,
        password_hash,
        contact_email,
        contact_phone,
        address,
        latitude,
        longitude,
        status,
        created_at,
        updated_at
    )
VALUES (
        'tenant-vendor-1',
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
        NOW(),
        NOW()
    ),
    (
        'tenant-client-1',
        'Mega Account',
        'client',
        'mega_admin',
        '$2a$10$YGl2KZOrJ8oAtuyu5l59JuLCAeHZMfm15blSCSLwGAkfIU04c.F6G',
        'contact@mega.com',
        '60000002',
        '15 Jurong Gateway Road',
        1.3329,
        103.7436,
        'active',
        NOW(),
        NOW()
    ),
    (
        'testt.ltd',
        'Test Tenant',
        'internal',
        'testt.ltd',
        '$2a$10$YGl2KZOrJ8oAtuyu5l59JuLCAeHZMfm15blSCSLwGAkfIU04c.F6G',
        'admin@testt.ltd',
        '90000000',
        'Singapore',
        1.3521,
        103.8198,
        'active',
        NOW(),
        NOW()
    );

-- ======================
-- Sites (2 Total)
-- ======================
INSERT INTO
    sites (
        site_id,
        tenant_id,
        site_name,
        location,
        latitude,
        longitude,
        created_at,
        updated_at
    )
VALUES (
        'site-uuid-1',
        'tenant-client-1',
        'Mega Jurong Site',
        'Jurong',
        1.3329,
        103.7436,
        NOW(),
        NOW()
    ),
    (
        'site-uuid-2',
        'tenant-client-1',
        'Mega Changi Site',
        'Changi',
        1.3644,
        103.9915,
        NOW(),
        NOW()
    );

-- ======================
-- Projects (2 Total)
-- ======================
INSERT INTO
    projects (
        project_id,
        site_id,
        tenant_id,
        project_reference_number,
        project_title,
        project_contract_number,
        project_contract_name,
        project_location_description,
        hdb_precinct_name,
        main_contractor_name,
        main_contractor_uen,
        offsite_fabricator_name,
        offsite_fabricator_uen,
        offsite_fabricator_location,
        worker_company_name,
        worker_company_uen,
        worker_company_client_name,
        worker_company_client_uen,
        status,
        created_at,
        updated_at
    )
VALUES (
        'project-uuid-1',
        'site-uuid-1',
        'tenant-client-1',
        'PRJ-2026-001',
        'Jurong Warehouse Expansion',
        'CNT-26-001',
        'Structural Phase A',
        'Jurong West St 21',
        'JW-N4-C1',
        'Mega Engineering',
        'MEGA12345X',
        'Delta Fabrication Ltd',
        'UEN-FAB-001',
        '10 Industrial Way, Singapore',
        'Mega Engineering',
        'MEGA12345X',
        'CA M&E (Vendor)',
        'CAMEN1234A',
        'active',
        NOW(),
        NOW()
    ),
    (
        'project-uuid-2',
        'site-uuid-2',
        'tenant-client-1',
        'PRJ-2026-002',
        'Changi Airport T5 Ground',
        'CNT-26-002',
        'Foundation Civil Works',
        'Changi East Drive',
        'CH-T5-P1',
        'Mega Engineering',
        'MEGA12345X',
        NULL,
        NULL,
        NULL,
        'Mega Engineering',
        'MEGA12345X',
        NULL,
        NULL,
        'active',
        NOW(),
        NOW()
    );

-- ======================
-- Users
-- ======================
INSERT INTO
    users (
        user_id,
        tenant_id,
        name,
        email,
        role,
        fin_nric,
        trade_code,
        current_project_id,
        company_name,
        status,
        created_at,
        updated_at
    )
VALUES (
        'user-vendor-admin',
        'tenant-vendor-1',
        'Robert CA Admin',
        'robert@came.com',
        'manager',
        'S1000001A',
        NULL,
        NULL,
        'CA M&E (Vendor)',
        'active',
        NOW(),
        NOW()
    ),
    (
        'user-client-admin-1',
        'tenant-client-1',
        'Mega Admin',
        'admin@mega.com',
        'manager',
        'S2000001A',
        NULL,
        NULL,
        'Mega Engineering',
        'active',
        NOW(),
        NOW()
    ),
    (
        'user-worker-1',
        'tenant-client-1',
        'Worker Ali',
        'ali@mega.com',
        'worker',
        'S8000001W',
        'ELEC',
        'project-uuid-1',
        'Mega Engineering',
        'active',
        NOW(),
        NOW()
    ),
    (
        'user-worker-2',
        'tenant-client-1',
        'Worker Bob',
        'bob@mega.com',
        'worker',
        'S8000002W',
        'PLUMB',
        'project-uuid-1',
        'Mega Engineering',
        'active',
        NOW(),
        NOW()
    ),
    (
        'user-worker-3',
        'tenant-vendor-1',
        'Worker Charlie',
        'charlie@came.com',
        'worker',
        'S8000003W',
        'WELD',
        'project-uuid-2',
        'CA M&E (Vendor)',
        'active',
        NOW(),
        NOW()
    ),
    (
        'user-worker-4',
        'tenant-vendor-1',
        'Worker David',
        'david@came.com',
        'worker',
        'S8000004W',
        'SCAF',
        'project-uuid-2',
        'CA M&E (Vendor)',
        'active',
        NOW(),
        NOW()
    ),
    (
        'user-worker-5',
        'tenant-vendor-1',
        'Worker Eve',
        'eve@came.com',
        'worker',
        'S8000005W',
        'PAINT',
        NULL,
        'CA M&E (Vendor)',
        'active',
        NOW(),
        NOW()
    ),
    (
        'user-pic-1',
        'tenant-client-1',
        'PIC John Doe',
        'john@mega.com',
        'pic',
        'S7000001P',
        'MGMT',
        'project-uuid-1',
        'Mega Engineering',
        'active',
        NOW(),
        NOW()
    );

-- ======================
-- Site Roles
-- ======================
INSERT INTO
    site_roles (
        site_role_id,
        site_id,
        user_id,
        role,
        is_primary,
        created_at,
        updated_at
    )
VALUES (
        'sr-1',
        'site-uuid-1',
        'user-worker-1',
        'sub_contractor',
        TRUE,
        NOW(),
        NOW()
    ),
    (
        'sr-2',
        'site-uuid-1',
        'user-worker-2',
        'sub_contractor',
        TRUE,
        NOW(),
        NOW()
    ),
    (
        'sr-3',
        'site-uuid-2',
        'user-worker-3',
        'sub_contractor',
        TRUE,
        NOW(),
        NOW()
    ),
    (
        'sr-4',
        'site-uuid-2',
        'user-worker-4',
        'sub_contractor',
        TRUE,
        NOW(),
        NOW()
    );

-- ======================
-- Devices (5 Total)
-- ======================
INSERT INTO
    devices (
        device_id,
        sn,
        tenant_id,
        site_id,
        model,
        status,
        last_heartbeat,
        created_at,
        updated_at
    )
VALUES (
        'device-001',
        'SN-GW-001',
        'tenant-client-1',
        'site-uuid-1',
        'Gateway-X1',
        'online',
        NOW(),
        NOW(),
        NOW()
    ),
    (
        'device-002',
        'SN-LOC-002',
        'tenant-client-1',
        'site-uuid-1',
        'Locate-P1',
        'online',
        NOW(),
        NOW(),
        NOW()
    ),
    (
        'device-003',
        'SN-GW-003',
        'tenant-client-1',
        'site-uuid-2',
        'Gateway-X1',
        'offline',
        NOW(),
        NOW(),
        NOW()
    ),
    (
        'device-004',
        'SN-ENV-004',
        'tenant-vendor-1',
        NULL,
        'Env-S1',
        'inactive',
        NOW(),
        NOW(),
        NOW()
    ),
    (
        'device-005',
        'SN-GW-005',
        'tenant-vendor-1',
        NULL,
        'Gateway-X2',
        'inactive',
        NOW(),
        NOW(),
        NOW()
    );

-- ======================
-- Attendance
-- ======================
INSERT INTO
    attendance (
        attendance_id,
        device_id,
        worker_id,
        site_id,
        tenant_id,
        time_in,
        time_out,
        direction,
        trade_code,
        status,
        submission_date,
        created_at,
        updated_at
    )
VALUES (
        'att-1',
        'device-001',
        'user-worker-1',
        'site-uuid-1',
        'tenant-vendor-1',
        '2026-02-10 07:34:47',
        '2026-02-10 10:34:47',
        'exit',
        'ELEC',
        'submitted',
        '2026-02-10',
        NOW(),
        NOW()
    ),
    (
        'att-2',
        'device-001',
        'user-worker-2',
        'site-uuid-1',
        'tenant-vendor-1',
        '2026-02-10 07:34:47',
        NULL,
        'entry',
        'PLUMB',
        'submitted',
        '2026-02-10',
        NOW(),
        NOW()
    ),
    (
        'att-3',
        'device-003',
        'user-worker-3',
        'site-uuid-2',
        'tenant-vendor-1',
        '2026-02-10 07:34:47',
        '2026-02-10 10:34:47',
        'exit',
        'WELD',
        'submitted',
        '2026-02-10',
        NOW(),
        NOW()
    ),
    (
        'att-4',
        'device-003',
        'user-worker-4',
        'site-uuid-2',
        'tenant-vendor-1',
        '2026-02-10 07:34:47',
        NULL,
        'entry',
        'SCAF',
        'submitted',
        '2026-02-10',
        NOW(),
        NOW()
    );

-- ======================
-- System Settings
-- ======================
INSERT INTO
    system_settings (
        id,
        device_sync_interval,
        cpd_submission_time,
        response_size_limit
    )
VALUES (
        1,
        '00:01:00',
        '09:00:00',
        1048576
    )
ON DUPLICATE KEY UPDATE
    device_sync_interval = '00:01:00',
    cpd_submission_time = '09:00:00',
    response_size_limit = 1048576;

SET FOREIGN_KEY_CHECKS = 1;