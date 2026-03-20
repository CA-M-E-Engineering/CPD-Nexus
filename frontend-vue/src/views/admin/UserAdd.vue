<script setup>
import { ref, onMounted, computed } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import { USER_TYPES, MAP_MODES } from '../../utils/constants.js';
import UnifiedMap from '../../components/ui/UnifiedMap.vue';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseTabs from '../../components/ui/BaseTabs.vue';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import ConfirmDialog from '../../components/ui/ConfirmDialog.vue';

const props = defineProps({
  id: [Number, String],
  mode: {
    type: String,
    default: 'add'
  }
});

const emit = defineEmits(['navigate']);

const isEdit = computed(() => props.mode === 'edit');
const isLoading = ref(false);
const isSaving = ref(false);

const formData = ref({
  user_name: '',
  user_type: USER_TYPES.CLIENT,
  username: '',
  password: '',
  email: '',
  phone: '',
  address: '',
  latitude: '',
  longitude: '',
  status: 'active',
  bridge_ws_url: '',
  bridge_auth_token: '',
  bridge_status: 'inactive',
  assigned_on_behalf_ofs: []
});

const availableOnBehalfOfs = ref([]);
const activeTab = ref('Overview');
const userSites = ref([]);
const userProjects = ref([]);
const allOrgs = ref([]);
const userDevices = ref([]);
const cpdAuthorisations = ref([]);
const siteDevicesMap = ref({});
const currentUser = ref(null);
const loadingResources = ref(false);

const tabs = computed(() => {
  const baseTabs = [
    { id: 'Overview', label: 'Company Overview' },
    { id: 'Sites', label: 'Sites' },
    { id: 'Devices', label: 'Devices' }
  ];
  return baseTabs;
});

