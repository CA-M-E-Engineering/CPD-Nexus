<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import { MAP_MODES } from '../../utils/constants.js';
import UnifiedMap from '../../components/ui/UnifiedMap.vue';
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
const fetchError = ref(null);

const formData = ref({
  company_name: '',
  uen: '',
  company_type: 'contractor',
  address: '',
  latitude: '',
  longitude: '',
  status: 'active'
});

const isEdit = computed(() => props.mode === 'edit');

const fetchCompany = async () => {
  if (!isEdit.value || !props.id) return;
  
  isLoading.value = true;
  try {
    const data = await api.getCompanyById(props.id);
    if (data) {
      formData.value = { 
        ...data,
        latitude: data.latitude || '',
        longitude: data.longitude || ''
      };
    }
  } catch (e) {
    console.error('Error fetching company:', e);
    notification.error("Failed to load company details.");
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchCompany);
watch(() => props.id, fetchCompany);

const handleSubmit = async () => {
  isSaving.value = true;
  try {
    const savedUser = localStorage.getItem('auth_user');
    let tenantId = '';
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            tenantId = user.tenant_id || user.id;
        } catch (e) { console.error(e); }
    }

    if (!tenantId) {
        throw new Error('Tenant identity missing. Please logout and login again.');
    }

    const payload = {
        ...formData.value,
        tenant_id: tenantId,
        latitude: parseFloat(formData.value.latitude) || 0,
        longitude: parseFloat(formData.value.longitude) || 0
    };

    console.log('[CompanyAdd] Submitting payload:', payload);

    if (isEdit.value) {
      await api.updateCompany(props.id, payload);
      notification.success('Company profile updated');
    } else {
      await api.createCompany(payload);
      notification.success('New company registered successfully');
    }
    emit('navigate', 'companies');
  } catch (err) {
    console.error('Failed to save company', err);
    notification.error(err.message || 'Failed to save company record');
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="company-add">
    <PageHeader 
      :title="isEdit ? 'Edit Company Profile' : 'Register New Company'" 
      :description="isEdit ? 'Update business identity and location parameters' : 'Onboard a new construction partner or fabricator'"
    />

    <div v-if="isLoading" class="loading-state">
      <p>Fetching company data...</p>
    </div>

    <div v-else>
      <form class="form-container" @submit.prevent="handleSubmit">
        <div class="form-grid">
          <BaseInput v-model="formData.company_name" label="Company Name" placeholder="e.g., BuildMaster Solutions" class="full-width" required />
          
          <BaseInput v-model="formData.uen" label="UEN / Reg No." placeholder="e.g., 202400000A" required />
          
          <div class="form-group">
            <label class="form-label">Entity Type <span class="required">*</span></label>
            <select v-model="formData.company_type" class="form-select">
                <option value="contractor">Contractor</option>
                <option value="offsite_fabricator">Offsite Fabricator</option>
            </select>
          </div>

          <BaseInput v-model="formData.address" label="Office Address" placeholder="e.g., 12 Raffles Place" class="full-width" />
          
          <BaseInput v-model="formData.latitude" label="Latitude" placeholder="e.g., 1.3521" type="number" step="any" />
          <BaseInput v-model="formData.longitude" label="Longitude" placeholder="e.g., 103.8198" type="number" step="any" />

          <!-- Map Integration -->
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
          <BaseButton variant="secondary" @click="$emit('navigate', 'companies')">Cancel</BaseButton>
          <BaseButton :loading="isSaving" type="submit">
            {{ isEdit ? 'Save Changes' : 'Register Company' }}
          </BaseButton>
        </div>
      </form>
    </div>
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

.loading-state {
  padding: 48px;
  text-align: center;
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
    color: #ef4444;
}

.full-width {
  grid-column: 1 / -1;
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

.map-section {
    margin-top: 16px;
    margin-bottom: 16px;
}

.map-wrapper {
    height: 400px;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    overflow: hidden;
}
</style>
