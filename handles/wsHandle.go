package handles

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"github.com/xiiiew/websocket-push-go/server"
)

var _conn_id_ uint64
var once sync.Once

func initConnId() {
	_conn_id_ = uint64(time.Now().Unix())
}

func getConnId() uint64 {
	once.Do(initConnId)
	atomic.AddUint64(&_conn_id_, 1)
	return _conn_id_
}

func WsHandle(w http.ResponseWriter, r *http.Request, q httprouter.Params) {
	wsUpgrader := websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 生成连接id
	connId := getConnId()
	// 初始化ws连接
	wsConn := server.InitWsConnection(connId, conn)
	// 开始处理ws消息
	wsConn.WsListen()
}
