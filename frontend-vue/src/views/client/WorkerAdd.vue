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
  company_id: '',
  status: 'active'
});

const companies = ref([]);

const isEdit = computed(() => props.mode === 'edit');



const fetchCompanies = async () => {
    try {
        const savedUser = localStorage.getItem('auth_user');
        let tenantId = null;
        if (savedUser) {
            try {
                const user = JSON.parse(savedUser);
                tenantId = user.tenant_id || user.id;
            } catch (e) {}
        }
        const data = await api.getCompanies({ tenant_id: tenantId });
        companies.value = data || [];
    } catch (err) {
        console.error('Failed to load companies', err);
        notification.error('Could not load company list');
    }
};

const fetchWorker = async () => {
  if (!isEdit.value || !props.id) return;
  isLoading.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let contextTenantId = null;
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            contextTenantId = user.tenant_id || user.id;
        } catch (e) {}
    }
    const data = await api.getWorkerById(props.id, { tenant_id: contextTenantId }); // Pass tenant_id to API call
    if (data) {
      formData.value = { 
        ...data,
        company_id: data.company_id || ''
      };
    }
  } finally {
    isLoading.value = false;
  }
};

onMounted(async () => {
    await fetchCompanies();
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
    // prepare payload: ensure current_project_id is synced with the form selection
    const savedUser = localStorage.getItem('auth_user');
    let tenantId = null;
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            tenantId = user.tenant_id || user.id;
        } catch (e) {}
    }

    const payload = { 
        ...formData.value,
        tenant_id: tenantId
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

        <div class="form-group">
          <label class="form-label">Employer Company</label>
          <select v-model="formData.company_id" class="form-select">
            <option value="">Select company</option>
            <option v-for="c in companies" :key="c.company_id" :value="c.company_id">
                {{ c.company_name }}
            </option>
          </select>
        </div>

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

