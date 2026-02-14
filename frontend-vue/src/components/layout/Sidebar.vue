<script setup>
defineProps({
  activeId: {
    type: String,
    required: true
  },
  sections: {
    type: Array,
    required: true
    // Each section: { title: string, items: Array<{ id: string, label: string, icon: string }> }
  },
  currentRole: {
    type: String,
    default: 'manager'
  }
});

defineEmits(['navigate', 'update:role']);
</script>

<template>
  <aside class="app-sidebar">
    <div class="sidebar-header">
      <div class="logo">
        <div class="logo-icon">S</div>
        <span>SGBuildex</span>
      </div>
    </div>

    <nav class="sidebar-nav">
      <div v-for="section in sections" :key="section.title" class="nav-section">
        <div class="nav-section-title">{{ section.title }}</div>
        <a
          v-for="item in section.items"
          :key="item.id"
          class="nav-item"
          :class="{ active: activeId === item.id }"
          @click="$emit('navigate', item.id)"
        >
          <div class="nav-icon"><i :class="item.icon"></i></div>
          <span>{{ item.label }}</span>
        </a>
      </div>
    </nav>

  </aside>
</template>

<style scoped>
.app-sidebar {
  width: 260px;
  background: var(--color-surface);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  position: sticky;
  top: 0;
  height: 100vh;
}

.sidebar-header {
  height: 64px;
  display: flex;
  align-items: center;
  padding: 0 24px;
  border-bottom: 1px solid var(--color-border);
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-weight: 700;
  font-size: 18px;
  color: var(--color-text-primary);
}

.logo-icon {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, var(--color-accent), #8b5cf6);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.sidebar-nav {
  flex: 1;
  padding: 24px 16px;
  overflow-y: auto;
}

.nav-section {
  margin-bottom: 32px;
}

.nav-section-title {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-muted);
  margin-bottom: 8px;
  padding: 0 12px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
  cursor: pointer;
  margin-bottom: 2px;
  font-size: 14px;
}

.nav-item:hover {
  background: var(--color-surface-hover);
  color: var(--color-text-primary);
}

.nav-item.active {
  background: var(--color-accent);
  color: white;
}

.nav-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.role-switcher {
  padding: 16px;
  border-top: 1px solid var(--color-border);
}

.role-select {
  width: 100%;
  padding: 10px 12px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 14px;
  cursor: pointer;
  transition: all var(--transition-fast);
  outline: none;
}

.role-select:hover {
  border-color: var(--color-accent);
}
</style>
