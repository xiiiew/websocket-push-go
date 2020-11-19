package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/xiiiew/websocket-push-go/handles"
	"github.com/xiiiew/websocket-push-go/middleware"
)

func registerWsRouter(router *httprouter.Router)  {
	router.GET("/ws", middleware.WsMiddle(handles.WsHandle))
}
