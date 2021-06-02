<template>
    <div class="login-container">
        <el-form ref="loginForm" class="login-form" :rules="loginRules" :model="loginForm">
            <div class="title-container">
                <h3 class="title">系统登录</h3>
            </div>
            <el-form-item prop="username">
                <span class="svg-container">
                    <i class="el-icon-user"></i>
                </span>
                <el-input ref="username" style="width: 280px" tabindex="1" v-model.trim="loginForm.username"
                          placeholder="请输入用户名"></el-input>
            </el-form-item>
            <el-tooltip v-model="capsTooltip" content="大写已开启" placement="right" manual>
                <el-form-item prop="password">
                    <span class="svg-container">
                        <i class="el-icon-lock"></i>
                    </span>
                    <el-input style="width: 280px"
                              @blur="capsTooltip = false"
                              tabindex="2"
                              :show-password="true"
                              v-model.trim="loginForm.password"
                              @keyup.native="checkCapslock"
                              placeholder="请输入密码"></el-input>
                </el-form-item>
            </el-tooltip>
            <el-form-item prop="captcha">
                <span class="svg-container">
                    <i class="el-icon-key"></i>
                </span>
                <el-input class="captcha"
                          style="width: 200px"
                          @keydown.native.enter="submitForm"
                          tabindex="3"
                          v-model.trim="loginForm.captcha"
                          placeholder="请输入验证码"></el-input>
                <img class="img-captcha" @click="refreshCaptcha" v-if="picPath !== '#'" :src="picPath">
            </el-form-item>
            <el-form-item>
                <el-button type="primary" style="width: 100%;" size="medium" @click="submitForm">登录
                </el-button>
            </el-form-item>
        </el-form>
    </div>
</template>

<script>
    import {captcha} from '@/api/user'
    import {mapActions} from "vuex";

    export default {
        data() {
            return {
                loginForm: {
                    username: '',
                    password: '',
                    captcha: '',
                    captchaId: ''
                },
                loginRules: {
                    username: [
                        {required: true, message: '请输入用户名', trigger: 'blur'},
                        {min: 3, max: 18, message: '长度在3到18个字符', trigger: 'blur'}
                    ],
                    password: [
                        {required: true, message: '请输入密码', trigger: 'blur'},
                        {min: 3, max: 18, message: '长度在3到18个字符', trigger: 'blur'}
                    ],
                    captcha: [
                        {required: true, type: "string", len: 4, message: "请输入4位验证码", trigger: 'blur'}
                    ]
                },
                capsTooltip: false, //密码框大小写提示
                picPath: '#'    //验证码路径,
            }
        },
        created() {
            this.refreshCaptcha()
        },
        mounted() {
            if (this.loginForm.username === '') {
                this.$refs.username.focus()
            }
        },
        methods: {
            ...mapActions("user", ["LoginIn"]),
            checkCapslock(e) {
                const {key} = e
                this.capsTooltip = key && key.length === 1 && (key >= 'A' && key <= 'Z')
            },
            refreshCaptcha() {
                captcha({}).then((res) => {
                    let basePath = process.env.NODE_ENV === 'development' ? process.env.VUE_APP_DEPLOY_BASE_API : ''
                    this.picPath = basePath + res.data.picPath
                    this.loginForm.captchaId = res.data.captchaId
                })
            },
            submitForm() {
                this.$refs.loginForm.validate((valid) => {
                    if (valid) {
                        this.LoginIn(this.loginForm).catch(() => {
                        });
                        this.refreshCaptcha();
                    } else {
                        this.$message({
                            type: "error",
                            message: "登录信息填写错误!",
                            showClose: true
                        });
                        return false;
                    }
                });
            }
        }
    }
</script>

<style>
    .login-container > .login-form .el-input > input {
        background: transparent;
        border: none;
        color: #eee;
        padding-left: 5px;
        height: 45px;
        line-height: 45px;
        caret-color: #fff; /* 光标颜色 */
    }
</style>
<style scoped>
    .login-container {
        background-color: #323a4c;
        width: 100%;
        height: 100%;
        margin: auto;
    }

    .title-container > .title {
        font-size: 26px;
        color: #eee;
        margin: 0px auto 30px auto;
        text-align: center;
        font-weight: bold;
    }

    .login-container > .login-form {
        position: relative;
        width: 320px;
        max-width: 100%;
        padding: 160px 35px 0;
        margin: 0 auto;
        overflow: hidden;
    }

    .login-container .el-form-item {
        border: 1px solid rgba(255, 255, 255, 0.1);
        background: rgba(0, 0, 0, 0.1);
        border-radius: 5px;
        color: #454545;
    }

    .svg-container {
        margin-left: 10px;
        color: #889aa4;
        vertical-align: middle;
        width: 25px;
        text-align: center;
        display: inline-block;
    }

    .captcha {
        position: relative;
    }

    .img-captcha {
        position: absolute;
        top: 5px;
        background-color: rgba(255, 255, 255, 0.7);
        border-radius: 2px;
    }
</style>
