<script setup>
  const number = defineModel();
  const role = defineModel('role');
  const emits = defineEmits(["corr-rate", "incorr-rate"])
</script>

<script>
  export default {
    emits: ['corr-rate', 'incorr-rate'],
    data() {
      return {
        isValid: true,
      }
    },
    methods: {
      validateNumber() {
        this.isValid = this.modelValue >= 0 && this.modelValue < 2_147_483_648;
        this.isValid &= this.modelValue !== '';
        if (this.isValid) {
          this.$emit('corr-rate');
        } else {
          this.$emit('incorr-rate');
        }
      }
    },
    mounted() {
      this.validateNumber();
    }
  }
</script>

<template>
  <div class="number">
    <input type="number" :class="{field: true, invalid: !isValid}" name="number" :disabled="role === 'expert'"
           placeholder="0.0" maxlength="10" v-model="number" @input="validateNumber">
  </div>
</template>

<style scoped>
  @import "../../style.css";

  p {
    font-family: "Inria Sans", sans-serif;
    font-weight: 300;
    font-size: 1.5vmin;
  }

  .number {
    display: flex;
    width: fit-content;
    padding: 0;
  }

  .field {
    width: 10vmin;
    height: 4vmin;
    margin: 0 0;
  }
</style>