package model

import (
    "deploy/model/request"
    "errors"
    "strings"
    "time"
)

type DeployTask struct {
    TaskId      int       `json:"task_id" gorm:"PRIMARY_KEY"`
    TaskName    string    `json:"task_name"`
    DeployType  int       `json:"deploy_type"`
    Description string    `json:"description"`
    Output      string    `json:"output"`
    EnvId       int       `json:"env_id"`
    Branch      string    `json:"branch"`
    Version     string    `json:"version"`
    AfterScript string    `json:"after_script"`
    Status      int       `json:"status"`
    Uuid        string    `json:"uuid"`
    CreateAt    time.Time `json:"create_at"`
    UpdateAt    time.Time `json:"update_at"`
    EnvCfg      EnvProServer
}

type EnvProServer struct {
    EnvId      int       `json:"env_id" gorm:"PRIMARY_KEY"`
    EnvName    string    `json:"env_name"`
    ProjectId  int       `json:"project_id"`
    ServerIds  string    `json:"server_ids"`
    JumpServer int       `json:"jump_server"`
    LastVer    string    `json:"last_ver"`
    Uuid       string    `json:"uuid"`
    CreateAt   time.Time `json:"create_at"`
    UpdateAt   time.Time `json:"update_at"`
    Servers    []Server
    Jumper     Server
    Project    Project
    //自定义shell需要时再增加字段
}

const (
    TaskPrePare  = 1 //初始状态
    TaskStarting = 2 //已开始发布

    TaskRunSuccess = 8 //发布成功
    TaskRunFail    = 9 // 发布失败

    DeployTypeIncrease = 1 //增量发布
    DeployTypeAll      = 2 //全量发布
)

//csv 响应前端结构
type CsvVersion struct {
    Version string
    Message string
}

//任务准备
type TaskPrepare struct {
    Task    DeployTask
    Env     EnvProServer
    Servers []Server
    Jumper  Server
    Project Project
}

//任务调度参数结构
type DeployTaskRunParams struct {
    Server      Server //当前发布目标机
    Jumper      Server
    Task        DeployTask
    Env         EnvProServer
    Project     Project
    ResChan     chan *DeployTaskResult
    PackagePath string //包全路径地址
    PackageName string //文件名 xxx.tar.gz
    PackageUuid string
    DstPath     string //目标机目标文件夹
    DstFilePath string
    DeployCmd   string //目标机脚本
}

//任务结果
type DeployTaskResult struct {
    Params    DeployTaskRunParams
    ResStatus int    //结果状态
    SwitchCmd string //切换软链接||php reload
    Output    string //错误输出
}

func DeleteEnvAndTask(projectId int) (err error) {
    var env []EnvProServer
    if err = mdb.Where("project_id = ?", projectId).Find(&env).Error; err != nil {
        return
    }
    for _, v := range env {
        mdb.Delete(&v)
        mdb.Delete(&DeployTask{}, "env_id = ?", v.EnvId)
    }
    return
}

func IsServerUsed(serverId int) (err error) {
    var total int
    if err = mdb.Model(&EnvProServer{}).Where("find_in_set(?, server_ids) or jump_server = ?", serverId, serverId).
        Count(&total).Error; err != nil {
        return
    }
    if total > 0 {
        err = errors.New("该服务已被占用, 不允许修改&删除")
    }
    return
}

func GetEnvCfgList(search *request.ComPageInfo) (envList []EnvProServer, total int, err error) {
    db := mdb
    if search.Condition != "" && search.SearchValue != "" {
        db = db.Where(search.Condition+" = ?", search.SearchValue)
    }
    if err = db.Model(&envList).Count(&total).Error; err != nil {
        return
    }
    if err = db.Limit(search.PageSize).Offset(search.PageSize * (search.CurrentPage - 1)).Find(&envList).Error; err != nil {
        return
    }
    for idx, li := range envList {
        mdb.Where("project_id = ?", li.ProjectId).First(&envList[idx].Project)
        mdb.Where("server_id in (?)", strings.Split(li.ServerIds, ",")).Find(&envList[idx].Servers)
        if li.JumpServer != 0 {
            mdb.Where("server_id = ?", li.JumpServer).First(&envList[idx].Jumper)
        }
    }
    return
}

