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
          backgroundColor: "#202124",
          alignTicks: false,
          marginRight: this.history.borderWidth,
        },
        credits: {
          enabled: false
        },
        scrollbar: {
          enabled: false,
        },
        legend: {
          enabled: true,
          align: 'center',
          verticalAlign: 'top',
          backgroundColor: "#202124",
          layout: 'horizontal',
          itemStyle: {
            color: '#D0D3D4'
          },
        },
        navigator: {
          series: {
            data: this.history.average,
            color: "#3498DB",
          },
          xAxis: {
            labels: {
              style: {
                color: '#D0D3D4',
              }
            },
            gridLineWidth: 0,
            minorGridLineWidth: 0,
            dateTimeLabelFormats: {
              millisecond: '%Y.%m.%d',
              second: '%Y.%m.%d',
              minute: '%Y.%m.%d',
              hour: '%Y.%m.%d',
              day: '%Y.%m.%d',
              week: '%Y.%m.%d',
              month: '%Y.%m.%d',
              year: '%Y.%m.%d'
            },
          },
          yAxis: {
            gridLineWidth: 0,
            minorGridLineWidth: 0,
          },
        },
        navigation: {
          buttonOptions: {
            theme: {
              states: {
                hover: {
                  fill: '#0D4579'
                },
                select: {
                  stroke: '#0D4579',
                  fill: '#0D4579'
                }
              }
            }
          },
          menuItemHoverStyle: {
            background: '#0D4579',
            color: '#D0D3D4',
          },
          menuItemStyle: {
            color: '#D0D3D4',
          },
          menuStyle: {
            background: '#202124',
          },
        },
        exporting: {
          enabled: true,
          buttons: {
            contextButton: {
              menuItems: ['viewFullscreen'],
              symbolStroke: '#D0D3D4',
              theme: {
                fill: '#202124'
              }
            },
          },
        },
        rangeSelector: {
          inputDateFormat: '%Y.%m.%d',
          inputEditDateFormat: '%Y.%m.%d',
          allButtonsEnabled: true,
          selected: 0,
          inputEnabled: true,
          inputStyle: {
            color: '#D0D3D4'
          },
          labelStyle: {
            color: '#D0D3D4',
            fontWeight: 'bold'
          },
          buttons: [{
            type: 'month',
            count: 1,
            text: this.history.label.rangeSelector.month,
          }, {
            type: 'month',
            count: 3,
            text: this.history.label.rangeSelector.threeMonths,
          }, {
            type: 'month',
            count: 6,
            text: this.history.label.rangeSelector.halfYear,
          }, {
            type: 'ytd',
            text: this.history.label.rangeSelector.yearToDay,
          }, {
            type: 'year',
            count: 1,
            text: this.history.label.rangeSelector.year,
          }, {
            type: 'all',
            text: this.history.label.rangeSelector.all,
          }],
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
                fill: "#0D4579",
                style: {
                  color: "#ECF0F1",
                },
              },
              select: {
                fill: "#0D4579",
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
          crosshair: {
            dashStyle: 'dot',
          },
          dateTimeLabelFormats: {
            millisecond: '%Y.%m.%d',
            second: '%Y.%m.%d',
            minute: '%Y.%m.%d',
            hour: '%Y.%m.%d',
            day: '%Y.%m.%d',
            week: '%Y.%m.%d',
            month: '%Y.%m.%d',
            year: '%Y.%m.%d'
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
          xDateFormat: "%Y.%m.%d",
          backgroundColor: "rgba(0,0,0,0.8)",
          shared: true,
          valueDecimals: 0,
          style: {
            fontSize: 13,
            color: "#D0D3D4",
          },
          headerFormat: '<span style="font-size: 13px">{point.key}</span><br/>',
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
            color: "#A6ACAF",
            opacity: 0.6,
            pointWidth: 3,
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
            color: "#0F5C70",
          },
        ],
      },
    };
  },
  watch: {
    "$i18n.locale"() {
      console.log(this.history)
      this.chartOptions.rangeSelector.buttons[0].text = this.history.label.rangeSelector.month;
      this.chartOptions.rangeSelector.buttons[1].text = this.history.label.rangeSelector.threeMonths;
      this.chartOptions.rangeSelector.buttons[2].text = this.history.label.rangeSelector.halfYear;
      this.chartOptions.rangeSelector.buttons[3].text = this.history.label.rangeSelector.yearToDay;
      this.chartOptions.rangeSelector.buttons[4].text = this.history.label.rangeSelector.year;
      this.chartOptions.rangeSelector.buttons[5].text = this.history.label.rangeSelector.all;
      this.chartOptions.yAxis[0].title.text = this.history.label.price;
      this.chartOptions.yAxis[1].title.text = this.history.label.volume;
      this.chartOptions.series[0].name = this.history.label.average;
      this.chartOptions.series[1].name = this.history.label.minAndmax;
      this.chartOptions.series[2].name = this.history.label.average5d;
      this.chartOptions.series[3].name = this.history.label.average20d;
      this.chartOptions.series[4].name = this.history.label.minAndmax5d;
      this.chartOptions.series[5].name = this.history.label.volume;
    }
  },
};
</script>
<style>
.stock {
  width: 100%;
  height: 100%;
}

input.highcharts-range-selector:focus {
  background-color: #0D4579;
}
</style>
