import { http } from './http';

/**
 * API module for Company-related operations.
 */
export const companiesApi = {
    /**
     * Fetch all companies
     */
    getCompanies: (params) => http.get('/companies', { params }),

    /**
     * Fetch a single company by ID
     */
    getCompanyById: (id) => http.get(`/companies/${id}`),

    /**
     * Create a new company
     */
    createCompany: (data) => http.post('/companies', data),

    /**
     * Update an existing company
     */
    updateCompany: (id, data) => http.put(`/companies/${id}`, data),

    /**
     * Delete a company
     */
    deleteCompany: (id) => http.delete(`/companies/${id}`),
};
