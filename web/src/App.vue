<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App',
  created() {
  },
  methods: {

  }
}
</script>

<style lang="scss">

.pagination {
  float: right;
}

/* 查询字段 样式定义 */
.query-name {
  padding-top: 10px;
  color: #595959;
  text-align: center;
  font-size: 14px;
}

.row {
  padding-bottom: 10px;
}

.bottom-line {
  border-bottom:  1px solid #EBEEF5;
  padding-top: 10px;
  padding-bottom: 10px;
}

.content-font-size {
  font-size: 15px;
}

.date-input {
  line-height: 20px;
}

.date-range {
  height: 32px;
  line-height: 32px;
}

.span-icon {
  color: gray;
}

</style>
