<script setup>
import { ref, computed, onMounted } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import DetailCard from '../../components/ui/DetailCard.vue';
import BaseTabs from '../../components/ui/BaseTabs.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import BaseModal from '../../components/ui/BaseModal.vue';
import ConfirmDialog from '../../components/ui/ConfirmDialog.vue';
import { USER_STATUS, DEVICE_STATUS } from '../../utils/constants.js';

const props = defineProps({
  id: [Number, String]
});

const emit = defineEmits(['navigate']);

const activeTab = ref('Overview');
const User = ref(null);
const isLoading = ref(true);
const UserUsers = ref([]);
const UserSites = ref([]);
const userProjects = ref([]);
const allOrgs = ref([]);
const cpdAuthorisations = ref([]);
const siteDevicesMap = ref({}); // { siteId: [devices] }
const currentUser = ref(null);

const tabs = [
  { id: 'Overview', label: 'Company Overview' },
  { id: 'Sites', label: 'Sites' },
  { id: 'Devices', label: 'Devices' }
];

const userDevices = ref([]); // Company-level inventory

const fetchError = ref(null);

const fetchUserData = async () => {
  isLoading.value = true;
  fetchError.value = null;
  try {
    const UserData = await api.getUserById(props.id);
    if (!UserData) throw new Error('Organization not found');
    User.value = UserData;

    // Fetch others but don't block main render if they fail
    try {
        const [workersData, sitesData, alluserDevices, projectsData, authsData] = await Promise.all([
            api.getWorkers({ user_id: props.id }),
            api.getSites({ user_id: props.id }),
            api.getDevices({ user_id: props.id }),
            api.getProjects({ user_id: props.id }),
            api.getPitstopAuthorisations()
        ]);
        UserUsers.value = workersData || [];
        UserSites.value = sitesData || [];
        userDevices.value = alluserDevices || [];
        userProjects.value = projectsData || [];
        cpdAuthorisations.value = authsData || [];

        // If Vendor, fetch all users for management view
        if (UserData.user_type === 'vendor') {
            const users = await api.getUsers();
            allOrgs.value = users || [];
        }

        // Fetch devices for each site
        const devicePromises = (sitesData || []).map(async (site) => {
            try {
                const devices = await api.getDevices({ site_id: site.site_id });
                siteDevicesMap.value[site.site_id] = devices;
            } catch (e) { console.error(`Failed for site ${site.site_id}`, e); }
        });
        await Promise.all(devicePromises);
    } catch (innerErr) {
        console.warn('Partial load failed', innerErr);
        // We still have User data, so we can show the page
    }

  } catch (err) {
    console.error('Failed to fetch User data', err);
    fetchError.value = err.message || 'Failed to load organization record';
  } finally {
    isLoading.value = false;
  }
};

onMounted(async () => {
    try {
        currentUser.value = await api.getUserProfile();
    } catch (e) { console.error('Failed to get profile', e); }
    fetchUserData();
});

const UserInfo = computed(() => {
  const fields = [
    { label: 'Company Name', value: User.value?.user_name },
    { label: 'System Username', value: User.value?.username },
    { label: 'Contact Phone', value: User.value?.phone },
    { label: 'Contact Email', value: User.value?.email },
    { label: 'Bridge WebSocket', value: User.value?.bridge_ws_url },
    { label: 'User Type', value: User.value?.user_type?.toUpperCase() }
  ];
  return fields.filter(f => f.value && f.value !== '---' && f.value !== '');
});

const userColumns = [
  { key: 'name', label: 'Full Name', bold: true },
  { key: 'email', label: 'Email Address' },
  { key: 'role', label: 'Designated Role' },
  { key: 'person_trade', label: 'Trade Code', mono: true },
  { key: 'status', label: 'Status' }
];

const projectColumns = [
  { key: 'title', label: 'Project Title', bold: true },
  { key: 'reference', label: 'Reference No.', mono: true },
  { key: 'status', label: 'Current Status' },
  { key: 'site_count', label: 'Sites' },
  { key: 'worker_count', label: 'Workers' }
];

const vendorManagedColumns = [
  { key: 'username', label: 'ID/Username', bold: true },
  { key: 'email', label: 'Primary Email' },
  { key: 'bridge_ws_url', label: 'Bridge URL', mono: true },
  { key: 'cpd_status', label: 'CPD Status' },
  { key: 'status', label: 'Account' }
];

const deviceColumns = [
  { key: 'device_id', label: 'ID', bold: true },
  { key: 'model', label: 'Device Model' },
  { key: 'sn', label: 'Device SN' },
  { key: 'last_heartbeat', label: 'Last Heartbeat' },
  { key: 'user_id', label: 'User ID' },
  { key: 'actions', label: '', size: 'sm' }
];

const getCPDStatus = (userId) => {
  const auth = cpdAuthorisations.value.find(a => a.user_id === userId);
  return auth ? 'linked' : 'not linked';
};

