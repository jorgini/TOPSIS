import Vuex from 'vuex';
import createPersistedState from 'vuex-persistedstate';
import auth from './auth.js';
import page from "./page.js";
import taskList from "./taskslist.js";
import error from "./error.js";
import task from "./task.js";

const store = new Vuex.Store({
    modules: {
        auth,
        page,
        task,
        error,
        taskList
    },
    plugins: [createPersistedState()],
});

if (!store.getters['isAuthenticated'] && store.getters['getPage'] !== 'Auth') {
    store.dispatch('showPage', 'Main');
}

if (store.getters['isAuthenticated']) {
    store.dispatch('takeTasks');
}

export default store;