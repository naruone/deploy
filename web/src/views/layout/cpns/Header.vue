<template>
    <div class="deploy-header">
        <span @click="goHome" class="title">
            <i class="el-icon-upload"></i>  Code 发布系统
        </span>
        <div @click="triggerCollapse" class="menu-trigger">
            <i class="el-icon-s-unfold" v-if="isCollapse"></i>
            <i class="el-icon-s-fold" v-else></i>
        </div>
        <div class="user-menu">
            <el-dropdown :hide-on-click="false">
                <span class="el-dropdown-link">
                    <span class="el-dropdown-link">
                    <img src="../../../assets/images/avatar.png" height="30" width="30"/>
                    {{$store.getters['user/userInfo'].nick_name}}
                    <i class="el-icon-arrow-down"></i>
                  </span>
                </span>
                <el-dropdown-menu slot="dropdown">
                    <el-dropdown-item><i class="el-icon-user-solid"></i>个人信息</el-dropdown-item>
                    <el-dropdown-item @click.native="showModifyPwd"><i class="el-icon-key"></i>修改密码</el-dropdown-item>
                    <el-dropdown-item @click.native="LoginOut"><i class="el-icon-cold-drink"></i>退出登陆</el-dropdown-item>
                </el-dropdown-menu>
            </el-dropdown>
        </div>
    </div>
</template>

<script>
    import {mapActions} from "vuex";

    export default {
        name: "Header",
        props: {
            isCollapse: {
                type: Boolean,
                default: false
            }
        },
        methods: {
            ...mapActions('user', ['LoginOut']),
            triggerCollapse() {
                this.$emit('triggerCollapse', !this.isCollapse)
            },
            goHome() {
                if (this.$route.name !== 'sys') {
                    this.$router.push({name: 'sys'});
                }
            },
            showModifyPwd() {
                this.$emit("showMdyPwd")
            },
        }
    }
</script>

<style scoped>
    .deploy-header {
        height: 60px;
        display: flex;
    }

    .deploy-header > .title {
        width: 200px;
        display: inline-flex;
        font-size: 1.2em;
        line-height: 60px;
        color: #555;
    }

    .deploy-header > .title {
        height: 50%;
        cursor: pointer;
    }

    .deploy-header > .title > i {
        margin-top: 18px;
        margin-right: 5px;
    }

    .menu-trigger {
        display: inline-flex;
        font-size: 1.5em;
        margin-top: 13px;
    }

    .user-menu {
        margin: 12px 25px auto auto;
    }

    .user-menu .el-dropdown-link {
        text-align: center;
        vertical-align: middle;
        cursor: default;
    }

    .user-menu .el-dropdown-link img {
        vertical-align: middle;
        border: 1px solid #ccc;
        border-radius: 15px;
    }

</style>