<script setup>
import { ref, onMounted } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import TableToolbar from '../../components/ui/TableToolbar.vue';

const selectedSite = ref('');
const selectedDate = ref('');

const sites = ref([]);
const attendance = ref([]);
const isLoading = ref(false);
const isExporting = ref(false);

const columns = [
  { key: 'worker_name', label: 'Worker', size: 'md', bold: true },
  { key: 'site_name', label: 'Site', size: 'md', muted: true },
  { key: 'time_in', label: 'Time In', size: 'md', muted: true },
  { key: 'time_out', label: 'Time Out', size: 'md', muted: true },
  { key: 'status', label: 'Status', size: 'sm' }
];

const fetchData = async () => {
  isLoading.value = true;
  try {
    const params = {};
    if (selectedSite.value) params.site_id = selectedSite.value;
    if (selectedDate.value) params.date = selectedDate.value;
    const savedUser = localStorage.getItem('auth_user');
    let tenantId = null;
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            tenantId = user.tenant_id || user.id;
        } catch (e) {
            console.error("Failed to parse auth_user", e);
        }
    }

    if (!tenantId) {
        isLoading.value = false;
        return;
    }
    
    // Add tenant_id to attendance params
    params.tenant_id = tenantId;

    const [sitesData, attendanceData] = await Promise.all([
      api.getSites({ tenant_id: tenantId }),
      api.getAttendance(params)
    ]);
    
    sites.value = sitesData || [];
    attendance.value = attendanceData || [];
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchData);

const handleFilter = () => {
  fetchData();
};

const handleExport = async () => {
  isExporting.value = true;
  await api.simulateExport('Attendance Logs');
  isExporting.value = false;
  notification.success('Attendance logs exported successfully!');
};

defineEmits(['navigate']);
</script>

<template>
  <div class="attendance-list">
    <PageHeader 
      title="Attendance Records" 
      description="View and audit worker check-in/out logs"
    />

    <div class="sync-notice">
      <i class="ri-information-line"></i>
      <span>CPD data is synchronized daily at 10:00 PM SGT</span>
    </div>

    <TableToolbar>
      <template #left>
        <div class="filter-group">
          <label>Site</label>
          <select v-model="selectedSite" @change="handleFilter" class="filter-select">
            <option value="">All Sites</option>
            <option v-for="site in sites" :key="site.site_id" :value="site.site_id">
              {{ site.site_name }}
            </option>
          </select>
        </div>
        <div class="filter-group">
          <label>Date</label>
          <input type="date" v-model="selectedDate" @change="handleFilter" class="filter-date" />
        </div>
      </template>
      <template #right>
        <BaseButton variant="secondary" icon="ri-download-line" :loading="isExporting" @click="handleExport">
          Export Logs
        </BaseButton>
      </template>
    </TableToolbar>

    <DataTable :loading="isLoading" :columns="columns" :data="attendance">
      <template #cell-status="{ item }">
        <BaseBadge :type="item.status === 'success' ? 'success' : 'warning'">
          {{ (item.status || 'PENDING').toUpperCase() }}
        </BaseBadge>
      </template>
    </DataTable>
  </div>
</template>

<style scoped>
.sync-notice {
  margin-bottom: 24px;
  padding: 12px 20px;
  background: rgba(59, 130, 246, 0.1);
  border-left: 4px solid var(--color-accent);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--color-text-primary);
  font-size: 14px;
  backdrop-filter: blur(8px);
}

.sync-notice i {
  color: var(--color-accent);
  font-size: 18px;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.filter-group label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--color-text-muted);
}

.filter-select, .filter-date {
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: 8px 12px;
  color: var(--color-text-primary);
  font-size: 14px;
  outline: none;
  min-width: 160px;
}

.filter-select:focus, .filter-date:focus {
  border-color: var(--color-accent);
}
</style>


