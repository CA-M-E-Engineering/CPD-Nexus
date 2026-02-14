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


const columns = [
  { key: 'site_name', label: 'Site Name', size: 'lg', bold: true },
  { key: 'location', label: 'Location', size: 'lg', muted: true },
  { key: 'worker_count', label: 'Workers', size: 'sm', muted: true },
  { key: 'device_count', label: 'Devices', size: 'sm', muted: true },
  { key: 'actions', label: 'Actions', width: '80px' }
];

const sites = ref([]);
const workerCount = ref(0);
const deviceCount = ref(0);
const isLoading = ref(false);
const error = ref(null);

const fetchSites = async () => {
  isLoading.value = true;
  try {
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
        console.warn("No tenant ID found, skipping fetch");
        isLoading.value = false;
        return;
    }

    const data = await api.getSites({ tenant_id: tenantId });
    sites.value = data || [];
    
    // Calculate totals
    workerCount.value = sites.value.reduce((sum, site) => sum + (site.worker_count || 0), 0);
    deviceCount.value = sites.value.reduce((sum, site) => sum + (site.device_count || 0), 0);
    
  } catch (err) {
    console.error(err);
    error.value = 'Failed to load sites';
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchSites);

const filteredSites = computed(() => {
  return sites.value; // No status filter anymore
});

const emit = defineEmits(['navigate']);

const handleRowClick = (site) => {
  emit('navigate', 'site-detail', { id: site.site_id });
};

const handleEdit = (site) => {
  emit('navigate', 'site-add', { id: site.site_id, mode: 'edit' });
};

const handleExport = async () => {
  isLoading.value = true;
  await api.simulateExport('Sites');
  isLoading.value = false;
};

// Delete logic
const showDeleteDialog = ref(false);
const siteToDelete = ref(null);
const isDeleting = ref(false);

const confirmDelete = (item) => {
  siteToDelete.value = item;
  showDeleteDialog.value = true;
};

const deleteSite = async () => {
  if (!siteToDelete.value) return;
  isDeleting.value = true;
  try {
    await api.deleteSite(siteToDelete.value.site_id);
    sites.value = sites.value.filter(s => s.site_id !== siteToDelete.value.site_id);
    showDeleteDialog.value = false;
  } catch (err) {
    console.error('Failed to delete site', err);
  } finally {
    isDeleting.value = false;
    siteToDelete.value = null;
  }
};
</script>

<template>
  <div class="site-list">
    <PageHeader 
      title="Site Management" 
      description="Manage your physical locations and assignments"
    >
      <template #stats>
        <div class="stat-item">
          <i class="ri-group-line stat-icon"></i>
          <span><span class="stat-value">{{ workerCount.toLocaleString() }}</span> workers</span>
        </div>
        <div class="stat-item">
          <i class="ri-cpu-line stat-icon"></i>
          <span><span class="stat-value">{{ deviceCount }}</span> devices</span>
        </div>
      </template>
      <template #actions>
        <BaseButton variant="secondary" :loading="isLoading" @click="handleExport">Export</BaseButton>
      </template>
    </PageHeader>

    <div v-if="error" class="error-notice">
      <i class="ri-error-warning-line"></i>
      <span>{{ error }}</span>
      <BaseButton variant="ghost" size="sm" @click="fetchSites">Retry</BaseButton>
    </div>

    <TableToolbar>
      <template #left>
        <!-- No status filters needed -->
      </template>
      <template #right>
        <BaseButton 
          icon="ri-add-line" 
          @click="$emit('navigate', 'site-add')"
        >
          Add Site
        </BaseButton>
      </template>
    </TableToolbar>

    <DataTable :loading="isLoading" :columns="columns" :data="filteredSites">
      <template #cell-site_name="{ item }">
        <div class="clickable-cell" @click="handleRowClick(item)">
          <strong>{{ item.site_name }}</strong>
        </div>
      </template>

      <template #cell-worker_count="{ item }">
        <div class="stat-cell">
          <i class="ri-group-line"></i>
          <span>{{ item.worker_count || 0 }}</span>
        </div>
      </template>

      <template #cell-device_count="{ item }">
        <div class="stat-cell">
          <i class="ri-cpu-line"></i>
          <span>{{ item.device_count || 0 }}</span>
        </div>
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
      title="Delete Site"
      :description="`Are you sure you want to delete ${siteToDelete?.site_name}? This action cannot be undone.`"
      @confirm="deleteSite"
      @cancel="showDeleteDialog = false"
    />
  </div>
</template>

<style scoped>
.clickable-cell { cursor: pointer; color: var(--color-accent); transition: opacity 0.2s; }
.clickable-cell:hover { opacity: 0.8; text-decoration: underline; }
.stat-cell { display: flex; align-items: center; gap: 8px; color: var(--color-text-secondary); }
.error-notice {
    margin-bottom: 16px;
    padding: 12px 16px;
    background: rgba(239, 68, 68, 0.1);
    border-left: 4px solid #ef4444;
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    gap: 12px;
    color: #ef4444;
    font-size: 14px;
}
</style>


