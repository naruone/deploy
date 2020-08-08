<template>
    <div class="task-info">
        <div v-for="(v,k) in dataRes">
            <el-divider content-position="left">{{k}}</el-divider>
            <div style="margin: 0 15px 10px 15px">
                <span>消息:</span>
                <div style="width: 100%;overflow: scroll;">
                    <pre>{{v.message}}</pre>
                </div>
                <span>状态:</span> {{status[v.status]}}
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        name: "TaskInfo",
        props: {
            data: String,
            status: Array
        },
        computed: {
            dataRes() {
                try {
                    return JSON.parse(this.data)
                } catch (e) {
                    return {}
                }
            }
        }
    }
</script>

<style scoped>
    .task-info > div {
        border: 1px solid #EEE;
        margin-bottom: 5px;
        border-radius: 3px;
        font-size: 0.8em;
    }

    .task-info span {
        font-weight: bolder;
        color: coral;
        font-size: 0.8em;
    }

    .task-info pre {
        font-size: 0.8em;
    }
</style>
