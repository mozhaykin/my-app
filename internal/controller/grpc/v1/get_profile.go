package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gitlab.golang-school.ru/potok-1/amozhaykin/my-app/gen/grpc/profile_v1"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
)

func (h Handlers) GetProfile(ctx context.Context, in *pb.GetProfileInput) (*pb.GetProfileOutput, error) {
	input := dto.GetProfileInput{
		ID: in.GetId(),
	}

	output, err := h.usecase.GetProfile(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())

		default:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	return &pb.GetProfileOutput{
		Id:        output.ID.String(),
		CreatedAt: timestamppb.New(output.CreatedAt),
		UpdatedAt: timestamppb.New(output.UpdatedAt),
		Name:      string(output.Name),
		Age:       int32(output.Age), //nolint:gosec
		Verified:  output.Verified,
		Status:    int32(output.Status), //nolint:gosec
		Contacts: &pb.GetProfileOutput_Contacts{
			Email: output.Contacts.Email,
			Phone: output.Contacts.Phone,
		},
	}, nil
}
