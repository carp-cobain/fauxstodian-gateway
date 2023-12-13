package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	gw "github.com/carp-cobain/fauxstodian-gateway/gen/go/fauxstodian/v1"
)

var grpcEndpoint = flag.String("grpc-endpoint", "localhost:50055", "gRPC server endpoint")

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := gw.RegisterFauxstodianServiceHandlerFromEndpoint(ctx, mux, *grpcEndpoint, opts); err != nil {
		return err
	}

	// Start HTTP server, proxying calls to gRPC the server
	return http.ListenAndServe(":8080", mux)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		grpclog.Fatal(err)
	}
}
