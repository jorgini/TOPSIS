<script setup>
const lingScale = defineModel();
const role = defineModel('role');
const emits = defineEmits(["corr-scale", "incorr-scale"]);
</script>

<script>
  import NumberForm from "./NumberForm.vue";
  import IntervalForm from "./IntervalForm.vue";
  import T1FSForm from "./T1FSForm.vue";
  import AIFSForm from "./AIFSForm.vue";
  import IT2FSForm from "./IT2FSForm.vue";
  export default {
    emits: ['corr-scale', 'incorr-scale'],
    components: {IntervalForm, NumberForm, T1FSForm, AIFSForm, IT2FSForm},
    data() {
      return {
        showWarning: false,
        warning: null,
        isDefaultLing: false,
        selectedType: null,
        curType: null,
        defNum: null,
        defInt: null,
        defT1FS: null,
        isValidRatings: [],
        isValidMarks: []
      }
    },
    methods: {
      deepEqual(obj1, obj2) {
        if (obj1 === obj2) {
          return true;
        }
        if (typeof obj1 !== 'object' || typeof obj2 !== 'object' || obj1 === null || obj2 === null) {
          return false;
        }

        let keysA = Object.keys(obj1), keysB = Object.keys(obj2);
        if (keysA.length !== keysB.length) {
          return false;
        }

        for (let key of keysA) {
          if (!keysB.includes(key) || !this.deepEqual(obj1[key], obj2[key])) {
            return false;
          }
        }

        return true;
      },
      chooseLingScale(event) {
        if (event.target.textContent.trim() === 'По умолчанию') {
          if (this.defNum === null || this.defInt === null || this.defT1FS === null) {
            return
          }

          if (this.selectedType === 'T1FS') {
            for (let k in this.defT1FS.marks) {
              this.modelValue.marks[k] = this.defT1FS.marks[k];
              this.modelValue.ratings[k] = this.defT1FS.ratings[k];
            }
          } else if (this.selectedType === 'Интервал') {
            for (let k in this.defInt.marks) {
              this.modelValue.marks[k] = this.defInt.marks[k];
              this.modelValue.ratings[k] = this.defInt.ratings[k];
            }
          } else {
            for (let k in this.defNum.marks) {
              this.modelValue.marks[k] = this.defNum.marks[k];
              this.modelValue.ratings[k] = this.defNum.ratings[k];
            }
          }
          this.isDefaultLing = true;
          this.isValidRatings = new Array(this.modelValue.ratings.length).fill(true);
          this.isValidMarks = new Array(this.modelValue.marks.length).fill(true);

          this.throw();
        } else {
          this.isDefaultLing = false;
        }
      },
      changeType() {
        if (this.selectedType === 'Число') {
          for (let k in this.modelValue.ratings) {
            this.modelValue.ratings[k] = this.defNum === null ? 0.0 : this.defNum.ratings[k];
          }
          this.isDefaultLing = this.deepEqual(this.modelValue, this.defNum);
        } else if (this.selectedType === 'Интервал') {
          for (let k in this.modelValue.ratings) {
            this.modelValue.ratings[k] = this.defInt === null ? {start: 0.0, end: 0.0} : this.defInt.ratings[k];
          }
          this.isDefaultLing = this.deepEqual(this.modelValue, this.defInt);
        } else if (this.selectedType === 'T1FS') {
          for (let k in this.modelValue.ratings) {
          this.modelValue.ratings[k] = this.defT1FS === null ? {vert: [0.0, 0.0, 0.0]} : this.defT1FS.ratings[k];
        }
        this.isDefaultLing = this.deepEqual(this.modelValue, this.defT1FS);
        } else if (this.selectedType === 'AIFS') {
          for (let k in this.modelValue.ratings) {
            this.modelValue.ratings[k] = {vert: [0.0, 0.0, 0.0], pi: 0.0};
          }
          this.isDefaultLing = false;
        } else if (this.selectedType === 'IT2FS') {
          for (let k in this.modelValue.ratings) {
            this.modelValue.ratings[k] = {bottom: [{start: 0.0, end: 0.0}, {start: 0.0, end: 0.0}], upward: [0.0]};
          }
          this.isDefaultLing = false;
        }
        this.curType = this.selectedType;
      },
      throw() {
        for (let k in this.isValidRatings) {
          if (!this.isValidRatings[k] || ! this.isValidMarks[k]) {
            this.$emit('incorr-scale');
            return
          }
        }

        this.$emit('corr-scale')
      },
      validateMark(i) {
        this.isDefaultLing = false;
        this.isValidMarks[i] = this.modelValue.marks[i].length > 0 && this.modelValue.marks[i].length < 51;

        this.throw();
      },
      validateRating(i, fact) {
        this.isDefaultLing = false;
        this.isValidRatings[i] = fact;

        this.throw();
      },
      addRating() {
        if (this.modelValue.ratings.length === 20) {
          return
        }

        this.isDefaultLing = false;
        this.modelValue.marks.push("");
        if (this.selectedType === 'Число') {
          this.modelValue.ratings.push(0.0);
        } else if (this.selectedType === 'Интервал') {
          this.modelValue.ratings.push({start: 0.0, end: 0.0});
        } else if (this.selectedType === 'T1FS') {
          this.modelValue.ratings.push({vert: [0.0, 0.0, 0.0]});
        } else if (this.selectedType === 'AIFS') {
          this.modelValue.ratings.push({vert: [0.0, 0.0, 0.0], pi: 0.0});
        } else if (this.selectedType === 'IT2FS') {
          this.modelValue.ratings.push({bottom: [{start: 0.0, end: 0.0}, {start: 0.0, end: 0.0}], upward: [0.0]});
        }
        this.isValidRatings.push(true);
        this.isValidMarks.push(false);
        this.throw();
      },
      popRating() {
        if (this.modelValue.marks.length === 1) {
          return
        }

        this.modelValue.marks.pop();
        this.modelValue.ratings.pop();
      }
    },
    mounted() {
      this.defNum = this.$store.getters['getDefaultLingNumScale'];
      this.defInt = this.$store.getters['getDefaultLingIntScale'];
      this.defT1FS = this.$store.getters['getDefaultLingT1FScale'];
      if (this.defNum === null || this.defInt === null || this.defT1FS === null) {
        this.showWarning = true;
        this.warning = 'Не удалось загрузить стандартные лингвистические шкалы';
      }

      if (this.modelValue.ratings[0] === null || typeof this.modelValue.ratings[0] === 'number') {
        this.selectedType = 'Число';
        if (this.deepEqual(this.modelValue, this.defNum)) {
          this.isDefaultLing = true;
        }
      } else if ('bottom' in this.modelValue.ratings[0]) {
        this.selectedType = 'IT2FS';
      } else if ('pi' in this.modelValue.ratings[0]) {
        this.selectedType = 'AIFS';
      } else if ('vert' in this.modelValue.ratings[0]) {
        this.selectedType = 'T1FS';
        if (this.deepEqual(this.modelValue, this.defT1FS)) {
          this.isDefaultLing = true;
        }
      } else if ('start' in this.modelValue.ratings[0]) {
        this.selectedType = 'Интервал';
        if (this.deepEqual(this.modelValue, this.defInt)) {
          this.isDefaultLing = true;
        }
      } else {
        console.log("warning:", this.modelValue.ratings[0]);
        this.showWarning = true;
        this.warning = 'Неизвестный тип в лингвистической шкале';
      }

      this.curType = this.selectedType;
      this.isValidRatings = new Array(this.modelValue.ratings.length).fill(true);
      this.isValidMarks = new Array(this.modelValue.marks.length).fill(true);
    }
  }
