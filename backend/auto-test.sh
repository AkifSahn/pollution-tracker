#!/bin/bash

usage(){
    programName=$(basename "$0")
    echo "Usage: $programName [--anomaly-chance=<0-100>] [--duration=<seconds>] [--rate=<number>] [--lat-min=] [--lat-max=] [--long-min=] [--long-max=]"
}

# Default values
duration=30
rate=5
chance=20

latMin=-90
latMax=90
longMin=-180
longMax=180

url="http://localhost:3000/api/pollutions"

# List of pollutant types
pollutants=("PM2.5" "PM10" "NO2" "SO2" "O3")

# Parse arguments
for arg in "$@"; do
    case $arg in
        -h|--help)
            usage
            exit 0
            ;;
        --duration=*)
            duration="${arg#*=}" ;;
        --rate=*)
            rate="${arg#*=}" ;;
        --anomaly-chance=*)
            chance="${arg#*=}" ;;
        --lat-min=*)
            latMin="${arg#*=}";;
        --lat-max=*)
            latMax="${arg#*=}";;
        --long-min=*)
            longMin="${arg#*=}";;
        --long-max=*)
            longMax="${arg#*=}";;
        *)
            echo "Unknown argument: $arg"
            usage
            exit 1
            ;;
    esac
done

# Generate random float between min and max
rand_float() {
    min=$1
    max=$2
    scale=${3:-6}
    echo "scale=$scale; $min + ($max - $min) * $(od -An -N4 -tu4 < /dev/urandom | tr -d ' \n') / 4294967295" | bc
}

# Random int 0-100
rand_int() {
    echo $(( RANDOM % 101 ))
}

# Choose random pollutant
rand_choice() {
    size=${#pollutants[@]}
    index=$(( RANDOM % size ))
    echo "${pollutants[$index]}"
}

generate_payload() {
    lat=$(rand_float $latMin $latMax 3)
    lon=$(rand_float $longMin $longMax 3)
    timestamp=$(date -Iseconds)
    pollutant=$(rand_choice)

    # Thresholds per pollutants
    # Also used in the backend sytem to detect anomalies
    case $pollutant in
        "PM2.5") normal_min=5; normal_max=40; anomaly_min=150; anomaly_max=250 ;;
        "PM10")  normal_min=10; normal_max=50; anomaly_min=180; anomaly_max=300 ;;
        "NO2")   normal_min=10; normal_max=60; anomaly_min=150; anomaly_max=250 ;;
        "SO2")   normal_min=5; normal_max=30; anomaly_min=100; anomaly_max=200 ;;
        "O3")    normal_min=20; normal_max=100; anomaly_min=200; anomaly_max=300 ;;
    esac

    roll=$(rand_int)
    if (( roll < chance )); then
        value=$(rand_float $anomaly_min $anomaly_max 2)
    else
        value=$(rand_float $normal_min $normal_max 2)
    fi

    cat <<EOF
{
  "time": "$timestamp",
  "latitude": $lat,
  "longitude": $lon,
  "value": $value,
  "pollutant": "$pollutant"
}
EOF
}

echo "Running for $duration seconds with $rate requests/sec and $chance% anomaly chance..."

end_time=$((SECONDS + duration))

while [ $SECONDS -lt $end_time ]; do
    for ((i=0; i<rate; i++)); do
        payload=$(generate_payload)
        echo "Sending payload:"
        echo "$payload"
        curl -s -X POST -H "Content-Type: application/json" -d "$payload" "$url" &
    done
    wait
    sleep 1
done

echo "Finished."

