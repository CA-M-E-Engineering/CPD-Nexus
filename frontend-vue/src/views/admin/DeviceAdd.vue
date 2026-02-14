<script setup>
import { ref, onMounted, computed } from 'vue';
import { api } from '../../services/api.js';
import { useDeviceStore } from '../../features/devices/store/deviceStore';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import BaseInput from '../../components/ui/BaseInput.vue';
import BaseButton from '../../components/ui/BaseButton.vue';

const props = defineProps({
  id: [Number, String],
  mode: { type: String, default: 'add' } // 'add' or 'edit'
});

const emit = defineEmits(['navigate']);
const deviceStore = useDeviceStore();

const isSaving = ref(false);
  const isLoading = ref(false);
  const Users = ref([]);

  const formData = ref({
    sn: '',
    model: '',
    status: 'offline',
    user_id: ''
  });

  const isEdit = computed(() => props.mode === 'edit');

  const fetchUsers = async () => {
    try {
      // Fetch all Users for assignment
      Users.value = await api.getUsers();
      if (Users.value.length > 0 && !formData.value.user_id) {
          // Default to vendor if available, or first user
          const vendor = Users.value.find(u => u.user_type === 'vendor');
          formData.value.user_id = vendor ? vendor.user_id : Users.value[0].user_id;
      }
    } catch (err) {
      console.error('Failed to fetch Users', err);
    }
  };

  const fetchDevice = async () => {
    if (!isEdit.value || !props.id) return;
    isLoading.value = true;
    try {
      const data = await deviceStore.getDeviceById(props.id);
      if (data) {
        formData.value = { 
          sn: data.sn,
          model: data.model,
          status: data.status,
          user_id: data.user_id || ''
        };
      }
    } catch (err) {
      console.error('Failed to fetch device', err);
    } finally {
      isLoading.value = false;
    }
  };

  onMounted(async () => {
    await fetchUsers();
    await fetchDevice();
  });

const handleSubmit = async () => {
  isSaving.value = true;
  try {
    if (isEdit.value) {
      await deviceStore.updateDevice(props.id, formData.value);
      notification.success('Device parameters updated');
    } else {
      await deviceStore.registerDevice(formData.value);
      notification.success('New device provisioned successfully');
    }
    emit('navigate', 'devices');
  } catch (err) {
    console.error('Failed to save device', err);
    notification.error(err.message || 'Failed to sync with device registry');
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="device-add">
    <PageHeader 
      :title="isEdit ? 'Edit Device Parameters' : 'Provision New Unit'" 
      :description="isEdit ? 'Update hardware serial and system status' : 'Register a new unit into the global device registry'"
    />

    <div v-if="isLoading" class="loading-state">
      <p>Fetching device parameters...</p>
    </div>

    <form v-else class="form-container" @submit.prevent="handleSubmit">
      <div class="form-grid">
        <BaseInput 
          v-model="formData.sn" 
          label="Serial Number (SN)" 
          placeholder="e.g., DS-K1T341AMF-001" 
          required 
        />
        <BaseInput 
          v-model="formData.model" 
          label="Model / Type" 
          placeholder="e.g., DS-K1T341AMF" 
          required 
        />
        
        <div class="form-group full-width">
          <label class="form-label">Assigned User</label>
          <select v-model="formData.user_id" class="form-select" required>
            <option v-if="Users.length === 0" value="" disabled>No users available</option>
            <option v-for="User in Users" :key="User.user_id" :value="User.user_id">
              {{ User.user_name }} ({{ User.user_type }})
            </option>
          </select>
        </div>


      </div>

      <div class="form-actions">
        <BaseButton variant="secondary" @click="$emit('navigate', 'devices')">Cancel</BaseButton>
        <BaseButton :loading="isSaving" type="submit">
          {{ isEdit ? 'Save Changes' : 'Register Device' }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<style scoped>
.form-container {
  max-width: 600px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 32px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
  margin-bottom: 32px;
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
    transition: border-color 0.2s;
}

.form-select:focus {
    border-color: var(--color-primary);
}

.loading-state {
  padding: 48px;
  text-align: center;
  color: var(--color-text-secondary);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
}

.full-width {
  grid-column: 1 / -1;
}
</style>
