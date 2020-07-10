<template>
    <div>
        <el-form ref="searchForm" :model="searchForm" :inline="true">
            <el-form-item label="筛选条件" prop="condition">
                <el-select style="width: 120px" v-model="searchForm.sCondition" placeholder="请选择">
                    <el-option key="env_name" label="名称" value="env_name"/>
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
            <el-button @click="editCfg({})"><i class="el-icon-plus"></i> 新增配置</el-button>
        </div>
        <el-table
                :data="tableData"
                border
                style="width: 100%">
            <el-table-column
                    prop="env_id"
                    label="ID"
                    width="40">
            </el-table-column>
            <el-table-column
                    prop="env_name"
                    label="环境名称">
            </el-table-column>
            <el-table-column
                    prop="Project.project_name"
                    label="项目">
            </el-table-column>
            <el-table-column
                    class-name="pre-line"
                    :formatter="serversFormatter"
                    label="服务器">
            </el-table-column>
            <el-table-column
                    :formatter="jumperFormatter"
                    label="跳板机">
            </el-table-column>
            <el-table-column
                    prop="last_ver"
                    label="最近版本">
            </el-table-column>
            <el-table-column
                    width="150"
                    prop="uuid"
                    label="UUID">
            </el-table-column>
            <el-table-column
                width="200"
                prop="create_at"
                label="创建时间">
            </el-table-column>
            <el-table-column align="center" label="操作">
                <template slot-scope="scope">
                    <el-button type="primary" icon="el-icon-edit" circle
                               @click="editCfg(scope.row)"
                               slot="reference"></el-button>
                    <el-button type="danger" icon="el-icon-delete" circle
                               @click="deleteCfg(scope.row)"
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
        <CfgEdit ref="cfgEditorFormDrawer"></CfgEdit>
    </div>
</template>

<script>
    import tableInfo from "../../plugins/mixins/tableInfo";
    import {delEnvCfg, getEnvCfgList} from "../../api/env_cfg";
    import CfgEdit from "./cpns/CfgEdit";

    export default {
        name: "EnvConfig",
        mixins: [tableInfo],
        components: {CfgEdit},
        created() {
            this.getTableData()
        },
        methods: {
            getList: getEnvCfgList,
            editCfg(row) {
                this.$refs.cfgEditorFormDrawer.setEditVal(row)
            },
            async deleteCfg(row) {
                this.$confirm('确认删除该配置吗, 此操作会删除项目对应所有发布记录?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(async () => {
                    await delEnvCfg({cfg_id: row.env_id}).then((res) => {
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
            serversFormatter(r) {
                return r['Servers'].reduce(function (o, n) {
                    return o ? (o + "\n" + n.ssh_addr) : n.ssh_addr;
                }, '')
            },
            jumperFormatter(r) {
                return r['Jumper'].ssh_addr ? r['Jumper'].ssh_addr : '-'
            }
        }
    }
</script>

<style scoped>

</style>
