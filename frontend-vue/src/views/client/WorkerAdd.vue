<script setup>
import { ref, onMounted, computed } from 'vue';
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
  fin: '',
  role: 'worker',
  trade_code: '',
  email: '',
  company_name: '',
  status: 'active'
});

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
    const data = await api.getWorkerById(props.id, { user_id: contextuserId });
    if (data) {
      formData.value = { 
        ...data,
        company_name: data.company_name || ''
      };
    }
  } finally {
    isLoading.value = false;
  }
};

onMounted(async () => {
    await fetchWorker();
    
    // Set default role/status if new
    if (!isEdit.value) {
        if (!formData.value.role) formData.value.role = 'worker';
        if (!formData.value.status) formData.value.status = 'active';
    }
});

const handleSubmit = async () => {
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
      <div class="form-grid">
        <BaseInput v-model="formData.name" label="Full Name" placeholder="e.g., John Smith" required />
        <BaseInput v-model="formData.fin" label="FIN/NRIC" placeholder="e.g., S1234567D" required />
        

        
        <div class="form-group">
            <label class="form-label">Designated Role</label>
            <select v-model="formData.role" class="form-select">
                <option value="worker">Worker</option>
                <option value="pic">PIC (Person In Charge)</option>
                <option value="manager">Manager</option>
            </select>
        </div>

        <BaseInput v-model="formData.company_name" label="Employer Company" placeholder="e.g., Mega Engineering Pte Ltd" />

        <BaseInput v-model="formData.trade_code" label="Trade Code" placeholder="e.g., STR-01" />
        <BaseInput v-model="formData.email" label="Email Address" type="email" placeholder="john@company.com" />
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
  margin-bottom: 32px;
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
