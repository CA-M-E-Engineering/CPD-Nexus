<script setup>
import { ref, onMounted, watch } from 'vue';
import { api } from '../../services/api.js';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseButton from '../../components/ui/BaseButton.vue';

const logs = ref([]);
const loading = ref(true);
const filters = ref({
  action: '',
  target_type: '',
});

const actionOptions = [
  { label: 'All Actions', value: '' },
  { label: 'Worker Actions', value: 'Worker' },
  { label: 'Project Actions', value: 'Project' },
  { label: 'Device Actions', value: 'Device' },
  { label: 'Submissions', value: 'Submission' },
  { label: 'Syncs', value: 'Sync' },
  { label: 'Logins', value: 'Login' },
];

const targetTypeOptions = [
  { label: 'All Types', value: '' },
  { label: 'Worker', value: 'worker' },
  { label: 'Project', value: 'project' },
  { label: 'Device', value: 'device' },
  { label: 'User', value: 'user' },
  { label: 'System', value: 'system' },
];

const fetchLogs = async () => {
  loading.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let userId = null;
    let isAdmin = false;
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            userId = user.id || user.user_id;
            // Only 'vendor' can see system-wide stats/logs
            isAdmin = user.role === 'vendor' || user.user_type === 'vendor';
        } catch (e) {
            console.error("Failed to parse user data from localStorage", e);
            // If parsing fails, userId remains null and isAdmin remains false,
            // which means the user will only see their own logs (or no logs if userId is null).
        }
    }

    const params = {
      user_id: isAdmin ? 'all' : userId,
      ...filters.value
    };

    const data = await api.getActivityLog(params);
    logs.value = data;
  } catch (err) {
    console.error("Failed to load activity logs", err);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchLogs);

// Refetch on filter change
watch(filters, fetchLogs, { deep: true });

const getIcon = (type) => {
  switch (type) {
    case 'worker': return 'ri-user-settings-line';
    case 'project': return 'ri-building-line';
    case 'device': return 'ri-cpu-line';
    case 'user': return 'ri-user-line';
    case 'system': return 'ri-settings-line';
    default: return 'ri-information-line';
  }
};

const getTypeClass = (type) => {
  switch (type) {
    case 'worker': return 'type-worker';
    case 'project': return 'type-project';
    case 'device': return 'type-device';
    case 'user': return 'type-user';
    case 'system': return 'type-system';
    default: return '';
  }
};
</script>

<template>
  <div class="activity-logs">
    <PageHeader 
      title="Activity Logs" 
      description="Historical record of all system interactions and events"
    >
      <template #actions>
        <BaseButton icon="ri-refresh-line" variant="secondary" @click="fetchLogs">Refresh</BaseButton>
      </template>
    </PageHeader>

    <div class="filter-bar">
      <div class="filter-group">
        <label>Action Filter</label>
        <select v-model="filters.action">
          <option v-for="opt in actionOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>
      <div class="filter-group">
        <label>Target Type</label>
        <select v-model="filters.target_type">
          <option v-for="opt in targetTypeOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>
    </div>

    <div class="logs-container">
      <div v-if="loading" class="loading-state">
        <i class="ri-loader-4-line spin"></i>
        <span>Loading activities...</span>
      </div>

      <div v-else-if="logs.length === 0" class="empty-state">
        <i class="ri-inbox-line"></i>
        <p>No activities found matching your criteria</p>
      </div>

      <div v-else class="logs-table-wrapper">
        <table class="logs-table">
          <thead>
            <tr>
              <th>Timestamp</th>
              <th>User</th>
              <th>Action</th>
              <th>Target</th>
              <th>Details</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in logs" :key="log.id">
              <td class="time-col">{{ log.time || log.created_at }}</td>
              <td class="user-col">
                <div class="user-badge">
                  <i class="ri-user-3-line"></i>
                  <span>{{ log.user_name || 'System' }}</span>
                </div>
              </td>
              <td>
                <span class="action-tag">{{ log.action }}</span>
              </td>
              <td class="target-col">
                <span :class="['target-type-tag', getTypeClass(log.target_type)]">
                  <i :class="getIcon(log.target_type)"></i>
                  {{ log.target_type }}
                </span>
                <span class="target-id">{{ log.target_id }}</span>
              </td>
              <td class="details-col">{{ log.details }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.filter-bar {
  display: flex;
  gap: 24px;
  background: var(--color-surface);
  padding: 16px 24px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  margin-bottom: 24px;
  align-items: flex-end;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-group label {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-muted);
  text-transform: uppercase;
}

.filter-group select {
  background: var(--color-background);
  border: 1px solid var(--color-border);
  color: var(--color-text);
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  min-width: 180px;
}

.logs-container {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  min-height: 400px;
}

.logs-table-wrapper {
  overflow-x: auto;
}

.logs-table {
  width: 100%;
  border-collapse: collapse;
}

.logs-table th {
  text-align: left;
  padding: 16px 24px;
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-muted);
  border-bottom: 1px solid var(--color-border);
  background: rgba(255, 255, 255, 0.02);
}

.logs-table td {
  padding: 16px 24px;
  border-bottom: 1px solid var(--color-border);
  vertical-align: middle;
}

.time-col {
  color: var(--color-text-muted);
  font-size: 13px;
  white-space: nowrap;
}

.user-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.user-badge i {
  color: var(--color-accent);
}

.action-tag {
  background: rgba(255, 255, 255, 0.05);
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 500;
}

.target-type-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  text-transform: capitalize;
  margin-right: 8px;
}

.target-id {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  color: var(--color-text-muted);
}

.type-worker { background: rgba(59, 130, 246, 0.1); color: #60a5fa; }
.type-project { background: rgba(16, 185, 129, 0.1); color: #34d399; }
.type-device { background: rgba(245, 158, 11, 0.1); color: #fbbf24; }
.type-user { background: rgba(139, 92, 246, 0.1); color: #a78bfa; }
.type-system { background: rgba(107, 114, 128, 0.1); color: #9ca3af; }

.details-col {
  font-size: 13px;
  color: var(--color-text-muted);
  max-width: 400px;
}

.loading-state, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px;
  gap: 16px;
  color: var(--color-text-muted);
}

.loading-state i, .empty-state i {
  font-size: 48px;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
