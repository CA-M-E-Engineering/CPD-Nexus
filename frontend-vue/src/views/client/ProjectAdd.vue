<script setup>
import { ref, onMounted, computed } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import PersonnelAssignment from '../../components/project/PersonnelAssignment.vue';

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
              <BaseInput v-model="formData.worker_company_name" label="Worker Company" placeholder="e.g., WorkForce Solutions" />
              <BaseInput v-model="formData.worker_company_uen" label="Worker Company UEN" placeholder="e.g., 201998765W" />
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

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
}
</style>
