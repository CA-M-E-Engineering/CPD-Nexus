<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import { MAP_MODES } from '../../utils/constants.js';
import UnifiedMap from '../../components/ui/UnifiedMap.vue';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseTabs from '../../components/ui/BaseTabs.vue';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';

const props = defineProps({
  id: [Number, String],
  mode: { type: String, default: 'add' } // 'add' or 'edit'
});

const emit = defineEmits(['navigate']);

const isSaving = ref(false);
const isLoading = ref(false);
const loadingSubData = ref(false);
const fetchError = ref(null);

const formData = ref({
  name: '',
  location: '',
  latitude: '', 
  longitude: '', 
});

const isEdit = computed(() => props.mode === 'edit');
const activeTab = ref('Workers');
const tabs = [
  { id: 'Workers', label: 'Assigned Workers' },
  { id: 'Devices', label: 'IoT Devices' },
  { id: 'Projects', label: 'Assigned Projects' }
];

const assignedWorkers = ref([]);
const assignedDevices = ref([]);
const assignedProjects = ref([]);

const availableDevices = ref([]);
const selectedDeviceToAssign = ref('');
const isAssigningDevice = ref(false);

const availableProjects = ref([]);
const selectedProjectToAssign = ref('');
const isAssigningProject = ref(false);

const workerColumns = [
  { key: 'name', label: 'Worker Name' },
  { key: 'role', label: 'Role' },
  { key: 'status', label: 'Status' },
  { key: 'actions', label: 'Actions', width: '100px' }
];

const deviceColumns = [
  { key: 'device_id', label: 'Device ID', bold: true },
  { key: 'model', label: 'Type' },
  { key: 'status', label: 'Status' },
  { key: 'battery', label: 'Battery' }
];

const projectColumns = [
  { key: 'title', label: 'Project Title' },
  { key: 'reference', label: 'Ref', muted: true },
  { key: 'status', label: 'Status' },
  { key: 'actions', label: 'Actions', width: '100px' }
];

const loadSubData = async () => {
    if (!isEdit.value || !props.id) return;
    const authUser = JSON.parse(localStorage.getItem('auth_user') || '{}');
    const userId = authUser.user_id || authUser.id;
    loadingSubData.value = true;
    try {
        const [workersData, devicesData, allDevices, allProjects] = await Promise.all([
          api.getWorkers({ user_id: userId, site_id: props.id }),
          api.getDevices({ user_id: userId, site_id: props.id }),
          api.getDevices({ user_id: userId }),
          api.getProjects({ user_id: userId })
        ]);
        assignedWorkers.value = workersData || [];
        assignedDevices.value = devicesData || [];
        
        // Projects explicitly assigned to this site
        const siteIdStr = String(props.id);
        assignedProjects.value = (allProjects || []).filter(p => p.site_id && String(p.site_id) === siteIdStr);

        // Compute available items
        availableDevices.value = (allDevices || []).filter(d => !d.site_id || String(d.site_id) !== siteIdStr);
        availableProjects.value = (allProjects || []).filter(p => !p.site_id || String(p.site_id) !== siteIdStr);
    } catch (e) {
        console.error('Failed to load sub resources', e);
    } finally {
        loadingSubData.value = false;
    }
};

