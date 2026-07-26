package server

import "net/http"

func Routes(srv *Server) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", srv.serveWs)
	return mux
}
