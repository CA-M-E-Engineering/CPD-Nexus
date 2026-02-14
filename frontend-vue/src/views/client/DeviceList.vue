<script setup>
import { ref, computed, onMounted } from 'vue';
import { api } from '../../services/api.js';
import { DATA_FILTERS } from '../../utils/constants';
import PageHeader from '../../components/ui/PageHeader.vue';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseFilterChip from '../../components/ui/BaseFilterChip.vue';
import TableToolbar from '../../components/ui/TableToolbar.vue';
import ConfirmDialog from '../../components/ui/ConfirmDialog.vue';

const emit = defineEmits(['navigate']);

const activeFilter = ref('All');
const filters = DATA_FILTERS.DEVICES;

const columns = [
  { key: 'device_id', label: 'Data Stream ID', size: 'md', bold: true },
  { key: 'sn', label: 'Serial Number', size: 'md' },
  { key: 'model', label: 'Hardware Model', size: 'md' },
  { key: 'last_heartbeat', label: 'Last Connection', size: 'md' },
  { key: 'site_name', label: 'Assigned Site', size: 'lg' },
  { key: 'actions', label: 'Actions', width: '120px' }
];

const devices = ref([]);
const totalDevices = ref(0);
const unassignedDevices = ref(0);
const isLoading = ref(true);
const sites = ref([]);
const showAssignModal = ref(false);
const selectedSiteId = ref('');
const selectedDeviceId = ref('');

const fetchDevices = async () => {
  isLoading.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let tenantId = null;
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            tenantId = user.tenant_id || user.id;
        } catch (e) {
            console.error('Failed to parse auth_user', e);
        }
    }

    if (!tenantId) {
        isLoading.value = false;
        return;
    }

    const [devicesData, sitesData] = await Promise.all([
      api.getDevices({ tenant_id: tenantId }),
      api.getSites({ tenant_id: tenantId })
    ]);
    
    
    devices.value = devicesData || [];
    sites.value = sitesData || [];

    // Calculate Stats
    totalDevices.value = devices.value.length;
    // Assuming 'active' status or similar for online check. 
    // The previous code had 320 workers which was wrong context.
    // Let's count assigned devices vs total.
    unassignedDevices.value = (devices.value || []).filter(d => !d.site_id).length;
    
  } catch (err) {
    console.error('Failed to load data', err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchDevices);

const filteredDevices = computed(() => {
  return devices.value.filter(device => {
    if (activeFilter.value === 'All') return true;
    const isAssigned = !!device.site_id;
    return activeFilter.value === 'Assigned' ? isAssigned : !isAssigned;
  });
});

const handleRowClick = (device) => {
  emit('navigate', 'device-detail', { id: device.device_id });
};

const handleExport = async () => {
  isLoading.value = true;
  await api.simulateExport('Devices');
  isLoading.value = false;
};

const openAssignSite = (device) => {
  selectedDeviceId.value = device.device_id;
  selectedSiteId.value = device.site_id || '';
  showAssignModal.value = true;
};

const confirmAssignSite = async () => {
  if (!selectedDeviceId.value || !selectedSiteId.value) return;
  isLoading.value = true;
  try {
    const siteId = selectedSiteId.value === 'unassign' ? 'unassign' : selectedSiteId.value;
    await api.assignDevicesToSite(siteId, [selectedDeviceId.value]);
    showAssignModal.value = false;
    await fetchDevices();
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="device-list">
    <PageHeader 
      title="IoT Devices" 
      description="Monitor and manage your site-deployed hardware"
    >
      <template #stats>
        <div class="stat-item">
          <i class="ri-cpu-line stat-icon"></i>
          <span><span class="stat-value">{{ totalDevices }}</span> devices</span>
        </div>
        <div class="stat-item">
          <i class="ri-link-unlink stat-icon"></i>
          <span><span class="stat-value">{{ unassignedDevices }}</span> unassigned</span>
        </div>
      </template>
      <template #actions>
        <BaseButton variant="secondary" @click="handleExport">Export</BaseButton>
      </template>
    </PageHeader>

    <TableToolbar>
      <template #left>
        <BaseFilterChip 
          v-for="filter in filters" 
          :key="filter"
          :label="filter"
          :active="activeFilter === filter"
          @click="activeFilter = filter"
        />
      </template>
      <template #right>
      </template>
    </TableToolbar>

    <DataTable 
      :loading="isLoading" 
      :columns="columns" 
      :data="filteredDevices"
      row-clickable
      @row-click="handleRowClick"
    >
      <template #cell-device_id="{ item }">
        <span class="mono-text">{{ item.device_id }}</span>
      </template>

      <template #cell-site_name="{ item }">
        <div class="assignment-info">
          <span v-if="item.site_name">{{ item.site_name }}</span>
          <span v-else class="text-placeholder">Unassigned</span>
        </div>
      </template>

      <template #cell-actions="{ item }">
        <div class="action-buttons">
          <BaseButton variant="ghost" size="sm" @click.stop="handleRowClick(item)">
            Manage
          </BaseButton>
          <BaseButton variant="ghost" size="sm" @click.stop="openAssignSite(item)">
            Assign Site
          </BaseButton>
        </div>
      </template>
    </DataTable>

    <div v-if="showAssignModal" class="modal-overlay" @click="showAssignModal = false">
      <div class="modal-content" @click.stop>
        <h3 class="modal-title">Assign Site</h3>
        <p class="modal-desc">Select the site to assign this device. ({{ sites.length }} sites loaded)</p>

        <div class="form-group">
          <select v-model="selectedSiteId" class="form-select">
            <option value="" disabled>Select site...</option>
            <option value="unassign" class="unassign-option">-- Unassign Device --</option>
            <option v-for="site in sites" :key="site.site_id" :value="site.site_id">
              {{ site.site_name }}
            </option>
          </select>
        </div>

        <div class="modal-footer">
          <BaseButton variant="secondary" @click="showAssignModal = false">Cancel</BaseButton>
          <BaseButton @click="confirmAssignSite">Confirm</BaseButton>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mono-text {
  font-family: var(--font-mono);
  font-size: 13px;
}

.text-placeholder {
  color: var(--color-text-secondary);
  font-style: italic;
  font-size: 13px;
}

.action-buttons {
  display: flex;
  gap: 4px;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 24px;
  width: 100%;
  max-width: 420px;
  box-shadow: var(--shadow-xl);
}

.modal-title {
  margin: 0 0 8px;
  font-size: 20px;
  font-weight: 600;
}

.modal-desc {
  margin: 0 0 24px;
  color: var(--color-text-secondary);
  font-size: 14px;
}

.form-group {
  margin-bottom: 24px;
}

.form-select {
  width: 100%;
  padding: 10px 12px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-primary);
  outline: none;
}

.form-select:focus {
  border-color: var(--color-primary);
}

.unassign-option {
  color: var(--color-error);
  font-weight: 600;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>

