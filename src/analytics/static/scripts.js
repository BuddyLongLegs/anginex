const TIME_UNITS = ["seconds", "minutes", "hours"];
const MAX_TIME_NUMBER = [60, 60, 24];
const TIME_MULTIPLIERS = [1 * 1000, 60 * 1000, 3600 * 1000];

const DEFAULT_TIME_STEP = 1;
const DEFAULT_TIME_STEP_UNIT = 1;
const DEFAULT_TIME_RANGE = 15;
const DEFAULT_TIME_RANGE_UNIT = 1;

let timeStep = DEFAULT_TIME_STEP;
let timeStepUnit = DEFAULT_TIME_STEP_UNIT;
let timeRange = DEFAULT_TIME_RANGE;
let timeRangeUnit = DEFAULT_TIME_RANGE_UNIT;

/*
  Utility functions
*/
function getQueryParams() {
  return `?q=${timeStep}+${TIME_UNITS[timeStepUnit]}&q=${timeRange}+${TIME_UNITS[timeRangeUnit]}`;
}

function updateChart(containerId, options, data) {
  if (data) {
    document.getElementById(containerId).classList.remove("loading");
  }
  if (!data.length) {
    document.getElementById(containerId).classList.add("no-data");
    return;
  } else {
    document.getElementById(containerId).classList.remove("no-data");
  }

  const chartData = google.visualization.arrayToDataTable(data);

  const chart = new google.visualization.LineChart(
    document.getElementById(containerId)
  );
  chart.draw(chartData, options);
}

function updateRankData(containerId, data) {
  if (data) {
    document.getElementById(containerId).classList.remove("loading");
  }

  const cont = document.getElementById(containerId);
  cont.innerHTML = '<div class="no-data-content"> No Data </div>';

  if (!data.length) {
    document.getElementById(containerId).classList.add("no-data");
    return;
  } else {
    document.getElementById(containerId).classList.remove("no-data");
  }

  const maxCount = Math.max(...data.map((d) => d.call_count));

  for (let i = 0; i < data.length; i++) {
    const d = data[i];
    const rank = document.createElement("div");
    rank.classList.add("rank-entry");
    rank.innerHTML = `
<div class="rank-main">
  <div
    class="rank-title"
    title="${d.host + d.path}"
  >
    <span class="method ${d.method}">${d.method}</span>
    ${d.host + d.path}
  </div>
  <div class="rank-count">${d.call_count}</div>
</div>
<div class="bar" style="width: ${(100 * d.call_count) / maxCount}%"></div>
    `;
    cont.appendChild(rank);
  }
}

/*
  RPM Chart
*/
let rpmData = null;

function updateRPMChart() {
  const options = {
    hAxis: { title: "Time", titleTextStyle: { color: "#333" } },
    vAxis: { minValue: 0 },
    legend: {
      position: "top",
    },
    timeline: {
      groupByRowLabel: true,
    },
    width: "100%",
    height: "100%",
  };
  updateChart("rpm-chart", options, rpmData);
}

async function getRPM() {
  let response = await fetch("/api/requests_per_unit_time" + getQueryParams());
  let data = await response.json();

  if (!data) {
    rpmData = [];
    updateRPMChart();
    return;
  }

  let maxDataTime = new Date();
  let minDataTime = new Date(
    new Date().setSeconds(0, 0) - timeRange * TIME_MULTIPLIERS[timeRangeUnit]
  );

  const timePoints = [];
  const timeStepInMilliSec = TIME_MULTIPLIERS[timeStepUnit] * timeStep;
  let currentTime = minDataTime;
  while (currentTime <= maxDataTime) {
    timePoints.push(currentTime);
    currentTime = new Date(currentTime.getTime() + timeStepInMilliSec);
  }

  const formattedData = Object.fromEntries(
    data.map((d) => [new Date(d.minute), d])
  );

  const chartData = [["Time", "Requests"]];
  for (let timePoint of timePoints) {
    chartData.push([timePoint, formattedData[timePoint]?.call_count || 0]);
  }

  rpmData = chartData;
  updateRPMChart();
}

