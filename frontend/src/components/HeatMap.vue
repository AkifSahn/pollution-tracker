<!-- HeatMap.vue -->
<template>
    <div class="max-w-7xl mx-auto px-4 py-8 my-10">

            <div id="map" class="flex h-96 w-full rounded-md"></div>

            <!-- Buttons -->
            <div class="flex justify-center gap-10 my-3">
                <button
                    @click="fetchData"
                    class="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600">
                    Fetch & Map Heatmap
                </button>

                <div class="flex flex-col justify-center items-center">
                    <p>Zaman aralığı</p>
                    <div class="flex justify-center items-center gap-3"> 
                        <label class="block w-4 text-center">{{ rangeValue }}h</label>
                        <input v-model="rangeValue" class="block" type="range" name="slider" value="24" min="1" max="24">
                    </div>
                </div>
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

                rangeValue: 24,
            };
        },
        mounted() {
            this.initMap();
            this.fetchData();
        },
        methods: {
            initMap() {
                this.map = L.map('map', {
                    worldCopyJump: false
                }).setView([41.0082,28.9784], 5); // Center on İstanbul

                L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                    attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
                    noWrap: true,
                }).addTo(this.map);

                const bounds = [[-90, -180], [90, 180]];
                this.map.setMaxBounds(bounds);
                this.map.options.minZoom = 2;
            },

            async fetchData() {
                this.loading = true;
                this.pollutions = [];

                try {
                    const now = new Date();
                    const yesterday = new Date(now);
                    yesterday.setHours(now.getHours()-this.rangeValue);

                    function formatDate(date) {
                        const year = date.getFullYear();
                        const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are 0-based
                        const day = String(date.getDate()).padStart(2, '0');
                        const hours = String(date.getHours()).padStart(2, '0');
                        const minutes = String(date.getMinutes()).padStart(2, '0');
                        const seconds = String(date.getSeconds()).padStart(2, '0');
                        return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
                    }

                    const toDate = formatDate(now);
                    const fromDate = formatDate(yesterday);
                    console.log(toDate)
                    console.log(fromDate)

                    const url = `http://127.0.0.1:3000/api/anomalies?from=${encodeURIComponent(fromDate)}&to=${encodeURIComponent(toDate)}`;
                    console.log(url)
                    const response = await fetch(url);
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

                // // Remove previous markers
                // this.map.eachLayer(layer => {
                //     if (layer instanceof L.Marker) {
                //         this.map.removeLayer(layer);
                //     }
                // });

                // this.pollutions.forEach(item => {
                //     const marker = L.marker([item.latitude, item.longitude])
                //         .addTo(this.map)
                //         .bindPopup(`lat: ${item.latitude}<br>lon: ${item.longitude}<br>val: ${item.value}`)
                // });

                const heatData = this.pollutions.map(item => [
                    item.latitude,
                    item.longitude,
                    item.value*100
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
        },
    };
</script>

