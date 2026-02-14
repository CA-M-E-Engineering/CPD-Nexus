<script setup>
defineProps({
  label: {
    type: String,
    required: true
  },
  active: {
    type: Boolean,
    default: false
  }
});

defineEmits(['click']);
</script>

<template>
  <button 
    class="filter-chip" 
    :class="{ active }"
    @click="$emit('click')"
  >
    {{ label }}
  </button>
</template>

<style scoped>
.filter-chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 16px;
  font-size: 13px;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.filter-chip:hover {
  border-color: var(--color-accent);
  background: var(--color-surface-hover);
  color: var(--color-text-primary);
}

.filter-chip.active {
  background: var(--color-accent);
  border-color: var(--color-accent);
  color: white;
}
</style>
