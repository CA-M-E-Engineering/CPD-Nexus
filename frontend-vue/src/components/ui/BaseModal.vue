<script setup>
import { onMounted, onUnmounted } from 'vue';
import BaseButton from './BaseButton.vue';

const props = defineProps({
    show: { type: Boolean, default: false },
    title: { type: String, default: '' },
    description: { type: String, default: '' },
    maxWidth: { type: String, default: '480px' },
    showFooter: { type: Boolean, default: true },
    confirmLabel: { type: String, default: 'Confirm' },
    cancelLabel: { type: String, default: 'Cancel' },
    loading: { type: Boolean, default: false },
    confirmDisabled: { type: Boolean, default: false }
});

const emit = defineEmits(['close', 'confirm']);

const handleEscape = (e) => {
    if (e.key === 'Escape' && props.show) {
        emit('close');
    }
};

onMounted(() => window.addEventListener('keydown', handleEscape));
onUnmounted(() => window.removeEventListener('keydown', handleEscape));
</script>

<template>
    <Transition name="fade">
        <div v-if="show" class="modal-overlay" @click="emit('close')">
            <div 
                class="modal-content" 
                :style="{ maxWidth }" 
                @click.stop
            >
                <div class="modal-header">
                    <h3 class="modal-title">{{ title }}</h3>
                    <p v-if="description" class="modal-desc">{{ description }}</p>
                </div>
                
                <div class="modal-body">
                    <slot></slot>
                </div>

                <div v-if="showFooter" class="modal-footer">
                    <slot name="footer">
                        <BaseButton variant="secondary" @click="emit('close')">
                            {{ cancelLabel }}
                        </BaseButton>
                        <BaseButton 
                            :loading="loading" 
                            :disabled="confirmDisabled"
                            @click="emit('confirm')"
                        >
                            {{ confirmLabel }}
                        </BaseButton>
                    </slot>
                </div>
            </div>
        </div>
    </Transition>
</template>

<style scoped>
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.82);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
}

.modal-content {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-lg);
    padding: 32px;
    width: 90%;
    box-shadow: var(--shadow-xl);
    animation: modalIn 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes modalIn {
    from { opacity: 0; transform: scale(0.95) translateY(10px); }
    to { opacity: 1; transform: scale(1) translateY(0); }
}

.modal-header { margin-bottom: 24px; }
.modal-title { font-size: 20px; font-weight: 600; margin-bottom: 8px; color: var(--color-text-primary); }
.modal-desc { font-size: 14px; color: var(--color-text-secondary); line-height: 1.5; }

.modal-body {
    margin-bottom: 32px;
}

.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding-top: 24px;
    border-top: 1px solid var(--color-border);
}

/* Transitions */
.fade-enter-active, .fade-leave-active { transition: opacity 0.25s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
