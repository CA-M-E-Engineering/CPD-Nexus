<template>
  <Transition name="fade">
    <div v-if="show" class="confirm-overlay" @click="$emit('cancel')">
      <div class="confirm-content" @click.stop>
        <div class="confirm-header">
          <div class="icon-box" :class="variant">
            <i :class="icon"></i>
          </div>
          <h3 class="confirm-title">{{ title }}</h3>
        </div>
        
        <p class="confirm-description">
          <slot>{{ description }}</slot>
        </p>

        <div class="confirm-actions">
          <BaseButton variant="secondary" @click="$emit('cancel')">
            {{ cancelLabel }}
          </BaseButton>
          <BaseButton :variant="variant" :loading="loading" @click="$emit('confirm')">
            {{ confirmLabel }}
          </BaseButton>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import BaseButton from './BaseButton.vue';

defineProps({
  show: Boolean,
  loading: Boolean,
  title: { type: String, default: 'Confirm Action' },
  description: { type: String, default: 'Are you sure you want to proceed?' },
  confirmLabel: { type: String, default: 'Confirm' },
  cancelLabel: { type: String, default: 'Cancel' },
  variant: { type: String, default: 'danger' }, // 'danger', 'primary'
  icon: { type: String, default: 'ri-alert-line' }
});

defineEmits(['confirm', 'cancel']);
</script>

<style scoped>
.confirm-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.confirm-content {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 32px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5);
}

.confirm-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  text-align: center;
}

.icon-box {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.icon-box.danger {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.icon-box.primary {
  background: rgba(59, 130, 246, 0.1);
  color: var(--color-accent);
}

.confirm-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
}

.confirm-description {
  font-size: 15px;
  line-height: 1.6;
  color: var(--color-text-secondary);
  text-align: center;
  margin: 0 0 32px 0;
}

.confirm-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

/* Animations */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
