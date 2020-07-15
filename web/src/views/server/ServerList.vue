<template>
    <div>
        <el-form ref="searchForm" :model="searchForm" :inline="true">
            <el-form-item label="项目状态" prop="condition">
                <el-select style="width: 120px" v-model.number="searchForm.type">
                    <el-option v-for="v in sType" :key="v.value" :label="v.label" :value="v.value"/>
                </el-select>
            </el-form-item>
            <el-form-item label="筛选条件" prop="condition">
                <el-select style="width: 120px" v-model="searchForm.sCondition" placeholder="请选择">
                    <el-option key="ssh_addr" label="Host" value="ssh_addr"/>
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
            <el-button @click="editServer({})"><i class="el-icon-plus"></i> 新增项目</el-button>
        </div>
        <el-table
            :data="tableData"
            border
            style="width: 100%">
            <el-table-column
                prop="server_id"
                label="ID"
                width="100">
            </el-table-column>
            <el-table-column
                label="服务器类型"
                width="200"
                :formatter="serverType">
            </el-table-column>
            <el-table-column
                prop="ssh_addr"
                width="120"
                label="服务器地址">
            </el-table-column>
            <el-table-column
                prop="ssh_port"
                width="100"
                label="服务器端口">
            </el-table-column>
            <el-table-column
                prop="ssh_key_path"
                label="私钥路径">
            </el-table-column>
            <el-table-column
                width="200"
                prop="create_at"
                label="创建时间">
            </el-table-column>
            <el-table-column align="center" label="操作" width="120">
                <template slot-scope="scope">
                    <el-button type="primary" icon="el-icon-edit" circle
                               @click="editServer(scope.row)"
                               slot="reference"></el-button>
                    <el-button type="danger" icon="el-icon-delete" circle
                               @click="deleteServer(scope.row)"
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
        <ServerEdit :s-type="sType" ref="serverEditFromDrawer"></ServerEdit>
    </div>
</template>

<script>
    import tableInfo from "../../plugins/mixins/tableInfo";
    import ServerEdit from "./cpns/ServerEdit";
    import {delServer, getServerList} from "../../api/server";

    export default {
        name: "ServerList",
        mixins: [tableInfo],
        components: {ServerEdit},
        data() {
            return {
                searchForm: {
                    type: ''
                },
                sType: [
                    {
                        value: 1,
                        label: '目标机'
                    },
                    {
                        value: 2,
                        label: '跳板机'
                    },
                ],
            }
        },
        created() {
            this.getTableData()
        },
        methods: {
            getList: getServerList,
            editServer(row) {
                this.$refs.serverEditFromDrawer.setEditVal(row)
            },
            async deleteServer(row) {
                this.$confirm('确认删除该服务器, 请先确保使用到本服务器的任务已被删除?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(async () => {
                    await delServer({server_id: row.server_id}).then((res) => {
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
            serverType(c) {
                return this.sType.reduce((o, v) => {
                    if (o !== '') return o
                    return c.type === v.value ? v.label : ''
                }, '')
            },
        }
    }
</script>

<style scoped>

</style>
