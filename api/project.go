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

func SaveProject(c *gin.Context) {
    var (
        p             model.Project
        projectVerify map[string][]string
        err           error
    )
    _ = c.ShouldBindJSON(&p)
    projectVerify = utils.Rules{
        "ProjectName": {utils.NotEmpty(), utils.Le("200"), utils.Ge("3")},
        "RepoUrl":     {utils.NotEmpty(), utils.Le("255"), utils.Ge("6")},
        "Dst":         {utils.NotEmpty(), utils.Le("200"), utils.Ge("3")},
        "WebRoot":     {utils.NotEmpty(), utils.Le("200"), utils.Ge("3")},
        "AfterScript": {utils.Le("5000")},
    }
    if err = utils.Verify(p, projectVerify); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if err = service.SaveProject(&p); err != nil {
        utils.FailWithMessage("保存失败,"+err.Error(), c)
        return
    }
    utils.OkWithMessage("保存成功", c)
}

//初始化项目git仓库
func InitProject(c *gin.Context) {
    var (
        projectId int
        err       error
    )
    if projectId, err = strconv.Atoi(c.Query("project_id")); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if err := service.InitProject(projectId); err != nil {
        utils.FailWithMessage(err.Error(), c)
    } else {
        utils.OkWithMessage("已开始初始化, 请稍后查看状态", c)
    }
}

func GetProjectList(c *gin.Context) {
    var (
        comPageInfo request.ComPageInfo
        list        []model.Project
        total       int
        err         error
    )
    _ = c.ShouldBindJSON(&comPageInfo)
    if list, total, err = service.ProjectList(&comPageInfo); err != nil {
        utils.FailWithMessage("获取失败, Message: "+err.Error(), c)
        return
    }
    utils.OkDetailed(response.PageResult{
        List:        list,
        Total:       total,
        PageSize:    comPageInfo.PageSize,
        CurrentPage: comPageInfo.CurrentPage,
    }, "获取成功", c)
}

func DelProject(c *gin.Context) {
    var (
        projectId int
        err       error
    )
    if projectId, err = strconv.Atoi(c.Query("project_id")); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if err = service.DelProject(projectId); err != nil {
        utils.FailWithMessage("删除失败, 原因:"+err.Error(), c)
        return
    }
    utils.OkWithMessage("删除成功", c)
}
