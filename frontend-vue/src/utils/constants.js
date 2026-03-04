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

export const TRADES = [
    { value: '1.1', label: '1.1 - Site Management (Ancillary Works)' },
    { value: '1.2', label: '1.2 - Site Support (Ancillary Works)' },
    { value: '1.3', label: '1.3 - General Machine Operation (Ancillary Works)' },
    { value: '1.4', label: '1.4 - Site Preparation (Ancillary Works)' },
    { value: '1.5', label: '1.5 - Scaffolding (Ancillary Works)' },
    { value: '2.1', label: '2.1 - Demolition (Civil & Structural Works)' },
    { value: '2.2', label: '2.2 - Earthworks (Civil & Structural Works)' },
    { value: '2.3', label: '2.3 - Foundation (Civil & Structural works)' },
    { value: '2.4', label: '2.4 - Tunnelling (Civil & Structural Works)' },
    { value: '2.5', label: '2.5 - Reinforced Concrete (Civil & Structural Works)' },
    { value: '2.6', label: '2.6 - Structural Steel (Civil & Structural Works)' },
    { value: '2.7', label: '2.7 - Mass Engineered Timber (Civil & Structural Works)' },
    { value: '2.8', label: '2.8 - Road & Drainage (Civil & Structural Works)' },
    { value: '3.1', label: '3.1 - Ceiling (Architectural Works)' },
    { value: '3.2', label: '3.2 - Partition Wall (Architectural Works)' },
    { value: '3.3', label: '3.3 - Floor (Architectural Works)' },
    { value: '3.4', label: '3.4 - Roofing (Architectural Works)' },
    { value: '3.5', label: '3.5 - Facade (Architectural Works)' },
    { value: '3.6', label: '3.6 - Door (Architectural Works)' },
    { value: '3.7', label: '3.7 - Window (Architectural Works)' },
    { value: '3.8', label: '3.8 - Finishes (Architectural Works)' },
    { value: '3.9', label: '3.9 - Waterproofing (Architectural Works)' },
    { value: '3.10', label: '3.10 - Joinery & Fixtures (Architectural Works)' },
    { value: '3.11', label: '3.11 - Landscaping (Architectural Works)' },
    { value: '4.1', label: '4.1 - Plumbing, Sanitary & Gas (Service Works)' },
    { value: '4.2', label: '4.2 - Fire Prevention & Protection (Service Works)' },
    { value: '4.3', label: '4.3 - Electrical (Service Works)' },
    { value: '4.4', label: '4.4 - Mechanical (Service Works)' },
    { value: '4.5', label: '4.5 - Lift & Escalator (Service Works)' },
    { value: '4.6', label: '4.6 - Prefab MEP (Service Works)' }
];

export const PASS_TYPES = [
    { value: 'SP', label: 'Singapore Pink IC (SP)' },
    { value: 'SB', label: 'Singapore Blue IC (SB)' },
    { value: 'EP', label: 'Employment Pass (EP)' },
    { value: 'SPASS', label: 'S Pass (SPASS)' },
    { value: 'WP', label: 'Work Permit (WP)' },
    { value: 'ENTREPASS', label: 'EntrePass' },
    { value: 'LTVP', label: 'Long-Term Visit Pass (LTVP)' }
];
