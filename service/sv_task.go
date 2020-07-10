package service

import (
    "deploy/model"
    "deploy/model/request"
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
