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
    }
};
