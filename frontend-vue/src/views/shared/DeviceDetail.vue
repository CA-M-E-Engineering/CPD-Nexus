<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { api } from '../../services/api.js';
import { useDeviceStore } from '../../features/devices/store/deviceStore';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import DetailCard from '../../components/ui/DetailCard.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';

const props = defineProps({
  id: [Number, String],
  role: {
    type: String,
    default: 'client'
  }
});

const emit = defineEmits(['navigate']);
const deviceStore = useDeviceStore();

const isAdmin = computed(() => props.role === 'manager');
const device = ref(null);
const isLoading = ref(true);
const isSaving = ref(false);
const fetchError = ref(null);
const sites = ref([]);
const tenants = ref([]);

const fetchDevice = async () => {
  if (!props.id) {
    console.error('[DeviceDetail] Missing device ID');
    isLoading.value = false;
    return;
  }

  isLoading.value = true;
  fetchError.value = null;
  
  try {
    const data = await deviceStore.getDeviceById(props.id);
    
    if (!data) {
      fetchError.value = `Device with ID ${props.id} not found in registry.`;
      device.value = null;
    } else {
      device.value = data;
    }
  } catch (err) {
    console.error('[DeviceDetail] Fetch failed:', err);
    fetchError.value = 'A system error occurred while synchronizing with the registry.';
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchDevice);
onMounted(async () => {
  try {
    const tenantsData = await api.getTenants();
    tenants.value = tenantsData;
  } catch (err) {
    console.error('[DeviceDetail] Failed to fetch tenants:', err);
  }

  if (!isAdmin.value) {
    try {
      const savedUser = localStorage.getItem('auth_user');
      let tenantId = 'tenant-client-1'; // Default strong fallback
      
      if (savedUser) {
          try {
              const user = JSON.parse(savedUser);
              // Only use user.tenant_id if present. Do NOT fallback to user.id which might be a worker UUID.
              if (user.tenant_id) {
                  tenantId = user.tenant_id;
              } else if (user.role === 'client' && user.id) {
                  // If role is client, id is likely tenant_id
                  tenantId = user.id;
              }
          } catch (e) {
              console.error("Failed to parse auth_user", e);
          }
      }
      
      sites.value = await api.getSites({ tenant_id: tenantId });
    } catch (err) {
      console.error('[DeviceDetail] Failed to fetch sites:', err);
    }
  }
});

watch(() => props.id, (newId) => {
  if (newId) {
    console.log(`[DeviceDetail] ID changed to ${newId}, refetching...`);
    fetchDevice();
  }
});

const deviceInfo = computed(() => [
  { label: 'Device ID', value: device.value?.device_id || '---' },
  { label: 'Hardware SN', value: device.value?.sn || '---' },
  { label: 'Device Type', value: device.value?.model || '---' },
  { label: 'Firmware', value: 'v2.4.1-stable' },
  { label: 'Registry Date', value: 'Jan 20, 2024' }
]);

const assignmentInfo = computed(() => [
  { label: 'Current Tenant', value: device.value?.tenant_name || 'Internal Inventory' },
  { label: 'Assigned Site', value: device.value?.site_name || 'Unassigned' },
  { label: 'Last Deployed', value: 'Jan 22, 2024' }
]);

const activityStatus = computed(() => [
  { label: 'Last Seen', value: device.value?.status === 'online' ? 'Just now' : '2 hours ago' },
  { label: 'Battery Level', value: (device.value?.battery || 0) + '%' },
  { label: 'Signal Quality', value: 'Good (-68 dBm)' }
]);

const showAssignModal = ref(false);
const selectedTenantId = ref('');

const onTenantChange = (event) => {
  selectedTenantId.value = event.target.value;
};

const showAssignSiteModal = ref(false);
const selectedSiteId = ref('');

const handleAssign = async () => {
  console.log('[DeviceDetail] handleAssign triggered. Selected Tenant:', selectedTenantId.value);
  if (!selectedTenantId.value || selectedTenantId.value === '') {
    notification.error('Please select an organization from the dropdown.');
    return;
  }
  if (!device.value || !device.value.device_id) {
    notification.error('System error: Device data is not loaded.');
    return;
  }

  isSaving.value = true;
  try {
    await deviceStore.assignDeviceToTenant(selectedTenantId.value, [device.value.device_id]);
    
    notification.success('Device reassigned successfully.');
    showAssignModal.value = false;
    selectedTenantId.value = ''; // Reset selection
    await fetchDevice();
  } catch (err) {
    console.error('[DeviceDetail] Reassignment failed:', err);
    notification.error(err.message || 'The system could not process this request.');
  } finally {
    isSaving.value = false;
  }
};

const handleAssignSite = async () => {
  if (!selectedSiteId.value || !device.value) {
    notification.error('Please select a site before confirming.');
    return;
  }
  isSaving.value = true;
  try {
    await deviceStore.updateDevice(device.value.device_id, {
      sn: device.value.sn || '',
      model: device.value.model || '',
      status: device.value.status || 'offline',
      site_id: selectedSiteId.value
    });
    showAssignSiteModal.value = false;
    notification.success('Device assigned to site successfully.');
    await fetchDevice();
  } catch (err) {
    console.error('[DeviceDetail] Site assignment failed:', err);
    notification.error('Failed to assign site: ' + (err.message || 'System error'));
  } finally {
    isSaving.value = false;
  }
};

const handleEdit = () => {
  if (device.value) {
    emit('navigate', 'device-add', { id: device.value.device_id, mode: 'edit' });
  }
};
</script>

<template>
  <div class="device-detail">
    <PageHeader 
      :title="'Device Detail: ' + (device?.device_id || '...')" 
      description="View real-time telemetry and management parameters"
      variant="detail"
    >
      <template #toolbar-left>
        <BaseButton variant="ghost" size="sm" @click="$emit('navigate', 'devices')">
          <template #icon><i class="ri-arrow-left-line"></i></template>
          Back to Registry
        </BaseButton>
      </template>
      <template #toolbar-right>
        <BaseButton v-if="isAdmin" variant="secondary" size="sm" icon="ri-edit-line" @click="handleEdit">Edit Device</BaseButton>
        <BaseButton size="sm" icon="ri-history-line">Activity Logs</BaseButton>
      </template>
    </PageHeader>

    <div v-if="isLoading" class="loading-state">
      <p>Synchronizing with device registry...</p>
    </div>

    <div v-else-if="fetchError" class="error-state">
      <i class="ri-error-warning-line error-icon"></i>
      <h3 class="error-title">Integration Error</h3>
      <p class="error-desc">{{ fetchError }}</p>
      <BaseButton variant="secondary" size="sm" @click="fetchDevice" class="mt-4">
        Try Re-sync
      </BaseButton>
    </div>

    <div v-else-if="device" class="content-body">
      <div class="detail-grid">
        <DetailCard 
          title="Hardware Parameters" 
          :badge-text="device.status?.toUpperCase()" 
          :badge-type="device.status === 'online' ? 'success' : 'inactive'"
          :rows="deviceInfo"
        />
        
        <DetailCard 
          title="Fleet Assignment" 
          :rows="assignmentInfo"
        >
          <div v-if="isAdmin" class="card-actions">
            <BaseButton variant="primary" size="sm" @click="showAssignModal = true">
              <template #icon><i class="ri-exchange-line"></i></template>
              Reassign Tenant
            </BaseButton>
          </div>
          <div v-else class="card-actions">
            <BaseButton variant="ghost" size="sm" @click="showAssignSiteModal = true">
              <template #icon><i class="ri-exchange-line"></i></template>
              Assign Site
            </BaseButton>
          </div>
        </DetailCard>

      </div>

      <div v-if="showAssignModal" class="modal-overlay" @click="showAssignModal = false">
        <div class="modal-content" @click.stop>
          <h3 class="modal-title">Reassign Tenant</h3>
          <p class="modal-desc">Select the organization that will manage this device.</p>
          
          <div class="form-group">
            <label class="form-label">Organization</label>
            <select :value="selectedTenantId" @change="onTenantChange" class="input">
              <option value="">-- Choose Tenant --</option>
              <option v-for="t in tenants" :key="t.tenant_id" :value="t.tenant_id">
                {{ t.tenant_name }}
              </option>
            </select>
          </div>

          <div class="modal-footer">
            <BaseButton variant="secondary" :disabled="isSaving" @click="showAssignModal = false">Cancel</BaseButton>
            <BaseButton :loading="isSaving" @click="handleAssign">Confirm Assignment</BaseButton>
          </div>
        </div>
      </div>

      <div v-if="showAssignSiteModal" class="modal-overlay" @click="showAssignSiteModal = false">
        <div class="modal-content" @click.stop>
          <h3 class="modal-title">Assign Site</h3>
          <p class="modal-desc">Select the site where this device will be deployed.</p>
          
          <div class="form-group">
            <select v-model="selectedSiteId" class="input">
              <option value="" disabled>Select site...</option>
              <option v-for="site in sites" :key="site.site_id" :value="site.site_id">{{ site.site_name }}</option>
            </select>
          </div>

          <div class="modal-footer">
            <BaseButton variant="secondary" :disabled="isSaving" @click="showAssignSiteModal = false">Cancel</BaseButton>
            <BaseButton :loading="isSaving" @click="handleAssignSite">Confirm Assignment</BaseButton>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.card-actions {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--color-border);
}

.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.modal-content {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 32px;
  width: 100%;
  max-width: 440px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.modal-title { font-size: 20px; font-weight: 700; margin-bottom: 8px; color: var(--color-text-primary); }
.modal-desc { font-size: 14px; color: var(--color-text-secondary); margin-bottom: 24px; }
.form-group { margin-bottom: 24px; }
.modal-footer { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }


.input {
  width: 100%;
  padding: 12px 14px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-primary);
  font-size: 14px;
  outline: none;
}

.loading-state, .error-state {
  padding: 80px;
  text-align: center;
  color: var(--color-text-secondary);
}

.error-icon {
  font-size: 48px;
  color: var(--color-danger);
  margin-bottom: 16px;
  opacity: 0.8;
}

.error-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--color-text-primary);
}

.error-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  max-width: 400px;
  margin: 0 auto 24px auto;
}

.mt-4 { margin-top: 16px; }
</style>
