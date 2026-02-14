<script setup>
import { computed } from 'vue';
import { notification } from '../../services/notification';

const props = defineProps({
  data: {
    type: Object,
    required: true
  }
});

const isVisible = computed(() => !!props.data);

const handleClose = () => {
  notification.clear();
};
</script>

<template>
  <Transition name="toast">
    <div v-if="isVisible" class="base-toast" :class="data.type" @click="handleClose">
      <div class="toast-icon">
        <i v-if="data.type === 'success'" class="ri-checkbox-circle-fill"></i>
        <i v-else-if="data.type === 'error'" class="ri-error-warning-fill"></i>
        <i v-else class="ri-information-fill"></i>
      </div>
      <div class="toast-message">
        {{ data.message }}
      </div>
      <button class="toast-close" @click.stop="handleClose">
        <i class="ri-close-line"></i>
      </button>
    </div>
  </Transition>
</template>

<style scoped>
.base-toast {
  position: fixed;
  bottom: 32px;
  right: 32px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  border-radius: var(--radius-md);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-lg);
  z-index: 9999;
  min-width: 320px;
  max-width: 480px;
  cursor: pointer;
}

.base-toast.success {
  border-left: 4px solid var(--color-success);
}

.base-toast.error {
  border-left: 4px solid var(--color-danger);
}

.toast-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.success .toast-icon {
  color: var(--color-success);
}

.error .toast-icon {
  color: var(--color-danger);
}

.toast-message {
  flex-grow: 1;
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.toast-close {
  background: none;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: var(--transition-fast);
}

.toast-close:hover {
  background: var(--color-surface-hover);
  color: var(--color-text-primary);
}

/* Transitions */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
