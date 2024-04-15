<script>
  import Task from "./dynamic/TaskShortcard.vue";
  export default {
    emits: ['show-component'],
    components: {Task},
    data() {
      return {
        cntTasks: null,
        tasks: null,
        isValidTitle: true,
        isValidSID: true,
        isValidPassword: true,
        title: '',
        sid: null,
        password: null,
      };
    },
    methods: {
      validateSID() {
        this.isValidSID = (this.sid > 0);
      },
      validatePassword() {
        this.isValidPassword = (this.password.length > 4 && this.password.length < 101);
      },
      validateTitle() {
        this.isValidTitle = (this.title.length > 0 && this.title.length < 101);
      },
      showMain() {
        return this.$emit('show-component', 'Main');
      },
      async confirmCreate() {
        this.validateTitle();
        if (this.isValidTitle) {
          await this.$store.dispatch('createTask', this.title);
          if (this.$store.getters['errorOccurred']) {
            return this.isValidTitle = false;
          } else {
            return this.$emit('show-component', 'TaskSet');
          }
        } else {
          this.$store.commit('setError', "Название не может быть пустым");
        }
      },
      showModal(event) {
        const id = event.target.className.split(' ')[1];
        const modal = document.getElementById(id);
        modal.showModal();
      },
      closeModal() {
        const connect = document.getElementById("connect");
        connect.close();
        const create = document.getElementById("create");
        create.close();
      },
      async showTask(i) {
        await this.$store.dispatch("showTask", this.tasks[i].sid);
        if (this.$store.getters['errorOccurred']) {
          console.log(this.$store.getters['errorOccurred']);
          document.getElementById('warning').innerHTML = "<p class='error'>" + this.$store.getters['errorOccurred'] + "</p>"
        } else {
          this.$emit('show-component', 'TaskSet');
        }
      },
      async deleteTask(i) {
        await this.$store.dispatch('deleteTask', this.tasks[i].sid);
        if (this.$store.getters['errorOccurred']) {
          console.log('Ошибка: ' + this.$store.getters['errorOccurred']);
          document.getElementById('warning').innerHTML = "<p class='error'>" + this.$store.getters['errorOccurred'] + "</p>"
        } else {
          location.reload();
        }
      },
      async confirmConnect() {
        this.validateSID()
        this.validatePassword()

        if (!this.isValidPassword || !this.isValidSID) {
          // todo show warning
          return
        }

        await this.$store.dispatch('connectToTask')
        if (this.$store.getters['errorOccurred']) {
          this.$emit('show-component', 'ErrorPage');
        } else {
          await this.$store.dispatch('showTask', this.sid);
          if (this.$store.getters['errorOccurred']) {
            this.$emit('show-component', 'ErrorPage');
            return
          }
          
          this.$emit('show-component', 'Ratings');
        }
      }
    },
    async mounted() {
      await this.$store.dispatch('takeTasks');
      this.cntTasks = this.$store.getters['getCntTasks'];
      this.tasks = this.$store.getters['getTasks'];
    }
  }
</script>

<template>
  <div class="content">
    <div class="header container-fluid">
      <h2 class="main">Decision Maker</h2>
      <h3>Личный кабинет</h3>
      <button class="cl-btn" @click="showMain">Главная</button>
    </div>
    <h1>Ваши задачи ({{ cntTasks }} / 30)</h1>
    <div id="warning"></div>
    <div id="list">
      <Task v-for="(_, i) in tasks" v-model="tasks[i]" @show-task="showTask(i)" @delete-task="deleteTask(i)"></Task>
    </div>
    <div class="btns">
      <button class="blk-btn create" @click="showModal">Создать задачу</button>
      <button class="blk-btn connect" @click="showModal">Присоединиться к задаче</button>
    </div>

    <dialog id="create">
      <p>Введите название задачи</p>
      <form>
        <input type="text" :class="{field: true, invalid: !isValidTitle}" name="title"
               placeholder="title" maxlength="100" v-model="title" @input="validateTitle" required/>
        <label v-if="this.$store.getters['errorOccurred']" class="error">{{ this.$store.getters['errorOccurred'] }}</label>
        <div class="btns">
          <button type="button" class="blk-btn conf" @click="confirmCreate">Подтвердить</button>
          <button type="reset" class="cl-btn reset" @click="closeModal">Отмена</button>
        </div>
      </form>
    </dialog>
    <dialog id="connect">
      <p>Введите идентификатор задачи и пароль</p>
      <form>
        <input type="number" :class="{field: true, invalid: !isValidSID}" name="sid"
               placeholder="sid" maxlength="20" v-model="sid" @input="validateSID" required/>
        <input type="text" :class="{field: true, invalid: !isValidPassword}" name="password"
               placeholder="password" maxlength="100" v-model="password" @input="validatePassword" required/>
        <label v-if="this.$store.getters['errorOccurred']" class="error">{{ this.$store.getters['errorOccurred'] }}</label>
        <div class="btns">
          <button type="button" class="blk-btn conf" @click="confirmConnect">Подтвердить</button>
          <button type="reset" class="cl-btn reset" @click="closeModal">Отмена</button>
        </div>
      </form>
    </dialog>
  </div>

  <div class="footer container-fluid bottom-0">
    <button class="about"><img alt="about" src="/about.png" width="100%"></button>
    <button class="settings"><img alt="settings" src="/settings.png" width="100%"></button>
  </div>
</template>

<style scoped>
  @import "../style.css";
  @import "../assets/personal.css";
</style>