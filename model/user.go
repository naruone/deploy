package model

import (
    "deploy/model/request"
    "deploy/utils"
    "errors"
    "strings"
    "time"
)

type User struct {
    UserId   uint      `json:"user_id" gorm:"PRIMARY_KEY"`
    UserName string    `json:"user_name"`
    Password string    `json:"password"`
    NickName string    `json:"nick_name"`
    CreateAt time.Time `json:"create_at"`
    Status   int       `json:"status"`
}

const (
    UserStatusNormal   = 1 //正常
    UserStatusDisabled = 2 //禁用
)

func AuthLogin(u *User) (userInter *User, err error) {
    var user User
    u.Password = utils.MD5V([]byte(u.Password))
    err = mdb.Where("user_name = ? AND password = ?", u.UserName, u.Password).First(&user).Error
    if err != nil {
        err = errors.New("用户名或者密码错误")
        return
    }
    if user.Status != UserStatusNormal {
        err = errors.New("该用户不允许登录")
        return
    }
    userInter = &user
    return
}

func GetUserList(search *request.ComPageInfo) (userList []User, total int, err error) {
    db := mdb
    if search.Condition != "" && search.SearchValue != "" {
        db = db.Where(search.Condition+" = ?", search.SearchValue)
    }
    if search.Status != 0 {
        db = db.Where("status = ?", search.Status)
    }
    if err = db.Model(&userList).Count(&total).Error; err != nil {
        return
    }
    err = db.Limit(search.PageSize).Offset(search.PageSize * (search.CurrentPage - 1)).Select([]string{
        "user_id", "user_name", "nick_name", "create_at", "status",
    }).Find(&userList).Error
    return
}

func SaveUser(u *User) (err error) {
    db := mdb
    if u.UserId != 0 { //更新
        var user User
        if err = db.Where("user_id = ?", u.UserId).First(&user).Error; err != nil {
            return
        }
        if u.Status == UserStatusNormal || u.Status == UserStatusDisabled {
            user.Status = u.Status
        }
        if u.UserName != "" {
            user.UserName = u.UserName
        }
        if u.Password != "" {
            user.Password = utils.MD5V([]byte(u.Password))
        }
        if u.NickName != "" {
            user.NickName = u.NickName
        }
        err = db.Save(&user).Error
    } else {
        if u.Password == "" {
            err = errors.New("新增用户密码必填")
            return
        }
        err = db.Save(&User{ //创建用户
            UserName: u.UserName,
            Password: utils.MD5V([]byte(u.Password)),
            NickName: u.NickName,
            CreateAt: time.Now(),
            Status:   UserStatusNormal,
        }).Error
    }
    if err != nil && strings.Index(err.Error(), "Duplicate entry") != -1 {
        err = errors.New("该用户名已存在")
    }
    return
}

func UpdatePwd(userId int, password string, oPwd string) (err error) {
    var user User
    if err = mdb.Where("user_id = ?", userId).Find(&user).Error; err != nil {
        return
    }
    if utils.MD5V([]byte(oPwd)) != user.Password {
        err = errors.New("旧密码错误")
        return
    }
    return mdb.Model(&user).Update("password", utils.MD5V([]byte(password))).Error
}
