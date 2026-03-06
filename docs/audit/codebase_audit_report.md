# Codebase Audit Report: SGBuildex / CPD Nexus

## 1. Project Overview
**System**: CPD Nexus (Construction Project Management & BCA Compliance)
**Backend**: Go (Hexagonal Architecture)
**Frontend**: Vue.js (Vite)
**Database**: MySQL

This audit covers security vulnerabilities, architectural consistency, code quality, and performance risks across the full repository.

---

## 2. Overall Evaluation Table

| Category | Rating | Summary |
| :--- | :--- | :--- |
| **Security** | 🔴 Critical | **Severe Broken Access Control** on administrative routes and lack of CSRF protection. |
| **Architecture** | 🟡 Medium | Generally follows Hexagonal patterns, but has leakage in Pitstop and Upload handling. |
| **Code Quality** | 🟢 Good | Parameterized SQL and clean separation of concerns in most areas. |
| **Maintainability** | 🟢 Good | Clear structure, though some DTO duplication exists. |
| **Performance** | 🟢 Good | Proper DB pooling and small payload handling. |

---

## 3. Critical Issues

### [Security] Broken Access Control (Bypassable Admin Routes)
- **File**: `internal/api/router.go:45-50`
- **Code Snippet**:
  ```go
  api.HandleFunc("/users", cfg.UsersHandler.GetUsers).Methods("GET")
  api.HandleFunc("/users", cfg.UsersHandler.CreateUser).Methods("POST")
  ```
- **Explanation**: These routes are registered in the `api` subrouter which only uses `UserScopeMiddleware`. This middleware extracts user info but **does not reject** unauthenticated or non-admin requests. Only the `scoped` subrouter uses `RequireUserScope`.
- **Why it is dangerous**: Any user with a valid JWT (even a low-privilege client) can list, create, update, or delete users, effectively granting themselves admin rights or accessing sensitive vendor data.
- **Recommended Fix**: Move administrative routes into a dedicated `adminOnly` subrouter with a middleware that enforces `user_type == 'admin'`.

### [Security] Context Spoofing in Pitstop Handler
- **File**: `internal/api/handlers/pitstop_handler.go:38-41`
- **Code Snippet**:
  ```go
  userID := r.Header.Get("X-User-ID")
  if userID == "" {
      userID = "Owner_001"
  }
  ```
- **Explanation**: The handler trusts the `X-User-ID` header provided by the client instead of using the authenticated `userID` from the request context.
- **Why it is dangerous**: An attacker can send any `userID` in the header to synchronize or test data "on behalf of" any other user in the system.
- **Recommended Fix**: Use `ports.GetUserID(r.Context())` and remove dependency on client-provided headers for identity.

---

## 4. High Priority Issues

### [Security] Missing CSRF Protection (Cookie-based Auth)
- **File**: `frontend-vue/src/api/http.js:53`
- **Code Snippet**: `credentials: 'include'`
- **Explanation**: The application uses HttpOnly cookies for authentication (`auth_token`). However, there is no Anti-CSRF token mechanism implemented.
- **Why it is dangerous**: A malicious site can trick a logged-in user into making state-changing requests (e.g., deleting a project) to the SGBuildex API.
- **Recommended Fix**: Implement Double Submit Cookie pattern or a dedicated CSRF middleware in the Go backend.

### [Security] Mass Assignment / IDOR in Worker Update
- **File**: `internal/api/handlers/workers.go:82`
- **Code Snippet**: 
  ```go
  if err := h.service.UpdateWorker(r.Context(), userID, id, &req); err != nil {
  ```
- **Explanation**: The `UpdateWorkerRequest` DTO contains fields like `UserID` and `IsSynced`. If the service layer doesn't explicitly ignore these fields, a user can "stolen" a worker by changing its `user_id` or artificially mark it as synced.
- **Why it is dangerous**: Cross-tenant data leakage and corruption of sync state with external BCA systems.
- **Recommended Fix**: Ensure the service layer explicitly ignores sensitive fields in the DTO or uses a separate "Safe" update DTO.

---

## 5. Medium Issues

### [Architecture] Implementation Leakage in Router
- **File**: `internal/api/router.go:124`
- **Code Snippet**: 
  ```go
  r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
  ```
- **Explanation**: The router directly interacts with the filesystem using a relative path. In Hexagonal architecture, file storage should be an `adapter`.
- **Why it is dangerous**: Makes the system harder to test and difficult to migrate to cloud storage (e.g., S3) in the future.
- **Recommended Fix**: Create a `FileServer` adapter and pull the storage path from configuration.

### [Maintainability] Hardcoded Fallbacks in Security Logic
- **File**: `internal/api/handlers/pitstop_handler.go:40`
- **Code Snippet**: `userID = "Owner_001"`
- **Explanation**: Hardcoded strings used for authorization fallbacks make the system brittle and difficult to audit.
- **Recommended Fix**: Always fail-closed. If no ID is found, return `401 Unauthorized`.

---

## 6. Low Priority Improvements

### [API Design] Lack of Rate Limiting
- **Location**: `main.go`
- **Observation**: There is no global or per-IP rate limiting on the REST API.
- **Risk**: Susceptibility to brute-force attacks or DoS.
- **Fix**: Integrate a simple rate-limiting middleware (e.g., `didip/tollbooth`).

### [Code Quality] Go Context Usage
- **Location**: Various Handlers
- **Observation**: Context is passed correctly, but some heavy operations don't check `ctx.Done()` during loop execution.

---

## 7. Architecture Analysis

### Strengths
- **Clear Separation**: Core domain is isolated from adapters.
- **Ports Pattern**: Use of interfaces in `internal/core/ports` allows for easy mocking in tests.
- **Dependency Direction**: Mostly correct; internal logic does not depend on HTTP/SQL implementation details.

### Weaknesses
- **Administrative Boundary**: The "Global Admin" vs "Tenant" boundary is blurred in the `router.go` and `UsersHandler`.
- **Infrastructure Leak**: File upload handling and auth secret initialization are slightly tightly coupled to the `api` package.

---

## 8. Suggested Fix Roadmap

1. **Phase 1 (Immediate)**: Fix Broken Access Control on `/users` and `/devices` routes.
2. **Phase 2 (Immediate)**: Remove `X-User-ID` header reliance in `PitstopHandler`.
3. **Phase 3 (Security)**: Implement CSRF protection in both Frontend and Backend.
4. **Phase 4 (Refactor)**: Moving administrative routes to a separate subrouter with strict role-based access control (RBAC).

---

## 9. File Statistics

| Layer | File Count | Major Component |
| :--- | :---: | :--- |
| **Adapters** | 22 | MySQL Repositories, External Clients |
| **API** | 16 | Handlers, Middleware, Router |
| **Core** | 35 | Services, Domain Entities, Ports |
| **Frontend** | 80 | Vue Components, API Bindings |

**End of Audit Report**
