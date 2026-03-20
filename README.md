# CPD-Nexus 🚀

**CPD-Nexus** (Construction Project & Data Nexus) is a unified, enterprise-grade platform for the construction industry. It handles project and workforce management, automates biometric attendance collection from IoT gateways via a WebSocket bridge, and submits Manpower Utilization records to BCA (Building and Construction Authority) through the SGTradeX Pitstop API.

---

## 🏗️ System Overview

```
┌──────────────┐      WebSocket       ┌───────────────────┐
│  IoT Bridge  │ ──────────────────▶  │  Backend Gateway  │
│  (Client)    │ ◀────────────────── │  (WS Server)      │
└──────────────┘                       └───────────┬───────┘
                                                    │
                                          Attendance Events
                                                    │
                                        ┌───────────▼───────────┐
  ┌────────────────┐   HTTP REST        │   Go Unified Backend   │
  │  Vue 3 Admin   │ ◀────────────────▶ │   (REST API + Logic)   │
  │   Dashboard    │                    └───────────┬────────────┘
  └────────────────┘                                │
                                                    │
                                         ┌──────────▼──────────┐
                                         │   MySQL Database    │
                                         └──────────┬──────────┘
                                                    │
                                    ┌───────────────▼──────────────┐
                                    │  SGTradeX Pitstop API (BCA)  │
                                    │  (Scheduled Daily Submission) │
                                    └──────────────────────────────┘
```

---

## ✨ Key Features

| Feature | Description |
|---|---|
| **Project Registry** | Manage construction sites, projects, and contractor metadata in one place |
| **Worker Directory** | Full lifecycle management of personnel with BCA-compliant field validation |
| **Biometric IoT Bridge** | Real-time bi-directional WebSocket connection to IoT device gateways |
| **Automated BCA Submission** | Scheduled daily submission of Manpower Utilization data to SGTradeX Pitstop |
| **CPD Submission Testing** | Manual per-project submission trigger for vendor testing and validation |
| **Analytics Dashboard** | Live operational metrics: attendance rates, sync status, device health |
| **Multi-Tenant Isolation** | All API operations are scoped to the authenticated user via secure JWT |
| **Input Validation** | BCA field rules enforced on both frontend and backend for all submissions |

---

## 🛠️ Project Structure

```
SGBuildex/
├── README.md                    # This file
├── ARCHITECTURE.md              # Layered architecture details
├── BRIDGE_COMMUNICATION.md      # IoT WebSocket protocol reference
│
├── backend/                     # Go Backend
│   ├── cmd/
│   │   └── server/main.go       # Application entry point & dependency wiring
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/        # HTTP request handlers (one per domain)
│   │   │   ├── middleware/      # Auth & scope middleware
│   │   │   └── router.go        # Route registration
│   │   ├── bridge/              # WebSocket bridge manager & message handlers
│   │   ├── core/
│   │   │   ├── domain/          # Pure domain models (no infrastructure deps)
│   │   │   ├── ports/           # Interface definitions (Repository & Service)
│   │   │   └── services/        # Business logic layer
│   │   ├── adapters/
│   │   │   ├── external/
│   │   │   │   └── sgbuildex/   # SGTradeX Pitstop API client & mapper
│   │   │   └── repository/
│   │   │       └── mysql/       # MySQL repository implementations
│   │   └── pkg/
│   │       ├── apperrors/       # Typed error helpers
│   │       ├── config/          # .env loader
│   │       ├── logger/          # Structured logging
│   │       ├── timeutil/        # Time formatting utilities
│   │       └── validation/      # BCA field validation rules
│   ├── migrate/                 # Ordered SQL migration scripts
│   ├── db.md                    # Database schema reference
│   └── .env                     # Backend environment config
│
└── frontend-vue/                # Vue.js 3 Frontend
    ├── src/
    │   ├── api/                 # Axios API bindings (one per domain)
    │   ├── components/          # Reusable UI components
    │   ├── stores/              # Pinia state stores
    │   ├── utils/
    │   │   ├── constants.js     # Shared enums (TRADES, PASS_TYPES)
    │   │   └── validation.js    # Frontend validation mirrors backend rules
    │   └── views/
    │       ├── client/          # Client-role pages (Workers, Projects, etc.)
    │       └── vendor/          # Vendor-role pages (Pitstop, Submissions)
    └── package.json
```

