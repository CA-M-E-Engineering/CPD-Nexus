<script setup>
import { onMounted, ref, onUnmounted, watch, nextTick } from 'vue';
import { api } from '../../services/api.js';
import { MAP_MODES, USER_STATUS } from '../../utils/constants.js';

// Add local constant if not in utils yet, or better, assumes it returns 'single-edit' 
// actually I should check constants.js first but I can extend the validator.

const props = defineProps({
    mode: {
        type: String,
        required: true,
        validator: (value) => [MAP_MODES.USERS, MAP_MODES.SITES, MAP_MODES.SINGLE_EDIT].includes(value)
    },
    // For SINGLE_EDIT mode
    lat: { type: Number, default: 1.3521 },
    lng: { type: Number, default: 103.8198 }
});

const emit = defineEmits(['update:lat', 'update:lng']);

const mapContainer = ref(null);
const isMapError = ref(false);
const errorMessage = ref('');

let map = null;
let tileLayer = null;
let resizeObserver = null;
let themeObserver = null;
let visibilityInterval = null;
let singleMarker = null; // Track single marker instance

const getThemeLayerUrl = (theme) => {
    return theme === 'light' 
        ? 'https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}.png'
        : 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png';
};

const updateMapTheme = () => {
    if (!map || !window.L) return;
    const currentTheme = document.documentElement.getAttribute('data-theme') || 'dark';
    const newUrl = getThemeLayerUrl(currentTheme);
    
    console.log(`[UnifiedMap] Updating tile layer to: ${currentTheme}`);
    
    if (tileLayer) {
        tileLayer.setUrl(newUrl);
    } else {
        tileLayer = window.L.tileLayer(newUrl, {
            maxZoom: 19,
            attribution: '&copy; CARTO'
        }).addTo(map);
    }
};

const initMap = async () => {
    try {
        isMapError.value = false;
        
        await nextTick();
        await new Promise(r => setTimeout(r, 500));
        
        if (!window.L) {
            // Retry logic or fail
             let attempts = 0;
            while (!window.L && attempts < 20) {
                await new Promise(r => setTimeout(r, 200));
                attempts++;
            }
             if (!window.L) {
                isMapError.value = true;
                errorMessage.value = 'Leaflet library failed to load.';
                return;
            }
        }

        if (!mapContainer.value) return;

        let items = [];
        try {
            // Retrieve User context
            const savedUser = localStorage.getItem('auth_user');
            let userId = null;
            if (savedUser) {
                try {
                    const user = JSON.parse(savedUser);
                    userId = user.id || user.user_id; // Handle both potential keys
                } catch (e) {
                    console.error('[UnifiedMap] Failed to parse auth_user', e);
                }
            }

            if (props.mode === MAP_MODES.USERS) {
                // Admin mode usually fetches all Users, but we can pass context if needed
                items = await api.getUsers();
            } else if (props.mode === MAP_MODES.SITES) {
                // Client mode requires user_id to filter sites
                items = await api.getSites({ user_id: userId });
            }
        } catch (apiErr) {
            console.error('[UnifiedMap] API Fetch Error:', apiErr);
            // Continue map initialization even if API fetch fails, just without markers
        }
        
        if (map) {
            map.remove();
            map = null;
            singleMarker = null;
            tileLayer = null;
        }

        // Default center
        let center = [1.3521, 103.8198];
        let zoom = 11;

        if (props.mode === MAP_MODES.SINGLE_EDIT) {
             center = [props.lat || 1.3521, props.lng || 103.8198];
             zoom = 15;
        }

        map = window.L.map(mapContainer.value, {
            zoomControl: false,
            attributionControl: true
        }).setView(center, zoom);

        updateMapTheme();
        window.L.control.zoom({ position: 'bottomright' }).addTo(map);

        if (props.mode === MAP_MODES.SINGLE_EDIT) {
            renderSingleEditMarker(center);
            
            // Map click to move marker
            map.on('click', (e) => {
                const { lat, lng } = e.latlng;
                updateSingleMarker(lat, lng);
            });
        } else if (items && items.length > 0) {
            items.forEach(item => {
                if (item.lat && item.lng) {
                    if (props.mode === MAP_MODES.USERS) {
                        renderUserMarker(item);
                    } else {
                        renderSiteMarker(item);
                    }
                }
            });
        }

        const invalidate = () => {
            if (map) {
                map.invalidateSize();
            }
        };

        invalidate();
        setTimeout(invalidate, 200);
        setTimeout(invalidate, 1000);
        setTimeout(invalidate, 3000);

    } catch (err) {
        isMapError.value = true;
        errorMessage.value = 'Map initialization failed.';
        isMapError.value = true;
        errorMessage.value = 'Map initialization failed.';
        console.error(`[UnifiedMap] Init Fail:`, err);
    }
};

