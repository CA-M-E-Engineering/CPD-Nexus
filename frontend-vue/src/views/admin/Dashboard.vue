<script setup>
import { ref, onMounted } from 'vue';
import { api } from '../../services/api.js';
import { MAP_MODES } from '../../utils/constants.js';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import UnifiedMap from '../../components/ui/UnifiedMap.vue';
import StatCard from '../../components/ui/StatCard.vue';

const stats = ref({
  total_workers: 0,
  active_projects: 0,
  active_sites: 0,
  total_devices: 0,
  compliance_rate: 0
});
const activities = ref([]);
const loading = ref(true);

const loadDashboardData = async () => {
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
            console.error('[Dashboard] Failed to parse auth_user', e);
        }
    }

    if (!userId) {
        console.warn('[Dashboard] No user ID found for analytics');
        return;
    }

    const queryId = isAdmin ? 'all' : userId;

    const [statsData, activityData] = await Promise.all([
      api.getDashboardStats({ user_id: queryId }),
      api.getActivityLog({ user_id: queryId })
    ]);

    stats.value = statsData;

    // Transform activities for UI
    activities.value = activityData.map(a => {
        let icon = 'ri-information-line';
        let type = 'info';

        if (a.action.includes('Created')) { icon = 'ri-add-circle-line'; type = 'success'; }
        else if (a.action.includes('Deleted')) { icon = 'ri-delete-bin-line'; type = 'warning'; }
        else if (a.action.includes('Submission')) { icon = 'ri-send-plane-line'; type = 'info'; }
        else if (a.action.includes('Login')) { icon = 'ri-user-follow-line'; type = 'info'; }

        return {
            title: a.action,
            subtitle: a.details,
            user: a.user_name,
            time: a.time || 'Recently',
            icon: icon,
            type: type
        };
    });

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
        label="Total Workers" 
        :value="stats.total_workers.toString()" 
        icon="ri-group-line" 
        color="blue" 
      />
      <StatCard 
        label="Active Projects" 
        :value="stats.active_projects.toString()" 
        icon="ri-folder-line" 
        color="green" 
      />
      <StatCard 
        label="Operational Sites" 
        :value="stats.active_sites.toString()" 
        icon="ri-map-pin-line" 
        color="yellow" 
      />
      <StatCard 
        label="IoT Devices" 
        :value="stats.total_devices.toString()" 
        icon="ri-cpu-line" 
        color="red" 
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
          <BaseButton variant="ghost" size="sm" @click="$emit('navigate', 'activity-log')">View All</BaseButton>
        </div>
        <div class="activity-list">
          <div v-for="(act, i) in activities" :key="i" class="activity-item">
            <div class="activity-icon" :class="act.type">
              <i :class="act.icon"></i>
            </div>
            <div class="activity-info">
              <div class="activity-top">
                <p class="activity-title">{{ act.title }}</p>
                <p class="activity-time">{{ act.time }}</p>
              </div>
              <p class="activity-subtitle">{{ act.subtitle }}</p>
              <p class="activity-user" v-if="act.user"><i class="ri-user-line"></i> {{ act.user }}</p>
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
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.activity-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2px;
}

.activity-title {
  font-size: 13px;
  color: var(--color-text-primary);
  font-weight: 600;
}

.activity-time {
  font-size: 11px;
  color: var(--color-text-muted);
}

.activity-subtitle {
  font-size: 12px;
  color: var(--color-text-secondary);
  line-height: 1.4;
  margin-bottom: 4px;
}

.activity-user {
  font-size: 11px;
  color: var(--color-accent);
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
}
</style>
