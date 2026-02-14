<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  badgeText: String,
  badgeType: String, // active, inactive, offline, pending
  rows: {
    type: Array,
    default: () => []
    // Each row: { label: string, value: string }
  }
});
</script>

<template>
  <div class="detail-card">
    <div class="detail-card-header">
      <h3 class="detail-card-title">{{ title }}</h3>
      <span v-if="badgeText" class="badge" :class="badgeType">
        <span class="badge-dot"></span>
        {{ badgeText }}
      </span>
    </div>
    <div class="detail-body">
      <div v-for="(row, index) in rows" :key="index" class="detail-row">
        <span class="detail-label">{{ row.label }}</span>
        <span class="detail-value">{{ row.value }}</span>
      </div>
      <slot></slot>
    </div>
  </div>
</template>

<style scoped>
.detail-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 24px;
  height: 100%;
}

.detail-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--color-border);
}

.detail-card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid var(--color-border);
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.detail-value {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

/* Badge styles replicated from sample.html logic */
.badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.badge.active {
  background: rgba(16, 185, 129, 0.15);
  color: var(--color-success);
}
.badge.active .badge-dot { background: var(--color-success); }

.badge.inactive {
  background: rgba(75, 85, 99, 0.15);
  color: var(--color-inactive);
}
.badge.inactive .badge-dot { background: var(--color-inactive); }

.badge.offline {
  background: rgba(239, 68, 68, 0.15);
  color: var(--color-danger);
}
.badge.offline .badge-dot { background: var(--color-danger); }

.badge.pending {
  background: rgba(245, 158, 11, 0.15);
  color: var(--color-warning);
}
.badge.pending .badge-dot { background: var(--color-warning); }
</style>
