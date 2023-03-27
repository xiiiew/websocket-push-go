package server

import (
	"errors"
	"github.com/xiiiew/websocket-push-go/common"
	"sync"
)

// 频道桶
type bucket struct {
	mu      sync.Mutex
	connMap map[uint64]*WsConnection // 所有连接(可不订阅频道)
	chMap   map[string]*Channel      // 所有频道连接
}

var (
	once           sync.Once
	bucketInstance *bucket
)

// 获取桶实例
func GetBucketInstance() *bucket {
	once.Do(initBucket)
	return bucketInstance
}

// 初始化桶
func initBucket() {
	bucketInstance = &bucket{
		connMap: make(map[uint64]*WsConnection),
		chMap:   make(map[string]*Channel),
	}
}

// 添加连接
func (b *bucket) AddConn(conn *WsConnection) {
	b.mu.Lock()
	b.mu.Unlock()

	b.connMap[conn.connId] = conn
}

// 删除连接
func (b *bucket) DelConn(conn *WsConnection) {
	// 删除所有频道中连接
	b.UnsubscribeAllCh(conn)

	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.connMap, conn.connId)
}

// 订阅
func (b *bucket) SubscribeCh(chName string, conn *WsConnection) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 频道不存在, 创建频道
	if ch, isExist := b.chMap[chName]; !isExist {
		ch = InitCh(chName)
		b.chMap[chName] = ch
	}

	b.chMap[chName].Subscribe(conn)
}

// 取消订阅
func (b *bucket) UnsubscribeCh(chName string, conn *WsConnection) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.chMap[chName].Unsubscribe(conn)

	// 频道为空,则删除
	if b.chMap[chName].Count() == 0 {
		delete(b.chMap, chName)
	}
}

// 取消订阅所有频道
func (b *bucket) UnsubscribeAllCh(conn *WsConnection) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ch := range b.chMap {
		ch.Unsubscribe(conn)
	}
}

// 推送给桶内所有用户
func (b *bucket) PushAll(wsMsg *common.WsMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, conn := range b.connMap {
		go conn.SendMessage(wsMsg)
	}
}

// 推送给指定用户
func (b *bucket) PushOne(connId uint64, wsMsg *common.WsMessage) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	conn, ok := b.connMap[connId]
	if !ok {
		return errors.New("connection not found")
	}

	return conn.SendMessage(wsMsg)
}

// 推送到频道内所有用户
func (b *bucket) PushCh(chName string, wsMsg *common.WsMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if ch, isExist := b.chMap[chName]; isExist {
		ch.PushAll(wsMsg)
	}
}

// 频道数
func (b *bucket) ChCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return len(b.chMap)
}
