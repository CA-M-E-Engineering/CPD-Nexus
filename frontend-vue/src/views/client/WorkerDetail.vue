<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { api } from '../../services/api.js';
import PageHeader from '../../components/ui/PageHeader.vue';
import DetailCard from '../../components/ui/DetailCard.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';

const props = defineProps({
  id: [Number, String]
});

const emit = defineEmits(['navigate']);

const worker = ref(null);
const isLoading = ref(true);

const fetchWorker = async () => {
  if (!props.id) return;
  isLoading.value = true;
  try {
    worker.value = await api.getWorkerById(props.id);
  } catch (err) {
    console.error('[WorkerDetail] Failed to fetch worker:', err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchWorker);

// Re-fetch if ID changes (wiring fix)
watch(() => props.id, fetchWorker);

const workerInfo = computed(() => [
  { label: 'Full Name', value: worker.value?.name || '---' },
  { label: 'FIN/NRIC', value: worker.value?.fin || '---' },
  { 
    label: 'Role', 
    value: worker.value?.role === 'pic' ? 'PIC (Person In Charge)' : 
           worker.value?.role === 'manager' ? 'Manager' : 
           worker.value?.role === 'worker' ? 'Worker' : 
           (worker.value?.role || '---') 
  },
  { label: 'Employer', value: worker.value?.company_name || '---' },
  { label: 'Trade Code', value: worker.value?.trade_code || '---' },
  { label: 'Email', value: worker.value?.email || '---' }
]);

const assignmentInfo = computed(() => [
  { label: 'Project Name', value: worker.value?.project_name || '---' },
  { label: 'Project ID', value: worker.value?.current_project_id || '---' },
  { label: 'Site Name', value: worker.value?.site_name || '---' },
  { label: 'Site Location', value: worker.value?.site_location || '---' },
  { label: 'Tenant Name', value: worker.value?.tenant_name || '---' },
  { label: 'Tenant Address', value: worker.value?.tenant_address || '---' }
]);

const handleEdit = () => {
  emit('navigate', 'worker-add', { id: props.id, mode: 'edit' });
};

const handleChangeProject = () => {
  emit('navigate', 'worker-assign-project', { id: props.id });
};
</script>

<template>
  <div class="worker-detail">
    <PageHeader 
      :title="worker?.name || 'Loading...'" 
      description="View worker details and activity history"
      variant="detail"
    >
      <template #toolbar-left>
        <BaseButton variant="ghost" size="sm" @click="$emit('navigate', 'workers')">
          <template #icon><i class="ri-arrow-left-line"></i></template>
          Back to Workers
        </BaseButton>
      </template>
      <template #toolbar-right>
        <BaseButton variant="secondary" size="sm" @click="handleEdit">Edit Details</BaseButton>
        <BaseButton variant="secondary" size="sm" @click="handleChangeProject">Change Project</BaseButton>
      </template>
    </PageHeader>

    <div v-if="isLoading" class="loading-state">
      <p>Loading worker profile...</p>
    </div>

    <div v-else-if="worker" class="detail-grid">
      <DetailCard 
        title="Worker Information" 
        :badge-text="worker.status" 
        :badge-type="worker.status === 'active' ? 'success' : 'inactive'"
        :rows="workerInfo"
      />
      <DetailCard 
        title="Assignment" 
        :rows="assignmentInfo"
      />
    </div>
  </div>
</template>

<style scoped>
.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 24px;
}
.loading-state {
  padding: 48px;
  text-align: center;
  color: var(--color-text-secondary);
}
</style>

