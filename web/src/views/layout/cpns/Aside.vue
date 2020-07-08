<template>
    <el-scrollbar style="height: calc(100vh - 60px)">
        <transition :duration="{ enter: 800, leave: 100 }" mode="out-in" name="el-fade-in-linear">
            <el-menu
                    :collapse="isCollapse"
                    class="menu-contain"
                    :collapse-transition="true"
                    :default-active="active"
                    @select="goToRoute"
                    unique-opened>

                <template v-for="r in $router.options.routes">
                    <template v-if="r.children && r.children.length">
                        <el-submenu v-if="r.meta.menu !== false" :index="r.name">
                            <template slot="title">
                                <i :class="r.meta.icon ? r.meta.icon : 'el-icon-document'"></i>
                                <span>{{r.meta.title}}</span>
                            </template>

                            <template v-for="cr in r.children">
                                <el-menu-item v-if="cr.meta.menu !== false" :index="r.name + '-' + cr.name">
                                    <i :class="cr.meta.icon ? cr.meta.icon : 'el-icon-document'"></i>
                                    {{cr.meta.title}}
                                </el-menu-item>
                            </template>
                        </el-submenu>
                    </template>
                    <el-menu-item v-else-if="r.meta.menu !== false" :index="r.name">
                        <i :class="r.meta.icon ? r.meta.icon : 'el-icon-document'"></i>
                        <span slot="title">{{r.meta.title}}</span>
                    </el-menu-item>
                </template>
            </el-menu>
        </transition>
    </el-scrollbar>
</template>

<script>
    export default {
        name: "Aside",
        props: {
            isCollapse: Boolean
        },
        data() {
            return {
                active: ''
            }
        },
        watch: {
            $route() {
                this.calcActive()
            }
        },
        created() {
            this.calcActive()
        },
        methods: {
            goToRoute(index) {
                const routeName = index.split('-').pop()
                if (this.$route.name !== routeName) {
                    this.$router.push({name: String(routeName)})
                }
            },
            calcActive() {
                this.active = this.$route.matched.reduce((preVal, n) => (preVal !== '' ? preVal + '-' : '') + n.name, '')
            }
        }
    }
</script>

<style scoped>
    .el-menu {
        height: calc(100vh - 60px);
    }

    .menu-contain > li.is-active {
        background-color: #eef5fe;
    }
</style>