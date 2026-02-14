<script setup>
import Sidebar from './Sidebar.vue';
import Header from './Header.vue';

defineProps({
  activeNavId: { type: String, required: true },
  navSections: { type: Array, required: true },
  userName: String,
  role: String,
  breadcrumbPath: String,
});

defineEmits(['navigate', 'update:role', 'logout']);
</script>

<template>
  <div class="app-layout">
    <Sidebar 
      :active-id="activeNavId" 
      :sections="navSections"
      :current-role="role"
      @navigate="$emit('navigate', $event)"
      @update:role="$emit('update:role', $event)"
    />
    
    <div class="main-container">
      <Header :user-name="userName" @logout="$emit('logout')">
        <template #breadcrumb>
          <span>{{ breadcrumbPath }}</span>
        </template>
      </Header>
      
      <main class="app-main fade-in">
        <slot></slot>
      </main>
    </div>
  </div>
</template>

<style scoped>
.app-layout { display: flex; min-height: 100vh; }
.main-container { flex: 1; display: flex; flex-direction: column; min-width: 0; }
.app-main { flex: 1; padding: 32px; background: var(--color-bg); }
</style>
