<template>
    <div>
        <div v-for="(v,k) in dataRes">
            <el-divider content-position="left">{{k}}</el-divider>
            <span>消息: {{v.message}}</span>
            <span>状态: {{v.status}}</span>
            <el-divider><i class="el-icon-s-promotion"></i></el-divider>
        </div>
    </div>
</template>

<script>
    export default {
        name: "TaskInfo",
        props: {
            data: String
        },
        computed: {
            dataRes() {
                return JSON.parse(this.data)
            }
        }

    }
</script>

<style scoped>

</style>
