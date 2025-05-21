package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	authgrpc "grpcAuth/interanl/grpc/auth"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, auth authgrpc.Auth) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, auth)

	return &App{
		log:        log,
		port:       port,
		gRPCServer: gRPCServer,
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const fn = "grpcapp.run"

	log := a.log.With(
		slog.String("fn", fn),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	log.Info(" gRPC server is runnig", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil

}

func (a *App) Stop() {
	const fn = "grpcapp.Stop"

	a.log.With(slog.String("fn", fn)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