const renderUserMarker = (user) => {
    const marker = window.L.circleMarker([user.lat, user.lng], {
        radius: 8,
        fillColor: user.status === USER_STATUS.ACTIVE ? '#10b981' : '#64748b',
        color: '#fff',
        weight: 2,
        opacity: 1,
        fillOpacity: 0.8
    }).addTo(map);

    marker.bindPopup(`
        <div class="map-popup">
            <h4 class="popup-title">${user.user_name}</h4>
            <p class="popup-info">Type: ${user.user_type}</p>
            <p class="popup-info">UEN: ${user.uen}</p>
            <div class="popup-status ${user.status}">${user.status.toUpperCase()}</div>
        </div>
    `, { className: 'custom-leaflet-popup' });
};

const renderSiteMarker = (site) => {
    const marker = window.L.circleMarker([site.lat, site.lng], {
        radius: 10,
        fillColor: '#3b82f6',
        color: '#fff',
        weight: 2,
        opacity: 1,
        fillOpacity: 0.8
    }).addTo(map);

    marker.bindPopup(`
        <div class="map-popup">
            <h4 class="popup-title">${site.site_name}</h4>
            <p class="popup-info">Location: ${site.location || 'N/A'}</p>
            <div class="popup-stats">
                <span class="p-stat"><strong>${site.worker_count || 0}</strong> Workers</span>
                <span class="p-stat"><strong>${site.device_count || 0}</strong> Devices</span>
            </div>
        </div>
    `, { className: 'custom-leaflet-popup' });
};

// --- SINGLE EDIT MODE HELPERS ---

const renderSingleEditMarker = (latlng) => {
    if (singleMarker) {
        singleMarker.setLatLng(latlng);
        return;
    }

    const icon = window.L.divIcon({
        className: 'custom-pin-icon',
        html: `<div style="background-color: #ef4444; width: 16px; height: 16px; border-radius: 50%; border: 2px solid white; box-shadow: 0 0 10px rgba(239, 68, 68, 0.5);"></div>`,
        iconSize: [20, 20],
        iconAnchor: [10, 10]
    });

    singleMarker = window.L.marker(latlng, {
        draggable: true,
        icon: icon
    }).addTo(map);

    singleMarker.on('dragend', (e) => {
        const { lat, lng } = e.target.getLatLng();
        updateSingleMarker(lat, lng);
    });
};

const updateSingleMarker = (lat, lng) => {
    // Round to reasonable precision
    const roundedLat = Number(lat.toFixed(6));
    const roundedLng = Number(lng.toFixed(6));
    
    emit('update:lat', roundedLat);
    emit('update:lng', roundedLng);

    // Update internal marker position if moved via click (not drag)
    if (singleMarker) {
        const cur = singleMarker.getLatLng();
        if (cur.lat !== lat || cur.lng !== lng) {
             singleMarker.setLatLng([lat, lng]);
        }
    }
};

// Watch for prop changes to update marker position (User typing in inputs)
watch(() => [props.lat, props.lng], ([newLat, newLng]) => {
    if (props.mode === MAP_MODES.SINGLE_EDIT && map && singleMarker) {
        const cur = singleMarker.getLatLng();
        // Only update if significantly different to avoid drag loops
        if (Math.abs(cur.lat - newLat) > 0.0001 || Math.abs(cur.lng - newLng) > 0.0001) {
            singleMarker.setLatLng([newLat, newLng]);
            map.panTo([newLat, newLng]);
        }
    }
});

onMounted(() => {
    initMap();

    // Theme Observer
    themeObserver = new MutationObserver(() => {
        console.log('[UnifiedMap] Theme change detected');
        updateMapTheme();
    });
    themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] });

    if (mapContainer.value) {
        resizeObserver = new ResizeObserver(() => {
            if (map) map.invalidateSize();
        });
        resizeObserver.observe(mapContainer.value);
    }

    // Keep-alive polling for visibility (handles complex grids)
    visibilityInterval = setInterval(() => {
        if (map) map.invalidateSize();
    }, 2000);
});

