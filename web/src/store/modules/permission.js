import { asyncRoutes, constantRoutes } from '@/router'
import { getRoutes } from '@/api/system/role'
import Layout from '@/layout'
import { getToken } from '@/utils/auth'

/**
  * 后台查询的菜单数据拼装成路由格式的数据
  * @param routes
  */
export function generaMenu(routes, data) {
  data.forEach(item => {
    const menu = {
      path: item.path,
      redirect: item.redirect,
      component: item.component === 'Layout' ? Layout : loadView(item.component),
      hidden: item.visible !== '0',
      children: [],
      name: item.menuName,
      meta: {
        title: item.title,
        icon: item.icon,
        noCache: true
      }
    }
    if (item.children) {
      generaMenu(menu.children, item.children)
    }
    routes.push(menu)
  })
}

export const loadView = (view) => { // 路由懒加载
  return (resolve) => require(['@/views' + view], resolve)
  // return () => import(`@/views${view}`)
}

const state = {
  routes: [],
  addRoutes: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  }
}

const actions = {
  generateRoutes({ commit }) {
    return new Promise(resolve => {
      const loadMenuData = []
      getRoutes().then(response => {
        let data = response
        if (response.code !== 200) {
          this.$message({
            message: '菜单数据加载异常',
            type: 0
          })
        } else {
          data = response.data
          Object.assign(loadMenuData, data)
          generaMenu(asyncRoutes, loadMenuData)
          asyncRoutes.push({ path: '*', redirect: '/404', hidden: true })
          commit('SET_ROUTES', asyncRoutes)
          // console.log(JSON.stringify(asyncRoutes))
          // console.log(asyncRoutes)
          resolve(asyncRoutes)
        }
      }).catch(error => {
        console.log(error)
      })
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
