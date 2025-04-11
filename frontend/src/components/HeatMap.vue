<!-- HeatMap.vue -->
<template>
    <div class="bg-white py-32">
        <div class="max-w-7xl mx-auto px-4">

            <p><i><b>{{ dataStartDate }}</b></i> tarihinden itibaren olan veriler gösteriliyor.</p>
            <div id="map" class="flex border-gray-600 border-2 shadow-lg h-96 w-full rounded-md z-10 "></div>

            <!-- Buttons -->
            <div class="flex justify-center gap-10 my-3">
                <div class="flex flex-col justify-center items-center">
                    <p>Zaman aralığı</p>
                    <div class="flex justify-center items-center gap-3"> 
                        <label class="block w-14 text-center">{{ rangeValueDay }} gün</label>
                        <input @input="fetchData" v-model="rangeValueDay" class="block accent-red-400" type="range" name="slider" value="0" min="0" max="30">

                        <label class="block w-14 text-center">{{ rangeValueHour }} saat</label>
                        <input @input="fetchData" v-model="rangeValueHour" class="block accent-red-400" type="range" name="slider" value="23" min="1" max="23">
                    </div>
                </div>
            </div>

        </div>
    </div>
</template>

<script>
    import L from 'leaflet';
    import 'leaflet/dist/leaflet.css';
    import 'leaflet.heat';

    import { useMapStore } from "../stores/mapStore";
    import { watch } from "vue";

    export default {
        data() {
            return {
                loading: false,
                pollutions: [],
                map: null,
                heatLayer: null,

                rangeValueDay: 0,
                rangeValueHour: 23,
                dataStartDate: null,

                topLeftPos: [41.807221, 26.371056],
                bottomRightPos: [36.887526, 44.700193],

                topLeftMarker: null,
                bottomRightMarker: null,
                mapStore: null,
            };
        },
        mounted() {
            this.mapStore = useMapStore()

            this.initMap();
            this.fetchData();

            watch(
                () => this.mapStore.markers,
                (newMarkers) => {
                    newMarkers.forEach((marker) => {
                        L.marker([this.marker.latitude, this.marker.longitude])
                            .bindPopup(`Value: ${this.marker.value}`)
                            .addTo(this.map)
                    })
                },
            )
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

                // Add draggable markers
                this.topLeftMarker = L.marker(this.topLeftPos, {
                    draggable: true,
                    // autoPan: true,
                }).addTo(this.map);

                this.bottomRightMarker = L.marker(this.bottomRightPos, {
                    draggable: true,
                    // autoPan: true,
                }).addTo(this.map);

                this.bottomRightMarker.on("drag", (e) => {
                    this.updateRect();
                })

                this.topLeftMarker.on("drag", (e) => {
                    this.updateRect();
                })

                this.updateRect();
            },

            updateRect(){
                const bounds = [this.topLeftMarker.getLatLng(), this.bottomRightMarker.getLatLng()];

                this.map.eachLayer(function(layer){
                    if (layer instanceof L.Rectangle) {
                        layer.remove();
                    }
                });

                L.rectangle(bounds, {
                    color:  "#ff7800",
                    weight: 1
                }).addTo(this.map);

                this.mapStore.graphBound = bounds;
            },

            async fetchData() {
                this.loading = true;
                this.pollutions = [];

                try {
                    const now = new Date();
                    const yesterday = new Date(now);
                    yesterday.setHours(now.getHours()-this.rangeValueHour-this.rangeValueDay*24);

                    function formatDate(date) {
                        const year = date.getFullYear();
                        const month = String(date.getMonth() + 1).padStart(2, '0');
                        const day = String(date.getDate()).padStart(2, '0');
                        const hours = String(date.getHours()).padStart(2, '0');
                        const minutes = String(date.getMinutes()).padStart(2, '0');
                        return `${year}-${month}-${day} ${hours}:${minutes}:00`;
                    }

                    const toDate = formatDate(now);
                    this.dataStartDate = formatDate(yesterday);
                    console.log(toDate)
                    console.log(this.dataStartDate)

                    const url = `http://127.0.0.1:3000/api/pollutions?from=${encodeURIComponent(this.dataStartDate)}&to=${encodeURIComponent(toDate)}`;
                    console.log(url)
                    const response = await fetch(url);
                    const jsonData = await response.json();
                    this.pollutions = jsonData.pollutions || [];
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
                    item.value*5
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

