# Comprehensive Codebase Audit Report: SGBuildex Platform

## 1. Project Overview

**Application:** CPD Nexus (SGBuildex)
**Architecture Style:** Hexagonal Architecture (Ports & Adapters)
**Tech Stack:** Go (Backend/REST), Vue.js (Frontend), MySQL
**Objective:** Formal code audit for production readiness, specifically assuming handling of sensitive government compliance data.

The project structure is mostly sound, with explicit separation of domains, ports, services, adapters, and handlers. However, several critical architectural violations and high-risk security practices were identified during this audit.

---

## 2. Overall Evaluation Table

| Category | Status | Primary Concern |
| :--- | :--- | :--- |
| **Security - Authentication** | 🔴 Needs Fix | JWT tokens stored insecurely in Vue `localStorage`. |
| **Security - SQL Injection** | 🟢 Good | Widespread use of parameterized queries; low risk. |
| **Security - File Uploads** | 🟢 Good | Strong validation (size, MIME type magic bytes, path sanitization). |
| **Architecture Boundaries** | 🔴 Needs Fix | Severe layering violations across Repositories and Services. |
| **Concurrency / Resources** | 🟡 Acceptable | Database connection pooling configured; some minor timer leaks in schedulers. |
| **Frontend Security** | 🟡 Acceptable | No `v-html` XSS vectors found, but token storage is a risk. |

---

## 3. Critical Issues

### [CRITICAL] 1. Authentication Token Storage vulnerability (XSS to Account Takeover)
- **File:** `frontend-vue/src/views/Login.vue` (and multiple other frontend files)
- **Line Number:** `29`
- **Code Snippet:** 
  ```javascript
  localStorage.setItem('auth_token', response.token);
  ```
- **Explanation:** The application stores the sensitive `auth_token` directly in the browser's `localStorage`.
- **Why it is dangerous:** `localStorage` is accessible by any JavaScript running on the page. In the event of an XSS attack (even a minor one through a dependency), attackers can easily steal the tokens and completely hijack user sessions containing sensitive government compliance data.
- **Recommended Fix:** The backend should issue the JWT as an `HttpOnly`, `Secure`, `SameSite=Strict` cookie during login. The frontend should rely on the browser to send this cookie with API requests automatically.

---

## 4. High Priority Issues

### [HIGH] 1. Architecture Layering Violation: Repositories depending on API Middleware
- **Files:** 
  - `backend/internal/adapters/repository/mysql/site_repo.go` (Line 7)
  - `backend/internal/adapters/repository/mysql/project_repo.go` (Line 7)
  - `backend/internal/adapters/repository/mysql/device_repository.go` (Line 7)
- **Code Snippet:**
  ```go
  import "sgbuildex/internal/api/middleware"
  // ...
  if middleware.IsVendor(ctx) { ... }
  ```
- **Explanation:** The Adapter layer (MySQL Repositories) directly imports and uses the HTTP middleware layer to check the context for `IsVendor`.
- **Why it is dangerous:** This completely breaks the Hexagonal Architecture pattern. Adapters must only depend on `core/domain` and `core/ports`. By bridging the adapter directly to the API transport layer, the repository becomes tightly coupled to HTTP semantics, making it impossible to reuse these repositories across other transports (e.g., CLI, gRPC, Message Queues).
- **Recommended Fix:** Define custom context keys inside the `core/ports` or `core/domain` layer. The middleware should set the boolean using the internal key, and the repository should extract it using that internal key (or better yet, pass `isVendor` as an explicit parameter from the service layer, as done in `WorkerRepository`).

### [HIGH] 2. Architecture Layering Violation: Services depending on API Layer
- **File:** `backend/internal/core/services/worker_service.go`
- **Line Number:** 6
- **Code Snippet:**
  ```go
  import "sgbuildex/internal/api/middleware"
  ```
- **Explanation:** The Application Service layer directly imports the HTTP middleware layer.
- **Why it is dangerous:** The Domain/Service layers must have **zero** dependencies on the outer adapters or API layers. This coupling prevents independent testing and breaks the isolation guarantee of Ports & Adapters.
- **Recommended Fix:** Shift context keys to a shared `pkg/contextutils` or `core/ports` package.

