import service from '../utils/request'

// 获取项目列表
export const getProjectList = (data) => {
    return service({
        url: "/project/getProjectList",
        method: 'post',
        data: data
    })
}

// 初始化项目
export const initProject = (data) => {
    return service({
        url: "/project/initProject",
        method: 'get',
        params: data
    })
}

export const saveProject = (data) => {
    return service({
        url: "/project/saveProject",
        method: 'post',
        data: data
    })
}

// 删除项目
export const delProject = (data) => {
    return service({
        url: "/project/delProject",
        method: 'get',
        params: data
    })
}