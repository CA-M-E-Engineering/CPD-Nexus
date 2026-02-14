<script setup>
import BaseSkeleton from './BaseSkeleton.vue';
import TableLoader from './TableLoader.vue';


const props = defineProps({
  columns: { type: Array, required: true },
  data: { type: Array, required: true },
  loading: { type: Boolean, default: false },
  rowClickable: { type: Boolean, default: false }
});

const emit = defineEmits(['row-click']);

const handleRowClick = (row) => {
  if (props.rowClickable) {
    emit('row-click', row);
  }
};
</script>


<template>
  <div class="table-container" :class="{ 'is-loading': loading }">
    <TableLoader v-if="loading" />
    <table class="table">
      <thead>
        <tr>
          <th 
            v-for="col in columns" 
            :key="col.key"
          >
            {{ col.label }}
          </th>


        </tr>
      </thead>
      <tbody class="table-body">
        <template v-if="loading">
          <tr v-for="i in 5" :key="i" class="skeleton-row">
            <td v-for="col in columns" :key="col.key">
              <BaseSkeleton height="20px" width="80%" />
            </td>

          </tr>
        </template>
        <template v-else-if="!data || data.length === 0">
          <tr class="empty-row">
            <td :colspan="columns.length" class="text-center py-16">
              <div class="empty-state-full">
                <div class="empty-icon">
                  <i class="ri-search-line"></i>
                </div>
                <h3 class="empty-title">No records found</h3>
                <p class="empty-description">Try adjusting your filters or adding a new entry.</p>
              </div>
            </td>
          </tr>
        </template>
        <template v-else>
          <tr 
            v-for="(row, index) in data" 
            :key="index" 
            class="data-row"
            :class="{ 'clickable': rowClickable }"
            @click="handleRowClick(row)"
          >
            <td 
              v-for="col in columns" 
              :key="col.key"
              :class="{ 'font-bold': col.bold, 'text-muted': col.muted }"
            >
              <slot :name="`cell-${col.key}`" :item="row" :row="row" :value="row[col.key]">
                {{ row[col.key] }}
              </slot>
            </td>


          </tr>
        </template>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.table-container { 
  background: var(--color-surface); 
  border: 1px solid var(--color-border); 
  border-radius: var(--radius-md); 
  overflow: hidden; 
  width: 100%; 
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

.table { 
  width: 100%; 
  border-collapse: separate; 
  border-spacing: 0; 
}

.table thead { 
  background: var(--color-bg-subtle); 
}

.table th { 
  padding: 12px 16px; 
  text-align: left; 
  font-size: 11px; 
  font-weight: 700; 
  text-transform: uppercase; 
  letter-spacing: 0.05em; 
  color: var(--color-text-muted); 
  border-bottom: 1px solid var(--color-border);
  border-right: 1px solid var(--color-border-light);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.table th:last-child {
  border-right: none;
}


.table td { 
  padding: 14px 16px; 
  border-bottom: 1px solid var(--color-border-light); 
  border-right: 1px solid var(--color-border-light);
  font-size: 14px; 
  color: var(--color-text-primary); 
  vertical-align: middle;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.table td:last-child {
  border-right: none;
}



.table tbody tr:last-child td { 
  border-bottom: none; 
}

.table tbody tr { 
  transition: background-color 0.2s ease; 
}

.table tbody tr:hover { 
  background-color: var(--color-bg-hover); 
}

.data-row.clickable {
  cursor: pointer;
}

.data-row.clickable:active {
  background-color: var(--color-surface-hover);
}

.font-bold {
  font-weight: 600;
  color: var(--color-text-strong);
}

.text-muted {
  color: var(--color-text-muted);
  font-weight: 400;
}


.table-container.is-loading { position: relative; }
.table-body { position: relative; }
.data-row { animation: fadeIn 0.3s ease-out; }

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(2px); }
  to { opacity: 1; transform: translateY(0); }
}

.empty-state-full {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: var(--color-text-secondary);
}
.empty-icon {
  width: 48px;
  height: 48px;
  background: var(--color-bg-subtle);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: var(--color-text-muted);
  margin-bottom: 8px;
}
.empty-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}
.empty-description {
  margin: 0;
  font-size: 14px;
  color: var(--color-text-muted);
  max-width: 300px;
}
.text-center { text-align: center; }
.py-16 { padding-top: 64px; padding-bottom: 64px; }
</style>



