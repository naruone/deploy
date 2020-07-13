package api

import (
    "deploy/model"
    "deploy/model/request"
    "deploy/model/response"
    "deploy/service"
    "deploy/utils"
    "github.com/gin-gonic/gin"
    "strconv"
)

func GetServerList(c *gin.Context) {
    var (
        pageInfo request.ComPageInfo
        list     []model.Server
        total    int
        err      error
    )
    _ = c.ShouldBindJSON(&pageInfo)
    if list, total, err = service.ServerList(&pageInfo); err != nil {
        utils.FailWithMessage("获取失败, Message: "+err.Error(), c)
        return
    }
    utils.OkDetailed(response.PageResult{
        List:        list,
        Total:       total,
        PageSize:    pageInfo.PageSize,
        CurrentPage: pageInfo.CurrentPage,
    }, "获取成功", c)
}

func SaveServer(c *gin.Context) {
    type sshKey struct {
        model.Server
        SshKey string `json:"ssh_key"`
    }
    var (
        server       model.Server
        serverVerify map[string][]string
        err          error
    )
    var _sshKey sshKey
    _ = c.ShouldBindJSON(&_sshKey)
    server = _sshKey.Server
    server.SshKey = _sshKey.SshKey
    serverVerify = utils.Rules{
        "Type":    {utils.NotEmpty(), utils.Le("2"), utils.Ge("1")},
        "SshAddr": {utils.NotEmpty()},
        "SshPort": {utils.NotEmpty(), utils.Ge("1"), utils.Le("65535")},
        "SshUser": {utils.NotEmpty(), utils.Ge("3"), utils.Le("200")},
        "SshKey":  {utils.Le("5000")},
    }
    if err = utils.Verify(server, serverVerify); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if err = service.SaveServer(&server); err != nil {
        utils.FailWithMessage("保存失败,"+err.Error(), c)
        return
    }
    utils.OkWithMessage("保存成功", c)
}

func DeleteServer(c *gin.Context) {
    var (
        serverId int
        err      error
    )
    if serverId, err = strconv.Atoi(c.Query("server_id")); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if err = service.DelServer(serverId); err != nil {
        utils.FailWithMessage("删除失败, 原因:"+err.Error(), c)
        return
    }
    utils.OkWithMessage("删除成功", c)
}
