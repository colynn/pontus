<template>
  <el-form ref="form" :model="user" :rules="rules" label-width="80px">
    <el-form-item label="用户名" prop="realName">
      <el-input v-model="user.realName" />
    </el-form-item>
    <el-form-item label="手机号码" prop="phone">
      <el-input v-model="user.phone" maxlength="11" />
    </el-form-item>
    <el-form-item label="邮箱" prop="email">
      <el-input v-model="user.email" maxlength="50" />
    </el-form-item>
    <!-- <el-form-item label="部门">
      <el-input v-model="user.dept" maxlength="50" />
    </el-form-item> -->
    <el-form-item>
      <el-button type="primary" size="mini" @click="submit">更新</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import { updateUserProfile } from '@/api/user'

export default {
  props: {
    // eslint-disable-next-line vue/require-default-prop
    user: { type: Object }
  },
  data() {
    return {
      // 表单校验
      rules: {
        realName: [
          { required: true, message: '用户名不能为空', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '邮箱地址不能为空', trigger: 'blur' },
          {
            type: 'email',
            message: "'请输入正确的邮箱地址",
            trigger: ['blur', 'change']
          }
        ],
        phone: [
          { required: false, message: '手机号码不能为空', trigger: 'blur' },
          {
            pattern: /^1[3|4|5|6|7|8|9][0-9]\d{8}$/,
            message: '请输入正确的手机号码',
            trigger: 'blur'
          }
        ]
      }
    }
  },
  methods: {
    submit() {
      this.$refs['form'].validate(valid => {
        if (valid) {
          updateUserProfile(this.user).then(response => {
            if (response.code === 200) {
              this.msgSuccess('修改成功')
            } else {
              this.msgError(response.msg)
            }
          })
        }
      })
    }
  }
}
</script>
