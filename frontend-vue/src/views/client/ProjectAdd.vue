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
  tenant_id: '',
  reference: '',
  title: '',
  site_id: '',
  location: '',
  contract: '',
  contract_name: '',
  hdb_precinct: '',
  main_contractor_id: '',
  offsite_fabricator_id: '',
  worker_company_id: '',
  worker_company_client_id: '',
  status: 'active'
});

const sites = ref([]);
const isEdit = computed(() => props.mode === 'edit');

const companies = ref([]);
const filteredContractors = computed(() => companies.value.filter(c => c.company_type === 'contractor'));
const filteredFabricators = computed(() => companies.value.filter(c => c.company_type === 'offsite_fabricator'));

const fetchData = async () => {
    try {
        const savedUser = localStorage.getItem('auth_user');
        let tenantId = 'tenant-client-1';
        if (savedUser) {
            try {
                const user = JSON.parse(savedUser);
                tenantId = user.tenant_id || user.id;
            } catch (e) {
                console.error("Failed to parse auth_user", e);
            }
        }
        formData.value.tenant_id = tenantId;
        
        const [sitesData, companiesData] = await Promise.all([
             api.getSites({ tenant_id: tenantId }),
             api.getCompanies({ tenant_id: tenantId })
        ]);
        
        sites.value = sitesData || [];
        companies.value = companiesData || [];
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
    }
  } finally {
    isLoading.value = false;
  }
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

    <form v-else class="form-container" @submit.prevent="handleSubmit">
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

        <div class="form-group">
            <label class="form-label">Main Contractor</label>
            <select v-model="formData.main_contractor_id" class="form-select">
                <option value="">Select Company</option>
                <option v-for="c in filteredContractors" :key="c.company_id" :value="c.company_id">
                    {{ c.company_name }} ({{ c.uen }})
                </option>
            </select>
        </div>

        <div class="form-group">
            <label class="form-label">Offsite Fabricator</label>
            <select v-model="formData.offsite_fabricator_id" class="form-select">
                <option value="">Select Fabricator</option>
                <option v-for="c in filteredFabricators" :key="c.company_id" :value="c.company_id">
                    {{ c.company_name }} ({{ c.uen }})
                </option>
            </select>
        </div>

        <div class="form-group">
            <label class="form-label">Worker Company</label>
            <select v-model="formData.worker_company_id" class="form-select">
                <option value="">Select Company</option>
                <option v-for="c in filteredContractors" :key="c.company_id" :value="c.company_id">
                    {{ c.company_name }} ({{ c.uen }})
                </option>
            </select>
        </div>

        <div class="form-group">
            <label class="form-label">Worker Company Client</label>
            <select v-model="formData.worker_company_client_id" class="form-select">
                <option value="">Select Company</option>
                <option v-for="c in filteredContractors" :key="c.company_id" :value="c.company_id">
                    {{ c.company_name }} ({{ c.uen }})
                </option>
            </select>
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
.form-container {
  max-width: 900px;
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

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
}
</style>

