import config from '../../config.yaml'

const state = {
    cntTasks: JSON.parse(localStorage.getItem('cntTasks')) || null,
    tasks: JSON.parse(localStorage.getItem('tasks')) || null
};

const mutations = {
    setTasks(state, tasks) {
        state.cntTasks = tasks.length
        state.tasks = tasks
    }
};

const actions = {
    async takeTasks({ commit, getters }) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/user', {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        })
            .then(response => {
                status = response.status
                return response.json()
            })
            .then(data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    commit('setTasks', data);
                }
            })
            .catch(error => {
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
                console.log('Ошибка:', error);
            });
    },
    async createTask({ commit, dispatch, getters }, title) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/user', {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({title: title})
        })
            .then(response => {
                status = response.status
                return response.json()
            })
            .then(async data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    await dispatch('showTask', data.sid, {root: true});
                }
            })
            .catch(error => {
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
                console.log('Ошибка:', error);
            });
    },
    async deleteTask({ commit, getters }, sid) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/user/tasks?sid=' + sid, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
        })
            .then(response => {
                status = response.status
                return response.json()
            })
            .then(data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                }
            })
            .catch(error => {
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
                console.log('Ошибка:', error);
            });
    },
    async connectToTask({ commit, getters }, {sid, pass}) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + "/solution/connect", {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({sid: sid, password: pass})
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root:true});
                    commit('setStatus', status, {root: true});
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                }
            })
            .catch(error => {
                console.log(error);
                commit('setError', error, {root: true});
                commit('setStatus', status, {root: true});
            })
    }
};

const getters = {
    getCntTasks: (state) => {
        return state.cntTasks
    },
    getTasks: (state) => {
        return state.tasks
    }
};

export default {
    state,
    mutations,
    actions,
    getters
}