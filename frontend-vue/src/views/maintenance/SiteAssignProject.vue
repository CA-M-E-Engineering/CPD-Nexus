<template>
  <MaintenanceLayout
    title="Project Assignment"
    description="Choose a project to link with this site. You can toggle projects below."
    action-label="Update Project Linking"
    :loading="isSaving"
    @back="$emit('navigate', 'site-detail', { id: props.id })"
    @action="handleSave"
  >
    <template #list>
      <div v-if="isLoading" class="loading-padding">
        <p>Fetching project directory...</p>
      </div>
      <div v-else class="selection-grid">
        <div 
          v-for="proj in allProjects" 
          :key="proj.project_id"
          class="selection-row"
          :class="{ selected: selectedProjectIds.includes(proj.project_id) }"
          @click="toggleSelection(proj.project_id)"
        >
          <div class="selection-checkbox">
            <i v-if="selectedProjectIds.includes(proj.project_id)" class="ri-checkbox-fill"></i>
            <i v-else class="ri-checkbox-blank-line"></i>
          </div>
          <div class="selection-info">
            <div class="selection-name">{{ proj.title }}</div>
            <div class="selection-meta">
              {{ proj.reference }} • 
              <span v-if="proj.site_id && String(proj.site_id) === String(props.id)" class="text-accent">Assigned to this site</span>
              <span v-else-if="proj.site_name">{{ proj.site_name }}</span>
              <span v-else>Unassigned</span>
            </div>
          </div>
          <div class="selection-status">
            <BaseBadge v-if="initialProjectIds.includes(proj.project_id)" type="info">Currently Linked</BaseBadge>
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
  id: [Number, String] // Site ID
});

const emit = defineEmits(['navigate']);

const isLoading = ref(true);
const isSaving = ref(false);
const allProjects = ref([]);
const initialProjectIds = ref([]);
const selectedProjectIds = ref([]);

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

    allProjects.value = await api.getProjects({ tenant_id: tenantId });
    
    // Find projects already linked to this site
    const linked = allProjects.value
      .filter(p => p.site_id && String(p.site_id) === String(props.id))
      .map(p => p.project_id);
      
    initialProjectIds.value = linked;
    selectedProjectIds.value = [...linked];
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchData);

const toggleSelection = (id) => {
  const idx = selectedProjectIds.value.indexOf(id);
  if (idx > -1) {
    selectedProjectIds.value.splice(idx, 1);
  } else {
    selectedProjectIds.value.push(id);
  }
};

const handleSave = async () => {
  isSaving.value = true;
  try {
    await api.assignProjectToSite(props.id, selectedProjectIds.value);
    notification.success('Project linking updated');
    emit('navigate', 'site-detail', { id: props.id });
  } catch (err) {
    console.error('Failed to assign project', err);
    notification.error(err.message || 'Failed to update project linking');
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

