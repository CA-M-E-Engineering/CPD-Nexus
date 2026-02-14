import { http } from './http';

/**
 * API module for Project-related operations.
 */
export const projectsApi = {
    /**
     * Fetch all projects
     */
    getProjects: (params) => {
        const query = new URLSearchParams(params).toString();
        return http.get(`/projects?${query}`);
    },

    /**
     * Fetch a single project by ID
     */
    getProjectById: (id) => http.get(`/projects/${id}`),

    /**
     * Create a new project
     */
    createProject: (data) => http.post('/projects', data),

    /**
     * Update an existing project
     */
    updateProject: (id, data) => http.put(`/projects/${id}`, data),

    /**
     * Delete a project
     */
    deleteProject: (id) => http.delete(`/projects/${id}`),
};
