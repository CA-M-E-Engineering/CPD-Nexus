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
  name: '',
  location: '',
  latitude: '', 
  longitude: '', 
});

const isEdit = computed(() => props.mode === 'edit');

const fetchSite = async () => {
  fetchError.value = null;
  if (!isEdit.value || !props.id) return;
  
  isLoading.value = true;
  try {
    let data = await api.getSiteById(props.id);
    
    // Robust parsing for potentially double-encoded or stringified responses
    let attempts = 0;
    while (typeof data === 'string' && attempts < 3) {
        try {
            data = data.trim();
            data = JSON.parse(data);
        } catch (e) {
            console.warn(`Parse attempt ${attempts} failed:`, e);
            break; 
        }
        attempts++;
    }

    if (data && typeof data === 'object') {
      formData.value = { 
        ...data,
        name: data.site_name || '', 
        location: data.location || '',
        latitude: data.lat || '',     
        longitude: data.lng || ''     
      };
    } else {
        fetchError.value = "Site data not found or invalid format.";
    }
  } catch (e) {
    console.error('Error fetching site:', e);
    fetchError.value = "Failed to load site details. Please check your connection.";
    notification.error("Failed to load site details.");
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchSite);
watch(() => props.id, fetchSite);

const handleSubmit = async () => {
  isSaving.value = true;
  try {
    // Get tenant_id from auth_user
    const authUser = JSON.parse(localStorage.getItem('auth_user') || '{}');
    const tenantId = authUser.tenant_id;

    if (!tenantId) {
        notification.error("Authentication error: Missing Tenant ID. Please re-login.");
        return;
    }

    const payload = {
        ...formData.value,
        site_name: formData.value.name, // Map name back to site_name
        tenant_id: tenantId,           // Add tenant_id
        // Convert lat/lng to numbers
        lat: parseFloat(formData.value.latitude) || 0,
        lng: parseFloat(formData.value.longitude) || 0
    };

    if (isEdit.value) {
      await api.updateSite(props.id, payload);
      notification.success('Site settings updated');
    } else {
      await api.createSite(payload);
      notification.success('New site created successfully');
    }
    emit('navigate', 'sites');
  } catch (err) {
    console.error('Failed to save site', err);
    notification.error(err.message || 'Failed to save site record');
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="site-add">
    <PageHeader 
      :title="isEdit ? 'Edit Site Settings' : 'Add New Site'" 
      :description="isEdit ? 'Update site parameters and contractor assignments' : 'Create a new physical location for operations'"
    />

    <div v-if="isLoading" class="loading-state">
      <p>Fetching site data...</p>
    </div>
    
    <div v-else-if="fetchError" class="error-state">
      <p class="error-text">{{ fetchError }}</p>
      <BaseButton size="sm" @click="fetchSite">Retry</BaseButton>
    </div>

    <div v-else>
      <form class="form-container" @submit.prevent="handleSubmit">
        <div class="form-grid">
          <BaseInput v-model="formData.name" label="Site Name" placeholder="e.g., Northshore Tunnel" required />
          <BaseInput v-model="formData.location" label="Location / Zone" placeholder="e.g., Woodlands" />
          
          <BaseInput v-model="formData.latitude" label="Latitude" placeholder="e.g., 1.3521" type="number" step="any" />
          <BaseInput v-model="formData.longitude" label="Longitude" placeholder="e.g., 103.8198" type="number" step="any" />

          <!-- Map Integration -->
          <div class="full-width map-section">
            <label class="form-label">Site Location</label>
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
          <BaseButton variant="secondary" @click="$emit('navigate', 'sites')">Cancel</BaseButton>
          <BaseButton :loading="isSaving" type="submit">
            {{ isEdit ? 'Save Changes' : 'Create Site' }}
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

.error-state {
  padding: 48px;
  text-align: center;
  background: var(--color-surface);
  border: 1px solid var(--color-danger);
  border-radius: var(--radius-md);
  margin-bottom: 24px;
}

.error-text {
  color: var(--color-danger);
  margin-bottom: 16px;
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

