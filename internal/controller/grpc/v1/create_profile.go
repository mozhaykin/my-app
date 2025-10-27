package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/grpc/profile_v1"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto/baggage"
)

func (h Handlers) CreateProfile(ctx context.Context, in *pb.CreateProfileInput) (*pb.CreateProfileOutput, error) {
	input := dto.CreateProfileInput{
		Name:  in.GetName(),
		Age:   int(in.GetAge()),
		Email: in.GetEmail(),
		Phone: in.GetPhone(),
	}

	output, err := h.usecase.CreateProfile(ctx, input)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	baggage.PutProfileID(ctx, output.ID.String())

	return &pb.CreateProfileOutput{
		Id: output.ID.String(),
	}, nil
}
