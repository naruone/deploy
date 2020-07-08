/** When your routing table is too long, you can split it into small modules **/
import Layout from '../../views/layout/Layout'

const serverRouter = {
    path: '/server',
    name: 'server',
    component: Layout,
    redirect: '/server/serverList',
    meta: {
        title: '服务器管理',
        icon: 'el-icon-cpu'
    },
    children: [
        {
            path: "serverList",
            name: 'serverList',
            meta: {
                title: '服务器列表'
            },
            component: () => import('../../views/server/ServerList')
        },
    ]
}
export default serverRouter
