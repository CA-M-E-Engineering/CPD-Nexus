<script setup>
import { ref, onMounted, computed } from 'vue';
import { api } from '../../services/api.js';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import DetailCard from '../../components/ui/DetailCard.vue';
import BaseTabs from '../../components/ui/BaseTabs.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import DataTable from '../../components/ui/DataTable.vue';
import BaseBadge from '../../components/ui/BaseBadge.vue';
import ConfirmDialog from '../../components/ui/ConfirmDialog.vue';
import { TENANT_STATUS } from '../../utils/constants.js';

const props = defineProps({
  id: [Number, String]
});

const emit = defineEmits(['navigate']);

const activeTab = ref('Workers');
const project = ref(null);
const isLoading = ref(true);

const tabs = [
  { id: 'Workers', label: 'Assigned Workers' },
  { id: 'Devices', label: 'IoT Devices' }
];

const fetchProject = async () => {
  isLoading.value = true;
  try {
    const response = await api.getProjectById(props.id);
    // Robust parsing in case Content-Type fix didn't catch everything
    project.value = typeof response === 'string' ? JSON.parse(response) : response;
  } catch (err) {
    console.error('Failed to fetch project:', err);
  } finally {
    isLoading.value = false;
  }
};

// fetchProject is called within the unified onMounted hook below

const projectInfo = computed(() => [
  { label: 'Project Ref', value: project.value?.reference || '---' },
  { label: 'Project Title', value: project.value?.title || '---' },
  { label: 'Site', value: project.value?.site_name || '---' },
  { label: 'Contract Number', value: project.value?.contract || '---' },
  { label: 'Contract Name', value: project.value?.contract_name || '---' },
  { label: 'Location', value: project.value?.location || '---' }
]);

const contractorInfo = computed(() => {
  const p = project.value;
  const fmt = (name, uen) => {
    if (!name && !uen) return '---';
    if (uen) return `${name || '---'} (${uen})`;
    return name;
  };
  return [
    { label: 'Main Contractor', value: fmt(p?.main_contractor_name, p?.main_contractor_uen) },
    { label: 'Offsite Fabricator', value: fmt(p?.offsite_fabricator_name, p?.offsite_fabricator_uen) },
    { label: 'Worker Company', value: fmt(p?.worker_company_name, p?.worker_company_uen) },
    { label: 'Worker Company Client', value: fmt(p?.worker_company_client_name, p?.worker_company_client_uen) },
    { label: 'HDB Precinct', value: p?.hdb_precinct || '---' }
  ];
});

const resourceStats = computed(() => [
  { label: 'Total Workers', value: project.value?.worker_count || '0' },
  { label: 'Active Today', value: project.value?.worker_count ? Math.floor(project.value.worker_count * 0.8) : '0' },
  { label: 'Devices', value: project.value?.device_count || '0' }
]);

const workerColumns = [
  { key: 'name', label: 'Worker Name' },
  { key: 'role', label: 'Role' },
  { key: 'trade_code', label: 'Trade Code' },
  { key: 'status', label: 'Status' },
  { key: 'actions', label: 'Actions', width: '100px' }
];

const assignedWorkers = ref([]);
const assignedDevices = ref([]);

const deviceColumns = [
  { key: 'device_id', label: 'Device ID', bold: true },
  { key: 'model', label: 'Type' },
  { key: 'status', label: 'Status' },
  { key: 'battery', label: 'Battery' }
];

const fetchWorkers = async () => {
    try {
        const data = await api.getWorkers();
        // Only show workers assigned to THIS project
        assignedWorkers.value = (data || []).filter(w => w.current_project_id === props.id);
    } catch (err) {
        console.error('Failed to fetch assigned workers', err);
    }
};

onMounted(async () => {
    // 1. Fetch project details
    await fetchProject();
    
    // 2. Resolve tenant context for scoped fetching
    const savedUser = localStorage.getItem('auth_user');
    let tenantId = 'tenant-client-1';
    if (savedUser) {
        try {
            const user = JSON.parse(savedUser);
            tenantId = user.tenant_id || user.id;
        } catch (e) {
            console.error('Failed to parse auth_user', e);
        }
    }

    // 3. Fetch related data in parallel with tenant scoping
    const [workersData, devicesData] = await Promise.all([
      api.getWorkers({ tenant_id: tenantId }),
      api.getDevices({ tenant_id: tenantId })
    ]);

    // 4. Filter workers assigned to THIS project
    // Ensure both IDs are strings for reliable comparison
    const currentProjectId = String(props.id);
    assignedWorkers.value = (workersData || []).filter(w => String(w.current_project_id) === currentProjectId);

    // 5. Filter devices assigned to the site associated with THIS project
    if (project.value && project.value.site_id) {
       const projectSiteId = String(project.value.site_id);
       assignedDevices.value = (devicesData || []).filter(d => {
         return d.site_id && String(d.site_id) === projectSiteId;
       });
    }
});

