import { reactive, computed } from 'vue';

const state = reactive({
    activeRequests: 0,
    isDelayed: false
});

let delayTimer = null;

export const loadingStore = {
    get isLoading() {
        return state.activeRequests > 0;
    },

    start() {
        state.activeRequests++;
        if (!state.isDelayed) {
            if (delayTimer) clearTimeout(delayTimer);
            // Small delay before showing the loading bar to prevent flickers on fast requests
            delayTimer = setTimeout(() => {
                if (state.activeRequests > 0) {
                    state.isDelayed = true;
                }
            }, 200);
        }
    },

    finish() {
        state.activeRequests = Math.max(0, state.activeRequests - 1);
        if (state.activeRequests === 0) {
            if (delayTimer) clearTimeout(delayTimer);
            state.isDelayed = false;
        }
    },

    get isDelayed() {
        return state.isDelayed;
    }
};
