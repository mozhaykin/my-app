package v1

import (
	pb "github.com/mozhaykin/my-app/gen/grpc/profile_v1"
	"github.com/mozhaykin/my-app/internal/usecase"
)

type Handlers struct {
	pb.UnimplementedProfileV1Server

	usecase *usecase.UseCase
}

func New(uc *usecase.UseCase) *Handlers {
	return &Handlers{
		usecase: uc,
	}
}
