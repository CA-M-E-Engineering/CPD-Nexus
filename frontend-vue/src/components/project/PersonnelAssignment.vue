<script setup>
import { ref, onMounted, computed } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import BaseButton from '../ui/BaseButton.vue';
import DataTable from '../ui/DataTable.vue';
import BaseBadge from '../ui/BaseBadge.vue';
import BaseModal from '../ui/BaseModal.vue';
import BaseInput from '../ui/BaseInput.vue';

const props = defineProps({
  projectId: { type: String, required: true },
  userId: { type: String, required: true }
});

const allWorkers = ref([]);
const isAdding = ref(false);
const searchQuery = ref('');
const isSaving = ref(false);

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'role', label: 'Role' },
  { key: 'company_name', label: 'Company' },
  { key: 'actions', label: 'Actions', width: '80px' }
];

const modalColumns = [
  { key: 'name', label: 'Name' },
  { key: 'role', label: 'Role' },
  { key: 'company_name', label: 'Company' },
  { key: 'status', label: 'Current Assignment' },
  { key: 'actions', label: '', width: '60px' }
];

const fetchWorkers = async () => {
  try {
    const response = await api.getWorkers({ user_id: props.userId });
    allWorkers.value = typeof response === 'string' ? JSON.parse(response) : response;
  } catch (err) {
    console.error('Failed to fetch workers:', err);
  }
};

const assignedPersonnel = computed(() => {
  return allWorkers.value.filter(w => String(w.current_project_id) === String(props.projectId));
});

const availablePersonnel = computed(() => {
  return allWorkers.value
    .filter(w => String(w.current_project_id) !== String(props.projectId))
    .filter(w => {
      if (!searchQuery.value) return true;
      const term = searchQuery.value.toLowerCase();
      return w.name.toLowerCase().includes(term) || (w.company_name && w.company_name.toLowerCase().includes(term));
    });
});

const handleAssign = async (workerId) => {
  isSaving.value = true;
  try {
    // Collect all currently assigned + the new one
    const workerIds = assignedPersonnel.value.map(w => w.worker_id);
    if (workerId) workerIds.push(workerId);

    await api.assignWorkersToProject(props.projectId, workerIds);
    notification.success('Personnel updated');
    await fetchWorkers();
  } catch (err) {
    notification.error('Failed to update personnel');
  } finally {
    isSaving.value = false;
  }
};

const handleRemove = async (worker) => {
  isSaving.value = true;
  try {
    const workerIds = assignedPersonnel.value
      .filter(w => w.worker_id !== worker.worker_id)
      .map(w => w.worker_id);

    await api.assignWorkersToProject(props.projectId, workerIds);
    notification.success(`${worker.name} removed from project`);
    await fetchWorkers();
  } catch (err) {
    notification.error('Failed to remove personnel');
  } finally {
    isSaving.value = false;
  }
};

onMounted(fetchWorkers);

const getRoleBadge = (role) => {
  switch (role?.toLowerCase()) {
    case 'pic': return 'warning';
    case 'manager': return 'info';
    default: return 'success';
  }
};
</script>

<template>
  <div class="personnel-assignment">
    <div class="header-row">
      <h3 class="panel-title">Project Personnel</h3>
      <BaseButton size="sm" icon="ri-user-add-line" @click="isAdding = true">Add Personnel</BaseButton>
    </div>

    <DataTable :columns="columns" :data="assignedPersonnel">
      <template #cell-name="{ item }">
        <div class="worker-name-cell">
          <strong>{{ item.name }}</strong>
          <span v-if="item.fin_nric" class="fin-text">{{ item.fin_nric }}</span>
        </div>
      </template>
      <template #cell-role="{ item }">
        <BaseBadge :type="getRoleBadge(item.role)">{{ item.role?.toUpperCase() }}</BaseBadge>
      </template>
      <template #cell-actions="{ item }">
        <BaseButton 
          variant="ghost" 
          size="sm" 
          icon="ri-delete-bin-line" 
          class="delete-btn"
          :loading="isSaving"
          @click="handleRemove(item)"
        />
      </template>
    </DataTable>

    <div v-if="assignedPersonnel.length === 0" class="empty-state">
      <p>No workers assigned to this project yet.</p>
    </div>

    <BaseModal v-model="isAdding" title="Assign Personnel" size="lg">
      <div class="modal-search">
        <BaseInput v-model="searchQuery" placeholder="Search by name or company..." icon="ri-search-line" />
      </div>

      <DataTable :columns="modalColumns" :data="availablePersonnel" class="modal-table">
        <template #cell-status="{ item }">
           <span v-if="item.project_title" class="assigned-status">
             Assigned to: <strong>{{ item.project_title }}</strong>
           </span>
           <span v-else class="unassigned-status">Unassigned</span>
        </template>
        <template #cell-actions="{ item }">
          <BaseButton 
            variant="secondary" 
            size="sm" 
            icon="ri-add-line"
            :loading="isSaving"
            @click="handleAssign(item.worker_id)"
          >
            Add
          </BaseButton>
        </template>
      </DataTable>
    </BaseModal>
  </div>
</template>

<style scoped>
.personnel-assignment {
  margin-top: 32px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 24px;
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.panel-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.worker-name-cell {
  display: flex;
  flex-direction: column;
}

.fin-text {
  font-size: 11px;
  color: var(--color-text-secondary);
}

.delete-btn :hover {
  color: #ef4444;
}

.empty-state {
  padding: 32px;
  text-align: center;
  color: var(--color-text-secondary);
  font-style: italic;
  font-size: 14px;
}

.modal-search {
  margin-bottom: 16px;
}

.modal-table {
  max-height: 400px;
  overflow-y: auto;
}

.assigned-status {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.unassigned-status {
  font-size: 12px;
  color: var(--color-success);
}
</style>
