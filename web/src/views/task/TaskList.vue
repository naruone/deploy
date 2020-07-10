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
            <el-button @click="editTask({})"><i class="el-icon-plus"></i> 新增任务</el-button>
        </div>
        <el-table
            :data="tableData"
            border
            style="width: 100%">
            <el-table-column
                prop="task_id"
                label="ID"
                width="40">
            </el-table-column>
            <el-table-column
                label="任务标题">
                <template slot-scope="scope">
                    <el-tooltip class="item" effect="dark" :content="scope.row.description" placement="right">
                        <span>{{scope.row.task_name}}</span>
                    </el-tooltip>
                </template>
            </el-table-column>
            <el-table-column
                prop="EnvCfg.env_name"
                label="环境">
            </el-table-column>
            <el-table-column
                :formatter="deployTypeFormatter"
                width="70"
                label="发布类型">
            </el-table-column>
            <el-table-column
                prop="version"
                width="80"
                label="版本">
            </el-table-column>
            <el-table-column
                prop="uuid"
                width="280"
                label="UUID">
            </el-table-column>
            <el-table-column
                width="80"
                label="状态">
                <template slot-scope="scope">
                    <el-tooltip v-if="scope.row.status === 9" placement="top-end" effect="dark">
                        <div slot="content">
                            失败原因:
                            <pre>{{scope.row.err_output}}</pre>
                        </div>
                        <el-tag type="danger" size="mini">{{ status[scope.row.status] }}</el-tag>
                    </el-tooltip>
                    <el-tag v-else-if="scope.row.status === 8" type="success" size="mini">
                        {{ status[scope.row.status] }}
                    </el-tag>
                    <el-tag v-else size="mini">{{ status[scope.row.status] }}</el-tag>
                </template>
            </el-table-column>
            <el-table-column
                class-name="pre-line"
                :formatter="timeDisplay"
                width="210"
                label="时间">
            </el-table-column>
            <el-table-column align="center" label="操作">
                <template slot-scope="scope">
                    <el-tooltip v-if="scope.row.status === 1" content="发布" placement="top" effect="dark">
                        <el-button type="primary" icon="el-icon-upload" circle></el-button>
                    </el-tooltip>
                    <el-button v-if="scope.row.status === 1" type="primary" icon="el-icon-edit" circle
                               @click="editTask(scope.row)"
                               slot="reference"></el-button>
                    <el-button v-if="scope.row.status === 1 || scope.row.status === 9" type="danger"
                               icon="el-icon-delete" circle
                               @click="delTask(scope.row)"
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
        <TaskEdit :d-type="dType" ref="taskEditorFormDrawer"></TaskEdit>
    </div>
</template>

<script>
    import tableInfo from "../../plugins/mixins/tableInfo";
    import TaskEdit from "./cpns/TaskEdit";
    import {getTaskList} from "../../api/task";

    export default {
        name: "TaskList",
        mixins: [tableInfo],
        components:{TaskEdit},
        created() {
            this.getTableData()
        },
        data() {
            return {
                dType: [
                    {
                        label: "全量",
                        value: 1
                    },
                    {
                        label: "增量",
                        value: 2
                    },
                ],
                status: {
                    '1': '未开始',
                    '2': '发布中',
                    '8': '成功',
                    '9': '失败'
                }
            }
        },
        methods: {
            getList: getTaskList,
            editTask(row) {
                this.$refs.taskEditorFormDrawer.setEditVal(row)
            },
            delTask(row) {

            },
            timeDisplay(row) {
                return "C: " + row.create_at + "\nU: " + row.update_at
            },
            deployTypeFormatter(row) {
                return this.dType.reduce((o, n) => {
                    return n.value === row.deploy_type ? n.label : o;
                }, '')
            },
        },
    }
</script>

<style scoped>

</style>
