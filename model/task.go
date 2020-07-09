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
			"update_at":   time.Now(),
		}).Error
	}
	if err != nil {
		switch {
		case strings.Index(err.Error(), "uniq-env_name") != -1:
			err = errors.New("该用配置名已存在")
		case strings.Index(err.Error(), "uniq-pro-sv-js") != -1:
			err = errors.New("该项目-目标机-跳板机已存在")
		}
	}
	return
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
