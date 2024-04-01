const state = {
    component: JSON.parse(localStorage.getItem('component')) || 'Main',
    curTask: JSON.parse(localStorage.getItem('curTask')) || null,
};

const mutations = {
    setPage(state, component) {
        state.component = component
    }
};

const actions = {
    showPage({ commit }, component) {
        commit('setPage', component)
    }
};

const getters = {
    getPage(state) {
        return state.component
    }
};

export default {
    state,
    mutations,
    actions,
    getters
}