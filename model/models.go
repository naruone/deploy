package model

import (
    "deploy/config"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "log"
    "os"
    "sync"
)

var (
    mdb    *gorm.DB
    dbLock sync.Mutex
)

func InitModels() {
    dbLock.Lock()
    defer dbLock.Unlock()
    if mdb != nil {
        return
    }
    var (
        _db      *gorm.DB
        dbDriver string
        err      error
    )
    //admin.Username+":"+admin.Password+"@("+admin.Path+")/"+admin.Dbname+"?"+admin.Config
    dbDriver = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%ds",
        config.GConfig.DbConfig.Username, config.GConfig.DbConfig.Password, config.GConfig.DbConfig.Host,
        config.GConfig.DbConfig.Port, config.GConfig.DbConfig.Database, config.GConfig.DbConfig.TimeOut)
    if _db, err = gorm.Open("mysql", dbDriver); err != nil {
        log.Fatal("数据库连接失败", err)
        os.Exit(0)
    }
    _db.DB().SetMaxOpenConns(config.GConfig.DbConfig.MaxOpenConnects)
    _db.DB().SetMaxIdleConns(config.GConfig.DbConfig.MaxIdleConnects)
    if config.GConfig.Env == "debug" {
        _db.LogMode(true)
    }
    mdb = _db
}
