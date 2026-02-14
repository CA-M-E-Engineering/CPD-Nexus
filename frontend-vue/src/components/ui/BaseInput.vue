<script setup>
defineProps({
  modelValue: [String, Number],
  label: String,
  placeholder: String,
  type: { type: String, default: 'text' },
  icon: String,
  error: String,
  hint: String,
  disabled: Boolean,
});

defineEmits(['update:modelValue']);
</script>

<template>
  <div class="form-group">
    <label v-if="label" class="form-label">{{ label }}</label>
    <div class="input-group">
      <span v-if="icon" class="input-icon">
        <i :class="icon"></i>
      </span>
      <input
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        class="input"
        :class="{ 'has-icon': icon, 'has-error': error }"
        @input="$emit('update:modelValue', $event.target.value)"
      />
    </div>
    <span v-if="error" class="error-text">{{ error }}</span>
    <span v-else-if="hint" class="form-hint">{{ hint }}</span>
  </div>
</template>

<style scoped>
.form-group { margin-bottom: 20px; width: 100%; }
.form-label { display: block; font-size: 14px; font-weight: 500; margin-bottom: 8px; color: var(--color-text-primary); }
.input-group { position: relative; display: flex; align-items: center; }
.input {
  width: 100%; padding: 10px 14px; background: var(--color-surface);
  border: 1px solid var(--color-border); border-radius: var(--radius-sm);
  color: var(--color-text-primary); font-size: 14px; transition: all var(--transition-fast);
  font-family: inherit;
}
.input:focus { outline: none; border-color: var(--color-accent); box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1); }
.input:disabled { opacity: 0.6; cursor: not-allowed; background: var(--color-bg); }
.input.has-icon { padding-left: 40px; }
.input-icon { position: absolute; left: 14px; color: var(--color-text-muted); font-size: 16px; display: flex; align-items: center; }
.input.has-error { border-color: var(--color-danger); }
.error-text { display: block; font-size: 12px; color: var(--color-danger); margin-top: 6px; }
.form-hint { display: block; font-size: 12px; color: var(--color-text-muted); margin-top: 6px; }
</style>
