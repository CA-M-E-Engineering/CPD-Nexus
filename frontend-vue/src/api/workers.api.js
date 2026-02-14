import { http } from './http';

/**
 * API module for Worker-related operations.
 */
export const workersApi = {
    /**
     * Fetch all workers
     */
    getWorkers: (params) => http.get('/workers', { params }),

    /**
     * Fetch a single worker by ID
     */
    getWorkerById: (id) => http.get(`/workers/${id}`),

    /**
     * Create a new worker
     */
    createWorker: (data) => http.post('/workers', data),

    /**
     * Update an existing worker
     */
    updateWorker: (id, data) => http.put(`/workers/${String(id)}`, data),

    /**
     * Delete a worker
     */
    deleteWorker: (id) => http.delete(`/workers/${id}`),
};
