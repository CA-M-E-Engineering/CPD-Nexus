# CPD-Nexus — Database Schema Reference

This document describes all tables in the CPD-Nexus MySQL schema, their purpose, and all field definitions. ID columns use human-readable prefixed strings (e.g. `user_001`, `site_003`) rather than UUIDs for traceability across logs and support tickets.

---

## 1. `users`

Represents accounts in the system. The `user_type` determines which UI pages and API scope the account is granted.

| Column | Type | Nullable | Description |
|---|---|---|---|
| `user_id` | VARCHAR(50) PK | ❌ | Unique ID, e.g. `user_001` |
| `user_name` | VARCHAR(255) | ❌ | Organisation or person display name |
| `user_type` | ENUM('client','vendor','internal') | ❌ | Role scope |
| `username` | VARCHAR(255) UNIQUE | ✅ | Login handle |
| `password_hash` | VARCHAR(255) | ✅ | Bcrypt hash |
| `contact_email` | VARCHAR(255) | ❌ | Primary contact email |
| `status` | ENUM('active','inactive') | ❌ | Account status |
| `bridge_url` | VARCHAR(500) | ✅ | WebSocket URL for IoT bridge connection |
| `bridge_auth_token` | VARCHAR(255) | ✅ | Auth token sent in bridge message `meta` |
| `bridge_is_synced` | TINYINT(1) | ❌ | Whether the bridge is actively connected |
| `created_at` | TIMESTAMP | ❌ | Auto-set on insert |

---

## 2. `sites`

A physical construction location. Devices and projects are assigned to sites.

| Column | Type | Nullable | Description |
|---|---|---|---|
| `site_id` | VARCHAR(50) PK | ❌ | Unique ID, e.g. `site_001` |
| `user_id` | VARCHAR(50) FK → users | ❌ | Owning tenant |
| `site_name` | VARCHAR(255) | ❌ | Human-readable site name |
| `location` | VARCHAR(255) | ✅ | Address or location description |
| `latitude` | FLOAT | ✅ | GPS latitude for geofencing |
| `longitude` | FLOAT | ✅ | GPS longitude for geofencing |

---

## 3. `projects`

A BCA-registered construction project within a site. Workers are assigned to projects; projects link to Pitstop authorisations for submission routing.

| Column | Type | Nullable | Description |
|---|---|---|---|
| `project_id` | VARCHAR(50) PK | ❌ | Unique ID, e.g. `p20260301120530` |
| `user_id` | VARCHAR(50) FK → users | ❌ | Owning tenant |
| `site_id` | VARCHAR(50) FK → sites | ✅ | Linked construction site |
| `project_title` | VARCHAR(1000) | ✅ | Project name (max 1000 chars) |
| `project_reference_number` | VARCHAR(50) | ✅ | BCA ref, e.g. `A1234-AB123-2022` |
| `project_location_description` | TEXT | ✅ | Location text (max 2000 chars) |
| `project_contract_number` | VARCHAR(50) | ✅ | HDB (`D/NNNNN/YY`) or LTA (max 20 chars) |
| `project_contract_name` | VARCHAR(100) | ✅ | Contract name (max 100 chars) |
| `hdb_precinct_name` | VARCHAR(40) | ✅ | HDB precinct (max 40 chars) |
| `main_contractor_name` | VARCHAR(255) | ✅ | Main contractor company name |
| `main_contractor_uen` | VARCHAR(20) | ✅ | Main contractor UEN |
| `worker_company_name` | VARCHAR(255) | ✅ | Employer company name |
| `worker_company_uen` | VARCHAR(20) | ✅ | Employer UEN |
| `worker_company_trade` | VARCHAR(100) | ✅ | Comma-separated BCA trade codes |
| `worker_company_client_name` | VARCHAR(255) | ✅ | Employer client company name |
| `worker_company_client_uen` | VARCHAR(20) | ✅ | Employer client UEN |
| `pitstop_auth_id` | VARCHAR(50) FK → pitstop_authorisations | ✅ | Linked Pitstop routing config |
| `status` | VARCHAR(20) | ✅ | `active`, `inactive`, `completed` |

---

## 4. `workers`