const fetchSite = async () => {
  fetchError.value = null;
  if (!isEdit.value || !props.id) return;
  
  isLoading.value = true;
  try {
    let data = await api.getSiteById(props.id);
    let attempts = 0;
    while (typeof data === 'string' && attempts < 3) {
        try {
            data = data.trim();
            data = JSON.parse(data);
        } catch (e) { break; }
        attempts++;
    }

    if (data && typeof data === 'object') {
      formData.value = { 
        ...data,
        name: data.site_name || '', 
        location: data.location || '',
        latitude: data.lat || '',     
        longitude: data.lng || ''     
      };
    } else {
        fetchError.value = "Site data not found or invalid format.";
    }

    await loadSubData();

  } catch (e) {
    console.error('Error fetching site:', e);
    fetchError.value = "Failed to load site details. Please check your connection.";
    notification.error("Failed to load site details.");
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchSite);
watch(() => props.id, fetchSite);

const handleSubmit = async () => {
  isSaving.value = true;
  try {
    const authUser = JSON.parse(localStorage.getItem('auth_user') || '{}');
    const userId = authUser.user_id;

    if (!userId) {
        notification.error("Authentication error: Missing User ID.");
        return;
    }

    const payload = {
        ...formData.value,
        site_name: formData.value.name,
        user_id: userId,
        lat: parseFloat(formData.value.latitude) || 0,
        lng: parseFloat(formData.value.longitude) || 0
    };

    if (isEdit.value) {
      await api.updateSite(props.id, payload);
      notification.success('Site settings updated');
    } else {
      await api.createSite(payload);
      notification.success('New site created successfully');
      emit('navigate', 'sites');
    }
  } catch (err) {
    notification.error(err.message || 'Failed to save site record');
  } finally {
    isSaving.value = false;
  }
};

const handleAssignDevice = async () => {
    if(!selectedDeviceToAssign.value) return;
    isAssigningDevice.value = true;
    try {
        await api.assignDevicesToSite(props.id, [selectedDeviceToAssign.value]);
        notification.success('Device assigned successfully');
        selectedDeviceToAssign.value = '';
        await loadSubData();
    } catch(e) {
        notification.error('Failed to assign device');
    } finally {
        isAssigningDevice.value = false;
    }
};

const handleAssignProject = async () => {
    if(!selectedProjectToAssign.value) return;
    isAssigningProject.value = true;
    try {
        await api.assignProjectToSite(props.id, [selectedProjectToAssign.value]);
        notification.success('Project assigned successfully');
        selectedProjectToAssign.value = '';
        await loadSubData();
    } catch(e) {
        notification.error('Failed to assign project');
    } finally {
        isAssigningProject.value = false;
    }
};

const handleManageWorker = (worker) => {
  emit('navigate', 'worker-add', { id: worker.user_id || worker.worker_id, mode: 'edit' });
};
const handleManageDevice = (device) => {
  emit('navigate', 'device-detail', { id: device.device_id });
};
const handleManageProject = (project) => {
  emit('navigate', 'project-detail', { id: project.project_id });
};
</script>

<template>
  <div class="site-add">
    <PageHeader 
      :title="isEdit ? 'Site Configuration Dashboard' : 'Add New Site'" 
      :description="isEdit ? 'Manage site properties and active assignments' : 'Create a new physical location for operations'"
    />

    <div v-if="isLoading" class="loading-state">
      <i class="ri-loader-4-line spin-icon"></i> Loading site data...
    </div>
    
    <div v-else-if="fetchError" class="error-state">
      <i class="ri-error-warning-line"></i>
      <p class="error-text">{{ fetchError }}</p>
      <BaseButton size="sm" @click="fetchSite">Retry</BaseButton>
    </div>

    <div v-else class="dashboard-layout">
      <!-- ── TOP SECTION (Split View) ── -->
      <div class="top-section">
        
        <!-- Site Details Form -->
        <div class="form-section-card details-card">
           <h3 class="section-header">Site Properties</h3>
           <form @submit.prevent="handleSubmit" class="form-col">
              <div class="form-grid">
                <BaseInput v-model="formData.name" label="Site Name" placeholder="e.g., Northshore Tunnel" required />
                <BaseInput v-model="formData.location" label="Location / Zone" placeholder="e.g., Woodlands" />
                <BaseInput v-model="formData.latitude" label="Latitude" placeholder="e.g., 1.3521" type="number" step="any" />
                <BaseInput v-model="formData.longitude" label="Longitude" placeholder="e.g., 103.8198" type="number" step="any" />
              </div>

              <div class="sticky-actions">
                <BaseButton variant="secondary" @click="$emit('navigate', 'sites')">Cancel</BaseButton>
                <BaseButton :loading="isSaving" type="submit">
                  {{ isEdit ? 'Save Changes' : 'Create Site' }}
                </BaseButton>
              </div>
           </form>
        </div>

        <!-- Site Map -->
        <div class="form-section-card map-card">
           <div class="profile-map-wrapper">
             <UnifiedMap 
                :mode="isEdit ? MAP_MODES.SINGLE_EDIT : MAP_MODES.SINGLE_EDIT"
                :lat="parseFloat(formData.latitude) || 1.3521"
                :lng="parseFloat(formData.longitude) || 103.8198"
                @update:lat="v => formData.latitude = v"
                @update:lng="v => formData.longitude = v"
             />
          </div>
        </div>
      </div>

      <!-- ── BOTTOM SECTION (Full Width Resources) ── -->
      <div class="bottom-section form-section-card" v-if="isEdit">
         <h3 class="section-header">Site Resources</h3>
         <BaseTabs v-model="activeTab" :tabs="tabs" />

         <!-- Workers Tab -->
         <div v-show="activeTab === 'Workers'" class="tab-content">
            <DataTable :loading="loadingSubData" :columns="workerColumns" :data="assignedWorkers" row-clickable @row-click="handleManageWorker">
              <template #cell-name="{ item }">
                  <strong>{{ item.name }}</strong>
              </template>
              <template #cell-status="{ item }">
                <BaseBadge :type="item.status === 'active' ? 'success' : 'inactive'">{{ item.status }}</BaseBadge>
              </template>
              <template #cell-actions="{ item }">
                <BaseButton variant="secondary" size="sm" @click.stop="handleManageWorker(item)">Manage</BaseButton>
              </template>
            </DataTable>
            <div v-if="!loadingSubData && assignedWorkers.length === 0" class="empty-state">
                <p>No workers currently assigned to this site.</p>
            </div>
         </div>

         <!-- Devices Tab -->
         <div v-show="activeTab === 'Devices'" class="tab-content">
            <!-- Inline Assignment Bar -->
            <div class="assign-bar">
               <select v-model="selectedDeviceToAssign" class="form-select inline-select">
                  <option value="" disabled selected>Select an available device...</option>
                  <option v-for="d in availableDevices" :key="d.device_id" :value="d.device_id">
                     {{ d.device_id }} - {{ d.model }}
                  </option>
               </select>
               <BaseButton :disabled="!selectedDeviceToAssign" :loading="isAssigningDevice" @click="handleAssignDevice">
                   Assign Device
               </BaseButton>
            </div>

            <DataTable :loading="loadingSubData" :columns="deviceColumns" :data="assignedDevices" row-clickable @row-click="handleManageDevice">
              <template #cell-status="{ item }">
                <BaseBadge :type="item.status === 'online' ? 'success' : 'inactive'">
                  {{ item.status ? item.status.toUpperCase() : 'UNKNOWN' }}
                </BaseBadge>
              </template>
              <template #cell-battery="{ item }">
                <span>{{ item.battery }}%</span>
              </template>
            </DataTable>
            <div v-if="!loadingSubData && assignedDevices.length === 0" class="empty-state">
                <p>No devices currently assigned to this site.</p>
            </div>
         </div>

         <!-- Projects Tab -->
         <div v-show="activeTab === 'Projects'" class="tab-content">
             <!-- Inline Assignment Bar -->
            <div class="assign-bar">
               <select v-model="selectedProjectToAssign" class="form-select inline-select">
                  <option value="" disabled selected>Select an available project...</option>
                  <option v-for="p in availableProjects" :key="p.project_id" :value="p.project_id">
                     {{ p.title }} ({{ p.reference }})
                  </option>
               </select>
               <BaseButton :disabled="!selectedProjectToAssign" :loading="isAssigningProject" @click="handleAssignProject">
                   Assign Project
               </BaseButton>
            </div>

            <DataTable :loading="loadingSubData" :columns="projectColumns" :data="assignedProjects" row-clickable @row-click="handleManageProject">
              <template #cell-title="{ item }">
                  <strong>{{ item.title }}</strong>
              </template>
              <template #cell-status="{ item }">
                <BaseBadge :type="item.status === 'active' ? 'info' : item.status === 'completed' ? 'success' : 'warning'">
                  {{ item.status.charAt(0).toUpperCase() + item.status.slice(1) }}
                </BaseBadge>
              </template>
              <template #cell-actions="{ item }">
                <BaseButton variant="secondary" size="sm" @click.stop="handleManageProject(item)">View</BaseButton>
              </template>
            </DataTable>
            <div v-if="!loadingSubData && assignedProjects.length === 0" class="empty-state">
                <p>No projects currently assigned to this site.</p>
            </div>
         </div>

      </div>

    </div>
  </div>
</template>

<style scoped>
/* ── Layout ── */
.dashboard-layout {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.top-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  align-items: stretch;
}

.bottom-section {
  width: 100%;
}

/* ── Cards ── */
.form-section-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 24px;
  box-shadow: var(--shadow-sm);
}

