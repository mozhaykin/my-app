package render

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto/baggage"
)

type Err struct {
	Error string `json:"error"`
}

func Error(ctx context.Context, w http.ResponseWriter, err error, status int, message string) {
	baggage.PutError(ctx, err) // кладем ошибку в контекст, для дальнейшего логирования

	err = unpack(err)
	err = fmt.Errorf("%s%w", message, err)

	JSON(w, Err{Error: err.Error()}, status)
}

func unpack(err error) error {
	for {
		e := errors.Unwrap(err)
		if e == nil {
			break
		}

		err = e
	}

	return err
}