---

## 🚀 Getting Started

### Prerequisites
- [Go](https://golang.org/dl/) 1.21+
- [Node.js](https://nodejs.org/) 18+
- [MySQL](https://www.mysql.com/) 8.0+

### 1. Database Setup
Apply migration scripts in order:
```bash
mysql -u root -p your_db_name < backend/migrate/001_initial_schema.sql
# Apply subsequent scripts in sequence
```

### 2. Backend Setup
```bash
cd backend
cp .env.example .env        # Edit with your values
go mod download
go run cmd/server/main.go
```
Server starts on `http://localhost:3000`

### 3. Frontend Setup
```bash
cd frontend-vue
npm install
npm run dev
```
Dashboard available at `http://localhost:5173`

---

## ⚙️ Environment Configuration

Create `backend/.env`:

```env
# Server
API_PORT=3000

# Database
DB_USER=your_db_user
DB_PASS=your_db_password
DB_HOST=127.0.0.1:3306
DB_NAME=cpd_nexus

# SGTradeX / Pitstop (BCA Submission)
SGTRADEX_API_KEY=your_api_key_here
INGRESS_URL=https://ingress.pitstop.uat.dextech.ai
PITSTOP_URL=https://ca-me-sgbuildex.pitstop.uat.dextech.ai

# Authentication Security
JWT_SECRET=uC77N3FGObzfI3iHVundm0d+Ai9Y8T2Zl1LODr8lmpE=
DEFAULT_USER_PASSWORD=Nexus@2026!ChangeMe

# Scheduler (HH:MM:SS format, 24-hour)
ATTENDANCE_SYNC_TIME=01:00:00
CPD_SUBMISSION_TIME=02:00:00
```

---

## 📋 Core Workflows

### Attendance Collection (Bridge → Nexus)
1. The **DailyScheduler** triggers `RequestAttendance` at the configured time.
2. The **Bridge RequestManager** sends `GET_ATTENDANCE` commands via WebSocket to connected devices.
3. Device responses come back as `GET_ATTENDANCE_RESPONSE` events.
4. The **AttendanceHandler** parses the response and writes records to the `attendance` table with `status = 'pending'`.

### BCA Submission (Nexus → SGTradeX)
1. The **DailyScheduler** triggers `PitstopService.SubmitPendingAttendance()` at the configured time.
2. The service fetches all `attendance` rows where `status != 'submitted'`.
3. Rows are mapped to `ManpowerUtilization` payloads via `MapAttendanceToManpower()`.
4. Payloads are batched respecting `MaxWorkersPerRequest` and `MaxPayloadSizeKB` limits.
5. Each batch POSTs to `POST /api/v1/data/push/manpower_utilization` with the `SGTRADEX-API-KEY` header.
6. On success, `attendance.status` is updated to `'submitted'`.

### Worker Sync (Nexus → IoT Bridge)
1. Worker is created/updated with biometric data → `is_synced` set to `pending_registration` or `pending_update`.
2. Admin triggers **Sync** from the dashboard.
3. Backend dispatches commands to the **RequestManager**.
4. Commands are sent over the persistent WebSocket connection established by the bridge.
5. On `REGISTER_USER_RESPONSE` with HTTP 200, `is_synced` is set to `synced`.

---

## 🔒 Security & Compliance

- All scoped API routes require a valid JWT (passed via HttpOnly cookie or Authorization header) enforced by `RequireUserScope` middleware.
- FIN/NRIC data is validated against Singapore government NRIC/FIN format before storage.
- BCA field rules (UEN, trade codes, work pass types, submission months) are enforced on both frontend input and backend service layers.
- The `SGTRADEX_API_KEY` is never exposed to the frontend — all external API calls are server-side.
- Multi-tenant isolation: all database queries are scoped to the requesting user's `user_id`.

---

## 📄 License
© 2026 CA-M&E Engineering. All rights reserved.
