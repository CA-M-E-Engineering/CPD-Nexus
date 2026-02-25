<script setup>
import { ref, computed, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useDeviceStore } from '../../features/devices/store/deviceStore';
import { api } from '../../services/api.js';
import { DATA_FILTERS } from '../../utils/constants';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseFilterChip from '../../components/ui/BaseFilterChip.vue';
import PageHeader from '../../components/ui/PageHeader.vue';
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
  { key: 'site_name', label: 'Current Assignment', size: 'lg' },
  { key: 'actions', label: 'Actions', width: '120px' }
];

const deviceStore = useDeviceStore();
const { devices, loading: isLoading } = storeToRefs(deviceStore);

// Computed stats from store
const siteCount = ref(0); // TODO: Move sites to store
const workerCount = ref(0); // TODO: Move workers to store

const fetchDevices = async () => {
  try {
    // Parallel fetch: Devices from store, others from API (for now)
    await Promise.all([
        deviceStore.fetchDevices(),
        api.getSites().then(data => siteCount.value = (data || []).length),
        api.getWorkers().then(data => workerCount.value = (data || []).length)
    ]);
  } catch (err) {
    console.error('Fetch Failed:', err);
  }
};

onMounted(fetchDevices);

const filteredDevices = computed(() => {
  if (!devices.value) return [];
  return devices.value.filter(device => {
    if (activeFilter.value === 'All') return true;
    const isAssigned = !!device.user_id && device.user_id !== 'tenant-vendor-1';
    return activeFilter.value === 'Assigned' ? isAssigned : !isAssigned;
  });
});

const showDeleteDialog = ref(false);
const deviceToDelete = ref(null);
const isDeleting = ref(false);

const confirmDelete = (device) => {
  deviceToDelete.value = device;
  showDeleteDialog.value = true;
};

const handleDelete = async () => {
  if (!deviceToDelete.value) return;
  isDeleting.value = true;
  try {
    await deviceStore.deleteDevice(deviceToDelete.value.device_id);
    showDeleteDialog.value = false;
  } finally {
    isDeleting.value = false;
    deviceToDelete.value = null;
  }
};

const handleRowClick = (device) => {
  emit('navigate', 'device-detail', { id: device.device_id });
};

const unassignedDeviceCount = computed(() => {
  if (!devices.value) return 0;
  return devices.value.filter(d => !d.user_id || d.user_id === 'tenant-vendor-1').length;
});

const formatDate = (dateStr) => {
  if (!dateStr) return '---';
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat('en-GB', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date);
};
</script>

<template>
  <div class="device-list">
    <PageHeader 
      title="IoT Infrastructure" 
      description="Monitor and manage system-wide device health and assignments"
    >
      <template #stats>
        <div class="stat-item">
          <i class="ri-cpu-line stat-icon"></i>
          <span><span class="stat-value">{{ devices.length }}</span> devices</span>
        </div>
        <div class="stat-item">
          <i class="ri-link-unlink stat-icon"></i>
          <span><span class="stat-value">{{ unassignedDeviceCount }}</span> unassigned</span>
        </div>
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
        <BaseButton icon="ri-add-line" @click="$emit('navigate', 'device-add')">
          Register Device
        </BaseButton>
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

      <template #cell-last_heartbeat="{ item }">
        <span v-if="item.last_heartbeat" class="connection-time">
          {{ formatDate(item.last_heartbeat) }}
        </span>
        <span v-else class="text-muted">Never</span>
      </template>
      
      <template #cell-site_name="{ item }">
        <div class="assignment-info">
          <div v-if="item.user_name && item.user_id !== 'tenant-vendor-1'" class="User-name">{{ item.user_name }}</div>
          <div v-if="item.site_name" class="site-name text-muted">{{ item.site_name }}</div>
          <div v-if="(!item.user_name || item.user_id === 'tenant-vendor-1') && !item.site_name" class="text-muted">Unassigned</div>
        </div>
      </template>

      <template #cell-actions="{ item }">
        <div class="action-buttons-table">
          <BaseButton variant="ghost" size="sm" @click.stop="handleRowClick(item)">
            Manage
          </BaseButton>
          <BaseButton variant="ghost" size="sm" class="delete-btn" @click.stop="confirmDelete(item)">
            Delete
          </BaseButton>
        </div>
      </template>
    </DataTable>

    <ConfirmDialog
      :show="showDeleteDialog"
      title="Decommission Device"
      :description="`Are you sure you want to remove device ${deviceToDelete?.device_id}? This will disable any active data streams.`"
      confirm-label="Delete Device"
      variant="danger"
      :loading="isDeleting"
      @confirm="handleDelete"
      @cancel="showDeleteDialog = false"
    />
  </div>
</template>

<style scoped>
.assignment-info {
  display: flex;
  flex-direction: column;
}
.User-name {
  font-weight: 500;
  color: var(--color-text-primary);
}
.site-name {
  font-size: 12px;
}
.battery-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.battery-info.low {
  color: var(--color-danger);
  font-weight: 500;
}

.action-buttons-table {
  display: flex;
  gap: 8px;
}

.delete-btn:hover {
  color: var(--color-danger);
}

.text-muted {
  color: var(--color-text-muted, #999);
  font-style: italic;
}
</style>



