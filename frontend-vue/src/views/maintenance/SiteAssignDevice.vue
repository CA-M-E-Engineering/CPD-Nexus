<template>
  <MaintenanceLayout
    :title="pageTitle"
    :description="pageDescription"
    :action-label="pageActionLabel"
    :loading="isSaving"
    @back="handleBack"
    @action="handleSave"
  >
    <template #list>
      <div v-if="isLoading" class="loading-padding">
        <p>Fetching device registry...</p>
      </div>
      <div v-else class="selection-grid">
        <div v-if="!hasSiteContext" class="site-picker">
          <label class="site-label">Select Site</label>
          <select v-model="selectedSiteId" class="site-select">
            <option value="" disabled>Select site...</option>
            <option v-for="site in sites" :key="site.id" :value="site.id">
              {{ site.name }}
            </option>
          </select>
        </div>
        <div v-if="!activeSiteId" class="loading-padding">
          <p>Select a site to start mapping devices.</p>
        </div>
        <div 
          v-else
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
            <div class="selection-name">{{ device.device_id }}</div>
            <div class="selection-meta">{{ device.model }} • {{ device.status }}</div>
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
import { ref, onMounted, computed, watch } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import MaintenanceLayout from '../../components/ui/MaintenanceLayout.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';

const props = defineProps({
  id: [Number, String] // Site ID
});

const emit = defineEmits(['navigate']);

const isLoading = ref(true);
const isSaving = ref(false);
const allDevices = ref([]);
const initialSelectedDevices = ref([]);
const selectedDevices = ref([]);
const sites = ref([]);
const selectedSiteId = ref('');

const hasSiteContext = computed(() => props.id !== undefined && props.id !== null && String(props.id) !== '');
const activeSiteId = computed(() => (hasSiteContext.value ? String(props.id) : selectedSiteId.value));
const pageTitle = computed(() => (hasSiteContext.value ? 'Device Assignment' : 'Request Deployment'));
const pageDescription = computed(() => (
  hasSiteContext.value
    ? 'Allocate IoT devices to this site. Devices checked below are assigned to the site.'
    : 'Select a site and map the devices to deploy.'
));
const pageActionLabel = computed(() => (hasSiteContext.value ? 'Save Device Allocation' : 'Save Deployment'));

const handleBack = () => {
  if (hasSiteContext.value) {
    emit('navigate', 'site-detail', { id: props.id });
  } else {
    emit('navigate', 'devices');
  }
};

const setSelectedForSite = (siteId) => {
  if (!siteId) {
    initialSelectedDevices.value = [];
    selectedDevices.value = [];
    return;
  }
  const assigned = allDevices.value
    .filter(device => String(device.site_id) === String(siteId))
    .map(device => device.device_id);
  initialSelectedDevices.value = assigned;
  selectedDevices.value = [...assigned];
};

const fetchData = async () => {
  isLoading.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let userId = 'User-client-1';
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            userId = user.user_id || user.id;
        } catch (e) {
            console.error("Failed to parse auth_user", e);
        }
    }
    const [sitesData, devicesData] = await Promise.all([
      api.getSites({ user_id: userId }),
      api.getDevices({ user_id: userId })
    ]);
    sites.value = sitesData || [];
    allDevices.value = devicesData || [];
    if (hasSiteContext.value) {
      selectedSiteId.value = String(props.id);
    }
    setSelectedForSite(activeSiteId.value);
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchData);

watch(selectedSiteId, (newId) => {
  if (!hasSiteContext.value) {
    setSelectedForSite(newId);
  }
});

const toggleSelection = (id) => {
  const idx = selectedDevices.value.indexOf(id);
  if (idx > -1) selectedDevices.value.splice(idx, 1);
  else selectedDevices.value.push(id);
};

const handleSave = async () => {
  if (!activeSiteId.value) return;
  isSaving.value = true;
  try {
    await api.assignDevicesToSite(activeSiteId.value, selectedDevices.value);
    notification.success('Site device allocation saved');
    if (hasSiteContext.value) {
      emit('navigate', 'site-detail', { id: props.id });
    } else {
      emit('navigate', 'devices');
    }
  } catch (err) {
    console.error('Failed to assign devices', err);
    notification.error(err.message || 'Failed to update site mapping');
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

.site-picker {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 16px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.site-label {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.site-select {
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-border);
  background: var(--color-bg);
  color: var(--color-text-primary);
}
</style>

