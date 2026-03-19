<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import { TRADES, PASS_TYPES } from '../../utils/constants.js';
import { validateNRICFIN, validateWorkPassType, validatePersonTrade, validateNRICWithPassType } from '../../utils/validation.js';

const props = defineProps({
  id: [Number, String],
  mode: { type: String, default: 'add' }
});

const emit = defineEmits(['navigate']);

const isSaving  = ref(false);
const isLoading = ref(false);

const formData = ref({
  name: '',
  role: 'worker',
  status: 'active',
  // SGBuildex compliance
  person_id_no: '',
  person_id_and_work_pass_type: 'WP',
  person_nationality: '',
  person_trade: '1.2'
});

const formErrors = ref({});

// ── Auth panel ──
const getTodayStr = () => {
  const d = new Date();
  return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')} 00:00:00`;
};

const authForm = ref({
  authType: 'face',
  userType: 'user',
  authStart: getTodayStr(),
  authEnd: '2037-12-31 23:59:59',
  cardType: 'normal',
  cardNo: ''
});
const fileName = ref('');

const handleFileUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) return;
  fileName.value = 'Uploading...';
  const fd = new FormData();
  fd.append('image', file);
  fd.append('trade', formData.value.person_trade || 'general');
  try {
    const res = await fetch('/api/upload/face', { method: 'POST', body: fd });
    if (res.ok) {
      const data = await res.json();
      fileName.value = data.url;
      notification.success('Image uploaded successfully');
    } else {
      throw new Error(`Upload returned status ${res.status}`);
    }
  } catch (err) {
    fileName.value = '';
    notification.error('Failed to upload image. ' + err.message);
  }
};

const isEdit = computed(() => props.mode === 'edit');

const fetchWorker = async () => {
  if (!isEdit.value || !props.id) return;
  isLoading.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let contextUserId = null;
    if (savedUser) {
      try { contextUserId = JSON.parse(savedUser)?.user_id || JSON.parse(savedUser)?.id; } catch (e) { /* */ }
    }
    const response = await api.getWorkerById(props.id, { user_id: contextUserId });
    const data = typeof response === 'string' ? JSON.parse(response) : response;
    if (data) {
      formData.value = {
        ...data,
        person_id_no:                data.person_id_no || '',
        person_id_and_work_pass_type:data.person_id_and_work_pass_type || 'WP',
        person_nationality:          data.person_nationality || '',
        person_trade:                data.person_trade || '1.2',
      };
      if (data.auth_start_time) authForm.value.authStart = data.auth_start_time;
      if (data.auth_end_time)   authForm.value.authEnd   = data.auth_end_time;
      if (data.card_number)     authForm.value.cardNo    = data.card_number;
      if (data.card_type)       authForm.value.cardType  = data.card_type;
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
  if (isEdit.value) await fetchWorker();
  else {
    formData.value.role   = formData.value.role   || 'worker';
    formData.value.status = formData.value.status || 'active';
  }
});

watch(() => props.id, async (newId) => {
  if (isEdit.value && newId) await fetchWorker();
});

