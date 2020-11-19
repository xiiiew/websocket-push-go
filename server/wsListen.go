package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/xiiiew/websocket-push-go/common"
	"time"
)

// 处理websocket消息
func (conn *WsConnection) WsListen() {
	// 将连接添加到桶
	GetBucketInstance().AddConn(conn)

	// 检查心跳
	go conn.checkPing()

	// 每20秒ping一次
	go conn.loopPing()

	// 处理消息
	for {
		message, err := conn.ReadMessage()
		if err != nil {
			conn.DelConn()
			return
		}

		if message.MsgType != websocket.TextMessage {
			continue
		}

		// 判断消息类型
		if common.IsPing(message.MsgData) {
			conn.pong()
			continue
		}
		if common.IsPong(message.MsgData) {
			conn.KeepAlive()
			continue
		}
		if isTrue, actionData := common.IsAction(message.MsgData); !isTrue {
			conn.KeepAlive()
			continue
		} else {
			switch actionData.Action {
			case common.SubAction: // 订阅
				conn.subscribe(actionData.Ch)
			case common.UnsubAction: // 取消订阅
				conn.unsubscribe(actionData.Ch)
			}
		}
	}
}

// 订阅
func (conn *WsConnection) subscribe(ch string) {
	conn.chs[ch] = true
	GetBucketInstance().SubscribeCh(ch, conn)
}

// 取消订阅
func (conn *WsConnection) unsubscribe(ch string) {
	delete(conn.chs, ch)
	GetBucketInstance().UnsubscribeCh(ch, conn)
}

// 取消订阅所有频道
//func (conn *WsConnection) unsubscribeAll() {
//	GetBucketInstance().UnsubscribeAllCh(conn)
//}

// 每秒检查心跳
func (conn *WsConnection) checkPing() {
	if !conn.IsAlive() {
		conn.Close()
		return
	}

	time.AfterFunc(time.Second, conn.checkPing)
}

// 每20秒Ping
func (conn *WsConnection) loopPing() {
	conn.ping()

	time.AfterFunc(conn.pingDuration*time.Second, conn.loopPing)
}

// ping
func (conn *WsConnection) ping() {
	message := new(common.WsMessage)
	msgData := common.WsPingData{Ping: int(time.Now().Unix())}
	message.MsgType = websocket.TextMessage
	message.MsgData, _ = json.Marshal(msgData)

	if err := conn.SendMessage(message); err != nil {
		conn.DelConn()
	}
}

// pong
func (conn *WsConnection) pong() {
	message := new(common.WsMessage)
	msgData := common.WsPongData{Pong: int(time.Now().Unix())}
	message.MsgType = websocket.TextMessage
	message.MsgData, _ = json.Marshal(msgData)

	if err := conn.SendMessage(message); err != nil {
		conn.DelConn()
		return
	}

	// 更新心跳
	conn.KeepAlive()
}

// 删除所有连接状态
func (conn *WsConnection) DelConn() {
	// 关闭ws连接
	conn.Close()
	// 从桶中删除连接
	GetBucketInstance().DelConn(conn)
}
