import { http } from './http';

export const pitstopApi = {
    /**
     * Fetch existing stored pitstop configuration routes
     */
    getAuthorisations() {
        return http.get('/pitstop/authorisations');
    },

    /**
     * Trigger a sync to Pitstop server and pull down latest config
     */
    syncAuthorisations() {
        return http.post('/pitstop/authorisations/sync');
    },

    /**
     * Assign specific pitstop "on behalf of" entities to a user account
     * @param {string} userId 
     * @param {string[]} onBehalfOfNames 
     */
    assignOnBehalfOfs(userId, onBehalfOfNames) {
        return http.post(`/users/${userId}/pitstop-on-behalf-of`, { on_behalf_of_names: onBehalfOfNames });
    }
};
