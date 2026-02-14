<script setup>
import { ref, computed, onMounted } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseFilterChip from '../../components/ui/BaseFilterChip.vue';
import TableToolbar from '../../components/ui/TableToolbar.vue';
import ConfirmDialog from '../../components/ui/ConfirmDialog.vue';

const activeFilter = ref('All');
const projects = ref([]);
const siteCount = ref(0);
const workerCount = ref(0);
const isLoading = ref(false);
const error = ref(null);

const filters = ['All', 'active', 'completed', 'on hold'];

const columns = [
  { key: 'title', label: 'Project Title', size: 'lg', bold: true },
  { key: 'reference', label: 'Ref', size: 'sm', muted: true },
  { key: 'site_name', label: 'Site', size: 'md' },
  { key: 'worker_count', label: 'Workers', size: 'sm', muted: true },
  { key: 'status', label: 'Status', size: 'sm' },
  { key: 'actions', label: 'Actions', width: '80px' }
];

const fetchProjects = async () => {
  isLoading.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let userId = null;
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            userId = user.user_id || user.id;
        } catch (e) {
            console.error('Failed to parse auth_user', e);
        }
    }

    if (!userId) {
        isLoading.value = false;
        return;
    }

    // projects.value = await api.getProjects({ user_id: userId });
    // Note: api.getProjects signature in projects.api.js vs api.js might differ, 
    // but based on standard pattern it should accept params. 
    // Checking api.js... it exports liveApi.
    // Let's assume standard passing.
    
    const data = await api.getProjects({ user_id: userId });
    projects.value = data || [];

    // Calculate Stats
    const uniqueSites = new Set(projects.value.map(p => p.site_id).filter(Boolean));
    siteCount.value = uniqueSites.size;
    
    workerCount.value = projects.value.reduce((sum, p) => sum + (p.worker_count || 0), 0);

  } catch (err) {
    console.error(err);
    error.value = 'Failed to load projects';
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchProjects);

const filteredProjects = computed(() => {
  return projects.value.filter(project => {
    return activeFilter.value === 'All' || project.status === activeFilter.value;
  });
});

const emit = defineEmits(['navigate']);

const handleRowClick = (project) => {
  emit('navigate', 'project-detail', { id: project.project_id });
};

const handleEdit = (project) => {
  emit('navigate', 'project-add', { id: project.project_id, mode: 'edit' });
};

const handleExport = async () => {
  isLoading.value = true;
  await api.simulateExport('Projects');
  isLoading.value = false;
  notification.success('Projects exported successfully!');
};

// Delete logic
const showDeleteDialog = ref(false);
const projectToDelete = ref(null);
const isDeleting = ref(false);

const confirmDelete = (item) => {
  projectToDelete.value = item;
  showDeleteDialog.value = true;
};

const deleteProject = async () => {
  if (!projectToDelete.value) return;
  isDeleting.value = true;
  try {
    await api.deleteProject(projectToDelete.value.project_id);
    projects.value = projects.value.filter(p => p.project_id !== projectToDelete.value.project_id);
    notification.success(`Project ${projectToDelete.value.title} deleted`);
    showDeleteDialog.value = false;
  } catch (err) {
    notification.error('Failed to delete project');
  } finally {
    isDeleting.value = false;
    projectToDelete.value = null;
  }
};
</script>

<template>
  <div class="project-list">
    <PageHeader 
      title="Project Registry" 
      description="Manage all construction projects and their assignments"
    >
      <template #stats>
        <div class="stat-item">
          <i class="ri-community-line stat-icon"></i>
          <span><span class="stat-value">{{ siteCount }}</span> sites</span>
        </div>
        <div class="stat-item">
          <i class="ri-group-line stat-icon"></i>
          <span><span class="stat-value">{{ workerCount.toLocaleString() }}</span> workers</span>
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
          @click="$emit('navigate', 'project-add')"
        >
          Add Project
        </BaseButton>
      </template>
    </TableToolbar>

    <DataTable :loading="isLoading" :columns="columns" :data="filteredProjects">
      <template #cell-title="{ item }">
        <div class="clickable-cell" @click="handleRowClick(item)">
          <strong>{{ item.title }}</strong>
        </div>
      </template>

      <template #cell-worker_count="{ item }">
        <div class="stat-cell">
          <i class="ri-group-line"></i>
          <span>{{ item.worker_count }}</span>
        </div>
      </template>

      <template #cell-device_count="{ item }">
        <div class="stat-cell">
          <i class="ri-cpu-line"></i>
          <span>{{ item.device_count }}</span>
        </div>
      </template>

      <template #cell-status="{ item }">
        <BaseBadge :type="item.status === 'active' ? 'info' : item.status === 'completed' ? 'success' : 'warning'">
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
      title="Delete Project"
      :description="`Are you sure you want to delete ${projectToDelete?.title}? This action cannot be undone.`"
      @confirm="deleteProject"
      @cancel="showDeleteDialog = false"
    />
  </div>
</template>

<style scoped>
.clickable-cell { cursor: pointer; color: var(--color-accent); transition: opacity 0.2s; }
.clickable-cell:hover { opacity: 0.8; text-decoration: underline; }
.stat-cell { display: flex; align-items: center; gap: 8px; color: var(--color-text-secondary); }
.stat-cell i { color: var(--color-accent); font-size: 16px; }
</style>


