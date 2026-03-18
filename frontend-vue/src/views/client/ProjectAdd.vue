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
  mode: { type: String, default: 'add' }
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
  submission_entity: 1, // 1: Onsite Builder, 2: Offsite Fabricator
  offsite_fabricator_name: '',
  offsite_fabricator_uen: '',
  offsite_fabricator_location: '',
  status: 'active'
});

const sites = ref([]);
const pitstopAuths = ref([]);
const formErrors = ref({});
const isEdit = computed(() => props.mode === 'edit');
const isOffsite = computed(() => formData.value.submission_entity === 2);
const selectedRegulator = computed(() => {
  if (!formData.value.pitstop_auth_id) return null;
  const auth = pitstopAuths.value.find(a => String(a.pitstop_auth_id) === String(formData.value.pitstop_auth_id));
  return auth ? auth.regulator_name : null;
});
const isBCA = computed(() => selectedRegulator.value && selectedRegulator.value.toUpperCase() === 'BCA');
const isHDB = computed(() => selectedRegulator.value && selectedRegulator.value.toUpperCase() === 'HDB');
const isLTA = computed(() => selectedRegulator.value && selectedRegulator.value.toUpperCase() === 'LTA');

// Null-safe string helper: converts null/undefined or literal "null" string to empty string
const ns = (v) => (v === null || v === undefined || v === 'null' ? '' : String(v));

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
            api.getPitstopAuthorisations().catch(() => [])
        ];

        const [sitesData, authsData] = await Promise.all(promises);
        sites.value = sitesData || [];
        pitstopAuths.value = authsData || [];
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
      // All string fields go through ns() to ensure no "null" strings in inputs
      formData.value = {
        ...formData.value,
        user_id:                   ns(data.user_id),
        reference:                 ns(data.reference),
        title:                     ns(data.title),
        site_id:                   ns(data.site_id),
        location:                  ns(data.location),
        contract:                  ns(data.contract),
        contract_name:             ns(data.contract_name),
        hdb_precinct:              ns(data.hdb_precinct),
        main_contractor_name:      ns(data.main_contractor_name),
        main_contractor_uen:       ns(data.main_contractor_uen),
        worker_company_name:       ns(data.worker_company_name),
        worker_company_uen:        ns(data.worker_company_uen),
        worker_company_client_name:ns(data.worker_company_client_name),
        worker_company_client_uen: ns(data.worker_company_client_uen),
        worker_company_trade:      ns(data.worker_company_trade),
        pitstop_auth_id:           ns(data.pitstop_auth_id),
        submission_entity:         data.submission_entity || 1,
        offsite_fabricator_name:   ns(data.offsite_fabricator_name),
        offsite_fabricator_uen:    ns(data.offsite_fabricator_uen),
        offsite_fabricator_location: ns(data.offsite_fabricator_location),
        status:                    ns(data.status) || 'active',
      };
      if (data.worker_company_trade) {
        selectedTrades.value = ns(data.worker_company_trade).split(',').map(s => s.trim()).filter(Boolean);
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

    if (!formData.value.title) {
        errors.title = 'Project title is required';
    } else if (formData.value.title.length > 1000) {
        errors.title = 'Title too long (max 1000 chars)';
    }

    // project_reference_number: Mandatory for Onsite Builder (entity=1), optional for Offsite Fabricator
    if (!isOffsite.value) {
        if (!formData.value.reference) {
            errors.reference = 'Project Reference Number is required for Onsite Builder';
        } else if (!validateProjectRef(formData.value.reference)) {
            errors.reference = 'Invalid format — expected: A1234-AB123-2022';
        }
        if (!formData.value.location) {
            errors.location = 'Project Location Description is required for Onsite Builder';
        }
    } else if (formData.value.reference && !validateProjectRef(formData.value.reference)) {
        errors.reference = 'Invalid format — expected: A1234-AB123-2022';
    }

    // contract number: optional, but validate format if provided
    if (formData.value.contract) {
        if (!validateHDBContract(formData.value.contract) && !validateLTAContract(formData.value.contract)) {
            errors.contract = 'Invalid contract format (HDB: D/NNNNN/YY max 10 chars, LTA: alphanumeric max 20 chars)';
        }
    }

    if (formData.value.main_contractor_uen && !validateUEN(formData.value.main_contractor_uen)) {
        errors.main_contractor_uen = 'Invalid UEN format';
    }

    // Main Contractor: Conditional Mandatory for BCA when Onsite (Fix 6)
    if (isBCA.value && !isOffsite.value) {
        if (!formData.value.main_contractor_name) {
            errors.main_contractor_name = 'Main Contractor name is mandatory for BCA Onsite submissions';
        }
        if (!formData.value.main_contractor_uen) {
            errors.main_contractor_uen = 'Main Contractor UEN is mandatory for BCA Onsite submissions';
        }
    }

    // person_employer_company_name & person_employer_company_uen are MANDATORY per API spec
    if (!formData.value.worker_company_name) {
        errors.worker_company_name = 'Person Employer Company Name is mandatory per API spec';
    }
    if (!formData.value.worker_company_uen) {
        errors.worker_company_uen = 'Person Employer Company UEN is mandatory per API spec';
    } else if (!validateUEN(formData.value.worker_company_uen)) {
        errors.worker_company_uen = 'Invalid UEN format';
    }

    if (formData.value.worker_company_client_uen && !validateUEN(formData.value.worker_company_client_uen)) {
        errors.worker_company_client_uen = 'Invalid UEN format';
    }

    // Offsite Fabricator fields: required only when entity = 2
    if (isOffsite.value) {
        if (!formData.value.offsite_fabricator_name) {
            errors.offsite_fabricator_name = 'Fabricator company name is required';
        }
        if (!formData.value.offsite_fabricator_uen) {
            errors.offsite_fabricator_uen = 'Fabricator UEN is required';
        } else if (!validateUEN(formData.value.offsite_fabricator_uen)) {
            errors.offsite_fabricator_uen = 'Invalid UEN format';
        }
        if (!formData.value.offsite_fabricator_location) {
            errors.offsite_fabricator_location = 'Fabricator location description is required';
        }
    } else {
        // Validate UEN format if provided even when Onsite
        if (formData.value.offsite_fabricator_uen && formData.value.offsite_fabricator_uen.trim() !== '' && !validateUEN(formData.value.offsite_fabricator_uen)) {
            errors.offsite_fabricator_uen = 'Invalid UEN format';
        }
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
        main_contractor_uen:       sanitizeUEN(formData.value.main_contractor_uen),
        worker_company_uen:        sanitizeUEN(formData.value.worker_company_uen),
        worker_company_client_uen: sanitizeUEN(formData.value.worker_company_client_uen),
        offsite_fabricator_uen:    sanitizeUEN(formData.value.offsite_fabricator_uen),
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
      :title="isEdit ? 'Edit Project Details' : 'Register New Project'"
      :description="isEdit ? 'Update existing construction project information' : 'Add a new construction project to your portfolio'"
    />

    <div v-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>Fetching project data...</p>
    </div>

    <form v-else class="project-form" @submit.prevent="handleSubmit">
      <div class="split-layout">
        <div class="form-side">

          <!-- ── SECTION 1: Project Identification ── -->
          <div class="form-section-card">
            <div class="section-header">
              <h3 class="section-title">Project Identification</h3>
              <p class="section-desc">Core project details and submission classification.</p>
            </div>

            <!-- Submission Entity (radio — at the top so it drives conditionality) -->
            <div class="form-group">
              <label class="form-label">Submission Entity <span class="required">*</span></label>
              <div class="entity-toggle">
                <label class="entity-option" :class="{ active: formData.submission_entity === 1 }">
                  <input type="radio" v-model.number="formData.submission_entity" :value="1" />
                  <div class="entity-icon"><i class="ri-building-2-line"></i></div>
                  <div>
                    <span class="entity-label">Onsite Builder</span>
                    <span class="entity-sub">Construction on project site</span>
                  </div>
                </label>
                <label class="entity-option" :class="{ active: formData.submission_entity === 2 }">
                  <input type="radio" v-model.number="formData.submission_entity" :value="2" />
                  <div class="entity-icon"><i class="ri-truck-line"></i></div>
                  <div>
                    <span class="entity-label">Offsite Fabricator</span>
                    <span class="entity-sub">Manufacturing at fabrication yard</span>
                  </div>
                </label>
              </div>
            </div>

            <div class="form-grid">
              <BaseInput
                v-model="formData.title"
                label="Project Title"
                placeholder="e.g., Proposed Construction of 50 Sty Mixed Commercial Building"
                :error="formErrors.title"
                required
              />
              <BaseInput
                v-model="formData.reference"
                label="Project Reference Number"
                :placeholder="isOffsite ? 'e.g., AE1234-AB123-2022 (optional for Offsite)' : 'e.g., AE1234-AB123-2022'"
                :error="formErrors.reference"
                :required="!isOffsite"
              />

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

              <BaseInput
                v-model="formData.location"
                label="Project Location Description"
                :placeholder="isOffsite ? 'e.g., 52 Jurong Gateway Road S608549 (optional for Offsite)' : 'e.g., 52 Jurong Gateway Road S608549'"
                :error="formErrors.location"
                :required="!isOffsite"
                class="full-width"
              />

              <BaseInput
                v-model="formData.contract"
                label="Project Contract Number"
                placeholder="e.g., D/12345/24 (HDB) or alphanumeric max 20 chars (LTA)"
                :error="formErrors.contract"
                :disabled="isBCA"
              >
                <template v-if="isBCA" #label-suffix>
                  <span class="opt-tag">Not needed for BCA</span>
                </template>
              </BaseInput>
              <BaseInput
                v-model="formData.contract_name"
                label="Project Contract Name"
                placeholder="e.g., Pasir Ris N2 C09"
                :error="formErrors.contract_name"
                :disabled="isBCA"
              >
                <template v-if="isBCA" #label-suffix>
                  <span class="opt-tag">Not needed for BCA</span>
                </template>
              </BaseInput>
              <BaseInput
                v-model="formData.hdb_precinct"
                label="HDB Precinct Name"
                placeholder="e.g., Golden Lily (HDB optional)"
                :error="formErrors.hdb_precinct"
                :disabled="!isHDB && formData.pitstop_auth_id"
              >
                <template v-if="!isHDB && formData.pitstop_auth_id" #label-suffix>
                  <span class="opt-tag">HDB only</span>
                </template>
              </BaseInput>
            </div>
          </div>

          <!-- ── SECTION 2: Offsite Fabricator (always visible) ── -->
          <div class="form-section-card" :class="{ 'offsite-active': isOffsite, 'offsite-inactive': !isOffsite }">
            <div class="section-header">
              <div>
                <h3 class="section-title">
                  <i class="ri-truck-line section-icon"></i>
                  Offsite Fabricator Details
                </h3>
                <p class="section-desc">
                  <span v-if="isOffsite" class="req-badge">Required — Offsite Fabricator selected</span>
                  <span v-else class="opt-badge">Optional — stored but not submitted when Onsite Builder is selected</span>
                </p>
              </div>
            </div>
            <div class="form-grid">
              <BaseInput
                v-model="formData.offsite_fabricator_name"
                label="Fabricator Company Name"
                placeholder="e.g., Prefab Solutions Pte Ltd"
                :error="formErrors.offsite_fabricator_name"
                :required="isOffsite"
              />
              <BaseInput
                v-model="formData.offsite_fabricator_uen"
                label="Fabricator UEN"
                placeholder="e.g., 202012345Z"
                :error="formErrors.offsite_fabricator_uen"
                :required="isOffsite"
              />
              <BaseInput
                v-model="formData.offsite_fabricator_location"
                label="Fabricator Location Description"
                placeholder="e.g., 10 Tuas South Ave 1 Singapore 637569"
                :error="formErrors.offsite_fabricator_location"
                :required="isOffsite"
                class="full-width"
              />
            </div>
          </div>

          <!-- ── SECTION 3: Contractor & Workforce ── -->
          <div class="form-section-card">
            <div class="section-header">
              <h3 class="section-title">Contractor &amp; Workforce</h3>
              <p class="section-desc">Employment and subcontracting details linked to submitted manpower records.</p>
            </div>

            <!-- API Mandatory block -->
            <div class="field-block mandatory-block">
              <div class="block-label">
                <i class="ri-error-warning-line"></i>
                Mandatory — required for every submission
              </div>
              <div class="form-grid">
                <BaseInput
                  v-model="formData.worker_company_name"
                  label="Person Employer Company Name"
                  placeholder="e.g., WorkForce Solutions Pte Ltd"
                  :error="formErrors.worker_company_name"
                  required
                />
                <BaseInput
                  v-model="formData.worker_company_uen"
                  label="Person Employer Company UEN"
                  placeholder="e.g., 201998765W"
                  :error="formErrors.worker_company_uen"
                  required
                />
              </div>
            </div>

            <!-- API Optional block -->
            <div class="field-block optional-block">
              <div class="block-label">
                <i class="ri-information-line"></i>
                Optional — but regulator-specific rules apply (see badges)
              </div>
              <div class="form-grid">
                <BaseInput 
                  v-model="formData.main_contractor_name" 
                  label="Main Contractor Company Name" 
                  placeholder="e.g., Mega Engineering Pte Ltd"
                  :required="isBCA && !isOffsite"
                  :error="formErrors.main_contractor_name"
                >
                  <template v-if="isBCA && !isOffsite" #label-suffix>
                    <span class="reg-badge bca">BCA Mandatory</span>
                  </template>
                </BaseInput>
                <BaseInput 
                  v-model="formData.main_contractor_uen" 
                  label="Main Contractor Company UEN" 
                  placeholder="e.g., 200012345X" 
                  :required="isBCA && !isOffsite"
                  :error="formErrors.main_contractor_uen"
                >
                  <template v-if="isBCA && !isOffsite" #label-suffix>
                    <span class="reg-badge bca">BCA Mandatory</span>
                  </template>
                </BaseInput>

                <BaseInput 
                  v-model="formData.worker_company_client_name" 
                  label="Person Employer Client Company Name"
                  placeholder="e.g., HDB Infrastructure Pte Ltd" 
                  :error="formErrors.worker_company_client_name" 
                >
                  <template #label-suffix>
                    <span class="reg-badge bca">BCA Mandatory</span>
                  </template>
                </BaseInput>
                <BaseInput 
                  v-model="formData.worker_company_client_uen" 
                  label="Person Employer Client Company UEN"
                  placeholder="e.g., 196100018G" 
                  :error="formErrors.worker_company_client_uen" 
                >
                  <template #label-suffix>
                    <span class="reg-badge bca">BCA Mandatory</span>
                  </template>
                </BaseInput>

                <div class="form-group full-width">
                  <label class="form-label">
                    Person Employer Company Trade(s)
                    <span class="reg-badge bca">BCA Mandatory</span>
                    <span class="reg-badge lta">LTA Mandatory</span>
                    <span class="opt-tag">HDB: not needed</span>
                  </label>
                  <div class="trade-selection-area">
                    <select v-model="currentTradeSelection" class="form-select" @change="addTrade">
                      <option value="" disabled>Select a trade to add...</option>
                      <option v-for="t in availableTrades" :key="t.value" :value="t.value">
                        {{ t.label }}
                      </option>
                    </select>
                    <div class="selected-trades-container" v-if="selectedTrades.length > 0">
                      <span v-for="trade in selectedTrades" :key="trade" class="trade-pill">
                        {{ getTradeLabel(trade) }}
                        <button type="button" class="remove-trade-btn" @click="removeTrade(trade)">
                          <i class="ri-close-line"></i>
                        </button>
                      </span>
                    </div>
                    <p v-else class="no-trades-hint">No trades selected yet. Use the dropdown above to add trades.</p>
                  </div>
                </div>
              </div>
            </div>
          </div>


        </div>

        <!-- Personnel panel (edit mode only) -->
        <div v-if="isEdit && props.id" class="personnel-side">
          <PersonnelAssignment
            :project-id="props.id"
            :user-id="formData.user_id"
          />
        </div>
      </div>

      <div class="form-actions">
        <BaseButton variant="secondary" @click="$emit('navigate', 'projects')" type="button">Cancel</BaseButton>
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
  gap: 20px;
  min-width: 0;
}

.personnel-side {
  width: 380px;
  flex-shrink: 0;
  position: sticky;
  top: 24px;
}

/* ── Section Cards ── */
.form-section-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 24px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.offsite-active {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(var(--accent-rgb), 0.08);
}

.offsite-inactive {
  border-style: dashed;
  opacity: 0.85;
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
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-icon {
  font-style: normal;
  font-size: 16px;
  color: var(--color-accent);
}

.section-desc {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin: 0;
}

.req-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-accent);
  background: rgba(var(--accent-rgb), 0.1);
  padding: 2px 8px;
  border-radius: 99px;
}
.req-badge::before { content: '●'; font-size: 8px; }

