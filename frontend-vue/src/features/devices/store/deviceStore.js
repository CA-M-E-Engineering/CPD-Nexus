import { defineStore } from 'pinia';
import { api } from '@/services/api';

export const useDeviceStore = defineStore('devices', {
    state: () => ({
        devices: [],
        loading: false,
        error: null,
    }),
    getters: {
        activeDevices: (state) => state.devices.filter(d => d.status !== 'inactive'),
        stats: (state) => {
            return {
                total: state.devices.length,
                online: state.devices.filter(d => d.status === 'online').length,
                offline: state.devices.filter(d => d.status === 'offline').length
            }
        }
    },
    actions: {
        async fetchDevices(tenantId = null) {
            this.loading = true;
            this.error = null;
            try {
                const params = tenantId ? { tenant_id: tenantId } : {};
                // Assuming api.getDevices is compatible or we import the specific api
                const response = await api.getDevices(params);
                this.devices = response || [];
            } catch (err) {
                this.error = err.message || 'Failed to fetch devices';
                console.error(err);
            } finally {
                this.loading = false;
            }
        },
        async registerDevice(payload) {
            try {
                const newDevice = await api.createDevice(payload);
                this.devices.push(newDevice);
                return newDevice;
            } catch (err) {
                throw err;
            }
        },
        async deleteDevice(id) {
            try {
                await api.deleteDevice(id);
                this.devices = this.devices.filter(d => d.device_id !== id);
            } catch (err) {
                throw err;
            }
        },
        async getDeviceById(id) {
            // Check cache first? For now always fetch fresh to be safe
            try {
                return await api.getDeviceById(id);
            } catch (err) {
                throw err;
            }
        },
        async updateDevice(id, payload) {
            try {
                await api.updateDevice(id, payload);
                // Update local state if exists
                const index = this.devices.findIndex(d => d.device_id === id);
                if (index !== -1) {
                    this.devices[index] = { ...this.devices[index], ...payload };
                }
            } catch (err) {
                throw err;
            }
        },
        async assignDeviceToTenant(tenantId, deviceIds) {
            try {
                await api.bulkAssign(tenantId, deviceIds);
                // Refresh devices or update local slice
                await this.fetchDevices();
            } catch (err) {
                throw err;
            }
        }
    },
});
