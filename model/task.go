package model

import (
    "errors"
    "time"
)

type DeployTask struct {
    TaskId      int       `json:"task_id" gorm:"PRIMARY_KEY"`
    TaskName    string    `json:"task_name"`
    DeployType  int       `json:"deploy_type"`
    Description string    `json:"description"`
    ErrOutput   string    `json:"err_output"`
    EnvId       int       `json:"env_id"`
    Version     string    `json:"version"`
    AfterScript string    `json:"after_script"`
    Status      int       `json:"status"`
    Uuid        string    `json:"uuid"`
    CreateAt    time.Time `json:"create_at"`
    UpdateAt    time.Time `json:"update_at"`
}

type EnvProServer struct {
    EnvId      int       `json:"env_id" gorm:"PRIMARY_KEY"`
    EnvName    string    `json:"env_name"`
    ProjectId  int       `json:"project_id"`
    ServerId   int       `json:"server_id"`
    JumpServer int       `json:"jump_server"`
    LastVer    string    `json:"last_ver"`
    Uuid       string    `json:"uuid"`
    CreateAt   time.Time `json:"create_at"`
    UpdateAt   time.Time `json:"update_at"`
}

const (
    TaskPrePare  = 1 //初始状态
    TaskStarting = 2 //已开始发布

    TaskRunSuccess = 8 //发布成功
    TaskRunFail    = 9 // 发布失败

    DeployTypeIncrease = 1 //增量发布
    DeployTypeAll      = 2 //全量发布
)

type CsvVersion struct {
    Version string
    Message string
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
    if err = mdb.Model(&EnvProServer{}).Where("server_id = ? or jump_server = ?", serverId, serverId).
        Count(&total).Error; err != nil {
        return
    }
    if total > 0 {
        err = errors.New("该服务已被占用, 请先删除对应配置")
    }
    return
}
