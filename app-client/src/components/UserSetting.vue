<script>
  export default {
    emits: ['show-component'],
    data() {
      return {
        user: {
          login: "",
          email: "",
          password: "",
        },
        isValidLogin: true,
        isValidEmail: true,
        isValidPassword: true,
        success: false
      }
    },
    methods: {
      showMain() {
        this.$emit('show-component', 'Main');
      },
      showPP() {
        this.$emit('show-component', 'Personal');
      },
      validateLogin() {
        this.isValidLogin = this.user.login.length > 0 && this.user.login.length < 101;
      },
      validateEmail() {
        this.isValidEmail = this.user.email.length > 0 && this.user.email.length < 101;
      },
      validatePassword() {
        this.isValidPassword = (this.user.password === '' || (this.user.password.length > 5 && this.user.password.length < 101));
      },
      async submit() {
        await this.$store.dispatch('updateUser', this.user);
        if (this.$store.getters['errorOccurred']) {
          this.$emit('show-component', 'ErrorPage');
        } else {
          this.success = true;
        }
      },
      logout() {
        this.$store.dispatch('logout');
        this.$emit('show-component', 'Main');
      },
      async deleteUser() {
        await this.$store.dispatch('deleteUser');
        if (this.$store.getters['errorOccurred']) {
          this.$emit('show-component', 'ErrorPage');
        } else {
          this.$emit('show-component', 'Main');
        }
      }
    },
    async mounted() {
      await this.$store.dispatch('reqUserInfo');
      if (this.$store.getters['errorOccurred']) {
        this.$emit('show-component', 'ErrorPage');
      } else {
        this.user = this.$store.getters['getUser'];
      }
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
    <h1>Настройки</h1>
    <div class="setting">
      <div v-if="success">
        <p>Данные успешно изменены</p>
      </div>
      <div>
        <p>Смена логина:</p>
        <input type="text" name="login" placeholder="login" :class="{field: true, invalid: !isValidLogin}"
               v-model="user.login" maxlength="80" @input="validateLogin"/>
      </div>
      <div>
        <p>Смена почты:</p>
        <input type="text" name="email" placeholder="email" :class="{field: true, invalid: !isValidEmail}"
               v-model="user.email" maxlength="80" @input="validateEmail"/>
      </div>
      <div>
        <p>Смена пароля:</p>
        <input type="text" name="password" placeholder="password" :class="{field: true, invalid: !isValidPassword}"
               v-model="user.password" maxlength="80" @input="validatePassword"/>
      </div>
    </div>
    <div class="btns">
      <button class="blk-btn" @click="submit">Изменить</button>
    </div>
    <div class="btns">
      <button class="cl-btn" @click="logout">Выйти из профиля</button>
      <button class="cl-btn del" @click="deleteUser">Удалить профиль</button>
    </div>
  </div>

  <div class="footer container-fluid bottom-0">
    <div style="display:flex; width:fit-content; height:100%; cursor: pointer" @click="showPP">
      <img alt="" src="/arrow.png" class="left-arrow">
      <p>Вернутся</p>
    </div>
  </div>
</template>

<style scoped>
@import "../style.css";
@import "../assets/usersetting.css";
</style>