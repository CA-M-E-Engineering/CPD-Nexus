<template>
  <MaintenanceLayout
    title="Device Registry"
    description="Manage system-wide IoT hardware and provision new units."
    action-label="Add New Device"
    @back="$emit('navigate', 'devices')"
    @action="showAddModal = true"
  >
    <template #list>
      <div v-if="isLoading" class="loading-padding">
        <p>Loading registry...</p>
      </div>
      <MaintenanceListItem
        v-for="device in devices"
        :key="device.id"
        :label="device.deviceId"
        :secondary="`${device.type} • ${device.site || 'Unassigned'}`"
        icon="ri-cpu-line"
        :status="device.status === DEVICE_STATUS.ONLINE ? 'Active' : 'Offline'"
      />
    </template>
  </MaintenanceLayout>

  <BaseModal
    :show="showAddModal"
    title="Register New Device"
    description="Enter the hardware details to provision a new unit."
    confirm-label="Register Device"
    :loading="isSaving"
    @close="showAddModal = false"
    @confirm="handleAdd"
  >
    <div class="form-body">
      <div class="form-group">
        <BaseInput v-model="formData.deviceId" label="Device ID" placeholder="e.g. GTW-00900" required />
      </div>
      <div class="form-group">
        <BaseInput v-model="formData.type" label="Model/Type" placeholder="e.g. Gateway Pro" required />
      </div>
    </div>
  </BaseModal>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import MaintenanceLayout from '../../components/ui/MaintenanceLayout.vue';
import MaintenanceListItem from '../../components/ui/MaintenanceListItem.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseModal from '../../components/ui/BaseModal.vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import { DEVICE_STATUS } from '../../utils/constants.js';

const devices = ref([]);
const showAddModal = ref(false);
const isLoading = ref(true);
const isSaving = ref(false);

const formData = ref({
  deviceId: '',
  type: ''
});

const fetchDevices = async () => {
  isLoading.value = true;
  try {
    devices.value = await api.getDevices();
  } catch (err) {
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

const handleAdd = async () => {
  if (!formData.value.deviceId || !formData.value.type) return;
  isSaving.value = true;
  try {
    // Map form fields to API expected fields if necessary
    // Backend expects 'sn' and 'type' (mapped to model in DB)
    const payload = {
      sn: formData.value.deviceId, // Using deviceId as SN for this form
      type: formData.value.type
    };
    
    await api.createDevice(payload);
    await fetchDevices(); // Reload the list from server
    notification.success('New device registered successfully');
    showAddModal.value = false;
    formData.value = { deviceId: '', type: '' };
  } catch (err) {
    console.error('Failed to register device', err);
    notification.error(err.message || 'Failed to provision hardware unit');
  } finally {
    isSaving.value = false;
  }
};

onMounted(fetchDevices);

defineEmits(['navigate']);
</script>

<style scoped>
.form-body {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.loading-padding {
  padding: 32px;
  text-align: center;
  color: var(--color-text-secondary);
}
</style>

