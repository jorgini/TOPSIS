<script setup>
const results = defineModel('results');
const alts = defineModel('alts');
</script>

<script>
  import Chart from "chart.js/auto";
  import distinctColors from 'distinct-colors';
  export default {
    props: {
      type: String,
    },
    mounted() {
      const ctx = document.getElementById('myChart');
      const palette = distinctColors({
        count: this.alts.length
      });

      const dataRanking = this.type === 'smart' ?
          this.results.result.order.map(ord => (ord + 1)) :
          this.results.result.coeffs
              .map(score => {
                if (typeof score === 'number')
                  return score
                else {
                  return (score.start + score.end) / 2
                }
              });

      new Chart(ctx, {
        type: 'bar',
        data: {
          labels: this.alts
              .map(alt => {
                return alt.title
              }),
          datasets: [{
            label: 'Alts Ranking',
            data: dataRanking,
            backgroundColor: palette,
            borderWidth: 1
          }]
        },
        options: {
          indexAxis: 'y',
          scales: {
            y: {
              beginAtZero: true
            }
          }
        }
      });

      const ctx2 = document.getElementById('myRadar');

      const dataSets = this.results.sens_analysis.Results
              .map(res => {
                if (this.type === 'smart') {
                  return res.order.map(ord => (ord + 1))
                } else {
                  return res.coeffs.map(score => {
                    if (typeof score === 'number') {
                      return score
                    } else {
                      return (score.start + score.end) / 2;
                    }
                  });
                }
              });

      let j = 0;
      const labels = Array.from({ length: this.results.sens_analysis.Results.length }, (_, i) => `sens-${i + 1}`);
      console.log(labels);
      new Chart(ctx2, {
        type: 'radar',
        data: {
          labels: labels,
          datasets: this.alts
              .map(alt => {
                const res = {
                  label: alt.title,
                  data: dataSets
                      .map(rank => rank[j]),
                  backgroundColor: palette[j],
                  fill: -1
                };
                j++;
                return res
              }),
        },
        options: {
          plugins: {
            filler: {
              propagate: false
            },
            'samples-filler-analyser': {
              target: 'chart-analyser'
            }
          },
          interaction: {
            intersect: false
          }
        }
      });
    }
  }
</script>

<template>
  <div class="info">
    <canvas id="myChart" style="width: 70%"></canvas>
    <canvas id="myRadar" style="width: 70%"></canvas>
  </div>
</template>

<style scoped>
  @import "../../style.css";

  .info {
    display: block;
  }

  .info > * {
    margin: 2vmin auto 2vmin auto;
  }

</style>