.details-card {
  display: flex;
  flex-direction: column;
}

.map-card {
  padding: 0;
  overflow: hidden;
  display: flex;
}

.profile-map-wrapper {
  flex: 1;
  width: 100%;
  min-height: 400px;
  background: var(--color-bg-subtle);
}

.section-header {
  margin-bottom: 20px;
  padding-bottom: 14px;
  border-bottom: 1px solid var(--color-border);
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-strong);
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.sticky-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
  margin-top: auto;
}

.assign-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  align-items: center;
}

.inline-select {
  flex: 1;
  max-width: 400px;
}

.form-select {
    padding: 10px 12px;
    background: var(--color-bg);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text-primary);
    font-size: 14px;
    outline: none;
}

.form-select:focus {
    border-color: var(--color-accent);
}

.tab-content {
  margin-top: 20px;
}

.loading-state {
  padding: 48px;
  text-align: center;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.spin-icon {
  animation: spin 1s linear infinite;
  display: inline-block;
}

@keyframes spin {
  100% { transform: rotate(360deg); }
}

.error-state {
  padding: 48px;
  text-align: center;
  background: var(--color-surface);
  border: 1px solid var(--color-danger);
  border-radius: var(--radius-md);
  margin-bottom: 24px;
}

.error-text {
  color: var(--color-danger);
  margin-bottom: 16px;
}

.empty-state {
  padding: 32px;
  text-align: center;
  color: var(--color-text-secondary);
  font-style: italic;
  font-size: 14px;
}

@media (max-width: 1024px) {
  .top-section { grid-template-columns: 1fr; }
}

@media (max-width: 640px) {
  .form-grid { grid-template-columns: 1fr; }
  .assign-bar { flex-direction: column; align-items: stretch; }
  .inline-select { max-width: none; }
}
</style>

