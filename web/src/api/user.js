import service from '../utils/request'

// 用户登录
export const login = (data) => {
    return service({
        url: "/auth/login",
        method: 'post',
        data: data
    })
}

//获取验证码
export const captcha = (data) => {
    return service({
        url: "/auth/captcha",
        method: 'get',
        params: data
    })
}

// 用户列表
export const getUserList = (data) => {
    return service({
        url: "/sys/getUserList",
        method: 'post',
        data: data
    })
}

//更新用户
export const saveUser = (data) => {
    return service({
        url: "/sys/saveUser",
        method: 'post',
        data: data
    })
}

//更新状态
export const setUserStatus = (data) => {
    return service({
        url: "/sys/setUserStatus",
        method: 'post',
        data: data
    })
}


//修改密码
export const modifyPwd = (data) => {
    return service({
        url: "/sys/modifyPwd",
        method: 'post',
        data: data
    })
}