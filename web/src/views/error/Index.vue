<template>
    <div class="error-page">
        <div class="inner">
            <h3><span style="color:#578ec5"><i class="el-icon-lightning"></i>404</span> 页面不存在!</h3>
            <p>Sorry! 系统找不到对应页面!</p>

            <div class="point">
                <span>你可以尝试以下操作:</span>
                <p><i class="el-icon-location-information"></i> 检查链接是否有误</p>
                <p @click="goHome"><i class="el-icon-location-information"></i> 跳转到主页</p>
            </div>
            <div class="auto-jump">
                {{sec}} 秒后自动跳转到 <span @click="goHome">主页</span>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        name: "Error",
        data() {
            return {
                sec: 5
            }
        },
        methods: {
            goHome() {
                this.$router.push({name: 'sys'})
            }
        },
        mounted() {
            const timer = setInterval(() => {
                this.sec--
                if (this.sec === 0){
                    clearInterval(timer)
                    this.goHome()
                }
            }, 1000)
        }
    }
</script>

<style scoped>
    .error-page {
        width: 100%;
        height: 100vh;
        color: #666;
        cursor: default;
        background-color: #eee;
    }

    .inner {
        width: 500px;
        height: 330px;
        margin: 100px auto;
        border: 1px solid #CCCCCC;
        border-radius: 3px;
        box-shadow: 2px 2px 3px #CCC;
    }

    .inner h3 {
        font-weight: normal;
        margin: 20px 30px;
        height: 65px;
        line-height: 65px;
        font-size: 2em;
        border-bottom: 1px solid #CCC;
    }

    .inner > p {
        margin: 20px 30px;
        font-size: 1.2em;
    }

    .inner .point {
        margin: 30px 30px auto 30px;
        border-bottom: 1px solid #CCC;
    }

    .inner .point p {
        font-size: 0.9em;
        padding: 5px;
    }

    .inner .point p:last-child, .auto-jump span {
        cursor: pointer;
        color: #578ec5;
    }

    .inner .point p:last-child:hover, .auto-jump span:hover {
        text-decoration: underline;
    }

    .auto-jump {
        margin-top: 17px;
        font-size: 0.8em;
        text-align: center;
    }
</style>