/** When your routing table is too long, you can split it into small modules **/
import Layout from '../../views/layout/Layout'

const sysRouter = {
    path: '/sys',
    name: 'sys',
    component: Layout,
    redirect: '/sys/console',
    meta: {
        title: '系统管理',
        icon: 'el-icon-setting'
    },
    children: [
        {
            path: "console",
            name: 'console',
            meta: {
                title: '仪表盘',
                menu: false
            },
            component: () => import('../../views/sys/Console')
        },
        {
            path: "userList",
            name: 'userList',
            meta: {
                title: '用户表'
            },
            component: () => import('../../views/sys/UserList')
        },
    ]
}
export default sysRouter