const handleEdit = () => {
  emit('navigate', 'user-add', { id: props.id, mode: 'edit' });
};

const showDeleteConfirm = ref(false);
const isDeleting = ref(false);

const handleDelete = async () => {
  isDeleting.value = true;
  try {
    await api.deleteUser(props.id);
    notification.success('Organization deleted successfully');
    emit('navigate', 'users');
  } catch (err) {
    console.error('Failed to delete User', err);
    notification.error(err.message || 'Failed to delete organization');
  } finally {
    isDeleting.value = false;
    showDeleteConfirm.value = false;
  }
};

const handleReturnToVendor = async (deviceId) => {
    // Only allow if viewing a client OR if viewer is vendor
    const vendorId = currentUser.value?.user_id || 'Owner_001';
    
    if (confirm(`Are you sure you want to return device ${deviceId} to the vendor pool? It will be removed from this organization.`)) {
        try {
            await api.assignDevicesToUser(vendorId, [deviceId]);
            notification.success(`Device ${deviceId} returned to vendor pool`);
            await fetchUserData();
        } catch (err) {
            console.error('Failed to return device', err);
            notification.error('Failed to return device to vendor');
        }
    }
};

const handleDeviceRemove = async (deviceId) => {
    if (confirm(`Are you sure you want to decommission device ${deviceId}? This will mark it as inactive.`)) {
        try {
            await api.deleteDevice(deviceId);
            notification.success(`Device ${deviceId} decommissioned`);
            await fetchUserData();
        } catch (err) {
            console.error('Failed to remove device', err);
            notification.error('Failed to decommission device');
        }
    }
};
</script>

<template>
  <div class="user-detail">
    <PageHeader 
      :title="User?.user_name || 'Loading organization...'" 
      description="Administrative view of organization profile and resources"
      variant="detail"
    >
      <template #toolbar-left>
        <BaseButton variant="ghost" size="sm" @click="$emit('navigate', 'users')">
          <template #icon><i class="ri-arrow-left-line"></i></template>
          Back to List
        </BaseButton>
      </template>
      <template #toolbar-right>
        <BaseButton variant="secondary" size="sm" icon="ri-edit-line" @click="handleEdit">Edit Profile</BaseButton>
        <BaseButton variant="danger" size="sm" icon="ri-delete-bin-line" @click="showDeleteConfirm = true">Delete</BaseButton>
      </template>
    </PageHeader>

    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Fetching organization record...</p>
    </div>

    <div v-else-if="fetchError" class="error-state">
        <i class="ri-error-warning-line"></i>
        <p>{{ fetchError }}</p>
        <BaseButton @click="fetchUserData">Retry</BaseButton>
    </div>

    <div v-else-if="User" class="content-body">
      <BaseTabs v-model="activeTab" :tabs="tabs" />

      <!-- Overview Tab -->
      <div v-show="activeTab === 'Overview'" class="tab-content">
        <div class="overview-layout">
          <DetailCard 
            title="Basic Organization Details" 
            :badge-text="User.status" 
            :badge-type="User.status === USER_STATUS.ACTIVE ? 'success' : 'inactive'"
            :rows="UserInfo"
            class="identity-card"
          />
          
          <div v-if="User.user_type === 'vendor'" class="user-list-section">
            <div class="section-header-row">
              <h3 class="section-title">System-wide Organizational Users</h3>
              <p class="section-subtitle">Admin view of all platform tenants and their bridge configurations</p>
            </div>
            <DataTable :columns="vendorManagedColumns" :data="allOrgs" no-data-text="No other users found in system">
              <template #cell-bridge_ws_url="{ value }">
                 <span v-if="value" class="mono-text link-style">{{ value }}</span>
                 <span v-else class="text-muted">None</span>
              </template>
              <template #cell-cpd_status="{ item }">
                 <BaseBadge :type="getCPDStatus(item.user_id) === 'linked' ? 'success' : 'inactive'">
                    {{ getCPDStatus(item.user_id).toUpperCase() }}
                 </BaseBadge>
              </template>
              <template #cell-status="{ item }">
                <BaseBadge :type="item.status === 'active' ? 'success' : 'warning'">{{ item.status.toUpperCase() }}</BaseBadge>
              </template>
            </DataTable>
          </div>

          <div v-else class="user-list-section">
            <div class="section-header-row">
              <h3 class="section-title">Associated Projects</h3>
              <p class="section-subtitle">Active and historical work authorizations for this organization</p>
            </div>
            <DataTable :columns="projectColumns" :data="userProjects">
              <template #cell-status="{ item }">
                <BaseBadge :type="item.status === 'active' ? 'success' : 'warning'">{{ item.status }}</BaseBadge>
              </template>
            </DataTable>
          </div>
        </div>
      </div>

      <!-- Sites Tab -->
      <div v-show="activeTab === 'Sites'" class="tab-content">
        <div v-if="UserSites.length === 0" class="empty-state-lite">
          <p>No sites registered for this organization.</p>
        </div>
        
        <div v-for="site in UserSites" :key="site.site_id" class="site-group">
          <div class="site-header-card">
            <div class="site-title-row">
              <h3 class="site-name">{{ site.site_name }}</h3>
              <BaseBadge :type="site.status === 'active' ? 'success' : 'warning'">{{ site.status }}</BaseBadge>
            </div>
            <div class="site-meta-grid">
              <div class="meta-item"><strong>District:</strong> {{ site.location || 'N/A' }}</div>
              <div class="meta-item"><strong>Address:</strong> {{ site.address || 'N/A' }}</div>
              <div class="meta-item"><strong>Manager:</strong> {{ site.manager || 'N/A' }}</div>
              <div class="meta-item"><strong>Workers:</strong> {{ site.worker_count }}</div>
            </div>
          </div>

          <div class="site-devices-section">
            <h4 class="sub-section-title">Assigned to Site</h4>
            <DataTable 
              :columns="deviceColumns.filter(c => c.key !== 'actions')" 
              :data="siteDevicesMap[site.site_id] || []"
              no-data-text="No devices assigned to this specific site"
            >
              <template #cell-status="{ item }">
                <BaseBadge :type="item.status === DEVICE_STATUS.ONLINE ? 'success' : 'danger'">{{ item.status }}</BaseBadge>
              </template>
              <template #cell-battery="{ item }">
                <div class="battery-indicator">
                  <span :class="['battery-text', item.battery < 20 ? 'critical' : '']">{{ item.battery }}%</span>
                </div>
              </template>
            </DataTable>
          </div>
        </div>
      </div>

      <!-- Devices Tab -->
      <div v-show="activeTab === 'Devices'" class="tab-content">
        <div class="tab-header-actions">
          <h3 class="section-title">Company Hardware Inventory</h3>
          <BaseButton variant="primary" size="sm" icon="ri-add-line" @click="$emit('navigate', 'user-assign-device', { id: props.id })">
            Assign Device
          </BaseButton>
        </div>


        <DataTable :columns="deviceColumns" :data="userDevices">
          <template #cell-device_id="{ value }">
            <span class="mono-text">{{ value }}</span>
          </template>
          <template #cell-user_id="{ value }">
            <span class="mono-text">{{ value }}</span>
          </template>
          <template #cell-actions="{ item }">
            <div class="action-buttons-group">
               <BaseButton 
                v-if="item.user_id !== (currentUser?.user_id || 'Owner_001')"
                variant="ghost" 
                size="sm" 
                v-tooltip="'Return to Vendor'"
                @click="handleReturnToVendor(item.device_id)"
              >
                <i class="ri-arrow-go-back-line"></i>
              </BaseButton>
              <BaseButton 
                variant="ghost" 
                size="sm" 
                v-tooltip="'Decommission'"
                class="delete-btn"
                @click="handleDeviceRemove(item.device_id)"
              >
                <i class="ri-delete-bin-line"></i>
              </BaseButton>
            </div>
          </template>
        </DataTable>
      </div>

      <ConfirmDialog
        :show="showDeleteConfirm"
        title="Delete Organization"
        description="Are you sure you want to delete this User? This action will mark the organization as inactive."
        confirm-label="Delete User"
        :loading="isDeleting"
        variant="danger"
        @close="showDeleteConfirm = false"
        @confirm="handleDelete"
      />
    </div>
  </div>
