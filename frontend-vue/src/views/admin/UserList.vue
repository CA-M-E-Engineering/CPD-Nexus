<script setup>
import { ref, computed, onMounted } from 'vue';
import { api } from '../../services/api.js';
import { DATA_FILTERS, USER_STATUS } from '../../utils/constants';
import PageHeader from '../../components/ui/PageHeader.vue';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseFilterChip from '../../components/ui/BaseFilterChip.vue';
import TableToolbar from '../../components/ui/TableToolbar.vue';
import ConfirmDialog from '../../components/ui/ConfirmDialog.vue';

const emit = defineEmits(['navigate']);

const activeFilter = ref('All');
const Users = ref([]);
const isLoading = ref(true);
const filters = DATA_FILTERS.USERS;

const columns = [
  { key: 'user_name', label: 'Company Name', size: 'lg', bold: true },
  { key: 'user_type', label: 'Type', size: 'sm' },
  { key: 'worker_count', label: 'Workers', size: 'md', align: 'center' },
  { key: 'device_count', label: 'Devices', size: 'md', align: 'center' },
  { key: 'email', label: 'Contact Email', size: 'md' },
  { key: 'phone', label: 'Phone', size: 'md' },
  { key: 'actions', label: 'Actions', width: '130px' }
];

const vendorColumns = computed(() => columns.filter(c => c.key !== 'actions'));

