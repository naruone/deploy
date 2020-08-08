<template>
    <div class="process-info">
        <div class="process-info-common">
            <el-divider content-position="left">公共</el-divider>
            <el-steps :space="100" finish-status="success"
                      :process-status="this.processData[this.taskId]['now_status']"
                      :active="this.processData[this.taskId]['now_step']"
                      align-center>
                <el-step title="开始"></el-step>
                <el-step v-for="v in comProcess" :title="v.v"></el-step>
            </el-steps>
        </div>
        <div class="process-info-server">
            <div v-for="(_act,k) in processData[taskId]['servers']">
                <el-divider content-position="left">{{k}}</el-divider>
                <el-steps :space="100" finish-status="success"
                          :process-status="_act['now_status']"
                          :active="_act['now_step']" align-center>
                    <el-step v-for="v in serverProcess" :title="v.v"></el-step>
                </el-steps>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        name: "ProcessInfo",
        props: {
            isJumper: {
                type: Boolean,
                default: false
            },
            processData: Object,
            comProcess: Array,
            serverProcess: Array,
            taskId: Number
        }
    }
</script>

<style scoped>
    .process-info > div {
        margin-bottom: 5px;
        border-radius: 3px;
        display: inline-block;
        alignment-baseline: top;
        vertical-align: middle;
    }

    .process-info-common {
        width: 300px;
    }

    .process-info-server {
        width: 300px;
        border-left: 1px solid #CCC;
    }
</style>
