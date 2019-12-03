package osp

import (
	"net/http"
)

func (rs *RestServer) registerRoutes() {
	rs.Mux.HandleFunc(
		"/intendSettle",
		postIntendSettleHandlerFn(rs),
	).Methods("POST")
}

type ()

func postIntendSettleHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
