import { http } from './http';

export const settingsApi = {
    getSettings: () => http.get('/settings'),
    updateSettings: (data) => http.put('/settings', data)
};