</script>

<template>
  <div>
    <div class="ling-scale row-cols-2">
      <div class="col-4">
        <p>Лингвистическая шкала:</p>
      </div>
      <div class="col-8">
        <div class="select">
          <button type="button" :class="{opt:true, non_active: !isDefaultLing, is_active: isDefaultLing}"
                  @click="chooseLingScale" :disabled="role==='expert'">По умолчанию
          </button>
          <button type="button" :class="{opt: true, non_active: isDefaultLing, is_active: !isDefaultLing}"
                  @click="chooseLingScale" :disabled="role==='expert'">Своя конфигурация
          </button>
        </div>

        <p>Тип оценкки:</p>
        <select v-model="selectedType" @change="changeType" :disabled="role==='expert'">
          <option>Число</option>
          <option>Интервал</option>
          <option>T1FS</option>
          <option>AIFS</option>
          <option>IT2FS</option>
        </select>
      </div>
    </div>

    <div class="row-cols-2">
      <div class="col-4">
        <div id="warning" v-if="showWarning===true">
          <p>{{ warning }}</p>
        </div>
      </div>
      <div class="col-8 scale">
        <table class="table table-bordered">
          <tbody>
          <tr v-for="(_, i) in lingScale.marks">
            <td class="mark">
              <input type="text" name="mark" :class="{field: true, invalid: !isValidMarks[i]}" placeholder="mark"
                     v-model="lingScale.marks[i]" :readonly="role==='expert'" @input="validateMark(i)"/>
            </td>
            <td class="eval">
              <NumberForm v-if="curType==='Число'" v-model="lingScale.ratings[i]" v-model:role="role"
                          @corr-rate="validateRating(i, true)" @incorr-rate="validateRating(i, false)"></NumberForm>
              <IntervalForm v-if="curType==='Интервал'" v-model="lingScale.ratings[i]" v-model:role="role"
                            @corr-rate="validateRating(i, true)" @incorr-rate="validateRating(i, false)"></IntervalForm>
              <T1FSForm v-if="curType==='T1FS'" v-model="lingScale.ratings[i]" v-model:role="role"
                        @corr-rate="validateRating(i, true)" @incorr-rate="validateRating(i, false)"></T1FSForm>
              <AIFSForm v-if="curType==='AIFS'" v-model="lingScale.ratings[i]" v-model:role="role"
                        @corr-rate="validateRating(i, true)" @incorr-rate="validateRating(i, false)"></AIFSForm>
              <IT2FSForm v-if="curType==='IT2FS'" v-model="lingScale.ratings[i]" v-model:role="role"
                         @corr-rate="validateRating(i, true)" @incorr-rate="validateRating(i, false)"></IT2FSForm>
            </td>
          </tr>
          </tbody>
        </table>
        <div>
          <button type="button" class="blk-btn" @click="addRating" :disabled="role==='expert'">+</button>
          <button type="button" class="blk-btn" @click="popRating" :disabled="role==='expert'">-</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  @import "../../style.css";
  @import "../../assets/tasksettings.css";

  .ling-scale {
    margin-bottom: 4vmin;
  }

  .blk-btn {
    margin-right: 1vmin;
  }

  #warning p {
    color: red;
    text-align: right;
  }

  .scale {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .col-8 > p {
    margin-left: 3vmin;
  }

  select {
    background-color: black;
    font-family: "Inria Sans", sans-serif;
    font-size: 2.5vmin;
    font-weight: 700;
    color: #ABF8F4;
    border-radius: 1em;
    border: 0;
    width: fit-content;
    margin-top: 1.5vmin;
  }

  table {
    border-collapse: separate;
    border-spacing: 0;
    border: 2px solid black;
    border-radius: 10px;
    width: max-content;
    overflow: hidden;
  }

  table tr {
    width: fit-content;
  }

  table tr {
    display: flex;
    flex-direction: row;
    width: fit-content;
  }

  td {
    background-color: #B8FCF8;
    display: flex;
    justify-content: center;
    width: fit-content;
  }

  table .mark {
    width: fit-content;
    align-items: center;
  }

  .mark .field {
    height: 5vmin;
    margin: 0;
  }
</style>