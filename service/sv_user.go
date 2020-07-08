package service

import (
    "deploy/model"
    "deploy/model/request"
)

func Login(u *model.User) (*model.User, error) {
    return model.AuthLogin(u)
}

func UserList(search *request.ComPageInfo) ([]model.User, int, error) {
    return model.GetUserList(search)
}

func SaveUser(u *model.User) error {
    return model.SaveUser(u)
}

func UpdatePwd(u *request.ModifyPwd) (err error) {
    return model.UpdatePwd(u.UserId, u.NewPwd, u.OldPwd)
}
