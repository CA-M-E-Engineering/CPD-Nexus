<script setup>
import { ref, computed } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';

const isDarkMode = ref(true);
const isSyncing = ref(false);

const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value;
  document.documentElement.setAttribute('data-theme', isDarkMode.value ? 'dark' : 'light');
};

const handleSync = async () => {
  if (isSyncing.value) return;

  // Get current user ID for scoped sync
  const savedUser = localStorage.getItem('auth_user');
  const userObj = savedUser ? JSON.parse(savedUser) : null;
  const userID = userObj?.user_id || '';

  isSyncing.value = true;
  try {
    const response = await api.syncUsers(userID);
    
    let sections = [];
    if (response.sent > 0) {
      sections.push(`✅ Sync Requests Sent: ${response.sent} worker(s)\n(Awaiting bridge verification)`);
    }

    if (response.unauthenticated_workers && response.unauthenticated_workers.length > 0) {
      const list = response.unauthenticated_workers.map(w => `  • ${w.name} (${w.worker_id})`).join('\n');
      sections.push(`⚠️ Missing Biometric/Card Data (${response.unauthenticated_workers.length}):\n${list}`);
    }

    if (response.invalid_workers && response.invalid_workers.length > 0) {
      const list = response.invalid_workers.map(w => `  • ${w.name} (${w.worker_id})`).join('\n');
      sections.push(`❌ Missing Site Devices (${response.invalid_workers.length}):\n${list}`);
    }

    if (sections.length === 0) {
      notification.success("No workers found that require synchronization.");
    } else {
      const message = sections.join('\n\n');
      const isError = response.invalid_workers?.length > 0 || response.unauthenticated_workers?.length > 0;
      
      if (isError) {
        notification.error(message, 10000, true);
      } else {
        notification.success(message, 5000, true);
      }
    }
    // Only trigger page refresh if the sync request was successful (200 OK)
    emit('sync');
  } catch (err) {
    console.error('Sync error:', err);
    const errorMsg = err.data?.error || err.data?.message || err.message || 'Synchronization failed. Please check bridge connection.';
    notification.error(errorMsg, 5000, true);
  } finally {
    isSyncing.value = false;
  }
};

const props = defineProps({
  userName: { type: String, default: 'Guest' },
  role: { type: String, default: 'manager' }
});

const emit = defineEmits(['logout', 'sync']);

const avatarInitials = computed(() => {
  if (!props.userName) return '??';
  return props.userName.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
});

const displayRole = computed(() => {
  return props.role === 'manager' ? 'System Manager' : 'Client Manager';
});
</script>

<template>
  <header class="app-header">
    <div class="breadcrumb">
      <div class="breadcrumb-item">
        <i class="ri-home-4-line"></i>
        <slot name="breadcrumb"></slot>
      </div>
    </div>

    <div class="header-actions">
      <button class="sync-btn" @click="handleSync" title="Sync Data" :disabled="isSyncing">
        <i class="ri-refresh-line" :class="{ 'spinning': isSyncing }"></i>
        <span>Sync</span>
      </button>
      
      <button class="theme-toggle" @click="toggleTheme">
        <i :class="isDarkMode ? 'ri-sun-line' : 'ri-moon-line'"></i>
      </button>

      <div class="user-menu">
        <div class="user-avatar">{{ avatarInitials }}</div>
        <div class="user-info">
          <span class="user-name">{{ userName }}</span>
          <span class="user-role">{{ displayRole }}</span>
        </div>
      </div>

      <button class="logout-btn" title="Sign Out" @click="$emit('logout')">
        <i class="ri-logout-box-r-line"></i>
      </button>
    </div>
  </header>
</template>

<style scoped>
.app-header { height: var(--header-height); background: var(--color-surface); border-bottom: 1px solid var(--color-border); display: flex; align-items: center; justify-content: space-between; padding: 0 32px; position: sticky; top: 0; z-index: 100; }
.breadcrumb { display: flex; align-items: center; gap: 8px; font-size: 14px; color: var(--color-text-secondary); }
.breadcrumb-item { display: flex; align-items: center; gap: 8px; }
.header-actions { display: flex; align-items: center; gap: 16px; }
.theme-toggle { width: 40px; height: 40px; border-radius: var(--radius-sm); background: var(--color-bg); border: 1px solid var(--color-border); color: var(--color-text-secondary); display: flex; align-items: center; justify-content: center; cursor: pointer; transition: all var(--transition-fast); font-size: 18px; outline: none; }
.theme-toggle:hover { background: var(--color-surface-hover); border-color: var(--color-border-light); color: var(--color-text-primary); }
.user-menu { display: flex; align-items: center; gap: 12px; padding: 6px 16px; border-radius: var(--radius-md); background: var(--color-bg); border: 1px solid var(--color-border); cursor: pointer; transition: all var(--transition-fast); }
.user-menu:hover { background: var(--color-surface-hover); }
.user-avatar { width: 32px; height: 32px; border-radius: 50%; background: linear-gradient(135deg, var(--color-accent), #8b5cf6); display: flex; align-items: center; justify-content: center; font-weight: 600; font-size: 14px; color: white; }
.user-info { display: flex; flex-direction: column; align-items: flex-start; }
.user-name { font-size: 13px; font-weight: 500; color: var(--color-text-primary); }
.user-role { font-size: 11px; color: var(--color-text-muted); }
.logout-btn { width: 40px; height: 40px; border-radius: var(--radius-sm); background: rgba(239, 68, 68, 0.1); border: 1px solid rgba(239, 68, 68, 0.2); color: #ef4444; display: flex; align-items: center; justify-content: center; cursor: pointer; transition: all var(--transition-fast); font-size: 18px; outline: none; }
.logout-btn:hover { background: #ef4444; color: white; border-color: #ef4444; box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3); }
.sync-btn { min-width: 80px; height: 40px; padding: 0 16px; border-radius: var(--radius-sm); background: var(--color-bg); border: 1px solid var(--color-border); color: var(--color-text-secondary); display: flex; align-items: center; justify-content: center; gap: 6px; cursor: pointer; transition: all var(--transition-fast); font-size: 14px; outline: none; }
.sync-btn:hover:not(:disabled) { background: var(--color-surface-hover); border-color: var(--color-border-light); color: var(--color-accent); }
.sync-btn:disabled { opacity: 0.7; cursor: not-allowed; }
.sync-btn i { font-size: 16px; transition: transform 0.3s ease; }
.sync-btn .spinning { animation: spin-continuous 1s linear infinite; }
@keyframes spin-continuous { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
</style>
