package grpc

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	pb "github.com/mozhaykin/my-app/gen/grpc/profile_v1"
	ver1 "github.com/mozhaykin/my-app/internal/controller/grpc/v1"
	"github.com/mozhaykin/my-app/internal/usecase"
	"github.com/mozhaykin/my-app/pkg/logger"
	"github.com/mozhaykin/my-app/pkg/otel"
)

type Config struct {
	Port string `default:"50051" envconfig:"GRPC_PORT"`
}

type Server struct {
	server *grpc.Server
}

func New(c Config, uc *usecase.UseCase) (*Server, error) {
	s := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(logger.Interceptor, otel.Interceptor),
	)

	// для просмотра через инсомнию или постман
	reflection.Register(s)

	v1 := ver1.New(uc)
	pb.RegisterProfileV1Server(s, v1)

	err := start(s, c.Port)
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	return &Server{server: s}, nil
}

func start(server *grpc.Server, port string) error {
	conn, err := net.Listen("tcp", net.JoinHostPort("", port)) //nolint: noctx
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	go func() {
		err = server.Serve(conn)
		if err != nil {
			log.Error().Err(err).Msg("grpc server: Serve")
		}
	}()

	log.Info().Msg("grpc server: started on port: " + port)

	return nil
}

func (s *Server) Close() {
	s.server.GracefulStop()

	log.Info().Msg("grpc server: closed")
}
