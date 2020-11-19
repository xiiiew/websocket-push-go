package common

import (
	"encoding/json"
	"strings"
)

// 推送的Message对象
type WsMessage struct {
	MsgType int
	MsgData []byte
}

// 订阅/取消订阅 消息结构
type WsAction struct {
	Action string
	Ch     string
}

// ping 消息结构 {"ping": 1600000000}
type WsPingData struct {
	Ping int `json:"ping"`
}

// pong 消息结构 {"ping": 1600000000}
type WsPongData struct {
	Pong int `json:"pong"`
}

// 解析消息
func WsMessageBuilder(msgType int, msgData []byte) (wsMessage *WsMessage) {
	return &WsMessage{
		MsgType: msgType,
		MsgData: msgData,
	}
}

// 是否是ping消息
func IsPing(msgData []byte) bool {
	pingData := new(WsPingData)
	if err := json.Unmarshal(msgData, pingData); err == nil {
		if pingData.Ping != 0 {
			return true
		}
	}
	return false
}

// 是否是pong消息
func IsPong(msgData []byte) bool {
	pongData := new(WsPongData)
	if err := json.Unmarshal(msgData, pongData); err == nil {
		if pongData.Pong != 0 {
			return true
		}
	}
	return false
}

// 是否是订阅/取消订阅消息
func IsAction(msgData []byte) (bool, *WsAction) {
	actionData := new(WsAction)
	if err := json.Unmarshal(msgData, actionData); err == nil {
		if strings.Trim(actionData.Ch, " ") == "" {
			return false, actionData
		}
		if actionData.Action == SubAction ||
			actionData.Action == UnsubAction {
			return true, actionData
		}
	}
	return false, actionData
}
