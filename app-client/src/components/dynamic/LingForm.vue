<script setup>
const ling = defineModel();
const scale = defineModel('scale');
const role = defineModel('role');
const emits = defineEmits(["corr-rate", "incorr-rate"]);
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
    validate() {
      this.isValid = this.scale.marks.includes(this.modelValue.mark.trim());
      if (this.isValid) {
        this.modelValue.eval = this.scale.ratings[this.scale.marks.indexOf(this.modelValue.mark.trim())];
        this.$emit('corr-rate');
      } else
        this.$emit('incorr-rate');
    }
  },
  mounted() {
    this.validate();
  }
}
</script>

<template>
  <div class="ling">
    <p>Оценка:</p>
    <input type="text" :class="{field: true, invalid: !isValid}" name="ling" :readonly="role==='expert'"
           placeholder="mark" maxlength="10" v-model="ling.mark" @input="validate">
  </div>
</template>

<style scoped>
@import "../../style.css";

.ling > p {
  font-family: "Inria Sans", sans-serif;
  font-weight: 700;
  font-size: 2vmin;
}

.ling {
  display: flex;
  width: fit-content;
  align-items: center;
}

.ling > * {
  display: inline-block;
  margin: auto 1.5vmin auto 0;
}

.ling > .field {
  width: 20vmin;
  height: 4vmin;
  margin: 0 1.5vmin 0 0;
}
</style>