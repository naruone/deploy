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

func GetEnvCfgList(c *gin.Context) {
    var (
        pageInfo request.ComPageInfo
        list     []model.EnvProServer
        total    int
        err      error
    )
    _ = c.ShouldBindJSON(&pageInfo)
    if list, total, err = service.EnvCfgList(&pageInfo); err != nil {
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

func SaveEnvCfg(c *gin.Context) {
    var (
        env       model.EnvProServer
        envVerify map[string][]string
        err       error
    )
    _ = c.ShouldBindJSON(&env)

    envVerify = utils.Rules{
        "EnvName":    {utils.NotEmpty(), utils.Le("200"), utils.Ge("3")},
        "ProjectId":  {utils.NotEmpty()},
        "ServerId":   {utils.NotEmpty()},
        "JumpServer": {utils.NotEmpty()},
    }
    if err = utils.Verify(env, envVerify); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if err = service.SaveEnvCfg(&env); err != nil {
        utils.FailWithMessage("保存失败,"+err.Error(), c)
        return
    }
    utils.OkWithMessage("保存成功", c)
}

func DelEnvCfg(c *gin.Context) {
    var (
        cfgId int
        err   error
    )
    if cfgId, err = strconv.Atoi(c.Query("cfg_id")); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if err = service.DelEnvCfg(cfgId); err != nil {
        utils.FailWithMessage("删除失败, 原因:"+err.Error(), c)
        return
    }
    utils.OkWithMessage("删除成功", c)
}
