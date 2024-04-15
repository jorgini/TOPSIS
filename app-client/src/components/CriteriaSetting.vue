<script>
import config from '../../config.yaml';
import CriterionForm from "./dynamic/CriterionForm.vue";
export default {
  emits: ['show-component'],
  components: {CriterionForm},
  data() {
    return {
      task: this.$store.getters['getTaskSettings'],
      default_calc: config.backend.default_calc,
      criteria: null,
      isValid: []
    }
  },
  methods: {
    showMain() {
      this.$emit('show-component', 'Main');
    },
    showAlts() {
      this.$emit('show-component', 'Alts');
    },
    async showRatings() {
      for (let i = 0; i < this.criteria.length; i++) {
        if (!this.isValid[i] || !this.validateTitle(i)) {
          // todo show warning
          return
        }
      }

      const prevCriteria = this.$store.getters['getCriteria'];
      if (prevCriteria.length !== this.criteria.length) {
        // todo show warning about deactivate statuses
      }

      await this.$store.dispatch('updateCriteria', {sid: this.task.sid, criteria: this.criteria})
      this.$emit('show-component', 'Ratings')
    },
    addCriteria() {
      this.criteria.push({title: "", description: "", type_of_criterion: true, weight: 0.0});
      this.isValid.push(true);
    },
    validateTitle(i) {
      return this.criteria[i].title.length > 0 && this.criteria[i].title.length < 101;
    },
    validateWeight(i, fact) {
      this.isValid[i] = fact;
    },
    deleteCriterion(i) {
      if (this.criteria.length === 1) {
        return
      }
      this.criteria.splice(i, 1);
    },
  },
  async mounted() {
    await this.$store.dispatch('showCriteria', this.task.sid);
    if (this.$store.getters['errorOccurred']) {
      console.log(this.$store.getters['errorOccurred']);
      this.$emit('show-component', 'ErrorPage')
      return
    }

    this.criteria = this.$store.getters['getCriteria'];
    if (this.criteria === null || this.criteria.length === 0) {
      this.criteria = [];
      this.isValid = [];
      this.addCriteria();
    } else {
      for (let i = 0; i < this.criteria.length; i++) {
        this.isValid.push(true);
      }
    }
  }
}
</script>

<template>
  <div class="content">
    <div class="header container-fluid">
      <h2 class="main">Decision Maker</h2>
      <h3>Задача: {{ task.title }}</h3>
      <button class="cl-btn" @click="showMain">Главная</button>
    </div>
    <div class="short-card container-fluid row-cols-3">
      <div class="col">
        <p class="annot">Используемый метод:</p>
        <p>{{ task.method.toUpperCase() }}</p>
      </div>
      <div class="col">
        <p class="annot">Настройки вычислений:</p>
        <p>{{ (task.calc_settings === default_calc) ? "По умолчанию" : "Своя конфигурация" }}</p>
      </div>
      <div class="col">
        <p class="annot">Тип задачи:</p>
        <p>{{ (task.task_type === 'individual') ? "Индивидуальная" : "Групповая" }}</p>
      </div>
    </div>
    <h1>Задайте критерии</h1>
    <div id="criteria" class="container-fluid">
      <CriterionForm v-for="(_, i) in criteria" v-model="criteria[i]" @delete-criterion="deleteCriterion(i)"
            @corr-rate="validateWeight(i, true)" @incorr-rate="validateWeight(i, false)"></CriterionForm>
    </div>
    <div class="add-criterion">
      <button class="blk-btn" @click="addCriteria">Добавить ещё</button>
    </div>
  </div>

  <footer class="footer" style="flex-shrink: 0">
    <div style="display:flex; width:fit-content; height:100%; cursor: pointer" @click="showAlts">
      <img alt="" src="/arrow.png" class="left-arrow">
      <p>Перейти к настройке альтернатив</p>
    </div>
    <div style="display: flex; width: fit-content; height: 100%; cursor: pointer" @click="showRatings">
      <p>Перейти к оценкам</p>
      <img alt="" src="/arrow.png" class="right-arrow">
    </div>
  </footer>
</template>

<style scoped>
@import "../style.css";
@import "../assets/personal.css";
@import "../assets/criteria.css";
</style>