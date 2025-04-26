<template>
    <div class="transition-colors duration-300 relative">
        <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-bold">Pollution Density Graph</h2>
            <button @click="showFullScreen = true"
                class="fullscreen-btn p-2 rounded-full transition-all duration-300 hover:bg-opacity-10 hover:bg-white cursor-pointer w-10 h-10 flex items-center justify-center">
                <font-awesome-icon icon="expand" class="h-4 w-4" />
            </button>
        </div>

        <div class="flex flex-col md:flex-row justify-between items-start gap-6">
            <!-- Pollutant options -->
            <div class="w-full md:w-auto">
                <h3 class="font-medium mb-3">Pollutant Type</h3>
                <div class="flex flex-wrap gap-2">
                    <button v-for="(item, index) in pollutantOptions" :key="index"
                        @click="selectedPollutant = item; fetchData()"
                        class="px-3 py-1.5 rounded-md transition-colors duration-300" :style="{
                            backgroundColor: selectedPollutant === item ? 'var(--secondary-color)' : 'var(--primary-color)',
                            color: 'var(--header-text)'
                        }">
                        {{ item }}
                    </button>
                </div>
            </div>

            <!-- Refresh button -->
            <button @click="fetchData" class="px-4 py-2 rounded-md transition-colors duration-300" :style="{
                backgroundColor: 'var(--primary-color)',
                color: 'var(--header-text)'
            }">
                Refresh Graph
            </button>
        </div>

        <!-- Chart -->
        <div class="mt-6">
            <div ref="chartContainer" class="w-full h-96 border-2 rounded-lg p-4 transition-colors duration-300"
                :style="{ borderColor: 'var(--primary-color)' }"></div>
        </div>

        <!-- Fullscreen Modal -->
        <ModalFullScreen :show="showFullScreen" title="Pollution Density Graph" @close="closeFullScreen">
            <div class="h-full flex flex-col p-4">
                <div class="flex flex-wrap justify-between items-center gap-4 mb-4">
                    <div>
                        <h3 class="font-medium mb-3">Pollutant Type</h3>
                        <div class="flex flex-wrap gap-2">
                            <button v-for="(item, index) in pollutantOptions" :key="index"
                                @click="selectedPollutant = item; fetchData()"
                                class="px-3 py-1.5 rounded-md transition-colors duration-300" :style="{
                                    backgroundColor: selectedPollutant === item ? 'var(--secondary-color)' : 'var(--primary-color)',
                                    color: 'var(--header-text)'
                                }">
                                {{ item }}
                            </button>
                        </div>
                    </div>

                    <button @click="fetchData" class="px-4 py-2 rounded-md transition-colors duration-300" :style="{
                        backgroundColor: 'var(--primary-color)',
                        color: 'var(--header-text)'
                    }">
                        Refresh Graph
                    </button>
                </div>

                <div ref="fullscreenChartContainer"
                    class="flex-grow border-2 rounded-lg p-4 transition-colors duration-300"
                    :style="{ borderColor: 'var(--primary-color)' }"></div>
            </div>
        </ModalFullScreen>
    </div>
</template>

<script>
import * as d3 from "d3";
import { useMapStore } from "../stores/mapStore";
import { watch, nextTick } from "vue";
import { fetchPollutants, fetchRegionDensity } from "../api";
import { useThemeStore } from "../stores/themeStore";
import ModalFullScreen from "./ModalFullScreen.vue";

