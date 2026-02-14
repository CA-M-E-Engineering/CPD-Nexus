<script setup>
import BaseButton from './BaseButton.vue';

defineProps({
  title: { type: String, required: true },
  description: String,
  showSearch: { type: Boolean, default: false },

  searchPlaceholder: { type: String, default: 'Search...' },
  primaryActionLabel: String,
  primaryActionIcon: String,
  modelValue: String,
  variant: { type: String, default: 'default' } // 'default' or 'detail'
});

defineEmits(['update:modelValue', 'primaryAction']);
</script>

<template>
  <div class="page-header" :class="{ 'is-detail': variant === 'detail' }">
    <div class="header-main">
      <h1 class="page-title">{{ title }}</h1>
      <p v-if="description" class="page-description">{{ description }}</p>
      <div v-if="$slots.stats" class="page-stats">
        <slot name="stats"></slot>
      </div>
    </div>
    
    <div class="page-toolbar">
      <div class="toolbar-left">
        <slot name="toolbar-left">
          <div v-if="showSearch && variant !== 'detail'" class="input-group search-group">
            <span class="input-icon">
              <i class="ri-search-2-line"></i>
            </span>
            <input
              type="text"
              :value="modelValue"
              :placeholder="searchPlaceholder"
              class="input search-input"
              @input="$emit('update:modelValue', $event.target.value)"
            />
          </div>
        </slot>
      </div>
      <div class="toolbar-right">
        <slot name="toolbar-right">
          <BaseButton v-if="primaryActionLabel" @click="$emit('primaryAction')">
            <template #icon v-if="primaryActionIcon"><i :class="primaryActionIcon"></i></template>
            {{ primaryActionLabel }}
          </BaseButton>
        </slot>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-header { margin-bottom: 32px; }
.header-main { margin-bottom: 24px; }
.page-title { font-size: 28px; font-weight: 700; margin-bottom: 8px; letter-spacing: -0.02em; color: var(--color-text-primary); }
.page-description { font-size: 15px; color: var(--color-text-secondary); }
.page-stats { display: flex; align-items: center; gap: 24px; margin-top: 16px; }
:deep(.stat-item) { display: flex; align-items: center; gap: 8px; font-size: 14px; color: var(--color-text-secondary); }
:deep(.stat-value) { font-weight: 600; color: var(--color-text-primary); }
:deep(.stat-icon) { color: var(--color-accent); font-size: 18px; }
.page-toolbar { display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-top: 24px; }

/* Detail Variant specific styling */
.is-detail .page-toolbar {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.toolbar-left, .toolbar-right { display: flex; align-items: center; gap: 12px; }
.search-group { position: relative; display: flex; align-items: center; }
.search-input { width: 320px; padding: 10px 14px 10px 40px; background: var(--color-surface); border: 1px solid var(--color-border); border-radius: var(--radius-sm); color: var(--color-text-primary); font-size: 14px; transition: all var(--transition-fast); outline: none; }
.search-input:focus { border-color: var(--color-accent); box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1); }
.input-icon { position: absolute; left: 14px; color: var(--color-text-muted); font-size: 16px; display: flex; align-items: center; }
</style>
