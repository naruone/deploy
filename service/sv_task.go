package service

import (
    "deploy/model"
    "deploy/model/request"
    "deploy/utils"
    "encoding/json"
    "errors"
    uuid "github.com/satori/go.uuid"
    "strconv"
    "strings"
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
        params     model.DeployTaskRunParams
        delFiles   []string //要删除的文件
        resultChan = make(chan *model.DeployTaskResult)
        err        error
    )
    model.UpdateTaskStatusAndOutput(prepareTask.Task.TaskId, map[string]interface{}{
        "status": model.TaskStarting,
    })
    params = model.DeployTaskRunParams{
        Task:        prepareTask.Task,
        Env:         prepareTask.Env,
        Project:     prepareTask.Project,
        ResChan:     resultChan,
        PackageUuid: uuid.NewV4().String(),
    }
    params.PackageName = params.PackageUuid + ".tar.gz"

    //全量则上一次版本号不传
    _lastVer := prepareTask.Env.LastVer
    if prepareTask.Task.DeployType == model.DeployTypeAll {
        _lastVer = ""
    }
    //打包
    if params.PackagePath, delFiles, err = GetRepository(&prepareTask.Project).
        Package(_lastVer, prepareTask.Task.Version, params.PackageName); err != nil {
        model.UpdateTaskStatusAndOutput(prepareTask.Task.TaskId, map[string]interface{}{
            "status": model.TaskRunFail,
            "output": "打包失败, 原因: " + err.Error(),
        })
        return
    }

    //after command 切换软链前执行
    if len(delFiles) > 0 {
        params.AfterScript = "rm -f " + strings.Join(delFiles, " ")
    }
    if prepareTask.Project.AfterScript != "" {
        if params.AfterScript == "" {
            params.AfterScript = prepareTask.Project.AfterScript
        } else {
            params.AfterScript = params.AfterScript + " && " + prepareTask.Project.AfterScript
        }
    }
    if prepareTask.Task.AfterScript != "" {
        if params.AfterScript == "" {
            params.AfterScript = prepareTask.Task.AfterScript
        } else {
            params.AfterScript = params.AfterScript + " && " + prepareTask.Task.AfterScript
        }
    }

    if prepareTask.Jumper.ServerId == 0 { //非跳板机
        go deployProcessHandle(resultChan, prepareTask) //监听发布任务结果
        for _, _srv := range prepareTask.Servers {
            params.Server = _srv
            go deployStart(params)
        }
    } else {
        //todo 有跳板逻辑
    }
}

func deployStart(params model.DeployTaskRunParams) {
    var (
        serverConn *utils.ServerConn
        dstFile    string //目标机文件上传路径
        result     model.DeployTaskResult
        dstDir     string //解压包路径
        deployCmd  string
        output     string
        err        error
    )
    result = model.DeployTaskResult{
        Server:      params.Server,
        PackagePath: params.PackagePath,
    }
    serverConn = utils.NewServerConn(params.Server.SshAddr+":"+strconv.Itoa(params.Server.SshPort), params.Server.SshUser, params.Server.SshKey)
    defer serverConn.Close()

    dstFile = strings.TrimRight(params.Server.WorkDir, "/") + "/" + params.PackageName
    //上传包
    if err = serverConn.CopyFile(params.PackagePath, dstFile); err != nil {
        result.ResStatus = model.TaskRunFail
        result.Output = "错误原因: " + err.Error()
        goto DepErr
    }
    //增量发布则尝试依赖上一次项目
    dstDir = strings.TrimRight(params.Server.WorkDir, "/") + "/" + params.PackageUuid
    if params.Task.DeployType == model.DeployTypeIncrease && params.Env.Uuid != "" {
        _resDir := strings.TrimRight(params.Server.WorkDir, "/") + "/" + params.Env.Uuid
        deployCmd = "([ ! -d " + _resDir + " ] || cp -r " + _resDir + " " + dstDir + ") && ([ ! -d " + dstDir + " ] || mkdir -p " + dstDir + ")"
    } else {
        deployCmd = "mkdir -p " + dstDir
    }
    deployCmd = deployCmd + " && tar -zx --no-same-owner -C " + dstDir + " -f " + dstFile + " && rm -f " + dstFile + " && cd " + dstDir
    params.AfterScript = strings.Replace(params.AfterScript, "{web_root}", dstDir, -1)
    params.AfterScript = strings.Replace(params.AfterScript, "{version}", params.Task.Version, -1)
    if params.AfterScript != "" {
        deployCmd = deployCmd + " && " + params.AfterScript
    }
    deployCmd = deployCmd + " && ln -snf " + dstDir + " " + params.Project.WebRoot //todo 切换软连后期改为结果里统一执行
    if output, err = serverConn.RunCmd(deployCmd); err != nil {
        result.ResStatus = model.TaskRunFail
        result.Output = "错误原因: " + err.Error() + "\nCommand: " + deployCmd
        if output != "" {
            result.Output = result.Output + "\nOutput: " + output
        }
        goto DepErr
    }
    result.ResStatus = model.TaskRunSuccess
    result.Output = "Command: " + deployCmd
    params.ResChan <- &result
    return
DepErr:
    params.ResChan <- &result
    _, _ = serverConn.RunCmd("rm -f " + dstFile)
    _, _ = serverConn.RunCmd("rm -rf " + dstDir)
}

