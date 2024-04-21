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
          new Array.from({ length: this.alts.length }, (_, i) => i + 1).reverse() :
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
          labels: this.results.result.order.map(ord => this.alts[ord].title),
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
      })

      const ctx2 = document.getElementById('myRadar');

      const labels = Array.from({ length: this.results.sens_analysis.Results.length }, (_, i) => `sens-${i + 1}`);
      const altsData = this.alts
          .map(alt => {
            return {
              label: alt.title,
              data: null,
              backgroundColor: null,
              fill: null
            }
          });

      console.log(this.results.sens_analysis)
      for (let j = 0; j < this.alts.length; j++) {
        altsData[j].data = this.results.sens_analysis.Results
            .map(stats => {
              let rank = stats.order.indexOf(j);
              if (this.type === 'smart') {
                return this.alts.length - rank;
              } else {
                if (typeof stats.coeffs[rank] === 'number') {
                  return stats.coeffs[rank];
                } else {
                  return (stats.coeffs[rank].start + stats.coeffs[rank].end) / 2;
                }
              }
            });

        const fill = palette[this.results.result.order.indexOf(j)].rgba();
        altsData[j].backgroundColor = `rgba(${fill[0]}, ${fill[1]}, ${fill[2]}, 0.7)`;
        altsData[j].fill = {
          target: "origin",
        };
      }

      new Chart(ctx2, {
        type: 'radar',
        data: {
          labels: labels,
          datasets: altsData,
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