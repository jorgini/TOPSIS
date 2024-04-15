<script>
  import config from '../../config.yaml';
  import AltForm from "./dynamic/AltForm.vue";
  export default {
    emits: ['show-component'],
    components: {AltForm},
    data() {
      return {
        task: this.$store.getters['getTaskSettings'],
        default_calc: config.backend.default_calc,
        alts: null
      }
    },
    methods: {
      showMain() {
        this.$emit('show-component', 'Main');
      },
      showTS() {
        this.$emit('show-component', 'TaskSet');
      },
      async showCriteria() {
        for (let i = 0; i < this.alts.length; i++) {
          if (!this.validateTitle(i)) {
            // todo show warning
            return
          }
        }

        const prevAlts = this.$store.getters['getAlts'];
        if (prevAlts.length !== this.alts.length) {
          // todo show warning about deactivate statuses
        }

        await this.$store.dispatch('updateAlts', {sid: this.task.sid, alts: this.alts})
        this.$emit('show-component', 'Criteria')
      },
      addAlt() {
        this.alts.push({title: "", description: ""});
      },
      validateTitle(i) {
        return this.alts[i].title.length > 0 && this.alts[i].title.length < 101;
      },
      deleteAlt(i) {
        if (this.alts.length === 1) {
          return
        }
        this.alts.splice(i, 1);
      },
    },
    async mounted() {
      await this.$store.dispatch('showAlts', this.task.sid);
      if (this.$store.getters['errorOccurred']) {
        console.log(this.$store.getters['errorOccurred']);
        this.$emit('show-component', 'ErrorPage')
        return
      }

      this.alts = this.$store.getters['getAlts'];
      if (this.alts === null || this.alts.length === 0) {
        this.alts = [];
        this.addAlt();
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
    <h1>Задайте альтерантивы</h1>
    <div id="alts" class="container-fluid">
      <AltForm v-for="(_, i) in alts" v-model="alts[i]" @delete-alt="deleteAlt(i)"></AltForm>
    </div>
    <div class="add-alts">
      <button class="blk-btn" @click="addAlt">Добавить ещё</button>
    </div>
  </div>

  <footer class="footer" style="flex-shrink: 0">
    <div style="display:flex; width:fit-content; height:100%; cursor: pointer" @click="showTS">
      <img alt="" src="/arrow.png" class="left-arrow">
      <p>Перейти к настройке задач</p>
    </div>
    <div style="display: flex; width: fit-content; height: 100%; cursor: pointer" @click="showCriteria">
      <p>Перейти к найстройке критериев</p>
      <img alt="" src="/arrow.png" class="right-arrow">
    </div>
  </footer>
</template>

<style scoped>
  @import "../style.css";
  @import "../assets/personal.css";
  @import "../assets/alts.css";

</style>