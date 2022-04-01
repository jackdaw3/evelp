<template>
  <div class="Stock">
    <highcharts
      class="stock"
      :constructor-type="'stockChart'"
      :options="chartOptions"
    ></highcharts>
  </div>
</template>

<script>
export default {
  props: {
    history: Object,
  },
  data() {
    return {
      chartOptions: {
        chart: {
          backgroundColor: "rgba(0,0,0,0)",
          alignTicks: false,
          marginRight: this.history.borderWidth,
        },
        scrollbar: {
          enabled: false,
        },
        navigator: {
          series: {
            color: "#3498DB",
          },
          xAxis: {
            gridLineWidth: 0,
            minorGridLineWidth: 0,
          },
          yAxis: {
            gridLineWidth: 0,
            minorGridLineWidth: 0,
          },
        },
        exporting: {
          enabled: false,
        },
        rangeSelector: {
          allButtonsEnabled: true,
          selected: 0,
          inputEnabled: false,
          buttonTheme: {
            fill: "none",
            stroke: "none",
            "stroke-width": 0,
            r: 8,
            style: {
              color: "#ECF0F1",
              fontWeight: "bold",
            },
            states: {
              hover: {
                fill: "#1C2937",
                style: {
                  color: "#ECF0F1",
                },
              },
              select: {
                fill: "#1C2937",
                style: {
                  color: "#ECF0F1",
                },
              },
            },
          },
        },
        xAxis: {
          gridLineWidth: 0,
          gridLineColor: "#505053",
          lineColor: "#505053",
          labels: {
            style: {
              color: "#ECF0F1",
            },
          },
        },
        yAxis: [
          {
            gridLineWidth: 0.3,
            gridLineColor: "#505053",
            minorGridLineWidth: 0,
            startOnTick: false,
            endOnTick: false,
            labels: {
              align: "right",
              x: -3,
              style: {
                color: "#ECF0F1",
              },
            },
            title: {
              text: this.history.label.price,
              style: {
                color: "#ECF0F1",
              },
            },
            height: "60%",
            lineWidth: 2,
            lineColor: "#505053",
            resize: {
              enabled: true,
            },
          },
          {
            gridLineWidth: 0.3,
            gridLineColor: "#505053",
            minorGridLineWidth: 0,
            labels: {
              align: "right",
              x: -3,
              style: {
                color: "#ECF0F1",
              },
            },
            title: {
              text: this.history.label.volume,
              style: {
                color: "#ECF0F1",
              },
            },
            top: "65%",
            height: "35%",
            offset: 0,
            lineWidth: 2,
            lineColor: "#505053",
          },
        ],
        tooltip: {
          split: false,
          backgroundColor: "rgba(0,0,0,0.8)",
          shared: true,
          valueDecimals: 0,
          style: {
            color: "#D0D3D4",
          },
        },
        series: [
          {
            name: this.history.label.average,
            data: this.history.average,
            yAxis: 0,
            color: "#D68910",
            lineWidth: 0,
            zIndex: 9,
            marker: {
              enabled: true,
              radius: 2.65,
            },
          },
          {
            name: this.history.label.minAndmax,
            data: this.history.minAndmax,
            type: "columnrange",
            yAxis: 0,
            zIndex: 8,
            color: "#3498DB",
            opacity: 0.6,
            pointWidth: 4.5,
          },
          {
            name: this.history.label.average5d,
            data: this.history.average5d,
            zIndex: 7,
            yAxis: 0,
            color: "#1D8348",
            lineWidth: 1.5,
            marker: {
              radius: 2.65,
            },
          },
          {
            name: this.history.label.average20d,
            data: this.history.average20d,
            zIndex: 7,
            yAxis: 0,
            color: "#C0392B",
            lineWidth: 1.5,
            marker: {
              radius: 1.65,
            },
          },
          {
            name: this.history.label.minAndmax5d,
            data: this.history.minAndmax5d,
            type: "arearange",
            zIndex: 6,
            yAxis: 0,
            color: "#85C1E9",
            opacity: 0.2,
          },
          {
            name: this.history.label.volume,
            type: "column",
            data: this.history.volume,
            yAxis: 1,
            color: "#003333",
          },
        ],
      },
    };
  },
};
</script>
<style>
.stock {
  width: 100%;
  height: 100%;
}
</style>



