<script setup>
import { ref, computed, onMounted } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import TableToolbar from '../../components/ui/TableToolbar.vue';

const selectedSite = ref('');
const selectedDate = ref('');
const selectedStatus = ref('');

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

const filteredAttendance = computed(() => {
  if (!selectedStatus.value) return attendance.value;
  return attendance.value.filter(a => {
    const s = (a.status || 'pending').toLowerCase();
    // Some status might be 'success' instead of submitted depending on old data, handle if needed, but user asked for pending/submitted
    return s === selectedStatus.value || (selectedStatus.value === 'submitted' && s === 'success');
  });
});

const fetchData = async () => {
  isLoading.value = true;
  try {
    const params = {};
    if (selectedSite.value) params.site_id = selectedSite.value;
    if (selectedDate.value) params.date = selectedDate.value;
    const savedUser = localStorage.getItem('auth_user');
    let userId = null;
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            userId = user.user_id || user.id;
        } catch (e) {
            console.error("Failed to parse auth_user", e);
        }
    }

    if (!userId) {
        isLoading.value = false;
        return;
    }
    
    // Add user_id to attendance params
    params.user_id = userId;

    const [sitesData, attendanceData] = await Promise.all([
      api.getSites({ user_id: userId }),
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
        <div class="filter-bar">
          <div class="filter-item">
            <i class="ri-map-pin-line filter-icon"></i>
            <select v-model="selectedSite" @change="handleFilter" class="filter-input">
              <option value="">All Sites</option>
              <option v-for="site in sites" :key="site.site_id" :value="site.site_id">
                {{ site.site_name }}
              </option>
            </select>
          </div>
          <div class="filter-item">
            <i class="ri-calendar-line filter-icon"></i>
            <input type="date" v-model="selectedDate" @change="handleFilter" class="filter-input" />
          </div>
          <div class="segmented-control">
            <button 
              class="seg-btn" 
              :class="{ active: selectedStatus === '' }" 
              @click="selectedStatus = ''"
            >All</button>
            <button 
              class="seg-btn" 
              :class="{ active: selectedStatus === 'pending' }" 
              @click="selectedStatus = 'pending'"
            >Pending</button>
            <button 
              class="seg-btn" 
              :class="{ active: selectedStatus === 'submitted' }" 
              @click="selectedStatus = 'submitted'"
            >Submitted</button>
          </div>
        </div>
      </template>
      <template #right>
        <BaseButton variant="secondary" icon="ri-download-line" :loading="isExporting" @click="handleExport">
          Export Logs
        </BaseButton>
      </template>
    </TableToolbar>

    <DataTable :loading="isLoading" :columns="columns" :data="filteredAttendance">
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

.filter-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  background: var(--color-surface);
  padding: 6px;
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
}

.filter-item {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: 6px 12px;
  transition: all 0.2s ease;
}

.filter-item:focus-within {
  border-color: var(--color-accent);
}

.filter-icon {
  color: var(--color-text-muted);
  font-size: 16px;
}

.filter-input {
  background: transparent;
  border: none;
  color: var(--color-text-primary);
  font-size: 13px;
  outline: none;
  min-width: 140px;
  cursor: pointer;
}

.filter-input[type="date"]::-webkit-calendar-picker-indicator {
  cursor: pointer;
  opacity: 0.6;
}

.segmented-control {
  display: flex;
  background: var(--color-bg);
  border-radius: var(--radius-sm);
  padding: 4px;
  border: 1px solid var(--color-border);
}

.seg-btn {
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  padding: 6px 16px;
  font-size: 13px;
  font-weight: 500;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.seg-btn:hover {
  color: var(--color-text-primary);
}

.seg-btn.active {
  background: var(--color-surface-hover);
  color: var(--color-accent);
  box-shadow: var(--shadow-sm);
}
</style>