//监听发布任务结果
func deployProcessHandle(resChan chan *model.DeployTaskResult, prepareTask *model.TaskPrepare) {
    var (
        _taskResult *model.DeployTaskResult
        resMap      []*model.DeployTaskResult
        updateRes   = make(map[string]interface{})
        messages    = make(map[string]interface{})
        _status     = model.TaskRunSuccess
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
        messages[v.Server.SshAddr] = map[string]interface{}{
            "status":  v.ResStatus,
            "message": v.Output,
        }
    }
    updateRes["status"] = _status
    updateRes["output"], _ = json.Marshal(messages)

    //todo 可以将切换软连接放在这里
    model.UpdateTaskStatusAndOutput(prepareTask.Task.TaskId, updateRes)
}

/*
   go func() {
       var (
           serverConn *utils.ServerConn
           nameUuid   string
           name       string   //包名 不带全路径
           filename   string   //全路径包名
           delFiles   []string //删除的文件
           dstFile    string   //目标服务器目录
           output     string
           _dstDir    string //目标服务器目录
           _dCmd      string //目标服务器拷贝文件夹command
           _afterCmd  string
       )
       nameUuid = uuid.NewV4().String()
       name = nameUuid + ".tar.gz"

       //全量则上一次版本号不传
       lastVer := env.LastVer
       if task.DeployType == model.DeployTypeAll {
           lastVer = ""
       }
       //打包
       if filename, delFiles, err = utils.GetRepository(&project).Package(lastVer, task.Version, name); err != nil {
           goto DepErr
       }
       //建立连接
       serverConn = utils.NewServerConn(server.SshAddr+":"+strconv.Itoa(server.SshPort), server.SshUser, server.SshKey)
       defer serverConn.Close()

       dstFile = strings.TrimRight(server.WorkDir, "/") + "/" + name
       //上传包
       if err = serverConn.CopyFile(filename, dstFile); err != nil {
           goto DepErr
       }
       //增量发布则尝试依赖上一次项目
       _dstDir = strings.TrimRight(server.WorkDir, "/") + "/" + nameUuid
       if task.DeployType == model.DeployTypeIncrease && env.Uuid != "" {
           _resDir := strings.TrimRight(server.WorkDir, "/") + "/" + env.Uuid
           _dCmd = "([ ! -d " + _resDir + " ] || cp -r " + _resDir + " " + _dstDir + ") && ([ ! -d " + _dstDir +
               " ] || mkdir -p " + _dstDir + ")"
           if len(delFiles) > 0 {
               _dCmd = _dCmd + " && cd " + _dstDir + " && rm -f " + strings.Join(delFiles, " ")
           }
       } else {
           _dCmd = "mkdir -p " + _dstDir
       }
       _dCmd = _dCmd + " && tar -zx --no-same-owner -C " + _dstDir + " -f " + dstFile + " && rm -f " + dstFile + " && cd " + _dstDir
       //after command 切换软链前执行
       if project.AfterScript != "" {
           _afterCmd = project.AfterScript
       }
       if task.AfterScript != "" {
           if _afterCmd == "" {
               _afterCmd = strings.Replace(task.AfterScript, "{web_root}", _dstDir, -1)
           } else {
               _afterCmd = _afterCmd + " && " + task.AfterScript
           }
       }
       _afterCmd = strings.Replace(_afterCmd, "{web_root}", _dstDir, -1)
       _afterCmd = strings.Replace(_afterCmd, "{version}", task.Version, -1)
       if _afterCmd != "" {
           _dCmd = _dCmd + " && " + _afterCmd
       }
       _dCmd = _dCmd + " && ln -snf " + _dstDir + " " + project.WebRoot
       fmt.Println(_dCmd)
       if output, err = serverConn.RunCmd(_dCmd); err != nil {
           goto DepErr
       }
       global.GDB.Model(&env).Updates(map[string]interface{}{
           "last_ver": task.Version,
           "uuid":     nameUuid,
       })
       global.GDB.Model(&task).Updates(map[string]interface{}{
           "uuid":   nameUuid,
           "status": model.TaskRunSuccess,
       })
       utils.DeleteFile(filename)
       return
   DepErr:
       global.GDB.Model(&task).Updates(map[string]interface{}{
           "err_output": err.Error(),
           "status":     model.ProjectInitFail,
       })
       fmt.Println(output)
       utils.DeleteFile(filename)
       if serverConn != nil {
           _, _ = serverConn.RunCmd("rm -f " + dstFile)
           _, _ = serverConn.RunCmd("rm -rf " + _dstDir)
       }
   }()
*/
