package service

import (
    "deploy/config"
    "deploy/model"
    "deploy/model/request"
    "deploy/utils"
    "encoding/json"
    "errors"
    uuid "github.com/satori/go.uuid"
    "path"
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
func Deploy(taskId int) (servers []string, err error) {
    var prepareTask model.TaskPrepare
    servers = nil
    if prepareTask, err = model.PrepareTask(taskId); err != nil {
        return
    }
    if prepareTask.Task.Status != model.TaskPrePare {
        err = errors.New("状态错误,只有未开始的任务才能做发布操作")
        return
    }
    model.UpdateTaskStatusAndOutput(prepareTask.Task.TaskId, map[string]interface{}{"status": model.TaskStarting})
    for _, s := range prepareTask.Servers {
        servers = append(servers, s.SshAddr)
    }
    go deployScheduleStart(&prepareTask)
    return
}

func RollBack(taskId int) (res interface{}, err error) {
    var (
        prepareTask model.TaskPrepare
        wg          sync.WaitGroup
        switchCmd   string
        mapLock     sync.Mutex
        result      = map[string]map[string]string{
            "success": {},
            "error":   {},
        }
    )
    if prepareTask, err = model.PrepareTask(taskId); err != nil {
        return
    }
    if prepareTask.Task.Status != model.TaskRunSuccess {
        err = errors.New("状态错误,只有发布成功的任务才能被回滚")
        return
    }

    if prepareTask.Task.Uuid == "" {
        err = errors.New("抱歉,该版本无法回退")
        return
    }

    if prepareTask.Task.Uuid == prepareTask.Env.Uuid {
        err = errors.New("该环境目前已经是当前版本, 无需切换")
        return
    }

    switchCmd = "ln -snf " + strings.TrimRight(config.GConfig.ServerWorkDir, "/") + "/" + prepareTask.Task.Uuid +
        " " + prepareTask.Project.WebRoot
    wg.Add(len(prepareTask.Servers))
    if prepareTask.Jumper.ServerId != 0 { //跳板机操作
        // 1. 连接跳板机.  2. [并发]执行目标机远程命令
        serverConn := utils.NewServerConn(prepareTask.Jumper.SshAddr+":"+strconv.Itoa(prepareTask.Jumper.SshPort),
            prepareTask.Jumper.SshUser, prepareTask.Jumper.SshKeyPath)
        for _, v := range prepareTask.Servers {
            go func(srv model.Server) {
                if _, err := serverConn.RunCmd(remoteGenCmd(srv, switchCmd)); err != nil {
                    mapLock.Lock()
                    result["error"][srv.SshAddr] = err.Error()
                } else {
                    mapLock.Lock()
                    result["success"][srv.SshAddr] = srv.SshAddr
                }
                mapLock.Unlock()
                wg.Done()
            }(v)
        }
        serverConn.Close()
    } else {
        for _, v := range prepareTask.Servers {
            go func(srv model.Server) {
                serverConn := utils.NewServerConn(srv.SshAddr+":"+strconv.Itoa(srv.SshPort), srv.SshUser, srv.SshKeyPath)
                if _, err := serverConn.RunCmd(switchCmd); err != nil {
                    mapLock.Lock()
                    result["error"][srv.SshAddr] = err.Error()
                } else {
                    mapLock.Lock()
                    result["success"][srv.SshAddr] = srv.SshAddr
                }
                mapLock.Unlock()
                serverConn.Close()
                wg.Done()
            }(v)
        }
    }
    wg.Wait()
    if len(result["error"]) == 0 { //无出错
        model.UpdateEnvRes(prepareTask.Env.EnvId, prepareTask.Task.Version, prepareTask.Task.Uuid)
    }
    res = result
    return
}

//准备工作 & 调度发布任务
func deployScheduleStart(prepareTask *model.TaskPrepare) {
    var (
        params   model.DeployTaskRunParams
        delFiles []string //要删除的文件
        err      error
    )
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
                Output:    "PackErr: " + err.Error(),
            }
        }
        ProcessListenChan <- &TaskProcessReport{
            Task:    prepareTask.Task,
            Process: TaskProcessPack,
            Result:  TaskFail,
        }
        return
    }
    ProcessListenChan <- &TaskProcessReport{
        Task:    prepareTask.Task,
        Process: TaskProcessPack,
        Result:  TaskSuccess,
    }

    //组合目标机script
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
    serverConn = utils.NewServerConn(params.Server.SshAddr+":"+strconv.Itoa(params.Server.SshPort), params.Server.SshUser, params.Server.SshKeyPath)
    defer serverConn.Close()
    //上传包
    if err = serverConn.CopyFile(params.PackagePath, params.DstFilePath); err != nil {
        params.ResChan <- &model.DeployTaskResult{
            Params:    params,
            ResStatus: model.TaskRunFail,
            Output:    "UploadErr: " + err.Error(),
        }
        ProcessListenChan <- &TaskProcessReport{
            Task:    params.Task,
            Server:  params.Server,
            Process: TaskProcessUploadDst,
            Result:  TaskFail,
        }
        return
    }
    ProcessListenChan <- &TaskProcessReport{
        Task:    params.Task,
        Server:  params.Server,
        Process: TaskProcessUploadDst,
        Result:  TaskSuccess,
    }

    if output, err = serverConn.RunCmd(params.DeployCmd); err != nil {
        _msg := "DeployErr: " + err.Error() + "\nCmd: " + params.DeployCmd
        if output != "" {
            _msg += "\nOutput: " + output
        }
        params.ResChan <- &model.DeployTaskResult{
            Params:    params,
            ResStatus: model.TaskRunFail,
            Output:    _msg,
        }
        ProcessListenChan <- &TaskProcessReport{
            Task:    params.Task,
            Server:  params.Server,
            Process: TaskProcessDeploy,
            Result:  TaskFail,
        }
        return
    }
    params.ResChan <- &model.DeployTaskResult{
        Params:    params,
        ResStatus: model.TaskRunSuccess,
        Output:    "Deploy: " + params.DeployCmd,
        SwitchCmd: "ln -snf " + params.DstPath + " " + params.Project.WebRoot,
    }
    ProcessListenChan <- &TaskProcessReport{
        Task:    params.Task,
        Server:  params.Server,
        Process: TaskProcessDeploy,
        Result:  TaskSuccess,
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
    serverConn = utils.NewServerConn(params.Jumper.SshAddr+":"+strconv.Itoa(params.Jumper.SshPort), params.Jumper.SshUser, params.Jumper.SshKeyPath)
    defer serverConn.Close()

    if err = serverConn.CopyFile(params.PackagePath, params.DstFilePath); err != nil {
        for _, _srv := range prepareTask.Servers {
            params.Server = _srv
            params.ResChan <- &model.DeployTaskResult{
                Params:    params,
                ResStatus: model.TaskRunFail,
                Output:    "UploadToJumperErr: " + err.Error(),
            }
        }
        ProcessListenChan <- &TaskProcessReport{
            Task:    params.Task,
            Process: TaskProcessUploadToJumper,
            Result:  TaskFail,
        }
        return
    }
    ProcessListenChan <- &TaskProcessReport{
        Task:    params.Task,
        Process: TaskProcessUploadToJumper,
        Result:  TaskSuccess,
    }

    wg.Add(len(prepareTask.Servers))
    for _, _sv := range prepareTask.Servers {
        params.Server = _sv
        go func(p model.DeployTaskRunParams, svrConn *utils.ServerConn) {
            var (
                commandLog string
                output     string
                _cmd       string
                err        error
            )

            //检测&创建目标机工作目录
            _cmd = remoteGenCmd(p.Server, "([ -d "+path.Dir(p.DstFilePath)+" ] || mkdir -p "+path.Dir(p.DstFilePath)+")")
            commandLog = "CheckDir: " + _cmd
            _, _ = svrConn.RunCmd(_cmd)

            //上传包到目标机
            _cmd = "scp -i " + p.Server.SshKeyPath + " " + p.DstFilePath + " " + p.Server.SshUser + "@" + p.Server.SshAddr + ":" + p.DstFilePath
            commandLog += "\nUpload: " + _cmd
            if output, err = svrConn.RunCmd(_cmd); err != nil {
                _msg := "UploadErr: " + err.Error() + "\nCmd: " + _cmd
                if output != "" {
                    _msg += "\nOutput: " + output
                }
                params.ResChan <- &model.DeployTaskResult{
                    Params:    p,
                    ResStatus: model.TaskRunFail,
                    Output:    _msg,
                }
                ProcessListenChan <- &TaskProcessReport{
                    Task:    params.Task,
                    Server:  params.Server,
                    Process: TaskProcessUploadDst,
                    Result:  TaskFail,
                }
                wg.Done()
                return
            }
            ProcessListenChan <- &TaskProcessReport{
                Task:    params.Task,
                Server:  params.Server,
                Process: TaskProcessUploadDst,
                Result:  TaskSuccess,
            }

            //执行发布
            _cmd = remoteGenCmd(p.Server, p.DeployCmd)
            commandLog += "\nDeploy: " + _cmd
            if output, err = svrConn.RunCmd(_cmd); err != nil {
                _msg := "DeployErr : " + err.Error() + "\nCmd: " + _cmd
                if output != "" {
                    _msg += "\nOutput: " + output
                }
                params.ResChan <- &model.DeployTaskResult{
                    Params:    p,
                    ResStatus: model.TaskRunFail,
                    Output:    _msg,
                }
                ProcessListenChan <- &TaskProcessReport{
                    Task:    params.Task,
                    Server:  params.Server,
                    Process: TaskProcessDeploy,
                    Result:  TaskFail,
                }
                wg.Done()
                return
            }
            params.ResChan <- &model.DeployTaskResult{
                Params:    p,
                ResStatus: model.TaskRunSuccess,
                Output:    commandLog,
                SwitchCmd: "ln -snf " + params.DstPath + " " + params.Project.WebRoot,
            }
            ProcessListenChan <- &TaskProcessReport{
                Task:    params.Task,
                Server:  params.Server,
                Process: TaskProcessDeploy,
                Result:  TaskSuccess,
            }
            wg.Done()
        }(params, serverConn)
    }
    wg.Wait()

    //删除跳板机文件
    _, _ = serverConn.RunCmd("rm -f " + params.DstFilePath)
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
            "message": v.Output + "\nSwitchDir: " + v.SwitchCmd,
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
    } else {
        collectResource(resMap)
    }
    model.UpdateTaskStatusAndOutput(prepareTask.Task.TaskId, updateRes)
    if len(resMap) > 0 { //删除本地打包 && 关闭此任务的ws监听
        utils.DeletePath(resMap[0].Params.PackagePath)
        CloseWsConnectByTaskId(resMap[0].Params.Task.TaskId)
    }
}

