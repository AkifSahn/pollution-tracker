// Base Url
const API_BASE_URL = 'http://127.0.0.1:3000/api';

// Get Pollutions
export async function fetchPollutions(fromDate, toDate, pollutant) {
    var url = `${API_BASE_URL}/pollutions?` +
        `from=${encodeURIComponent(fromDate)}&` +
        `to=${encodeURIComponent(toDate)}`
    if (pollutant) {
        url += `&pollutant=${encodeURIComponent(pollutant)}`;
    }

    const response = await fetch(url);
    const data = await response.json();
    return data;
}

// Get Pollutants
export async function fetchPollutants() {
    const url = `${API_BASE_URL}/pollutants`;
    const response = await fetch(url);
    const data = await response.json();
    return data;
}

// Get Anomalies
export async function fetchAnomaliesOfRange(fromDate, toDate) {
    const url = `${API_BASE_URL}/anomalies` +
        `?from=${encodeURIComponent(fromDate)}&` +
        `to=${encodeURIComponent(toDate)}`;

    const response = await fetch(url);
    const data = await response.json();
    return data;
}

// Get Region Density
export async function fetchRegionDensity(pollutant, params) {
    const { latFrom, latTo, longFrom, longTo, from, to } = params;

    const url = `${API_BASE_URL}/pollutions/density/rect` +
        `?pollutant=${encodeURIComponent(pollutant)}&` +
        `latFrom=${encodeURIComponent(latFrom)}&` +
        `latTo=${encodeURIComponent(latTo)}&` +
        `longFrom=${encodeURIComponent(longFrom)}&` +
        `longTo=${encodeURIComponent(longTo)}&` +
        `from=${encodeURIComponent(from)}&` +
        `to=${encodeURIComponent(to)}`;

    const response = await fetch(url);
    const data = await response.json();
    return data;
}