const fetchUsers = async () => {
  isLoading.value = true;
  try {
    const data = await api.getUsers();
    Users.value = data || [];
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchUsers);

const vendorUsers = computed(() => {
  return Users.value.filter(u => u.user_type === 'vendor' && (activeFilter.value === 'All' || activeFilter.value === 'vendor' || u.status === activeFilter.value));
});

const clientUsers = computed(() => {
  return Users.value.filter(u => u.user_type !== 'vendor' && (activeFilter.value === 'All' || activeFilter.value === 'client' || u.status === activeFilter.value));
});

const handleRowClick = (User) => {
  emit('navigate', 'user-detail', { id: User.user_id });
};

const handleEdit = (User) => {
  emit('navigate', 'user-add', { id: User.user_id, mode: 'edit' });
};

const handleExport = async () => {
  isLoading.value = true;
  await api.simulateExport('users');
  isLoading.value = false;
};

const showDeleteDialog = ref(false);
const UserToDelete = ref(null);
const isDeleting = ref(false);

const confirmDelete = (item) => {
  UserToDelete.value = item;
  showDeleteDialog.value = true;
};

const deleteUser = async () => {
  if (!UserToDelete.value) return;
  isDeleting.value = true;
  try {
    await api.deleteUser(UserToDelete.value.user_id);
    Users.value = Users.value.filter(t => t.user_id !== UserToDelete.value.user_id);
    showDeleteDialog.value = false;
  } finally {
    isDeleting.value = false;
    UserToDelete.value = null;
  }
};
</script>

<template>
  <div class="User-list">
    <PageHeader 
      title="Access Control" 
      description="System-wide organizational management and authority mapping"
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
        <BaseButton icon="ri-add-line" @click="$emit('navigate', 'user-add')">
          Add Organization
        </BaseButton>
      </template>
    </TableToolbar>

    <div class="users-content-stack">
      <!-- Section: System Administrators -->
      <section v-if="vendorUsers.length > 0" class="management-section vendor-section">
        <div class="section-header">
          <div class="header-main">
            <div class="section-icon admin"><i class="ri-shield-check-line"></i></div>
            <div class="section-text">
              <h3 class="section-title">System Administrators</h3>
              <p class="section-desc">Organizations with root-level system management and auditing authority</p>
            </div>
          </div>
          <div class="section-count">{{ vendorUsers.length }}</div>
        </div>

        <DataTable 
          :loading="isLoading" 
          :columns="vendorColumns" 
          :data="vendorUsers"
          row-clickable
          class="vendor-table"
          @row-click="handleRowClick"
        >
          <template #cell-user_name="{ item }">
            <div class="user-name-cell vendor-highlight">
              <div class="vendor-avatar">
                <i class="ri-government-line"></i>
              </div>
              <div class="name-info">
                <span class="strong-text">{{ item.user_name }}</span>
                <span v-if="item.user_name === 'CA M&E Account'" class="system-badge">Primary Owner</span>
              </div>
            </div>
          </template>

          <template #cell-user_type="{ item }">
             <div class="authority-badge">
               <i class="ri-key-2-line"></i>
               <span>ROOT ACCESS</span>
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
        </DataTable>
      </section>

      <!-- Section: Client Organizations -->
      <section v-if="clientUsers.length > 0" class="management-section">
        <div class="section-header">
           <div class="header-main">
            <div class="section-icon client"><i class="ri-building-2-line"></i></div>
            <div class="section-text">
              <h3 class="section-title">Client Organizations</h3>
              <p class="section-desc">Managed tenants, contractors, and project stakeholders</p>
            </div>
          </div>
          <div class="section-count">{{ clientUsers.length }}</div>
        </div>

        <DataTable 
          :loading="isLoading" 
          :columns="columns" 
          :data="clientUsers"
          row-clickable
          @row-click="handleRowClick"
        >
          <template #cell-user_name="{ item }">
            <div class="user-name-cell">
              <div class="client-avatar">
                {{ item.user_name.charAt(0).toUpperCase() }}
              </div>
              <span>{{ item.user_name }}</span>
            </div>
          </template>

          <template #cell-user_type="{ item }">
            <BaseBadge type="info" size="sm">CLIENT</BaseBadge>
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
            <BaseBadge :type="item.status === USER_STATUS.ACTIVE ? 'success' : item.status === USER_STATUS.PENDING ? 'warning' : 'danger'">
              {{ item.status.toUpperCase() }}
            </BaseBadge>
          </template>

          <template #cell-actions="{ item }">
            <div class="action-buttons-group">
              <BaseButton variant="danger" size="sm" icon="ri-delete-bin-line" @click.stop="confirmDelete(item)">
                Delete
              </BaseButton>
            </div>
          </template>
        </DataTable>
      </section>
    </div>

    <ConfirmDialog
      :show="showDeleteDialog"
      :loading="isDeleting"
      title="Delete Organization"
      :description="`Are you sure you want to delete ${UserToDelete?.user_name}? This action will archive all associated projects and device allocations.`"
      @confirm="deleteUser"
      @cancel="showDeleteDialog = false"
    />
  </div>
</template>

<style scoped>
.users-content-stack {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.management-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 4px;
}

.header-main {
  display: flex;
  align-items: center;
  gap: 16px;
}

.section-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
}

.section-icon.admin {
  background: var(--color-accent-dim);
  color: var(--color-accent);
}

.section-icon.client {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.section-text {
  display: flex;
  flex-direction: column;
}

.section-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
}

.section-desc {
  font-size: 13px;
  color: var(--color-text-muted);
  margin: 2px 0 0 0;
}

.section-count {
  padding: 4px 12px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.user-name-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.vendor-avatar {
  width: 32px;
  height: 32px;
  background: var(--color-accent);
  color: white;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.2);
}

.client-avatar {
  width: 32px;
  height: 32px;
  background: var(--color-bg-subtle);
  border: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
}

.name-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.strong-text {
  font-weight: 700;
  color: var(--color-text-primary);
}

.system-badge {
  font-size: 10px;
  background: linear-gradient(135deg, var(--color-accent), #8b5cf6);
  color: white;
  padding: 1px 6px;
  border-radius: 4px;
  text-transform: uppercase;
  font-weight: 700;
  width: fit-content;
}

.authority-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--color-accent);
  font-size: 11px;
  font-weight:700;
  background: var(--color-accent-dim);
  padding: 4px 10px;
  border-radius: 6px;
  width: fit-content;
}

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
