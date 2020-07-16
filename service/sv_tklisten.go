package service

import (
    "deploy/router/middleware"
    "encoding/json"
    "fmt"
    "github.com/gorilla/websocket"
    "net/http"
    "time"
)

type TaskListen struct {
    TaskId int //任务ID
    Client *websocket.Conn
}

//websocket 请求参数
type WsReqData struct {
    Type   string
    Data   string
    TaskId int
}

type WsRespData struct {
    Status  int //0 成功, 999失败
    Message string
    Data    interface{}
}

type TaskProcessReport struct {
    TaskId  int
    Process int
}

var (
    taskListens []*TaskListen
    wsUpGrader  = websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }
    ProcessListenChan = make(chan *TaskProcessReport)
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    var (
        conn      *websocket.Conn
        msgType   int
        recData   []byte
        wsReqData WsReqData

        err error
    )
    if conn, err = wsUpGrader.Upgrade(w, r, nil); err != nil {
        http.NotFound(w, r)
        return
    }
    defer func() {
        _ = conn.Close()
        removeClient(conn)
    }()
    for {
        if msgType, recData, err = conn.ReadMessage(); err != nil {
            removeClient(conn)
            break
        }
        if err = json.Unmarshal(recData, &wsReqData); err != nil || wsReqData.TaskId == 0 {
            removeClient(conn)
            break
        }
        if wsReqData.Type == "authentication" {
            myJwt := middleware.NewJWT()
            if _, err := myJwt.ParseToken(wsReqData.Data); err != nil {
                removeClient(conn)
                _ = conn.WriteMessage(msgType, getWsRespData(999, "token error", nil))
                _ = conn.Close()
                return
            }
            _ = conn.WriteMessage(msgType, getWsRespData(0, "auth success", nil))
            taskListens = append(taskListens, &TaskListen{
                TaskId: wsReqData.TaskId,
                Client: conn,
            })
        }
        if wsReqData.Type == "keep-alive" {
            _ = conn.WriteMessage(msgType, getWsRespData(0, time.Now().String(), nil))
        }
        fmt.Println("------------", len(taskListens))
    }
}

//删除关闭的连接
func removeClient(conn *websocket.Conn) {
    for idx, v := range taskListens {
        if v.Client == conn {
            _ = conn.Close()
            taskListens = append(taskListens[:idx], taskListens[idx+1:]...)
        }
    }
}

func getWsRespData(code int, msg string, data interface{}) (res []byte) {
    res, _ = json.Marshal(WsRespData{
        Status:  code,
        Message: msg,
        Data:    data,
    })
    return
}

func sendMsg(report *TaskProcessReport) {
    for _, v := range taskListens {
        if v.TaskId == report.TaskId {
            _ = v.Client.WriteMessage(1, getWsRespData(0, "获取", report))
        }
    }
}

func SchedulerDeployInfo() {
    for {
        select {
        case taskProcess := <-ProcessListenChan:
            sendMsg(taskProcess)
        }
    }
}
