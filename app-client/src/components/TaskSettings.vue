<script>
  import config from '../../config.yaml';
  import LingScale from "./dynamic/LingScale.vue";
  export default {
    component: { LingScale },
    emits: ['show-component'],
    data() {
      return {
        task: {
          sid: null,
          title: null,
          description: null,
          method: null,
          task_type: null,
          calc_settings: null,
        },
        prevType: null,
        isVisibleDescription: false,
        isValidTitle: true,
        isValidPassword: true,
        isPasswordInstalled: false,
        password: "",
        defaultCalc: config.backend.default_calc,
        isDefaultCalc: true,
        isCalcVisible: false,
        calc: {
          isDefaultNormValue: true,
          isDefaultNormWeight: true,
          isDefaultRanking: true,
          isDefaultFsDist: true,
          isDefaultNumDist: true,
          isDefaultAggregate: true,
        },
        role: null
      }
    },
    methods: {
      showPP() {
        this.$emit('show-component', 'Personal');
      },
      showMain() {
        this.$emit('show-component', 'Main');
      },
      showPassDialog() {
        const modal = document.getElementById('info');
        modal.showModal();
      },
      async showAlts() {
        if (!this.isValidTitle) {
          return;
        }

        if (this.role === 'expert') {
          this.$emit('show-component', 'Alts');
          return
        }

        if (this.prevType === 'individual' && this.task.task_type === 'group' && !this.isPasswordInstalled) {
          this.showPassDialog();
          return;
        }

        if (this.prevType === 'group' && this.task.task_type === 'individual') {
          const modal = document.getElementById('warning');
          modal.showModal();
          return
        }

        await this.updateTask();
      },
      async updateTask() {
        await this.$store.dispatch('updateTask', this.task);
        if (this.$store.getters['errorOccurred']) {
          console.log(this.$store.getters['errorOccurred']);
          this.$emit('show-component', 'ErrorPage');
        } else {
          console.log(this.task.calc_settings);
          this.$emit('show-component', 'Alts');
        }
      },
      async submitPass() {
        this.validatePassword()
        if (!this.isValidPassword) {
          // todo show warning
          return
        }

        await this.$store.dispatch('setPass', {sid: this.task.sid, password: this.password});
        if (this.$store.getters['errorOccurred']) {
          this.$emit('show-component', 'ErrorPage');
          return
        }
        this.closeModal()
        this.isPasswordInstalled = true;
      },
      closeModal() {
        const modal = document.getElementById('info');
        modal.close();
        const modal2 = document.getElementById('warning');
        modal2.close();
      },
      validateTitle() {
        this.isValidTitle = this.task.title.length > 0 && this.task.title.length < 101;
      },
      validatePassword() {
        this.isValidPassword = this.password.length >= 4 && this.password.length < 101;
      },
      switchDescription() {
        this.isVisibleDescription = !this.isVisibleDescription;
      },
      switchViewCalc() {
        this.isCalcVisible = !this.isCalcVisible;
      },
      chooseMethod(event) {
        if (this.role === 'expert')
          return
        if (event.target.textContent.trim() === 'SMART') {
          this.task.method = 'smart';
        } else {
          this.task.method = 'topsis';
        }
      },
      chooseType(event) {
        if (this.role === 'expert')
          return
        if (event.target.textContent.trim() === 'Индивидуальная') {
          this.task.task_type = 'individual';
        } else {
          this.task.task_type = 'group';
        }
      },
      chooseCalc(event) {
        if (this.role === 'expert')
          return
        if (event.target.textContent.trim() === 'По умолчанию') {
          this.task.calc_settings = this.defaultCalc;
          this.isDefaultCalc = true;
          this.calc.isDefaultNormValue = true;
          this.calc.isDefaultNormWeight = true;
          this.calc.isDefaultRanking = true;
          this.calc.isDefaultFsDist = true;
          this.calc.isDefaultNumDist = true;
          this.calc.isDefaultAggregate = true;
        } else {
          this.isDefaultCalc = false;
        }
      },
      chooseNormValue(event) {
        if (this.role === 'expert')
          return
        if (event.target.textContent.trim() === 'По сумме') {
          this.calc.isDefaultNormValue = true;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111));
          this.task.calc_settings |= config.calc.normalizeWithSum;
        } else {
          this.calc.isDefaultNormValue = false;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111));
          this.task.calc_settings |= config.calc.normalizeValueWithMax;
          this.isDefaultCalc = false;
        }
      },
      chooseNormWeight(event) {
        if (this.role === 'expert')
          return
        if (event.target.textContent.trim() === 'По сумме') {
          this.calc.isDefaultNormWeight = true;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 4));
          this.task.calc_settings |= (config.calc.normalizeWithSum << 4);
        } else {
          this.calc.isDefaultNormWeight = false;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 4));
          this.task.calc_settings |= (config.calc.normalizeWeightsByMidPoint << 4);
          this.isDefaultCalc = false;
        }
      },
      chooseRanking(event) {
        if (this.role === 'expert')
          return
        if (event.target.textContent.trim() === 'По умолчанию') {
          this.calc.isDefaultRanking = true;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 8));
          this.task.calc_settings |= (config.calc.default << 8);

          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 16));
          this.task.calc_settings |= (config.calc.default << 16);
        } else {
          this.calc.isDefaultRanking = false;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 8));
          this.task.calc_settings |= (config.calc.sengupta << 8);
          this.isDefaultCalc = false;

          this.calc.isDefaultFsDist = false;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 12));
          this.task.calc_settings |= (config.calc.alphaSlices << 12);

          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 16));
          this.task.calc_settings |= (config.calc.sengupta << 16);
        }
      },
      chooseFsDist(event) {
        if (this.role === 'expert')
          return
        if (event.target.textContent.trim() === 'По умолчанию') {
          this.calc.isDefaultFsDist = true;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 12));
          this.task.calc_settings |= (config.calc.default << 12);

          this.calc.isDefaultRanking = true;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 8));
          this.task.calc_settings |= (config.calc.default << 8);
        } else {
          this.calc.isDefaultFsDist = false;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 12));
          this.task.calc_settings |= (config.calc.alphaSlices << 12);
          this.isDefaultCalc = false;
        }
      },
      chooseNumDist(event) {
        if (this.role === 'expert')
          return
        if (event.target.textContent.trim() === 'Квадратичная метрика') {
          this.calc.isDefaultNumDist = true;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 20));
          this.task.calc_settings |= (config.calc.sqrtDistance << 20);
        } else {
          this.calc.isDefaultNumDist = false;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 20));
          this.task.calc_settings |= (config.calc.cbrtDistance << 20);
          this.isDefaultCalc = false;
        }
      },
      chooseAggregate(event) {
        if (this.role === 'expert')
          return
        if (event.target.textContent.trim() === 'Агрегация матриц') {
          this.calc.isDefaultAggregate = true;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 24));
          this.task.calc_settings |= (config.calc.aggregateMatrix << 24);
        } else {
          this.calc.isDefaultAggregate = false;
          this.task.calc_settings = this.task.calc_settings & (~(0b1111 << 24));
          this.task.calc_settings |= (config.calc.aggregateFinals << 24);
          this.isDefaultCalc = false;
        }
      }
    },
    async mounted() {
      this.task =  this.$store.getters['getTaskSettings'];
      this.prevType = this.task.task_type;
      if (this.prevType === 'group') {
        this.isPasswordInstalled = true;
      }

      this.role = await this.$store.dispatch('getRole', this.task.sid);
      if (this.$store.getters['errorOccurred']) {
        console.log(this.$store.getters['errorOccurred']);
        this.$emit('show-component', 'ErrorPage');
      }

      if (this.task.calc_settings !== this.defaultCalc) {
        this.isDefaultCalc = false;
        this.calc.isDefaultNormValue = (((this.task.calc_settings & 0b1111)) === 0b0000);
        this.calc.isDefaultNormWeight = ((((this.task.calc_settings >> 4) & 0b1111)) === 0b0000);
        this.calc.isDefaultRanking = ((((this.task.calc_settings >> 8) & 0b1111)) === 0b1010);
        this.calc.isDefaultFsDist = ((((this.task.calc_settings >> 12) & 0b1111)) === 0b1010);
        this.calc.isDefaultNumDist = ((((this.task.calc_settings >> 20) & 0b1111)) === 0b0101);
        this.calc.isDefaultAggregate = ((((this.task.calc_settings >> 24) & 0b1111)) === 0b0111);
      }
    }
  }
