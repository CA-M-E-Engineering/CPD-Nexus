<template>
  <MaintenanceLayout
    title="Hardware Allocation"
    description="Allocate IoT devices to this tenant. Devices checked are currently assigned to the organization."
    action-label="Update Hardware Allocation"
    :loading="isSaving"
    @back="$emit('navigate', 'tenant-detail', { id: props.id })"
    @action="handleSave"
  >
    <template #list>
      <div v-if="isLoading" class="loading-padding">
        <p>Fetching hardware registry...</p>
      </div>
      <div v-else class="selection-grid">
        <div 
          v-for="device in allDevices" 
          :key="device.device_id"
          class="selection-row"
          :class="{ selected: selectedDevices.includes(device.device_id) }"
          @click="toggleSelection(device.device_id)"
        >
          <div class="selection-checkbox">
            <i v-if="selectedDevices.includes(device.device_id)" class="ri-checkbox-fill"></i>
            <i v-else class="ri-checkbox-blank-line"></i>
          </div>
            <div class="selection-info">
            <div class="selection-name">{{ device.sn }}</div>
            <div class="selection-meta">{{ device.model }} • <span :class="{'text-success': device.status === 'online', 'text-danger': device.status === 'offline'}">{{ device.status }}</span></div>
          </div>
          <div class="selection-status">
            <BaseBadge v-if="initialSelectedDevices.includes(device.device_id)" type="info">Currently Allocated</BaseBadge>
          </div>
        </div>
      </div>
    </template>
  </MaintenanceLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import MaintenanceLayout from '../../components/ui/MaintenanceLayout.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';

const props = defineProps({
  id: [Number, String] // Tenant ID
});

const emit = defineEmits(['navigate']);

const isLoading = ref(true);
const isSaving = ref(false);
const allDevices = ref([]);
const initialSelectedDevices = ref([]);
const selectedDevices = ref([]);

const fetchData = async () => {
  isLoading.value = true;
  try {
    const [devices, tenantDevices] = await Promise.all([
      api.getDevices(),
      api.getDevices({ tenant_id: props.id })
    ]);
    
    allDevices.value = devices;
    
    // Extract IDs from tenant devices
    const assignedIds = tenantDevices.map(d => d.device_id);
    initialSelectedDevices.value = assignedIds;
    selectedDevices.value = [...assignedIds];
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchData);

const toggleSelection = (id) => {
  const idx = selectedDevices.value.indexOf(id);
  if (idx > -1) selectedDevices.value.splice(idx, 1);
  else selectedDevices.value.push(id);
};

const handleSave = async () => {
  isSaving.value = true;
  try {
    await api.bulkAssign(props.id, selectedDevices.value);
    notification.success('Hardware allocation updated');
    emit('navigate', 'tenant-detail', { id: props.id });
  } catch (err) {
    console.error('Failed to update allocation', err);
    notification.error(err.message || 'Failed to sync hardware allocation');
  } finally {
    isSaving.value = false;
  }
};
</script>

<style scoped>
.selection-grid {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.selection-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
}

.selection-row:hover {
  border-color: var(--color-accent);
  background: var(--color-surface-hover);
}

.selection-row.selected {
  border-color: var(--color-accent);
  background: rgba(59, 130, 246, 0.05);
}

.selection-checkbox {
  font-size: 20px;
  color: var(--color-text-muted);
}

.selected .selection-checkbox {
  color: var(--color-accent);
}

.selection-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.selection-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.selection-meta {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.loading-padding {
  padding: 32px;
  text-align: center;
  color: var(--color-text-secondary);
}
</style>

