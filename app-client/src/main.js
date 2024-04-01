import { createApp } from 'vue'
import store from './store/init.js'
import App from './App.vue'
import Main from './components/Main.vue'
import Auth from './components/Auth.vue'
import PersonalPage from './components/PersonalPage.vue'
import TaskSet from './components/TaskSettings.vue'
import {useIntervalFn} from "@vueuse/core";
import AltsSetting from "./components/AltsSetting.vue";
import CriteriaSetting from "./components/CriteriaSetting.vue";
import RatingSetting from "./components/RatingSetting.vue";
import Final from "./components/Final.vue";
import ErrorPage from "./components/ErrorPage.vue";

const { pause, resume, isActive } = useIntervalFn(() => {
    if (Date.now() > store.getters['getExpiredAt']) {
        store.dispatch('refreshToken');
    }
}, 1000);

createApp(App)
    .use(store)
    .component('Main', Main)
    .component('Auth', Auth)
    .component('Personal', PersonalPage)
    .component('TaskSet', TaskSet)
    .component('Alts', AltsSetting)
    .component('Criteria', CriteriaSetting)
    .component('Ratings', RatingSetting)
    .component('Final', Final)
    .component('ErrorPage', ErrorPage)
    .mount('#app');
