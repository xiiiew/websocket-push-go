package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/xiiiew/websocket-push-go/handles"
	"github.com/xiiiew/websocket-push-go/middleware"
)

func registerHttpRouter(router *httprouter.Router)  {
	router.POST("/push/ch/:ch", middleware.HttpMiddle(handles.PushCh))
	router.POST("/push/all", middleware.HttpMiddle(handles.PushAll))
}