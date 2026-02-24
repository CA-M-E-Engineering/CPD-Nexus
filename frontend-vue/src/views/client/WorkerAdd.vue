<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';

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

const passTypes = [
    { value: 'SP', label: 'Singapore Pink IC (SP)' },
    { value: 'SB', label: 'Singapore Blue IC (SB)' },
    { value: 'EP', label: 'Employment Pass (EP)' },
    { value: 'SPASS', label: 'S Pass (SPASS)' },
    { value: 'WP', label: 'Work Permit (WP)' },
    { value: 'ENTREPASS', label: 'EntrePass' },
    { value: 'LTVP', label: 'Long-Term Visit Pass (LTVP)' }
];

const bcaTrades = [
    { value: '1.1', label: '1.1 - Site Management (Ancillary)' },
    { value: '1.2', label: '1.2 - Site Support (Ancillary)' },
    { value: '1.3', label: '1.3 - General Machine Operation' },
    { value: '1.4', label: '1.4 - Site Preparation' },
    { value: '1.5', label: '1.5 - Scaffolding' },
    { value: '2.1', label: '2.1 - Demolition (Civil/Structural)' },
    { value: '2.2', label: '2.2 - Earthworks' },
    { value: '2.3', label: '2.3 - Foundation' },
    { value: '2.4', label: '2.4 - Tunnelling' },
    { value: '2.5', label: '2.5 - Reinforced Concrete' },
    { value: '2.6', label: '2.6 - Structural Steel' },
    { value: '3.1', label: '3.1 - Ceiling (Architectural)' },
    { value: '3.2', label: '3.2 - Partition Wall' },
    { value: '4.1', label: '4.1 - Plumbing, Sanitary & Gas' },
    { value: '4.3', label: '4.3 - Electrical' }
];

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

const validateId = (val) => {
    const regex = /^[STFGM]\d{7}[A-Z0-9]$/;
    return regex.test(val);
};

const handleSubmit = async () => {
  if (!validateId(formData.value.person_id_no)) {
      notification.error('Invalid Person Identity Number format (e.g. S1234567D)');
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
        user_id: userId
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

    <form v-else class="form-container" @submit.prevent="handleSubmit">
      
      <div class="form-section">
          <h3 class="section-title">Personal Information</h3>
          <div class="form-grid">
            <BaseInput v-model="formData.name" label="Full Name" placeholder="e.g., John Smith" required />
            <BaseInput v-model="formData.email" label="Email Address" type="email" placeholder="john@company.com" />
            
            <div class="form-group">
                <label class="form-label">Designated Role</label>
                <select v-model="formData.role" class="form-select">
                    <option value="worker">Worker</option>
                    <option value="pic">PIC (Person In Charge)</option>
                    <option value="manager">Manager</option>
                </select>
            </div>

            <div class="form-group">
                <label class="form-label">System Status</label>
                <select v-model="formData.status" class="form-select">
                    <option value="active">Active</option>
                    <option value="inactive">Inactive</option>
                    <option value="suspended">Suspended</option>
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
            />
            
            <div class="form-group">
                <label class="form-label">ID / Work Pass Type</label>
                <select v-model="formData.person_id_and_work_pass_type" class="form-select">
                    <option v-for="type in passTypes" :key="type.value" :value="type.value">
                        {{ type.label }}
                    </option>
                </select>
            </div>

            <BaseInput v-model="formData.person_nationality" label="Nationality (ISO Code)" placeholder="e.g., SG, BD, IN" maxlength="2" required />
            
            <div class="form-group">
                <label class="form-label">Designated Trade</label>
                <select v-model="formData.person_trade" class="form-select">
                    <option v-for="trade in bcaTrades" :key="trade.value" :value="trade.value">
                        {{ trade.label }}
                    </option>
                </select>
            </div>
          </div>
      </div>



      <div class="form-actions">
        <BaseButton variant="secondary" @click="$emit('navigate', 'workers')">Cancel</BaseButton>
        <BaseButton :loading="isSaving" type="submit">
          {{ isEdit ? 'Save Changes' : 'Register Worker' }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<style scoped>
.form-container {
  max-width: 800px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 32px;
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

.required {
    color: var(--color-danger, #ef4444);
}

.loading-state {
  padding: 48px;
  text-align: center;
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
