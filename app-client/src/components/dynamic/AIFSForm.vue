<script setup>
const aifs = defineModel();
const role = defineModel('role');
const emits = defineEmits(["corr-rate", "incorr-rate"]);
</script>

<script>
export default {
  emits: ['corr-rate', 'incorr-rate'],
  data() {
    return {
      isValidVert: [true, true, true, true],
      isValidPi: true,
      selectedFig: null,
    }
  },
  methods: {
    changeFigure() {
      if (this.selectedFig === 'Треугольник' && this.modelValue.vert.length === 4) {
        this.modelValue.vert.pop();
        this.isValidVert[3] = true;
        this.validateVert(0);
      } else if (this.selectedFig === 'Трапеция' && this.modelValue.vert.length === 3) {
        this.modelValue.vert.push(this.modelValue.vert[2]);
        this.isValidVert[3] = true;
      }
    },
    normalize(i) {
      for (let k = i + 1; k < this.modelValue.vert.length; ++k) {
        if (this.modelValue.vert[i] > this.modelValue.vert[k]) {
          this.modelValue.vert[k] = this.modelValue.vert[i];
          this.isValidVert[k] = true;
        }
      }
    },
    validateVert(i) {
      this.isValidVert[i] = this.modelValue.vert[i] >= 0 && this.modelValue.vert[i] < 2_147_483_648;
      this.isValidVert[i] &= this.modelValue.vert[i] !== '';
      for (let k = 0; k < i; ++k) {
        this.isValidVert[i] &= (this.modelValue.vert[k] <= this.modelValue.vert[i]);
      }

      if (this.isValidVert[i])
        this.normalize(i);

      if (this.isValidVert[0] && this.isValidVert[1] && this.isValidVert[2] && this.isValidVert[3] && this.isValidPi) {
        this.$emit('corr-rate');
      } else {
        this.$emit('incorr-rate');
      }
    },
    validatePi() {
      this.isValidPi = this.modelValue.pi >= 0 && this.modelValue.pi <= 1;
      this.isValidPi &= this.modelValue.pi !== '';
      if (this.isValidVert[0] && this.isValidVert[1] && this.isValidVert[2] && this.isValidVert[3] && this.isValidPi) {
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
      this.validateVert(k);
    this.validatePi();
  }
}
</script>

<template>
  <div class="aifs">
    <select @change="changeFigure" v-model="selectedFig" :disabled="role==='expert'">
      <option>Треугольник</option>
      <option>Трапеция</option>
    </select>
    <p>a:</p>
    <input type="number" :class="{field: true, invalid: !isValidVert[0]}" name="vert" :readonly="role==='expert'"
           placeholder="0.0" maxlength="10" v-model="aifs.vert[0]" @input="validateVert(0)">
    <p>b:</p>
    <input type="number" :class="{field: true, invalid: !isValidVert[1]}" name="vert" :readonly="role==='expert'"
           placeholder="0.0" maxlength="10" v-model="aifs.vert[1]" @input="validateVert(1)">
    <p>c:</p>
    <input type="number" :class="{field: true, invalid: !isValidVert[2]}" name="vert" :readonly="role==='expert'"
           placeholder="0.0" maxlength="10" v-model="aifs.vert[2]" @input="validateVert(2)">
    <p v-if="aifs.vert.length > 3">d:</p>
    <input v-if="aifs.vert.length > 3" type="number" :class="{field: true, invalid: !isValidVert[3]}" name="vert"
           :readonly="role==='expert'" placeholder="0.0" maxlength="10" v-model="aifs.vert[3]" @input="validateVert(3)">
    <p>pi:</p>
    <input type="number" :class="{field: true, invalid: !isValidPi}" name="pi" :readonly="role==='expert'"
           placeholder="0.0" maxlength="10" v-model="aifs.pi" @input="validatePi">
  </div>
</template>

<style scoped>
@import "../../style.css";

.aifs p {
  font-family: "Inria Sans", sans-serif;
  font-weight: 700;
  font-size: 2vmin;
}

.aifs {
  display: flex;
  width: fit-content;
  align-items: center;
}

.aifs > * {
  display: inline-block;
  margin: auto 1.5vmin auto 0;
}

.aifs > .field {
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
  width: fit-content;
  margin-top: 1.5vmin;
  box-shadow: 0 4px 4px 0 rgba(0, 0, 0, 0.25);
}
</style>