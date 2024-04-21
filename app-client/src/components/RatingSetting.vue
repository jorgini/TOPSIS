<script>
import config from '../../config.yaml';
import RatingForm from "./dynamic/RatingForm.vue";
export default {
  emits: ['show-component'],
  components: {RatingForm},
  data() {
    return {
      task: this.$store.getters['getTaskSettings'],
      default_calc: config.backend.default_calc,
      criteria: [],
      alts: [{title: ""}],
      matrix : [],
      ord: 0,
      isValid: [],
      curAlt: null
    }
  },
  methods: {
    showMain() {
      this.$store.commit('setCriteria', this.criteria);
      this.$emit('show-component', 'Main');
    },
    async goPrev() {
      if (this.ord === 0) {
        this.$store.commit('setMatrix', this.matrix);
        this.$emit('show-component', 'Criteria');
      } else {
        this.ord--;
        await this.$store.dispatch('changeOrd', this.ord);
        await this.$store.dispatch('takeMatrix', {sid: this.task.sid, ord: this.ord});
        if (this.$store.getters['errorOccurred']) {
          console.log(this.$store.getters['errorOccurred']);
          this.$emit('show-component', 'ErrorPage')
          return
        }
        location.reload();
      }
    },
    async goNext() {
      for (const fact of this.isValid) {
        if (!fact) {
          // todo show warning
          return
        }
      }

      await this.$store.dispatch('updateMatrix', {sid: this.task.sid, ord: this.ord, ratings: this.matrix});
      if (this.$store.getters['errorOccurred']) {
        console.log(this.$store.getters['errorOccurred']);
        this.$emit('show-component', 'ErrorPage')
        return
      }

      if (this.ord === (this.alts.length - 1)) {
        await this.$store.dispatch('completeTask', this.task.sid);
        if (this.$store.getters['errorOccurred']) {
          console.log(this.$store.getters['errorOccurred']);
          this.$emit('show-component', 'ErrorPage')
          return
        }
        this.$emit('show-component', 'Final');
      } else {
        this.ord++;
        await this.$store.dispatch('changeOrd', this.ord);
        await this.$store.dispatch('takeMatrix', {sid: this.task.sid, ord: this.ord});
        if (this.$store.getters['errorOccurred']) {
          console.log(this.$store.getters['errorOccurred']);
          this.$emit('show-component', 'ErrorPage')
          return
        }
        location.reload();
      }
    },
    validateRate(i, fact) {
      this.isValid[i] = fact;
    },
    updateRate(newVal, i) {
      this.matrix[i] = newVal;
    }
  },
  async mounted() {
    this.ord = this.$store.getters['getOrd'];

    await this.$store.dispatch('takeMatrix', {sid: this.task.sid, ord: this.ord});
    if (this.$store.getters['errorOccurred']) {
      console.log(this.$store.getters['errorOccurred']);
      this.$emit('show-component', 'ErrorPage')
      return
    }

    await this.$store.dispatch('showAlts', this.task.sid);
    await this.$store.dispatch('showCriteria', this.task.sid);

    if (this.$store.getters['errorOccurred']) {
      console.log(this.$store.getters['errorOccurred']);
      this.$emit('show-component', 'ErrorPage')
      return
    }

    this.criteria = this.$store.getters['getCriteria'];
    this.alts = this.$store.getters['getAlts'];
    this.curAlt = this.alts[this.ord].title;

    this.matrix = this.$store.getters['getRatings'];
    this.isValid = new Array(this.matrix.length).fill(true);
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

    <div class="short-card container-fluid row-cols-2" style="min-height: 10vmin">
      <div class="col-3" style="display: flex; align-items: center; justify-content: right">
        <p>Лингвистическая шкала:</p>
      </div>
      <div class="col-9" style="display: flex; align-items: center; justify-content: center">
        <table class="table table-bordered">
          <tbody>
            <tr>
              <td v-for="(_, i) in task.ling_scale.marks">
                <p>{{ task.ling_scale.marks[i] }}</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <h1>Альтернатива: {{ curAlt }}</h1>
    <div id="ratings" class="container-fluid">
      <RatingForm v-for="(_, i) in matrix" :rating="{rate: matrix[i]}" :criterion="criteria[i]" :scale="this.task.ling_scale"
          @corr-rate="validateRate(i, true)" @incorr-rate="validateRate(i, false)"
          @update-rate="updateRate($event, i)"></RatingForm>
    </div>
  </div>

  <footer class="footer" style="flex-shrink: 0">
    <div style="display:flex; width:fit-content; height:100%; cursor: pointer" @click="goPrev">
      <img alt="" src="/arrow.png" class="left-arrow">
      <p>{{ ord > 0 ? 'Перейти к предыдущей альтернативе' : 'Перейти к настройке критериев' }}</p>
    </div>
    <div style="display: flex; width: fit-content; height: 100%; cursor: pointer" @click="goNext">
      <p>{{ ord < alts.length - 1 ? 'Перейти к следующей альтернативе' : 'Перейти к финальному отчету' }}</p>
      <img alt="" src="/arrow.png" class="right-arrow">
    </div>
  </footer>
</template>

<style scoped>
@import "../style.css";
@import "../assets/personal.css";
@import "../assets/ratings.css";
</style>