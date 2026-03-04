<script setup>
import { ref, onMounted, computed } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import { MAP_MODES, USER_TYPES } from '../../utils/constants.js';
import UnifiedMap from '../../components/ui/UnifiedMap.vue';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';

const props = defineProps({
  id: [Number, String],
  mode: {
    type: String,
    default: 'add'
  }
});

const emit = defineEmits(['navigate']);

const isEdit = computed(() => props.mode === 'edit');
const isLoading = ref(false);
const isSaving = ref(false);

const formData = ref({
  user_name: '',
  user_type: USER_TYPES.CLIENT,
  username: '',
  password: '',
  email: '',
  phone: '',
  address: '',
  latitude: '',
  longitude: '',
  status: 'active',
  bridge_ws_url: '',
  bridge_auth_token: '',
  bridge_status: 'inactive',
  assigned_on_behalf_ofs: []
});

const availableOnBehalfOfs = ref([]);

const fetchUserAndAuthorisations = async () => {
  isLoading.value = true;
  try {
    const promises = [api.getPitstopAuthorisations()];
    if (isEdit.value && props.id) {
        promises.push(api.getUserById(props.id));
    }
    
    const results = await Promise.all(promises);
    const authsData = results[0] || [];
    // Extract unique onBehalfOfs from pitstop authorisations
    const onBehalfOfSet = new Set();
    const assignedToUser = new Set();
    
    authsData.forEach(auth => {
        if (auth.on_behalf_of_name) {
            onBehalfOfSet.add(auth.on_behalf_of_name);
            if (isEdit.value && props.id && auth.user_id === props.id) {
                assignedToUser.add(auth.on_behalf_of_name);
            }
        }
    });
    
    availableOnBehalfOfs.value = Array.from(onBehalfOfSet).sort();
    
    // Setup form data
    if (isEdit.value && results[1]) {
      const data = results[1];
      formData.value = {
        ...formData.value,
        user_name: data.user_name,
        user_type: data.user_type || USER_TYPES.CLIENT,
        username: data.username || '',
        email: data.email || '',
        phone: data.phone || '',
        address: data.address || '',
        latitude: data.lat || '',
        longitude: data.lng || '',
        status: data.status || 'active',
        bridge_ws_url: data.bridge_ws_url || '',
        bridge_auth_token: data.bridge_auth_token || '',
        bridge_status: data.bridge_status || 'inactive',
        assigned_on_behalf_ofs: Array.from(assignedToUser)
      };
    }
  } catch (err) {
      console.error('Failed to init user form:', err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchUserAndAuthorisations);

const handleSubmit = async () => {
  isSaving.value = true;
  try {
    const payload = {
        ...formData.value,
        lat: parseFloat(formData.value.latitude) || 0,
        lng: parseFloat(formData.value.longitude) || 0
    };

    let savedUserId = props.id;

    if (isEdit.value) {
      await api.updateUser(savedUserId, payload);
      notification.success('Organization updated successfully');
    } else {
      const response = await api.createUser(payload);
      savedUserId = response.user_id || response.id; // handle id retrieval based on create response mapping
      notification.success('New organization registered successfully');
    }
    
    // Now dispatch pitstop onBehalfOfs assignment array
    if (savedUserId && formData.value.assigned_on_behalf_ofs) {
        try {
            await api.assignPitstopOnBehalfOfs(savedUserId, formData.value.assigned_on_behalf_ofs);
        } catch (authErr) {
            console.error('Pitstop assignment failed:', authErr);
            notification.error('Organization saved, but Pitstop On Behalf Of assignments failed to sync.');
        }
    }

    emit('navigate', 'users');
  } catch (err) {
    console.error('[UserAdd] Save failed:', err);
    notification.error(err.message || 'Failed to save organization details');
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="user-add">
    <PageHeader
      :title="isEdit ? 'Edit Organization' : 'Register New User'"
      :description="isEdit ? 'Update details for this organization account' : 'Add a new client organization to the system'"
    />

    <div v-if="isLoading" class="loading-state">
      <p>Loading record details...</p>
    </div>

    <form v-else class="form-container" @submit.prevent="handleSubmit">
      <div class="form-info-banner">
        <i class="ri-information-line"></i>
        <span>Creating a User will automatically generate a <strong>Login Account</strong> and a corresponding <strong>Primary Business Entity</strong>.</span>
      </div>

      <!-- Three-column layout -->
      <div class="three-col-layout">

        <!-- LEFT: User info -->
        <div class="form-panel">
          <div class="form-grid">
            <div class="section-title full-width">Account &amp; Identity</div>
            <BaseInput v-model="formData.user_name" label="User Name" placeholder="e.g., Acme Logistics" class="full-width" required />
            <BaseInput v-model="formData.email" label="Contact Email" type="email" placeholder="contact@user.com" required />
            <BaseInput v-model="formData.username" label="Username" placeholder="e.g., acme_admin" required />
            <BaseInput v-model="formData.password" label="Password" type="password" placeholder="At least 6 characters" class="full-width" :required="!isEdit" />

            <div class="section-title full-width">Business Details</div>
            <BaseInput v-model="formData.phone" label="Contact Phone" placeholder="e.g., +65 1234 5678" class="full-width" />
            <BaseInput v-model="formData.address" label="Office Address" placeholder="e.g., 120 Lower Delta Road" class="full-width" />
            <BaseInput v-model="formData.latitude" label="Latitude" placeholder="e.g., 1.3521" type="number" step="any" />
            <BaseInput v-model="formData.longitude" label="Longitude" placeholder="e.g., 103.8198" type="number" step="any" />
          </div>
        </div>

        <!-- MIDDLE: Bridge config -->
        <div class="form-panel bridge-panel">
            <div class="section-title">Bridge Connection</div>
            <p class="bridge-description">Configure the IoT Bridge WebSocket connection for this organization. When active, the system will connect automatically on startup.</p>

            <div class="bridge-status-indicator" :class="formData.bridge_status === 'active' ? 'status-active' : 'status-inactive'">
              <i :class="formData.bridge_status === 'active' ? 'ri-wifi-line' : 'ri-wifi-off-line'"></i>
              <span>{{ formData.bridge_status === 'active' ? 'Bridge Active' : 'Bridge Inactive' }}</span>
            </div>

            <div class="bridge-fields">
              <BaseInput
                v-model="formData.bridge_ws_url"
                label="WebSocket URL"
                placeholder="ws://192.168.1.100:8081/ws"
              />
              <BaseInput
                v-model="formData.bridge_auth_token"
                label="Auth Token (Optional)"
                placeholder="Leave blank if not required"
              />
              <div class="form-group">
                <label class="form-label">Connection Status</label>
                <select v-model="formData.bridge_status" class="form-select">
                  <option value="active">Active — Connect on startup</option>
                  <option value="inactive">Inactive — Do not connect</option>
                </select>
                <span class="field-help">Set to Active to enable the IoT Bridge WebSocket connection for this organization's devices.</span>
              </div>
            </div>
          </div>
          
          <!-- CPD Data Integrations -->
          <div class="form-panel integration-panel">
              <div class="section-title">Submit CPD Data On Behalf Of</div>
              <p class="bridge-description">Assign one or more 'On Behalf Of' entities synced via configuration to assign datasets directly to this client organization.</p>
              <div class="form-group" style="margin-top: 16px;">
                  <label class="form-label">Assign Entities</label>
                  
                  <div class="dropdown-wrapper">
                    <details class="custom-dropdown">
                        <summary class="dropdown-summary">
                            <span>
                                {{ formData.assigned_on_behalf_ofs.length > 0 
                                   ? `${formData.assigned_on_behalf_ofs.length} entity(s) selected` 
                                   : 'Select entities...' }}
                            </span>
                            <i class="ri-arrow-down-s-line"></i>
                        </summary>
                        <div class="dropdown-content contractor-list">
                            <div v-if="availableOnBehalfOfs.length === 0" class="no-contractors">
                                No entities fetched yet. Go to Configuration -> Submit CPD Data to sync.
                            </div>
                            <label v-for="entity in availableOnBehalfOfs" :key="entity" class="checkbox-label">
                                <input type="checkbox" :value="entity" v-model="formData.assigned_on_behalf_ofs" />
                                <span>{{ entity }}</span>
                            </label>
                        </div>
                    </details>
                  </div>
                  <span class="field-help" style="margin-top: 8px; display: block;">Selected entities will automatically map their dataset authorizations to this user ID.</span>
              </div>
          </div>
      </div>

      <!-- Map full-width below both panels -->
      <div class="map-section">
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

      <div class="form-actions">
        <BaseButton variant="secondary" @click="$emit('navigate', 'users')">Cancel</BaseButton>
        <BaseButton :loading="isSaving" @click="handleSubmit" variant="primary">
          {{ isEdit ? 'Save Changes' : 'Register User' }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<style scoped>
/* Styles updated for premium look */
.form-container {
  max-width: 1600px;
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-lg);
  padding: 40px;
  box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.3);
}

.form-info-banner {
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: var(--radius-md);
  padding: 16px;
  margin-bottom: 32px;
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--color-text-primary);
  font-size: 14px;
}

.form-info-banner i {
  color: var(--color-accent);
  font-size: 1.25rem;
  flex-shrink: 0;
}

/* Three-column layout */
.three-col-layout {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 32px;
  margin-bottom: 32px;
  align-items: start;
}

@media (max-width: 1300px) {
  .three-col-layout {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 900px) {
  .three-col-layout {
    grid-template-columns: 1fr;
  }
}

.form-panel {
  display: flex;
  flex-direction: column;
  gap: 0;
}

/* Bridge panel styling */
.bridge-panel {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: var(--radius-lg);
  padding: 28px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.bridge-description {
  font-size: 13px;
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin: 0;
}

.bridge-status-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 600;
}

.status-active {
  background: rgba(34, 197, 94, 0.12);
  border: 1px solid rgba(34, 197, 94, 0.3);
  color: #22c55e;
}

.status-inactive {
  background: rgba(156, 163, 175, 0.1);
  border: 1px solid rgba(156, 163, 175, 0.2);
  color: var(--color-text-secondary);
}

.bridge-status-indicator i {
  font-size: 18px;
}

.bridge-fields {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Left panel inner grid */
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}

.section-title {
  font-size: 13px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--color-accent);
  letter-spacing: 0.05em;
  margin-top: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  margin-bottom: 4px;
}

.full-width {
  grid-column: 1 / -1;
}

/* Map */
.map-section {
  margin-bottom: 32px;
}

.form-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-secondary);
  display: block;
  margin-bottom: 10px;
}

