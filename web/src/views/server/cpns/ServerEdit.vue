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
                    <el-radio-group v-model="serverForm.type">
                        <el-radio v-for="v in sType" :key="v.value" :label="v.value">{{ v.label }}</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-form-item label="服务器地址" prop="ssh_addr">
                    <el-input ref="pInput" v-model="serverForm.ssh_addr" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="SSH端口" prop="ssh_port">
                    <el-input v-model.number="serverForm.ssh_port" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="登录名称" prop="ssh_user">
                    <el-input v-model="serverForm.ssh_user" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="秘钥路径" prop="ssh_key_path">
                    <el-input v-model="serverForm.ssh_key_path" autocomplete="off"></el-input>
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
import {saveServer} from "@/api/server";

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
                ssh_key_path: '',
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
                ssh_key_path: [
                    {required: true, message: 'ssh key 路径必填', trigger: 'blur'},
                    {max: 300, message: '长度小于300个字符', trigger: 'blur'}
                ],
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
