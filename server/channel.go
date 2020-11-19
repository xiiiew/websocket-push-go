package server

import (
	"sync"
	"github.com/xiiiew/websocket-push-go/common"
)

// 频道
type Channel struct {
	mu      sync.Mutex
	chName  string
	connMap map[uint64]*WsConnection
}

// 初始化频道
func InitCh(ch string) *Channel {
	return &Channel{
		chName:  ch,
		connMap: make(map[uint64]*WsConnection),
	}
}

// 订阅频道
func (c *Channel) Subscribe(conn *WsConnection) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connMap[conn.connId] = conn
}

// 取消订阅频道
func (c *Channel) Unsubscribe(conn *WsConnection) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.connMap, conn.connId)
}

// 频道中用户数
func (c *Channel) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.connMap)
}

// 发送消息到所有用户
func (c *Channel) PushAll(wsMsg *common.WsMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, conn := range c.connMap {
		go conn.SendMessage(wsMsg)
	}
}

// 发送消息到指定用户
//func (c *Channel) PushOne(connId uint64, wsMsg *common.WsMessage) error {
//	c.mu.Lock()
//	defer c.mu.Unlock()
//
//	return c.connMap[connId].SendMessage(wsMsg)
//}
