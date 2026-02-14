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

export const USER_STATUS = {
    ACTIVE: 'active',
    PENDING: 'pending',
    INACTIVE: 'inactive'
};

export const USER_TYPES = {
    VENDOR: 'vendor',
    CLIENT: 'client'
};

export const DATA_FILTERS = {
    USERS: ['All', USER_TYPES.VENDOR, USER_TYPES.CLIENT, USER_STATUS.PENDING, USER_STATUS.INACTIVE],
    DEVICES: ['All', 'Assigned', 'Unassigned']
};

export const API_ENDPOINTS = {
    DEVICES: '/devices',
    USERS: '/users',
    SITES: '/sites',
    WORKERS: '/workers',
    PROJECTS: '/projects',
    ATTENDANCE: '/attendance'
};

export const MAP_MODES = {
    USERS: 'users',
    SITES: 'sites',
    SINGLE_EDIT: 'single-edit'
};
