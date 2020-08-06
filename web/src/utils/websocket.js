import {store} from "../store/index";

let ws = null
let timer = null

export const InitWebSocket = (taskId, updateStatus, callback) => {
    ws = new WebSocket("ws://" + process.env.VUE_APP_DEPLOY_WS_URI + "/ws")
    ws.onmessage = function (e) {
        let _data = JSON.parse(e.data)
        if (_data.Status !== 0) {
            console.log("ws connect error: ", e.data);
            return
        }
        switch (_data.Message) {
            case 'auth-success':
                timer = setInterval(() => {
                    SendMsg({type: "keep-alive"})
                }, 5000);
                break;
            case 'auto-push':
                if (typeof callback === "function") {
                    callback({
                        taskId: taskId,
                        data: _data.Data
                    })
                }
                break;
        }
    }
    ws.onclose = function (e) {
        console.log('close ws from server');
        if (typeof updateStatus === "function") {
            updateStatus()
        }
        clearInterval(timer)
    }
    ws.onopen = function () {
        const token = store.getters['user/token']
        ws.send(JSON.stringify({
            type: "authentication",
            taskId: taskId,
            data: token
        }))
    }
    ws.onerror = function () {
        console.log('ws connect error')
    }
}

export const SendMsg = (data) => {
    if (ws.readyState === ws.OPEN) {
        ws.send(JSON.stringify(data))
    } else {
        console.log("send msg fail")
    }
}
