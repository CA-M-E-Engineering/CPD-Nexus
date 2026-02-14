import { http } from './http';

/**
 * API module for managing relationships and assignments.
 */
export const assignmentsApi = {
    /**
     * Assign workers to a project
     */
    assignWorkersToProject: (projectId, workerIds) =>
        http.post(`/projects/${projectId}/assign-workers`, { workerIds }),

    /**
     * Assign devices to a site
     */
    assignDevicesToSite: (siteId, deviceIds) =>
        http.post(`/sites/${siteId}/assign-devices`, { deviceIds }),

    /**
     * Assign projects to a site
     */
    assignProjectToSite: (siteId, projectIds) =>
        http.post(`/sites/${siteId}/assign-projects`, { projectIds }),
};
