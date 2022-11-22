const getters = {
  sidebar: state => state.app.sidebar,
  device: state => state.app.device,
  token: state => state.user.token,
  avatar: state => state.user.avatar,
  name: state => state.user.name,
  roles: state => state.user.roles,
  permisaction: state => state.user.permisaction,
  permission_routes: state => state.permission.routes,
  snippet: state => state.snippet.snippet
}
export default getters
