import { http } from './http';

/**
 * API module for Bridge-related operations.
 */
export const bridgeApi = {
    /**
     * Trigger manual user synchronization
     */
    syncUsers: (userID) => http.post('/bridge/sync-users', null, {
        headers: { 'X-User-ID': userID }
    }),
};
