<template>
  <MaintenanceLayout
    title="Change Worker Project"
    :description="`Update the project assignment for ${worker?.name || 'Worker'}`"
    action-label="Update Project Assignment"
    :loading="isSaving"
    @back="$emit('navigate', 'worker-detail', { id: props.id })"
    @action="handleSave"
  >
    <template #list>
      <div v-if="isLoading" class="loading-padding">
        <p>Fetching project directory...</p>
      </div>
      <div v-else-if="projects.length === 0" class="loading-padding">
        <p>No active projects found for this tenant.</p>
        <p><small>Create a project first to assign workers.</small></p>
      </div>
      <div v-else class="selection-grid">
        <div 
          v-for="project in projects" 
          :key="project.project_id"
          class="selection-row"
          :class="{ selected: selectedProjectId === project.project_id }"
          @click="selectedProjectId = project.project_id"
        >
          <div class="selection-checkbox">
            <i v-if="selectedProjectId === project.project_id" class="ri-checkbox-circle-fill"></i>
            <i v-else class="ri-checkbox-blank-circle-line"></i>
          </div>
          <div class="selection-info">
            <div class="selection-name">{{ project.title }}</div>
            <div class="selection-meta">
              {{ project.reference }} • Site: {{ project.site_name || 'N/A' }}
            </div>
          </div>
          <div class="selection-status">
            <BaseBadge v-if="worker?.current_project_id === project.project_id" type="info">Current Project</BaseBadge>
          </div>
        </div>

        <div 
          class="selection-row"
          :class="{ selected: selectedProjectId === '' }"
          @click="selectedProjectId = ''"
        >
          <div class="selection-checkbox">
            <i v-if="selectedProjectId === ''" class="ri-checkbox-circle-fill"></i>
            <i v-else class="ri-checkbox-blank-circle-line"></i>
          </div>
          <div class="selection-info">
            <div class="selection-name">Unassigned</div>
            <div class="selection-meta">Remove worker from all projects</div>
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
  id: [Number, String] // Worker ID
});

const emit = defineEmits(['navigate']);

const isLoading = ref(true);
const isSaving = ref(false);
const worker = ref(null);
const projects = ref([]);
const selectedProjectId = ref(null);

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

    const [workerData, projectsData] = await Promise.all([
      api.getWorkerById(props.id, { tenant_id: tenantId }),
      api.getProjects({ tenant_id: tenantId })
    ]);
    worker.value = workerData;
    projects.value = projectsData || [];
    selectedProjectId.value = workerData.current_project_id || '';
  } catch (err) {
    console.error('Failed to fetch data', err);
    notification.error('Failed to load project information');
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchData);

const handleSave = async () => {
  // Normalize comparisons (handling null/undefined as empty string)
  // Ensure both are treated as strings to avoid type mismatches (e.g. "1" vs 1)
  const currentId = String(worker.value?.current_project_id || '');
  const newId = String(selectedProjectId.value || '');

  console.log('[WorkerAssignProject] Comparing IDs:', { current: currentId, new: newId });

  if (newId === currentId) {
    console.log('[WorkerAssignProject] No change in project selection');
    emit('navigate', 'worker-detail', { id: props.id });
    return;
  }

  isSaving.value = true;
  try {
    console.log('[WorkerAssignProject] Current worker state:', worker.value);
    console.log('[WorkerAssignProject] New project selection:', newId);
    console.log('[WorkerAssignProject] Updating project assignment:', { from: currentId, to: newId });
    
    // We only need to send the field we want to update
    const updateData = {
      current_project_id: newId
    };
    
    console.log('[WorkerAssignProject] Calling updateWorker with payload:', updateData);
    await api.updateWorker(props.id, updateData);
    notification.success('Project assignment updated');
    emit('navigate', 'worker-detail', { id: props.id });
  } catch (err) {
    console.error('[WorkerAssignProject] Failed to update project', err);
    notification.error(err.message || 'Failed to update project assignment');
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
