package v1

import (
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/domain"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/render"
)

func (h *Handlers) GetProfiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input := dto.GetProfilesInput{
		Sort:   r.URL.Query().Get("sort"),
		Order:  r.URL.Query().Get("order"),
		Offset: atoi(r.URL.Query().Get("offset")),
		Limit:  atoi(r.URL.Query().Get("limit")),
	}

	output, err := h.usecase.GetProfiles(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			render.Error(ctx, w, err, http.StatusNotFound, "")

		default:
			render.Error(ctx, w, err, http.StatusBadRequest, "request failed: ")
		}

		return
	}

	render.JSON(w, output, http.StatusOK)
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return n
}
