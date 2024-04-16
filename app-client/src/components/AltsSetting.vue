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
        alts: [{title: "", description: ""}],
        prevAltsLength: null,
        role: null
      }
    },
    methods: {
      showMain() {
        this.$emit('show-component', 'Main');
      },
      showTS() {
        this.$emit('show-component', 'TaskSet');
      },
      showWarning() {
        const modal = document.getElementById('warning');
        modal.showModal();
      },
      async showCriteria() {
        for (let i = 0; i < this.alts.length; i++) {
          if (!this.validateTitle(i)) {
            // todo show warning
            return
          }
        }

        if (this.role === 'expert') {
          this.$emit('show-component', 'Criteria');
          return
        }

        if (this.prevAltsLength !== this.alts.length && this.prevAltsLength !== 0) {
          this.showWarning();
          return
        }

        await this.updateAlts();
      },
      async updateAlts() {
        await this.$store.dispatch('updateAlts', {sid: this.task.sid, alts: this.alts});
        if (this.$store.getters['errorOccurred']) {
          console.log(this.$store.getters['errorOccurred']);
          this.$emit('show-component', 'ErrorPage');
        } else {
          this.$emit('show-component', 'Criteria');
        }
      },
      closeModal() {
        const modal = document.getElementById('warning');
        modal.close();
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
        this.prevAltsLength = 0;
        this.alts = [];
        this.addAlt();
      } else {
        this.prevAltsLength = this.alts.length;
      }

      this.role = await this.$store.dispatch('getRole', this.task.sid);
      if (this.$store.getters['errorOccurred']) {
        console.log(this.$store.getters['errorOccurred']);
        this.$emit('show-component', 'ErrorPage')
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
      <AltForm v-for="(_, i) in alts" v-model:alt="alts[i]" v-model:role="role" @delete-alt="deleteAlt(i)"></AltForm>
    </div>
    <div class="add-alts">
      <button class="blk-btn" :disabled="role==='expert'" @click="addAlt">Добавить ещё</button>
    </div>
  </div>

  <dialog id="warning">
    <h3>Предупреждение</h3>
    <p>Вы хотите изменить количество альтернатив.</p>
    <p>Если вы продолжите, то это приведет к обнулению оценок всех экспертов.</p>
    <div class="btns">
      <button class="blk-btn" @click="updateAlts">Подтвердить</button>
      <button class="cl-btn" @click="closeModal">Отмена</button>
    </div>
  </dialog>

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
  @import "../assets/alts.css";
</style>