import { http } from './http';

/**
 * API module for Tenant-related operations.
 */
export const tenantsApi = {
    /**
     * Fetch all tenants
     */
    getTenants: () => http.get('/tenants'),

    /**
     * Fetch a single tenant by ID
     */
    getTenantById: (id) => http.get(`/tenants/${id}`),

    /**
     * Create a new tenant
     */
    createTenant: (data) => http.post('/tenants', data),

    /**
     * Update an existing tenant
     */
    updateTenant: (id, data) => http.put(`/tenants/${id}`, data),

    /**
     * Delete a tenant
     */
    deleteTenant: (id) => http.delete(`/tenants/${id}`),
};
