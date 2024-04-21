import config from '../../config.yaml'

const state = {
    defaultLingNumScale: JSON.parse(localStorage.getItem('defaultLingNumScale')) || null,
    defaultLingIntScale: JSON.parse(localStorage.getItem('defaultLingIntScale')) || null,
    defaultLingT1FScale: JSON.parse(localStorage.getItem('defaultLingT1FScale')) || null,
};

const mutations = {
    setDefaultLingNumScale(state, scale) {
        state.defaultLingNumScale = scale
    },
    setDefaultLingIntScale(state, scale) {
        state.defaultLingIntScale = scale
    },
    setDefaultLingT1FScale(state, scale) {
        state.defaultLingT1FScale = scale
    }
};

const actions = {
    async takeDefaultLingScale({ commit }) {
        let status;
        await fetch(config.backend.url + '/solution/defaults', {
            method: "GET",
            headers: {
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
                    commit('setDefaultLingNumScale', null);
                    commit('setDefaultLingIntScale', null);
                    commit('setDefaultLingT1FScale', null);
                } else {
                    commit('setError', null, {root: true});
                    commit('setStatus', 200, {root: true});
                    commit('setDefaultLingNumScale', data.number);
                    commit('setDefaultLingIntScale', data.interval);
                    commit('setDefaultLingT1FScale', data.t1fs);
                }
            })
            .catch(error => {
                console.log(error)
                commit('setError', error, {root: true});
                commit('setStatus', 500, {root: true});
            });
    }
};

const getters = {
    getDefaultLingT1FScale(state) {
        return state.defaultLingT1FScale
    },
    getDefaultLingIntScale(state) {
        return state.defaultLingIntScale
    },
    getDefaultLingNumScale(state) {
        return state.defaultLingNumScale
    }
};

export default {
    state,
    mutations,
    actions,
    getters
}