package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"grpcAuth/interanl/app"
	"grpcAuth/interanl/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()
	fmt.Printf("%v", cfg)

	log := setupLogger(cfg.Env)

	log.Debug("Config loaded", "cfg", cfg)

	application := app.New(log, cfg.GRPC.Port, cfg.Storage, cfg.TokenTTL)
	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping app", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()

	log.Info("app stopped")
	// TODO: logger

	// TODO: storage

	// TODO: run

	// TODO: Grpc инициировать

	// TODO: shutdown

}

func setupLogger(env string) *slog.Logger {

	zapL := zap.Must(zap.NewProduction())
	defer zapL.Sync()

	log := slog.New(
		zapslog.NewHandler(zapL.Core(), &zapslog.HandlerOptions{AddSource: true}),
	)

	return log
}
