package handles

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/xiiiew/websocket-push-go/common"
	"github.com/xiiiew/websocket-push-go/response"
	"github.com/xiiiew/websocket-push-go/server"
	"io/ioutil"
	"net/http"
	"strconv"
)

/*
推送到频道
uri: /push/ch/:ch
*/
func PushCh(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ch := p.ByName("ch")
	if ch == "" {
		w.Write(response.HttpErrorResponseBuilder("channel cannot be empty"))
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(response.HttpErrorResponseBuilder(err.Error()))
		return
	}
	if bytes == nil {
		w.Write(response.HttpErrorResponseBuilder("message cannot be empty"))
		return
	}
	fmt.Println("HTTP", ch, string(bytes))
	msg := common.WsMessageBuilder(websocket.TextMessage, bytes)
	server.GetBucketInstance().PushCh(ch, msg)
	w.Write(response.HttpSuccessResponseBuilder(nil))
}

/*
推送给所有用户
uri: /push/all
*/
func PushAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(response.HttpErrorResponseBuilder(err.Error()))
		return
	}
	if bytes == nil {
		w.Write(response.HttpErrorResponseBuilder("message cannot be empty"))
		return
	}
	msg := common.WsMessageBuilder(websocket.TextMessage, bytes)
	server.GetBucketInstance().PushAll(msg)
	w.Write(response.HttpSuccessResponseBuilder(nil))
}

/*
推送给指定连接
uri: /push/one/:connId
*/
func PushOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	connId := p.ByName("connId")
	if connId == "" {
		w.Write(response.HttpErrorResponseBuilder("connId cannot be empty"))
		return
	}
	connIdInt, err := strconv.Atoi(connId)
	if err != nil {
		w.Write(response.HttpErrorResponseBuilder("connId error"))
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(response.HttpErrorResponseBuilder(err.Error()))
		return
	}
	if bytes == nil {
		w.Write(response.HttpErrorResponseBuilder("message cannot be empty"))
		return
	}
	msg := common.WsMessageBuilder(websocket.TextMessage, bytes)
	err = server.GetBucketInstance().PushOne(uint64(connIdInt), msg)
	if err != nil {
		w.Write(response.HttpErrorResponseBuilder(err.Error()))
		return
	}
	w.Write(response.HttpSuccessResponseBuilder(nil))
}
