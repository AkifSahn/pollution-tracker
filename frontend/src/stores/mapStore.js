import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useMapStore = defineStore('mapStore', () => {
    const markers = ref([]) // [{ lat, lng, value }]
    var graphBound = ref([])

    function addMarker(marker) {
        markers.value.push(marker)
    }

    function clearMarkers() {
        markers.value = []
    }

    return {
        markers,
        addMarker,
        clearMarkers,

        graphBound,
    }
})

