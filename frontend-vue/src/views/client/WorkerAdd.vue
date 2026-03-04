<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import { TRADES, PASS_TYPES } from '../../utils/constants.js';
import { validateNRICFIN, validateWorkPassType, validatePersonTrade } from '../../utils/validation.js';

const props = defineProps({
  id: [Number, String],
  mode: { type: String, default: 'add' } // 'add' or 'edit'
});

const emit = defineEmits(['navigate']);

const isSaving = ref(false);
const isLoading = ref(false);

const formData = ref({
  name: '',
  role: 'worker',
  status: 'active',
  email: '',
  // compliance fields
  person_id_no: '',
  person_id_and_work_pass_type: 'WP',
  person_nationality: '',
  person_trade: '1.2'
});

const formErrors = ref({});

const getTodayStr = () => {
    const d = new Date();
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    return `${year}-${month}-${day} 00:00:00`;
};
const getEndStr = () => {
    return '2037-12-31 23:59:59';
};

const authForm = ref({
  authType: 'face',
  userType: 'user',
  authStart: getTodayStr(),
  authEnd: getEndStr(),
  cardType: 'normal',
  cardNo: ''
});
const fileName = ref('');

const handleFileUpload = async (event) => {
  const file = event.target.files[0];
  if (file) {
    fileName.value = "Uploading...";
    const uploadFormData = new FormData();
    uploadFormData.append('image', file);
    // 1. Pass the trade to the backend for directoy organization
    uploadFormData.append('trade', formData.value.person_trade || 'general');

    try {
        const response = await fetch('/api/upload/face', { method: 'POST', body: uploadFormData });
        if (response.ok) {
            const data = await response.json();
            // 2. Store the full URL address as requested
            fileName.value = data.url; 
            notification.success("Image uploaded and stored in trade folder");
        } else {
            throw new Error(`Upload returned status ${response.status}`);
        }
    } catch (err) {
        fileName.value = "";
        notification.error("Failed to upload image. " + err.message);
    }
  }
};

// Sync and Deployment handlers removed as backend handles saving together




const isEdit = computed(() => props.mode === 'edit');

const fetchWorker = async () => {
  if (!isEdit.value || !props.id) return;
  isLoading.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let contextuserId = null;
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            contextuserId = user.user_id || user.id;
        } catch (e) {}
    }
    const response = await api.getWorkerById(props.id, { user_id: contextuserId });
    const data = typeof response === 'string' ? JSON.parse(response) : response;
    
    if (data) {
      formData.value = { 
        ...data,
        person_id_no: data.person_id_no || '',
        person_id_and_work_pass_type: data.person_id_and_work_pass_type || 'WP',
        person_trade: data.person_trade || '1.2'
      };

      if (data.auth_start_time) authForm.value.authStart = data.auth_start_time;
      if (data.auth_end_time) authForm.value.authEnd = data.auth_end_time;
      if (data.card_number) authForm.value.cardNo = data.card_number;
      if (data.card_type) authForm.value.cardType = data.card_type;
      
      if (data.face_img_loc) {
        fileName.value = data.face_img_loc;
        authForm.value.authType = 'face';
      } else if (data.card_number) {
        authForm.value.authType = 'card';
      }
      if (data.user_type) authForm.value.userType = data.user_type;
    }
  } finally {
    isLoading.value = false;
  }
};

onMounted(async () => {
    if (isEdit.value) {
        await fetchWorker();
    } else {
        if (!formData.value.role) formData.value.role = 'worker';
        if (!formData.value.status) formData.value.status = 'active';
    }
});

watch(() => props.id, async (newId) => {
    if (isEdit.value && newId) {
        await fetchWorker();
    }
});

const validateForm = () => {
    const errors = {};
    if (!formData.value.name) {
        errors.name = 'Full name is required';
    }

    if (!formData.value.person_id_no) {
        errors.person_id_no = 'Person Identity Number is required';
    } else if (!validateNRICFIN(formData.value.person_id_no)) {
        errors.person_id_no = 'Invalid NRIC/FIN format (e.g. S1234567D)';
    }

    if (!formData.value.person_nationality) {
        errors.person_nationality = 'Nationality is required';
    } else if (formData.value.person_nationality.length !== 2) {
        errors.person_nationality = 'Must be 2-character ISO code (e.g. SG)';
    }

    if (formData.value.person_id_and_work_pass_type && !validateWorkPassType(formData.value.person_id_and_work_pass_type)) {
        errors.person_id_and_work_pass_type = 'Invalid work pass type';
    }

    if (formData.value.person_trade && !validatePersonTrade(formData.value.person_trade)) {
        errors.person_trade = 'Invalid trade code';
    }

    formErrors.value = errors;
    return Object.keys(errors).length === 0;
};

