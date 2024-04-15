import config from '../../config.yaml'

const state = {
    token: JSON.parse(localStorage.getItem('token')) || null,
    exp_at: JSON.parse(localStorage.getItem('exp_at')) || null,
};

const mutations = {
    setToken(state, token) {
        state.token = token;
        if (token)
            state.exp_at = Date.now() + (14 * 60 * 1000);
        else
            state.exp_at = Date.now() + (30 * 1000);
    }
};

const actions = {
    async signup({ commit, dispatch }, credentials) {
            let status;
            await fetch(config.backend.url + "/auth/sign-up", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(credentials)
            })
                .then(response => {
                    status = response.status;
                    return response.json();
                })
                .then(async data => {
                    console.log(data)
                    if (status !== 200) {
                        commit('setError', data.message, {root: true});
                    } else {
                        commit('setError', null, {root: true});
                        await dispatch('login', credentials)
                    }
                })
                .catch((error) => {
                    commit('setError', error, {root: true})
                    console.error('Ошибка:', error);
                });
        },
    async login({ commit }, credentials) {
            let status, headers;
            await fetch(config.backend.url + "/auth/log-in", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(credentials),
                credentials: 'include'
            })
                .then(response => {
                    status = response.status;
                    headers = response.headers;
                    return response.json()
                })
                .then(data => {
                    console.log(data)
                    if (status !== 200) {
                        commit('setToken', null);
                        commit('setError', data.message, {root: true});
                    } else {
                        commit('setToken', data.token);
                        commit('setError', null, {root: true})
                    }
                })
                .catch((error) => {
                    commit('setError', error, {root: true})
                    console.error('Ошибка:', error);
                });
        },
    logout({ commit }) {
            commit('setToken', null);
        },
    refreshToken({ commit }) {
            let status, headers;
            fetch(config.backend.url + "/auth/refresh", {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include'
            })
                .then(response => {
                    status = response.status;
                    headers = response.headers;
                    return response.json()
                })
                .then(data => {
                    console.log(data)
                    if (status !== 200) {
                        commit('setToken', null);
                        commit('setError', data.message, {root: true});
                    } else {
                        commit('setToken', data.token);
                        commit('setError', null, {root: true});
                    }
                })
                .catch((error) => {
                    commit('setError', error, {root: true});
                    console.error('Ошибка:', error);
                });
        }
};

const getters = {
    isAuthenticated: (state) => {
        return !!state.token
    },
    getToken: (state) => {
        return state.token
    },
    getExpiredAt: (state) => {
        return state.exp_at
    }
};

export default {
    state,
    mutations,
    actions,
    getters
}
