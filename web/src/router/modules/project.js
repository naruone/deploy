/** When your routing table is too long, you can split it into small modules **/
import Layout from '../../views/layout/Layout'

const projectRouter = {
    path: '/project',
    name: 'project',
    component: Layout,
    redirect: '/project/projectList',
    meta: {
        title: '项目管理',
        icon: 'el-icon-wallet'
    },
    children: [
        {
            path: "projectList",
            name: 'projectList',
            meta: {
                title: '项目列表'
            },
            component: () => import('../../views/project/Project')
        },
    ]
}
export default projectRouter
