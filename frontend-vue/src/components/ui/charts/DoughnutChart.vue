<script>
import { Doughnut } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, ArcElement)

export default {
  name: 'DoughnutChart',
  components: { Doughnut },
  props: {
    chartData: {
      type: Object,
      required: true
    },
    chartOptions: {
      type: Object,
      default: () => ({
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { 
            position: 'right',
            labels: { color: '#9ca3af', boxWidth: 12 }
          }
        },
        cutout: '60%'
      })
    }
  }
}
</script>

<template>
  <div class="chart-container">
    <Doughnut :data="chartData" :options="chartOptions" />
  </div>
</template>

<style scoped>
.chart-container {
  position: relative;
  height: 100%;
  width: 100%;
}
</style>
