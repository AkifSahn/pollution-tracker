<!-- HeatMap.vue -->
<template>
    <div class="transition-colors duration-300 relative">
        <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-bold">Pollution Heat Map</h2>
            <button @click="showFullScreen = true"
                class="fullscreen-btn p-2 rounded-full transition-all duration-300 hover:bg-opacity-10 hover:bg-white cursor-pointer w-10 h-10 flex items-center justify-center">
                <font-awesome-icon icon="expand" class="h-5 w-5" />
            </button>
        </div>

        <p class="mb-4"><i><b>{{ dataStartDate }}</b></i> tarihinden itibaren olan veriler gösteriliyor.</p>

        <div id="map" class="border-2 shadow-lg h-96 w-full rounded-lg z-10 transition-colors duration-300"
            :style="{ borderColor: 'var(--primary-color)' }"></div>

        <div class="mt-6 space-y-4">
            <div class="flex flex-col md:flex-row md:justify-between md:items-center gap-4">
                <div class="flex-1">
                    <p class="font-medium mb-2">Zaman aralığı</p>
                    <div class="flex flex-col sm:flex-row items-start sm:items-center gap-4">
                        <div class="flex items-center gap-2">
                            <label class="block w-14 text-center">{{ rangeValueDay }} gün</label>
                            <input @input="fetchData" v-model="rangeValueDay" class="block w-full" type="range"
                                name="daySlider" min="0" max="30" :style="{ accentColor: 'var(--primary-color)' }">
                        </div>

                        <div class="flex items-center gap-2">
                            <label class="block w-14 text-center">{{ rangeValueHour }} saat</label>
                            <input @input="fetchData" v-model="rangeValueHour" class="block w-full" type="range"
                                name="hourSlider" min="1" max="23" :style="{ accentColor: 'var(--primary-color)' }">
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <ModalFullScreen :show="showFullScreen" title="Pollution Heat Map" @close="closeFullScreen">
            <div class="h-full flex flex-col p-4">
                <p class="mb-3"><i><b>{{ dataStartDate }}</b></i> tarihinden itibaren olan veriler gösteriliyor.</p>

                <div id="fullscreen-map" class="flex-grow border-2 rounded-lg transition-colors duration-300"
                    :style="{ borderColor: 'var(--primary-color)' }"></div>

                <div class="mt-4">
                    <p class="font-medium mb-2">Zaman aralığı</p>
                    <div class="flex flex-wrap items-center gap-4">
                        <div class="flex items-center gap-2">
                            <label class="block w-14 text-center">{{ rangeValueDay }} gün</label>
                            <input @input="fetchData" v-model="rangeValueDay" class="block w-40" type="range" min="0"
                                max="30" :style="{ accentColor: 'var(--primary-color)' }">
                        </div>

                        <div class="flex items-center gap-2">
                            <label class="block w-14 text-center">{{ rangeValueHour }} saat</label>
                            <input @input="fetchData" v-model="rangeValueHour" class="block w-40" type="range" min="1"
                                max="23" :style="{ accentColor: 'var(--primary-color)' }">
                        </div>
                    </div>
                </div>
            </div>
        </ModalFullScreen>
    </div>
</template>

<script>
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import 'leaflet.heat';

import { useMapStore } from "../stores/mapStore";
import { watch, nextTick } from "vue";
import { fetchPollutions, fetchAnomaliesOfRange } from "../api";
import ModalFullScreen from "./ModalFullScreen.vue";

