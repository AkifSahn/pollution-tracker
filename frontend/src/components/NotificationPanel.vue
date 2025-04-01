<template>
    <nav class="flex bg-linear-to-b from-cyan-400 to-cyan-600 justify-between items-center px-6 py-3 h-16 shadow-lg ">
        <!-- Logo -->
        <div>
            Logo
        </div>
        <h1>Pollution Tracker</h1>
        <div>
            <button @click="toggleNotifications" 
                class="bg-green-200 hover:bg-green-300 duration-300 shadow-lg px-3 py-1 rounded-sm hover:cursor-pointer" 
                type="submit">
                <span v-if="notifications.length === 0">
                    <font-awesome-icon :icon="['far', 'bell']"  />
                </span>
                <span v-else>
                    <font-awesome-icon :icon="['fas', 'bell']"/>
                </span>
                Bildirimler {{ notifications.length }} 
            </button>
        </div>

        <div v-if="showNotifications"
            class="absolute top-12 right-0 bg-gray-300 opacity-90 w-64 h-48 rounded-md shadow-lg overflow-y-auto z-11 my-5 mx-5"
            >
            <ul class="p-2 space-y-2">
                <li 
                    v-for="(notification, index) in notifications" :key="index"
                    class="p-2 bg-gray-100 rounded text-sm"
                    >{{ notification }}</li>
                <li v-if="notifications.length === 0"
                    class="p-2 text-gray-500 text-sm"
                    >Henüz bir bildirim yok.
                </li>
            </ul>

        </div>

    </nav>
</template>

<script>
    export default{
        data(){
            return {
                showNotifications: false,
                notifications: [],
            };
        },
        methods: {
            toggleNotifications(){
                this.showNotifications = !this.showNotifications;
                console.log(this.showNotifications);
            },
        }

    }
</script>