</script>

<template>
  <div class="content">
    <div class="container-fluid header">
      <h2 class="main">Decision Maker</h2>
      <h3>Настройки задачи</h3>
      <button class="cl-btn" @click="showMain">Главная</button>
    </div>

    <form class="params">
      <div class="title row-cols-2">
        <div class="col-4">
          <p>Название:</p>
        </div>
        <div class="col-8">
          <input type="text" :class="{field: true, invalid: !isValidTitle}" name="title" :readonly="role==='expert'"
                 placeholder="title" maxlength="100" v-model="task.title" @input="validateTitle" required/>
        </div>
      </div>

      <div class="description row-cols-2">
        <div class="col-4">
          <img alt="flag" class="flag close" src="/arrow.png" @click="switchDescription"/>
          <p>Описание:</p>
        </div>
        <div class="col-8">
          <textarea v-if="isVisibleDescription" type="text" class="field" name="description" placeholder="description"
                    maxlength="1000" v-model="task.description" :readonly="role==='expert'"/>
        </div>
      </div>

      <div class="method row-cols-2">
        <div class="col-4">
          <p>Используемый метод:</p>
        </div>
        <div class="col-8">
          <div class="select">
            <button type="button" :disabled="role==='expert'"
                    :class="{opt: true, non_active: task.method !== 'smart', is_active: task.method === 'smart'}"
                    @click="chooseMethod">SMART</button>
            <button type="button" :disabled="role==='expert'"
                    :class="{opt:true, non_active:task.method !== 'topsis', is_active: task.method === 'topsis'}"
                    @click="chooseMethod">TOPSIS</button>
          </div>
        </div>
      </div>

      <div class="type row-cols-2">
        <div class="col-4">
          <p>Тип задачи:</p>
        </div>
        <div class="col-8">
          <div class="select">
            <button type="button" :disabled="role==='expert'"
                    :class="{opt: true, non_active: task.task_type !== 'individual', is_active: task.task_type === 'individual'}"
                    @click="chooseType">Индивидуальная
            </button>
            <button type="button" :disabled="role==='expert'"
                    :class="{opt: true, non_active: task.task_type !== 'group', is_active: task.task_type === 'group'}"
                    @click="chooseType">Групповая
            </button>
          </div>
          <img v-if="prevType==='group' && task.task_type==='group' && role==='maintainer'" alt="" src="/about.png"
               class="about" @click="showPassDialog">
        </div>
      </div>

      <div class="calc row-cols-2">
        <div class="col-4">
          <img alt="flag" :class="{flag: true, close: !isCalcVisible, open: isVisibleDescription}" src="/arrow.png"
               @click="switchViewCalc"/>
          <p>Настройка вычислений:</p>
        </div>
        <div class="col-8">
          <div class="select">
            <button type="button" :class="{opt:true, non_active: !isDefaultCalc, is_active: isDefaultCalc}"
                    @click="chooseCalc" :disabled="role==='expert'">По умолчанию
            </button>
            <button type="button" :class="{opt: true, non_active: isDefaultCalc, is_active: !isDefaultCalc}"
                    @click="chooseCalc" :disabled="role==='expert'">Своя конфигурация
            </button>
          </div>
        </div>
      </div>

      <div v-if="isCalcVisible" id="options">
        <div class="calc row-cols-2">
          <div class="col-4">
            <p>Нормализация оценок:</p>
          </div>
          <div class="col-8">
            <div class="select">
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: !calc.isDefaultNormValue, is_active: calc.isDefaultNormValue}"
                      @click="chooseNormValue">По сумме
              </button>
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: calc.isDefaultNormValue, is_active: !calc.isDefaultNormValue}"
                      @click="chooseNormValue">По максимуму
              </button>
            </div>
          </div>
        </div>

        <div class="calc row-cols-2">
          <div class="col-4">
            <p>Нормализация весов:</p>
          </div>
          <div class="col-8">
            <div class="select">
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: !calc.isDefaultNormWeight, is_active: calc.isDefaultNormWeight}"
                      @click="chooseNormWeight">По сумме
              </button>
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: calc.isDefaultNormWeight, is_active: !calc.isDefaultNormWeight}"
                      @click="chooseNormWeight">По средней точке
              </button>
            </div>
          </div>
        </div>

        <div class="calc row-cols-2">
          <div class="col-4">
            <p>Алгоритм ранжирования:</p>
          </div>
          <div class="col-8">
            <div class="select">
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: !calc.isDefaultRanking, is_active: calc.isDefaultRanking}"
                      @click="chooseRanking">По умолчанию
              </button>
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: calc.isDefaultRanking, is_active: !calc.isDefaultRanking}"
                      @click="chooseRanking">Сенгупта
              </button>
            </div>
          </div>
        </div>

        <div class="calc row-cols-2">
          <div class="col-4">
            <p>Расстояния для нечетких множеств:</p>
          </div>
          <div class="col-8">
            <div class="select">
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: !calc.isDefaultFsDist, is_active: calc.isDefaultFsDist}"
                      @click="chooseFsDist">По умолчанию
              </button>
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: calc.isDefaultFsDist, is_active: !calc.isDefaultFsDist}"
                      @click="chooseFsDist">Альфа-срезы
              </button>
            </div>
          </div>
        </div>

        <div class="calc row-cols-2">
          <div class="col-4">
            <p>Расстояния для чисел:</p>
          </div>
          <div class="col-8">
            <div class="select">
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: !calc.isDefaultNumDist, is_active: calc.isDefaultNumDist}"
                      @click="chooseNumDist">Квадратичная метрика
              </button>
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: calc.isDefaultNumDist, is_active: !calc.isDefaultNumDist}"
                      @click="chooseNumDist">Кубическая метрика
              </button>
            </div>
          </div>
        </div>

        <div class="calc row-cols-2">
          <div class="col-4">
            <p>Агрегация оценок:</p>
          </div>
          <div class="col-8">
            <div class="select">
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: !calc.isDefaultAggregate, is_active: calc.isDefaultAggregate}"
                      @click="chooseAggregate">Агрегация матриц
              </button>
              <button type="button" :disabled="role==='expert'"
                      :class="{opt: true, non_active: calc.isDefaultAggregate, is_active: !calc.isDefaultAggregate}"
                      @click="chooseAggregate">Агрегация результатов
              </button>
            </div>
          </div>
        </div>
      </div>
