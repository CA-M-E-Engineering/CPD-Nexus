<template>
  <MaintenanceLayout
    title="Worker Assignment"
    description="Select workers to assign to this project. Pre-selected workers are currently assigned."
    action-label="Save Worker Assignments"
    :loading="isSaving"
    @back="$emit('navigate', 'project-detail', { id: props.id })"
    @action="handleSave"
  >
    <template #list>
      <div v-if="isLoading" class="loading-padding">
        <p>Fetching worker directory...</p>
      </div>
      <div v-else class="selection-grid">
        <div 
          v-for="worker in allWorkers" 
          :key="worker.user_id"
          class="selection-row"
          :class="{ selected: selectedWorkers.includes(worker.user_id) }"
          @click="toggleSelection(worker.user_id)"
        >
          <div class="selection-checkbox">
            <i v-if="selectedWorkers.includes(worker.user_id)" class="ri-checkbox-fill"></i>
            <i v-else class="ri-checkbox-blank-line"></i>
          </div>
          <div class="selection-info">
            <div class="selection-name">{{ worker.name }}</div>
            <div class="selection-meta">
              {{ worker.role === 'pic' ? 'PIC' : 'Worker' }} • {{ worker.trade_code || 'General' }}
              <div v-if="worker.site_name" class="current-loc">
                Currently at: {{ worker.site_name }} ({{ worker.tenant_name }})
              </div>
              <div v-else class="current-loc unassigned">
                Currently Unassigned
              </div>
            </div>
          </div>
          <div class="selection-status">
            <BaseBadge v-if="initialSelectedWorkers.includes(worker.user_id)" type="info">Currently Assigned</BaseBadge>
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
  id: [Number, String] // Project ID
});

const emit = defineEmits(['navigate']);

const isLoading = ref(true);
const isSaving = ref(false);
const allWorkers = ref([]);
const initialSelectedWorkers = ref([]); 
const selectedWorkers = ref([]);

const fetchData = async () => {
  isLoading.value = true;
  try {
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
    const data = await api.getWorkers({ tenant_id: tenantId });
    allWorkers.value = data || [];
    
    // Auto-detect workers currently assigned to this project ID
    initialSelectedWorkers.value = allWorkers.value
        .filter(w => w.current_project_id === props.id)
        .map(w => w.user_id);
        
    selectedWorkers.value = [...initialSelectedWorkers.value];
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchData);

const toggleSelection = (id) => {
  const idx = selectedWorkers.value.indexOf(id);
  if (idx > -1) selectedWorkers.value.splice(idx, 1);
  else selectedWorkers.value.push(id);
};

const handleSave = async () => {
  isSaving.value = true;
  try {
    await api.assignWorkersToProject(props.id, selectedWorkers.value);
    notification.success('Worker assignments updated');
    emit('navigate', 'project-detail', { id: props.id });
  } catch (err) {
    console.error('Failed to assign workers', err);
    notification.error(err.message || 'Failed to update project assignments');
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
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.current-loc {
  font-size: 11px;
  color: var(--color-accent);
  opacity: 0.8;
}

.unassigned {
  color: var(--color-text-muted);
}

.loading-padding {
  padding: 32px;
  text-align: center;
  color: var(--color-text-secondary);
}
</style>

