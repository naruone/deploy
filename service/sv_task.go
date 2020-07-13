package service

import (
    "deploy/config"
    "deploy/model"
    "deploy/model/request"
    "deploy/utils"
    "encoding/json"
    "errors"
    uuid "github.com/satori/go.uuid"
    "strconv"
    "strings"
    "sync"
)

func DeployTaskList(search *request.ComPageInfo) ([]model.DeployTask, int, error) {
    return model.GetDeployTaskList(search)
}

func SaveTask(task *model.DeployTask) error {
    return model.SaveTask(task)
}

func DelTask(taskId int) error {
    return model.DelTask(taskId)
}

//获取分支
func GetBranches(projectId int) (res []string, err error) {
    var project model.Project
    if project, err = model.GetProjectById(projectId, false); err != nil {
        return
    }
    res, err = GetRepository(&project).GetBranches()
    return
}

//获取版本
func GetVersions(projectId int, branch string) (res []model.CsvVersion, err error) {
    var (
        project model.Project
    )
    if project, err = model.GetProjectById(projectId, false); err != nil {
        return
    }
    res, err = GetRepository(&project).GetVersions(branch)
    return
}

//发布
func Deploy(taskId int) (err error) {
    var prepareTask model.TaskPrepare
    if prepareTask, err = model.PrepareTask(taskId); err != nil {
        return
    }
    if prepareTask.Task.Status != model.TaskPrePare {
        err = errors.New("状态错误,只有未开始的任务才能做发布操作")
        return
    }
    go deployScheduleStart(&prepareTask)
    return
}

//准备工作 & 调度发布任务
func deployScheduleStart(prepareTask *model.TaskPrepare) {
    var (
        params   model.DeployTaskRunParams
        delFiles []string //要删除的文件
        err      error
    )
    model.UpdateTaskStatusAndOutput(prepareTask.Task.TaskId, map[string]interface{}{
        "status": model.TaskStarting,
    })
    params = model.DeployTaskRunParams{
        Jumper:      prepareTask.Jumper,
        Task:        prepareTask.Task,
        Env:         prepareTask.Env,
        Project:     prepareTask.Project,
        ResChan:     make(chan *model.DeployTaskResult),
        PackageUuid: uuid.NewV4().String(),
    }
    params.PackageName = params.PackageUuid + ".tar.gz"
    params.DstPath = strings.TrimRight(config.GConfig.ServerWorkDir, "/") + "/" + params.PackageUuid
    params.DstFilePath = strings.TrimRight(config.GConfig.ServerWorkDir, "/") + "/" + params.PackageName
    //全量则上一次版本号不传
    _lastVer := prepareTask.Env.LastVer
    if prepareTask.Task.DeployType == model.DeployTypeAll {
        _lastVer = ""
    }
    //监听发布任务结果
    go deployProcessHandle(params.ResChan, prepareTask)

    //打包
    if params.PackagePath, delFiles, err = GetRepository(&prepareTask.Project).
        Package(_lastVer, prepareTask.Task.Version, params.PackageName); err != nil {
        for _, _srv := range prepareTask.Servers {
            params.Server = _srv
            params.ResChan <- &model.DeployTaskResult{
                Params:    params,
                ResStatus: model.TaskRunFail,
                Output:    "打包失败, 原因: " + err.Error(),
            }
        }
        return
    }

    //组合目标机
    getDeployCmd(&params, delFiles)
    if prepareTask.Jumper.ServerId == 0 { //非跳板机
        for _, _srv := range prepareTask.Servers {
            params.Server = _srv
            go deployStartDirect(params)
        }
    } else { //有跳板逻辑
        go deployStartByJumper(params, prepareTask)
    }
}

//直接发布(非跳板机)
func deployStartDirect(params model.DeployTaskRunParams) {
    var (
        serverConn *utils.ServerConn
        output     string
        err        error
    )
    serverConn = utils.NewServerConn(params.Server.SshAddr+":"+strconv.Itoa(params.Server.SshPort), params.Server.SshUser, params.Server.SshKey)
    defer serverConn.Close()
    //上传包
    if err = serverConn.CopyFile(params.PackagePath, params.DstFilePath); err != nil {
        params.ResChan <- &model.DeployTaskResult{
            Params:    params,
            ResStatus: model.TaskRunFail,
            Output:    "上传包到目标机错误: " + err.Error(),
        }
        return
    }

    if output, err = serverConn.RunCmd(params.DeployCmd); err != nil {
        _msg := "错误原因: " + err.Error() + "\nCommand: " + params.DeployCmd
        if output != "" {
            _msg += "\nOutput: " + output
        }
        params.ResChan <- &model.DeployTaskResult{
            Params:    params,
            ResStatus: model.TaskRunFail,
            Output:    _msg,
        }
        return
    }

    params.ResChan <- &model.DeployTaskResult{
        Params:    params,
        ResStatus: model.TaskRunSuccess,
        Output:    "Command: " + params.DeployCmd,
        SwitchCmd: "ln -snf " + params.DstPath + " " + params.Project.WebRoot,
    }
    return
}