### [HIGH] 3. Architecture Layering Violation: Services depending directly on External Adapters
- **File:** `backend/internal/core/services/pitstop_service.go`
- **Line Number:** 6
- **Code Snippet:**
  ```go
  import "sgbuildex/internal/adapters/external/sgbuildex"
  ```
- **Explanation:** The `pitstop_service` directly imports a concrete external API adapter, bypassing the defined port.
- **Why it is dangerous:** The core business logic should only depend on interfaces (Ports). By depending directly on a concrete adapter, you lose the ability to mock this external service for testing, and swapping to a different API client will require modifying core business logic.
- **Recommended Fix:** Create an interface in `internal/core/ports/sgbuildex_port.go` and inject it into the `pitstop_service` constructor.

---

## 5. Medium Issues

### [MEDIUM] 1. Minor Resource Leak in Schedulers
- **File:** `backend/internal/core/services/scheduler.go`
- **Code Snippet:**
  ```go
  case <-time.After(durationUntilNext):
      // ...
  ```
- **Explanation:** Go's `time.After` leaks the underlying `time.Timer` until it fires. If `s.reset` is repeatedly hit, multiple timers will accumulate in memory until their durations expire.
- **Why it is dangerous:** For deeply nested or rapidly changing schedulers, this can cause a slow memory leak of timer objects.
- **Recommended Fix:** Use `time.NewTimer(duration)` and call `timer.Stop()` when resetting, ensuring the channel is safely drained.

---

## 6. Low Priority Improvements

### [LOW] 1. Hardcoded Allowed Origins in `AllowedOrigins` configuration string
- **File:** `backend/internal/pkg/config/config.go`
- **Explanation:** The fallback string for `ALLOWED_ORIGINS` has hardcoded localhosts. Since `getEnv` allows fallbacks, an incorrectly configured production environment might inadvertently allow local origins.

### [LOW] 2. Implicit User Type Fallback
- **File:** `backend/internal/adapters/repository/mysql/worker_repo.go`
- **Explanation:** `userType := w.UserType; if userType == "" { userType = "user" }` fallback exists inside the database repository. Default values should be handled exclusively in the Domain layer or the Service layer, not deep within SQL update strings.

---

## 7. Architecture Analysis (Strengths and Weaknesses)

**Strengths:**
- High baseline security in critical paths: SQL queries are thoroughly parameterized using `database/sql` driver mechanics (`?` execution), entirely preventing common SQL injections.
- Strong implementation of API routing and structured context propagation.
- Secure by default on file uploads: using UUID for file names to prevent path traversal, and strict MIME type detection using magic bytes, not just extensions.

**Weaknesses:**
- The team has misunderstood the dependency rule of Hexagonal Architecture context propagation. `worker_service.go`, `site_repo.go`, `project_repo.go`, and `device_repository.go` have severe inverted dependency relations by statically compiling the HTTP API package inside the Core and Adapter layers.

---

## 8. Suggested Fix Roadmap

1. **Phase 1 (Critical Security):** Rewrite the Auth flow to use `HttpOnly` Secure cookies instead of `localStorage`. Update the Vue interceptors to handle `credentials: true`.
2. **Phase 2 (Architectural Rescue):** create `ports.ContextKeys` or a dedicated package for Context values. Delete ALL imports of `sgbuildex/internal/api/middleware` from `services/` and `adapters/repository/`. Inject parameters functionally where it makes sense.
3. **Phase 3 (Core Isolation):** Refactor `pitstop_service.go` to depend on an interface instead of the concrete `sgbuildex` adapter package.
4. **Phase 4 (Refinement):** Refactor the `scheduler.go` to use reusable `time.NewTimer` to prevent goroutine memory bloat.

---

## 9. File Statistics

- **Total Backend Files Reviewed:** \~150 files (`internal/`, `cmd/`, `pkg/`)
- **Total Frontend Files Reviewed:** \~100 files (`src/views/`, `src/components/`, `src/api/`)
- **Security Hotspots Identified:** Frontend Authentication Storage, API-to-DB Boundary leakage.
