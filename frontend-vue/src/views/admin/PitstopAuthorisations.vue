<script setup>
import { ref, onMounted } from 'vue';
import { pitstopApi } from '../../api/pitstop.api';
import DataTable from '../../components/ui/DataTable.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import PageHeader from '../../components/ui/PageHeader.vue';
import { notification } from '../../services/notification';

const isLoading = ref(false);
const isSyncing = ref(false);
const authorisations = ref([]);

const columns = [
    { key: 'dataset_name', label: 'Dataset Name', sortable: true },
    { key: 'regulator_name', label: 'Regulator', sortable: true },
    { key: 'maincon_name', label: 'Main Contractor', sortable: true },
    { key: 'status', label: 'Status' },
    { key: 'last_synced_at', label: 'Last Synced' },
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

onMounted(() => {
    fetchAuthorisations();
});
</script>

<template>
    <div class="page-container">
        <PageHeader 
            title="Pitstop Authorisations" 
            description="Manage configurations and data routing IDs fetched from Pitstop mapping APIs."
        >
            <template #actions>
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
                :columns="columns"
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
                <template #cell-actions="{ item }">
                    <BaseButton 
                        variant="ghost" 
                        size="sm" 
                        icon="ri-refresh-line"
                        @click.stop="handleSync"
                        title="Sync this configuration"
                    />
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
