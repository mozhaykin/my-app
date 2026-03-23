package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/mozhaykin/my-app/gen/grpc/profile_v1"
	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/internal/dto/baggage"
)

func (h Handlers) GetProfile(ctx context.Context, in *pb.GetProfileInput) (*pb.GetProfileOutput, error) {
	input := dto.GetProfileInput{
		ID: in.GetId(),
	}

	baggage.PutProfileID(ctx, input.ID)

	output, err := h.usecase.GetProfile(ctx, input)
	if err != nil {
		baggage.PutError(ctx, err)

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
