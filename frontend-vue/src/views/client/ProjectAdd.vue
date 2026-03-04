<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import PersonnelAssignment from '../../components/project/PersonnelAssignment.vue';
import { TRADES } from '../../utils/constants.js';
import { validateProjectRef, validateUEN, validateHDBContract, validateLTAContract, sanitizeUEN } from '../../utils/validation.js';


const props = defineProps({
  id: [Number, String],
  mode: { type: String, default: 'add' } // 'add' or 'edit'
});

const emit = defineEmits(['navigate']);

const isSaving = ref(false);
const isLoading = ref(false);

const formData = ref({
  user_id: '',
  reference: '',
  title: '',
  site_id: '',
  location: '',
  contract: '',
  contract_name: '',
  hdb_precinct: '',
  main_contractor_name: '',
  main_contractor_uen: '',

  worker_company_name: '',
  worker_company_uen: '',
  worker_company_client_name: '',
  worker_company_client_uen: '',
  worker_company_trade: '',
  pitstop_auth_id: '',
  status: 'active'
});

const sites = ref([]);
const pitstopAuths = ref([]);
const formErrors = ref({});
const isEdit = computed(() => props.mode === 'edit');

const fetchData = async () => {
    try {
        const savedUser = localStorage.getItem('auth_user');
        let userId = 'User-client-1';
        if (savedUser) {
            try {
                const user = JSON.parse(savedUser);
                userId = user.user_id || user.id;
            } catch (e) {
                console.error("Failed to parse auth_user", e);
            }
        }
        formData.value.user_id = userId;
        
        const promises = [
            api.getSites({ user_id: userId }),
            api.getPitstopAuthorisations().catch(() => []) // Fallback on fail
        ];
        
        const [sitesData, authsData] = await Promise.all(promises);
        sites.value = sitesData || [];
        
        // Filter pitstop authorisations for this specific user
        pitstopAuths.value = (authsData || []).filter(auth => auth.user_id === userId);
    } catch (err) {
        console.error('Failed to load dependency data', err);
    }
};

const fetchProject = async () => {
  if (!isEdit.value || !props.id) return;
  isLoading.value = true;
  try {
    const data = await api.getProjectById(props.id);
    if (data) {
      formData.value = { 
        ...data,
        site_id: data.site_id || '',
        pitstop_auth_id: data.pitstop_auth_id || ''
      };
      if (data.worker_company_trade) {
        selectedTrades.value = data.worker_company_trade.split(',').map(s => s.trim()).filter(Boolean);
      }
    }
  } finally {
    isLoading.value = false;
  }
};

const selectedTrades = ref([]);

watch(selectedTrades, (newVal) => {
  formData.value.worker_company_trade = newVal.join(', ');
}, { deep: true });

const currentTradeSelection = ref('');

const availableTrades = computed(() => {
  return TRADES.filter(t => !selectedTrades.value.includes(t.value));
});

const getTradeLabel = (val) => {
  const t = TRADES.find(t => t.value === val);
  return t ? t.label : val;
};

const addTrade = () => {
  if (currentTradeSelection.value && !selectedTrades.value.includes(currentTradeSelection.value)) {
    selectedTrades.value.push(currentTradeSelection.value);
  }
  currentTradeSelection.value = '';
};

const removeTrade = (val) => {
  selectedTrades.value = selectedTrades.value.filter(t => t !== val);
};

onMounted(async () => {
    await fetchData();
    await fetchProject();
});

