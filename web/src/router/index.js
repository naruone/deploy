import Vue from 'vue'
import Router from 'vue-router'
import {store} from "@/store";
/**
 *  导入子模块路由
 *  !!! 路由name 不要用中划线命名, 菜单输出时用了-作为分隔符
 */
import sysRouter from "./modules/sys"
import serverRouter from "./modules/server"
import projectRouter from "./modules/project";
import taskRouter from "./modules/task";

Vue.use(Router)

const routes = [
    {
        path: '/',
        name: 'login',
        component: () => import('../views/sys/Login'),
        meta: {
            needLogin: false,
            title: '登录',
            menu: false
        }
    },
    sysRouter,
    projectRouter,
    serverRouter,
    taskRouter,
    {
        path: '*',
        name: '404',
        meta: {
            title: 'page not found',
            menu: false
        },
        component: () => import('../views/error/Index')
    }
]

const router = new Router({
    mode: 'hash',
    base: process.env.BASE_URL,
    routes
})

router.beforeEach(async (to, from, next) => {
    if (store.getters['user/expiresAt'] <= (new Date()).getTime()) {
       await store.dispatch('user/LoginOut', false)
    }
    const token = store.getters['user/token']
    document.title = process.env.VUE_APP_TITLE + (to.meta.title ? '-' + to.meta.title : '')
    if (to.meta.needLogin !== false) {  //需要登录
        if (token) {
            next()
        } else {
            next({
                name: "login",
                query: {redirect: to.fullPath}
            })
        }
    } else {    //无需登录
        if (token) {
            next({name: 'sys', replace: true});
        } else {
            next()
        }
    }
});

export default router
