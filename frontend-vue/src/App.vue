<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import AppLayout from './components/layout/AppLayout.vue';
import LoadingBar from './components/ui/LoadingBar.vue';
import { useNavigation } from './composables/useNavigation';
import { ROLES } from './utils/constants';
import { api } from './services/api.js';
import { notification } from './services/notification';
import BaseToast from './components/ui/BaseToast.vue';

// Core
import Login from './views/Login.vue';

// Shared/Unified Views
import DeviceDetail from './views/shared/DeviceDetail.vue';

// Admin Views
import AdminDashboard from './views/admin/Dashboard.vue';
import AdminUserList from './views/admin/UserList.vue';
import AdminUserDetail from './views/admin/UserDetail.vue';
import AdminUserAdd from './views/admin/UserAdd.vue';
import AdminDeviceList from './views/admin/DeviceList.vue';
import AdminDeviceAdd from './views/admin/DeviceAdd.vue';
import AdminSettings from './views/admin/Settings.vue';

// Client Views
import ClientDashboard from './views/client/Dashboard.vue';
import ClientSiteList from './views/client/SiteList.vue';
import ClientSiteDetail from './views/client/SiteDetail.vue';
import ClientSiteAdd from './views/client/SiteAdd.vue';
import ClientProjectList from './views/client/ProjectList.vue';
import ClientProjectAdd from './views/client/ProjectAdd.vue';
import ClientProjectDetail from './views/client/ProjectDetail.vue';
import ClientWorkerList from './views/client/WorkerList.vue';
import ClientWorkerDetail from './views/client/WorkerDetail.vue';
import ClientWorkerAdd from './views/client/WorkerAdd.vue';

import ClientDeviceList from './views/client/DeviceList.vue';
import ClientAttendanceList from './views/client/AttendanceList.vue';

// Maintenance Views
import DeviceAddMaintenance from './views/maintenance/DeviceAddMaintenance.vue';
import SiteAssignProject from './views/maintenance/SiteAssignProject.vue';
import ProjectAssignWorkers from './views/maintenance/ProjectAssignWorkers.vue';
import UserAssignDevice from './views/maintenance/UserAssignDevice.vue';
import SiteAssignDevice from './views/maintenance/SiteAssignDevice.vue';
import WorkerAssignProject from './views/maintenance/WorkerAssignProject.vue';

const { 
  activeNavId, 
  contextData, 
  currentRole, 
  navSections, 
  breadcrumbPath,
  setRole,
  navigate 
} = useNavigation();

const isAuthenticated = ref(false);
const user = ref(null);

const handleLoginSuccess = (data) => {
  user.value = data.user;
  setRole(data.role);
  isAuthenticated.value = true;
  localStorage.setItem('auth_user', JSON.stringify(data.user));
};

const handleLogout = () => {
  isAuthenticated.value = false;
  user.value = null;
  localStorage.removeItem('auth_token'); // Clear token
  localStorage.removeItem('auth_user'); // Clear user data
  setRole(ROLES.MANAGER);
};

// Check for existing session
const checkAuth = async () => {
  const token = localStorage.getItem('auth_token');
  const savedUser = localStorage.getItem('auth_user');
  
  if (token && savedUser) {
    try {
      const userData = JSON.parse(savedUser);
      // For MVP, we trust the local user data since backend Me is a stub
      // const userData = await api.getUserProfile(); 
      handleLoginSuccess({ user: userData, role: userData.role || ROLES.MANAGER });
    } catch (e) {
      console.error("Session restore failed", e);
      localStorage.removeItem('auth_token');
      localStorage.removeItem('auth_user');
    }
  } else if (token) {
    // Fallback if no user data but token exists (legacy/edge case)
    try {
        const userData = await api.getUserProfile();
        handleLoginSuccess({ user: userData, role: userData.role });
    } catch (e) {
        localStorage.removeItem('auth_token');
    }
  }
};

onMounted(() => {
  checkAuth();
});

const componentProps = computed(() => {
  const props = { ...contextData.value };
  if (activeNavId.value === 'device-detail') {
    props.role = currentRole.value;
  }
  return props;
});

const currentComponent = computed(() => {
  if (currentRole.value === ROLES.MANAGER) {
    const componentMap = {
      'dashboard': AdminDashboard,
      'users': AdminUserList,
      'user-detail': AdminUserDetail,
      'user-add': AdminUserAdd,
      'devices': AdminDeviceList,
      'device-detail': DeviceDetail,
      'device-add': AdminDeviceAdd,
      'user-assign-device': UserAssignDevice,
      'settings': AdminSettings
    };
    return componentMap[activeNavId.value];
  } else {
    const componentMap = {
      'dashboard': ClientDashboard,
      'sites': ClientSiteList,
      'site-detail': ClientSiteDetail,
      'site-add': ClientSiteAdd,
      'device-request': SiteAssignDevice,
      'projects': ClientProjectList,
      'project-detail': ClientProjectDetail,
      'project-add': ClientProjectAdd,
      'workers': ClientWorkerList,
      'worker-detail': ClientWorkerDetail,
      'worker-add': ClientWorkerAdd,

      'worker-assign-project': WorkerAssignProject,
      'project-assign-workers': ProjectAssignWorkers,
      'site-assign-project': SiteAssignProject,
      'site-assign-device': SiteAssignDevice,
      'devices': ClientDeviceList,
      'device-detail': DeviceDetail,
      'attendance': ClientAttendanceList
    };
    return componentMap[activeNavId.value];
  }
});
</script>

<template>
  <Login v-if="!isAuthenticated" @login-success="handleLoginSuccess" />
  
  <AppLayout 
    v-else
    :role="currentRole" 
    :active-nav-id="activeNavId" 
    :user-name="user?.name"
    :nav-sections="navSections"
    :breadcrumb-path="breadcrumbPath"
    @navigate="navigate"
    @logout="handleLogout"
    @update:role="setRole"
  >
    <template v-if="currentComponent">
      <component 
        :is="currentComponent" 
        v-bind="componentProps"
        @navigate="navigate" 
      />
    </template>
    
    <div v-else class="empty-state">
      <div class="empty-content">
        <h2 class="empty-title">Page Not Found</h2>
        <p class="empty-description">
          The requested page <strong>{{ activeNavId }}</strong> could not be localized.
        </p>
        <BaseButton @click="navigate('dashboard')">Back to Dashboard</BaseButton>
      </div>
    </div>

  </AppLayout>

  <BaseToast v-if="notification.state.value" :data="notification.state.value" />
</template>

<style>
.empty-state { 
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  padding: 64px 24px; 
}
.empty-content {
  max-width: 400px;
}
.empty-title { font-size: 18px; font-weight: 600; margin-bottom: 8px; color: var(--color-text-primary); }
.empty-description { font-size: 14px; color: var(--color-text-secondary); margin-bottom: 24px; }
</style>
