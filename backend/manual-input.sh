#!/bin/bash

ENDPOINT="http://localhost:3000/api/ingest/manual"

if ! command -v curl &> /dev/null; then
    echo "curl is not installed!"
    exit 1
fi

if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <latitude> <longitude> <parameter> <value>"
    exit 1
fi

latitude=$1
longitude=$2
parameter=$3
value=$4

BODY="{\"latitude\":$latitude, \"longitude\":$longitude, \"value\": $value, \"pollutant\": \"$parameter\"}"

echo "Body: $BODY"

STATUS=$(curl -s -o response.json -w "%{http_code}" \
    -X POST \
    -H "Content-Type: application/json" \
    -d "$BODY" \
    "$ENDPOINT")

echo "Response: $STATUS"

rm -f response.json

