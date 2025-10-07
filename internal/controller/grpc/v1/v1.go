package v1

import (
	pb "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/grpc/profile_v1"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/usecase"
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