const loadResources = async () => {
    if (!isEdit.value || !props.id) return;
    loadingResources.value = true;
    try {
        const [sitesData, allDevices, projectsData, authsData] = await Promise.all([
            api.getSites({ user_id: props.id }),
            api.getDevices({ user_id: props.id }),
            api.getProjects({ user_id: props.id }),
            api.getPitstopAuthorisations()
        ]);
        userSites.value = sitesData || [];
        userDevices.value = allDevices || [];
        userProjects.value = projectsData || [];
        cpdAuthorisations.value = authsData || [];

        // If viewing/editing a vendor-type user
        if (formData.value.user_type === 'vendor') {
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
    } catch (e) {
        console.error('Failed to load sub-resources', e);
    } finally {
        loadingResources.value = false;
    }
};

const fetchUserAndAuthorisations = async () => {
  isLoading.value = true;
  try {
    const promises = [api.getPitstopAuthorisations()];
    if (isEdit.value && props.id) {
        promises.push(api.getUserById(props.id));
    }
    
    const results = await Promise.all(promises);
    const authsData = results[0] || [];
    // Extract unique onBehalfOfs from pitstop authorisations
    const onBehalfOfSet = new Set();
    const assignedToUser = new Set();
    
    authsData.forEach(auth => {
        if (auth.on_behalf_of_name) {
            onBehalfOfSet.add(auth.on_behalf_of_name);
            if (isEdit.value && props.id && auth.user_id === props.id) {
                assignedToUser.add(auth.on_behalf_of_name);
            }
        }
    });
    
    availableOnBehalfOfs.value = Array.from(onBehalfOfSet).sort();
    
    // Setup form data
    if (isEdit.value && results[1]) {
      const data = results[1];
      formData.value = {
        ...formData.value,
        user_name: data.user_name,
        user_type: data.user_type || USER_TYPES.CLIENT,
        username: data.username || '',
        email: data.email || '',
        phone: data.phone || '',
        address: data.address || '',
        latitude: data.lat || '',
        longitude: data.lng || '',
        status: data.status || 'active',
        bridge_ws_url: data.bridge_ws_url || '',
        bridge_auth_token: data.bridge_auth_token || '',
        bridge_status: data.bridge_status || 'inactive',
        assigned_on_behalf_ofs: Array.from(assignedToUser)
      };
    }
    
    await loadResources();
  } catch (err) {
      console.error('Failed to init user form:', err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(async () => {
    try {
        currentUser.value = await api.getUserProfile();
    } catch (e) { console.error('Failed to get profile', e); }
    fetchUserAndAuthorisations();
});

const handleSubmit = async () => {
  isSaving.value = true;
  try {
    const payload = {
        ...formData.value,
        bridge_ws_url: formData.value.bridge_ws_url || serverBridgeUrl.value,
        lat: parseFloat(formData.value.latitude) || 0,
        lng: parseFloat(formData.value.longitude) || 0
    };

    let savedUserId = props.id;

    if (isEdit.value) {
      await api.updateUser(savedUserId, payload);
      notification.success('Organization updated successfully');
    } else {
      const response = await api.createUser(payload);
      savedUserId = response.user_id || response.id; // handle id retrieval based on create response mapping
      notification.success('New organization registered successfully');
    }
    
    // Now dispatch pitstop onBehalfOfs assignment array
    if (savedUserId && formData.value.assigned_on_behalf_ofs) {
        try {
            await api.assignPitstopOnBehalfOfs(savedUserId, formData.value.assigned_on_behalf_ofs);
        } catch (authErr) {
            console.error('Pitstop assignment failed:', authErr);
            notification.error('Organization saved, but Pitstop On Behalf Of assignments failed to sync.');
        }
    }

    emit('navigate', 'users');
  } catch (err) {
    console.error('[UserAdd] Save failed:', err);
    notification.error(err.message || 'Failed to save organization details');
  } finally {
    isSaving.value = false;
  }
};

// Resource Management Actions
const showDeleteConfirm = ref(false);
const isDeleting = ref(false);

const handleDelete = async () => {
  isDeleting.value = true;
  try {
    await api.deleteUser(props.id);
    notification.success('Organization deleted successfully');
    emit('navigate', 'users');
  } catch (err) {
    notification.error(err.message || 'Failed to delete organization');
  } finally {
    isDeleting.value = false;
    showDeleteConfirm.value = false;
  }
};

const handleReturnToVendor = (deviceId) => {
    deviceToReturn.value = deviceId;
    showReturnConfirm.value = true;
};

const confirmReturnToVendor = async () => {
    if (!deviceToReturn.value) return;
    isReturning.value = true;
    try {
        await api.assignDevicesToUser(null, [deviceToReturn.value]);
        notification.success('Device successfully returned to vendor pool');
        showReturnConfirm.value = false;
        await loadResources();
    } catch (e) {
        notification.error('Failed to return device to pool');
    } finally {
        isReturning.value = false;
        deviceToReturn.value = null;
    }
};

// Table Columns
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
  { key: 'actions', label: 'Actions', size: 'sm' }
];

const getCPDStatus = (userId) => {
  const auth = cpdAuthorisations.value.find(a => a.user_id === userId);
  return auth ? 'linked' : 'not linked';
};

const serverBridgeUrl = computed(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.hostname;
    const port = host === 'localhost' || host === '127.0.0.1' ? ':3010' : '';
    return `${protocol}//${host}${port}/api/v1/bridge/connect`;
});

const fullConnectionString = computed(() => {
    if (!isEdit.value || !props.id || !formData.value.bridge_auth_token) {
        return `${serverBridgeUrl.value}?user_id=[ID]&token=[TOKEN]`;
    }
    return `${serverBridgeUrl.value}?user_id=${props.id}&token=${formData.value.bridge_auth_token}`;
});

const copyToClipboard = (text, label) => {
    navigator.clipboard.writeText(text).then(() => {
        notification.success(`${label} copied to clipboard`);
    }).catch(() => {
        notification.error('Failed to copy to clipboard');
    });
};

const generateToken = () => {
    const chars = 'abcdef0123456789';
    let result = '';
    for (let i = 0; i < 32; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    formData.value.bridge_auth_token = result;
    notification.success('New secret token generated (remember to save)');
};
</script>

<template>
  <div class="user-add">
    <PageHeader
      :title="isEdit ? formData.user_name : 'Register New User'"
      :description="isEdit ? 'Update details for this organization account' : 'Add a new client organization to the system'"
      :variant="isEdit ? 'detail' : 'default'"
    >
      <template #toolbar-left v-if="isEdit">
        <BaseButton variant="ghost" size="sm" @click="$emit('navigate', 'users')">
          <template #icon><i class="ri-arrow-left-line"></i></template>
          Back to List
        </BaseButton>
      </template>
      <template #toolbar-right v-if="isEdit">
        <BaseButton variant="danger" size="sm" icon="ri-delete-bin-line" @click="showDeleteConfirm = true">Delete Organization</BaseButton>
      </template>
    </PageHeader>

    <div v-if="isLoading" class="loading-state">
      <p>Loading record details...</p>
    </div>

    <form v-else class="form-container" @submit.prevent="handleSubmit">
      <div v-if="!isEdit" class="form-info-banner">
        <i class="ri-information-line"></i>
        <span>Creating a User will automatically generate a <strong>Login Account</strong> and a corresponding <strong>Primary Business Entity</strong>.</span>
      </div>

      <!-- Dashboard Top Section -->
      <div class="dashboard-grid">
         
         <!-- COLUMN 1: Account, Identity, Bridge & integration -->
         <div class="left-main-stack">
            <!-- Account & Identity -->
            <div class="form-panel glass-card">
               <h3 class="panel-title"><i class="ri-user-settings-line"></i> Account & Identity</h3>
               <div class="row-grid">
                  <BaseInput v-model="formData.user_name" label="Organization Name" placeholder="e.g., Acme Corp" class="full-width" required />
                  <BaseInput v-model="formData.email" label="Primary Email" type="email" placeholder="admin@acme.com" required />
                  <BaseInput v-model="formData.username" label="System Username" placeholder="e.g., acme_admin" required />
                  <BaseInput v-model="formData.password" label="Access Password" type="password" placeholder="Min 6 chars" :required="!isEdit" />
                  <BaseInput v-model="formData.phone" label="Contact Phone" placeholder="+65 8888 8888" />
                  <div class="form-group" v-if="isEdit">
                     <label class="form-label">Account Status</label>
                     <select v-model="formData.status" class="form-select">
                       <option value="active">Active</option>
                       <option value="inactive">Inactive</option>
                       <option value="pending">Pending</option>
                     </select>
                  </div>
               </div>
            </div>

            <!-- Bridge Panel -->
            <div class="form-panel glass-card bridge-card">
               <div class="panel-header">
                  <h3 class="panel-title"><i class="ri-wifi-line"></i> IoT Bridge Connection</h3>
                  <BaseBadge :type="formData.bridge_status === 'active' ? 'success' : 'inactive'">
                     {{ formData.bridge_status.toUpperCase() }}
                  </BaseBadge>
               </div>
               
               <div class="bridge-info-row">
                  <div class="info-item">
                     <span class="info-label">Full Connection String (Copy to Bridge)</span>
                     <div class="info-value-group">
                        <code class="mono-text">{{ fullConnectionString }}</code>
                        <BaseButton variant="ghost" size="sm" icon="ri-file-copy-line" @click="copyToClipboard(fullConnectionString, 'Connection String')" title="Copy Full String" />
                     </div>
                  </div>
                  <div class="info-item mt-2">
                     <span class="info-label">Individual Credentials</span>
                     <div class="creds-stack">
                        <div class="cred-pill" @click="copyToClipboard(serverBridgeUrl, 'Gateway URL')">
                           <span class="l">URL:</span> <code class="v">{{ serverBridgeUrl }}</code>
                        </div>
                        <div class="cred-pill" v-if="isEdit" @click="copyToClipboard(props.id, 'User ID')">
                           <span class="l">ID:</span> <code class="v">{{ props.id }}</code>
                        </div>
                        <div class="cred-pill" v-if="formData.bridge_auth_token" @click="copyToClipboard(formData.bridge_auth_token, 'Secret Token')">
                           <span class="l">TOKEN:</span> <code class="v">{{ formData.bridge_auth_token }}</code>
                        </div>
                     </div>
                  </div>
               </div>

               <div class="row-grid mt-3">
                  <div class="form-group full-width">
                     <label class="form-label">Secret Authentication Token</label>
                     <div class="input-with-action">
                        <BaseInput 
                           v-model="formData.bridge_auth_token" 
                           type="text" 
                           placeholder="Generates automatically on save if empty" 
                           class="flex-grow no-margin"
                        />
                        <div class="action-buttons">
                           <BaseButton variant="secondary" size="sm" icon="ri-refresh-line" @click="generateToken" title="Regenerate Token" />
                           <BaseButton variant="ghost" size="sm" icon="ri-file-copy-line" @click="copyToClipboard(formData.bridge_auth_token, 'Secret Token')" title="Copy Token" />
                        </div>
                     </div>
                     <span class="form-hint">Used by the bridge for secure authentication.</span>
                  </div>
                  
                  <div class="form-group full-width">
                     <BaseInput v-model="formData.bridge_ws_url" label="Bridge Software Version" placeholder="e.g. v1.2.0 (Optional reference)" />
                  </div>

                  <div class="form-group full-width">
                     <label class="form-label">Auto-Connect Priority</label>
                     <select v-model="formData.bridge_status" class="form-select">
                       <option value="active">Enabled — Connect and sync automatically</option>
                       <option value="inactive">Disabled — Manual toggle only</option>
                     </select>
                  </div>
               </div>
            </div>

            <!-- Integration Panel -->
            <div class="form-panel glass-card integration-card">
               <h3 class="panel-title"><i class="ri-link-m"></i> CPD Integration</h3>
               <div class="integration-body">
                  <p class="panel-hint">Assign 'On Behalf Of' entities to map datasets directly to this client.</p>
                  <div class="dropdown-wrapper">
                    <details class="custom-dropdown">
                        <summary class="dropdown-summary">
                            <span>{{ formData.assigned_on_behalf_ofs.length > 0 ? `${formData.assigned_on_behalf_ofs.length} Selected` : 'Select sync entities...' }}</span>
                            <i class="ri-arrow-down-s-line"></i>
                        </summary>
                        <div class="dropdown-content contractor-list">
                            <label v-for="entity in availableOnBehalfOfs" :key="entity" class="checkbox-label">
                                <input type="checkbox" :value="entity" v-model="formData.assigned_on_behalf_ofs" />
                                <span>{{ entity }}</span>
                            </label>
                        </div>
                    </details>
                  </div>
               </div>
            </div>
         </div>

         <!-- COLUMN 2: Map & Address -->
         <div class="map-side-panel glass-card">
            <h3 class="panel-title"><i class="ri-map-pin-line"></i> Registered Address & Geofence</h3>
            <div class="address-fields">
               <BaseInput v-model="formData.address" label="Street Address" placeholder="123 Discovery Way" class="full-width" />
               <div class="coord-grid">
                  <BaseInput v-model="formData.latitude" label="Lat" type="number" step="any" />
                  <BaseInput v-model="formData.longitude" label="Lng" type="number" step="any" />
               </div>
            </div>
            <div class="map-frame-lg">
               <UnifiedMap
                 :mode="MAP_MODES.SINGLE_EDIT"
                 :lat="parseFloat(formData.latitude) || 1.3521"
                 :lng="parseFloat(formData.longitude) || 103.8198"
                 @update:lat="v => formData.latitude = v"
                 @update:lng="v => formData.longitude = v"
               />
            </div>
         </div>
      </div>

      <!-- Resource Dashboard (Edit Mode Only) -->
      <div v-if="isEdit" class="resources-section glass-card">
         <div class="resources-header">
            <h3 class="panel-title"><i class="ri-archive-line"></i> User Resource Inventory</h3>
            <BaseTabs v-model="activeTab" :tabs="tabs" />
         </div>

         <!-- Overview/Sub-users Tab (For Vendors) -->
         <div v-show="activeTab === 'Overview'" class="tab-pane">
            <div v-if="formData.user_type === 'vendor'" class="sub-resource-group">
               <div class="resource-meta">
                  <h4>Managed System Groups</h4>
                  <p>System-wide view of platform tenants and their configurations</p>
               </div>
               <DataTable :columns="vendorManagedColumns" :data="allOrgs" no-data-text="No managed organizations found">
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
            <div v-else class="sub-resource-group">
               <div class="resource-meta">
                  <h4>Associated Project Portfolio</h4>
                  <p>Active and historical work authorizations for this account</p>
               </div>
               <DataTable :columns="projectColumns" :data="userProjects" />
            </div>
         </div>

         <!-- Sites Tab -->
         <div v-show="activeTab === 'Sites'" class="tab-pane">
            <div class="sites-grid">
               <div v-for="site in userSites" :key="site.site_id" class="site-summary-card">
                  <div class="site-card-header">
                     <div class="site-info">
                        <strong>{{ site.site_name }}</strong>
                        <span class="location-hint">{{ site.location || 'No location set' }}</span>
                     </div>
                     <BaseBadge>{{ site.status }}</BaseBadge>
                  </div>
                  <div class="site-card-body">
                     <div class="site-stat"><span>Workers:</span> <strong>{{ site.worker_count }}</strong></div>
                     <div class="site-stat"><span>Devices:</span> <strong>{{ siteDevicesMap[site.site_id]?.length || 0 }}</strong></div>
                  </div>
               </div>
            </div>
            <div v-if="userSites.length === 0" class="empty-hint">No registered sites found for this user.</div>
         </div>

         <!-- Devices Tab -->
         <div v-show="activeTab === 'Devices'" class="tab-pane">
            <DataTable :columns="deviceColumns" :data="userDevices">
               <template #cell-actions="{ item }">
                  <div class="btn-stack sm">
                     <BaseButton 
                       v-if="item.user_id !== (currentUser?.user_id || 'Owner_001')"
                       variant="primary" size="sm" @click="handleReturnToVendor(item.device_id)">
                       Return to Pool
                     </BaseButton>
                  </div>
               </template>
            </DataTable>
         </div>
      </div>

      <div class="form-actions">
        <BaseButton variant="secondary" @click="$emit('navigate', 'users')">Cancel</BaseButton>
        <BaseButton :loading="isSaving" type="submit" variant="primary">
          {{ isEdit ? 'Save Changes' : 'Register User' }}
        </BaseButton>
      </div>
    </form>

    <ConfirmDialog
      :show="showDeleteConfirm"
      title="Delete Organization"
      description="Permanently mark this organization as inactive? Access will be revoked immediately."
      confirm-label="Confirm Deletion"
      :loading="isDeleting"
      @confirm="handleDelete"
      @cancel="showDeleteConfirm = false"
    />

    <ConfirmDialog
      :show="showReturnConfirm"
      title="Return Device to Pool"
      description="Returning this device will immediately revoke this user's access to the device and move it back to the vendor's unassigned inventory."
      confirm-label="Return to Pool"
      :loading="isReturning"
      @confirm="confirmReturnToVendor"
      @cancel="showReturnConfirm = false"
    />
  </div>
</template>

<style scoped>
.form-container {
  max-width: 1400px;
  display: flex;
  flex-direction: column;
  gap: 32px;
  padding-bottom: 64px;
}

.glass-card {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-lg);
  padding: 32px;
  box-shadow: 0 8px 32px -8px rgba(0, 0, 0, 0.3);
}

.panel-title {
  font-size: 16px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--color-accent);
  letter-spacing: 0.05em;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.panel-title i {
  font-size: 20px;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
}

.left-main-stack {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.map-side-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.map-side-panel .address-fields {
  margin-bottom: 24px;
}

.map-frame-lg {
  flex-grow: 1;
  min-height: 650px;
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.row-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
}

.full-width {
  grid-column: 1 / -1;
}

/* Bridge Card */
.bridge-card .panel-title {
  margin-bottom: 0;
}

.bridge-info-row {
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: var(--radius-sm);
  padding: 16px;
  margin-bottom: 24px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.info-label {
  font-size: 11px;
  text-transform: uppercase;
  color: var(--color-text-muted);
  font-weight: 600;
}

.info-value-group {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.mono-text {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--color-accent);
  background: rgba(139, 92, 246, 0.1);
  padding: 4px 8px;
  border-radius: 4px;
  word-break: break-all;
  line-height: 1.4;
}

.creds-stack {
   display: flex;
   flex-wrap: wrap;
   gap: 8px;
   margin-top: 4px;
}

.cred-pill {
   background: rgba(255, 255, 255, 0.05);
   border: 1px solid rgba(255, 255, 255, 0.08);
   border-radius: 4px;
   padding: 4px 8px;
   font-size: 11px;
   cursor: pointer;
   display: flex;
   gap: 6px;
   transition: all 0.2s ease;
}

.cred-pill:hover {
   background: rgba(255, 255, 255, 0.1);
   border-color: var(--color-accent);
}

.cred-pill .l { color: var(--color-text-muted); font-weight: 600; }
.cred-pill .v { color: #fff; font-family: var(--font-mono); }

.mt-2 { margin-top: 8px; }

.input-with-action {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.input-with-action :deep(.form-group) {
  margin-bottom: 0;
  flex: 1;
}

.action-buttons {
  display: flex;
  gap: 8px;
  padding-top: 4px; /* Align with the inner input height */
}

.no-margin :deep(.form-group) {
  margin-bottom: 0;
}

.mt-3 {
  margin-top: 24px;
}

/* Integration Card */
.panel-hint {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin-bottom: 16px;
  line-height: 1.5;
}

/* Map Panel */
.map-layout {
  display: grid;
  grid-template-columns: 350px 1fr;
  gap: 32px;
  align-items: start;
}

.address-fields {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.coord-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.map-frame {
  height: 300px;
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

/* Resources Section */
.resources-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.resources-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.tab-pane {
  animation: fadeIn 0.3s ease-out;
}

.resource-meta {
  margin-bottom: 24px;
}

.resource-meta h4 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 4px;
}

.resource-meta p {
  font-size: 13px;
  color: var(--color-text-secondary);
}

/* Sites Grid */
.sites-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.site-summary-card {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-md);
  padding: 20px;
  transition: transform 0.2s ease;
}

.site-summary-card:hover {
  transform: translateY(-4px);
  background: rgba(255, 255, 255, 0.04);
}

.site-card-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
}

.site-info {
  display: flex;
  flex-direction: column;
}

.location-hint {
  font-size: 12px;
  color: var(--color-text-muted);
}

.site-card-body {
  display: flex;
  gap: 20px;
  font-size: 13px;
}

.site-stat span {
  color: var(--color-text-muted);
  margin-right: 4px;
}

/* Action Styles */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  padding-top: 32px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.btn-stack {
  display: flex;
  gap: 8px;
}

.btn-stack.sm {
  gap: 4px;
}

.form-select {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-md);
  padding: 10px 12px;
  color: var(--color-text-primary);
  font-size: 14px;
  outline: none;
  width: 100%;
}

.form-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: 6px;
  display: block;
}

.empty-hint {
  padding: 64px;
  text-align: center;
  color: var(--color-text-muted);
  font-style: italic;
  background: rgba(0,0,0,0.1);
  border-radius: var(--radius-md);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Dropdown */
.dropdown-summary {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 14px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: var(--radius-md);
    cursor: pointer;
    font-size: 14px;
}

.dropdown-content {
    background: #1a1f2e;
    border: 1px solid rgba(255,255,255,0.1);
    border-radius: var(--radius-md);
    padding: 8px;
    margin-top: 4px;
    max-height: 200px;
    overflow-y: auto;
}

.checkbox-label {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px;
    cursor: pointer;
    border-radius: 4px;
}

.checkbox-label:hover {
    background: rgba(255,255,255,0.05);
}

@media (max-width: 1100px) {
  .dashboard-grid, .map-layout {
    grid-template-columns: 1fr;
  }
}
</style>
