import { ref } from 'vue';

const activeNotification = ref(null);

export const notification = {
    state: activeNotification,

    success(message, duration = 3000) {
        this.notify(message, 'success', duration);
    },

    error(message, duration = 4000) {
        this.notify(message, 'error', duration);
    },

    notify(message, type = 'success', duration = 3000) {
        activeNotification.value = {
            id: Date.now(),
            message,
            type
        };

        if (duration > 0) {
            setTimeout(() => {
                if (activeNotification.value?.id === activeNotification.value?.id) {
                    this.clear();
                }
            }, duration);
        }
    },

    clear() {
        activeNotification.value = null;
    }
};
