<script setup>
import HeatMap from './components/HeatMap.vue'
import PollutionGraph from './components/PollutionGraph.vue'
import NotificationPanel from './components/NotificationPanel.vue'
import RegionView from './components/RegionView.vue'
</script>

<template>
    <div class="bg-gradient-to-b from-slate-100 to-slate-400 h-screen">
        <NotificationPanel/>
        <HeatMap/>
    </div>
</template>
