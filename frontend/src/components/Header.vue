<template>
    <header class="sticky top-0 z-20 transition-colors duration-300" :style="{
        backgroundColor: 'var(--header-bg)',
        color: 'var(--header-text)'
    }">
        <div class="container mx-auto px-4 py-3">
            <div class="flex justify-between items-center">
                <!-- Logo -->
                <div class="flex items-center space-x-2">
                    <font-awesome-icon icon="map-location-dot" class="h-8 w-8" />
                    <h1 class="text-xl font-bold">Pollution Tracker</h1>
                </div>

                <div class="flex items-center space-x-4">
                    <button @click="toggleDarkMode"
                        class="p-2 rounded-full cursor-pointer border border-transparent transition-all duration-300 hover:border-white dark-mode-toggle">
                        <font-awesome-icon v-if="themeStore.isDarkMode" icon="sun" class="h-6 w-6" />
                        <font-awesome-icon v-else icon="moon" class="h-6 w-6" />
                    </button>

                    <button @click="toggleNotifications"
                        class="flex items-center transition-colors duration-300 shadow-lg px-4 py-2 rounded-sm hover:cursor-pointer"
                        :style="{ backgroundColor: 'var(--secondary-color)' }" type="button">
                        <span v-if="notifications.length === 0">
                            <font-awesome-icon :icon="['far', 'bell']" class="mr-2" />
                        </span>
                        <span v-else>
                            <font-awesome-icon :icon="['fas', 'bell']" class="mr-2" />
                        </span>
                        Bildirimler {{ notifications.length }}
                    </button>
                </div>
            </div>
        </div>

        <div v-if="showNotifications"
            class="absolute top-12 right-0 mt-2 dark:bg-gray-800 opacity-90 w-64 h-48 rounded-md shadow-lg overflow-y-auto z-11 my-5 mx-5">
            <div class="p-3 font-semibold border-b border-gray-200 dark:border-gray-700"
                :style="{ backgroundColor: 'var(--primary-color)', color: 'var(--header-text)' }">
                Bildirimler
            </div>
            <ul>
                <li v-for="(notification, index) in notifications" :key="index"
                    class="p-3 hover:bg-gray-700 transition-colors duration-200">
                    <div class="flex justify-between">
                        <span class="text-sm text-gray-200" v-html="notification"></span>
                        <button @click="popNotification(index)"
                            class="text-gray-300 hover:text-red-400 transition-colors">
                            <font-awesome-icon icon="times" class="h-5 w-5" />
                        </button>
                    </div>
                </li>
                <li v-if="notifications.length === 0" class="p-2 text-center text-gray-400">
                    Henüz bir bildirim yok.
                </li>
            </ul>
        </div>
    </header>
</template>

<script>
import { useMapStore } from '../stores/mapStore'
import { useThemeStore } from '../stores/themeStore'

export default {
    data() {
        return {
            showNotifications: false,
            notifications: [],
            ws: null,
        };
    },
    setup() {
        const themeStore = useThemeStore()
        return { themeStore }
    },
    mounted() {
        this.initWs();
    },
    methods: {
        toggleDarkMode() {
            this.themeStore.toggleDarkMode();
        },
        toggleNotifications() {
            this.showNotifications = !this.showNotifications;
        },
        popNotification(index) {
            this.notifications.splice(index, 1)
        },
        initWs() {
            const mapStore = useMapStore();
            this.ws = new WebSocket("ws://localhost:3000/ws")
            this.ws.onopen = (e) => {
                console.log("Websocket connection established")
            };
            this.ws.onmessage = (e) => {
                let d = JSON.parse(e.data)
                this.notifications.push(`${d.message}<br>pollutant: ${d.pollutant}<br>lat: ${d.latitude}, lng: ${d.longitude}<br>val: ${d.value}`)

                // Push the notification to Pinia store
                if (d.latitude && d.longitude && d.value) {
                    mapStore.addMarker({
                        latitude: d.latitude,
                        longitude: d.longitude,
                        value: d.value,
                        pollutant: d.pollutant,
                    })
                }
            };

            this.ws.onclose = (e) => {
                this.ws = null
                console.log("Ws connection closed")
            };
            this.ws.onerror = (e) => {
                console.log("An error occured", e)
            };
        },
    }
}
</script>

<style scoped>
.dark-mode-toggle:hover {
    transform: scale(1.05);
    box-shadow: 0 0 5px rgba(255, 255, 255, 0.5);
}
</style>
