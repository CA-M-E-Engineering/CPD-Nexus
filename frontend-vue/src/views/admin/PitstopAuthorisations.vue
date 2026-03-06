<script setup>
import { ref, onMounted } from 'vue';
import { pitstopApi } from '../../api/pitstop.api';
import DataTable from '../../components/ui/DataTable.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import PageHeader from '../../components/ui/PageHeader.vue';
import { notification } from '../../services/notification';

const isLoading = ref(false);
const isSyncing = ref(false);
const isSubmitting = ref({});
const authorisations = ref([]);
const projects = ref([]);

const authColumns = [
    { key: 'dataset_name', label: 'Dataset Name', sortable: true },
    { key: 'regulator_name', label: 'Regulator', sortable: true },
    { key: 'on_behalf_of_name', label: 'On Behalf Of', sortable: true },
    { key: 'status', label: 'Status' },
    { key: 'last_synced_at', label: 'Last Synced' }
];

const projectColumns = [
    { key: 'reference', label: 'Ref Num', sortable: true },
    { key: 'title', label: 'Title', sortable: true },
    { key: 'status', label: 'Status' },
    { key: 'actions', label: 'Actions', align: 'right' }
];

const fetchAuthorisations = async () => {
    isLoading.value = true;
    try {
        authorisations.value = await pitstopApi.getAuthorisations() || [];
    } catch (error) {
        console.error('Failed to load pitstop authorisations:', error);
        notification.error('Failed to load pitstop configuration.');
    } finally {
        isLoading.value = false;
    }
};

const handleSync = async () => {
    isSyncing.value = true;
    try {
        await pitstopApi.syncAuthorisations();
        notification.success('Pitstop configuration synced successfully');
        await fetchAuthorisations();
    } catch (error) {
        console.error('Failed to sync pitstop config:', error);
        notification.error('Sync failed. Please check network or API keys.');
    } finally {
        isSyncing.value = false;
    }
};

const fetchProjects = async () => {
    try {
        const response = await pitstopApi.getTestingProjects();
        // Backend returns { "data": [...] }
        projects.value = response?.data || [];
    } catch (error) {
        console.error('Failed to load testing projects:', error);
        notification.error('Failed to load projects for testing.');
    }
};

const handleTestSubmission = async (projectId) => {
    isSubmitting.value[projectId] = true;
    try {
        const result = await pitstopApi.testSubmission(projectId);
        // Backend returns { "metrics": { "payloads_submitted": N, "validation_failed": M } }
        const submitted = result?.metrics?.payloads_submitted || 0;
        const failed = result?.metrics?.validation_failed || 0;

        if (failed > 0) {
            notification.warning(`Test complete. Submitted ${submitted}, but ${failed} failed validation.`);
        } else {
            notification.success(`Test complete. Submitted ${submitted} payloads.`);
        }
        // Refresh the list to reflect updated statuses in DB
        await fetchProjects();
    } catch (error) {
        console.error('Test submission failed:', error);
        notification.error('Test submission failed. Check console for details.');
    } finally {
        isSubmitting.value[projectId] = false;
    }
};

onMounted(() => {
    fetchAuthorisations();
    fetchProjects();
});
</script>

<template>
    <div class="page-container">
        <PageHeader 
            title="Pitstop Authorisations" 
            description="Manage configurations and data routing IDs fetched from Pitstop mapping APIs."
        >
            <template #toolbar-right>
                <BaseButton 
                    variant="primary" 
                    icon="ri-refresh-line"
                    :loading="isSyncing"
                    @click="handleSync"
                >
                    {{ isSyncing ? 'Syncing...' : 'Sync Configuration' }}
                </BaseButton>
            </template>
        </PageHeader>

        <div class="content-section">
            <DataTable 
                :columns="authColumns"
                :data="authorisations"
                :loading="isLoading"
                empty-message="No configurations found. Run a sync to fetch data."
            >
                <template #cell-status="{ value }">
                    <span class="status-badge" :class="value?.toLowerCase() || 'inactive'">
                        {{ value || 'UNKNOWN' }}
                    </span>
                </template>
                <template #cell-last_synced_at="{ value }">
                    <span v-if="value">{{ new Date(value).toLocaleString() }}</span>
                    <span v-else class="text-muted">Never synced</span>
                </template>
            </DataTable>
        </div>

        <div class="section-title">
            <h2 style="margin-top: 2rem; margin-bottom: 1rem;">CPD Submission Testing</h2>
            <p class="text-muted" style="margin-bottom: 1rem;">Manually trigger an external API push containing the latest un-synced attendance for a specific project.</p>
        </div>

        <div class="content-section">
            <DataTable 
                :columns="projectColumns"
                :data="projects"
                :loading="isLoading"
                empty-message="No active projects available for testing."
            >
                <template #cell-status="{ value }">
                    <span class="status-badge" :class="value?.toLowerCase() || 'inactive'">
                        {{ value || 'UNKNOWN' }}
                    </span>
                </template>
                <template #cell-actions="{ row }">
                    <BaseButton 
                        variant="secondary" 
                        size="sm"
                        icon="ri-send-plane-fill"
                        :loading="isSubmitting[row.project_id]"
                        @click="handleTestSubmission(row.project_id)"
                    >
                        {{ isSubmitting[row.project_id] ? 'Submitting...' : 'Submit Test' }}
                    </BaseButton>
                </template>
            </DataTable>
        </div>
    </div>
</template>

<style scoped>
.page-container {
    padding: 24px;
    max-width: 1400px;
    margin: 0 auto;
}

.content-section {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    overflow: hidden;
}

.status-badge {
    display: inline-flex;
    align-items: center;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
}
.status-badge.active { background: rgba(16, 185, 129, 0.1); color: #10b981; }
.status-badge.inactive { background: rgba(107, 114, 128, 0.1); color: #6b7280; }
.text-muted { color: var(--color-text-secondary); }
</style>