const handleSubmit = async () => {
  if (!validateForm()) {
      notification.error('Please fix the errors in the form before submitting');
      return;
  }

  isSaving.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let userId = null;
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            userId = user.user_id || user.id;
        } catch (e) {}
    }

    const payload = { 
        ...formData.value,
        user_id: userId,
        user_type: authForm.value.userType,
        auth_start_time: authForm.value.authStart,
        auth_end_time: authForm.value.authEnd,
        card_number: authForm.value.cardNo,
        card_type: authForm.value.cardType,
        face_img_loc: fileName.value
    };

    if (isEdit.value) {
      await api.updateWorker(props.id, payload);
      notification.success('Worker profile updated');
    } else {
      await api.createWorker(payload);
      notification.success('New worker registered successfully');
    }
    emit('navigate', 'workers');
  } catch (err) {
    console.error('Failed to save worker', err);
    notification.error(err.message || 'Failed to save worker record');
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="worker-add">
    <PageHeader 
      :title="isEdit ? 'Edit Worker Detail' : 'Register New Worker'" 
      :description="isEdit ? 'Update workforce profile information' : 'Add a new worker to your organization workforce'"
    />

    <div v-if="isLoading" class="loading-state">
      <p>Fetching worker data...</p>
    </div>

    <div v-else class="split-layout">
      
      <form class="form-container" @submit.prevent="handleSubmit">
        <h3 class="panel-title">Worker Details</h3>
        <p class="panel-desc">Core workforce demographic and assignment info.</p>

      <div class="form-section">
          <h3 class="section-title">Personal Information</h3>
          <div class="form-grid">
            <BaseInput v-model="formData.name" label="Full Name" placeholder="e.g., John Smith" :error="formErrors.name" required />
            <BaseInput v-model="formData.email" label="Email Address" type="email" placeholder="john@company.com" :error="formErrors.email" />
            
            <div class="form-group">
                <label class="form-label">Designated Role</label>
                <select v-model="formData.role" class="form-select">
                    <option value="worker">Worker</option>
                    <option value="pic">PIC (Person In Charge)</option>
                    <option value="manager">Manager</option>
                </select>
            </div>

          </div>
      </div>

      <div class="form-section">
          <h3 class="section-title">Identification & Pass Info</h3>
          <div class="form-grid">
            <BaseInput 
                v-model="formData.person_id_no" 
                label="Person Identity Number (NRIC/FIN)" 
                placeholder="e.g., S1234567D" 
                required 
                maxlength="9"
                :error="formErrors.person_id_no"
            />
            
            <div class="form-group">
                <label class="form-label">ID / Work Pass Type</label>
                <select v-model="formData.person_id_and_work_pass_type" class="form-select" :class="{ 'has-error': formErrors.person_id_and_work_pass_type }">
                    <option v-for="type in PASS_TYPES" :key="type.value" :value="type.value">
                        {{ type.label }}
                    </option>
                </select>
                <span v-if="formErrors.person_id_and_work_pass_type" class="error-text">{{ formErrors.person_id_and_work_pass_type }}</span>
            </div>

            <BaseInput v-model="formData.person_nationality" label="Nationality (ISO Code)" placeholder="e.g., SG, BD, IN" maxlength="2" :error="formErrors.person_nationality" required />
            
            <div class="form-group">
                <label class="form-label">Designated Trade</label>
                <select v-model="formData.person_trade" class="form-select" :class="{ 'has-error': formErrors.person_trade }">
                    <option v-for="trade in TRADES" :key="trade.value" :value="trade.value">
                        {{ trade.label }}
                    </option>
                </select>
                <span v-if="formErrors.person_trade" class="error-text">{{ formErrors.person_trade }}</span>
            </div>
          </div>
      </div>



      <div class="form-actions">
        <BaseButton variant="secondary" @click="$emit('navigate', 'workers')" type="button">Cancel</BaseButton>
        <BaseButton :loading="isSaving" type="submit">
          {{ isEdit ? 'Save Changes' : 'Register Worker' }}
        </BaseButton>
      </div>
      </form>

      <!-- Device Auth Panel -->
      <div class="panel-right">
        <div class="auth-form-container">
          <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 24px;">
            <div>
              <h3 class="panel-title" style="margin-bottom: 4px;">IoT Authentication Setup</h3>
              <p class="panel-desc" style="margin-bottom: 0;">Push biometric/access credentials required for site entry.</p>
            </div>
            <BaseBadge v-if="isEdit" :type="formData.is_synced === 1 ? 'success' : 'warning'" style="margin-top: 4px;">
                {{ formData.is_synced === 1 ? 'Synced' : 'Pending Sync' }}
            </BaseBadge>
          </div>
          
          <div v-if="formData.worker_id" class="form-group">
            <label class="form-label">System Worker ID</label>
            <BaseInput v-model="formData.worker_id" disabled />
          </div>

          <div class="form-group">
            <label class="form-label">User Type</label>
            <select v-model="authForm.userType" class="form-select">
              <option value="user">User</option>
              <option value="visitor">Visitor</option>
              <option value="blocklist">Blocklist</option>
            </select>
          </div>

          <BaseInput 
            v-model="authForm.authStart" 
            label="Authorization Start Time (YYYY-MM-DD HH:MM:SS)" 
            placeholder="e.g., 2026-02-25 00:00:00"
          />

          <BaseInput 
            v-model="authForm.authEnd" 
            label="Authorization End Time (YYYY-MM-DD HH:MM:SS)" 
            placeholder="e.g., 2037-12-31 23:59:59"
          />

          <div class="form-group">
            <label class="form-label">Authentication Method</label>
            <select v-model="authForm.authType" class="form-select">
              <option value="face">Face Recognition</option>
              <option value="card">Access Card (NFC/RFID)</option>
            </select>
          </div>

          <div v-if="authForm.authType === 'face'" class="form-group">
             <label class="form-label">Face Image Scan</label>
             <div class="upload-area">
                <i class="ri-image-add-line upload-icon"></i>
                <span>Click or drag image here</span>
                <input type="file" accept="image/*" class="file-input" @change="handleFileUpload" />
             </div>
             <p v-if="fileName" class="file-name">Staged: {{ fileName }}</p>
          </div>

          <template v-if="authForm.authType === 'card'">
              <BaseInput 
                v-model="authForm.cardNo" 
                label="Hardware Access Card Number" 
                placeholder="e.g., 0011223344" 
                required 
              />
              <div class="form-group">
                <label class="form-label">Card Type</label>
                <select v-model="authForm.cardType" class="form-select">
                  <option value="normal">Normal Card</option>
                  <option value="super">Super Card</option>
                </select>
              </div>
          </template>
        </div>

      </div>

    </div>
  </div>
</template>

<style scoped>
.split-layout {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 32px;
  align-items: start;
}

.panel-title {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.panel-desc {
  margin: 0 0 24px;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.form-container {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 32px;
}

.auth-form-container {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.upload-area {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px;
  border: 2px dashed var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg);
  cursor: pointer;
  transition: border-color 0.2s;
}

.upload-area:hover {
  border-color: var(--color-accent);
}

.upload-icon {
  font-size: 32px;
  color: var(--color-text-muted);
  margin-bottom: 8px;
}

.file-input {
  position: absolute;
  inset: 0;
  opacity: 0;
  cursor: pointer;
}

.file-name {
    margin-top: 8px;
    font-size: 12px;
    color: var(--color-accent);
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.form-section {
    margin-bottom: 32px;
    padding-bottom: 24px;
    border-bottom: 1px solid var(--color-border);
}

.form-section:last-of-type {
    border-bottom: none;
    margin-bottom: 0;
}

.section-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--color-text-primary);
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    gap: 8px;
}

.span-full {
    grid-column: 1 / -1;
}

.helper-text {
    font-size: 11px;
    color: var(--color-text-secondary);
    margin-top: 4px;
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
}

.form-select.has-error {
    border-color: #ef4444;
}

.error-text {
    font-size: 11px;
    color: #ef4444;
    margin-top: 4px;
}

.required {
    color: var(--color-danger, #ef4444);
}

.loading-state {
  padding: 48px;
  text-align: center;
}

@media (max-width: 1024px) {
  .split-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
}
</style>
