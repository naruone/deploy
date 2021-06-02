<template>
    <el-container>
        <el-header>
            <Header @triggerCollapse="triggerMainCollapse" @showMdyPwd="mdyPwdShow=true"
                    :isCollapse="isCollapse"></Header>
        </el-header>
        <el-container>
            <el-aside :width="isCollapse ? 'auto': '200px'">
                <Aside :isCollapse="isCollapse"/>
            </el-aside>
            <el-main>
                <div class="main-breadcrumb">
                    <bread-crumb/>
                </div>
                <div class="main-content">
                    <keep-alive>
                        <router-view/>
                    </keep-alive>
                </div>
            </el-main>
        </el-container>

        <el-dialog :visible.sync="mdyPwdShow" @close="clearPwd" title="修改密码" width="360px">
            <el-form :model="pwdForm" :rules="pwdFormRule" label-width="80px" ref="modifyPwdForm">
                <el-form-item label="原密码" prop="old_pwd">
                    <el-input show-password v-model="pwdForm.old_pwd"></el-input>
                </el-form-item>
                <el-form-item label="新密码" prop="new_pwd">
                    <el-input show-password v-model="pwdForm.new_pwd"></el-input>
                </el-form-item>
                <el-form-item label="确认密码" prop="c_new_pwd">
                    <el-input show-password v-model="pwdForm.c_new_pwd"></el-input>
                </el-form-item>
            </el-form>
            <div class="dialog-footer" slot="footer">
                <el-button @click="mdyPwdShow=false">取 消</el-button>
                <el-button @click="modifyPwd" type="primary">确 定</el-button>
            </div>
        </el-dialog>
    </el-container>
</template>

<script>
    import Aside from "./cpns/Aside";
    import Header from "./cpns/Header";
    import BreadCrumb from "./cpns/BreadCrumb";
    import {modifyPwd} from "@/api/user";

    export default {
        name: "Layout",
        components: {
            Aside, Header, BreadCrumb
        },
        data() {
            return {
                isCollapse: false,
                mdyPwdShow: false,
                pwdForm: {
                    old_pwd: '',
                    new_pwd: '',
                    c_new_pwd: ''
                },
                pwdFormRule: {
                    old_pwd: [
                        {required: true, message: '请输入原密码', trigger: 'blur'},
                        {min: 6, max: 18, message: '长度在6到18个字符', trigger: 'blur'}
                    ],
                    new_pwd: [
                        {required: true, message: '请输入新密码', trigger: 'blur'},
                        {min: 6, max: 18, message: '长度在6到18个字符', trigger: 'blur'}
                    ],
                    c_new_pwd: [
                        {required: true, type: "string", message: "请输入确认密码", trigger: 'blur'},
                        {
                            validator: (rule, value, callback) => {
                                if (value !== this.pwdForm.new_pwd) {
                                    return callback(new Error('确认密码不对'));
                                }
                                callback()
                            }, trigger: 'blur'
                        }
                    ]
                }
            }
        },
        methods: {
            triggerMainCollapse(status) {
                this.isCollapse = status
            },
            modifyPwd() {
                this.$refs.modifyPwdForm.validate(async (valid) => {
                    if (valid) {
                        await modifyPwd(this.pwdForm).then((res) => {
                            this.$message({
                                type: 'success',
                                message: res.msg,
                                showClose: true
                            });
                            this.mdyPwdShow = false
                        }).catch(() => {
                        })
                    }
                })
            },
            clearPwd() {
                this.pwdForm = {
                    old_pwd: '',
                    new_pwd: '',
                    c_new_pwd: ''
                }
                this.$refs.modifyPwdForm.clearValidate()
            },
        }
    }
</script>

<style scoped>
    .el-container {
        height: 100vh;
        width: 100vw;
    }

    .el-header {
        background-color: #ffffff;
        box-shadow: 0 1px 1px #CCC;
        z-index: 2;
    }

    .el-main {
        background-color: #f0f2f5;
        padding: 12px;
    }

    .el-main > .main-breadcrumb {
        background-color: #ffffff;
        padding: 15px;
        border-radius: 3px 3px 0 0;
        border-bottom: #eeeeee solid 1px;
        box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)
    }

    .el-main > .main-content {
        min-height: 300px;
        background-color: #FFF;
        padding: 15px 15px 60px 15px;
        border-top: 3px double #EEE;
        box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    }
</style>
