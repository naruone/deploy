<template>
    <el-drawer
            title="配置添加/编辑"
            :before-close="handleClose"
            :visible.sync="showDialog"
            @opened="()=>{this.$refs.pInput.focus()}"
            direction="rtl">
        <div style="height: calc(100vh - 75px);overflow-y: scroll">
            <el-form :model="cfgForm" :rules="cfgRules" class="drawerForm" label-width="100px"
                     ref="cfgForm">
                <el-form-item label="环境名称" prop="env_name">
                    <el-input ref="pInput" v-model="cfgForm.env_name" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="项目" prop="project_id">
                    <el-select v-model="cfgForm.project_id">
                        <el-option v-for="v in projectOptions" :key="v.value" :label="v.label" :value="v.value"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="目标机" prop="server_ids">
                    <el-select multiple v-model="cfgForm.server_ids" placeholder="请选择目标机">
                        <el-option v-for="v in serverOptions" :key="v.value" :label="v.label" :value="v.value"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="跳板机" prop="jump_server">
                    <el-select v-model="cfgForm.jump_server">
                        <el-option v-for="v in jumperOptions" :key="v.value" :label="v.label" :value="v.value"/>
                    </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button @click="handleClose">取 消</el-button>
                    <el-button type="primary" @click="cfgSave">确 定</el-button>
                </el-form-item>
            </el-form>
        </div>
    </el-drawer>
</template>

<script>
    import {getCfgSelectOptions, saveEnvCfg} from "../../../api/env_cfg";

    export default {
        name: "CfgEdit",
        data() {
            return {
                showDialog: false,
                fieldDisabled: false,
                cfgForm: {
                    env_id: 0,
                    env_name: '',
                    project_id: 0,
                    server_ids: [],
                    jump_server: 0
                },
                cfgRules: {
                    env_name: [
                        {required: true, message: '环境名称必填', trigger: 'blur'},
                        {min: 3, max: 200, message: '长度在3到200个字符', trigger: 'blur'}
                    ],
                    project_id: [
                        {required: true, message: '项目必选', trigger: 'change'},
                    ],
                    server_ids: [
                        {required: true, message: '目标机必选', trigger: 'blur'},
                    ]
                },
                projectOptions: [],
                serverOptions: [],
                jumperOptions: [{
                    label: '请选择',
                    value: 0
                }]
            }
        },
        watch: {
            showDialog(n) {
                if (n) this.getSelVal()
            }
        },
        methods: {
            cfgSave() {
                this.$refs.cfgForm.validate(async (valid) => {
                    if (valid) {
                        let formData = Object.assign({}, this.cfgForm);
                        formData.server_ids = formData.server_ids.join(',')
                        await saveEnvCfg(formData).then((res) => {
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
                for (let k in this.cfgForm) {
                    this.$set(this.cfgForm, k, row[k] ? row[k] : '')
                }
                if (typeof this.cfgForm.server_ids === 'string' && String(this.cfgForm.server_ids) !== '') {
                    this.cfgForm.server_ids = this.cfgForm.server_ids.split(',').map((v) => {
                        return Number(v)
                    })
                }
                this.showDialog = true
            },
            handleClose() {
                for (let k in this.cfgForm) {
                    this.$set(this.cfgForm, k, '')
                }
                this.cfgForm.server_ids = []
                this.$refs.cfgForm.clearValidate()
                this.showDialog = false
                this.projectOptions = []
                this.serverOptions = []
                this.jumperOptions = [{
                    label: '请选择',
                    value: 0
                }]
            },
            async getSelVal() {
                await getCfgSelectOptions().then((res) => {
                    for (let v of res.data['projects']) {
                        this.projectOptions.push({
                            label: v.project_name,
                            value: v.project_id
                        })
                    }
                    for (let v of res.data['servers']) {
                        if (v.type === 1) { //目标机
                            this.serverOptions.push({
                                label: v.ssh_addr,
                                value: v.server_id
                            })
                        } else {//跳板机
                            this.jumperOptions.push({
                                label: v.ssh_addr,
                                value: v.server_id
                            })
                        }
                    }
                }).catch(() => {
                })
            }
        }
    }
</script>
