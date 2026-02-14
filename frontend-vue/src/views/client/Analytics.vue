<script setup>
import { ref, onMounted } from 'vue';
import { api } from '../../services/api.js';
import PageHeader from '../../components/ui/PageHeader.vue';
import StatCard from '../../components/ui/StatCard.vue';
import BaseFilterChip from '../../components/ui/BaseFilterChip.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import DetailCard from '../../components/ui/DetailCard.vue';
import BarChart from '../../components/ui/charts/BarChart.vue';
import DoughnutChart from '../../components/ui/charts/DoughnutChart.vue';

const activeRange = ref('Last 7 Days');
const ranges = ['Last 7 Days', 'Last 30 Days', 'Last Quarter'];

const stats = ref([]);
const loading = ref(true);

// Chart Data Refs
const attendanceChartData = ref({ labels: [], datasets: [] });
const tradeChartData = ref({ labels: [], datasets: [] });
const deviceChartData = ref({ labels: [], datasets: [] });

// Colors
const colors = ['#3b82f6', '#10b981', '#8b5cf6', '#f59e0b', '#ec4899', '#6366f1'];

const fetchAnalytics = async () => {
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
      console.warn("No User context found for analytics");
      loading.value = false;
      return; 
    }

    const [statsData, detailedData] = await Promise.all([
      api.getDashboardStats({ user_id: userId }),
      api.getDetailedAnalytics({ user_id: userId })
    ]);
    
    // Top Stats
    stats.value = [
      { label: 'Active Sites', value: statsData.active_sites.toString(), trend: 'Currently active', trendType: 'neutral', icon: 'ri-map-pin-line', color: 'green' },
      { label: 'Active Projects', value: statsData.active_projects.toString(), trend: 'Ongoing contracts', trendType: 'neutral', icon: 'ri-building-2-line', color: 'blue' },
      { label: 'Total Workers', value: statsData.total_workers.toString(), trend: 'Counted in real-time', trendType: 'positive', icon: 'ri-group-line', color: 'orange' },
      { label: 'Total Devices', value: statsData.total_devices.toString(), trend: 'Hardware registry', trendType: 'neutral', icon: 'ri-cpu-line', color: 'purple' },
    ];

    // 1. Attendance Trends (Bar Chart)
    const trendLabels = Object.keys(detailedData.attendance_trends); // e.g., Mon, Tue
    // Sort logic if keys are dates, but here they are Mon-Sun map for MVP. For real implementation, sort by date.
    // Assuming API returned Mon..Sun or similar. Mock is fixed map.
    const sortedDays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
    const trendValues = sortedDays.map(d => detailedData.attendance_trends[d] || 0);

    attendanceChartData.value = {
      labels: sortedDays,
      datasets: [{
        label: 'Check-ins',
        backgroundColor: '#3b82f6',
        borderRadius: 4,
        data: trendValues
      }]
    };

    // 2. Trade Distribution (Doughnut)
    const tradeKeys = Object.keys(detailedData.workers_by_trade);
    tradeChartData.value = {
      labels: tradeKeys,
      datasets: [{
        backgroundColor: colors,
        borderWidth: 0,
        data: tradeKeys.map(k => detailedData.workers_by_trade[k])
      }]
    };

    // 3. Device Status (Doughnut)
    const deviceKeys = Object.keys(detailedData.devices_by_status);
    const deviceColors = deviceKeys.map(k => k === 'online' ? '#10b981' : '#ef4444'); // Green for online
    deviceChartData.value = {
      labels: deviceKeys.map(k => k.toUpperCase()),
      datasets: [{
        backgroundColor: deviceColors, // Simple logic
        borderWidth: 0,
        data: deviceKeys.map(k => detailedData.devices_by_status[k])
      }]
    };

  } catch (err) {
    console.error("Failed to load analytics", err);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchAnalytics);

defineEmits(['navigate']);
</script>

<template>
  <div class="analytics-view">
    <PageHeader 
      title="Analytics" 
      description="Insights and trends across your operations"
    >
      <template #actions>
        <BaseButton variant="secondary" icon="ri-download-line">Export Report</BaseButton>
      </template>
    </PageHeader>

    <div class="table-toolbar">
      <div class="filters">
        <BaseFilterChip 
          v-for="range in ranges" 
          :key="range"
          :label="range"
          :active="activeRange === range"
          @click="activeRange = range"
        />
      </div>
    </div>

    <div class="stats-grid">
      <StatCard 
        v-for="stat in stats" 
        :key="stat.label"
        v-bind="stat"
      />
    </div>

    <!-- Charts Grid -->
    <div class="charts-grid" v-if="!loading">
      
      <!-- Row 1: Trends & Distribution -->
      <div class="chart-card wide">
        <div class="card-header">
          <h3>Weekly Attendance Trends</h3>
        </div>
        <div class="chart-wrapper">
          <BarChart :chartData="attendanceChartData" />
        </div>
      </div>

      <div class="chart-card">
        <div class="card-header">
          <h3>Workforce trade Distribution</h3>
        </div>
        <div class="chart-wrapper">
          <DoughnutChart :chartData="tradeChartData" />
        </div>
      </div>

      <!-- Row 2: Device Health & More -->
      <div class="chart-card">
        <div class="card-header">
          <h3>Device Health Status</h3>
        </div>
        <div class="chart-wrapper">
          <DoughnutChart :chartData="deviceChartData" />
        </div>
      </div>

      <div class="chart-card">
        <div class="card-header">
           <h3>Site Compliance (Mock)</h3>
        </div>
        <div class="kpi-wrapper">
           <div class="kpi-circle">
             <span class="score">95%</span>
             <span class="label">Pass</span>
           </div>
           <p class="kpi-desc">Your sites consistently meet safety regulations.</p>
        </div>
      </div>

    </div>

    <div v-else class="loading-state">
      <p>Loading analytics data...</p>
    </div>
  </div>
</template>

<style scoped>
.table-toolbar {
  margin-bottom: 24px;
}

.filters {
  display: flex;
  gap: 8px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-bottom: 32px;
}

.chart-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 24px;
  display: flex;
  flex-direction: column;
}

.chart-card.wide {
  grid-column: span 2;
}

.card-header {
  margin-bottom: 20px;
}

.card-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.chart-wrapper {
  flex: 1;
  min-height: 250px;
  position: relative;
}

.loading-state {
  text-align: center;
  padding: 64px;
  color: var(--color-text-secondary);
}

.kpi-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.kpi-circle {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 4px solid #10b981;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.kpi-circle .score {
  font-size: 32px;
  font-weight: bold;
  color: var(--color-text-primary);
}

.kpi-circle .label {
  font-size: 12px;
  color: #10b981;
  text-transform: uppercase;
  font-weight: 600;
}

.kpi-desc {
  text-align: center;
  color: var(--color-text-secondary);
  font-size: 14px;
}

@media (max-width: 1024px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
  .chart-card.wide {
    grid-column: span 1;
  }
}
</style>