Personnel registered to the system. Biometric and BCA compliance fields are stored here.

| Column | Type | Nullable | Description |
|---|---|---|---|
| `worker_id` | VARCHAR(50) PK | ❌ | Unique ID, e.g. `w20260301120530` |
| `user_id` | VARCHAR(50) FK → users | ❌ | Employer organisation |
| `name` | VARCHAR(100) | ❌ | Full legal name |
| `email` | VARCHAR(255) | ✅ | Contact email |
| `role` | ENUM('worker','manager','pic') | ❌ | Functional role on site |
| `user_type` | VARCHAR(20) | ✅ | Bridge user type (e.g. `normal`) |
| `status` | ENUM('active','inactive') | ❌ | Employment status |
| `person_id_no` | VARCHAR(20) | ✅ | FIN / NRIC (validated with regex) |
| `person_id_and_work_pass_type` | VARCHAR(10) | ✅ | Enum: `WP`, `SP`, `EP`, `SPASS`, `ENTREPASS`, `LTVP`, `SB` |
| `person_trade` | VARCHAR(10) | ✅ | BCA trade code (e.g. `2.5`, `3.11`) |
| `person_nationality` | CHAR(2) | ✅ | ISO 3166 country code |
| `site_id` | VARCHAR(50) FK → sites | ✅ | Assigned construction site |
| `current_project_id` | VARCHAR(50) FK → projects | ✅ | Current project assignment |
| `face_img_loc` | VARCHAR(500) | ✅ | URL to face image |
| `card_number` | VARCHAR(50) | ✅ | Access card number |
| `card_type` | VARCHAR(20) | ✅ | Card type for bridge |
| `fdid` | INT | ✅ | Face data ID on device |
| `auth_start_time` | VARCHAR(50) | ✅ | Bridge validity start (RFC3339) |
| `auth_end_time` | VARCHAR(50) | ✅ | Bridge validity end (RFC3339) |
| `is_synced` | TINYINT | ❌ | `0` = Pending Update, `1` = Synced, `2` = Pending Registration |

> **`is_synced` lifecycle:** Set to `2` (Pending Registration) on create if biometric data exists. Set to `0` (Pending Update) when biometric or auth fields change. Set to `1` (Synced) only after receiving a successful bridge response.

---

## 5. `devices`

IoT biometric gateways. Each device is identified by its hardware Serial Number.

| Column | Type | Nullable | Description |
|---|---|---|---|
| `device_id` | VARCHAR(50) PK | ❌ | Unique ID, e.g. `device_001` |
| `sn` | VARCHAR(100) UNIQUE | ❌ | Physical serial number |
| `model` | VARCHAR(50) | ✅ | Hardware model |
| `user_id` | VARCHAR(50) FK → users | ❌ | Owner |
| `site_id` | VARCHAR(50) FK → sites | ✅ | Current deployment site |
| `status` | ENUM('online','offline','error','inactive') | ❌ | Bridge connectivity status |
| `terminal_id` | VARCHAR(100) | ✅ | Bridge terminal identifier |
| `last_heartbeat` | TIMESTAMP | ✅ | Last event timestamp from bridge |

---

## 6. `attendance`

Primary log of all time-in/time-out events. Records start as `pending` and transition to `submitted` or `failed` after SGTradeX submission.

| Column | Type | Nullable | Description |
|---|---|---|---|
| `attendance_id` | VARCHAR(50) PK | ❌ | Format: `ATT-YYYYMMDD-NNNN` |
| `device_id` | VARCHAR(50) FK → devices | ❌ | Source device (`BRIDGE_AGGREGATED` for bridge-sourced) |
| `worker_id` | VARCHAR(50) FK → workers | ❌ | Recognised worker |
| `site_id` | VARCHAR(50) FK → sites | ❌ | Site of the event |
| `user_id` | VARCHAR(50) FK → users | ❌ | Tenant scope |
| `time_in` | DATETIME | ✅ | Entry timestamp (RFC3339) |
| `time_out` | DATETIME | ✅ | Exit timestamp — nullable if not yet checked out |
| `direction` | VARCHAR(20) | ✅ | `in`, `out`, or `unknown` |
| `trade_code` | VARCHAR(10) | ✅ | BCA trade code at time of event |
| `status` | ENUM('pending','submitted','failed') | ❌ | SGTradeX submission state |
| `submission_date` | DATE | ✅ | Calendar date of the event (used for monthly grouping) |
| `response_payload` | TEXT | ✅ | Raw JSON from last submission attempt |
| `error_message` | TEXT | ✅ | Error detail if status is `failed` |
| `created_at` | TIMESTAMP | ❌ | Record creation time |
| `updated_at` | TIMESTAMP | ❌ | Last status update time |

