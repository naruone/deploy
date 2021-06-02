import router from '@/router/index'
import {login} from '@/api/user'

export const user = {
    namespaced: true,
    state: {
        userInfo: {
            user_id: "",
            user_name: "",
            nick_name: "",
            email: "",
            header_img: "",
            authority: "",
        },
        token: "",
        expiresAt: ""
    },
    mutations: {
        setUserInfo(state, userInfo) {
            state.userInfo = userInfo
        },
        setToken(state, token) {
            state.token = token
        },
        setExpiresAt(state, expiresAt) {
            state.expiresAt = expiresAt
        },
        LoginOut(state) {
            state.userInfo = {}
            state.token = ""
            state.expiresAt = ""
        },
    },
    actions: {
        async LoginIn({commit}, loginInfo) {
            const res = await login(loginInfo)
            commit('setUserInfo', res.data.user)
            commit('setToken', res.data.token)
            commit('setExpiresAt', res.data.expiresAt)
            let redirect = router.history.current.query.redirect
            if (redirect === undefined) {
                redirect = '/sys';
            }
            await router.push({path: String(redirect)})
        },
        async LoginOut({commit}, redirect = true) {
            commit("LoginOut")
            if (redirect) await router.push({name: 'login', replace: true})
        }
    },
    getters: {
        userInfo(state) {
            return state.userInfo
        },
        token(state) {
            return state.token
        },
        expiresAt(state) {
            return state.expiresAt
        }
    }
}
