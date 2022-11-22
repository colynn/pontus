<template>
  <div class="app-container">
    <el-row :gutter="20">
      <!--用户数据-->
      <el-col>
        <el-form ref="queryForm" :model="queryParams" :inline="true" label-width="68px">
          <el-form-item label="帐号" prop="username">
            <el-input
              v-model="queryParams.username"
              placeholder="请输入帐号"
              clearable
              size="small"
              style="width: 240px"
              @keyup.enter.native="handleQuery"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="el-icon-search" size="mini" @click="handleQuery">搜索</el-button>
            <el-button icon="el-icon-refresh" size="mini" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>

        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button
              v-permisaction="['system:user:add']"
              type="primary"
              icon="el-icon-plus"
              size="mini"
              @click="handleAdd"
            >新增</el-button>
          </el-col>
          <!-- <el-col :span="1.5">
            <el-button

              type="success"
              icon="el-icon-edit"
              size="mini"
              :disabled="single"
              @click="handleUpdate"
            >修改</el-button>
          </el-col> -->
          <el-col :span="3">
            <upload :actions="['system:user:import']" url="/api/v1/user/upload" />
          </el-col>
          <!--
          <el-col :span="1.5">
            <el-button
              v-permisaction="['system:user:export']"
              type="warning"
              icon="el-icon-download"
              size="mini"
              @click="handleExport"
            >导出</el-button>
          </el-col> -->
        </el-row>

        <el-table
          v-loading="loading"
          :data="userData"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="45" align="center" />
          <el-table-column label="帐号" align="center" prop="username" :show-overflow-tooltip="true" />
          <el-table-column label="用户名" align="center" prop="realName" :show-overflow-tooltip="true" />
          <el-table-column label="邮箱" align="center" prop="email" width="180" :show-overflow-tooltip="true" />
          <!-- <el-table-column label="角色" align="center" prop="role" width="180" :show-overflow-tooltip="true" /> -->
          <el-table-column label="角色" align="center" prop="roleName" :show-overflow-tooltip="true" />
          <el-table-column label="创建时间" align="center" prop="createdAt" width="160">
            <template slot-scope="scope">
              <span>{{ parseTime(scope.row.createdAt) }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            align="center"
            width="220"
            class-name="small-padding fixed-width"
          >
            <template slot-scope="scope">
              <el-button
                v-permisaction="['system:user:edit']"
                size="mini"
                type="text"
                icon="el-icon-edit"
                @click="handleUpdate(scope.row)"
              >修改</el-button>
              <el-button
                v-if="scope.row.userId !== 1"
                v-permisaction="['system:user:remove']"
                size="mini"
                type="text"
                icon="el-icon-delete"
                @click="handleDelete(scope.row)"
              >删除</el-button>
              <!-- <el-button
                size="mini"
                type="text"
                icon="el-icon-key"
                @click="handleResetPwd(scope.row)"
              >重置</el-button> -->
            </template>
          </el-table-column>
        </el-table>

        <Pagination
          v-show="total>0"
          :total="total"
          :page.sync="queryParams.pageIndex"
          :limit.sync="queryParams.pageSize"
          @pagination="getList"
        />
      </el-col>
    </el-row>

    <!-- 添加或修改参数配置对话框 -->
    <el-dialog :title="title" :visible.sync="open" width="600px">
      <el-form ref="form" :model="form" :rules="rules" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户名" prop="realName">
              <el-input v-model="form.realName" placeholder="请输入用户名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号码" prop="phone">
              <el-input v-model="form.phone" placeholder="请输入手机号码" maxlength="11" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" maxlength="50" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="帐号" prop="username">
              <el-input v-model="form.username" placeholder="请输入用户名" :disabled="canNotEdit" />
            </el-form-item>
          </el-col>
          <!-- <el-col :span="12">
            <el-form-item v-if="form.userId == undefined" label="用户密码" prop="password">
              <el-input v-model="form.password" placeholder="请输入用户密码" type="password" />
            </el-form-item>
          </el-col> -->
          <el-col :span="24">
            <el-form-item label="角色">
              <!-- <el-select v-model="form.roleIDs" multiple placeholder="请选择" style="width: 100%" @change="$forceUpdate(); addUserRoleRule()"> -->
              <el-select v-model="form.roleID" placeholder="请选择" style="width: 100%" @change="$forceUpdate(); addUserRoleRule()">
                <el-option
                  v-for="item in roleOptions"
                  :key="item.roleID"
                  :label="item.roleName"
                  :value="item.roleID"
                  :disabled="item.status == 1"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input v-model="form.remark" type="textarea" placeholder="请输入内容" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="cancel">取 消</el-button>
        <el-button type="primary" @click="submitForm">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { listUser, getUser, delUser, addUser, updateUser } from '@/api/user'
import { listAllRole } from '@/api/system/role'
import Pagination from '@/components/Pagination'
import Upload from '@/components/Upload'

export default {
  name: 'User',
  components: { Pagination, Upload },
  data() {
    return {
      // 遮罩层
      loading: true,
      // 选中数组
      ids: [],
      // 非单个禁用
      single: true,
      // 非多个禁用
      multiple: true,
      // 总条数
      total: 0,
      // 用户表格数据
      userList: null,

      userData: [],
      // 弹出层标题
      title: '',
      // 部门树选项
      deptOptions: undefined,
      // 是否显示弹出层
      open: false,
      // 用户名是否允许编辑
      canNotEdit: true,

      // 部门名称
      deptName: undefined,
      // 默认密码
      initPassword: undefined,
      // 日期范围
      dateRange: [],
      // 状态数据字典
      statusOptions: [],
      // 性别状态字典
      sexOptions: [],
      // 岗位选项
      postOptions: [],
      // 角色选项
      roleOptions: [],

      // 表单参数
      form: {
      },
      allowToSQLReviewer: false,
      defaultProps: {
        children: 'children',
        label: 'deptName'
      },
      // 查询参数
      queryParams: {
        pageIndex: 1,
        pageSize: 10,
        username: undefined,
        phone: undefined,
        status: undefined,
        deptId: undefined
      },
      // 表单校验
      rules: {
        username: [
          { required: true, message: '用户名不能为空', trigger: 'blur' }
        ],
        roleID: [
          { required: true }
        ],
        realName: [
          { required: false, message: '用户昵称不能为空', trigger: 'blur' }
        ],
        deptId: [
          { required: false, message: '归属部门不能为空', trigger: 'blur' }
        ],
        password: [
          { required: false, message: '用户密码不能为空', trigger: 'blur' }
        ],
        email: [
          { required: false, message: '邮箱地址不能为空', trigger: 'blur' },
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
  watch: {
    // 根据名称筛选部门树
    deptName(val) {
      this.$refs.tree.filter(val)
    }
  },
  created() {
    this.getList()
  },
  methods: {
    /** 查询用户列表 */
    getList() {
      this.loading = true
      listUser(this.addDateRange(this.queryParams, this.dateRange)).then(response => {
        this.userList = response.data.items
        listAllRole().then(res => {
          this.roleOptions = res.data
          this.userList.forEach((item, key) => {
            for (let i = 0; i < this.roleOptions.length; i++) {
              if (item.roleId === this.roleOptions[i].roleID) {
                this.userList[key].roleName = this.roleOptions[i].roleName
                break
              }
            }
          })
          this.userData = this.userList
        })

        this.total = response.data.count
        this.loading = false
      }
      )
    },
    // TODO: current for sql audit use, later will be combined.
    addUserRoleRule() {
      if (this.form.roleID !== undefined) {
        this.roleOptions.forEach((item) => {
          if (item.roleID === this.form.roleID) {
            this.form.rule = item.roleKey
          }
        })
      }
    },
    // 取消按钮
    cancel() {
      this.open = false
      this.reset()
    },
    // 表单重置
    reset() {
      this.form = {
        userId: undefined,
        roleID: undefined,
        rule: undefined,
        deptId: undefined,
        username: undefined,
        realName: undefined,
        phone: undefined,
        email: undefined,
        status: '0',
        remark: undefined,
        postIds: undefined,
        roleIds: undefined
      }
      this.allowToSQLReviewer = false
      this.resetForm('form')
    },
    /** 搜索按钮操作 */
    handleQuery() {
      this.queryParams.page = 1
      this.getList()
    },
    /** 重置按钮操作 */
    resetQuery() {
      this.dateRange = []
      this.resetForm('queryForm')
      this.handleQuery()
    },
    // 多选框选中数据
    handleSelectionChange(selection) {
      this.ids = selection.map(item => item.userId)
      this.single = selection.length !== 1
      this.multiple = !selection.length
    },
    /** 新增按钮操作 */
    handleAdd() {
      this.reset()
      listAllRole().then(response => {
        // this.postOptions = response.data.posts
        this.roleOptions = response.data
        this.open = true
        this.canNotEdit = false
        this.title = '添加用户'
        this.form.password = this.initPassword
        this.form.roleID = undefined
      })
    },
    /** 修改按钮操作 */
    handleUpdate(row) {
      this.reset()
      const userId = row.userId || this.ids
      getUser(userId).then(response => {
        this.form = response.data
        this.postOptions = response.posts
        this.roleOptions = response.roles
        this.form.roleID = response.roleIds[0]
        if (this.form.rule === 'admin') {
          this.allowToSQLReviewer = true
        }
        this.open = true
        this.canNotEdit = true
        this.title = '修改用户'
        this.form.password = ''
        this.form.originRoleID = response.roleIds[0]
      })
    },
    /** 提交按钮 */
    submitForm: function() {
      this.$refs['form'].validate(valid => {
        if (valid) {
          if (this.allowToSQLReviewer === true) {
            this.form.rule = 'admin'
          }
          if (this.form.userId !== undefined) {
            updateUser(this.form.userId, this.form).then(response => {
              if (response.code === 200) {
                // if update role, need update user token
                this.triggerUpdateUserRole()
                this.open = false
                this.getList()
              } else {
                this.msgError(response.msg)
              }
            })
          } else {
            if (this.form.roleID === undefined) {
              this.$message.warning('请选择角色后重试')
              return
            }
            addUser(this.form).then(response => {
              if (response.code === 200) {
                this.msgSuccess('新增成功')
                this.open = false
                this.getList()
              } else {
                this.msgError(response.msg)
              }
            })
          }
        }
      })
    },
    triggerUpdateUserRole() {
      if (this.form.roleID !== this.form.originRoleID) {
        // this.$store.dispatch('user/changeRoles').then(() => { this.$message.success('角色更新成功') })
        this.$message.success('修改成功，角色变更需重新登录后生效')
      } else {
        this.msgSuccess('修改成功')
      }
    },
    /** 删除按钮操作 */
    handleDelete(row) {
      const userIds = row.userId || this.ids
      this.$confirm('是否确认删除用户编号为"' + userIds + '"的数据项?', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(function() {
        return delUser(userIds)
      }).then(() => {
        this.getList()
        this.msgSuccess('删除成功')
      }).catch(function() {})
    },
    /** 导入按钮操作 */
    handleImport() {
      this.upload.title = '用户导入'
      this.upload.open = true
    },
    // 文件上传中处理
    handleFileUploadProgress(event, file, fileList) {
      this.upload.isUploading = true
    },
    // 文件上传成功处理
    handleFileSuccess(response, file, fileList) {
      this.upload.open = false
      this.upload.isUploading = false
      this.$refs.upload.clearFiles()
      this.$alert(response.msg, '导入结果', { dangerouslyUseHTMLString: true })
      this.getList()
    },
    // 提交上传文件
    submitFileForm() {
      this.$refs.upload.submit()
    }
  }
}
</script>
