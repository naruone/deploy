<template>
    <div>
        <el-form ref="searchForm" :model="searchForm" :inline="true">
            <el-form-item label="项目状态" prop="condition">
                <el-select style="width: 120px" v-model.number="searchForm.status">
                    <el-option v-for="v in pStatus" :key="v.value" :label="v.label" :value="v.value"/>
                </el-select>
            </el-form-item>
            <el-form-item label="筛选条件" prop="condition">
                <el-select style="width: 120px" v-model="searchForm.sCondition" placeholder="请选择">
                    <el-option key="project_name" label="项目名" value="project_name"/>
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
            <el-button @click="editProject({})"><i class="el-icon-plus"></i> 新增项目</el-button>
        </div>
        <el-table
                :data="tableData"
                border
                style="width: 100%">
            <el-table-column
                    prop="project_id"
                    label="ID"
                    width="50">
            </el-table-column>
            <el-table-column
                    prop="project_name"
                    label="项目名称"
                    width="120">
            </el-table-column>
            <el-table-column
                    prop="repo_url"
                    width="240"
                    label="仓库地址">
            </el-table-column>
            <el-table-column
                    prop="dst"
                    label="初始化目录">
            </el-table-column>
            <el-table-column
                    width="200"
                    prop="web_root"
                    label="Web目录">
            </el-table-column>
            <el-table-column
                    width="80"
                    prop="status"
                    label="INIT 状态">
                <template slot-scope="scope">
                    <el-popover v-if="scope.row.status === 3" trigger="hover" placement="top">
                        <p>失败原因: {{ scope.row.err_msg }}</p>
                        <div slot="reference" class="name-wrapper">
                            <el-tag type="danger" size="mini">{{ projectStatus(scope.row) }}</el-tag>
                        </div>
                    </el-popover>
                    <div v-else>
                        <el-tag size="mini">{{ projectStatus(scope.row) }}</el-tag>
                    </div>
                </template>

            </el-table-column>
            <el-table-column
                    width="200"
                    prop="create_at"
                    label="创建时间">
            </el-table-column>
            <el-table-column align="center" label="操作">
                <template slot-scope="scope">
                    <el-button type="primary" icon="el-icon-edit" circle
                               @click="editProject(scope.row)"
                               slot="reference"></el-button>
                    <el-tooltip class="item" v-if="scope.row.status === 1 || scope.row.status === 3" effect="light"
                                content="初始化" placement="top">
                        <el-button circle type="primary"
                                   @click="initProject(scope.row)"
                                   icon="el-icon-connection"></el-button>
                    </el-tooltip>
                    <el-button type="danger" icon="el-icon-delete" circle
                               @click="deleteProject(scope.row)"
                               slot="reference"></el-button>
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
        <ProjectEdit ref="projectEditFormDrawer"></ProjectEdit>
    </div>
</template>

<script>
    import tableInfo from "../../plugins/mixins/tableInfo";
    import {delProject, getProjectList, initProject} from "../../api/project"
    import ProjectEdit from "./cpns/ProjectEdit";

    export default {
        name: "Project",
        mixins: [tableInfo],
        data() {
            return {
                searchForm: {
                    status: 0   //默认请选择
                },
                pStatus: [
                    {
                        value: 0,
                        label: '请选择'
                    },
                    {
                        value: 1,
                        label: '未开始'
                    },
                    {
                        value: 2,
                        label: '成功'
                    },
                    {
                        value: 3,
                        label: '失败'
                    },
                    {
                        value: 9,
                        label: '进行中'
                    }
                ],
            }
        },
        created() {
            this.getTableData()
        },
        components: {
            ProjectEdit
        },
        methods: {
            getList: getProjectList,
            projectStatus(c) {
                return this.pStatus.reduce((o, v) => {
                    if (o !== '') return o
                    return c.status === v.value ? v.label : ''
                }, '')
            },
            async initProject(row) {
                this.$confirm('确认初始化该项目?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(async () => {
                    const res = await initProject({project_id: row.project_id})
                    if (res.code === 0) {
                        this.$message({
                            type: 'success',
                            message: res.msg
                        })
                        await this.getTableData()
                    }
                }).catch(() => {
                    this.$message({
                        type: 'info',
                        message: '已取消'
                    })
                })
            },
            async deleteProject(row) {
                this.$confirm('确认删除该项目, 此操作会删除对应的环境配置和已存在的发布任务?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(async () => {
                    await delProject({project_id: row.project_id}).then((res) => {
                        this.$message({
                            type: 'success',
                            message: res.msg
                        })
                        this.getTableData()
                    }).catch(() => {
                    })
                }).catch(() => {
                    this.$message({
                        type: 'info',
                        message: '已取消'
                    })
                })
            },
            editProject(row) {
                this.$refs.projectEditFormDrawer.setEditVal(row)
            }
        },
    }
</script>