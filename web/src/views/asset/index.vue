<template>
  <div class="dashboard-container">
    <div class="headerNotice">
      您好，{{ name }}
    </div>

    <el-row>
      <div class="headerNotice">系统更新日志</div>
    </el-row>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
  computed: {
    ...mapGetters([
      'name'
    ])
  }
}
</script>

<style lang="scss" scoped>

.headerNotice {
  color: #595959;
  padding: 20px;
}
</style>

