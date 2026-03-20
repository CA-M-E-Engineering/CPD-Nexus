# CPD-Nexus — System Architecture

This document describes the layered architecture of the **CPD-Nexus** backend and how its components interact.

---

## 1. Architectural Pattern: Hexagonal (Ports & Adapters)

The Go backend follows the **Hexagonal Architecture** (also known as Ports & Adapters or Clean Architecture). The goal is to keep the core business logic completely isolated from infrastructure details.

```
┌─────────────────────────────────────────────────────────┐
│                        Adapters (IN)                     │
│              HTTP Handlers  │  Bridge Message Handlers   │
└───────────────────┬─────────────────────────────────────┘
                    │  depends on
┌───────────────────▼─────────────────────────────────────┐
│                   Ports (Interfaces)                      │
│   WorkerService  │  AttendanceRepository  │  ...          │
└───────────────────┬─────────────────────────────────────┘
                    │  implemented by
┌───────────────────▼─────────────────────────────────────┐
│                   Core (Domain + Services)                │
│   domain/worker.go  │  services/pitstop_service.go  │ ... │
└───────────────────┬─────────────────────────────────────┘
                    │  depends on
┌───────────────────▼─────────────────────────────────────┐
│                   Adapters (OUT)                          │
│   repository/mysql/  │  external/sgbuildex/  │  bridge/  │
└─────────────────────────────────────────────────────────┘
```

### Layer Rules (strictly enforced)

| Layer | Package | Can import | Cannot import |
|---|---|---|---|
| **Domain** | `core/domain` | stdlib only | anything else |
| **Ports** | `core/ports` | `domain` | adapters, services |
| **Services** | `core/services` | `domain`, `ports`, adapters (needed for DI) | HTTP packages |
| **Handlers (IN)** | `api/handlers`, `bridge/handlers` | `ports` (interfaces) | `services` (concrete types) |
| **Repositories (OUT)** | `adapters/repository/mysql` | `domain`, `ports`, `database/sql` | `services`, handlers |
| **External (OUT)** | `adapters/external/sgbuildex` | `domain`, `ports`, stdlib | `database/sql`, handlers |

---

## 2. Package Responsibilities

### `internal/core/domain/`
Pure business entities. **No infrastructure imports** (`database/sql`, `net/http`, etc.).

| File | Contents |
|---|---|
| `worker.go` | `Worker` struct with sync status constants |
| `attendance.go` | `Attendance` struct for API responses |
| `sgbuildex.go` | `AttendanceRow` — join result for SGBuildex mapping (uses `*time.Time`, not `sql.NullTime`) |
| `project_site.go` | `Project` and `Site` structs |
| `settings.go` | `SystemSettings` — scheduler times, batch limits |
| `pitstop.go` | `PitstopAuthorisation` — cached Pitstop routing config |
| `constants.go` | Shared `SyncStatus*` integer constants |

### `internal/core/ports/`
Interface definitions that decouple layers. Every service and repository is accessed through a port.

| File | Interfaces Defined |
|---|---|
| `worker.go` | `WorkerRepository`, `WorkerService` |
| `project.go` | `ProjectRepository`, `ProjectService` |
| `attendance.go` | `AttendanceRepository`, `AttendanceService` |
| `pitstop_repo.go` | `PitstopRepository`, `PitstopService` |
| `submission.go` | `SubmissionRepository` |
| `settings.go` | `SettingsRepository` |
| `bridge_repo.go` | `BridgeRepository` |

### `internal/core/services/`
Business logic. Each service depends only on port interfaces.

| File | Responsibility |
|---|---|
| `worker_service.go` | Worker CRUD, validation, sync status transitions |
| `project_service.go` | Project CRUD with BCA field validation |
| `attendance_service.go` | Bridge attendance processing and ID generation |
| `pitstop_service.go` | Pitstop config sync, BCA submission, per-project test submission |
| `settings.go` | System settings management |
| `scheduler.go` | `DailyScheduler` — clock-based task runner, resets on settings change |

### `internal/adapters/repository/mysql/`
All database access. The only layer that uses `database/sql`.
Scans `sql.NullTime` → `*time.Time` before returning domain objects.

### `internal/adapters/external/sgbuildex/`
The SGTradeX Pitstop API adapter.

| File | Role |
|---|---|
| `client.go` | HTTP client; loads `SGTRADEX_API_KEY` at construction; 30s timeout |
| `config.go` | `FetchConfig()` — pulls routing config from `/api/v1/config` |
| `submitter.go` | Generic `SubmitPayloads[T Submittable]()` — batches payloads with size/count limits |
| `mappers.go` | `MapAttendanceToManpower()` — converts `domain.AttendanceRow` → `ManpowerUtilization` |
| `request.go` | `PushRequest`, `ParticipantWrapper`, `OnBehalfWrapper` structs |
| `utils.go` | `Ptr()`, `FormatOptionalTime()` — shared nilable helpers |
| `payloads/` | `ManpowerUtilization` struct matching BCA API schema |

