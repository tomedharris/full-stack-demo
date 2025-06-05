package main

import (
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-playground/form/v4"
)

// application acts as a container for all required services.
type application struct {
	config      config
	logger      *slog.Logger
	formDecoder *form.Decoder
	wg          sync.WaitGroup
}

// decodePostForm decodes a form from a POST request and populates the destination struct.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError

		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}

		return err
	}

	return nil
}
