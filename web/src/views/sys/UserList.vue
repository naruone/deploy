<template>
    <div>
        <el-form ref="searchForm" :model="searchForm" :inline="true">
            <el-form-item label="用户状态" prop="condition">
                <el-select style="width: 120px" v-model.number="searchForm.status">
                    <el-option v-for="v in uStatus" :key="v.value" :label="v.label" :value="v.value"/>
                </el-select>
            </el-form-item>
            <el-form-item label="筛选条件" prop="condition">
                <el-select style="width: 120px" v-model="searchForm.sCondition" placeholder="请选择">
                    <el-option key="user_name" label="用户名" value="user_name"/>
                </el-select>
            </el-form-item>
            <el-form-item prop="searchValue">
                <el-input v-model="searchForm.sValue" placeholder="搜索关键词"></el-input>
            </el-form-item>
            <el-form-item>
                <el-button @click="getTableData(1)" type="primary">查询</el-button>
            </el-form-item>
        </el-form>
        <div class="button-box">
            <el-button @click="showUserEdit({})"><i class="el-icon-plus"></i> 新增用户</el-button>
        </div>
        <el-table :data="tableData" border style="width: 100%">
            <el-table-column
                    prop="user_id"
                    label="ID"
                    width="120">
            </el-table-column>
            <el-table-column
                    prop="user_name"
                    label="登录名"
                    width="180">
            </el-table-column>
            <el-table-column
                    prop="nick_name"
                    label="用户昵称">
            </el-table-column>
            <el-table-column
                    :formatter="userStatus"
                    label="用户状态">
            </el-table-column>
            <el-table-column
                    prop="create_at"
                    label="创建时间">
            </el-table-column>
            <el-table-column align="center" label="操作">
                <template slot-scope="scope">
                    <el-button circle type="primary" @click="showUserEdit(scope.row)" icon="el-icon-edit"></el-button>
                    <el-tooltip class="item" v-if="scope.row.status === 2 && scope.row.user_id !== 1" effect="light"
                                content="启用" placement="top">
                        <el-button circle
                                   @click="setUserStatus(scope.row,1)"
                                   type="primary" icon="el-icon-unlock"></el-button>
                    </el-tooltip>
                    <el-tooltip class="item" v-else-if="scope.row.status === 1 && scope.row.user_id !== 1"
                                effect="light"
                                content="禁用" placement="top">
                        <el-button circle
                                   @click="setUserStatus(scope.row,2)"
                                   type="warning" icon="el-icon-lock"></el-button>
                    </el-tooltip>
                </template>
            </el-table-column>
        </el-table>
        <div class="page-content">
            <el-pagination
                    @size-change="handleSizeChange"
                    @current-change="handleCurrentChange"
                    :current-page="currentPage"
                    :page-sizes="[10, 20, 50, 100]"
                    :page-size="pageSize"
                    layout="total, sizes, prev, pager, next, jumper"
                    :total="total">
            </el-pagination>
        </div>

        <el-drawer title="用户添加/编辑" :visible.sync="userEditFormShow"
                   direction="rtl"
                   @opened="()=>{this.$refs.mInput.focus()}"
                   :before-close="handleClose"
                   size="30%">
            <el-form ref="uForm" :model="userForm" class="drawerForm" :rules="userFormRule" label-width="70px">
                <el-form-item label="用户名" prop="user_name">
                    <el-input v-model.trim="userForm.user_name" ref="mInput" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="密 码" prop="password">
                    <el-input v-model.trim="userForm.password" type="password" show-password
                              autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="昵 称" prop="nick_name">
                    <el-input v-model.trim="userForm.nick_name" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button @click="handleClose">取 消</el-button>
                    <el-button type="primary" @click="userSave">确 定</el-button>
                </el-form-item>
            </el-form>
        </el-drawer>
    </div>
</template>

<script>
    import tableInfo from '@/plugins/mixins/tableInfo'
    import {getUserList, saveUser, setUserStatus} from "@/api/user";

    export default {
        name: "UserList",
        mixins: [tableInfo],
        data() {
            return {
                userEditFormShow: false,
                userForm: {
                    user_id: 0,
                    user_name: '',
                    password: '',
                    nick_name: ''
                },
                userFormRule: {
                    user_name: [
                        {required: true, message: '请输入用户名', trigger: 'blur'},
                        {min: 3, max: 18, message: '长度在3到18个字符', trigger: 'blur'}
                    ],
                    password: [
                        {message: '请输入密码', trigger: 'blur'},
                        {min: 3, max: 18, message: '长度在3到18个字符', trigger: 'blur'}
                    ],
                    nick_name: [
                        {required: true, type: "string", min: 3, max: 20, message: "请输入昵称", trigger: 'blur'}
                    ]
                },
                uStatus: [
                    {
                        value: 1,
                        label: '正常'
                    },
                    {
                        value: 2,
                        label: '禁用'
                    }
                ],
                searchForm: {
                    status: 1   //默认查状态正常的用户
                }
            }
        },
        created() {
            this.getTableData()
        },
        methods: {
            getList: getUserList,
            userStatus(c) { //状态展示
                return this.uStatus.reduce((o, v) => {
                    if (o !== '') return o
                    return c.status === v.value ? v.label : ''
                }, '')
            },
            setUserStatus(row, status) {
                this.$confirm('确定此操作?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(async () => {
                    await setUserStatus({
                        user_id: row.user_id,
                        status: status
                    }).then((res) => {
                        this.$message({
                            type: 'success',
                            message: res.msg
                        })
                        this.getTableData()
                    }).catch(() => {
                    })
                }, () => {
                    this.$message({
                        type: 'info',
                        message: '已取消'
                    })
                })
            },
            userSave() {
                this.$refs.uForm.validate(async (valid) => {
                    if (valid) {
                        await saveUser(this.userForm).then((res) => {
                            this.$message({
                                type: 'success',
                                message: res.msg,
                                showClose: true
                            });
                            this.getTableData()
                            this.userEditFormShow = false
                        }).catch(() => {
                        })
                    }
                })
            },
            showUserEdit(row) {
                row.password = ""
                for (let k in this.userForm) {
                    this.$set(this.userForm, k, row[k] ? row[k] : '')
                }
                this.userEditFormShow = true
            },
            handleClose() {
                this.$refs.uForm.clearValidate()
                this.userEditFormShow = false
            }
        }
    }
</script>
