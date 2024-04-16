<script>
  import Task from "./dynamic/TaskShortcard.vue";
  import axios from "axios";
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
        this.isValidPassword = (this.password.length >= 4 && this.password.length < 101);
      },
      validateTitle() {
        this.isValidTitle = (this.title.length > 0 && this.title.length < 101);
      },
      showMain() {
        return this.$emit('show-component', 'Main');
      },
      showSetting() {
        this.$emit('show-component', 'UserSet');
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

        await this.$store.dispatch('connectToTask', {sid: this.sid, pass: this.password})
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
      },
      async downloadFiles() {
        try {
          const response = await axios.get('/Вычислительная схема.pdf', {
            responseType: 'blob',
          });

          const url = window.URL.createObjectURL(new Blob([response.data]));
          const link = document.createElement('a');
          link.href = url;
          link.setAttribute('download', 'Вычислительная схема.pdf');
          document.body.appendChild(link);
          link.click();

          const response2 = await axios.get('/Требования к входным данным.pdf', {
            responseType: 'blob',
          });

          const url2 = window.URL.createObjectURL(new Blob([response2.data]));
          const link2 = document.createElement('a');
          link2.href = url2;
          link2.setAttribute('download', 'Требования к входным данным.pdf');
          document.body.appendChild(link2);
          link2.click();
        } catch (error) {
          console.error('Ошибка при скачивании файла:', error);
          this.$store.commit('setError', error);
          this.$store.commit('setStatus', 500);
          this.$emit('show-component', 'ErrorPage');
        }
      },
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
    <img alt="about" src="/about.png" class="about" @click="downloadFiles">
    <img alt="settings" src="/settings.png" class="settings" @click="showSetting">
  </div>
</template>

<style scoped>
  @import "../style.css";
  @import "../assets/personal.css";
</style>