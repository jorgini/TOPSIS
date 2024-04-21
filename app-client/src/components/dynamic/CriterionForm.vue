<script setup>
  const criterion = defineModel('criterion');
  const role = defineModel('role');
  const emits = defineEmits(['delete-criterion', 'corr-rate', 'incorr-rate']);
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
        this.isValidTitle = this.criterion.title.length > 0 && this.criterion.title.length < 101;
      },
      deleteAlt() {
        if (this.role === 'expert')
          return
        this.$emit('delete-criterion', null);
      },
      switchDisc() {
        this.isVisitableDisc = !this.isVisitableDisc
      },
      changeWeight() {
        if (this.selectedOpt === 'Число') {
          this.criterion.weight = 0.0
        } else {
          this.criterion.weight = {start: 0.0, end: 0.0}
        }
      },
      changeType() {
        if (this.selectedType === 'Преимущественный') {
          this.criterion.type_of_criterion = true;
        } else {
          this.criterion.type_of_criterion = false;
        }
      }
    },
    mounted() {
      if (this.criterion.weight === null|| typeof this.criterion.weight === 'number') {
        this.selectedOpt = 'Число';
      } else {
        this.selectedOpt = 'Интервал';
      }

      if (!this.criterion.type_of_criterion) {
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
        <input type="text" :class="{field: true, invalid: !isValidTitle}" name="title" :readonly="role==='expert'"
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
                  :readonly="role==='expert'" placeholder="description" maxlength="1000" v-model="criterion.description"/>
      </div>
      <div class="col-3"></div>
    </div>
    <div class="row-cols-3">
      <div class="col-3">
        <p>Вес:</p>
      </div>
      <div class="col-9">
        <select v-model="selectedOpt" :disabled="role==='expert'" @change="changeWeight">
          <option>Число</option>
          <option>Интервал</option>
        </select>
        <NumberForm v-if="selectedOpt==='Число'" v-model:role="role" v-model="criterion.weight"
                    @corr-rate="this.$emit('corr-rate')" @incorr-rate="this.$emit('incorr-rate')"></NumberForm>
        <IntervalForm v-if="selectedOpt==='Интервал'" v-model:role="role" v-model="criterion.weight"
                      @corr-rate="this.$emit('corr-rate')" @incorr-rate="this.$emit('incorr-rate')"></IntervalForm>
      </div>
    </div>
    <div class="row-cols-3">
      <div class="col-3">
        <p>Тип критерия:</p>
      </div>
      <div class="col-9">
        <select v-model="selectedType" :disabled="role==='expert'" @change="changeType">
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