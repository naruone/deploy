<template>
    <div id="console_contain">
        <el-card class="box-card">
            <div slot="header">
                <span>发布系统环境</span>
            </div>
            <div class="item">
                系统环境: {{data['system'] }}
            </div>
            <div class="item">
                Go 版本: {{data['go_version'] }}
            </div>
            <div class="item">
                CPU核心数: {{data['cpu'] }}
            </div>
        </el-card>
        <el-card class="box-card">
            <div slot="header">
                <span>发布统计</span>
            </div>
            <div class="item">
                xxxx: xxxx
            </div>
            <div class="item">
                xxxx: xxxx
            </div>
            <div class="item">
                xxxx: xxxx
            </div>
        </el-card>
        <el-card class="box-card">
            <div slot="header">
                <span>项目统计</span>
            </div>
            <div class="item">
                xxxx: xxxx
            </div>
            <div class="item">
                xxxx: xxxx
            </div>
            <div class="item">
                xxxx: xxxx
            </div>
        </el-card>
        <el-card class="box-card">
            <div slot="header">
                <span>XXXX展示</span>
            </div>
            <div class="item">
                xxxx: xxxx
            </div>
            <div class="item">
                xxxx: xxxx
            </div>
            <div class="item">
                xxxx: xxxx
            </div>
        </el-card>
    </div>
</template>

<script>
    import {getDeployInfo} from '../../api/env_cfg'

    export default {
        name: "Console",
        data() {
            return {
                data: ''
            }
        },
        created() {
            getDeployInfo().then((res) => {
                console.log(res);
                this.data = res.data
            }).catch((e) => {
            })
        },
    }
</script>

<style scoped>
    #console_contain .box-card {
        display: inline-block;
        margin: 10px 20px 10px auto;
    }

    .box-card {
        width: 45%;
    }

    .item {
        margin-bottom: 18px;
        color: #555555;
        font-size: 14px;
    }
</style>
