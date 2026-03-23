package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/mozhaykin/my-app/gen/grpc/profile_v1"
	"github.com/mozhaykin/my-app/internal/domain"
	"github.com/mozhaykin/my-app/internal/dto"
	"github.com/mozhaykin/my-app/internal/dto/baggage"
)

func (h Handlers) UpdateProfile(ctx context.Context, in *pb.UpdateProfileInput) (*emptypb.Empty, error) {
	input := dto.UpdateProfileInput{
		ID:    in.GetId(),
		Name:  in.Name,
		Age:   parseAge(in.Age),
		Email: in.Email,
		Phone: in.Phone,
	}

	baggage.PutProfileID(ctx, input.ID)

	err := h.usecase.UpdateProfile(ctx, input)
	if err != nil {
		baggage.PutError(ctx, err)

		switch {
		case errors.Is(err, domain.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())

		default:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}

func parseAge(age *int32) *int {
	if age == nil {
		return nil
	}

	a := int(*age)

	return &a
}
