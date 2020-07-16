let webSocket = null

function initWebSocket() {
    webSocket = new WebSocket('ws://127.0.0.1')
    webSocket.onmessage = function (e) {
        websocketonmessage(e)
    }
    webSocket.onclose = function (e) {
        websocketclose(e)
    }
    webSocket.onopen = function () {
        websocketOpen()
    }

    // 连接发生错误的回调方法
    webSocket.onerror = function () {
        console.log('WebSocket连接发生错误')
    }
}

// 实际调用的方法
function sendSock(agentData, callback) {
    globalCallback = callback
    if (websock.readyState === websock.OPEN) {
        // 若是ws开启状态
        websocketsend(agentData)
    } else if (websock.readyState === websock.CONNECTING) {
        // 若是 正在开启状态，则等待1s后重新调用
        setTimeout(function () {
            sendSock(agentData, callback)
        }, 1000)
    } else {
        // 若未开启 ，则等待1s后重新调用
        setTimeout(function () {
            sendSock(agentData, callback)
        }, 1000)
    }
}

// 数据接收
function websocketonmessage(e) {
    globalCallback(JSON.parse(e.data))
}

// 数据发送
function websocketsend(agentData) {
    websock.send(JSON.stringify(agentData))
}

// 关闭
function websocketclose(e) {
    console.log('connection closed (' + e.code + ')')
}

// 创建 websocket 连接
function websocketOpen(e) {
    console.log('连接成功')
}

initWebSocket()

// 将方法暴露出去
export {
    sendSock
}
