<template>
  <label class="base-radio" :class="{ 'is-checked': modelValue === value, 'is-disabled': disabled }">
    <input
      type="radio"
      :checked="modelValue === value"
      :disabled="disabled"
      class="radio-input"
      @change="$emit('update:modelValue', value)"
    />
    <div class="radio-box">
      <div class="radio-inner"></div>
    </div>
    <span v-if="label" class="radio-label">{{ label }}</span>
  </label>
</template>

<script setup>
defineProps({
  modelValue: [String, Number, Boolean, Object],
  value: [String, Number, Boolean, Object],
  label: String,
  disabled: {
    type: Boolean,
    default: false
  }
});

defineEmits(['update:modelValue']);
</script>

<style scoped>
.base-radio {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  user-select: none;
  transition: all var(--transition-fast);
}

.radio-input {
  position: absolute;
  opacity: 0;
  cursor: pointer;
  height: 0;
  width: 0;
}

.radio-box {
  width: 20px;
  height: 20px;
  background: var(--color-bg-subtle);
  border: 1px solid var(--color-border);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.base-radio:hover .radio-box {
  border-color: var(--color-accent);
}

.is-checked .radio-box {
  border-color: var(--color-accent);
}

.radio-inner {
  width: 10px;
  height: 10px;
  background: var(--color-accent);
  border-radius: 50%;
  transform: scale(0);
  transition: transform 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.is-checked .radio-inner {
  transform: scale(1);
}

.radio-label {
  font-size: 14px;
  color: var(--color-text-primary);
}

.is-disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.is-disabled .radio-input {
  cursor: not-allowed;
}
</style>
