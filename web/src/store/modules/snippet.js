
const state = {
  snippet: []
}

const mutations = {
  snippetTag(state, vm) {
    console.log(vm)
    state.snippet.push({ 'title': vm.title, 'text': vm.text })
  },
  snippetTagToJson(state) {
    localStorage.setItem('snippet', JSON.stringify(state.snippet))
  },
  snippetTagFromJson(state) {
    state.snippet = JSON.parse(localStorage.getItem('snippet'))
  },
  snippetRemoveTag(state, vm) {
    const index = state.snippet.indexOf(vm)
    state.snippet.splice(index, 1)
  }
}

export default {
  namespaced: true,
  state,
  mutations
}
