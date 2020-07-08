package model

import (
    "deploy/config"
    "deploy/model/request"
    "deploy/utils"
    "errors"
    "strings"
    "time"
)

type Project struct {
    ProjectId   uint      `json:"project_id" gorm:"PRIMARY_KEY"`
    ProjectName string    `json:"project_name"`
    RepoUrl     string    `json:"repo_url"`
    Dst         string    `json:"dst"`
    WebRoot     string    `json:"web_root"`
    AfterScript string    `json:"after_script"`
    ErrMsg      string    `json:"err_msg"`
    CreateAt    time.Time `json:"create_at"`
    UpdateAt    time.Time `json:"update_at"`
    Status      int       `json:"status"`
}

const (
    ProjectNotInit        = 1
    ProjectInitSuccess    = 2
    ProjectInitFail       = 3
    ProjectInitProcessing = 9
)

func GetProjectList(search *request.ComPageInfo) (projectList []Project, total int, err error) {
    db := mdb
    if search.Condition != "" && search.SearchValue != "" {
        db = db.Where(search.Condition+" = ?", search.SearchValue)
    }
    if search.Status != 0 {
        db = db.Where("status = ?", search.Status)
    }
    if err = db.Model(&projectList).Count(&total).Error; err != nil {
        return
    }
    err = db.Limit(search.PageSize).Offset(search.PageSize * (search.CurrentPage - 1)).Find(&projectList).Error
    return
}

func DelProject(projectId int) (err error) {
    var project Project
    if err = mdb.Where("project_id = ?", projectId).First(&project).Error; err != nil {
        return
    }
    if project.Status == ProjectInitProcessing {
        err = errors.New("该项目正在初始化, 请等待初始化完成后删除")
    }
    if err = mdb.Delete(&project).Error; err != nil {
        return
    }
    utils.DeletePath(strings.TrimRight(config.GConfig.Repo, "/") + "/" + project.Dst)
    return
}

func SaveProject(p *Project) (err error) {
    if p.ProjectId == 0 {
        p.Status = ProjectNotInit
        p.CreateAt = time.Now()
        p.UpdateAt = time.Now()
        err = mdb.Save(p).Error
    } else {
        err = mdb.Model(p).Updates(map[string]interface{}{
            "project_name": p.ProjectName,
            "repo_url":     p.RepoUrl,
            "dst":          p.Dst,
            "web_root":     p.WebRoot,
            "after_script": p.AfterScript,
        }).Error
    }
    if err != nil {
        switch {
        case strings.Index(err.Error(), "uni-project-name") != -1:
            err = errors.New("该用项目名已存在")
        case strings.Index(err.Error(), "uniq-repo") != -1:
            err = errors.New("该仓库地址已存在")
        case strings.Index(err.Error(), "uniq-dst") != -1:
            err = errors.New("初始DIR已存在")
        }
    }
    return
}

func GetProjectById(projectId int, isInit bool) (project Project, err error) {
    err = mdb.Where("project_id = ?", projectId).First(&project).Error
    if isInit && project.Status == ProjectInitSuccess {
        err = errors.New("该项目已经初始化成功")
        return
    }
    if !isInit && project.Status == ProjectNotInit {
        err = errors.New("该项目还未初始化")
        return
    }
    if project.Status == ProjectInitProcessing {
        err = errors.New("该项目正在初始化")
        return
    }

    if !isInit && project.Status == ProjectInitFail {
        err = errors.New("该项目初始化失败, 请先处理, 原因: " + project.ErrMsg)
        return
    }
    return
}
