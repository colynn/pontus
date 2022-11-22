<template>
  <div class="app-container">
    <div>
      <div style="padding-top:20px;">
        <el-row class="row">
          <el-col :span="4" class="info-title">
            物理机信息
          </el-col>
          <el-col :span="3" :xl="2" style="float:right">
            <el-button v-permisaction="['asset:pc:edit']" size="mini" type="primary" @click.native="showEdit">编辑</el-button>
            <el-popconfirm
              title="你确认删除吗？"
              @onConfirm="deletePcItem()"
            >
              <el-button slot="reference" v-permisaction="['asset:pc:remove']" size="mini" type="danger" style="margin-left: 10px">删除</el-button>
            </el-popconfirm>
          </el-col>
        </el-row>
        <el-row class="row content-font-size">
          <el-col :span="7" class="info-item">资产标识： <span class="value">{{ form.InstanceID }}</span></el-col>
          <el-col :span="5" class="info-item">序列号：<span class="value">{{ form.SerialNumber }}</span></el-col>
          <el-col :span="6" class="info-item">品牌： {{ form.Manufactory }}</el-col>
          <el-col :span="6" class="info-item">操作系统： {{ form.OSName }}</el-col>
          <el-col :span="7" class="info-item">CPU： {{ form.CPU }} CPU</el-col>
          <el-col :span="5" class="info-item">内存: {{ form.Memory }}G</el-col>
          <el-col :span="6" class="info-item">硬盘: {{ form.Disk }}G</el-col>

          <el-col :span="7" class="info-item">地区： {{ form.Region }}</el-col>
          <!-- <el-col :span="5" class="info-item">状态： {{ form.StatusValue }}</el-col> -->
          <el-col :span="5" class="info-item">领用人： {{ form.Recipient }}</el-col>
          <el-col :span="6" class="info-item">一级部门： {{ form.Department }}</el-col>
          <el-col :span="6" class="info-item">二级部门： {{ form.SecondDept }}</el-col>
          <el-col :span="7" class="info-item">采购类型：<span>{{ form.ProcurementType }}</span></el-col>
          <el-col :span="5" class="info-item">入库日期： <span>{{ parseDate(form.DeliveryDate) }}</span></el-col>
          <el-col :span="6" class="info-item">盘点时间：<span>{{ parseDate(form.InventoryTime) }}</span></el-col>
          <el-col :span="6" class="info-item">发票号码： <span>{{ form.InvoiceNumber }}</span></el-col>
          <el-col :span="7" class="info-item">发票日期：<span>{{ parseTime(form.InvoiceTime) }}</span></el-col>
          <el-col :span="5" class="info-item">税前金额：<span>{{ form.PretaxAmount }}</span></el-col>
          <el-col :span="6" class="info-item">总金额：<span>{{ form.PretaxGrossAmount }}</span></el-col>
          <el-col :span="6" class="info-item">保修期：<span>{{ form.WarrantyPeriod }}月</span></el-col>
          <el-col :span="24" class="info-item">备注：{{ form.Description }}</el-col>
        </el-row>
      </div>

    </div>

    <div v-show="true" class="portlet-body">
      <div class="setTitle">变更历史</div>
      <div class="content-font-size">
        <template v-for="(item) in logList">
          <el-row v-if="item.type == 1" :key="item.id" class="row bottom-line">{{ item.username }} {{ item.content }} - {{ item.created_at }}</el-row>
          <el-row v-else :key="item.id" class="row bottom-line">
            <el-row class="row">{{ item.username }} 进行了改变 - {{ parseTime(item.created_at) }}</el-row>
            <el-table border :data="item.content">
              <el-table-column prop="Field" label="字段" min-width="8%" />
              <el-table-column prop="OriginValue" label="原值" min-width="20%">
                <template slot-scope="scope">
                  <span>{{ scope.row.OriginValue }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="NewValue" label="新值" min-width="8%">
                <template slot-scope="scope">
                  <span> {{ scope.row.NewValue }}</span>
                </template>
              </el-table-column>
            </el-table>
          </el-row>
        </template>
      </div>
      <Pagination
        v-show="total>0"
        :total="total"
        :page.sync="queryParams.PageIndex"
        :limit.sync="queryParams.PageSize"
        @pagination="getAuditLog"
      />
    </div>

    <!-- 编辑对话框 -->
    <el-dialog title="编辑信息" :visible.sync="open" width="800px">
      <el-form ref="form" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="资产标识" prop="InstanceID">
          <el-input v-model="form.InstanceID" placeholder="请输入资产标识" :disabled="isEdit" />
        </el-form-item>

        <el-form-item label="品牌" prop="Manufactory">
          <el-input v-model="form.Manufactory" placeholder="请输入电脑品牌" />
        </el-form-item>

        <el-form-item label="操作系统" prop="OSName">
          <el-input v-model="form.OSName" placeholder="" />
        </el-form-item>

        <el-form-item label="配置" prop="Configuration">
          <el-input v-model="form.Configuration" placeholder="配置" />
        </el-form-item>

        <el-form-item label="领用人" prop="Useranme" :span="12">
          <el-select v-model="form.RecipientID" clearable filterable placeholder="请选择" class="el-input--small">
            <el-option
              v-for="item in userList"
              :key="item.userId"
              :label="item.realName == '' ? item.username:item.realName"
              :value="item.userId"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="地区" prop="Region">
          <el-input v-model="form.Region" placeholder="上海/北京" />
        </el-form-item>

        <el-form-item label="一级部门" prop="Department">
          <el-select v-model="form.primaryDept" clearable filterable placeholder="请选择" class="el-input--small">
            <el-option
              v-for="item in primaryDeptList"
              :key="item.deptId"
              :label="item.deptName"
              :value="item.deptId"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="采购类型" prop="ProcurementType">
          <el-input v-model="form.ProcurementType" placeholder="请输入" />
        </el-form-item>

        <el-form-item label="入库日期" prop="DeliveryDate">
          <el-date-picker
            v-model="form.DeliveryDate"
            value-format="yyyy-MM-dd"
            type="date"
            placeholder="选择日期"
          />
        </el-form-item>

        <el-form-item label="盘点时间" prop="InventoryTime">
          <el-date-picker
            v-model="form.InventoryTime"
            value-format="yyyy-MM-dd"
            type="date"
            placeholder="选择日期"
          />
        </el-form-item>

        <el-form-item label="发票号码" prop="InvoiceNumber">
          <el-input v-model="form.InvoiceNumber" placeholder="请输入" />
        </el-form-item>

        <el-form-item label="发票日期" prop="InvoiceTime">
          <el-date-picker
            v-model="form.InvoiceTime"
            value-format="yyyy-MM-dd HH:mm:ss"
            type="datetime"
            placeholder="选择日期"
          />
        </el-form-item>

        <el-form-item label="税前金额" prop="PretaxAmount">
          <el-input-number v-model="form.PretaxAmount" placeholder="请输入" />
        </el-form-item>

        <el-form-item label="总金额" prop="PretaxGrossAmount">
          <el-input-number v-model="form.PretaxGrossAmount" placeholder="请输入" />
        </el-form-item>

        <el-form-item label="保修期" prop="WarrantyPeriod">
          <el-input-number v-model="form.WarrantyPeriod" />
        </el-form-item>

        <el-form-item label="备注" prop="Description">
          <el-input v-model="form.Description" type="textarea" placeholder="请输入内容" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="cancel">取 消</el-button>
        <el-button type="primary" @click="submitForm">确 定</el-button>
      </div>
    </el-dialog>

  </div>
</template>
<script>
import { pcDetail, pcUpdate, pcDelete, pcAuditLog } from '@/api/cmdb/tangible'
import { TangibleAssetStatus } from '@/utils/constant'
import { parseTime, parseDate } from '@/utils/utils'
import { listUser } from '@/api/user'
export default {

  components: { },
  data() {
    return {
      fileList: [],
      loading: true,
      total: 0,
      // 查询参数
      queryParams: {
        PageIndex: 1,
        PageSize: 10
      },

      instanceID: '',
      // 资产状态：
      TangibleAssetStatus: TangibleAssetStatus,
      status: '',

      projectInfo: {

      },
      open: false,
      hovered: false,
      // 表单参数
      isEdit: true,
      form: {},
      // 表单校验
      rules: {
        // TODO: rules update pending
        roleName: [
          { required: true, message: '角色名称不能为空', trigger: 'blur' }
        ],
        roleKey: [
          { required: true, message: '权限字符不能为空', trigger: 'blur' }
        ],
        roleSort: [
          { required: true, message: '角色顺序不能为空', trigger: 'blur' }
        ]
      },
      logList: [],
      // 日期范围
      timerange: null,
      pcInfo: {}, // 详情信息

      // user 查询参数
      userQueryParams: {
        pageIndex: 1,
        pageSize: 10000
      },
      userList: [],
      primaryDeptList: []
    }
  },
  created() {
    this.getPCDetail()
    this.getAuditLog()
  },

  methods: {
    getPCDetail() {
      this.instanceID = this.$route.params.instanceID
      pcDetail(this.instanceID).then(response => {
        this.form = response.data
        this.parseStatusValue()
        this.form.InvoiceTime = parseTime(this.form.InvoiceTime)
        this.form.DeliveryDate = parseDate(this.form.DeliveryDate)
        this.form.InventoryTime = parseDate(this.form.InventoryTime)
        this.loading = false
      })
    },
    // TODO: generate component
    getAuditLog() {
      const query = this.queryParams
      pcAuditLog(this.instanceID, query).then(response => {
        this.logList = response.data.items
        this.total = response.data.count
      })
    },
    getUserList() {
      this.loading = true
      listUser(this.userQueryParams).then(response => {
        this.userList = response.data.items
        this.loading = false
      }
      )
    },
    showEdit() {
      this.open = true
      this.getUserList()
    },
    cancel() {
      this.open = false
    },
    submitForm() {
      // TODO: add verify for data object
      const data = this.form
      pcUpdate(this.$route.params.instanceID, data).then(response => {
        this.form = response.data
        this.parseStatusValue()
        this.form.InvoiceTime = parseTime(this.form.InvoiceTime)
        this.form.DeliveryDate = parseDate(this.form.DeliveryDate)
        this.form.InventoryTime = parseDate(this.form.InventoryTime)
        this.$message.success('更新成功')
        this.open = false
        this.getAuditLog()
      })
    },
    deletePcItem() {
      pcDelete(this.$route.params.instanceID).then(response => {
        this.$message.success('删除成功')
        this.$router.push('/asset/list/pc/list')
      })
    },
    parseStatusValue() {
      switch (this.form.Status) {
        case 1:
          this.form.StatusValue = '使用中'
          break
        case 2:
          this.form.StatusValue = '闲置中'
          break
      }
    }
  }
}
</script>

<style scoped>
.info-title {
  font-size: 16px;
}
.info-item {
  line-height: 40px;
  color: #545353 !important;
}
.setTitle {
  padding-top: 20px;
  padding-bottom: 20px;
  font-size: 16px;
  border-top: 1px solid  #EBEEF5;
}

.value {
  padding: 5px;
}

/* .value:hover {
  border: 1px solid #ccc;
} */
</style>

