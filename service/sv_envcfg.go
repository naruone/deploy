package service

import (
    "deploy/model"
    "deploy/model/request"
)

//项目列表
func EnvCfgList(search *request.ComPageInfo) ([]model.EnvProServer, int, error) {
    return model.GetEnvCfgList(search)
}

func SaveEnvCfg(env *model.EnvProServer) (err error) {
    return model.SaveEnvCfg(env)
}

func DelEnvCfg(cfgId int) (err error) {
    if err = model.CheckDelTask(cfgId); err != nil {
        return
    }
    if err = model.DelTaskByEnvId(cfgId); err != nil {
        return
    }
    return model.DelEnvCfg(cfgId)
}
