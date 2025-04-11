<template>
    <nav
        class="flex bg-red-700 sticky top-0 z-20 opacity-90 text-white justify-between items-center px-6 py-3 h-16 shadow-lg ">
        <div>
            Logo
        </div>
        <h1>Pollution Tracker</h1>
        <div>
            <button @click="toggleNotifications"
                class="bg-orange-500 w-36 hover:bg-orange-400 duration-300 shadow-lg px-3 py-1 rounded-sm hover:cursor-pointer"
                type="submit">
                <span v-if="notifications.length === 0">
                    <font-awesome-icon :icon="['far', 'bell']" />
                </span>
                <span v-else>
                    <font-awesome-icon :icon="['fas', 'bell']" />
                </span>
                Bildirimler {{ notifications.length }}
            </button>
        </div>

        <div v-if="showNotifications"
            class="absolute top-12 right-0 bg-gray-300 opacity-90 w-64 h-48 rounded-md shadow-lg overflow-y-auto z-11 my-5 mx-5">
            <ul class="p-2 space-y-2">
                <li v-for="(notification, index) in notifications" :key="index"
                    class="p-2 bg-gray-100 text-black rounded text-sm">
                    <div class="flex justify-between items-center">
                        <span v-html="notification"></span>
                        <button @click="popNotification(index)"
                            class="bg-orange-500 hover:bg-orange-400 hover:cursor-pointer px-2 py-1 h-10 rounded-sm"
                            type="submit">X</button>
                    </div>
                </li>
                <li v-if="notifications.length === 0" class="p-2 text-gray-500 text-sm">Henüz bir bildirim yok.
                </li>
            </ul>

        </div>

    </nav>
</template>

<script>
import { useMapStore } from '../stores/mapStore'

export default {
    data() {
        return {
            showNotifications: false,
            notifications: [],
            ws: null,
        };
    },
    mounted() {
        this.initWs();
    },
    methods: {
        toggleNotifications() {
            this.showNotifications = !this.showNotifications;
            console.log(this.showNotifications);
        },
        popNotification(index) {
            this.notifications.splice(index, 1)
        },

        initWs() {
            const mapStore = useMapStore();
            this.ws = new WebSocket("ws://127.0.0.1:3000/ws")
            this.ws.onopen = (e) => {
                console.log("Websocket connection established")
            };
            this.ws.onmessage = (e) => {
                let d = JSON.parse(e.data)
                this.notifications.push(`${d.message}<br>pollutant: ${d.pollutant}<br>lat: ${d.latitude}, lng: ${d.longitude}<br>val: ${d.value}`)

                // Push the notification to Pinia store
                if (d.latitude && d.longitude && d.value) {
                    console.log("Adding to store")
                    mapStore.addMarker({
                        latitude: d.latitude,
                        longitude: d.longitude,
                        value: d.value,
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