//切换工作目录
func switchSymbol(resMap []*model.DeployTaskResult) {
    var (
        _oneRes   *model.DeployTaskResult
        wg        sync.WaitGroup
        rmVersion string
    )
    _oneRes = resMap[0] //随便取一条结果用做判断
    rmVersion = keepVersionDeal(_oneRes.Params.Task, _oneRes.Params.Env.KeepVersionCnt)

    wg.Add(len(resMap))
    if _oneRes.Params.Jumper.ServerId != 0 { //跳板机操作
        // 1. 连接跳板机.  2. [并发]执行目标机远程命令
        serverConn := utils.NewServerConn(_oneRes.Params.Jumper.SshAddr+":"+strconv.Itoa(_oneRes.Params.Jumper.SshPort),
            _oneRes.Params.Jumper.SshUser, _oneRes.Params.Jumper.SshKeyPath)
        for _, v := range resMap {
            go func(res *model.DeployTaskResult) {
                if _, err := serverConn.RunCmd(remoteGenCmd(res.Params.Server, res.SwitchCmd)); err != nil {
                    ProcessListenChan <- &TaskProcessReport{
                        Task:    res.Params.Task,
                        Server:  res.Params.Server,
                        Process: TaskProcessChangeWorkDir,
                        Result:  TaskFail,
                    }
                } else {
                    ProcessListenChan <- &TaskProcessReport{
                        Task:    res.Params.Task,
                        Server:  res.Params.Server,
                        Process: TaskProcessChangeWorkDir,
                        Result:  TaskSuccess,
                    }
                }
                if rmVersion != "" {
                    _, _ = serverConn.RunCmd(remoteGenCmd(res.Params.Server, rmVersion))
                }
                wg.Done()
            }(v)
        }
        serverConn.Close()
    } else {
        for _, r := range resMap {
            go func(_res *model.DeployTaskResult) {
                serverConn := utils.NewServerConn(_res.Params.Server.SshAddr+":"+strconv.Itoa(_res.Params.Server.SshPort),
                    _res.Params.Server.SshUser, _res.Params.Server.SshKeyPath)
                if _, err := serverConn.RunCmd(_res.SwitchCmd); err != nil {
                    ProcessListenChan <- &TaskProcessReport{
                        Task:    _res.Params.Task,
                        Server:  _res.Params.Server,
                        Process: TaskProcessChangeWorkDir,
                        Result:  TaskFail,
                    }
                } else {
                    ProcessListenChan <- &TaskProcessReport{
                        Task:    _res.Params.Task,
                        Server:  _res.Params.Server,
                        Process: TaskProcessChangeWorkDir,
                        Result:  TaskSuccess,
                    }
                }
                if rmVersion != "" {
                    _, _ = serverConn.RunCmd(rmVersion)
                }
                serverConn.Close()
                wg.Done()
            }(r)
        }
    }
    wg.Wait()
}

