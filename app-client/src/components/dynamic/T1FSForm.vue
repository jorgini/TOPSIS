<script setup>
const t1fs = defineModel();
</script>

<script>
  export default {
    emits: ['corr-rate', 'incorr-rate'],
    data() {
      return {
        isValid: [true, true, true, true],
        selectedFig: null,
      }
    },
    methods: {
      changeFigure() {
        if (this.selectedFig === 'Треугольник' && this.modelValue.vert.length === 4) {
          this.modelValue.vert.pop();
          this.isValid[3] = true;
          this.validateNumber(0);
        } else if (this.selectedFig === 'Трапеция' && this.modelValue.vert.length === 3) {
          this.modelValue.vert.push(this.modelValue.vert[2]);
          this.isValid[3] = true;
        }
      },
      normalize(i) {
        for (let k = i + 1; k < this.modelValue.vert.length; ++k) {
          if (this.modelValue.vert[i] > this.modelValue.vert[k]) {
            this.modelValue.vert[k] = this.modelValue.vert[i];
            this.isValid[k] = true;
          }
        }
      },
      validateNumber(i) {
        this.isValid[i] = this.modelValue.vert[i] >= 0 && this.modelValue.vert[i] < 2_147_483_648;
        this.isValid[i] &= this.modelValue.vert[i] !== '';
        for (let k = 0; k < i; ++k) {
          this.isValid[i] &= (this.modelValue.vert[k] <= this.modelValue.vert[i]);
        }

        if (this.isValid[i])
          this.normalize(i);

        if (this.isValid[0] && this.isValid[1] && this.isValid[2] && this.isValid[3]) {
          this.$emit('corr-rate');
        } else {
          this.$emit('incorr-rate');
        }
      }
    },
    mounted() {
      if (this.modelValue.vert.length === 3) {
        this.selectedFig = 'Треугольник';
      } else {
        this.selectedFig = 'Трапеция';
      }

      for (let k = 0; k < this.modelValue.vert.length; k++)
        this.validateNumber(k);
    }
  }
</script>

<template>
  <div class="t1fs">
    <select @change="changeFigure" v-model="selectedFig">
      <option>Треугольник</option>
      <option>Трапеция</option>
    </select>
    <p>a:</p>
    <input type="number" :class="{field: true, invalid: !isValid[0]}" name="vert"
           placeholder="0.0" maxlength="10" v-model="t1fs.vert[0]" @input="validateNumber(0)">
    <p>b:</p>
    <input type="number" :class="{field: true, invalid: !isValid[1]}" name="vert"
           placeholder="0.0" maxlength="10" v-model="t1fs.vert[1]" @input="validateNumber(1)">
    <p>c:</p>
    <input type="number" :class="{field: true, invalid: !isValid[2]}" name="vert"
           placeholder="0.0" maxlength="10" v-model="t1fs.vert[2]" @input="validateNumber(2)">
    <p v-if="t1fs.vert.length > 3">d:</p>
    <input v-if="t1fs.vert.length > 3" type="number" :class="{field: true, invalid: !isValid[3]}" name="vert"
           placeholder="0.0" maxlength="10" v-model="t1fs.vert[3]" @input="validateNumber(3)">
  </div>
</template>

<style scoped>
@import "../../style.css";

.t1fs > p {
  font-family: "Inria Sans", sans-serif;
  font-weight: 700;
  font-size: 2vmin;
}

.t1fs {
  display: flex;
  width: fit-content;
  align-items: center;
}

.t1fs > * {
  display: inline-block;
  margin: auto 1.5vmin auto 0;
}

.t1fs > .field {
  width: 8vmin;
  height: 4vmin;
  margin: 0 1.5vmin 0 0;
}

select {
  background-color: #ABF8F4;
  font: inherit;
  font-family: "Inria Sans", sans-serif;
  font-size: 1.8vmin;
  font-weight: 700;
  color: black;
  border-radius: 1em;
  border: 1px solid black;
  width: min-content;
  margin-top: 1.5vmin;
  box-shadow: 0 4px 4px 0 rgba(0, 0, 0, 0.25);
}
</style>