<script setup>
import { loadingStore } from '../../store/loading';
</script>

<template>
  <div 
    class="loading-bar-container" 
    :class="{ 'is-active': loadingStore.isLoading && loadingStore.isDelayed }"
  >
    <div class="loading-bar-progress"></div>
  </div>
</template>

<style scoped>
.loading-bar-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  z-index: 9999;
  opacity: 0;
  transition: opacity 0.3s ease;
  pointer-events: none;
}

.loading-bar-container.is-active {
  opacity: 1;
}

.loading-bar-progress {
  height: 100%;
  background: var(--color-accent);
  width: 100%;
  transform-origin: left;
  animation: loading-progress 2s infinite ease-in-out;
}

@keyframes loading-progress {
  0% { transform: translateX(-100%) scaleX(0.1); }
  50% { transform: translateX(0) scaleX(0.4); }
  100% { transform: translateX(100%) scaleX(0.1); }
}
</style>
