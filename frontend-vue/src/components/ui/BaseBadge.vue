<script setup>
defineProps({
  type: { type: String, default: 'info' }, // success, warning, danger, info, inactive
  dot: { type: Boolean, default: true },
});
</script>

<template>
  <span class="badge" :class="`badge-${type}`">
    <span v-if="dot" class="badge-dot"></span>
    <slot></slot>
  </span>
</template>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.badge-success { background: rgba(16, 185, 129, 0.15); color: var(--color-success); }
.badge-success .badge-dot { background: var(--color-success); }

.badge-warning { background: rgba(245, 158, 11, 0.15); color: var(--color-warning); }
.badge-warning .badge-dot { background: var(--color-warning); }

.badge-danger { background: rgba(239, 68, 68, 0.15); color: var(--color-danger); }
.badge-danger .badge-dot { background: var(--color-danger); }

.badge-info { background: rgba(59, 130, 246, 0.15); color: var(--color-accent); }
.badge-info .badge-dot { background: var(--color-accent); }

.badge-inactive { background: rgba(75, 85, 99, 0.15); color: var(--color-inactive); }
.badge-inactive .badge-dot { background: var(--color-inactive); }
</style>
