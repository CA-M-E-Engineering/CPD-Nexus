<script setup>
defineProps({
  label: { type: String, required: true },
  value: { type: String, required: true },
  trend: String, // The trend text (e.g., "↑ 12%")
  trendType: { type: String, default: 'positive' }, // positive, negative, neutral
  icon: String,
  color: { type: String, default: 'blue' }
});
</script>

<template>
  <div class="stat-card">
    <div class="stat-header">
      <span class="stat-label">{{ label }}</span>
      <div class="stat-icon" :class="color">
        <i :class="icon"></i>
      </div>
    </div>
    <div class="stat-value">{{ value }}</div>
    <div v-if="trend" class="stat-trend" :class="trendType">
      <i v-if="trendType === 'positive'" class="ri-arrow-up-line"></i>
      <i v-else-if="trendType === 'negative'" class="ri-arrow-down-line"></i>
      {{ trend }}
    </div>
  </div>
</template>

<style scoped>
.stat-card { background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); padding: 24px; transition: all var(--transition-fast); }
.stat-card:hover { border-color: var(--color-border-light); box-shadow: var(--shadow-sm); }
.stat-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.stat-label { font-size: 13px; color: var(--color-text-secondary); font-weight: 500; }
.stat-icon { width: 36px; height: 36px; border-radius: var(--radius-sm); display: flex; align-items: center; justify-content: center; font-size: 18px; }
.stat-icon.blue { background: rgba(59, 130, 246, 0.15); color: var(--color-accent); }
.stat-icon.green { background: rgba(16, 185, 129, 0.15); color: var(--color-success); }
.stat-icon.yellow { background: rgba(245, 158, 11, 0.15); color: var(--color-warning); }
.stat-icon.red { background: rgba(239, 68, 68, 0.15); color: var(--color-danger); }
.stat-value { font-size: 32px; font-weight: 700; letter-spacing: -0.02em; margin-bottom: 4px; }
.stat-trend { font-size: 13px; display: flex; align-items: center; gap: 4px; }
.stat-trend.positive { color: var(--color-success); }
.stat-trend.negative { color: var(--color-danger); }
.stat-trend.neutral { color: var(--color-text-muted); }
</style>
