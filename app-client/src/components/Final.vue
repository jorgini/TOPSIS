<script>
  import NumberForm from "./dynamic/NumberForm.vue";
  import IntervalForm from "./dynamic/IntervalForm.vue";
  import Result from "./dynamic/Result.vue";

  export default {
    emits: ['show-component'],
    components: {IntervalForm, NumberForm, Result},
    data() {
      return {
        task: {sid: null, title: "", task_type: ""},
        alts: [],
        isReady: false,
        isVisibleSetting: false,
        isVisibleExperts: false,
        role: null,
        experts: [{login: "", status:""}],
        weights: null,
        weightType: null,
        curType: null,
        isValidWeight: null,
        threshold: null,
        isValidThreshold: true,
        final: null
      }
    },
    methods: {
      showMain() {
        this.$emit('show-component', 'Main');
      },
      showExperts() {
        this.isVisibleExperts = true;
        const modal = document.getElementById('experts');
        modal.showModal();
      },
      showSetting() {
        this.isVisibleSetting = true;
        const modal = document.getElementById('setting');
        modal.showModal();
      },
      showTS() {
        this.$emit('show-component', 'TaskSet');
      },
      showRatings() {
        this.$emit('show-component', 'Ratings');
      },
      showPP() {
        this.$emit('show-component', 'Personal');
      },
      submitWeights() {
        for (const fact of this.isValidWeight) {
          if (!fact) {
            // todo show warning
            return
          }
        }

        this.$store.dispatch('setExpertsWeights', {sid: this.task.sid, weights: this.weights});
        if (this.$store.getters['errorOccurred']) {
          console.log((this.$store.getters['errorOccurred']))
          this.$emit('show-component', 'ErrorPage')
          return
        }
        this.isVisibleExperts = false;
        const modal = document.getElementById('experts');
        modal.close();
      },
      defaultWeights() {
        this.$store.dispatch('setExpertsWeights', {sid: this.task.sid, weights: []});
        if (this.$store.getters['errorOccurred']) {
          console.log((this.$store.getters['errorOccurred']))
          this.$emit('show-component', 'ErrorPage')
          return
        }
        this.isVisibleExperts = false;
        const modal = document.getElementById('experts');
        modal.close();
      },
      submitThreshold() {
        this.validateThreshold();
        if (!this.isValidThreshold) {
          // todo show warning
        } else {
          this.isVisibleSetting = false;
          const modal = document.getElementById('setting');
          modal.close();
        }
      },
      defaultThreshold() {
        this.threshold = -1;
        this.isVisibleSetting = false;
        const modal = document.getElementById('setting');
        modal.close();
      },
      validateWeight(i, fact) {
        this.isValidWeight[i] = fact;
      },
      validateThreshold() {
        this.isValidThreshold = this.threshold >= 0 && this.threshold <= 1;
      },
      changeWeightType(i) {
        if (this.weightType[i] === 'Число') {
          this.weights[i] = 1.0;
          this.curType[i] = 'Число';
        } else {
          this.weights[i] = {start: 1.0, end: 1.0};
          this.curType[i] = 'Интервал';
        }
      },
      async calcFinal() {
        if (this.threshold === null || this.threshold === -1) {
          await this.$store.dispatch('takeFinal', {sid: this.task.sid, threshold: {}});
        } else {
          await this.$store.dispatch('takeFinal', {sid: this.task.sid, threshold: {threshold: this.threshold}});
        }

        if (this.$store.getters['errorOccurred']) {
          if (this.$store.getters['errorOccurred'] === "doesn't complete") {
            // todo show warning
          } else {
            console.log(this.$store.getters['errorOccurred']);
            this.$emit('show-component', 'ErrorPage')
          }
          return
        }
        this.final = this.$store.getters['getFinal'];
        this.isReady = true;
      },
    },
    async mounted() {
      this.task = this.$store.getters['getTaskSettings'];
      this.alts = this.$store.getters['getAlts'];
      this.final = this.$store.getters['getFinal'];

      if (this.final === null || this.final.fid !== this.task.sid) {
        this.showSetting();
        if (this.task.task_type === 'individual') {
          this.defaultWeights();
        } else {
          this.experts = await this.$store.dispatch('getExperts', this.task.sid);
          if (this.$store.getters['errorOccurred']) {
            console.log((this.$store.getters['errorOccurred']))
            this.$emit('show-component', 'ErrorPage')
            return
          }

          this.weights = new Array(this.experts.length).fill(1.0);
          this.weightType = new Array((this.experts.length)).fill('Число');
          this.curType = new Array(this.experts.length).fill('Число');
          this.isValidWeight = new Array(this.experts.length).fill(true);

          this.role = await this.$store.dispatch('getRole', this.task.sid);
          if (this.$store.getters['errorOccurred']) {
            console.log((this.$store.getters['errorOccurred']))
            this.$emit('show-component', 'ErrorPage')
            return
          }
          this.showExperts()
        }
      } else {
        this.isReady = true;
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
    <div class="short-card container-fluid">
      <div v-if="task.task_type==='group'" class="col">
        <button class="cl-btn" @click="showExperts">Эксперты</button>
      </div>
      <div class="col">
        <button class="cl-btn" @click="showSetting">Настройки финального отчета</button>
      </div>
    </div>

    <dialog id="experts">
      <div class="exp" v-for="(_, i) in experts">
        <p>{{ experts[i].login }}</p>
        <p>{{ experts[i].status }}</p>
        <div v-if="role==='maintainer'" class="weight" @change="changeWeightType(i)">
          <select v-model="weightType[i]">
            <option>Число</option>
            <option>Интервал</option>
          </select>
          <NumberForm v-if="curType[i]==='Число'" v-model="weights[i]" @corr-rate="validateWeight(i, true)"
              @incorr-rate="validateWeight(i, false)"></NumberForm>
          <IntervalForm v-if="curType[i]==='Интервал'" v-model="weights[i]" @corr-rate="validateWeight(i, true)"
              @incorr-rate="validateWeight(i, false)"></IntervalForm>
        </div>
      </div>
      <div class="btns">
        <button class="blk-btn" @click="submitWeights">Подтвердить</button>
        <button class="cl-btn" @click="defaultWeights">Пропустить</button>
      </div>
    </dialog>

    <dialog id="setting">
      <p>Порог для анализа чувствительности:</p>
      <input type="number" :class="{field: true, invalid: !isValidThreshold}" name="threshold"
             placeholder="0.0" maxlength="100" v-model="threshold" @input="validateThreshold" required/>
      <div class="btns">
        <button class="cl-btn" @click="showTS">Вернутся к настройкам вычислений</button>
        <button class="blk-btn" @click="submitThreshold">Подтвердить</button>
        <button class="cl-btn" @click="defaultThreshold">Пропустить</button>
      </div>
    </dialog>

    <h1>Финальный отчет</h1>
    <div class="btns">
      <button class="blk-btn" @click="calcFinal">Вычислить</button>
    </div>

    <Result v-if="isReady" :results="final" :alts="alts"></Result>
  </div>

  <footer class="footer" style="flex-shrink: 0">
    <div style="display:flex; width:fit-content; height:100%; cursor: pointer" @click="showRatings">
      <img alt="" src="/arrow.png" class="left-arrow">
      <p>Вернутся к оценкам</p>
    </div>
    <div style="display: flex; width: fit-content; height: 100%; cursor: pointer" @click="showPP">
      <p>Перейти в личный кабинет</p>
      <img alt="" src="/arrow.png" class="right-arrow">
    </div>
  </footer>
</template>

<style scoped>
  @import "../style.css";
  @import "../assets/personal.css";
  @import "../assets/final.css";
</style>