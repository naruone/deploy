<template>
    <div id="console_contain">
        <el-card class="box-card">
            <div slot="header">
                <span>发布系统环境</span>
            </div>
            <div class="item">
                系统信息: {{data['system'] }}
            </div>
            <div class="item">
                Go 版本: {{data['go_version'] }}
            </div>
            <div class="item">
                CPU 核心: {{data['cpu'] }}
            </div>
            <div class="item">
                Go 协程数: {{data['goroutine_num'] }}
            </div>
            <div class="item">
                Mem 占用: {{data['sys_mem_used'] }}
            </div>
            <div class="item">
                (GC)内存: {{data['next_gc_mem'] }}
            </div>
            <div class="item">
                启动时间: {{data['start_time'] }}
            </div>
        </el-card>
        <el-card class="box-card">
            <div slot="header">
                <span>XXXX统计数据(暂不支持-懒)</span>
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
        width: 100%;
    }

    .item {
        width: 50%;
        padding-top: 10px;
        text-indent: 1.2em;
        padding-bottom: 10px;
        color: #555555;
        font-size: 14px;
        display: inline-block;
        outline: none;
        box-sizing: border-box;
        border-right: 1px #CCC solid;
    }
</style>
