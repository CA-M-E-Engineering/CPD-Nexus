<script setup>
const props = defineProps({
  variant: { type: String, default: 'primary' }, // primary, secondary, ghost, danger
  size: { type: String, default: 'md' }, // sm, md, lg
  type: { type: String, default: 'button' },
  icon: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
});

defineEmits(['click']);
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    class="btn"
    :class="[`btn-${variant}`, `btn-${size}`, { 'is-loading': loading }]"
    @click="$emit('click', $event)"
  >
    <span v-if="loading" class="spinner"></span>
    <template v-else>
      <slot name="icon">
        <i v-if="icon" :class="icon"></i>
      </slot>
      <slot></slot>
    </template>
  </button>
</template>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border-radius: var(--radius-sm);
  font-family: inherit;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
  border: 1px solid transparent;
  outline: none;
  white-space: nowrap;
  user-select: none;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Variants */
.btn-primary {
  background: var(--color-accent);
  color: white;
}
.btn-primary:hover:not(:disabled) {
  background: var(--color-accent-hover);
  box-shadow: var(--shadow-sm);
}

.btn-secondary {
  background: var(--color-surface);
  color: var(--color-text-primary);
  border-color: var(--color-border);
}
.btn-secondary:hover:not(:disabled) {
  background: var(--color-surface-hover);
  border-color: var(--color-border-light);
}

.btn-ghost {
  background: transparent;
  color: var(--color-text-secondary);
}
.btn-ghost:hover:not(:disabled) {
  background: var(--color-surface-hover);
  color: var(--color-text-primary);
}

.btn-danger {
  background: var(--color-danger);
  color: white;
}
.btn-danger:hover:not(:disabled) {
  opacity: 0.9;
}

/* Sizes */
.btn-sm { padding: 6px 12px; font-size: 13px; }
.btn-md { padding: 10px 18px; }
.btn-lg { padding: 14px 24px; font-size: 16px; }

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: rotate 0.8s linear infinite;
}

@keyframes rotate { to { transform: rotate(360deg); } }
</style>
