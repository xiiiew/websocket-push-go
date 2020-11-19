package middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// ws中间件
func WsMiddle(router httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// TODO: 鉴权等
		router(w, r, p)
	}
}
