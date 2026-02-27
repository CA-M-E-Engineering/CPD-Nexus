<script setup>
import { ref, onMounted, computed } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import { MAP_MODES, USER_TYPES } from '../../utils/constants.js';
import UnifiedMap from '../../components/ui/UnifiedMap.vue';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';

const props = defineProps({
  id: [Number, String],
  mode: {
    type: String,
    default: 'add'
  }
});

const emit = defineEmits(['navigate']);

const isEdit = computed(() => props.mode === 'edit');
const isLoading = ref(false);
const isSaving = ref(false);

const formData = ref({
  user_name: '',
  user_type: USER_TYPES.CLIENT,
  username: '',
  password: '',
  email: '',
  phone: '',
  address: '',
  latitude: '',
  longitude: '',
  status: 'active',
  bridge_ws_url: '',
  bridge_auth_token: '',
  bridge_status: 'inactive'
});

const fetchUser = async () => {
  if (!isEdit.value || !props.id) return;
  isLoading.value = true;
  try {
    const data = await api.getUserById(props.id);
    if (data) {
      formData.value = {
        user_name: data.user_name,
        user_type: data.user_type || USER_TYPES.CLIENT,
        username: data.username || '',
        email: data.email || '',
        phone: data.phone || '',
        address: data.address || '',
        latitude: data.lat || '',
        longitude: data.lng || '',
        status: data.status || 'active',
        bridge_ws_url: data.bridge_ws_url || '',
        bridge_auth_token: data.bridge_auth_token || '',
        bridge_status: data.bridge_status || 'inactive'
      };
    }
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchUser);

const handleSubmit = async () => {
  isSaving.value = true;
  try {
    const payload = {
        ...formData.value,
        lat: parseFloat(formData.value.latitude) || 0,
        lng: parseFloat(formData.value.longitude) || 0
    };

    if (isEdit.value) {
      await api.updateUser(props.id, payload);
      notification.success('Organization updated successfully');
    } else {
      await api.createUser(payload);
      notification.success('New organization registered successfully');
    }
    emit('navigate', 'users');
  } catch (err) {
    console.error('[UserAdd] Save failed:', err);
    notification.error(err.message || 'Failed to save organization details');
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="user-add">
    <PageHeader
      :title="isEdit ? 'Edit Organization' : 'Register New User'"
      :description="isEdit ? 'Update details for this organization account' : 'Add a new client organization to the system'"
    />

    <div v-if="isLoading" class="loading-state">
      <p>Loading record details...</p>
    </div>

    <form v-else class="form-container" @submit.prevent="handleSubmit">
      <div class="form-info-banner">
        <i class="ri-information-line"></i>
        <span>Creating a User will automatically generate a <strong>Login Account</strong> and a corresponding <strong>Primary Business Entity</strong>.</span>
      </div>

      <!-- Two-column layout -->
      <div class="two-col-layout">

        <!-- LEFT: User info -->
        <div class="form-panel">
          <div class="form-grid">
            <div class="section-title full-width">Account &amp; Identity</div>
            <BaseInput v-model="formData.user_name" label="User Name" placeholder="e.g., Acme Logistics" class="full-width" required />
            <BaseInput v-model="formData.email" label="Contact Email" type="email" placeholder="contact@user.com" required />
            <BaseInput v-model="formData.username" label="Username" placeholder="e.g., acme_admin" required />
            <BaseInput v-model="formData.password" label="Password" type="password" placeholder="At least 6 characters" class="full-width" :required="!isEdit" />

            <div class="section-title full-width">Business Details</div>
            <BaseInput v-model="formData.phone" label="Contact Phone" placeholder="e.g., +65 1234 5678" class="full-width" />
            <BaseInput v-model="formData.address" label="Office Address" placeholder="e.g., 120 Lower Delta Road" class="full-width" />
            <BaseInput v-model="formData.latitude" label="Latitude" placeholder="e.g., 1.3521" type="number" step="any" />
            <BaseInput v-model="formData.longitude" label="Longitude" placeholder="e.g., 103.8198" type="number" step="any" />
          </div>
        </div>

        <!-- RIGHT: Bridge config -->
        <div class="form-panel bridge-panel">
          <div class="section-title">Bridge Connection</div>
          <p class="bridge-description">Configure the IoT Bridge WebSocket connection for this organization. When active, the system will connect automatically on startup.</p>

          <div class="bridge-status-indicator" :class="formData.bridge_status === 'active' ? 'status-active' : 'status-inactive'">
            <i :class="formData.bridge_status === 'active' ? 'ri-wifi-line' : 'ri-wifi-off-line'"></i>
            <span>{{ formData.bridge_status === 'active' ? 'Bridge Active' : 'Bridge Inactive' }}</span>
          </div>

          <div class="bridge-fields">
            <BaseInput
              v-model="formData.bridge_ws_url"
              label="WebSocket URL"
              placeholder="ws://192.168.1.100:8081/ws"
            />
            <BaseInput
              v-model="formData.bridge_auth_token"
              label="Auth Token (Optional)"
              placeholder="Leave blank if not required"
            />
            <div class="form-group">
              <label class="form-label">Connection Status</label>
              <select v-model="formData.bridge_status" class="form-select">
                <option value="active">Active — Connect on startup</option>
                <option value="inactive">Inactive — Do not connect</option>
              </select>
              <span class="field-help">Set to Active to enable the IoT Bridge WebSocket connection for this organization's devices.</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Map full-width below both panels -->
      <div class="map-section">
        <label class="form-label">Office Location</label>
        <div class="map-wrapper">
          <UnifiedMap
            :mode="MAP_MODES.SINGLE_EDIT"
            :lat="parseFloat(formData.latitude) || 1.3521"
            :lng="parseFloat(formData.longitude) || 103.8198"
            @update:lat="v => formData.latitude = v"
            @update:lng="v => formData.longitude = v"
          />
        </div>
      </div>

      <div class="form-actions">
        <BaseButton variant="secondary" @click="$emit('navigate', 'users')">Cancel</BaseButton>
        <BaseButton :loading="isSaving" @click="handleSubmit" variant="primary">
          {{ isEdit ? 'Save Changes' : 'Register User' }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<style scoped>
/* Styles updated for premium look */
.form-container {
  max-width: 1200px;
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-lg);
  padding: 40px;
  box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.3);
}

.form-info-banner {
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: var(--radius-md);
  padding: 16px;
  margin-bottom: 32px;
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--color-text-primary);
  font-size: 14px;
}

.form-info-banner i {
  color: var(--color-accent);
  font-size: 1.25rem;
  flex-shrink: 0;
}

/* Two-column layout */
.two-col-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
  margin-bottom: 32px;
  align-items: start;
}

