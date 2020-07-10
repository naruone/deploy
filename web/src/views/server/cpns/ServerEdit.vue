<template>
    <el-drawer
        title="服务器添加/编辑"
        :before-close="handleClose"
        :visible.sync="showDialog"
        @opened="()=>{this.$refs.pInput.focus()}"
        size="35%"
        direction="rtl">
        <div style="height: calc(100vh - 75px);overflow-y: scroll">
            <el-form :model="serverForm" :rules="serverRules" class="drawerForm" label-width="100px"
                     ref="serverForm">
                <el-form-item label="服务器类型" prop="type">
                    <el-select v-model="serverForm.type">
                        <el-option v-for="v in sType" :key="v.value" :label="v.label" :value="v.value"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="服务器地址" prop="ssh_addr">
                    <el-input ref="pInput" v-model="serverForm.ssh_addr" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="端口" prop="ssh_port">
                    <el-input v-model.number="serverForm.ssh_port" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="用户名" prop="ssh_user">
                    <el-input v-model="serverForm.ssh_user" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="秘钥" prop="ssh_key">
                    <el-input :autosize="{maxRows:10,minRows:4}" type="textarea" v-model="serverForm.ssh_key"
                              autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="工作目录" prop="work_dir">
                    <el-input v-model="serverForm.work_dir" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button @click="handleClose">取 消</el-button>
                    <el-button type="primary" @click="serverSave">确 定</el-button>
                </el-form-item>
            </el-form>
        </div>
    </el-drawer>
</template>

<script>
    import {saveServer} from "../../../api/server";

    export default {
        name: "ServerEdit",
        props: {
            sType: {
                required: true,
                type: Array
            }
        },
        data() {
            return {
                showDialog: false,
                fieldDisabled: false,
                serverForm: {
                    server_id: 0,
                    type: 1,
                    ssh_addr: '',
                    ssh_port: 22,
                    ssh_user: '',
                    ssh_key: '',
                    work_dir: '/opt/_repo'
                },
                serverRules: {
                    type: [
                        {required: true, message: '服务器类型必选', trigger: 'blur'},
                        {type: 'enum', enum: [1, 2], message: '请选择正确的服务器类型', trigger: 'change'},
                    ],
                    ssh_addr: [
                        {required: true, message: '请输入服务器地址', trigger: 'blur'},
                    ],
                    ssh_port: [
                        {required: true, type: 'integer', message: '请输入正确的端口', trigger: 'blur'},
                    ],
                    ssh_user: [
                        {required: true, message: '用户名必填', trigger: 'blur'},
                        {min: 3, max: 200, message: '长度在3到200个字符', trigger: 'blur'}
                    ],
                    ssh_key: [
                        {max: 5000, message: '长度小于5000个字符', trigger: 'blur'}
                    ],
                    work_dir: [
                        {required: true, message: '工作目录必填', trigger: 'blur'},
                    ]
                }
            }
        },
        methods: {
            serverSave() {
                this.$refs.serverForm.validate(async (valid) => {
                    if (valid) {
                        await saveServer(this.serverForm).then((res) => {
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
                for (let k in this.serverForm) {
                    this.$set(this.serverForm, k, row[k] ? row[k] : '')
                }
                if (Object.keys(row).length === 0) {
                    this.serverForm.type = 1
                    this.serverForm.ssh_port = 22
                    this.serverForm.work_dir = '/opt/_repo'
                }
                this.showDialog = true
            },
            handleClose() {
                for (let k in this.serverForm) {
                    this.$set(this.serverForm, k, '')
                }
                this.$refs.serverForm.clearValidate()
                this.showDialog = false
            }
        }
    }
</script>