//通过跳板机发布
func deployStartByJumper(params model.DeployTaskRunParams, prepareTask *model.TaskPrepare) {
    var (
        serverConn *utils.ServerConn
        wg         sync.WaitGroup
        err        error
    )
    serverConn = utils.NewServerConn(params.Jumper.SshAddr+":"+strconv.Itoa(params.Jumper.SshPort), params.Jumper.SshUser, params.Jumper.SshKey)
    defer serverConn.Close()

    if err = serverConn.CopyFile(params.PackagePath, params.DstFilePath); err != nil {
        for _, _srv := range prepareTask.Servers {
            params.Server = _srv
            params.ResChan <- &model.DeployTaskResult{
                Params:    params,
                ResStatus: model.TaskRunFail,
                Output:    "上传包到跳板机错误: " + err.Error(),
            }
        }
        return
    }
    wg.Add(len(prepareTask.Servers))
    for _, _sv := range prepareTask.Servers {
        params.Server = _sv
        go func(p model.DeployTaskRunParams) {
            // ssh -Tq -p 22 -i ~/.ssh/id_rsa root@ip 'shell'

            //if output, err = serverConn.RunCmd(p.DeployCmd); err != nil {
            //    result.ResStatus = model.TaskRunFail
            //    result.Output = "错误原因: " + err.Error() + "\nCommand: " + deployCmd
            //    if output != "" {
            //        result.Output = result.Output + "\nOutput: " + output
            //    }
            //    goto DepErr
            //}
            //result.ResStatus = model.TaskRunSuccess
            //result.Output = "Command: " + deployCmd
            //result.SwitchCmd = "ln -snf " + dstDir + " " + params.Project.WebRoot
            //params.ResChan <- &result

            wg.Done()
        }(params)
    }
    wg.Wait()

    return
}

//监听发布任务结果
func deployProcessHandle(resChan chan *model.DeployTaskResult, prepareTask *model.TaskPrepare) {
    var (
        _taskResult *model.DeployTaskResult
        resMap      []*model.DeployTaskResult
        updateRes   = make(map[string]interface{})
        messages    = make(map[string]interface{})
        _status     = model.TaskRunSuccess
        _uuid       string
    )
    for i := 0; i < len(prepareTask.Servers); i++ {
        select {
        case _taskResult = <-resChan:
            resMap = append(resMap, _taskResult)
        }
    }
    for _, v := range resMap {
        if _status == model.TaskRunSuccess && v.ResStatus != model.TaskRunSuccess {
            _status = model.TaskRunFail
        }
        messages[v.Params.Server.SshAddr] = map[string]interface{}{
            "status":  v.ResStatus,
            "message": v.Output + "\n切换: " + v.SwitchCmd,
        }
    }
    updateRes["status"] = _status
    updateRes["output"], _ = json.Marshal(messages)
    if _status == model.TaskRunSuccess {
        if len(resMap) > 0 {
            _uuid = resMap[0].Params.PackageUuid
            updateRes["uuid"] = _uuid
        }
        switchSymbol(resMap)
        model.UpdateEnvRes(prepareTask.Env.EnvId, prepareTask.Task.Version, _uuid)
    }
    model.UpdateTaskStatusAndOutput(prepareTask.Task.TaskId, updateRes)
}

//生成目标机脚本
func getDeployCmd(params *model.DeployTaskRunParams, delFiles []string) {
    var deployCmd string
    if params.Task.DeployType == model.DeployTypeIncrease && params.Env.Uuid != "" {
        _resDir := strings.TrimRight(config.GConfig.ServerWorkDir, "/") + "/" + params.Env.Uuid
        deployCmd = "([ ! -d " + _resDir + " ] || cp -r " + _resDir + " " + params.DstPath + ") && ([ -d " + params.DstPath + " ] || mkdir -p " + params.DstPath + ")"
    } else {
        deployCmd = "mkdir -p " + params.DstPath
    }
    deployCmd += " && tar -zx --no-same-owner -C " + params.DstPath + " -f " + params.DstFilePath + " && rm -f " +
        params.DstFilePath + " && cd " + params.DstPath

    //after command 切换软链前执行
    if len(delFiles) > 0 {
        deployCmd += " && rm -f " + strings.Join(delFiles, " ")
    }
    _aftScript := ""
    if params.Project.AfterScript != "" {
        _aftScript = params.Project.AfterScript
    }
    if params.Task.AfterScript != "" {
        if _aftScript == "" {
            _aftScript = params.Task.AfterScript
        } else {
            _aftScript += " && " + params.Task.AfterScript
        }
    }
    _aftScript = strings.Replace(_aftScript, "{web_root}", params.DstPath, -1)
    _aftScript = strings.Replace(_aftScript, "{version}", params.Task.Version, -1)
    if _aftScript != "" {
        deployCmd += " && " + _aftScript
    }
    params.DeployCmd = deployCmd
    return
}

func switchSymbol(resMap []*model.DeployTaskResult) {
    var (
        serverConn *utils.ServerConn
    )
    if resMap[0].Params.Jumper.ServerId != 0 { //跳板机操作
        // todo 跳板机切换
        // 1. 连接跳板机.  2. [并发]执行目标机远程命令
        return
    }
    for _, r := range resMap {
        serverConn = utils.NewServerConn(r.Params.Server.SshAddr+":"+strconv.Itoa(r.Params.Server.SshPort),
            r.Params.Server.SshUser, r.Params.Server.SshKey)
        _, _ = serverConn.RunCmd(r.SwitchCmd)
        serverConn.Close()
    }
}
