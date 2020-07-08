import service from '../utils/request'

// 获取服务器列表
export const getServerList = (data) => {
    return service({
        url: "/server/getServerList",
        method: 'post',
        data: data
    })
}

export const saveServer = (data) => {
    return service({
        url: "/server/saveServer",
        method: 'post',
        data: data
    })
}

// 删除项目
export const delServer = (data) => {
    return service({
        url: "/server/deleteServer",
        method: 'get',
        params: data
    })
}
