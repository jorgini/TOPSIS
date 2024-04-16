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
      prevCriteriaLength: null,
      isValid: [],
      role: null
    }
  },
  methods: {
    showMain() {
      this.$emit('show-component', 'Main');
    },
    showAlts() {
      this.$emit('show-component', 'Alts');
    },
    showWarning() {
      const modal = document.getElementById('warning');
      modal.showModal();
    },
    async showRatings() {
      for (let i = 0; i < this.criteria.length; i++) {
        if (!this.isValid[i] || !this.validateTitle(i)) {
          // todo show warning
          return
        }
      }

      if (this.role === 'expert') {
        this.$emit('show-component', 'Ratings');
        return
      }

      if (this.prevCriteriaLength !== this.criteria.length && this.prevCriteriaLength !== 0) {
        this.showWarning();
        return;
      }

      await this.updateCriteria();
    },
    async updateCriteria() {
      await this.$store.dispatch('updateCriteria', {sid: this.task.sid, criteria: this.criteria});
      if (this.$store.getters['errorOccurred']) {
        console.log(this.$store.getters['errorOccurred']);
        this.$emit('show-component', 'ErrorPage');
      } else {
        this.$store.dispatch('changeOrd', 0);
        this.$emit('show-component', 'Ratings');
      }
    },
    closeModal() {
      const modal = document.getElementById('warning');
      modal.close();
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

    this.role = await this.$store.dispatch('getRole', this.task.sid);
    if (this.$store.getters['errorOccurred']) {
      console.log(this.$store.getters['errorOccurred']);
      this.$emit('show-component', 'ErrorPage')
      return
    }

    this.criteria = this.$store.getters['getCriteria'];
    if (this.criteria === null || this.criteria.length === 0) {
      this.prevCriteriaLength = 0;
      this.criteria = [];
      this.isValid = [];
      this.addCriteria();
    } else {
      this.prevCriteriaLength = this.criteria.length;
      this.isValid = new Array(this.criteria.length).fill(true);
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
      <CriterionForm v-for="(_, i) in criteria" v-model:criterion="criteria[i]" v-model:role="role"
                     @delete-criterion="deleteCriterion(i)" @corr-rate="validateWeight(i, true)"
                     @incorr-rate="validateWeight(i, false)"></CriterionForm>
    </div>
    <div class="add-criterion">
      <button class="blk-btn" :disabled="role==='expert'" @click="addCriteria">Добавить ещё</button>
    </div>
  </div>

  <dialog id="warning">
    <h3>Предупреждение</h3>
    <p>Вы хотите изменить количество критериев.</p>
    <p>Если вы продолжите, то это приведет к обнулению оценок всех экспертов.</p>
    <div class="btns">
      <button class="blk-btn" @click="updateCriteria">Подтвердить</button>
      <button class="cl-btn" @click="closeModal">Отмена</button>
    </div>
  </dialog>

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
@import "../assets/criteria.css";
</style>