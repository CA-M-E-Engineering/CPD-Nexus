import { http } from './http';

/**
 * API module for Construction Site operations.
 */
export const sitesApi = {
    /**
     * Fetch all sites
     */
    getSites: (params) => {
        console.log('[SitesAPI] getSites called with:', params);
        return http.get('/sites', { params });
    },

    /**
     * Fetch a single site by ID
     */
    getSiteById: (id) => http.get(`/sites/${id}`),

    /**
     * Create a new site
     */
    createSite: (data) => http.post('/sites', data),

    /**
     * Update an existing site
     */
    updateSite: (id, data) => http.put(`/sites/${id}`, data),

    /**
     * Delete a site
     */
    deleteSite: (id) => http.delete(`/sites/${id}`),
};
