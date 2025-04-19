<template>
    <div class="flex-row justify-center bg-red-400 gap-x-14 py-10">

        <div class="place-self-center">
            <button @click="fetchData"
                class="block px-4 py-2 bg-blue-700 hover:bg-blue-600 hover:cursor-pointer text-white rounded-md"
                type="submit">Graph</button>
        </div>

        <div class="flex justify-center">
            <ul class="flex flex-col space-y-3 mx-10 p-4">
                <button v-for="(item, index) in pollutantOptions" :key="index"
                    @click="selectedPollutant = item; fetchData()" :class="[
                        'px-4 py-2 rounded-md text-white hover:cursor-pointer',
                        selectedPollutant === item ? 'bg-orange-500 hover:bg-orange-400' : 'bg-blue-700 hover:bg-blue-600'
                    ]">
                    {{ item }}
                </button>
            </ul>
            <div ref="chartContainer" class="w-full max-w-7xl h-96 place-self-center"></div>
        </div>

    </div>
</template>


<script>

import * as d3 from "d3";
import { useMapStore } from "../stores/mapStore";
import { watch } from "vue";

export default {
    data() {
        return {
            latitudeFrom: -90,
            latitudeTo: 90,
            longitudeFrom: -180,
            longitudeTo: 180,
            data: [],
            mapStore: null,

            pollutantOptions: [], // ['PM2.5', 'PM10', 'NO2', 'SO2', 'O3'],
            selectedPollutant: null,
        };
    },
    mounted() {
        this.mapStore = useMapStore();

        watch(
            () => this.mapStore.graphBound,
            (newBound) => {
                this.latitudeFrom = Math.min(newBound[0].lat, newBound[1].lat);
                this.longitudeFrom = Math.min(newBound[0].lng, newBound[1].lng);
                this.latitudeTo = Math.max(newBound[0].lat, newBound[1].lat);
                this.longitudeTo = Math.max(newBound[0].lng, newBound[1].lng);
            }
        );

        this.fetchAvailablePollutants()
            .then(() => this.fetchData())
            .then(() => this.drawChart());

    },
    methods: {

        async fetchAvailablePollutants() {
            try {
                const url = `http://127.0.0.1:3000/api/pollutants`;
                const response = await fetch(url);
                const jsonData = await response.json();

                this.pollutantOptions = jsonData.pollutants || [];
                this.selectedPollutant = this.pollutantOptions[0];
            } catch (error) {
                console.error("Error fetching data: ", error);
            }
        },

        async fetchData() {
            // Fetch the pollution data from backend
            // Let the user decide the lat-start, lat-end and long-start, long-end
            // User can choose these values from the map and see the graph for that rectangle. 
            // Density graph for overall region for the time range also given by user

            try {
                const url = `http://127.0.0.1:3000/api/region/density/${this.selectedPollutant}?` +
                    `latFrom=${encodeURIComponent(this.latitudeFrom)}&` +
                    `latTo=${encodeURIComponent(this.latitudeTo)}&` +
                    `longFrom=${encodeURIComponent(this.longitudeFrom)}&` +
                    `longTo=${encodeURIComponent(this.longitudeTo)}&` +
                    `from=${encodeURIComponent("2025-03-01 00:00:00")}&` +
                    `to=${encodeURIComponent("2025-04-28 23:00:00")}`;

                console.log("Fetching:", url);

                const response = await fetch(url);
                const jsonData = await response.json();

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
        },
        drawChart() {
            const container = this.$refs.chartContainer;

            d3.select(this.$refs.chartContainer).selectAll("*").remove();

            const margin = { top: 20, right: 30, bottom: 30, left: 40 };
            const width = container.clientWidth - margin.left - margin.right;
            const height = container.clientHeight - margin.top - margin.bottom;

            // Clear old chart
            d3.select(container).selectAll("*").remove();

            const svg = d3.select(container)
                .append('svg')
                .attr('width', width + margin.left + margin.right)
                .attr('height', height + margin.top + margin.bottom)
                .append('g')
                .attr('transform', `translate(${margin.left},${margin.top})`);

            // X scale (time)
            const x = d3.scaleBand()
                .domain(this.data.map(d => d.time))
                .range([0, width])
                .padding(0.1); // spacing between bars

            // Y scale (value)
            const y = d3.scaleLinear()
                .domain([0, 100])//d3.max(this.data, d => d.value)])
                .nice()
                .range([height, 0]);

            // X Axis
            svg.append('g')
                .attr('transform', `translate(0,${height})`)
                .call(
                    d3.axisBottom(x)
                        .tickFormat(d3.timeFormat('%m-%d %H:%M')) // format timestamps
                        .tickValues(this.data.map(d => d.time))
                )
                .selectAll("text") // rotate x labels for readability
                .attr("transform", "rotate(-10)")
                .style("text-anchor", "end");

            // Y Axis
            svg.append('g')
                .call(d3.axisLeft(y));

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
                .attr('fill', 'steelblue');
        },
    },
};
</script>
