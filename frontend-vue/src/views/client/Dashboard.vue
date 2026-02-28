<script setup>
import { ref, onMounted } from 'vue';
import { api } from '../../services/api.js';
import PageHeader from '../../components/ui/PageHeader.vue';
import StatCard from '../../components/ui/StatCard.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import UnifiedMap from '../../components/ui/UnifiedMap.vue';


const stats = ref([]);
const recentProjects = ref([]);
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
               console.error("Failed to parse auth_user", e);
           }
        }

        if (!userId) {
             console.warn("No User context found for dashboard");
             loading.value = false;
             return;
        }

        const [statsData, projectsData] = await Promise.all([
            api.getDashboardStats({ user_id: userId }),
            api.getProjects({ user_id: userId })
        ]);

        stats.value = [
            { label: 'Active Sites', value: statsData.active_sites.toString(), trend: '', trendType: 'neutral', icon: 'ri-map-pin-line', color: 'blue' },
            { label: 'Total Workers', value: statsData.total_workers.toLocaleString(), trend: '', trendType: 'positive', icon: 'ri-group-line', color: 'green' },
            { label: 'Total Devices', value: statsData.total_devices.toString(), trend: '', trendType: 'neutral', icon: 'ri-cpu-line', color: 'purple' }
        ];

        recentProjects.value = (projectsData || []).slice(0, 5).map(p => ({
          title: p.title,
          reference: p.reference || 'No ref',
          status: p.status,
          icon: 'ri-building-2-line',
          type: 'info'
        }));

    } catch (err) {
        console.error("Failed to load client dashboard", err);
    } finally {
        loading.value = false;
    }
};

onMounted(() => {
    loadDashboardData();
});

const quickActions = [
  { label: 'Add Site', icon: 'ri-map-pin-add-line', target: 'site-add' },
  { label: 'Register Worker', icon: 'ri-user-add-line', target: 'worker-add' },
  { label: 'Create Project', icon: 'ri-building-2-line', target: 'project-add' }
];

defineEmits(['navigate']);
</script>

<template>
  <div class="client-dashboard">
    <PageHeader 
      title="Global Build Group Dashboard" 
      description="Real-time operational overview of your construction sites and workforce"
    >
      <template #actions>
        <BaseButton icon="ri-add-line" @click="$emit('navigate', 'sites')">Add New Site</BaseButton>
      </template>
    </PageHeader>

    <div class="stats-grid">
      <StatCard 
        v-for="stat in stats" 
        :key="stat.label"
        v-bind="stat"
      />
    </div>

    <div class="dashboard-content-grid">
      <div class="content-card main-map-card">
        <div class="card-header">
          <h3 class="card-title">Site Operational Map</h3>
        </div>
        <div class="map-wrapper">
           <UnifiedMap mode="sites" />
        </div>
      </div>

      <div class="side-content-stack">
        <div class="content-card quick-actions-card">
          <div class="card-header">
            <h3 class="card-title">Quick Actions</h3>
          </div>
          <div class="actions-grid">
            <button 
              v-for="action in quickActions" 
              :key="action.label" 
              class="action-btn"
              @click="$emit('navigate', action.target)"
            >
              <i :class="action.icon"></i>
              <span>{{ action.label }}</span>
            </button>
          </div>
        </div>

        <div class="content-card">
          <div class="card-header">
            <h3 class="card-title">Recent Projects</h3>
          </div>
          <div class="activity-list">
            <div v-if="recentProjects.length === 0" class="empty-state">No projects found.</div>
            <div v-for="(project, i) in recentProjects" :key="i" class="activity-item">
              <div class="activity-icon" :class="project.type">
                <i :class="project.icon"></i>
              </div>
              <div class="activity-info">
                <p class="activity-title">{{ project.title }}</p>
                <p class="activity-time">Ref: {{ project.reference }} &bull; {{ project.status }}</p>
              </div>
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
  margin-bottom: 24px;
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

.side-content-stack {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.actions-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.action-btn:hover {
  background: var(--color-surface-hover);
  border-color: var(--color-accent);
  color: var(--color-accent);
  transform: translateY(-2px);
}

.action-btn i {
  font-size: 20px;
}

.action-btn span {
  font-size: 11px;
  font-weight: 500;
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

.empty-state {
  font-size: 13px;
  color: var(--color-text-secondary);
  text-align: center;
  padding: 12px;
}
</style>