### `internal/bridge/`
Host for the WebSocket gateway and the manager of persistent IoT bridge connections.

| File | Role |
|---|---|
| `manager.go` | `RequestManager` — maintains per-user transport map, dispatches commands synchronously |
| `transport.go` | Low-level WebSocket server transport (Upgrade & Heartbeat) |
| `types.go` | Message envelope structs (`BridgeMessage`, `BridgeMeta`) matching the protocol |
| `handlers/attendance.go` | Processes `GET_ATTENDANCE_RESPONSE` events from bridges |
| `handlers/user_sync.go` | Builds `REGISTER_USER` / `UPDATE_USER` command payloads |
| `handlers/user_sync_response.go` | Processes device acknowledgment; updates worker `is_synced` flag |

---

## 3. Request Flow Examples

### REST API Request (e.g. GET /api/workers)
```
HTTP Request
    → UserScopeMiddleware (extracts user_id from JWT cookie/header)
    → RequireUserScope (validates user identity is present)
    → WorkersHandler.GetWorkers()
    → ports.WorkerService.ListWorkers()     ← interface call
    → services.WorkerService.ListWorkers()  ← implementation
    → ports.WorkerRepository.List()         ← interface call
    → mysql.WorkerRepository.List()         ← SQL execution
    → []domain.Worker (returned up chain)
    → JSON response
```

### Scheduled BCA Submission
```
DailyScheduler fires at configured CPD_SUBMISSION_TIME
    → PitstopService.SubmitPendingAttendance()
    → AttendanceRepository.ExtractPendingAttendance()  [all non-submitted rows]
    → MapAttendanceToManpower(rows)                    [domain → payload]
    → SubmitPayloads(...)                              [batched HTTP POST]
        → POST /api/v1/data/push/manpower_utilization
        → SubmissionRepository.UpdateAttendanceStatus() [mark submitted/failed]
```

---

## 4. Scheduler Design

`DailyScheduler` is a reusable time-based task runner:
- Reads scheduled time (`HH:MM:SS`) from `SystemSettings` on every loop iteration.
- Waits until the next occurrence of that time.
- Exposes a `Reset()` channel — when settings are updated, the handler calls `Reset()` to immediately re-evaluate the next run time without waiting for the current sleep to expire.
- Two instances run in the system: `AttendanceSync` and `CPDSubmission`.

---

## 5. Multi-Tenant Isolation

All user-owned data (workers, projects, sites, devices, attendance) is scoped by `user_id`:
- HTTP layer: User context extracted from secure JWT by `UserScopeMiddleware`, enforced by `RequireUserScope`.
- Service layer: `userID` parameter passed through every operation and validated.
- Repository layer: every query includes `WHERE ... AND user_id = ?`.

Cross-tenant operations (e.g. assigning a project that belongs to a different user) are detected and rejected in the service layer with descriptive errors.

---

## 6. Validation Strategy

Field validation is applied at two points with identical rules:

| Concern | Backend | Frontend |
|---|---|---|
| NRIC / FIN format | `validation.ValidateNRICFIN()` | `validateNRICFIN()` in `validation.js` |
| Work pass type enum | `validation.ValidateWorkPassType()` | `validateWorkPassType()` |
| BCA trade code | `validation.ValidatePersonTrade()` | `validatePersonTrade()` |
| UEN format | `validation.ValidateUEN()` | `validateUEN()` |
| Project reference | `validation.ValidateProjectReferenceNumber()` | `validateProjectRef()` |
| HDB/LTA contract | `ValidateHDBContractNumber/LTA` | `validateHDBContract/LTA` |

The frontend validates on form submission; the backend re-validates in the service layer regardless of the calling source.

---

## 7. Maintenance Guide

| Task | Where |
|---|---|
| Add a new API endpoint | `api/router.go` + new handler method |
| Add a new service method | Declare in `ports/`, implement in `services/` |
| Update database schema | Add new `NNN_description.sql` file to `migrate/`, never modify existing |
| Update BCA field rules | `pkg/validation/sgbuildex_rules.go` AND `frontend-vue/src/utils/validation.js` |
| Change scheduler time | Update `SystemSettings` via `PUT /api/settings` |
| Update frontend styles | Global tokens in `frontend-vue/src/assets/styles/index.css` |
| Add new shared frontend constant | `frontend-vue/src/utils/constants.js` |
