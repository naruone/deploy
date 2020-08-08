package config

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "sync"
)

var (
    GConfig *SysConfig
    cfgLock sync.Mutex
)

type SysConfig struct {
    Env               string    `json:"env"`
    Port              int       `json:"port"`
    ViewPath          string    `json:"view_path"`
    StaticDir         string    `json:"static_dir"`
    DbConfig          *DbConfig `json:"db"`
    JwtSigningKey     string    `json:"jwt-signing-key"`
    Captcha           *Captcha  `json:"captcha"`
    SshConnectTimeout int       `json:"ssh_connect_timeout"`
    Repo              string    `json:"repo"`
    RepoPackage       string    `json:"repo_package"`
    ServerWorkDir     string    `json:"srv_work_dir"`
}

type DbConfig struct {
    Host            string `json:"host"`
    Port            int    `json:"port"`
    Database        string `json:"database"`
    Username        string `json:"username"`
    Password        string `json:"password"`
    Charset         string `json:"charset"`
    TimeOut         int    `json:"timeout"`
    MaxOpenConnects int    `json:"maxOpenConnects"` //数据库连接池最大连接数
    MaxIdleConnects int    `json:"maxIdleConnects"` //连接池最大允许的空闲连接数, 超过的会被关闭
}

type Captcha struct {
    Long   int `json:"long"`
    Width  int `json:"width"`
    Height int `json:"height"`
}

func InitConfig() {
    cfgLock.Lock()
    defer cfgLock.Unlock()
    if GConfig != nil {
        return
    }
    var (
        bytes     []byte
        sysConfig SysConfig
        err       error
    )
    if bytes, err = ioutil.ReadFile("./config/config.json"); err != nil {
        log.Fatal(err)
    }
    if err = json.Unmarshal(bytes, &sysConfig); err != nil {
        log.Fatal(err)
    }
    GConfig = &sysConfig
}
