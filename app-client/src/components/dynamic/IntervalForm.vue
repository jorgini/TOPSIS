<script setup>
const int = defineModel();
</script>

<script>
  export default {
    emits: ['corr-rate', 'incorr-rate'],
    data() {
      return {
        isValidStart: true,
        isValidEnd: true,
      }
    },
    methods: {
      normalize() {
        if (this.modelValue.end < this.modelValue.start) {
          this.modelValue.end = this.modelValue.start;
          this.isValidEnd = true;
        }
      },
      validateStart() {
        this.isValidStart = (this.modelValue.start >= 0 && this.modelValue.start < 2_147_483_648);
        if (this.isValidStart) {
          this.normalize();
          this.$emit('corr-rate');
        } else {
          this.$emit('incorr-rate');
        }
      },
      validateEnd() {
        this.isValidEnd = (this.modelValue.end >= 0 && this.modelValue.end < 2_147_483_648 && (this.modelValue.end >= this.modelValue.start));
        if (this.isValidStart && this.isValidEnd) {
          this.$emit('corr-rate');
        } else {
          this.$emit('incorr-rate');
        }
      }
    }
  }
</script>

<template>
  <div class="interval">
    <p>a:</p>
    <input type="number" :class="{field: true, invalid: !isValidStart}" name="start"
           placeholder="0.0" maxlength="10" v-model="int.start" @input="validateStart">
    <p>b:</p>
    <input type="number" :class="{field: true, invalid: !isValidEnd}" name="end"
           placeholder="0.0" maxlength="10" v-model="int.end" @input="validateEnd">
  </div>
</template>

<style>
@import "../../style.css";

.interval > p {
  font-family: "Inria Sans", sans-serif;
  font-weight: 700;
  font-size: 2vmin;
}

.interval {
  display: flex;
  width: fit-content;
  align-items: center;
}

.interval > * {
  display: inline-block;
  margin: auto 1.5vmin auto 0;
}

.interval > .field {
  width: 10vmin;
  height: 4vmin;
  margin: 0 1.5vmin 0 0;
}
</style>