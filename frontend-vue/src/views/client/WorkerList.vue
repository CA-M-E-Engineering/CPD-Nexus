<script setup>
import { ref, computed, onMounted } from 'vue';

import { api } from '../../services/api.js';

import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseFilterChip from '../../components/ui/BaseFilterChip.vue';
import PageHeader from '../../components/ui/PageHeader.vue';
import TableToolbar from '../../components/ui/TableToolbar.vue';
import ConfirmDialog from '../../components/ui/ConfirmDialog.vue';


const activeFilter = ref('All');
const filters = ['All', 'active', 'inactive'];


const columns = [
  { key: 'name', label: 'Name', size: 'lg', bold: true },
  { key: 'role', label: 'Role', size: 'md', muted: true },
  { key: 'person_trade', label: 'Trade', size: 'sm', muted: true },
  { key: 'project_name', label: 'Project Name', size: 'md' },
  { key: 'status', label: 'Status', size: 'sm' },
  { key: 'actions', label: 'Actions', width: '80px' }
];

const workers = ref([]);
const totalWorkers = ref(0);
const unassignedCount = ref(0);
const isLoading = ref(false);
const error = ref(null);

const fetchWorkers = async () => {
  console.log('[WorkerList] Fetching workers...');
  isLoading.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let userId = null;
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            userId = user.user_id || user.id;
        } catch (e) {
            console.error('[WorkerList] Failed to parse auth_user', e);
        }
    }
    
    if (!userId) {
        console.warn("[WorkerList] No User ID found");
        isLoading.value = false;
        return;
    }
    
    console.log('[WorkerList] Requesting workers for User:', userId);
    const data = await api.getWorkers({ user_id: userId });
    console.log('[WorkerList] Received workers:', data);
    workers.value = data || [];

    // Calculate Stats
    totalWorkers.value = workers.value.length;
    console.log('[WorkerList] Total workers set to:', totalWorkers.value);
    unassignedCount.value = workers.value.filter(w => !w.current_project_id).length;

  } catch (err) {
    console.error('[WorkerList] Fetch Error:', err);
    error.value = 'Failed to load workers';
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchWorkers);

const filteredWorkers = computed(() => {
  return workers.value.filter(worker => {
    return activeFilter.value === 'All' || worker.status === activeFilter.value;
  });
});

const emit = defineEmits(['navigate']);

const handleRowClick = (worker) => {
  emit('navigate', 'worker-detail', { id: worker.worker_id });
};

const handleEdit = (worker) => {
  emit('navigate', 'worker-add', { id: worker.worker_id, mode: 'edit' });
};

const handleExport = async () => {
  isLoading.value = true;
  await api.simulateExport('Workers');
  isLoading.value = false;
};

// Delete logic
const showDeleteDialog = ref(false);
const workerToDelete = ref(null);
const isDeleting = ref(false);

const confirmDelete = (item) => {
  workerToDelete.value = item;
  showDeleteDialog.value = true;
};

const deleteWorker = async () => {
  if (!workerToDelete.value) return;
  isDeleting.value = true;
  try {
    await api.deleteWorker(workerToDelete.value.worker_id);
    workers.value = workers.value.filter(w => w.worker_id !== workerToDelete.value.worker_id);
    showDeleteDialog.value = false;
  } catch (err) {
     console.error('Failed to delete worker', err);
  } finally {
    isDeleting.value = false;
    workerToDelete.value = null;
  }
};

</script>

<template>
  <div class="worker-list">
    <PageHeader 
      title="Worker Management" 
      description="Manage your workforce across all sites"
    >

      <template #stats>
        <div class="stat-item">
          <i class="ri-group-line stat-icon"></i>
          <span><span class="stat-value">{{ totalWorkers }}</span> total</span>
        </div>
        <div class="stat-item">
          <i class="ri-user-unfollow-line stat-icon"></i>
          <span><span class="stat-value">{{ unassignedCount }}</span> unassigned</span>
        </div>
      </template>

      <template #actions>
        <BaseButton variant="secondary" :loading="isLoading" @click="handleExport">Export</BaseButton>
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
        <BaseButton 
          icon="ri-add-line" 
          @click="$emit('navigate', 'worker-add')"
        >
          Add Worker
        </BaseButton>
      </template>
    </TableToolbar>



    <DataTable :loading="isLoading" :columns="columns" :data="filteredWorkers">
      <template #cell-name="{ item }">
        <div class="clickable-cell" @click="handleRowClick(item)">
          <strong>{{ item.name }}</strong>
        </div>
      </template>


      <template #cell-role="{ item }">
        {{ 
          item.role === 'pic' ? 'PIC (Person In Charge)' : 
          item.role === 'manager' ? 'Manager' : 
          item.role === 'worker' ? 'Worker' : 
          item.role 
        }}
      </template>

      <template #cell-project_name="{ item }">
        <span v-if="item.project_name">{{ item.project_name }}</span>
        <span v-else class="text-unassigned">Not Assigned</span>
      </template>

      <template #cell-status="{ item }">
        <BaseBadge :type="item.status === 'active' ? 'success' : 'inactive'">
          {{ item.status.charAt(0).toUpperCase() + item.status.slice(1) }}
        </BaseBadge>
      </template>

      <template #cell-actions="{ item }">
        <div class="action-buttons">
          <BaseButton variant="ghost" size="sm" @click="handleEdit(item)">
            <i class="ri-edit-line"></i> Edit
          </BaseButton>
          <BaseButton variant="ghost" size="sm" class="delete-btn" @click="confirmDelete(item)">
            <i class="ri-delete-bin-line"></i> Delete
          </BaseButton>
        </div>
      </template>
    </DataTable>

    <ConfirmDialog
      :show="showDeleteDialog"
      :loading="isDeleting"
      title="Delete Worker"
      :description="`Are you sure you want to delete ${workerToDelete?.name}? This action cannot be undone.`"
      @confirm="deleteWorker"
      @cancel="showDeleteDialog = false"
    />
  </div>
</template>


<style scoped>
.stat-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text-secondary);
}

.stat-cell i {
  color: var(--color-accent);
  font-size: 16px;
}

.text-unassigned {
  color: var(--color-text-secondary);
  font-style: italic;
  font-size: 13px;
}
</style>