.opt-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--color-text-secondary);
  background: var(--color-bg);
  padding: 2px 8px;
  border-radius: 99px;
  border: 1px solid var(--color-border);
}
.opt-badge::before { content: '○'; font-size: 8px; }

/* ── Submission Entity Toggle ── */
.entity-toggle {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-top: 8px;
}

.entity-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s;
  background: var(--color-bg);
  user-select: none;
}

.entity-option:hover {
  border-color: var(--color-accent);
  background: rgba(var(--accent-rgb), 0.04);
}

.entity-option.active {
  border-color: var(--color-accent);
  background: rgba(var(--accent-rgb), 0.08);
}

.entity-option input[type="radio"] {
  display: none;
}

.entity-icon {
  font-size: 22px;
  color: var(--color-accent);
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: rgba(var(--accent-rgb), 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.entity-label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.entity-sub {
  display: block;
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-top: 2px;
}

/* ── Form Grid ── */
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
}

.full-width {
  grid-column: 1 / -1;
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
}

.required {
  color: #ef4444;
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
  transition: border-color 0.15s;
}

.form-select:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(var(--accent-rgb), 0.1);
}

.help-text {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-top: 2px;
}

/* ── Mandatory / Optional field blocks ── */
.field-block {
  border-radius: var(--radius-sm);
  padding: 16px;
  margin-bottom: 16px;
}

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