export default {
    components: {
        ModalFullScreen
    },
    data() {
        return {
            latitudeFrom: -90,
            latitudeTo: 90,
            longitudeFrom: -180,
            longitudeTo: 180,
            data: [],
            mapStore: null,
            themeStore: null,
            showFullScreen: false,

            pollutantOptions: [],
            selectedPollutant: null,
        };
    },
    mounted() {
        this.mapStore = useMapStore();
        this.themeStore = useThemeStore();

        watch(
            () => this.mapStore.graphBound,
            (newBound) => {
                this.latitudeFrom = Math.min(newBound[0].lat, newBound[1].lat);
                this.longitudeFrom = Math.min(newBound[0].lng, newBound[1].lng);
                this.latitudeTo = Math.max(newBound[0].lat, newBound[1].lat);
                this.longitudeTo = Math.max(newBound[0].lng, newBound[1].lng);
                this.fetchData()
                    .then(() => this.drawChart());
            }
        );

        watch(
            () => [this.mapStore.timeTo, this.mapStore.timeFrom],
            () => {
                this.fetchData()
                    .then(() => this.drawChart());
            }
        );

        // Redraw chart when theme changes
        watch(
            () => this.themeStore.isDarkMode,
            () => {
                this.drawChart();
                if (this.showFullScreen) {
                    this.drawFullScreenChart();
                }
            }
        );

        // Watch for fullscreen state change
        watch(
            () => this.showFullScreen,
            (newVal) => {
                if (newVal) {
                    // Draw fullscreen chart after modal is shown
                    nextTick(() => {
                        this.drawFullScreenChart();
                    });
                }
            }
        );

        // Handle window resize
        window.addEventListener('resize', this.handleResize);

        this.fetchAvailablePollutants();
    },
    beforeUnmount() {
        window.removeEventListener('resize', this.handleResize);
    },
    methods: {
        handleResize() {
            this.drawChart();
            if (this.showFullScreen) {
                this.drawFullScreenChart();
            }
        },

        closeFullScreen() {
            this.showFullScreen = false;
        },

        async fetchAvailablePollutants() {
            try {
                const jsonData = await fetchPollutants();
                this.pollutantOptions = jsonData.pollutants || [];
                this.selectedPollutant = this.pollutantOptions[0];

                if (this.selectedPollutant) {
                    await this.fetchData();
                }
            } catch (error) {
                console.error("Error fetching pollutants:", error);
            }
        },

        async fetchData() {
            try {
                const params = {
                    latFrom: this.latitudeFrom,
                    latTo: this.latitudeTo,
                    longFrom: this.longitudeFrom,
                    longTo: this.longitudeTo,
                    from: this.mapStore.timeFrom,
                    to: this.mapStore.timeTo
                };

                const jsonData = await fetchRegionDensity(this.selectedPollutant, params);

                if (jsonData.densities == null) {
                    this.data = [];
                } else {
                    this.data = jsonData.densities.map(e => ({
                        time: new Date(e.time),
                        value: e.density
                    }));
                }

            } catch (error) {
                console.error('Error fetching data:', error);
            }
            this.drawChart();
            if (this.showFullScreen) {
                this.drawFullScreenChart();
            }
        },

        drawChart() {
            this.drawChartInContainer(this.$refs.chartContainer);
        },

        drawFullScreenChart() {
            if (this.$refs.fullscreenChartContainer) {
                this.drawChartInContainer(this.$refs.fullscreenChartContainer);
            }
        },

        drawChartInContainer(container) {
            if (!container) return;

            d3.select(container).selectAll("*").remove();

            const margin = { top: 20, right: 30, bottom: 60, left: 40 };
            const width = container.clientWidth - margin.left - margin.right;
            const height = container.clientHeight - margin.top - margin.bottom;

            // Theme colors
            const isDark = this.themeStore.isDarkMode;
            const textColor = getComputedStyle(document.documentElement).getPropertyValue('--text-color').trim();
            const primaryColor = getComputedStyle(document.documentElement).getPropertyValue('--primary-color').trim();
            const secondaryColor = getComputedStyle(document.documentElement).getPropertyValue('--secondary-color').trim();

            const svg = d3.select(container)
                .append('svg')
                .attr('width', width + margin.left + margin.right)
                .attr('height', height + margin.top + margin.bottom)
                .append('g')
                .attr('transform', `translate(${margin.left},${margin.top})`);

            // No data message
            if (!this.data || this.data.length === 0) {
                svg.append('text')
                    .attr('x', width / 2)
                    .attr('y', height / 2)
                    .attr('text-anchor', 'middle')
                    .attr('fill', textColor)
                    .text('No data available for the selected criteria');
                return;
            }

            // X scale (time)
            const x = d3.scaleBand()
                .domain(this.data.map(d => d.time))
                .range([0, width])
                .padding(0.1);

            // Y scale (value)
            const y = d3.scaleLinear()
                .domain([0, Math.max(100, d3.max(this.data, d => d.value))])
                .nice()
                .range([height, 0]);

            // X Axis
            svg.append('g')
                .attr('transform', `translate(0,${height})`)
                .call(
                    d3.axisBottom(x)
                        .tickFormat(d3.timeFormat('%m-%d %H:%M'))
                        .tickValues(x.domain().filter((_, i) => i % Math.ceil(this.data.length / 10) === 0))
                )
                .selectAll("text")
                .attr("transform", "rotate(-25)")
                .style("text-anchor", "end")
                .attr("fill", textColor);

            // Y Axis
            svg.append('g')
                .call(d3.axisLeft(y))
                .selectAll("text")
                .attr("fill", textColor);

            // Grid lines
            svg.append('g')
                .attr('class', 'grid')
                .call(d3.axisLeft(y)
                    .tickSize(-width)
                    .tickFormat(''))
                .selectAll("line")
                .attr("stroke", isDark ? "rgba(255,255,255,0.1)" : "rgba(0,0,0,0.1)");

            // Title
            svg.append("text")
                .attr("x", width / 2)
                .attr("y", -margin.top / 2)
                .attr("text-anchor", "middle")
                .attr("fill", textColor)
                .style("font-size", "14px")
                .text(`${this.selectedPollutant} Density Over Time`);

            // X Axis label
            svg.append("text")
                .attr("x", width / 2)
                .attr("y", height + margin.bottom - 5)
                .attr("text-anchor", "middle")
                .attr("fill", textColor)
                .text("Time");

            // Y Axis label
            svg.append("text")
                .attr("transform", "rotate(-90)")
                .attr("x", -height / 2)
                .attr("y", -margin.left + 15)
                .attr("text-anchor", "middle")
                .attr("fill", textColor)
                .text("Density Value");

            // Bars
            svg.selectAll('.bar')
                .data(this.data)
                .enter()
                .append('rect')
                .attr('class', 'bar')
                .attr('x', d => x(d.time))
                .attr('y', d => y(d.value))
                .attr('width', x.bandwidth())
                .attr('height', d => height - y(d.value))
                .attr('fill', primaryColor)
                .attr('opacity', 0.8)
                .on('mouseover', function () {
                    d3.select(this).attr('opacity', 1);
                })
                .on('mouseout', function () {
                    d3.select(this).attr('opacity', 0.8);
                });
        }
    },
};
</script>

<style scoped>
/* Ensure svg is responsive */
svg {
    max-width: 100%;
    height: auto;
}

.fullscreen-btn {
    color: var(--primary-color);
    border: 1px solid var(--primary-color);
}

.fullscreen-btn:hover {
    transform: scale(1.05);
    box-shadow: 0 0 3px var(--primary-color);
}
</style>
