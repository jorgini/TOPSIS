import config from '../../config.yaml'

const state = {
    token: JSON.parse(localStorage.getItem('token')) || null,
    exp_at: JSON.parse(localStorage.getItem('exp_at')) || null,
    userInfo: JSON.parse(localStorage.getItem('userInfo')) || null
};

const mutations = {
    setToken(state, token) {
        state.token = token;
        if (token)
            state.exp_at = Date.now() + (14 * 60 * 1000);
        else
            state.exp_at = Date.now() + (30 * 1000);
    },
    setUser(state, user) {
        state.userInfo = user
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
        commit('setUser', null);
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
        },
    async reqUserInfo({ commit, getters }) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/user/settings', {
            method: "GET",
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-type': 'application/json'
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
                    commit('setUser', {login: data.login, email: data.email, password: ""});
                }
            })
            .catch(error => {
                console.log(error);
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    async updateUser({ commit, getters }, user) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/user/settings', {
            method: "PATCH",
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-type': 'application/json'
            },
            body: JSON.stringify(user),
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
                }
            })
            .catch(error => {
                console.log(error);
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    async deleteUser({ commit, dispatch, getters }) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/user/settings', {
            method: "DELETE",
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-type': 'application/json'
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
                    commit('setUser', null);
                    dispatch('logout');
                }
            })
            .catch(error => {
                console.log(error);
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
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
    },
    getUser: (state) => {
        return state.userInfo
    }
};

export default {
    state,
    mutations,
    actions,
    getters
}
