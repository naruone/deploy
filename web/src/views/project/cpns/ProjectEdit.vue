<template>
    <el-drawer
        title="项目添加/编辑"
        :before-close="handleClose"
        :visible.sync="showDialog"
        @opened="()=>{this.$refs.pInput.focus()}"
        direction="rtl">
        <div style="height: calc(100vh - 75px);overflow-y: scroll">
            <el-form :model="projectForm" :rules="projectRules" class="drawerForm" label-width="80px"
                     ref="projectForm">
                <el-form-item label="项目名称" prop="project_name">
                    <el-input ref="pInput" v-model="projectForm.project_name" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="Git 仓库" prop="repo_url">
                    <el-input :disabled="fieldDisabled" v-model="projectForm.repo_url" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="项目Dir" prop="dst">
                    <el-input :disabled="fieldDisabled" v-model="projectForm.dst" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="WebRoot" prop="web_root">
                    <el-input v-model="projectForm.web_root" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="AftScript" prop="after_script">
                    <el-input :autosize="{maxRows:10,minRows:4}" type="textarea" v-model="projectForm.after_script"
                              autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button @click="handleClose">取 消</el-button>
                    <el-button type="primary" @click="projectSave">确 定</el-button>
                </el-form-item>
            </el-form>
        </div>
    </el-drawer>
</template>

<script>
import {saveProject} from '@/api/project'

export default {
    name: "ProjectEdit",
    data() {
        return {
            showDialog: false,
            fieldDisabled: false,
            projectForm: {
                project_id: 0,
                project_name: '',
                repo_url: '',
                dst: '',
                web_root: '',
                after_script: ''
            },
            projectRules: {
                project_name: [
                    {required: true, message: '请输入项目名', trigger: 'blur'},
                    {min: 3, max: 200, message: '长度在3到200个字符', trigger: 'blur'}
                ],
                repo_url: [
                    {required: true, message: '请输入仓库地址', trigger: 'blur'},
                    {min: 6, max: 255, message: '长度在6到255个字符', trigger: 'blur'}
                ],
                dst: [
                    {required: true, message: '输入项目初始化目录', trigger: 'blur'},
                    {min: 3, max: 200, message: '长度在3到200个字符', trigger: 'blur'}
                ],
                web_root: [
                    {required: true, message: '输入WebRoot', trigger: 'blur'},
                    {min: 3, max: 200, message: '长度在3到200个字符', trigger: 'blur'}
                ],
                after_script: [
                    {max: 5000, message: '长度小于5000个字符', trigger: 'blur'}
                ]
            }
        }
    },
    methods: {
        projectSave() {
            this.$refs.projectForm.validate(async (valid) => {
                if (valid) {
                    await saveProject(this.projectForm).then((res) => {
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
        setEditVal(row) {
            for (let k in this.projectForm) {
                this.$set(this.projectForm, k, row[k] ? row[k] : '')
            }
            this.fieldDisabled = Object.keys(row).length !== 0 && row.status === 2
            this.showDialog = true
        },
        handleClose() {
            for (let k in this.projectForm) {
                this.$set(this.projectForm, k, '')
            }
            this.$refs.projectForm.clearValidate()
            this.showDialog = false
            this.fieldDisabled = false
        }
    }
}
</script>
