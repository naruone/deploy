package utils

import (
    "bytes"
    "crypto/md5"
    "deploy/config"
    "encoding/hex"
    "fmt"
    "math"
    "os"
    "os/exec"
    "reflect"
    "runtime"
    "time"
)

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
    s, err := os.Stat(path)
    if err != nil {
        return false
    }
    return s.IsDir()
}

func CreateDir(path string) (err error) {
    if !IsDir(path) {
        err = os.MkdirAll(path, os.ModePerm)
    }
    return
}

func DeletePath(path string) {
    _ = os.RemoveAll(path)
}

func MD5V(str []byte) string {
    h := md5.New()
    h.Write(str)
    return hex.EncodeToString(h.Sum(nil))
}

// 利用反射将结构体转化为map
func StructToMap(obj interface{}) map[string]interface{} {
    obj1 := reflect.TypeOf(obj)
    obj2 := reflect.ValueOf(obj)

    var data = make(map[string]interface{})
    for i := 0; i < obj1.NumField(); i++ {
        data[obj1.Field(i).Name] = obj2.Field(i).Interface()
    }
    return data
}

func RunCmd(dir, cmdName string, args ...string) (out string, errOut string, err error) {
    bufOut := &bytes.Buffer{}
    bufErr := &bytes.Buffer{}

    cmd := exec.Command(cmdName, args...)
    cmd.Dir = dir
    cmd.Stdout = bufOut
    cmd.Stderr = bufErr

    err = cmd.Run()
    out = string(bufOut.Bytes())
    errOut = string(bufErr.Bytes())
    return
}

func SystemInfo() map[string]interface{} {
    var (
        afterLastGC string
        mstat       runtime.MemStats
    )
    runtime.ReadMemStats(&mstat)

    costTime := int(time.Since(config.GConfig.StartTime).Seconds())

    if mstat.LastGC != 0 {
        afterLastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(mstat.LastGC))/1000/1000/1000)
    } else {
        afterLastGC = "0"
    }

    return map[string]interface{}{
        "system":          runtime.GOOS + " " + runtime.GOARCH,
        "run_time":        fmt.Sprintf("%d天%d小时%d分%d秒", costTime/(3600*24), costTime%(3600*24)/3600, costTime%3600/60, costTime%(60)),
        "goroutine_num":   runtime.NumGoroutine(),
        "go_version":      runtime.Version(),
        "cpu":             runtime.NumCPU(),
        "stack_mem_used":  FileSizeFormat(int64(mstat.Alloc)),      //堆上对象占用的内存大小
        "stack_alloc_mem": FileSizeFormat(int64(mstat.TotalAlloc)), //堆上总共分配出的内存大小
        "sys_mem_used":    FileSizeFormat(int64(mstat.Sys)),        //程序从操作系统总共申请的内存大小
        "pointer_find":    mstat.Lookups,
        "last_gc_time":    afterLastGC,
        "next_gc_mem":     FileSizeFormat(int64(mstat.NextGC)),
    }
}

//计算文件大小B转换成 KB, GB等
func FileSizeFormat(s int64) string {
    sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
    if s < 10 {
        return fmt.Sprintf("%d B", s)
    }

    e := math.Floor(math.Log(float64(s)) / math.Log(1024))
    suffix := sizes[int(e)]
    val := float64(s) / math.Pow(1024, math.Floor(e))
    f := "%.0f"
    if val < 10 {
        f = "%.1f"
    }
    return fmt.Sprintf(f+" %s", val, suffix)
}
