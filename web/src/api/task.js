import service from '../utils/request'

export const getTaskList = (data) => {
    return service({
        url: "/task/getTaskList",
        method: 'post',
        data: data
    })
}
