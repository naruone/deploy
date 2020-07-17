<template>
    <div class="process-info">
        <div class="process-info-common">
            <el-divider content-position="left">公共</el-divider>
            <el-steps :space="100" :active="common_active" align-center>
                <el-step title="开始"></el-step>
                <el-step v-for="v in processes" :title="v.v"></el-step>
            </el-steps>
        </div>
        <div class="process-info-server" v-if="processData[task.task_id]">
            <div v-for="(_,k) in processData[task.task_id]['servers']">
                <el-divider content-position="left">{{k}}</el-divider>
                <el-steps :space="100" :active="servers_active[k]" align-center>
                    <el-step v-for="v in server" :title="v.v"></el-step>
                </el-steps>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        name: "ProcessInfo",
        props: {
            task: Object,
            processData: Object
        },
        watch: {
            processData: {
                immediate: true,
                handler(n) {
                    if (n[this.task.task_id] === undefined) {
                        return
                    }
                    let _data = n[this.task.task_id]
                    let _comActive = 1
                    for (let v of this.processes) {
                        if (_data[v.k] !== undefined) {
                            _comActive++
                        }
                    }
                    this.common_active = _comActive
                    for (let k in _data['servers']) {
                        this.$set(this.servers_active, k, 0)
                        for (let v of this.server) {
                            if (_data['servers'][k][v.k] !== undefined) {
                                this.servers_active[k] += 1
                            }
                        }
                    }
                    console.log(this.servers_active)
                    // console.log(n[this.task.task_id]);

                }
            }
        },
        data() {
            return {
                direct: [
                    {
                        k: 'pack',
                        v: '打包'
                    },
                ],
                jumper: [
                    {
                        k: 'pack',
                        v: '打包'
                    },
                    {
                        k: 'upload_jumper',
                        v: '上传跳板'
                    },
                ],
                server: [
                    {
                        k: 'upload_dst',
                        v: '上传包'
                    },
                    {
                        k: 'deploy',
                        v: '发布'
                    },
                    {
                        k: 'change_dir',
                        v: '切目录'
                    },
                ],
                common_active: 0,
                servers_active: {
                    '122.51.244.208': 1
                }
            }
        },
        computed: {
            processes() {
                if (this.task.EnvCfg.Jumper.server_id === 0) {
                    return this.direct
                } else {
                    return this.jumper
                }
            }
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
