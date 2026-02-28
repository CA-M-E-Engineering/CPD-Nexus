# Enterprise Database Architecture

This document describes the unified database schema for **CPD-Nexus**. The design relies on string-based auto-incrementing IDs with specific prefixes (e.g., `user_000`, `site_000`) instead of UUIDs, optimizing for readability while sustaining relational integrity in MySQL.

---

## 1️⃣ Users Table

**Purpose:** Represents the core hierarchy. Depending on `user_type`, a user acts as an administrative client, vendor, or internal operator. This replaces the legacy `tenants` table.

| Field           | Type                                 | Mandatory | Description                                      |
| --------------- | ------------------------------------ | --------- | ------------------------------------------------ |
| user_id         | VARCHAR(50) / PK                     | ✅        | Unique ID (e.g., `user_000`)                     |
| user_name       | VARCHAR(255)                         | ✅        | Formal organization or person name               |
| user_type       | ENUM('client','vendor','internal')   | ✅        | Role/scope of the account                        |
| username        | VARCHAR(255) / UNIQUE                | ❌        | Optional login handle                            |
| password_hash   | VARCHAR(255)                         | ❌        | Security payload for authentication              |
| contact_email   | VARCHAR(255)                         | ✅        | Admin email                                      |
| status          | ENUM('active','inactive')            | ✅        | Account status                                   |
| created_at      | TIMESTAMP                            | ✅        | Auto-generated creation time                     |

---

## 2️⃣ Sites Table

**Purpose:** Logical abstraction for real-world construction zones. Devices and projects are directly mapped to a specific site.

| Field       | Type                 | Mandatory | Description                           |
|------------ | -------------------- | --------- | ------------------------------------- |
| site_id     | VARCHAR(50) / PK     | ✅        | Unique ID (e.g., `site_000`)          |
| user_id     | VARCHAR(50) / FK     | ✅        | Governing client/vendor (`users`)     |
| site_name   | VARCHAR(255)         | ✅        | Human-readable site name              |
| location    | VARCHAR(255)         | ❌        | Geological map reference              |
| latitude    | FLOAT                | ❌        | Geofencing coord                      |
| longitude   | FLOAT                | ❌        | Geofencing coord                      |

---

## 3️⃣ Projects Table

**Purpose:** Sub-divisions within a site. A single construction site can have multiple concurrent internal projects over its lifecycle.

| Field                | Type                 | Mandatory | Description                             |
|--------------------- | -------------------- | --------- | --------------------------------------- |
| project_id           | VARCHAR(50) / PK     | ✅        | Unique ID (e.g., `project_000`)         |
| user_id              | VARCHAR(50) / FK     | ✅        | Governing vendor/client                 |
| site_id              | VARCHAR(50) / FK     | ❌        | Which physical site it belongs to       |
| project_name         | VARCHAR(255)         | ✅        | Designated project name                 |
| project_ref_no       | VARCHAR(50)          | ✅        | Official reference marker               |
| status               | ENUM                 | ✅        | active, inactive, completed             |

---

## 4️⃣ Workers Table

**Purpose:** Stores personnel registered to the deployment.

| Field                        | Type             | Mandatory | Description                                          |
|----------------------------- | ---------------- | --------- | ---------------------------------------------------- |
| worker_id                    | VARCHAR(50) / PK | ✅        | Unique ID (e.g., `worker_000`)                       |
| user_id                      | VARCHAR(50) / FK | ✅        | Employer organization                                |
| name                         | VARCHAR(100)     | ✅        | Full legal name                                      |
| person_id_no                 | VARCHAR(20)      | ✅        | FIN / NRIC parameter                                 |
| person_id_and_work_pass_type | VARCHAR(10)      | ✅        | WP, SP, EP, etc.                                     |
| person_trade                 | VARCHAR(10)      | ❌        | BCA designated trade index (e.g. 1.2, 2.5)          |
| role                         | ENUM             | ✅        | 'worker', 'manager', 'pic'                           |
| current_project_id           | VARCHAR(50) / FK | ❌        | Currently assigned project                           |
| is_synced                    | TINYINT          | ✅        | 0=Pending Update, 1=Synced, 2=Pending Registration. Set to 1 only upon receiving HTTP 200 OK from the Bridge asynchronously |

---

## 5️⃣ Devices Table

**Purpose:** IoT gateway/biometric registries responsible for fetching real-time attendance hardware streams.

| Field           | Type                 | Mandatory | Description                                 |
|---------------- | -------------------- | --------- | ------------------------------------------- |
| device_id       | VARCHAR(50) / PK     | ✅        | Unique ID (e.g., `device_000`)              |
| sn              | VARCHAR(100) / UNIQUE| ✅        | Physical Serial Number on gateway           |
| model           | VARCHAR(50)          | ✅        | Hardware Model                              |
| user_id         | VARCHAR(50) / FK     | ✅        | Owner of the device                         |
| site_id         | VARCHAR(50) / FK     | ❌        | Current active deployment location          |
| status          | ENUM                 | ✅        | 'online', 'offline', 'error', 'inactive'    |
| last_heartbeat  | TIMESTAMP            | ❌        | Automatically touched via polling event     |

---

## 6️⃣ Attendance Records Table

**Purpose:** Master log for temporal footprints. A centralized table mapping workers to devices.

| Field             | Type                 | Mandatory | Description                                      |
|-----------------  | -------------------- | --------- | ------------------------------------------------ |
| id                | VARCHAR(50) / PK     | ✅        | Format: `ATT-YYYYMMDD-0000`                      |
| device_id         | VARCHAR(50) / FK     | ✅        | Linked hardware                                  |
| worker_id         | VARCHAR(50) / FK     | ✅        | Recognized worker                                |
| site_id           | VARCHAR(50) / FK     | ✅        | Derived location string                          |
| time_in           | TIMESTAMP            | ✅        | Recorded entry timestamp                         |
| time_out          | TIMESTAMP            | ❌        | Recorded exit timestamp (if batched)             |
| status            | ENUM                 | ✅        | 'pending', 'submitted', 'failed'                 |

### Security & Scaling

This normalized approach guarantees horizontal scalability while maintaining strict relational constraints across Sites, Workers, and their corresponding hardware assignments. All tables incorporate implicit indexing logic to enhance BCA queries filtering by active device bounds and worker assignment trees.
