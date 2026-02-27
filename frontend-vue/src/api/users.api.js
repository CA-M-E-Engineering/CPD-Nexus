import { http } from './http';

/**
 * API module for User-related operations.
 */
export const usersApi = {
    /**
     * Fetch all users
     */
    getUsers: () => http.get('/users'),

    /**
     * Fetch a single user by ID
     */
    getUserById: (id) => http.get(`/users/${id}`),

    /**
     * Create a new user
     */
    createUser: (data) => http.post('/users', data),

    /**
     * Update an existing user
     */
    updateUser: (id, data) => http.put(`/users/${id}`, data),

    /**
     * Delete a user
     */
    deleteUser: (id) => http.delete(`/users/${id}`),

    /**
     * Update a user's bridge connection configuration
     */
    updateBridgeConfig: (id, data) => http.put(`/users/${id}/bridge`, data),
};
