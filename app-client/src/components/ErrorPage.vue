<script>
  export default {
    emits: ['show-component'],
    data() {
      return {
        message: this.$store.getters['errorOccurred'],
        status: this.$store.getters['getStatusCode'],
        recom400: 'Проверьте корректность заполняемых данных в соответсвии с требовании к входным данным',
        recom: 'Рекомендуем вернутся на главную страницу и повторить действия',
        badReq: 'Некорректные входные данные',
        unauth: 'Необходима аутентификация',
        forbidden: 'Этот ресурс недоступен вам',
        notFound: 'Данный путь не найден',
        ise: 'Ошибка на сервере',
        unknown: 'Неизвестная ошибка',
        curError: null
      }
    },
    methods: {
      showMain() {
        this.$store.commit('setError', null);
        this.$emit('show-component', 'Main');
      },
      showPP() {
        this.$store.commit('setError', null);
        this.$emit('show-component', 'Personal');
      },
      showPrev() {
        this.$store.commit('setError', null);
        this.$emit('show-component', this.$store.getters['getPage']);
      }
    },
    mounted() {
      if (this.status === 400) {
        this.curError = this.badReq;
      } else if (this.status === 401) {
        this.curError = this.unauth;
      } else if (this.status === 403) {
        this.curError = this.forbidden;
      } else if (this.status === 404) {
        this.curError = this.notFound;
      } else if (this.status === 500) {
        this.curError = this.ise;
      } else {
        this.curError = this.unknown;
      }
    }
  }
</script>

<template>
  <div class="content">
    <div class="header">
      <h2 class="main">Decision Maker</h2>
      <button v-if="this.$store.getters['isAuthenticated']" class="cl-btn" @click="showPP">Личный кабинет</button>
    </div>
    <div class="info">
      <h1>Ошибка</h1>
      <h3>Произошла ошибка: {{ curError }}</h3>
      <p class="bold">{{ status === 400 ? recom400 : recom }}</p>
      <p>{{status}} Error: {{ message }}</p>
      <button class="cl-btn" @click="showMain">Главная</button>
    </div>
  </div>

  <footer class="footer" style="flex-shrink: 0">
    <div style="display:flex; width:fit-content; height:100%; cursor: pointer" @click="showPrev">
      <img alt="" src="/arrow.png" class="left-arrow">
      <p>Вернутся</p>
    </div>
  </footer>
</template>

<style scoped>
  @import "../style.css";

  .info {
    display: block;
    text-align: center;
    width: fit-content;
    margin: 20vmin auto auto auto;
  }

  .info p {
    font-family: "Inria Sans", sans-serif;
    font-weight: 300;
    font-size: 2.5vmin;
  }

  .info .bold {
    font-weight: 700;
  }

  .info .cl-btn {
    margin: 6vmin auto auto auto;
  }

  .left-arrow {
    margin: 1.5vmin 1vmin auto 1.5vmin;
    align-items: start;
    mix-blend-mode: multiply;
    border: 0;
    padding: 0;
    width: 15px;
    height: auto;
    rotate: 90deg;
  }

  .footer {
    justify-content: space-between;
  }

  .footer p {
    font-family: "Inria Sans", sans-serif;
    font-size: 2.5vmin;
    font-weight: 700;
    margin-top: 1vmin;
  }
</style>