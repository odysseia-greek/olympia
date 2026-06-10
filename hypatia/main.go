package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/attike/aristophanes/comedy"
	"github.com/odysseia-greek/olympia/hypatia/mouseion"
	pb "github.com/odysseia-greek/olympia/hypatia/proto/v1"
	"google.golang.org/grpc"
)

const standardPort = ":50061"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=HYPATIA
	logging.System(`
 __ __  __ __  ____   ____  ______  ____   ____
|  |  ||  |  ||    \ /    ||      ||    | /    |
|  |  ||  |  ||  o  )  o  ||      | |  | |  o |
|  _  ||  ~  ||   _/|     ||_|  |_| |  | |     |
|  |  ||___, ||  |  |  _  |  |  |   |  | |  _  |
|  |  ||     ||  |  |  |  |  |  |   |  | |  |  |
|__|__||____/ |__|  |__|__|  |__|  |____||__|__|
`)

	logging.System("\"Ὑπατία\"")
	logging.System("Hypatia - mathematician, astronomer, philosopher")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	ctx := context.Background()
	cfg, err := mouseion.CreateNewConfig(ctx)
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			comedy.UnaryServerInterceptor(
				cfg.Streamer,
				comedy.WithHeaderKey(config.HeaderKey),
				comedy.WithContextKeyName(config.DefaultTracingName),
				comedy.WithCloseHop(),
			),
		),
	)

	pb.RegisterHypatiaServer(server, cfg)

	logging.Info(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