func SaveEnvCfg(env *EnvProServer) (err error) {
    if env.EnvId == 0 {
        env.CreateAt = time.Now()
        env.UpdateAt = time.Now()
        err = mdb.Save(env).Error
    } else {
        err = mdb.Model(env).Updates(map[string]interface{}{
            "env_name":    env.EnvName,
            "project_id":  env.ProjectId,
            "server_ids":  env.ServerIds,
            "jump_server": env.JumpServer,
        }).Error
    }
    if err != nil {
        switch {
        case strings.Index(err.Error(), "uniq-env_name") != -1:
            err = errors.New("该配置名已存在")
        case strings.Index(err.Error(), "uniq-pro-sv-js") != -1:
            err = errors.New("该项目-目标机-跳板机已存在")
        }
    }
    return
}

func UpdateEnvRes(envId int, lastVer string, Uuid string) {
    mdb.Model(&EnvProServer{}).Where("env_id = ?", envId).Updates(map[string]interface{}{
        "last_ver": lastVer,
        "uuid":     Uuid,
    })
}

func DelEnvCfg(cfgId int) (err error) {
    return mdb.Delete(EnvProServer{}, "env_id = ?", cfgId).Error
}

func CheckDelTask(envId int) (err error) {
    var c int
    if err = mdb.Where("env_id = ?", envId).
        Where("status != ? and status != ? and status != ?", TaskPrePare, TaskRunSuccess, TaskRunFail).
        Model(&DeployTask{}).Count(&c).Error; err != nil {
        return
    }
    if c > 0 {
        err = errors.New("当前有发布中的任务, 请处理掉再删除")
    }
    return
}

func DelTaskByEnvId(envId int) (err error) {
    return mdb.Delete(DeployTask{}, "env_id = ?", envId).Error
}

func GetDeployTaskList(search *request.ComPageInfo) (taskList []DeployTask, total int, err error) {
    db := mdb
    if search.Condition != "" && search.SearchValue != "" {
        db = db.Where(search.Condition+" = ?", search.SearchValue)
    }
    if err = db.Model(&taskList).Count(&total).Error; err != nil {
        return
    }
    if err = db.Order("task_id desc").Limit(search.PageSize).Offset(search.PageSize * (search.CurrentPage - 1)).
        Find(&taskList).Error; err != nil {
        return
    }
    for idx, tl := range taskList {
        mdb.Where("env_id = ?", tl.EnvId).First(&taskList[idx].EnvCfg)
        if tl.Status == TaskStarting {
            mdb.Where("find_in_set(server_id, ?)", taskList[idx].EnvCfg.ServerIds).Find(&taskList[idx].EnvCfg.Servers)
        }
        if taskList[idx].EnvCfg.JumpServer != 0 {
            mdb.Where("server_id = ?", taskList[idx].EnvCfg.JumpServer).First(&taskList[idx].EnvCfg.Jumper)
        }
    }
    return
}

func SaveTask(task *DeployTask) (err error) {
    if task.TaskId == 0 {
        task.Status = TaskPrePare
        task.CreateAt = time.Now()
        task.UpdateAt = time.Now()
        err = mdb.Save(task).Error
    } else {
        err = errors.New("任务不支持修改, 如要修改, 请删除重建")
    }
    return
}

func UpdateTaskStatusAndOutput(taskId int, params map[string]interface{}) {
    mdb.Model(&DeployTask{TaskId: taskId}).Updates(params)
}

func DelTask(taskId int) (err error) {
    var task DeployTask
    if err = mdb.Where("task_id = ?", taskId).First(&task).Error; err != nil {
        return
    }
    if task.Status != TaskPrePare && task.Status != TaskRunFail {
        err = errors.New("该状态的任务不允许被删除")
    }
    err = mdb.Delete(&task).Error
    return
}

func PrepareTask(taskId int) (taskPrepare TaskPrepare, err error) {
    if err = mdb.Where("task_id = ?", taskId).First(&taskPrepare.Task).Error; err != nil {
        return
    }
    if err = mdb.Where("env_id = ?", taskPrepare.Task.EnvId).First(&taskPrepare.Env).Error; err != nil {
        return
    }
    if err = mdb.Where("project_id = ?", taskPrepare.Env.ProjectId).First(&taskPrepare.Project).Error; err != nil {
        return
    }
    if err = mdb.Where("find_in_set(server_id, ?)", taskPrepare.Env.ServerIds).Find(&taskPrepare.Servers).Error; err != nil {
        return
    }
    if taskPrepare.Env.JumpServer != 0 {
        err = mdb.Where("server_id = ?", taskPrepare.Env.JumpServer).First(&taskPrepare.Jumper).Error
    }
    return
}

func GetTaskById(taskId int) (task DeployTask) {
    mdb.Where("task_id = ?", taskId).First(&task)
    return
}
