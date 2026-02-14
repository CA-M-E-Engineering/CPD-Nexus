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
  { key: 'worker_count', label: 'Workers', size: 'md', align: 'center' },
  { key: 'device_count', label: 'Devices', size: 'md', align: 'center' },
  { key: 'email', label: 'Contact Email', size: 'md' },
  { key: 'phone', label: 'Phone', size: 'md' },
  { key: 'actions', label: 'Actions', width: '100px' }
];

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

const filteredUsers = computed(() => {
  return Users.value.filter(User => {
    return activeFilter.value === 'All' || User.status === activeFilter.value || User.user_type === activeFilter.value;
  });
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
      title="User Management" 
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
        <BaseButton icon="ri-add-line" @click="$emit('navigate', 'user-add')">
          Add User
        </BaseButton>
      </template>
    </TableToolbar>

    <DataTable 
      :loading="isLoading" 
      :columns="columns" 
      :data="filteredUsers"
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
        <BaseBadge :type="item.status === USER_STATUS.ACTIVE ? 'success' : item.status === USER_STATUS.PENDING ? 'warning' : 'danger'">
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
      title="Delete User"
      :description="`Are you sure you want to delete ${UserToDelete?.user_name}? All associated records will be archived.`"
      @confirm="deleteUser"
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
