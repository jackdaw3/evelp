<template>
  <div class="Stock">
    <highcharts class="stock" :constructor-type="'stockChart'" :options="chartOptions"></highcharts>
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
          backgroundColor: "#FFFFFF",
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
        },
        rangeSelector: {
          allButtonsEnabled: true,
          selected: 0,
        },
        title: {
          text: this.history.title,
        },
        yAxis: [
          {
            startOnTick: false,
            endOnTick: false,
            labels: {
              align: "right",
              x: -3,
            },
            title: {
              text: this.history.label.price,
            },
            height: "60%",
            lineWidth: 2,
            resize: {
              enabled: true,
            },
          },
          {
            labels: {
              align: "right",
              x: -3,
            },
            title: {
              text: this.history.label.volume,
            },
            top: "65%",
            height: "35%",
            offset: 0,
            lineWidth: 2,
          },
        ],
        tooltip: {
          split: false,
          shared: true,
          valueDecimals: 0,
        },
        series: [
          {
            name: this.history.label.average,
            data: this.history.average,
            yAxis: 0,
            color: "#F1C40F",
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
            color: "#3498DB",
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



