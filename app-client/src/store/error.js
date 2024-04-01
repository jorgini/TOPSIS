const state = {
    statusCode: JSON.parse(localStorage.getItem('statusCode')) || 200,
    error: JSON.parse(localStorage.getItem('error')) || null,
};

const mutations = {
    setError(state, error) {
        state.error = error
    },
    setStatus(state, status) {
        state. statusCode = status
    }
};

const actions = {
};

const getters = {
    errorOccurred: (state) => {
        return state.error
    },
    getStatusCode: (state) => {
        return state.statusCode
    }
};

export default {
    state,
    mutations,
    actions,
    getters
}