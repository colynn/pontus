<template>
  <div class="app-container">
    <el-row class="row">
      <el-col :span="2" class="query-name">一级部门</el-col>
      <el-col :span="3">
        <el-select v-model="primary" clearable placeholder="请选择" class="el-input--small">
          <el-option
            v-for="(item, index) in deptOptions"
            :key="index"
            :label="item.showName"
            :value="item.deptId"
          />
        </el-select>
      </el-col>
      <el-col :span="2" class="query-name">资产状态</el-col>
      <el-col :span="3">
        <el-select v-model="status" clearable placeholder="请选择" class="el-input--small" @change="updateItem">
          <el-option
            v-for="item in TangibleAssetStatus"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-col>
      <el-col :xl="2" :span="1" class="query-name" style="margin-left:10px; margin-right: 10px;">领用人</el-col>
      <el-col :span="3">
        <el-input v-model="username" clearable placeholder="请输入领用人" class="el-input--small" />
      </el-col>
      <el-col :span="2" class="query-name" style="margin-left: 20px;">入库日期</el-col>
      <el-col :span="4">
        <el-date-picker
          v-model="timerange"
          type="daterange"
          class="date-range"
          value-format="yyyy-MM-dd"
          unlink-panels
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          :picker-options="pickerOptions"
        />
      </el-col>
    </el-row>
    <el-row class="">
      <el-col :span="3" :sm="3" :xl="2">
        <el-button type="primary" class="el-button--mini" @click="handleQuery">搜索</el-button>
        <el-button type="warning el-button--mini" @click="resetQuery">重置</el-button>
        <!-- <el-button type="success" class="el-button--mini" @click="exportData">导出</el-button> -->
      </el-col>
      <el-col :span="3">
        <upload :actions="['asset:pc:import']" url="/api/v1/assets/pc/upload" />
      </el-col>
    </el-row>
    <el-table
      v-loading="loading"
      :data="tableData"
      style="width: 100%"
      :default-sort="{prop: 'date', order: 'descending'}"
    >
      <el-table-column prop="Manufactory" label="品牌" sortable width="100" />
      <el-table-column prop="Type" label="类型" sortable width="120" />
      <el-table-column prop="Configuration" show-overflow-tooltip label="配置" sortable />
      <el-table-column prop="StatusValue" label="状态" sortable width="90" />
      <el-table-column prop="Region" sortable label="地区" width="90" />
      <el-table-column prop="Department" show-overflow-tooltip label="归属部门" sortable width="240" />
      <el-table-column prop="Recipient" show-overflow-tooltip label="领用人" sortable width="90" />
      <el-table-column prop="ProcurementType" sortable label="采购类型" width="100" />
      <el-table-column label="操作" width="80">
        <template slot-scope="scope">
          <router-link :to="{name: 'pcDetail', params: {instanceID: scope.row.InstanceID} }">
            <el-button type="text" size="small">详情</el-button>
          </router-link>
        </template>
      </el-table-column>
    </el-table>
    <pagination
      v-show="total>0"
      :total="total"
      :page.sync="queryParams.PageIndex"
      :limit.sync="queryParams.PageSize"
      @pagination="getPCAssetList"
    />
  </div>
</template>

<script>
import { Message } from 'element-ui'
import { PCList } from '@/api/cmdb/tangible'
import { TangibleAssetStatus } from '@/utils/constant'
import Upload from '@/components/Upload'

export default {
  name: 'Pclist',
  components: { Upload },
  data() {
    return {
      fileList: [],
      total: 10,
      loading: true,
      // 查询参数
      queryParams: {
        PageIndex: 1,
        PageSize: 10
      },
      // 部门-选择项
      primaryDeptID: 0,
      secondDeptID: 0,
      deptOptions: [],
      primary: undefined,

      // 资产状态：
      TangibleAssetStatus: TangibleAssetStatus,
      status: undefined,

      // 领用人
      username: '',

      // 日期范围
      timerange: null,
      value1: '',
      pickerOptions: {
        shortcuts: [{
          text: '最近一周',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近一个月',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近三个月',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
            picker.$emit('pick', [start, end])
          }
        }]
      },
      tableData: []
    }
  },
  created() {
    this.getPCAssetList()
  },

  methods: {
    getPCAssetList() {
      this.queryParams.UserName = this.username
      this.queryParams.Status = this.status
      this.queryParams.PrimaryDeptID = this.primary
      if (this.timerange !== undefined && this.timerange !== null) {
        this.queryParams.StartAt = this.timerange[0]
        this.queryParams.EndAt = this.timerange[1]
      }
      PCList(this.queryParams).then(response => {
        this.tableData = response.data.items
        this.tableData.forEach((_, index) => {
          switch (this.tableData[index].Status) {
            case 1:
              this.tableData[index].StatusValue = '使用中'
              break
            case 2:
              this.tableData[index].StatusValue = '闲置中'
              break
          }
        })
        this.total = response.data.count
        this.loading = false
      })
    },
    /** 搜索按钮操作 */
    handleQuery() {
      this.getPCAssetList()
    },
    alterChosePrimary() {
      if (this.primary === undefined || this.primary === '') {
        Message({
          message: '请您先选择一级部门',
          type: 'info',
          duration: 2 * 1000
        })
        return
      }
    },
    updateItem() {
      if (this.status === '') {
        this.queryParams.Status = ''
      }
    },
    resetQuery() {
      this.queryParams = {
        PageIndex: 1,
        PageSize: 10
      }
      this.primary = undefined
      this.status = undefined
      this.timerange = undefined
      this.username = undefined
      this.getPCAssetList()
    },
    exportData() {
      Message({
        message: '即将上线',
        type: 'info',
        duration: 2 * 1000
      })
    }
  }
}
</script>

<style scoped>
</style>
