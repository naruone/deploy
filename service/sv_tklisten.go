package service

import (
    "deploy/model"
    "deploy/router/middleware"
    "encoding/json"
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
    Task    model.DeployTask
    Server  model.Server
    Process string
    Result  string
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
    deployResults     = make(map[int]map[string]interface{})

    //任务阶段
    TaskProcessPack           = "pack"
    TaskProcessUploadToJumper = "upload_jumper"
    TaskProcessUploadDst      = "upload_dst"
    TaskProcessDeploy         = "deploy"
    TaskProcessChangeWorkDir  = "change_dir"
    //end 任务阶段

    TaskSuccess = "success"
    TaskFail    = "fail"
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
                _ = conn.WriteMessage(msgType, getWsRespData(999, "token-error", nil))
                _ = conn.Close()
                return
            }
            if (model.GetTaskById(wsReqData.TaskId)).Status != model.TaskStarting {
                removeClient(conn)
                _ = conn.WriteMessage(msgType, getWsRespData(999, "task not processing", nil))
                _ = conn.Close()
                return
            }
            _ = conn.WriteMessage(msgType, getWsRespData(0, "auth-success", nil))
            taskListens = append(taskListens, &TaskListen{
                TaskId: wsReqData.TaskId,
                Client: conn,
            })
            if deployResults[wsReqData.TaskId] != nil {
                //连接成功时如果有消息则发送一次
                sendMsg(wsReqData.TaskId, deployResults[wsReqData.TaskId]) //连接成功时发一次
            }
        }
        if wsReqData.Type == "keep-alive" {
            _ = conn.WriteMessage(msgType, getWsRespData(0, time.Now().String(), nil))
        }
        if wsReqData.Type == "get-process" { //目前没用到, 连接时会主动推送(暂时没有场景)
            sendMsg(wsReqData.TaskId, deployResults[wsReqData.TaskId])
        }
    }
}

//删除关闭的连接
func removeClient(conn *websocket.Conn) {
    var _taskListens []*TaskListen
    for _, v := range taskListens {
        if v.Client == conn {
            _ = conn.Close()
        } else {
            _taskListens = append(_taskListens, v)
        }
    }
    taskListens = _taskListens
}

func getWsRespData(code int, msg string, data interface{}) (res []byte) {
    res, _ = json.Marshal(WsRespData{
        Status:  code,
        Message: msg,
        Data:    data,
    })
    return
}

func sendMsg(taskId int, report interface{}) {
    for _, v := range taskListens {
        if v.TaskId == taskId {
            _ = v.Client.WriteMessage(1, getWsRespData(0, "auto-push", report))
        }
    }
}

func CloseWsConnectByTaskId(taskId int) {
    var _taskListens []*TaskListen
    for _, v := range taskListens {
        if v.TaskId == taskId {
            _ = v.Client.WriteMessage(1, getWsRespData(0, "auto-close", nil))
            _ = v.Client.Close()
        } else {
            _taskListens = append(_taskListens, v)
        }
    }
    taskListens = _taskListens
    if deployResults[taskId] != nil {
        delete(deployResults, taskId)
    }
}

func SchedulerDeployInfo() {
    var serverResult map[string]map[string]string
    for {
        select {
        case taskProcess := <-ProcessListenChan:
            if deployResults[taskProcess.Task.TaskId] == nil {
                deployResults[taskProcess.Task.TaskId] = map[string]interface{}{}
            }
            if deployResults[taskProcess.Task.TaskId]["servers"] == nil {
                deployResults[taskProcess.Task.TaskId]["servers"] = map[string]map[string]string{}
            }
            if taskProcess.Process == TaskProcessPack || taskProcess.Process == TaskProcessUploadToJumper {
                deployResults[taskProcess.Task.TaskId][taskProcess.Process] = taskProcess.Result
            } else {
                serverResult = deployResults[taskProcess.Task.TaskId]["servers"].(map[string]map[string]string)
                if serverResult[taskProcess.Server.SshAddr] == nil {
                    serverResult[taskProcess.Server.SshAddr] = map[string]string{}
                }
                serverResult[taskProcess.Server.SshAddr][taskProcess.Process] = taskProcess.Result
                deployResults[taskProcess.Task.TaskId]["servers"] = serverResult
            }
            sendMsg(taskProcess.Task.TaskId, deployResults[taskProcess.Task.TaskId])
        }
    }
}
