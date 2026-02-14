<script setup>
defineProps({
  tabs: {
    type: Array,
    required: true,
    // Each tab: { id: string, label: string }
  },
  modelValue: {
    type: String,
    required: true
  }
});

defineEmits(['update:modelValue']);
</script>

<template>
  <div class="tabs">
    <div class="tab-list">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        class="tab"
        :class="{ active: modelValue === tab.id }"
        @click="$emit('update:modelValue', tab.id)"
      >
        {{ tab.label }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.tabs {
  border-bottom: 1px solid var(--color-border);
  margin-bottom: 32px;
}

.tab-list {
  display: flex;
  gap: 4px;
}

.tab {
  padding: 12px 20px;
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all var(--transition-fast);
}

.tab:hover {
  color: var(--color-text-primary);
}

.tab.active {
  color: var(--color-accent);
  border-bottom-color: var(--color-accent);
}
</style>
