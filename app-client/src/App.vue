<script>
import Main from "./components/Main.vue";
import Auth from "./components/Auth.vue";
import PersonalPage from "./components/PersonalPage.vue";
import TaskSettings from "./components/TaskSettings.vue";
import TaskList from "./components/dynamic/TaskShortcard.vue";
export default {
  data() {
    return {
      currentPage: this.$store.getters['getPage'],
    }
  },
  components: {
    Main,
    Auth,
    PersonalPage,
    TaskSettings,
    TaskList
  },
  methods: {
    switchComponent(name) {
      this.currentPage = name;
      if (name === 'ErrorPage')
        return
      this.$store.dispatch('showPage', name);
    }
  }
};
</script>

<template>
  <component :is="currentPage" @show-component="switchComponent" />
</template>