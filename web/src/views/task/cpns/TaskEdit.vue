<template>
    <el-drawer
        title="任务添加/编辑"
        :before-close="handleClose"
        :visible.sync="showDialog"
        @opened="()=>{this.$refs.pInput.focus()}"
        direction="rtl">
        <div style="height: calc(100vh - 75px);overflow-y: scroll">
            <el-form :model="taskForm" :rules="taskRules" class="drawerForm" label-width="100px"
                     ref="taskForm">
                <el-form-item label="任务名称" prop="task_name">
                    <el-input ref="pInput" v-model="taskForm.task_name" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="发布环境" prop="env_id">
                    <el-select @change="loadBranch(taskForm.env_id)" v-model="taskForm.env_id">
                        <el-option v-for="v in envOptions" :key="v.value" :label="v.label" :value="v.value"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="发布类型" prop="deploy_type">
                    <el-radio-group v-model="taskForm.deploy_type">
                        <el-radio v-for="v in dType" :key="v.value" :label="v.value">{{ v.label }}</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-form-item label="分支" prop="branch">
                    <el-select @change="loadVersion(taskForm.env_id,taskForm.branch)" v-model="taskForm.branch"
                               placeholder="注意: 选择目标环境当前版本号,相当于全量发">
                        <el-option v-for="v in branchOptions" :key="v" :label="v" :value="v"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="版本" prop="version">
                    <el-select v-model="taskForm.version">
                        <el-option v-for="v in versionOptions" :key="v.Version" :label="v.Message" :value="v.Version"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="描述" prop="description">
                    <el-input type="textarea" v-model="taskForm.description" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="AftScript" prop="after_script">
                    <el-input type="textarea" v-model="taskForm.after_script" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button @click="handleClose">取 消</el-button>
                    <el-button type="primary" @click="taskSave">确 定</el-button>
                </el-form-item>
            </el-form>
        </div>
    </el-drawer>
</template>

<script>
import {getBranches, getEnvOptions, getVersions, saveTask} from "@/api/task";

export default {
    name: "TaskEdit",
    props: {
        dType: {
            required: true,
            type: Array
        }
    },
    data() {
        return {
            showDialog: false,
            taskForm: {
                task_id: 0,
                task_name: '',
                env_id: '',
                deploy_type: 1,
                branch: '',
                version: '',
                description: '',
                after_script: '',
            },
            taskRules: {
                task_name: [
                    {required: true, message: '任务名称必填', trigger: 'blur'},
                    {min: 3, max: 200, message: '长度在3到200个字符', trigger: 'blur'}
                ],
                env_id: [
                    {required: true, message: '环境必选', trigger: 'change'},
                ],
                deploy_type: [
                    {required: true, message: '类型必选', trigger: 'blur'},
                ],
                branch: [
                    {required: true, message: '分支必选', trigger: 'blur'},
                ],
                version: [
                    {required: true, message: '版本必选', trigger: 'blur'},
                ],
                description: [
                    {max: 2000, message: '长度小于2000个字符', trigger: 'blur'}
                ],
                after_script: [
                    {max: 2000, message: '长度小于2000个字符', trigger: 'blur'}
                ]
            },
            envOptions: [],
            branchOptions: [],
            versionOptions: []
        }
    },
    watch: {
        showDialog(n) {
            if (n) this.getSelVal()
        }
    },
    methods: {
        taskSave() {
            this.$refs.taskForm.validate(async (valid) => {
                if (valid) {
                    await saveTask(this.taskForm).then((res) => {
                        this.$message({
                            type: 'success',
                            message: res.msg,
                            showClose: true
                        });
                        this.$parent.getTableData()
                        this.handleClose()
                    }).catch((_) => {
                    })
                }
            })
        },
        setEditVal() {
            this.showDialog = true
        },
        handleClose() {
            for (let k in this.taskForm) {
                this.$set(this.taskForm, k, '')
            }
            this.taskForm.deploy_type = 1
            this.$refs.taskForm.clearValidate()
            this.envOptions = []
            this.branchOptions = []
            this.versionOptions = []
            this.showDialog = false
        },
        async getSelVal() {
            await getEnvOptions().then((res) => {
                for (let v of res.data.envList) {
                    this.envOptions.push({
                        label: v.env_name,
                        value: v.env_id,
                        info: v
                    })
                }
            }).catch(() => {
            })
        },
        async loadBranch(envId) {
            this.taskForm.branch = ''
            this.taskForm.version = ''
            let envObj = this.envOptions.reduce((o, n) => {
                return n.value === envId ? n.info : o
            }, '')
            await getBranches({project_id: envObj.project_id}).then((res) => {
                this.branchOptions = res.data
            }).catch(() => {
            })
        },
        async loadVersion(envId, branch) {
            let envObj = this.envOptions.reduce((o, n) => {
                return n.value === envId ? n.info : o
            }, '')
            await getVersions({project_id: envObj.project_id, branch: branch}).then((res) => {
                this.versionOptions = res.data
            }).catch(() => {
            })
        }
    }
}
</script>
