package server

import "net/http"

func Routes(srv *Server) http.Handler {
	mux := http.NewServeMux()

	fsPublic := http.FileServer(http.Dir("./web/public"))
	fsDist := http.FileServer(http.Dir("./web/dist/"))

	mux.Handle("/", fsPublic)
	mux.Handle("/dist/", http.StripPrefix("/dist/", fsDist))

	mux.HandleFunc("GET /api/v1/health", srv.HealthCheck)
	mux.HandleFunc("POST /api/v1/auth/guest", srv.RegisterGuest)
	mux.HandleFunc("QUERY /api/v1/validate", srv.CheckShipPlacements)

	mux.Handle("POST /api/v1/matchmaking", srv.Auth(http.HandlerFunc(srv.EnterMatchmaking)))
	mux.Handle("GET /api/v1/matchmaking/status", srv.Auth(http.HandlerFunc(srv.CheckStatus)))
	mux.HandleFunc("/ws", srv.serveWs)

	return mux
}