onUnmounted(() => {
    if (map) {
        map.remove();
        map = null;
    }
    if (resizeObserver) resizeObserver.disconnect();
    if (themeObserver) themeObserver.disconnect();
    if (visibilityInterval) clearInterval(visibilityInterval);
});

// React to mode changes (though unlikely for dashboard usage)
watch(() => props.mode, initMap);
</script>

<template>
    <div class="unified-map-wrapper">
        <div ref="mapContainer" class="map-container"></div>
        
        <!-- Error Overlay -->
        <div v-if="isMapError" class="map-error-overlay">
            <div class="error-content">
                <i class="ri-error-warning-line"></i>
                <p>{{ errorMessage }}</p>
                <button @click="initMap" class="retry-btn">Retry Initialization</button>
            </div>
        </div>

        <div v-else class="map-overlay">
            <div class="legend">
                <div v-if="mode === MAP_MODES.USERS" class="legend-content">
                    <div class="legend-item"><span class="dot active"></span> Active User</div>
                    <div class="legend-item"><span class="dot inactive"></span> Pending/Inactive</div>
                </div>
                <div v-else class="legend-content">
                    <div class="legend-item"><span class="dot site"></span> Site Location</div>
                </div>
            </div>
        </div>
    </div>
</template>

<style>
/* Global Leaflet Custom Styles */
.custom-leaflet-popup .leaflet-popup-content-wrapper {
    background: var(--color-surface);
    color: var(--color-text-primary);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: 12px;
}

.custom-leaflet-popup .leaflet-popup-tip {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
}

.map-popup .popup-title {
    margin: 0 0 8px 0;
    font-size: 15px;
    font-weight: 700;
}

.map-popup .popup-info {
    margin: 0 0 4px 0;
    font-size: 12px;
    color: var(--color-text-secondary);
}

.popup-status {
    display: inline-block;
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 10px;
    font-weight: 700;
    margin-top: 12px;
}

.popup-status.active { background: rgba(16, 185, 129, 0.1); color: #10b981; }
.popup-status.pending { background: rgba(100, 116, 139, 0.1); color: #64748b; }

.popup-stats {
    display: flex;
    gap: 12px;
    margin-top: 12px;
    padding-top: 12px;
    border-top: 1px solid var(--color-border);
}

.p-stat {
    font-size: 11px;
    color: var(--color-text-secondary);
}

.p-stat strong {
    color: var(--color-text-primary);
}

.leaflet-container {
    background: var(--color-bg) !important;
}
</style>

<style scoped>
.unified-map-wrapper {
    position: relative;
    width: 100%;
    height: 100%;
    min-height: 480px;
    border-radius: var(--radius-sm);
    overflow: hidden;
}

.map-container {
    width: 100%;
    height: 100%;
    background: var(--color-bg);
}

.map-overlay {
    position: absolute;
    top: 16px;
    right: 16px;
    z-index: 1000;
    pointer-events: none;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 12px;
}

.legend {
    background: var(--color-surface);
    backdrop-filter: blur(8px);
    border: 1px solid var(--color-border);
    padding: 12px;
    border-radius: 8px;
}
.legend-content {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.legend-item {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 11px;
    font-weight: 500;
    color: var(--color-text-primary);
}

.dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
}

.dot.active { background: #10b981; box-shadow: 0 0 8px #10b981; }
.dot.inactive { background: #64748b; }
.dot.site { background: #3b82f6; box-shadow: 0 0 8px #3b82f6; }

.map-error-overlay {
    position: absolute;
    inset: 0;
    background: rgba(13, 17, 23, 0.9);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2000;
    backdrop-filter: blur(4px);
    text-align: center;
}

.error-content {
    padding: 32px;
}

.error-content i {
    font-size: 48px;
    color: #ef4444;
    margin-bottom: 16px;
}

.error-content p {
    color: var(--color-text-secondary);
    margin-bottom: 24px;
    font-size: 14px;
}

.retry-btn {
    padding: 10px 20px;
    background: var(--color-accent);
    color: white;
    border: none;
    border-radius: var(--radius-sm);
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
}

.retry-btn:hover {
    background: var(--color-accent-hover);
    transform: translateY(-1px);
}
</style>
