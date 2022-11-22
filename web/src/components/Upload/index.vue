<template>
  <div>
    <el-upload
      style="width: 200px"
      class="upload"
      :action="uploadExcel()"
      :show-file-list="false"
      :on-success="uploadSuccess"
      :headers="uploadHeader()"
      :on-error="uploadError"
      :on-exceed="beforeUpload"
    >
      <el-button v-permisaction="actions" class="el-button--mini" type="info" icon="el-icon-upload2" style="float: left">导入</el-button>
    </el-upload>
  </div>

</template>

<script>
import { getToken } from '@/utils/auth'

export default {
  name: 'Upload',
  props: {
    actions: {
      type: Array,
      default() {
        return ['']
      }
    },
    url: {
      type: String,
      default: '/'
    }
  },
  data() {
    return {}
  },
  created() {
    console.log(this.actions)
  },
  methods: {
    uploadExcel() {
      return this.url
    },
    uploadHeader() {
      return {
        'Authorization': `Bearer ` + getToken()
      }
    },
    /* 上传http状态200,校验 */
    uploadSuccess(res) {
      if (typeof res === 'object') {
        if (res.code === 200) {
          this.$message.success('导入成功')
          return
        } else if (res.code === 401) {
          this.$message.info('请刷新页面后重试')
          return
        }
      }
      this.$message.error('导入失败-请联系管理员')
    },
    /* 上传导入失败 */
    uploadError() {
      this.$message.error('导入失败-请联系管理员')
    },
    beforeUpload(file) {
      const isLt1M = file.size / 1024 / 1024 < 1
      console.log('file.szie: ' + file.size + 'isLt1M: ' + isLt1M)
      if (isLt1M) {
        return true
      }
      this.$message({
        message: 'Please do not upload files larger than 1m in size.',
        type: 'warning'
      })
      return false
    }
  }
}
</script>
