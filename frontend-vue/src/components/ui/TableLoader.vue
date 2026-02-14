<template>
  <div class="table-loader">
    <div class="progress-base">
      <div class="progress-fill"></div>
    </div>
  </div>
</template>

<style scoped>
.table-loader {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  overflow: hidden;
  z-index: 10;
}

.progress-base {
  width: 100%;
  height: 100%;
  background: transparent;
}

.progress-fill {
  height: 100%;
  width: 100%;
  background: var(--color-accent);
  transform: translateX(-100%);
  animation: progress-slide 1.5s infinite ease-in-out;
}

@keyframes progress-slide {
  0% { transform: translateX(-100%); }
  50% { transform: translateX(0); }
  100% { transform: translateX(100%); }
}
</style>