/*
  AVG and MAX Latency Chart
*/

let latencyData = null;

function updateLatencyChart() {
  const options = {
    hAxis: { title: "Time", titleTextStyle: { color: "#333" } },
    vAxis: { minValue: 0, title: "Latency (ms)" },
    legend: {
      position: "top",
    },
    timeline: {
      groupByRowLabel: true,
    },
    width: "100%",
    height: "100%",
  };
  updateChart("latency-chart", options, latencyData);
}

async function getLatency() {
  let response = await fetch(
    "/api/avg_and_max_latency_per_unit_time" + getQueryParams()
  );
  let data = await response.json();

  if (!data) {
    latencyData = [];
    updateLatencyChart();
    return;
  }
  let maxDataTime = new Date();
  let minDataTime = new Date(
    new Date().setSeconds(0, 0) - timeRange * TIME_MULTIPLIERS[timeRangeUnit]
  );

  const timePoints = [];
  const timeStepInMilliSec = TIME_MULTIPLIERS[timeStepUnit] * timeStep;
  let currentTime = minDataTime;
  while (currentTime <= maxDataTime) {
    timePoints.push(currentTime);
    currentTime = new Date(currentTime.getTime() + timeStepInMilliSec);
  }

  const formattedData = Object.fromEntries(
    data.map((d) => [new Date(d.minute), d])
  );

  const chartData = [["Time", "AVG Latency", "MAX Latency"]];
  for (let timePoint of timePoints) {
    chartData.push([
      timePoint,
      formattedData[timePoint]?.avg_latency || -1,
      formattedData[timePoint]?.max_latency || -1,
    ]);
  }

  latencyData = chartData;
  updateLatencyChart();
}

/*
  System Metrics Chart
*/
let metricsData = null;

function updateMetricsChart() {
  const options = {
    hAxis: { title: "Time", titleTextStyle: { color: "#333" } },
    vAxis: { minValue: 0, title: "Usage%" },
    legend: {
      position: "top",
    },
    timeline: {
      groupByRowLabel: true,
    },
    width: "100%",
    height: "100%",
  };
  updateChart("system-metric-chart", options, metricsData);
}

async function getSystemMetrics() {
  let response = await fetch(
    "/api/system_usage_per_unit_time" + getQueryParams()
  );
  let data = await response.json();

  if (!data) {
    metricsData = [];
    updateMetricsChart();
    return;
  }
  let maxDataTime = new Date();
  let minDataTime = new Date(
    new Date().setSeconds(0, 0) - timeRange * TIME_MULTIPLIERS[timeRangeUnit]
  );

  const timePoints = [];
  const timeStepInMilliSec = TIME_MULTIPLIERS[timeStepUnit] * timeStep;
  let currentTime = minDataTime;
  while (currentTime <= maxDataTime) {
    timePoints.push(currentTime);
    currentTime = new Date(currentTime.getTime() + timeStepInMilliSec);
  }

  const formattedData = Object.fromEntries(
    data.map((d) => [new Date(d.minute), d])
  );

  const chartData = [["Time", "Memory Usage%", "CPU Usage%"]];
  for (let timePoint of timePoints) {
    chartData.push([
      timePoint,
      formattedData[timePoint]?.avg_memory || 0,
      formattedData[timePoint]?.avg_cpu || 0,
    ]);
  }

  metricsData = chartData;
  updateMetricsChart();
}

/*
  Most Hit Endpoints
*/
let mostHitData = null;
function updateMostHitRankChart() {
  updateRankData("most-hit-rank-chart", mostHitData);
}

async function getMostHit() {
  let response = await fetch("/api/most_hit_endpoints" + getQueryParams());
  let data = await response.json();

  if (!data) {
    mostHitData = [];
    updateMostHitRankChart();
    return;
  }

  mostHitData = data;
  updateMostHitRankChart();
}

/*
  Most Server Errored Endpoints
*/
let mostServerErrData = null;
function updateServerErrRankChart() {
  updateRankData("most-server-err-rank-chart", mostServerErrData);
}

