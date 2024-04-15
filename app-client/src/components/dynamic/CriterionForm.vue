<script setup>
  const criterion = defineModel();
</script>

<script>
  import NumberForm from "./NumberForm.vue";
  import IntervalForm from "./IntervalForm.vue";
  export default {
    emits: ['delete-criterion', 'corr-rate', 'incorr-rate'],
    components: {NumberForm, IntervalForm},
    data() {
      return {
        isValidTitle: true,
        isVisitableDisc: false,
        selectedOpt: 'Число',
        selectedType: 'Преимущественный'
      }
    },
    methods: {
      validate() {
        this.isValidTitle = this.modelValue.title.length > 0 && this.modelValue.title.length < 101;
      },
      deleteAlt() {
        this.$emit('delete-criterion', null);
      },
      switchDisc() {
        this.isVisitableDisc = !this.isVisitableDisc
      },
      changeWeight() {
        if (this.selectedOpt === 'Число') {
          this.modelValue.weight = 0.0
        } else {
          this.modelValue.weight = {start: 0.0, end: 0.0}
        }
      },
      changeType() {
        if (this.selectedType === 'Преимущественный') {
          this.modelValue.type_of_criterion = true;
        } else {
          this.modelValue.type_of_criterion = false;
        }
      }
    },
    mounted() {
      if (this.modelValue.weight === null|| typeof this.modelValue.weight === 'number') {
        this.selectedOpt = 'Число';
      } else {
        this.selectedOpt = 'Интервал';
      }

      if (!this.modelValue.type_of_criterion) {
        this.selectedType = 'Затратный';
      }
    }
  }
</script>

<template>
  <div class="criterion">
    <div class="row-cols-3">
      <div class="col-3">
        <p>Название:</p>
      </div>
      <div class="col-6">
        <input type="text" :class="{field: true, invalid: !isValidTitle}" name="title"
               placeholder="title" maxlength="100" v-model="criterion.title" @input="validate" required/>
      </div>
      <div class="col-3 right-col">
        <img alt="" src="/cancel.png" class="cancel" @click="deleteAlt">
      </div>
    </div>
    <div class="row-cols-3">
      <div class="col-3">
        <img alt="" src="/arrow.png" :class="{flag: true, close: !isVisitableDisc}" @click="switchDisc">
        <p>Описание:</p>
      </div>
      <div class="col-6">
        <textarea type="text" :class="{field: true, invisible: !isVisitableDisc}" name="description"
                  placeholder="description" maxlength="1000" v-model="criterion.description"/>
      </div>
      <div class="col-3"></div>
    </div>
    <div class="row-cols-3">
      <div class="col-3">
        <p>Вес:</p>
      </div>
      <div class="col-9">
        <select v-model="selectedOpt" @change="changeWeight">
          <option>Число</option>
          <option>Интервал</option>
        </select>
        <NumberForm v-if="selectedOpt==='Число'" v-model="criterion.weight" @corr-rate="this.$emit('corr-rate')"
                    @incorr-rate="this.$emit('incorr-rate')"></NumberForm>
        <IntervalForm v-if="selectedOpt==='Интервал'" v-model="criterion.weight" @corr-rate="this.$emit('corr-rate')"
                      @incorr-rate="this.$emit('incorr-rate')"></IntervalForm>
      </div>
    </div>
    <div class="row-cols-3">
      <div class="col-3">
        <p>Тип критерия:</p>
      </div>
      <div class="col-9">
        <select v-model="selectedType" @change="changeType">
          <option>Преимущественный</option>
          <option>Затратный</option>
        </select>
      </div>
    </div>
  </div>
</template>

<style scoped>
@import "../../style.css";
@import "../../assets/criteria.css";

.col-9 {
  display: flex;
  flex-direction: row;
  align-items: center;
}
.col-9 > * {
  margin: auto 2vmin auto 2vmin;
}
</style>