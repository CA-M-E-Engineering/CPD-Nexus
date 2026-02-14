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
const filters = ['All', 'contractor', 'offsite_fabricator'];


const columns = [
  { key: 'company_name', label: 'Company Name', size: 'lg', bold: true },
  { key: 'uen', label: 'UEN / Reg No.', size: 'md', mono: true },
  { key: 'company_type', label: 'Entity Type', size: 'md' },
  { key: 'address', label: 'Office Address', size: 'lg', muted: true },
  { key: 'status', label: 'Status', size: 'sm' },
  { key: 'actions', label: 'Actions', width: '80px' }
];

const companies = ref([]);
const isLoading = ref(false);
const error = ref(null);

const fetchCompanies = async () => {
  isLoading.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let tenantId = '';
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            tenantId = user.tenant_id || user.id;
        } catch (e) { console.error(e); }
    }
    if (!tenantId) {
        isLoading.value = false;
        return;
    }
    const data = await api.getCompanies({ tenant_id: tenantId });
    companies.value = data || [];
  } catch (err) {
    console.error(err);
    error.value = 'Failed to load companies';
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchCompanies);

const filteredCompanies = computed(() => {
  return companies.value.filter(c => {
    return activeFilter.value === 'All' || c.company_type === activeFilter.value;
  });
});

const emit = defineEmits(['navigate']);

const handleEdit = (company) => {
  emit('navigate', 'company-add', { id: company.company_id, mode: 'edit' });
};

const handleExport = async () => {
  isLoading.value = true;
  await api.simulateExport('Companies');
  isLoading.value = false;
};

// Delete logic
const showDeleteDialog = ref(false);
const companyToDelete = ref(null);
const isDeleting = ref(false);

const confirmDelete = (item) => {
  companyToDelete.value = item;
  showDeleteDialog.value = true;
};

const deleteCompany = async () => {
  if (!companyToDelete.value) return;
  isDeleting.value = true;
  try {
    await api.deleteCompany(companyToDelete.value.company_id);
    companies.value = companies.value.filter(s => s.company_id !== companyToDelete.value.company_id);
    showDeleteDialog.value = false;
  } catch (err) {
    console.error('Failed to delete company', err);
  } finally {
    isDeleting.value = false;
    companyToDelete.value = null;
  }
};
</script>

<template>
  <div class="company-list">
    <PageHeader 
      title="Company Registry" 
      description="Manage construction entities and offsite fabrication partners"
    />

    <div v-if="error" class="error-notice">
      <i class="ri-error-warning-line"></i>
      <span>{{ error }}</span>
      <BaseButton variant="ghost" size="sm" @click="fetchCompanies">Retry</BaseButton>
    </div>

    <TableToolbar>
      <template #left>
        <BaseFilterChip 
          v-for="filter in filters" 
          :key="filter"
          :label="filter === 'offsite_fabricator' ? 'Offsite Fabricator' : filter"
          :active="activeFilter === filter"
          @click="activeFilter = filter"
        />
      </template>
      <template #right>
        <BaseButton 
          icon="ri-add-line" 
          @click="$emit('navigate', 'company-add')"
        >
          Add Company
        </BaseButton>
      </template>
    </TableToolbar>

    <DataTable :loading="isLoading" :columns="columns" :data="filteredCompanies">
      <template #cell-company_name="{ item }">
        <strong>{{ item.company_name }}</strong>
      </template>

      <template #cell-company_type="{ value }">
        <span class="type-badge">
          {{ value === 'offsite_fabricator' ? 'Offsite Fabricator' : (value === 'contractor' ? 'Contractor' : value) }}
        </span>
      </template>

      <template #cell-status="{ item }">
        <BaseBadge :type="item.status === 'active' ? 'success' : 'danger'">
          {{ item.status }}
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
      title="Delete Company Entity"
      :description="`Are you sure you want to remove ${companyToDelete?.company_name}? This will detach the entity from your records.`"
      @confirm="deleteCompany"
      @cancel="showDeleteDialog = false"
    />
  </div>
</template>

<style scoped>
.type-badge {
    text-transform: capitalize;
    font-size: 13px;
    color: var(--color-text-secondary);
}
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
