package model

import (
    "deploy/model/request"
    "errors"
    "strings"
    "time"
)

type Server struct {
    ServerId   int       `json:"server_id" gorm:"PRIMARY_KEY"`
    Type       int       `json:"type"`
    SshAddr    string    `json:"ssh_addr"`
    SshPort    int       `json:"ssh_port"`
    SshUser    string    `json:"ssh_user"`
    SshKeyPath string    `json:"ssh_key_path"`
    CreateAt   time.Time `json:"create_at"`
    UpdateAt   time.Time `json:"update_at"`
}

const (
    ServerTypeServer = 1 //服务器
    ServerTypeJumper = 2 //跳板机
)

func GetServerList(search *request.ComPageInfo) (serverList []Server, total int, err error) {
    db := mdb
    if search.Condition != "" && search.SearchValue != "" {
        db = db.Where(search.Condition+" = ?", search.SearchValue)
    }
    if search.Type != 0 {
        db = db.Where("type = ?", search.Type)
    }
    if err = db.Model(&serverList).Count(&total).Error; err != nil {
        return
    }
    err = db.Limit(search.PageSize).Offset(search.PageSize * (search.CurrentPage - 1)).Find(&serverList).Error
    return
}

func DelServer(serverId int) (err error) {
    return mdb.Delete(Server{}, "server_id = ?", serverId).Error
}

func SaveServer(p *Server) (err error) {
    if p.Type != ServerTypeServer && p.Type != ServerTypeJumper {
        err = errors.New("服务器类型错误, 不支持的类型")
        return
    }
    if p.ServerId == 0 {
        p.CreateAt = time.Now()
        p.UpdateAt = time.Now()
        err = mdb.Save(p).Error
    } else {
        _updateParams := map[string]interface{}{
            "type":         p.Type,
            "ssh_addr":     p.SshAddr,
            "ssh_port":     p.SshPort,
            "ssh_user":     p.SshUser,
            "ssh_key_path": p.SshKeyPath,
            "update_at":    time.Now(),
        }
        err = mdb.Model(p).Updates(_updateParams).Error
    }
    if err != nil {
        if strings.Index(err.Error(), "uniq-type_ip") != -1 {
            err = errors.New("该: 类型-IP 的服务器已存在")
        }
    }
    return
}

func GetServerById(serverId int) (server Server) {
    mdb.Where("server_id = ?", serverId).First(&server)
    return
}
