import { http } from './http';

/**
 * API module for Attendance-related operations.
 */
export const attendanceApi = {
    /**
     * Fetch attendance records with filtering
     * @param {Object} params - { site_id, worker_id, date }
     */
    getAttendance: (params) => http.get('/attendance', { params }),

    /**
     * Update an attendance record
     */
    updateAttendance: (id, data) => http.put(`/attendance/${id}`, data),

    /**
     * Fetch a single attendance record by ID
     */
    getAttendanceById: (id) => http.get(`/attendance/${id}`),
};
