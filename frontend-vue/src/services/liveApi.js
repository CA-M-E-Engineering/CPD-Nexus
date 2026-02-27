import { authApi } from '../api/auth.api';
import { workersApi } from '../api/workers.api';
import { projectsApi } from '../api/projects.api';
import { sitesApi } from '../api/sites.api';
import { devicesApi } from '../api/devices.api';
import { usersApi } from '../api/users.api';
import { assignmentsApi } from '../api/assignments.api';
import { analyticsApi } from '../api/analytics.api';
import { attendanceApi } from '../api/attendance.api';
import { settingsApi } from '../api/settings.api';
import { bridgeApi } from '../api/bridge.api';

import { http } from '../api/http';

export const liveApi = {
    // --- Auth & User ---
    login: authApi.login,
    getUserProfile: authApi.getUserProfile,
    logout: authApi.logout,


    // --- Workers ---
    getWorkers: workersApi.getWorkers,
    getWorkerById: workersApi.getWorkerById,
    createWorker: workersApi.createWorker,
    updateWorker: workersApi.updateWorker,
    deleteWorker: workersApi.deleteWorker,

    // --- Projects ---
    getProjects: projectsApi.getProjects,
    getProjectById: projectsApi.getProjectById,
    createProject: projectsApi.createProject,
    updateProject: projectsApi.updateProject,
    deleteProject: projectsApi.deleteProject,

    // --- Sites ---
    getSites: sitesApi.getSites,
    getSiteById: sitesApi.getSiteById,
    createSite: sitesApi.createSite,
    updateSite: sitesApi.updateSite,
    deleteSite: sitesApi.deleteSite,

    // --- Devices ---
    getDevices: devicesApi.getDevices,
    getDeviceById: devicesApi.getDeviceById,
    createDevice: devicesApi.createDevice,
    updateDevice: devicesApi.updateDevice,
    deleteDevice: devicesApi.deleteDevice,
    bulkAssign: devicesApi.bulkAssign,

    // --- Users ---
    getUsers: usersApi.getUsers,
    getUserById: usersApi.getUserById,
    createUser: usersApi.createUser,
    updateUser: usersApi.updateUser,
    deleteUser: usersApi.deleteUser,
    updateBridgeConfig: usersApi.updateBridgeConfig,

    // --- Assignments (Relations) ---
    assignWorkersToProject: assignmentsApi.assignWorkersToProject,
    assignDevicesToSite: assignmentsApi.assignDevicesToSite,
    assignProjectToSite: assignmentsApi.assignProjectToSite,

    // --- Analytics & Dashboard ---
    getDashboardStats: analyticsApi.getDashboardStats,
    getActivityLog: analyticsApi.getActivityLog,
    getDetailedAnalytics: analyticsApi.getDetailedAnalytics,
    getAttendance: attendanceApi.getAttendance,

    // --- System Settings ---
    getSettings: settingsApi.getSettings,
    updateSettings: settingsApi.updateSettings,

    // --- Bridge ---
    syncUsers: bridgeApi.syncUsers,

    // --- Utilities ---
    async simulateExport(label) {
        // For live API, this might trigger a download or job
        return http.post('/export', { type: label });
    }
};
