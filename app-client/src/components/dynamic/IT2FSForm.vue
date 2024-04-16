<script setup>
const it2fs = defineModel();
</script>

<script>
export default {
  emits: ['corr-rate', 'incorr-rate'],
  data() {
    return {
      isValidBottom: [true, true, true, true],
      isValidUpward: [true, true],
      selectedFig: null,
    }
  },
  methods: {
    changeFigure() {
      if (this.selectedFig === 'Треугольник' && this.modelValue.upward.length === 2) {
        this.modelValue.upward.pop();
        this.isValidUpward[1] = true;
        this.validateUpward(0);
      } else if (this.selectedFig === 'Трапеция' && this.modelValue.upward.length === 1) {
        this.modelValue.upward.push(this.modelValue.upward[0]);
        this.isValidUpward[2] = true;
      }
    },
    normalize(i) {
      const verts = [this.modelValue.bottom[0].start, this.modelValue.bottom[0].end, this.modelValue.upward[0]];
      if (this.modelValue.upward.length > 1)
        verts.push(this.modelValue.upward[1]);
      verts.push(this.modelValue.bottom[1].start);
      verts.push(this.modelValue.bottom[1].end);

      const x = verts[i];
      if (i <= 1 && x > verts[1]) {
        this.modelValue.bottom[0].end = x;
        this.isValidBottom[1] = true;
      }

      if (i <= 2 && x > verts[2]) {
        this.modelValue.upward[0] = x;
        this.isValidUpward[0] = true;
      }

      if (i <= 3 && x > verts[3]) {
        if (this.modelValue.upward.length > 1) {
          this.modelValue.upward[1] = x;
          this.isValidUpward[1] = true;
        } else {
          this.modelValue.bottom[1].start = x;
          this.isValidBottom[2] = true;
        }
      }

      if (i <= 4 && x > verts[4]) {
        if (this.modelValue.upward.length > 1) {
          this.modelValue.bottom[1].start = x;
          this.isValidUpward[2] = true;
        } else {
          this.modelValue.bottom[1].end = x;
          this.isValidBottom[3] = true;
        }
      }

      if (i <= 5 && verts.length === 6 && x > verts[5]) {
        this.modelValue.bottom[1].end = x;
        this.isValidBottom[3] = true;
      }
    },
    validateBottom(i) {
      const x = i % 2 === 0 ? this.modelValue.bottom[Math.floor(i / 2)].start : this.modelValue.bottom[Math.floor(i / 2)].end;
      this.isValidBottom[i] = x >= 0 && x < 2_147_483_648;
      this.isValidBottom[i] &= x !== '';
      if (i > 1) {
        this.isValidBottom[i] &= this.modelValue.bottom[1].start <= x &&
            this.modelValue.bottom[0].end <= x && this.modelValue.bottom[0].start <= x;

        for (let k = 0; k < this.modelValue.upward.length; ++k) {
          this.isValidBottom[i] &= this.modelValue.upward[k] <= x;
        }
      } else {
        this.isValidBottom[i] &= this.modelValue.bottom[0].start <= x;
      }

      if (this.isValidBottom[i])
        this.normalize((i > 1 ? i + this.modelValue.upward.length : i));

      if (this.isValidBottom[0] && this.isValidBottom[1] && this.isValidBottom[2] && this.isValidBottom[3]
          && this.isValidUpward[0] && this.isValidUpward[1]) {
        this.$emit('corr-rate');
      } else {
        this.$emit('incorr-rate');
      }
    },
    validateUpward(i) {
      const x = this.modelValue.upward[i];
      this.isValidUpward[i] = x >= 0 && x <= 2_147_483_648;
      this.isValidUpward[i] &= x !== '';
      this.isValidUpward[i] &= this.modelValue.upward[0] <= x && this.modelValue.bottom[0].end <= x && this.modelValue.bottom[0].start <= x;

      if (this.isValidUpward[i])
        this.normalize(i + 2);

      if (this.isValidBottom[0] && this.isValidBottom[1] && this.isValidBottom[2] && this.isValidBottom[3]
          && this.isValidUpward[0] && this.isValidUpward[1]) {
        this.$emit('corr-rate');
      } else {
        this.$emit('incorr-rate');
      }
    }
  },
  mounted() {
    if (this.modelValue.upward.length === 1) {
      this.selectedFig = 'Треугольник';
    } else {
      this.selectedFig = 'Трапеция';
    }
  }
}
</script>

<template>
  <div class="it2fs">
    <select @change="changeFigure" v-model="selectedFig">
      <option>Треугольник</option>
      <option>Трапеция</option>
    </select>
    <p>_a:</p>
    <input type="number" :class="{field: true, invalid: !isValidBottom[0]}" name="vert"
           placeholder="0.0" maxlength="10" v-model="it2fs.bottom[0].start" @input="validateBottom(0)">

    <p>^a:</p>
    <input type="number" :class="{field: true, invalid: !isValidBottom[1]}" name="vert"
           placeholder="0.0" maxlength="10" v-model="it2fs.bottom[0].end" @input="validateBottom(1)">

    <p>b:</p>
    <input type="number" :class="{field: true, invalid: !isValidUpward[0]}" name="vert"
           placeholder="0.0" maxlength="10" v-model="it2fs.upward[0]" @input="validateUpward(0)">

    <p v-if="it2fs.upward.length > 1">c:</p>
    <input v-if="it2fs.upward.length > 1" type="number" :class="{field: true, invalid: !isValidUpward[1]}" name="vert"
           placeholder="0.0" maxlength="10" v-model="it2fs.upward[1]" @input="validateUpward(1)">

    <p>{{ (it2fs.upward.length === 1) ? "_c" : "_d" }}:</p>
    <input type="number" :class="{field: true, invalid: !isValidBottom[2]}" name="vert"
           placeholder="0.0" maxlength="10" v-model="it2fs.bottom[1].start" @input="validateBottom(2)">

    <p>{{ (it2fs.upward.length === 1) ? "^c" : "^d" }}:</p>
    <input type="number" :class="{field: true, invalid: !isValidBottom[3]}" name="vert"
           placeholder="0.0" maxlength="10" v-model="it2fs.bottom[1].end" @input="validateBottom(3)">
  </div>
</template>

<style scoped>
@import "../../style.css";

.it2fs p {
  font-family: "Inria Sans", sans-serif;
  font-weight: 700;
  font-size: 2vmin;
}

.it2fs {
  display: flex;
  width: fit-content;
  align-items: center;
}

.it2fs > * {
  display: inline-block;
  margin: auto 1.5vmin auto 0;
}

.it2fs > .field {
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