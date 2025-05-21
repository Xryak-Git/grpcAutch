package app

import (
	"grpcAuth/interanl/app/grpcapp"
	"grpcAuth/interanl/services/auth"
	"grpcAuth/interanl/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {

	// TODO: инициализировать хранилище
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	// TODO: инициализировать сервис auth
	authService := auth.New(log, tokenTTL, storage, storage, storage)

	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcApp,
	}

}
