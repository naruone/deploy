import service from '../utils/request'

export const getTaskList = (data) => {
    return service({
        url: "/task/getTaskList",
        method: 'post',
        data: data
    })
}

export const getEnvOptions = (data) => {
    return service({
        url: "/task/getEnvOptions",
        method: 'get',
        params: data
    })
}

export const saveTask = (data) => {
    return service({
        url: "/task/saveTask",
        method: 'post',
        data: data
    })
}

export const deleteTask = (data) => {
    return service({
        url: "/task/deleteTask",
        method: 'get',
        params: data
    })
}

export const getBranches = (data) => {
    return service({
        url: "/task/getBranches",
        method: 'get',
        params: data
    })
}

export const getVersions = (data) => {
    return service({
        url: "/task/getVersions",
        method: 'get',
        params: data
    })
}
