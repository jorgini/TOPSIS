<script>
import axios from "axios";

export default {
  emits: ['show-component'],
  data() {
    return {
      formData: {
        email: '',
        login: '',
        password: ''
      },
      type: 'login',
      message: '',
      isValid: true,
      isValidLogin: true,
      isValidEmail: true,
      isValidPassword: true
    };
  },
  methods: {
    showMain() {
      return this.$emit('show-component', 'Main');
    },
    validateEmail() {
      const pattern = new RegExp('[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,4}$');
      this.isValid = true;
      this.isValidEmail = pattern.test(this.formData.email);
    },
    validateLogin() {
      this.isValid = true;
      this.isValidLogin = this.formData.login.length > 0 && this.formData.login.length < 80;
    },
    validatePassword() {
      this.isValid = true;
      this.isValidPassword = this.formData.password.length > 5 && this.formData.password.length < 80;
    },
    async login() {
      if (this.type === 'login') {
        this.validateLogin();
        this.validatePassword();
        if ((this.isValidLogin && this.isValidPassword) === false) {
          this.message = 'Неверные входные данные';
          this.isValid = false;
          return;
        } else {
          this.isValid = true;
        }

        await this.$store.dispatch('login', this.formData);
        if (this.$store.getters['errorOccurred']) {
          console.log("warning: " + this.$store.getters['errorOccurred']);
          this.isValid = false;
          this.message = 'Ошибка: ' + this.$store.getters['errorOccurred'];
        } else {
          this.$emit('show-component', 'Personal');
        }
      } else {
        this.type = 'login';
      }
    },
    async signup() {
      if (this.type === 'signup') {
        this.validateEmail();
        this.validateLogin();
        this.validatePassword();
        if ((this.isValidLogin && this.isValidEmail && this.isValidPassword) === false) {
          this.message = 'Неверные входные данные';
          this.isValid = false;
          return;
        } else {
          this.isValid = true;
        }

        await this.$store.dispatch('signup', this.formData);
        if (this.$store.getters['errorOccurred']) {
          console.log("warning: " + this.$store.getters['errorOccurred']);
          this.isValid = false;
          this.message = 'Ошибка: ' + this.$store.getters['errorOccurred'];
        } else {
          this.$emit('show-component', 'Personal');
        }
      } else {
        this.type = 'signup';
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
    }
  }
}
</script>

<template>
  <div class="content">
    <div class="container-fluid header">
      <h2 class="main">Decision Maker</h2>
      <button class="cl-btn" @click="showMain">Главная</button>
    </div>
    <div class="auth container-fluid">
      <h1>Аутентификация</h1>
      <form id="register-form">
        <div class="error" v-if="!isValidLogin">Логин не может быть пустым</div>
        <input type="text" :class="{field: true, invalid: (!isValidLogin || !isValid)}" name="login"
               placeholder="login" maxlength="80" v-model="formData.login" @input="validateLogin" required />
        <div class="error" v-if="!isValidEmail && type==='signup'">Неверный формат ввода</div>
        <input v-if="type==='signup'" type="text" :class="{field: true, invalid: (!isValidEmail || !isValid)}" name="email"
               placeholder="email" maxlength="100" v-model="formData.email" @input="validateEmail" required />
        <div class="error" v-if="!isValidPassword">Минимум 6 символов</div>
        <input type="text" :class="{field: true, invalid: (!isValidPassword || !isValid)}" name="password"
               placeholder="password" maxlength="80" v-model="formData.password" @input="validatePassword" required/>
      </form>

      <div class="buttons">
        <div class="error err-server" v-if="!isValid">{{ message }}</div>
        <button v-if="type==='login'" class="blk-btn log-in" @click="login">Войти</button>
        <button v-if="type==='signup'" class="cl-btn register" @click="signup">Зарегистрироваться</button>

        <button v-if="type==='signup'" class="blk-btn log-in" @click="login">Войти</button>
        <button v-if="type==='login'" class="cl-btn register" @click="signup">Зарегистрироваться</button>
      </div>
    </div>
  </div>

  <div class="footer">
    <img alt="about" src="/about.png" class="about" @click="downloadFiles">
  </div>
</template>

<style scoped>
  @import "../style.css";
  @import "../assets/auth_style.css";
</style>