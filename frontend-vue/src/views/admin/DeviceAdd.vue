<script setup>
import { ref, onMounted, computed } from 'vue';
import { api } from '../../services/api.js';
import { useDeviceStore } from '../../features/devices/store/deviceStore';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';

const props = defineProps({
  id: [Number, String],
  mode: { type: String, default: 'add' }, // 'add' or 'edit'
  role: { type: String, default: 'manager' }
});

const isAdmin = computed(() => props.role === 'manager');

const emit = defineEmits(['navigate']);
const deviceStore = useDeviceStore();

const isSaving = ref(false);
  const isLoading = ref(false);
  const Users = ref([]);

  const formData = ref({
    sn: '',
    model: '',
    status: 'offline',
    user_id: ''
  });

  const isEdit = computed(() => props.mode === 'edit');

  const fetchUsers = async () => {
    try {
      // Fetch all Users for assignment
      Users.value = await api.getUsers();
      if (Users.value.length > 0 && !formData.value.user_id) {
          // Default to vendor if available, or first user
          const vendor = Users.value.find(u => u.user_type === 'vendor');
          formData.value.user_id = vendor ? vendor.user_id : Users.value[0].user_id;
      }
    } catch (err) {
      console.error('Failed to fetch Users', err);
    }
  };

  const fetchDevice = async () => {
    if (!isEdit.value || !props.id) return;
    isLoading.value = true;
    try {
      const data = await deviceStore.getDeviceById(props.id);
      if (data) {
        formData.value = { 
          sn: data.sn,
          model: data.model,
          status: data.status,
          user_id: data.user_id || ''
        };
        device.value = data;
      }
    } catch (err) {
      console.error('Failed to fetch device', err);
    } finally {
      isLoading.value = false;
    }
  };

  const device = ref(null);
  
  const deviceInfo = computed(() => [
    { label: 'Firmware Version', value: 'v2.4.1-stable' },
    { label: 'Registry Date', value: 'Jan 20, 2024' },
    { label: 'Last Deployed', value: 'Jan 22, 2024' }
  ]);

  const telemetry = computed(() => [
    { label: 'Last Connection', value: device.value?.last_heartbeat ? formatDate(device.value.last_heartbeat) : 'Never' },
    { label: 'Battery Level', value: (device.value?.battery || 0) + '%' },
    { label: 'Signal Quality', value: 'Good (-68 dBm)' }
  ]);

  const formatDate = (dateStr) => {
    if (!dateStr) return '---';
    const date = new Date(dateStr);
    return new Intl.DateTimeFormat('en-GB', {
      day: '2-digit', month: 'short', year: 'numeric',
      hour: '2-digit', minute: '2-digit'
    }).format(date);
  };

  onMounted(async () => {
    await fetchUsers();
    await fetchDevice();
  });

