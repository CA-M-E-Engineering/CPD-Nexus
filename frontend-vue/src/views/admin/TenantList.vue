<script setup>
import { ref, computed, onMounted } from 'vue';
import { api } from '../../services/api.js';
import { DATA_FILTERS, TENANT_STATUS } from '../../utils/constants';
import PageHeader from '../../components/ui/PageHeader.vue';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseFilterChip from '../../components/ui/BaseFilterChip.vue';
import TableToolbar from '../../components/ui/TableToolbar.vue';
import ConfirmDialog from '../../components/ui/ConfirmDialog.vue';

const emit = defineEmits(['navigate']);

const activeFilter = ref('All');
const tenants = ref([]);
const isLoading = ref(true);
const filters = DATA_FILTERS.TENANTS;

const columns = [
  { key: 'tenant_name', label: 'Company Name', size: 'lg', bold: true },
  { key: 'worker_count', label: 'Workers', size: 'md', align: 'center' },
  { key: 'device_count', label: 'Devices', size: 'md', align: 'center' },
  { key: 'email', label: 'Contact Email', size: 'md' },
  { key: 'phone', label: 'Phone', size: 'md' },
  { key: 'actions', label: 'Actions', width: '100px' }
];

const fetchTenants = async () => {
  isLoading.value = true;
  try {
    const data = await api.getTenants();
    tenants.value = data || [];
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchTenants);

const filteredTenants = computed(() => {
  return tenants.value.filter(tenant => {
    return activeFilter.value === 'All' || tenant.status === activeFilter.value || tenant.tenant_type === activeFilter.value;
  });
});

const handleRowClick = (tenant) => {
  emit('navigate', 'tenant-detail', { id: tenant.tenant_id });
};

const handleEdit = (tenant) => {
  emit('navigate', 'tenant-add', { id: tenant.tenant_id, mode: 'edit' });
};

const handleExport = async () => {
  isLoading.value = true;
  await api.simulateExport('Tenants');
  isLoading.value = false;
};

const showDeleteDialog = ref(false);
const tenantToDelete = ref(null);
const isDeleting = ref(false);

const confirmDelete = (item) => {
  tenantToDelete.value = item;
  showDeleteDialog.value = true;
};

const deleteTenant = async () => {
  if (!tenantToDelete.value) return;
  isDeleting.value = true;
  try {
    await api.deleteTenant(tenantToDelete.value.tenant_id);
    tenants.value = tenants.value.filter(t => t.tenant_id !== tenantToDelete.value.tenant_id);
    showDeleteDialog.value = false;
  } finally {
    isDeleting.value = false;
    tenantToDelete.value = null;
  }
};
</script>

<template>
  <div class="tenant-list">
    <PageHeader 
      title="Tenant Management" 
      description="Manage all organizations and their access levels"
    >
      <template #actions>
        <BaseButton variant="secondary" @click="handleExport">Export</BaseButton>
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
        <BaseButton icon="ri-add-line" @click="$emit('navigate', 'tenant-add')">
          Add Tenant
        </BaseButton>
      </template>
    </TableToolbar>

    <DataTable 
      :loading="isLoading" 
      :columns="columns" 
      :data="filteredTenants"
      row-clickable
      @row-click="handleRowClick"
    >
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
        <BaseBadge :type="item.status === TENANT_STATUS.ACTIVE ? 'success' : item.status === TENANT_STATUS.PENDING ? 'warning' : 'danger'">
          {{ item.status.toUpperCase() }}
        </BaseBadge>
      </template>

      <template #cell-actions="{ item }">
        <div class="action-buttons-group">
          <BaseButton variant="ghost" size="sm" @click.stop="handleEdit(item)">
            <i class="ri-edit-line"></i>
          </BaseButton>
          <BaseButton variant="ghost" size="sm" class="delete-btn" @click.stop="confirmDelete(item)">
            <i class="ri-delete-bin-line"></i>
          </BaseButton>
        </div>
      </template>
    </DataTable>

    <ConfirmDialog
      :show="showDeleteDialog"
      :loading="isDeleting"
      title="Delete Tenant"
      :description="`Are you sure you want to delete ${tenantToDelete?.tenant_name}? All associated records will be archived.`"
      @confirm="deleteTenant"
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

.action-buttons-group {
  display: flex;
  gap: 4px;
}

.delete-btn:hover {
  color: var(--color-danger);
}
</style>



