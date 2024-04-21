import Vuex from 'vuex';
import createPersistedState from 'vuex-persistedstate';
import auth from './auth.js';
import page from "./page.js";
import taskList from "./taskslist.js";
import error from "./error.js";
import task from "./task.js";
import defaults from "./defaults.js";

const store = new Vuex.Store({
    modules: {
        auth,
        page,
        task,
        error,
        taskList,
        defaults
    },
    plugins: [createPersistedState()],
});

if (!store.getters['isAuthenticated'] && store.getters['getPage'] !== 'Auth') {
    store.dispatch('showPage', 'Main');
}

if (store.getters['isAuthenticated']) {
    await store.dispatch('takeTasks');
}

await store.dispatch('takeDefaultLingScale');

export default store;