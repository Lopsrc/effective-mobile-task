package savecar

import (
	"context"
	"effectiveM-test-task/internal/handlers/api"
	"effectiveM-test-task/internal/models/car"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const (
	op = "Save Car"
)

type Response struct {
}

type CarSaver interface {
	Save(ctx context.Context, car car.CarSave) error
}

func New(log *slog.Logger, carSaver CarSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: saving", op))
		ctx := r.Context()

		// API request.
		api.SomeApi(context.Background(), log)

		ct := r.Header.Get("Content-Type")
		if ct != "" {
			mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
			if mediaType != "application/json" {
				msg := "Content-Type header is not application/json"
				http.Error(w, msg, http.StatusUnsupportedMediaType)
				return
			}
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var c car.CarSave
		err := dec.Decode(&c)
		if err != nil {
			var syntaxError *json.SyntaxError
			var unmarshalTypeError *json.UnmarshalTypeError

			switch {
			case errors.As(err, &syntaxError):
				msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
				log.Error(msg)
				http.Error(w, msg, http.StatusBadRequest)

			case errors.Is(err, io.ErrUnexpectedEOF):
				msg := "Request body contains badly-formed JSON"
				log.Error(msg)
				http.Error(w, msg, http.StatusBadRequest)

			case errors.As(err, &unmarshalTypeError):
				msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
				log.Error(msg)
				http.Error(w, msg, http.StatusBadRequest)

			case strings.HasPrefix(err.Error(), "json: unknown field "):
				fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
				msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
				log.Error(msg)
				http.Error(w, msg, http.StatusBadRequest)

			case errors.Is(err, io.EOF):
				msg := "Request body must not be empty"
				log.Error(msg)
				http.Error(w, msg, http.StatusBadRequest)

			case err.Error() == "http: request body too large":
				msg := "Request body must not be larger than 1MB"
				log.Error(msg)
				http.Error(w, msg, http.StatusRequestEntityTooLarge)

			default:
				log.Error("Internal error")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		if err := carSaver.Save(ctx, c); err != nil {
			log.Error("%s: saving error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		if err := json.NewEncoder(w).Encode(Response{}); err != nil {
			log.Error("%s: encoding error: %v", op, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%s: saved", op))
	}
}
