<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import PersonnelAssignment from '../../components/project/PersonnelAssignment.vue';

const TRADES = [
  { value: '1.1', label: '1.1 - Site Management (Ancillary Works)' },
  { value: '1.2', label: '1.2 - Site Support (Ancillary Works)' },
  { value: '1.3', label: '1.3 - General Machine Operation (Ancillary Works)' },
  { value: '1.4', label: '1.4 - Site Preparation (Ancillary Works)' },
  { value: '1.5', label: '1.5 - Scaffolding (Ancillary Works)' },
  { value: '2.1', label: '2.1 - Demolition (Civil & Structural Works)' },
  { value: '2.2', label: '2.2 - Earthworks (Civil & Structural Works)' },
  { value: '2.3', label: '2.3 - Foundation (Civil & Structural works)' },
  { value: '2.4', label: '2.4 - Tunnelling (Civil & Structural Works)' },
  { value: '2.5', label: '2.5 - Reinforced Concrete (Civil & Structural Works)' },
  { value: '2.6', label: '2.6 - Structural Steel (Civil & Structural Works)' },
  { value: '2.7', label: '2.7 - Mass Engineered Timber (Civil & Structural Works)' },
  { value: '2.8', label: '2.8 - Road & Drainage (Civil & Structural Works)' },
  { value: '3.1', label: '3.1 - Ceiling (Architectural Works)' },
  { value: '3.2', label: '3.2 - Partition Wall (Architectural Works)' },
  { value: '3.3', label: '3.3 - Floor (Architectural Works)' },
  { value: '3.4', label: '3.4 - Roofing (Architectural Works)' },
  { value: '3.5', label: '3.5 - Facade (Architectural Works)' },
  { value: '3.6', label: '3.6 - Door (Architectural Works)' },
  { value: '3.7', label: '3.7 - Window (Architectural Works)' },
  { value: '3.8', label: '3.8 - Finishes (Architectural Works)' },
  { value: '3.9', label: '3.9 - Waterproofing (Architectural Works)' },
  { value: '3.10', label: '3.10 - Joinery & Fixtures (Architectural Works)' },
  { value: '3.11', label: '3.11 - Landscaping (Architectural Works)' },
  { value: '4.1', label: '4.1 - Plumbing, Sanitary & Gas (Service Works)' },
  { value: '4.2', label: '4.2 - Fire Prevention & Protection (Service Works)' },
  { value: '4.3', label: '4.3 - Electrical (Service Works)' },
  { value: '4.4', label: '4.4 - Mechanical (Service Works)' },
  { value: '4.5', label: '4.5 - Lift & Escalator (Service Works)' },
  { value: '4.6', label: '4.6 - Prefab MEP (Service Works)' }
];

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
  offsite_fabricator_name: '',
  offsite_fabricator_uen: '',
  offsite_fabricator_location: '',
  worker_company_name: '',
  worker_company_uen: '',
  worker_company_client_name: '',
  worker_company_client_uen: '',
  worker_company_trade: '',
  status: 'active'
});

const sites = ref([]);
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
        
        const sitesData = await api.getSites({ user_id: userId });
        sites.value = sitesData || [];
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
        site_id: data.site_id || ''
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

const handleSubmit = async () => {
  isSaving.value = true;
  try {
    if (isEdit.value) {
      await api.updateProject(props.id, formData.value);
      notification.success('Project details updated');
    } else {
      await api.createProject(formData.value);
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
              <BaseInput v-model="formData.reference" label="Project Reference Number" placeholder="e.g., PRJ-2024-001" required />
              <BaseInput v-model="formData.title" label="Project Title" placeholder="e.g., Marina Bay Tower" required />
              
              <div class="form-group">
                  <label class="form-label">Assign to Site <span class="required">*</span></label>
                  <select v-model="formData.site_id" class="form-select" required>
                      <option value="">Select site</option>
                      <option v-for="s in sites" :key="s.site_id" :value="s.site_id">
                          {{ s.site_name }}
                      </option>
                  </select>
              </div>

              <BaseInput v-model="formData.location" label="Project Location Description" placeholder="e.g., Marina Bay, Central Singapore" required />
              <BaseInput v-model="formData.contract" label="Project Contract Number" placeholder="e.g., CNT-2024-MB" required />
              <BaseInput v-model="formData.contract_name" label="Project Contract Name" placeholder="e.g., Marina Bay Development" required />
              <BaseInput v-model="formData.hdb_precinct" label="HDB Precinct Name" placeholder="e.g., Marina Precinct (if applicable)" />
            </div>
          </div>

          <div class="form-section-card">
            <h3 class="section-title">Contractor & Workforce</h3>
            <div class="form-grid">
              <BaseInput v-model="formData.main_contractor_name" label="Main Contractor" placeholder="e.g., Mega Engineering Pte Ltd" />
              <BaseInput v-model="formData.main_contractor_uen" label="Main Contractor UEN" placeholder="e.g., 200012345X" />
              <BaseInput v-model="formData.worker_company_name" label="Worker Company Name" placeholder="e.g., WorkForce Solutions Pte Ltd" />
              <BaseInput v-model="formData.worker_company_uen" label="Worker Company UEN" placeholder="e.g., 201998765W" />
              <BaseInput v-model="formData.worker_company_client_name" label="Worker Company Client Name" placeholder="e.g., HDB Infrastructure" />
              <BaseInput v-model="formData.worker_company_client_uen" label="Worker Company Client UEN" placeholder="e.g., 196100018G" />
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

          <div class="form-section-card">
            <h3 class="section-title">Offsite Fabricator</h3>
            <div class="form-grid">
              <BaseInput v-model="formData.offsite_fabricator_name" label="Company Name" placeholder="e.g., Delta Fabrication Ltd" />
              <BaseInput v-model="formData.offsite_fabricator_uen" label="UEN" placeholder="e.g., UEN-FAB-001" />
              <BaseInput v-model="formData.offsite_fabricator_location" label="Location" placeholder="e.g., 10 Industrial Way, Singapore" class="full-width" />
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