const handleSubmit = async () => {
  isSaving.value = true;
  try {
    if (isEdit.value) {
      await deviceStore.updateDevice(props.id, formData.value);
      notification.success('Device parameters updated');
    } else {
      await deviceStore.registerDevice(formData.value);
      notification.success('New device provisioned successfully');
    }
    emit('navigate', 'devices');
  } catch (err) {
    console.error('Failed to save device', err);
    notification.error(err.message || 'Failed to sync with device registry');
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="device-add">
    <PageHeader 
      :title="isEdit ? 'Device Management: ' + (device?.device_id || props.id) : 'Provision New Unit'" 
      :description="isEdit ? 'Update hardware parameters and monitor telemetry' : 'Register a new unit into the global device registry'"
      :variant="isEdit ? 'detail' : 'default'"
    >
      <template #toolbar-left v-if="isEdit">
        <BaseButton variant="ghost" size="sm" @click="$emit('navigate', 'devices')">
          <template #icon><i class="ri-arrow-left-line"></i></template>
          Back to Infrastructure
        </BaseButton>
      </template>
    </PageHeader>

    <div v-if="isLoading" class="loading-state">
      <p>Synchronizing with device registry...</p>
    </div>

    <div v-else class="dashboard-grid">
      <!-- Hardware Form -->
      <div class="form-side">
        <form class="form-container" @submit.prevent="handleSubmit">
          <div class="section-badge">Hardware Identity</div>
          <div class="form-grid">
            <BaseInput 
              v-model="formData.sn" 
              label="Serial Number (SN)" 
              placeholder="e.g., DS-K1T341AMF-001" 
              required 
              :disabled="!isAdmin"
            />
            <BaseInput 
              v-model="formData.model" 
              label="Model / Type" 
              placeholder="e.g., DS-K1T341AMF" 
              required 
              :disabled="!isAdmin"
            />
            
            <div class="form-group full-width" v-if="isAdmin">
              <label class="form-label">Deployment Assignment</label>
              <select v-model="formData.user_id" class="form-select" required>
                <option v-if="Users.length === 0" value="" disabled>No users available</option>
                <option v-for="User in Users" :key="User.user_id" :value="User.user_id">
                  {{ User.user_name }} ({{ User.user_type === 'vendor' ? 'Service Provider' : 'Fleet Owner' }})
                </option>
              </select>
            </div>

            <div v-else class="form-group full-width">
               <label class="form-label">Active Assignment</label>
               <div class="read-only-value">{{ device?.user_name || 'Unassigned' }}</div>
            </div>

            <div class="form-group full-width" v-if="isEdit && isAdmin">
               <label class="form-label">System Status</label>
               <div class="status-toggle-group">
                 <BaseButton 
                   v-for="s in ['online', 'offline', 'error']" 
                   :key="s"
                   type="button"
                   size="sm"
                   :variant="formData.status === s ? (s === 'error' ? 'danger' : 'primary') : 'secondary'"
                   @click="formData.status = s"
                 >
                   {{ s.toUpperCase() }}
                 </BaseButton>
               </div>
            </div>
          </div>

          <div class="form-actions" v-if="isAdmin">
            <BaseButton variant="secondary" @click="$emit('navigate', 'devices')">Cancel</BaseButton>
            <BaseButton :loading="isSaving" type="submit">
              {{ isEdit ? 'Update Parameters' : 'Register Device' }}
            </BaseButton>
          </div>
        </form>
      </div>

      <!-- Telemetry Side -->
      <div v-if="isEdit" class="telemetry-side">
        <div class="telemetry-card">
          <div class="card-header">
            <h3 class="card-title">Real-time Telemetry</h3>
            <BaseBadge :type="formData.status === 'online' ? 'success' : formData.status === 'error' ? 'danger' : 'inactive'">
                {{ formData.status.toUpperCase() }}
            </BaseBadge>
          </div>
          
          <div class="telemetry-list">
            <div v-for="item in telemetry" :key="item.label" class="telemetry-item">
              <span class="item-label">{{ item.label }}</span>
              <span class="item-value" :class="{ 'highlight': item.label === 'Battery Level' && parseInt(item.value) < 20 }">
                {{ item.value }}
              </span>
            </div>
          </div>

          <div class="divider"></div>

          <div class="info-list">
            <div v-for="info in deviceInfo" :key="info.label" class="info-item">
              <span class="info-label">{{ info.label }}</span>
              <span class="info-value">{{ info.value }}</span>
            </div>
          </div>

          <div class="card-footer">
             <BaseButton variant="ghost" size="sm" icon="ri-history-line">Event Timeline</BaseButton>
             <BaseButton variant="ghost" size="sm" icon="ri-settings-3-line">Advanced Diagnostics</BaseButton>
          </div>
        </div>

        <div class="status-note">
           <i class="ri-information-line"></i>
           <p>Last heartbeat recorded from encrypted endpoint. All data streams are synchronized with global registry.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-grid {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 32px;
  align-items: start;
}

.form-container {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 32px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.section-badge {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--color-primary);
  letter-spacing: 0.05em;
  margin-bottom: 24px;
  display: block;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  margin-bottom: 32px;
}

.status-toggle-group {
    display: flex;
    gap: 8px;
}

.telemetry-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.card-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.telemetry-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.telemetry-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.item-label { font-size: 13px; color: var(--color-text-muted); }
.item-value { font-size: 14px; font-weight: 600; color: var(--color-text-secondary); }
.item-value.highlight { color: var(--color-danger); }

.divider {
  height: 1px;
  background: var(--color-border);
  margin: 16px 0;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}

.info-item {
  display: flex;
  justify-content: space-between;
}

.info-label { font-size: 12px; color: var(--color-text-muted); }
.info-value { font-size: 13px; color: var(--color-text-secondary); }

.card-footer {
  display: flex;
  flex-direction: column;
  gap: 8px;
  border-top: 1px solid var(--color-border);
  padding-top: 16px;
}

.status-note {
  margin-top: 16px;
  display: flex;
  gap: 12px;
  background: rgba(59, 130, 246, 0.05);
  border: 1px solid rgba(59, 130, 246, 0.1);
  padding: 16px;
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  font-size: 12px;
  line-height: 1.5;
}

.status-note i {
  color: var(--color-primary);
  font-size: 16px;
}

@media (max-width: 1024px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.form-label {
    font-size: 14px;
    font-weight: 500;
    color: var(--color-text-secondary);
}

.form-select {
    width: 100%;
    padding: 10px 12px;
    background: var(--color-bg);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text-primary);
    font-size: 14px;
    outline: none;
    transition: border-color 0.2s;
}

.read-only-value {
    padding: 10px 12px;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text-secondary);
    font-size: 14px;
    font-weight: 500;
}

.loading-state {
  padding: 80px;
  text-align: center;
  color: var(--color-text-secondary);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
}

.full-width {
  grid-column: 1 / -1;
}
</style>
