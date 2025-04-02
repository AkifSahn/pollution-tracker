<script setup>
import HeatMap from './components/HeatMap.vue'
import PollutionGraph from './components/PollutionGraph.vue'
import NotificationPanel from './components/NotificationPanel.vue'
import RegionView from './components/RegionView.vue'
</script>

<template>
    <div class="bg-slate-100 h-screen">
        <NotificationPanel/>
        <HeatMap/>
    </div>
</template>
