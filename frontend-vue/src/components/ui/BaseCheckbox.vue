<template>
  <label class="base-checkbox" :class="{ 'is-checked': modelValue, 'is-disabled': disabled }">
    <input
      type="checkbox"
      :checked="modelValue"
      :disabled="disabled"
      class="checkbox-input"
      @change="$emit('update:modelValue', $event.target.checked)"
    />
    <div class="checkbox-box">
      <i class="ri-check-line check-icon"></i>
    </div>
    <span v-if="label" class="checkbox-label">{{ label }}</span>
  </label>
</template>

<script setup>
defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  label: String,
  disabled: {
    type: Boolean,
    default: false
  }
});

defineEmits(['update:modelValue']);
</script>

<style scoped>
.base-checkbox {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  user-select: none;
  transition: all var(--transition-fast);
}

.checkbox-input {
  position: absolute;
  opacity: 0;
  cursor: pointer;
  height: 0;
  width: 0;
}

.checkbox-box {
  width: 20px;
  height: 20px;
  background: var(--color-bg-subtle);
  border: 1px solid var(--color-border);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  color: transparent;
}

.base-checkbox:hover .checkbox-box {
  border-color: var(--color-accent);
}

.is-checked .checkbox-box {
  background: var(--color-accent);
  border-color: var(--color-accent);
  color: white;
}

.check-icon {
  font-size: 14px;
  stroke-width: 2;
  transform: scale(0.5);
  transition: transform 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.is-checked .check-icon {
  transform: scale(1);
}

.checkbox-label {
  font-size: 14px;
  color: var(--color-text-primary);
}

.is-disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.is-disabled .checkbox-input {
  cursor: not-allowed;
}
</style>
