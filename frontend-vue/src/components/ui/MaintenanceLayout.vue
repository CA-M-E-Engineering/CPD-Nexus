<template>
  <div class="maintenance-layout">
    <div class="maintenance-header">
      <button class="back-link" @click="$emit('back')">
        <i class="ri-arrow-left-line"></i> Back to Detail
      </button>
      <h1 class="maintenance-title">{{ title }}</h1>
      <p v-if="description" class="maintenance-subtitle">{{ description }}</p>
    </div>

    <div class="maintenance-content">
      <div class="list-container">
        <slot name="list"></slot>
      </div>
    </div>

    <div class="maintenance-footer">
      <div class="footer-content">
        <slot name="action">
          <BaseButton 
            v-if="actionLabel" 
            :loading="loading" 
            @click="$emit('action')"
            class="primary-action-btn"
          >
            {{ actionLabel }}
          </BaseButton>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import BaseButton from './BaseButton.vue';

defineProps({
  title: String,
  description: String,
  actionLabel: String,
  loading: Boolean
});

defineEmits(['back', 'action']);
</script>

<style scoped>
.maintenance-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--color-bg);
  position: relative;
}

.maintenance-header {
  padding: 40px 48px 24px 48px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
}

.back-link {
  display: flex;
  align-items: center;
  gap: 8px;
  background: none;
  border: none;
  color: var(--color-text-tertiary);
  font-size: 13px;
  cursor: pointer;
  padding: 0;
  margin-bottom: 16px;
  transition: color var(--transition-fast);
}

.back-link:hover {
  color: var(--color-accent);
}

.maintenance-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 8px 0;
  letter-spacing: -0.5px;
}

.maintenance-subtitle {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0;
}

.maintenance-content {
  flex: 1;
  overflow-y: auto;
  padding: 32px 48px 120px 48px; /* Extra bottom padding for floating footer */
}

.list-container {
  max-width: 800px;
  margin: 0 auto;
}

.maintenance-footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 24px 48px;
  background: rgba(18, 18, 18, 0.8);
  backdrop-filter: blur(20px);
  border-top: 1px solid var(--color-border);
  display: flex;
  justify-content: center;
  z-index: 100;
}

.footer-content {
  width: 100%;
  max-width: 800px;
  display: flex;
  justify-content: flex-end;
}

.primary-action-btn {
  min-width: 200px;
  height: 48px;
  font-size: 15px;
  font-weight: 600;
}

/* Scrollbar styling */
.maintenance-content::-webkit-scrollbar {
  width: 6px;
}

.maintenance-content::-webkit-scrollbar-track {
  background: transparent;
}

.maintenance-content::-webkit-scrollbar-thumb {
  background: var(--color-border);
  border-radius: 3px;
}
</style>
