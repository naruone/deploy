<template>
    <el-breadcrumb separator-class="el-icon-arrow-right">
        <template v-for="i in this.$route.matched">
            <el-breadcrumb-item v-if="checkUrl(i)" :to="i.path">{{ i.meta.title }}
            </el-breadcrumb-item>
            <el-breadcrumb-item v-else>{{ i.meta.title }}</el-breadcrumb-item>
        </template>

    </el-breadcrumb>
</template>

<script>
export default {
    name: "BreadCrumb",
    methods: {
        checkUrl(i) {
            return this.$route.path !== (i.redirect ? i.redirect : i.path)
        }
    }
}
</script>

<style scoped>

</style>