<!--      <LingScale></LingScale>-->
    </form>

    <dialog id="info">
      <h3>Найстройки групповой задачи</h3>
      <p>Идентификатор данной задачи для подключения: {{ task.sid }}</p>
      <p>{{ prevType==='group' ? 'Можете поменять пароль, если вы забыли его (это не повлияет на уже подключенных экспертов):' :
        'Установите пароль от данной задачи:' }}</p>
      <input type="text" name="password" :class="{field: true, invalid: !isValidPassword}" placeholder="password"
             v-model="password" @input="validatePassword" required/>
      <div class="btns">
        <button class="blk-btn" @click="submitPass">Подтвердить</button>
        <button class="cl-btn" @click="closeModal">Отмена</button>
      </div>
    </dialog>

    <dialog id="warning">
      <h3>Предупреждение</h3>
      <p>Вы хотите поменять тип задачи с групповой на индивидуальную.</p>
      <p>Если вы продолжите, то пользователи, которые уже подсоединились к задаче автоматически отстранятся от неё.</p>
      <div class="btns">
        <button class="blk-btn" @click="updateTask">Подтвердить</button>
        <button class="cl-btn" @click="closeModal">Отмена</button>
      </div>
    </dialog>
  </div>

  <footer class="footer" style="flex-shrink: 0">
    <div style="display:flex; width:fit-content; height:100%; cursor: pointer" @click="showPP">
      <img alt="" src="/arrow.png" class="left-arrow">
      <p>Отмена</p>
    </div>
    <div style="display: flex; width: fit-content; height: 100%; cursor: pointer" @click="showAlts">
      <p>Перейти к альтернативам</p>
      <img alt="" src="/arrow.png" class="right-arrow">
    </div>
  </footer>
</template>

<style scoped>
  @import "../style.css";
  @import "../assets/tasksettings.css";
</style>