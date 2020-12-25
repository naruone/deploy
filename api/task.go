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

func GetTaskList(c *gin.Context) {
    var (
        pageInfo request.ComPageInfo
        list     []model.DeployTask
        total    int
        err      error
    )
    _ = c.ShouldBindJSON(&pageInfo)
    if list, total, err = service.DeployTaskList(&pageInfo); err != nil {
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

func DeleteTask(c *gin.Context) {
    var (
        taskId int
        err    error
    )
    if taskId, err = strconv.Atoi(c.Query("task_id")); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if err = service.DelTask(taskId); err != nil {
        utils.FailWithMessage("删除失败, 原因:"+err.Error(), c)
        return
    }
    utils.OkWithMessage("删除成功", c)
}

func GetEnvOptions(c *gin.Context) {
    var (
        search  *request.ComPageInfo
        envList []model.EnvProServer
        err     error
    )
    search = &request.ComPageInfo{
        BasePageInfo: request.BasePageInfo{
            CurrentPage: 1,
            PageSize:    10000,
            SearchValue: "",
            Condition:   "",
        },
        Status: 0,
        Type:   0,
    }
    if envList, _, err = service.EnvCfgList(search); err != nil {
        utils.FailWithMessage("获取失败,"+err.Error(), c)
        return
    }
    utils.OkWithData(map[string]interface{}{
        "envList": envList,
    }, c)
}

func GetVersions(c *gin.Context) {
    var (
        res       []model.CsvVersion
        projectId int
        branch    string
        err       error
    )
    if projectId, err = strconv.Atoi(c.Query("project_id")); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    branch = c.Query("branch")
    if res, err = service.GetVersions(projectId, branch); err != nil {
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
        servers []string
        err     error
        taskId  int
    )
    if taskId, err = strconv.Atoi(c.Query("task_id")); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }

    if servers, err = service.Deploy(taskId); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    utils.OkWithData(servers, c)
}

func SaveTask(c *gin.Context) {
    var (
        task       model.DeployTask
        taskVerify map[string][]string
        err        error
    )
    _ = c.ShouldBindJSON(&task)
    taskVerify = utils.Rules{
        "TaskName":       {utils.NotEmpty(), utils.Le("200"), utils.Ge("3")},
        "Branch":         {utils.NotEmpty()},
        "Version":        {utils.NotEmpty(), utils.Ge("1"), utils.Le("65535")},
        "Description":    {utils.Le("5000")},
        "AfterScript":    {utils.Le("5000")},
        "KeepVersionCnt": {utils.NotEmpty(), utils.Le("20"), utils.Ge("2")},
    }
    if err = utils.Verify(task, taskVerify); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if err = service.SaveTask(&task); err != nil {
        utils.FailWithMessage("保存失败,"+err.Error(), c)
        return
    }
    utils.OkWithMessage("保存成功", c)
}

func RollBack(c *gin.Context) {
    var (
        taskId int
        res    interface{}
        err    error
    )
    if taskId, err = strconv.Atoi(c.Query("task_id")); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    if res, err = service.RollBack(taskId); err != nil {
        utils.FailWithMessage(err.Error(), c)
        return
    }
    utils.OkWithData(res, c)
}

func DeployInfo(c *gin.Context) {
    service.WebSocketHandler(c.Writer, c.Request)
}
