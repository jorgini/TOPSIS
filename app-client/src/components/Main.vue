<script>
  import axios from 'axios';

  export default {
    methods: {
      showAuth() {
        return this.$emit('show-component', 'Auth');
      },
      showPersonal() {
        return this.$emit('show-component', 'Personal');
      },
      async downloadTech() {
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
        } catch (error) {
          console.error('Ошибка при скачивании файла:', error);
          this.$store.commit('setError', error);
          this.$store.commit('setStatus', 500);
          this.$emit('show-component', 'ErrorPage');
        }
      },
      async downloadInputInfo() {
        try {
          const response = await axios.get('/Требования к входным данным.pdf', {
            responseType: 'blob',
          });

          const url = window.URL.createObjectURL(new Blob([response.data]));
          const link = document.createElement('a');
          link.href = url;
          link.setAttribute('download', 'Требования к входным данным.pdf');
          document.body.appendChild(link);
          link.click();
        } catch (error) {
          console.error('Ошибка при скачивании файла:', error);
          this.$store.commit('setError', error);
          this.$store.commit('setStatus', 500);
          this.$emit('show-component', 'ErrorPage');
        }
      },
    },
  }
</script>

<template>
  <div class="content">
    <div class="container-fluid header">
      <h2 class="main">Decision Maker</h2>
      <button v-if="!this.$store.getters['isAuthenticated']" class="cl-btn sign-up" @click="showAuth">Авторизоваться</button>
      <button v-if="this.$store.getters['isAuthenticated']" class="cl-btn sign-up" @click="showPersonal">Личный кабинет</button>
    </div>
    <h1>Веб-приложение для помощи в решении задач многокритериального принятия решений </h1>
    <div class="info container-fluid">
      <p>Данный продукт реализует методы многокритериального принятия решений (MCDM) TOPSIS и SMART, осуществляющее
        помощь в принятии решений ипредоставляющее аналитический отчет по произведенным вычислениям.</p>
      <p>Уникальные функции и возможности:</p>
      <ul>
        <li>Выставление нечетких оценок</li>
        <li>Создание групповых задач для решения несколькими пользователями</li>
        <li>Анализ чувствительности полученных результатов</li>
        <li>Гибкая настройка вычислений</li>
      </ul>
    </div>
    <button v-if="!this.$store.getters['isAuthenticated']" class="blk-btn start" @click="showAuth">Начать использование</button>
    <button v-if="this.$store.getters['isAuthenticated']" class="blk-btn start" @click="showPersonal">Начать использование</button>
    <div class="subinfo container-fluid">
      <div class="tech">
        <h2>Техническая справка о реализуемых методах</h2>
        <button class="blk-btn download" @click="downloadTech">Скачать</button>
      </div>
      <div class="input">
        <h2>Требования к входным данным</h2>
        <button class="blk-btn download" @click="downloadInputInfo">Скачать</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
  @import "https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css";
  @import "../style.css";
  @import "../assets/main.css";
</style>