> **Submission rule:** Records where `status = 'submitted'` are excluded from all extraction queries and will never be re-submitted.

---

## 7. `pitstop_authorisations`

Caches the Pitstop routing configuration fetched from the `/api/v1/config` endpoint. Each row maps a dataset (e.g. `manpower_utilization`) to a regulator and an on-behalf-of contractor. Projects link to a specific row to resolve submission routing.

| Column | Type | Nullable | Description |
|---|---|---|---|
| `pitstop_auth_id` | VARCHAR(50) PK | ❌ | Unique ID, e.g. `pa202603011205001` |
| `dataset_id` | VARCHAR(50) | ❌ | API path segment, e.g. `manpower_utilization` |
| `dataset_name` | VARCHAR(255) | ✅ | Human-readable name |
| `user_id` | VARCHAR(50) FK → users | ✅ | Assigned contractor/vendor |
| `regulator_id` | CHAR(36) | ❌ | BCA UUID → `participants[].id` |
| `regulator_name` | VARCHAR(255) | ✅ | BCA display name |
| `on_behalf_of_id` | CHAR(36) | ❌ | Contractor UUID → `on_behalf_of[].id` |
| `on_behalf_of_name` | VARCHAR(255) | ✅ | Contractor display name |
| `status` | VARCHAR(20) | ❌ | `ACTIVE` or `INACTIVE` |
| `last_synced_at` | DATETIME | ✅ | Last sync from Pitstop API |

---

## 8. `submission_logs`

Audit trail for every SGTradeX API call. Written after each batch submission regardless of outcome.

| Column | Type | Nullable | Description |
|---|---|---|---|
| `id` | INT AUTO_INCREMENT PK | ❌ | Log entry ID |
| `data_element_id` | VARCHAR(50) | ❌ | E.g. `manpower_utilization` |
| `internal_id` | VARCHAR(50) | ❌ | Attendance ID or other internal reference |
| `status` | VARCHAR(20) | ❌ | `submitted` or `failed` |
| `payload` | TEXT | ✅ | Full JSON payload sent |
| `error_message` | TEXT | ✅ | API error detail if failed |
| `created_at` | TIMESTAMP | ❌ | Auto-set on insert |

---

## 9. `system_settings`

Single-row table for global system configuration. Used by the scheduler and submission worker.

| Column | Type | Nullable | Description |
|---|---|---|---|
| `id` | INT PK | ❌ | Always `1` (single row) |
| `max_workers_per_request` | INT | ❌ | Max items per SGTradeX batch (default: 100) |
| `max_payload_size_kb` | INT | ❌ | Max batch payload size in KB (default: 256) |
| `max_requests_per_minute` | INT | ❌ | Rate limit for SGTradeX calls (default: 150) |
| `attendance_sync_time` | VARCHAR(8) | ❌ | Bridge fetch trigger time `HH:MM:SS` |
| `cpd_submission_time` | VARCHAR(8) | ❌ | SGTradeX submission trigger time `HH:MM:SS` |

---

## Schema Design Principles

- **Prefixed IDs**: All primary keys use domain-prefixed strings (`user_`, `site_`, etc.) for immediate human readability in logs and debugs.
- **Multi-Tenant by Design**: Every user-owned table carries a `user_id` FK. All queries are scoped to prevent cross-tenant data access.
- **Soft Deletes via Status**: Entities are deactivated by setting `status = 'inactive'` rather than being hard-deleted, preserving audit history.
- **Attendance Immutability**: Once an attendance record is `submitted`, it is never modified. A new record must be created for corrections.
