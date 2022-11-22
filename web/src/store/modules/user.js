import { login, logout, getInfo, refreshtoken } from '@/api/user'
import { getToken, setToken, removeToken } from '@/utils/auth'
import router, { resetRouter } from '@/router'

const getDefaultState = () => {
  return {
    token: getToken(),
    name: '',
    avatar: '',
    roles: [],
    permissions: [],
    permisaction: []
  }
}

const state = getDefaultState()

const mutations = {
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_NAME: (state, name) => {
    state.name = name
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  },
  SET_PERMISSIONS: (state, permisaction) => {
    state.permisaction = permisaction
  }
}

const actions = {
  // user login
  login({ commit }, userInfo) {
    return new Promise((resolve, reject) => {
      login(userInfo).then(response => {
        const { token } = response
        commit('SET_TOKEN', token)
        setToken(token)
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // get user info
  getInfo({ commit }) {
    return new Promise((resolve, reject) => {
      getInfo().then(response => {
        if (!response || !response.data) {
          commit('SET_TOKEN', '')
          removeToken()
          resolve()
        }

        const { username, realName, avatar, roles, permissions } = response.data

        // roles must be a non-empty array
        if (!roles || roles.length <= 0) {
          reject('getInfo: roles must be a non-null array!')
        }
        // TODO: permissions, roles
        commit('SET_PERMISSIONS', permissions)
        commit('SET_ROLES', roles)
        if (realName === '') {
          commit('SET_NAME', username)
        } else {
          commit('SET_NAME', realName)
        }
        commit('SET_AVATAR', avatar)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // user logout
  logout({ commit, state }) {
    return new Promise((resolve, reject) => {
      logout(state.token).then(() => {
        removeToken() // must remove  token  first
        commit('SET_ROLES', [])
        commit('SET_PERMISSIONS', [])
        commit('SET_ROUTES', [])
        resetRouter()
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // 刷新token
  refreshToken({ commit, state }) {
    return new Promise((resolve, reject) => {
      refreshtoken(state.token).then(response => {
        const { token } = response
        commit('SET_TOKEN', token)
        setToken(token)
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // remove token
  resetToken({ commit }) {
    return new Promise(resolve => {
      removeToken() // must remove  token  first
      commit('RESET_STATE')
      resolve()
    })
  },
  // dynamically modify permissions
  changeRoles({ commit, dispatch }, role) {
    return new Promise(async resolve => {
      await dispatch('getInfo')

      resetRouter()

      // generate accessible routes map based on roles
      const accessRoutes = await dispatch('permission/generateRoutes')

      // dynamically add accessible routes
      router.addRoutes(accessRoutes)

      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
