package server

import "net/http"

func Routes(srv *Server) http.Handler {
	mux := http.NewServeMux()

	fsPublic := http.FileServer(http.Dir("./web/public"))
	fsDist := http.FileServer(http.Dir("./web/dist"))
	mux.Handle("/", fsPublic)
	mux.Handle("/dist/", http.StripPrefix("/dist/", fsDist))

	mux.HandleFunc("/ws", srv.serveWs)
	return mux
}
