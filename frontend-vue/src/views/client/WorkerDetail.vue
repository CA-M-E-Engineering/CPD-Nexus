<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import DetailCard from '../../components/ui/DetailCard.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import BaseInput from '../../components/ui/BaseInput.vue';

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

const passTypes = [
    { value: 'SP', label: 'Singapore Pink IC (SP)' },
    { value: 'SB', label: 'Singapore Blue IC (SB)' },
    { value: 'EP', label: 'Employment Pass (EP)' },
    { value: 'SPASS', label: 'S Pass (SPASS)' },
    { value: 'WP', label: 'Work Permit (WP)' },
    { value: 'ENTREPASS', label: 'EntrePass' },
    { value: 'LTVP', label: 'Long-Term Visit Pass (LTVP)' }
];

const bcaTrades = [
    { value: '1.1', label: '1.1 - Site Management (Ancillary)' },
    { value: '1.2', label: '1.2 - Site Support (Ancillary)' },
    { value: '1.3', label: '1.3 - General Machine Operation' },
    { value: '1.4', label: '1.4 - Site Preparation' },
    { value: '1.5', label: '1.5 - Scaffolding' },
    { value: '2.1', label: '2.1 - Demolition (Civil/Structural)' },
    { value: '2.2', label: '2.2 - Earthworks' },
    { value: '2.3', label: '2.3 - Foundation' },
    { value: '2.4', label: '2.4 - Tunnelling' },
    { value: '2.5', label: '2.5 - Reinforced Concrete' },
    { value: '2.6', label: '2.6 - Structural Steel' },
    { value: '3.1', label: '3.1 - Ceiling (Architectural)' },
    { value: '3.2', label: '3.2 - Partition Wall' },
    { value: '4.1', label: '4.1 - Plumbing, Sanitary & Gas' },
    { value: '4.3', label: '4.3 - Electrical' }
];

const getPassLabel = (val) => passTypes.find(p => p.value === val)?.label || val || '---';
const getTradeLabel = (val) => bcaTrades.find(t => t.value === val)?.label || val || '---';

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
  { label: 'Pass Type', value: getPassLabel(worker.value?.person_id_and_work_pass_type) },
  { label: 'Designated Trade', value: getTradeLabel(worker.value?.person_trade) }
]);

const assignmentInfo = computed(() => [
  { label: 'Project', value: worker.value?.project_name || '---' },
  { label: 'Site Name', value: worker.value?.site_name || '---' },
  { label: 'Site Location', value: worker.value?.site_location || '---' }
]);

const authenticationInfo = computed(() => [
  { label: 'Face Recognition', value: 'Enrolled' },
  { label: 'Access Card (NFC)', value: 'Pending Provision' },
  { label: 'Fingerprint ID', value: 'Not Supported' }
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
      <DetailCard 
        title="IoT Authentication" 
        :badge-text="'Active'"
        :badge-type="'success'"
        :rows="authenticationInfo"
      />
    </div>
  </div>
</template>

<style scoped>
.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}
.loading-state {
  padding: 48px;
  text-align: center;
  color: var(--color-text-secondary);
}

@media (max-width: 1024px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>

