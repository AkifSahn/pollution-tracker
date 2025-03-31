<!-- Content.vue -->
<template>
    <div class="max-w-7xl mx-auto px-4 py-8">
        <button
            @click="fetchData"
            class="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 mb-4"
            >
            Fetch & Map Pollution Heatmap
        </button>


            <div class="flex justify-between items-center">
                <!-- Map Container -->
                <div id="map" class="flex h-96 w-5/6 rounded-md"></div>

                <form class="w-full max-w-md mx-auto p-4 bg-white rounded-xl shadow space-y-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700">Latitude</label>
                        <input type="text" class="mt-1 block w-full border rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" placeholder="latitude" />
                    </div>

                    <div>
                        <label class="block text-sm font-medium text-gray-700">Longitude</label>
                        <input type="text" class="mt-1 block w-full border rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" placeholder="longitude" />
                    </div>

                    <button @click="fetchData" type="submit" class="w-full bg-blue-600 text-white font-semibold py-2 px-4 rounded-md hover:bg-blue-700 transition">
                        Ekle
                    </button>
                </form>
            </div>

            <!-- Optional: List of data below the map -->
            <div v-if="pollutions.length > 0" class="mt-4 bg-gray-100 p-4 rounded-md">
                <h3 class="text-lg font-bold mb-2">Pollution Data</h3>
                <ul class="space-y-2">
                    <li v-for="(item, index) in pollutions" :key="index" class="border-b pb-2">
                        <p><strong>Time:</strong> {{ item.time }}</p>
                        <p><strong>Location:</strong> Lat {{ item.latitude }}, Lon {{ item.longitude }}</p>
                        <p><strong>Value (Intensity):</strong> {{ item.value }} ({{ item.pollutant }})</p>
                        <p><strong>Anomaly:</strong> {{ item.is_anomaly ? 'Yes' : 'No' }}</p>
                    </li>
                </ul>
            </div>
    </div>
</template>

<script>
    import L from 'leaflet';
    import 'leaflet/dist/leaflet.css';
    import 'leaflet.heat';

    export default {
        data() {
            return {
                loading: false,
                pollutions: [],
                map: null,
                heatLayer: null,
            };
        },
        mounted() {
            this.initMap();
            this.fetchData();
        },
        methods: {
            initMap() {
                this.map = L.map('map').setView([41.0082,28.9784], 5); // Center on İstanbul
                L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                    attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
                }).addTo(this.map);

            },

            async fetchData() {
                this.loading = true;
                this.pollutions = [];

                try {
                    const response = await fetch(
                        'http://127.0.0.1:3000/api/anomalies?from=2025-03-01%2010%3A10%3A10&to=2025-03-31%2023%3A10%3A10'
                    );
                    const jsonData = await response.json();
                    this.pollutions = jsonData.pollutions || []; // Fallback to empty array if pollutions is missing
                    this.updateHeatmap();
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
                    item.value*10 // Adjust divisor based on your value range
                ]);

                this.pollutions.forEach(item => {
                    const marker = L.marker([item.latitude, item.longitude])
                        .addTo(this.map)
                        .bindPopup("test")
                });

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
        },
    };
</script>

