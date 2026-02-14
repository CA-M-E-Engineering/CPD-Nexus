import { loadingStore } from '../store/loading';

/**
 * Utility for making HTTP requests with standardized error handling and configuration.
 */

const BASE_URL = '/api';

/**
 * Custom error class for API failures
 */
class ApiError extends Error {
    constructor(message, status, data) {
        super(message);
        this.name = 'ApiError';
        this.status = status;
        this.data = data;
    }
}

/**
 * Shared request helper
 */
async function request(endpoint, options = {}) {
    let url = `${BASE_URL}${endpoint}`;

    if (options.params) {
        const queryParams = new URLSearchParams();
        Object.entries(options.params).forEach(([key, value]) => {
            if (value !== undefined && value !== null) {
                queryParams.append(key, value);
            }
        });
        const queryString = queryParams.toString();
        if (queryString) {
            url += (url.includes('?') ? '&' : '?') + queryString;
        }
    }
    const defaultHeaders = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    };

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 30000); // 30s timeout

    const config = {
        ...options,
        signal: controller.signal,
        headers: {
            ...defaultHeaders,
            ...options.headers,
        },
    };

    const token = localStorage.getItem('auth_token');
    if (token) {
        config.headers['Authorization'] = `Bearer ${token}`;
    }

    if (config.body && typeof config.body === 'object') {
        config.body = JSON.stringify(config.body);
    }

    loadingStore.start();
    try {
        const response = await fetch(url, config);
        clearTimeout(timeoutId);
        console.log(`[HTTP Debug] ${config.method || 'GET'} ${url} -> Status ${response.status}`);
        const contentType = response.headers.get('content-type');
        let data = null;

        if (contentType && contentType.includes('application/json')) {
            data = await response.json().catch(() => null);
        } else {
            data = await response.text().catch(() => null);
        }
        console.log(`[HTTP Debug] Data:`, data);

        if (!response.ok) {
            const errorMessage = (typeof data === 'object' ? data?.message : data) || `Request failed with status ${response.status}`;
            throw new ApiError(errorMessage, response.status, data);
        }

        return data;
    } catch (error) {
        if (error instanceof ApiError) throw error;
        throw new ApiError(error.message || 'Network error occurred', 500);
    } finally {
        loadingStore.finish();
    }
}

export const http = {
    get: (url, options) => request(url, { ...options, method: 'GET' }),
    post: (url, body, options) => request(url, { ...options, method: 'POST', body }),
    put: (url, body, options) => request(url, { ...options, method: 'PUT', body }),
    patch: (url, body, options) => request(url, { ...options, method: 'PATCH', body }),
    delete: (url, options) => request(url, { ...options, method: 'DELETE' }),
};
