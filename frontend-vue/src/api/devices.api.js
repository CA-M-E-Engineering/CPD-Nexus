import { http } from './http';

/**
 * API module for Device-related operations.
 */
export const devicesApi = {
    /**
     * Fetch all devices
     */
    getDevices: (params) => http.get('/devices', { params }),

    /**
     * Fetch a single device by ID
     */
    getDeviceById: (id) => http.get(`/devices/${id}`),

    /**
     * Create a new device
     */
    createDevice: (data) => http.post('/devices', data),

    /**
     * Update an existing device
     */
    updateDevice: (id, data) => http.put(`/devices/${id}`, data),

    /**
     * Delete a device
     */
    deleteDevice: (id) => http.delete(`/devices/${id}`),

    /**
     * Bulk assign devices to a User
     */
    bulkAssign: (userId, deviceIds) => http.post(`/users/${userId}/devices/bulk`, { deviceIds }),
};