@media (max-width: 900px) {
  .two-col-layout {
    grid-template-columns: 1fr;
  }
}

.form-panel {
  display: flex;
  flex-direction: column;
  gap: 0;
}

/* Bridge panel styling */
.bridge-panel {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: var(--radius-lg);
  padding: 28px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.bridge-description {
  font-size: 13px;
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin: 0;
}

.bridge-status-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 600;
}

.status-active {
  background: rgba(34, 197, 94, 0.12);
  border: 1px solid rgba(34, 197, 94, 0.3);
  color: #22c55e;
}

.status-inactive {
  background: rgba(156, 163, 175, 0.1);
  border: 1px solid rgba(156, 163, 175, 0.2);
  color: var(--color-text-secondary);
}

.bridge-status-indicator i {
  font-size: 18px;
}

.bridge-fields {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Left panel inner grid */
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}

.section-title {
  font-size: 13px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--color-accent);
  letter-spacing: 0.05em;
  margin-top: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  margin-bottom: 4px;
}

.full-width {
  grid-column: 1 / -1;
}

/* Map */
.map-section {
  margin-bottom: 32px;
}

.form-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
  display: block;
  margin-bottom: 10px;
}

.map-wrapper {
  height: 350px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-md);
  overflow: hidden;
}

/* Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  padding-top: 28px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

/* Misc */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-select {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-md);
  padding: 12px;
  color: var(--color-text-primary);
  font-size: 14px;
  outline: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.form-select:focus {
  border-color: var(--color-accent);
  background: rgba(255, 255, 255, 0.08);
}

.field-help {
  font-size: 11px;
  color: var(--color-text-muted);
  line-height: 1.5;
}

.loading-state {
  padding: 80px;
  text-align: center;
  color: var(--color-text-secondary);
}
</style>
