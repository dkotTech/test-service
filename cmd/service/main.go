package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/dkotTech/shutdown"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"test-service/balances"
	"test-service/events"
	"test-service/leaderboard"

	httpbalances "test-service/balances/http"
	httpleaderboard "test-service/leaderboard/http"

	grpcbalances "test-service/balances/grpc"
)

func main() {
	// setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	_ = context.Background()

	// repositories section
	balancesRepo := balances.NewRepositoryMock()
	leaderboardRepo := leaderboard.NewRepositoryMock()

	// services section
	eventsSrv := events.NewEventsService()
	balancesSrv := balances.NewService(balancesRepo, eventsSrv, balances.NewMWValidate, balances.NewMWErrors)
	leaderboardSrv := leaderboard.NewService(leaderboardRepo, eventsSrv, leaderboard.NewMWValidate, leaderboard.NewMWErrors)

	// http routing section
	mux := chi.NewRouter()

	httpLogger := logger.With("server", "http")
	mux.Mount("/api/wallet", httpbalances.Handlers(httpLogger, balancesSrv))
	mux.Mount("/api/leaderboard", httpleaderboard.Handlers(httpLogger, leaderboardSrv))

	mux.Mount("/ws", httpleaderboard.WsHandlers(httpLogger.With("ws", "true"), eventsSrv))

	// setup http server
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		logger.Info("[main-http] start", "host", "localhost:8080")
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Error("[main-http] stopped with err", "err", err)
		}
	}()

	// grpc servers section
	grpcLogger := logger.With("server", "grpc")
	grpcSrv := grpc.NewServer()

	walletSrv := grpcbalances.NewServer(grpcLogger, balancesSrv)
	grpcbalances.RegisterBalancesServiceServer(grpcSrv, walletSrv)

	reflection.Register(grpcSrv)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		logger.Info("[main-grpc] start", "host", "localhost:50051")
		log.Fatal(grpcSrv.Serve(lis))
	}()

	// graceful shutdown callbacks
	wait := shutdown.Graceful(context.Background(), map[string]shutdown.Operation{
		"main-http": func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		},
		"main-grpc": func(ctx context.Context) error {
			grpcSrv.GracefulStop()
			return nil
		},
	})

	<-wait
	logger.Info("bye")
}
