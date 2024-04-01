<script>
  export default {
    props: {
      alts: Array,
      results: Object
    },
    data() {
      return {

      }
    },
    methods: {

    }
  }
</script>

<template>
  <div class="info">
    <div v-for="(_, i) in results.result.coeffs" class="place">
      <p>{{i + 1}}. {{ alts[results.result.order[i]].title }} - {{ results.result.coeffs[i] }}</p>
    </div>
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

  .place {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
  }

  .place p {
    font-family: "Inria Sans", sans-serif;
    font-weight: 700;
    font-size: 3vmin;
    margin: auto;
  }
</style>