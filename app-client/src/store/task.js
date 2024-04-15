import config from '../../config.yaml'

const state = {
    taskSettings: JSON.parse(localStorage.getItem('taskSettings')) || null,
    taskAlts: JSON.parse(localStorage.getItem('taskAlts')) || null,
    taskCriteria: JSON.parse(localStorage.getItem('taskCriteria')) || null,
    taskMatrix: JSON.parse((localStorage.getItem('taskMatrix'))) || null,
    ord: JSON.parse((localStorage.getItem('ord'))) || 0,
    final: JSON.parse(localStorage.getItem('final')) || null,
};

const mutations = {
    setTaskSettings(state, task) {
        state.taskSettings = task
    },
    setAlts(state, alts) {
        state.taskAlts = alts
    },
    setCriteria(state, criteria) {
        state.taskCriteria = criteria
    },
    setMatrix(state, matrix) {
        state.taskMatrix = matrix
    },
    setOrd(state, ord) {
        state.ord = ord;
    },
    setFinal(state, final) {
        state.final = final;
    }
};

const actions = {
    async showTask({ commit, getters }, sid) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/solution/settings?sid=' + sid, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            }
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                    commit('setTaskSettings', null);
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    commit('setTaskSettings', data);
                }
            })
            .catch(error => {
                console.log(error)
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    async updateTask({ commit, dispatch, getters }, task) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/solution/settings?sid=' + task.sid, {
            method: "PUT",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            },
            body: JSON.stringify(task),
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                    commit('setTaskSettings', null);
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    commit('setTaskSettings', task);
                }
            })
            .catch(error => {
                console.log(error)
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    async showAlts({ commit, getters }, sid) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/solution/alternatives?sid=' + sid, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            }
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                    commit('setAlts', null);
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    commit('setAlts', data);
                }
            })
            .catch(error => {
                console.log(error)
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    async updateAlts({ commit, getters }, {sid, alts}) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/solution/alternatives?sid=' + sid, {
            method: "PUT",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            },
            body: JSON.stringify(alts),
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                    commit('setAlts', null);
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    commit('setAlts', alts);
                }
            })
            .catch(error => {
                console.log(error)
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    async showCriteria({ commit, getters }, sid) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/solution/criteria?sid=' + sid, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            }
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                    commit('setCriteria', null);
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    commit('setCriteria', data);
                }
            })
            .catch(error => {
                console.log(error)
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    async updateCriteria({ commit, getters }, {sid, criteria}) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/solution/criteria?sid=' + sid, {
            method: "PUT",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            },
            body: JSON.stringify(criteria),
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                    commit('setCriteria', null);
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    commit('setCriteria', criteria);
                }
            })
            .catch(error => {
                console.log(error)
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    async createMatrix({ commit, dispatch, getters }, sid) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/solution/rating?sid=' + sid, {
            method: 'POST',
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            }
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(async data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                    commit('setMatrix', null);
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    commit('setOrd', 0);
                    await dispatch('takeMatrix', {sid: sid, ord: 0})
                }
            })
            .catch(error => {
                console.log(error)
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    async takeMatrix({ commit, dispatch, getters }, { sid, ord }) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/solution/rating?sid=' + sid + "&ord=" + ord, {
            method: 'GET',
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            }
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(async data => {
                console.log(data);
                if (status === 500) {
                    await dispatch('createMatrix', sid);
                } else if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                    commit('setMatrix', null);
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    commit('setMatrix', data.ratings);
                }
            })
            .catch(error => {
                console.log(error)
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    async updateMatrix({ commit, dispatch, getters }, { sid, ord, ratings }) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/solution/rating?sid=' + sid + "&ord=" + ord, {
            method: 'PUT',
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ratings: ratings}),
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(async data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                    commit('setMatrix', null);
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                }
            })
            .catch(error => {
                console.log(error)
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    async completeTask({ commit, dispatch, getters }, sid) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + '/solution/rating/complete?sid=' + sid, {
            method: 'PATCH',
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            }
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(async data => {
                console.log(data);
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                    commit('setMatrix', null);
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                }
            })
            .catch(error => {
                console.log(error)
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    },
    changeOrd({commit}, newOrd) {
        commit('setOrd', newOrd);
    },
    async getRole({ commit, getters }, sid) {
        let status;
        let role;
        const token = getters['getToken'];
        await fetch(config.backend.url + "/solution/experts/role?sid=" + sid, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            }
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(data => {
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    role = data.message;
                }
            })
            .catch(error => {
                console.log(error);
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            })
        return role;
    },
    async getExperts({ commit, getters }, sid) {
        let status;
        let experts;
        const token = getters['getToken'];
        await fetch(config.backend.url + "/solution/experts?sid=" + sid, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            }
        })
            .then(response => {
                status = response.status;
                return response.json()
            })
            .then(data => {
                if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    experts = data;
                }
            })
            .catch(error => {
                console.log(error);
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            })
        return experts;
    },
    async setExpertsWeights({ commit, getters }, {sid, weights}) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + "/solution/experts?sid=" + sid, {
            method: "PUT",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            },
            body: JSON.stringify(weights),
        })
            .then(response => {
                status = response.status;
                return response.json();
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
                commit('setStatus', 500, {root: true});
            })
    },
    async takeFinal({ commit, getters }, {sid, threshold}) {
        let status;
        const token = getters['getToken'];
        await fetch(config.backend.url + "/solution/final?sid=" + sid, {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            },
            body: JSON.stringify(threshold),
        })
            .then(response => {
                status = response.status;
                return response.json();
            })
            .then(data => {
                console.log(data);
                if (status === 400) {
                    commit('setError', "doesn't complete", {root: true});
                    commit('setStatus', 400, {root: true});
                } else if (status !== 200) {
                    commit('setError', data.message, {root: true});
                    commit('setStatus', status, {root: true});
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    commit('setFinal', data);
                }
            })
            .catch(error => {
                console.log(error);
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            })
    }
};

const getters = {
    getTaskSettings(state) {
        return state.taskSettings
    },
    getAlts(state) {
        return state.taskAlts
    },
    getCriteria(state) {
        return state.taskCriteria
    },
    getRatings(state) {
        return state.taskMatrix
    },
    getOrd(state) {
        return state.ord
    },
    getFinal(state) {
        return state.final
    }
};

export default {
    state,
    mutations,
    actions,
    getters
}