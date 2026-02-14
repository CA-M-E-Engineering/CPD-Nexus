<script setup>
import { computed } from 'vue';

const props = defineProps({
  type: {
    type: String,
    default: 'text', // 'text', 'avatar', 'rect', 'circle'
    validator: (value) => ['text', 'avatar', 'rect', 'circle'].includes(value)
  },
  width: { type: String, default: '100%' },
  height: { type: String, default: '1rem' },
  animated: { type: Boolean, default: true }
});

const skeletonClasses = computed(() => [
  'base-skeleton',
  `skeleton-${props.type}`,
  { 'skeleton-animated': props.animated }
]);

const styles = computed(() => ({
  width: props.width,
  height: props.type === 'avatar' || props.type === 'circle' ? props.width : props.height
}));
</script>

<template>
  <div :class="skeletonClasses" :style="styles"></div>
</template>

<style scoped>
.base-skeleton {
  background: var(--color-border-light);
  border-radius: var(--radius-sm);
  position: relative;
  overflow: hidden;
}

.skeleton-text {
  border-radius: 4px;
}

.skeleton-rect {
  border-radius: var(--radius-sm);
}

.skeleton-avatar, .skeleton-circle {
  border-radius: 50%;
  flex-shrink: 0;
}

.skeleton-animated::after {
  content: "";
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  transform: translateX(-100%);
  background: linear-gradient(
    90deg,
    transparent 0,
    rgba(255, 255, 255, 0.05) 20%,
    rgba(255, 255, 255, 0.1) 60%,
    transparent 100%
  );
  animation: shimmer 1.5s infinite;
}

[data-theme="light"] .skeleton-animated::after {
  background: linear-gradient(
    90deg,
    transparent 0,
    rgba(255, 255, 255, 0.4) 20%,
    rgba(255, 255, 255, 0.6) 60%,
    transparent 100%
  );
}

@keyframes shimmer {
  100% {
    transform: translateX(100%);
  }
}
</style>
