<script setup>
import { ref, onMounted, computed } from 'vue';
import { api } from '../../services/api.js';
import PageHeader from '../../components/ui/PageHeader.vue';
import DetailCard from '../../components/ui/DetailCard.vue';
import BaseTabs from '../../components/ui/BaseTabs.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';

const props = defineProps({
  id: [Number, String]
});

const emit = defineEmits(['navigate']);

const activeTab = ref('Workers');
const site = ref(null);
const isLoading = ref(true);
const loadingSubData = ref(false);

const tabs = [
  { id: 'Workers', label: 'Assigned Workers' },
  { id: 'Devices', label: 'IoT Devices' }
];

const fetchSite = async () => {
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
    const response = await api.getSiteById(props.id, { tenant_id: tenantId });
    site.value = typeof response === 'string' ? JSON.parse(response) : response;
  } catch (err) {
    console.error('Failed to fetch site:', err);
  } finally {
    isLoading.value = false;
  }
};

const siteInfo = computed(() => [
  { label: 'Site Name', value: site.value?.site_name || '---' },
  { label: 'Location', value: site.value?.location || '---' },
  { label: 'Coordinates', value: site.value?.lat && site.value?.lng ? `${site.value.lat}, ${site.value.lng}` : '---' },
]);

const workerColumns = [
  { key: 'name', label: 'Worker Name' },
  { key: 'role', label: 'Role' },
  { key: 'status', label: 'Status' },
  { key: 'actions', label: 'Actions', width: '100px' }
];

const assignedWorkers = ref([]);
const assignedDevices = ref([]);

onMounted(async () => {
    await fetchSite();

    const savedUser = localStorage.getItem('auth_user');
    let tenantId = 'tenant-client-1';
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            tenantId = user.tenant_id || user.id;
        } catch (e) {
            console.error('Failed to parse auth_user', e);
        }
    }

    loadingSubData.value = true;
    try {
        const [workersData, devicesData] = await Promise.all([
          api.getWorkers({ tenant_id: tenantId, site_id: props.id }),
          api.getDevices({ tenant_id: tenantId, site_id: props.id })
        ]);
        
        assignedWorkers.value = workersData || [];
        assignedDevices.value = devicesData || [];
    } finally {
        loadingSubData.value = false;
    }
});

const deviceColumns = [
  { key: 'device_id', label: 'Device ID', bold: true },
  { key: 'model', label: 'Type' },
  { key: 'status', label: 'Status' },
  { key: 'battery', label: 'Battery' }
];

const handleManageDevice = (device) => {
  emit('navigate', 'device-detail', { id: device.device_id });
};

const resourceStats = computed(() => [
  { label: 'Total Workers', value: site.value?.worker_count || '0' },
  { label: 'Active Today', value: Math.floor(site.value?.worker_count * 0.7) || '0' },
  { label: 'Devices', value: site.value?.device_count || '0' }
]);

const handleEdit = () => {
  emit('navigate', 'site-add', { id: props.id, mode: 'edit' });
};

const handleManageWorker = (worker) => {
  emit('navigate', 'worker-detail', { id: worker.user_id });
};
</script>

<template>
  <div class="site-detail">
    <PageHeader 
      :title="site?.site_name || 'Loading...'" 
      description="View site details, devices, and workers"
      variant="detail"
    >
      <template #toolbar-left>
        <BaseButton variant="ghost" size="sm" @click="$emit('navigate', 'sites')">
          <template #icon><i class="ri-arrow-left-line"></i></template>
          Back to Sites
        </BaseButton>
      </template>
      <template #toolbar-right>
        <BaseButton variant="secondary" size="sm" @click="handleEdit">Edit Site</BaseButton>
        <BaseButton variant="secondary" size="sm" @click="$emit('navigate', 'site-assign-device', { id: props.id })">Assign Device</BaseButton>
        <BaseButton size="sm" @click="$emit('navigate', 'site-assign-project', { id: props.id })">Assign Project</BaseButton>
      </template>
    </PageHeader>

    <div v-if="isLoading" class="loading-state">
      <p>Loading site details...</p>
    </div>

    <div v-else-if="site" class="content-body">
      <div class="detail-grid">
          <DetailCard 
            title="Site Information" 
            :rows="siteInfo"
          />
          <DetailCard 
            title="Resources" 
            :rows="resourceStats"
          />
      </div>

      <BaseTabs v-model="activeTab" :tabs="tabs" />

      <div v-show="activeTab === 'Workers'" class="tab-content">
        <DataTable :loading="loadingSubData" :columns="workerColumns" :data="assignedWorkers">
          <template #cell-name="{ item }">
              <strong>{{ item.name }}</strong>
          </template>
          <template #cell-status="{ item }">
            <BaseBadge :type="item.status === 'active' ? 'success' : 'inactive'">{{ item.status }}</BaseBadge>
          </template>
          <template #cell-actions="{ item }">
            <BaseButton variant="secondary" size="sm" @click="handleManageWorker(item)">Manage</BaseButton>
          </template>
        </DataTable>
      </div>

      <div v-show="activeTab === 'Devices'" class="tab-content">
          <DataTable :loading="loadingSubData" :columns="deviceColumns" :data="assignedDevices" row-clickable @row-click="handleManageDevice">
            <template #cell-status="{ item }">
              <BaseBadge :type="item.status === 'online' ? 'success' : 'inactive'">
                {{ item.status.toUpperCase() }}
              </BaseBadge>
            </template>
            <template #cell-battery="{ item }">
              <span>{{ item.battery }}%</span>
            </template>
          </DataTable>
          <div v-if="assignedDevices.length === 0" class="empty-state">
              <p class="empty-description">No devices currently assigned to this site.</p>
          </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.loading-state {
  padding: 64px;
  text-align: center;
  color: var(--color-text-secondary);
}

.tab-content {
  margin-top: 24px;
}
</style>

