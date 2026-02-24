import { ref, computed } from 'vue';
import { ROLES } from '../utils/constants';

// Global state acts as a singleton store
const activeNavId = ref('dashboard');
const contextData = ref(null);
const currentRole = ref(ROLES.MANAGER);

const navigation = {
    [ROLES.MANAGER]: [
        {
            title: 'Overview',
            items: [{ id: 'dashboard', label: 'Dashboard', icon: 'ri-dashboard-line' }]
        },
        {
            title: 'Management',
            items: [
                { id: 'users', label: 'Users', icon: 'ri-building-line' },
                { id: 'devices', label: 'IoT Devices', icon: 'ri-cpu-line' }
            ]
        },
        {
            title: 'Configuration',
            items: [{ id: 'settings', label: 'System Settings', icon: 'ri-settings-4-line' }]
        }
    ],
    [ROLES.CLIENT]: [
        {
            title: 'Overview',
            items: [{ id: 'dashboard', label: 'Dashboard', icon: 'ri-dashboard-line' }]
        },
        {
            title: 'Management',
            items: [
                { id: 'sites', label: 'Sites', icon: 'ri-map-pin-line' },
                { id: 'projects', label: 'Projects', icon: 'ri-folder-line' },
                { id: 'workers', label: 'Workers', icon: 'ri-group-line' },

                { id: 'devices', label: 'Devices', icon: 'ri-cpu-line' },
                { id: 'attendance', label: 'Attendance', icon: 'ri-calendar-check-line' }
            ]
        }
    ]
};

const breadcrumbMapping = {
    'dashboard': 'Dashboard',
    'users': 'Users / List',
    'user-detail': 'Users / Detail',
    'user-add': 'Users / Add New',
    'devices': 'Devices / Registry',
    'device-detail': 'Devices / Detail',
    'device-add': 'Devices / Provision',
    'device-request': 'Devices / Request Deployment',
    'sites': 'Sites / Management',
    'site-detail': 'Sites / Detail',
    'site-add': 'Sites / New Site',
    'projects': 'Projects / Management',
    'project-detail': 'Projects / Detail',
    'project-add': 'Projects / New Project',
    'workers': 'Workers / Directory',
    'worker-detail': 'Workers / Profile',
    'worker-add': 'Workers / New Worker',

    'attendance': 'Attendance Records',
    'analytics': 'Operational Insights',
    'settings': 'System Settings',
    'project-assign-workers': 'Workers / Assignment',
    'site-assign-project': 'Projects / Site Assignment',
    'site-assign-device': 'Devices / Site Assignment',
    'user-assign-device': 'Devices / User Allocation'
};

export function useNavigation() {
    const navSections = computed(() => navigation[currentRole.value]);

    const breadcrumbPath = computed(() => {
        const rolePrefix = currentRole.value === ROLES.MANAGER ? 'Manager' : 'Client';
        const path = breadcrumbMapping[activeNavId.value] || activeNavId.value;
        return `${rolePrefix} / ${path}`;
    });

    const setRole = (role) => {
        console.log("useNavigation: Switching role to", role);
        currentRole.value = role;
        activeNavId.value = 'dashboard';
        contextData.value = null;
    };

    const navigate = (id, data = null) => {
        activeNavId.value = id;
        contextData.value = data;
    };

    return {
        activeNavId,
        contextData,
        currentRole,
        navSections,
        breadcrumbPath,
        setRole,
        navigate
    };
}