const handleEdit = () => {
  emit('navigate', 'project-add', { id: props.id, mode: 'edit' });
};

const handleAssignWorker = () => {
  emit('navigate', 'project-assign-workers', { id: props.id });
};

const handleManageWorker = (worker) => {
  emit('navigate', 'worker-detail', { id: worker.user_id });
};

const showDeleteConfirm = ref(false);
const isDeleting = ref(false);

const handleDelete = async () => {
  isDeleting.value = true;
  try {
    await api.deleteProject(props.id);
    notification.success('Project deleted successfully');
    emit('navigate', 'projects');
  } catch (err) {
    console.error('Failed to delete project', err);
    notification.error(err.message || 'Failed to delete project');
  } finally {
    isDeleting.value = false;
    showDeleteConfirm.value = false;
  }
};
</script>

<template>
  <div class="project-detail">
    <PageHeader 
      :title="project?.title || (isLoading ? 'Loading...' : 'Project Detail')" 
      description="Detailed project view and resource management"
      variant="detail"
    >
      <template #toolbar-left>
        <BaseButton variant="ghost" size="sm" @click="$emit('navigate', 'projects')">
          <template #icon><i class="ri-arrow-left-line"></i></template>
          Back to Projects
        </BaseButton>
      </template>
      <template #toolbar-right>
        <BaseButton variant="secondary" size="sm" icon="ri-edit-line" @click="handleEdit">Edit</BaseButton>
        <BaseButton size="sm" icon="ri-add-line" @click="handleAssignWorker">Assign Worker</BaseButton>
      </template>
    </PageHeader>

    <div v-if="isLoading" class="loading-state">
      <p>Loading project details...</p>
    </div>

    <div v-else-if="project" class="content-body">
      <div class="detail-grid">
          <DetailCard 
            title="Project Information" 
            :badge-text="project.status" 
            :badge-type="project.status === 'active' ? 'success' : 'inactive'"
            :rows="projectInfo"
          />
          <DetailCard 
            title="Contractor Details" 
            :rows="contractorInfo"
          />
          <DetailCard 
            title="Resources" 
            :rows="resourceStats"
          />
      </div>

      <BaseTabs v-model="activeTab" :tabs="tabs" />

    <div v-show="activeTab === 'Workers'" class="tab-content">
        <DataTable :columns="workerColumns" :data="assignedWorkers">
          <template #cell-name="{ item }">
              <strong>{{ item.name }}</strong>
          </template>
          <template #cell-status="{ item }">
            <BaseBadge :type="item.status === 'active' ? 'success' : 'inactive'">{{ item.status }}</BaseBadge>
          </template>
          <template #cell-actions="{ item }">
            <BaseButton variant="secondary" size="sm" @click="handleManageWorker(item)">Manage</BaseButton>
          </template>
        </DataTable>
      </div>

    <div v-show="activeTab === 'Devices'" class="tab-content">
          <DataTable :columns="deviceColumns" :data="assignedDevices" row-clickable @row-click="(d) => $emit('navigate', 'device-detail', { id: d.device_id })">
            <template #cell-status="{ item }">
              <BaseBadge :type="item.status === 'online' ? 'success' : 'inactive'">
                {{ item.status.toUpperCase() }}
              </BaseBadge>
            </template>
            <template #cell-battery="{ item }">
              <span>{{ item.battery }}%</span>
            </template>
          </DataTable>
          <div v-if="assignedDevices.length === 0" class="empty-state">
              <p class="empty-description">No devices currently assigned to this project's site.</p>
          </div>
      </div>

      <ConfirmDialog
        :show="showDeleteConfirm"
        title="Delete Project"
        description="Are you sure you want to delete this project? This will significantly impact resource planning and history."
        confirm-label="Delete Project"
        :loading="isDeleting"
        variant="danger"
        @close="showDeleteConfirm = false"
        @confirm="handleDelete"
      />
    </div>
  </div>
</template>

<style scoped>
.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.tab-content {
  margin-top: 24px;
}



.loading-state {
  padding: 64px;
  text-align: center;
  color: var(--color-text-secondary);
}
</style>