//回收资源
func collectResource(resMap []*model.DeployTaskResult) {
    var (
        serverConn *utils.ServerConn
        _oneRes    *model.DeployTaskResult
        wg         sync.WaitGroup
    )
    _oneRes = resMap[0]                      //随便取一条结果用做判断
    if _oneRes.Params.Jumper.ServerId != 0 { //跳板机操作
        // 1. 连接跳板机.  2. [并发]执行目标机远程命令
        serverConn = utils.NewServerConn(_oneRes.Params.Jumper.SshAddr+":"+strconv.Itoa(_oneRes.Params.Jumper.SshPort),
            _oneRes.Params.Jumper.SshUser, _oneRes.Params.Jumper.SshKeyPath)
        wg.Add(len(resMap))
        for _, v := range resMap {
            go func(res *model.DeployTaskResult) {
                _, _ = serverConn.RunCmd(remoteGenCmd(res.Params.Server, "rm -rf "+res.Params.DstPath))
                _, _ = serverConn.RunCmd(remoteGenCmd(res.Params.Server, "rm -f "+res.Params.DstFilePath))
                wg.Done()
            }(v)
        }
        wg.Wait()
        serverConn.Close()
        return
    }

    for _, r := range resMap {
        serverConn = utils.NewServerConn(r.Params.Server.SshAddr+":"+strconv.Itoa(r.Params.Server.SshPort),
            r.Params.Server.SshUser, r.Params.Server.SshKeyPath)
        _, _ = serverConn.RunCmd("rm -rf " + r.Params.DstPath)
        _, _ = serverConn.RunCmd("rm -f " + r.Params.DstFilePath)
        serverConn.Close()
    }
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
    deployCmd += " && cd " + params.DstPath
    if len(delFiles) > 0 {
        deployCmd += " && rm -f " + strings.Join(delFiles, " ")
    }
    deployCmd += " && tar -zx --no-same-owner -C " + params.DstPath + " -f " + params.DstFilePath + " && rm -f " + params.DstFilePath

    //after command 切换软链前执行
    _aftScript := ""
    if params.Project.AfterScript != "" { //项目上after_script
        _aftScript = params.Project.AfterScript
    }
    if params.Env.AfterScript != "" { //环境上的after_script
        if _aftScript == "" {
            _aftScript = params.Env.AfterScript
        } else {
            _aftScript += " && " + params.Env.AfterScript
        }
    }
    if params.Task.AfterScript != "" { //任务上的after_script
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

func remoteGenCmd(server model.Server, cmd string) (remoteCmd string) {
    //ssh -T -q -i ~/.ss/key_rsa www@123.123.123.123 '[ ! -d ~/a/b/ ] && mkdir -p ~/a/b/'
    remoteCmd = "ssh -Tq -p " + strconv.Itoa(server.SshPort) + " -i " + server.SshKeyPath + " " +
        server.SshUser + "@" + server.SshAddr + " '" + cmd + "'"
    return
}

//处理保留版本
func keepVersionDeal(task model.DeployTask, keepVersionCnt int) (rmVersions string) {
    rmVersions = ""
    uuids := model.GetNotKeepTask(task.EnvId, task.TaskId, keepVersionCnt)
    if len(uuids) > 0 {
        rmVersions = "rm -rf " + strings.Join(uuids, " ")
    }
    return
}
