import { http } from './http';

/**
 * API module for Analytics and Reporting.
 */
export const analyticsApi = {
    /**
     * Get dashboard summary stats
     */
    getDashboardStats: (params) => http.get('/analytics/dashboard', { params }),

    /**
     * Get activity log
     */
    getActivityLog: (params) => http.get('/analytics/activity-log', { params }),

    /**
     * Get detailed analytics for charts
     */
    getDetailedAnalytics: (params) => http.get('/analytics/detailed', { params }),
};