.mandatory-block .block-label {
  color: #ea580c;
}

.optional-block .block-label {
  color: var(--color-text-secondary);
}

.opt-tag {
  display: inline-block;
  font-size: 10px;
  font-weight: 500;
  color: var(--color-text-secondary);
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  padding: 1px 6px;
  border-radius: 4px;
  margin-left: 4px;
  text-transform: lowercase;
  letter-spacing: 0;
}

.reg-badge {
  display: inline-block;
  font-size: 10px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 4px;
  margin-left: 4px;
  letter-spacing: 0;
  text-transform: none;
}

.reg-badge.bca {
  background: rgba(234, 88, 12, 0.1);
  color: #ea580c;
  border: 1px solid rgba(234, 88, 12, 0.3);
}

.reg-badge.lta {
  background: rgba(37, 99, 235, 0.08);
  color: #2563eb;
  border: 1px solid rgba(37, 99, 235, 0.25);
}


/* ── Trades ── */
.trade-selection-area {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.selected-trades-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.trade-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: rgba(var(--accent-rgb), 0.08);
  border: 1px solid rgba(var(--accent-rgb), 0.25);
  color: var(--color-accent);
  padding: 5px 10px;
  border-radius: 99px;
  font-size: 12px;
  font-weight: 500;
}

.remove-trade-btn {
  background: none;
  border: none;
  padding: 0;
  display: inline-flex;
  align-items: center;
  color: inherit;
  opacity: 0.6;
  cursor: pointer;
  border-radius: 50%;
  transition: opacity 0.15s;
}

.remove-trade-btn:hover { opacity: 1; }

.no-trades-hint {
  font-size: 12px;
  color: var(--color-text-secondary);
  font-style: italic;
  margin: 0;
}

/* ── Loading ── */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 48px;
  gap: 16px;
}

.loading-spinner {
  width: 36px;
  height: 36px;
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
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
}

/* ── Responsive ── */
@media (max-width: 1200px) {
  .split-layout { flex-direction: column; }
  .personnel-side { width: 100%; position: static; }
}

@media (max-width: 768px) {
  .entity-toggle { grid-template-columns: 1fr; }
  .form-grid { grid-template-columns: 1fr; }
}
</style>
