/**
 * Global Constants & Configurations
 */

export const ROLES = {
    MANAGER: 'manager',
    PIC: 'pic',
    WORKER: 'worker',
    CLIENT: 'client'
};

export const DEVICE_TYPES = {
    GATEWAY: 'Gateway',
    LOCATOR: 'Locator',
    ENV_SENSOR: 'Env Sensor'
};

export const DEVICE_STATUS = {
    ONLINE: 'online',
    OFFLINE: 'offline',
    READY: 'ready',
    MAINTENANCE: 'maintenance'
};

export const TENANT_STATUS = {
    ACTIVE: 'active',
    PENDING: 'pending',
    INACTIVE: 'inactive'
};

export const TENANT_TYPES = {
    VENDOR: 'vendor',
    CLIENT: 'client'
};

export const DATA_FILTERS = {
    TENANTS: ['All', TENANT_TYPES.VENDOR, TENANT_TYPES.CLIENT, TENANT_STATUS.PENDING, TENANT_STATUS.INACTIVE],
    DEVICES: ['All', 'Assigned', 'Unassigned']
};

export const API_ENDPOINTS = {
    DEVICES: '/devices',
    TENANTS: '/tenants',
    SITES: '/sites',
    WORKERS: '/workers',
    PROJECTS: '/projects',
    ATTENDANCE: '/attendance'
};

export const MAP_MODES = {
    TENANTS: 'tenants',
    SITES: 'sites',
    SINGLE_EDIT: 'single-edit'
};
