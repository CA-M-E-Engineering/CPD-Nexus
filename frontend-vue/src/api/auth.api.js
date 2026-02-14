import { http } from './http';

/**
 * API module for Authentication and User Profile.
 */
export const authApi = {
    /**
     * Login with username and password
     */
    login: (username, password) => http.post('/auth/login', { username, password }),

    /**
     * Get current user profile
     */
    getUserProfile: () => http.get('/auth/me'),

    /**
     * Logout
     */
    logout: () => http.post('/auth/logout'),
};
