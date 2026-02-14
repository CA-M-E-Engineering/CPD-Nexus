# MVP Database Design (Updated)

This document describes the database schema for the MVP system, including tenants, sites, devices, users, and attendance records. Each table includes field descriptions, types, and constraints.

---

## 1️⃣ Tenants Table

**Purpose:** Represents each company/vendor using the system.

| Field         | Type                               | Mandatory | Description                                        |
| ------------- | --------------------------------- | --------- | -------------------------------------------------- |
| tenant_id     | UUID / PK                          | ✅        | Unique ID for the company/vendor                  |
| tenant_name   | VARCHAR                             | ✅        | Company name                                      |
| tenant_type   | ENUM('client','vendor','internal') | ✅        | Type of tenant                                     |
| contact_email | VARCHAR                             | ✅        | Main contact email                                |
| contact_phone | VARCHAR                             | ❌        | Optional phone number                             |
| is_bca       | BOOLEAN                             | ✅        | Whether the company is BCA-registered             |
| status        | ENUM('active','inactive')           | ✅        | Tenant status                                     |
| created_at    | TIMESTAMP                           | ✅        | Record creation timestamp                         |
| updated_at    | TIMESTAMP                           | ✅        | Last modification timestamp                       |

---

## 2️⃣ Sites Table

**Purpose:** Represents projects or construction sites under a tenant.

| Field       | Type      | Mandatory | Description                         |
|------------ | -------- | --------- | ----------------------------------- |
| site_id     | UUID / PK| ✅        | Unique project/site ID              |
| tenant_id   | UUID / FK| ✅        | References tenants                  |
| site_name   | VARCHAR  | ✅        | Human-readable site name            |
| location    | VARCHAR  | ❌        | Optional address                    |
| latitude    | VARCHAR  | ❌        | Optional latitude                   |
| longitude   | VARCHAR  | ❌        | Optional longitude                  |
| status      | ENUM('active','inactive') | ✅ | Site status                        |
| created_at  | TIMESTAMP| ✅        | Record creation timestamp           |
| updated_at  | TIMESTAMP| ✅        | Last modification timestamp         |

---

## 3️⃣ Devices Table

**Purpose:** Stores biometric device registration and runtime status.

| Field           | Type                           | Mandatory | Description                             |
|---------------- | ----------------------------- | --------- | --------------------------------------- |
| device_id       | UUID / PK                      | ✅        | Unique device ID                         |
| sn              | VARCHAR                        | ✅        | Serial number of device                  |
| tenant_id       | UUID / FK                      | ✅        | Client that owns the device              |
| site_id         | UUID / FK                      | ✅        | Site/project installed at                |
| model           | VARCHAR                        | ✅        | Device model (e.g., DS-K1T341AMF)      |
| status          | ENUM('online','offline','unknown') | ✅    | Current device status                    |
| last_heartbeat  | TIMESTAMP                      | ❌        | Last heartbeat timestamp (nullable)      |
| created_at      | TIMESTAMP                      | ✅        | Device registration timestamp            |
| updated_at      | TIMESTAMP                      | ✅        | Last modification timestamp              |

---

## 4️⃣ Users Table

**Purpose:** Stores login-enabled users (admin/operator) and worker identities.

| Field          | Type                                           | Mandatory | Description                                           |
|--------------- | --------------------------------------------- | --------- | ---------------------------------------------------- |
| user_id        | UUID / PK                                     | ✅        | Unique user ID                                       |
| tenant_id      | UUID / FK                                     | ✅        | Tenant/company this user belongs to                 |
| name           | VARCHAR                                       | ✅        | Full name                                           |
| email          | VARCHAR UNIQUE                                | ✅        | Login email                                         |
| password_hash  | VARCHAR                                       | ❌        | Hashed password (optional for non-login workers)   |
| role           | ENUM('super_admin','admin','operator','worker') | ✅      | Role controlling permissions                        |
| fin_nric       | VARCHAR (encrypted)                           | ❌        | FIN / NRIC (only for workers)                       |
| trade_code     | VARCHAR                                       | ❌        | Trade code (only for workers)                       |
| status         | ENUM('active','inactive','suspended','archived') | ✅      | Account or worker status                             |
| created_at     | TIMESTAMP                                     | ✅        | Record creation timestamp                             |
| updated_at     | TIMESTAMP                                     | ✅        | Last modification timestamp                           |

---

## 5️⃣ Attendance Records Table

**Purpose:** Stores attendance logs from devices and tracks submission status. Each record represents a single time-in/time-out.

| Field             | Type                          | Mandatory | Description                                   |
|----------------- | ---------------------------- | --------- | -------------------------------------------- |
| attendance_id     | UUID / PK                   | ✅        | Unique attendance record ID                  |
| device_id         | UUID / FK                   | ✅        | Device that generated the record            |
| user_id           | UUID / FK                   | ✅        | Worker associated with the record           |
| site_id           | UUID / FK                   | ✅        | Site/project where log occurred             |
| tenant_id         | UUID / FK                   | ✅        | Tenant reference                             |
| time_in           | TIMESTAMP                   | ✅        | Exact time-in timestamp                      |
| time_out          | TIMESTAMP                   | ✅        | Exact time-out timestamp                     |
| status            | ENUM('pending','submitted','failed') | ✅  | Submission status to SGBuilderx             |
| submission_date   | DATE                         | ✅        | Date of the submission batch                |
| batch_id          | UUID                         | ❌        | Optional batch grouping for submissions     |
| response_payload  | JSON / TEXT                  | ❌        | API response from SGBuilderx if submitted  |
| retry_count       | INT (default 0)              | ✅        | Number of retries attempted                 |
| error_message     | TEXT                         | ❌        | Error reason if submission failed           |
| created_at        | TIMESTAMP                     | ✅        | Record creation timestamp                    |
| updated_at        | TIMESTAMP                     | ✅        | Last modification timestamp                  |

---

**Notes:**

- `ManpowerUtilization` payload is derived from **Users, Sites, Devices, and Attendance records**.  
- `Participants` and `OnBehalfOf` are stored in the request metadata for API submission, not persisted separately.  
- Each **attendance record** corresponds to a single time-in/time-out; multiple records per worker per day are allowed.  
- Companies now include `is_bca` boolean to identify BCA-registered tenants.  
- No separate `manpower` table is needed; data can be aggregated from existing tables.

