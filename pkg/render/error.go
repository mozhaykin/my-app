package render

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Err struct {
	Error string `json:"error"`
}

func Error(w http.ResponseWriter, err error, status int) {
	log.Error().Err(err).Msg("")

	err = unpack(err)
	err = fmt.Errorf("%w", err)

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
