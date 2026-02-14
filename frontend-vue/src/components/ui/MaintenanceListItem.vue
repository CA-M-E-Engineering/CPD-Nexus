<template>
  <div class="maintenance-item">
    <div class="item-icon">
      <i :class="icon || 'ri-checkbox-blank-circle-line'"></i>
    </div>
    <div class="item-info">
      <div class="item-primary">{{ label }}</div>
      <div v-if="secondary" class="item-secondary">{{ secondary }}</div>
    </div>
    <div v-if="status" class="item-status">
      <span class="status-dot" :class="status.toLowerCase()"></span>
      {{ status }}
    </div>
  </div>
</template>

<script setup>
defineProps({
  label: String,
  secondary: String,
  icon: String,
  status: String
});
</script>

<style scoped>
.maintenance-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  margin-bottom: 12px;
  transition: all var(--transition-fast);
}

.maintenance-item:hover {
  border-color: var(--color-border-hover);
  background: var(--color-surface-hover);
}

.item-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  border-radius: var(--radius-sm);
  color: var(--color-accent);
  font-size: 20px;
}

.item-info {
  flex: 1;
}

.item-primary {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 2px;
}

.item-secondary {
  font-size: 13px;
  color: var(--color-text-tertiary);
}

.item-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  font-weight: 500;
  color: var(--color-text-secondary);
  text-transform: capitalize;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-text-tertiary);
}

.status-dot.active, .status-dot.online, .status-dot.operational { background: var(--color-success); }
.status-dot.inactive, .status-dot.offline { background: var(--color-danger); }
.status-dot.pending, .status-dot.maintenance { background: var(--color-warning); }
</style>
