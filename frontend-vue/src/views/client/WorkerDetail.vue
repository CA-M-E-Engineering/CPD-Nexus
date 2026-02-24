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
    const response = await api.getWorkerById(props.id);
    worker.value = typeof response === 'string' ? JSON.parse(response) : response;
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
  { label: 'NRIC / FIN', value: worker.value?.person_id_no || '---' },
  { label: 'Nationality', value: worker.value?.person_nationality || '---' },
  { 
    label: 'System Role', 
    value: worker.value?.role === 'pic' ? 'PIC (Person In Charge)' : 
           worker.value?.role === 'manager' ? 'Manager' : 
           worker.value?.role === 'worker' ? 'Worker' : 
           (worker.value?.role || '---') 
  },
  { label: 'Email', value: worker.value?.email || '---' }
]);

const complianceInfo = computed(() => [
  { label: 'Pass Type', value: worker.value?.person_id_and_work_pass_type || '---' },
  { label: 'Designated Trade', value: worker.value?.person_trade || '---' }
]);

const assignmentInfo = computed(() => [
  { label: 'Project', value: worker.value?.project_name || '---' },
  { label: 'Site Name', value: worker.value?.site_name || '---' },
  { label: 'Site Location', value: worker.value?.site_location || '---' }
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
        title="Compliance Info" 
        :rows="complianceInfo"
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