</template>

<style scoped>
.tab-content {
  padding-top: 24px;
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.overview-layout {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.identity-card {
  max-width: 600px;
}

.section-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.section-title, .sub-section-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 0px;
}

.section-subtitle {
  font-size: 13px;
  color: var(--color-text-muted);
  margin-top: 4px;
}

.sub-section-title {
  font-size: 16px;
  margin-top: 16px;
}

.site-group {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 24px;
  margin-bottom: 32px;
  box-shadow: var(--shadow-sm);
}

.site-header-card {
  border-bottom: 1px solid var(--color-border);
  padding-bottom: 20px;
  margin-bottom: 20px;
}

.site-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.site-name {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.site-meta-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.meta-item strong {
  color: var(--color-text-primary);
  margin-right: 4px;
}

.action-buttons-group {
    display: flex;
    gap: 4px;
}

.delete-btn:hover {
    color: var(--color-danger);
}

.battery-text.critical {
  color: var(--color-danger);
  font-weight: 700;
}

.link-style {
  color: var(--color-accent);
  text-decoration: underline;
  text-underline-offset: 4px;
}

.mono-text {
  font-family: var(--font-mono, monospace);
  font-size: 12px;
  background: var(--color-bg-subtle);
  padding: 2px 6px;
  border-radius: 4px;
}

.loading-state {
  padding: 64px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: var(--color-text-secondary);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state-lite {
  padding: 48px;
  text-align: center;
  background: var(--color-bg);
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
}

.error-state {
  padding: 64px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: var(--color-danger);
}

.error-state i {
  font-size: 48px;
}
</style>

