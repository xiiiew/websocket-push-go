package routes

import "github.com/julienschmidt/httprouter"

func Register() *httprouter.Router {
	router := httprouter.New()

	registerWsRouter(router)
	registerHttpRouter(router)

	return router
}