async function getMostServerErr() {
  let response = await fetch(
    "/api/most_server_errored_endpoints_in_time_range" + getQueryParams()
  );
  let data = await response.json();

  if (!data) {
    mostServerErrData = [];
    updateServerErrRankChart();
    return;
  }

  mostServerErrData = data;
  updateServerErrRankChart();
}

/*
  Most User Errored Endpoints
*/
let mostUserErrData = null;
function updateUserErrRankChart() {
  updateRankData("most-user-err-rank-chart", mostUserErrData);
}

async function getMostUserErr() {
  let response = await fetch(
    "/api/most_user_errored_endpoints_in_time_range" + getQueryParams()
  );
  let data = await response.json();

  if (!data) {
    mostUserErrData = [];
    updateUserErrRankChart();
    return;
  }

  mostUserErrData = data;
  updateUserErrRankChart();
}

/*
  Most Success Endpoints
*/
let mostSuccessData = null;
function updateSuccessRankChart() {
  updateRankData("most-success-rank-chart", mostSuccessData);
}

async function getMostSuccess() {
  let response = await fetch(
    "/api/most_successful_endpoints_in_time_range" + getQueryParams()
  );
  let data = await response.json();

  if (!data) {
    mostSuccessData = [];
    updateSuccessRankChart();
    return;
  }

  mostSuccessData = data;
  updateSuccessRankChart();
}

/*
  Polling and overall setup
*/
function fetchAllCharts() {
  getRPM();
  getLatency();
  getSystemMetrics();
  getMostHit();
  getMostServerErr();
  getMostUserErr();
  getMostSuccess();
}

function updateAllCharts() {
  updateRPMChart();
  updateLatencyChart();
  updateMetricsChart();
}

google.charts.load("current", { packages: ["corechart"] });
google.charts.setOnLoadCallback(fetchAllCharts);

window.addEventListener("resize", updateAllCharts);

// Filter Inputs
const timeStepInput = document.getElementById("time-step");
const timeStepUnitInput = document.getElementById("time-step-unit");
const timeRangeInput = document.getElementById("time-range");
const timeRangeUnitInput = document.getElementById("time-range-unit");
const refetchButton = document.getElementById("refetch");

timeStepUnitInput.innerHTML = TIME_UNITS.map(
  (unit, i) => `<option value="${i}">${unit}</option>`
).join("");
timeRangeUnitInput.innerHTML = TIME_UNITS.map(
  (unit, i) => `<option value="${i}">${unit}</option>`
).join("");

timeStepInput.value = DEFAULT_TIME_STEP;
timeStepUnitInput.value = DEFAULT_TIME_STEP_UNIT;
timeRangeInput.value = DEFAULT_TIME_RANGE;
timeRangeUnitInput.value = DEFAULT_TIME_RANGE_UNIT;

timeStepInput.addEventListener("change", (e) => {
  timeStep = parseInt(e.target.value);
});

timeStepUnitInput.addEventListener("change", (e) => {
  timeStepUnit = parseInt(e.target.value);
  timeStepInput.max = MAX_TIME_NUMBER[timeStepUnit];
  if (timeStep > MAX_TIME_NUMBER[timeStepUnit]) {
    timeStep = MAX_TIME_NUMBER[timeStepUnit];
    timeStepInput.value = timeStep;
  }
});

timeRangeInput.addEventListener("change", (e) => {
  timeRange = parseInt(e.target.value);
});

timeRangeUnitInput.addEventListener("change", (e) => {
  timeRangeUnit = parseInt(e.target.value);
});

timeRangeUnitInput.addEventListener("change", (e) => {
  timeRangeUnit = parseInt(e.target.value);
  timeRangeInput.max = MAX_TIME_NUMBER[timeRangeUnit];
  if (timeRange > MAX_TIME_NUMBER[timeRangeUnit]) {
    timeRange = MAX_TIME_NUMBER[timeRangeUnit];
    timeRangeInput.value = timeRange;
  }
});

refetchButton.addEventListener("click", fetchAllCharts);
