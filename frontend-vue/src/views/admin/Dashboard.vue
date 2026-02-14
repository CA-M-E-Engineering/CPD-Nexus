<script setup>
import { ref, onMounted } from 'vue';
import StatCard from '../../components/ui/StatCard.vue';
import { api } from '../../services/api.js';
import { MAP_MODES } from '../../utils/constants.js';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import UnifiedMap from '../../components/ui/UnifiedMap.vue';

const stats = ref([]);
const activities = ref([]);
const loading = ref(true);

const loadDashboardData = async () => {
  loading.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let userId = null;
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            userId = user.id || user.user_id;
        } catch (e) {
            console.error('[Dashboard] Failed to parse auth_user', e);
        }
    }

    if (!userId) {
        console.warn('[Dashboard] No user ID found for analytics');
        return;
    }

    const [statsData, activityData] = await Promise.all([
      api.getDashboardStats({ user_id: userId }),
      api.getActivityLog({ user_id: userId })
    ]);

    // Transform stats object to array for UI
    stats.value = [
        { label: 'Total Companies', value: statsData.active_sites + 4, trend: '↑ 2% from last month', trendType: 'positive', icon: 'ri-building-line', color: 'blue' }, // Mock logic for total companies based on active sites + others
        { label: 'Total Workers', value: statsData.total_workers.toLocaleString(), trend: '↑ 12% from last month', trendType: 'positive', icon: 'ri-group-line', color: 'blue' },
        { label: 'Active Devices', value: statsData.total_devices.toLocaleString(), trend: '↑ 18% from last month', trendType: 'positive', icon: 'ri-cpu-line', color: 'green' },
        { label: 'System Uptime', value: '99.98%', trend: 'Stable', trendType: 'neutral', icon: 'ri-shield-check-line', color: 'purple' },
        { label: 'Compliance Rate', value: statsData.compliance_rate + '%', trend: '↓ 1% from last week', trendType: 'negative', icon: 'ri-task-line', color: 'yellow' },
    ];

    activities.value = activityData.map(a => ({
        title: `${a.action}: ${a.target}`,
        time: a.time,
        icon: 'ri-information-line', // Generic icon or map based on action
        type: 'info'
    }));

  } catch (err) {
    console.error("Failed to load dashboard data", err);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
    loadDashboardData();
});

defineEmits(['navigate']);
</script>

<template>
  <div class="admin-dashboard">
    <PageHeader 
      title="System Overview" 
      description="Monitor system-wide performance and management metrics"
    >
      <template #actions>
        <BaseButton icon="ri-download-line" variant="secondary">Generate Report</BaseButton>
      </template>
    </PageHeader>

    <div class="stats-grid">
      <StatCard 
        v-for="stat in stats" 
        :key="stat.label"
        v-bind="stat"
      />
    </div>

    <div class="quick-actions">
      <div class="action-card" @click="$emit('navigate', 'device-add')">
        <div class="action-icon"><i class="ri-cpu-line"></i></div>
        <div class="action-info">
          <h4>Provision Device</h4>
          <p>Add new hardware to the system</p>
        </div>
        <i class="ri-arrow-right-line action-arrow"></i>
      </div>
      <div class="action-card" @click="$emit('navigate', 'user-add')">
        <div class="action-icon"><i class="ri-building-line"></i></div>
        <div class="action-info">
          <h4>Register User</h4>
          <p>Add new contractor or client</p>
        </div>
        <i class="ri-arrow-right-line action-arrow"></i>
      </div>
    </div>

    <div class="dashboard-content-grid">
      <div class="content-card main-map-card">
        <div class="card-header">
          <h3 class="card-title">Global User Distribution</h3>
        </div>
        <div class="map-wrapper">
          <UnifiedMap :mode="MAP_MODES.USERS" />
        </div>
      </div>

      <div class="content-card">
        <div class="card-header">
          <h3 class="card-title">Recent System Activity</h3>
        </div>
        <div class="activity-list">
          <div v-for="(act, i) in activities" :key="i" class="activity-item">
            <div class="activity-icon" :class="act.type">
              <i :class="act.icon"></i>
            </div>
            <div class="activity-info">
              <p class="activity-title">{{ act.title }}</p>
              <p class="activity-time">{{ act.time }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.action-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  transition: all 0.2s;
}

.action-card:hover {
  border-color: var(--color-accent);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: var(--color-accent-dim);
  color: var(--color-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.action-info h4 {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 4px 0;
}

.action-info p {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin: 0;
}

.action-arrow {
  margin-left: auto;
  color: var(--color-text-tertiary);
  font-size: 20px;
  transition: transform 0.2s;
}

.action-card:hover .action-arrow {
  transform: translateX(4px);
  color: var(--color-accent);
}

.dashboard-content-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
}

@media (max-width: 1024px) {
  .dashboard-content-grid {
    grid-template-columns: 1fr;
  }
}

.content-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 24px;
}

.card-header {
  margin-bottom: 20px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.main-map-card {
  display: flex;
  flex-direction: column;
}

.map-wrapper {
  flex: 1;
  width: 100%;
  height: 100%;
  min-height: 480px;
  background: var(--color-bg);
  border-radius: var(--radius-sm);
  overflow: hidden;
  position: relative;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.activity-item {
  display: flex;
  gap: 12px;
  align-items: center;
}

.activity-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
}

.activity-icon.info { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
.activity-icon.success { background: rgba(16, 185, 129, 0.1); color: #10b981; }
.activity-icon.warning { background: rgba(245, 158, 11, 0.1); color: #f59e0b; }

.activity-info p {
  margin: 0;
}

.activity-title {
  font-size: 13px;
  color: var(--color-text-primary);
  font-weight: 500;
}

.activity-time {
  font-size: 11px;
  color: var(--color-text-muted);
  margin-top: 2px;
}
</style>
