package server

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"github.com/xiiiew/websocket-push-go/common"
)

type WsConnection struct {
	mutex             sync.Mutex
	connId            uint64
	conn              *websocket.Conn
	inChan            chan *common.WsMessage
	outChan           chan *common.WsMessage
	closeChan         chan byte
	isClosed          bool
	lastHeartbeatTime time.Time       // 最近一次心跳时间
	chs               map[string]bool // 订阅了哪些频道
	pingDuration      time.Duration
	closeDuration     time.Duration
}

// 初始化ws连接
func InitWsConnection(connId uint64, conn *websocket.Conn) (wsConnection *WsConnection) {
	wsConnection = &WsConnection{
		conn:              conn,
		connId:            connId,
		inChan:            make(chan *common.WsMessage, common.Config.InChanLength),
		outChan:           make(chan *common.WsMessage, common.Config.OutChanLength),
		closeChan:         make(chan byte),
		lastHeartbeatTime: time.Now(),
		chs:               make(map[string]bool),
		pingDuration:      time.Duration(common.Config.PingDuration),
		closeDuration:     time.Duration(common.Config.CloseDuration),
	}

	go wsConnection.readLoop()
	go wsConnection.writeLoop()

	return
}

// 发送消息
func (conn *WsConnection) SendMessage(message *common.WsMessage) (err error) {
	select {
	case conn.outChan <- message:
	case <-conn.closeChan:
		err = common.CONNECTION_IS_CLOSED
	}
	return
}

// 读取消息
func (conn *WsConnection) ReadMessage() (message *common.WsMessage, err error) {
	select {
	case message = <-conn.inChan:
	case <-conn.closeChan:
		err = common.CONNECTION_IS_CLOSED
	}
	return
}

// 关闭连接
func (conn *WsConnection) Close() {
	conn.conn.Close()

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	if !conn.isClosed {
		conn.isClosed = true
		close(conn.closeChan)
	}
}

// 检查心跳
func (conn *WsConnection) IsAlive() bool {
	now := time.Now()

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	// 连接已关闭 或者 太久没有心跳
	if conn.isClosed || now.Sub(conn.lastHeartbeatTime) > conn.closeDuration*time.Second {
		return false
	}
	return true
}

// 更新心跳
func (conn *WsConnection) KeepAlive() {
	now := time.Now()

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	conn.lastHeartbeatTime = now
}

// 读websocket
func (conn *WsConnection) readLoop() {
	for {
		if msgType, msgData, err := conn.conn.ReadMessage(); err != nil {
			goto ERR
		} else {

			message := common.WsMessageBuilder(msgType, msgData)

			select {
			case conn.inChan <- message:
			case <-conn.closeChan:
				goto CLOSED
			}
		}
	}

ERR:
	conn.Close()
CLOSED:
}

// 写websocket
func (conn *WsConnection) writeLoop() {
	for {
		select {
		case message := <-conn.outChan:
			if err := conn.conn.WriteMessage(message.MsgType, message.MsgData); err != nil {
				goto ERR
			}
		case <-conn.closeChan:
			goto CLOSED
		}
	}
ERR:
	conn.Close()
CLOSED:
}

