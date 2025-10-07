package grpcclientv1

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/grpc/profile_v1"
)

type Client struct {
	client pb.ProfileV1Client
	conn   *grpc.ClientConn
}

func New(host string) (*Client, error) {
	conn, err := grpc.NewClient(host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(timeoutInterceptor),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	client := pb.NewProfileV1Client(conn)

	return &Client{client: client, conn: conn}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Error().Err(err).Msg("grpc client: c.conn.Close")
	}
}

func timeoutInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	// Устанавливаем таймаут для каждого вызова
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Вызываем метод
	return invoker(ctx, method, req, reply, cc, opts...)
}