.map-wrapper {
  height: 350px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-md);
  overflow: hidden;
}

/* Actions */
.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  padding-top: 28px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

/* Misc */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-select {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-md);
  padding: 12px;
  color: var(--color-text-primary);
  font-size: 14px;
  outline: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.form-select:focus {
  border-color: var(--color-accent);
  background: rgba(255, 255, 255, 0.08);
}

.field-help {
  font-size: 11px;
  color: var(--color-text-muted);
  line-height: 1.5;
}

.loading-state {
  padding: 80px;
  text-align: center;
  color: var(--color-text-secondary);
}

/* Integration Panel */
.integration-panel {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: var(--radius-lg);
  padding: 28px;
  display: flex;
  flex-direction: column;
}

.contractor-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 250px;
  overflow-y: auto;
  padding: 12px;
}

/* Custom Dropdown Styling */
.dropdown-wrapper {
    position: relative;
    width: 100%;
}

.custom-dropdown {
    width: 100%;
    position: relative;
}

.custom-dropdown[open] .dropdown-summary i {
    transform: rotate(180deg);
}

.dropdown-summary {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 14px;
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: var(--radius-md);
    color: var(--color-text-primary);
    font-size: 14px;
    cursor: pointer;
    list-style: none;
    transition: all 0.2s ease;
}

.dropdown-summary::-webkit-details-marker {
    display: none;
}

.dropdown-summary:hover {
    background: rgba(255, 255, 255, 0.08);
}

.dropdown-summary i {
    transition: transform 0.2s ease;
    color: var(--color-text-muted);
}

.dropdown-content {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    right: 0;
    background: #1e293b; /* Solid dark background to avoid transparency issues over other elements */
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: var(--radius-md);
    box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.5);
    z-index: 50;
}

.no-contractors {
    font-size: 13px;
    color: var(--color-text-muted);
    font-style: italic;
    text-align: center;
    padding: 12px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  transition: all 0.2s ease;
  background: rgba(255, 255, 255, 0.03);
}

.checkbox-label:hover {
  background: rgba(255, 255, 255, 0.08);
}

.checkbox-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
  accent-color: var(--color-accent);
  cursor: pointer;
}

.checkbox-label span {
  font-size: 14px;
  color: var(--color-text-primary);
  user-select: none;
}
</style>
