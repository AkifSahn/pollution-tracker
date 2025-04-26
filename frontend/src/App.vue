<script setup>
import { onMounted } from 'vue'
import HeatMap from './components/HeatMap.vue'
import PollutionGraph from './components/PollutionGraph.vue'
import Header from './components/Header.vue'
import Footer from './components/Footer.vue'
import { useThemeStore } from './stores/themeStore'

const themeStore = useThemeStore()

onMounted(() => {
    // Initialize theme based on user preference
    themeStore.initTheme()
})
</script>

<template>
    <div class="app-container">
        <Header />

        <main class="main-content container mx-auto px-4 py-6">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div class="card p-4">
                    <HeatMap />
                </div>
                <div class="card p-4">
                    <PollutionGraph />
                </div>
            </div>
        </main>

        <Footer />
    </div>
</template>

