package common

import "errors"

// 操作类型
const (
	SubAction = "sub" // 订阅

	UnsubAction = "unsub" // 取消订阅
)

// ws连接错误
var (
	CONNECTION_IS_CLOSED = errors.New("CONNECTION_IS_CLOSED")
)
