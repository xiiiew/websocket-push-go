package middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// http中间件
func HttpMiddle(router httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// TODO: 鉴权等
		//token := r.Header.Get("token")
		//if strings.EqualFold(token, "your_token") {
			router(w, r, p)
		//	return
		//}
		//w.WriteHeader(401)
		//w.Write([]byte("401"))
	}
}