const validateForm = () => {
  const errors = {};

  if (!formData.value.name) {
    errors.name = 'Full name is required';
  }

  // API Mandatory: person_id_no
  if (!formData.value.person_id_no) {
    errors.person_id_no = 'Person Identity Number is required (API mandatory)';
  } else if (!validateNRICFIN(formData.value.person_id_no)) {
    errors.person_id_no = 'Invalid NRIC/FIN format (e.g. S1234567D)';
  }

  // API Mandatory: person_id_and_work_pass_type
  if (!formData.value.person_id_and_work_pass_type) {
    errors.person_id_and_work_pass_type = 'Work pass type is required (API mandatory)';
  } else if (!validateWorkPassType(formData.value.person_id_and_work_pass_type)) {
    errors.person_id_and_work_pass_type = 'Invalid work pass type';
  }

  // NRIC prefix vs pass type cross-check (ICA/MOM spec)
  if (formData.value.person_id_no && formData.value.person_id_and_work_pass_type) {
    if (!validateNRICWithPassType(formData.value.person_id_no, formData.value.person_id_and_work_pass_type)) {
      const isForeign = ['EP','SPASS','WP','ENTREPASS','LTVP'].includes(formData.value.person_id_and_work_pass_type);
      errors.person_id_no = isForeign
        ? 'FIN must start with F, G, or M for this pass type'
        : 'NRIC must start with S or T for Singapore ID card holders';
    }
  }

  // API Mandatory: person_trade
  if (!formData.value.person_trade) {
    errors.person_trade = 'Trade is required (API mandatory)';
  } else if (!validatePersonTrade(formData.value.person_trade)) {
    errors.person_trade = 'Invalid trade code';
  }

  // API Optional: person_nationality — validate format only if provided
  if (formData.value.person_nationality && formData.value.person_nationality.length !== 2) {
    errors.person_nationality = 'Must be 2-character ISO code (e.g. SG)';
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
      try { userId = JSON.parse(savedUser)?.user_id || JSON.parse(savedUser)?.id; } catch (e) { /* */ }
    }
    const payload = {
      ...formData.value,
      user_id:        userId,
      user_type:      authForm.value.userType,
      auth_start_time:authForm.value.authStart,
      auth_end_time:  authForm.value.authEnd,
      card_number:    authForm.value.cardNo,
      card_type:      authForm.value.cardType,
      face_img_loc:   fileName.value
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
      :description="isEdit ? 'Update workforce profile and SGBuildex compliance fields' : 'Add a new worker — ensure all API-mandatory fields are filled for successful CPD submission'"
    />

    <div v-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>Fetching worker data...</p>
    </div>

    <div v-else class="profile-layout">
      <!-- ── LEFT: Live Profile Sidebar ── -->
      <div class="profile-sidebar">
        <div class="profile-card">
          <div class="profile-image-wrapper" :class="{ 'has-image': !!fileName && fileName !== 'Uploading...' }">
             <img v-if="fileName && fileName !== 'Uploading...'" :src="fileName" alt="Profile" class="profile-image" />
             <div v-else class="profile-placeholder">
               <i class="ri-user-smile-line"></i>
             </div>
             
             <!-- Upload Overlay -->
             <div class="image-upload-overlay">
               <i v-if="fileName === 'Uploading...'" class="ri-loader-4-line spin-icon"></i>
               <i v-else class="ri-camera-lens-line"></i>
               <span v-if="fileName === 'Uploading...'">Uploading...</span>
               <span v-else>{{ fileName && fileName !== 'Uploading...' ? 'Change Picture' : 'Upload Picture' }}</span>
               <input v-if="fileName !== 'Uploading...'" type="file" accept="image/*" class="file-input" @change="handleFileUpload" />
             </div>
          </div>
          
          <div class="profile-header-info">
             <h2 class="profile-name">{{ formData.name || 'New Worker' }}</h2>
             <div class="profile-status">
               <span class="profile-role">{{ formData.role || 'Worker' }}</span>
             </div>
          </div>
          
          <div class="profile-divider"></div>
          
          <div class="profile-stats">
             <div class="stat-item">
               <span class="stat-label">Nationality</span>
               <span class="stat-value">{{ formData.person_nationality || '---' }}</span>
             </div>
             <div class="stat-item">
               <span class="stat-label">NRIC / FIN</span>
               <span class="stat-value">{{ formData.person_id_no || '---' }}</span>
             </div>
             <div class="stat-item">
               <span class="stat-label">Sync Status</span>
               <BaseBadge v-if="isEdit" :type="formData.is_synced === 1 ? 'success' : 'warning'">
                 {{ formData.is_synced === 1 ? 'Synced' : 'Sync Pending' }}
               </BaseBadge>
               <span v-else class="stat-value" style="color:var(--color-text-muted)">Unsaved</span>
             </div>
          </div>
        </div>
      </div>

      <!-- ── RIGHT: Form Main Area ── -->
      <form class="profile-main form-col" @submit.prevent="handleSubmit">

        <!-- Section 1: Personal Info -->
        <div class="form-section-card">
          <div class="section-header">
            <h3 class="section-title">Personal Information</h3>
            <p class="section-desc">Basic profile and role within your organisation.</p>
          </div>
          <div class="form-grid">
            <BaseInput v-model="formData.name" label="Full Name" placeholder="e.g., John Smith" :error="formErrors.name" required class="full-width" />
          </div>
        </div>

        <!-- Section 2: SGBuildex Compliance -->
        <div class="form-section-card">
          <div class="section-header">
            <h3 class="section-title">SGBuildex Compliance Fields</h3>
            <p class="section-desc">Identification, pass, and trade details required for manpower utilization submission.</p>
          </div>

          <div class="field-block mandatory-block">
            <div class="block-label">
              <i class="ri-error-warning-line"></i>
              Mandatory — required for every API submission
            </div>
            <div class="form-grid">
              <BaseInput v-model="formData.person_id_no" label="Person Identity Number (NRIC / FIN)" placeholder="e.g., S1234567D" required maxlength="9" :error="formErrors.person_id_no" />
              <div class="form-group">
                <label class="form-label">ID / Work Pass Type <span class="required">*</span></label>
                <select v-model="formData.person_id_and_work_pass_type" class="form-select" :class="{ 'has-error': formErrors.person_id_and_work_pass_type }">
                  <option v-for="type in PASS_TYPES" :key="type.value" :value="type.value">{{ type.label }}</option>
                </select>
                <span v-if="formErrors.person_id_and_work_pass_type" class="error-text">{{ formErrors.person_id_and_work_pass_type }}</span>
              </div>
              <div class="form-group full-width">
                <label class="form-label">Person Trade <span class="required">*</span></label>
                <select v-model="formData.person_trade" class="form-select" :class="{ 'has-error': formErrors.person_trade }">
                  <option v-for="trade in TRADES" :key="trade.value" :value="trade.value">{{ trade.label }}</option>
                </select>
                <span v-if="formErrors.person_trade" class="error-text">{{ formErrors.person_trade }}</span>
              </div>
            </div>
          </div>

          <div class="field-block optional-block">
            <div class="block-label">
              <i class="ri-information-line"></i>
              Optional — submitted when available
            </div>
            <div class="form-grid">
              <BaseInput v-model="formData.person_nationality" label="Person Nationality (ISO 2-char)" placeholder="e.g., SG, BD, IN, CN" maxlength="2" :error="formErrors.person_nationality">
                <template #label-suffix>
                  <span class="reg-badge hdb">HDB Mandatory</span>
                  <span class="opt-tag">BCA/LTA: not needed</span>
                </template>
              </BaseInput>
              <div class="info-note full-width">
                <i class="ri-building-2-line"></i>
                <div>
                  <strong>Person Employer Company Name &amp; UEN</strong> are configured at the <span class="inline-link" @click="$emit('navigate', 'projects')">Project level</span>. <br>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Section 3: IoT Auth panel -->
        <div class="form-section-card auth-panel-inline">
          <div class="section-header">
            <h3 class="section-title">IoT Authentication Setup</h3>
            <p class="section-desc">Push biometric/access credentials for site entry.</p>
          </div>

          <div class="form-grid">
            <div v-if="formData.worker_id" class="form-group full-width">
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

            <div class="form-group">
              <label class="form-label">Authentication Method</label>
              <select v-model="authForm.authType" class="form-select">
                <option value="face">Face Recognition</option>
                <option value="card">Access Card (NFC/RFID)</option>
              </select>
            </div>

            <BaseInput v-model="authForm.authStart" label="Auth Start (YYYY-MM-DD HH:MM:SS)" placeholder="e.g., 2026-02-25 00:00:00" />
            <BaseInput v-model="authForm.authEnd"   label="Auth End (YYYY-MM-DD HH:MM:SS)"   placeholder="e.g., 2037-12-31 23:59:59" />

            <template v-if="authForm.authType === 'card'">
              <BaseInput v-model="authForm.cardNo" label="Hardware Access Card Number" placeholder="e.g., 0011223344" required />
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

        <div class="form-actions sticky-actions">
          <BaseButton variant="secondary" @click="$emit('navigate', 'workers')" type="button">Cancel</BaseButton>
          <BaseButton :loading="isSaving" type="submit">
            {{ isEdit ? 'Save Changes' : 'Register Worker' }}
          </BaseButton>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
/* ── Layout ── */
.profile-layout {
  display: grid;
  grid-template-columns: 340px 1fr;
  gap: 28px;
  align-items: start;
}

.profile-sidebar {
  display: flex;
  flex-direction: column;
  gap: 24px;
  position: sticky;
  top: 24px;
}

.profile-main {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ── Profile Card ── */
.profile-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: transform 0.2s, box-shadow 0.2s;
}

.profile-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.profile-image-wrapper {
  width: 100%;
  aspect-ratio: 1;
  background: var(--color-bg-subtle);
  display: flex;
  justify-content: center;
  align-items: center;
  border-bottom: 1px solid var(--color-border);
  position: relative;
  overflow: hidden;
}

.profile-image-wrapper.has-image {
  background: transparent;
}

.profile-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.profile-placeholder {
  font-size: 72px;
  color: var(--color-text-muted);
}

.image-upload-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  opacity: 0;
  transition: opacity 0.2s;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
}

.profile-image-wrapper:hover .image-upload-overlay {
  opacity: 1;
}

.image-upload-overlay .upload-icon {
  font-size: 28px;
  color: #fff;
  margin-bottom: 4px;
}

.spin-icon {
  animation: spin 1s linear infinite;
}

.file-input {
  position: absolute;
  inset: 0;
  opacity: 0;
  cursor: pointer;
}

.profile-header-info {
  padding: 24px 24px 18px;
  text-align: center;
}

.profile-name {
  margin: 0 0 10px;
  font-size: 22px;
  font-weight: 700;
  color: var(--color-text-primary);
  letter-spacing: -0.01em;
}

.profile-status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.profile-role {
  font-size: 14px;
  color: var(--color-text-secondary);
  font-weight: 500;
  text-transform: capitalize;
}

.profile-divider {
  height: 1px;
  background: var(--color-border);
  margin: 0 24px;
}

.profile-stats {
  padding: 20px 24px 28px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.stat-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

/* ── Section Cards ── */
.form-section-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 24px;
}

.section-header {
  margin-bottom: 20px;
  padding-bottom: 14px;
  border-bottom: 1px solid var(--color-border);
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 4px 0;
}

.section-desc {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin: 0;
}

/* ── Field blocks ── */
.field-block {
  border-radius: var(--radius-sm);
  padding: 16px;
  margin-bottom: 16px;
}

.field-block:last-child { margin-bottom: 0; }

.mandatory-block {
  background: rgba(234, 88, 12, 0.05);
  border: 1px solid rgba(234, 88, 12, 0.25);
}

.optional-block {
  background: var(--color-bg);
  border: 1px dashed var(--color-border);
}

.block-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-bottom: 14px;
}

