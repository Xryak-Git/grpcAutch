package suite

import (
	"context"
	"github.com/Xryak-Git/grpcAuthProto/gen/go/grpcAuth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpcAuth/interanl/config"
	"net"
	"strconv"
	"testing"
)

var grpcHost = "localhost"

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient grpcAuth.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/config.yml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(
		context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connetcion failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: grpcAuth.NewAuthClient(cc),
	}

}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}

//withDialer := func(l *bufconn.Listener) grpc.DialOption {
//	return grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
//		return l.Dial()
//	})
//}
//
//var l *bufconn.Listener
//
//conn, err := grpc.NewClient(
//"localhost:8080",
//withDialer(l),
//grpc.WithTransportCredentials(insecure.NewCredentials()),
//)
