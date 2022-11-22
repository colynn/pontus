<!-- thanks by wangjianhui2464-->
<template>
  <div style="width: 100%;" :style="{height:height}" />
</template>

<script>
require(['emmet/emmet'], function(data) {
  window.emmet = data.emmet
})
const ace = require('brace')
export default {
  name: 'Editor',
  props: {
    value: {
      type: String,
      required: true
    },
    isRead: {
      type: Boolean,
      default: false
    },
    height: {
      type: String,
      default: '300px'
    }
  },
  data() {
    return {
      editor: null,
      contentBackup: ''
    }
  },
  watch: {
    value(val) {
      if (this.contentBackup !== val) {
        this.editor.setValue(val, 1)
      }
    },
    theme: function(newTheme) {
      this.editor.setTheme('ace/theme/' + newTheme)
    },
    lang: function(newLang) {
      this.editor.getSession().setMode('ace/mode/' + newLang)
    }
  },
  mounted() {
    const vm = this
    require('brace/ext/emmet')
    require('brace/ext/language_tools')
    // require('brace/mode/mysql')
    // require('brace/theme/xcode')
    const editor = vm.editor = ace.edit(this.$el)
    this.$emit('init', editor)
    const staticWordCompleter = {
      getCompletions: function(editor, session, pos, prefix, callback) {
        vm.$emit('setCompletions', editor, session, pos, prefix, callback)
      }
    }
    editor.completers = [staticWordCompleter]
    editor.setOptions({
      enableBasicAutocompletion: true,
      enableLiveAutocompletion: true
    })
    editor.$blockScrolling = Infinity
    editor.setFontSize(14)
    editor.setOption('enableEmmet', true)
    editor.getSession().setMode('ace/mode/mysql')
    editor.setTheme('ace/theme/xcode')
    editor.setValue(this.value, 1)
    editor.setReadOnly(this.is_read)
    editor.on('change', function() {
      const content = editor.getValue()
      vm.$emit('input', content)
      vm.contentBackup = content
    })
  }
}
</script>
