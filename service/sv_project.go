package service

import (
    "deploy/model"
    "deploy/model/request"
)

//项目列表
func ProjectList(search *request.ComPageInfo) ([]model.Project, int, error) {
    return model.GetProjectList(search)
}

//初始化项目
func InitProject(projectId int) (err error) {
    var project model.Project
    if project, err = model.GetProjectById(projectId, true); err != nil {
        return
    }
    go func(pObj *model.Project) {
        pObj.Status = model.ProjectInitProcessing
        _ = model.SaveProject(pObj)
        Repo := GetRepository(pObj)
        if errOut, err, processing := Repo.CloneRepo(); err != nil {
            if processing {
                return
            }
            pObj.Status = model.ProjectInitFail
            pObj.ErrMsg = err.Error() + " " + errOut
        } else {
            pObj.Status = model.ProjectInitSuccess
            pObj.ErrMsg = ""
        }
        _ = model.SaveProject(pObj)
    }(&project)
    return
}

func SaveProject(p *model.Project) (err error) {
    return model.SaveProject(p)
}

func DelProject(projectId int) (err error) {
    if err = model.DelProject(projectId); err != nil {
        return
    }
    _ = model.DeleteEnvAndTask(projectId)
    return
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
    /*
       var (
           task    model.DeployTask
           env     model.EnvProServer
           server  model.Server
           jumper  model.Server
           project model.Project
       )
       if err = global.GDB.Where("task_id = ?", taskId).Find(&task).Error; err != nil {
           return
       }
       if err = global.GDB.Model(&task).Related(&env, "EnvId").Error; err != nil {
           return
       }
       if err = global.GDB.Model(&env).Related(&server, "ServerId").
           Related(&project, "ProjectId").Error; err != nil {
           return
       }
       if err = global.GDB.Model(&env).Related(&jumper, "JumpServer").Error; err != nil {
           jumper = model.Server{}
           err = nil
       }
       if task.Status != model.TaskPrePare {
           err = errors.New("请勿重复操作")
           return
       }
       global.GDB.Model(&task).Update("status", model.TaskStarting)

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
    return
}
