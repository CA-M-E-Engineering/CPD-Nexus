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
  status: 'active' // Default status
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
        status: data.status || 'active'
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

      <div class="form-grid">
        <div class="section-title full-width">Account & Identity</div>
        <BaseInput v-model="formData.user_name" label="User Name" placeholder="e.g., Acme Logistics" class="full-width" required />
        
        <BaseInput v-model="formData.email" label="Contact Email" type="email" placeholder="contact@User.com" required />
        <BaseInput v-model="formData.username" label="Username" placeholder="e.g., acme_admin" required />
        
        <BaseInput v-model="formData.password" label="Password" type="password" placeholder="At least 6 characters" class="full-width" :required="!isEdit" />
        
        <div class="section-title full-width">Business Details</div>
        <BaseInput v-model="formData.phone" label="Contact Phone" placeholder="e.g., +65 1234 5678" class="full-width" />
        
        <BaseInput v-model="formData.address" label="Office Address" placeholder="e.g., 120 Lower Delta Road" class="full-width" />
        
        <BaseInput v-model="formData.latitude" label="Latitude" placeholder="e.g., 1.3521" type="number" step="any" />
        <BaseInput v-model="formData.longitude" label="Longitude" placeholder="e.g., 103.8198" type="number" step="any" />
        
        <div class="full-width map-section">
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
  max-width: 900px;
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
}

.section-title {
  font-size: 14px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--color-accent);
  letter-spacing: 0.05em;
  margin-top: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  margin-bottom: 12px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  margin-bottom: 32px;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  padding-top: 32px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
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
}

.loading-state {
  padding: 80px;
  text-align: center;
  color: var(--color-text-secondary);
}

.full-width {
  grid-column: 1 / -1;
}

.map-section {
    margin-top: 16px;
}

.map-wrapper {
    height: 400px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: var(--radius-md);
    overflow: hidden;
    margin-top: 12px;
}

.checkbox-wrapper {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 8px;
    font-size: 14px;
    color: var(--color-text-primary);
}
</style>

