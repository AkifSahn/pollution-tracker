
<template>
    <div class="flex-row justify-center bg-red-400 gap-x-14 py-10"> 

        <div class="flex place-self-center flex-col items-center gap-4 bg-gray-500 p-5 rounded-md max-w-96">
            <div class="flex gap-12">
                <div class="rounded-sm">
                    <label class="w-14 text-center">Latitude From: {{ latitudeFrom }}</label>
                    <input v-model="latitudeFrom" class="block accent-red-400" type="range" min="-90" max="90">

                    <label class="w-14 text-center">Latitude To: {{ latitudeTo }}</label>
                    <input v-model="latitudeTo" class="block accent-red-400" type="range" min="-90" max="90">
                </div>

                <div class="rounded-sm">
                    <label class="w-14 text-center">Longitude From: {{ longitudeFrom }}</label>
                    <input v-model="longitudeFrom" class="block accent-red-400" type="range" min="-180" max="180">

                    <label class="w-14 text-center">Longitude To: {{ longitudeTo }}</label>
                    <input v-model="longitudeTo" class="block accent-red-400" type="range" min="-180" max="180">
                </div>
            </div>

            <button @click="fetchData" class="block mt-4 px-4 py-2 bg-blue-600 hover:bg-blue-700 hover:cursor-pointer text-white rounded-md" type="submit">Graph</button>
        </div>
        <div ref="chartContainer" class="w-full max-w-7xl h-96 place-self-center"></div>
    </div>
</template>


<script>

    import * as d3 from "d3";
    import { useMapStore } from "../stores/mapStore";
    import { watch } from "vue";

    export default{
        data() {
            return {
                latitudeFrom: -90,
                latitudeTo: 90,
                longitudeFrom: -180,
                longitudeTo: 180,
                data: [ ],
                mapStore: null,
            };
        },
        mounted() {
            this.mapStore = useMapStore();

            watch(
                () => this.mapStore.graphBound,
                (newBound) => {
                    this.latitudeFrom = Math.min(newBound[0].lat,newBound[1].lat);
                    this.longitudeFrom = Math.min(newBound[0].lng,newBound[1].lng);

                    this.latitudeTo = Math.max(newBound[0].lat, newBound[1].lat);
                    this.longitudeTo = Math.max(newBound[0].lng, newBound[1].lng);
                },
            )

            this.fetchData().then(() =>{
                this.drawChart();
            });
        },
        methods: {
            async fetchData(){
                // Fetch the pollution data from backend
                // Let the user decide the lat-start, lat-end and long-start, long-end
                // User can choose these values from the map and see the graph for that rectangle. 
                // Density graph for overall region for the time range also given by user

                try{
                    const url = `http://127.0.0.1:3000/api/region/density/SO2?` +
                        `latFrom=${encodeURIComponent(this.latitudeFrom)}&` +
                        `latTo=${encodeURIComponent(this.latitudeTo)}&` +
                        `longFrom=${encodeURIComponent(this.longitudeFrom)}&` +
                        `longTo=${encodeURIComponent(this.longitudeTo)}&` +
                        `from=${encodeURIComponent("2025-03-01 00:00:00")}&` +
                        `to=${encodeURIComponent("2025-04-10 23:00:00")}`;

                    console.log("Fetching:", url);

                    const response = await fetch(url);
                    const jsonData = await response.json();

                    if (jsonData.densities == null) {
                        this.data = [];
                    }else{
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
