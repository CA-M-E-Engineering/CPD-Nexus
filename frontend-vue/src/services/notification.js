import { ref } from 'vue';

const activeNotification = ref(null);

export const notification = {
    state: activeNotification,

    success(message, duration = 3000, isLarge = false) {
        this.notify(message, 'success', duration, isLarge);
    },

    error(message, duration = 4000, isLarge = false) {
        this.notify(message, 'error', duration, isLarge);
    },

    notify(message, type = 'success', duration = 3000, isLarge = false) {
        activeNotification.value = {
            id: Date.now(),
            message,
            type,
            isLarge
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
