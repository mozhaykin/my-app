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

func Error(w http.ResponseWriter, err error, status int, message string) {
	log.Error().Err(err).Msg(message)

	err = unpack(err)
	err = fmt.Errorf("%s: %s", message, err) //nolint:errorlint, err113

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
