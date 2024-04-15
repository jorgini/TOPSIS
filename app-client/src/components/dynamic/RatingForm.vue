<script>
import NumberForm from "./NumberForm.vue";
import IntervalForm from "./IntervalForm.vue";
import T1FSForm from "./T1FSForm.vue";
import AIFSForm from "./AIFSForm.vue";
import IT2FSForm from "./IT2FSForm.vue";
export default {
  emits:['corr-rate', 'incorr-rate', 'update-rate'],
  props: {
    rating: Object,
    criterion: Object
  },
  components: {NumberForm, IntervalForm, T1FSForm, AIFSForm, IT2FSForm},
  data() {
    return {
      isValid: true,
      selectedType: 'Число',
      curType: 'Число',
      localCopy: this.rating.rate
    }
  },
  watch: {
    localCopy: {
      handler(newVal) {
        this.$emit('update-rate', newVal);
      },
      deep: true,
    }
  },
  methods: {
    changeRating() {
      if (this.selectedType === 'Число') {
        this.localCopy = 0.0;
      } else if (this.selectedType === 'Интервал') {
        this.localCopy = {start: 0.0, end: 0.0};
      } else if (this.selectedType === 'T1FS') {
        this.localCopy = {vert: [0.0, 0.0, 0.0]};
      } else if (this.selectedType === 'AIFS') {
        this.localCopy = {vert: [0.0, 0.0, 0.0], pi: 0.0};
      } else if (this.selectedType === 'IT2FS') {
        this.localCopy = {bottom: [{start: 0.0, end: 0.0}, {start: 0.0, end: 0.0}], upward: [0.0]};
      }
      this.curType = this.selectedType;
    },
  },
  mounted() {
    if (this.localCopy === null || typeof this.localCopy === 'number') {
      this.selectedType = 'Число';
    } else if ('bottom' in this.localCopy) {
      this.selectedType = 'IT2FS';
    } else if ('pi' in this.localCopy) {
      this.selectedType = 'AIFS';
    } else if ('vert' in this.localCopy) {
      this.selectedType = 'T1FS';
    } else if ('start' in this.localCopy) {
      this.selectedType = 'Интервал';
    } else {
      console.log("warning:", this.localCopy)
    }
    this.curType = this.selectedType;
  }
}
</script>

<template>
  <div class="rating">
    <div class="row-cols-3">
      <div class="col-3">
        <p>Критерий: {{ criterion.title }}</p>
      </div>
      <div class="col-2">
        <p>Оценка:</p>
        <select v-model="selectedType" @change="changeRating">
          <option>Число</option>
          <option>Интервал</option>
          <option>T1FS</option>
          <option>AIFS</option>
          <option>IT2FS</option>
        </select>
      </div>
      <div class="col-7">
        <NumberForm v-if="curType==='Число'" v-model="localCopy" @corr-rate="this.$emit('corr-rate')"
                    @incorr-rate="this.$emit('incorr-rate')"></NumberForm>
        <IntervalForm v-if="curType==='Интервал'" v-model="localCopy" @corr-rate="this.$emit('corr-rate')"
                      @incorr-rate="this.$emit('incorr-rate')"></IntervalForm>
        <T1FSForm v-if="curType==='T1FS'" v-model="localCopy" @corr-rate="this.$emit('corr-rate')"
                  @incorr-rate="this.$emit('incorr-rate')"></T1FSForm>
        <AIFSForm v-if="curType==='AIFS'" v-model="localCopy" @corr-rate="this.$emit('corr-rate')"
                  @incorr-rate="this.$emit('incorr-rate')"></AIFSForm>
        <IT2FSForm v-if="curType==='IT2FS'" v-model="localCopy" @corr-rate="this.$emit('corr-rate')"
                   @incorr-rate="this.$emit('incorr-rate')"></IT2FSForm>
        <!--todo linguistic rate-->
      </div>
    </div>
  </div>
</template>

<style scoped>
  @import "../../style.css";
  .rating {
    padding: 2vmin 0;
    border-bottom: 2px solid black;
  }

  .row-cols-3 {
    display: flex;
    flex-direction: row;
  }

  .col-2, .col-3 {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
  }

  .col-7 {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
  }

  .col-2 > * {
    margin: auto 1vmin auto 1vmin;
  }

  .col-3 > * {
    margin: auto 1vmin auto 1vmin;
  }

  p {
    font-family: "Inria Sans", sans-serif;
    font-size: 2.5vmin;
    font-weight: 300;
  }

  .col-2 select {
    background-color: black;
    font-family: "Inria Sans", sans-serif;
    font-size: 2.5vmin;
    font-weight: 700;
    color: #ABF8F4;
    border-radius: 1em;
    border: 0;
    width: fit-content;
  }
</style>