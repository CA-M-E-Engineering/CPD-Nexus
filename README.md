# CPD-Nexus 🚀

**CPD-Nexus** (Construction Project & Data Nexus) is a unified, high-performance platform designed for the construction industry. It streamlines project management, automates worker attendance tracking through IoT device integration using a WebSocket bridge, and ensures seamless compliance with BCA (Building and Construction Authority) reporting requirements.

---

## 🏗️ Core Architecture

The project is built with a decoupled architecture for maximum scalability and performance:

### Frontend
A premium, glassmorphism-inspired dashboard built with **Vue.js 3** and **Vite**.
* **Role-Based Access**: Specialized interfaces and routing for `Manager` (administrative control) and `Client` (operational overview).
* **Responsive UI Design**: High-end UX with modern typography, smooth animations, and data grids.
* **Key Modules**: Device Registry, Project Management, Site Allocation, Worker Directory, and Attendance tracking.

### Backend
A robust, high-currency unified server written in **Go (Golang)** featuring:
* **REST API**: Serving the vue.js management dashboard via standard CRUD endpoints.
* **Bridge Connector (WebSocket)**: A dedicated `RequestManager` that manages real-time bi-directional communication with biometric IoT hardware. It actively issues `FETCH_ATTENDANCE` commands and processes the inbound attendance event streams securely.
* **Submission Worker**: Automated background processor scheduled to run daily to submit Manpower Utilization (MU) and Distribution (MD) records to BCA.
* **Database**: MySQL schema optimized for rapid project-site-worker-device mapping and lookups.

---

## ✨ Key Features

* **Unified Project Registry**: Manage construction sites, projects, and workforce profiles in one place.
* **Biometric IoT Integration**: Automated attendance fetching from remote device gateways via the Bridge module.
* **Real-Time Device Allocation**: Map, unassign, and redeploy IoT devices seamlessly across multiple construction sites.
* **BCA Compliance (CPD)**: Automated daily submission of worker attendance data tailored to government API standards.
* **Dynamic Trade Categorization**: Supports detailed BCA-compliant designated trade mapping for both local and foreign workers.
* **User-Scoped Manual Sync**: Strict bi-directional data flow ensuring users only synchronize workforce data belonging to their own organization.
* **Granular Sync Intelligence**: Real-time validation layer that identifies and reports missing biometrics or unassigned site hardware before attempting IoT deployment.

---

## ⚙️ Core Workflows

### 🔄 Manual Synchronization Protocol
Synchronization is triggered manually via the UI to ensure administrative control. The backend implements a granular validation pipeline:
1. **Validation**: Checks for biometric availability (Face/Card) and active site hardware.
2. **Categorization**: Reports workers missing biometrics or devices without blocking the entire sync batch.
3. **Deployment**: Issues `REGISTER_USER` or `UPDATE_USER` commands to all online devices at the worker's assigned site.
4. **Security**: All sync requests require a valid `X-User-ID` header to maintain multi-tenant isolation.

### 👥 Worker Lifecycle Management
The system enforces strict data integrity rules for personnel:
- **Deactivation**: Setting a worker to `inactive` automatically wipes their `current_project_id`.
- **Global Visibility**: Inactive workers are automatically omitted from the Worker Management list, site personnel Overviews, and all Analytics dashboards.
- **Project Assignment**: Only workers with an `active` status can be redeployed to construction projects.

---

## 🛠️ Project Structure

```bash
SGBuildex/
├── backend/            # Go Backend (API, Bridge, Workers)
│   ├── cmd/            # Entry points and tools
│   │   └── server/     # Main Application Entrypoint
│   │       └── main.go
│   ├── internal/       # Core business logic (Adapters, Domain, Ports, Services)
│   │   ├── bridge/     # WebSocket connection logic and Attendance Handlers
│   │   ├── api/        # REST API Routes and Controllers
│   │   └── core/       # Business Domains and Services
│   ├── migrate/        # SQL Migration scripts
│   └── .env            # Backend environment configuration
├── frontend-vue/       # Vue.js 3 Frontend application
│   ├── src/            # Core source files
│   │   ├── api/        # Axios API configurations
│   │   ├── components/ # Reusable UI components (Modals, Badges, Tables)
│   │   └── views/      # Page-level components (Dashboards, Resource Lists)
│   └── package.json    # Node dependencies and build scripts
```

---

## 🚀 Getting Started

### 1. Prerequisites
* [Go](https://golang.org/dl/) (1.21+)
* [Node.js](https://nodejs.org/) (18+)
* [MySQL](https://www.mysql.com/) (For production-ready storage)

### 2. Backend Setup
```bash
cd backend
go mod download
go run cmd/server/main.go
```
*The server will start on `http://localhost:3000`*

### 3. Frontend Setup
```bash
cd frontend-vue
npm install
npm run dev
```
*The dashboard will be available at `http://localhost:5173`*

---

## ⚙️ Configuration
Configure your `.env` file in the `backend/` directory:

```env
API_PORT=3000
FRONTEND_PORT=5173
DB_USER=your_db_username
DB_PASS=your_db_password
DB_HOST=127.0.0.1:3306
DB_NAME=your_db_name

# SGTRADEX / Pitstop Configuration
SGTRADEX_API_KEY=your_sgtradex_api_key_here
```

---

## 🔒 Security & Compliance
CPD-Nexus handles sensitive FIN/NRIC data securely and strictly follows the necessary compliance outlines for API transmission across the BCA and external IoT endpoints.

---

## 📄 License
© 2026 CA-M&E Engineering. All rights reserved.