.mandatory-block .block-label { color: #ea580c; }
.optional-block  .block-label { color: var(--color-text-secondary); }

/* ── Info note ── */
.info-note {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  background: rgba(var(--accent-rgb), 0.06);
  border: 1px solid rgba(var(--accent-rgb), 0.2);
  border-radius: var(--radius-sm);
  padding: 12px 14px;
  font-size: 13px;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.info-note i {
  color: var(--color-accent);
  font-size: 16px;
  flex-shrink: 0;
  margin-top: 1px;
}

.inline-link {
  color: var(--color-accent);
  cursor: pointer;
  font-weight: 600;
  text-decoration: underline;
  text-underline-offset: 2px;
}

/* ── Form helpers ── */
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
}

.full-width { grid-column: 1 / -1; }

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.required { color: #ef4444; }

.form-select {
  width: 100%;
  padding: 10px 12px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 14px;
  outline: none;
  transition: border-color 0.15s;
}

.form-select:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(var(--accent-rgb), 0.1);
}

.form-select.has-error { border-color: #ef4444; }

.error-text {
  font-size: 11px;
  color: #ef4444;
}

/* Removed unneeded UI */

/* ── Loading ── */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 64px 48px;
  gap: 16px;
}

.loading-spinner {
  width: 34px;
  height: 34px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* ── Actions ── */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid var(--color-border);
}

/* ── Responsive ── */
@media (max-width: 1024px) {
  .profile-layout { grid-template-columns: 1fr; }
  .profile-sidebar { position: static; }
}

@media (max-width: 640px) {
  .form-grid { grid-template-columns: 1fr; }
}
</style>