const validateForm = () => {
    const errors = {};
    if (!formData.value.reference) {
        errors.reference = 'Project reference is required';
    } else if (!validateProjectRef(formData.value.reference)) {
        errors.reference = 'Invalid format. Expected: A1234-12345-2022';
    }

    if (!formData.value.title) {
        errors.title = 'Project title is required';
    } else if (formData.value.title.length > 1000) {
        errors.title = 'Title too long (max 1000)';
    }

    if (formData.value.contract) {
        const isHDB = validateHDBContract(formData.value.contract);
        const isLTA = validateLTAContract(formData.value.contract);
        if (!isHDB && !isLTA) {
            errors.contract = 'Invalid contract format (HDB: D/NNNNN/YY, LTA: Max 20 chars)';
        }
    }

    if (formData.value.main_contractor_uen && !validateUEN(formData.value.main_contractor_uen)) {
        errors.main_contractor_uen = 'Invalid UEN format';
    }
    if (formData.value.worker_company_uen && !validateUEN(formData.value.worker_company_uen)) {
        errors.worker_company_uen = 'Invalid UEN format';
    }
    if (formData.value.worker_company_client_uen && !validateUEN(formData.value.worker_company_client_uen)) {
        errors.worker_company_client_uen = 'Invalid UEN format';
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
    const dataToSave = {
        ...formData.value,
        main_contractor_uen: sanitizeUEN(formData.value.main_contractor_uen),
        worker_company_uen: sanitizeUEN(formData.value.worker_company_uen),
        worker_company_client_uen: sanitizeUEN(formData.value.worker_company_client_uen),
    };
    
    if (isEdit.value) {
      await api.updateProject(props.id, dataToSave);
      notification.success('Project details updated');
    } else {
      await api.createProject(dataToSave);
      notification.success('New project created successfully');
    }
    emit('navigate', 'projects');
  } catch (err) {
    console.error('Failed to save project', err);
    notification.error(err.message || 'Failed to save project record');
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="project-add">
    <PageHeader 
      :title="isEdit ? 'Edit Project details' : 'Register New Project'" 
      :description="isEdit ? 'Update existing construction project information' : 'Add a new construction project to your portfolio'"
    />

    <div v-if="isLoading" class="loading-state">
      <p>Fetching project data...</p>
    </div>

    <form v-else class="project-form" @submit.prevent="handleSubmit">
      <div class="split-layout">
        <div class="form-side">
          <div class="form-section-card">
            <h3 class="section-title">Project Details</h3>
            <div class="form-grid">
              <BaseInput v-model="formData.reference" label="Project Reference Number" placeholder="e.g., A1234-12345-2022" :error="formErrors.reference" required />
              <BaseInput v-model="formData.title" label="Project Title" placeholder="e.g., Marina Bay Tower" :error="formErrors.title" required />
              
              <div class="form-group">
                  <label class="form-label">Assign to Site <span class="required">*</span></label>
                  <select v-model="formData.site_id" class="form-select" required>
                      <option value="">Select site</option>
                      <option v-for="s in sites" :key="s.site_id" :value="s.site_id">
                          {{ s.site_name }}
                      </option>
                  </select>
              </div>

              <div class="form-group">
                  <label class="form-label">Submit CPD Data On Behalf Of</label>
                  <select v-model="formData.pitstop_auth_id" class="form-select">
                      <option value="">None</option>
                      <option v-for="pa in pitstopAuths" :key="pa.pitstop_auth_id" :value="pa.pitstop_auth_id">
                          {{ pa.on_behalf_of_name }}
                      </option>
                  </select>
                  <span class="help-text">Links this project to CPD Data submission if selected.</span>
              </div>

              <BaseInput v-model="formData.location" label="Project Location Description" placeholder="e.g., Marina Bay, Central Singapore" :error="formErrors.location" required />
              <BaseInput v-model="formData.contract" label="Project Contract Number" placeholder="HDB Case: D/12345/24" :error="formErrors.contract" required />
              <BaseInput v-model="formData.contract_name" label="Project Contract Name" placeholder="e.g., Marina Bay Development" :error="formErrors.contract_name" required />
              <BaseInput v-model="formData.hdb_precinct" label="HDB Precinct Name" placeholder="e.g., Marina Precinct (if applicable)" :error="formErrors.hdb_precinct" />
            </div>
          </div>

          <div class="form-section-card">
            <h3 class="section-title">Contractor & Workforce</h3>
            <div class="form-grid">
              <BaseInput v-model="formData.main_contractor_name" label="Main Contractor" placeholder="e.g., Mega Engineering Pte Ltd" />
              <BaseInput v-model="formData.main_contractor_uen" label="Main Contractor UEN" placeholder="e.g., 200012345X" :error="formErrors.main_contractor_uen" />
              <BaseInput v-model="formData.worker_company_name" label="Worker Company Name" placeholder="e.g., WorkForce Solutions Pte Ltd" />
              <BaseInput v-model="formData.worker_company_uen" label="Worker Company UEN" placeholder="e.g., 201998765W" :error="formErrors.worker_company_uen" />
              <BaseInput v-model="formData.worker_company_client_name" label="Worker Company Client Name" placeholder="e.g., HDB Infrastructure" />
              <BaseInput v-model="formData.worker_company_client_uen" label="Worker Company Client UEN" placeholder="e.g., 196100018G" :error="formErrors.worker_company_client_uen" />
              <div class="form-group full-width">
                  <label class="form-label">Worker Company Trade(s)</label>
                  <div class="trade-selection-area">
                    <select v-model="currentTradeSelection" class="form-select" @change="addTrade">
                        <option value="" disabled>Select a trade to add...</option>
                        <option v-for="t in availableTrades" :key="t.value" :value="t.value">
                            {{ t.label }}
                        </option>
                    </select>
                    <div class="selected-trades-container">
                      <span v-for="trade in selectedTrades" :key="trade" class="trade-pill">
                        {{ getTradeLabel(trade) }}
                        <button type="button" class="remove-trade-btn" @click="removeTrade(trade)">
                          <i class="ri-close-line"></i>
                        </button>
                      </span>
                    </div>
                  </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="isEdit && props.id" class="personnel-side">
          <PersonnelAssignment 
            :project-id="props.id" 
            :user-id="formData.user_id" 
          />
        </div>
      </div>

      <div class="form-actions">
        <BaseButton variant="secondary" @click="$emit('navigate', 'projects')">Cancel</BaseButton>
        <BaseButton :loading="isSaving" type="submit">
          {{ isEdit ? 'Save Changes' : 'Create Project' }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<style scoped>
.project-form {
  width: 100%;
}

.split-layout {
  display: flex;
  gap: 24px;
  align-items: flex-start;
}

.form-side {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.personnel-side {
  flex: 1;
  position: sticky;
  top: 24px;
}

.form-section-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 24px;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 20px 0;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--color-border);
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.loading-state {
  padding: 48px;
  text-align: center;
}

@media (max-width: 1200px) {
  .split-layout {
    flex-direction: column;
  }
  .personnel-side {
    position: static;
    width: 100%;
  }
}

@media (max-width: 640px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 0;
}

.required {
  color: #ef4444;
}

.full-width {
  grid-column: 1 / -1;
}

.form-select {
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: 10px 12px;
  color: var(--color-text-primary);
  font-size: 14px;
  outline: none;
  cursor: pointer;
}

.form-select:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(var(--accent-rgb), 0.1);
}

.trade-selection-area {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.selected-trades-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.trade-pill {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: var(--color-surface-hover, #f3f4f6);
  border: 1px solid var(--color-border);
  padding: 6px 12px;
  border-radius: 16px;
  font-size: 13px;
  color: var(--color-text-primary);
}

.remove-trade-btn {
  background: none;
  border: none;
  padding: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.2s;
}

.remove-trade-btn:hover {
  color: #ef4444;
}

.help-text {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-top: 4px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
}
</style>
