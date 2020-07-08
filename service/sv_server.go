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
    return model.SaveServer(s)
}

func DelServer(serverId int) (err error) {
    if err = model.IsServerUsed(serverId); err != nil {
        return
    }
    err = model.DelServer(serverId)
    return
}
