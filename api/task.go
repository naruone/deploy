package api

import (
    "deploy/model"
    "deploy/utils"
    "github.com/gin-gonic/gin"
    "strconv"
)

//任务相关
import (
    "deploy/service"
)

func GetVersions(c *gin.Context) {
    var (
        res       []model.CsvVersion
        projectId int
        err       error
    )
    if projectId, err = strconv.Atoi(c.Query("project_id")); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if res, err = service.GetVersions(projectId, "master"); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    utils.OkDetailed(res, "获取版本成功", c)
}

func GetBranches(c *gin.Context) {
    var (
        res       []string
        projectId int
        err       error
    )
    if projectId, err = strconv.Atoi(c.Query("project_id")); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if res, err = service.GetBranches(projectId); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    utils.OkDetailed(res, "获取分支成功", c)
}

func Deploy(c *gin.Context) {
    var (
        err       error
        projectId int
    )
    if projectId, err = strconv.Atoi(c.Query("project_id")); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }

    if err = service.Deploy(projectId); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    utils.Ok(c)
}
