package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jauharfz/battleship/internal/auth"
	"github.com/jauharfz/battleship/internal/lobby"
	"github.com/jauharfz/battleship/internal/server"
)

const (
	addr   = ":8080"
	secret = "inilahmysecretygy"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticketStore := auth.NewTicketStore(6*time.Second, 4*time.Second)
	matchmaker := lobby.NewMatchmaker(ticketStore, 4*time.Second)
	jwtManager := auth.NewJWTManager(secret, 2000*time.Second)
	handler := server.NewServer(jwtManager, ticketStore, matchmaker)
	srv := http.Server{
		Addr:    addr,
		Handler: server.Routes(handler),
	}

	var wg sync.WaitGroup
	wg.Go(func() { handler.Matchmaker.Matching(&wg, ctx) })
	wg.Go(func() { handler.TicketStore.TicketSweeper(ctx) })

	wg.Go(func() {
		slog.Info("server starting", "addr", addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error("the server shutdown abnormally", "error", err)
		}
	})
	<-ctx.Done()
	slog.Info("shutting down server gracefully....")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}
	slog.Info("server exited cleanly")
	wg.Wait()
}
