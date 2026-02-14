# CPD-Nexus 🚀

**CPD-Nexus** (Construction Project & Data Nexus) is a unified, high-performance platform designed for the construction industry. it streamlines project management, automates worker attendance tracking through IoT device integration, and ensures seamless compliance with BCA (Building and Construction Authority) reporting requirements.

---

## 🏗️ Core Architecture

The project is built with a decoupled architecture for maximum scalability and performance:

*   **Frontend**: A premium, glassmorphism-inspired dashboard built with **Vue.js 3** and **Vite**.
*   **Backend**: A robust, high-currency unified server written in **Go (Golang)** featuring:
    *   **REST API**: Serving the management dashboard.
    *   **Bridge Connector**: WebSocket-based real-time communication with biometric IoT devices.
    *   **Submission Worker**: Automated background processor for BCA compliance reporting.
*   **Database**: MySQL/SQLite (configurable) with an optimized schema for rapid project-site-worker lookups.

---

## ✨ Key Features

*   **Unified Project Registry**: Manage construction sites, projects, and contractor assignments in one place.
*   **Biometric IoT Integration**: Automated attendance fetching from remote biometric devices via the Bridge module.
*   **BCA Compliance (CPD)**: Automated daily submission of Manpower Utilization (MU) and Distribution (MD) data.
*   **Real-time Analytics**: Insights into worker density across different sites.
*   **Modern UI**: High-end UX with support for dark mode and fluid animations.

---

## 🛠️ Project Structure

```bash
SGBuildex/
├── backend/            # Go Backend (API, Bridge, Workers)
│   ├── cmd/            # Entry points and tools
│   ├── internal/       # Core business logic (Adapters, Domain, Ports, Services)
│   └── migrate/        # SQL Migration scripts
├── frontend-vue/       # Vue.js 3 Frontend application
│   └── src/            # Components, Views, and Services
└── .env                # Shared environment configuration
```

---

## 🚀 Getting Started

### 1. Prerequisites
*   [Go](https://golang.org/dl/) (1.21+)
*   [Node.js](https://nodejs.org/) (18+)
*   [MySQL](https://www.mysql.com/) (For production-ready storage)

### 2. Backend Setup
```bash
cd backend
go mod download
go run main.go
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
Configure your `.env` file in the root directory:

```env
API_PORT=3000
FRONTEND_PORT=5173
DB_USER=root
DB_PASS=your_password
DB_HOST=127.0.0.1:3306
DB_NAME=bas_mvp
```

---

## 🔒 Security & Compliance
CPD-Nexus handles sensitive FIN/NRIC data using encryption at rest and follows the SGBuildex specification for secure API transmission.

---

## 📄 License
© 2024 CA-M-E Engineering. All rights reserved.
