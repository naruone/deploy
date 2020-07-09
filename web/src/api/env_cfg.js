import service from '../utils/request'

export const getEnvCfgList = (data) => {
    return service({
        url: "/envCfg/getEnvCfgList",
        method: 'post',
        data: data
    })
}


export const saveEnvCfg = (data) => {
    return service({
        url: "/envCfg/saveEnvCfg",
        method: 'post',
        data: data
    })
}

export const delEnvCfg = (data) => {
    return service({
        url: "/envCfg/delEnvCfg",
        method: 'get',
        params: data
    })
}

export const getCfgSelectOptions = (data) => {
    return service({
        url: "/envCfg/delEnvCfg",
        method: 'get',
        params: data
    })
}

