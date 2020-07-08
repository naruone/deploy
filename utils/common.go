package utils

import (
    "crypto/md5"
    "encoding/hex"
    "os"
    "reflect"
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