export default {
    components: {
        ModalFullScreen
    },
    data() {
        return {
            loading: false,
            pollutions: [],
            map: null,
            fullScreenMap: null,
            heatLayer: null,
            fullScreenHeatLayer: null,
            showFullScreen: false,

            rangeValueDay: 0,
            rangeValueHour: 23,
            dataStartDate: null,

            topLeftPos: [41.807221, 26.371056],
            bottomRightPos: [36.887526, 44.700193],

            topLeftMarker: null,
            bottomRightMarker: null,
            topLeftMarkerFullScreen: null,
            bottomRightMarkerFullScreen: null,
            mapStore: null,
        };
    },
    mounted() {
        this.mapStore = useMapStore()

        this.initMap();
        this.fetchData();
        this.fetchAnomalies();

        window.addEventListener('resize', this.handleResize);

        watch(
            () => this.mapStore.markers,
            (newMarkers) => {
                createAnomalyMarkers(newMarkers);
            },
            { deep: true }
        )

        watch(
            () => this.showFullScreen,
            (newVal) => {
                if (newVal) {
                    nextTick(() => {
                        this.initFullScreenMap();
                    });
                }
            }
        );
    },
    beforeUnmount() {
        window.removeEventListener('resize', this.handleResize);
    },
    methods: {
        formatDate(date) {
            const year = date.getFullYear();
            const month = String(date.getMonth() + 1).padStart(2, '0');
            const day = String(date.getDate()).padStart(2, '0');
            const hours = String(date.getHours()).padStart(2, '0');
            const minutes = String(date.getMinutes()).padStart(2, '0');
            return `${year}-${month}-${day} ${hours}:${minutes}:00`;
        },
        createAnomalyMarkers(markers) {
            markers.forEach((marker) => {
                L.marker([marker.latitude, marker.longitude])
                    .bindPopup(`Value: ${marker.value}<br>Pollutant: ${marker.pollutant}`)
                    .addTo(this.map)

                if (this.fullScreenMap) {
                    L.marker([marker.latitude, marker.longitude])
                        .bindPopup(`Value: ${marker.value}<br>Pollutant: ${marker.pollutant}`)
                        .addTo(this.fullScreenMap)
                }
            })

        },

        handleResize() {
            if (this.map) {
                this.map.invalidateSize();
            }
            if (this.fullScreenMap) {
                this.fullScreenMap.invalidateSize();
            }
        },
        closeFullScreen() {
            this.showFullScreen = false;
            this.fullScreenMap = null;
            this.fullScreenHeatLayer = null;
        },
        initMap() {
            this.map = L.map('map', {
                worldCopyJump: false
            }).setView([41.0082, 28.9784], 5); // Center on İstanbul

            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
                noWrap: true,
            }).addTo(this.map);

            const bounds = [[-90, -180], [90, 180]];
            this.map.setMaxBounds(bounds);
            this.map.options.minZoom = 2;

            const customDivIcon = L.divIcon({
                html: "<div style='background-color: red; width: 15px; height: 15px; border-radius: 50%; opacity: 70%;'></div>",
                className: "",
                iconSize: [15, 15],
            });

            // Add draggable markers
            this.topLeftMarker = L.marker(this.topLeftPos, {
                draggable: true,
                icon: customDivIcon,
            }).addTo(this.map);

            this.bottomRightMarker = L.marker(this.bottomRightPos, {
                draggable: true,
                icon: customDivIcon,
            }).addTo(this.map);

            this.bottomRightMarker.on("drag", (e) => {
                this.updateRect();
            })

            this.bottomRightMarker.on("dragend", (e) => {
                this.updateGraphBounds();
            })

            this.topLeftMarker.on("drag", (e) => {
                this.updateRect();
            })

            this.topLeftMarker.on("dragend", (e) => {
                this.updateGraphBounds();
            })

            this.updateRect();
        },
        initFullScreenMap() {
            if (!document.getElementById('fullscreen-map')) {
                return;
            }

            this.fullScreenMap = L.map('fullscreen-map', {
                worldCopyJump: false
            });

            // Copy view from the regular map
            this.fullScreenMap.setView(this.map.getCenter(), this.map.getZoom());

            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
                noWrap: true,
            }).addTo(this.fullScreenMap);

            const bounds = [[-90, -180], [90, 180]];
            this.fullScreenMap.setMaxBounds(bounds);
            this.fullScreenMap.options.minZoom = 2;

            this.topLeftMarkerFullScreen = L.marker(this.topLeftMarker.getLatLng(), {
                draggable: true,
            }).addTo(this.fullScreenMap);

            this.bottomRightMarkerFullScreen = L.marker(this.bottomRightMarker.getLatLng(), {
                draggable: true,
            }).addTo(this.fullScreenMap);

            this.bottomRightMarkerFullScreen.on("drag", (e) => {
                this.updateRectFullScreen();
                // Also update original map
                this.bottomRightMarker.setLatLng(this.bottomRightMarkerFullScreen.getLatLng());
                this.updateRect();
            });

            this.bottomRightMarkerFullScreen.on("dragend", (e) => {
                this.updateGraphBounds();
            });

            this.topLeftMarkerFullScreen.on("drag", (e) => {
                this.updateRectFullScreen();
                // Also update original map
                this.topLeftMarker.setLatLng(this.topLeftMarkerFullScreen.getLatLng());
                this.updateRect();
            });

            this.topLeftMarkerFullScreen.on("dragend", (e) => {
                this.updateGraphBounds();
            });

            this.updateRectFullScreen();
            this.updateFullScreenHeatmap();

            // Give the map a moment to render before calling invalidateSize
            setTimeout(() => {
                this.fullScreenMap.invalidateSize();
            }, 100);
        },
        updateRect() {
            const bounds = [this.topLeftMarker.getLatLng(), this.bottomRightMarker.getLatLng()];

            this.map.eachLayer(function (layer) {
                if (layer instanceof L.Rectangle) {
                    layer.remove();
                }
            });

            const rect = L.rectangle(bounds, {
                color: getComputedStyle(document.documentElement).getPropertyValue('--primary-color').trim() || "#3b82f6",
                weight: 1
            });

            rect.addTo(this.map);
        },
        updateGraphBounds() {
            var bounds = [this.topLeftMarker.getLatLng(), this.bottomRightMarker.getLatLng()];
            if (this.fullScreenMap) {
                bounds = [this.topLeftMarkerFullScreen.getLatLng(), this.bottomRightMarkerFullScreen.getLatLng()];
            }
            this.mapStore.graphBound = bounds
            console.log(bounds)
        },
        updateRectFullScreen() {
            if (!this.fullScreenMap) return;

            const bounds = [this.topLeftMarkerFullScreen.getLatLng(), this.bottomRightMarkerFullScreen.getLatLng()];

            this.fullScreenMap.eachLayer(function (layer) {
                if (layer instanceof L.Rectangle) {
                    layer.remove();
                }
            });

            const rect = L.rectangle(bounds, {
                color: getComputedStyle(document.documentElement).getPropertyValue('--primary-color').trim() || "#3b82f6",
                weight: 1
            });

            rect.addTo(this.fullScreenMap);
        },
        async fetchAnomalies() {
            try {
                const now = new Date();
                const before = new Date(now);
                before.setHours(now.getHours() - 4); // TODO: fix the database time difference issue. DB is 3 hours back
                const data = await fetchAnomaliesOfRange(this.formatDate(before), this.formatDate(now));
                console.log(data)
                const pollutions = data.pollutions || [];
                this.createAnomalyMarkers(pollutions)
            } catch (error) {
                console.log("An error occured while fetching anomalies!");
            }
        },
        async fetchData() {
            this.loading = true;
            this.pollutions = [];

            try {
                const now = new Date();
                const yesterday = new Date(now);
                yesterday.setHours(now.getHours() - this.rangeValueHour - this.rangeValueDay * 24);


                const toDate = this.formatDate(now);
                this.dataStartDate = this.formatDate(yesterday);
                this.mapStore.timeFrom = this.dataStartDate
                this.mapStore.timeTo = toDate

                const data = await fetchPollutions(this.dataStartDate, toDate);
                this.pollutions = data.pollutions || [];
                this.updateHeatmap();

                // Update fullscreen heatmap if it exists
                if (this.fullScreenMap) {
                    this.updateFullScreenHeatmap();
                }
            } catch (error) {
                console.error('Error fetching data:', error);
                this.pollutions = [];
            } finally {
                this.loading = false;
            }
        },
        updateHeatmap() {
            if (this.heatLayer) {
                this.map.removeLayer(this.heatLayer);
            }

            const heatData = this.pollutions.map(item => [
                item.latitude,
                item.longitude,
                item.value * 5
            ]);

            this.heatLayer = L.heatLayer(heatData, {
                radius: 25,
                blur: 15,
                max: 1.0,
            }).addTo(this.map);

            if (this.pollutions.length > 0) {
                const bounds = L.latLngBounds(this.pollutions.map(item => [item.latitude, item.longitude]));
                this.map.fitBounds(bounds, { padding: [50, 50] });
            }
        },
        updateFullScreenHeatmap() {
            if (!this.fullScreenMap) return;

            if (this.fullScreenHeatLayer) {
                this.fullScreenMap.removeLayer(this.fullScreenHeatLayer);
            }

            const heatData = this.pollutions.map(item => [
                item.latitude,
                item.longitude,
                item.value * 5
            ]);

            this.fullScreenHeatLayer = L.heatLayer(heatData, {
                radius: 25,
                blur: 15,
                max: 1.0,
            }).addTo(this.fullScreenMap);

            if (this.pollutions.length > 0) {
                const bounds = L.latLngBounds(this.pollutions.map(item => [item.latitude, item.longitude]));
                this.fullScreenMap.fitBounds(bounds, { padding: [50, 50] });
            }
        }
    },
};
</script>

<style>
.leaflet-control-attribution {
    background-color: rgba(255, 255, 255, 0.8) !important;
    color: #333 !important;
}

.leaflet-popup-content-wrapper,
.leaflet-popup-tip {
    background-color: var(--card-bg) !important;
    color: var(--text-color) !important;
}

.fullscreen-btn {
    color: var(--primary-color);
    border: 1px solid var(--primary-color);
}

.fullscreen-btn:hover {
    transform: scale(1.05);
    box-shadow: 0 0 3px var(--primary-color);
}
</style>
