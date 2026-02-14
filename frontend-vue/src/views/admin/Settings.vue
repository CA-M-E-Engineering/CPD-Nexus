<script setup>
// Force HMR update
import { ref, onMounted, computed } from 'vue';
import { api } from '../../services/api';
import { notification } from '../../services/notification';
import PageHeader from '../../components/ui/PageHeader.vue';
import DetailCard from '../../components/ui/DetailCard.vue';
import BaseButton from '../../components/ui/BaseButton.vue';
import BaseInput from '../../components/ui/BaseInput.vue';

const isLoading = ref(false);
const isSaving = ref(false);

const settings = ref({
  device_sync_interval: '00:01:00',
  cpd_submission_time: '09:00:00',
  response_size_limit: 1048576
});

const stats = ref({
  total_devices: 0,
  online_devices: 0
});

const sampleJson = computed(() => {
  return JSON.stringify({
    "participants": [
      {
        "id": "S8000001W",
        "name": "Worker Ali",
        "meta": {
          "data_ref_id": "PROJ-JUR-01"
        },
        "on_behalf_of": {
          "id": "S8000001W",
          "name": "Worker Ali"
        }
      }
    ],
    "payload": [
      {
        "submission_entity": 1,
        "submission_month": "2026-02",
        "project_reference_number": "PROJ-JUR-01",
        "project_title": "Mega Jurong Site",
        "project_location_description": "Jurong",
        "main_contractor_company_name": "Mega Engineering",
        "main_contractor_company_unique_entity_number": "MEGA12345X",
        "person_id_no": "S8000001W",
        "person_name": "Worker Ali",
        "person_id_and_work_pass_type": "NRIC",
        "person_trade": "ELEC",
        "person_employer_company_name": "Mega Engineering",
        "person_employer_company_unique_entity_number": "MEGA12345X",
        "person_employer_company_trade": [
          "ELEC"
        ],
        "person_employer_client_company_name": "Mega Engineering",
        "person_employer_client_company_unique_entity_number": "MEGA12345X",
        "person_attendance_date": "2026-02-10",
        "person_attendance_details": [
          {
            "time_in": "2026-02-10T07:34:47Z",
            "time_out": "2026-02-10T10:34:47Z"
          }
        ]
      }
    ],
    "on_behalf_of": [
      {
        "id": "MEGA12345X"
      }
    ]
  }, null, 2);
});

const fetchSettings = async () => {
  isLoading.value = true;
  try {
    const response = await api.getSettings();
    if (response) {
      settings.value = response.settings;
      stats.value = {
        total_devices: response.total_devices,
        deployed_devices: response.deployed_devices
      };
    }
  } catch (err) {
    console.error('Failed to load settings', err);
    notification.error('Failed to load system settings');
  } finally {
    isLoading.value = false;
  }
};

const updateSettings = async (section) => {
  isSaving.value = true;
  try {
    await api.updateSettings(settings.value);
    notification.success(`${section} settings updated successfully`);
  } catch (err) {
    console.error('Failed to save settings', err);
    notification.error('Failed to save settings');
  } finally {
    isSaving.value = false;
  }
};

onMounted(fetchSettings);
</script>

<template>
  <div class="admin-settings">
    <PageHeader 
      title="System Settings" 
      description="Configure global application parameters and monitor infrastructure health"
    />

    <div v-if="isLoading" class="loading-state">
      <p>Loading configuration...</p>
    </div>

    <div v-else class="settings-grid">
      <!-- Panel 1: General Configuration -->
      <div class="settings-section">
        <DetailCard title="General Configuration">
          <div class="setting-item stats-row">
            <div class="stat-box">
              <span class="stat-label">Total Devices</span>
              <span class="stat-value">{{ stats.total_devices }}</span>
            </div>
            <div class="stat-box">
              <span class="stat-label">Deployed</span>
              <span class="stat-value success">{{ stats.deployed_devices }}</span>
            </div>
          </div>

          <div class="setting-item">
            <BaseInput 
              label="Device Sync Interval (HH:MM:SS)" 
              type="text" 
              v-model="settings.device_sync_interval"
              placeholder="00:01:00"
            />
            <p class="help-text">Frequency of device heartbeat synchronization.</p>
          </div>

          <div class="setting-actions">
            <BaseButton :loading="isSaving" @click="updateSettings('General')">Update Configuration</BaseButton>
          </div>
        </DetailCard>
      </div>

      <!-- Panel 2: CPD Customization -->
      <div class="settings-section">
        <DetailCard title="CPD Customization">
          <div class="setting-item">
            <BaseInput 
              label="Response Size Limit (Bytes)" 
              type="number" 
              v-model.number="settings.response_size_limit" 
            />
            <p class="help-text">Maximum allowed payload size for API responses.</p>
          </div>

          <div class="setting-item">
            <BaseInput 
              label="Submission Time (HH:MM:SS)" 
              type="text" 
              v-model="settings.cpd_submission_time" 
              placeholder="09:00:00"
            />
            <p class="help-text">Scheduled daily time for CPD data submission.</p>
          </div>

          <div class="setting-actions">
            <BaseButton :loading="isSaving" @click="updateSettings('CPD')">Update CPD Settings</BaseButton>
          </div>
        </DetailCard>
      </div>

      <!-- Panel 3: Data Format -->
      <div class="settings-section">
        <DetailCard title="CPD JSON Format">
          <div class="code-block-wrapper">
            <pre><code>{{ sampleJson }}</code></pre>
          </div>
          <p class="help-text mt-2">Sample structure of the daily submission payload.</p>
        </DetailCard>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 24px;
  align-items: start;
}

.setting-item {
  margin-bottom: 20px;
}

.setting-actions {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}

.help-text {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-top: 6px;
}

.stats-row {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  background: var(--color-bg);
  padding: 12px;
  border-radius: var(--radius-sm);
}

.stat-box {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--color-text-secondary);
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.stat-value.success {
  color: var(--color-success);
}

.code-block-wrapper {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 16px;
  border-radius: var(--radius-sm);
  font-family: 'Fira Code', monospace;
  font-size: 12px;
  overflow-x: auto;
  border: 1px solid var(--color-border);
}

.code-block-wrapper pre {
  margin: 0;
}

.mt-2 {
  margin-top: 8px;
}

.loading-state {
  padding: 48px;
  text-align: center;
  color: var(--color-text-secondary);
}
</style>
