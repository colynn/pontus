<template>
  <div class="app-container">
    <el-row :gutter="20">
      <el-col :span="18" :xs="24">
        <el-card>
          <div slot="header" class="clearfix">
            <span>个人信息</span>
          </div>
          <el-tabs v-model="activeTab">
            <el-tab-pane label="基本资料" name="userinfo">
              <userInfo :user="user" />
            </el-tab-pane>
            <!-- <el-tab-pane label="修改密码" name="resetPwd">
              <resetPwd :user="user" />
            </el-tab-pane> -->
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import userInfo from './userInfo'
import { getUserProfile } from '@/api/user'

export default {
  name: 'Profile',
  components: { userInfo },
  data() {
    return {
      user: {},
      roleGroup: {},
      postGroup: {},
      activeTab: 'userinfo'
    }
  },
  created() {
    this.getUser()
  },
  methods: {
    getUser() {
      getUserProfile().then(response => {
        this.user = response.data
        this.roleGroup = response.roleGroup
        this.postGroup = response.postGroup
      })
    }
  }
}
</script>
