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
    DelRepository(uint(projectId))
    return
}
