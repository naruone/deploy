package service

import (
    "deploy/model"
    "deploy/model/request"
)

//项目列表
func ServerList(search *request.ComPageInfo) ([]model.Server, int, error) {
    return model.GetServerList(search)
}

func SaveServer(s *model.Server) (err error) {
    if s.ServerId != 0 {
        server := model.GetServerById(s.ServerId)
        if server.Type != s.Type { //如果修改了服务器类型则需要判断当前类型是否被使用
            if err = model.IsServerUsed(s.ServerId); err != nil {
                return
            }
        }
    }
    return model.SaveServer(s)
}

func DelServer(serverId int) (err error) {
    if err = model.IsServerUsed(serverId); err != nil {
        return
    }
    err = model.DelServer(serverId)
    return
}
