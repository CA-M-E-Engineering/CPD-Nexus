<script setup>
</script>

<template>
  <div class="table-toolbar">
    <div class="toolbar-left">
      <slot name="left"></slot>
    </div>
    <div class="toolbar-right">
      <slot name="right"></slot>
    </div>
  </div>
</template>

<style scoped>
.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.toolbar-left {
  display: flex;
  gap: 12px;
  align-items: center;
}

.toolbar-right {
  display: flex;
  gap: 12px;
  align-items: center;
}
</style>
